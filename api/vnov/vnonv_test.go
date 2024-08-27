package vnov

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
)

// TestGenJsons generates the jsons to be used by demo 09
func TestGenJsons(_ *testing.T) {
	// mrxsunbath_7680x4320_yuv420p
	names := []string{"mrxsunbath_7680x4320_yuv420p", "mrxfish_7680x4320_yuv420p", "mrxfire_7680x4320_yuv420p", "mrxroar_7680x4320_yuv420p"}
	for _, n := range names {

		// generate the .vc6 information
		base := quality{
			Loq:      0,
			W:        7680,
			H:        4320,
			Size:     sizes[n],
			FileName: n + ".vc6",
		}

		b, _ := json.MarshalIndent(base, "", "     ")
		f, _ := os.Create("../../demodata/demo09/" + n + ".vc6.json")
		f.Write(b)

		// generate thh decoded file data
		// starting at a level of quality of 1
		for i := 1; i < 5; i++ {
			base := quality{
				Loq:      i,
				W:        7680 / int(math.Pow(2, float64(i))),
				H:        4320 / int(math.Pow(2, float64(i))),
				FileName: fmt.Sprintf("dec_%v.yuv.loq%v", n, i),
			}
			// Size:     int(float64((7680*4320)/((i+1)*(i+1))*72) * 1.5),
			base.Size = int(float64((base.H*base.W)*72) * 1.5)

			b, _ := json.MarshalIndent(base, "", "     ")
			f, _ := os.Create(fmt.Sprintf("../../demodata/demo09/dec_%v.yuv.loq%v.json", n, i))
			f.Write(b)
		}

	}
}

// manual file sizes
var sizes = map[string]int{
	"mrxsunbath_7680x4320_yuv420p": 449484954,
	"mrxfish_7680x4320_yuv420p":    444653378,
	"mrxfire_7680x4320_yuv420p":    449198101,
	"mrxroar_7680x4320_yuv420p":    446362538}

type quality struct {
	Loq      int    `json:"levelOfQuality"`
	W        int    `json:"width"`
	H        int    `json:"height"`
	Size     int    `json:"fileSizeBytes"`
	FileName string `json:"fileName"`
}
