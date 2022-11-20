package genworldvoronoi

import (
	"math/rand"
	"regexp"
	"strings"
)

var approaches []string

func init() {
	approaches = weightedToArray(approach)
}

// Religion represents a religion in the world.
//
// TODO: Ensure we can infer symbolisms from events and other things.
//
// For example, if they worship the 99 beer bottles on the wall, we should
// be able to infer that they highly value beer and the number 99, as well
// as walls. They might be fearful of the number 100, and might have a
// taboo against the number 1.
// They might look kindly on people who can drink 99 beers in a row.
//
// Another example: If they worship the sun, we should be able to infer
// that they highly value the sun, and that they might be fearful of the
// moon. They might have a celebration during the summer solstice and consider
// a total eclipse of the sun to be a bad omen and a moon eclipse to be a good
// omen.
//
// # DEITIES AND SYMBOLS
//
// Folk religions that are purely based on the culture might warship
// nature itself, such as the sun, summer, the rain, a particular animal,
// or a particular plant. They might worship one or multiple deities that
// represent nature, like the sun god, the rain god, the god of the forest.
//
// Organized religions might either worship one or multiple gods, or a single
// person that is considered to be a god (or chosen).
//
// # GRAPH
//
// All these themes and connections could be represented as a graph, which
// would allow us to infer the relationships between deities and symbols and
// if mundane events hold any significance for a particular religion.
type Religion struct {
	ID           int
	Origin       int
	Name         string
	Culture      *Culture
	Type         string
	Form         string
	Deity        string
	Expansion    string
	Expansionism int
	Parent       *Religion
}

// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/religions-generator.js
func (m *Civ) genFolkReligion(c *Culture) *Religion {
	form := rw(forms[ReligionGroupFolk])
	r := &Religion{
		Origin:  c.ID,
		Name:    c.Name + " " + rw(types[form]),
		Culture: c,
		Type:    "Folk",
		Form:    form,
	}
	if form != ReligionFormAnimism {
		r.Deity = getDeityName(c)
	}
	m.Religions = append(m.Religions, r)
	return r
}

/*
func (m *Map) genReligionOrganized() []*Religion {
	var religions []*Religion
	cities := m.cities_r
	sort.Slice(cities, func(i, j int) bool {
		return cities[i].Score > cities[j].Score
	})
	for _, c := range cities {
		form := rw(forms[ReligionGroupOrganized])
		//const state = cells.state[center];
		culture := m.getRCulture(c.R)

		var deity string
		if form != ReligionFormNontheism {
			deity = getDeityName(culture)
		}
		name, expansion := getReligionName(form, deity, center)
		//if (expansion === "state" && !state) expansion = "global";
		//if (expansion === "culture" && !culture) expansion = "global";

		//if (expansion === "state" && Math.random() > 0.5) center = states[state].center;
		//if (expansion === "culture" && Math.random() > 0.5) center = cultures[culture].center;

		//if (!cells.burg[center] && cells.c[center].some(c => cells.burg[c]))
		//  center = cells.c[center].find(c => cells.burg[c]);
		//const [x, y] = cells.p[center];

		//const s = spacing * gauss(1, 0.3, 0.2, 2, 2); // randomize to make the placement not uniform
		//if (religionsTree.find(x, y, s) !== undefined) continue; // to close to existing religion

		// add "Old" to name of the folk religion on this culture
		//isFolkBased := expansion == "culture" || P(0.5)
		//folk := isFolkBased && religions.find(r => r.culture === culture && r.type === "Folk");
		//if (folk && expansion === "culture" && folk.name.slice(0, 3) !== "Old") folk.name = "Old " + folk.name;

		//const origins = folk ? [folk.i] : getReligionsInRadius({x, y, r: 150 / count, max: 2});
		//const expansionism = rand(3, 8);
		//const baseColor = religions[culture]?.color || states[state]?.color || getRandomColor();
		//const color = getMixedColor(baseColor, 0.3, 0);

		religions.push({
		  i: religions.length,
		  name,
		  color,
		  culture,
		  type: "Organized",
		  form,
		  deity,
		  expansion,
		  expansionism,
		  center,
		  origins
		});
	}
	return religions
}*/

