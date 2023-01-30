package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Fail(expected interface{}, result interface{}, test *testing.T) {
	test.Fatalf("\nexpected = %#v \nresult = %#v \n", expected, result)
}

// Made my own basic test assertion helper. Takes in an expected & result object of any type,
// and Asserts that their Go syntax representations (%#v) are the same
func Assert(expected interface{}, result interface{}, test *testing.T) {
	// fmt.Printf("%#v %#v\n", expected, result)
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", result) {
		Fail(expected, result, test)
	}
}

func AssertContains[T any](collection []T, item T, test *testing.T) {
	found := false
	for _, el := range collection {
		if reflect.DeepEqual(el, item) {
			found = true
			break
		}
	}

	if !found {
		Fail(collection, item, test)
	}
}

// Flattens a given json string, removing all tabs, spaces and newlines
func FlattenJSON(json string) string {
	json = strings.Replace(json, "\n", "", -1)
	json = strings.Replace(json, "\t", "", -1)
	json = strings.Replace(json, " ", "", -1)
	return json
}
