package mapping

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*
func TestCamera(t *testing.T) {

	in, _ := os.ReadFile("/workspace/gl-mrx-demo-svc/demodata/demo06/fish_camera.json")

	mapDefs := Mapping{MappingDefinitions: map[string][]string{
		"intrinsicMatrix": {"intrinsics"},
		"cameraRotation":  {"rotation"}}}
	t2 := MappingAction{Mapping: mapDefs,
		InputFormat: "application/json", OutputFormat: "application/xml",
		OutputSchema: "/workspace/gl-mrx-demo-svc/demodata/schema/demo06output.xsd",
	}

	o, e := t2.Transform(in, nil)
	fmt.Println(e)
	fmt.Println(string(o[0]))
} */

var mapDefs = Mapping{MappingDefinitions: map[string][]string{
	"string":        {"string"},
	"integer":       {"integer"},
	"float":         {"float"},
	"bool":          {"bool"},
	"stringArray":   {"stringArray"},
	"integerArray":  {"integerArray"},
	"boolArray":     {"boolArray"},
	"floatArray":    {"floatArray"},
	"field1":        {"objectField1Arr"},
	"field2":        {"objectField2Arr"},
	"Reversefield1": {"Reversefield1"},
	"Reversefield2": {"Reversefield2"},
	"float2d":       {"float2d"},
	"float1d":       {"float1d"},
	"floats2d":      {"floats2d"},
},
	MissedTags: "metadataTag",
}

