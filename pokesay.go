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
)

func main() {
	width := 40
	if len(os.Args) > 1 {
		width, _ = strconv.Atoi(os.Args[1])
	}

	scanner := bufio.NewScanner(os.Stdin)
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
	cats := fpathParts[1 : len(fpathParts)-1]

	fmt.Println("choice:", fchoice, "/", "categories:", cats)
}
