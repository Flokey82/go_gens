package genworldvoronoi

import "math/rand"

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
}

// TODO: Add properties like: has fruits, seeds, flowers, etc.
// What parts are edible, what parts are poisonous, etc.
// TODO: Add generic numeric properties like ratio of area to edible weight, etc.
// Potentially we could use the genetics package for all of this.
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

func NewKingdom(name string) *BioLevel {
	return &BioLevel{
		Name:  name,
		Level: BioLevelKingdom,
	}
}

func (b *BioLevel) NewChild(name string, hereditary ...BioProperty) *BioLevel {
	child := &BioLevel{
		Name:       name,
		Level:      b.Level + 1,
		Hereditary: append(hereditary, b.Hereditary...),
		Parent:     b,
	}

	// Pick a random number of inheritable properties.
	n := rand.Intn(len(b.Inheritable) + 1)
	for i, prop := range rand.Perm(len(b.Inheritable)) {
		if i >= n {
			break
		}
		child.Inheritable = append(child.Inheritable, b.Inheritable[prop])
	}
	b.Children = append(b.Children, child)
	return child
}

func (b *BioLevel) appendHereditaries(props ...BioProperty) *BioLevel {
	b.Hereditary = append(b.Hereditary, props...)
	return b
}

func (b *BioLevel) appendInheritables(props ...BioProperty) *BioLevel {
	b.Inheritable = append(b.Inheritable, props...)
	return b
}

