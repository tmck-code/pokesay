package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/tmck-code/pokesay/src/pokesay"
)

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
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "         "},
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "▄▄"},
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "          "},
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "▄▄"},
				},
			},
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			result := pokesay.TokeniseANSIString(tc.input)
			if Debug() {
				fmt.Printf("input: 	  '%v\x1b[0m'\n", tc.input)
				fmt.Printf("expected: '%v\x1b[0m'\n", tc.expected)
				fmt.Printf("result:   '%v\x1b[0m'\n", result)
			}
			for i, line := range tc.expected {
				Assert(line, result[i], t)
			}
			Assert(tc.expected, result, t)
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
		test.Run(tc.name, func(t *testing.T) {
			result := pokesay.ReverseUnicodeString(tc.input)
			if Debug() {
				fmt.Printf("input: 	  '%v\x1b[0m'\n", tc.input)
				fmt.Printf("expected: '%v\x1b[0m'\n", tc.expected)
				fmt.Printf("result:   '%v\x1b[0m'\n", result)
			}
			Assert(tc.expected, result, t)
		})
	}
}

func TestANSITokenise(test *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]pokesay.ANSILineToken
	}{
		{
			// purple fg, red bg
			name:  "Single line with fg and bg",
			input: "\x1b[38;5;129mAAA\x1b[48;5;160mXX",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{Text: "AAA", FGColour: "\x1b[38;5;129m", BGColour: ""},
					pokesay.ANSILineToken{Text: "XX", FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m"},
				},
			},
		},
		{
			// purple fg, red bg
			name:  "Longer single line with fg and bg",
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{Text: "AAA    ", FGColour: "\x1b[38;5;129m", BGColour: ""},
					pokesay.ANSILineToken{Text: " XX ", FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m"},
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
				},
				{ // Line 2
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;46m", BGColour: "", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;190m", BGColour: "", Text: "▄"},
				},
			},
		},
		// purple fg, red bg
		{
			name:  "Single line with spaces",
			input: "\x1b[38;5;129mAAA  \x1b[48;5;160m  XX  \x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{Text: "AAA  ", FGColour: "\x1b[38;5;129m", BGColour: ""},
					pokesay.ANSILineToken{Text: "  XX  ", FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m"},
				},
			},
		},
		{
			name:  "Single line with existing ANSI reset",
			input: "\x1b[38;5;129mAAA\x1b[48;5;160mXX\x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{Text: "AAA", FGColour: "\x1b[38;5;129m", BGColour: ""},
					pokesay.ANSILineToken{Text: "XX", FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m"},
				},
			},
		},
		{
			name: "Top of Egg",
			input: "    \x1b[49m   \x1b[38;5;16m▄▄\x1b[48;5;16m\x1b[38;5;142m▄▄▄\x1b[49m\x1b[38;5;16m▄▄\n" +
				"     ▄\x1b[48;5;16m\x1b[38;5;58m▄\x1b[48;5;58m\x1b[38;5;70m▄\x1b[48;5;70m \x1b[48;5;227m    \x1b[48;5;237m\x1b[38;5;227m▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "    "},
					pokesay.ANSILineToken{FGColour: "", BGColour: "\x1b[49m", Text: "   "},

					pokesay.ANSILineToken{FGColour: "\u001b[38;5;16m", BGColour: "\x1b[49m", Text: "▄▄"},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;142m", BGColour: "\u001b[48;5;16m", Text: "▄▄▄"},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;16m", BGColour: "\x1b[49m", Text: "▄▄"},
				},
				{
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;16m", BGColour: "\x1b[49m", Text: "     ▄"},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;58m", BGColour: "\u001b[48;5;16m", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;70m", BGColour: "\u001b[48;5;58m", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;70m", BGColour: "\u001b[48;5;70m", Text: " "},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;70m", BGColour: "\u001b[48;5;227m", Text: "    "},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;227m", BGColour: "\u001b[48;5;237m", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;237m", BGColour: "\u001b[48;5;16m", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\u001b[38;5;16m", BGColour: "\u001b[49m", Text: "▄"},
				},
			},
		},
		{
			name: "Multi-line with trailing spaces",
			// The AAA has a purple fg
			// The XX has a red bg
			input: "  \x1b[38;5;129mAAA \x1b[48;5;160m XY \x1b[0m     ",
			// The AAA should still have a purple fg
			// The XX should still have a red bg
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "  "},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA "},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: " XY "},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: "     "},
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
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA    "},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: " XX "},
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
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: "AAA    "},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: " XX "},
				},
			},
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			result := pokesay.TokeniseANSIString(tc.input)

			fmt.Printf("input: 	  '\n%s\x1b[0m'\n", tc.input)
			fmt.Printf("expected:   '\n%s\n", pokesay.BuildANSIString(tc.expected))
			fmt.Printf("result:   '\n%s\n", pokesay.BuildANSIString(result))
			for i, line := range tc.expected {
				if Debug() {
					fmt.Printf("expected: %+v\x1b[0m\n", line)
					fmt.Printf("  result: %+v\x1b[0m\n", result[i])
					rb, err := json.MarshalIndent(result[i], "", "  ")
					if err != nil {
						fmt.Println("error:", err)
					}
					fmt.Printf("  result: %+v\x1b[0m\n", string(rb))
					eb, err := json.MarshalIndent(line, "", "  ")
					if err != nil {
						fmt.Println("error:", err)
					}
					fmt.Printf("expected: %+v\x1b[0m\n", string(eb))
					for j, token := range result[i] {
						Assert(line[j], token, t)
						if (j + 1) < len(line) {
							break
						}
					}
					Assert(line, result[i], t)
				}
			}
			Assert(tc.expected, result, t)
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
		expected [][]pokesay.ANSILineToken
	}{
		{
			name: "Single line with ANSI colours",
			// The AAA has a purple fg, and the XX has a red bg
			input: "\x1b[38;5;129mAAA \x1b[48;5;160m XX \x1b[0m",
			// expected: "\x1b[0m\x1b[38;5;129m\x1b[48;5;160m XX \x1b[49m AAA\x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					{"", "", ""},
					{"\x1b[38;5;129m", "\x1b[48;5;160m", " XX "},
					{"\x1b[38;5;129m", "", " AAA"},
				},
			},
		},
		{
			name: "Multi-line with ANSI colours",
			// purple fg, red bg
			// the 4 spaces after AAA should have a purple fg, and no bg
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			// expected: "\x1b[0m\x1b[38;5;129m\x1b[48;5;160m XX \x1b[0m\x1b[38;5;129m    AAA\x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					{"", "", ""},
					{"\x1b[38;5;129m", "\x1b[48;5;160m", " XX "},
					{"\x1b[38;5;129m", "", "    AAA"},
				},
			},
		},
		{
			name: "Multi-line with trailing spaces",
			// The AAA has a purple fg, the XX has a red bg
			input: "  \x1b[38;5;129mAAA \x1b[48;5;160m XY \x1b[0m  ",
			// The AAA should still have a purple fg, and the XX should still have a red bg
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: ""},
					pokesay.ANSILineToken{FGColour: "\x1b[0m", BGColour: "", Text: "  "},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "\x1b[48;5;160m", Text: " YX "},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;129m", BGColour: "", Text: " AAA"},
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: "  "},
				},
			},
		},
		{
			name:  "Multi-line with colour continuation",
			input: "\x1b[38;5;160m▄ \x1b[38;5;46m▄\n▄ \x1b[38;5;190m▄",
			// expected: "\x1b[0m\x1b[38;5;46m▄\x1b[38;5;160m ▄\n\x1b[38;5;190m▄\x1b[38;5;46m ▄\x1b[0m",
			expected: [][]pokesay.ANSILineToken{
				{
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: ""},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;46m", BGColour: "", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;160m", BGColour: "", Text: " ▄"},
				},
				{
					pokesay.ANSILineToken{FGColour: "", BGColour: "", Text: ""},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;190m", BGColour: "", Text: "▄"},
					pokesay.ANSILineToken{FGColour: "\x1b[38;5;46m", BGColour: "", Text: " ▄"},
				},
			},
		},
		{
			// 	// Test ANSI pokemon reversal --------------------------------------------------
			// 	//
			// 	// These are larger "integration" tests for reversing ANSI strings.
			// 	// - reverse pokemon sprite (with & without ANSI colours)
			name: "flip pikachu without colour",
			input: strings.Join(
				[]string{
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
				},
				"\n",
			),
			expected: [][]pokesay.ANSILineToken{
				{{"", "", "  "}, {"", "", "▄▄          ▄▄         "}},
				{{"", "", ""}, {"", "", "▄▄ ▄▄▄▄▄▄     ▄▄▄        "}},
				{{"", "", ""}, {"", "", "▀▄   ▄▄  ▄▄▄ ▀▄  ▄       "}},
				{{"", "", "  "}, {"", "", "▀▄    ▄▄  ▄▄   ▄▄▄     "}},
				{{"", "", "   "}, {"", "", "▄▄  ▄▀ ▄  ▄▄▄   ▄▄    "}},
				{{"", "", "   "}, {"", "", "▀▄ ▄▄▄   ▄▄▄   ▄▄▀    "}},
				{{"", "", "    "}, {"", "", "▄▄▄▄▄   ▄▄ ▄▄▄▄▄▀    "}},
				{{"", "", "     "}, {"", "", "▀▄▄  ▄▄▄▄           "}},
				{{"", "", "      "}, {"", "", "▀▄    ▄▄▄▀         "}},
				{{"", "", "        "}, {"", "", "▀▄▀▀             "}},
			},
		},
		{
			name: "flip pikachu with colour",
			input: strings.Join(
				[]string{
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
				},
				"\n",
			),
			// expected: strings.Join(
			// 	[]string{
			// 		"  \x1b[38;5;16m▄▄\x1b[0m         \x1b[0m\x1b[48;5;16m\x1b[38;5;232m ▄\x1b[0m\x1b[38;5;16m▄\x1b[0m     \x1b[0m    ",
			// 		"\x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m \x1b[0m\x1b[48;5;214m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[38;5;16m▄▄▄▄\x1b[0m    \x1b[0m\x1b[48;5;16m\x1b[38;5;94m \x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m▄        ",
			// 		"\x1b[49m▀\x1b[0m\x1b[48;5;94m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m   \x1b[0m\x1b[48;5;214m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;232m\x1b[38;5;214m  \x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄▄\x1b[0m▄ ▀\x1b[0m\x1b[48;5;58m▄\x1b[0m\x1b[48;5;94m \x1b[0m\x1b[48;5;16m \x1b[0m▄       ",
			// 		"  \x1b[49m▀\x1b[0m\x1b[48;5;94m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m    \x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;232m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;58m  \x1b[0m\x1b[48;5;16m\x1b[38;5;58m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m   \x1b[0m\x1b[48;5;94m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m▄     ",
			// 		"   \x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;16m  ▄\x1b[0m\x1b[38;5;16m▀\x1b[0m\x1b[48;5;16m\x1b[38;5;214m \x1b[0m\x1b[48;5;58m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;16m  ▄\x1b[0m\x1b[48;5;94m\x1b[38;5;231m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m   \x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m▄    ",
			// 		"   \x1b[49m▀\x1b[0m\x1b[48;5;94m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;94m\x1b[38;5;94m \x1b[0m\x1b[48;5;16m\x1b[38;5;94m▄\x1b[0m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;196m   ▄\x1b[0m\x1b[48;5;232m\x1b[38;5;196m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m   \x1b[0m\x1b[48;5;58m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m▄\x1b[0m▀    ",
			// 		"    \x1b[38;5;16m▄\x1b[0m\x1b[48;5;16m\x1b[38;5;52m▄\x1b[0m\x1b[48;5;52m\x1b[38;5;232m▄\x1b[0m\x1b[48;5;232m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;88m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;232m   \x1b[0m\x1b[48;5;196m\x1b[38;5;232m▄\x1b[0m\x1b[48;5;196m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;94m ▄▄\x1b[0m\x1b[48;5;94m▄\x1b[0m\x1b[48;5;232m▄\x1b[0m\x1b[48;5;94m▄\x1b[0m▀    ",
			// 		"     \x1b[49m▀\x1b[0m\x1b[48;5;232m\x1b[38;5;16m▄\x1b[0m\x1b[48;5;88m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;214m  \x1b[0m\x1b[48;5;232m\x1b[38;5;214m▄\x1b[0m\x1b[48;5;214m\x1b[38;5;232m▄▄\x1b[0m\x1b[48;5;232m\x1b[38;5;94m▄\x1b[0m\x1b[48;5;94m  \x1b[0m\x1b[48;5;16m \x1b[0m        ",
			// 		"      \x1b[49m▀\x1b[0m\x1b[48;5;232m▄\x1b[0m\x1b[48;5;214m    ▄\x1b[0m\x1b[48;5;94m▄▄\x1b[0m▀         ",
			// 		"        \x1b[39m\x1b[0m▀\x1b[0m\x1b[48;5;214m▄\x1b[0m▀▀             ",
			// 	},
			// 	"\n",
			// ),
			expected: [][]pokesay.ANSILineToken{
				{{"", "", "  "}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄"}, {"\x1b[38;5;232m", "\x1b[49m", "         "}, {"\x1b[38;5;232m", "\x1b[48;5;16m", " ▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"", "\x1b[49m", "     "}, {"", "", "    "}},
				{{"", "", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", " "}, {"\x1b[38;5;94m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄▄▄"}, {"\x1b[38;5;94m", "\x1b[49m", "    "}, {"\x1b[38;5;94m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄        "}},
				{{"", "", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", "   "}, {"\x1b[38;5;94m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;232m", "  "}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄ ▀"}, {"\x1b[38;5;16m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", " "}, {"\x1b[38;5;16m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;16m", "\x1b[49m", "▄       "}},
				{{"", "", "  "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", "    "}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;214m", "  "}, {"\x1b[38;5;58m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;214m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄     "}},
				{{"", "", "   "}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "  ▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;214m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "  ▄"}, {"\x1b[38;5;231m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄    "}},
				{{"", "", "   "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", " "}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;196m", "\x1b[48;5;214m", "   ▄"}, {"\x1b[38;5;196m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;214m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀    "}},
				{{"", "", "    "}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;52m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;52m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;88m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;232m", "\x1b[48;5;196m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;196m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;214m", " ▄▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀    "}},
				{{"", "", "     "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;88m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "  "}, {"\x1b[38;5;214m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;214m", "▄▄"}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "  "}, {"\x1b[38;5;16m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;16m", "\x1b[49m", "        "}},
				{{"", "", "      "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "    ▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀         "}},
				{{"", "", "        "}, {"\x1b[39m", "\x1b[49m", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀▀             "}},
			},
		},
		{
			name: "flip egg with colour",
			input: strings.Join(
				[]string{
					"     \x1b[49m   \x1b[38;5;16m▄▄\x1b[48;5;16m\x1b[38;5;142m▄▄▄\x1b[49m\x1b[38;5;16m▄▄",
					"      ▄\x1b[48;5;16m\x1b[38;5;58m▄\x1b[48;5;58m\x1b[38;5;70m▄\x1b[48;5;70m \x1b[48;5;227m    \x1b[48;5;237m\x1b[38;5;227m▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄",
					"     ▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[48;5;70m\x1b[38;5;227m▄▄\x1b[48;5;227m    \x1b[38;5;70m▄▄\x1b[48;5;142m \x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄",
					"     \x1b[48;5;16m \x1b[48;5;227m       \x1b[48;5;70m\x1b[38;5;227m▄\x1b[38;5;58m▄\x1b[48;5;58m \x1b[48;5;142m \x1b[48;5;16m \x1b[49m",
					"     \x1b[48;5;16m \x1b[48;5;142m\x1b[38;5;237m▄\x1b[48;5;227m\x1b[38;5;142m▄\x1b[48;5;70m  \x1b[48;5;227m▄▄\x1b[38;5;58m▄\x1b[48;5;142m▄▄ \x1b[38;5;237m▄\x1b[48;5;16m \x1b[49m",
					"      \x1b[48;5;16m \x1b[48;5;142m▄   \x1b[48;5;58m    \x1b[38;5;234m▄\x1b[48;5;16m \x1b[49m",
					"       \x1b[38;5;16m▀▀\x1b[48;5;142m▄▄▄\x1b[48;5;58m▄▄\x1b[49m▀▀\x1b[39m\x1b[39m",
				},
				"\n",
			),
			// expected: strings.Join(
			// 	[]string{
			// 		"        \x1b[49m\x1b[38;5;16m\x1b[49m▄▄\x1b[48;5;16m\x1b[38;5;142m\x1b[48;5;16m▄▄▄\x1b[38;5;16m\x1b[49m▄▄\x1b[49m    \x1b[49m \x1b[0m",
			// 		"      \x1b[49m\x1b[38;5;16m\x1b[49m▄\x1b[48;5;16m\x1b[38;5;237m\x1b[48;5;16m▄\x1b[48;5;237m\x1b[38;5;227m\x1b[48;5;237m▄\x1b[48;5;58m\x1b[38;5;70m\x1b[48;5;227m    \x1b[48;5;58m\x1b[38;5;70m\x1b[48;5;70m \x1b[48;5;58m\x1b[38;5;70m\x1b[48;5;58m▄\x1b[48;5;16m\x1b[38;5;58m\x1b[48;5;16m▄\x1b[49m\x1b[38;5;16m\x1b[49m▄\x1b[49m \x1b[0m",
			// 		"     \x1b[49m\x1b[38;5;16m\x1b[49m▄\x1b[48;5;16m\x1b[38;5;237m\x1b[48;5;16m▄\x1b[38;5;70m\x1b[48;5;142m \x1b[38;5;70m\x1b[48;5;227m▄▄\x1b[48;5;70m\x1b[38;5;227m\x1b[48;5;227m    \x1b[48;5;70m\x1b[38;5;227m\x1b[48;5;70m▄▄\x1b[48;5;16m\x1b[38;5;237m\x1b[48;5;16m▄\x1b[49m\x1b[38;5;16m\x1b[49m▄\x1b[49m \x1b[0m",
			// 		"     \x1b[38;5;58m\x1b[49m\x1b[38;5;58m\x1b[48;5;16m \x1b[38;5;58m\x1b[48;5;142m \x1b[38;5;58m\x1b[48;5;58m \x1b[38;5;58m\x1b[48;5;70m▄\x1b[48;5;70m\x1b[38;5;227m\x1b[48;5;70m▄\x1b[49m\x1b[38;5;16m\x1b[48;5;227m       \x1b[49m\x1b[38;5;16m\x1b[48;5;16m \x1b[49m\x1b[38;5;16m\x1b[49m \x1b[49m \x1b[0m",
			// 		"     \x1b[38;5;237m\x1b[49m\x1b[38;5;237m\x1b[48;5;16m \x1b[38;5;237m\x1b[48;5;142m▄\x1b[38;5;58m\x1b[48;5;142m ▄▄\x1b[38;5;58m\x1b[48;5;227m▄\x1b[48;5;227m\x1b[38;5;142m\x1b[48;5;227m▄▄\x1b[48;5;227m\x1b[38;5;142m\x1b[48;5;70m  \x1b[48;5;227m\x1b[38;5;142m\x1b[48;5;227m▄\x1b[48;5;142m\x1b[38;5;237m\x1b[48;5;142m▄\x1b[38;5;58m\x1b[48;5;16m \x1b[38;5;58m\x1b[49m \x1b[49m \x1b[0m",
			// 		"      \x1b[38;5;234m\x1b[49m\x1b[38;5;234m\x1b[48;5;16m \x1b[38;5;234m\x1b[48;5;58m▄\x1b[38;5;237m\x1b[48;5;58m    \x1b[38;5;237m\x1b[48;5;142m   ▄\x1b[38;5;237m\x1b[48;5;16m \x1b[38;5;237m\x1b[49m \x1b[49m \x1b[0m",
			// 		"       \x1b[39m\x1b[39m\x1b[49m\x1b[38;5;16m\x1b[49m▀▀\x1b[38;5;16m\x1b[48;5;58m▄▄\x1b[38;5;16m\x1b[48;5;142m▄▄▄\x1b[38;5;16m\x1b[49m▀▀\x1b[38;5;234m\x1b[49m \x1b[49m \x1b[0m",
			// 		"     \x1b[49m                 \x1b[0m",
			// 	},
			// 	"\n",
			// ),
			expected: [][]pokesay.ANSILineToken{
				{{"", "", "   "}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄"}, {"\x1b[38;5;142m", "\x1b[48;5;16m", "▄▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄"}, {"", "\x1b[49m", "   "}, {"", "", "     "}},
				{{"", "", " "}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;227m", "\x1b[48;5;237m", "▄"}, {"\x1b[38;5;70m", "\x1b[48;5;227m", "    "}, {"\x1b[38;5;70m", "\x1b[48;5;70m", " "}, {"\x1b[38;5;70m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄      "}},
				{{"", "", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;70m", "\x1b[48;5;142m", " "}, {"\x1b[38;5;70m", "\x1b[48;5;227m", "▄▄"}, {"\x1b[38;5;227m", "\x1b[48;5;227m", "    "}, {"\x1b[38;5;227m", "\x1b[48;5;70m", "▄▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄     "}},
				{{"", "", ""}, {"\x1b[38;5;58m", "\x1b[49m", ""}, {"\x1b[38;5;58m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;58m", "\x1b[48;5;142m", " "}, {"\x1b[38;5;58m", "\x1b[48;5;58m", " "}, {"\x1b[38;5;58m", "\x1b[48;5;70m", "▄"}, {"\x1b[38;5;227m", "\x1b[48;5;70m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;227m", "       "}, {"\x1b[38;5;16m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;16m", "\x1b[49m", "     "}},
				{{"", "", ""}, {"\x1b[38;5;237m", "\x1b[49m", ""}, {"\x1b[38;5;237m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;237m", "\x1b[48;5;142m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;142m", " ▄▄"}, {"\x1b[38;5;58m", "\x1b[48;5;227m", "▄"}, {"\x1b[38;5;142m", "\x1b[48;5;227m", "▄▄"}, {"\x1b[38;5;142m", "\x1b[48;5;70m", "  "}, {"\x1b[38;5;142m", "\x1b[48;5;227m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;142m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;58m", "\x1b[49m", "     "}},
				{{"", "", " "}, {"\x1b[38;5;234m", "\x1b[49m", ""}, {"\x1b[38;5;234m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;234m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;58m", "    "}, {"\x1b[38;5;237m", "\x1b[48;5;142m", "   ▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;237m", "\x1b[49m", "      "}},
				{{"", "", "  "}, {"\x1b[39m", "\x1b[49m", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▀▀"}, {"\x1b[38;5;16m", "\x1b[48;5;58m", "▄▄"}, {"\x1b[38;5;16m", "\x1b[48;5;142m", "▄▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀▀"}, {"\x1b[38;5;234m", "\x1b[49m", "       "}},
			},
		},
	}
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			result := pokesay.ReverseANSIString(pokesay.TokeniseANSIString(tc.input))

			fmt.Printf("input: 	  '\n%s\x1b[0m'\n", tc.input)
			fmt.Printf("expected:   '\n%s\n", pokesay.BuildANSIString(tc.expected))
			fmt.Printf("result:   '\n%s\n", pokesay.BuildANSIString(result))
			for i, line := range tc.expected {
				if Debug() {
					fmt.Printf("expected: %+v\x1b[0m\n", line)
					fmt.Printf("  result: %+v\x1b[0m\n", result[i])

					eb, err := json.MarshalIndent(line, "", "  ")
					if err != nil {
						fmt.Println("error:", err)
					}
					fmt.Printf("expected: %+v\x1b[0m\n", string(eb))
					rb, err := json.MarshalIndent(result[i], "", "  ")
					if err != nil {
						fmt.Println("error:", err)
					}
					fmt.Printf("  result: %+v\x1b[0m\n", string(rb))
					for j, token := range result[i] {
						Assert(line[j], token, t)
						if (j + 1) < len(line) {
							break
						}
					}
					Assert(line, result[i], t)
				}
			}
			Assert(tc.expected, result, t)
		})
	}
}
