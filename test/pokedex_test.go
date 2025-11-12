package test

import (
	"embed"
	"fmt"
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

func TestCreateNameMetadataMew(test *testing.T) {
	result := pokedex.CreateNameMetadata(
		fmt.Sprintf("%04d", 0),
		"mew",
		pokedex.PokemonName{English: "Mew", Japanese: "ミュウ", JapanesePhonetic: "myuu", Slug: "mew"},
		"data/cows/similar_names/",
		[]string{"data/cows/similar_names/gen7x/mew.cow", "data/cows/similar_names/gen7x/mewtwo.cow"},
	)
	expected := &pokedex.PokemonMetadata{
		Idx:              "0000",
		Name:             "Mew",
		JapaneseName:     "ミュウ",
		JapanesePhonetic: "myuu",
		Entries: []pokedex.PokemonEntryMapping{
			{EntryIndex: 0, Categories: []string{"medium", "gen7x"}},
		},
	}

	Assert(expected, result, test)
}

func TestCreateNameMetadataNatu(test *testing.T) {
	result := pokedex.CreateNameMetadata(
		fmt.Sprintf("%04d", 0),
		"natu",
		pokedex.PokemonName{English: "Natu", Japanese: "ネイティ", JapanesePhonetic: "neiti", Slug: "natu"},
		"data/cows/similar_names/",
		[]string{
			"data/cows/similar_names/gen8/natu.cow",
			"data/cows/similar_names/gen8/eternatus.cow",
			"data/cows/similar_names/gen8/eternatus-eternamax.cow",
		},
	)
	expected := &pokedex.PokemonMetadata{
		Idx:              "0000",
		Name:             "Natu",
		JapaneseName:     "ネイティ",
		JapanesePhonetic: "neiti",
		Entries: []pokedex.PokemonEntryMapping{
			{EntryIndex: 0, Categories: []string{"small", "gen8"}},
		},
	}

	Assert(expected, result, test)
}

func TestCreateNameMetadataEternatus(test *testing.T) {
	result := pokedex.CreateNameMetadata(
		fmt.Sprintf("%04d", 0),
		"eternatus",
		pokedex.PokemonName{English: "Eternatus", Japanese: "ムゲンダイナ", JapanesePhonetic: "mugendaina", Slug: "eternatus"},
		"data/cows/similar_names/",
		[]string{
			"data/cows/similar_names/gen8/natu.cow",
			"data/cows/similar_names/gen8/eternatus.cow",
			"data/cows/similar_names/gen8/eternatus-eternamax.cow",
		},
	)
	expected := &pokedex.PokemonMetadata{
		Idx:              "0000",
		Name:             "Eternatus",
		JapaneseName:     "ムゲンダイナ",
		JapanesePhonetic: "mugendaina",
		Entries: []pokedex.PokemonEntryMapping{
			{EntryIndex: 1, Categories: []string{"big", "gen8"}},
			{EntryIndex: 2, Categories: []string{"big", "gen8"}},
		},
	}

	Assert(expected, result, test)
}