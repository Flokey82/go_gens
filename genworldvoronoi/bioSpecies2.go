package genworldvoronoi

import (
	"log"
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/genbiome"
)

func (b *Bio) placeAllSpecies(kingdom *BioLevel) {
	// Place all species.
	for _, species := range kingdom.getAllByLevel(BioLevelSpecies) {
		if species == nil {
			log.Println("no species found for " + kingdom.Name)
			continue
		}
		b.placeSpeciesFromLevel(species)
	}
	log.Println("Placed", len(b.Species), "species")
}

func (b *Bio) placeSpeciesFromLevel(level *BioLevel) {
	// Get the fitness function for the biolevel tolerances.
	sf := b.getToleranceScoreFunc(level.SpeciesTolerances)

	scoreNbs := func(r int) float64 {
		score := sf(r)
		var count int
		for _, n := range b.GetRegNeighbors(r) {
			score += sf(n)
			count++
		}
		return score / float64(count+1)
	}

	// TODO: Use directly competing species as seeds to maximize
	// distance between species that compete for the same resources.
	distSeedFunc := func() []int {
		var res []int
		for _, s := range b.Species {
			res = append(res, s.Origin)
		}
		return res
	}

	// Score all regions, pick highest score.
	newspecies := -1
	lastMax := math.Inf(-1)
	for i, val := range b.CalcFitnessScore(scoreNbs, distSeedFunc) {
		if val > lastMax && val != -1 {
			newspecies = i
			lastMax = val
		}
	}
	if newspecies == -1 {
		log.Println("no newspecies found for " + level.Name)
		return
	}
	log.Println("Placing species", level.Name, "at", newspecies, "with score", lastMax)
	s := level.ToSpecies()
	s.Origin = newspecies
	b.Species = append(b.Species, s)
}

const (
	BioLevelKingdom = iota
	BioLevelFamily
	BioLevelGenus
	BioLevelSpecies
	BioLevelSubspecies
)

type BioProperty int

const (
	BioPropertyEdibleSeed BioProperty = iota
	BioPropertyEdibleFruit
	BioPropertyEdibleFlower
	BioPropertyEdibleLeafs
)

type BioLevel struct {
	Name        string
	Level       int
	Parent      *BioLevel
	Children    []*BioLevel
	Hereditary  []BioProperty // Guaranteed to be inherited.
	Inheritable []BioProperty // May be inherited.
	SpeciesProperties
	SpeciesTolerances
}

// TODO: Add properties like: has fruits, seeds, flowers, etc.
// What parts are edible, what parts are poisonous, etc.
//
// TODO: Add generic numeric properties like ratio of area to edible weight, etc.
// Potentially we could use the genetics package for all of this.
//
// TODO: Add restrictions that act as binary masks.
// The environment can be encoded as a number of bits and the species can have
// a mask that specifies which bits are required to be set for the species to
// survive. This could be used to encode things like temperature, humidity,
// soil type, etc
// Each restriction type could be a separate byte. That way we only need to
// apply the mask to each byte and then check if the result is zero.
// We have currently 15 biomes in the extended whittaker model, so we could
// encode the biomes that a plant can survive in as a 15 bit mask or a 16 bit
// number.
//
// TODO: Allow generation of random species and sub-species.

func NewKingdom(name string, options ...BioLevelOption) *BioLevel {
	kingdom := &BioLevel{
		Name:  name,
		Level: BioLevelKingdom,
	}

	// Apply options.
	for _, option := range options {
		option(kingdom)
	}
	return kingdom
}

func (b *BioLevel) NewChild(name string, options ...BioLevelOption) *BioLevel {
	child := &BioLevel{
		Name:              name,
		Level:             b.Level + 1,
		Hereditary:        b.Hereditary,
		Parent:            b,
		SpeciesProperties: b.SpeciesProperties,
		SpeciesTolerances: b.SpeciesTolerances,
	}

	// Pick a random number of inheritable properties.
	n := rand.Intn(len(b.Inheritable) + 1)
	for i, prop := range rand.Perm(len(b.Inheritable)) {
		if i >= n {
			break
		}
		child.Inheritable = append(child.Inheritable, b.Inheritable[prop])
	}

	// Apply options.
	for _, option := range options {
		option(child)
	}

	b.Children = append(b.Children, child)
	return child
}

