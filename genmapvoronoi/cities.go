package genmapvoronoi

import (
	"math"

	"github.com/Flokey82/go_gens/vmesh"
)

// cityScore calculates the fitness value for settlements for all regions.
func cityScore(render *Terrain) *vmesh.Heightmap {
	h := render.h
	cities := render.cities

	// Initialize fitness score with the suqare root of the flux value.
	// This will favor placing cities along (and at the end of)
	// large rivers.
	score := getFlux(h).MapF(math.Sqrt)
	for i := 0; i < h.Len(); i++ {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if h.Values[i] <= 0 || h.IsNearEdge(i) {
			score.Values[i] = -999999.0
			continue
		}

		// TODO: Add bonus if near ocean or lake.
		// TODO: Consider sediment/fertility of land.

		// Prefer points towards the middle of the map.
		score.Values[i] += 0.01 / (1e-9 + math.Abs(h.Vertices[i].X) - h.Extent.Width/2)
		score.Values[i] += 0.01 / (1e-9 + math.Abs(h.Vertices[i].Y) - h.Extent.Height/2)

		// Penalty for proximity / bonus for higher distance to other cities.
		for j := 0; j < len(cities); j++ {
			score.Values[i] -= 0.02 / (h.Distance(cities[j], i) + 1e-9)
		}

		// Visit all neighbors and modify the score based on their properties.
		for _, nbs := range h.Neighbours(i) {
			if h.Values[nbs] <= 0 {
				// If a neighbor is below (or at) sea level, or a lake,
				// we increase the fitness value and reduce it by a fraction,
				// depending on the size of the lake or ocean it is part of.
				score.Values[nbs] += 0.09
			}
		}
	}
	return score
}

// placeCity places another city at the region with the highest fitness score.
func placeCity(render *Terrain) {
	score := cityScore(render)
	var newcity int
	lastMax := -999999.0
	for i, val := range score.Values {
		if val > lastMax {
			newcity = i
			lastMax = val
		}
	}
	render.cities = append(render.cities, newcity)
}

func placeCities(render *Terrain) {
	params := render.params
	n := params.NumCities
	for i := 0; i < n; i++ {
		placeCity(render)
	}
}
