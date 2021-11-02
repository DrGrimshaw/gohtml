package gohtml

import (
	"reflect"
	"testing"
)

type internalStruct struct {
	Name string `html:"l=Name,e=span,c=name"`
}

func TestParseTag(t *testing.T) {
	sut := internalStruct{Name: "John"}

	v := reflect.ValueOf(sut)
	typeOfS := v.Type()

	hfs, err := parseTag(typeOfS.Field(0))

	assertNil(t, err)
	assertEqual(t, "Name", hfs.Label)
	assertEqual(t, "span", hfs.Element)
	assertEqual(t, "name", hfs.Class)
}
