package pokedex

import (
	"errors"
	"fmt"
)

type PokemonEntry struct {
	Name  string
	Index int
}

func (p PokemonEntry) String() string {
	return fmt.Sprintf("{Index: %d, Name: %s}", p.Index, p.Name)
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

func (t PokemonTrie) GetCategoryPaths(s string) ([][]string, error) {
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
	if len(matches) == 0 {
		return nil, errors.New(fmt.Sprintf("Could not find category: %s", s))
	}
	return matches, nil
}