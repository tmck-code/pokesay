#!/bin/bash

set -euo pipefail

GOOS="$1"
GOARCH="$2"
VERSION="${3:-latest}"

function get_latest_version() {
   curl -s https://api.github.com/repos/tmck-code/pokesay/releases/latest \
     | grep tag_name \
     | cut -d "\"" -f 4- \
     | sed 's/",$//g'
}

# Downloads the release corresponding to the VERSION, GOOS and GOARCH arguments
function download_bin() {
  local pokesay_bin="pokesay-$GOOS-$GOARCH"
  local release_version="$VERSION"

  if [ "$VERSION" == "latest" ]; then
    echo "- No specific release version given, finding latest version..."
    release_version=$(get_latest_version)
  fi
  url="https://github.com/tmck-code/pokesay/releases/download/${release_version}/${pokesay_bin}"
  echo "- VERSION: $release_version, GOOS=$GOOS, GOARCH=$GOARCH"
  echo "- Downloading $url"
  curl -sLo pokesay "$url"
  echo -n "- Downloaded! "
  file pokesay
}

# Move the downloaded bin into the $HOME/bin directory
# Check for any root perms, only use sudo if required
function install_bin() {
  echo "- Installing to $HOME/bin/pokesay... "
  if [ "${TERMUX:-}" == 1 ]; then
    mv -v ./pokesay "$HOME/bin/"
    chmod u+wrx "$HOME/bin/pokesay"
  else
    # Use sudo in case $HOME/bin is root-owned
    if [ -w "$HOME/bin" ]; then
      mv ./pokesay "$HOME/bin/"
      chmod u+wrx "$HOME/bin/pokesay"
    else
      echo "- $HOME/bin is not writable, requires sudo permission"
      sudo mv ./pokesay "$HOME/bin/"
      sudo chmod u+wrx "$HOME/bin/pokesay"
    fi
    echo -n "- Installed to $HOME/bin/pokesay: "
    ls -lh "$HOME/bin/pokesay"
  fi
}

# Checks if the destination dir is in the users PATH, provides instructions if not
function check_path() {
  if [[ ":$PATH:" == *":$HOME/bin:"* ]] ; then
    echo "- Your path already contains $HOME/bin!"
  else
    echo "\
Your path is missing $HOME/bin, you can add it by

  echo '\$PATH=\"\$PATH:\$HOME/bin\"' >> \$HOME/.bashrc

Or, you can always just use the program with the full location: '$HOME/bin/pokesay'
  "
  fi
}

# Provides post-install instructions
# - How to add to bashrc to display pokemon on new shell
function post_install_instructions() {
  echo "
To have a new pokemon in every new shell session, either
- edit your .bashrc file to add 'fortune | pokesay' at the bottom of the file
- or, run this command

    echo 'fortune | pokesay' >> \$HOME/.bashrc
"
}

mkdir -p "$HOME/bin"

echo "1. Downloading"
download_bin
echo "2. Installing"
install_bin

echo "3. Demo"
export PATH="$HOME/bin:$PATH"
echo "hello world!" | pokesay

echo "4. Check \$PATH"
check_path
echo "5. Post-Installation Instructions"
post_install_instructions

