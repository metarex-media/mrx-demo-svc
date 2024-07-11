// package ninjs handles all news metadata and their transformations, such as ninjs and newsml
package ninjs

import (
	"encoding/json"
)

// https://github.com/iptc/newsinjson/tree/main/examples - ninjs example data here

// NinJSToMD takes a ninjs file and translates it to a json format that
// can be used to write markdown files.
func NinJSToMD(ninput []byte, _ ...string) ([]byte, error) {

	var ninjs NinJS
	err := json.Unmarshal(ninput, &ninjs)
	if err != nil {
		return nil, err
	}

	MarkDown := MDFormat{
		Slug: ninjs.Slugline,
	}

	// extract main headlines only
	// no different formats
	for _, headline := range ninjs.Headlines {
		if headline.Role == "main" {
			MarkDown.Title = "# " + headline.Value
		}
	}

	// only use text bodies
	for _, body := range ninjs.Bodies {
		if body.Role == "text" {
			MarkDown.ShortSummary = body.Value

		}
	}

	// if no body was found try again with
	// a=finding a summary
	if MarkDown.ShortSummary == "" {
		for _, desc := range ninjs.Descriptions {
			if desc.Role == "text" {
				MarkDown.ShortSummary = desc.Value

			}
		}

	}

	return json.MarshalIndent(MarkDown, "", "    ")
}

// MDFormat is the markdown format
type MDFormat struct {
	Title        string `json:"title"`
	ShortSummary string `json:"shortSummary"`
	LongSummary  string `json:"longSummary,omitempty"`
	Slug         string `json:"slug"`
}

// NinJS is the ninjs go layout.
// not every field in NinJS is present
// in this struct/
type NinJS struct {
	URI                string         `json:"uri"`
	Altids             []AltID        `json:"altids"`
	Type               string         `json:"type"`
	Representationtype string         `json:"representationtype"`
	Genres             []Genre        `json:"genres"`
	Profile            string         `json:"profile"`
	Version            string         `json:"version"`
	Versioncreated     string         `json:"versioncreated"`
	Contentcreated     string         `json:"contentcreated"`
	Embargoed          string         `json:"embargoed"`
	Urgency            int            `json:"urgency"`
	Slugline           string         `json:"slugline"`
	Headlines          []Headline     `json:"headlines"`
	Descriptions       []Description  `json:"descriptions"`
	Bodies             []Body         `json:"bodies"`
	Copyrightholder    string         `json:"copyrightholder"`
	Copyrightnotice    string         `json:"copyrightnotice"`
	Usageterms         string         `json:"usageterms"`
	Ednote             string         `json:"ednote"`
	Language           string         `json:"language"`
	Eventdetails       []EventDetails `json:"eventdetails"`
	Pubstatus          string         `json:"pubstatus"`
	Subjects           []Subjects     `json:"subjects"`
}

// Ninjs Subjects
type Subjects struct {
	Literal    string `json:"literal"`
	Name       string `json:"name"`
	URI        string `json:"uri"`
	Rel        string `json:"rel"`
	Creator    string `json:"creator"`
	Relevance  int    `json:"relevance"`
	Confidence int    `json:"confidence"`
}

// Ninjs Headline
type Headline struct {
	Role  string `json:"role"`
	Value string `json:"value"`
}

// Ninjs Body
type Body struct {
	Role        string `json:"role"`
	Charcount   int    `json:"charcount"`
	Contenttype string `json:"contenttype"`
	Value       string `json:"value"`
}

// Ninjs Description
type Description struct {
	Role  string `json:"role"`
	Value string `json:"value"`
}

// NinJs AltID
type AltID struct {
	Role  string `json:"role"`
	Value string `json:"value"`
}

// Ninjs Genre
type Genre struct {
	Literal string `json:"literal"`
	URI     string `json:"uri"`
	Name    string `json:"name"`
}

// Ninjs EventDetails
type EventDetails struct {
	Eventstatus string    `json:"eventstatus"`
	Dates       Dates     `json:"dates"`
	Organiser   Organiser `json:"organiser"`
}

// Ninjs Organiser
type Organiser struct {
	Name    string    `json:"name"`
	Rel     string    `json:"rel"`
	URI     string    `json:"uri"`
	Literal string    `json:"literal"`
	Symbols []Symbols `json:"symbols"`
}

// Ninjs Symbols
type Symbols struct {
	Ticker     string `json:"ticker"`
	Exchange   string `json:"exchange"`
	Symboltype string `json:"symboltype"`
	Symbol     string `json:"symbol"`
}

// Ninjs Dates
type Dates struct {
	Startdate string `json:"startdate"`
	Enddate   string `json:"enddate"`
}
