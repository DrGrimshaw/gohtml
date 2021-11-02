package gohtml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type simpleStruct struct {
	Name string `html:"l=Name,e=span,c=name"`
}

type differentTypeStruct struct {
	Name string `html:"l=Name,e=span,c=name"`
	Number int `html:"l=Number,e=span,c=number"`
	Boolean bool `html:"l=Boolean,e=span,c=boolean"`
	Decimal float64 `html:"l=Decimal,e=span,c=decimal"`
}

type subTypeStruct struct {
	Name    string  `html:"l=Name,e=span,c=name"`
	Address address `html:"l=address,c=address"`
}

type simpleStructContainer struct {
	Names []simpleStruct `html:"e=li"`
}

type arrStruct struct {
	NameList simpleStructContainer `html:"l=Names,e=ul"`
}

type complexArrStruct struct {
	NameList struct{
		Names []string `html:"e=li"`
	} `html:"l=Names,e=ul"`
}

type tableStruct struct {
	People []person `html:"l=Data,row"`
}

type person struct{
	Name string `html:"e=span,l=Name,omitempty"`
	Age string `html:"e=span,l=Age"`
	Location string `html:"e=span,l=location"`
}

type address struct {
	Street string `html:"l=Street,e=span,c=street"`
	County string `html:"l=County,e=span,c=county"`
	Country string `html:"l=Country,e=span,c=country"`
	Postcode string `html:"l=Postcode,e=span,c=postcode"`
}

type complexTableStruct struct {
	Data []complexRowStruct `html:"l=Data,row"`
}

type complexRowStruct struct {
	Label            string `html:"l=Label,e=span"`
	Type             string `html:"l=Type,e=span"`
	Description string `html:"l=Description,e=span"`
	NameList struct {
		Names []string `html:"e=li"`
	} `html:"l=Name List,e=ul"`
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
			Names []string `html:"e=li"`
		}{
			Names: []string{"John Doe","Jane Doe"},
		},
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t, "<div><span>Names</span><ul><li>John Doe</li><li>Jane Doe</li></ul></div>", result)
}

func TestEncode_Table(t *testing.T) {
	sut := tableStruct{
		People: []person{
			{Name:"John Doe", Age: "29", Location: "Washington DC"},
			{Name:"Jane Doe", Age: "25", Location: "London"},
		},
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t, "<table><thead><tr><td>Name</td><td>Age</td><td>location</td></tr></thead><tr><td><span>John Doe</span></td><td><span>29</span></td><td><span>Washington DC</span></td></tr><tr><td><span>Jane Doe</span></td><td><span>25</span></td><td><span>London</span></td></tr></table>", result)
}

func TestEncode_TableComplex(t *testing.T) {
	sut := complexTableStruct{
		Data: []complexRowStruct{{
			Label:       "EX1",
			Type:        "Random",
			Description: "This is a random example",
			NameList: struct {
				Names []string `html:"e=li"`
			}{
				Names: []string{"John Doe", "Jane Doe"},
			},
		}},
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t, "<table><thead><tr><td>Label</td><td>Type</td><td>Description</td><td>Name List</td></tr></thead><tr><td><span>EX1</span></td><td><span>Random</span></td><td><span>This is a random example</span></td><td><ul><li>John Doe</li><li>Jane Doe</li></ul></td></tr></table>", result)
}

func TestEncode_OmitEmpty(t *testing.T) {
	sut := person{
		Age:      "22",
		Location: "Washington DC",
	}

	result, err := Encode(sut)

	assert.Nil(t, err)
	assert.Equal(t,"<div><span>Age</span><span>22</span></div><div><span>location</span><span>Washington DC</span></div>",result)
}
