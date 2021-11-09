package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
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

func pickRandomPokemon() string {
	count := 0
	choice := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(_bindata))
	fpath := ""
	for path, _ := range _bindata {
		if count == choice {
			fpath = path
			break
		}
		count += 1
	}
	return fpath
}

func printPokemon(t *timer.Timer) {
	fpath := pickRandomPokemon()
	t.Mark("printPokemon.randomChoice")
	data, err := Asset(fpath)
	if err != nil {
		log.Println(fpath, err)
	}
	binary.Write(os.Stdout, binary.LittleEndian, data)

	t.Mark("printPokemon.findAndPrint")
	fpathParts := strings.Split(fpath, "/")
	fchoice := strings.Split(fpathParts[len(fpathParts)-1], ".")[0]
	cats := fpathParts[1 : len(fpathParts)-1]

	fmt.Println("choice:", fchoice, "/", "categories:", cats)
	t.Mark("printPokemon.summarise")
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
	t.Mark("printPokemon")

	t.StopTimer()
	t.PrintJson()
}
