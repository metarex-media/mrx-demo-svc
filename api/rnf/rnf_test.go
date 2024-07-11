package rnf

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFFmpeg(t *testing.T) {

	fs, _ := os.ReadDir("../../demodata/demo12")

	keys := map[string]string{
		"IET.csv":         "IET",
		"springwatch.csv": "SW",
		"lostpast.csv":    "LP",
		"laundromat.csv":  "LM",
	}

	for _, f := range fs {

		if strings.Contains(f.Name(), ".csv") {
			b, _ := os.ReadFile("../../demodata/demo12/" + f.Name())
			genBytes, err := GenerateFFmpeg(b, keys[f.Name()])

			expecBytes, _ := os.ReadFile(fmt.Sprintf("./testdata/%v.txt", f.Name()))

			got := sha256.New()
			got.Write(genBytes)
			expected := sha256.New()
			expected.Write(expecBytes)

			//
			Convey("Checking ffmpeg scripts are generated", t, func() {
				Convey(fmt.Sprintf("Generating the ffmpeg of %s", f.Name()), func() {
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
func YamlMap() {

	//elapsed,in,oute,out,chapter,segment,subject,storyline-importance,shot,Rosana Prada,Elfed Howell,Werner Bleisteiner,Johan Bolin,Lucas Zwicker,primary-topic,secondary-topic,audience question,predation,procreation,threat

	c, _ := os.Open("/workspace/gl-mrx-demo-svc/xslxer/mrx-rnf-logie-baird-r0.csv")
	re := csv.NewReader(c)
	header, _ := re.Read()

	body, berr := re.Read()

	var output []IET

	for berr == nil {

		// / 21	L := Laundromat{}
		w := warningsYaml{
			Predation:   isTrue(body[17]),
			Threat:      isTrue(body[19]),
			Procreation: isTrue(body[18]),
		}
		t := topicYaml{
			Shot:      body[8],
			Subject:   body[6],
			Primary:   body[14],
			Secondary: body[15],
			Audience:  body[16],
		}
		in, _ := strconv.Atoi(body[1])
		out, _ := strconv.Atoi(body[2])
		imp, _ := strconv.Atoi(body[7])

		characters := make([]string, 0)
		for i, name := range body[9:14] {

			if name != "" {
				characters = append(characters, header[9+i])
			}
		}

		output = append(output, IET{
			Frames: frames{In: in, Out: out}, Importance: imp,
			Chapter:  body[4],
			Warnings: &w, Topics: &t,
			Characters: characters})

		body, berr = re.Read()
	}

	fmt.Println(body, output, len(body))
	b, _ := yaml.Marshal(output)
	f, _ := os.Create("./demodata/inputs/IET.yaml")
	f.Write(b)
}

func JsonMap() {

	//elapsed,in,oute,out,chapter,segment,subject,storyline-importance,shot,character1,character2,character3,character4,character5,primary-topic,secondary-topic,violence,threat,drugs,suicide

	c, _ := os.Open("/workspace/gl-mrx-demo-svc/api/rnf/demodata/inputs/cosmos-launderomat.csv")
	re := csv.NewReader(c)
	re.Read()

	body, berr := re.Read()

	var output []Laundromat

	for berr == nil {

		// /	L := Laundromat{}
		w := warnings{
			Violence: isTrue(body[16]),
			Threat:   isTrue(body[17]),
			Drugs:    isTrue(body[18]),
			Suicide:  isTrue(body[19]),
		}
		t := topic{
			Shot:      body[8],
			Subject:   body[6],
			Primary:   body[14],
			Secondary: body[15],
		}
		in, _ := strconv.Atoi(body[1])
		out, _ := strconv.Atoi(body[2])
		imp, _ := strconv.Atoi(body[7])
		output = append(output, Laundromat{
			In: in, Out: out, Importance: imp,
			Chapter:  body[4],
			Warnings: &w, Topics: &t})

		body, berr = re.Read()
	}

	fmt.Println(body, output, len(body))
	b, _ := json.MarshalIndent(output, "", "    ")
	f, _ := os.Create("./demodata/inputs/cosmos-laundromat.json")
	f.Write(b)
}

func isTrue(is string) bool {
	return strings.ToLower(is) == "true"
}

type Laundromat struct {
	In         int `json:"In"`
	Out        int `json:"Out"`
	Chapter    string
	Segment    string    `json:"Segment,omitempty"`
	Topics     *topic    `json:"Topics,omitempty"`
	Warnings   *warnings `json:"Warnings,omitempty"`
	Importance int       `json:"Importance,omitempty"`
}

type topic struct {
	Subject, Shot      string
	Primary, Secondary string
}
type warnings struct {
	Violence, Threat, Drugs, Suicide bool
}

//elapsed,in,oute,out,chapter,segment,subject,storyline-importance,shot,Rosana Prada,Elfed Howell,Werner Bleisteiner,Johan Bolin,Lucas Zwicker,primary-topic,secondary-topic,audience question,predation,procreation,threat

type IET struct {
	Frames     frames `yaml:"frames"`
	Chapter    string
	Segment    string        `json:"Segment,omitempty"`
	Topics     *topicYaml    `json:"Topics,omitempty"`
	Warnings   *warningsYaml `json:"Trigger-Warnings,omitempty"`
	Characters []string      `yaml:"Speaker,omitempty"`
	Importance int           `yaml:"Story,omitempty"`
}

type frames struct {
	In  int `yaml:"in(f)"`
	Out int `yaml:"out(f)"`
}

type warningsYaml struct {
	Predation, Procreation, Threat bool
}

type topicYaml struct {
	Subject, Shot      string
	Primary, Secondary string
	Audience           string `yaml:"Audience_Question,omitempty"`
}
*/
