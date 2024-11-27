package pokesay

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	"github.com/tmck-code/pokesay/src/pokedex"
)

type BoxCharacters struct {
	HorizontalEdge    string
	VerticalEdge      string
	TopRightCorner    string
	TopLeftCorner     string
	BottomRightCorner string
	BottomLeftCorner  string
	BalloonString     string
	BalloonTether     string
	Separator         string
	RightArrow        string
	CategorySeparator string
}

type Args struct {
	Width          int
	NoWrap         bool
	DrawBubble     bool
	TabSpaces      string
	NoTabSpaces    bool
	NoCategoryInfo bool
	ListCategories bool
	ListNames      bool
	Category       string
	NameToken      string
	JapaneseName   bool
	BoxCharacters  *BoxCharacters
	DrawInfoBorder bool
	Help           bool
	Verbose        bool
}

var (
	textStyleItalic    *color.Color   = color.New(color.Italic)
	textStyleBold      *color.Color   = color.New(color.Bold)
	resetColourANSI    string         = "\033[0m"
	AsciiBoxCharacters *BoxCharacters = &BoxCharacters{
		HorizontalEdge:    "-",
		VerticalEdge:      "|",
		TopRightCorner:    "\\",
		TopLeftCorner:     "/",
		BottomRightCorner: "/",
		BottomLeftCorner:  "\\",
		BalloonString:     "\\",
		BalloonTether:     "¡",
		Separator:         "|",
		RightArrow:        ">",
		CategorySeparator: "/",
	}
	UnicodeBoxCharacters *BoxCharacters = &BoxCharacters{
		HorizontalEdge:    "─",
		VerticalEdge:      "│",
		TopRightCorner:    "╮",
		TopLeftCorner:     "╭",
		BottomRightCorner: "╯",
		BottomLeftCorner:  "╰",
		BalloonString:     "╲",
		BalloonTether:     "╲",
		Separator:         "│",
		RightArrow:        "→",
		CategorySeparator: "/",
	}
	SingleWidthChars map[string]bool = map[string]bool{
		"♀": true,
		"♂": true,
	}
)

func DetermineBoxCharacters(unicodeBox bool) *BoxCharacters {
	if unicodeBox {
		return UnicodeBoxCharacters
	} else {
		return AsciiBoxCharacters
	}
}

// The main print function! This uses a chosen pokemon's index, names and categories, and an
// embedded filesystem of cowfile data
// 1. The text received from STDIN is printed inside a speech bubble
// 2. The cowfile data is retrieved using the matching index, decompressed (un-gzipped),
// 3. The pokemon is printed along with the name & category information
func Print(args Args, choice int, names []string, categories []string, cows embed.FS, drawBubble bool) {
	printSpeechBubble(args.BoxCharacters, bufio.NewScanner(os.Stdin), args.Width, args.NoTabSpaces, args.TabSpaces, args.NoWrap, drawBubble)
	printPokemon(args, choice, names, categories, cows)
}

// Prints text from STDIN, surrounded by a speech bubble.
func printSpeechBubble(boxCharacters *BoxCharacters, scanner *bufio.Scanner, width int, noTabSpaces bool, tabSpaces string, noWrap bool, drawBubble bool) {
	if drawBubble {
		fmt.Printf(
			"%s%s%s\n",
			boxCharacters.TopLeftCorner,
			strings.Repeat(boxCharacters.HorizontalEdge, width+2),
			boxCharacters.TopRightCorner,
		)
	}

	for scanner.Scan() {
		line := scanner.Text()

		if !noTabSpaces {
			line = strings.Replace(line, "\t", tabSpaces, -1)
		}
		if noWrap {
			printSpeechBubbleLine(boxCharacters, line, width, drawBubble)
		} else {
			printWrappedText(boxCharacters, line, width, tabSpaces, drawBubble)
		}
	}

	bottomBorder := strings.Repeat(boxCharacters.HorizontalEdge, 6) +
		boxCharacters.BalloonTether +
		strings.Repeat(boxCharacters.HorizontalEdge, width+2-7)

	if drawBubble {
		fmt.Printf("%s%s%s\n", boxCharacters.BottomLeftCorner, bottomBorder, boxCharacters.BottomRightCorner)
	} else {
		fmt.Printf(" %s \n", bottomBorder)
	}
	for i := 0; i < 4; i++ {
		fmt.Printf("%s%s\n", strings.Repeat(" ", i+8), boxCharacters.BalloonString)
	}
}

