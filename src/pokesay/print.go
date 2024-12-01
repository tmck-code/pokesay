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
	// a sentinel value that won't collide with a regular string (hopefully)
	// TODO: find a better way to do this
	Sentinel string = "\x00"
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
	var b strings.Builder
	c := make(chan string)

	// pass the buffer to the functions
	go drawPokemon(args, choice, names, categories, cows, c)
	go drawSpeechBubble(args.BoxChars, bufio.NewScanner(os.Stdin), args, c)

	sentinels := 0
	for s := range c {
		if s == Sentinel {
			sentinels++
			if sentinels == 2 {
				close(c)
				break
			}
		}
		b.WriteString(s)
	}
	fmt.Print(b.String())
}

// Prints text from STDIN, surrounded by a speech bubble.
func drawSpeechBubble(boxChars *BoxChars, scanner *bufio.Scanner, args Args, c chan string) {
	if args.DrawBubble {
		// c <- boxChars.TopLeftCorner
		// c <- strings.Repeat(boxChars.HorizontalEdge, args.Width+2)
		// c <- boxChars.TopRightCorner + "\n"
		c <- boxChars.TopLeftCorner + strings.Repeat(boxChars.HorizontalEdge, args.Width+2) + boxChars.TopRightCorner + "\n"
	}

	for scanner.Scan() {
		line := scanner.Text()

		if !args.NoTabSpaces {
			line = strings.Replace(line, "\t", args.TabSpaces, -1)
		}
		if args.NoWrap {
			drawSpeechBubbleLine(boxChars, line, args, c)
		} else {
			drawWrappedText(boxChars, line, args, c)
		}
	}

	bottomBorder := strings.Repeat(boxChars.HorizontalEdge, 6) +
		boxChars.BalloonTether +
		strings.Repeat(boxChars.HorizontalEdge, args.Width+2-7)

	if args.DrawBubble {
		// c <- boxChars.BottomLeftCorner
		// c <- bottomBorder
		// c <- boxChars.BottomRightCorner + "\n"
		c <- boxChars.BottomLeftCorner + bottomBorder + boxChars.BottomRightCorner + "\n"
	} else {
		// c <- " "
		// c <- bottomBorder
		// c <- " \n"
		c <- bottomBorder + "\n"
	}
	for i := 0; i < 4; i++ {
		// c <- strings.Repeat(" ", i+8)
		// c <- boxChars.BalloonString + "\n"
		c <- strings.Repeat(" ", i+8) + boxChars.BalloonString + "\n"
	}
	c <- Sentinel
}

// Prints a single speech bubble line
func drawSpeechBubbleLine(boxChars *BoxChars, line string, args Args, c chan string) {
	if !args.DrawBubble {
		c <- line + "\n"
		return
	}

	lineLen := UnicodeStringLength(line)
	if lineLen <= args.Width {
		// c <- boxChars.VerticalEdge
		// c <- line
		// c <- resetColourANSI
		// c <- strings.Repeat(" ", args.Width-lineLen)
		// c <- boxChars.VerticalEdge + "\n"
		c <- boxChars.VerticalEdge + line + resetColourANSI + strings.Repeat(" ", args.Width-lineLen) + boxChars.VerticalEdge + "\n"
	} else if lineLen > args.Width {
		// c <- boxChars.VerticalEdge
		// c <- line
		// c <- resetColourANSI + "\n"
		c <- boxChars.VerticalEdge + line + resetColourANSI + "\n"
	}
}

// Prints line of text across multiple lines, wrapping it so that it doesn't exceed the desired width.
func drawWrappedText(boxChars *BoxChars, line string, args Args, c chan string) {
	for _, wline := range strings.Split(wordwrap.WrapString(strings.Replace(line, "\t", args.TabSpaces, -1), uint(args.Width)), "\n") {
		drawSpeechBubbleLine(boxChars, wline, args, c)
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
func drawPokemon(args Args, index int, names []string, categoryKeys []string, GOBCowData embed.FS, c chan string) {
	d, _ := GOBCowData.ReadFile(pokedex.EntryFpath("build/assets/cows", index))
	c <- string(pokedex.Decompress(d))

	width := nameLength(names)
	namesFmt := make([]string, 0)
	for _, name := range names {
		namesFmt = append(namesFmt, textStyleBold.Sprint(name))
	}
	// count name separators
	width += (len(names) - 1) * 3
	width += 2     // for the arrow
	width += 2 + 2 // for the end box characters

	if args.DrawInfoBorder {
		// topBorder := fmt.Sprintf(
		// 	"%s%s%s",
		// 	args.BoxChars.TopLeftCorner, strings.Repeat(args.BoxChars.HorizontalEdge, width-2), args.BoxChars.TopRightCorner,
		// )
		c <- args.BoxChars.TopLeftCorner
		c <- strings.Repeat(args.BoxChars.HorizontalEdge, width-2)
		c <- args.BoxChars.TopRightCorner + "\n"
		// c <- args.BoxChars.TopLeftCorner + strings.Repeat(args.BoxChars.HorizontalEdge, width-2) + args.BoxChars.TopRightCorner + "\n"
	}

	infoSep := " " + args.BoxChars.Separator + " "
	if args.DrawInfoBorder {
		c <- args.BoxChars.VerticalEdge + " "
	}
	if args.NoCategoryInfo {
		// infoLine = fmt.Sprintf(
		// 	"%s %s",
		// 	args.BoxChars.RightArrow, strings.Join(namesFmt, fmt.Sprintf(" %s ", args.BoxChars.Separator)),
		// )
		// c <- args.BoxChars.RightArrow
		// c <- strings.Join(namesFmt, infoSep) + "\n"
		c <- args.BoxChars.RightArrow + strings.Join(namesFmt, infoSep) + "\n"
	} else {
		// infoLine = fmt.Sprintf(
		// 	"%s %s %s %s",
		// 	args.BoxChars.RightArrow,
		// 	strings.Join(namesFmt, fmt.Sprintf(" %s ", args.BoxChars.Separator)),
		// 	args.BoxChars.Separator,
		// 	textStyleItalic.Sprint(strings.Join(categoryKeys, args.BoxChars.CategorySeparator)),
		// )
		// c <- args.BoxChars.RightArrow
		// c <- strings.Join(namesFmt, infoSep)
		// c <- infoSep
		// c <- textStyleItalic.Sprint(strings.Join(categoryKeys, args.BoxChars.CategorySeparator)) + "\n"
		c <- args.BoxChars.RightArrow + strings.Join(namesFmt, infoSep) + infoSep + textStyleItalic.Sprint(strings.Join(categoryKeys, args.BoxChars.CategorySeparator)) + "\n"

		for _, category := range categoryKeys {
			width += len(category)
		}
		width += len(categoryKeys) - 1 + 1 + 2 // lol why did I do this
	}

	if args.DrawInfoBorder {
		// c <- " " + args.BoxChars.VerticalEdge + "\n"
		// c <- args.BoxChars.BottomLeftCorner
		// c <- strings.Repeat(args.BoxChars.HorizontalEdge, width-2)
		// c <- args.BoxChars.BottomRightCorner + "\n"
		c <- " " + args.BoxChars.VerticalEdge + "\n" + args.BoxChars.BottomLeftCorner + strings.Repeat(args.BoxChars.HorizontalEdge, width-2) + args.BoxChars.BottomRightCorner + "\n"
	}

	c <- Sentinel
	// close(c)
}
