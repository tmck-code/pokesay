#!/bin/sh

# Define install directories and names
pokesay_bin="pokesay"
install_path="$HOME/.${pokesay_bin}"
bin_path="$HOME/bin"

# Make sure the directories exist
mkdir -p $install_path/
mkdir -p $install_path/cows/
mkdir -p $bin_path/

if [ ! -f ${bin_path}/cowsay ]; then
  go get -u -v github.com/msmith491/go-cowsay/...
  ln -sv $HOME/go/src/github.com/msmith491/go-cowsay/cowsay $HOME/bin/
fi

# Remove any previously installed cows
rm -rf $install_path/cows/*

# Copy the cows and the main script to the install path.
tar xzf cows.tar.gz -C $install_path/
N_POKEMON=$(find $install_path/ -type f -name *.cow | wc -l)
cp $pokesay_bin $bin_path/$pokesay_bin
chmod +x "$bin_path/$pokesay_bin"
echo "\nCopied $N_POKEMON Pokémon to install path '$install_path'"

echo "\

- The files were installed to '$install_path/'.
- A '$pokesay_bin' script was created in '$bin_path/'.

It may be necessary to logout and login back again in order to have the '$pokesay_bin' available in your path.
"
