package mapping

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"slices"
	"strings"

	"aqwari.net/xml/xsd"
)

func jsonSchemaExtract(jsonSchema string) (map[string]map[int][]mdProperties, error) {
	b, err := httpRead(jsonSchema)
	if err != nil {
		return nil, fmt.Errorf("error extracting the schema for the output data: %v", err)
	}

	// @TODO move the schmea to a function to handle several schema types
	var schema map[string]any
	err = json.Unmarshal(b, &schema)

	if err != nil {
		return nil, err
	}

	// fmt.Println(schema)
	mdPaths := ExtractJsonMetaDatapaths(schema)
	return mdPaths, nil
}

// ExractMDPaths gets all the metadata paths from a json schema.
// These are flat dotpaths along with their corresponding type.
// Needs to have all array types implemented
func ExtractJsonMetaDatapaths(schema map[string]any) map[string]map[int][]mdProperties {
	s := schemaJson{WholeSchema: schema}
	// says any but should only be int, because of type constrains
	return s.validMDPaths(schema, "", make(map[string]map[int][]mdProperties), []any{}, 0)
}

// constants to make traversing the
// schema easier
const (
	schemaType = "type"
)

// schemaJson is the layout for the schmea to be decoded
type schemaJson struct {
	// The whole schema,
	// or at least the segments as it is being
	// recursed down
	WholeSchema map[string]any
	// The definitions of the schema
	// to be called for references
	Defs map[string]map[string]any
}

// validMDPaths gets all the metadata paths from a json schema.
// These are flat dotpaths along with their corresponding type.
// Needs to have arrays implemented
// path, depth then properties
func (s *schemaJson) validMDPaths(schema map[string]any, parent string, found map[string]map[int][]mdProperties, dimensions []any, depth int) map[string]map[int][]mdProperties {

	// arbitary depth to stop recursion
	if depth > 10 {
		return found
	}

	depth++
	// @TODO break this array into smaller more manageable functions
	for k, v := range schema {

		switch {
		case k == "properties":

			// recursively search the next layer to find the property values
			found = s.validMDPaths(v.(map[string]any), parent, found, dimensions, depth)
			// @TODO fix things matching schema names
		case k == "required":
			// store the required variables
			required := v.([]any)

			for _, r := range required {
				k := r.(string)
				depth := strings.Count(parent, ".") - len(dimensions)
				if _, ok := found[k]; !ok {
					found[k] = make(map[int][]mdProperties)
				}
				results := found[k][depth]

				results = append(results, mdProperties{fullpath: parent + k, required: true,
					dimension: slices.Clone(dimensions)})
				found[k][depth] = results
			}

		case k == "$ref":

			if path, ok := v.(string); ok {

				// extract the path then recurse through the child
				ref := s.ExtractReference(path)
				found = s.validMDPaths(ref, parent, found, dimensions, depth)

			}
			/*
				search the list of references

				if not build from the main schema

				which is search for the first path then make your own map to be outputted
				then attach it to the path

			*/
		case k == "$defs":
			// skip the defs
			// as these will be come back to
		case k == "items":

			// items are skipped to prevent them being saved as a field
		default:
			// extract the properties recursively from the child
			// such as array properties
			found = s.extractProperties(schema, parent, found, dimensions, depth, k, v)

		}
	}

	return found
}

