package gohtml

import (
	"fmt"
	"testing"
)

const assertEqualErrStr = "expectations of test not met. expected: %s, actual: %s"

func assertEqual(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Fatal(fmt.Errorf(assertEqualErrStr, expected, actual))
	}
}

func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
