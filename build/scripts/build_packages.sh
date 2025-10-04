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
        "$pkg_name/pokesay/usr/share/man/man1" \
        "$pkg_name/pokesay/usr/share/bash-completion/completions" \
        "$pkg_name/pokesay/usr/share/zsh/site-functions" \
        "$pkg_name/pokesay/usr/share/fish/vendor_completions.d" \
        "$pkg_name/pokesay/usr/share/pokesay"

    cp "$bin" "$pkg_name/pokesay/usr/bin/pokesay"
    cat build/packages/pokesay.1 | \
      sed -e "s/DATE/$(date '+%B %Y')/g" \
          -e "s/VERSION/$VERSION/g" | \
      gzip -c > "$pkg_name/pokesay/usr/share/man/man1/pokesay.1.gz"

    # Add completions and data files
    cp build/packages/pokesay-completion.bash "$pkg_name/pokesay/usr/share/bash-completion/completions/pokesay"
    cp build/packages/pokesay-completion.zsh "$pkg_name/pokesay/usr/share/zsh/site-functions/_pokesay"
    cp build/packages/pokesay-completion.fish "$pkg_name/pokesay/usr/share/fish/vendor_completions.d/pokesay.fish"
    cp build/packages/pokesay-names.txt "$pkg_name/pokesay/usr/share/pokesay/pokesay-names.txt"
    cp build/packages/pokesay-ids.txt "$pkg_name/pokesay/usr/share/pokesay/pokesay-ids.txt"

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

    cd /home/u/pokesay

    ARCH_DIR="/home/u/pokesay/build/arch"
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

    cp build/packages/pokesay-completion.bash "$ARCH_DIR/"
    cp build/packages/pokesay-completion.zsh "$ARCH_DIR/"
    cp build/packages/pokesay-completion.fish "$ARCH_DIR/"
    cp build/packages/pokesay-names.txt "$ARCH_DIR/"
    cp build/packages/pokesay-ids.txt "$ARCH_DIR/"

    su - u -c "\
      cd \"$ARCH_DIR\" && \
      makepkg --printsrcinfo > .SRCINFO && \
      makepkg -f --noconfirm"

    cp "$ARCH_DIR"/*.pkg.tar.zst dist/packages/
    rm -rf "$ARCH_DIR"
}

case "${1}" in
    deb)   build_deb linux amd64; build_deb android arm64 ;;
    arch)  build_arch linux amd64 x86_64 ;;
    *)     echo "Usage: $0 {deb|arch}" ;;
esac
