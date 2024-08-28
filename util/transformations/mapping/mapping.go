// Package mapping is for best guess data transforms, where
// some sort of mapping has been provided to help fill in the blanks.
package mapping

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// Options is the json layout for the mapping information
// in the Metarex register.
type Options struct {
	// Are types converted - yet to be implemented
	ConvertTypes bool `json:"convertTypes"`
	// The key to store missed fields.
	// if its empty the fields are not saved
	MissedTags string `json:"MissedFieldsKey"`
	// MappingDefinitions has the format : name of field [list of alternative field names]
	MappingDefinitions map[string][]string `json:"mappingDefinitions"`
	// MaxDepth may also be required
}

// Action contains all the information for the
// mapping metadata from one type to another action.
type Action struct {
	// The location of the output schema
	OutputSchema string
	// Any other useful information
	MrxID                     string
	InputFormat, OutputFormat string // data Types
	InputTiming               string // how to process the data generically - embedded or clocked
	Mapping                   Options
}

// mdProperties gives the information about a target field
type mdProperties struct {
	// the array properties
	dimension []any
	// the full path of the field
	fullpath string
	// the data
	dataType string
	// any bonus stuff like enum - later
	required bool
}

// ActionType describes the mapping's action for logging
func (m *Action) ActionType() string {
	return fmt.Sprintf("mapping to %v", m.MrxID)
}

// DataID gives the metarex ID of the data being transformed
func (m *Action) DataID() string {
	return m.MrxID
}

/*
Transform maps the metadata from one type to another,
using the destination metadata schema as a translation tool.

It follows these steps for transforming data:

  - flattens the input data.
  - it finds the flat destination layout from the schema.
  - it matches the source fields with any destination fields.
  - if multiple destination fields are found then the matching is based on similar depth of the data.
  - Fields are handled in alphabetical order, to preserve the array order of objects

Default behaviours:

  - 1 to 1 mapping of source to destination fields.
  - objects are built as they go, no object is transformed from one to another due to the flat nature of the mapping.
  - fields which aren't translated are discarded, they can be saved as an object with the "MissedTags" field.
  - mapping dictionary is expected field key with an array of potential field names.

Default transformations type:

  - everything is transformed into a string e.g. fmt.Sprintf("%v", data)
  - both integers and floats become floats.
  - floats are floor() to ints, ints are not floor().
  - boolean values are not transformed.
  - single values are appended to an array, of same type.
  - arrays types follow the same rules as the single values. e.g. floats are floor() into an integer array.
*/
func (m *Action) Transform(input []byte, _ url.Values) ([]byte, error) {

	// convert each data point at a time

	// extract the metadata in a generic format
	var metaDataInput []map[string]any
	var err error
	// Get the input metadata
	switch m.InputFormat {
	case "text/csv":
		metaDataInput, err = csvDecode(input)
	case "application/json":
		metaDataInput, err = jsonDecode(input, m.InputTiming)
	case "application/yaml":
		metaDataInput, err = yamlDecode(input, m.InputTiming)
	default:
		return nil, fmt.Errorf("unknown or invalid mime type of input data: %v", m.InputFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("error extracting the input data: %v", err)
	}

	// Extract the destination properties
	var mdPaths map[string]map[int][]mdProperties
	var xmlProperties *xMLEncoderInformation
	schemaPath := strings.Split(m.OutputSchema, ".")

	// @TODO include some data formating like header order for csvs
	// extract the file layout based on a schema
	switch strings.ToLower(schemaPath[len(schemaPath)-1]) {
	case "json":
		mdPaths, err = jsonSchemaExtract(m.OutputSchema)
	case "xml", "xsd":
		mdPaths, xmlProperties, err = xmlSchemaExtract(m.OutputSchema)
	default:
		return nil, fmt.Errorf("unknown schema extension: %v", schemaPath[len(schemaPath)-1])
	}

	if err != nil {
		return nil, fmt.Errorf("error extracting the layout of the output data: %v", err)
	}

	// translate each row of metadata into the end result
	translatedData := make([]map[string]any, len(metaDataInput))

	for j, input := range metaDataInput {
		translatedData[j] = dataTranslate(input, mdPaths, m.Mapping.MappingDefinitions, m.Mapping.MissedTags)

		// check that required properties all have values
		// @todo include for json schemas as well
		// if xmlProperties != nil {
		translatedData[j] = requiredPropertiesCheck(mdPaths, translatedData[j])

		if xmlProperties != nil {
			translatedData[j][xmlProperties.rootElement+".xmlns"] = xmlProperties.targetNameSpace
		}
	}

	// convert the generic format into the bytes of the destination type
	var outputBytes []byte
	switch m.OutputFormat {
	case "text/csv":
		outputBytes, err = csvEncode(translatedData, mdPaths)
	case "application/json":
		outputBytes, err = jsonEncode(translatedData)
	case "application/yaml":
		outputBytes, err = yamlEncode(translatedData)
	case "application/xml":
		if xmlProperties == nil {
			xmlProperties = &xMLEncoderInformation{}
		}
		outputBytes, err = xmlBuild(translatedData, *xmlProperties)
	default:
		return nil, fmt.Errorf("unknown or invalid mime type of output data %v", m.OutputFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("error saving the data as the output format: %v", err)
	}

	return outputBytes, nil
}

// constants that describe the various types
// to be used with the mapping
// based off the types used by json schema
const (
	String  = "string"
	Number  = "number"
	Integer = "integer"
	Boolean = "boolean"

	StringArray  = "stringarray"
	IntegerArray = "integerarray"
	FloatArray   = "numberarray"
	ObjectArray  = "objectarray"
	BooleanArray = "booleanarray"
)

// a map to save files that are missed
// with the ability to save as a string
type overrun struct {
	leftOver map[string]any
}

// string returns the json version of the string
func (o overrun) String() string {
	b, _ := json.Marshal(o.leftOver)
	return fmt.Sprintf(`%v`, string(b))
}

// MarhsalJson Exposes the hidden leftOver to be exposed
func (o overrun) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.leftOver)
}

