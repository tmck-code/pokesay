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
	"github.com/tmck-code/pokesay-go/internal/timer"
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

func pickRandomPokemon() []byte {
	idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(_bindata))
	current := 0
	var choice []byte
	for data, _ := range _bindata {
		if current == idx {
			choice, _ = Asset(data)
			break
		}
		current += 1
	}
	return choice
}

func pickRandomPokemon2() []byte {
	idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(_bindatalist))
	choice, _ := Asset(_bindatalist[idx])
	return choice
}

func printPokemon(t *timer.Timer) {
	// data := pickRandomPokemon()
	// t.Mark("printPokemon.choose")
	data := pickRandomPokemon2()
	t.Mark("printPokemon.choose2")

	binary.Write(os.Stdout, binary.LittleEndian, data)
	t.Mark("printPokemon.print")
}

func main() {
	width := 40
	if len(os.Args) > 1 {
		width, _ = strconv.Atoi(os.Args[1])
	}
	t := timer.NewTimer()

	printSpeechBubble(bufio.NewScanner(os.Stdin), width)
	t.Mark("printSpeechBubble")

	printPokemon(t)

	t.StopTimer()
	t.PrintJson()
}
