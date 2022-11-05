package genworldvoronoi

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Flokey82/naming-language-gen/naming"
	"github.com/s0rg/fantasyname"
)

func GenLanguage(seed int64) *Language {
	rand.Seed(seed)
	lang := naming.RandomLanguage(true, true)

	wordConfig := &naming.WordParams{
		MinSyllables: 1,
		MaxSyllables: 3,
		Structure:    naming.DefaultSyllableStructures,
	}

	nameConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 20),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(2, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "words",
	}

	cityConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 15),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(3, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "city",
	}

	firstNameConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 15),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(3, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "firstname",
	}

	lastNameConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 15),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(3, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "lastname",
	}

	return &Language{
		lang:            lang,
		wordConfig:      wordConfig,
		nameConfig:      nameConfig,
		cityConfig:      cityConfig,
		firstNameConfig: firstNameConfig,
		lastNameConfig:  lastNameConfig,
	}
}

type Language struct {
	lang            *naming.Language
	wordConfig      *naming.WordParams
	nameConfig      *naming.NameParams
	cityConfig      *naming.NameParams
	firstNameConfig *naming.NameParams
	lastNameConfig  *naming.NameParams
}

func (l *Language) GetWord(group string) string {
	return l.lang.GetWord(l.wordConfig, group)
}

func (l *Language) MakeName() string {
	return l.lang.MakeName(l.nameConfig)
}

func (l *Language) MakeCityName() string {
	return l.lang.MakeName(l.cityConfig)
}

func (l *Language) MakeFirstName() string {
	return l.lang.MakeName(l.firstNameConfig)
}

func (l *Language) MakeLastName() string {
	return l.lang.MakeName(l.lastNameConfig)
}

func newFantasyName() (fmt.Stringer, error) {
	return fantasyname.Compile("")
}

var largeAreaSuffix = []string{
	"land",
	"plains",
	"expanse",
	"region",
}

var fertileLandPrefix = []string{
	"green",
	"lush",
	"bountiful",
	"fruitful",
	"rich",
	"abundant",
	"fertile",
}

func NewDesertGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, []string{
		"charred",
		"desolate",
		"empty",
		"arid",
		"bleak",
		"scorched",
		"burnt",
		"forsaken",
	}, []string{
		"desert",
		"wasteland",
		"sands",
		"barrens",
		"expanse",
		"region",
	}, WordPair{
		A: []string{
			"burned",
			"scorched",
			"charred",
			"lost",
			"ashen",
		},
		B: DangerousSuffixB,
	})
}

var mountainPrefix = []string{
	"rocky",
	"mountainous",
	"spiked",
	"steep",
	"rough",
	"craggy",
	"toothy",
	"jagged",
	"broken",
}

var mountainDangerPrefix = []string{
	"shattered",
	"lost",
	"petrified",
	"forgotten",
	"abandoned",
	"broken",
	"buried",
	"fallen",
}

func NewMountainRangeGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, mountainPrefix, []string{
		"mountains",
		"teeth",
		"spikes",
		"peaks",
		"rocks",
		"thorns",
		"jags",
		"spurs",
		"spires",
		"pinnacles",
	}, WordPair{
		A: mountainDangerPrefix,
		B: DangerousSuffixB,
	})
}

func NewMountainGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, mountainPrefix, []string{
		"mountain",
		"tooth",
		"spike",
		"peak",
		"rock",
		"thorn",
		"jag",
		"spur",
		"spire",
		"pinnacle",
		"stone",
	}, WordPair{
		A: mountainDangerPrefix,
		B: DangerousSuffixB,
	})
}

// dangerousTerrainDescriptor returns a random descriptor for a dangerous terrain.
// This suffix is supposed to be used with "... of ...".
var dangerousTerrainDescriptor = []string{
	"death",
	"doom",
	"despair",
	"darkness",
	"evil",
	"chaos",
	"madness",
	"loss",
	"pain",
	"anguish",
	"terror",
	"horror",
	"lost souls",
	"the dead",
	"the damned",
	"the cursed",
	"the forsaken",
	"the lost",
	"the forgotten",
	"the abandoned",
	"the unknown",
	"the doomed",
}

var DangerousPrefixes = []string{
	"cursed",
	"forsaken",
	"forbidden",
	"lost",
	"abandoned",
	"unknown",
	"doomed",
	"haunted",
	"dark",
	"evil",
	"chaotic",
	"mad",
	"lost",
	"lost souls",
	"the dead",
	"the damned",
	"the cursed",
	"the forsaken",
	"the lost",
	"the forgotten",
	"the abandoned",
	"the unknown",
	"the doomed",
}

var DangerousSuffixB = []string{
	"hope",
	"souls",
	"fate",
	"dreams",
	"love",
	"life",
	"joy",
	"peace",
	"serenity",
	"calm",
	"tranquility",
	"corpses",
	"heroes",
	"villains",
	"princes",
	"princesses",
	"queens",
	"kings",
	"emperors",
	"empresses",
	"lords",
	"ladies",
	"knights",
	"lovers",
}

func NewForestGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, []string{
		"dark",
		"shadowy",
		"twisted",
		"broad",
		"thick",
		"dense",
		"overgrown",
		"lush",
		"green",
		"mossy",
		"moldy",
		"mold-covered",
		"leafy",
		"leaf-covered",
		"leaf-strewn",
		"leaf-littered",
		"woody",
		"wooded",
		"wood-strewn",
		"wood-littered",
		"wood-covered",
		"wooden",
		"bark-covered",
		"bark-strewn",
		"bark-littered",
		"barky",
	}, []string{
		"forest",
		"woods",
		"wood",
		"grove",
		"groves",
		"thicket",
		"thickets",
	}, WordPair{
		A: []string{
			"rotten",
			"lost",
			"petrified",
			"forgotten",
			"abandoned",
			"consumed",
			"buried",
		},
		B: DangerousSuffixB,
	})
}

func NewSwampGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, []string{
		"muddy",
		"mud-covered",
		"swampy",
		"marshy",
		"sticky",
		"humid",
		"muggy",
		"wet",
		"waterlogged",
		"water-covered",
		"moist",
		"mossy",
		"moldy",
	}, []string{
		"swamp",
		"marsh",
		"mire",
		"bog",
		"quagmire",
		"quag",
		"porridge",
		"muck",
		"mud",
		"phlegm",
		"sewer",
		"sludge",
		"puddles",
	}, WordPair{
		A: []string{
			"drowned",
			"lost",
			"sunken",
			"forgotten",
			"abandoned",
			"flushed",
			"buried",
		},
		B: DangerousSuffixB,
	})
}

type RiverGenerator struct {
	*namer
	Prefix          []string
	SuffixSmall     []string
	SuffixLarge     []string
	DangerousSuffix WordPair
}

func NewRiverGenerator(seed int64) *RiverGenerator {
	return &RiverGenerator{
		namer: newNamer(seed),
		Prefix: []string{
			"clear",
			"clean",
			"pure",
			"fresh",
			"cold",
			"cool",
			"bracing",
			"refreshing",
			"crisp",
			"fast",
			"swift",
			"running",
			"rushing",
			"flowing",
			"rippling",
			"lively",
			"snaking",
			"meandering",
			"bubbling",
			"sparkling",
			"glittering",
			"shimmering",
			"shining",
		},
		SuffixSmall: []string{
			"brook",
			"run",
			"brooklet",
			"waters",
			"flow",
			"wee",
		},
		SuffixLarge: []string{
			"river",
			"stream",
			"current",
			"torrent",
			"flow",
		},
		DangerousSuffix: WordPair{
			A: []string{
				"drowned",
				"lost",
				"sunken",
				"forgotten",
				"abandoned",
				"flushed",
			},
			B: DangerousSuffixB,
		},
	}
}

func (g *RiverGenerator) Generate(seed int64, small, dangerous bool) string {
	// TODO: Add chance to generate a custom name using fantasyname.
	g.resetToSeed(seed)
	prefix := g.randomString(g.Prefix)
	var suffix string
	if small {
		suffix = g.randomString(g.SuffixSmall)
	} else {
		suffix = g.randomString(g.SuffixLarge)
	}
	// There is a chance that we simply merge the two words.
	var name string
	if g.randomChance(0.5) && !dangerous {
		name = prefix + suffix
	} else {
		name = prefix + " " + suffix
	}
	if !dangerous {
		return name
	}
	return name + " of " + g.randomPair(g.DangerousSuffix.A, g.DangerousSuffix.B)
}

type BasicGenerator struct {
	*namer
	Prefix          []string
	Suffix          []string
	DangerousSuffix WordPair
}

func NewBasicGenerator(seed int64, prefix, suffix []string, danger WordPair) *BasicGenerator {
	return &BasicGenerator{
		namer:           newNamer(seed),
		Prefix:          prefix,
		Suffix:          suffix,
		DangerousSuffix: danger,
	}
}

func (g *BasicGenerator) Generate(seed int64, dangerous bool) string {
	// TODO: Add chance to generate a custom name using fantasyname.
	g.resetToSeed(seed)
	prefix := g.randomString(g.Prefix)
	suffix := g.randomString(g.Suffix)
	// There is a chance that we simply merge the two words.
	var name string
	if g.randomChance(0.5) && !dangerous {
		name = prefix + suffix
	} else {
		name = prefix + " " + suffix
	}
	if !dangerous {
		return name
	}
	return name + " of " + g.randomPair(g.DangerousSuffix.A, g.DangerousSuffix.B)
}

type namer struct {
	rand *rand.Rand
}

func newNamer(seed int64) *namer {
	return &namer{
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (n *namer) resetToSeed(seed int64) {
	n.rand.Seed(seed)
}

func (n *namer) randomPair(a, b []string) string {
	return n.randomString(a) + " " + n.randomString(b)
}

func (n *namer) randomString(list []string) string {
	return list[n.rand.Intn(len(list))]
}

func (n *namer) randomChance(chance float64) bool {
	return math.Abs(n.rand.Float64()) < chance
}

type WordPair struct {
	A []string
	B []string
}
