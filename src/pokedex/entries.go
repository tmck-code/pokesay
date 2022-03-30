package pokedex

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type PokemonEntry struct {
	Name  string
	Index int
}

type Node struct {
	Children map[string]*Node
	Data     []*PokemonEntry
}

type PokemonTrie struct {
	Root *Node
	Len  int
	Keys [][]string
}

func NewPokemonEntry(idx int, name string) *PokemonEntry {
	return &PokemonEntry{
		Index: idx,
		Name:  name,
	}
}

func NewNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
	}
}

func NewTrie() *PokemonTrie {
	return &PokemonTrie{
		Len:  0,
		Root: NewNode(),
		Keys: make([][]string, 0),
	}
}

func EntryFpath(idx int) string {
	return fmt.Sprintf("build/%d.cow", idx)
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func (t *PokemonTrie) Insert(s []string, data *PokemonEntry) {
	current := t.Root
	found := false
	for _, k := range t.Keys {
		if Equal(k, s) {
			found = true
			break
		}
	}
	if !found {
		t.Keys = append(t.Keys, s)
	}
	for _, char := range s {
		if _, ok := current.Children[char]; ok {
			current = current.Children[char]
		} else {
			current.Children[char] = NewNode()
			current = current.Children[char]
			current.Data = make([]*PokemonEntry, 0)
		}
	}
	current.Data = append(current.Data, data)
	t.Len += 1
}

func (t PokemonTrie) GetCategoryPaths(s string) [][]string {
	matches := [][]string{}
	for _, k := range t.Keys {
		for i, el := range k {
			if el == s {
				if i == 0 {
					matches = append(matches, []string{s})
				} else {
					matches = append(matches, k)
				}
			}
		}
	}
	return matches
}

func (t PokemonTrie) GetCategory(s []string) ([]*PokemonEntry, bool) {
	current := t.Root
	matches := make([]*PokemonEntry, 0)
	for _, char := range s {
		if _, ok := current.Children[char]; ok {
			current = current.Children[char]
			for _, p := range current.Children {
				matches = append(matches, p.Data...)
			}
		} else {
			return nil, false
		}
	}
	return matches, true
}

func TokenizeName(name string) []string {
	return strings.Split(name, "-")
}

func matchPokemon(p []*PokemonEntry, nameToken string) ([]*PokemonEntry, bool) {
	matches := make([]*PokemonEntry, 0)
	for _, pk := range p {
		fmt.Println(nameToken, pk)
		for _, tk := range TokenizeName(pk.Name) {
			if tk == nameToken {
				matches = append(matches, pk)
			}
		}
	}
	if len(matches) > 0 {
		return matches, true
	} else {
		return nil, false
	}
}

func (t PokemonTrie) MatchNameToken(s string) []*PokemonEntry {
	current := t.Root
	matches := make([]*PokemonEntry, 0)
	for _, key := range t.Keys {
		fmt.Println("checking", key)
		for _, char := range key {
			if _, ok := current.Children[char]; ok {
				current = current.Children[char]
				if m, ok := matchPokemon(current.Data, s); ok {
					matches = append(matches, m...)
				}
			}
		}
	}
	return matches
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

func WriteStructToFile(data interface{}, fpath string) {
	ostream, err := os.Create(fpath)
	check(err)

	writer := bufio.NewWriter(ostream)
	enc := gob.NewEncoder(writer)
	enc.Encode(data)
	writer.Flush()
	ostream.Close()
}

func ReadStructFromBytes(data []byte) PokemonTrie {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	categories := &PokemonTrie{}

	err := dec.Decode(&categories)
	check(err)

	return *categories
}

func WriteCompressedToFile(data []byte, fpath string) {
	ostream, err := os.Create(fpath)
	check(err)

	writer := bufio.NewWriter(ostream)
	writer.WriteString(string(Compress(data)))
	writer.Flush()
	ostream.Close()
}
