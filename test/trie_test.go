package test

import (
	"fmt"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

func TestNewPokemonEntry(test *testing.T) {
	p := pokedex.NewPokemonEntry(1, "yo")
	Assert(1, p.Index, p, test)
}

func TestTrieInsert(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(1, "bulbasaur"))

	result, err := t.GetCategory([]string{"p", "g1"})
	pokesay.Check(err)

	Assert(2, len(result), result, test)
	Assert(
		"[{Index: 0, Name: pikachu} {Index: 1, Name: bulbasaur}]",
		fmt.Sprintf("%s", result),
		result, test,
	)
}

func TestTrieFind(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewPokemonEntry(1, "bulbasaur"))
	t.Insert([]string{"x", "g1", "l"}, pokedex.NewPokemonEntry(2, "pikachu-other"))

	expected := []pokedex.PokemonMatch{
		{
			Entry:      &pokedex.PokemonEntry{Index: 0, Name: "pikachu"},
			Categories: []string{"p", "g1", "r"},
		},
		{
			Entry:      &pokedex.PokemonEntry{Index: 2, Name: "pikachu-other"},
			Categories: []string{"x", "g1", "l"},
		},
	}

	results, err := t.Find("pikachu")
	pokesay.Check(err)

	for i := 0; i <= len(results)-1; i++ {
		Assert(expected[i].Entry, results[i].Entry, results[i], test)
		Assert(expected[i].Categories, results[i].Categories, results[i], test)
	}
}
