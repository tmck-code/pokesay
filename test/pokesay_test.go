package test

import (
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
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

	Assert(expected, result, test)
}

func TestChooseByName(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewEntry(2, "charmander"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewEntry(4, "charmander"))

	result := pokesay.ChooseByName("charmander", *t)
	expected := []pokedex.PokemonMatch{
		{Entry: &pokedex.Entry{Index: 2, Value: "charmander"}, Keys: []string{"medium", "g1", "o"}},
		{Entry: &pokedex.Entry{Index: 4, Value: "charmander"}, Keys: []string{"big", "g1"}},
	}

	AssertContains(expected, *result, test)
}

func TestChooseByCategory(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewEntry(4, "charmander"))

	result, keys := pokesay.ChooseByCategory("big", *t)

	expectedKeys := [][]string{{"big", "g1", "o"}, {"big", "g1"}}
	expectedResult := []pokedex.Entry{
		{Index: 3, Value: "bulbasaur"},
		{Index: 4, Value: "charmander"},
	}
	AssertContains(expectedResult, *result, test)
	AssertContains(expectedKeys, keys, test)
}
