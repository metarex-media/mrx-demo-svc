// Package register is the local version of the metarex register (metarex.media/ui/reg)
package register

import (
	"encoding/json"

	"github.com/metarex-media/mrx-demo-svc/util/transformations/mapping"

	"gorm.io/gorm"
)

// MetarexReg is the go layout of the
// metarex register format
type MetarexReg struct {
	ID string `json:"metarexId"`
	// more metarex information
	MediaType   string `json:"mediaType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Timing      string `json:"timingIs"`
	TreatAs     string `json:"treatAs"`
	Version     string `json:"registerEntrySemVer,omitempty"`
	Mrx         Mrx    `json:"mrx"`
	Extra       any    `json:"any,omitempty"`
}

// Mrx contains the transformation properties
// and the information about the specification
type Mrx struct {
	Spec     string           `json:"specification"`
	Mapping  *mapping.Options `json:"mapping,omitempty"`
	Services []Services       `json:"services,omitempty"`
	Schema   []string         `json:"schema,omitempty"`
}

// Services contains all the information to
// make a request to a service.
type Services struct {
	API         string      `json:"API"`
	Method      string      `json:"method"`
	ID          string      `json:"metarexId"`
	Spec        string      `json:"APISchema"`
	Output      string      `json:"output"`
	Description string      `json:"description"`
	ServiceID   string      `json:"serviceID"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

// Parameter contains the any parameters
// a service has, including if they are required or not
type Parameter struct {
	Key         string `json:"key"`
	Optional    bool   `json:"optional"`
	Description string `json:"description"`
}

var (
	register map[string]MetarexReg
)

// MetarexRegister is the SQL databse entry for
// a metarex id.
type MetarexRegister struct {
	gorm.Model
	MRXID   string `gorm:"column:MRXID"`
	Reg     string `gorm:"column:RegisterValue"`
	OwnerID int    `gorm:"column:OwnerID"`
}

func init() {

	// utils.JSON(&slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}, "./tmp/")
	register = make(map[string]MetarexReg)

	// set up the register values to be saved locally as a map
	registerValues := []string{
		regNinJS, regNewsMl, battery, regBirdWav,
		regNewsMD, regW3CGPS, regGPX,
		regMRX,
		regRNF, regRNFInputCsv, regRNFInputYaml,
		regRNFInputJSON, regMRXC2p,
		regQC, regVnova,
		regUnrealCamera, regOtherCamera,
		regMXF,
	}

	for _, rv := range registerValues {
		var mrxID MetarexReg
		json.Unmarshal([]byte(rv), &mrxID)
		/*	if err != nil {
			fmt.Println(rv)
		}*/
		register[mrxID.ID] = mrxID

	}

}

// GetRegEntry returns the register entry for an ID
// as well as a bool declaring if was found or not
// like when you check maps for if the entry exists.
//
// reg, ok := register[mrxID]
func GetRegEntry(mrxID string) (MetarexReg, bool) {
	reg, ok := register[mrxID]
	return reg, ok
}

// Register values to be put
// into the database

var regUnrealCamera = `
{
    "metarexId": "MRX.123.456.789.ghj",
    "name": "Unreal Engine Camera Position",
    "description": "Camera location information from Unreal Engine as JSON. Sample available: Rexy's sunbathing scene.",
    "mediaType": "application/json",
    "timingIs": "clocked",
    "treatAs": "text",
    "registerEntrySemVer": "0.2.0",
    "mrx": {
      "specification": "https://docs.unrealengine.com/4.27/en-US/API/Runtime/Engine/Camera/UCameraComponent/"
    }
  }`
var regOtherCamera = `
  {
      "metarexId": "MRX.123.456.789.cxm",
      "name": "Demo Camera Position",
      "description": "Example data of camera position as XML.",
      "mediaType": "application/xml",
      "timingIs": "clocked",
      "treatAs": "text",
      "mrx": {
        "specification": "http://localhost:8080/schema/demo06output.xsd",
        "mapping": {
              "mappingDefinitions": {
                "intrinsicMatrix": ["intrinsics"],
		        "cameraRotation":  ["rotation"]
             }
         }
      }
    }`

var regNinJS = `{
    "metarexId": "MRX.123.456.789.njs",
    "name": "NAB NinJS",
    "description": "NinJS data to be used in the Metarex NAB Demo",
    "mediaType": "application/json",
    "timingIs": "clocked",
    "treatAs": "text",
    "mrx": {
        "specification": "https://iptc.org/std/ninjs/ninjs-schema_2.1.json",
        "services": [
        {
            "serviceID": "toNewsMD",
            "API": "http://localhost:9000/ninjsToMD",
            "APISchema": "./DemoAPI.yaml",
            "output": "MRX.123.456.789.nmd",
            "description":"Convert Ninjs to a simple markdown via an API"
        },
        {
            "serviceID": "toNewsML",
            "API": "http://localhost:9000/ninjsToNewsml",
            "APISchema": "./DemoAPI.yaml",
            "output": "MRX.123.456.789.nmj",
            "description":"Convert to newsml"
        }
        ]
    }
}`

var regBirdWav = `{
    "metarexId": "MRX.123.456.789.wav",
    "name": "Wav audio file",
    "description": "Wav audio file of birdsong",
    "mediaType": "audio/wav",
    "timingIs": "embedded",
    "treatAs": "binary",
    "mrx": {
        "services": [{
            "serviceID": "ToWaveform",
            "API": "http://localhost:9000/waveform",
            "APISchema": "./DemoAPI.yaml",
            "output": "image/png",
            "description":"Convert wave files to a png visualisation of the audio"
        }
        ]
    }
}`

var regNewsMl = `{
    "metarexId": "MRX.123.456.789.nmj",
    "name": "NAB NewsML",
    "description": "The full newsItem of the Newsml specification",
    "mediaType": "application/xml",
    "timingIs": "clocked",
    "treatAs": "text",
    "mrx": {
        "specification": "https://www.iptc.org/std/NewsML-G2/2.33/specification/individual/NewsML-G2_2.33-spec-NewsItem-Power.xsd",
        "mapping": {
            "MissedFieldsKey":"tag",
              "mappingDefinitions": {
                "name":            ["name", "literal"],
                "start":           ["startdate"],
                "end":             ["enddate"],
                "versionCreated":  ["versioncreated"],
                "contentCreated":  ["contentcreated"],
                "note":            ["ednote"],
                "lang":            ["language"],
                "standardversion": ["version"]
             }
         }
    }
}`

var regNewsMD = `{
    "metarexId": "MRX.123.456.789.nmd",
    "name": "NAB News Markdown",
    "description": "A simple version of ninjs that can be used to write markdown",
    "mediaType": "application/json",
    "timingIs": "clocked",
    "treatAs": "text",
    "mrx": {
        "specification":"http://localhost:8080/schema/demo01output.json"
    }
}`

var battery = `{
    "metarexId": "MRX.123.456.789.bat",
    "name": "NAB Battery",
    "description": "A simple battery format",
    "mediaType": "application/json",
    "timingIs": "clocked",
    "treatAs": "text",
    "mrx": {
        "specification": "http://localhost:8080/schema/demo07input.json",
        "services": [{
            "serviceID": "fill",
            "API": "http://localhost:9000/battery",
            "APISchema": "./DemoAPI.yaml",
            "output": "image/png",
            "description":"A service that converts the battery percentage to a scaled image"
        },
        {
            "serviceID": "fault",
            "API": "http://localhost:9000/batteryFault",
            "APISchema": "./DemoAPI.yaml",
            "output": "image/jpeg",
            "description":"A service that converts the battery information into a visualisation of the fault"
        },
        {
            "serviceID":"stagger",
            "API": "http://localhost:9000/batteryStagger",
            "APISchema": "./DemoAPI.yaml",
            "output":"image/png",
            "description":"A service that converts the battery percentage to a staggered image"
        }
        ]
    }
}`

var regGPX = `
{
    "metarexId": "MRX.123.456.789.gpx",
    "name": "topografix gps exchange format",
    "description": "Open source ",
    "mediaType": "application/xml",
    "timingIs": "embedded",
    "treatAs": "text",
    "mrx": {
        "specification": "https://www.topografix.com/GPX/1/1/gpx.xsd",
        "services": [{
            "serviceID": "toW3C",
            "API": "http://localhost:9000/gps",
            "APISchema": "./DemoAPI.yaml",
            "output": "MRX.123.456.789.gps",
            "description":"A service tha converts gps exchange xml to w3c json"
        }]
    }
}`

var regQC = `
{
    "metarexId": "MRX.123.456.789.vqc",
    "name": "Venera quality control report",
    "description": "quality",
    "mediaType": "application/xml",
    "timingIs": "embedded",
    "treatAs": "text",
    "mrx": {
        "specification": "http://localhost:8080/schema/Pulsar_Report.xsd",
        "services": [{
            "serviceID": "tograph",
            "API": "http://localhost:9000/qcToGraph",
            "APISchema": "./DemoAPI.yaml",
            "output": "image/png",
            "description":"A service tha converts venera qc reports to a graphical representation"
        }]
    }
}`

var regW3CGPS = `
{
    "metarexId": "MRX.123.456.789.gps",
    "name": "W3C GPS",
    "description": "The gps location of a device, provided by the internet browser",
    "mediaType": "application/json",
    "timingIs": "embedded",
    "treatAs": "text",
    "mrx": {
        "specification": "http://localhost:8080/schema/demo03output.json"

    }
}
`

var regMRX = `{
    "metarexId": "MRX.123.456.789.jpg",
    "name": "NAB MRX",
    "description": "A pseudo identity for jpeg files loaded with c2pa header metadata",
    "mediaType": "application/octet-stream",
    "timingIs": "embedded",
    "treatAs": "binary",
    "mrx": {
        "services": [
        {
            "serviceID": "extractHeaderc2pa",
            "API": "http://localhost:9000/C2PAExtract",
            "APISchema": "./DemoAPI.yaml",
            "output": "MRX.123.456.789.c2p",
            "description":"Extract the header data of an mrx as json"
        }
        ]
    }
}`

var regMRXC2p = `{
    "metarexId": "MRX.123.456.789.c2p",
    "name": "C2PA manifest",
    "description": "The c2pa manifest",
    "mediaType": "application/json",
    "timingIs": "clocked",
    "treatAs": "text",
    "mrx": {
        "specification":"https://github.com/contentauth/c2patool/blob/main/schemas/manifest-definition.json"
    }
}`

var regVnova = `{
    "metarexId": "MRX.123.456.789.brq",
    "name": "Level Of Quality",
    "description": "A level of quality log, describing the contents if a .loq or .yuv file(s)",
    "mediaType": "application/json",
    "timingIs": "clocked",
    "treatAs": "text",
    "mrx": {
        "specification":"http://localhost:8080/schema/demo09schema.json"
    }
}`

var regRNFInputCsv = `{
    "metarexId": "MRX.123.456.789.rnc",
    "name": "NAB RNF csv",
    "description": "An input csv of tagged media metadata",
    "mediaType": "text/csv",
    "timingIs": "embedded",
    "treatAs": "text",
    "mrx": {
        "specification":"http://localhost:8080/schema/demo04rnc.json"
    }
}`

var regRNFInputYaml = `{
    "metarexId": "MRX.123.456.789.rny",
    "name": "NAB RNF yaml",
    "description": "An input yaml of tagged media metadata",
    "mediaType": "application/yaml",
    "timingIs": "embedded",
    "treatAs": "text",
    "mrx": {
        "specification":"http://localhost:8080/schema/demo04rny.json"
    }
}`

var regRNFInputJSON = `{
    "metarexId": "MRX.123.456.789.rnj",
    "name": "NAB RNF json",
    "description": "An input json of tagged media metadata",
    "mediaType": "application/json",
    "timingIs": "embedded",
    "treatAs": "text",
    "mrx": {
        "specification":"http://localhost:8080/schema/demo04rny.json"
    }
}`

var regRNF = `{
    "metarexId": "MRX.123.456.789.rnf",
    "name": "NAB RNF csv ",
    "description": "The rnf normalised format, which is used to generate the rnf segments",
    "mediaType": "text/csv",
    "timingIs": "embedded",
    "treatAs": "text",
    "mrx": {
        "specification": "http://localhost:8080/schema/rnfSchema.json",
        "services": [
        {
            "serviceID": "generateFFmpeg",
            "API": "http://localhost:9000/ffmpeg",
            "APISchema": "./DemoAPI.yaml",
            "output": "text/plain",
            "description":"Generate the ffmpeg script to build the rnf segments",
            "parameters": [{
                "key": "title",
                "optional":false,
                "description" : "the title of film to be segemented for RNF"
            }]
        }
        ],
        "mapping": {
            "MissedFieldsKey":"metadataTags",
              "mappingDefinitions": {
                "in":                   ["in", "In", "in(f)"],
                "out":                  ["out", "Out", "out(f)"],
                "chapter":              ["chapter", "Chapter"],
                "storyline-importance": ["storyline-importance", "Storyline-importance", "Importance", "Story"]
             }
         }
    }
}`

var regMXF = `{
    "metarexId": "MRX.123.456.789.mxf",
    "name": "MXF metadata file",
    "description": "An mxf file",
    "mediaType": "application/mxf",
    "timingIs": "embedded",
    "treatAs": "binary",
    "mrx": {
        "services": [
            {
                "serviceID": "MXFToGraph",
                "API": "http://localhost:9000/mxfGraph",
                "APISchema": "./DemoAPI.yaml",
                "output": "image/png",
                "description": "Convert an mxf file to a graph detailing its test report"
            },
            {
                "serviceID": "MXFToReport",
                "API": "http://localhost:9000/mxfReport",
                "APISchema": "./DemoAPI.yaml",
                "output": "application/yaml",
                "description": "Run a series of tests on the mxf file, with the report given as a yaml"
            }
        ]
    }
}`
