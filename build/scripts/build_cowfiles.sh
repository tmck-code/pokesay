#!/bin/bash

set -euo pipefail

mkdir -p /tmp/convert/

# Convert all of the pokesprite .pngs -> cowfiles for the terminal
go run /usr/local/src/src/bin/convert/png_convert.go \
  -from /tmp/original/pokesprite/ \
  -tmpDir /tmp/convert/ \
  -to /tmp/cows/ \
  -padding 4 \
  -skipDuplicates \
  -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]' \
  && mv -v /tmp/cows/pokemon-gen8 /tmp/cows/gen8 > /dev/null \
  && mv -v /tmp/cows/pokemon-gen7x /tmp/cows/gen7x > /dev/null \
  && cat /tmp/original/pokesprite/data/pokemon.json | jq -c .[] > /tmp/cows/pokemon.json
