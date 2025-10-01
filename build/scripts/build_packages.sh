#!/bin/bash

set -euo pipefail

function get_latest_version() {
   curl -s https://api.github.com/repos/tmck-code/pokesay/releases/latest \
     | grep tag_name \
     | cut -d "\"" -f 4- \
     | sed -E 's/",$|v//g'
}

VERSION="${VERSION:-$(get_latest_version)}"
MAINTAINER="Tom McKeesick <tmck01@gmail.com>"
DESCRIPTION="Print pokemon in the CLI! An adaptation of the classic 'cowsay'"

mkdir -p /usr/local/src/build/packages

function build_deb() {
    local os=$1
    local arch=$2
    local suffix=${3:-}

    OUTPUT_DIR="build/deb"
    mkdir -p "$OUTPUT_DIR"

    local bin="build/bin/pokesay-${os}-${arch}${suffix}"
    local pkg_name="pokesay-${os}-${arch}"

    mkdir -p "$pkg_name/pokesay/DEBIAN" "$pkg_name/pokesay/usr/bin" "$pkg_name/pokesay/usr/share/man/man1"

    cp "$bin" "$pkg_name/pokesay/usr/bin/pokesay"

    # Compress and install the man page
    gzip -c "build/packages/pokesay.1" > "$pkg_name/pokesay/usr/share/man/man1/pokesay.1.gz"

    cat > "$pkg_name/pokesay/DEBIAN/control" <<EOF
Package: pokesay
Version: $VERSION
Standards-Version: $VERSION
Section: utils
Priority: optional
Architecture: $arch
Maintainer: $MAINTAINER
Description: $DESCRIPTION
EOF

    dpkg-deb --build "$pkg_name/pokesay/" "$OUTPUT_DIR/pokesay_${VERSION}_${arch}.deb"

    rm -rf "$pkg_name"
    mv -v "$OUTPUT_DIR/pokesay_${VERSION}_${arch}.deb" /usr/local/src/build/packages/
}

function build_arch() {
    local os=$1
    local arch=$2
    local arch_arch=$3
    local suffix=${4:-}

    cd /usr/local/src

    ARCH_DIR="/usr/local/src/build/arch"
    mkdir -p "$ARCH_DIR"

    cp "build/bin/pokesay-${os}-${arch}${suffix}" "$ARCH_DIR/"
    cp "build/packages/pokesay.1" "$ARCH_DIR/"

    cat > "$ARCH_DIR/PKGBUILD" <<EOF
# Maintainer: $MAINTAINER
pkgname=pokesay
gitname=pokesay
pkgver=$VERSION
pkgrel=1
pkgdesc="$DESCRIPTION"
arch=('$arch_arch')
url="https://github.com/tmck-code/pokesay"
license=('BSD-3-Clause')
depends=()
source=("pokesay-linux-amd64")
sha256sums=('SKIP')

package() {
    install -Dm755 "\$srcdir/pokesay-linux-amd64" "\$pkgdir/usr/bin/pokesay"
    install -Dm644 "\$srcdir/../pokesay.1" "\$pkgdir/usr/share/man/man1/pokesay.1"
}
EOF
    cd "$ARCH_DIR"
    makepkg -f --noconfirm

    mv -v "$ARCH_DIR"/*.pkg.tar.zst /usr/local/src/build/packages/
    rm -rf "$ARCH_DIR"
}

case "${1}" in
    deb)   build_deb linux amd64; build_deb android arm64 ;;
    arch)  build_arch linux amd64 x86_64 ;;
    *)     echo "Usage: $0 {deb|arch}" ;;
esac
