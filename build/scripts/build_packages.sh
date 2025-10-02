#!/bin/bash

set -euo pipefail
[ -n "${DEBUG:-}" ] && set -x
[ -z "${VERSION:-}" ] && echo "no VERSION!" && exit 1

MAINTAINER="Tom McKeesick <tmck01@gmail.com>"
DESCRIPTION="Print pokemon in the CLI! An adaptation of the classic 'cowsay'"

function build_deb() {
    local os=$1
    local arch=$2
    local suffix=${3:-}

    OUTPUT_DIR="dist/packages"
    mkdir -p "$OUTPUT_DIR"

    local bin="dist/bin/pokesay-${VERSION}-${os}-${arch}${suffix}"
    local pkg_name="pokesay-${VERSION}-${os}-${arch}"

    mkdir -p \
        "$pkg_name/pokesay/DEBIAN" \
        "$pkg_name/pokesay/usr/bin" \
        "$pkg_name/pokesay/usr/share/man/man1"

    cp "$bin" "$pkg_name/pokesay/usr/bin/pokesay"
    cat build/packages/pokesay.1 | \
      sed -e "s/DATE/$(date '+%B %Y')/g" \
          -e "s/VERSION/$VERSION/g" | \
      gzip -c > "$pkg_name/pokesay/usr/share/man/man1/pokesay.1.gz"

    cat build/packages/DEBIAN/control | \
      sed -e "s/VERSION/$VERSION/g" \
          -e "s/ARCH/$arch/g" \
          -e "s/MAINTAINER/$MAINTAINER/g" \
          -e "s/DESCRIPTION/$DESCRIPTION/g" \
      > "$pkg_name/pokesay/DEBIAN/control"

    dpkg-deb --build "$pkg_name/pokesay/" "$OUTPUT_DIR/pokesay-${VERSION}-${os}-${arch}.deb"
    rm -rf "$pkg_name"
}

function build_arch() {
    local os=$1
    local arch=$2
    local arch_arch=$3
    local suffix=${4:-}

    cd /usr/local/src

    ARCH_DIR="/usr/local/src/build/arch"
    BIN_FILE="pokesay-${VERSION}-${os}-${arch}${suffix}"
    mkdir -p "$ARCH_DIR"

    cp "dist/bin/$BIN_FILE" "$ARCH_DIR/"
    cat "build/packages/pokesay.1" | \
      sed -e "s/DATE/$(date '+%B %Y')/g" \
          -e "s/VERSION/$VERSION/g" \
      > "$ARCH_DIR/pokesay.1"
    cp LICENSE "$ARCH_DIR/"

    SHA256_SUM=$(sha256sum "$ARCH_DIR/$BIN_FILE" | cut -d' ' -f1)

    cat build/packages/arch/PKGBUILD | \
      sed -e "s/VERSION/$VERSION/g" \
          -e "s/ARCH_ARCH/$arch_arch/g" \
          -e "s/BIN_FILE/$BIN_FILE/g" \
          -e "s/MAINTAINER/$MAINTAINER/g" \
          -e "s/DESCRIPTION/$DESCRIPTION/g" \
          -e "s/SHA256_SUM/$SHA256_SUM/g" \
      > "$ARCH_DIR/PKGBUILD"

    cd "$ARCH_DIR"
    makepkg --printsrcinfo > .SRCINFO
    makepkg -f --noconfirm

    cp "$ARCH_DIR"/*.pkg.tar.zst /usr/local/src/dist/packages/
    rm -rf "$ARCH_DIR"
}

case "${1}" in
    deb)   build_deb linux amd64; build_deb android arm64 ;;
    arch)  build_arch linux amd64 x86_64 ;;
    *)     echo "Usage: $0 {deb|arch}" ;;
esac
