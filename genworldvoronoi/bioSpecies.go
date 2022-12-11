package genworldvoronoi

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func (b *Bio) genNRandomSpecies(n int) {
	b.resetRand()

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
	s := b.newSpecies(newspecies, SpeciesKingdoms[b.rand.Intn(len(SpeciesKingdoms))])
	b.Species = append(b.Species, s)
	return s
}

// PlaceSpeciesAt places a species at the given region.
// TODO: Allow specifying the species type/subtype?
func (b *Bio) PlaceSpeciesAt(r int) *Species {
	// TODO: Pick species type based on biome through a weighted random array.
	s := b.newSpecies(r, SpeciesKingdoms[b.rand.Intn(len(SpeciesKingdoms))])
	b.Species = append(b.Species, s)
	return s
}

func (b *Bio) newSpecies(r int, t SpeciesKingdom) *Species {
	rnd := rand.New(rand.NewSource(b.Seed + int64(r)))

	// TODO: Get culture and language from the region and use it to generate the name.
	s := &Species{
		Origin:    r,
		Size:      SpeciesSizes[rnd.Intn(len(SpeciesSizes))],
		Kingdom:   t,
		Ecosphere: b.getEcosphere(r),
	}

	// minMaxRange returns a min and max range for the given value with the given variance.
	minMaxRange := func(val, minVal, maxVal, variation float64) [2]float64 {
		nVar := rnd.Float64() * variation
		varMin := nVar * rnd.Float64()
		newMin := math.Max(minVal, val-varMin)
		newMax := math.Min(maxVal, val+(nVar-varMin))
		return [2]float64{newMin, newMax}
	}

	// Prefered elevation range.
	minElev, maxElev := minMax(b.Elevation)
	if s.Ecosphere == EcosphereTypeOcean {
		s.ElevRange = minMaxRange(b.Elevation[r], minElev, 0, 0.2)
	} else {
		s.ElevRange = minMaxRange(b.Elevation[r], 0, maxElev, 0.2)
	}

	// Preferred temperature range.
	s.TempRange = minMaxRange(b.getRTemperature(r, maxElev), float64(minTemp), float64(maxTemp), 0.2)

	// Preferred humidity range.
	minHum, maxHum := minMax(b.Moisture)
	s.HumRange = minMaxRange(b.Moisture[r], minHum, maxHum, 0.2)

	// Pick subtype and mode of locomotion.
	switch s.Ecosphere {
	case EcosphereTypeOcean, EcosphereTypeLake, EcosphereTypeRiver:
		subTypes := speciesKingdomToFamiliesWater[s.Kingdom]
		s.Family = subTypes[rnd.Intn(len(subTypes))]
		s.Locomotion = SpeciesFamiliesToLocomotion[s.Family]
		// There is further a remote chance that we have another way of locomotion.
		if rnd.Float64() < 0.01 {
			s.Locomotion |= LocomotionTypesWater[rnd.Intn(len(LocomotionTypesWater))]
		}
	default:
		subTypes := speciesKingdomToFamiliesLand[s.Kingdom]
		s.Family = subTypes[rnd.Intn(len(subTypes))]
		s.Locomotion = SpeciesFamiliesToLocomotion[s.Family]
		// There is further a remote chance that we have another way of locomotion.
		if rnd.Float64() < 0.02 {
			s.Locomotion |= LocomotionTypesLand[rnd.Intn(len(LocomotionTypesLand))]
		}
	}

	// Pick a random type of prey.
	digestiveSystems := SpeciesKingdomToDigestiveSystems[s.Kingdom]
	s.Digestion = digestiveSystems[rnd.Intn(len(digestiveSystems))]

	// If we are not in the ocean, we probably have a preferred biome.
	if s.Ecosphere != EcosphereTypeOcean && rnd.Float64() < 0.7 {
		s.PreferredBiomes = []int{b.getRWhittakerModBiomeFunc()(r)}
	}

	return s
}

func (b *Bio) getSpeciesScores(s *Species) []float64 {
	scores := make([]float64, b.mesh.numRegions)
	_, maxElev := minMax(b.Elevation)
	bf := b.getRWhittakerModBiomeFunc()
	for i := range scores {
		scores[i] = b.getSpeciesScore(s, i, maxElev, bf)
	}
	return scores
}

