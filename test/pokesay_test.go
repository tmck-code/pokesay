package test

import (
	"embed"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

var (
	//go:embed data/total.txt
	GOBTotal []byte
	//go:embed data/cows/*.metadata
	GOBCowNames embed.FS
)

func TestChooseByName(test *testing.T) {
	names := make(map[string][]int)
	names["hoothoot"] = []int{4}
	result, _ := pokesay.ChooseByName(
		names,
		"hoothoot",
		GOBCowNames,
		"data/cows",
	)

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

func TestChooseByCategory(test *testing.T) {
}

func TestChooseByRandomIndex(test *testing.T) {
	resultTotal, result := pokesay.ChooseByRandomIndex(GOBTotal)
	Assert(9, resultTotal, test)
	Assert(0 <= result, true, test)
	Assert(9 >= result, true, test)
}