var (
	KingdomFauna = NewKingdom("Fauna")
	KingdomFlora = NewKingdom("Flora")
	KingdomFungi = NewKingdom("Fungi")

	// Trees.
	FamilyTree      = KingdomFauna.NewChild("Tree")
	GenusLeafyTree  = FamilyTree.NewChild("Leafy Tree")  // Trees that produce leaves.
	GenusNeedleTree = FamilyTree.NewChild("Needle Tree") // Trees that produce needles.

	// Shrubs.
	FamilyShrub      = KingdomFlora.NewChild("Shrub")
	GenusLeafyShrub  = FamilyShrub.NewChild("Leafy Shrub")  // Shrubs that produce leaves.
	GenusNeedleShrub = FamilyShrub.NewChild("Needle Shrub") // Shrubs that produce needles.

	// Grasses.
	FamilyGrass = KingdomFlora.NewChild("Grass")
	GenusCereal = FamilyGrass.NewChild("Cereal").
			appendHereditaries(BioPropertyEdibleSeed) // Cereals are grasses that produce seeds.
	SpeciesWheat     = GenusCereal.NewChild("Wheat")
	SpeciesRice      = GenusCereal.NewChild("Rice")
	SpeciesBarley    = GenusCereal.NewChild("Barley")
	SpeciesOats      = GenusCereal.NewChild("Oats")
	SpeciesCorn      = GenusCereal.NewChild("Corn")
	SpeciesRye       = GenusCereal.NewChild("Rye")
	SpeciesBuckwheat = GenusCereal.NewChild("Buckwheat")
	SpeciesMillet    = GenusCereal.NewChild("Millet")
	SpeciesSorghum   = GenusCereal.NewChild("Sorghum")
	GenusReed        = FamilyGrass.NewChild("Reed")  // Reeds are grasses that produce stems.
	GenusGrass       = FamilyGrass.NewChild("Grass") // Grasses are grasses that produce leaves.

	// Herbs.
	FamilyHerb = KingdomFlora.NewChild("Herb")

	// Flowers.
	FamilyFlower    = KingdomFlora.NewChild("Flower")
	GenusLeafyGreen = FamilyFlower.NewChild("Leafy Green").
			appendHereditaries(BioPropertyEdibleLeafs)
	SpeciesCabbage = GenusLeafyGreen.NewChild("Cabbage")
	SpeciesLettuce = GenusLeafyGreen.NewChild("Lettuce")
	SpeciesSpinach = GenusLeafyGreen.NewChild("Spinach")
	GenusFlower    = FamilyFlower.NewChild("Flower")
	SpeciesRose    = GenusFlower.NewChild("Rose")
	SpeciesTulip   = GenusFlower.NewChild("Tulip")
	SpeciesDaisy   = GenusFlower.NewChild("Daisy")

	// Ferns.
	FamilyFern = KingdomFlora.NewChild("Fern")

	// Mosses.
	FamilyMoss = KingdomFlora.NewChild("Moss")

	// Vines.
	FamilyVine        = KingdomFlora.NewChild("Vine")
	GenusGrapeVine    = FamilyVine.NewChild("Grape Vine")
	GenusClimbingVine = FamilyVine.NewChild("Climbing Vine")

	// Cacti.
	FamilyCactus    = KingdomFlora.NewChild("Cactus")
	FamilySucculent = KingdomFlora.NewChild("Succulent")

	// Insects.
	FamilyInsect     = KingdomFauna.NewChild("Insect")
	GenusAnt         = FamilyInsect.NewChild("Ant")
	GenusBeetle      = FamilyInsect.NewChild("Beetle")
	GenusButterfly   = FamilyInsect.NewChild("Butterfly")
	GenusCaterpillar = FamilyInsect.NewChild("Caterpillar")
	GenusDragonfly   = FamilyInsect.NewChild("Dragonfly")
	GenusFly         = FamilyInsect.NewChild("Fly")
	GenusGrasshopper = FamilyInsect.NewChild("Grasshopper")
	GenusMoth        = FamilyInsect.NewChild("Moth")

	// Arachnids.
	FamilyArachnid = KingdomFauna.NewChild("Arachnid")
	GenusSpider    = FamilyArachnid.NewChild("Spider")
	GenusScorpion  = FamilyArachnid.NewChild("Scorpion")
	GenusTick      = FamilyArachnid.NewChild("Tick")

	// Mammals.
	FamilyMammal = KingdomFauna.NewChild("Mammal")
	GenusWhale   = FamilyMammal.NewChild("Whale")
	GenusCanine  = FamilyMammal.NewChild("Canine")
	GenusFeline  = FamilyMammal.NewChild("Feline")
	GenusRodent  = FamilyMammal.NewChild("Rodent")

	// Birds.
	FamilyBird = KingdomFauna.NewChild("Bird")

	// Fishes.
	FamilyFish = KingdomFauna.NewChild("Fish")

	// Crustaceans.
	FamilyCrustacean = KingdomFauna.NewChild("Crustacean")
	GenusCrab        = FamilyCrustacean.NewChild("Crab")
	GenusShrimp      = FamilyCrustacean.NewChild("Shrimp")
	GenusLobster     = FamilyCrustacean.NewChild("Lobster")

	// Mollusks.
	FamilyMollusk = KingdomFauna.NewChild("Mollusk")
	GenusClam     = FamilyMollusk.NewChild("Clam")
	GenusSnail    = FamilyMollusk.NewChild("Snail")
	GenusOctopus  = FamilyMollusk.NewChild("Octopus")

	// Amphibians.
	FamilyAmphibian = KingdomFauna.NewChild("Amphibian")
	GenusFrog       = FamilyAmphibian.NewChild("Frog")
	GenusToad       = FamilyAmphibian.NewChild("Toad")
	GenusSalamander = FamilyAmphibian.NewChild("Salamander")

	// Reptiles.
	FamilyReptile = KingdomFauna.NewChild("Reptile")
	GenusSerpent  = FamilyReptile.NewChild("Serpent")
	GenusLizard   = FamilyReptile.NewChild("Lizard")
	GenusTurtle   = FamilyReptile.NewChild("Turtle")

	// Worms.
	FamilyWorm = KingdomFauna.NewChild("Worm")

	// Fungi.
	FamilyMushroom = KingdomFungi.NewChild("Mushroom")
	FamilyMold     = KingdomFungi.NewChild("Mold")
)
