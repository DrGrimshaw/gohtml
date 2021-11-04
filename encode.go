package gohtml

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// htmlFieldStructure is the struct form of the tag result
type htmlFieldStructure struct {
	Label     string
	Element   string
	Class     string
	IsRow     bool
	OmitEmpty bool
}

// Initial root tag searched for in metadata `html:"<opts...>"`
const rootTag = "html"

// Sub options, following tags to be used as follows `html:"l=Example,e=span,c=customclass"`
const labelTag = "label"
const shortLabelTag = "l"
const elementTag = "element"
const shortElementTag = "e"
const classTag = "class"
const shortClassTag = "c"

// Sub options without values to be used as follows `html:"row,omitempty"`
const optRowTag = "row"
const optOmitEmptyTag = "omitempty"

// Different patterns used for formatting output throughout
const labelPattern = "<span>%s</span>"

// Stripped pattern used for slices
const strippedPattern = "<%s%s>%v</%s>"

// Wrap pattern used by default
const wrapPattern = "<div>%s<%s%s>%v</%s></div>"

// These patterns are used in the row option
const tablePattern = "<table>%s</table>"
const titleRowPattern = "<thead><tr>%v</tr></thead>"
const colPattern = "<td>%v</td>"

// Encode finds and extracts data for HTML output
// The only error that returns is due to a malformed struct tags
func Encode(structure interface{}) (string, error) {
	return encode(structure, false, "")
}

func encode(structure interface{}, isRow bool, wrap string) (string, error) {
	v := reflect.ValueOf(structure)
	typeOfS := v.Type()

	var output strings.Builder

	for i := 0; i < v.NumField(); i++ {
		hfs, err := parseTag(typeOfS.Field(i))
		if err != nil {
			return "", err
		}

		subOutput, err := parseByType(v.Field(i).Interface(), hfs, isRow, wrap)
		if err != nil {
			return "", err
		}

		output.WriteString(subOutput)
	}

	return output.String(), nil
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
		label = fmt.Sprintf(labelPattern, hfs.Label)
	}

	if hfs.OmitEmpty && isEmptyValue(reflect.ValueOf(i)) {
		return "", nil
	}

	var output strings.Builder
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(i)

		var arrOutput strings.Builder

		for j := 0; j < s.Len(); j++ {
			if hfs.IsRow &&
				j == 0 &&
				s.Index(j).Kind() == reflect.Struct {
				var titleRowOutput strings.Builder
				for k := 0; k < s.Index(j).NumField(); k++ {
					arrItemTags, err := parseTag(s.Index(j).Type().Field(k))
					if err != nil {
						return "", err
					}
					titleRowOutput.WriteString(fmt.Sprintf(colPattern, arrItemTags.Label))
				}
				arrOutput.WriteString(fmt.Sprintf(titleRowPattern, titleRowOutput.String()))
			}

			sliceOutput, err := parseByType(s.Index(j).Interface(), hfs, true, "")
			if err != nil {
				return "", err
			}

			arrOutput.WriteString(sliceOutput)
		}

		if hfs.IsRow {
			output.WriteString(fmt.Sprintf(tablePattern, arrOutput.String()))
		} else {
			output.WriteString(arrOutput.String())
		}

	case reflect.Struct:
		wrap := ""
		if hfs.IsRow {
			wrap = colPattern
		}

		subOutput, err := encode(i, hfs.IsRow, wrap)
		if err != nil {
			return "", err
		}

		output.WriteString(sprintOutput(stripElement, label, hfs.Element, elementOptions, subOutput))
	default:
		output.WriteString(sprintOutput(stripElement, label, hfs.Element, elementOptions, i))
	}

	return fmt.Sprintf(wrapElement, output.String()), nil
}

func parseTag(field reflect.StructField) (htmlFieldStructure, error) {
	hfs := htmlFieldStructure{
		Label:     "",
		Element:   "div",
		Class:     "",
		IsRow:     false,
		OmitEmpty: false,
	}

	var subOptSearch = regexp.MustCompile(`.*=.*`)

	if htmlTag, ok := field.Tag.Lookup(rootTag); ok {
		for _, v := range strings.Split(htmlTag, ",") {
			if subOptSearch.MatchString(v) {
				subOpts := strings.Split(v, "=")

				if len(subOpts) != 2 {
					continue
				}

				switch subOpts[0] {
				case labelTag, shortLabelTag:
					hfs.Label = subOpts[1]
				case elementTag, shortElementTag:
					hfs.Element = subOpts[1]
				case classTag, shortClassTag:
					hfs.Class = subOpts[1]
				}

				continue
			}

			switch v {
			case optRowTag:
				hfs.IsRow = true
				hfs.Element = "tr"
			case optOmitEmptyTag:
				hfs.OmitEmpty = true
			}
		}
	}

	return hfs, nil
}

// sprintOutput decides on the method of formatting to use depending on whether it's a slice
func sprintOutput(isSlice bool, label string, element string, elementOptions string, object interface{}) string {
	if isSlice {
		return fmt.Sprintf(strippedPattern, element, elementOptions, object, element)
	}

	return fmt.Sprintf(wrapPattern, label, element, elementOptions, object, element)
}

// isEmptyValue func was lifted wholesale from the JSON encoding package as
// it holds the exact same functionality required for omitempty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		return v.Interface() == reflect.New(reflect.TypeOf(v.Interface())).Elem().Interface()
	}
	return false
}
