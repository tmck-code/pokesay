package test

import (
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

func TestCategoryPaths(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewEntry(4, "charmander"))

	expected := [][]string{
		{"small", "g1", "r"},
		{"small", "g1", "o"},
		{"medium", "g1", "o"},
		{"big", "g1", "o"},
		{"big", "g1"},
	}
	Assert(expected, t.Keys, t, test)

	expected = [][]string{
		{"big", "g1", "o"},
		{"big", "g1"},
	}
	result, err := t.FindKeys("big")
	pokesay.Check(err)
	Assert(expected, result, result, test)
}

func TestReadNames(test *testing.T) {
	result := pokedex.ReadNames("./data/pokemon.json")

	expected := map[string]pokedex.PokemonName{
		"bulbasaur": {English: "Bulbasaur", Japanese: "フシギダネ", JapanesePhonetic: "fushigidane"},
		"ivysaur":   {English: "Ivysaur", Japanese: "フシギソウ", JapanesePhonetic: "fushigisou"},
		"venusaur":  {English: "Venusaur", Japanese: "フシギバナ", JapanesePhonetic: "fushigibana"},
	}

	Assert(expected, result, result, test)
}
