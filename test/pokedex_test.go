package test

import (
	"embed"
	"os"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
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
		"bulbasaur": {English: "Bulbasaur", Japanese: "フシギダネ", JapanesePhonetic: "fushigidane", Slug: "bulbasaur"},
		"ivysaur":   {English: "Ivysaur", Japanese: "フシギソウ", JapanesePhonetic: "fushigisou", Slug: "ivysaur"},
		"venusaur":  {English: "Venusaur", Japanese: "フシギバナ", JapanesePhonetic: "fushigibana", Slug: "venusaur"},
	}

	Assert(expected, result, test)
}

func TestReadEntry(test *testing.T) {
	result := pokedex.ReadPokemonCow(GOBCowData, "data/cows/1.cow")

	expected, err := os.ReadFile("data/cows/egg.cow")
	pokedex.Check(err)

	Assert(expected, result, test)
}

func TestReadMetadataFromEmbedded(test *testing.T) {
	result := pokedex.ReadMetadataFromEmbedded(GOBMetadata, "data/cows/4.metadata")

	expected := pokedex.PokemonMetadata{
		Name:             "Hoothoot",
		JapaneseName:     "ホーホー",
		JapanesePhonetic: "ho-ho-",
		Entries: []pokedex.PokemonEntryMapping{
			{EntryIndex: 1586, Categories: []string{"small", "gen7x", "shiny"}},
			{EntryIndex: 2960, Categories: []string{"small", "gen8", "regular"}},
			{EntryIndex: 4285, Categories: []string{"small", "gen8", "shiny"}},
			{EntryIndex: 428, Categories: []string{"small", "gen7x", "regular"}},
		},
	}

	Assert(expected, result, test)
}
