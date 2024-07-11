package mrx

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJPG2C2PA(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo08")

	for _, f := range fs {

		if strings.Contains(f.Name(), ".jpg") || strings.Contains(f.Name(), ".jpeg") {
			b, _ := os.ReadFile("../../demodata/demo08/" + f.Name())
			genBytes, err := JpegC2PA(b)

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.json", f.Name()))

			got := sha256.New()
			got.Write(genBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			//
			Convey("Checking C2PA metadata is extracted out of the jpeg", t, func() {
				Convey(fmt.Sprintf("Generating the c2pa from %s", f.Name()), func() {
					Convey("The hash matches the expected output", func() {

						So(err, ShouldBeNil)
						So(got.Sum(nil), ShouldResemble, expected.Sum(nil))

					})
				})
			})
		}
	}
}

/*
func TestSetup(t *testing.T) {
	files, _ := os.ReadDir("./c2pa-jpeg")

	hashmap := make(map[uint64]string)

	for _, f := range files {
		if !f.IsDir() {
			b, _ := os.ReadFile("./c2pa-jpeg/" + f.Name())
			in := xxhash.Sum64(b)
			hashmap[in] = f.Name()

			db, de := os.ReadFile("./c2pa-jpeg/manifests/" + C2PAHash[in] + "/detailed.json")
			//	fmt.Println(string(db), de)
			if de != nil {
				fmt.Println(len(db), de)
			}
		}
	}
	/*
	   fmt.Println(hashmap)
	   fmt.Println(len(hashmap), len(files))
	   // Only pass t into top-level Convey calls

	   	Convey("Given some integer with a starting value", t, func() {
	   		x := 1

	   		Convey("When the integer is incremented", func() {
	   			x++

	   			Convey("The value should be greater by one", func() {
	   				So(hashmap, ShouldEqual, 2)
	   			})
	   		})
	   	})

}

/*
// go get github.com/mmTristan/mrx-tool@b182b33c3f0c682ddd26c1a47ebdca028bf3a340
func TestXxx(t *testing.T) {
	// f, _ := os.Open("./testdata/velDemo.mrx")
	// f, _ := os.Open("./testdata/freeMXF.mxf")
	f, _ := os.Open("./testdata/RDD48.mxf")
	out, _ := os.Create("./testdata/out.json")
	// generated.DecodeTAS_07_DMS_IdentifierRoleCode([]byte("hello"))
	e := mxfparse.MXFParse(f, out, &mxfparse.MrxParseOptions{})
	//e := decode.StreamDecode(f, os.Stdout, []int{1}, true)

	fmt.Println(e)
}*/
