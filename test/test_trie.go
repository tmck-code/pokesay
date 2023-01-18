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
