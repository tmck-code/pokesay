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

  # Pokemon
  img2xterm "${f}" | sed '/^\ $/d' >> "${ofpath}"
  echo -e "\n" >> "{ofpath}"

  [ -d "${dir_path}" ] || mkdir -p "${dir_path}"
  mv ${ofpath} /tmp/cows/${local_path}
}
export -f generateCowfile

find ${DOCKER_BUILD_DIR} -type d | parallel -I {} mkdir -p "${DOCKER_OUTPUT_DIR}{}"
total=$(find ${DOCKER_BUILD_DIR} -type f -iname *.png | wc -l)


BAR='▇▇▇▇▇▇▇▇▇▇' # this is full bar
EMP='----------' # This is an empty bar

for i in {1..20}; do
    sleep .1                 # wait 100ms between "frames"
done

i=0
sp='/-\|'
for f in $(find ${DOCKER_BUILD_DIR} -type f -name *.png); do
  generateCowfile "${f}" &
  if ! (( i % 10 )); then
    index=$[ i / 100 ]
    echo -ne "\r $i ${BAR:0:$index}${EMP:$index:10}" # print $i chars of $BAR from 0 position
  fi
  ((i=i+1))
done

wait
echo -e "\n- All jobs have finished..."

echo "- Rearranging files"
shopt -s extglob
ls -alh ${DOCKER_BUILD_DIR}
mv ${DOCKER_BUILD_DIR}/icons/pokemon ${DOCKER_OUTPUT_DIR}/pokemon/
mkdir ${DOCKER_OUTPUT_DIR}/items/
mv ${DOCKER_BUILD_DIR}/icons/* ${DOCKER_OUTPUT_DIR}/items/
echo "- Finished building ${i} cowfiles -> ${DOCKER_OUTPUT_DIR}"

