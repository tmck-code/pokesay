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

// Test pokemon selection algorithms -------------------------------------------

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

// Test unicode helpers --------------------------------------------------------

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

// Test ANSI tokenisation ------------------------------------------------------

func TestUnicodeTokenise(test *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]pokesay.ANSILineToken
	}{
		{
			name:  "Single line with no colour",
			input: "         ▄▄          ▄▄",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "         ", Reset: true},
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "▄▄"},
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "          ", Reset: true},
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "▄▄"},
				},
			},
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(test *testing.T) {
			result := pokesay.TokeniseANSIString(tc.input)
			if Debug() {
				fmt.Printf("input: 	  '%v\x1b[0m'\n", tc.input)
				fmt.Printf("expected: '%v\x1b[0m'\n", tc.expected)
				fmt.Printf("result:   '%v\x1b[0m'\n", result)
			}
			for i, line := range tc.expected {
				Assert(line, result[i], test)
			}
			Assert(tc.expected, result, test)
		})
	}
}

func TestUnicodeReverse(test *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Reverse basic ANSI string with no colour",
			input:    "         ▄▄          ▄▄",
			expected: "▄▄          ▄▄         ",
		},
		{
			name:     "Reverse basic ANSI string with no colour and trailing spaces",
			input:    "         ▄▄          ▄▄      ",
			expected: "      ▄▄          ▄▄         ",
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(test *testing.T) {
			result := pokesay.ReverseUnicodeString(tc.input)
			if Debug() {
				fmt.Printf("input: 	  '%v\x1b[0m'\n", tc.input)
				fmt.Printf("expected: '%v\x1b[0m'\n", tc.expected)
				fmt.Printf("result:   '%v\x1b[0m'\n", result)
			}
			Assert(tc.expected, result, test)
		})
	}
}

func TestANSITokenise(test *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]pokesay.ANSILineToken
	}{
		// purple fg, red bg
		{
			name:  "Single line with fg and bg",
			input: "\x1b[38;5;129mAAA\x1b[48;5;160mXX",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: "XX"},
				},
			},
		},
		{
			name: "Multi-line",
			// line 1 : purple fg,                  line 2: red bg
			input: "\x1b[38;5;160m▄\x1b[38;5;46m▄\n▄\x1b[38;5;190m▄",
			expected: [][]pokesay.ANSILineToken{
				{ // Line 1
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;160m", BGColour: "", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;46m", BGColour: "", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: ""},
				},
				{ // Line 2
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;46m", BGColour: "", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;190m", BGColour: "", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: ""},
				},
			},
		},
		// purple fg, red bg
		{
			name:  "Single line with spaces",
			input: "\x1b[38;5;129mAAA  \x1b[48;5;160m  XX  \x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: "XX"},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: ""},
				},
			},
		},
		{
			name:  "Single line with existing ANSI reset",
			input: "\x1b[38;5;129mAAA\x1b[48;5;160mXX\x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: "XX"},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: ""},
				},
			},
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(test *testing.T) {
			result := pokesay.TokeniseANSIString(tc.input)
			if Debug() {
				fmt.Printf("input: 	  '%v\x1b[0m'\n", tc.input)
				fmt.Printf("expected: '%v\x1b[0m'\n", tc.expected)
				fmt.Printf("result:   '%v\x1b[0m'\n", result)
			}
			Assert(tc.expected, result, test)
		})
	}
}

func TestANSIWithSpaces(test *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]pokesay.ANSILineToken
	}{
		{
			name: "Multi-line with trailing spaces",
			// The AAA has a purple fg
			// The XX has a red bg
			input: "  \x1b[38;5;129mAAA \x1b[48;5;160m XY \x1b[0m     ",
			// The AAA should still have a purple fg
			// The XX should still have a red bg
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "  ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: " XY ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: "     ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: ""},
				},
			},
		},
		{
			name: "Lines with FG continuation (spaces)",
			// purple fg, red bg
			// the 4 spaces after AAA should have a purple fg, and no bg
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			// expected := "\x1b[0m\x1b[48;5;160m\x1b[38;5;129m XX \x1b[38;5;129m\x1b[49m    AAA\x1b[0m"
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA    ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: " XX ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: ""},
				},
			},
		},
		{
			name: "Lines with BG continuation (spaces)",
			// purple fg, red bg
			// the 4 spaces after AAA should have a purple fg, and no bg
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			// expected := "\x1b[0m\x1b[48;5;160m\x1b[38;5;129m XX \x1b[38;5;129m\x1b[49m    AAA\x1b[0m"
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA    ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: " XX ", Reset: true},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: ""},
				},
			},
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(test *testing.T) {
			result := pokesay.TokeniseANSIString(tc.input)
			if Debug() {
				fmt.Printf("input: 	  '%v\x1b[0m'\n", tc.input)
				fmt.Printf("expected: '%v\x1b[0m'\n", tc.expected)
				fmt.Printf("result:   '%v\x1b[0m'\n", result)
			}
			Assert(tc.expected, result, test)
		})
	}
}

