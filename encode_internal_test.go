package gohtml

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type internalStruct struct {
	Name string `label:"Name" element:"span" class:"name"`
}

func TestParseTag(t *testing.T) {
	sut := internalStruct{Name: "John"}

	v := reflect.ValueOf(sut)
	typeOfS := v.Type()

	hfs, err := parseTag(typeOfS.Field(0))

	assert.Nil(t, err)
	assert.Equal(t, "Name", hfs.Label)
	assert.Equal(t, "span", hfs.Element)
	assert.Equal(t, "name", hfs.Class)
}
