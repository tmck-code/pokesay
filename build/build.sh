#!/bin/bash

set -euo pipefail

function build() {
    echo "building $1 / $2"
    GOOS=$1 GOARCH=$2 go build -o bin/pokesay-${1}-${2}${3:-} pokesay.go
}

mkdir bin/

build darwin  amd64 &
build darwin  arm64 &
build linux   amd64 &
build windows amd64 .exe &
build android arm64 &
wait