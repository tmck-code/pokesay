package pokesay

import (
	"bufio"
	"embed"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-wordwrap"
	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	textStyleItalic *color.Color = color.New(color.Italic)
	textStyleBold   *color.Color = color.New(color.Bold)
)

func printSpeechBubbleLine(line string, width int) {
	if len(line) > width {
		fmt.Printf("| %s\n", line)
	} else if len(line) == width {
		fmt.Printf("| %s |\n", line)
	} else {
		fmt.Printf("| %s%s |\n", line, strings.Repeat(" ", width-len(line)))
	}
}

func printWrappedText(line string, width int, tabSpaces string) {
	for _, wline := range strings.Split(wordwrap.WrapString(strings.Replace(line, "\t", tabSpaces, -1), uint(width)), "\n") {
		printSpeechBubbleLine(wline, width)
	}
}

func PrintSpeechBubble(scanner *bufio.Scanner, width int, noTabSpaces bool, tabSpaces string, noWrap bool) {
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

func PrintPokemon(index int, names []string, categoryKeys []string, GOBCowData embed.FS) {
	d, _ := GOBCowData.ReadFile(pokedex.EntryFpath("build/assets/cows", index))

	fmt.Printf(
		"%s> %s | %s\n",
		pokedex.Decompress(d),
		textStyleBold.Sprint(strings.Join(names, " / ")),
		textStyleItalic.Sprint(strings.Join(categoryKeys, "/")),
	)
}
