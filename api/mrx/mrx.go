// package mrx handles all the file metadata and their extraction from files
package mrx

import (
	"embed"
	"fmt"

	_ "embed"

	"github.com/cespare/xxhash/v2"
)

//go:embed c2pa-jpeg/manifests/*
var mdSystem embed.FS

// The hash map of all te jpeg files used in the demo
// each one relates to the folder name where the manifest
// is stored.
var c2PAHash = map[uint64]string{
	206660612399611532: "adobe-20220124-CAICAI", 1779168395016079233: "adobe-20220124-C", 2922633521276313851: "adobe-20220124-CACAICAICICA", 3491551332456534290: "adobe-20220124-CICA", 3799586296415158911: "adobe-20220124-E-sig-CA", 4844496264397448694: "truepic-20230212-camera", 5618526863113647572: "adobe-20220124-E-dat-CA", 6852618351824180727: "adobe-20220124-XCA", 7326278477035881403: "truepic-20230212-library", 8731300505378100326: "nikon-20221019-building", 8924334609209574706: "adobe-20220124-CII", 9262652483651377359: "truepic-20230212-landscape", 9281016480572015694: "adobe-20220124-CAIAIIICAICIICAIICICA", 10306458505631521614: "adobe-20220124-XCI", 10361118825179637513: "adobe-20220124-E-clm-CAICAI", 11475358003335232434: "adobe-20220124-CAI", 13574538305397050894: "adobe-20220124-CICACACA", 14511528915343955556: "adobe-20220124-E-uri-CIE-sig-CA", 15233633124507241297: "adobe-20220124-CI", 16456076172579578992: "adobe-20220124-CIE-sig-CA", 16466799697891799827: "adobe-20221004-ukraine_building", 17514532766207691289: "adobe-20220124-CAICA", 17582052511769936677: "adobe-20220124-CA", 17782787895197731367: "adobe-20220124-CACA", 18284171616483865122: "adobe-20220124-E-uri-CA"}

// Jpeg2CPA pseduo extracts metadata by comparing the hashes of the input
// and extracting a matching file if available
func JpegC2PA(input []byte, _ ...string) ([]byte, error) {

	// get the input hash straight away
	in := xxhash.Sum64(input)

	// is it a known file
	filename, ok := c2PAHash[in]
	if !ok {
		return nil, fmt.Errorf("no C2PA metadata found")
	}
	db, de := mdSystem.ReadFile("c2pa-jpeg/manifests/" + filename + "/detailed.json")
	//	fmt.Println(string(db), de)

	if de != nil {
		return nil, fmt.Errorf("error extracting metadata")
	}

	return db, nil
}
