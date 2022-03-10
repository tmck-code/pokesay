package main

import (
	"bufio"
	"encoding/binary"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"google.golang.org/protobuf/proto"
)

func printSpeechBubbleLine(line string, width int) {
	if len(line) > width {
		fmt.Println("|", line)
	} else {
		fmt.Println("|", line, strings.Repeat(" ", width-len(line)), "|")
	}
}

func printWrappedText(line string, width int, tabSpaces string) {
	for _, wline := range strings.Split(wordwrap.WrapString(strings.Replace(line, "\t", tabSpaces, -1), uint(width)), "\n") {
		printSpeechBubbleLine(wline, width)
	}
}

func printSpeechBubble(scanner *bufio.Scanner, args Args) {
	border := strings.Repeat("-", args.Width+2)
	fmt.Println("/" + border + "\\")

	for scanner.Scan() {
		line := scanner.Text()

		if !args.NoTabSpaces {
			line = strings.Replace(line, "\t", args.TabSpaces, -1)
		}
		if args.NoWrap {
			printSpeechBubbleLine(line, args.Width)
		} else {
			printWrappedText(line, args.Width, args.TabSpaces)
		}
	}
	fmt.Println("\\" + border + "/")
	for i := 0; i < 4; i++ {
		fmt.Println(strings.Repeat(" ", i+8), "\\")
	}
}

func randomInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

func printPokemon(list *PokemonSerialList) {
	choice := list.Pokemon[randomInt(len(list.Pokemon))]
	binary.Write(os.Stdout, binary.LittleEndian, choice.Data)
	fmt.Printf("choice: %s", choice.Name)
}

type Args struct {
	Width       int
	NoWrap      bool
	TabSpaces   string
	NoTabSpaces bool
}

func parseFlags() Args {
	width := flag.Int("width", 80, "the max speech bubble width")
	noWrap := flag.Bool("nowrap", false, "disable text wrapping (fastest)")
	tabWidth := flag.Int("tabwidth", 4, "replace any tab characters with N spaces")
	noTabSpaces := flag.Bool("notabspaces", false, "do not replace tab characters (fastest)")
	fastest := flag.Bool("fastest", false, "run with the fastest possible configuration (-nowrap -notabspaces)")

	flag.Parse()
	var args Args

	if *fastest {
		args = Args{
			Width:       *width,
			NoWrap:      true,
			TabSpaces:   "    ",
			NoTabSpaces: true,
		}
	} else {
		args = Args{
			Width:       *width,
			NoWrap:      *noWrap,
			TabSpaces:   strings.Repeat(" ", *tabWidth),
			NoTabSpaces: *noTabSpaces,
		}
	}
	return args
}

func loadPokemon(fpath string) *PokemonSerialList {
	pokemon := &PokemonSerialList{}
	data, _ := os.ReadFile(fpath)
	err := proto.Unmarshal(data, pokemon)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	return pokemon

}

func main() {
	args := parseArgs()
	pokemon := loadPokemon("data.txt")

	printSpeechBubble(bufio.NewScanner(os.Stdin), args.Width)
	printPokemon(pokemon)
}