// dataTranslate, takes a flat map compares each field with a mapping dictionary
// it then finds the destination fields that match and do a best guess for the full path based on the
// depths of the source and destination.
func dataTranslate(source map[string]any, destination map[string]map[int][]mdProperties, mappingPaths map[string][]string, missedFields string) map[string]any {

	// this is the output
	transformedData := make(map[string]any)
	keys := make([]string, len(source))

	i := 0
	for k := range source {
		keys[i] = k
		i++
	}
	slices.Sort(keys)
	// get the keys in alphabetical order
	// so that it can

	leftover := make(map[string]any)
	// @TODO break up the search so fields have multiple paths
	// then do some width and depth searches to have some best guess approximations
	for _, key := range keys {

		fullfield := key
		data := source[key]

		// get field from the full path name
		kfields := strings.Split(fullfield, ".")
		kfield := kfields[len(kfields)-1]

		_, err := strconv.Atoi(kfield)
		i := 2
		for err == nil {
			kfield = kfields[len(kfields)-i]
			_, err = strconv.Atoi(kfield)
			i++

			if i > len(kfields) {
				err = fmt.Errorf("checked all fields")
			}
		}

		depth := -1 // depth starts at 0 and any field will push the -1 to 0
		for _, field := range kfields {
			// integers do not count as part of the depth
			if _, err := strconv.Atoi(field); err != nil {
				depth++
			}
		}

		var props mdProperties
		//	this path is only used when the fields exactly match
		paths, ok := destination[kfield][depth]

		if !ok {
			// @TODO add a way for the user
			// to limit the depth searches
			for i := 1; i < 5; i++ {
				// check shallower
				paths, ok = destination[kfield][depth-i]
				if ok {
					break
				}
				// check deeper
				paths, ok = destination[kfield][depth+i]
				if ok {
					break
				}
				// alternate depths between less and more
			}
		}

		if ok {

			//	find the closest path that matches the field
			props = fieldPathMatch(fullfield, paths) // paths[0]
			//	fmt.Println(props, "PROPS")
		}

		for path, mp := range mappingPaths {

			// check if the field or the fullpath match
			if slices.Contains(mp, fullfield) || slices.Contains(mp, kfield) {
				// use that path
				paths, ok := destination[path][depth]
				if !ok {
					// @TODO add a way for the user
					// to limit the depth searches
					for i := 1; i < 10; i++ {
						// check shallower
						paths, ok = destination[path][depth-i]
						if ok {
							break
						}
						// check deeper
						paths, ok = destination[path][depth+i]
						if ok {
							break
						}
						// alternate depths between less and more
					}
				}

				// @TODO if paths have a length of more than 1 make a decision on which one
				if len(paths) == 1 {

					props = paths[0]
				} else {

					props = fieldPathMatch(fullfield, paths)
				}
			}
		}

		// if properties are found apply them
		if !reflect.DeepEqual(props, mdProperties{}) { // (props != mdProperties{}) {
			fullpath := props.fullpath

			// build the arrays
			if len(props.dimension) > 0 {
				//
				buildNDimensionalData(transformedData, data, props, 0)
			} else {
				// assign the single data point to the field
				dataAssign(data, props, transformedData, fullpath)
			}

		} else if data != nil {
			// save missed fields for later
			leftover[fullfield] = data
		}

	}

	// add the missed fields if required
	if missedFields != "" && len(leftover) != 0 {
		transformedData[missedFields] = overrun{leftOver: leftover}
	}

	return transformedData

}

