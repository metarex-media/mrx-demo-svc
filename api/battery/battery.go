// Package battery handles all the battery metadata to image transformations
package battery

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"

	_ "embed"
)

//go:embed data/*
var imgSystem embed.FS

const (
	/*
		energy0  = "data/0energy.png"
		energy20 = "data/20energy.png"
		energy50 = "data/50energy.png"
		energy60 = "data/60energy.png"
		energy80 = "data/80energy.png"*/

	baseBattery = "data/battery.png"

	fire     = "data/fire.jpg"
	puncture = "data/puncture.jpg"
	mystery  = "data/mystery.jpg"
	empty    = "data/empty.jpg"
	good     = "data/good.jpg"
)

type batteryELT struct {
	Percentage float64 `json:"percentage"`
	Status     string  `json:"status"`
}

/*
// BatteryToPNGStagger returns a different png
// for each 20% interval of the battery values.
// the images are returned scaled to a width of 200
func BatteryToPNGStagger(input []byte, _ ...string) ([]byte, error) {

	var batt batteryELT
	err := json.Unmarshal(input, &batt)

	if err != nil {
		return nil, fmt.Errorf("error handling data:" + err.Error())
	}

	var inputFile string

	switch {
	case batt.Percentage < 20:
		inputFile = energy0
	case batt.Percentage < 40:
		inputFile = energy20
	case batt.Percentage < 60:
		inputFile = energy50
	case batt.Percentage < 80:
		inputFile = energy60
	default:
		inputFile = energy80
	}

	f, err := imgSystem.Open(inputFile)

	if err != nil {
		return nil, fmt.Errorf("error extracting image file:" + err.Error())
	}

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("error extracting image:" + err.Error())
	}

	// fit to scale
	m := resize.Resize(200, 0, img, resize.Lanczos3)

	outputBytes := bytes.NewBuffer([]byte{})

	err = png.Encode(outputBytes, m)
	if err != nil {
		return nil, fmt.Errorf("error encoding the png: %v", err)
	}

	return outputBytes.Bytes(), nil
}*/

// ToPNG draws a scaled visual representation
// of the percentage. With different colors for different
// percent ranges.
// Red is < 30
// 30 < yellow < 60
// green is > 60
// It returns the png bytes
func ToPNG(input []byte, _ ...string) ([]byte, error) {

	var batt batteryELT
	err := json.Unmarshal(input, &batt)

	if err != nil {
		return nil, fmt.Errorf("error handling data:" + err.Error())
	}

	return batteryImage(batt)
}

// batteryImage draws a scaled visual representation
// of the percentage. With different colors for different
// percent ranges.
// It returns the png bytes
func batteryImage(battery batteryELT) ([]byte, error) {

	battbytes, err := imgSystem.Open(baseBattery)
	if err != nil {
		return nil, fmt.Errorf("error getting battery base image file: %v", err)
	}

	batteryImage, err := png.Decode(battbytes)
	if err != nil {
		return nil, fmt.Errorf("error getting battery base image: %v", err)
	}

	if battery.Percentage > 100 || battery.Percentage < 0 {
		return nil, fmt.Errorf("battery percentage of %v is not between 0 and 100", battery.Percentage)
	}

	// calculate the percentage width of the battery
	// the battery space is always 417 pixels wide
	width := (417 * battery.Percentage / 100)
	x1, y1, y2 := 63, 160, 352

	var baseColour color.Color
	switch {
	case battery.Percentage < 30: // red
		baseColour = color.NRGBA{186, 33, 9, 0xff}
	case battery.Percentage < 60: // yellow
		baseColour = color.NRGBA{194, 166, 73, 0xff}
	default: // green
		baseColour = color.NRGBA{145, 182, 69, 0xff}
	}
	draw.Draw(batteryImage.(*image.NRGBA), image.Rect(x1, y1, x1+int(width), y2), &image.Uniform{baseColour}, image.Point{}, draw.Over)

	outputBytes := bytes.NewBuffer([]byte{})

	err = png.Encode(outputBytes, batteryImage)
	if err != nil {
		return nil, fmt.Errorf("error encoding the png: %v", err)
	}

	return outputBytes.Bytes(), nil
}

const (
	// Firecode is the error code for the battery being on fire
	Firecode = "E00F"
	// Puncturecode is the error code for the battery being punctured
	Puncturecode = "E00P"
)

// FaultToJPEG returns an error image based on the error code.
func FaultToJPEG(input []byte, _ ...string) ([]byte, error) {

	var batt batteryELT
	err := json.Unmarshal(input, &batt)

	if err != nil {
		return nil, fmt.Errorf("error handling data:" + err.Error())
	}

	var faultfile string

	switch batt.Status {
	case Firecode:
		faultfile = fire
	case Puncturecode:
		faultfile = puncture

	default:
		switch {
		case len(batt.Status) > 0:
			faultfile = mystery
		case batt.Percentage == 0:
			faultfile = empty
		default:
			faultfile = good
		}
		// no fault file
	}

	f, _ := imgSystem.Open(faultfile)
	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("error decoding the jpeg: %v", err)
	}

	m := resize.Resize(200, 0, img, resize.Lanczos3)

	outputBytes := bytes.NewBuffer([]byte{})

	err = jpeg.Encode(outputBytes, m, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}

	return outputBytes.Bytes(), nil
}