func (m *Map) getReligionName(form, deity string, center int) (string, string) {
	//const {cells, cultures, burgs, states} = pack;

	c := m.Cultures[m.RegionToCulture[center]]
	random := func() string {
		return c.Language.MakeName()
	}
	rType := func() string {
		return rw(types[form])
	}
	deitySplit := regexp.MustCompile(`/[ ,]+/`)
	supreme := func() string {
		return deitySplit.Split(deity, -1)[0]
	}
	culture := func() string {
		return c.Name
	}
	place := func(adj string) string {
		/*	e:=m.r_empires[center]
			const burgId = cells.burg[center];
			const stateId = cells.state[center];

			const base = burgId ? burgs[burgId].name : states[stateId].name;
			let name = trimVowels(base.split(/[ ,]+/)[0]);
			return adj ? getAdjective(name) : name;*/
		return "TODO_PLACE"
	}

	switch rw(methods) {
	case MethodRandomType:
		return random() + " " + rType(), "global"
	case MethodRandomIsm:
		return trimVowels(random(), 3) + "ism", "global"
	case MethodSurpremeIsm:
		if deity != "" {
			return trimVowels(supreme(), 3) + "ism", "global"
		}
	case MethodFaithOfSupreme:
		if deity != "" {
			return ra([]string{"Faith", "Way", "Path", "Word", "Witnesses"}) + " of " + supreme(), "global"
		}
	case MethodPlaceIsm:
		return place("") + "ism", "state"
	case MethodCultureIsm:
		return trimVowels(culture(), 3) + "ism", "culture"
	case MethodPlaceIanType:
		return place("adj") + " " + rType(), "state"
	case MethodCultureType:
		return culture() + " " + rType(), "culture"
	}
	return trimVowels(random(), 3) + "ism", "global"
}

// chars that serve as vowels
const VOWELS = `aeiouyɑ'əøɛœæɶɒɨɪɔɐʊɤɯаоиеёэыуюяàèìòùỳẁȁȅȉȍȕáéíóúýẃőűâêîôûŷŵäëïöüÿẅãẽĩõũỹąęįǫųāēīōūȳăĕĭŏŭǎěǐǒǔȧėȯẏẇạẹịọụỵẉḛḭṵṳ`

func vowel(c rune) bool {
	return strings.IndexRune(VOWELS, c) != -1
}

// remove vowels from the end of the string
func trimVowels(str string, minLength int) string {
	for len(str) > minLength && vowel(rune(str[len(str)-1])) {
		str = str[:len(str)-1]
	}
	return str
}

// get supreme deity name
func getDeityName(culture *Culture) string {
	if culture == nil {
		return "ERROR"
	}
	meaning := generateMeaning()
	cultureName := culture.Language.MakeName() // Names.getCulture(culture, nil, nil, "", 0.8)
	return cultureName + ", The " + meaning
}

func generateMeaning() string {
	switch ra(approaches) { // select generation approach
	case ApproachNumber:
		return ra(base["number"])
	case ApproachBeing:
		return ra(base["being"])
	case ApproachAdjective:
		return ra(base["adjective"])
	case ApproachColorAnimal:
		return ra(base["color"]) + " " + ra(base["animal"])
	case ApproachAdjectiveAnimal:
		return ra(base["adjective"]) + " " + ra(base["animal"])
	case ApproachAdjectiveBeing:
		return ra(base["adjective"]) + " " + ra(base["being"])
	case ApproachAdjectiveGenitive:
		return ra(base["adjective"]) + " " + ra(base["being"])
	case ApproachColorBeing:
		return ra(base["color"]) + " " + ra(base["being"])
	case ApproachColorGenitive:
		return ra(base["color"]) + " " + ra(base["genitive"])
	case ApproachBeingOfGenitive:
		return ra(base["being"]) + " of " + ra(base["genitive"])
	case ApproachBeingOfTheGenitive:
		return ra(base["being"]) + " of the " + ra(base["genitive"])
	case ApproachAnimalOfGenitive:
		return ra(base["animal"]) + " of " + ra(base["genitive"])
	case ApproachAdjectiveBeingOfGenitive:
		return ra(base["adjective"]) + " " + ra(base["being"]) + " of " + ra(base["genitive"])
	case ApproachAdjectiveAnimalOfGenitive:
		return ra(base["adjective"]) + " " + ra(base["animal"]) + " of " + ra(base["genitive"])
	default:
		return "ERROR"
	}
}

// name generation approach and relative chance to be selected
var approach = map[string]int{
	ApproachNumber:                    1,
	ApproachBeing:                     3,
	ApproachAdjective:                 5,
	ApproachColorAnimal:               5,
	ApproachAdjectiveAnimal:           5,
	ApproachAdjectiveBeing:            5,
	ApproachAdjectiveGenitive:         1,
	ApproachColorBeing:                3,
	ApproachColorGenitive:             3,
	ApproachBeingOfGenitive:           2,
	ApproachBeingOfTheGenitive:        1,
	ApproachAnimalOfGenitive:          1,
	ApproachAdjectiveBeingOfGenitive:  2,
	ApproachAdjectiveAnimalOfGenitive: 2,
}

