package genworldvoronoi

import (
	"log"
	"math"
)

// rCityScore calculates the fitness value for settlements for all regions.
func (m *Map) rCityScore() []float64 {
	cities := m.cities_r
	_, maxFlux := minMax(m.r_flux)
	score := make([]float64, m.mesh.numRegions)
	steepness := m.getRSteepness()

	// Initialize fitness score with the normalized flux value.
	// This will favor placing cities along (and at the end of)
	// large rivers.
	for i, fl := range m.r_flux {
		// Skip all regions below sea level.
		if m.r_elevation[i] <= 0 {
			continue
		}
		score[i] = math.Sqrt(fl / maxFlux) // originally: math.Sqrt(fl / maxFlux)
	}

	// TODO: Create different fitness functions for different types of settlement.
	//   - Capital
	//   - Cities / Settlements
	//     ) Proximity to capital!
	//   - Agricultural
	//   - Mining
	//   - ...

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

		// Visit all neighbors and modify the score based on their properties.
		var hasWaterBodyBonus bool
		nbs := m.rNeighbors(i)
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
					score[i] += 0.55 * (1 - 1/(float64(wbSize)+1e-9))
				}
			} else {
				// If the sourrounding terrain is flat, we get a bonus.
				stp := steepness[nb]
				score[i] += 0.5 * (1.0 - stp*stp) / float64(len(nbs))
			}

			// TODO:
			// - Consider biome
			// - Consider sediment/fertility of land.
			// - Add bonus for mountain proximity (mines, resources)
		}

		// Penalty for proximity / bonus for higher distance to other cities.
		//
		// We multiply the score by the distance to other cities, amplifying
		// positive scores.
		//
		// NOTE: Originally this was done with some constant values, which might be better
		// since we are here dependent on the current score we have assigned and cannot
		// recover an initially bad score caused by a low water flux.
		if math.IsInf(r_distance_c[i], 0) {
			continue
		}
		dist := (r_distance_c[i] / maxDistC)
		score[i] *= dist * dist // originally: -= 0.02 / (float64(r_distance_c[i]) + 1e-9)

		// The steeper the terrain, the less likely it is to be settled.
		stp := steepness[i]
		score[i] *= 1.0 - (stp * stp)

		// TODO: Bonus for trade routes.
	}
	return score
}

// rPlaceCity places another city at the region with the highest fitness score.
func (m *Map) rPlaceCity() {
	var newcity int
	lastMax := math.Inf(-1)
	for i, val := range m.rCityScore() {
		if val > lastMax {
			newcity = i
			lastMax = val
		}
	}
	m.cities_r = append(m.cities_r, newcity)
}

// rPlaceNCities places n cities with the highest fitness scores.
func (m *Map) rPlaceNCities(n int) {
	for i := 0; i < n; i++ {
		log.Println("placing city", i)
		m.rPlaceCity()
	}
}
