#!/bin/bash

set -euo pipefail

FROM="${1:-/tmp/cows/}"

go run ./src/bin/pokedex/pokedex.go \
  -from "${FROM}" \
  -fromMetadata "${FROM}/pokemon.json" \
  -to ./build/assets/ \
  -toDataSubDir cows/ \
  -toMetadataSubDir metadata/ \
  -toTotalFname total.txt

rm -rf cows
tree -L 1 build/assets
