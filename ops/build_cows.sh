#!/bin/bash

set -euo pipefail

if [ -z "${DOCKER_BUILD_DIR:-}" ] || [ -z "${DOCKER_OUTPUT_DIR:-}" ]; then
  echo "DOCKER_BUILD_DIR and DOCKER_OUTPUT_DIR env vars must be defined"
  exit 1
fi

function generateCowfile() {
  local worker_n="${1}"
  local f="${2}"

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
  echo -n '.'
  # printf "%9s %s\n" "${worker_n}" "Generated cowfile: ${f}"
}
export -f generateCowfile

find ${DOCKER_BUILD_DIR} -type d | parallel -I {} mkdir -p "${DOCKER_OUTPUT_DIR}{}"
total=$(find ${DOCKER_BUILD_DIR} -type f -iname *.png | wc -l)

i=0
for f in $(find ${DOCKER_BUILD_DIR} -type f -name *.png); do
  generateCowfile "${i}/${total}" "${f}" &
  ((i=i+1))
done

wait
echo -e "\n- All jobs have finished..."

echo "- Rearranging files"
shopt -s extglob
mv ${DOCKER_BUILD_DIR}/icons/pokemon ${DOCKER_OUTPUT_DIR}/pokemon/
mkdir ${DOCKER_OUTPUT_DIR}/items/
mv ${DOCKER_BUILD_DIR}/icons/!(items) ${DOCKER_OUTPUT_DIR}/items/

