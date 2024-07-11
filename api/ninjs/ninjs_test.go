package ninjs

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	xsdvalidate "github.com/terminalstatic/go-xsd-validate"
)

func TestToNewsML(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo10")

	for _, f := range fs {

		if strings.Contains(f.Name(), ".json") {
			b, _ := os.ReadFile("../../demodata/demo10/" + f.Name())
			genBytes, err := NinjsToNewsml(b)

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.xml", f.Name()))

			got := sha256.New()
			got.Write(genBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			validation := validate(genBytes, "./NewsItem.xsd")

			//
			Convey("Checking ninjs is converted to newsml", t, func() {
				Convey(fmt.Sprintf("Generating the newsml of %s", f.Name()), func() {
					Convey("The hash matches the expected output and it validates the xsd schema", func() {

						So(err, ShouldBeNil)
						So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
						So(validation, ShouldBeNil)
					})
				})
			})
		}
	}
}

func TestToMD(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo10")

	for _, f := range fs {

		if strings.Contains(f.Name(), ".json") {
			b, _ := os.ReadFile("../../demodata/demo10/" + f.Name())
			genBytes, err := NinJSToMD(b)

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.json", f.Name()))

			got := sha256.New()
			got.Write(genBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			//
			Convey("Checking ninjs is converted to markdown json", t, func() {
				Convey(fmt.Sprintf("Generating the markdown of %s", f.Name()), func() {
					Convey("The hash matches the expected output", func() {

						So(err, ShouldBeNil)
						So(got.Sum(nil), ShouldResemble, expected.Sum(nil))

					})
				})
			})
		}
	}
}

var parseOptions = xsdvalidate.ParsErrDefault

func validate(xmldoc []byte, schemaFile string) error {
	xsdvalidate.Init()
	defer xsdvalidate.Cleanup()

	in, err := os.Open(schemaFile)
	if err != nil {
		return fmt.Errorf("header: %s\nError: %s\n", xml.Header, err)
	}

	schema, _ := io.ReadAll(in)
	in.Close()

	xmlHandler, err := xsdvalidate.NewXmlHandlerMem(xmldoc, parseOptions)
	defer xmlHandler.Free()
	if err != nil {
		return fmt.Errorf("validation err :%v", err.Error())
	}

	xsdHandler, err := xsdvalidate.NewXsdHandlerMem(schema, parseOptions)
	defer xsdHandler.Free()
	if err != nil {
		return fmt.Errorf("validation err :%v", err.Error())
	}

	err = xsdHandler.Validate(xmlHandler, parseOptions)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	err = xsdHandler.ValidateMem(xmldoc, parseOptions)

	if err != nil {
		return fmt.Errorf("header: %s\nError: %s\n", xml.Header, err)
	}

	return nil
}
