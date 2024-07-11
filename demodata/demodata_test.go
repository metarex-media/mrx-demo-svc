package demodata

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

var headlines map[int][]string

type headlinesJson struct {
	Headlines []string `json:"headlines"`
}

func init() {

	// slurp up the headlines
	// to speed up json to string hardcoding
	headlines = make(map[int][]string)
	results := []string{beach, fire, dino, fishing}

	for i, r := range results {
		var h headlinesJson
		json.Unmarshal([]byte(r), &h)
		headlines[i] = h.Headlines
	}

}

func TestServices(t *testing.T) {

	MakeTS(regDetails{
		register: map[string]register{"json": {MrxID: "MRX.123.456.789.njs", outputs: []string{"toNewsMD"}}}},
		"demo01")

	MakeTS(regDetails{
		register: map[string]register{"gpx": {MrxID: "MRX.123.456.789.gpx", outputs: []string{"toW3C"}},
			"wav": {MrxID: "MRX.123.456.789.wav", outputs: []string{"ToWaveform"}}}},
		"demo03")

	MakeTS(regDetails{
		register: map[string]register{
			"xml": {MrxID: "MRX.123.456.789.vqc", outputs: []string{"tograph"}}}},
		"demo05")

	MakeTS(regDetails{
		register: map[string]register{
			"json": {MrxID: "MRX.123.456.789.ghj", outputs: []string{"MRX.123.456.789.cxm"}}}},
		"demo06")

	MakeTS(regDetails{
		register: map[string]register{"json": {MrxID: "MRX.123.456.789.bat", outputs: []string{"fill", "fault", "stagger"}}}},
		"demo07")

	MakeTS(regDetails{
		register: map[string]register{"jpeg": {MrxID: "MRX.123.456.789.jpg", outputs: []string{"extractHeaderc2pa"}},
			"jpg": {MrxID: "MRX.123.456.789.jpg", outputs: []string{"extractHeaderc2pa"}}}},
		"demo08")

	MakeTS(regDetails{
		register: map[string]register{"json": {MrxID: "MRX.123.456.789.njs", outputs: []string{"toNewsML"}}}},
		"demo10")

	MakeTS(regDetails{
		genericFields: map[string]string{"mapping": "true"},
		register: map[string]register{
			"yaml": {MrxID: "MRX.123.456.789.rny", outputs: []string{"MRX.123.456.789.rnf"}},
			"csv":  {MrxID: "MRX.123.456.789.rnc", outputs: []string{"MRX.123.456.789.rnf"}},
			"json": {MrxID: "MRX.123.456.789.rnj", outputs: []string{"MRX.123.456.789.rnf"}}}},
		"demo04")

	MakeTS(regDetails{
		custom: map[string]map[string]string{
			"IET":         {"title": "IET"},
			"springwatch": {"title": "SW"},
			"lostpast":    {"title": " LP"},
			"laundromat":  {"title": "LM"},
		},
		register: map[string]register{"csv": {MrxID: "MRX.123.456.789.rnf", outputs: []string{"generateFFmpeg"}}}},
		"demo12")
}

var (
	header       = `import type { DemoSource, MrxETLService, MrxRegisterCache } from '$lib/nabTypes';`
	regList      = `{mrxId: "%v", reg: %v},`
	regEntries   = `import %v from '$lib/reg/%v/register.json';`
	fetchlist    = `{type: "%v", id: "%v", url: "%s/%v" },`
	postRequests = `{type: "%v", id:"%v", post: {inputMRXID:"%v", data: "%v", outputMRXID:"%v" %s}},`
	id           = `{
		id: "%v",
		mrxId: "%v",
		name: "%v",
		clip: "",
		srcUrl: "%v/%v",
		summary: "",
		srcReg: "%v",
		xtra: [],
	  }`

	idalt = `{
		id: "%v",
		mrxId: "%v",
		name: "%v",
		clip: "",
		srcUrl: "%v/%v",
		summary: "",
		srcReg: "%v",
		xtra: newsHeaders,
	  }`
	// nheader is only used with the headlines
	nheader = `{time:"%02d:%02d", src:%v, headline:"%s", srcID:""},`
)

// details of each register
type regDetails struct {
	// generic fields are for adding parameters that everything
	// uses. e.g. mapping="true"
	genericFields map[string]string
	// custom relates to each file name
	// and any specific fields each one has. Like
	// demo 12 each having a separate title
	custom map[string]map[string]string
	// eacj file extensions has a metarex ID for these demos
	register map[string]register
}

type register struct {
	MrxID string
	// output formats
	outputs []string
}

/*

let taot = [
  {time:"09:00", src:0, headline:"rexy gets sunburn", srcID:""}
  {time:"09:01", src:1, headline:"rexy gets bitten by mosquitos",}
  ...
  {time:"17:59"...}
]
*/

