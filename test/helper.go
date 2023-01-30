package test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	successMark = "✓"
	failMark    = "❌"
)

// Checks if the test is running in debug mode, i.e. has been run with the ENV var DEBUG=true.
// To do this, either first run `export DEBUG=true`, and then run the test command,
// or do it all at once with `DEBUG=true go test -v ./test“
func debug() bool {
	return os.Getenv("DEBUG") == "true"
}

// Fails a test with a formatted message showing the expected vs. result. (These are both printed in %#v form)
func Fail(expected interface{}, result interface{}, test *testing.T) {
	test.Fatalf("%s items don't match!\n> expected:\t%#v\n>   result:\t%#v\n", failMark, expected, result)
}

// Takes in an expected & result object, of any type.
// Asserts that their Go syntax representations (%#v) are the same.
// Prints a message on success if the ENV var DEBUG is set to "true".
// Fails the test if this is not true.
func Assert(expected interface{}, result interface{}, test *testing.T) {
	expectedString, resultString := fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", result)
	if expectedString == resultString {
		if debug() {
			fmt.Printf("%s items match!\n> expected:\t%s\n>   result:\t%s\n", successMark, expected, result)
		}
		return
	}
	Fail(expectedString, resultString, test)
}

// Takes in an expected slice of objects and an 'item' object, of any type
// Asserts that the 'item' is contained within the slice.
// Prints a message on success if the ENV var DEBUG is set to "true".
// Fails the test if this is not true.
func AssertContains[T any](slice []T, item T, test *testing.T) {
	for _, el := range slice {
		if reflect.DeepEqual(el, item) {
			if debug() {
				fmt.Printf("%s found expected item!\n>  item:\t%v\n> slice:\t%v\n", successMark, item, slice)
			}
			return
		}
	}
	Fail(slice, item, test)
}

// Flattens a given json string, removing all tabs, spaces and newlines
func FlattenJSON(json string) string {
	json = strings.Replace(json, "\n", "", -1)
	json = strings.Replace(json, "\t", "", -1)
	json = strings.Replace(json, " ", "", -1)
	return json
}
