package genworldvoronoi

import (
	"fmt"
	"strings"
)

// SpeciesProperties defines the properties of a species.
type SpeciesProperties struct {
	Kingdom    SpeciesKingdom  // General type of the species.
	Family     SpeciesFamily   // Subtype of the species.
	Digestion  DigestiveSystem // What kind of food the species can eat.
	Size       SpeciesSize     // Size of the species.
	Locomotion Locomotion      // How the species moves. (TODO: Primary locomotion)
}

// CompetitorHash returns a hash that can be used to compare two species and
// determine if they are competitors.
// This will be helpful when spreading the species and avoiding competing species.
func (s *SpeciesProperties) CompetitorHash() int64 {
	return int64(s.Size) | (int64(s.Locomotion) << 8) | (int64(s.Digestion) << 16)
}

func (s *SpeciesProperties) String() string {
	return fmt.Sprintf("%s %s %s %s %s", s.Kingdom, s.Family, s.Digestion, s.Size, s.Locomotion)
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

type SpeciesSize byte

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

type DigestiveSystem byte

const (
	DigestiveSystemCarnivore DigestiveSystem = iota
	DigestiveSystemHerbivore
	DigestiveSystemOmnivore
	DigestivePhotosynthetic
	DigestiveDecomposer
	DigestiveParasitic
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
	case DigestiveParasitic:
		return "parasitic"
	}
	return "unknown"
}

var DigestiveSystems = []DigestiveSystem{
	DigestiveSystemCarnivore,
	DigestiveSystemHerbivore,
	DigestiveSystemOmnivore,
	DigestivePhotosynthetic,
	DigestiveDecomposer,
	DigestiveParasitic,
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