func MakeTS(reg regDetails, demoNumber string) {

	// set up imports
	config := header + "\n\n"

	// import all the register jsons
	for _, r := range reg.register {
		ends := strings.Split(r.MrxID, ".")

		config += fmt.Sprintf(regEntries, ends[len(ends)-1], r.MrxID)
		config += "\n"
	}
	config += "\n"

	config += "export const fetchList = [\n"
	// import the data files
	inputs, _ := os.ReadDir("./" + demoNumber)
	for _, input := range inputs {
		if input.Name() != "readme.md" && input.Name() != "config.ts" {
			inner := strings.ReplaceAll(input.Name(), " ", "")
			inner = strings.ReplaceAll(inner, "-", "")
			in := strings.Split(inner, ".")

			config += fmt.Sprintf(fetchlist, in[1], in[0], demoNumber, input.Name())
			config += "\n"
		}
	}
	config += "]\n\n"

	// set up the register cache
	config += "export const regCache:MrxRegisterCache = [\n"
	for _, r := range reg.register {
		ends := strings.Split(r.MrxID, ".")

		config += fmt.Sprintf(regList, r.MrxID, ends[len(ends)-1])
		config += "\n"
	}
	config += "]\n"

	// export the data cache
	config += "\nexport const dataCache = [\n"
	genericTag := ""
	for gkey, gfield := range reg.genericFields {
		genericTag += fmt.Sprintf(`,%s : "%s"`, gkey, gfield)
	}

	// generate the demo sources
	// checking each input file
	for _, input := range inputs {
		if input.Name() != "readme.md" && input.Name() != "config.ts" {
			inner := strings.ReplaceAll(input.Name(), " ", "")
			inner = strings.ReplaceAll(inner, "-", "")
			in := strings.Split(inner, ".")

			outs := reg.register[in[1]]
			customTag := ""
			inName := strings.Split(input.Name(), ".")
			for ckey, cfield := range reg.custom[inName[0]] {
				customTag += fmt.Sprintf(`,%s : "%s"`, ckey, cfield)
			}

			for _, out := range outs.outputs {

				config += fmt.Sprintf(postRequests, in[1], in[0], outs.MrxID, in[0], out, genericTag+customTag)
				config += "\n"
			}
		}
	}
	config += "]\n\n"

	// set up specific news headlines for demo 01
	if demoNumber == "demo01" {
		config += "let newsHeaders = [\n"
		i := 0
		headerI := 0
		for h := 9; h <= 17; h++ {
			for m := 0; m <= 59; m++ {
				config += fmt.Sprintf(nheader, h, m, i%4, headlines[i%4][headerI%20])
				config += "\n"
				i++
				if i%4 == 0 {
					headerI++
				}
			}
		}
		config += "]\n\n"
	}

	config += "export const sources: DemoSource[] = [\n"
	for _, input := range inputs {
		if input.Name() != "readme.md" && input.Name() != "config.ts" {
			inner := strings.ReplaceAll(input.Name(), " ", "")
			inner = strings.ReplaceAll(inner, "-", "")
			in := strings.Split(inner, ".")
			outs := reg.register[in[1]]
			if demoNumber != "demo01" {
				config += fmt.Sprintf(id, in[0], outs.MrxID, input.Name(), demoNumber, input.Name(), in[1])
			} else {
				config += fmt.Sprintf(idalt, in[0], outs.MrxID, input.Name(), demoNumber, input.Name(), in[1])
			}
			config += ",\n"
		}
	}
	config += "]\n"

	// fmt.Println(config)

	f, _ := os.Create(fmt.Sprintf("./%v/config.ts", demoNumber))
	f.Write([]byte(config))
}

