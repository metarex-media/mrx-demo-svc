// Package ninjs handles all news metadata and their transformations, such as ninjs and newsml
package ninjs

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

// ToNewsMl converts a ninjs byte stream to a newsml []byte stream
func ToNewsMl(input []byte, _ ...string) ([]byte, error) {

	var newItem NinJS
	err := json.Unmarshal(input, &newItem)
	if err != nil {
		return nil, err
	}

	defaultTime := time.Now().Format(timeFormat)

	// set root item attributes
	var out NewsItem
	out.Xmlns = "http://iptc.org/std/nar/2006-10-01/"
	out.Standard = "NewsML-G2"
	out.Standardversion = "2.24"
	out.Version = "11"

	// extract copyright info
	c := RightsInfo{
		CopyrightHolder: copyright{Name: newItem.Copyrightholder},
		CopyrightNotice: newItem.Copyrightnotice,
	}

	// get the item metadata
	im := itemMeta{VersionCreated: newItem.Versioncreated,
		Provider:  Provider{URI: newItem.URI},
		PubStatus: PubStatus{Qcode: "stat:" + newItem.Pubstatus}}

	// set it so it passes the schema
	if im.PubStatus.Qcode == "stat:" {
		im.PubStatus.Qcode = "stat:unknown"
	}

	if im.VersionCreated == "" {
		im.VersionCreated = defaultTime
	}

	// extract the subjects
	subs := make([]Subject, len(newItem.Subjects))
	for i, s := range newItem.Subjects {
		subs[i] = Subject{Name: s.Name, Qcode: s.URI}
	}

	cm := Contentmeta{Slugline: newItem.Slugline, Subject: subs, Language: Lang{Tag: newItem.Language},
		ContentCreated: newItem.Contentcreated, ContentModified: newItem.Versioncreated}

	if cm.Language.Tag == "" {
		cm.Language.Tag = "unknown"
	}

	/*if cm.ContentCreated == "" {
		cm.ContentCreated = defaultTime
	}*/

	if cm.ContentModified == "" {
		cm.ContentModified = defaultTime
	}

	// extract the content bodies
	// with the same method as toNewsMD
	var bods []Bodies
	title := ""
	for _, headline := range newItem.Headlines {
		if headline.Role == "main" {
			title = headline.Value
		}
	}
	contentBody := ""
	for _, body := range newItem.Bodies {
		if body.Role == "text" {
			contentBody = body.Value
		}
	}

	if contentBody == "" {
		for _, desc := range newItem.Descriptions {
			if desc.Role == "text" {
				contentBody = desc.Value
			}
		}
	}

	bods = append(bods, Bodies{
		BodyHead:    Head{Hedline: Hedline{Hl1: title}},
		BodyContent: Content{P: []string{contentBody}},
	})

	cs := ContentSet{InlineXML: Inline{Contenttype: "application/nitf+xml",
		Nitf: Nitf{Body: bods}}}

	out.RightsInfo = c
	out.ItemMeta = im
	out.ContentMeta = cm
	out.ContentSet = cs

	return xml.MarshalIndent(out, "", "    ")
}

// NewsItem is the Newsml struct
type NewsItem struct {
	XMLName         xml.Name    `xml:"newsItem"`
	Text            string      `xml:",chardata"`
	Xmlns           string      `xml:"xmlns,attr"` // "http://iptc.org/std/nar/2006-10-01/"
	GUID            string      `xml:"guid,attr"`
	Version         string      `xml:"version,attr"`
	Standard        string      `xml:"standard,attr"`
	Standardversion string      `xml:"standardversion,attr"`
	RightsInfo      RightsInfo  `xml:"rightsInfo"`
	ItemMeta        itemMeta    `xml:"itemMeta"`
	ContentMeta     Contentmeta `xml:"contentMeta"`
	ContentSet      ContentSet  `xml:"contentSet"`
}

