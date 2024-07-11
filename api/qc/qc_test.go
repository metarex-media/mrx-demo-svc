package qc

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQC(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo05")

	for _, f := range fs {

		if strings.Contains(f.Name(), ".xml") {
			b, _ := os.ReadFile("../../demodata/demo05/" + f.Name())
			genBytes, err := QCBarChart(b)

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.png", f.Name()))

			got := sha256.New()
			got.Write(genBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			//
			Convey("Checking qc bar charts are generated", t, func() {
				Convey(fmt.Sprintf("Generating the barchart of %s", f.Name()), func() {
					Convey("The hash matches the expected output", func() {

						So(err, ShouldBeNil)
						So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
					})
				})
			})
		}
	}
}