// fieldPathMatch matches fields to the closest path
func fieldPathMatch(path string, choices []mdProperties) mdProperties {

	if len(choices) == 1 {
		return choices[0]
	}

	pieces := strings.Split(path, ".")
	filteredPath := ""
	for _, piece := range pieces {

		if _, err := strconv.Atoi(piece); err != nil {
			if len(filteredPath) == 0 {
				filteredPath += piece
			} else {
				filteredPath += "." + piece
			}
		} else {
			filteredPath += ".%v"
		}
	}
	// a quick version of this
	// https://en.wikipedia.org/wiki/Levenshtein_distance

	// make the score max so something always matches
	score := 0xfffffffffffffff
	var base mdProperties

	for _, choice := range choices {

		// skip values that are required placeholders
		// but if they are actual values don't ignore.
		// Json and xsd schemas parsed differently causing this
		if choice.required && choice.dataType == "" {
			continue
		}

		// @TODO rework this to go from back to front so the fields match
		/*
			// working version of going from field back
			pathBytes := []byte(filteredPath)
			for i := range []byte(filteredPath) {
				char := pathBytes[len(pathBytes)-1-i]
				if i+1 > len(choice.fullpath) {
					break
				}
				// @TODO update by having string lowercase and measuring the distance between them
				if choice.fullpath[i] != char {
					compareScore += 1
				}
			}
		*/
		compareScore := 0
		for i, char := range []byte(filteredPath) {
			if i+1 > len(choice.fullpath) {
				break
			}
			// @TODO update by having string lowercase and measuring the distance between them
			if choice.fullpath[i] != char {
				compareScore++
			}
		}

		compareScore += int(math.Abs(float64(len(filteredPath) - len(choice.fullpath))))

		if compareScore < score {
			base = choice
			score = compareScore
		}

	}

	return base
}

// buildNDimensionalData
func buildNDimensionalData(destination map[string]any, data any, props mdProperties, depth int) {
	// @ TODO preserve the origins data if neccessary, at the moment a new base path is generated for each type
	/* of data. Should look at sending the dimensions it had and comparing with the target data.
	get current dimensions of the data and compare
	*/
	switch da := data.(type) {

	// @TODO add recursiveness for multidimensional arrays
	// @TODO add more ways to traverse through the multidimensional objects
	case []any:

		// get length
		dataDimension := 1
		if len(da) > 1 {
			depth := true
			mid := da[0]
			for depth {
				switch m := mid.(type) {
				case []any:
					dataDimension++
					mid = m[0]
				default:
					depth = false
				}
			}
		}
		// only place the array if the target is an array
		if dataDimension == 1 && strings.Contains(props.dataType, "array") {
			dataAssign(da, props, destination, fmt.Sprintf(props.fullpath, props.dimension...))
		} else {
			for _, d := range da {

				fullpath := fmt.Sprintf(props.fullpath, props.dimension...)

				p := props.dimension[depth].(int)
				p++
				props.dimension[depth] = p
				switch d.(type) {
				case []any:

					buildNDimensionalData(destination, data, props, depth+1)
					// reset the depth after the array
				default:
					dataAssign(d, props, destination, fullpath)

				}
			}
		}
		// reset to 0 after movign along the array
		//	p := props.dimension[depth].(int)
		p := 0
		props.dimension[depth] = p
	default:
		// update the path name regardless
		fullpath := fmt.Sprintf(props.fullpath, props.dimension...)
		p := props.dimension[depth].(int)
		p++
		props.dimension[depth] = p
		dataAssign(data, props, destination, fullpath)
	}
}

