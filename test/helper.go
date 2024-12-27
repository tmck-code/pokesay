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

// Checks if the test is running in Debug mode, i.e. has been run with the ENV var DEBUG=true.
// To do this, either first run `export DEBUG=true`, and then run the test command,
// or do it all at once with `DEBUG=true go test -v ./test“
func Debug() bool {
	return os.Getenv("DEBUG") == "true"
}

// Fails a test with a formatted message showing the expected vs. result. (These are both printed in %#v form)
func Fail(expected interface{}, result interface{}, test *testing.T) {
	test.Fatalf("\n\x1b[38;5;196m%s items don't match!\x1b[0m\n> expected:\t%#v\x1b[0m\n>   result:\t%#v\x1b[0m\n\n", failMark, expected, result)
}

// Takes in an expected & result object, of any type.
// Asserts that their Go syntax representations (%#v) are the same.
// Prints a message on success if the ENV var DEBUG is set to "true".
// Fails the test if this is not true.
func Assert(expected interface{}, result interface{}, test *testing.T) {
	expectedString, resultString := fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", result)
	if expectedString == resultString {
		if Debug() {
			fmt.Printf("\x1b[38;5;46m%s items match!\x1b[0m\n> expected:\t%#v\x1b[0m\n>   result:\t%#v\x1b[0m\n\n", successMark, expected, result)
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
			if Debug() {
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
