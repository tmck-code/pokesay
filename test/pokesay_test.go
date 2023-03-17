package test

import (
	"embed"
	"fmt"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

var (
	//go:embed data/total.txt
	GOBTotal []byte
	//go:embed data/cows/*.metadata
	GOBCowNames embed.FS
	//go:embed all:data/categories
	GOBCategories embed.FS
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
	// ChooseByCategory(
	//	category string,
	//	categoryDir []fs.DirEntry,
	//	categoryFiles embed.FS,
	//	categoryRootDir string,
	//	metadataFiles embed.FS,
	//	metadataRootDir string
	//)
	//(pokedex.PokemonMetadata, pokedex.PokemonEntryMapping)

	dir, _ := GOBCategories.ReadDir("data/categories/small")

	fmt.Printf("%#v\n", dir)

	metadata, entry := pokesay.ChooseByCategory(
		"small",
		dir,
		GOBCategories,
		"data/categories",
		GOBCowNames,
		"data/cows",
	)

	Assert(metadata, "", test)
	Assert(entry, "", test)
}

func TestChooseByRandomIndex(test *testing.T) {
	resultTotal, result := pokesay.ChooseByRandomIndex(GOBTotal)
	Assert(9, resultTotal, test)
	Assert(0 <= result, true, test)
	Assert(9 >= result, true, test)
}
