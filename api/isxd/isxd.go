// Package isxd converts wave files to images
package isxd

import (
	"bytes"

	"github.com/metarex-media/mrx-tool/mrxUnitTest"
)

// Visualise returns a 700 x 200 image png visualising the soundwave.
// The louder the sound the more opaque that peak is.
func Visualise(input []byte, _ ...string) ([]byte, error) {

	var report bytes.Buffer
	byteStream := bytes.NewReader(input)
	// github.com/metarex-media/mrx-tool

	err := mrxUnitTest.MRXTest(byteStream, &report)

	if err != nil {
		return nil, err
	}
	// now save the image
	imageBuf := bytes.NewBuffer([]byte{})
	err = mrxUnitTest.DrawGraph(&report, imageBuf)

	if err != nil {
		return nil, err
	}

	return imageBuf.Bytes(), nil
}

// Visualise returns a 700 x 200 image png visualising the soundwave.
// The louder the sound the more opaque that peak is.
func Report(input []byte, _ ...string) ([]byte, error) {

	var report bytes.Buffer
	byteStream := bytes.NewReader(input)
	// github.com/metarex-media/mrx-tool

	mrxUnitTest.MRXTest(byteStream, &report)

	return report.Bytes(), nil
}