// Test ANSI line reversal -----------------------------------------------------
//
// These are smaller "unit" tests for ANSI line reversal.
// - reverse individual lines
// - reverse multiple (newline separated) lines

func TestReverseANSIString(test *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Single line with ANSI colours",
			// The AAA has a purple fg, and the XX has a red bg
			input:    "\x1b[38;5;129mAAA \x1b[48;5;160m XX \x1b[0m",
			expected: "\x1b[0m\x1b[38;5;129m\x1b[48;5;160m XX \x1b[38;5;129m AAA\x1b[0m",
		},
		{
			name: "Multi-line with ANSI colours",
			// purple fg, red bg
			// the 4 spaces after AAA should have a purple fg, and no bg
			input:    "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			expected: "\x1b[0m\x1b[38;5;129m\x1b[48;5;160m XX \x1b[0m\x1b[38;5;129m    AAA\x1b[0m",
		},
		{
			name: "Multi-line with trailing spaces",
			// The AAA has a purple fg, the XX has a red bg
			input: "  \x1b[38;5;129mAAA \x1b[48;5;160m XY \x1b[0m  ",
			// The AAA should still have a purple fg, and the XX should still have a red bg
			expected: "\x1b[0m\x1b[0m  \x1b[38;5;129m\x1b[48;5;160m YX \x1b[0m\x1b[38;5;129m AAA  \x1b[0m",
		},
		{
			name:     "Multi-line with colour continuation",
			input:    "\x1b[38;5;160m▄ \x1b[38;5;46m▄\n▄ \x1b[38;5;190m▄",
			expected: "\x1b[0m\x1b[38;5;46m▄\x1b[38;5;160m ▄\x1b[0m\n\x1b[0m\x1b[38;5;190m▄\x1b[38;5;46m ▄\x1b[0m",
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(test *testing.T) {
			result := pokesay.ReverseANSIString(tc.input)
			if Debug() {
				fmt.Printf("input: 	  '%v\x1b[0m'\n", tc.input)
				fmt.Printf("expected: '%v\x1b[0m'\n", tc.expected)
				fmt.Printf("result:   '%v\x1b[0m'\n", result)
			}
			Assert(tc.expected, result, test)
		})
	}
}

// Test ANSI pokemon reversal --------------------------------------------------
//
// These are larger "integration" tests for reversing ANSI strings.
// - reverse pokemon sprite (with & without ANSI colours)

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
		"  ▄▄          ▄▄         \x1b[0m",
		"▄▄ ▄▄▄▄▄▄     ▄▄▄        \x1b[0m",
		"▀▄   ▄▄  ▄▄▄ ▀▄  ▄       \x1b[0m",
		"  ▀▄    ▄▄  ▄▄   ▄▄▄     \x1b[0m",
		"   ▄▄  ▄▀ ▄  ▄▄▄   ▄▄    \x1b[0m",
		"   ▀▄ ▄▄▄   ▄▄▄   ▄▄▀    \x1b[0m",
		"    ▄▄▄▄▄   ▄▄ ▄▄▄▄▄▀    \x1b[0m",
		"     ▀▄▄  ▄▄▄▄           \x1b[0m",
		"      ▀▄    ▄▄▄▀         \x1b[0m",
		"        ▀▄▀▀             \x1b[0m",
	}
	results := pokesay.ReverseANSIString(strings.Join(msg, "\n"))

	splitResults := strings.Split(results, "\n")
	for i := 0; i < len(expected); i++ {
		Assert(expected[i], splitResults[i], test)
	}

	Assert(strings.Join(expected, "\n"), results, test)
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
	results := pokesay.ReverseANSIString(strings.Join(msg, "\n"))

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
	splitResults := strings.Split(results, "\n")
	for i, line := range splitResults {
		fmt.Println("results:", i, line)
	}
	for i := 0; i < len(expected); i++ {
		Assert(expected[i], splitResults[i], test)
	}
	Assert(expected, results, test)
}
