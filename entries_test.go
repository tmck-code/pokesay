package main

import (
	"testing"
	"fmt"
	"github.com/tmck-code/pokesay/src/pokedex"
)

func assert(expected interface{}, result interface{}, obj interface{}, test *testing.T) {
	if fmt.Sprintf("%v", expected) != fmt.Sprintf("%v", result) {
		test.Fatalf(`expected = %+v, result = %+v, obj = %+v`, expected, result, obj)
	}
}

func TestNewPokemonEntry(test *testing.T) {
	p := pokedex.NewPokemonEntry(1, "yo")
	assert(1, p.Index, p, test)
}

func TestTrieInsert(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert(
		[]string{"pokemon", "gen1", "regular"},
		pokedex.NewPokemonEntry(0, "pikachu"),
	)
	t.Insert(
		[]string{"pokemon", "gen1", "regular"},
		pokedex.NewPokemonEntry(1, "bulbasaur"),
	)

	result, err := t.GetCategory([]string{"pokemon", "gen1"})
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
	t.Insert([]string{"small", "gen1", "regular"}, pokedex.NewPokemonEntry(0, "pikachu"))
	t.Insert([]string{"small", "gen1", "other"}, pokedex.NewPokemonEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "gen1", "other"}, pokedex.NewPokemonEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "gen1", "other"}, pokedex.NewPokemonEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "gen1"}, pokedex.NewPokemonEntry(4, "charmander"))
	
	expected := [][]string{
		[]string{"small", "gen1", "regular"},
		[]string{"small", "gen1", "other"},
		[]string{"medium", "gen1", "other"},
		[]string{"big", "gen1", "other"},
		[]string{"big", "gen1"},
	}
	assert(expected, t.Keys, t, test)

	expected = [][]string{
		[]string{"big", "gen1", "other"},
		[]string{"big", "gen1"},
	}
	result, err := t.GetCategoryPaths("big")
	check(err)
	assert(expected, result, result, test)
}
