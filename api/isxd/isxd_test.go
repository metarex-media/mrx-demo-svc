package isxd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVisualisation(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo13")

	for _, f := range fs {

		if strings.Contains(f.Name(), ".mxf") {
			b, _ := os.ReadFile("../../demodata/demo13/" + f.Name())
			graphBytes, err := Visualise(b)

			// f, _ := os.Create(fmt.Sprintf("./testdata/%v.png", f.Name()))
			// f.Write(genBytes)

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.png", f.Name()))

			got := sha256.New()
			got.Write(graphBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			//
			Convey("Checking mxf graph pngs are generated", t, func() {
				Convey(fmt.Sprintf("Generating the visual of %s", f.Name()), func() {
					Convey("The hash matches the expected output", func() {

						So(err, ShouldBeNil)
						So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
					})
				})
			})

			reportBytes, err := Report(b)

			//	f, _ := os.Create(fmt.Sprintf("./testdata/%v.yaml", f.Name()))
			//	f.Write(reportBytes)

			expecBytes, _ = os.ReadFile(fmt.Sprintf("./testdata/%v.yaml", f.Name()))

			gotReport := sha256.New()
			gotReport.Write(reportBytes)
			expectedReport := sha256.New()
			expectedReport.Write(expecBytes)

			//
			Convey("Checking mxf yaml reports are generated", t, func() {
				Convey(fmt.Sprintf("Generating the report of %s", f.Name()), func() {
					Convey("The hash matches the expected output", func() {

						So(err, ShouldBeNil)
						So(gotReport.Sum(nil), ShouldResemble, expectedReport.Sum(nil))
					})
				})
			})
		}
	}

}
