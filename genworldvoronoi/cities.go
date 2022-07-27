package genworldvoronoi

import (
	"log"
)

// rCityScore calculates the fitness value for settlements for all regions.
func (m *Map) rCityScore() []float64 {
	cities := m.cities_r
	_, maxFlux := minMax(m.r_flux)
	score := make([]float64, m.mesh.numRegions)
	for i, fl := range m.r_flux {
		score[i] = fl / maxFlux //math.Sqrt(fl / maxFlux)
	}

	// TODO: Create different fitness functions for different types of settlement.
	//   - Capital
	//   - Cities / Settlements
	//     ) Proximity to capital!
	//   - Agricultural
	//   - Mining
	//   - ...
	r_distance_c := m.assignDistanceField(cities, make(map[int]bool))
	for i := 0; i < m.mesh.numRegions; i++ {
		if m.r_elevation[i] <= 0 || m.r_pool[i] > 0 {
			score[i] = -999999.0
			continue
		}

		// TODO: Add bonus if near ocean or lake.
		// TODO: Consider biome
		// TODO: Consider sediment/fertility of land.
		// TODO: Add bonus for mountain proximity (mines, resources)

		for _, nbs := range m.rNeighbors(i) {
			// TODO: Take in account size of ocean / lake.
			if m.r_elevation[nbs] <= 0 || m.r_pool[nbs] > 0 {
				score[i] += 0.09
				if m.r_waterbodies[nbs] >= 0 {
					score[i] -= 0.05 / (float64(m.r_waterbody_size[m.r_waterbodies[nbs]]) + 1e-9)
				} else if m.r_drainage[nbs] >= 0 {
					score[i] -= 0.05 / (float64(m.r_lake_size[m.r_drainage[nbs]]) + 1e-9)
				}
				break
			}
		}

		// Penalty for proximity to other cities.
		score[i] *= (float64(r_distance_c[i]) + 1e-9)
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
