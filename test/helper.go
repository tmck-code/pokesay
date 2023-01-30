package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Fails a test with a formatted message showing the expected vs. result.
// (These are both printed in %#v form)
func Fail(expected interface{}, result interface{}, test *testing.T) {
	test.Fatalf("\nexpected = %#v \nresult = %#v \n", expected, result)
}

// Takes in an expected & result object, of any type.
// Asserts that their Go syntax representations (%#v) are the same.
// Fails the test if this is not true.
func Assert(expected interface{}, result interface{}, test *testing.T) {
	expectedString, resultString := fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", result)
	if expectedString != resultString {
		Fail(expectedString, resultString, test)
	}
}

// Takes in an expected collection of objects and an 'item' object, of any type
// Asserts that the 'item' is contained within the collection.
// Fails the test if this is not true.
func AssertContains[T any](collection []T, item T, test *testing.T) {
	for _, el := range collection {
		if reflect.DeepEqual(el, item) {
			return
		}
	}
	Fail(collection, item, test)
}

// Flattens a given json string, removing all tabs, spaces and newlines
func FlattenJSON(json string) string {
	json = strings.Replace(json, "\n", "", -1)
	json = strings.Replace(json, "\t", "", -1)
	json = strings.Replace(json, " ", "", -1)
	return json
}
