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
)

func run() []string {
	fromDir := flag.String("from", ".", "from dir")
	toDir := flag.String("to", ".", "to dir")

	flag.Parse()
	args := []string{*fromDir, *toDir}
	fmt.Println(args)
	return args
}

func main() {
	args := run()
	fileList := []string{}
	fmt.Println("starting at", args[0])
	err := filepath.Walk(args[0], func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".png" {
			fileList = append(fileList, path)
		}
		return err
	})
	if err != nil {
		fmt.Println("Fatal!")
		log.Fatal(err)
	}
	fmt.Println("found", string(len(fileList)))

	for _, file := range fileList {
		to_dir := filepath.Join(args[1], filepath.Dir(strings.ReplaceAll(file, args[0], "")))
		to_fpath := filepath.Join(to_dir, filepath.Base(file) + ".cow")
		os.MkdirAll(to_dir, 0755)

		cmd := exec.Command("/usr/local/bin/img2xterm", string(file))

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
		fmt.Println(file, "->", to_fpath)
	}
}
