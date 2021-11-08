#!/bin/bash

set -euo pipefail

VERSION=0.1.0
GOOS=$1
GOARCH=$2

mkdir -p $HOME/bin

curl -sL -o pokesay https://github.com/tmck-code/pokesay-go/releases/download/v${VERSION}/pokesay-$GOOS-$GOARCH
# Use sudo in case someone's $HOME/bin dir is root-owned
sudo mv -v ./pokesay $HOME/bin/
sudo chmod u+wrx $HOME/bin/pokesay

export PATH="$HOME/bin:$PATH"
echo "hello world!" | pokesay
echo -e "\nInstall complete! Location: $HOME/bin/pokesay"
