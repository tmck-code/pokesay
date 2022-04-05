package pokedex

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
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

type PokemonMetadata struct {
	Categories string
	Name       string
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

func MetadataFpath(idx int) string {
	return fmt.Sprintf("build/%d.metadata", idx)
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

func (t PokemonTrie) GetCategoryPaths(s string) ([][]string, error) {
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
	if len(matches) == 0 {
		return nil, errors.New(fmt.Sprintf("Category not found: %s", s))
	}
	return matches, nil
}

func (t PokemonTrie) GetCategory(s []string) ([]*PokemonEntry, error) {
	current := t.Root
	matches := make([]*PokemonEntry, 0)
	for _, char := range s {
		if _, ok := current.Children[char]; ok {
			current = current.Children[char]
			for _, p := range current.Children {
				matches = append(matches, p.Data...)
			}
		} else {
			return nil, errors.New(fmt.Sprintf("Could not find category: %s", s))
		}
	}
	return matches, nil
}

func TokenizeName(name string) []string {
	return strings.Split(name, "-")
}

type PokemonMatch struct {
	Entry      *PokemonEntry
	Categories []string
}

func (t PokemonTrie) MatchNameToken(s string) ([]*PokemonMatch, error) {
	matches := t.Root.MatchNameToken(s, []string{})
	if len(matches) > 0 {
		return matches, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Could not find name: %s", s))
	}
}

func (current Node) MatchNameToken(s string, keys []string) []*PokemonMatch {
	matches := []*PokemonMatch{}

	for _, entry := range current.Data {
		for _, tk := range TokenizeName(entry.Name) {
			if tk == s {
				matches = append(matches, &PokemonMatch{Entry: entry, Categories: keys})
			}
		}
	}
	for k, node := range current.Children {
		more := node.MatchNameToken(s, append(keys, k))
		if len(more) > 0 {
			matches = append(matches, more...)
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

func ReadTrieFromBytes(data []byte) PokemonTrie {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	d := &PokemonTrie{}

	err := dec.Decode(&d)
	check(err)

	return *d
}

func ReadMetadataFromBytes(data []byte) PokemonMetadata {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	d := &PokemonMetadata{}

	err := dec.Decode(&d)
	check(err)

	return *d
}

func WriteBytesToFile(data []byte, fpath string, compress bool) {
	ostream, err := os.Create(fpath)
	check(err)

	writer := bufio.NewWriter(ostream)
	if compress {
		writer.WriteString(string(Compress(data)))
	} else {
		writer.WriteString(string(data))
	}
	writer.Flush()
	ostream.Close()
}
