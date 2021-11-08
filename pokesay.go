package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
)

type finishedTimer struct {
	OrderedDurations []int64
	StageDurations   map[string]int64
}

type timer struct {
	Start          time.Time
	StageNames     []string
	StageTimes     map[string]time.Time
	StageDurations map[string]time.Duration
	FinishedTimer  finishedTimer
}

func newTimer() *timer {
	now := time.Now()
	return &timer{
		Start:          now,
		StageNames:     []string{"start"},
		StageTimes:     map[string]time.Time{"start": now},
		StageDurations: map[string]time.Duration{"start": now.Sub(now)},
	}
}

func (t *timer) lastStageTime() time.Time {
	return t.StageTimes[t.StageNames[len(t.StageNames)-1]]
}

func (t *timer) mark(stage string) {
	now := time.Now()
	t.StageDurations[stage] = now.Sub(t.lastStageTime())
	t.StageTimes[stage] = now
	t.StageNames = append(t.StageNames, stage)
}

func (t *timer) stopTimer() {
	t.FinishedTimer = finishedTimer{StageDurations: map[string]int64{}}
	var total int64 = 0
	for _, stage := range t.StageNames {
		n := t.StageDurations[stage].Nanoseconds()
		t.FinishedTimer.OrderedDurations = append(t.FinishedTimer.OrderedDurations, n)
		t.FinishedTimer.StageDurations[stage] = n
		total += n
	}
}

func (t *timer) PrintJson() {
	json.NewEncoder(os.Stdout).Encode(t)
}

func printSpeechBubble(scanner *bufio.Scanner, width int, timer *timer) {
	fmt.Println("/" + strings.Repeat("-", width+2) + "\\")
	for scanner.Scan() {
		line := wordwrap.WrapString(strings.Replace(scanner.Text(), "\t", "    ", -1), uint(width))
		for _, wline := range strings.Split(wordwrap.WrapString(line, uint(width)), "\n") {
			if len(wline) > width {
				fmt.Println("| ", wline, len(wline))
			} else {
				fmt.Println("|", wline, strings.Repeat(" ", width-len(wline)), "|")
			}
		}
	}
	fmt.Println("\\" + strings.Repeat("-", width+2) + "/")
	for i := 0; i < 4; i++ {
		fmt.Println(strings.Repeat(" ", i+8), "\\")
	}
}

func printPokemon(timer *timer) {
	count := 0
	choice := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(_bindata))
	fpath := ""
	timer.mark("printPokemon.randomChoice")

	for path, _ := range _bindata {
		if count == choice {
			data, err := Asset(path)
			fpath = path
			if err != nil {
				log.Println(path, err)
				break
			}
			binary.Write(os.Stdout, binary.LittleEndian, data)
			// break
		}
		count += 1
	}
	timer.mark("printPokemon.findAndPrint")
	fpathParts := strings.Split(fpath, "/")
	fchoice := strings.Split(fpathParts[len(fpathParts)-1], ".")[0]
	cats := fpathParts[1 : len(fpathParts)-1]

	fmt.Println("choice:", fchoice, "/", "categories:", cats)
	timer.mark("printPokemon.summarise")
}

func main() {
	width := 40
	if len(os.Args) > 1 {
		width, _ = strconv.Atoi(os.Args[1])
	}
	t := newTimer()

	printSpeechBubble(bufio.NewScanner(os.Stdin), width, t)
	t.mark("printSpeechBubble")

	printPokemon(t)
	t.mark("printPokemon")

	t.stopTimer()
	t.PrintJson()
}
