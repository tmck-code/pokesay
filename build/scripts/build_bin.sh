#!/bin/bash

set -euo pipefail
[ -n "${DEBUG:-}" ] && set -x
[ -z "${VERSION:-}" ] && echo "no VERSION!" && exit 1

test -f "/usr/local/opt/coreutils/libexec/gnubin" && \
  export PATH="/usr/local/opt/coreutils/libexec/gnubin:${PATH}"

OUTPUT_DIR="dist/bin"
mkdir -p "$OUTPUT_DIR"

cat build/packages/pokesay.1 | \
  sed -e "s/DATE/$(date '+%B %Y')/g" \
      -e "s/VERSION/$VERSION/g" > pokesay.1

function build() {
    echo -e "  - building $1 / $2"
    GOOS=$1 GOARCH=$2 go build \
      -o "${OUTPUT_DIR}/pokesay-${VERSION}-${1}-${2}${3:-}" \
      pokesay.go > /dev/null 2>&1
    echo -e "  \e[1;32m✔ built as ${OUTPUT_DIR}/pokesay-${VERSION}-${1}-${2}${3:-}\e[0m"
}

function tarball() {
    echo "  - tarballing $1 / $2"
    local binfile="pokesay-${VERSION}-${1}-${2}${3:-}"

    cp \
        "$OUTPUT_DIR"/pokesay-* \
        build/packages/pokesay-completion.bash \
        build/packages/pokesay-completion.zsh \
        build/packages/pokesay-completion.fish \
        build/packages/pokesay-names.txt \
        build/packages/pokesay-ids.txt \
        .

    # create a tarball for the linux/amd64 binary (used for AUR package)
    tar czf \
        "dist/tarballs/${binfile}.tar.gz" \
        "$binfile" \
        LICENSE \
        pokesay.1 \
        pokesay-completion.bash \
        pokesay-completion.zsh \
        pokesay-completion.fish \
        pokesay-names.txt \
        pokesay-ids.txt

    mkdir -p usr/share/pokesay
    mv pokesay-names.txt pokesay-ids.txt usr/share/pokesay/

    # create a full tarball with all binaries and files (used for homebrew formula)
    tar czf \
        "dist/tarballs/pokesay-$VERSION.tar.gz" \
        pokesay-$VERSION-* \
        LICENSE \
        pokesay.1 \
        pokesay-completion.bash \
        pokesay-completion.zsh \
        pokesay-completion.fish \
        usr/share/pokesay/
    echo -e "  \e[1;32m✔ tarballed as dist/tarballs/${binfile}.tar.gz\e[0m"

    rm -rf pokesay.1 pokesay-* usr/
}

build darwin  amd64
build darwin  arm64
build linux   amd64
build windows amd64 .exe
build android arm64
wait

# just create a tarball for the linux/amd64 (used for AUR package)
tarball linux amd64

rm -f pokesay.1
