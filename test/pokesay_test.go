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

func TestChooseRandomCategory(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewEntry(4, "charmander"))

	resultKeys, err := t.FindKeyPaths("big")
	pokesay.Check(err)

	expectedKeys := [][]string{{"big", "g1", "o"}, {"big", "g1"}}
	Assert(expectedKeys, resultKeys, test)

	resultKey, result := pokesay.ChooseRandomCategory(resultKeys, *t)
	AssertContains(expectedKeys, resultKey, test)

	expectedResult := []pokedex.Entry{
		{Index: 3, Value: "bulbasaur"},
		{Index: 4, Value: "charmander"},
	}
	AssertContains(expectedResult, *result[0], test)
}
