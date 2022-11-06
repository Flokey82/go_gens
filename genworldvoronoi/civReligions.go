package genworldvoronoi

import (
	"log"
	"math/rand"
)

var approaches []string

func init() {
	approaches = weightedToArray(approach)
}

// turn weighted array into simple array
func weightedToArray(weighted map[string]int) []string {
	var res []string
	for key, weight := range weighted {
		for j := 0; j < weight; j++ {
			res = append(res, key)
		}
	}
	return res
}

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
	// Parents   *Religion
}

// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/religions-generator.js
func (m *Map) genReligions() []*Religion {
	var religions []*Religion
	// For every culture, there is a folk religion.
	for _, c := range m.cultures_r {
		form := rw(forms[ReligionGroupFolk])
		name := c.Name + " " + rw(types[form])
		var deity string
		if form == ReligionFormAnimism {
			deity = ""
		} else {
			deity = getDeityName(c)
		}
		religions = append(religions, &Religion{
			Origin:  c.ID,
			Name:    name,
			Culture: c,
			Type:    "Folk",
			Form:    form,
			Deity:   deity,
		})
	}

	for _, r := range religions {
		log.Println(r.Name)
		log.Println(r.Deity)
	}
	return religions
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
	a := ra(approaches) // select generation approach
	switch a {
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
