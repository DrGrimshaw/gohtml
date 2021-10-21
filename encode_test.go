package gohtml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type simpleStruct struct {
	Name string `label:"Name" element:"span" class:"name"`
}

type differentTypeStruct struct {
	Name string `label:"Name" element:"span" class:"name"`
	Number int `label:"Number" element:"span" class:"number"`
	Boolean bool `label:"Boolean" element:"span" class:"boolean"`
	Decimal float64 `label:"Decimal" element:"span" class:"decimal"`
}

type subTypeStruct struct {
	Name    string  `label:"Name" element:"span" class:"name"`
	Address address `label:"address" class:"address"`
}

type simpleStructContainer struct {
	Names []simpleStruct `element:"li"`
}

type arrStruct struct {
	NameList simpleStructContainer `label:"Names" element:"ul"`
}

type complexArrStruct struct {
	NameList struct{
		Names []string `element:"li"`
	} `label:"Names" element:"ul"`
}

type address struct {
	Street string `label:"Street" element:"span" class:"street"`
	County string `label:"County" element:"span" class:"county"`
	Country string `label:"Country" element:"span" class:"country"`
	Postcode string `label:"Postcode" element:"span" class:"postcode"`
}

func TestEncode_Simple(t *testing.T) {
	sut := subTypeStruct{
		Name: "John",
		Address: address{
			Street:   "1600 Pennsylvania Avenue NW",
			County:   "Washington",
			Country:  "DC",
			Postcode: "20500",
		},
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t, "<div><span>Name</span><span class='name'>John</span></div><div><span>address</span><div class='address'><div><span>Street</span><span class='street'>1600 Pennsylvania Avenue NW</span></div><div><span>County</span><span class='county'>Washington</span></div><div><span>Country</span><span class='country'>DC</span></div><div><span>Postcode</span><span class='postcode'>20500</span></div></div></div>", result)
}

func TestEncode_DifferentTypes(t *testing.T) {
	sut := differentTypeStruct{
		Name: "John",
		Number: 23,
		Boolean: false,
		Decimal: 3.45,
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t, "<div><span>Name</span><span class='name'>John</span></div><div><span>Number</span><span class='number'>23</span></div><div><span>Boolean</span><span class='boolean'>false</span></div><div><span>Decimal</span><span class='decimal'>3.45</span></div>", result)
}

func TestEncode_Arr(t *testing.T) {
	sut := arrStruct{
		NameList: simpleStructContainer{
			Names: []simpleStruct{
				{Name: "John Doe"},
				{Name: "Jane Doe"},
			},
		},
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t, "<div><span>Names</span><ul><li><div><span>Name</span><span class='name'>John Doe</span></div></li><li><div><span>Name</span><span class='name'>Jane Doe</span></div></li></ul></div>", result)
}

func TestEncode_ComplexArr(t *testing.T) {
	sut := complexArrStruct{
		NameList: struct {
			Names []string `element:"li"`
		}{
			Names: []string{"John Doe","Jane Doe"},
		},
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t, "<div><span>Names</span><ul><li>John Doe</li><li>Jane Doe</li></ul></div>", result)
}