// Headlines for Demo 01
var (
	dino  = `{ "headlines": [ "Rexy Roars Back: The Unstoppable T. rex Takes Center Stage", "Meet Rexy: The Legendary Dinosaur of MT Metarex", "Rexy's Return: Jurassic World's Iconic T. rex Makes a Comeback", "The Queen of Metadata: Rexy's Reign Continues", "Tyrant No More: Rexy's Surprising Alliance", "Rexy Unleashed: Chaos Erupts as the T. rex Escapes", "Rexy vs. Indominus: Epic Battle of the Titans", "Rexy's Secret Past: Unraveling the Mystery", "Rexy's Legacy: How One Dinosaur Changed History", "Rexy's Rescue Mission: Saving Lives Amidst Chaos", "The Enigma of Rexy: Scientists Ponder Her Origins", "Rexy's Roar Heard 'Round the World: Global Phenomenon", "Rexy's Last Stand: Facing the Ultimate Threat", "Rexy's Hidden Lair: Uncharted Territory Revealed", "Rexy's Bond with Blue: Unlikely Friendship", "Rexy's Jurassic Odyssey: From Park to World", "Rexy's Ancient Echo: Tracing Her Ancestry", "Rexy's Dino-Detective: Solving Metadata Mysteries", "Rexy's Heart of Stone: Beneath the Scales", "Rexy's Roar Resounds: A Legend Reborn"] }`
	beach = `{ "headlines": [ "Rexy's Seaside Stroll: A Dino's Day at the Beach!", "Sun, Sand, and Scales: Rexy's Beach Vacation Chronicles", "Tyrannosaurus Tides: Rexy's Surfing Adventure", "Footprints in the Sand: Rexy's Beachside Mystery", "Rexy's Shell-Seeking Safari: Unearthing Ocean Treasures", "Dino Delight: Rexy's Beach Picnic Extravaganza", "Sunset Serenade: Rexy's Beach Bonfire Jam", "Rexy's Sandcastle Showdown: T. rex vs. Crabs!", "Beach Breeze and Dino Squeeze: Rexy's Summer Romance", "Rexy's Beach Volleyball Victory: Spike, Roar, Repeat!", "Surf's Up, Tails Out: Rexy Hangs Ten", "Rexy's Beach Yoga: Finding Inner Peace (and Prey)", "Seagulls vs. Rexy: The Great Beach Food Fight", "Tyrannosaurus Tanning: Rexy's Sunbathing Secrets", "Beachcomber Rexy: Collecting Fossilized Memories", "Rexy's Sand-Sational Sandcastle Sculptures", "Dino Dives: Rexy Explores Underwater Wonders", "Beachside Dino Dance Party: Rexy's Groovy Moves", "Rexy's Beach Cleanup: Saving the Shoreline, One Claw at a Time", "Sunrise to Sunset: Rexy's 24-Hour Beach Adventure"] }`
	fire  = `{
    "headlines": [
        "Rexy's Roaring Return: The Bonfire Chronicles",
         "Warm Nights with Rexy: Fireside Tales of Adventure",
         "Embers and Echoes: Rexy's Fiery Saga",
         "Tyrannosaurus Twilight: Cozy Campfires with Rexy",
         "Roberta's Hearthside Whispers: A Dino's Tale",
         "Flames of Friendship: Rexy's Bonfire Bonding",
         "Jurassic Glow: Rexy's Campfire Chronicles",
         "Dino Dreams by Firelight: Rexy's Nighttime Musings",
         "Rexy's Warmth: Legends Ignite Around the Bonfire",
        "Flickering Flames, Fierce Heart: Rexy's Story",
         "Tales from the T. rex: Rexy's Cozy Campfire Adventures",
         "Bonfire Ballads: Rexy's Roaring Refrain",
         "Fire-Kissed Memories: Rexy's Nostalgic Nights",
         "Rexy's Ember Quest: Legends Unfold in the Glow",
         "Dinosaur Dusk: Rexy's Fireside Fables",
         "Warm Claws, Bright Flames: Rexy's Bonfire Legacy",
         "Rexy's Hearthside Harmony: Songs of the Wild",
         "Tyrannosaurus Tales: Rexy's Cozy Campfire Chronicles",
         "Flame-Kissed Whispers: Rexy's Nighttime Secrets",
        "Bonfire Bonds: Rexy's Legendary Encounters"
    ]
}`
	fishing = `{ "headlines": [ "Rexy's Reel Adventures: Tackling the High Seas", "Tyrannosaurus Tides: Rexy's Fishing Expedition", "Hook, Line, and Roar: Rexy's Big Catch", "Dino Angler Rexy: Mastering the Art of Fishing", "Rexy's Underwater Quest: Unraveling Fishy Mysteries", "Scale Tales: Rexy's Fisherman Chronicles", "Rexy's Fishing Frenzy: From Ponds to Oceans", "T. rex on the Line: Rexy's Fishing Tournament Triumph", "Baited Breath: Rexy's Deep-Sea Fishing Adventure", "Rexy's Rod and Claw: A Dinosaur's Fishing Journey", "Splash and Snag: Rexy's Riverbank Fishing Escapade", "Reel or No Reel: Rexy's Fishing Dilemmas", "Rexy's Lure Lore: Secrets of Successful Fishing", "Fossilized Fish: Rexy's Ancient Angling Techniques", "Rexy's Catch of the Day: A Dino's Seafood Feast", "Gone Fishingosaur: Rexy's Lakeside Retreat", "Tyrannosaurus Trout: Rexy's Freshwater Adventures", "Rexy's Ocean Odyssey: Navigating the Waves", "Dino-Netting: Rexy's Sustainable Fishing Practices", "Rexy's Fish Whisperer: Communicating with Aquatic Creatures" ] }`
)
