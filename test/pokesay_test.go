package test

import (
	"embed"
	"fmt"
	"strings"
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

func TestTokeniseANSIString(test *testing.T) {
	line := "\x1b[38;5;129mAAA \x1b[48;5;160m XX \x1b[0m"

	expected := [][]pokesay.ANSILineToken{
		{
			pokesay.ANSILineToken{Colour: "\x1b[38;5;129m", Text: "AAA "},
			pokesay.ANSILineToken{Colour: "\x1b[48;5;160m\x1b[38;5;129m", Text: " XX "},
			pokesay.ANSILineToken{Colour: "\x1b[0m", Text: ""},
		},
	}
	result := pokesay.TokeniseANSIString(line)
	Assert(expected, result, test)
}

func TestTokeniseANSIStringLines(test *testing.T) {
	msg := "\x1b[38;5;160m▄ \x1b[38;5;46m▄\n▄ \x1b[38;5;190m▄"

	expected := [][]pokesay.ANSILineToken{
		{ // Line 1
			pokesay.ANSILineToken{Colour: "\x1b[38;5;160m", Text: "▄ "},
			pokesay.ANSILineToken{Colour: "\x1b[38;5;46m", Text: "▄"},
			pokesay.ANSILineToken{Colour: "\x1b[0m", Text: ""},
		},
		{ // Line 2
			pokesay.ANSILineToken{Colour: "\x1b[38;5;46m", Text: "▄ "},
			pokesay.ANSILineToken{Colour: "\x1b[38;5;190m", Text: "▄"},
			pokesay.ANSILineToken{Colour: "\x1b[0m", Text: ""},
		},
	}
	result := pokesay.TokeniseANSIString(msg)
	Assert(expected, result, test)

}

func TestFlipHorizontalLine(test *testing.T) {
	// The AAA has a purple fg
	// The XX has a red bg
	line := "\x1b[38;5;129mAAA \x1b[48;5;160m XY \x1b[0m"

	// The AAA should still have a purple fg
	// The XX should still have a red bg
	expected := "\x1b[0m\x1b[48;5;160m\x1b[38;5;129m YX \x1b[0m\x1b[38;5;129m AAA\x1b[0m"
	result := pokesay.ReverseANSIString(line)

	Assert(expected, result, test)
}

func TestTokeniseANSIStringWithNoColour(test *testing.T) {
	msg := "         ▄▄          ▄▄"
	expected := []pokesay.ANSILineToken{
		pokesay.ANSILineToken{Colour: "", Text: "         ▄▄          ▄▄"},
	}
	result := pokesay.TokeniseANSIString(msg)
	Assert(expected, result, test)
}

func TestReverseUnicodeString(test *testing.T) {
	msg := "         ▄▄          ▄▄"
	expected := "▄▄          ▄▄         "
	result := pokesay.ReverseUnicodeString(msg)
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
	expected := []string{
		"  ▄▄          ▄▄         ",
		"▄▄ ▄▄▄▄▄▄     ▄▄▄        ",
		"▀▄   ▄▄  ▄▄▄ ▀▄  ▄       ",
		"  ▀▄    ▄▄  ▄▄   ▄▄▄     ",
		"   ▄▄  ▄▀ ▄  ▄▄▄   ▄▄    ",
		"   ▀▄ ▄▄▄   ▄▄▄   ▄▄▀    ",
		"    ▄▄▄▄▄   ▄▄ ▄▄▄▄▄▀    ",
		"     ▀▄▄  ▄▄▄▄           ",
		"      ▀▄    ▄▄▄▀         ",
		"        ▀▄▀▀             ",
	}
	results := pokesay.ReverseANSIStrings(msg)

	for i := 0; i < len(expected); i++ {
		Assert(expected[i], results[i], test)
	}

	Assert(expected, results, test)
}

func TestFlipHorizontalWithColourContinuation(test *testing.T) {
	msg := []string{
		"\x1b[38;5;160m▄ \x1b[38;5;46m▄",
		"▄ \x1b[38;5;190m▄",
	}

	result := pokesay.ReverseANSIStrings(msg)

	expected := []string{
		"\x1b[38;5;46m▄ \x1b[38;5;160m▄",
		"\x1b[38;5;190m▄ \x1b[38;5;46m▄",
	}

	data := map[string]string{
		"msg":      strings.Join(msg, "\n"),
		"expected": strings.Join(expected, "\n"),
		"result":   strings.Join(result, "\n"),
	}

	for msg, d := range data {
		fmt.Printf("%s:\n%s\x1b[0m\n%#v\n", msg, d, d)
	}

	Assert(expected, result, test)
}

func TestFlipHorizontal(test *testing.T) {
	msg := []string{
		"    \x1b[49m     \x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;232m▄ \x1b[49m         \x1b[38;5;16m▄▄",
		"        ▄\x1b[48;5;16m\x1b[38;5;94m▄\x1b[48;5;232m▄\x1b[48;5;16m \x1b[49m    \x1b[38;5;16m▄▄▄▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;214m\x1b[38;5;94m▄\x1b[48;5;94m \x1b[48;5;16m▄\x1b[49m\x1b[38;5;16m▄",
		"       ▄\x1b[48;5;16m \x1b[48;5;94m \x1b[48;5;58m▄\x1b[49m▀ ▄\x1b[48;5;16m\x1b[38;5;214m▄▄\x1b[48;5;232m  \x1b[38;5;94m▄\x1b[48;5;214m▄\x1b[48;5;94m   \x1b[38;5;16m▄\x1b[49m▀",
		"     ▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;94m▄\x1b[48;5;214m   \x1b[48;5;16m▄\x1b[38;5;58m▄\x1b[48;5;214m  \x1b[38;5;232m▄\x1b[48;5;232m\x1b[38;5;94m▄\x1b[48;5;94m    \x1b[38;5;16m▄\x1b[49m▀",
		"    ▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;214m   \x1b[38;5;94m▄\x1b[48;5;94m\x1b[38;5;231m▄\x1b[48;5;214m\x1b[38;5;16m▄  \x1b[48;5;58m\x1b[38;5;214m▄\x1b[48;5;16m \x1b[49m\x1b[38;5;16m▀\x1b[48;5;94m▄  \x1b[48;5;16m\x1b[38;5;94m▄\x1b[49m\x1b[38;5;16m▄",
		"    ▀\x1b[48;5;214m▄\x1b[48;5;58m\x1b[38;5;214m▄\x1b[48;5;214m   \x1b[48;5;16m▄\x1b[48;5;232m\x1b[38;5;196m▄\x1b[48;5;214m▄   \x1b[48;5;16m\x1b[38;5;214m▄\x1b[49m\x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;94m▄\x1b[48;5;94m \x1b[38;5;16m▄\x1b[49m▀",
		"    ▀\x1b[48;5;94m▄\x1b[48;5;232m▄\x1b[48;5;94m▄\x1b[48;5;214m\x1b[38;5;94m▄▄ \x1b[48;5;196m\x1b[38;5;214m▄\x1b[38;5;232m▄\x1b[48;5;214m   \x1b[48;5;88m\x1b[38;5;214m▄\x1b[48;5;232m▄\x1b[48;5;52m\x1b[38;5;232m▄\x1b[48;5;16m\x1b[38;5;52m▄\x1b[49m\x1b[38;5;16m▄",
		"        \x1b[48;5;16m \x1b[48;5;94m  \x1b[48;5;232m\x1b[38;5;94m▄\x1b[48;5;214m\x1b[38;5;232m▄▄\x1b[48;5;232m\x1b[38;5;214m▄\x1b[48;5;214m  \x1b[48;5;88m▄\x1b[48;5;232m\x1b[38;5;16m▄\x1b[49m▀",
		"         ▀\x1b[48;5;94m▄▄\x1b[48;5;214m▄    \x1b[48;5;232m▄\x1b[49m▀",
		"             ▀▀\x1b[48;5;214m▄\x1b[49m▀\x1b[39m\x1b[39m",
	}
	fmt.Println("msg:", msg)
	results := pokesay.ReverseANSIStrings(msg)

	expected := []string{
		"  \x1b[38;5;16m▄▄\x1b[0m         \x1b[0m\x1b[48;5;16m\x1b[38;5;232m ▄\x1b[0m\x1b[38;5;16m▄\x1b[0m     \x1b[0m    ",
		"\x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m \x1b[0m\x1b[48;5;214m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[38;5;16m▄▄▄▄\x1b[0m    \x1b[0m\x1b[48;5;16m\x1b[38;5;94m \x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m▄        ",
		"\x1b[49m▀\x1b[0m\x1b[48;5;94m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m   \x1b[0m\x1b[48;5;214m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;232m\x1b[38;5;214m  \x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄▄\x1b[0m▄ ▀\x1b[0m\x1b[48;5;58m▄\x1b[0m\x1b[48;5;94m \x1b[0m\x1b[48;5;16m \x1b[0m▄       ",
		"  \x1b[49m▀\x1b[0m\x1b[48;5;94m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m    \x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;232m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;58m  \x1b[0m\x1b[48;5;16m\x1b[38;5;58m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m   \x1b[0m\x1b[48;5;94m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m▄     ",
		"   \x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;16m  ▄\x1b[0m\x1b[38;5;16m▀\x1b[0m\x1b[48;5;16m\x1b[38;5;214m \x1b[0m\x1b[48;5;58m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;16m  ▄\x1b[0m\x1b[48;5;94m\x1b[38;5;231m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m   \x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m▄    ",
		"   \x1b[49m▀\x1b[0m\x1b[48;5;94m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m \x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;196m   ▄\x1b[0m\x1b[48;5;232m\x1b[38;5;196m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m   \x1b[0m\x1b[48;5;58m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m▄\x1b[0m▀    ",
		"    \x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;52m▄\x1b[0m\x1b[48;5;52m\x1b[38;5;232m▄\x1b[0m\x1b[48;5;232m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;88m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;232m   \x1b[0m\x1b[48;5;196m\x1b[38;5;232m▄\x1b[0m\x1b[48;5;196m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;94m ▄▄\x1b[0m\x1b[48;5;94m▄\x1b[0m\x1b[48;5;232m▄\x1b[0m\x1b[48;5;94m▄\x1b[0m▀    ",
		"     \x1b[49m▀\x1b[0m\x1b[48;5;232m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;88m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m  \x1b[0m\x1b[48;5;232m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;232m▄▄\x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;94m  \x1b[0m\x1b[48;5;16m \x1b[0m        ",
		"      \x1b[49m▀\x1b[0m\x1b[48;5;232m▄\x1b[0m\x1b[48;5;214m    ▄\x1b[0m\x1b[48;5;94m▄▄\x1b[0m▀         ",
		"        \x1b[39m\x1b[0m▀\x1b[0m\x1b[48;5;214m▄\x1b[0m▀▀             ",
	}
	for i, line := range msg {
		fmt.Println("msg:", i, line)
	}
	for i, line := range expected {
		fmt.Println("expected:", i, line)
	}
	for i, line := range results {
		fmt.Println("results:", i, line)
	}
	for i := 0; i < len(expected); i++ {
		Assert(expected[i], results[i], test)
	}
	Assert(expected, results, test)
}
