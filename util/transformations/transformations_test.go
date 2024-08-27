package transformations

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/metarex-media/mrx-demo-svc/util/transformations/api"
	"github.com/metarex-media/mrx-demo-svc/util/transformations/mapping"
	. "github.com/smartystreets/goconvey/convey"
)

func TestXxx(t *testing.T) {
	// out, e := MrxPathFinder("MRX.123.456.789.njs", "MRX.123.456.789.c2p", "", false)

	ins := []string{"MRX.123.456.789.njs", "MRX.123.456.789.rny", "MRX.123.456.789.bat"}
	outs := []string{"MRX.123.456.789.nmj", "MRX.123.456.789.rnf", "image/jpeg"}
	expecURL := []string{"localhost:8080/autoelt?inputMRXID=MRX.123.456.789.njs&outputMRXID=toNewsML",
		"localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rny&outputMRXID=MRX.123.456.789.rnf&mapping=true",
		"localhost:8080/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=fault"}

	for i, in := range ins {
		out, err := ServicesMrxPathFinder(in, outs[i], "")

		Convey("Checking service paths are found", t, func() {
			Convey(fmt.Sprintf("Transforming %s to %s", in, outs[i]), func() {
				Convey("Generates an apicall of "+expecURL[i], func() {
					So(err, ShouldBeNil)
					So(out[0].APICall, ShouldResemble, expecURL[i])

				})
			})
		})

	}

	toMap := []bool{false, true}
	actions := [][][]Action{
		{{&api.Action{API: "http://localhost:9000/ninjsToNewsml", MrxID: "MRX.123.456.789.njs", ResponseMIMEType: "application/json", APISchemaLocation: "./DemoAPI.yaml", APIParams: []api.Parameter{}}}},
		{{&mapping.Action{OutputSchema: "http://localhost:8080/schema/rnfSchema.json", MrxID: "MRX.123.456.789.rnf", InputFormat: "application/yaml", OutputFormat: "text/csv", InputTiming: "embedded", Mapping: mapping.Options{ConvertTypes: false, MissedTags: "metadataTags", MappingDefinitions: map[string][]string{"chapter": {"chapter", "Chapter"}, "in": {"in", "In", "in(f)"}, "out": {"out", "Out", "out(f)"}, "storyline-importance": {"storyline-importance", "Storyline-importance", "Importance", "Story"}}}}}},
	}

	for i, in := range ins[:len(ins)-1] {
		out, err := MrxPathFinder("", in, outs[i], toMap[i])

		Convey("Checking service paths are found and the correct actions are produced", t, func() {
			Convey(fmt.Sprintf("Transforming %s to %s", in, outs[i]), func() {
				Convey(fmt.Sprintf("Generates an %v action of %v ", reflect.TypeOf(actions[i][0][0]), actions[i][0][0]), func() {
					So(err, ShouldBeNil)
					So(out.Actions, ShouldResemble, actions[i])

				})
			})
		})

	}

	// register EXPORt

}
