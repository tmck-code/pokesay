package pokedex

import (
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
	Root     *Node      `json:"root"`
	Len      int        `json:"len"`
	KeyPaths [][]string `json:"keys"`
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
		Len:      0,
		Root:     NewNode(),
		KeyPaths: make([][]string, 0),
	}
}

func NewTrieFromBytes(data []byte) Trie {
	return ReadStructFromBytes[Trie](data)
}

func (t Trie) WriteToFile(fpath string) {
	WriteStructToFile(t, fpath)
}

func (t *Trie) ToString(indentation ...int) string {
	return StructToJSON(*t, indentation...)
}

func (t *Trie) Insert(keyPath []string, data *Entry) {
	// fmt.Println("adding", keyPath, data)
	current := t.Root
	found := false
	for _, k := range t.KeyPaths {
		if KeyPathsEqual(k, keyPath) {
			found = true
			break
		}
	}
	if !found {
		t.KeyPaths = append(t.KeyPaths, keyPath)
	}
	for _, key := range keyPath {
		if _, ok := current.Children[key]; ok {
			current = current.Children[key]
		} else {
			current.Children[key] = NewNode()
			current = current.Children[key]
		}
	}
	if current.Data == nil {
		current.Data = []*Entry{data}
	} else {
		current.Data = append(current.Data, data)
	}
	t.Len += 1
}

func KeyPathsEqual(a, b []string) bool {
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

type PokemonMatch struct {
	Entry *Entry
	Keys  []string
}

// Search every node in a trie for a given value.
// Returns a slice of all nodes that have a matching value.
// Returns an error if the value wasn't found
func (t Trie) Find(value string, findFirst bool) ([]*PokemonMatch, error) {
	matches := t.Root.Find(value, []string{}, findFirst)
	if len(matches) > 0 {
		return matches, nil
	} else {
		return nil, fmt.Errorf("value not found: %s", value)
	}
}

func (current Node) Find(value string, keys []string, findFirst bool) []*PokemonMatch {
	matches := []*PokemonMatch{}

	for _, entry := range current.Data {
		if strings.ToLower(entry.Value) == strings.ToLower(value) {
			if findFirst {
				return []*PokemonMatch{{Entry: entry, Keys: keys}}
			}
			// fmt.Printf("%s matches %s: %v - %d\n", value, entry.Value, keys, len(current.Data))
			matches = append(matches, &PokemonMatch{Entry: entry, Keys: keys})
		}
	}
	for k, node := range current.Children {
		more := node.Find(value, append(keys, k), findFirst)
		if len(more) > 0 {
			if findFirst {
				return more
			}
			matches = append(matches, more...)
		}
	}
	return matches
}

func (t Trie) FindKeyPaths(key string) ([][]string, error) {
	matches := [][]string{}
	for _, k := range t.KeyPaths {
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