func TestBaseTypes(t *testing.T) {

	/*

		Test layout

		compress into a mapping function that spits our the errors etc

		have a master parent schmea

		design changes with clear outcomes:

		- Arrays are preserved
		- objects to arrays
		- arrays to objects
		- do a change of all types

		test required fields etc
	*/

	/*

		have a type in
		type out with matching jsons etc

	*/

	inTypes := []string{"strings", "integer", "float", "bool",
		"stringArray", "floatArray", "integerArray",
		"all"}

	t2 := MappingAction{Mapping: mapDefs,
		InputFormat: "application/json", OutputFormat: "application/json",
		OutputSchema: "./testdata2/basetypes/baseTypesSchema.json",
	}

	for _, in := range inTypes {
		inBytes, errFile := os.ReadFile(fmt.Sprintf("./testdata2/basetypes/%vin.json", in))

		outBytes, err := t2.Transform(inBytes, nil)

		fbase, errExpec := os.ReadFile(fmt.Sprintf("./testdata2/basetypes/%vout.json", in))

		got := sha256.New()
		got.Write(outBytes)
		expected := sha256.New()
		expected.Write(fbase)

		// to generate a file to compare outputs when things go wrong
		/*if fmt.Sprintf("%032x", got.Sum(nil)) != fmt.Sprintf("%032x", expected.Sum(nil)) {

			fbaseo, _ := os.Create(fmt.Sprintf("./testdata2/basetypes/%vout.json", in))
			fbaseo.Write(outBytes)
		}*/

		//
		Convey(fmt.Sprintf("Checking  types of %s are correctly handled", in), t, func() {
			Convey("the type is mapped without changing", func() {
				Convey("The hash matches the expected output", func() {
					So(errFile, ShouldBeNil)
					So(errExpec, ShouldBeNil)
					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
				})
			})
		})
	}

	/*
				two arrays to objects and vice versa
				two d arrays - no need to validate every type

				2d arrays to 2d arrays @TODO fix these
				2 d array to 1d array

		$
	*/
	t2.OutputSchema = "./testdata2/objects/objectSchema.json"
	objects := []string{"arraysToObj", "objToArrays", "float2d", "float2ds"} //, "float2dTo1d"}

	for _, in := range objects {
		inBytes, errFile := os.ReadFile(fmt.Sprintf("./testdata2/objects/%v.json", in))

		outBytes, err := t2.Transform(inBytes, nil)

		fbase, errExpec := os.ReadFile(fmt.Sprintf("./testdata2/objects/%vout.json", in))
		fmt.Println(err)
		got := sha256.New()
		got.Write(outBytes)
		expected := sha256.New()
		expected.Write(fbase)

		// to generate a file to compare outputs when things go wrong
		if fmt.Sprintf("%032x", got.Sum(nil)) != fmt.Sprintf("%032x", expected.Sum(nil)) {

			fbaseo, _ := os.Create(fmt.Sprintf("./testdata2/objects/%vout.json", in))
			fbaseo.Write(outBytes)
		}

		//
		Convey(fmt.Sprintf("Checking objects that %s are correctly handled", in), t, func() {
			Convey("the type is mapped as intended", func() {
				Convey("The hash matches the expected output", func() {
					So(errFile, ShouldBeNil)
					So(errExpec, ShouldBeNil)
					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
				})
			})
		})
	}

	/*

		in, _ := os.ReadFile("./testdata/mrx.csv")

		mapDefs := Mapping{MappingDefinitions: map[string][]string{
			"string":      {"string"},
			"out":     {"out"},
			"segment": {"segment", "slice", "piece"},
			"source":  {"source"}},
			MissedTags: "metadataTag",
		}
		t2 := MappingAction{Mapping: mapDefs,
			InputFormat: "text/csv", OutputFormat: "text/csv",
			OutputSchema: "./testdata/csvschema.json",
		}

		o, _ := t2.Transform([][]byte{in}, nil)
		fout, _ := os.Create("./testdata/out.csv")
		fout.Write(o[0])

		fmt.Println(string(o[0]))

		t2.OutputFormat = "application/json"

		ojson, _ := t2.Transform([][]byte{in}, nil)
		foutjson, _ := os.Create("./testdata/out.json")
		foutjson.Write(ojson[0])

		/*
			/localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnc&outputMRXID=MRX.123.456.789.rnf&mapping=true
					mapping behvaiours to consider:
					- variable inputs like rnf
					- multiple fields having the same name
					- mapping integers to strings etc */
	/*
		mpats := `{"mappingDefinitions": {
						"chapter": [
							"chapter",
							"Chapter"
						],
						"in": [
							"in",
							"In",
							"in(f)"
						],
						"out": [
							"out",
							"Out",
							"out(f)"
						],
						"storyline-importance": [
							"storyline-importance",
							"Storyline-importance",
							"Importance",
							"Story"
						]
					}}`
		var mpatsMap Mapping
		json.Unmarshal([]byte(mpats), &mpatsMap)
		mpatsMap.MissedTags = "metadataTags"
		t2o := MappingAction{Mapping: mpatsMap,
			InputFormat: "text/csv", OutputFormat: "text/csv",
			OutputSchema: "/workspace/gl-mrx-demo-svc/api/rnf/rnfSchema.json",
		}
		inPast, _ := os.ReadFile("/workspace/gl-mrx-demo-svc/demodata/demo04/lostpast.csv")
		opast, _ := t2o.Transform([][]byte{inPast}, nil)
		//fout, _ := os.Create("./testdata/out.csv")
		//fout.Write(o[0])
		fmt.Println(string(opast[0]))*/
}

func TestOutput(t *testing.T) {

	t2 := MappingAction{Mapping: mapDefs,
		InputFormat: "application/json", OutputFormat: "application/json",
		OutputSchema: "./testdata2/basetypes/baseTypesSchema.json",
	}

	types := []string{"json", "yaml", "csv", "xml"}

	for _, ty := range types {
		inBytes, errFile := os.ReadFile("./testdata2/basetypes/allin.json")
		if ty != "csv" {
			t2.OutputFormat = fmt.Sprintf("application/%v", ty)
		} else {
			t2.OutputFormat = "text/csv"
		}

		if ty == "xml" {
			t2.OutputSchema = "./testdata2/basetypes/baseTypesSchema.xsd"
		}
		outBytes, err := t2.Transform(inBytes, nil)

		fbase, errExpec := os.ReadFile(fmt.Sprintf("./testdata2/basetypes/allin.%v", ty))
		fmt.Println(err)
		got := sha256.New()
		got.Write(outBytes)
		expected := sha256.New()
		expected.Write(fbase)

		// to generate a file to compare outputs when things go wrong
		/*if fmt.Sprintf("%032x", got.Sum(nil)) != fmt.Sprintf("%032x", expected.Sum(nil)) {

			fbaseo, _ := os.Create(fmt.Sprintf("./testdata2/basetypes/allin.%v", ty))
			fbaseo.Write(outBytes)
		}*/

		//
		Convey(fmt.Sprintf("Checking that %s formats are correctly written to", ty), t, func() {
			Convey("the file type is written", func() {
				Convey("The hash matches the expected output", func() {
					So(errFile, ShouldBeNil)
					So(errExpec, ShouldBeNil)
					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
				})
			})
		})
	}
	// test
}

