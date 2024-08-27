// Package wavdraw converts wave files to images
package wavdraw

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"

	"github.com/cespare/xxhash/v2"
	"github.com/go-audio/wav"
)

// base the colours off of the file given
var colours = map[uint64]color.NRGBA{
	15798483879803493932: {0xC2, 0xA6, 0x49, 0xff},
	6227423963629637186:  {0x91, 0xB6, 0x45, 0xff},
	14840129082908086855: {0x9A, 0x3A, 0x73, 0xff},
	16896718556509083039: {0x43, 0x3F, 0x87, 0xff},
	7082774861642532824:  {191, 33, 33, 0xff},
	825898258475879677:   {14, 227, 227, 0xff},
}

// Visualise returns a 700 x 200 image png visualising the soundwave.
// The louder the sound the more opaque that peak is.
func Visualise(input []byte, _ ...string) ([]byte, error) {

	// decode the wave file
	inner := wav.NewDecoder(bytes.NewReader(input))

	// is it actually a wav file?
	buff, err := inner.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	pos := 0
	width := 700

	// just round up the size and don't worry about integer cut off
	size := len(buff.Data) / width
	levels := make([]int, width)
	max := 0

	// get the average sound for each level
	for i := 0; i < width; i++ {
		level := averageSound(buff.Data[pos : pos+size : pos+size])
		pos += size
		levels[i] = level

		if level > max {
			max = level
		}
	}

	baseImage := image.NewNRGBA(image.Rect(0, 0, width, 200))
	draw.Draw(baseImage,
		baseImage.Bounds(),
		&image.Uniform{color.White}, image.Point{}, draw.Over)
	var baseColor color.NRGBA

	// get the colour of known maps
	col, ok := colours[xxhash.Sum64(input)]
	if ok {
		baseColor = col
	} else {
		baseColor = color.NRGBA{A: 0xff}
	}

	// calculate the square root of the level
	// to make the transparent visually scale
	// if its linear then you can't see the low levels
	for i, level := range levels {

		opacity := (255 * 255 * level) / max
		scaledOpac := math.Sqrt(float64(opacity))
		baseColor.A = uint8(scaledOpac)

		height := ((level * 200) / max)
		start := (200 - height) / 2

		draw.Draw(baseImage,
			image.Rect(i, start, i+1, start+height),
			&image.Uniform{baseColor}, image.Point{}, draw.Over)
	}

	output := bytes.NewBuffer([]byte{})
	err = png.Encode(output, baseImage)

	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func averageSound(levels []int) int {
	total := 0
	for _, level := range levels {
		total += int(math.Abs(float64(level)))
	}

	return total / len(levels)
}
