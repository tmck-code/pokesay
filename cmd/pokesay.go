package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
)

func printSpeechBubble(scanner *bufio.Scanner, width int) {
	border := strings.Repeat("-", width+2)
	fmt.Println("/" + border + "\\")
	for scanner.Scan() {
		for _, wline := range strings.Split(wordwrap.WrapString(strings.Replace(scanner.Text(), "\t", "    ", -1), uint(width)), "\n") {
			if len(wline) > width {
				fmt.Println("| ", wline, len(wline))
			} else {
				fmt.Println("|", wline, strings.Repeat(" ", width-len(wline)), "|")
			}
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

func printPokemon() {
	choice :=  PokemonList[randomInt(len(PokemonList))]
	binary.Write(os.Stdout, binary.LittleEndian, choice.Data)
	fmt.Printf("choice: %s / categories: %s\n", choice.Name.Name, choice.Name.Categories)
}

type Args struct {
	Width int
}

func parseArgs() Args {
	if len(os.Args) <= 1 {
		return Args{Width: 40}
	}
	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return Args{Width: 40}
	}
	return Args{Width: width}
}

func main() {
	args := parseArgs()
	printSpeechBubble(bufio.NewScanner(os.Stdin), args.Width)
	printPokemon()
}
