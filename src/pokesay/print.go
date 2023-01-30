package pokesay

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-wordwrap"
	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	textStyleItalic *color.Color = color.New(color.Italic)
	textStyleBold   *color.Color = color.New(color.Bold)
)

type Args struct {
	Width          int
	NoWrap         bool
	TabSpaces      string
	NoTabSpaces    bool
	NoCategoryInfo bool
	ListCategories bool
	ListNames      bool
	Category       string
	NameToken      string
	JapaneseName   bool
}

// The main print function! This uses a chosen pokemon's index, names and categories, and an
// embedded filesystem of cowfile data
// 1. The text received from STDIN is printed inside a speech bubble
// 2. The cowfile data is retrieved using the matching index, decompressed (un-gzipped),
// 3. The pokemon is printed along with the name & category information
func Print(args Args, choice int, names []string, categories []string, cows embed.FS) {
	printSpeechBubble(bufio.NewScanner(os.Stdin), args.Width, args.NoTabSpaces, args.TabSpaces, args.NoWrap)
	printPokemon(args, choice, names, categories, cows)
}

// Prints text from STDIN, surrounded by a speech bubble.
func printSpeechBubble(scanner *bufio.Scanner, width int, noTabSpaces bool, tabSpaces string, noWrap bool) {
	border := strings.Repeat("-", width+2)
	fmt.Println("/" + border + "\\")

	for scanner.Scan() {
		line := scanner.Text()

		if !noTabSpaces {
			line = strings.Replace(line, "\t", tabSpaces, -1)
		}
		if noWrap {
			printSpeechBubbleLine(line, width)
		} else {
			printWrappedText(line, width, tabSpaces)
		}
	}
	fmt.Println("\\" + border + "/")
	for i := 0; i < 4; i++ {
		fmt.Println(strings.Repeat(" ", i+8), "\\")
	}
}

// Prints a single speech bubble line
func printSpeechBubbleLine(line string, width int) {
	if len(line) > width {
		fmt.Printf("| %s\n", line)
	} else if len(line) == width {
		fmt.Printf("| %s |\n", line)
	} else {
		fmt.Printf("| %s%s |\n", line, strings.Repeat(" ", width-len(line)))
	}
}

// Prints line of text across multiple lines, wrapping it so that it doesn't exceed the desired width.
func printWrappedText(line string, width int, tabSpaces string) {
	for _, wline := range strings.Split(wordwrap.WrapString(strings.Replace(line, "\t", tabSpaces, -1), uint(width)), "\n") {
		printSpeechBubbleLine(wline, width)
	}
}

// Prints a pokemon with its name & category information.
func printPokemon(args Args, index int, names []string, categoryKeys []string, GOBCowData embed.FS) {
	d, _ := GOBCowData.ReadFile(pokedex.EntryFpath("build/assets/cows", index))
	delimiter := "|"

	namesFmt := make([]string, 0)
	for _, name := range names {
		namesFmt = append(namesFmt, textStyleBold.Sprint(name))
	}

	if args.NoCategoryInfo {
		fmt.Printf(
			"%s> %s\n",
			pokedex.Decompress(d),
			strings.Join(namesFmt, fmt.Sprintf(" %s ", delimiter)),
		)
	} else {
		fmt.Printf(
			"%s> %s %s %s\n",
			pokedex.Decompress(d),
			strings.Join(namesFmt, fmt.Sprintf(" %s ", delimiter)),
			delimiter,
			textStyleItalic.Sprint(strings.Join(categoryKeys, "/")),
		)
	}
}
