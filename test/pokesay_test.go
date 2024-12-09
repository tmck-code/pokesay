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
	dir, _ := GOBCategories.ReadDir("data/categories/small")

	metadata, entry := pokesay.ChooseByCategory(
		"small",
		dir,
		GOBCategories,
		"data/categories",
		GOBCowNames,
		"data/cows",
	)

	expectedMetadata := pokedex.PokemonMetadata{
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

	expectedEntry := pokedex.PokemonEntryMapping{
		EntryIndex: 2960,
		Categories: []string{"small", "gen8", "regular"},
	}

	Assert(expectedMetadata, metadata, test)
	Assert(expectedEntry, entry, test)
}

func TestChooseByNameAndCategory(test *testing.T) {
	names := make(map[string][]int)
	names["hoothoot"] = []int{4}
	metadata, entry := pokesay.ChooseByNameAndCategory(
		names,
		"hoothoot",
		GOBCowNames,
		"data/cows",
		"small",
	)

	Assert("small", entry.Categories[0], test)
	Assert("Hoothoot", metadata.Name, test)
}

func TestChooseByRandomIndex(test *testing.T) {
	resultTotal, result := pokesay.ChooseByRandomIndex(GOBTotal)
	Assert(9, resultTotal, test)
	Assert(0 <= result, true, test)
	Assert(9 >= result, true, test)
}

func TestUnicodeStringLength(test *testing.T) {
	msg := []string{
		" ▄  █ ▄███▄   █    █    ████▄       ▄ ▄   ████▄ █▄▄▄▄ █     ██▄",   // 63
		"█   █ █▀   ▀  █    █    █   █      █   █  █   █ █  ▄▀ █     █  █",  // 64
		"██▀▀█ ██▄▄    █    █    █   █     █ ▄   █ █   █ █▀▀▌  █     █   █", // 65
		"█   █ █▄   ▄▀ ███▄ ███▄ ▀████     █  █  █ ▀████ █  █  ███▄  █  █",  // 64
		"   █  ▀███▀       ▀    ▀           █ █ █          █       ▀ ███▀",  // 64
		"  ▀                                 ▀ ▀          ▀",                // 50
	}
	expected := []int{63, 64, 65, 64, 64, 50}
	results := make([]int, len(msg))
	for i, line := range msg {
		results[i] = pokesay.UnicodeStringLength(line)
	}
	Assert(expected, results, test)
}

func TestTokeniseAnsiLine(test *testing.T) {
	line := "[35mAAA [41m XX [0m"
	expected := []string{"[35m", "AAA ", "[41m", " XX ", "[0m"}
	result := pokesay.TokeniseANSILine(line)
	Assert(expected, result, test)
}

func TestFlipHorizontalLine(test *testing.T) {
	// The AAA has a purple fg
	// The XX has a red bg
	line := "[35mAAA [41m XX [0m"

	// The AAA should still have a purple fg
	// The XX should still have a red bg
	expected := "[41m XX [35m AAA[0m"
	result := pokesay.ReverseString(line)

	Assert(expected, result, test)
}

func TestFlipHorizontalWithoutColour(test *testing.T) {
	msg := []string{
		"         ▄▄          ▄▄",
		"        ▄▄▄     ▄▄▄▄▄▄ ▄▄",
		"       ▄  ▄▀ ▄▄▄  ▄▄   ▄▀",
		"     ▄▄▄   ▄▄  ▄▄    ▄▀",
		"    ▄▄   ▄▄▄  ▄ ▀▄  ▄▄",
		"    ▀▄▄   ▄▄▄   ▄▄▄ ▄▀",
		"    ▀▄▄▄▄▄ ▄▄   ▄▄▄▄▄",
		"           ▄▄▄▄  ▄▄▀",
		"         ▀▄▄▄    ▄▀",
		"             ▀▀▄▀",
	}
	fmt.Println("msg:", msg)
	expected := []string{""}
	results := []string{""}
	Assert(expected, results, test)
}

func TestFlipHorizontal(test *testing.T) {
	msg := []string{
		"    [49m     [38;5;16m▄[48;5;16m[38;5;232m▄ [49m         [38;5;16m▄▄",
		"        ▄[48;5;16m[38;5;94m▄[48;5;232m▄[48;5;16m [49m    [38;5;16m▄▄▄▄[48;5;16m[38;5;214m▄[48;5;214m[38;5;94m▄[48;5;94m [48;5;16m▄[49m[38;5;16m▄",
		"       ▄[48;5;16m [48;5;94m [48;5;58m▄[49m▀ ▄[48;5;16m[38;5;214m▄▄[48;5;232m  [38;5;94m▄[48;5;214m▄[48;5;94m   [38;5;16m▄[49m▀",
		"     ▄[48;5;16m[38;5;214m▄[48;5;94m▄[48;5;214m   [48;5;16m▄[38;5;58m▄[48;5;214m  [38;5;232m▄[48;5;232m[38;5;94m▄[48;5;94m    [38;5;16m▄[49m▀",
		"    ▄[48;5;16m[38;5;214m▄[48;5;214m   [38;5;94m▄[48;5;94m[38;5;231m▄[48;5;214m[38;5;16m▄  [48;5;58m[38;5;214m▄[48;5;16m [49m[38;5;16m▀[48;5;94m▄  [48;5;16m[38;5;94m▄[49m[38;5;16m▄",
		"    ▀[48;5;214m▄[48;5;58m[38;5;214m▄[48;5;214m   [48;5;16m▄[48;5;232m[38;5;196m▄[48;5;214m▄   [48;5;16m[38;5;214m▄[49m[38;5;16m▄[48;5;16m[38;5;94m▄[48;5;94m [38;5;16m▄[49m▀",
		"    ▀[48;5;94m▄[48;5;232m▄[48;5;94m▄[48;5;214m[38;5;94m▄▄ [48;5;196m[38;5;214m▄[38;5;232m▄[48;5;214m   [48;5;88m[38;5;214m▄[48;5;232m▄[48;5;52m[38;5;232m▄[48;5;16m[38;5;52m▄[49m[38;5;16m▄",
		"        [48;5;16m [48;5;94m  [48;5;232m[38;5;94m▄[48;5;214m[38;5;232m▄▄[48;5;232m[38;5;214m▄[48;5;214m  [48;5;88m▄[48;5;232m[38;5;16m▄[49m▀",
		"         ▀[48;5;94m▄▄[48;5;214m▄    [48;5;232m▄[49m▀",
		"             ▀▀[48;5;214m▄[49m▀[39m[39m",
	}
	fmt.Println("msg:", msg)
	expected := []string{""}
	results := []string{""}
	Assert(expected, results, test)
}
