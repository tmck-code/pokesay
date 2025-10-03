#!/bin/bash
#
# Run with:
#
# docker build -f docs/Dockerfile.bench -t poke-bench .
# mkdir output
# docker run -it \
#   -v $PWD/docs/bench.sh:/home/poke-bench/bench.sh \
#   -v $PWD/output/:/home/poke-bench/output/ \
#   poke-bench:latest \
#   bash -c "./bench.sh"

set -euo pipefail

hyperfine \
    -m 200 \
    -n 'tmck-code/pokesay' 'echo w | pokesay -n pikachu' \
    -n 'yannjor/krabby' 'krabby name pikachu' \
    -n 'luna-altair/kingler' 'kingler name pikachu' \
    -n 'xiota/pokemon-colorscripts' 'pokemon-colorscripts --name pikachu' \
    -n 'rubiin/pokego' 'pokego -n pikachu' \
    -n 'talwat/pokeget' 'pokeget pikachu' \
    -n 'possatti/pokemonsay' 'echo w | ./bin/pokemonsay-bash -p Pikachu' \
    -n 'HRKings/pokemonsay-newgenerations' 'echo w | pokemonsay-newgenerations -p pikachu' \
    -n 'dfrankland/pokemonsay' 'echo w | pokemonsay-npm' \
    --export-markdown output/bench-results.1.md

hyperfine -i \
    -m 200 \
    -n 'tmck-code/pokesay' 'echo w | pokesay' \
    -n 'yannjor/krabby' 'krabby random' \
    -n 'luna-altair/kingler' 'kingler random' \
    -n 'xiota/pokemon-colorscripts' 'pokemon-colorscripts --random' \
    -n 'rubiin/pokego' 'pokego -r 1,2,3,4,5' \
    -n 'talwat/pokeget' 'pokeget random' \
    -n 'possatti/pokemonsay' 'echo w | ./bin/pokemonsay-bash' \
    -n 'HRKings/pokemonsay-newgenerations' 'echo w | pokemonsay-newgenerations' \
    -n 'dfrankland/pokemonsay' 'echo w | pokemonsay-npm' \
    --export-markdown output/bench-results.2.md
