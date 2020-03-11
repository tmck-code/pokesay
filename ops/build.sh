#!/bin/bash

set -euo pipefail

IDIR="${1}"
ODIR="${2}"


function generateCowfile() {
  local worker_n="${1}"
  local f="${2}"

  local ofpath="${ODIR}${f%.png}.cow"
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
  printf "%9s %s\n" "${worker_n}" "Generated cowfile: ${f}"
}
export -f generateCowfile

find ${IDIR} -type d | parallel -I {} mkdir -p "${ODIR}{}"
total=$(find ${IDIR} -type f -iname *.png | wc -l)

i=0
for f in $(find ${IDIR} -type f -name *.png); do
  generateCowfile "${i}/${total}" "${f}" &
  ((i=i+1))
done

echo "- Waiting for all jobs to finish..."
wait

echo "- Rearranging files"
shopt -s extglob
mv ${IDIR}/pokemon ${ODIR}/pokemon/
mkdir ${ODIR}/items/
mv ${IDIR}/!(items) ${ODIR}/items/
