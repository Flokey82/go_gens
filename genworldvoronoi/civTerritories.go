package genworldvoronoi

import (
	"container/heap"
	"math"
)

// getTerritoryCultureWeightFunc returns a weight function which returns a penalty
// for expanding from a region into a region with a different culture.
// An additional penalty is applied if the destination region has a different
// culture than the origin region.
func (m *Civ) getTerritoryCultureWeightFunc() func(o, u, v int) float64 {
	return func(o, u, v int) float64 {
		var penalty float64
		// TODO: Compare culture expansionism?
		// If the destination has a higher culture expansionism than the
		// origin culture, then it's less likely to expand into that territory.
		if m.RegionToCulture[o] != m.RegionToCulture[v] {
			penalty += 0.25
		}
		if m.RegionToCulture[u] != m.RegionToCulture[v] {
			penalty += 0.75
		}
		return penalty
	}
}

// getTerritoryBiomeWeightFunc returns a weight function which returns a penalty
// for expanding from a region into a region with a different biome.
func (m *Civ) getTerritoryBiomeWeightFunc() func(o, u, v int) float64 {
	biomeFunc := m.getRegWhittakerModBiomeFunc()
	climatFunc := m.getFitnessClimate()
	return func(o, u, v int) float64 {
		var penalty float64

		// Try to stick with original biome?
		// if biomeFunc(o) != biomeFunc(v) {
		//	// Penalty is higher for inhospitable climates.
		//	penalty += 0.05 * (1 - (climatFunc(o)+climatFunc(v))/2)
		// }

		// Changes in biomes are also considered natural boundaries.
		if biomeFunc(u) != biomeFunc(v) {
			// Penalty is higher for inhospitable climates.
			penalty += 0.1 + 0.9*(1-(climatFunc(u)+climatFunc(v))/2)
		}
		return penalty
	}
}

// getTerritoryWeightFunc returns a weight function which returns a penalty
// depending on the slope of the terrain, the distance, and changes in
// flux (river crossings).
func (m *Civ) getTerritoryWeightFunc() func(o, u, v int) float64 {
	// Get maxFlux and maxElev for normalizing.
	_, maxFlux := minMax(m.Flux)
	_, maxElev := minMax(m.Elevation)

	return func(o, u, v int) float64 {
		// Don't cross from water to land and vice versa,
		// don't do anything below or at sea level.
		if (m.Elevation[u] > 0) != (m.Elevation[v] > 0) || m.Elevation[v] <= 0 {
			return -1
		}

		// Calculate horizontal distance.
		ulat := m.LatLon[u][0]
		ulon := m.LatLon[u][1]
		vlat := m.LatLon[v][0]
		vlon := m.LatLon[v][1]
		horiz := haversine(ulat, ulon, vlat, vlon) / (2 * math.Pi)

		// TODO: Maybe add a small penalty based on distance from the capital?
		// oLat := m.r_latLon[o][0]
		// oLon := m.r_latLon[o][1]
		// originDist := haversine(vlat, vlon, oLat, oLon) / (2 * math.Pi)

		// Calculate vertical distance.
		vert := (m.Elevation[v] - m.Elevation[u]) / maxElev
		if vert > 0 {
			vert /= 10
		}
		diff := 1 + 0.25*math.Pow(vert/horiz, 2)
		diff += 100 * math.Sqrt(m.Flux[u]/maxFlux)
		if m.Elevation[u] <= 0 {
			diff = 100
		}
		return horiz * diff
	}
}

// NOTE: The weight function takes three parameters:
// o: The origin/seed region
// u: The region we expand from
// v: The region we expand to
func (m *Civ) regPlaceNTerritoriesCustom(seedPoints []int, weight func(o, u, v int) float64) []int {
	var queue ascPriorityQueue
	heap.Init(&queue)
	outReg := make([]int, 0, 8)

	// 'terr' will hold a mapping of region to territory.
	// The territory ID is the region number of the capital city.
	terr := initRegionSlice(m.mesh.numRegions)
	for i := 0; i < len(seedPoints); i++ {
		terr[seedPoints[i]] = seedPoints[i]
		for _, v := range m.mesh.r_circulate_r(outReg, seedPoints[i]) {
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
		for _, v := range m.mesh.r_circulate_r(outReg, u.destination) {
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

func (m *Civ) rRelaxTerritories(terr []int, n int) {
	outReg := make([]int, 0, 8)
	for i := 0; i < n; i++ {
		// TODO: Make sure that we can put some type of constraints on
		// how much a territory can move.
		for r, t := range terr {
			if t < 0 {
				continue
			}
			var nbCountOtherTerr, nbCountSameTerr int
			otherTerr := -1
			for _, v := range m.mesh.r_circulate_r(outReg, r) {
				if v < 0 {
					continue
				}
				if terr[v] != t {
					nbCountOtherTerr++
					otherTerr = terr[v]
				} else {
					nbCountSameTerr++
				}
			}
			if nbCountOtherTerr > nbCountSameTerr && otherTerr >= 0 {
				terr[r] = otherTerr
			}
		}
	}
}

// getTerritoryNeighbors returns a list of territories neighboring the
// territory with the ID 'r' based on the provided slice of len NumRegions
// which maps the index (region id) to their respective territory ID.
func (m *Civ) getTerritoryNeighbors(r int, r_terr []int) []int {
	var res []int
	seenTerritories := make(map[int]bool)
	outReg := make([]int, 0, 8)
	for i, rg := range r_terr {
		if rg != r {
			continue
		}
		for _, nb := range m.mesh.r_circulate_r(outReg, i) {
			// Determine territory ID.
			terrID := r_terr[nb]
			if terrID < 0 || terrID == r || seenTerritories[terrID] {
				continue
			}
			seenTerritories[terrID] = true
			res = append(res, terrID)
		}
	}
	return res
}
