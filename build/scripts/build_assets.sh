#!/bin/bash

set -euo pipefail

extras=""
if [ "${ENV:-}" == "dev" ]; then
  extras="-newlineMode"
fi

tar xzf build/cows.tar.gz
go run ./src/bin/pokedex/pokedex.go \
  -from ./cows/ \
  -to ./build/assets/ \
  -toCategoryFpath pokedex.gob \
  -toDataSubDir cows/ \
  -toMetadataSubDir metadata/ \
  -toTotalFname total.txt \
  "${extras}"

rm -rf cows
ls -alh build/assets
