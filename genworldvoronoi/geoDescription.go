package genworldvoronoi

import (
	"math"
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/genbiome"
)

func (m *Geo) GenerateRegPropertyDescription(p RegProperty) string {
	str := "The region"
	if p.OnIsland {
		// TODO: Note if the region is located on an island or peninsula.
		str += ", located on an island,"
	}
	str += " is covered by " + genbiome.WhittakerModBiomeToString(p.Biome) + ".\n"

	// Add info on the potential dangers of the region.
	if p.DistanceToVolcano < 3 {
		if p.DangerVolcano > 0.2 {
			if p.DistanceToVolcano < 1 {
				str += " It's location on a volcano"
			} else {
				str += " The proximity to a volcano"
			}
			if p.DangerVolcano > 0.5 {
				str += " results in a constant danger of destruction by a volcanic eruption"
			} else {
				str += " poses a looming threat of a possible volcanic eruption"
			}
			str += ". "
		} else {
			if p.DistanceToVolcano < 1 {
				str += " It is located on a volcano"
			} else {
				str += " It is close to a volcano"
			}
			str += ". "
		}
	} else if p.DistanceToMountain < 3 {
		if p.DangerRockslide > 0.2 {
			if p.DistanceToMountain < 1 {
				str += " The exposed location on a mountain"
			} else {
				str += " The proximity to a mountain"
			}
			if p.DangerRockslide > 0.5 {
				str += " poses a constant danger of a deadly rockslide"
			} else {
				str += " results in occasional rockslides that threaten anyone nearby"
			}
			str += ". "
		}
	} else if p.DistanceToFaultline < 3 {
		if p.DangerEarthquake > 0.2 {
			if p.DistanceToFaultline < 1 {
				str += " The exposed location on a faultline"
			} else {
				str += " The proximity to a faultline"
			}
			if p.DangerEarthquake > 0.5 {
				str += " results in a constant danger of a deadly earthquake"
			} else {
				str += " poses a looming threat of a possible earthquake"
			}
			str += ". "
		}
	}
	if p.DistanceToRiver < 1 {
		str += " The nearby river provides access to fresh water"
		if p.DangerFlood > 0.2 {
			if p.DangerFlood > 0.5 {
				str += " and is infamous for frequent floods"
			} else {
				str += " and might cause occasional flooding"
			}
		}
		str += ". \n"
	}
	return str
}

type BiomeDescription struct {
	Adjectives []string
	Nouns      []string
	Part1      []string
	Part2      []string
	Part3      []string
	Part4      []string
}

func GenerateFlavorText(desc BiomeDescription) string {
	text := "The " + desc.Adjectives[rand.Intn(len(desc.Adjectives))] + " " + desc.Nouns[rand.Intn(len(desc.Nouns))] + " stretches out as far as the eye can see.\n"
	if len(desc.Part1) > 0 {
		text += desc.Part1[rand.Intn(len(desc.Part1))] + " \n"
	}
	if len(desc.Part2) > 0 {
		text += desc.Part2[rand.Intn(len(desc.Part2))] + " \n"
	}
	if len(desc.Part3) > 0 {
		text += desc.Part3[rand.Intn(len(desc.Part3))] + " \n"
	}
	if len(desc.Part4) > 0 {
		text += desc.Part4[rand.Intn(len(desc.Part4))] + " \n"
	}
	return text
}

func generateFlavorTextForBiome(seed int64, biome int) string {
	rand.Seed(seed)
	switch biome {
	case genbiome.WhittakerModBiomeSubtropicalDesert:
		return GenerateFlavorText(desertDescription)
	case genbiome.WhittakerModBiomeColdDesert:
		return GenerateFlavorText(coldDesertDescription)
	case genbiome.WhittakerModBiomeTropicalRainforest:
		return GenerateFlavorText(tropicalRainforestDescription)
	case genbiome.WhittakerModBiomeTropicalSeasonalForest:
		return GenerateFlavorText(tropicalSeasonalForestDescription)
	case genbiome.WhittakerModBiomeTemperateRainforest:
		return GenerateFlavorText(temperateRainforestDescription)
	case genbiome.WhittakerModBiomeTemperateSeasonalForest:
		return GenerateFlavorText(temperateSeasonalForestDescription)
	case genbiome.WhittakerModBiomeWoodlandShrubland:
		return GenerateFlavorText(shrublandDescription)
	case genbiome.WhittakerModBiomeBorealForestTaiga:
		return GenerateFlavorText(borealForestDescription)
	case genbiome.WhittakerModBiomeTundra:
		return GenerateFlavorText(tundraDescription)
	case genbiome.WhittakerModBiomeHotSwamp:
		return GenerateFlavorText(hotSwampDescription)
	case genbiome.WhittakerModBiomeWetlands:
		return GenerateFlavorText(temperateWetlandDescription)
	case genbiome.WhittakerModBiomeSavannah:
		return GenerateFlavorText(savannahDescription)
	case genbiome.WhittakerModBiomeSnow:
		return GenerateFlavorText(snowDescription)
	default:
		return "The " + genbiome.WhittakerModBiomeToString(biome) + " stretches out as far as the eye can see."
	}
}

