#!/bin/bash

set -euo pipefail

# Convert all of the pokesprite .pngs -> cowfiles for the terminal
go run /usr/local/src/src/bin/convert/png_convert.go \
  -from /tmp/original/pokesprite/ \
  -to /tmp/cows/ \
  -padding 4 \
  -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]' \
  && mv -v /tmp/cows/pokemon-gen8 /tmp/cows/gen8 \
  && mv -v /tmp/cows/pokemon-gen7x /tmp/cows/gen7x \
  && cat /tmp/original/pokesprite/data/pokemon.json | jq -c .[] > /tmp/cows/pokemon.json
