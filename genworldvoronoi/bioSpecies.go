package genworldvoronoi

import (
	"fmt"
	"math"
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

func (b *Bio) placeSpeciesAt(r int, tf func(int) SpeciesTolerances) *Species {
	// TODO: Pick species type based on biome through a weighted random array.
	b.rand.Seed(b.Seed + int64(r))
	s := b.newSpecies(r, SpeciesKingdoms[b.rand.Intn(len(SpeciesKingdoms))], tf)
	b.Species = append(b.Species, s)
	return s
}

func (b *Bio) newSpecies(r int, t SpeciesKingdom, tf func(int) SpeciesTolerances) *Species {
	// TODO: Get culture and language from the region and use it to generate the name.
	s := &Species{
		Origin: r,
		SpeciesProperties: SpeciesProperties{
			Kingdom: t,
			Size:    SpeciesSizes[b.rand.Intn(len(SpeciesSizes))],
		},
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
	Name   string
	Origin int // The region where the species originated, acts as a seed.
	SpeciesProperties
	SpeciesTolerances
}

func (s *Species) String() string {
	return fmt.Sprintf("%s (%s) %s\n", s.Name, s.SpeciesProperties.String(), s.SpeciesTolerances.String())
}
