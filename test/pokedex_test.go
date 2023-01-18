package test

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
		test.Fatalf("\nexpected = %#v \nresult = %#v \nobj = %#v", expected, result, obj)
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
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewPokemonEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewPokemonEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewPokemonEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewPokemonEntry(4, "charmander"))

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
		"bulbasaur": {English: "Bulbasaur", Japanese: "フシギダネ", JapanesePhonetic: "fushigidane"},
		"ivysaur":   {English: "Ivysaur", Japanese: "フシギソウ", JapanesePhonetic: "fushigisou"},
		"venusaur":  {English: "Venusaur", Japanese: "フシギバナ", JapanesePhonetic: "fushigibana"},
	}

	assert(expected, result, result, test)
}
