// package qc handles Venera qc metadata
package qc

import (
	"bytes"
	"fmt"
	"math"

	"github.com/antchfx/xmlquery"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

// QCBarChart converts a venera qc report into a simple bar chart with
// Critical Alert, warning and error fields.
// if all are 0 then the report has passed
func QCBarChart(input []byte, _ ...string) ([]byte, error) {

	doc, err := xmlquery.Parse(bytes.NewReader(input))

	if err != nil {
		return nil, err
	}

	// xml can be comprised of several PulsarReports
	PR := xmlquery.Find(doc, "//PulsarReport")
	fields := []string{"CriticalAlert", "Warning", "Error"}

	//testPass := colour.CNRGBA64{0x91 << 8, 0xB6 << 8, 0x45 << 8, 0xffff, colour.ColorSpace{}}
	//cb := context.Background()
	//img := image.NewNRGBA(image.Rect(0, 0, 800, 300))
	max := 0
	counts := make([]float64, len(fields))
	messSample := make([]string, len(fields)*3)
	pass := "pass"
	for i, f := range fields {
		messages := getField(PR, f)

		// if any messages are found then the report is failed
		if len(messages) > max {
			max = len(messages)
			pass = "fail"
		}

		counts[i] = float64(len(messages))
		if len(messages) < 3 {
			messages = append(messages, make([]string, 4-len(messages))...)
		}
		for j := 0; j < 3; j++ {
			if len(messages[j]) > 70 {
				messSample[i*3+j] = messages[j][:70]
			} else {
				messSample[i*3+j] = messages[j]
			}
		}

	}

	// want a max y axis height of 5
	if max < 5 {
		max = 5
	}

	// generate the chart ticks
	ticks := make([]chart.Tick, 6)
	step := int(math.RoundToEven(float64(max) / 5))
	for i := 0; i < 6; i++ {
		ticks[i] = chart.Tick{float64(i * step), fmt.Sprintf("%v", i*step)}
	}

	graph := chart.BarChart{
		Title: fmt.Sprintf("Qc Report results : %s", pass),
		YAxis: chart.YAxis{
			Ticks: ticks, Range: &chart.ContinuousRange{Min: 0, Max: float64(max)}},

		Width:  500,
		Height: 300,
	}

	red := drawing.Color{0xff, 0, 0, 0xff}
	orange := drawing.Color{0xff, 0x80, 0, 0xff}
	yellow := drawing.Color{0xff, 0xff, 0, 0xff}
	colours := []drawing.Color{red, orange, yellow}
	empty := drawing.Color{0xff, 0xff, 0, 0x00}
	values := make([]chart.Value, 3)

	for i, c := range counts {

		if c == 0 {
			values[i] = chart.Value{
				Label: fields[i], Value: c,
				Style: chart.Style{FillColor: empty, StrokeColor: empty},
			}
		} else {
			values[i] = chart.Value{
				Label: fields[i], Value: c,
				Style: chart.Style{FillColor: colours[i], StrokeColor: colours[i]},
			}
		}

	}

	graph.Bars = values

	buf := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buf)
	fmt.Println(err)

	return buf.Bytes(), nil

}

func getField(pulsarReports []*xmlquery.Node, field string) []string {
	fields := []string{}
	for _, p := range pulsarReports {

		//
		res := messageExtractor(p, field)
		fields = append(fields, res...)

	}

	return fields
}

// messageExtractor searches for the number of fields in a document.
// this fits the venera xml schema and probably won't work with other data formats.
func messageExtractor(results *xmlquery.Node, field string) []string {

	out := xmlquery.Find(results, "//@"+field+"sNum")

	messages := []string{}
	for _, o := range out {
		ebody := o.Parent
		for _, attr := range ebody.Attr {
			if attr.Name.Local == field+"sNum" {

				if attr.Value != "0" {
					res := xmlquery.Find(ebody, fmt.Sprintf(`//*[contains(local-name(), "%s") and not(contains(local-name(), "%ss"))]`, field, field))

					for _, r := range res {
						for _, rattr := range r.Attr {
							// fmt.Println(rattr)
							if rattr.Name.Local == "Message" {
								messages = append(messages, fmt.Sprintf(" %s\":\" %s", r.Data, rattr.Value))
								// fail = append(fail, Check{Name: r.Data, ErrorMessage: rattr.Value})
							}
						}
					}
				}
			}
		}
	}

	return messages
}
