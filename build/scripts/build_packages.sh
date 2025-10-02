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

    OUTPUT_DIR="dist/packages"
    mkdir -p "$OUTPUT_DIR"

    local bin="dist/bin/pokesay-${VERSION}-${os}-${arch}${suffix}"
    local pkg_name="pokesay-${VERSION}-${os}-${arch}"

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
    mkdir -p "$ARCH_DIR"

    cp "dist/bin/pokesay-${VERSION}-${os}-${arch}${suffix}" "$ARCH_DIR/"
    cp "build/packages/pokesay.1" "$ARCH_DIR/"
    cp LICENSE "$ARCH_DIR/"

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
source=("pokesay-${VERSION}-${os}-${arch}${suffix}")
sha256sums=('SKIP')

package() {
    install -Dm755 "\$srcdir/pokesay-${VERSION}-${os}-${arch}${suffix}" "\$pkgdir/usr/bin/pokesay"
    install -Dm644 "\$srcdir/../pokesay.1" "\$pkgdir/usr/share/man/man1/pokesay.1"
    install -Dm644 "\$srcdir/../LICENSE" "\$pkgdir/usr/share/licenses/pokesay/LICENSE"
}
EOF
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
