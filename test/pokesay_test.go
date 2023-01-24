package test

import (
	"embed"
	"os"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

var (
	//go:embed data/cows/*cow
	GOBCowData embed.FS
)

func TestListCategories(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewEntry(4, "charmander"))

	result := pokesay.ListCategories(*t)
	expected := []string{"big", "g1", "medium", "o", "r", "small"}

	Assert(expected, result, result, test)
}

func TestReadEntry(test *testing.T) {
	result := pokesay.ReadPokemonCow(GOBCowData, "data/cows/1.cow")

	expected, err := os.ReadFile("data/cows/egg.cow")
	pokesay.Check(err)

	Assert(string(expected), string(result), string(result), test)
}