const (
	ApproachNumber                    = "Number"
	ApproachBeing                     = "Being"
	ApproachAdjective                 = "Adjective"
	ApproachColorAnimal               = "Color + Animal"
	ApproachAdjectiveAnimal           = "Adjective + Animal"
	ApproachAdjectiveBeing            = "Adjective + Being"
	ApproachAdjectiveGenitive         = "Adjective + Genitive"
	ApproachColorBeing                = "Color + Being"
	ApproachColorGenitive             = "Color + Genitive"
	ApproachBeingOfGenitive           = "Being + of + Genitive"
	ApproachBeingOfTheGenitive        = "Being + of the + Genitive"
	ApproachAnimalOfGenitive          = "Animal + of + Genitive"
	ApproachAdjectiveBeingOfGenitive  = "Adjective + Being + of + Genitive"
	ApproachAdjectiveAnimalOfGenitive = "Adjective + Animal + of + Genitive"
)

const (
	MethodRandomType     = "Random + type"
	MethodRandomIsm      = "Random + ism"
	MethodSurpremeIsm    = "Supreme + ism"
	MethodFaithOfSupreme = "Faith of + Supreme"
	MethodPlaceIsm       = "Place + ism"
	MethodCultureIsm     = "Culture + ism"
	MethodPlaceIanType   = "Place + ian + type"
	MethodCultureType    = "Culture + type"
)

var methods = map[string]int{
	MethodRandomType:     3,
	MethodRandomIsm:      1,
	MethodSurpremeIsm:    5,
	MethodFaithOfSupreme: 5,
	MethodPlaceIsm:       1,
	MethodCultureIsm:     2,
	MethodPlaceIanType:   6,
	MethodCultureType:    4,
}

func ra(array []string) string {
	return array[rand.Intn(len(array))]
}

func rw(mp map[string]int) string {
	return ra(weightedToArray(mp))
}

var base = map[string][]string{
	"number": {"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve"},
	"being": {
		"Ancestor",
		"Ancient",
		"Brother",
		"Chief",
		"Council",
		"Creator",
		"Deity",
		"Elder",
		"Father",
		"Forebear",
		"Forefather",
		"Giver",
		"God",
		"Goddess",
		"Guardian",
		"Lady",
		"Lord",
		"Maker",
		"Master",
		"Mother",
		"Numen",
		"Overlord",
		"Reaper",
		"Ruler",
		"Sister",
		"Spirit",
		"Virgin",
	},
	"animal": {
		"Antelope",
		"Ape",
		"Badger",
		"Basilisk",
		"Bear",
		"Beaver",
		"Bison",
		"Boar",
		"Buffalo",
		"Camel",
		"Cat",
		"Centaur",
		"Chimera",
		"Cobra",
		"Crane",
		"Crocodile",
		"Crow",
		"Cyclope",
		"Deer",
		"Dog",
		"Dragon",
		"Eagle",
		"Elk",
		"Falcon",
		"Fox",
		"Goat",
		"Goose",
		"Hare",
		"Hawk",
		"Heron",
		"Horse",
		"Hound",
		"Hyena",
		"Ibis",
		"Jackal",
		"Jaguar",
		"Kraken",
		"Lark",
		"Leopard",
		"Lion",
		"Mantis",
		"Marten",
		"Moose",
		"Mule",
		"Narwhal",
		"Owl",
		"Ox",
		"Panther",
		"Pegasus",
		"Phoenix",
		"Rat",
		"Raven",
		"Rook",
		"Scorpion",
		"Serpent",
		"Shark",
		"Sheep",
		"Snake",
		"Sphinx",
		"Spider",
		"Swan",
		"Tiger",
		"Turtle",
		"Unicorn",
		"Viper",
		"Vulture",
		"Walrus",
		"Wolf",
		"Wolverine",
		"Worm",
		"Wyvern",
	},
	"adjective": {
		"Aggressive",
		"Almighty",
		"Ancient",
		"Beautiful",
		"Benevolent",
		"Big",
		"Blind",
		"Blond",
		"Bloody",
		"Brave",
		"Broken",
		"Brutal",
		"Burning",
		"Calm",
		"Cheerful",
		"Crazy",
		"Cruel",
		"Dead",
		"Deadly",
		"Devastating",
		"Distant",
		"Disturbing",
		"Divine",
		"Dying",
		"Eternal",
		"Evil",
		"Explicit",
		"Fair",
		"Far",
		"Fat",
		"Fatal",
		"Favorable",
		"Flying",
		"Friendly",
		"Frozen",
		"Giant",
		"Good",
		"Grateful",
		"Great",
		"Happy",
		"High",
		"Holy",
		"Honest",
		"Huge",
		"Hungry",
		"Immutable",
		"Infallible",
		"Inherent",
		"Last",
		"Latter",
		"Lost",
		"Loud",
		"Lucky",
		"Mad",
		"Magical",
		"Main",
		"Major",
		"Marine",
		"Naval",
		"New",
		"Old",
		"Patient",
		"Peaceful",
		"Pregnant",
		"Prime",
		"Proud",
		"Pure",
		"Sacred",
		"Sad",
		"Scary",
		"Secret",
		"Selected",
		"Severe",
		"Silent",
		"Sleeping",
		"Slumbering",
		"Strong",
		"Sunny",
		"Superior",
		"Sustainable",
		"Troubled",
		"Unhappy",
		"Unknown",
		"Waking",
		"Wild",
		"Wise",
		"Worried",
		"Young",
	},
	"genitive": {
		"Cold",
		"Day",
		"Death",
		"Doom",
		"Fate",
		"Fire",
		"Fog",
		"Frost",
		"Gates",
		"Heaven",
		"Home",
		"Ice",
		"Justice",
		"Life",
		"Light",
		"Lightning",
		"Love",
		"Nature",
		"Night",
		"Pain",
		"Snow",
		"Springs",
		"Summer",
		"Thunder",
		"Time",
		"Victory",
		"War",
		"Winter",
	},
	"theGenitive": {
		"Abyss",
		"Blood",
		"Dawn",
		"Earth",
		"East",
		"Eclipse",
		"Fall",
		"Harvest",
		"Moon",
		"North",
		"Peak",
		"Rainbow",
		"Sea",
		"Sky",
		"South",
		"Stars",
		"Storm",
		"Sun",
		"Tree",
		"Underworld",
		"West",
		"Wild",
		"Word",
		"World",
	},
	"color": {
		"Amber",
		"Black",
		"Blue",
		"Bright",
		"Brown",
		"Dark",
		"Golden",
		"Green",
		"Grey",
		"Light",
		"Orange",
		"Pink",
		"Purple",
		"Red",
		"White",
		"Yellow",
	},
}