var desertDescription = BiomeDescription{
	Adjectives: []string{"arid", "bleak", "dry", "hot", "scorching", "sizzling", "sunny", "torrid"},
	Nouns:      []string{"desert", "dune", "oasis", "sand", "sand dune", "sandstorm", "sands", "wasteland"},
	Part1: []string{
		"The hot sun beats down relentlessly, baking the earth and sapping the energy from all who dare to venture into the wilderness.",
		"The dry air is thin and crisp, with a biting wind that cuts through clothing and chaps the skin.",
	},
	Part2: []string{
		"Mirages dance on the horizon, tempting travelers with the promise of water and shade.",
		"The only respite from the searing heat is found in the scarce oases that dot the landscape, offering a brief moment of coolness before the journey continues.",
	},
	Part3: []string{
		"The subtropical desert is home to a diverse array of plants and animals, including cacti, small mammals, and reptiles.",
		"The subtropical desert is a harsh and unforgiving place, with only the most resilient species able to survive the extreme conditions.",
	},
}

var coldDesertDescription = BiomeDescription{
	Adjectives: []string{"arid", "bleak", "cold", "frigid", "frozen", "icy", "snowy", "windy"},
	Nouns:      []string{"desert", "dune", "tundra", "wasteland"},
	Part1: []string{
		"The ground is dry and rocky, with patches of scrubby grass and scattered clusters of bushes and cacti dotting the landscape.",
		"The ground is cracked and parched, with only a few hardy plants able to survive the harsh conditions.",
	},
	Part2: []string{
		"The air is thin and biting, making it hard to catch one's breath.",
		"The frigid landscape is a stark and barren wasteland, with barely a sign of life for miles.",
	},
	Part3: []string{
		"The cold desert is home to a small array of plants and animals, including succulents, small mammals, and reptiles.",
		"The cold desert is a harsh and unforgiving place, with only the most resilient species able to survive the extreme conditions.",
	},
}

var temperateWetlandDescription = BiomeDescription{
	Adjectives: []string{"damp", "humid", "marshy", "moist", "muddy", "soggy", "swampy"},
	Nouns:      []string{"marsh", "swamp", "wetland", "landscape"},
	Part1: []string{
		"The air is heavy with the smell of mud and decay, as waterlogged plants and animals rot in the sun.",
		"The ground is soft and spongy, giving way beneath the feet with each step.",
	},
	Part2: []string{
		"The sky is overcast and gray, with a constant drizzle that keeps everything damp and humid.",
		"The sun shines brightly, illuminating the lush greenery and sparkling waters of the wetlands.",
	},
	Part3: []string{
		"The wetlands are home to a diverse array of plants and animals, including tall reeds, swaying grasses, and brightly-colored flowers.",
		"The wetlands are teeming with life, from the smallest insects to the largest predators.",
	},
}

var hotSwampDescription = BiomeDescription{
	Adjectives: []string{"damp", "humid", "marshy", "moist", "muddy", "soggy", "swampy", "tropical"},
	Nouns:      []string{"marsh", "swamp", "wetland", "bog"},
	Part1: []string{
		"The air is heavy with the smell of mud and decay, as waterlogged plants and animals rot in the hot sun.",
		"The ground is soft and spongy, giving way beneath the feet with each step.",
	},
	Part2: []string{
		"The sun beats down mercilessly, turning the swamp into a steamy, mosquito-infested hellscape.",
		"The humidity is suffocating, making it hard to catch one's breath.",
	},
	Part3: []string{
		"The swamp is home to a diverse array of plants and animals, including tall reeds, swaying grasses, and brightly-colored flowers.",
		"The swamp is teeming with life, from the smallest insects to the largest predators.",
	},
}

var temperateRainforestDescription = BiomeDescription{
	Adjectives: []string{"damp", "humid", "lush", "moist", "rainy", "verdant"},
	Nouns:      []string{"forest", "rainforest", "jungle"},
	Part1: []string{
		"The air is thick with the sound of dripping water, as rain patters against the leaves and falls to the ground.",
		"The ground is soft and spongy, muffling the footsteps of those who walk among the trees.",
	},
	Part2: []string{
		"The canopy overhead is dense and impenetrable, blocking out most of the sunlight and creating a dim and shadowy realm.",
		"The sunlight filters through the leaves in shafts of gold, casting a warm and peaceful glow over the forest.",
	},
	Part3: []string{
		"The rainforest is home to a diverse array of plants and animals, including tall trees, colourful flowers, and exotic birds.",
		"The rainforest teems with life, from the tiniest insects to the largest mammals.",
	},
}

