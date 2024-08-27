// Package mapping is for best guess data transforms, where
// some sort of mapping has been provided to help fill in the blanks.
package mapping

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// jsonEncode takes a flat dotpath map and
// produces the json object.
// if it is an array then an array of json is generated,
// instead of several files json
func jsonEncode(toBuild []map[string]any) ([]byte, error) {
	if len(toBuild) == 1 {

		dest := make(map[string]any)
		for buildKey, b := range toBuild[0] {
			path := strings.Split(buildKey, ".")
			recurseMapBuild(dest, path, b)
		}

		return json.Marshal(dest)
	}

	destArr := make([]map[string]any, len(toBuild))
	for i, build := range toBuild {
		dest := make(map[string]any)

		for buildKey, b := range build {
			path := strings.Split(buildKey, ".")
			recurseMapBuild(dest, path, b)
		}
		destArr[i] = dest
	}

	return json.MarshalIndent(destArr, "", "    ")
	// create an array
}

// yamlEncode takes a flat dotpath map and
// produces the json object.
// if it is an array then an array of json is generated,
// instead of several files json
func yamlEncode(toBuild []map[string]any) ([]byte, error) {
	if len(toBuild) == 1 {

		dest := make(map[string]any)
		for buildKey, b := range toBuild[0] {
			path := strings.Split(buildKey, ".")
			recurseMapBuild(dest, path, b)
		}

		outs, err := yaml.Marshal(dest)
		return append([]byte("---\n"), outs...), err
	}

	destArr := make([]map[string]any, len(toBuild))
	for i, build := range toBuild {
		dest := make(map[string]any)

		for buildKey, b := range build {
			path := strings.Split(buildKey, ".")
			recurseMapBuild(dest, path, b)
		}
		destArr[i] = dest
	}

	outs, err := yaml.Marshal(destArr)
	return append([]byte("---\n"), outs...), err
	// create an array

}

func csvEncode(metaDataInput []map[string]any, paths map[string]map[int][]mdProperties) ([]byte, error) {

	if len(metaDataInput) == 0 {
		return nil, nil
	}

	// get the headers in a uniform order
	headerOrder := make([]string, len(paths))
	i := 0
	for k := range paths {
		headerOrder[i] = k
		i++
	}

	// @TODO improve csv output order
	// this is just to keep it uniform for demos
	slices.Sort(headerOrder)

	// write the new layout to csv
	bwrite := bytes.NewBuffer([]byte{})
	cwrite := csv.NewWriter(bwrite)

	err := cwrite.Write(headerOrder)
	if err != nil {
		return nil, err
	}

	for i, dataLine := range metaDataInput {
		line := make([]string, len(headerOrder))
		for i, head := range headerOrder {
			line[i] = fmt.Sprintf("%v", dataLine[head])

		}
		err := cwrite.Write(line)
		if err != nil {
			return nil, err
		}

		// flush all the values every 100 lines
		if i%100 == 0 {
			cwrite.Flush()
		}
	}
	// flush all the values
	cwrite.Flush()

	return bwrite.Bytes(), nil
}

func recurseMapBuild(body map[string]any, path []string, value any) {

	if len(path) > 1 {

		// if the next one is an array
		if pos, err := strconv.Atoi(path[1]); err == nil {
			// make the array as big as the position then fill in the blanks
			// later
			if _, ok := body[path[0]]; !ok {
				body[path[0]] = make([]any, pos)
			}

			// check for array positions
			/*
				arr := body[path[0]].([]any)
				size := len(arr)
				if pos+1 > size {
					arr = append(arr, make([]any, pos-size+1)...)
				}*/
			// reassign as arrays do not work like maps
			body[path[0]] = recurseArrayBuild(body[path[0]].([]any), path[1:], value)
		} else {

			// if value is a number then check some lengths but append
			if _, ok := body[path[0]]; !ok {
				body[path[0]] = make(map[string]any)
			}

			recurseMapBuild(body[path[0]].(map[string]any), path[1:], value)
		}
	} else {

		body[path[0]] = value
	}
}

