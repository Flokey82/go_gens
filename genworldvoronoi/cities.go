package genworldvoronoi

import (
	"log"
	"math"
)

const (
	TownTypeDefault = "town"
	TownTypeTrading = "trading"
	TownTypeMining  = "mining"
	TownTypeFarming = "agricultural"
)

type City struct {
	R     int     // Region where the city is located
	Score float64 // Score of the fitness function
	Type  string  // Type of city
}

// rPlaceNCities places n cities with the highest fitness scores.
func (m *Map) rPlaceNCities(n int, cType string) {
	var sf func(int) float64
	switch cType {
	case TownTypeDefault:
		sf = m.getFitnessCityDefault()
	case TownTypeTrading:
		sf = m.getFitnessTradingTowns()
	case TownTypeMining:
		sf = m.getFitnessSteepMountains()
	case TownTypeFarming:
		sf = m.getFitnessArableLand()
	default:
		return
	}

	for i := 0; i < n; i++ {
		log.Println("placing "+cType+" city", i)
		m.rPlaceCity(cType, sf)
	}
}

// rPlaceCity places another city at the region with the highest fitness score.
func (m *Map) rPlaceCity(cType string, sf func(int) float64) {
	var newcity int
	lastMax := math.Inf(-1)
	for i, val := range m.rCityScore(sf) {
		if val > lastMax {
			newcity = i
			lastMax = val
		}
	}
	m.cities_r = append(m.cities_r, &City{
		R:     newcity,
		Score: lastMax,
		Type:  cType,
	})
}

// rCityScore calculates the fitness value for settlements for all regions.
func (m *Map) rCityScore(sf func(int) float64) []float64 {
	// TODO: Create different fitness functions for different types of settlement.
	//   - Capital
	//   - Cities / Settlements
	//     ) Proximity to capital!
	//   - Agricultural
	//   - Mining
	//   - ...

	var cities []int
	for _, c := range m.cities_r {
		cities = append(cities, c.R)
	}
	score := make([]float64, m.mesh.numRegions)

	// Get distance to other cities.
	r_distance_c := m.assignDistanceField(cities, make(map[int]bool))
	_, maxDistC := minMax(r_distance_c)

	// Calculate the fitness score for each region
	for i := 0; i < m.mesh.numRegions; i++ {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.r_elevation[i] <= 0 || m.r_pool[i] > 0 {
			score[i] = -1.0
			continue
		}
		score[i] = sf(i)

		// Penalty for proximity / bonus for higher distance to other cities.
		//
		// We multiply the score by the distance to other cities, amplifying
		// positive scores.
		//
		// NOTE: Originally this was done with some constant values, which might be better
		// since we are here dependent on the current score we have assigned and cannot
		// recover an initially bad score caused by a low water flux.

		// TODO: Allow for distance to other cities / city types to be a parameter.
		if math.IsInf(r_distance_c[i], 0) {
			continue
		}
		dist := (r_distance_c[i] / maxDistC)
		score[i] *= dist * dist // originally: -= 0.02 / (float64(r_distance_c[i]) + 1e-9)
	}
	return score
}

func (m *Map) getFitnessTradingTowns() func(int) float64 {
	// TODO: Fix this.
	// I think this function should avoid the penalty wrt.
	// proximity to towns of other types.
	_, connecting := m.getTradeRoutes()
	return func(r int) float64 {
		return float64(len(connecting[r]))
	}
}

func (m *Map) getFitnessCityDefault() func(int) float64 {
	_, maxFlux := minMax(m.r_flux)
	steepness := m.getRSteepness()

	return func(i int) float64 {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.r_elevation[i] <= 0 || m.r_pool[i] > 0 {
			return -1.0
		}

		// Visit all neighbors and modify the score based on their properties.
		var hasWaterBodyBonus bool
		nbs := m.rNeighbors(i)

		// Initialize fitness score with the normalized flux value.
		// This will favor placing cities along (and at the end of)
		// large rivers.
		score := math.Sqrt(m.r_flux[i] / maxFlux)
		for _, nb := range nbs {
			// Add bonus if near ocean or lake.
			if m.isRBelowOrAtSeaLevelOrPool(nb) {
				// We only apply this bonus once.
				if hasWaterBodyBonus {
					continue
				}
				// If a neighbor is below (or at) sea level, or a lake,
				// we increase the fitness value and reduce it by a fraction,
				// depending on the size of the lake or ocean it is part of.
				//
				// TODO: Improve this.

				// If nb is part of a waterbody (ocean) or lake, we reduce the score by a constant factor.
				// The larger the waterbody/lake, the smaller the penalty, which will favor larger waterbodies.
				if wbSize := m.getRLakeOrWaterBodySize(nb); wbSize > 0 {
					hasWaterBodyBonus = true
					score += 0.55 * (1 - 1/(float64(wbSize)+1e-9))
				}
			} else {
				// If the sourrounding terrain is flat, we get a bonus.
				stp := steepness[nb]
				score += 0.5 * (1.0 - stp*stp) / float64(len(nbs))
			}

			// TODO:
			// - Consider biome
			// - Consider sediment/fertility of land.
			// - Add bonus for mountain proximity (mines, resources)
		}

		// The steeper the terrain, the less likely it is to be settled.
		// TODO: Bonus for trade routes.
		stp := steepness[i]
		score *= 1.0 - (stp * stp)
		return score
	}
}
