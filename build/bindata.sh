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

var PokemonList []PokemonData = []PokemonData{" > tmp

grep 'var _cows' bindata.go | sed 's/var _cows/\t{PokemonName{[]string{"/g' | sed 's/Cow = /Cow"}, ""}, []byte(/g' | sed 's/$/)},/g' >> tmp

echo "}" >> tmp

mv tmp cmd/bindata.go