// build multi dimensional arrays
func recurseArrayBuild(body []any, path []string, value any) []any {

	// get the position of 0
	if pos, err := strconv.Atoi(path[0]); err == nil {
		// make the array as big as the position then fill in the blanks
		// later

		// check for array size matching the cooridnates0
		size := len(body)
		if pos+1 > size {

			body = append(body, make([]any, pos-size+1)...)

		}

		// if the path ends as an array
		if len(path) == 1 {
			body[pos] = value
			// if there are more arrays
		} else if _, err := strconv.Atoi(path[1]); err == nil {
			// no nil maps
			if body[pos] == nil {
				body[pos] = make([]any, 0)
			}

			body[pos] = recurseArrayBuild(body[pos].([]any), path[1:], value)
		} else { // treat is back as a map

			if body[pos] == nil {
				body[pos] = make(map[string]any)
			}

			recurseMapBuild(body[pos].(map[string]any), path[1:], value)

		}

		// recurseMapBuild2(arr[pos], path[2:], value)
	}
	return body
}

// OrderedMap is a map that retains the order
// items were set. As well as logging any XML attributes
type OrderedMap struct {
	m          map[string]any
	order      []string
	attributes map[string][]xml.Attr
}

// Set the Key and value
// if the key already exists then the value is preserved
// and the original value overwritten.
func (o *OrderedMap) Set(key string, value any) {
	if !slices.Contains(o.order, key) {
		o.order = append(o.order, key)
	}

	o.m[key] = value
}

// AddAttr adds attributes to a parent key, with a key and value for the attribute
func (o *OrderedMap) AddAttr(parentKey, attrKey string, attrValue any) {
	input := ""
	if attrValue != nil {
		input = fmt.Sprintf("%v", attrValue)
	}

	o.attributes[parentKey] = append(o.attributes[parentKey], xml.Attr{Name: xml.Name{Local: attrKey}, Value: input})
}

// Get a key value
func (o *OrderedMap) Get(key string) (any, bool) {
	val, ok := o.m[key]
	return val, ok
}

// NewMap generates a new ordered map
func NewMap() *OrderedMap {
	return &OrderedMap{
		m:          make(map[string]any),
		order:      make([]string, 0),
		attributes: make(map[string][]xml.Attr),
	}
}

// XMLMapEntry is part of the xml encoder structure,
// it's used for encoding values that are not arrays or maps
type XMLMapEntry struct {
	XMLName xml.Name
	Value   any `xml:",chardata"`
}

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (o *OrderedMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(o.m) == 0 {
		return nil
	}

	//	start.Attr = []xml.Attr{{xml.Name{Local: "test"}, "tval"}, {xml.Name{Local: "test"}, "tval"}}
	if !reflect.DeepEqual(start, xml.StartElement{Name: xml.Name{Local: "OrderedMap"}}) {
		err := e.EncodeToken(start)
		if err != nil {
			return err
		}
	}

	/*
		insert some attribute gleaning code for adding too the start
	*/
	/*
	   prefix should be soemwhere saying where the start is.
	   ns: is a shorthand fir the the URI of the namespace
	*/
	for _, k := range o.order {
		// add any namespace stuff here
		v := o.m[k]

		switch input := v.(type) {
		case *OrderedMap:
			var err error
			if attr, ok := o.attributes[k]; ok {
				err = input.MarshalXML(e, xml.StartElement{Name: xml.Name{Local: k}, Attr: attr})
			} else {
				err = input.MarshalXML(e, xml.StartElement{Name: xml.Name{Local: k}})
			}
			if err != nil {
				return err
			}

		case []any:

			inputArray := XMLArray(input)
			inputArray.MarshalXML(e, xml.StartElement{Name: xml.Name{Local: k}})
		/*
			create an xml marshall array for arrays
		*/
		/*
			for _, v := range input {

			}
		*/
		// @TODO include more arrays that are not just any or float e.g. int
		case []float64:
			arrEncoder(e, input, k)
		case []int:
			arrEncoder(e, input, k)
		case []string:
			arrEncoder(e, input, k)

		default:
			e.Encode(XMLMapEntry{XMLName: xml.Name{Local: k}, Value: v})
		}
	}
	if !reflect.DeepEqual(start, xml.StartElement{Name: xml.Name{Local: "OrderedMap"}}) {
		return e.EncodeToken(start.End())
	}

	return nil
}

// encodes an Array of any type, to be used for when type != any ironically.
func arrEncoder[T any](e *xml.Encoder, input []T, k string) {
	for _, v := range input {
		e.Encode(XMLMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}
}

// XMLArray  implements a XML encoder
// it's used in conjunction with the ordered map
type XMLArray []any

// MarshalXML marshals the map to XML, with each key in the map being a
// tag and it's corresponding value being it's contents.
func (xa XMLArray) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(xa) == 0 {
		return nil
	}

	for _, v := range xa {
		switch input := v.(type) {
		case *OrderedMap:
			input.MarshalXML(e, start)
		case []any:
			inputArray := XMLArray(input)

			err := e.EncodeToken(start)
			if err != nil {
				return err
			}

			err = inputArray.MarshalXML(e, start)
			if err != nil {
				return err
			}
		/*
			create an xml marshall array for arrays
		*/
		/*
			for _, v := range input {

			}
		*/
		default:
			e.Encode(XMLMapEntry{XMLName: start.Name, Value: v})
		}
	}

	return e.EncodeToken(start.End())
}

