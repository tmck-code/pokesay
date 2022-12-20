package main

import (
	"fmt"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

// Made my own basic test helper. Takes in an expected & result object of any type, and asserts
// that their Go syntax representations (%#v) are the same
func assert(expected interface{}, result interface{}, obj interface{}, test *testing.T) {
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", result) {
		test.Fatalf(`expected = %#v, result = %#v, obj = %#v`, expected, result, obj)
	}
}

func TestNewPokemonEntry(test *testing.T) {
	p := pokedex.NewPokemonEntry(1, "yo", "スイクン")
	assert(1, p.Index, p, test)
}

func TestTrieInsert(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu", "ピカチュウ"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(1, "bulbasaur", "フシギダネ"))

	result, err := t.GetCategory([]string{"p", "g1"})
	pokesay.Check(err)

	assert(2, len(result), result, test)
	assert(
		"[{Index: 0, Name: pikachu} {Index: 1, Name: bulbasaur}]",
		fmt.Sprintf("%s", result),
		result, test,
	)
}

func TestCategoryPaths(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu", "ピカチュウ"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewPokemonEntry(1, "bulbasaur", "フシギダネ"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewPokemonEntry(2, "bulbasaur", "フシギダネ"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewPokemonEntry(3, "bulbasaur", "フシギダネ"))
	t.Insert([]string{"big", "g1"}, pokedex.NewPokemonEntry(4, "charmander", "ヒトカゲ"))

	expected := [][]string{
		{"small", "g1", "r"},
		{"small", "g1", "o"},
		{"medium", "g1", "o"},
		{"big", "g1", "o"},
		{"big", "g1"},
	}
	assert(expected, t.Keys, t, test)

	expected = [][]string{
		{"big", "g1", "o"},
		{"big", "g1"},
	}
	result, err := t.GetCategoryPaths("big")
	pokesay.Check(err)
	assert(expected, result, result, test)
}

func TestReadNames(test *testing.T) {
	result := pokedex.ReadNames("./data/pokemon.json")

	expected := map[string]pokedex.PokemonName{
		"bulbasaur": {English: "Bulbasaur", Japanese: "フシギダネ", JapaneseRomaji: "fushigidane"},
		"ivysaur":   {English: "Ivysaur", Japanese: "フシギソウ", JapaneseRomaji: "fushigisou"},
		"venusaur":  {English: "Venusaur", Japanese: "フシギバナ", JapaneseRomaji: "fushigibana"},
	}

	assert(expected, result, result, test)
}
