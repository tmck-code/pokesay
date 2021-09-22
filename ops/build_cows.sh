#!/bin/bash

set -euo pipefail

if [ -z "${DOCKER_BUILD_DIR:-}" ] || [ -z "${DOCKER_OUTPUT_DIR:-}" ]; then
  echo "DOCKER_BUILD_DIR and DOCKER_OUTPUT_DIR env vars must be defined"
  exit 1
fi

function generateCowfile() {
  local f="${1}"

  local ofpath="${DOCKER_OUTPUT_DIR}${f%.png}.cow"
  local local_path=${ofpath##*icons/}
  local dir_path=$(dirname /tmp/cows/${local_path})

  # Cowsay bubble header
  cat <<EOF > "${ofpath}"
binmode STDOUT, ":utf8";
\$the_cow =<<EOC;
     \$thoughts
      \$thoughts
       \$thoughts
        \$thoughts
EOF

  # Pokemon
  img2xterm "${f}" >> "${ofpath}"

  # Cowsay footer
  cat <<EOF >> "${ofpath}"
EOC
EOF

  [ -d "${dir_path}" ] || mkdir -p "${dir_path}"
  mv ${ofpath} /tmp/cows/${local_path}
}
export -f generateCowfile

find ${DOCKER_BUILD_DIR} -type d | parallel -I {} mkdir -p "${DOCKER_OUTPUT_DIR}{}"
total=$(find ${DOCKER_BUILD_DIR} -type f -iname *.png | wc -l)


BAR='▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇' # this is full bar
EMP='                             ' # This is an empty bar

for i in {1..20}; do
    sleep .1                 # wait 100ms between "frames"
done

i=0
sp='/-\|'
for f in $(find ${DOCKER_BUILD_DIR} -type f -name *.png); do
  generateCowfile "${f}" &
  if ! (( i % 10 )); then
    index=$[ i / 100 ]
    end=$[ 30 - $[ i / 100 ] ]
    echo -ne "\r${BAR:0:$index}${EMP:0:$end}" # print $i chars of $BAR from 0 position
  fi
  ((i=i+1))
done

wait
echo -e "\n- All jobs have finished..."

echo "- Rearranging files"
shopt -s extglob
mv ${DOCKER_BUILD_DIR}/icons/pokemon ${DOCKER_OUTPUT_DIR}/pokemon/
mkdir ${DOCKER_OUTPUT_DIR}/items/
mv ${DOCKER_BUILD_DIR}/icons/* ${DOCKER_OUTPUT_DIR}/items/
echo "- Finished building cowfiles"