var tropicalRainforestDescription = BiomeDescription{
	Adjectives: []string{"damp", "humid", "lush", "moist", "rainy", "verdant", "tropical"},
	Nouns:      []string{"forest", "rainforest", "jungle"},
	Part1: []string{
		"The air is thick with the sound of dripping water, as rain patters against the leaves and falls to the ground.",
		"The ground is soft and spongy, muffling the footsteps of those who walk among the trees.",
	},
	Part2: []string{
		"The canopy overhead is dense and impenetrable, blocking out most of the sunlight and creating a dim and shadowy realm.",
		"The sunlight filters through the leaves in shafts of gold, casting a warm and peaceful glow over the forest.",
	},
	Part3: []string{
		"The rainforest is home to a diverse array of plants and animals, including tall trees, colourful flowers, and exotic birds.",
		"The rainforest teems with life, from the tiniest insects to the largest mammals.",
	},
}

var shrublandDescription = BiomeDescription{
	Adjectives: []string{"brushy", "scrubby", "thick", "verdant"},
	Nouns:      []string{"scrubland", "shrubland", "thicket"},
	Part1: []string{
		"The ground is covered in a thick layer of brush and scrub, making it difficult to traverse.",
		"The shrubs and bushes are dense and impenetrable, blocking out most of the sunlight and creating a dim and shadowy realm.",
	},
	Part2: []string{
		"The air is dry and dusty, with a constant breeze rustling the leaves of the shrubs.",
		"The air is heavy and still, with the sound of buzzing insects filling the air.",
	},
	Part3: []string{
		"The shrubland is home to a diverse array of plants and animals, including hardy grasses, colorful flowers, and small mammals.",
		"The shrubland is dotted with clusters of hardy bushes and low-lying trees, providing shelter and food for a variety of species.",
	},
}

var temperateGrasslandDescription = BiomeDescription{
	Adjectives: []string{"flat", "grassy", "rolling", "verdant"},
	Nouns:      []string{"grassland", "meadow", "prairie"},
	Part1: []string{
		"The grass sways gently in the breeze, creating a sea of green that stretches out to the horizon.",
		"The ground is flat and open, with only a few scattered clusters of bushes and trees dotting the landscape.",
	},
	Part2: []string{
		"The air is fresh and clean, with a light breeze carrying the scent of wildflowers.",
		"The sun beats down mercilessly, baking the earth and turning the grassland into a dry, parched wasteland.",
	},
	Part3: []string{
		"The grassland is home to a diverse array of plants and animals, including tall grasses, colourful wildflowers, and small mammals.",
		"The grassland is teeming with life, from the smallest insects to the largest predators.",
	},
}

var tundraDescription = BiomeDescription{
	Adjectives: []string{"barren", "frozen", "icy", "snowy", "treeless"},
	Nouns:      []string{"tundra", "wasteland", "waste"},
	Part1: []string{
		"The ground is frozen solid, with a layer of snow covering the surface.",
		"The snow is deep and drifts are piled high, making it difficult to traverse the landscape.",
	},
	Part2: []string{
		"The air is crisp and cold, biting at the exposed skin.",
		"The air is still and silent, with a deathly chill that seems to seep into the bones.",
	},
	Part3: []string{
		"The tundra is home to a hardy array of plants and animals, including lichens, mosses, and small mammals.",
		"The tundra is a barren and lifeless place, with only a few scattered patches of vegetation surviving in the harsh conditions.",
	},
}

var savannahDescription = BiomeDescription{
	Adjectives: []string{"arid", "dry", "grassy", "savannah"},
	Nouns:      []string{"plain", "landscape", "expanse"},
	Part1: []string{
		"The ground is dry and dusty, with patches of grass and scattered clusters of bushes and trees dotting the landscape.",
		"The grass is tall and waving, with a few scattered trees providing shade and shelter.",
	},
	Part2: []string{
		"The sun beats down mercilessly, baking the earth and turning the savannah into a dry, parched wasteland.",
		"The air is dry and hot, with a light breeze providing some relief from the scorching sun.",
	},
	Part3: []string{
		"The savannah is home to a diverse array of plants and animals, including tall grasses, colourful wildflowers, and large mammals.",
		"The savannah teems with life, from the smallest insects to the largest predators.",
	},
}

