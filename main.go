package main

import (
   // "io/ioutil"
	"path/filepath"
    "fmt"
    "log"
    "os"
)

func main() {
    argsWithProg := os.Args
    argsWithoutProg := os.Args[1:]

    arg := os.Args[1]

    fmt.Println(argsWithProg)
    fmt.Println(argsWithoutProg)
    fmt.Println(arg)
    // body, err := ioutil.ReadFile(arg)
    // if err != nil {
    //     log.Fatalf("unable to read file: %v", err)
    // }
    // fmt.Println(string(body))

	err := filepath.Walk(arg,
		func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path, info.Size())
		return nil
	})
	if err != nil {
		log.Println(err)
	}

}

