package demotest

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDemo(t *testing.T) {

	_, e := http.Get("http://localhost:8080/")

	if e != nil {
		panic(fmt.Sprintf("error connecting to autoelt server at http://localhost:8080/: %v", e))
	}

	_, e = http.Get("http://localhost:9000/")

	if e != nil {
		panic(fmt.Sprintf("error connecting to services server at http://localhost:9000/: %v", e))
	}
	// MAKE sure both servers are running first

	tests := []testConditions{
		// DEMO 1
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.njs&outputMRXID=toNewsMD", input: "../demodata/demo01/SIPA_image_2.json",
			expected: "./testdata/expected/SIPA_image_2_out.json", dataType: echo.MIMEApplicationJSON},
		// Demo 3
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.gpx&outputMRXID=toW3C", input: "../demodata/demo03/Newhaven_Brighton.gpx",
			expected: "./testdata/expected/Newhaven_Brighton_out.json", dataType: echo.MIMEApplicationXML},
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.wav&outputMRXID=ToWaveform", input: "../demodata/demo03/European Robin - short.wav",
			expected: "./testdata/expected/European Robin - short_out.png", dataType: "audio/wav"},

		//Demo 4
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnc&outputMRXID=MRX.123.456.789.rnf&mapping=true", input: "../demodata/demo04/lostpast.csv",
			expected: "./testdata/expected/lostpast_out.csv", dataType: "text/csv"},
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnj&outputMRXID=MRX.123.456.789.rnf&mapping=true", input: "../demodata/demo04/cosmos-laundromat.json",
			expected: "./testdata/expected/cosmos-laundromat.csv", dataType: echo.MIMEApplicationJSON},
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rny&outputMRXID=MRX.123.456.789.rnf&mapping=true", input: "../demodata/demo04/IET.yaml",
			expected: "./testdata/expected/IET_out.csv", dataType: "application/yaml"},

		// Demo 5
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.vqc&outputMRXID=tograph", input: "../demodata/demo05/PulsarReport2Fail.xml",
			expected: "./testdata/expected/PulsarReport2Fail.png", dataType: "application/xml"},

		// Demo 6
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.ghj&outputMRXID=MRX.123.456.789.cxm&mapping=true", input: "../demodata/demo06/beach_camera.json",
			expected: "./testdata/expected/beach_camera.xml", dataType: "application/xml"},

		// Demo 7
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=fill", input: "../demodata/demo07/good.json",
			expected: "./testdata/expected/good_out.png", dataType: echo.MIMEApplicationJSON},
		/*	{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=stagger", input: "../demodata/demo07/red.json",
			expected: "./testdata/expected/red_out.png", dataType: echo.MIMEApplicationJSON},*/
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.bat&outputMRXID=fault", input: "../demodata/demo07/fire.json",
			expected: "./testdata/expected/fire_out.jpg", dataType: echo.MIMEApplicationJSON},

		//demo 8
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.jpg&outputMRXID=extractHeaderc2pa", input: "../demodata/demo08/truepic-20230212-library.jpg",
			expected: "./testdata/expected/truepic-20230212-library_out.json", dataType: "image/jpeg"},

		// demo 10
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.njs&outputMRXID=toNewsML", input: "../demodata/demo10/ntb_text.json",
			expected: "./testdata/expected/ntb_text_out.xml", dataType: echo.MIMEApplicationJSON},

		// demo 12
		{url: "http://localhost:8080/autoelt?inputMRXID=MRX.123.456.789.rnf&outputMRXID=generateFFmpeg&title=IET", input: "../demodata/demo12/IET.csv",
			expected: "./testdata/expected/IET_out.txt", dataType: "text/csv"},
	}

	for _, test := range tests {

		// get the test file
		inputBytes, errFile := os.ReadFile(test.input)

		// make the post request
		outputBytes, err := request(test.url, test.dataType, inputBytes)

		// check it against the expected outcome
		fbase, errExpec := os.ReadFile(test.expected)

		got := sha256.New()
		got.Write(outputBytes)
		expected := sha256.New()
		expected.Write(fbase)

		// to generate a file to compare outputs when things go wrong
		if fmt.Sprintf("%032x", got.Sum(nil)) != fmt.Sprintf("%032x", expected.Sum(nil)) {

			fbaseo, _ := os.Create(test.expected + ".out")
			fbaseo.Write(outputBytes)
		}

		//
		Convey(fmt.Sprintf("Checking %s runs as intended", test.input[12:18]), t, func() {
			Convey(fmt.Sprintf("Posting %s to %s", test.input, test.url), func() {
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

type testConditions struct {
	url             string
	input, expected string
	dataType        string
}

func request(posturl, dataType string, body []byte) ([]byte, error) {
	/*
		adress, data to post

	*/

	r, err := http.Post(posturl, dataType, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return io.ReadAll(r.Body)

}
