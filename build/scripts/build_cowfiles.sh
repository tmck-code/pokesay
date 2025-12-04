#!/bin/bash

set -euo pipefail

skip=(resources/ misc/ icons/ items/ items-outline)
declare -A h

cd /tmp/original/cows/
fpaths=$(find . -iname '*.cow' -type f | sort -h)

nOldFpaths=$(echo "$fpaths" | wc -l | tr -d '\n')
i=1
echo "$fpaths" | while read f; do
  progress=$(( ((i * 100)) / nOldFpaths))
  for s in "${skip[@]}"; do
    if [[ "$f" == *"$s"* ]]; then
      echo -e "[$i/$nOldFpaths $progress%] \e[0;33m skipping:  $f\e[0m"
      ((i++))
      continue 2
    fi
  done

  s=$(sha256sum "$f" | cut -d' ' -f1)
  if [[ ${h[$s]:-} ]]; then
    echo -e "[$i/$nOldFpaths $progress%] \e[0;31m duplicate: $f\e[0m"
  else
    echo -e "[$i/$nOldFpaths $progress%] \e[0;32m copying:   $f > /tmp/cows/$f\e[0m"
    cp --parents "$f" /tmp/cows/
    h[$s]=1
  fi
  ((i++))
done

# rename dirs & generate pokemon.json metadata file
mv -v /tmp/cows/pokemon-gen8 /tmp/cows/gen8
mv -v /tmp/cows/pokemon-gen7x /tmp/cows/gen7x
cat /tmp/original/cows/data/pokemon.json | jq -c .[] > /tmp/cows/pokemon.json

nNewFpaths=$(find /tmp/cows/ -iname '*.cow' | wc -l | tr -d '\n')
echo -e "\n\e[1;32m✔ all done, total files: $nNewFpaths / $nOldFpaths\e[0m"
