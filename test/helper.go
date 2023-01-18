package test

import (
	"fmt"
	"strings"
	"testing"
)

func Fail(expected interface{}, result interface{}, obj interface{}, test *testing.T) {
	test.Fatalf("\nexpected = %#v \nresult = %#v \nobj = %#v", expected, result, obj)
}

// Made my own basic test assertion helper. Takes in an expected & result object of any type,
// and Asserts that their Go syntax representations (%#v) are the same
func Assert(expected interface{}, result interface{}, obj interface{}, test *testing.T) {
	// fmt.Printf("%#v %#v\n", expected, result)
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", result) {
		Fail(expected, result, obj, test)
	}
}

// Flattens a given json string, removing all tabs, spaces and newlines
func FlattenJSON(json string) string {
	json = strings.Replace(json, "\n", "", -1)
	json = strings.Replace(json, "\t", "", -1)
	json = strings.Replace(json, " ", "", -1)
	return json
}
