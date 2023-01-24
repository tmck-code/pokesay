package test

import (
	"embed"
	"os"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

var (
	//go:embed data/cows/*cow
	GOBCowData embed.FS
	//go:embed data/cows/*metadata
	GOBMetadata embed.FS
)

func TestReadNames(test *testing.T) {
	result := pokedex.ReadNames("./data/pokemon.json")

	expected := map[string]pokedex.PokemonName{
		"bulbasaur": {English: "Bulbasaur", Japanese: "フシギダネ", JapanesePhonetic: "fushigidane"},
		"ivysaur":   {English: "Ivysaur", Japanese: "フシギソウ", JapanesePhonetic: "fushigisou"},
		"venusaur":  {English: "Venusaur", Japanese: "フシギバナ", JapanesePhonetic: "fushigibana"},
	}

	Assert(expected, result, result, test)
}

func TestReadEntry(test *testing.T) {
	result := pokedex.ReadPokemonCow(GOBCowData, "data/cows/1.cow")

	expected, err := os.ReadFile("data/cows/egg.cow")
	pokesay.Check(err)

	Assert(expected, result, result, test)
}

func TestReadMetadataFromEmbedded(test *testing.T) {
	result := pokedex.ReadMetadataFromEmbedded(GOBMetadata, "data/cows/4.metadata")

	expected := pokedex.PokemonMetadata{
		Categories:       "small/gen7x/regular",
		Name:             "abomasnow",
		JapaneseName:     "ユキノオー",
		JapanesePhonetic: "yukinoo-",
	}

	Assert(expected, result, result, test)
}
