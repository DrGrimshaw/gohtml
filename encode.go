package gohtml

import (
	"fmt"
	"github.com/fatih/structtag"
	"reflect"
)

type htmlFieldStructure struct {
	Label string
	Element string
	Class string
}

const labelTag = "label"
const elementTag = "element"
const classTag = "class"

// Encode finds and extracts data for HTML output
// The only error that returns is due to a malformed struct tags
func Encode(structure interface{}) (string, error) {
	v := reflect.ValueOf(structure)
	typeOfS := v.Type()

	output := ""

	for i := 0; i < v.NumField(); i++ {
		hfs, err := parseTag(typeOfS.Field(i))
		if err != nil {
			return "", err
		}

		subOutput, err := parseByType(v.Field(i).Interface(), hfs, false)
		if err != nil {
			return "", err
		}

		output += subOutput
	}

	return output, nil
}

func parseByType(i interface{}, hfs htmlFieldStructure, isSlice bool) (string, error) {
	elementOptions := ""
	if hfs.Class != "" {
		elementOptions = fmt.Sprintf(" class='%s'", hfs.Class)
	}

	label := ""
	if hfs.Label != "" {
		label = fmt.Sprintf("<span>%s</span>", hfs.Label)
	}

	output := ""
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(i)
		for j := 0; j < s.Len(); j++ {
			sliceOutput, err := parseByType(s.Index(j).Interface(), hfs, true)
			if err != nil {
				return "", err
			}

			output += sliceOutput
		}
	case reflect.Struct:
		subOutput, err := Encode(i)
		if err != nil {
			return "", err
		}

		if isSlice {
			output += fmt.Sprintf("<%s%s>%s</%s>", hfs.Element, elementOptions, subOutput, hfs.Element)
		} else {
			output += fmt.Sprintf("<div>%s<%s%s>%s</%s></div>", label, hfs.Element, elementOptions, subOutput, hfs.Element)
		}
	default:
		output += fmt.Sprintf("<div>%s<%s%s>%v</%s></div>", label, hfs.Element, elementOptions, i, hfs.Element)
	}

	return output, nil
}

func parseTag(field reflect.StructField) (htmlFieldStructure, error) {
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return htmlFieldStructure{}, err
	}

	label := ""
	element := "div"
	class := ""

	if htmlLabelTag, err := tags.Get(labelTag); err == nil {
		label = htmlLabelTag.Name
	}

	if htmlElementTag, err := tags.Get(elementTag); err == nil {
		element = htmlElementTag.Name
	}

	if htmlClassTag, err := tags.Get(classTag); err == nil {
		class = htmlClassTag.Name
	}

	return htmlFieldStructure{
		Label: label,
		Element: element,
		Class: class,
	}, nil
}
