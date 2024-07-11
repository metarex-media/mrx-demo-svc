package gps

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJPG2C2PA(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo03")

	for _, f := range fs {

		if strings.Contains(f.Name(), ".gpx") {
			b, _ := os.ReadFile("../../demodata/demo03/" + f.Name())
			genBytes, err := ConvertGPX(b)

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.json", f.Name()))

			got := sha256.New()
			got.Write(genBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			//
			Convey("Checking gpx is converted to w3c", t, func() {
				Convey(fmt.Sprintf("Generating w3c metadata from %s", f.Name()), func() {
					Convey("The hash matches the expected output", func() {

						So(err, ShouldBeNil)
						So(got.Sum(nil), ShouldResemble, expected.Sum(nil))

					})
				})
			})
		}
	}
}
