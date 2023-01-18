package pokedex

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Entry struct {
	Value string `json:"value"`
	Index int    `json:"index"`
}

func (p Entry) String() string {
	return fmt.Sprintf("{Index: %d, Value: %s}", p.Index, p.Value)
}

type Node struct {
	Children map[string]*Node `json:"children"`
	Data     []*Entry         `json:"data"`
}

type Trie struct {
	Root *Node      `json:"root"`
	Len  int        `json:"len"`
	Keys [][]string `json:"keys"`
}

func NewEntry(idx int, value string) *Entry {
	return &Entry{
		Index: idx,
		Value: value,
	}
}

func NewNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
	}
}

func NewTrie() *Trie {
	return &Trie{
		Len:  0,
		Root: NewNode(),
		Keys: make([][]string, 0),
	}
}

func NewTrieFromBytes(data []byte) Trie {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	d := &Trie{}

	err := dec.Decode(&d)
	Check(err)

	return *d
}

func (t Trie) WriteToFile(fpath string) {
	ostream, err := os.Create(fpath)
	Check(err)

	writer := bufio.NewWriter(ostream)
	enc := gob.NewEncoder(writer)
	enc.Encode(t)
	writer.Flush()
	ostream.Close()
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

func (t *Trie) ToString(indentation ...int) string {
	if len(indentation) == 1 {
		json, err := json.MarshalIndent(t, "", strings.Repeat(" ", indentation[0]))
		Check(err)
		return string(json)
	} else {
		json, err := json.Marshal(t)
		Check(err)
		return string(json)
	}
}

func (t *Trie) Insert(keys []string, data *Entry) {
	current := t.Root
	found := false
	for _, k := range t.Keys {
		if Equal(k, keys) {
			found = true
			break
		}
	}
	if !found {
		t.Keys = append(t.Keys, keys)
	}
	for _, key := range keys {
		if _, ok := current.Children[key]; ok {
			current = current.Children[key]
		} else {
			current.Children[key] = NewNode()
			current = current.Children[key]
			current.Data = make([]*Entry, 0)
		}
	}
	current.Data = append(current.Data, data)
	t.Len += 1
}

type PokemonMatch struct {
	Entry *Entry
	Keys  []string
}

func (t Trie) Find(s string) ([]*PokemonMatch, error) {
	matches := t.Root.Find(s, []string{})
	if len(matches) > 0 {
		return matches, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Could not find value: %s", s))
	}
}

func (current Node) Find(s string, keys []string) []*PokemonMatch {
	matches := []*PokemonMatch{}

	for _, entry := range current.Data {
		for _, tk := range TokenizeValue(entry.Value) {
			if tk == s {
				matches = append(matches, &PokemonMatch{Entry: entry, Keys: keys})
			}
		}
	}
	for k, node := range current.Children {
		more := node.Find(s, append(keys, k))
		if len(more) > 0 {
			matches = append(matches, more...)
		}
	}
	return matches
}

func TokenizeValue(value string) []string {
	return strings.Split(value, "-")
}

func (t Trie) FindKeys(s string) ([][]string, error) {
	matches := [][]string{}
	for _, k := range t.Keys {
		for _, el := range k {
			if el == s {
				matches = append(matches, k)
			}
		}
	}
	if len(matches) == 0 {
		return nil, errors.New(fmt.Sprintf("Category not found: %s", s))
	}
	return matches, nil
}

// given a
func (t Trie) FindKeyEntries(s []string) ([]*Entry, error) {
	current := t.Root
	matches := make([]*Entry, 0)
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
	if len(matches) == 0 {
		return nil, errors.New(fmt.Sprintf("Could not find category: %s", s))
	}
	return matches, nil
}
