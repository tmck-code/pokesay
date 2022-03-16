package pokedex

import (
	"os"
	"bufio"
	"bytes"
	"encoding/gob"
	"log"
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
	NCategories int
}

func WriteToFile(categories PokemonEntryMap, fpath string) {
	ostream, err := os.Create(fpath)
	check(err)

	writer := bufio.NewWriter(ostream)
	enc := gob.NewEncoder(writer)
	enc.Encode(categories)
	writer.Flush()
	ostream.Close()
}

func ReadFromBytes(data []byte) PokemonEntryMap {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	categories := &PokemonEntryMap{}

	err := dec.Decode(&categories)
	check(err)

	return *categories
}

func ReadFromFile(fpath string) PokemonEntryMap {
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

