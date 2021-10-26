package main

import (
	// "io/ioutil"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	walk(os.Args[1])
}

func walk(dirPath string) {
	count := 0
	choice := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(_bindata))
	fmt.Println(choice)

	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// fmt.Println(path, info.Size())
			if count == choice {
				data, err := Asset(path)
				if err != nil {
					return err
				}
				binary.Write(os.Stdout, binary.LittleEndian, data)
			}
			count += 1
			return nil
		})
	fmt.Println(count)
	if err != nil {
		log.Println(err)
	}
}
