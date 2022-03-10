#!/bin/bash

set -euo pipefail

echo "\
package main

type PokemonName struct {
        Categories []string
        Name       string
}

type PokemonData struct {
        Name PokemonName
        Data []byte
}

var PokemonList []PokemonData = []PokemonData{" > header

cp bindata.go bindata.go.bak

grep 'var _cows' bindata.go > pokemon_data
grep '^.*"cows/.*cow",$' bindata.go > pokemon_names

rm -f left centre right data tmp tmp_bytes tmp_end

echo "Generating data"
while read line; do
    echo $line | cut -d '=' -f 2 | tr -d ' ' >> data
    echo $line | sed 's/var _cows/\t{PokemonName{[]string{"/g' | sed 's/Cow = /Cow"}, "/g' | sed 's/, "".*/, "/g' >> left
    echo '"}, []byte(' >> tmp_bytes
    echo ')},' >> tmp_end
done < pokemon_data

echo "Generating names"
while read line; do
    echo $line | sed 's/^.*cows\///g' | sed 's/\.cow",$//g' | rev | cut -d '/' -f 1 | rev >> centre
done < pokemon_names

n_pokemon=$(wc -l pokemon_data)


cat header > tmp

paste -d '' left centre tmp_bytes data tmp_end >> tmp

echo "}" >> tmp

mv tmp cmd/bindata.go
