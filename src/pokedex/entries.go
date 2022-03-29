package pokedex

import (
	"os"
	"fmt"
	"bufio"
	"bytes"
	"encoding/gob"
	"compress/gzip"
	"log"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type PokemonEntry struct {
	Name       	string
	NameTokens	[]string
	Index       string
}

type PokemonEntryMap struct {
	Categories map[string][]*PokemonEntry
}

func NewPokemonEntry(idx int, name string, nameTokens []string) *PokemonEntry {
	return &PokemonEntry{
		Index: fmt.Sprintf("build/%d.cow", idx),
		Name: name,
		NameTokens: nameTokens,
	}
}

func Compress(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err := gz.Write(data)
	check(err)

	err = gz.Close()
	check(err)

	return b.Bytes()
}

func Decompress(data []byte) []byte {
	buf := bytes.NewBuffer(data)

	reader, err := gzip.NewReader(buf)
	check(err)

	var resB bytes.Buffer

	_, err = resB.ReadFrom(reader)
	check(err)

	return resB.Bytes()
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

func WriteByteToFile(pokemon [][]byte, fpath string) {
	for i, entry := range pokemon {
		ostream, err := os.Create(fmt.Sprintf("build/%d.cow", i))
		check(err)

		writer := bufio.NewWriter(ostream)
		fmt.Println(i, string(entry))
		writer.WriteString(string(entry))
		writer.Flush()
		ostream.Close()
	}
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

func ReadDataFromBytes(data []byte) [][]byte {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	pokemon := make([][]byte, 0)
	idx := 0
	for scanner.Scan() {
		dat := scanner.Text()
		fmt.Println(idx, string(dat))
		pokemon = append(pokemon, []byte(dat))
		idx++
	}

	return pokemon
}

