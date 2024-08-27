// Package rnf handles the transformations to the RNF metadata and the outputs
package rnf

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var clips = map[string]clipPropertiesAPI{
	"BBB": {Title: "BBB", Framerate: 24, OutputFolder: "./rnf"},
	"SW":  {Title: "SpringWatch", Framerate: 24, OutputFolder: "./rnf"},
	"LP":  {Title: "LostPast", Framerate: 25, OutputFolder: "./rnf"},
	"LM":  {Title: "Cosmos-Laundromat", Framerate: 25, OutputFolder: "./rnf"},
	"IET": {Title: "IET_Panel", Framerate: 25, OutputFolder: "./rnf"},
}

// GetFFmpegParams return the api parameters in the order
// they are handled in the GenerateFFmpeg function.
func GetFFmpegParams() []string {
	return []string{"title"}
}

// GenerateFFmpeg converts an RNF csv into the ffmpeg script for generating the segments
func GenerateFFmpeg(base []byte, destination ...string) ([]byte, error) {

	reader := csv.NewReader(bytes.NewReader(base))
	headers, err := reader.Read()

	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(headers, []string{"chapter", "in", "metadataTags", "out", "storyline-importance"}) {
		return nil, fmt.Errorf("invalid csv data format")
	}
	// chapter,in,metadataTags,out,storyline-importance
	body, berr := reader.Read()

	if berr != nil {
		return nil, fmt.Errorf("error reading data %v", berr)
	}

	i := 0
	outputScript := bytes.NewBuffer([]byte{})

	if len(destination) != 1 {
		return nil, fmt.Errorf("invalid number of title parameters passed, wanted 1 got %v", len(destination))
	}

	dest := destination[0]

	props := clips[dest]

	for berr == nil {

		tags := ""

		var metadata map[string]any

		json.Unmarshal([]byte(body[2]), &metadata)

		keys := make([]string, len(metadata))
		j := 0
		for k := range metadata {
			keys[j] = k
			j++
		}

		slices.Sort(keys)

		for _, k := range keys {
			v := metadata[k]
			// if it is not demo data
			if k != "elapsed" && k != "oute" {
				key := strings.Split(k, ".")
				k = key[len(key)-1]
				// handle metadata to be formed into part of the rnf name
				switch md := v.(type) {
				case string:

					if strings.ToLower(md) == "true" {
						if len(tags) != 0 {
							tags += ","
						}
						tags += strings.ReplaceAll(k, " ", "-")
					} else if md != "" {
						if len(tags) != 0 {
							tags += ","
						}
						tags += strings.ReplaceAll(v.(string), " ", "-")
					}
				case bool:

					if md {
						if len(tags) != 0 {
							tags += ","
						}
						tags += strings.ReplaceAll(k, " ", "-")
					}
				case []any:
					for _, k := range md {
						if len(tags) != 0 {
							tags += ","
						}
						kmid := fmt.Sprintf("%v", k)
						tags += strings.ReplaceAll(kmid, " ", "-")
					}
					// @TODO add mroe values as they need to be handled
				default:
					// fmt.Println(k, v)
				}

			}
		}
		in, _ := strconv.Atoi(body[1])
		out, _ := strconv.Atoi(body[3])
		tagger, err := ffmpegSegmentScript(key{TracksTag: tags}, props, in, out, i)
		if err != nil {
			return nil, fmt.Errorf("error generating the ffmpeg script: %v", err)
		}

		outputScript.Write([]byte(tagger))
		body, berr = reader.Read()

		i++
	}

	return outputScript.Bytes(), nil
}

const (
	outputVideoTag = "%s/%s_%04v_%s.mp4"
	// output2     = "%s/%s_01_01_%04v.mp4"
	outputVid   = "%s/segment%04v.mp4"
	outputAudio = "%s/segment%04v.wav"
	outputCSV   = "%s/segment%04v.csv"
	outputVTT   = "%s/%s_%04v_%s.vtt"
)

// framesToDur changes the time in frames to
// time in hours:miuntes:seconds:milliseconds
func framesToDur(frame, fps int) string {
	// mod number, divisor
	hourSize := fps * 60 * 60
	minuteSize := fps * 60

	hour := frame / hourSize
	if hour >= 1 {
		frame -= (hour * hourSize)
	}

	minute := frame / minuteSize
	if minute >= 1 {
		frame -= (minute * minuteSize)
	}

	second := frame / fps
	if second >= 1 {
		frame -= (second * fps)
	}

	subSecond := int((float64(frame) - 1) * float64(1000.0/float64(fps)))

	if frame == 0 {
		return fmt.Sprintf("%02d:%02d:%02d.000", hour, minute, second)
	}

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hour, minute, second, subSecond)
}

type inputs struct {
	InputFrames string `json:"InputFrames" yaml:"InputFrames"`
	InputAudio  string `json:"InputAudio" yaml:"InputAudio"`
}

func getInputFiles(title string) inputs {

	return inputs{InputFrames: "./rnf/" + title + "/frame%05d.jpg",
		InputAudio: "./rnf/" + title + "/audio.wav"} // config[title]
}

type clipPropertiesAPI struct {
	Size          int    `json:"size" yaml:"size"`
	Framerate     int    `json:"framerate" yaml:"framerate"`
	Title         string `json:"title" yaml:"title"`
	OutputFolder  string `json:"OutputFolder" yaml:"OutputFolder"` // output folder is rnf/bbb or springwatch1 etc
	SegmentFolder string
	SegmentsToGen []int `json:"SegmentsToGen"`
}
type key struct {
	ChapterTag, SegmentTag, TracksTag string
}

// Automatically build the tags when printing the key
// utilises the go String()  format.
func (k key) String() string {
	var place bool
	var sep = "-"
	var tags string
	if k.ChapterTag != "" {
		tags += k.ChapterTag
		place = true
	}

	if k.TracksTag != "" {
		if place {
			tags += sep
		}

		tags += k.TracksTag
		place = true
	}

	if k.SegmentTag != "" {
		if place {
			tags += sep
		}

		tags += k.SegmentTag
	}

	return tags
}

// generate and run the segmentation ffmpeg script
func ffmpegSegmentScript(key key, clipInfo clipPropertiesAPI, start, end, segment int) (string, error) {

	inputs := getInputFiles(clipInfo.Title)

	starter := framesToDur(start, clipInfo.Framerate)
	durere := framesToDur(end+1, clipInfo.Framerate)

	// fmt.Println(starter, durere, start, end, segment)
	ff2 := fmt.Sprintf("ffmpeg -y -i %v -ss %v -t %v -acodec pcm_s16le -ac 1 -ar 16000 %v\n",
		inputs.InputAudio, starter, durere, fmt.Sprintf(outputAudio, clipInfo.OutputFolder, segment))

	ff2 += fmt.Sprintf("ffmpeg -y -start_number %v -framerate %v -i %s -i %s -frames:v %v -vcodec mpeg4 -r %v -q:v 0 %s \n",
		start, clipInfo.Framerate, inputs.InputFrames, fmt.Sprintf(outputAudio, clipInfo.SegmentFolder, segment), end, clipInfo.Framerate, fmt.Sprintf(outputVideoTag, clipInfo.OutputFolder, clipInfo.Title, segment, key)) // fmt.Sprintf(outputVid, clipInfo.OutputFolder, segment))
	// fmt.Println(cmd.Args)

	return ff2, nil
}
