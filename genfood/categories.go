package genfood

import "github.com/Flokey82/go_gens/genstory"

// Types of baked goods
// See: https://en.wikipedia.org/wiki/List_of_baked_goods
const (
	BakedGoodBiscuit   = "biscuit"
	BakedGoodBread     = "bread"
	BakedGoodCake      = "cake"
	BakedGoodCookie    = "cookie"
	BakedGoodCracker   = "cracker"
	BakedGoodCroissant = "croissant"
	BakedGoodDonut     = "donut"
	BakedGoodMuffin    = "muffin"
	BakedGoodPastry    = "pastry"
	BakedGoodPie       = "pie"
	BakedGoodPretzel   = "pretzel"
	BakedGoodRoll      = "roll"
	BakedGoodTart      = "tart"
	BakedGoodTorte     = "torte"
)

// Types of dishes.
// TODO: Find authoritative list of main dishes.
const (
	DishPudding  = "pudding"
	DishCreame   = "creame"
	DishSoup     = "soup"
	DishStew     = "stew"
	DishSalad    = "salad"
	DishSandwich = "sandwich"
	DishRoast    = "roast"
	DishPasta    = "pasta"
	DishSteak    = "steak"
	DishSausage  = "sausage"
	DishRibs     = "ribs"
	DishStuffed  = "stuffed"
)

// TODO: Find a way to derive genstory.Rules from some form of ruleset for each dish.

var SandwichRules = &genstory.Rules{
	Expansions: map[string][]string{
		"bread": {
			"toastbread",
			"[grain]-bread",
			"[color] bread",
		},
		"grain": {
			"barley",
			"wheat",
		},
		"color": {
			"white",
			"brown",
			"black",
			"red",
			"green",
			"blue",
		},
		"filling": {
			"[meat]-[vegetable]",
			"[vegetable]-[meat]",
			"[vegetable]-[cheese]",
			"[cheese]-[meat]",
		},
		"vegetable": {
			"artichoke",
			"asparagus",
			"aubergine",
			"beet",
		},
		"meat": {
			"beef",
			"chicken",
			"duck",
			"fish",
			"goat",
			"lamb",
			"pork",
			"turkey",
		},
		"cheese": {
			"butter cheese",
			"[cheese_animal] cheese",
		},
		"cheese_animal": {
			"cow",
			"pig",
			"sheep",
			"goat",
			"horse",
			"mouse",
		},
		"starts": {
			"[preparation] [filling] sandwich in [bread] and [special]",
		},
		"special": {
			"[treatment] [special_veg/vegetable]",
			"[treatment] [special_meat/meat]",
			"[treatment] [special_ingredient]",
		},
		"special_ingredient": {
			"sock",
			"tooth",
			"nail",
			"hair",
			"ear",
			"eye",
			"nose",
		},
		"treatment": {
			"gilded",
			"pickled",
			"stretched",
			"massaged",
			"fermented",
			"rotten",
			"burnt",
			"blessed",
			"poisoned",
			"spiced",
			"spiked",
			"digested",
		},
		"preparation": preparationMethods,
	},
	Start: "[starts]",
}
