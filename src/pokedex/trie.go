package pokedex

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
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

func (t Trie) WriteToFile(fpath string) {
	WriteStructToFile(t, fpath)
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

// Search every node in a trie for a given value.
// Returns a slice of all nodes that have a matching value.
// Returns an error if the value wasn't found
func (t Trie) Find(value string) ([]*PokemonMatch, error) {
	matches := t.Root.Find(value, []string{})
	if len(matches) > 0 {
		return matches, nil
	} else {
		return nil, fmt.Errorf("value not found: %s", value)
	}
}

func (current Node) Find(value string, keys []string) []*PokemonMatch {
	matches := []*PokemonMatch{}

	for _, entry := range current.Data {
		for _, tk := range TokenizeValue(entry.Value) {
			if tk == value {
				matches = append(matches, &PokemonMatch{Entry: entry, Keys: keys})
			}
		}
	}
	for k, node := range current.Children {
		more := node.Find(value, append(keys, k))
		if len(more) > 0 {
			matches = append(matches, more...)
		}
	}
	return matches
}

func TokenizeValue(value string) []string {
	return strings.Split(value, "-")
}

func (t Trie) FindKeyPaths(key string) ([][]string, error) {
	matches := [][]string{}
	for _, k := range t.Keys {
		for _, el := range k {
			if el == key {
				matches = append(matches, k)
			}
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return matches, nil
}

// given a
func (t Trie) FindByKeyPath(keyPath []string) ([]*Entry, error) {
	current := t.Root
	matches := make([]*Entry, 0)
	for _, char := range keyPath {
		if _, ok := current.Children[char]; ok {
			current = current.Children[char]
			for _, p := range current.Children {
				matches = append(matches, p.Data...)
			}
		} else {
			return nil, fmt.Errorf("key path not found: %s", keyPath)
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("key path not found: %s", keyPath)
	}
	return matches, nil
}
