package battery

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBatteryOut(t *testing.T) {

	good := []string{"red.json", "yellow.json", "good.json"}

	for _, g := range good {

		b, _ := os.ReadFile("../../demodata/demo07/" + g)
		genBytes, err := ToPNG(b)

		expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.png", g))

		got := sha256.New()
		got.Write(genBytes)
		expected := sha256.New()
		expected.Write(expecBytes)

		//
		Convey("Checking battery json is convertd to png", t, func() {
			Convey(fmt.Sprintf("Generating the battery from %s", g), func() {
				Convey("The hash matches the expected output", func() {

					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))

				})
			})
		})

	}

	bad := []string{"fire.json", "puncture.json", "unknown.json", "good.json"}

	for _, g := range bad {

		b, _ := os.ReadFile("../../demodata/demo07/" + g)

		genBytes, err := FaultToJPEG(b)

		expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.jpg", g))

		got := sha256.New()
		got.Write(genBytes)
		expected := sha256.New()
		expected.Write(expecBytes)

		//
		Convey("Checking battery json with faults generate the correct jpg", t, func() {
			Convey(fmt.Sprintf("Generating the fault from %s", g), func() {
				Convey("The hash matches the expected output", func() {

					So(err, ShouldBeNil)
					So(got.Sum(nil), ShouldResemble, expected.Sum(nil))

				})
			})
		})

	}
}
