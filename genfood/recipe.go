package genfood

import (
	"github.com/Flokey82/go_gens/genstory"
)

var ExampleRules = &genstory.Rules{
	Expansions: map[string][]string{
		"prep": {
			"finely chop",
			"mince",
			"grate",
			"slice",
			"chop",
		},
		"age": {
			"old",
			"shriveled",
			"dated",
			"smelly",
			"fresh",
			"young",
		},
		"veg": vegetables,
		"fire": {
			"get a fire going",
			"stoke your hearth",
			"ensure your hearth is lit",
		},
		"serve_timing": {
			"on a new moon",
			"on a full moon",
			"during a solar eclipse",
			"during a lunar eclipse",
			"on a birthday",
			"on a wedding day",
			"on a funeral day",
			"on a Thursday",
			"on a Friday",
		},
		"wrapping": {
			"autum leaves",
			"leaves",
			"bark",
			"chicken skin",
			"wedding dress",
		},
		"serve_ritual": {
			"light a candle",
			"ring a bell",
			"stomp your feet",
			"clap your hands",
			"whistle",
		},
		"serve_side": {
			"pretzles",
			"bread",
			"wine",
			"mud",
			"mud and wine",
			"mud and pretzles",
			"mud and bread",
			"blue cheese",
		},
		"serve": {
			"Wrap in [wrapping] and serve [serve_timing].",
			"Fill in an old shoe and [serve_ritual] when [serve_timing].",
			"Serve in a bowl and cover in [wrapping].",
			"Spread on a plate with mud, serve with [serve_side].",
			"Serve in a bowl with a spoon, a fork and a knife.",
		},
		"preparation": {
			"Prepare a pot filled with water and bring to a boil. Add the [prep_veg/prep:past] [veg_veg/veg] and [boil].",
			"Take a pan and add the [prep_veg/prep:past] [veg_veg/veg] and [fry].",
		},
		"boil": {
			"boil [duration]",
			"stew [duration]",
			"simmer [duration]",
		},
		"fry": {
			"fry [duration]",
			"roast [duration]",
		},
		"duration": {
			"for [duration_time]",
			"over night",
			"unil blackened",
			"and stop when discolored",
		},
		"duration_time": {
			"10 minutes",
			"30 minutes",
			"1 hour",
			"2 hours",
			"3 hours",
			"4 hours",
			"5 hours",
			"9 days",
			"10 days",
		},
		"meat": meats,
		"strange_ingredient": {
			"goose feathers",
			"dragon scales",
			"unicorn horn",
			"troll cheese",
		},
	},
	Start: "[prep_veg/prep:capitalize] [age_veg/age:a] [veg_veg/veg] and [fire]. [preparation] Add [prep:past] [meat_meat/meat] and [strange_ingredient]. [serve]",
}
