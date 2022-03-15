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

	"internal/timer"
	// "timer"
	// "github.com/tmck-code/pokesay-go/internal/timer"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type PokemonEntry struct {
	Name       string
	Data       []byte
	Categories []string
}

type PokemonEntryMap struct {
	Categories map[string][]PokemonEntry
	Total int
}


func findFiles(dirpath string, ext string, skip []string) PokemonEntryMap {
	categories := &PokemonEntryMap{Categories: make(map[string][]PokemonEntry)}
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
				p := PokemonEntry{Name: fpath, Data: data}
				if val, ok := categories.Categories[c]; ok {
					val = append(val, p)
				} else {
					categories.Categories[c] = []PokemonEntry{p}
				}
				// fmt.Println("categories:", categories.Pokemon[c], c)

				categories.Categories[c] = append(categories.Categories[c], PokemonEntry{Name: fpath, Data: data})
			}
		}
		return err
	})
	check(err)

	return *categories
}

func createCategories(fpath string) []string {
	return strings.Split(fpath, "/")
}

func writeToFile(categories PokemonEntryMap, fpath string) {
	ostream, err := os.Create(fpath)
	check(err)

	writer := bufio.NewWriter(ostream)
	enc := gob.NewEncoder(writer)
	enc.Encode(categories)
	writer.Flush()
	ostream.Close()
}

func readFromFile(fpath string) PokemonEntryMap {
	istream, err := os.Open(fpath)
	check(err)

	reader := bufio.NewReader(istream)
	dec := gob.NewDecoder(reader)

	categories := &PokemonEntryMap{}

	err = dec.Decode(&categories)
	check(err)
	istream.Close()

	return *categories
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
	if n <= 0 {
		log.Fatal("n <= 0! -_-")
	}
	return rand.Intn(n-1) + 1
}

func main() {
	args := parseArgs()
	fmt.Println("starting at", args.FromDir)

	t := timer.NewTimer()

	categories := findFiles(args.FromDir, ".cow", args.SkipDirs)

	writeToFile(categories, args.ToFpath)

	categories = readFromFile(args.ToFpath)

	t.Mark("writeProtobuf")

	total := 0
	for _, _ = range categories.Categories {
		total += 1
	}
	rand.Seed(time.Now().UnixNano())
	randomCategory := randInt(total)
	t.Mark("randomCategory")

	idx := 0
	choice := ""
	for cat, _ := range categories.Categories {
		idx += 1
		if idx == randomCategory {
			choice = cat
		}
	}
	t.Mark("randomCategoryChoice")

	// fmt.Println(choice, categories[choice])
	fmt.Println("choosing pokemon", "cat:", randomCategory, "choice:", choice, " - ", len(categories.Categories[choice]))
	randPokemon := randInt(len(categories.Categories[choice]))
	t.Mark("randomPokemon")
	// fmt.Println("category choice", choice, "pokemon choice", randPokemon, "/", len(categories[choice]))
	pokemon := categories.Categories[choice][randPokemon]
	t.Mark("randomPokemonChoice")

	fmt.Printf("%s\n", pokemon.Data)
	t.Mark("print")
	t.StopTimer()
	t.PrintJson()
}
