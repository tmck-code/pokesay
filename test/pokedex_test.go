package test

import (
	"os"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

func TestReadNames(test *testing.T) {
	result := pokedex.ReadNames("./data/pokemon.json")

	expected := map[string]pokedex.PokemonName{
		"bulbasaur": {English: "Bulbasaur", Japanese: "フシギダネ", JapanesePhonetic: "fushigidane"},
		"ivysaur":   {English: "Ivysaur", Japanese: "フシギソウ", JapanesePhonetic: "fushigisou"},
		"venusaur":  {English: "Venusaur", Japanese: "フシギバナ", JapanesePhonetic: "fushigibana"},
	}

	Assert(expected, result, result, test)
}

func TestReadStructFromBytes(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(1, "bulbasaur"))

	t.WriteToFile("test.txt")

	data, err := os.ReadFile("test.txt")
	pokesay.Check(err)
	result := pokedex.ReadStructFromBytes[pokedex.Trie](data)
	Assert(t.ToString(), result.ToString(), result.ToString(), test)
}