// extractProperties sorts if the schema value is an array or object and processes
// accordingly
func (s *schemaJson) extractProperties(schema map[string]any, parent string, found map[string]map[int][]mdProperties, dimensions []any, depth int, field string, value any) map[string]map[int][]mdProperties {

	switch children := value.(type) {
	case map[string]any:

		// figure out the child  type here

		dataType, ok := children[schemaType]
		if !ok {
			if path, ok := children["$ref"]; ok {
				ref := s.ExtractReference(path.(string))
				if ref != nil {
					dataType = ref["type"].(string)
					children = ref
				}
			} else {
				dataType = "any"
			}
		}

		// find the depth of the field
		depth := strings.Count(parent, ".") - len(dimensions) // as dimensions add depths with the array space
		results := found[field][depth]
		if _, ok := found[field]; !ok {
			found[field] = make(map[int][]mdProperties)
		}

		if dataType != "array" {

			// if its an object just ignore as it has no finished path?

			results = append(results, mdProperties{fullpath: parent + field, dataType: fmt.Sprintf("%v", dataType),
				dimension: slices.Clone(dimensions)}) // clone the dimensions to not cross the streams

			// add the results and continue searching if there are more child entries
			found[field][depth] = results
			found = s.validMDPaths(children, parent+field+".", found, dimensions, depth)

		} else {
			// @TODO some errors about types not being found

			// find the items in the array
			ArrProperties, ok := children["items"]
			if ok {

				// are the properties and a map or an array?
				ArrPropertiesValid, ok := ArrProperties.(map[string]any)
				if !ok {
					// use an array of map as the fallback
					ArrPropertiesValidArr, _ := ArrProperties.([]any)

					if len(ArrPropertiesValidArr) != 0 {

						ArrPropertiesValid = ArrPropertiesValidArr[0].(map[string]any)
					}
				}

				// get the item type
				arrType := ArrPropertiesValid[schemaType]
				switch arrType {
				case "array":
					// find the items recursively
					switch arrProps := ArrPropertiesValid["items"].(type) {
					case []any:
						for _, arrProp := range arrProps {
							//@ TODO fix the immediate problem of
							// several items but only one key

							// in case of an array of arrays
							found = s.validMDPaths(arrProp.(map[string]any), parent+field+".%v.%v.", found, append(slices.Clone(dimensions), 0, 0), depth)
						}
					case map[string]any:

						found = s.validMDPaths(arrProps, parent+field+".%v.%v.", found, append(slices.Clone(dimensions), 0, 0), depth)
					default:

					}

				case "object":
					// keep the search moving
					found = s.validMDPaths(ArrPropertiesValid, parent+field+".%v.", found, append(slices.Clone(dimensions), 0), depth)

				default:

					// else its a field
					// add to the array and keep on searching
					props := mdProperties{fullpath: parent + field, dataType: fmt.Sprintf("%v%v", arrType, dataType)}

					if len(dimensions) != 0 {
						props.dimension = slices.Clone(dimensions)
					}
					results = append(results, props)
					found[field][depth] = results

					found = s.validMDPaths(children, parent+field+".%v.", found, append(slices.Clone(dimensions), 0), depth)
				}

			}

			// continue the search if there are more child entries
			// found = validMDPaths(children, parent+k+".", found)
		}
	// for when it is a type not an object just the type
	case string:

		// skip objects and arrays
		if children != "object" && children != "array" && parent != "" {
			fields := strings.Split(parent, ".")

			// find the field for arrays of values
			// as the field will not be the most
			// recent name in the list
			for i := len(fields) - 2; i >= 0; i-- {

				if fields[i] != `%v` {
					field = fields[i]
					break
				}
			}

			depth := strings.Count(parent, ".") - len(dimensions) // as dimensions add depths with the array space
			results := found[field][depth]

			results = append(results, mdProperties{fullpath: parent[:len(parent)-1], dataType: fmt.Sprintf("%v", children),
				dimension: slices.Clone(dimensions)}) // clone the dimensions to not cross the streams

			if _, ok := found[field]; !ok {
				found[field] = make(map[int][]mdProperties)
			}
			// add the results and continue searching if there are more child entries
			found[field][depth] = results
		}

	}

	return found
}

func xmlSchemaExtract(xsdSchema string) (map[string]map[int][]mdProperties, *xMLEncoderInformation, error) {
	bs, err := httpRead(xsdSchema)
	if err != nil {
		return nil, nil, err
	}
	sc, err := xsd.Parse(bs)

	if err != nil {
		return nil, nil, err
	}

	tree, err := xsd.Normalize(bs)

	if err != nil {
		return nil, nil, err
	}

	elements := []string{}
	for _, tr := range tree {
		child := tr.Children
		for _, c := range child {
			//fmt.Println(c.Copy().Name.Local, "COPYER")
			if c.Copy().Name.Local == "element" {

				properties := make(map[string]string)
				for _, p := range c.StartElement.Attr {
					properties[p.Name.Local] = p.Value
				}

				// make sure its a standalone type for the moment
				// and not just an array or something
				if properties["type"] == "" || !strings.Contains(properties["type"], ":") {
					elements = append(elements, properties["name"])
				}
			}
		}
	}

	if len(elements) == 0 {
		return nil, nil, fmt.Errorf("no root elements found in the schema, could no build the metadata object")
	}

	// @TODO loop through each element
	TargetSpace := sc[0].TargetNS

	xmlProperties := xMLEncoderInformation{attr: make(map[string]bool), keyorder: []string{elements[0] + ".xmlns"}}

	root := sc[0].FindType(xml.Name{Space: TargetSpace, Local: elements[0]}).(*xsd.ComplexType)

	mdPaths := make(map[string]map[int][]mdProperties)
	mdPaths, _ = xmlProperties.xmlDataPaths(sc[0], root, elements[0]+".", elements[0]+".", mdPaths, []any{})

	xmlProperties.attr[elements[0]+".xmlns"] = true
	xmlProperties.targetNameSpace = TargetSpace
	xmlProperties.rootElement = elements[0]

	return mdPaths, &xmlProperties, nil
}