func TestExtraction(t *testing.T) {

	t2 := MappingAction{Mapping: mapDefs,
		InputFormat: "application/json", OutputFormat: "application/json",
		OutputSchema: "./testdata2/basetypes/baseTypesSchema.json",
	}

	types := []string{"json", "yaml"}

	for _, ty := range types {
		inBytes, errFile := os.ReadFile(fmt.Sprintf("./testdata2/basetypes/allin.%v", ty))

		t2.InputFormat = fmt.Sprintf("application/%v", ty)

		outBytes, err := t2.Transform(inBytes, nil)

		fbase, errExpec := os.ReadFile("./testdata2/basetypes/allin.json")
		fmt.Println(err)
		got := sha256.New()
		got.Write(outBytes)
		expected := sha256.New()
		expected.Write(fbase)

		//
		Convey(fmt.Sprintf("Checking %s files can be extracted from", ty), t, func() {
			Convey("the file type is written", func() {
				Convey("The hash matches the expected output", func() {
					So(errFile, ShouldBeNil)
					So(errExpec, ShouldBeNil)
					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
				})
			})
		})
	}

	// embedded timing

	t2.InputTiming = "embedded"
	for _, ty := range types {
		inBytes, errFile := os.ReadFile(fmt.Sprintf("./testdata2/basetypes/allinembedded.%v", ty))

		t2.InputFormat = fmt.Sprintf("application/%v", ty)

		outBytes, err := t2.Transform(inBytes, nil)

		fbase, errExpec := os.ReadFile("./testdata2/basetypes/allinembedded.json")

		got := sha256.New()
		got.Write(outBytes)
		expected := sha256.New()
		expected.Write(fbase)

		//
		Convey(fmt.Sprintf("Checking %s files with embedded timing can be extracted from", ty), t, func() {
			Convey("the file type is written", func() {
				Convey("The hash matches the expected output", func() {
					So(errFile, ShouldBeNil)
					So(errExpec, ShouldBeNil)
					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
				})
			})
		})
	}
	// test

	inBytes, errFile := os.ReadFile("./testdata2/basetypes/simple.csv")

	t2.InputTiming = ""
	t2.InputFormat = "text/csv"

	outBytes, err := t2.Transform(inBytes, nil)

	fbase, errExpec := os.ReadFile("./testdata2/basetypes/simpleout.json")
	fmt.Println(err)
	got := sha256.New()
	got.Write(outBytes)
	expected := sha256.New()
	expected.Write(fbase)

	//
	Convey("Checking simple csv files can be extracted from", t, func() {
		Convey("the file type is written", func() {
			Convey("The hash matches the expected output", func() {
				So(errFile, ShouldBeNil)
				So(errExpec, ShouldBeNil)
				So(err, ShouldBeNil)
				So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
			})
		})
	})
}