// ContentSet field of NewsML
type ContentSet struct {
	Text      string `xml:",chardata"`
	InlineXML Inline `xml:"inlineXML"`
}

// Inline field of NewsML
type Inline struct {
	Text        string `xml:",chardata"`
	Contenttype string `xml:"contenttype,attr"`
	Nitf        Nitf   `xml:"nitf"`
}

// Nitf field of NewsML
type Nitf struct {
	Text string   `xml:",chardata"`
	Body []Bodies `xml:"body"`
}

// Bodies field of NewsML
type Bodies struct {
	Text        string  `xml:",chardata"`
	BodyHead    Head    `xml:"body.head"`
	BodyContent Content `xml:"body.content"`
}

// Content field of NewsML
type Content struct {
	Text string   `xml:",chardata"`
	P    []string `xml:"p"`
}

// Head field of NewsML
type Head struct {
	Text    string  `xml:",chardata"`
	Hedline Hedline `xml:"hedline"`
	/*Byline  struct {
		Text  string `xml:",chardata"`
		Byttl string `xml:"byttl"`
	} `xml:"byline"`*/
}

// Hedline field of NewsML
type Hedline struct {
	Text string `xml:",chardata"`
	Hl1  string `xml:"hl1"`
}

// Contentmeta field of NewsML
type Contentmeta struct {
	Text            string     `xml:",chardata"`
	ContentCreated  string     `xml:"contentCreated,omitempty"`
	ContentModified string     `xml:"contentModified"`
	Located         Located    `xml:"located"`
	Creator         Creator    `xml:"creator"`
	InfoSource      InfoSource `xml:"infoSource"`
	Language        Lang       `xml:"language"`
	Subject         []Subject  `xml:"subject"`
	Slugline        string     `xml:"slugline"`
	Headline        string     `xml:"headline"`
}

// Subject field of NewsML
type Subject struct {
	Text  string `xml:",chardata"`
	Qcode string `xml:"qcode,attr,omitempty"`
	Name  string `xml:"name"`
}

// Lang field of NewsML
type Lang struct {
	Text string `xml:",chardata"`
	Tag  string `xml:"tag,attr"`
}

// InfoSource field of NewsML
type InfoSource struct {
	Text  string `xml:",chardata"`
	Qcode string `xml:"qcode,attr,omitempty"`
	Name  string `xml:"name"`
}

// Creator field of NewsML
type Creator struct {
	Text string `xml:",chardata"`
	URI  string `xml:"uri,attr"`
	Name string `xml:"name"`
}

// Located field of NewsML
type Located struct {
	Text  string `xml:",chardata"`
	Qcode string `xml:"qcode,attr,omitempty"`
	Name  string `xml:"name"`
}

// RightsInfo field of NewsML
type RightsInfo struct {
	Text            string    `xml:",chardata"`
	CopyrightHolder copyright `xml:"copyrightHolder"`
	CopyrightNotice string    `xml:"copyrightNotice"`
}

// ItemMeta field of NewsML
type itemMeta struct {
	Text           string    `xml:",chardata"`
	ItemClass      ItemClass `xml:"itemClass"`
	Provider       Provider  `xml:"provider"`
	VersionCreated string    `xml:"versionCreated"`
	PubStatus      PubStatus `xml:"pubStatus"`
}

// ItemClass field of NewsML
type ItemClass struct {
	Text  string `xml:",chardata"`
	Qcode string `xml:"qcode,attr,omitempty"`
}

// Provider field of NewsML
type Provider struct {
	Text string `xml:",chardata"`
	URI  string `xml:"uri,attr"`
}

// PubStatus field of NewsML
type PubStatus struct {
	Text  string `xml:",chardata"`
	Qcode string `xml:"qcode,attr,omitempty"`
}

// Copyright field of NewsML
type copyright struct {
	Text string `xml:",chardata"`
	URI  string `xml:"uri,attr"`
	Name string `xml:"name"`
}
