package wavdraw

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVisualisation(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo03")

	for _, f := range fs {

		if strings.Contains(f.Name(), ".wav") {
			b, _ := os.ReadFile("../../demodata/demo03/" + f.Name())
			genBytes, err := Visualise(b)

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.png", f.Name()))

			got := sha256.New()
			got.Write(genBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			//
			Convey("Checking audio pngs are generated", t, func() {
				Convey(fmt.Sprintf("Generating the visual of %s", f.Name()), func() {
					Convey("The hash matches the expected output", func() {

						So(err, ShouldBeNil)
						So(got.Sum(nil), ShouldResemble, expected.Sum(nil))
					})
				})
			})
		}
	}

}
