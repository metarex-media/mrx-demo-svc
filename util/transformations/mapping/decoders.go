// Package mapping is for best guess data transforms, where
// some sort of mapping has been provided to help fill in the blanks.
package mapping

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// jsonDecode returns the data as a flat map
// where a dot path of the name is generated e.g. path.to.data
func jsonDecode(input []byte, timing string) ([]map[string]any, error) {

	// if the json is a self timed array
	if timing == "embedded" {

		var jsonLayout []map[string]any
		err := json.Unmarshal(input, &jsonLayout)

		if err != nil {
			return nil, err
		}
		// do some recusrive searching
		for i, m := range jsonLayout {
			jsonLayout[i] = mapExtract(m, make(map[string]any), "")
		}

		return jsonLayout, nil

	}

	var jsonLayout map[string]any
	err := json.Unmarshal(input, &jsonLayout)

	if err != nil {
		return nil, err
	}
	// do some recursive searching
	flatMap := mapExtract(jsonLayout, make(map[string]any), "")

	return []map[string]any{flatMap}, nil

}

// yamlDecode returns the data as a flat map
func yamlDecode(input []byte, timing string) ([]map[string]any, error) {

	if timing == "embedded" {

		var yamlLayout []map[string]any
		err := yaml.Unmarshal(input, &yamlLayout)

		if err != nil {
			return nil, err
		}
		// do some recusrive searching
		for i, m := range yamlLayout {
			yamlLayout[i] = mapExtract(m, make(map[string]any), "")
		}

		return yamlLayout, nil

	}

	var yamlLayout map[string]any
	err := yaml.Unmarshal(input, &yamlLayout)

	if err != nil {
		return nil, err
	}
	// fmt.Println(yamlLayout)
	// do some recursive searching
	flatMap := mapExtract(yamlLayout, make(map[string]any), "")
	// fmt.Println(flatMap)
	return []map[string]any{flatMap}, nil

}

// mapExtract recursively searches the data for all non object and array values
func mapExtract(inputLayout, found map[string]any, parent string) map[string]any {

	for k, v := range inputLayout {

		// fmt.Println(v, reflect.TypeOf(v))
		switch dest := v.(type) {
		case map[string]any:
			found = mapExtract(dest, found, parent+k+".")
		case []any:
			// if its really []any

			if len(dest) != 0 {
				switch dest[0].(type) {

				case map[string]any:
					// if it map[string]any then assume everything in the array is for the moment
					for i, d := range dest {

						found = mapExtract(d.(map[string]any), found, fmt.Sprintf("%v.%v.", parent+k, i))
					}
				case []any:

					found = arrayExtract(dest, found, parent+k)
				default: // else assume []int or the likes

					found[parent+k] = v
				}
			}
		case []float64:
		//	fmt.Println(k, v)
		default:

			found[parent+k] = v
		}
	}

	return found
}

// ArrayExtract recursively searches through arrays of type any to extract the values
// and their relative positions.
// expects parent without the dot path at the end e.g. full.path not full.path.
func arrayExtract(inputArray []any, found map[string]any, parent string) map[string]any {

	if len(inputArray) != 0 {

		switch inputArray[0].(type) {

		case map[string]any:
			// if it map[string]any then assume everything in the array is for the moment
			for i, d := range inputArray {
				found = mapExtract(d.(map[string]any), found, fmt.Sprintf("%v.%v.", parent, i))
			}
		case []any:

			for i, d := range inputArray {

				// recurse into the array
				found = arrayExtract(d.([]any), found, fmt.Sprintf("%v.%v", parent, i))
			}
			// we go again recurively
		default: // else assume []int or the likes

			for i, d := range inputArray {
				pos := fmt.Sprintf("%v.%v", parent, i)
				found[pos] = d
			}
			//	fmt.Println(parent, inputArray, reflect.TypeOf(inputArray), "default")
			//	found[parent] = inputArray
		}
	}

	return found
}

// csvDecode extracts each layer of csv as a new map
func csvDecode(input []byte) ([]map[string]any, error) {
	r := csv.NewReader(bytes.NewReader(input))
	// get headers
	headers, err := r.Read()
	if err != nil {
		return nil, err
	}

	bases := make([]map[string]any, 0)

	body, err := r.Read()
	if err != nil {
		return nil, err
	}

	for err == nil {
		base := make(map[string]any)

		for i, b := range body {
			if b != "" { // slip nil values
				base[headers[i]] = b
			}
		}

		bases = append(bases, base)
		body, err = r.Read()

		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	return bases, nil
}
