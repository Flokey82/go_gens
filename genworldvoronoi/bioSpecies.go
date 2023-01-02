package genworldvoronoi

import (
	"fmt"
	"math"
	"strings"
)

func (b *Bio) genNRandomSpecies(n int) {
	// TODO: Expand species from their origin until they encounter a competing
	// species or they can't exist in the climate/environment.

	// TODO: Ensure that we favor subtypes that are found in a particular biome.
	// Like trees in forests, or grass in grasslands.

	// TODO: Select prey and predators and add species if species are missing
	// that are needed for the food chain / balance the ecosystem.

	// TODO: Proper fitness function per species or per ecosphere.
	sf := func(r int) float64 {
		return 1.0
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

	// Place n species on the map.
	for i := 0; i < n; i++ {
		b.PlaceSpecies(sf, distSeedFunc)
	}

	// DEBUG: Print all species.
	for _, s := range b.Species {
		fmt.Println(s)
	}
}

// PlaceSpecies places another species on the map at the region with the highest fitness score.
func (b *Bio) PlaceSpecies(sf func(int) float64, distSeedFunc func() []int) *Species {
	// Score all regions, pick highest score.
	var newspecies int
	lastMax := math.Inf(-1)
	for i, val := range b.CalcFitnessScore(sf, distSeedFunc) {
		if val > lastMax {
			newspecies = i
			lastMax = val
		}
	}
	tf := b.getTolerancesForRegionFunc()
	return b.placeSpeciesAt(newspecies, tf)
}

// PlaceSpeciesAt places a species at the given region.
// TODO: Allow specifying the species type/subtype?
func (b *Bio) PlaceSpeciesAt(r int) *Species {
	tf := b.getTolerancesForRegionFunc()
	return b.placeSpeciesAt(r, tf)
}

func (b *Bio) placeSpeciesAt(r int, tf func(int) *SpeciesTolerances) *Species {
	// TODO: Pick species type based on biome through a weighted random array.
	b.rand.Seed(b.Seed + int64(r))
	s := b.newSpecies(r, SpeciesKingdoms[b.rand.Intn(len(SpeciesKingdoms))], tf)
	b.Species = append(b.Species, s)
	return s
}

func (b *Bio) newSpecies(r int, t SpeciesKingdom, tf func(int) *SpeciesTolerances) *Species {
	// TODO: Get culture and language from the region and use it to generate the name.
	s := &Species{
		Origin:            r,
		Kingdom:           t,
		Size:              SpeciesSizes[b.rand.Intn(len(SpeciesSizes))],
		SpeciesTolerances: tf(r),
	}

	// Pick subtype and mode of locomotion.
	if s.Ecosphere.IsWater() {
		subTypes := speciesKingdomToFamiliesWater[s.Kingdom]
		s.Family = subTypes[b.rand.Intn(len(subTypes))]
		s.Locomotion = s.Family.Locomotion()
		// There is further a remote chance that we have another way of locomotion.
		if b.rand.Float64() < 0.01 {
			s.Locomotion |= LocomotionTypesWater[b.rand.Intn(len(LocomotionTypesWater))]
		}
	} else {
		subTypes := speciesKingdomToFamiliesLand[s.Kingdom]
		s.Family = subTypes[b.rand.Intn(len(subTypes))]
		s.Locomotion = s.Family.Locomotion()
		// There is further a remote chance that we have another way of locomotion.
		if b.rand.Float64() < 0.02 {
			s.Locomotion |= LocomotionTypesLand[b.rand.Intn(len(LocomotionTypesLand))]
		}
	}

	// Pick a random type of prey.
	digestiveSystems := s.Kingdom.DigestiveSystems()
	s.Digestion = digestiveSystems[b.rand.Intn(len(digestiveSystems))]

	return s
}

func (b *Bio) getSpeciesScores(s *Species) []float64 {
	scores := make([]float64, b.mesh.numRegions)
	tsf := b.getToleranceScoreFunc(s.SpeciesTolerances)
	for i := range scores {
		scores[i] = tsf(i)
	}
	return scores
}

type Species struct {
	Name       string
	Origin     int             // The region where the species originated, acts as a seed.
	Kingdom    SpeciesKingdom  // General type of the species.
	Family     SpeciesFamily   // Subtype of the species.
	Digestion  DigestiveSystem // What kind of food the species can eat.
	Size       SpeciesSize     // Size of the species.
	Locomotion Locomotion      // How the species moves. (TODO: Primary locomotion)
	*SpeciesTolerances
}

func (s *Species) String() string {
	str := fmt.Sprintf("%s (%s, %s %s), lives at %s", s.Name, s.Kingdom, s.Size, s.Family, s.Ecosphere)
	if s.Locomotion != LocomotionNone {
		str += fmt.Sprintf(", can: %s", s.Locomotion)
	}
	str += fmt.Sprintf(", digestion: %s", s.Digestion)
	str += fmt.Sprintf("\n%s", s.SpeciesTolerances.String())
	return str
}

type Locomotion byte

func (l Locomotion) isSet(b Locomotion) bool {
	return l&b != 0
}

const (
	LocomotionNone Locomotion = 0
	LocomotionFly  Locomotion = 1 << iota
	LocomotionBurrow
	LocomotionWalk
	LocomotionSwim
	LocomotionClimb
	LocomotionSlither
)

func (l Locomotion) String() string {
	var strs []string
	if l.isSet(LocomotionFly) {
		strs = append(strs, "fly")
	}
	if l.isSet(LocomotionBurrow) {
		strs = append(strs, "burrow")
	}
	if l.isSet(LocomotionWalk) {
		strs = append(strs, "walk")
	}
	if l.isSet(LocomotionSwim) {
		strs = append(strs, "swim")
	}
	if l.isSet(LocomotionClimb) {
		strs = append(strs, "climb")
	}
	if l.isSet(LocomotionSlither) {
		strs = append(strs, "slither")
	}
	return strings.Join(strs, ", ")
}

var LocomotionTypes = []Locomotion{
	LocomotionFly,
	LocomotionBurrow,
	LocomotionWalk,
	LocomotionSwim,
	LocomotionClimb,
	LocomotionSlither,
}

var LocomotionTypesLand = []Locomotion{
	LocomotionFly,
	LocomotionBurrow,
	LocomotionWalk,
	LocomotionClimb,
	LocomotionSlither,
}

var LocomotionTypesWater = []Locomotion{
	LocomotionBurrow,
	LocomotionSwim,
	LocomotionClimb,
	LocomotionSlither,
}

type SpeciesKingdom int

const (
	SpeciesKingdomFlora SpeciesKingdom = iota
	SpeciesKingdomFauna                // Maybe split this up into different types of fauna?
	SpeciesKingdomFunga
)

func (s SpeciesKingdom) String() string {
	switch s {
	case SpeciesKingdomFlora:
		return "flora"
	case SpeciesKingdomFauna:
		return "fauna"
	case SpeciesKingdomFunga:
		return "funga"
	}
	return "unknown"
}

func (s SpeciesKingdom) DigestiveSystems() []DigestiveSystem {
	switch s {
	case SpeciesKingdomFlora:
		return []DigestiveSystem{
			// TODO: Allow weighted selection. Some plants can eat other plants or animals.
			DigestivePhotosynthetic,
		}
	case SpeciesKingdomFauna:
		return []DigestiveSystem{
			DigestiveSystemCarnivore,
			DigestiveSystemHerbivore,
			DigestiveSystemOmnivore,
		}
	case SpeciesKingdomFunga:
		return []DigestiveSystem{
			DigestivePhotosynthetic,
			DigestiveDecomposer,
			DigestiveSystemCarnivore, // rare
		}
	}
	return nil
}

var SpeciesKingdoms = []SpeciesKingdom{
	SpeciesKingdomFauna,
	SpeciesKingdomFlora,
	SpeciesKingdomFunga,
}

type DigestiveSystem int

const (
	DigestiveSystemCarnivore DigestiveSystem = iota
	DigestiveSystemHerbivore
	DigestiveSystemOmnivore
	DigestivePhotosynthetic
	DigestiveDecomposer
)

func (d DigestiveSystem) String() string {
	switch d {
	case DigestiveSystemCarnivore:
		return "carnivore"
	case DigestiveSystemHerbivore:
		return "herbivore"
	case DigestiveSystemOmnivore:
		return "omnivore"
	case DigestivePhotosynthetic:
		return "photosynthetic"
	case DigestiveDecomposer:
		return "decomposer"
	}
	return "unknown"
}

var DigestiveSystems = []DigestiveSystem{
	DigestiveSystemCarnivore,
	DigestiveSystemHerbivore,
	DigestiveSystemOmnivore,
	DigestivePhotosynthetic,
	DigestiveDecomposer,
}

type SpeciesFamily int

const (
	SpeciesFamilyNone SpeciesFamily = iota
	SpeciesFamilyTree
	SpeciesFamilyShrub
	SpeciesFamilyGrass
	SpeciesFamilyReed
	SpeciesFamilyHerb
	SpeciesFamilyFlower
	SpeciesFamilyFern
	SpeciesFamilyMoss
	SpeciesFamilyVine
	SpeciesFamilyCactus
	SpeciesFamilySucculent
	SpeciesFamilyInsect
	SpeciesFamilyArachnid
	SpeciesFamilyMammal
	SpeciesFamilyBird
	SpeciesFamilyFish
	SpeciesFamilyCrustacean
	SpeciesFamilyMollusk
	SpeciesFamilyMolluskClam
	SpeciesFamilyMolluskSnail
	SpeciesFamilyAmphibian
	SpeciesFamilyReptileSerpent
	SpeciesFamilyReptileLizard
	SpeciesFamilyRodent
	SpeciesFamilyWorm
	SpeciesFamilyMushroom
	SpeciesFamilyMold
)

func (s SpeciesFamily) String() string {
	switch s {
	case SpeciesFamilyNone:
		return "none"
	case SpeciesFamilyTree:
		return "tree"
	case SpeciesFamilyShrub:
		return "shrub"
	case SpeciesFamilyGrass:
		return "grass"
	case SpeciesFamilyReed:
		return "reed"
	case SpeciesFamilyHerb:
		return "herb"
	case SpeciesFamilyFlower:
		return "flower"
	case SpeciesFamilyFern:
		return "fern"
	case SpeciesFamilyMoss:
		return "moss"
	case SpeciesFamilyVine:
		return "vine"
	case SpeciesFamilyCactus:
		return "cactus"
	case SpeciesFamilySucculent:
		return "succulent"
	case SpeciesFamilyInsect:
		return "insect"
	case SpeciesFamilyArachnid:
		return "arachnid"
	case SpeciesFamilyMammal:
		return "mammal"
	case SpeciesFamilyBird:
		return "bird"
	case SpeciesFamilyFish:
		return "fish"
	case SpeciesFamilyCrustacean:
		return "crustacean"
	case SpeciesFamilyMollusk:
		return "mollusk"
	case SpeciesFamilyMolluskClam:
		return "clam"
	case SpeciesFamilyMolluskSnail:
		return "snail"
	case SpeciesFamilyAmphibian:
		return "amphibian"
	case SpeciesFamilyReptileSerpent:
		return "serpent"
	case SpeciesFamilyReptileLizard:
		return "lizard"
	case SpeciesFamilyRodent:
		return "rodent"
	case SpeciesFamilyWorm:
		return "worm"
	case SpeciesFamilyMushroom:
		return "mushroom"
	case SpeciesFamilyMold:
		return "mold"
	}
	return "unknown"
}

func (s SpeciesFamily) Locomotion() Locomotion {
	switch s {
	case SpeciesFamilyInsect:
		return LocomotionWalk
	case SpeciesFamilyArachnid:
		return LocomotionWalk
	case SpeciesFamilyMammal:
		return LocomotionWalk
	case SpeciesFamilyBird:
		return LocomotionFly
	case SpeciesFamilyFish:
		return LocomotionSwim
	case SpeciesFamilyCrustacean:
		return LocomotionSwim
	case SpeciesFamilyMollusk:
		return LocomotionSwim | LocomotionWalk
	case SpeciesFamilyMolluskSnail:
		return LocomotionSlither | LocomotionClimb
	case SpeciesFamilyAmphibian:
		return LocomotionWalk | LocomotionSwim
	case SpeciesFamilyReptileSerpent:
		return LocomotionSlither
	case SpeciesFamilyReptileLizard:
		return LocomotionWalk | LocomotionClimb
	case SpeciesFamilyRodent:
		return LocomotionWalk | LocomotionClimb | LocomotionBurrow
	case SpeciesFamilyWorm:
		return LocomotionSlither | LocomotionBurrow
	}
	return LocomotionNone
}

var speciesKingdomToFamiliesLand = map[SpeciesKingdom][]SpeciesFamily{
	SpeciesKingdomFlora: {
		SpeciesFamilyTree,
		SpeciesFamilyShrub,
		SpeciesFamilyGrass,
		SpeciesFamilyReed,
		SpeciesFamilyHerb,
		SpeciesFamilyFlower,
		SpeciesFamilyFern,
		SpeciesFamilyMoss,
		SpeciesFamilyVine,
		SpeciesFamilyCactus,
		SpeciesFamilySucculent,
	},
	SpeciesKingdomFauna: {
		SpeciesFamilyInsect,
		SpeciesFamilyArachnid,
		SpeciesFamilyMammal,
		SpeciesFamilyBird,
		SpeciesFamilyAmphibian,
		SpeciesFamilyReptileSerpent,
		SpeciesFamilyReptileLizard,
		SpeciesFamilyMolluskSnail,
		SpeciesFamilyRodent,
		SpeciesFamilyMollusk,
	},
	SpeciesKingdomFunga: {
		SpeciesFamilyMushroom,
		SpeciesFamilyMold,
	},
}

var speciesKingdomToFamiliesWater = map[SpeciesKingdom][]SpeciesFamily{
	SpeciesKingdomFlora: {
		SpeciesFamilyGrass,
		SpeciesFamilyHerb,
	},
	SpeciesKingdomFauna: {
		SpeciesFamilyFish,
		SpeciesFamilyCrustacean,
		SpeciesFamilyMollusk,
		SpeciesFamilyMolluskClam,
		SpeciesFamilyReptileSerpent,
	},
	SpeciesKingdomFunga: {
		SpeciesFamilyMushroom,
	},
}

type SpeciesSize int

const (
	SpeciesSizeDefault SpeciesSize = iota
	SpeciesSizeTiny
	SpeciesSizeSmall
	SpeciesSizeMedium
	SpeciesSizeLarge
	SpeciesSizeHuge
)

func (s SpeciesSize) String() string {
	switch s {
	case SpeciesSizeDefault:
		return "default"
	case SpeciesSizeTiny:
		return "tiny"
	case SpeciesSizeSmall:
		return "small"
	case SpeciesSizeMedium:
		return "medium"
	case SpeciesSizeLarge:
		return "large"
	case SpeciesSizeHuge:
		return "huge"
	}
	return "unknown"
}

var SpeciesSizes = []SpeciesSize{
	SpeciesSizeTiny,
	SpeciesSizeSmall,
	SpeciesSizeMedium,
	SpeciesSizeLarge,
	SpeciesSizeHuge,
}

func easeInOutCubic(x float64) float64 {
	if x < 0.5 {
		return 4 * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 3)/2
}

type EcosphereType int

const (
	EcosphereTypeOcean EcosphereType = iota
	EcosphereTypeRiver
	EcosphereTypeLake
	EcosphereTypeLand
)

func (e EcosphereType) IsWater() bool {
	return e == EcosphereTypeOcean || e == EcosphereTypeRiver || e == EcosphereTypeLake
}

// getEcosphere returns the ecosphere of the given region.
func (b *Bio) getEcosphere(r int) EcosphereType {
	// Get the ecosphere we are in.
	if b.Elevation[r] <= 0.0 {
		return EcosphereTypeOcean
	}
	if b.isRegRiver(r) {
		return EcosphereTypeRiver
	}
	if b.isRegLake(r) {
		return EcosphereTypeLake
	}
	return EcosphereTypeLand
}

func (e EcosphereType) String() string {
	switch e {
	case EcosphereTypeOcean:
		return "ocean"
	case EcosphereTypeRiver:
		return "river"
	case EcosphereTypeLake:
		return "lake"
	case EcosphereTypeLand:
		return "land"
	}
	return "unknown"
}
