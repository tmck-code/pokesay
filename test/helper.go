package test

import (
	"fmt"
	"testing"
)

// Made my own basic test helper. Takes in an expected & result object of any type, and Asserts
// that their Go syntax representations (%#v) are the same
func Assert(expected interface{}, result interface{}, obj interface{}, test *testing.T) {
	// fmt.Printf("%#v %#v\n", expected, result)
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", result) {
		test.Fatalf("\nexpected = %#v \nresult = %#v \nobj = %#v", expected, result, obj)
	}
}
