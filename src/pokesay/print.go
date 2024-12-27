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

type BoxChars struct {
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
	BoxChars       *BoxChars
	DrawInfoBorder bool
	Help           bool
	Verbose        bool
}

var (
	textStyleItalic *color.Color = color.New(color.Italic)
	textStyleBold   *color.Color = color.New(color.Bold)
	resetColourANSI string       = "\033[0m"
	AsciiBoxChars   *BoxChars    = &BoxChars{
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
	UnicodeBoxChars *BoxChars = &BoxChars{
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

func DetermineBoxChars(unicodeBox bool) *BoxChars {
	if unicodeBox {
		return UnicodeBoxChars
	} else {
		return AsciiBoxChars
	}
}

// The main print function! This uses a chosen pokemon's index, names and categories, and an
// embedded filesystem of cowfile data
// 1. The text received from STDIN is printed inside a speech bubble
// 2. The cowfile data is retrieved using the matching index, decompressed (un-gzipped),
// 3. The pokemon is printed along with the name & category information
func Print(args Args, choice int, names []string, categories []string, cows embed.FS) {
	printSpeechBubble(args.BoxChars, bufio.NewScanner(os.Stdin), args)
	printPokemon(args, choice, names, categories, cows)
}

// Prints text from STDIN, surrounded by a speech bubble.
func printSpeechBubble(boxChars *BoxChars, scanner *bufio.Scanner, args Args) {
	if args.DrawBubble {
		fmt.Printf(
			"%s%s%s\n",
			boxChars.TopLeftCorner,
			strings.Repeat(boxChars.HorizontalEdge, args.Width+2),
			boxChars.TopRightCorner,
		)
	}

	for scanner.Scan() {
		line := scanner.Text()

		if !args.NoTabSpaces {
			line = strings.Replace(line, "\t", args.TabSpaces, -1)
		}
		if args.NoWrap {
			printSpeechBubbleLine(boxChars, line, args)
		} else {
			printWrappedText(boxChars, line, args)
		}
	}

	bottomBorder := strings.Repeat(boxChars.HorizontalEdge, 6) +
		boxChars.BalloonTether +
		strings.Repeat(boxChars.HorizontalEdge, args.Width+2-7)

	if args.DrawBubble {
		fmt.Printf("%s%s%s\n", boxChars.BottomLeftCorner, bottomBorder, boxChars.BottomRightCorner)
	} else {
		fmt.Printf(" %s \n", bottomBorder)
	}
	for i := 0; i < 4; i++ {
		fmt.Printf("%s%s\n", strings.Repeat(" ", i+8), boxChars.BalloonString)
	}
}

// Prints a single speech bubble line
func printSpeechBubbleLine(boxChars *BoxChars, line string, args Args) {
	if !args.DrawBubble {
		fmt.Println(line)
		return
	}

	lineLen := UnicodeStringLength(line)
	if lineLen <= args.Width {
		// print the line with padding, the most common case
		fmt.Printf(
			"%s %s%s%s %s\n",
			boxChars.VerticalEdge, // left-hand side of the bubble
			line, resetColourANSI, // the text
			strings.Repeat(" ", args.Width-lineLen), // padding
			boxChars.VerticalEdge,                   // right-hand side of the bubble
		)
	} else if lineLen > args.Width {
		// print the line without padding or right-hand side of the bubble if the line is too long
		fmt.Printf(
			"%s %s%s\n",
			boxChars.VerticalEdge, // left-hand side of the bubble
			line, resetColourANSI, // the text
		)
	}
}

// Prints line of text across multiple lines, wrapping it so that it doesn't exceed the desired width.
func printWrappedText(boxChars *BoxChars, line string, args Args) {
	for _, wline := range strings.Split(wordwrap.WrapString(strings.Replace(line, "\t", args.TabSpaces, -1), uint(args.Width)), "\n") {
		printSpeechBubbleLine(boxChars, wline, args)
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
			// e.g. "\x1b[38;5;196m"
			//       ^^^ start    ^ end
			if s[i:i+2] == "\x1b[" {
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

func ReverseUnicodeString(s string) string {
	runes := []rune(s)
	reversed := make([]rune, len(runes))

	for i, r := range runes {
		reversed[len(runes)-1-i] = r
	}
	return string(reversed)
}

type ANSILineToken struct {
	FGColour string
	BGColour string
	Text     string
}

func TokeniseANSIString(msg string) [][]ANSILineToken {
	var isColour bool
	var colour string
	var fg string
	var bg string

	var tokens []ANSILineToken
	var lines [][]ANSILineToken
	for _, line := range strings.Split(msg, "\n") {
		var text string

		for _, ch := range line {
			if ch == '\033' {
				if text != "" {
					tokens = append(tokens, ANSILineToken{fg, bg, text})
					colour = ""
					text = ""
				}
				isColour = true
				colour = string(ch)
			} else if isColour {
				colour += string(ch)
				if ch == 'm' {
					isColour = false
					if strings.Contains(colour, "38;5;") {
						fg = colour
					} else if strings.Contains(colour, "48;5;") {
						bg = colour
					} else {
						fg = colour
						bg = ""
					}
				}
			} else {
				text += string(ch)
			}
		}
		if text != "" {
			tokens = append(tokens, ANSILineToken{fg, bg, text})
		}
		if (colour != "") && len(tokens) > 0 {
			tokens = append(tokens, ANSILineToken{"\033[0m", "", ""})
		}
		lines = append(lines, tokens)
		tokens = nil
	}
	return lines
}

func ReverseANSIString(line string) string {
	lines := TokeniseANSIString(line)
	reversed := ""

	maxWidth := 0
	widths := make([]int, len(lines))
	for idx, l := range strings.Split(line, "\n") {
		ln := UnicodeStringLength(l)
		if ln > maxWidth {
			maxWidth = ln
		}
		widths[idx] = ln
	}

	for idx, tokens := range lines {
		// ensure vertical alignment
		reversed += strings.Repeat(" ", maxWidth-widths[idx])
		for i := len(tokens) - 1; i >= 0; i-- {
			reversed += tokens[i].FGColour + tokens[i].BGColour + ReverseUnicodeString(tokens[i].Text)
		}
		if idx < len(lines)-1 {
			reversed += "\n"
		}
	}
	if reversed[len(reversed)-2:len(reversed)-1] != "\x1b[0m" {
		reversed += "\033[0m"
	}
	if line[len(line)-1] == '\n' {
		reversed += "\n"
	}
	return reversed
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
			args.BoxChars.RightArrow, strings.Join(namesFmt, fmt.Sprintf(" %s ", args.BoxChars.Separator)),
		)
	} else {
		infoLine = fmt.Sprintf(
			"%s %s %s %s",
			args.BoxChars.RightArrow,
			strings.Join(namesFmt, fmt.Sprintf(" %s ", args.BoxChars.Separator)),
			args.BoxChars.Separator,
			textStyleItalic.Sprint(strings.Join(categoryKeys, args.BoxChars.CategorySeparator)),
		)
		for _, category := range categoryKeys {
			width += len(category)
		}
		width += len(categoryKeys) - 1 + 1 + 2 // lol why did I do this
	}

	if args.DrawInfoBorder {
		topBorder := fmt.Sprintf(
			"%s%s%s",
			args.BoxChars.TopLeftCorner, strings.Repeat(args.BoxChars.HorizontalEdge, width-2), args.BoxChars.TopRightCorner,
		)
		bottomBorder := fmt.Sprintf(
			"%s%s%s",
			args.BoxChars.BottomLeftCorner, strings.Repeat(args.BoxChars.HorizontalEdge, width-2), args.BoxChars.BottomRightCorner,
		)
		infoLine = fmt.Sprintf(
			"%s\n%s %s %s\n%s\n",
			topBorder, args.BoxChars.VerticalEdge, infoLine, args.BoxChars.VerticalEdge, bottomBorder,
		)
	} else {
		infoLine = fmt.Sprintf("%s\n", infoLine)
	}
	fmt.Printf("%s%s", pokedex.Decompress(d), infoLine)
}
