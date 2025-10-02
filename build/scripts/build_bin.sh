#!/bin/bash

set -euo pipefail
[ -n "${DEBUG:-}" ] && set -x
[ -z "${VERSION:-}" ] && echo "no VERSION!" && exit 1

OUTPUT_DIR="dist/bin"
mkdir -p "$OUTPUT_DIR"

cp build/packages/pokesay.1 .

function build() {
    echo "- building $1 / $2"
    GOOS=$1 GOARCH=$2 go build -o "${OUTPUT_DIR}/pokesay-${VERSION}-${1}-${2}${3:-}" pokesay.go
    echo "- built as ${OUTPUT_DIR}/pokesay-${VERSION}-${1}-${2}${3:-}"
}

function tarball() {
    echo "- tarballing $1 / $2"
    cp "${OUTPUT_DIR}/pokesay-${VERSION}-${1}-${2}${3:-}" .

    tar czf \
        "dist/tarballs/pokesay-${VERSION}-${1}-${2}${3:-}.tar.gz" \
        "pokesay-${VERSION}-${1}-${2}${3:-}" LICENSE pokesay.1

    rm -f "pokesay-${VERSION}-${1}-${2}${3:-}"
}

build darwin  amd64 &
build darwin  arm64 &
build linux   amd64 &
build windows amd64 .exe &
build android arm64 &
wait

# just create a tarball for the linux/amd64 (used for AUR package)
tarball linux amd64

rm -f pokesay.1