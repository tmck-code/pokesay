#!/bin/sh

# Define install directories and names
pokesay_bin="pokesay"
install_path="$HOME/.${pokesay_bin}"
bin_path="$HOME/bin"

# Make sure the install paths exist
mkdir -p ${install_path}/cows/ ${bin_path}
# Remove any previously installed cowfiles and binaries
rm -rf ${install_path}/cows/* ${bin_path}/cowsay

# Copy the cows and the main script to the install path.
tar xzf cows.tar.gz -C ${install_path}/
N_POKEMON=$(find ${install_path}/ -type f -name *.cow | wc -l)
# Copy the executable to the install path and ensure it has +x permissions
cp -v ${pokesay_bin} cowsay ${bin_path}/
chmod +x "${bin_path}/${pokesay_bin}"

cat <<EOF

Copied $N_POKEMON Pokémon to install path '${install_path}'

- The files were installed to '${install_path}/'.
- A '${pokesay_bin}' script was created in '${bin_path}/'.

It may be necessary to logout and login back again in order to have the '${pokesay_bin}' available in your path.
EOF