func TestSchemaReading(t *testing.T) {

	t2 := MappingAction{Mapping: mapDefs,
		InputFormat: "application/json", OutputFormat: "application/json",
		OutputSchema: "./testdata2/basetypes/baseTypesSchema.json",
	}

	schemas := []string{"baseTypesSchema.json", "baseTypeSchemaRef.json", "baseTypesSchemaReq.json", "baseTypesSchema.xsd"}
	dest := "./testdata2/basetypes/allin.json"

	for _, schema := range schemas {
		inBytes, errFile := os.ReadFile("./testdata2/basetypes/allin.json")

		t2.OutputSchema = fmt.Sprintf("./testdata2/basetypes/%v", schema)

		outBytes, err := t2.Transform(inBytes, nil)

		if schema == "baseTypesSchema.xsd" {
			dest = "./testdata2/basetypes/allinxsd.json"
		}

		fbase, errExpec := os.ReadFile(dest)

		got := sha256.New()
		got.Write(outBytes)
		expected := sha256.New()
		expected.Write(fbase)

		//
		Convey("Checking the different types of schema properties can be utilised", t, func() {
			Convey(fmt.Sprintf("the file type is written using a schema of %v", schema), func() {
				Convey("The hash matches the expected output", func() {
					So(errFile, ShouldBeNil)
					So(errExpec, ShouldBeNil)
					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
				})
			})
		})
	}
}