func (b *Bio) getSpeciesScore(s *Species, r int, maxElev float64, bf func(int) int) float64 {
	// Check what ecosphere we are in and if it matches the species.
	if s.Ecosphere != b.getEcosphere(r) {
		return -1.0
	}

	// Check if we require a specific biome.
	if len(s.PreferredBiomes) > 0 && !isInIntList(s.PreferredBiomes, bf(r)) {
		return -1
	}

	// Check how much we diverge from the preferred temperature range.
	tempDiv := getDiversionFromRange(b.getRTemperature(r, maxElev), s.TempRange)
	tempScore := easeInOutCubic(1 - tempDiv/(s.TempRange[1]-s.TempRange[0]))

	// Check how much we diverge from the preferred humidity range.
	humDiv := getDiversionFromRange(b.Moisture[r], s.HumRange)
	humScore := easeInOutCubic(1 - humDiv/(s.HumRange[1]-s.HumRange[0]))

	// Check how much we diverge from the preferred elevation range.
	elevDiv := getDiversionFromRange(b.Elevation[r], s.ElevRange)
	elevScore := easeInOutCubic(1 - elevDiv/(s.ElevRange[1]-s.ElevRange[0]))

	return tempScore * humScore * elevScore
}

type Species struct {
	Name            string
	Origin          int             // The region where the species originated, acts as a seed.
	Kingdom         SpeciesKingdom  // General type of the species.
	Family          SpeciesFamily   // Subtype of the species.
	Digestion       DigestiveSystem // What kind of food the species can eat.
	Size            SpeciesSize     // Size of the species.
	Locomotion      Locomotion      // How the species moves. (TODO: Primary locomotion)
	Ecosphere       EcosphereType   // Ocean, River, Land, Lake
	TempRange       [2]float64      // Min, Max temperature
	HumRange        [2]float64      // Min, Max humidity
	ElevRange       [2]float64      // Min, Max elevation
	PreferredBiomes []int           // Only applies to non-marine species.
}

func (s *Species) String() string {
	str := fmt.Sprintf("%s (%s, %s %s), lives at %s", s.Name, s.Kingdom, s.Size, s.Family, s.Ecosphere)
	if s.Locomotion != LocomotionNone {
		str += fmt.Sprintf(", can: %s", s.Locomotion)
	}
	if len(s.PreferredBiomes) > 0 {
		str += fmt.Sprintf(", preferred biomes: %v", s.PreferredBiomes)
	}
	str += fmt.Sprintf(", temperature: %.2f-%.2f, humidity: %.2f-%.2f, elevation: %.2f-%.2f",
		s.TempRange[0], s.TempRange[1],
		s.HumRange[0], s.HumRange[1],
		s.ElevRange[0], s.ElevRange[1])
	str += fmt.Sprintf(", digestion: %s", s.Digestion)
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

var SpeciesKingdoms = []SpeciesKingdom{
	SpeciesKingdomFauna,
	SpeciesKingdomFlora,
	SpeciesKingdomFunga,
}

var SpeciesKingdomToDigestiveSystems = map[SpeciesKingdom][]DigestiveSystem{
	SpeciesKingdomFlora: {
		// TODO: Allow weighted selection. Some plants can eat other plants or animals.
		DigestivePhotosynthetic,
	},
	SpeciesKingdomFauna: {
		DigestiveSystemCarnivore,
		DigestiveSystemHerbivore,
		DigestiveSystemOmnivore,
	},
	SpeciesKingdomFunga: {
		DigestivePhotosynthetic,
		DigestiveDecomposer,
		DigestiveSystemCarnivore, // rare
	},
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

var SpeciesFamiliesToLocomotion = map[SpeciesFamily]Locomotion{
	SpeciesFamilyInsect:         LocomotionWalk,
	SpeciesFamilyArachnid:       LocomotionWalk,
	SpeciesFamilyMammal:         LocomotionWalk,
	SpeciesFamilyBird:           LocomotionFly,
	SpeciesFamilyFish:           LocomotionSwim,
	SpeciesFamilyCrustacean:     LocomotionSwim,
	SpeciesFamilyMollusk:        LocomotionSwim | LocomotionWalk,
	SpeciesFamilyMolluskSnail:   LocomotionSlither | LocomotionClimb,
	SpeciesFamilyAmphibian:      LocomotionWalk | LocomotionSwim,
	SpeciesFamilyReptileSerpent: LocomotionSlither,
	SpeciesFamilyReptileLizard:  LocomotionWalk | LocomotionClimb,
	SpeciesFamilyRodent:         LocomotionWalk | LocomotionClimb | LocomotionBurrow,
	SpeciesFamilyWorm:           LocomotionSlither | LocomotionBurrow,
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

type EcosphereType int

const (
	EcosphereTypeOcean EcosphereType = iota
	EcosphereTypeRiver
	EcosphereTypeLake
	EcosphereTypeLand
)

// getEcosphere returns the ecosphere of the given region.
func (b *Bio) getEcosphere(r int) EcosphereType {
	// Get the ecosphere we are in.
	if b.Elevation[r] <= 0.0 {
		return EcosphereTypeOcean
	}
	if b.isRiver(r) {
		return EcosphereTypeRiver
	}
	if b.isLake(r) {
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