// @TODO include the target namespace
func xmlBuild(toBuild []map[string]any, xmlProperties xMLEncoderInformation) ([]byte, error) {

	// @TODO format this so non xml schemas can still create xml documents

	keys := make([]string, len(toBuild[0]))
	i := 0
	for k := range toBuild[0] {
		keys[i] = k
		i++
	}
	slices.Sort(keys)
	translations := make(map[string][]string)

	// ASSIGN atributes here
	for _, k := range keys {
		mids := strings.Split(k, ".")
		out := ""
		for _, m := range mids {
			if _, err := strconv.Atoi(m); err != nil {
				if len(out) == 0 {
					out += m
				} else {
					out += "." + m
				}
			}
		}

		translations[out] = append(translations[out], k)
	}

	type pathAndAttribute struct {
		path      string
		attribute bool
	}
	// get Key order
	orderKeys := []pathAndAttribute{}

	for _, order := range xmlProperties.keyorder {
		matches, ok := translations[order]

		if ok {
			var at bool

			if _, ok := xmlProperties.attr[order]; ok {
				// then assign some attribute informatoin
				at = true
			}
			for _, m := range matches {
				orderKeys = append(orderKeys, pathAndAttribute{path: m, attribute: at})
			}
		}
	}

	dest := NewMap()
	/*
		ordered. Each order follows a map path.
		But each map path is only one item deep

	*/
	for _, buildKey := range orderKeys {
		path := strings.Split(buildKey.path, ".")
		val := toBuild[0][buildKey.path]
		// fmt.Println(path)
		recurseMapBuildOrder(dest, path, val, buildKey.attribute)
	}

	// assign the root namespace as well

	// group attriubtes by parent string.
	// create the array then assign

	// fmt.Println(dest.m["knowledgeItem"].(*OrderedMap).m["conceptSet"])

	// fmt.Println(string(b))
	return xml.MarshalIndent(dest, "", "    ")
}

func recurseMapBuildOrder(body *OrderedMap, path []string, value any, attribute bool) {

	if len(path) > 1 {

		// if the next one is an array
		if pos, err := strconv.Atoi(path[1]); err == nil {
			// make the array as big as the position then fill in the blanks
			// later
			if _, ok := body.Get(path[0]); !ok {
				body.Set(path[0], make([]any, pos))
			}

			// check for array positions
			/*
				arr := body[path[0]].([]any)
				size := len(arr)
				if pos+1 > size {
					arr = append(arr, make([]any, pos-size+1)...)
				}*/
			// reassign as arrays do not work like maps
			arr, _ := body.Get(path[0])

			body.Set(path[0], recurseArrayBuildOrder(arr.([]any), path[1:], value, attribute))
		} else {

			if len(path[1:]) == 1 && attribute {
				body.AddAttr(path[0], path[1], value)

			} else {

				// if value is a number then check some lengths but append
				if _, ok := body.Get(path[0]); !ok {
					body.Set(path[0], NewMap())
				}
				recurseMap, _ := body.Get(path[0])
				recurseMapBuildOrder(recurseMap.(*OrderedMap), path[1:], value, attribute)
			}
		}
	} else {

		body.Set(path[0], value)
	}
}

// build multi dimensional arrays
func recurseArrayBuildOrder(body []any, path []string, value any, attribute bool) []any {

	// get the position of 0
	if pos, err := strconv.Atoi(path[0]); err == nil {
		// make the array as big as the position then fill in the blanks
		// later

		// check for array size matching the cooridnates
		size := len(body)
		if pos+1 > size {

			body = append(body, make([]any, pos-size+1)...)

		}

		// no nil maps
		if _, err := strconv.Atoi(path[1]); err == nil {
			if body[pos] == nil {
				body[pos] = make([]any, 0)
			}

			body[pos] = recurseArrayBuildOrder(body[pos].([]any), path[1:], value, attribute)
		} else { // treat is back as a map

			if body[pos] == nil {
				body[pos] = NewMap()
			}

			recurseMapBuildOrder(body[pos].(*OrderedMap), path[1:], value, attribute)

		}

		// recurseMapBuild2(arr[pos], path[2:], value)
	}
	return body
}