/*
func TestAddress(t *testing.T) {
	fmt.Println("start")
	in, _ := os.ReadFile("./testdata/address.json")
	fmt.Println("starting")
	mapDefs := Mapping{MappingDefinitions: map[string][]string{
		"Name": {"firstName"},
		"role": {"primary job"},
		"year": {"year"},
	}}
	t2 := MappingAction{Mapping: mapDefs,
		InputFormat: "application/json", OutputFormat: "application/json",
		OutputSchema: "./testdata/adressSchema.json",
	}

	o, err := t2.Transform(in, nil)
	fmt.Println(err, len(o))
	fout, _ := os.Create("./testdata/role.json")
	fout.Write(o)
}

func TestVariables(t *testing.T) {

	var schemaDataOut2 = []byte(`{
		"$id": "https://qri.io/schema/",
		"$comment" : "sample comment",
		"title": "Person",
		"type": "object",
		"properties": {
			"valid": {
				"type" : "boolean"
			},
			"Name": {
				"type": "string"
			},
			"year": {
				"type": "integer"
			},
			"numberCheck": {
				"type": "number"
			},
			"people": {
				"type": "array",
				"items": {
					"type": "object" ,
					"properties": {
						"age": {
							"type": "integer"
						},
						"Name": {
							"type": "string"
						},
						"roles" : {
							"type": "array",
							"items": {
								"type": "object" ,
								"properties": {
									"task": {
										"type": "string"
									}
								}
							}
						}
					}
				}
			},
			"action" :{
				"type": "object",
				"properties": {
					"year": {
						"type": "integer"
					},
					"role": {
						"type": "string"
					}
				}
			}
		},
		"required": ["firstName", "lastName"]
	  }`)

	input2 := `{
		"firstName": "Tristan",
		"lastName": "Larke",
		"primary job": "video systems architect",
		"secondary job": "whatever bruce delegates",
		"year":2024,
		"valid": true,
		"check": 5.43,
		"in" : {"age": [1, 2], "firstName":  "Tristan", "name" : ["bruce", "bruce2"]},
		"tasks":["in", "out", "check","cheker"],
		"jobs" : {
				"primary job": "lighting assistant",
				"year":2020
		}
	  }`
	in2, _ := jsonBaseExtract([]byte(input2), "clocked")
	var schema2 map[string]any

	/*
		needs to be fields with paths.
		Then match input fields with the closest paths possible
		so field, paths []string doesn't seem that useful, depth? build unique paths that
		can't be easily searched.

		Have field[]properties- depth, full path etc.

*/
/*
	json.Unmarshal([]byte(schemaDataOut2), &schema2)

	props2 := ExtractMDpaths(schema2)
	mapDefs2 := Mapping{MappingDefinitions: map[string][]string{
		"Name":        {"firstName", "name"},
		"role":        {"primary job"},
		"year":        {"year"},
		"test":        {"test"},
		"age":         {"age"},
		"task":        {"tasks"},
		"valid":       {"valid"},
		"numberCheck": {"check"},
	}}
	fmt.Println(props2, "PATHS")
	made2 := dataTranslate(in2[0], props2, mapDefs2.MappingDefinitions, "miss")

	b2, e2 := jsonBuild([]map[string]any{made2})
	fmt.Println(string(b2), e2, made2)

	fout, _ := os.Create("./testdata/role2.json")
	fout.Write(b2)

	// develop a test to check each input with an expected input struct
}

func TestVariables2(t *testing.T) {
	in, _ := os.ReadFile("./testdata/adress2.json")

	mapDefs := Mapping{MappingDefinitions: map[string][]string{
		"stringCheck":       {"stringCheck"},
		"intCheck":          {"intCheck"},
		"floatCheck":        {"floatCheck"},
		"boolCheck":         {"boolCheck"},
		"numberCheck":       {"check"},
		"stringArrayCheck":  {"stringArrayCheck"},
		"integerArrayCheck": {"integerArrayCheck"},
		"floatArrayCheck":   {"floatArrayCheck"},
		"boolArrayCheck":    {"boolArrayCheck"},
		"field1":            {"obejctField1Arr"},
		"field2":            {"obejctField2Arr"},
		"Reversefield1":     {"Reversefield1"},
		"Reversefield2":     {"Reversefield2"},
		"integerArrays":     {"integerArrays"},
		"float2d":           {"float2d"},
	},
		MissedTags: "missed"}

	applications := []string{"json", "yaml"}

	for _, ext := range applications {
		t2 := MappingAction{Mapping: mapDefs,
			InputFormat: fmt.Sprintf("application/%v", ext), OutputFormat: fmt.Sprintf("application/%v", ext),
			OutputSchema: "./testdata/adressSchema2.json",
		}

		o, err := t2.Transform(in, nil)

		fmt.Println(err, len(o))
		//	fout, _ := os.Create("./testdata/role2.json")
		//	fout.Write(o[0])
		type objectTest struct {
			Field1 string `json:"field1,omitempty"`
			Field2 int    `json:"field2,omitempty"`
		}

		type testOutput struct {
			StringCheck       string       `json:"stringCheck" yaml:"stringCheck" `
			BoolCheck         bool         `json:"boolCheck" yaml:"boolCheck"`
			IntCheck          int          `json:"intCheck" yaml:"intCheck"`
			FloatCheck        float64      `json:"floatCheck" yaml:"floatCheck"`
			Year              int          `json:"year" yaml:"year"`
			StringArrayCheck  []string     `json:"stringArrayCheck" yaml:"stringArrayCheck"`
			IntegerArrayCheck []int        `json:"integerArrayCheck" yaml:"integerArrayCheck"`
			FloatArrayCheck   []float64    `json:"floatArrayCheck" yaml:"floatArrayCheck"`
			BoolArrayCheck    []bool       `json:"boolArrayCheck" yaml:"boolArrayCheck"`
			ObjectArrayCheck  []objectTest `json:"objectArrayCheck" yaml:"objectArrayCheck"`
			ReverseField1     []string     `json:"Reversefield1,omitempty" yaml:"Reversefield1,omitempty"`
			ReverseField2     []int        `json:"Reversefield2,omitempty" yaml:"Reversefield2,omitempty"`
			IntegerArrays     []int        `json:"integerArrays" yaml:"integerArrays"`
			Float2d           []float64    `json:"float2d" yaml:"float2d"`
		}

		expectedOutput := testOutput{
			BoolCheck:         true,
			StringCheck:       "john",
			StringArrayCheck:  []string{"one", "two", "three", "four"},
			IntegerArrayCheck: []int{1, 2, 3, 4},
			FloatArrayCheck:   []float64{1.5, 2.5, 3.5, 4.5},
			BoolArrayCheck:    []bool{true, false, true, false},
			IntCheck:          100,
			FloatCheck:        50.5,
			ObjectArrayCheck:  []objectTest{{"Matthew", 2}, {"Mark", 7}, {"Luke", 9}, {"John", 4}},
			ReverseField1:     []string{"Matthew", "Mark", "Luke", "John"},
			ReverseField2:     []int{2, 7, 9, 4},
			IntegerArrays:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}

		var generatedOutput testOutput
		yaml.Unmarshal(o, &generatedOutput)
		fout, _ := os.Create(fmt.Sprintf("./testdata/role2.%v", ext))
		fout.Write(o)
		Convey("Checking that the noise is generated", t, func() {
			Convey(fmt.Sprintf("Comparing the generated image to %v ", "hi"), func() {
				Convey("No error is returned and the file matches exactly", func() {
					So(err, ShouldBeNil)

				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed boolean types", func() {
				Convey("The valid field of type boolean is successfully transformed", func() {
					So(generatedOutput.BoolCheck, ShouldResemble, expectedOutput.BoolCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type integer is successfully transformed", func() {
					So(generatedOutput.IntCheck, ShouldResemble, expectedOutput.IntCheck)
				})
			})
		})
		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type integer is successfully transformed", func() {
					So(generatedOutput.StringCheck, ShouldResemble, expectedOutput.StringCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.StringArrayCheck, ShouldResemble, expectedOutput.StringArrayCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.IntegerArrayCheck, ShouldResemble, expectedOutput.IntegerArrayCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.FloatArrayCheck, ShouldResemble, expectedOutput.FloatArrayCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.FloatCheck, ShouldResemble, expectedOutput.FloatCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.BoolArrayCheck, ShouldResemble, expectedOutput.BoolArrayCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.ObjectArrayCheck, ShouldResemble, expectedOutput.ObjectArrayCheck)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.ReverseField1, ShouldResemble, expectedOutput.ReverseField1)
					So(generatedOutput.ReverseField2, ShouldResemble, expectedOutput.ReverseField2)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.IntegerArrays, ShouldResemble, expectedOutput.IntegerArrays)
				})
			})
		})

		Convey("Checking that the output fields match the intended value", t, func() {
			Convey("Comparing the transformed output to expected", func() {
				Convey("The year field of type stringArray is successfully transformed", func() {
					So(generatedOutput.Float2d, ShouldResemble, expectedOutput.Float2d)
				})
			})
		})
	}
	x := 8.9
	fmt.Println(int(x))
}

/*

// name extract gets
func nameExtract(schema map[string]any, parent string, found map[string]mdProperties) map[string]mdProperties {
	for k, v := range schema {
		switch {
		case k == "properties":
			found = nameExtract(v.(map[string]any), parent, found)
		case !slices.Contains(schemaKey, k):

			if children, ok := v.(map[string]any); ok {
				mdType, ok := children["type"]
				if !ok {
					mdType = "any"
				}

				// if its an object just ignore as it has no finished path?
				found[parent+k] = mdProperties{fullpath: parent + k, mdType: fmt.Sprintf("%v", mdType)}
				found = nameExtract(children, parent+k+".", found)
			}
		}
	}

	return found
}
*/
/*
func TestLPX(t *testing.T) {
	fmt.Println("start")
	in, _ := os.ReadFile("./testdata/lpx.json")
	fmt.Println("starting")
	mapDefs := Mapping{
		MissedTags: "Missed",
		MappingDefinitions: map[string][]string{
			"uri":             {"uri"},
			"contenttype":     {"contenttype"},
			"name":            {"name", "literal"},
			"start":           {"startdate"},
			"end":             {"enddate"},
			"versionCreated":  {"versioncreated"},
			"contentCreated":  {"contentcreated"},
			"note":            {"ednote"},
			"lang":            {"language"},
			"standardversion": {"version"},
		}}
	t2 := MappingAction{Mapping: mapDefs,
		InputFormat: "application/json", OutputFormat: "application/json",
		OutputSchema: "./testdata/lpxschema.json",
	}

	o, err := t2.Transform(in, nil)
	fmt.Println(err)
	fmt.Println(err, len(o))
	fout, e := os.Create("./testdata/lpxout.json")
	fmt.Println(e, err)
	fout.Write(o)

	insmall, _ := os.ReadFile("./testdata/lpxsmall.json")
	t2.OutputSchema = "./testdata/lpxsmallschema.json"
	osmall, err := t2.Transform(insmall, nil)
	fmt.Println(err)
	fmt.Println(err, len(osmall))
	foutsmall, e := os.Create("./testdata/lpxsmallout.json")
	fmt.Println(e, err)
	foutsmall.Write(osmall)

	inninjs, _ := os.ReadFile("./testdata/ninjs-demo.json")
	t2.OutputSchema = "./testdata/ninjsSchema.json"
	osninjs, err := t2.Transform(inninjs, nil)
	fmt.Println(err)
	fmt.Println(err, len(osninjs))
	foutninjs, e := os.Create("./testdata/ninjsout.json")
	fmt.Println(e, err)
	foutninjs.Write(osninjs)
	/*
	   // might be worth using his method and seeing how it is handled
	   b, _ := os.Open("./testdata/newsitem.xml")
	   decoder := xml2map.NewDecoderWithPrefix(b, "", "text")
	   result, err := decoder.Decode()
	   fmt.Println(result, err)
	   be, _ := json.MarshalIndent(result, "", "    ")
	   fmt.Println(string(be))
*/
/*
	innewsml, _ := os.ReadFile("./testdata/ninjs-demo.json")
	t2.OutputSchema = "./testdata/eventschema.json"
	osnewsml, err := t2.Transform(innewsml, nil)
	fmt.Println(err)
	fmt.Println(err, len(osnewsml))
	foutnewsml, e := os.Create("./testdata/newsmlout.json")
	fmt.Println(e, err)
	foutnewsml.Write(osnewsml)
	// XSD TEST
	t2.OutputSchema = "./testdata/eventdetails1.xsd"
	t2.OutputFormat = "application/xml"
	brr, _ := t2.Transform(innewsml, nil)
	frr, e := os.Create("./testdata/newsmlout.xml")

	frr.Write(brr)

	// fmt.Println(sc)

	// @TODO implement an array to compare names to for base types
	// ANd get the correct NAMESPACE

	/*for i, t := range tree {
		child := t.Children
		for _, c := range child {
			/*	if c.Copy().Name.Local == "element" {
				fmt.Println(c.Copy(), i, c.Copy().Attr[0].Name, c)
				bsd, _ := json.Marshal(c.Copy().Attr)
				fmt.Println(string(bsd))


			// fmt.Println(c.Copy(), i)
			if c.Copy().Name.Local == "import" {
				//fmt.Println(c.Copy(), i, c.Copy().Attr[0].Name, c)
				//	bsd, _ := json.Marshal(c.Copy().Attr)
				//fmt.Println(string(bsd))
			}
		}
	}*/