var types = map[string]map[string]int{
	ReligionFormShamanism: {
		"Beliefs":   3,
		"Shamanism": 2,
		"Spirits":   1,
	},
	ReligionFormAnimism: {
		"Spirits": 1,
		"Beliefs": 1,
	},
	ReligionFormAncestorWorship: {
		"Beliefs":     1,
		"Forefathers": 2,
		"Ancestors":   2,
	},
	ReligionFormPolytheism: {
		"Deities":  3,
		"Faith":    1,
		"Gods":     1,
		"Pantheon": 1,
	},
	ReligionFormDualism: {
		"Religion": 3,
		"Faith":    1,
		"Cult":     1,
	},
	ReligionFormMonotheism: {
		"Religion": 1,
		"Church":   1,
	},
	ReligionFormNontheism: {
		"Beliefs": 3,
		"Spirits": 1,
	},

	ReligionFormCult: {
		"Cult":    4,
		"Sect":    4,
		"Arcanum": 1,
		"Coterie": 1,
		"Order":   1,
		"Worship": 1,
	},
	ReligionFormDarkCult: {
		"Cult":      2,
		"Sect":      2,
		"Blasphemy": 1,
		"Circle":    1,
		"Coven":     1,
		"Idols":     1,
		"Occultism": 1,
	},

	ReligionFormHeresy: {
		"Heresy":      3,
		"Sect":        2,
		"Apostates":   1,
		"Brotherhood": 1,
		"Circle":      1,
		"Dissent":     1,
		"Dissenters":  1,
		"Iconoclasm":  1,
		"Schism":      1,
		"Society":     1,
	},
}

const (
	ReligionGroupFolk      = "Folk"
	ReligionGroupOrganized = "Organized"
	ReligionGroupCult      = "Cult"
	ReligionGroupHeresy    = "Heresy"
	ReligionFormShamanism  = "Shamanism"
	ReligionFormAnimism    = "Animism"
	//ReligionFormNature = "Nature"
	ReligionFormAncestorWorship = "Ancestor worship"
	ReligionFormPolytheism      = "Polytheism"
	ReligionFormDualism         = "Dualism"
	ReligionFormMonotheism      = "Monotheism"
	ReligionFormNontheism       = "Non-theism"
	ReligionFormCult            = "Cult"
	ReligionFormDarkCult        = "Dark Cult"
	ReligionFormHeresy          = "Heresy"
)

var forms = map[string]map[string]int{
	ReligionGroupFolk: {
		ReligionFormShamanism:       2,
		ReligionFormAnimism:         2,
		ReligionFormAncestorWorship: 1,
		ReligionFormPolytheism:      2,
	},
	ReligionGroupOrganized: {
		ReligionFormPolytheism: 5,
		ReligionFormDualism:    1,
		ReligionFormMonotheism: 4,
		ReligionFormNontheism:  1,
	},
	ReligionFormCult: {
		ReligionFormCult:     1,
		ReligionFormDarkCult: 1,
	},
	ReligionFormHeresy: {ReligionFormHeresy: 1},
}
