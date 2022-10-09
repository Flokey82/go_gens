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
	_, maxSteepness := minMax(steepness)

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
			if m.r_elevation[nb] <= 0 || m.r_pool[nb] > 0 {
				// We only apply this bonus once.
				if hasWaterBodyBonus {
					continue
				}
				hasWaterBodyBonus = true
				// If a neighbor is below (or at) sea level, or a lake,
				// we increase the fitness value and reduce it by a fraction,
				// depending on the size of the lake or ocean it is part of.
				//
				// TODO: Improve this.
				scoreDelta := 0.5
				if wbIdx := m.r_waterbodies[nb]; wbIdx >= 0 && m.r_waterbody_size[wbIdx] > 0 {
					// If nb is part of a waterbody (ocean), we reduce the score by a constant factor.
					// The larger the waterbody, the smaller the penalty, which will favor larger waterbodies.
					scoreDelta -= 0.5 / (float64(m.r_waterbody_size[wbIdx]) + 1e-9)
				} else if drIdx := m.r_drainage[nb]; drIdx >= 0 && m.r_lake_size[drIdx] > 0 {
					// If a drainage is set for nb, it is part of a lake.
					// So we reduce the score by a constant factor, which is smaller, the larger the lake.
					scoreDelta -= 0.5 / (float64(m.r_lake_size[drIdx]) + 1e-9)
				} else {
					scoreDelta = 0.0
				}
				score[i] += scoreDelta
			} else {
				// If the sourrounding terrain is flat, we get a bonus.
				score[i] += 0.5 * (1.0 - steepness[nb]) / float64(len(nbs))
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
		score[i] *= r_distance_c[i] / maxDistC // originally: -= 0.02 / (float64(r_distance_c[i]) + 1e-9)

		// The steeper the terrain, the less likely it is to be settled.
		score[i] *= 1.0 - (steepness[i] / maxSteepness)
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