func (x *xMLEncoderInformation) xmlDataPaths(parentSchema xsd.Schema, root xsd.Type, parentPath, parent string, properties map[string]map[int][]mdProperties, dimensions []any) (map[string]map[int][]mdProperties, []string) {
	// fmt.Println(reflect.TypeOf(root))

	depth := strings.Count(parentPath, ".")

	output := []string{}
	switch field := root.(type) {
	case *xsd.ComplexType:

		attr := field.Attributes
		for _, a := range attr {

			// attributes aren't going to be complex types
			// they are always simple
			output = append(output, parentPath+a.Name.Local)
			x.keyorder = append(x.keyorder, parentPath+a.Name.Local)
			if _, ok := properties[a.Name.Local]; !ok {
				properties[a.Name.Local] = make(map[int][]mdProperties)
			}
			// get the properties
			props := properties[a.Name.Local][depth]
			props = append(props, mdProperties{
				dimension: slices.Clone(dimensions), fullpath: parent + a.Name.Local,
				required: !a.Optional, dataType: xmlTypeStringSwitcher(fmt.Sprintf("%v", a.Type), a.Plural),
			})
			properties[a.Name.Local][depth] = props

			x.attr[parentPath+a.Name.Local] = true

		}

		fields := field.Elements
		for _, f := range fields {

			found := parentSchema.FindType(f.Name)

			if found != nil {
				var outputMid []string
				// search again at the next layer
				properties, outputMid = x.xmlDataPaths(parentSchema, found, parentPath+f.Name.Local+".", parent+f.Name.Local+".%v.", properties, append(slices.Clone(dimensions), 0))
				output = append(output, outputMid...)
			} else {

				if _, ok := properties[f.Name.Local]; !ok {
					properties[f.Name.Local] = make(map[int][]mdProperties)
				}

				props := properties[f.Name.Local][depth]
				props = append(props, mdProperties{
					dimension: slices.Clone(dimensions), fullpath: parent + f.Name.Local,
					required: !f.Optional, dataType: xmlTypeStringSwitcher(fmt.Sprintf("%v", f.Type), f.Plural),
				})

				properties[f.Name.Local][depth] = props
				x.keyorder = append(x.keyorder, parentPath+f.Name.Local)
				output = append(output, parentPath+f.Name.Local)
			}
		}

	case *xsd.SimpleType:

		x.keyorder = append(x.keyorder, parentPath+field.Name.Local)
		output = append(output, parentPath+field.Name.Local)

	}

	return properties, output
}

// ExtractReference gets the object at the end of a reference path
func (s *schemaJson) ExtractReference(path string) map[string]any {

	// @TODO include http extraction
	// return a complete path and update s if at a new location
	if path[0] == '#' {

		paths := strings.Split(path, "/")
		//find the end map
		dest := s.WholeSchema
		// recursively search the map
		// @TODO add error checking to prevent
		// paths to nowhere
		for _, p := range paths[1:] {

			dest = dest[p].(map[string]any)

		}
		return dest
	}
	return nil
}

func xmlTypeStringSwitcher(xmlString string, Array bool) string {
	base := ""
	switch xmlString {
	case "String":
		base = String
	case "Decimal", "Number":
		base = Number
	case "Integer":
		base = Integer
	case "boolean":
		base = Boolean
	default:
		return "any"
	}

	if Array {
		return base + "array"
	}

	return base
}