// Prints a single speech bubble line
func printSpeechBubbleLine(boxCharacters *BoxCharacters, line string, width int, drawBubble bool) {
	if !drawBubble {
		fmt.Println(line)
	}

	lineLen := UnicodeStringLength(line)
	if lineLen <= width {
		// print the line with padding, the most common case
		fmt.Printf(
			"%s %s%s%s %s\n",
			boxCharacters.VerticalEdge, // left-hand side of the bubble
			line, resetColourANSI,      // the text
			strings.Repeat(" ", width-lineLen), // padding
			boxCharacters.VerticalEdge,         // right-hand side of the bubble
		)
	} else if lineLen > width {
		// print the line without padding or right-hand side of the bubble if the line is too long
		fmt.Printf(
			"%s %s%s\n",
			boxCharacters.VerticalEdge, // left-hand side of the bubble
			line, resetColourANSI,      // the text
		)
	}
}

// Prints line of text across multiple lines, wrapping it so that it doesn't exceed the desired width.
func printWrappedText(boxCharacters *BoxCharacters, line string, width int, tabSpaces string, drawBubble bool) {
	for _, wline := range strings.Split(wordwrap.WrapString(line, uint(width)), "\n") {
		printSpeechBubbleLine(boxCharacters, wline, width, drawBubble)
	}
}

func nameLength(names []string) int {
	totalLen := 0

	for _, name := range names {
		for _, c := range name {
			// check if ascii or single-width unicode
			if (c < 128) || (SingleWidthChars[string(c)]) {
				totalLen++
			} else {
				totalLen += 2
			}
		}

	}
	return totalLen
}

// Returns the length of a string, taking into account Unicode characters and ANSI escape codes.
func UnicodeStringLength(s string) int {
	nRunes, totalLen, ansiCode := len(s), 0, false

	for i, r := range s {
		if i < nRunes-1 {
			// detect the beginning of an ANSI escape code
			// e.g. "\033[38;5;196m"
			//       ^^^ start    ^ end
			if s[i:i+2] == "\033[" {
				ansiCode = true
			}
		}
		if ansiCode {
			// detect the end of an ANSI escape code
			if r == 'm' {
				ansiCode = false
			}
		} else {
			if r < 128 {
				// if ascii, then use width of 1. this saves some time
				totalLen++
			} else {
				totalLen += runewidth.RuneWidth(r)
			}
		}
	}
	return totalLen
}

// Prints a pokemon with its name & category information.
func printPokemon(args Args, index int, names []string, categoryKeys []string, GOBCowData embed.FS) {
	d, _ := GOBCowData.ReadFile(pokedex.EntryFpath("build/assets/cows", index))

	width := nameLength(names)
	namesFmt := make([]string, 0)
	for _, name := range names {
		namesFmt = append(namesFmt, textStyleBold.Sprint(name))
	}
	// count name separators
	width += (len(names) - 1) * 3
	width += 2     // for the arrow
	width += 2 + 2 // for the end box characters

	infoLine := ""

	if args.NoCategoryInfo {
		infoLine = fmt.Sprintf(
			"%s %s",
			args.BoxCharacters.RightArrow, strings.Join(namesFmt, fmt.Sprintf(" %s ", args.BoxCharacters.Separator)),
		)
	} else {
		infoLine = fmt.Sprintf(
			"%s %s %s %s",
			args.BoxCharacters.RightArrow,
			strings.Join(namesFmt, fmt.Sprintf(" %s ", args.BoxCharacters.Separator)),
			args.BoxCharacters.Separator,
			textStyleItalic.Sprint(strings.Join(categoryKeys, args.BoxCharacters.CategorySeparator)),
		)
		for _, category := range categoryKeys {
			width += len(category)
		}
		width += len(categoryKeys) - 1 + 1 + 2 // lol why did I do this
	}

	if args.DrawInfoBorder {
		topBorder := fmt.Sprintf(
			"%s%s%s",
			args.BoxCharacters.TopLeftCorner, strings.Repeat(args.BoxCharacters.HorizontalEdge, width-2), args.BoxCharacters.TopRightCorner,
		)
		bottomBorder := fmt.Sprintf(
			"%s%s%s",
			args.BoxCharacters.BottomLeftCorner, strings.Repeat(args.BoxCharacters.HorizontalEdge, width-2), args.BoxCharacters.BottomRightCorner,
		)
		infoLine = fmt.Sprintf(
			"%s\n%s %s %s\n%s\n",
			topBorder, args.BoxCharacters.VerticalEdge, infoLine, args.BoxCharacters.VerticalEdge, bottomBorder,
		)
	} else {
		infoLine = fmt.Sprintf("%s\n", infoLine)
	}
	fmt.Printf("%s%s", pokedex.Decompress(d), infoLine)
}
