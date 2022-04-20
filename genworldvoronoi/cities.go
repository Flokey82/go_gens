package genworldvoronoi

import (
	"math"
)

func (m *Map) rCityScore() []float64 {
	cities := m.cities_r
	_, maxFlux := minMax(m.r_flux)
	score := make([]float64, m.mesh.numRegions)
	for i, fl := range m.r_flux {
		score[i] = math.Sqrt(fl / maxFlux)
	}

	r_distance_c := m.assignDistanceField(cities, make(map[int]bool))
	for i := 0; i < m.mesh.numRegions; i++ {
		if m.r_elevation[i] <= 0 {
			score[i] = -999999.0
			continue
		}

		// TODO: Add bonus if near ocean or lake.
		// TODO: Consider biome
		// TODO: Consider sediment/fertility of land.

		// Penalty for proximity to other cities.
		for j := 0; j < len(cities); j++ {
			score[i] -= 0.05 / (float64(r_distance_c[i]) + 1e-9)
		}
		for _, nbs := range m.rNeighbors(i) {
			if m.r_elevation[nbs] <= 0 {
				score[nbs] += 0.09
			}
		}
	}
	return score
}

func (m *Map) rPlaceCity() {
	score := m.rCityScore()
	var newcity int
	lastMax := -999999.0
	for i, val := range score {
		if val > lastMax {
			newcity = i
			lastMax = val
		}
	}
	m.cities_r = append(m.cities_r, newcity)
}

func (m *Map) rPlaceNCities(n int) {
	for i := 0; i < n; i++ {
		m.rPlaceCity()
	}
}
