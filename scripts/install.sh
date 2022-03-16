#!/bin/bash

set -euo pipefail

VERSION=0.6
GOOS=$1
GOARCH=$2

mkdir -p $HOME/bin
curl -sL -o pokesay https://github.com/tmck-code/pokesay-go/releases/download/v${VERSION}/pokesay-$GOOS-$GOARCH

if [ "${TERMUX:-}" == 1 ]; then
  mv -v ./pokesay $HOME/bin/
  chmod u+wrx $HOME/bin/pokesay
else
  # Use sudo in case $HOME/bin is root-owned
  sudo mv -v ./pokesay $HOME/bin/
  sudo chmod u+wrx $HOME/bin/pokesay
fi

export PATH="$HOME/bin:$PATH"
echo "hello world!" | pokesay
echo -e "\nInstall complete! Location: $HOME/bin/pokesay"
