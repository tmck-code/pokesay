#!/bin/bash

set -euo pipefail

function get_latest_version() {
   curl -s https://api.github.com/repos/tmck-code/pokesay/releases/latest \
     | grep tag_name \
     | cut -d "\"" -f 4- \
     | sed -E 's/",$|v//g'
}

VERSION="$(get_latest_version)"
MAINTAINER="Tom McKeesick <tmck01@gmail.com>"
DESCRIPTION="Print pokemon in the CLI! An adaptation of the classic 'cowsay'"

OUTPUT_DIR="build/deb"
mkdir -p "$OUTPUT_DIR"

function build_deb() {
    local os=$1
    local arch=$2
    local suffix=${3:-}

    local bin="build/bin/pokesay-${os}-${arch}${suffix}"
    local pkg_name="pokesay-${os}-${arch}"

    mkdir -p "$pkg_name" \
        "$pkg_name/pokesay/DEBIAN" \
        "$pkg_name/pokesay/usr/bin" \
        "$pkg_name/pokesay/usr/share/man/man1"

    cp "$bin" "$pkg_name/pokesay/usr/bin/pokesay"
    gzip -c "docs/pokesay.1" > "$pkg_name/pokesay/usr/share/man/man1/pokesay.1.gz"

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
}

build_deb linux   amd64 &
build_deb android arm64 &
wait
