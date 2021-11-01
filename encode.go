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
	IsRow bool
}

const labelTag = "label"
const elementTag = "element"
const classTag = "class"
const rowTag = "row"

// Encode finds and extracts data for HTML output
// The only error that returns is due to a malformed struct tags
func Encode(structure interface{}) (string, error) {
	return encode(structure, false, "")
}

func encode(structure interface{}, isRow bool, wrap string) (string, error) {
	v := reflect.ValueOf(structure)
	typeOfS := v.Type()

	output := ""

	for i := 0; i < v.NumField(); i++ {
		hfs, err := parseTag(typeOfS.Field(i))
		if err != nil {
			return "", err
		}

		subOutput, err := parseByType(v.Field(i).Interface(), hfs, isRow, wrap)
		if err != nil {
			return "", err
		}

		output += subOutput
	}

	return output, nil
}

func parseByType(i interface{}, hfs htmlFieldStructure, stripElement bool, wrapElement string) (string, error) {
	if wrapElement == "" {
		wrapElement = "%s"
	}

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

		arrOutput := ""

		for j := 0; j < s.Len(); j++ {
			if hfs.IsRow &&
				j == 0 &&
				s.Index(j).Kind() == reflect.Struct {
				titleRowOutput := ""
				for k := 0; k < s.Index(j).NumField(); k++ {
					arrItemTags, err := parseTag(s.Index(j).Type().Field(k))
					if err != nil {
						return "", err
					}
					titleRowOutput += fmt.Sprintf(colPattern, arrItemTags.Label)
				}
				arrOutput += fmt.Sprintf(titleRowPattern, titleRowOutput)
			}

			sliceOutput, err := parseByType(s.Index(j).Interface(), hfs, true, "")
			if err != nil {
				return "", err
			}

			arrOutput += sliceOutput
		}

		if hfs.IsRow {
			arrOutput = fmt.Sprintf("<table>%s</table>", arrOutput)
		}

		output += arrOutput

	case reflect.Struct:
		wrap := ""
		if hfs.IsRow {
			wrap = 	colPattern
		}

		subOutput, err := encode(i, hfs.IsRow, wrap)
		if err != nil {
			return "", err
		}

		output += sprintOutput(stripElement, label, hfs.Element, elementOptions, subOutput)
	default:
		output += sprintOutput(stripElement, label, hfs.Element, elementOptions, i)
	}

	output = fmt.Sprintf(wrapElement, output)

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
	isRow := false

	if htmlLabelTag, err := tags.Get(labelTag); err == nil {
		label = htmlLabelTag.Name
	}

	if htmlElementTag, err := tags.Get(elementTag); err == nil {
		element = htmlElementTag.Name
	}

	if htmlClassTag, err := tags.Get(classTag); err == nil {
		class = htmlClassTag.Name
	}

	if htmlRowTag, err := tags.Get(rowTag); err == nil && htmlRowTag.Name == "true" {
		isRow = true
		element = "tr"
	}

	return htmlFieldStructure{
		Label: label,
		Element: element,
		Class: class,
		IsRow: isRow,
	}, nil
}

const strippedPattern = "<%s%s>%v</%s>"
const rowPattern = "<tr>%v</tr>"
const titleRowPattern = "<thead><tr>%v</tr></thead>"
const colPattern = "<td>%v</td>"
const wrapPattern = "<div>%s<%s%s>%v</%s></div>"

func sprintOutput(isSlice bool, label string, element string, elementOptions string, object interface{}) string {
	if isSlice {
		return fmt.Sprintf(strippedPattern, element, elementOptions, object, element)
	} else {
		return fmt.Sprintf(wrapPattern, label, element, elementOptions, object, element)
	}
}
