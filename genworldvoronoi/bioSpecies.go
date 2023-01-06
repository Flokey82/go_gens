package genworldvoronoi

import (
	"container/heap"
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

func (b *Bio) expandSpecies() []int {
	// For now, let's just do this the dumb way.
	// TODO: Species with different competition hashes should be able to coexist in
	// the same region?
	// We might need to create a full index of all regions for each unique
	// competition hash.... or, which is more wasteful, per species.
	var seedPoints []int
	originToSpecFit := make(map[int]func(int) float64)
	for _, s := range b.Species {
		seedPoints = append(seedPoints, s.Origin)
		originToSpecFit[s.Origin] = b.getToleranceScoreFunc(s.SpeciesTolerances)
	}
	var queue ascPriorityQueue
	heap.Init(&queue)
	outReg := make([]int, 0, 8)

	// Get maxFlux and maxElev for normalizing.
	_, maxFlux := minMax(b.Flux)
	_, maxElev := minMax(b.Elevation)

	// TODO: Move this to a generic function.
	terrainWeight := func(o, u, v int) float64 {
		// Don't cross from water to land and vice versa.
		if (b.Elevation[u] > 0) != (b.Elevation[v] > 0) {
			return -1
		}

		// Calculate horizontal distance.
		ulat := b.LatLon[u][0]
		ulon := b.LatLon[u][1]
		vlat := b.LatLon[v][0]
		vlon := b.LatLon[v][1]
		horiz := haversine(ulat, ulon, vlat, vlon) / (2 * math.Pi)

		// Calculate vertical distance.
		vert := (b.Elevation[v] - b.Elevation[u]) / maxElev
		if vert > 0 {
			vert /= 10
		}
		diff := 1 + 0.25*math.Pow(vert/horiz, 2)

		// NOTE: Flux should only apply to animals since plants and fungi
		// don't need to worry about drowning.
		diff += 100 * math.Sqrt(b.Flux[u]/maxFlux)
		if b.Elevation[u] <= 0 {
			diff = 100
		}
		return horiz * diff
	}

	weight := func(o, u, v int) float64 {
		// Call species specific fitness function.
		sFit := originToSpecFit[o](v)
		if sFit < 0 {
			return -1
		}
		// Call terrain specific fitness function.
		tFit := terrainWeight(o, u, v)
		if tFit < 0 {
			return -1
		}
		return tFit * sFit
	}

	// 'terr' will hold a mapping of region to species.
	// The territory ID is the region number of the species origin.
	terr := initRegionSlice(b.mesh.numRegions)
	for i := 0; i < len(seedPoints); i++ {
		terr[seedPoints[i]] = seedPoints[i]
		for _, v := range b.mesh.r_circulate_r(outReg, seedPoints[i]) {
			newdist := weight(seedPoints[i], seedPoints[i], v)
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueEntry{
				score:       newdist,
				origin:      seedPoints[i],
				destination: v,
			})
		}
	}

	// Extend territories until the queue is empty.
	for queue.Len() > 0 {
		u := heap.Pop(&queue).(*queueEntry)
		if terr[u.destination] >= 0 {
			continue
		}
		terr[u.destination] = u.origin
		for _, v := range b.mesh.r_circulate_r(outReg, u.destination) {
			if terr[v] >= 0 {
				continue
			}
			newdist := weight(u.origin, u.destination, v)
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueEntry{
				score:       u.score + newdist,
				origin:      u.origin,
				destination: v,
			})
		}
	}
	return terr
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
