package genmapvoronoi

import (
	"math"

	"github.com/Flokey82/go_gens/vmesh"
)

func cityScore(render *Terrain) *vmesh.Heightmap {
	h := render.h
	cities := render.cities
	score := getFlux(h).MapF(math.Sqrt)
	for i := 0; i < h.Len(); i++ {
		if h.Values[i] <= 0 || h.IsNearEdge(i) {
			score.Values[i] = -999999.0
			continue
		}

		// TODO: Add bonus if near ocean or lake.
		// TODO: Consider sediment/fertility of land.
		score.Values[i] += 0.01 / (1e-9 + math.Abs(h.Vertices[i].X) - h.Extent.Width/2)
		score.Values[i] += 0.01 / (1e-9 + math.Abs(h.Vertices[i].Y) - h.Extent.Height/2)
		for j := 0; j < len(cities); j++ {
			score.Values[i] -= 0.02 / (h.Distance(cities[j], i) + 1e-9)
		}
		for _, nbs := range h.Neighbours(i) {
			if h.Values[nbs] <= 0 {
				score.Values[nbs] += 0.09
			}
		}
	}
	return score
}

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