// dataAssign handles the data converting it into the required data type of the destination
func dataAssign(data any, props mdProperties, transformedData map[string]any, fullpath string) {

	// find the data type and then handle using the go methods
	switch props.dataType {

	case String:
		transformedData[fullpath] = fmt.Sprintf("%v", data)
	case Integer:

		transformedData[fullpath] = AnyToInt(data)

	case Boolean:
		switch b := data.(type) {
		case bool:
			transformedData[fullpath] = b
		// case reflect.String:
		//	data.
		case string:
			if strings.ToLower(b) == "true" {
				transformedData[fullpath] = true
			} else if strings.ToLower(b) == "false" {
				transformedData[fullpath] = false
			}
		}
	case BooleanArray:

		dataArr := data.([]any)
		boolArr := make([]bool, len(dataArr))
		for i, d := range dataArr {
			if boo, ok := d.(bool); ok {

				boolArr[i] = boo
			}
		}

		transformedData[fullpath] = boolArr
	case Number:

		transformedData[fullpath] = AnyToFloat64(data)

	case StringArray:
		// arrays are []interface from json and yaml
		dataArr, ok := data.([]any)

		//	check if its an array or individual to add to the array
		out, outOk := transformedData[fullpath]
		if !outOk { // create a base
			out = []string{}
		}

		outArr := out.([]string)
		// if its an array append
		// else just overwrite
		if ok {
			stringArr := make([]string, len(dataArr))
			for i, d := range dataArr {
				stringArr[i] = fmt.Sprintf("%v", d)
			}

			outArr = append(outArr, stringArr...)
		} else {
			outArr = append(outArr, fmt.Sprintf("%v", data))
		}
		transformedData[fullpath] = outArr
	case IntegerArray:
		// arrays are []interface from json and yaml

		dataArr, ok := data.([]any)

		//	check if its an array or individual to add to the array
		out, outOk := transformedData[fullpath]
		if !outOk {
			out = []int{}
		}
		outArr := out.([]int)
		if ok {
			integerArr := make([]int, len(dataArr))
			for i, d := range dataArr {
				integerArr[i] = AnyToInt(d)
			}
			outArr = append(outArr, integerArr...)
		} else {
			outArr = append(outArr, AnyToInt(data))
		}
		transformedData[fullpath] = outArr
	case FloatArray:
		dataArr, ok := data.([]any)

		out, outOk := transformedData[fullpath]
		if !outOk {
			out = []float64{}
		}
		outArr := out.([]float64)

		if ok {
			floatArr := make([]float64, len(dataArr))
			for i, d := range dataArr {

				floatArr[i] = AnyToFloat64(d)

			}
			outArr = append(outArr, floatArr...)
		} else {
			outArr = append(outArr, AnyToFloat64(data))
		}

		transformedData[fullpath] = outArr

	case ObjectArray:
		//
		// @TODO check and append data
	case "any":
		transformedData[fullpath] = data
	default:

		// @TODO return an error missed and unhandled data types
		// or maybe a warning

	}
}

// XMLEncoderInformation is extracted from xsd and used for building
// the xml otuput
type xMLEncoderInformation struct {
	// path of parent and array of children
	attr     map[string]bool
	keyorder []string
	// @TODO add namespaces
	// namespaces      xml.Attr
	targetNameSpace string
	rootElement     string
}

// httpRead tries reading http before going
// local files.
func httpRead(input string) ([]byte, error) {
	resp, err := http.Get(input)
	if err == nil {

		return io.ReadAll(resp.Body)
	}

	return os.ReadFile(input)
}

// requiredPropertiesCheck checks all required properties are included,
// if not add an empty value.
func requiredPropertiesCheck(mdPaths map[string]map[int][]mdProperties, translatedData map[string]any) map[string]any {
	for _, paths := range mdPaths {
		for _, path := range paths {
			for _, p := range path {
				if p.required {
					// check bare dimensions
					fullpath := p.fullpath
					if len(p.dimension) != 0 {
						baredim := make([]any, len(p.dimension))
						for i := range baredim {
							baredim[i] = 0
						}
						fullpath = fmt.Sprintf(fullpath, baredim...)
					}

					_, ok := translatedData[fullpath]
					// @TODO fix this for JSON
					if !ok {
						//	fmt.Println(p.dataType, "dataType", p)
						translatedData[fullpath] = nil
					}
				}
			}
		}
	}

	return translatedData
}
