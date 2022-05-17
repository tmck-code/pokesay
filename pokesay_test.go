package main

import (
	"fmt"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
)

// Made my own basic test helper. Takes in an expected & result object of any type, and asserts
// that their Go syntax representations (%#v) are the same
func assert(expected interface{}, result interface{}, obj interface{}, test *testing.T) {
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", result) {
		test.Fatalf(`expected = %#v, result = %#v, obj = %#v`, expected, result, obj)
	}
}

func TestNewPokemonEntry(test *testing.T) {
	p := pokedex.NewPokemonEntry(1, "yo")
	assert(1, p.Index, p, test)
}

func TestTrieInsert(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(1, "bulbasaur"))

	result, err := t.GetCategory([]string{"p", "g1"})
	check(err)

	assert(2, len(result), result, test)
	assert(
		"[{Index: 0, Name: pikachu} {Index: 1, Name: bulbasaur}]",
		fmt.Sprintf("%s", result),
		result, test,
	)
}

func TestCategoryPaths(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewPokemonEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewPokemonEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewPokemonEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewPokemonEntry(4, "charmander"))

	expected := [][]string{
		[]string{"small", "g1", "r"},
		[]string{"small", "g1", "o"},
		[]string{"medium", "g1", "o"},
		[]string{"big", "g1", "o"},
		[]string{"big", "g1"},
	}
	assert(expected, t.Keys, t, test)

	expected = [][]string{
		[]string{"big", "g1", "o"},
		[]string{"big", "g1"},
	}
	result, err := t.GetCategoryPaths("big")
	check(err)
	assert(expected, result, result, test)
}

func TestListCategories(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewPokemonEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewPokemonEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewPokemonEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewPokemonEntry(4, "charmander"))

	result := ListCategories(*t)
	assert([]string{"big", "g1", "medium", "o", "r", "small"}, result, result, test)
}
