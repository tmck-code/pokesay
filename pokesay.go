package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"strings"
	"github.com/mitchellh/go-wordwrap"
	"strconv"
)

func main() {
	width := 40
	if len(os.Args) > 1 {
		width, _ = strconv.Atoi(os.Args[1])
	}

	lines := make([]string, 0, width)
	lines = append(lines, strings.Repeat("-", width))
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, strings.Replace(scanner.Text(), "\t", "    ", -1))
	}
	lines = append(lines, strings.Repeat("-", width))

	for _, l := range(strings.Split(wordwrap.WrapString(strings.Join(lines, "\n"), uint(width)), "\n")) {
		if (len(l) > width+2) {
			fmt.Println("| " + l)
		} else {
			fmt.Println("| " + l + strings.Repeat(" ", (width+2)-len(l)) + "|")
		}
	}
	for i := 0; i<4; i++ {
		fmt.Println(strings.Repeat(" ", i+8) + "\\")
	}

	count := 0
	choice := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(_bindata))
	fpath := ""

	for path, _ := range _bindata {
		if count == choice {
			data, err := Asset(path)
			fpath = path
			if err != nil {
				log.Println(err)
				break
			}
			binary.Write(os.Stdout, binary.LittleEndian, data)
			if err != nil {
				log.Println("Asset %s can't read by error: %v", path, err)
				break
			}
		}
		count += 1
	}
	fpathParts := strings.Split(fpath, "/")
	fchoice := strings.Split(fpathParts[len(fpathParts)-1], ".")[0]
	cats := fpathParts[1:len(fpathParts)-1]

	fmt.Println("choice:", fchoice, "/", "categories:", cats)
}
