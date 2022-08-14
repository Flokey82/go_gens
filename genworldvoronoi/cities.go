package genworldvoronoi

import (
	"log"
)

// rCityScore calculates the fitness value for settlements for all regions.
func (m *Map) rCityScore() []float64 {
	cities := m.cities_r
	_, maxFlux := minMax(m.r_flux)
	score := make([]float64, m.mesh.numRegions)

	// Initialize fitness score with the normalized flux value.
	// This will favor placing cities along (and at the end of)
	// large rivers.
	for i, fl := range m.r_flux {
		score[i] = fl / maxFlux // originally: math.Sqrt(fl / maxFlux)
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

	// Calculate the fitness score for each region
	for i := 0; i < m.mesh.numRegions; i++ {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.r_elevation[i] <= 0 || m.r_pool[i] > 0 {
			score[i] = -999999.0
			continue
		}

		// Visit all neighbors and modify the score based on their properties.
		for _, nbs := range m.rNeighbors(i) {
			if m.r_elevation[nbs] <= 0 || m.r_pool[nbs] > 0 {
				// If a neighbor is below (or at) sea level, or a lake,
				// we increase the fitness value and reduce it by a fraction,
				// depending on the size of the lake or ocean it is part of.
				//
				// TODO: Improve this.
				score[i] += 0.09
				if m.r_waterbodies[nbs] >= 0 {
					// If nbs is part of a waterbody (ocean), we reduce the score by a constant factor.
					// The larger the waterbody, the smaller the penalty, which will favor larger waterbodies.
					score[i] -= 0.05 / (float64(m.r_waterbody_size[m.r_waterbodies[nbs]]) + 1e-9)
				} else if m.r_drainage[nbs] >= 0 {
					// If a drainage is set for nbs, it is part of a lake.
					// So we reduce the score by a constant factor, which is smaller, the larger the lake.
					score[i] -= 0.05 / (float64(m.r_lake_size[m.r_drainage[nbs]]) + 1e-9)
				}
				break
			}

			// TODO: Add bonus if near ocean or lake.
			// TODO: Consider biome
			// TODO: Consider sediment/fertility of land.
			// TODO: Add bonus for mountain proximity (mines, resources)
		}

		// Penalty for proximity / bonus for higher distance to other cities.
		//
		// We multiply the score by the distance to other cities, amplifying
		// positive scores.
		//
		// NOTE: Originally this was done with some constant values, which might be better
		// since we are here dependent on the current score we have assigned and cannot
		// recover an initially bad score caused by a low water flux.
		score[i] *= (float64(r_distance_c[i]) + 1e-9) // originally: -= 0.02 / (float64(r_distance_c[i]) + 1e-9)
	}
	return score
}

// rPlaceCity places another city at the region with the highest fitness score.
func (m *Map) rPlaceCity() {
	var newcity int
	lastMax := -999999.0
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
