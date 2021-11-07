#!/bin/bash

set -euxo pipefail

VERSION=0.1.0
GOOS=$1
GOARCH=$2

mkdir -p $HOME/bin
curl -L --output $HOME/bin/pokesay https://github.com/tmck-code/pokesay-go/releases/download/v${VERSION}/pokesay-$GOOS-$GOARCH
export PATH="$HOME/bin:$PATH"
echo "hello world!" | pokesay
