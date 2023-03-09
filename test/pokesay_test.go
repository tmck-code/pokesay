package test

import (
	_ "embed"
	"testing"

	"github.com/tmck-code/pokesay/src/pokesay"
)

var (
	//go:embed data/total.txt
	GOBTotal []byte
)

func TestChooseByName(test *testing.T) {
}

func TestChooseByCategory(test *testing.T) {
}

func TestChooseByRandomIndex(test *testing.T) {
	resultTotal, result := pokesay.ChooseByRandomIndex(GOBTotal)
	Assert(9, resultTotal, test)
	Assert(0 <= result, true, test)
	Assert(9 >= result, true, test)
}