// get all the roots - ones that are element
//	fmt.Println(sc[0].TargetNS, sc[0].Doc, sc[0].FindType(xml.Name{Space: sc[0].TargetNS, Local: "knowledgeItem"}))
// target namespace
//}

/*
// roo2t := sc[0].FindType(xml.Name{Space: "http://www.w3.org/2001/XMLSchema", Local: "knowledgeItem"})
// fmt.Println("ROOT:", roo2t)

func (m *MappingAction) TransformTEST(in [][]byte) ([][]byte, error) {

	//output := make([][]byte, len(in))
	// convert each data point at a time
	for _, input := range in {

		// extract the metadata in a generic format
		var metaDataInput []map[string]any
		var err error

		switch m.InputFormat {
		case "text/csv":
			metaDataInput, err = csvBaseExtract(input)
		case "application/json":
			metaDataInput, err = jsonBaseExtract(input)
		case "application/yaml":
			metaDataInput, err = yamlBaseExtract(input)
		default:
			return nil, fmt.Errorf("unknown or invalid mime type of input data: %v", m.InputFormat)
		}

		if err != nil {
			return nil, fmt.Errorf("error extracting the input data: %v", err)
		}
		//	fmt.Println(metaDataInput)
		// @TODO add more robust schema getting
		// @TODO add more schema types
		// @TODO include some data formating like header order for csvs
		// depth can be used for schema width in CSVs
		//	b, err := os.ReadFile(m.OutputSchema)
		if err != nil {
			return nil, fmt.Errorf("error extracting the schema for the output data: %v", err)
		}

		// @TODO move the schmea to a function to handle several schema types
		//	var schema map[string]any
		//	json.Unmarshal(b, &schema)
		// fmt.Println(schema)
		//	mdPaths2 := ExtractMDpaths(schema)
		//fmt.Println(mdPaths2, "1")
		mdPaths := make(map[string]map[int][]mdProperties)

		//	bs, _ := os.ReadFile("./testdata/eventdetails1.xsd")
		//sc, _ := xsd.Parse(bs)
		//xmlProperties := XMLEncoderInformation{attr: make(map[string]bool), keyorder: []string{"knowledgeItem.xlmns"}}
		//	root := sc[0].FindType(xml.Name{Space: "http://iptc.org/std/nar/2006-10-01/", Local: "knowledgeItem"}).(*xsd.ComplexType)

		//mdPaths, _ = xmlProperties.xmlDataPaths(sc[0], root, "knowledgeItem.", "knowledgeItem.", mdPaths, []any{})
		var xmlProperties *XMLEncoderInformation
		mdPaths, xmlProperties, _ = xmlSchemaExtract("./testdata/eventdetails1.xsd")
		xmlProperties.attr["knowledgeItem.xmlns"] = true
		//fmt.Println(mdPaths, "2")
		//	fmt.Println(mdPaths, len(b), len(schema), e)
		//	fmt.Println(s2, "s2", len(s2))
		//	fmt.Println(mdPaths, "MAPATHS", len(mdPaths))
		// translate each row of metadata into the end result
		translatedData := make([]map[string]any, len(metaDataInput))

		for j, input := range metaDataInput {
			translatedData[j] = dataTranslate(input, mdPaths, m.Mapping.MappingDefinitions, m.Mapping.MissedTags)
			translatedData[j]["knowledgeItem.xmlns"] = "http://iptc.org/std/nar/2006-10-01/"
			// check for missed tags

			// Check for required objects being added
			// maybe add an error or warning
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
							fmt.Println(fullpath)
							_, ok := translatedData[j][fullpath]
							if !ok {
								fmt.Println(fullpath)
								translatedData[j][fullpath] = nil
							}
						}
					}
				}
			}
		}

		//	fmt.Println(oo)
		xmlProperties.namespaces = xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: xmlProperties.targetNameSpace}
		//	fmt.Println(root)
		//	fmt.Println(translatedData)
		bout, _ := xmlBuild(translatedData, *xmlProperties)
		//	fmt.Println(xmlProperties.attr)
		return [][]byte{bout}, nil
	}
	return nil, nil
}

/*
func recuse(tree *xmltree.Element, parent string, id int) []string {
	parents := []string{}
	//	fmt.Println(parent, tree.Name, len(tree.Children), id)
	children := tree.Children
	if len(children) == 0 {
		fmt.Println(tree, "TREER")
	}
	for _, child := range children {
		if strings.Contains(parent, "eventDetails") {
			fmt.Println(child, tree.Name.Local)
			fmt.Println(child.StartElement.Attr, "ATTIR", len(child.Children))
			fmt.Println(child, child.Scope, "LINER")
		}
		at := child.StartElement.Attr
		childName := ""
		if len(at) > 0 {
			for _, a := range at {
				if a.Name.Local == "name" {
					childName = a.Value
					break
				}
			}
		} else {

			childRun := recuse(&child, parent, id)
			parents = append(parents, childRun...)
			//fmt.Println("ATTIR", parents)
			continue
		}
		if childName == "" && strings.Contains(parent, "eventDetails") {
			//fmt.Println(child.Copy().Attr, "NAMELESS", parent)
			//	fmt.Println(child.Children)

		}
		if len(child.Children) != 0 && childName != "" {
			childRun := recuse(&child, parent+childName+".", id)
			parents = append(parents, childRun...)
		} else if childName == "" {
			childRun := recuse(&child, parent, id)
			parents = append(parents, childRun...)
		} else {

			parents = append(parents, parent+childName)
		}
	}
	return parents
}
*/
