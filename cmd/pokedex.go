package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PokemonEntry struct {
	Name       string
	Data       []byte
	Categories []string
}

func findFiles(dirpath string, ext string, skip []string) map[string][]PokemonEntry {
	categories := make(map[string][]PokemonEntry)
	err := filepath.Walk(dirpath, func(fpath string, f os.FileInfo, err error) error {
		for _, s := range skip {
			if strings.Contains(fpath, s) {
				return err
			}
		}
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			data, err := os.ReadFile(fpath)
			check(err)
			pokemonCategories := createCategories(fpath)
			for _, c := range pokemonCategories[3 : len(pokemonCategories)-1] {

				if val, ok := categories[c]; ok {
					categories[c] = append(val, PokemonEntry{Name: fpath, Data: data})
				} else {
					categories[c] = []PokemonEntry{PokemonEntry{Name: fpath, Data: data}}
				}
			}
		}
		return err
	})
	check(err)

	return categories
}

func createCategories(fpath string) []string {
	return strings.Split(fpath, "/")
}

func createPokemonEntry(pokemon *PokemonEntry, fpath string) *PokemonEntry {
	return &PokemonEntry{
		Name:       pokemon.Name,
		Data:       pokemon.Data,
		Categories: createCategories(pokemon.Name),
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func writeToFile(categories map[string][]PokemonEntry, fpath string) {
	ostream, err := os.Create(fpath)
	check(err)

	writer := bufio.NewWriter(ostream)
	enc := gob.NewEncoder(writer)
	enc.Encode(categories)
	writer.Flush()
	ostream.Close()

	fmt.Println("->", fpath)
}

func readFromFile(fpath string) map[string][]PokemonEntry {
	istream, err := os.Open(fpath)
	check(err)

	reader := bufio.NewReader(istream)
	dec := gob.NewDecoder(reader)

	categories := make(map[string][]PokemonEntry)

	err = dec.Decode(&categories)
	check(err)
	istream.Close()

	return categories
}

type CowBuildArgs struct {
	FromDir  string
	ToFpath  string
	SkipDirs []string
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	toFpath := flag.String("to", ".", "to fpath")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, ToFpath: *toFpath}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	return args
}

func randInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

func main() {
	args := parseArgs()
	fmt.Println("starting at", args.FromDir)

	// categories := findFiles(args.FromDir, ".cow", args.SkipDirs)

	// writeToFile(categories, args.ToFpath)

	categories := readFromFile(args.ToFpath)

	total := 0
	for _, _ = range categories {
		total += 1
	}
	randomCategory := randInt(total)

	idx := 0
	choice := ""
	for cat, _ := range categories {
		idx += 1
		if idx == randomCategory {
			choice = cat
		}
	}
	fmt.Println("category choice", choice)

	pokemon := categories[choice][randInt(len(categories[choice]))]

	fmt.Printf("%s\n", pokemon.Data)
}