var borealForestDescription = BiomeDescription{
	Adjectives: []string{"dense", "frozen", "snowy"},
	Nouns:      []string{"forest", "trees", "taiga"},
	Part1: []string{
		"The canopy overhead is dense and impenetrable, blocking out most of the sunlight and creating a dim and shadowy realm.",
		"The sunlight filters through the leaves in shafts of gold, casting a warm and peaceful glow over the forest.",
	},
	Part2: []string{
		"The ground is frozen solid, with a layer of snow covering the surface.",
		"The snow is deep and drifts are piled high, making it difficult to traverse the landscape.",
	},
	Part3: []string{
		"The air is crisp and cold, biting at the exposed skin.",
		"The air is still and silent, with a deathly chill that seems to seep into the bones.",
	},
	Part4: []string{
		"The boreal forest is home to a hardy array of plants and animals, including conifers, small mammals, and migratory birds.",
		"The boreal forest is a harsh and unforgiving place, with only the most resilient species able to survive the long, cold winters.",
	},
}

var temperateSeasonalForestDescription = BiomeDescription{
	Adjectives: []string{"dense", "lush", "temperate"},
	Nouns:      []string{"forest", "woods"},
	Part1: []string{
		"The canopy overhead is dense and verdant, with a riot of leaves and branches creating a green, leafy ceiling.",
		"The sunlight filters through the leaves in dappled patterns, casting a warm and peaceful glow over the forest.",
	},
	Part2: []string{
		"The ground is soft and spongy, with a thick layer of leaves and debris covering the surface.",
		"The ground is hard and rocky, with only a few patches of moss and lichen providing a splash of color.",
	},
	Part3: []string{
		"The air is fresh and clean, with a light breeze carrying the scent of pine and wood smoke.",
		"The air is humid and thick, with the sound of buzzing insects and chirping birds filling the air.",
	},
	Part4: []string{
		"The forest is home to a diverse array of plants and animals, including deciduous trees, colourful wildflowers, and small mammals.",
		"The woods are teeming with life, from the smallest insects to the largest predators.",
	},
}

var tropicalSeasonalForestDescription = BiomeDescription{
	Adjectives: []string{"dense", "lush", "tropical"},
	Nouns:      []string{"forest", "jungle", "rainforest"},
	Part1: []string{
		"The canopy overhead is dense and verdant, with a riot of leaves and branches creating a green, leafy ceiling.",
		"The sunlight filters through the leaves in dappled patterns, casting a warm and peaceful glow over the forest.",
	},
	Part2: []string{
		"The ground is soft and spongy, with a thick layer of leaves and debris covering the surface.",
		"The ground is hard and rocky, with only a few patches of moss and lichen providing a splash of color.",
	},
	Part3: []string{
		"The air is hot and humid, with the sound of buzzing insects and chirping birds filling the air.",
		"The air is heavy with the scent of flowers and the sound of distant thunder, as a tropical rainstorm approaches.",
	},
	Part4: []string{
		"The forest is home to a diverse array of plants and animals, including tall trees, colourful wildflowers, and exotic creatures.",
		"The jungle is teeming with life, from the smallest insects to the largest predators.",
	},
}

var snowDescription = BiomeDescription{
	Adjectives: []string{"frozen", "glacial", "icy", "snowy"},
	Nouns:      []string{"plain", "wasteland", "landscape"},
	Part1: []string{
		"The ground is frozen solid, with a layer of snow covering the surface.",
		"The snow is deep and drifts are piled high, making it difficult to traverse the terrain.",
	},
	Part2: []string{
		"The air is crisp and cold, biting at the exposed skin.",
		"The air is still and silent, with a deathly chill that seems to seep into the bones.",
	},
	Part3: []string{
		"The snowy environment is barren and lifeless, with only a few hardy plants and animals able to survive the harsh conditions.",
		"The snowy landscape is home to a hardy array of plants and animals, including lichens, mosses, and small mammals.",
	},
}

type NameGenerators struct {
	Desert        *BasicGenerator
	Mountain      *BasicGenerator
	MountainRange *BasicGenerator
	Forest        *BasicGenerator
	Swamp         *BasicGenerator
	River         *RiverGenerator
}

func NewNameGenerators(seed int64) *NameGenerators {
	return &NameGenerators{
		Desert:        NewDesertGenerator(seed),
		Mountain:      NewMountainGenerator(seed),
		MountainRange: NewMountainRangeGenerator(seed),
		Forest:        NewForestGenerator(seed),
		Swamp:         NewSwampGenerator(seed),
		River:         NewRiverGenerator(seed),
	}
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
		"decaying",
		"rotten",
		"rotting",
		"putrid",
		"stinking",
		"stenchy",
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
	name := "The " + g.randomPair(g.Prefix, g.Suffix)
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
	// Make sure that the string chosen from a is not contained in b and vice versa.
	// This is to avoid names like "The muddy mud" or "The swampy swamp".
	for i := 0; i < 100; i++ {
		s1 := n.randomString(a)
		s2 := n.randomString(b)
		if !strings.Contains(s2, s1) && !strings.Contains(s1, s2) {
			return s1 + " " + s2
		}
	}
	// If we can't find a pair that doesn't contain the other, just return a random pair.
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
