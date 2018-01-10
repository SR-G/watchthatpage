package core

import (
	"fmt"
	"testing"
)

// AssertEquals allows to compare expected and actual values. Will print default error message in case of found differences.
func AssertEquals(t *testing.T, a interface{}, b interface{}) {
	AssertEqualsMsg(t, a, b, "")
}

// AssertEqualsMsg allows to compare expected and actual values.
// Will print provided message in case of found differences.
func AssertEqualsMsg(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("\nexpected : [%v]\nactual   : [%v]", a, b)
	}
	t.Fatal(message)
}
