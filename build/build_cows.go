package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func run() []string {
	fromDir := flag.String("from", ".", "from dir")
	toDir := flag.String("to", ".", "to dir")

	flag.Parse()
	args := []string{*fromDir, *toDir}
	fmt.Println(args)
	return args
}

func findFiles(dirpath string, ext string, ch chan<- string) {
	fmt.Println("starting at", dirpath)
	count := 0
	err := filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			ch <- path
			count += 1
		}
		return err
	})
	if err != nil {
		fmt.Println("Fatal!")
		log.Fatal(err)
	}
	ch <- "END"
	close(ch)
	fmt.Println("exiting", count, len(ch))
}

func convertPngToCow(source_fpath string, source_dirpath string, destination_dirpath string, wg *sync.WaitGroup) {
	defer wg.Done()
	to_dir := filepath.Join(destination_dirpath, filepath.Dir(strings.ReplaceAll(source_fpath, source_dirpath, "")))
	to_fpath := filepath.Join(to_dir, filepath.Base(source_fpath)+".cow")
	os.MkdirAll(to_dir, 0755)

	cmd := exec.Command("/usr/local/bin/img2xterm", string(source_fpath))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	}
	err = os.WriteFile(to_fpath, out.Bytes(), 0644)
	fmt.Println(source_fpath, "->", to_fpath)
}

func main() {
	args := run()

	ch := make(chan string)
	go findFiles(args[0], ".png", ch)

	fmt.Println("channel has", len(ch), "items")
	var wg sync.WaitGroup
	opened := true
	var file string
	count := 0
	for opened {
		file = <-ch
		if file == "END" {
			fmt.Println("reached END, exiting")
			break
		}
		wg.Add(1)
		go convertPngToCow(file, args[0], args[1], &wg)
		count += 1
	}
	wg.Wait()
	fmt.Println("final", count)
}
