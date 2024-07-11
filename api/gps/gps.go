// package battery handles all the gps style metadata and their transformations
package gps

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// ConvertGPX converts GPX coordinates into Json W3C coords
func ConvertGPX(input []byte, _ ...string) ([]byte, error) {
	var gpx Gpx
	xml.Unmarshal(input, &gpx)
	w3cPoints := make([]W3C, len(gpx.Wpt)+len(gpx.Trk.Trackseg.Trkpt))

	// handle waypoints and trackpoints in the same way
	for i, wpt := range gpx.Wpt {
		w3cPoints[i] = waypointToW3C(wpt)
	}

	for i, wpt := range gpx.Trk.Trackseg.Trkpt {
		w3cPoints[i+len(gpx.Wpt)] = waypointToW3C(wpt)
	}

	return json.MarshalIndent(w3cPoints, "", "    ")
}

const (
	timeLayout = "2006-01-02T15:04:05Z"
)

func waypointToW3C(w Waypoint) W3C {

	c := Coords{
		Latitude:  w.Lat,
		Longitude: w.Lon,
		Altitude:  w.Ele,
	}

	t := int64(0)

	if w.Time != "" {
		// parse the string into the golang format
		// so it can be converted to unix
		tmid, _ := time.Parse(timeLayout, w.Time)
		t = tmid.Unix()
	}
	return W3C{
		Coords:    c,
		Timestamp: t,
	}
}

// GPX xml layout structs
type Gpx struct {
	XMLName xml.Name   `xml:"gpx"`
	Text    string     `xml:",chardata"`
	Version string     `xml:"version,attr"`
	Name    string     `xml:"name"`
	Wpt     []Waypoint `xml:"wpt"`
	Trk     Track      `xml:"trk"`
}

// Waypoint is the gpx waypoint
type Waypoint struct {
	Text string  `xml:",chardata"`
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Ele  float64 `xml:"ele"`
	Time string  `xml:"time"`
	Name string  `xml:"name"`
	Sym  string  `xml:"sym"`
}

// Track is the GPX track
type Track struct {
	Text     string       `xml:",chardata"`
	Name     string       `xml:"name"`
	Number   string       `xml:"number"`
	Trackseg TrackSegment `xml:"trkseg"`
}

// TrackSegment is the GPX trkseg
type TrackSegment struct {
	Text  string     `xml:",chardata"`
	Trkpt []Waypoint `xml:"trkpt"`
}

// w3c json output structs
type W3C struct {
	Coords    Coords `json:"coords"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// Coords are the w3c json coordinates
type Coords struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Altitude         float64 `json:"altitude,omitempty"`
	Accuracy         float64 `json:"accuracy,omitempty"`
	AltitudeAccuracy float64 `json:"altitudeAccuracy,omitempty"`
	Heading          float64 `json:"heading,omitempty"`
	Speed            float64 `json:"speed,omitempty"`
}
