package genreligion

import "github.com/Flokey82/go_gens/genlanguage"

var DeityMeaningApproaches []string

func init() {
	DeityMeaningApproaches = weightedToArray(GenMeaningApproaches)
}

// Deity represents a deity name.
type Deity struct {
	Name     string
	Meaning  string
	Approach string
}

// FullName returns the full name of the deity (including the meaning, if any).
func (d *Deity) FullName() string {
	if d == nil {
		return ""
	}
	if d.Meaning == "" {
		return d.Name
	}
	return d.Name + ", The " + d.Meaning
}

// GetDeity returns a deity name for the given culture.
// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/religions-generator.js
func (g *Generator) GetDeity(lang *genlanguage.Language, approach string) *Deity {
	if lang == nil {
		return nil
	}
	return &Deity{
		Name:     lang.MakeName(),
		Meaning:  g.GenerateDeityMeaning(approach),
		Approach: approach,
	}
}

// RandDeityMeaningGenMethod returns a random deity meaning generation method.
func (g *Generator) RandDeityGenMethod() string {
	return g.Ra(DeityMeaningApproaches)
}

// GenerateDeityMeaning generates a meaning for a deity name.
// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/religions-generator.js
func (g *Generator) GenerateDeityMeaning(approach string) string {
	switch approach { // select generation approach
	case ApproachNumber:
		return g.Ra(GenBase[GenBaseNumber])
	case ApproachBeing:
		return g.Ra(GenBase[GenBaseBeing])
	case ApproachAdjective:
		return g.Ra(GenBase[GenBaseAdjective])
	case ApproachColorAnimal:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseAnimal])
	case ApproachAdjectiveAnimal:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseAnimal])
	case ApproachAdjectiveBeing:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseBeing])
	case ApproachAdjectiveGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachColorBeing:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseBeing])
	case ApproachColorGenitive:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachBeingOfGenitive:
		return g.Ra(GenBase[GenBaseBeing]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachBeingOfTheGenitive:
		return g.Ra(GenBase[GenBaseBeing]) + " of the " + g.Ra(GenBase[GenBaseTheGenitive])
	case ApproachAnimalOfGenitive:
		return g.Ra(GenBase[GenBaseAnimal]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachAdjectiveBeingOfGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseBeing]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachAdjectiveAnimalOfGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseAnimal]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	default:
		return "ERROR"
	}
}

// GenMeaningApproaches contains a map of name generation
// approaches and their relative chance to be selected.
var GenMeaningApproaches = map[string]int{
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
	GenBaseNumber      = "number"
	GenBaseBeing       = "being"
	GenBaseAnimal      = "animal"
	GenBaseAdjective   = "adjective"
	GenBaseGenitive    = "genitive"
	GenBaseTheGenitive = "theGenitive"
	GenBaseColor       = "color"
)

// GenBase contains a map of word lists used for name generation.
// TODO: Group individual entries by logical categories.
// So we can build up a pantheon of gods, each associated with different domains.
// For example:
// {North, South, East, West} -> {Direction}
// {Bride, Groom, Widow, Widower, Wife, Husband} -> {Marriage}
// {Giver, Taker, Destroyer, Creator, Maker, Breaker} -> {Action}
// {Sky, Earth, Water, Fire, Air, Spirit} -> {Elements}
// {Light, Dark, Bright, Shining, Shadow, Darkness} -> {Light}
var GenBase = map[string][]string{
	GenBaseNumber: {"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve"},
	GenBaseBeing: {
		"Ancestor",
		"Ancient",
		"Angel",
		"Betrayer",
		"Bride",
		"Brother",
		"Chief",
		"Child",
		"Council",
		"Creator",
		"Deity",
		"Elder",
		"Father",
		"Forebearer",
		"Forefather",
		"Foremother",
		"Forgiven",
		"Forgotten",
		"Giver",
		"God",
		"Goddess",
		"Golem",
		"Groom",
		"Guardian",
		"Guide",
		"Keeper",
		"King",
		"Lady",
		"Lord",
		"Lover",
		"Maker",
		"Master",
		"Mistress",
		"Mother",
		"Numen",
		"Orphan",
		"Overlord",
		"Reaper",
		"Ruler",
		"Seducer",
		"Seductress",
		"Servant",
		"Sister",
		"Spirit",
		"Virgin",
		"Warrior",
		"Watcher",
		"Widow",
		"Widower",
		"Wife",
		"Witch",
		"Wizard",
	},
	GenBaseAnimal: {
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
	GenBaseAdjective: {
		"Aggressive",
		"Almighty",
		"Ancient",
		"Angry",
		"Anxious",
		"Awful",
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
		"Careful",
		"Charming",
		"Cheerful",
		"Chosen",
		"Clean",
		"Crazy",
		"Cruel",
		"Dead",
		"Deadly",
		"Deaf",
		"Deathless",
		"Deep",
		"Defiant",
		"Delicate",
		"Delightful",
		"Desperate",
		"Devastating",
		"Distant",
		"Disturbing",
		"Divine",
		"Dying",
		"Enchanting",
		"Ephemeral",
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
		"Greedy",
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
		"Lonely",
		"Lost",
		"Loud",
		"Loving",
		"Lucky",
		"Lustful",
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
		"Spiteful",
		"Strong",
		"Sunny",
		"Superior",
		"Sustainable",
		"Troubled",
		"Undying",
		"Unhappy",
		"Unknown",
		"Waking",
		"Wild",
		"Wise",
		"Worried",
		"Young",
	},
	GenBaseGenitive: {
		"Cold",
		"Darkness",
		"Dawn",
		"Day",
		"Death",
		"Doom",
		"Dreams",
		"Dusk",
		"Fate",
		"Fire",
		"Fog",
		"Frost",
		"Gates",
		"Heaven",
		"Home",
		"Hope",
		"Ice",
		"Justice",
		"Kings",
		"Life",
		"Light",
		"Lightning",
		"Love",
		"Nature",
		"Night",
		"Pain",
		"Snow",
		"Springs",
		"Stars",
		"Summer",
		"Sun",
		"Sunset",
		"Thunder",
		"Time",
		"Victory",
		"War",
		"Wealth",
		"Winter",
		"Wisdom",
	},
	GenBaseTheGenitive: {
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
	GenBaseColor: {
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