func (b *BioLevel) getAllByLevel(level int) []*BioLevel {
	if b.Level > level {
		return nil
	}
	if b.Level == level {
		return []*BioLevel{b}
	}
	var result []*BioLevel
	for _, child := range b.Children {
		result = append(result, child.getAllByLevel(level)...)
	}
	return result
}

func (b *BioLevel) ToSpecies() *Species {
	return &Species{
		Name:              b.Name,
		SpeciesProperties: b.SpeciesProperties,
		SpeciesTolerances: b.SpeciesTolerances,
	}
}

var (
	KingdomFauna = NewKingdom("Fauna")
	KingdomFlora = NewKingdom("Flora",
		BioLevelDigestion(DigestivePhotosynthetic))
	KingdomFungi = NewKingdom("Fungi",
		BioLevelDigestion(DigestiveDecomposer),
		BioLevelHumidityRange(0.2, 1))

	// Trees.
	FamilyTree = KingdomFauna.NewChild("Tree",
		BioLevelEcosphere(EcosphereTypeLand),
		BioLevelSteepRange(0, 0.5),
		BioLevelTempRange(0, 35),
		BioLevelHumidityRange(0.4, 0.9),
		BioLevelElevRange(0, 0.7))

	// Leaf trees.
	GenusLeafyTree = FamilyTree.NewChild("Leafy Tree") // Trees that produce leaves.
	SpeciesOak     = GenusLeafyTree.NewChild("Oak",
		BioLevelHumidityRange(0.4, 0.8))

	// Needle trees.
	GenusNeedleTree = FamilyTree.NewChild("Needle Tree") // Trees that produce needles.
	SpeciesPine     = GenusNeedleTree.NewChild("Pine",
		BioLevelHumidityRange(0.4, 0.8))

	// Shrubs.
	FamilyShrub = KingdomFlora.NewChild("Shrub",
		BioLevelEcosphere(EcosphereTypeLand))
	GenusLeafyShrub  = FamilyShrub.NewChild("Leafy Shrub")  // Shrubs that produce leaves.
	GenusNeedleShrub = FamilyShrub.NewChild("Needle Shrub") // Shrubs that produce needles.

	// Grasses.
	FamilyGrass = KingdomFlora.NewChild("Grass",
		BioLevelTempRange(0, 35),
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver)) // Grasses can grow alongside rivers.
	GenusCereal = FamilyGrass.NewChild("Cereal",
		BioLevelTempRange(14, 28),
		BioLevelHumidityRange(0.1, 0.8),
		BioLevelAppendHereditary(BioPropertyEdibleSeed), // Cereals are grasses that produce seeds.
	)
	SpeciesWheat = GenusCereal.NewChild("Wheat",
		BioLevelTempRange(20, 25),
		BioLevelSteepRange(0, 0.7))
	SpeciesRice = GenusCereal.NewChild("Rice",
		BioLevelTempRange(21, 37),
		BioLevelHumidityRange(0.6, 1),
		BioLevelSteepRange(0, 0.4))
	SpeciesBarley = GenusCereal.NewChild("Barley",
		BioLevelTempRange(12, 25))
	SpeciesOats = GenusCereal.NewChild("Oats",
		BioLevelTempRange(15, 23))
	SpeciesCorn = GenusCereal.NewChild("Corn",
		BioLevelTempRange(20, 30))
	SpeciesRye = GenusCereal.NewChild("Rye",
		BioLevelTempRange(25, 35))
	SpeciesBuckwheat = GenusCereal.NewChild("Buckwheat",
		BioLevelTempRange(17, 27))
	SpeciesMillet = GenusCereal.NewChild("Millet",
		BioLevelTempRange(26, 29))
	SpeciesSorghum = GenusCereal.NewChild("Sorghum",
		BioLevelTempRange(26, 34))

	GenusReed        = FamilyGrass.NewChild("Reed")
	SpeciesSugarCane = GenusReed.NewChild("Sugar Cane",
		BioLevelHumidityRange(0.6, 1),
		BioLevelTempRange(20, 35))
	SpeciesBamboo = GenusReed.NewChild("Bamboo",
		BioLevelTempRange(15, 26))

	// Generic grasses.
	GenusGrass       = FamilyGrass.NewChild("Grass")
	SpeciesGrass     = GenusGrass.NewChild("Grass")
	SpeciesCrabGrass = GenusGrass.NewChild("Crab Grass")

	// Herbs.
	FamilyHerb = KingdomFlora.NewChild("Herb",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver))

	// Flowers.
	FamilyFlower = KingdomFlora.NewChild("Flower",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver))

	// Leafy greens.
	GenusLeafyGreen = FamilyFlower.NewChild("Leafy Green",
		BioLevelAppendHereditary(BioPropertyEdibleLeafs))
	SpeciesCabbage = GenusLeafyGreen.NewChild("Cabbage")
	SpeciesLettuce = GenusLeafyGreen.NewChild("Lettuce")
	SpeciesSpinach = GenusLeafyGreen.NewChild("Spinach")

	// Generic flowers.
	GenusFlower  = FamilyFlower.NewChild("Flower")
	SpeciesRose  = GenusFlower.NewChild("Rose")
	SpeciesTulip = GenusFlower.NewChild("Tulip")
	SpeciesDaisy = GenusFlower.NewChild("Daisy")

	// Ferns.
	FamilyFern = KingdomFlora.NewChild("Fern",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver))

	// Mosses.
	FamilyMoss = KingdomFlora.NewChild("Moss",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver))

	// Vines.
	FamilyVine = KingdomFlora.NewChild("Vine",
		BioLevelEcosphere(EcosphereTypeLand))
	GenusGrapeVine = FamilyVine.NewChild("Grape Vine",
		BioLevelAppendHereditary(BioPropertyEdibleFruit))
	SpeciesWineGrape = GenusGrapeVine.NewChild("Wine Grape",
		BioLevelTempRange(12, 20),
		BioLevelHumidityRange(0.4, 0.8))
	SpeciesKiwiVine = GenusGrapeVine.NewChild("Kiwi Vine",
		BioLevelTempRange(18, 28),
		BioLevelHumidityRange(0.4, 0.8))

	GenusClimbingVine   = FamilyVine.NewChild("Climbing Vine")
	SpeciesMorningGlory = GenusClimbingVine.NewChild("Morning Glory")
	SpeciesClimbingBean = GenusClimbingVine.NewChild("Climbing Bean")

	// Cacti.
	FamilyCactus = KingdomFlora.NewChild("Cactus",
		BioLevelHumidityRange(0, 0.2),
		BioLevelTempRange(10, 35),
		BioLevelEcosphere(EcosphereTypeLand))
	FamilySucculent = KingdomFlora.NewChild("Succulent",
		BioLevelHumidityRange(0, 0.2),
		BioLevelTempRange(10, 35),
		BioLevelEcosphere(EcosphereTypeLand))

	// Insects.
	FamilyInsect = KingdomFauna.NewChild("Insect",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelDigestion(DigestiveSystemHerbivore),
		BioLevelLocomotion(LocomotionWalk|LocomotionClimb))
	GenusAnt         = FamilyInsect.NewChild("Ant")
	GenusBeetle      = FamilyInsect.NewChild("Beetle")
	GenusCaterpillar = FamilyInsect.NewChild("Caterpillar")
	GenusGrasshopper = FamilyInsect.NewChild("Grasshopper")
	GenusButterfly   = FamilyInsect.NewChild("Butterfly",
		BioLevelLocomotion(LocomotionFly))
	GenusDragonfly = FamilyInsect.NewChild("Dragonfly",
		BioLevelLocomotion(LocomotionFly))
	GenusFly = FamilyInsect.NewChild("Fly",
		BioLevelLocomotion(LocomotionFly))
	GenusMoth = FamilyInsect.NewChild("Moth",
		BioLevelLocomotion(LocomotionFly))

	// Arachnids.
	FamilyArachnid = KingdomFauna.NewChild("Arachnid",
		BioLevelEcosphere(EcosphereTypeLand),
		BioLevelLocomotion(LocomotionWalk|LocomotionClimb))
	GenusSpider   = FamilyArachnid.NewChild("Spider")
	GenusScorpion = FamilyArachnid.NewChild("Scorpion")
	GenusTick     = FamilyArachnid.NewChild("Tick")

	// Mammals.
	FamilyMammal = KingdomFauna.NewChild("Mammal")
	GenusWhale   = FamilyMammal.NewChild("Whale",
		BioLevelEcosphere(EcosphereTypeOcean),
		BioLevelDigestion(DigestiveSystemCarnivore),
		BioLevelSize(SpeciesSizeLarge),
		BioLevelLocomotion(LocomotionSwim))
	SpeciesBlueWhale     = GenusWhale.NewChild("Blue Whale")
	SpeciesHumpbackWhale = GenusWhale.NewChild("Humpback Whale")

	GenusCanine = FamilyMammal.NewChild("Canine",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelDigestion(DigestiveSystemCarnivore),
		BioLevelLocomotion(LocomotionWalk))
	SpeciesDog = GenusCanine.NewChild("Dog",
		BioLevelSize(SpeciesSizeSmall),
		BioLevelTempRange(21, 26))
	SpeciesWolf = GenusCanine.NewChild("Wolf",
		BioLevelSize(SpeciesSizeMedium),
		BioLevelTempRange(-50, 48.8))
	SpeciesFox = GenusCanine.NewChild("Fox",
		BioLevelSize(SpeciesSizeSmall),
		BioLevelTempRange(-15, 15))

	// Felines.
	GenusFeline = FamilyMammal.NewChild("Feline",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelDigestion(DigestiveSystemCarnivore),
		BioLevelLocomotion(LocomotionWalk))
	SpeciesCat = GenusFeline.NewChild("Cat",
		BioLevelSize(SpeciesSizeSmall),
		BioLevelTempRange(25, 30),
	)
	SpeciesLion = GenusFeline.NewChild("Lion",
		BioLevelSize(SpeciesSizeMedium),
		BioLevelTempRange(20, 30),
	)
	SpeciesTiger = GenusFeline.NewChild("Tiger",
		BioLevelSize(SpeciesSizeMedium))
	SpeciesLeopard = GenusFeline.NewChild("Leopard",
		BioLevelSize(SpeciesSizeMedium),
		BioLevelTempRange(25, 32),
	)
	SpeciesSandCat = GenusFeline.NewChild("Sand Cat",
		BioLevelSize(SpeciesSizeTiny),
		BioLevelPreferredBiomes(genbiome.WhittakerModBiomeSubtropicalDesert),
		BioLevelTempRange(25, 35),
	)
	// ... etc.

	// Rodents.
	GenusRodent = FamilyMammal.NewChild("Rodent",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelDigestion(DigestiveSystemHerbivore),
		BioLevelSize(SpeciesSizeTiny),
		BioLevelLocomotion(LocomotionBurrow|LocomotionClimb|LocomotionSwim|LocomotionWalk))
	SpeciesMouse     = GenusRodent.NewChild("Mouse")
	SpeciesRat       = GenusRodent.NewChild("Rat")
	SpeciesHamster   = GenusRodent.NewChild("Hamster")
	SpeciesGuineaPig = GenusRodent.NewChild("Guinea Pig")
	SpeciesSquirrel  = GenusRodent.NewChild("Squirrel")
	SpeciesChipmunk  = GenusRodent.NewChild("Chipmunk")
	SpeciesRaccoon   = GenusRodent.NewChild("Raccoon")
	SpeciesBeaver    = GenusRodent.NewChild("Beaver",
		BioLevelEcosphere(EcosphereTypeRiver))
	SpeciesPorcupine = GenusRodent.NewChild("Porcupine")
	SpeciesGopher    = GenusRodent.NewChild("Gopher")
	SpeciesGroundhog = GenusRodent.NewChild("Groundhog")

	// Birds.
	FamilyBird = KingdomFauna.NewChild("Bird",
		BioLevelEcosphere(EcosphereTypeLand),
		BioLevelLocomotion(LocomotionFly))

	// Fishes.
	FamilyFish = KingdomFauna.NewChild("Fish",
		BioLevelLocomotion(LocomotionSwim))
	GenusSaltwaterFish = FamilyFish.NewChild("Saltwater Fish",
		BioLevelEcosphere(EcosphereTypeOcean),
		BioLevelDigestion(DigestiveSystemCarnivore))
	SpeciesShark   = GenusSaltwaterFish.NewChild("Shark")
	SpeciesTuna    = GenusSaltwaterFish.NewChild("Tuna")
	SpeciesCod     = GenusSaltwaterFish.NewChild("Cod")
	SpeciesHaddock = GenusSaltwaterFish.NewChild("Haddock")
	SpeciesHerring = GenusSaltwaterFish.NewChild("Herring")

	GenusFreshwaterFish = FamilyFish.NewChild("Freshwater Fish",
		BioLevelEcosphere(EcosphereTypeRiver|EcosphereTypeLake))
	SpeciesTrout  = GenusFreshwaterFish.NewChild("Trout")
	SpeciesSalmon = GenusFreshwaterFish.NewChild("Salmon")
	SpeciesBass   = GenusFreshwaterFish.NewChild("Bass")

	// Crustaceans.
	FamilyCrustacean = KingdomFauna.NewChild("Crustacean",
		BioLevelLocomotion(LocomotionWalk|LocomotionSwim))

	// Crabs.
	GenusCrab        = FamilyCrustacean.NewChild("Crab")
	SpeciesRiverCrab = GenusCrab.NewChild("River Crab",
		BioLevelEcosphere(EcosphereTypeRiver),
		BioLevelElevRange(-0.01, 0.01),
	)
	SpeciesRockCrab = GenusCrab.NewChild("Rock Crab",
		BioLevelEcosphere(EcosphereTypeOcean|EcosphereTypeLand),
		BioLevelElevRange(-0.3, -0.1),
	)

	// Shrimp.
	GenusShrimp = FamilyCrustacean.NewChild("Shrimp",
		BioLevelEcosphere(EcosphereTypeOcean),
		BioLevelTempRange(14, 29),
	)

	// Lobsters.
	GenusLobster    = FamilyCrustacean.NewChild("Lobster")
	SpeciesCrayfish = GenusLobster.NewChild("Crayfish",
		BioLevelEcosphere(EcosphereTypeRiver|EcosphereTypeLake),
		BioLevelTempRange(18, 25),
	)
	SpeciesLobster = GenusLobster.NewChild("Lobster",
		BioLevelEcosphere(EcosphereTypeOcean),
		BioLevelTempRange(16, 18),
	)

	// Mollusks.
	FamilyMollusk = KingdomFauna.NewChild("Mollusk")

	// Clams.
	GenusClam     = FamilyMollusk.NewChild("Clam")
	SpeciesMussel = GenusClam.NewChild("Mussel",
		BioLevelEcosphere(EcosphereTypeOcean|EcosphereTypeRiver|EcosphereTypeLake))
	SpeciesOyster = GenusClam.NewChild("Oyster",
		BioLevelEcosphere(EcosphereTypeOcean))

	// Snails.
	GenusSnail = FamilyMollusk.NewChild("Snail",
		BioLevelDigestion(DigestiveSystemHerbivore),
		BioLevelLocomotion(LocomotionSlither))
	SpeciesSlug = GenusSnail.NewChild("Slug",
		BioLevelEcosphere(EcosphereTypeLand),
		BioLevelHumidityRange(0.2, 1),
		BioLevelTempRange(15, 20),
	)
	SpeciesSeaSnail = GenusSnail.NewChild("Sea Snail",
		BioLevelEcosphere(EcosphereTypeOcean),
		BioLevelElevRange(-0.01, 0.0),
	)

	// Octopuses.
	GenusOctopus = FamilyMollusk.NewChild("Octopus",
		BioLevelEcosphere(EcosphereTypeOcean),
		BioLevelDigestion(DigestiveSystemCarnivore),
		BioLevelLocomotion(LocomotionSwim))
	SpeciesSquid = GenusOctopus.NewChild("Squid")

	// Amphibians.
	FamilyAmphibian = KingdomFauna.NewChild("Amphibian",
		BioLevelLocomotion(LocomotionSwim|LocomotionWalk))

	// Frogs.
	GenusFrog = FamilyAmphibian.NewChild("Frog",
		BioLevelHumidityRange(0.2, 1))
	SpeciesBullfrog = GenusFrog.NewChild("Bullfrog",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelElevRange(0, 0.1),
	)
	SpeciesTreeFrog = GenusFrog.NewChild("Tree Frog",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelElevRange(0.1, 0.4),
		BioLevelPreferredBiomes(
			genbiome.WhittakerModBiomeTemperateSeasonalForest,
			genbiome.WhittakerModBiomeTemperateRainforest,
			genbiome.WhittakerModBiomeTropicalSeasonalForest,
			genbiome.WhittakerModBiomeTropicalRainforest,
		),
	)

	// Toads.
	GenusToad   = FamilyAmphibian.NewChild("Toad")
	SpeciesBufo = GenusToad.NewChild("Bufo",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelElevRange(0, 0.1),
	)
	GenusSalamander = FamilyAmphibian.NewChild("Salamander")

	// Reptiles.
	FamilyReptile = KingdomFauna.NewChild("Reptile",
		BioLevelDigestion(DigestiveSystemCarnivore))

	// Snakes.
	GenusSerpent = FamilyReptile.NewChild("Serpent",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelTempRange(15, 30),
		BioLevelElevRange(0, 0.4),
		BioLevelLocomotion(LocomotionSlither),
		BioLevelSize(SpeciesSizeSmall))
	SpeciesSnake = GenusSerpent.NewChild("Snake",
		BioLevelHumidityRange(0.2, 1),
	)
	SpeciesCobra = GenusSerpent.NewChild("Cobra",
		BioLevelHumidityRange(0.01, 1),
		BioLevelPreferredBiomes(
			genbiome.WhittakerModBiomeTemperateRainforest,
			genbiome.WhittakerModBiomeTropicalRainforest,
		),
	)
	SpeciesRattlesnake = GenusSerpent.NewChild("Rattlesnake",
		BioLevelHumidityRange(0.01, 0.2),
		BioLevelPreferredBiomes(
			genbiome.WhittakerModBiomeSubtropicalDesert,
		),
	)
	SpeciesAnaconda = GenusSerpent.NewChild("Anaconda",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelSize(SpeciesSizeLarge),
		BioLevelLocomotion(LocomotionSlither|LocomotionClimb),
		BioLevelPreferredBiomes(
			genbiome.WhittakerModBiomeTropicalSeasonalForest,
			genbiome.WhittakerModBiomeTropicalRainforest,
		),
	)

	// Crocodiles.
	GenusCrocodilia = FamilyReptile.NewChild("Crocodilia",
		BioLevelTempRange(20, 30),
		BioLevelLocomotion(LocomotionWalk|LocomotionSwim),
		BioLevelSize(SpeciesSizeLarge))
	SpeciesCrocodile = GenusCrocodilia.NewChild("Crocodile",
		BioLevelEcosphere(EcosphereTypeRiver|EcosphereTypeLake|EcosphereTypeOcean),
		BioLevelElevRange(0, 0.01),
	)
	SpeciesAlligator = FamilyReptile.NewChild("Alligator",
		BioLevelEcosphere(EcosphereTypeRiver),
		BioLevelElevRange(0, 0.02),
	)

	// Lizards.
	GenusLizard = FamilyReptile.NewChild("Lizard",
		BioLevelEcosphere(EcosphereTypeLand|EcosphereTypeRiver),
		BioLevelElevRange(0, 0.2),
		BioLevelHumidityRange(0.2, 1),
		BioLevelTempRange(15, 30),
		BioLevelLocomotion(LocomotionWalk),
	)
	SpeciesIguana = GenusLizard.NewChild("Iguana")
	SpeciesGecko  = GenusLizard.NewChild("Gecko")

	// Turtles.
	GenusTurtle = FamilyReptile.NewChild("Turtle",
		BioLevelLocomotion(LocomotionWalk|LocomotionSwim))
	SpeciesTortoise = GenusTurtle.NewChild("Tortoise",
		BioLevelEcosphere(EcosphereTypeLand),
		BioLevelElevRange(0, 0.2),
		BioLevelHumidityRange(0.2, 1),
		BioLevelTempRange(10, 20),
	)
	SpeciesTurle = GenusTurtle.NewChild("Turle",
		BioLevelEcosphere(EcosphereTypeOcean),
		BioLevelElevRange(-0.01, 0.0),
	)

	// Worms.
	FamilyWorm = KingdomFauna.NewChild("Worm",
		BioLevelDigestion(DigestiveSystemOmnivore),
		BioLevelLocomotion(LocomotionSlither|LocomotionBurrow),
	)
	GenusEarthworm = FamilyWorm.NewChild("Earthworm",
		BioLevelEcosphere(EcosphereTypeLand),
		BioLevelElevRange(0, 0.2),
		BioLevelHumidityRange(0.2, 1),
		BioLevelTempRange(10, 20),
	)

	// Fungi.
	FamilyMushroom    = KingdomFungi.NewChild("Mushroom")
	GenusMushroom     = FamilyMushroom.NewChild("Mushroom")
	SpeciesChampignon = GenusMushroom.NewChild("Champignon",
		BioLevelEcosphere(EcosphereTypeLand),
		BioLevelElevRange(0, 0.2),
		BioLevelTempRange(10, 20),
		BioLevelPreferredBiomes(genbiome.WhittakerModBiomeTemperateSeasonalForest, genbiome.WhittakerModBiomeTemperateRainforest),
	)
	FamilyMold = KingdomFungi.NewChild("Mold")
)
