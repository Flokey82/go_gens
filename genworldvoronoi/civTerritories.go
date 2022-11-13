package genworldvoronoi

import (
	"container/heap"
	"math"
)

func (m *Map) rPlaceNTerritories(n int) {
	m.resetRand()
	// Territories are based on cities acting as their capital.
	// Since the algorithm places the cities with the highes scores
	// first, we use the top 'n' cities as the capitals for the
	// territories.
	var seedCities []int
	for i, c := range m.cities_r {
		if i >= n {
			break
		}
		seedCities = append(seedCities, c.ID)
	}
	weight := m.getTerritoryWeightFunc()
	biomeWeight := m.getTerritoryBiomeWeightFunc()
	cultureWeight := m.getTerritoryCultureWeightFunc()

	m.r_territory = m.rPlaceNTerritoriesCustom(seedCities, func(o, u, v int) float64 {
		if (m.r_elevation[u] > 0) != (m.r_elevation[v] > 0) || m.r_elevation[v] <= 0 {
			return -1
		}
		return weight(o, u, v) + biomeWeight(o, u, v) + cultureWeight(o, u, v)
	})
	m.rRelaxTerritories(m.r_territory, 15)
}

func (m *Map) rPlaceNCityStates(n int) {
	m.resetRand()
	// Territories are based on cities acting as their capital.
	// Since the algorithm places the cities with the highes scores
	// first, we use the top 'n' cities as the capitals for the
	// territories.
	var seedCities []int
	for i, c := range m.cities_r {
		if i >= n {
			break
		}
		seedCities = append(seedCities, c.ID)
	}
	weight := m.getTerritoryWeightFunc()
	biomeWeight := m.getTerritoryBiomeWeightFunc()
	cultureWeight := m.getTerritoryCultureWeightFunc()

	m.r_city = m.rPlaceNTerritoriesCustom(seedCities, func(o, u, v int) float64 {
		if m.r_territory[u] != m.r_territory[v] {
			return -1
		}
		// TODO: Make sure we take in account expansionism, wealth, score, and culture.
		return weight(o, u, v) + biomeWeight(o, u, v) + cultureWeight(o, u, v)
	})

	// Before relaxing the territories, we'd need to ensure that we only
	// relax without changing the borders of the empire...
	// So we'd only re-assign IDs that belong to the same territory.
	// m.rRelaxTerritories(cityStates, 5)
}

func (m *Map) getTerritoryCultureWeightFunc() func(o, u, v int) float64 {
	return func(o, u, v int) float64 {
		var penalty float64
		// TODO: Compare culture expansionism.
		// If the destination has a higher culture expansionism than the
		// origin culture, then it's less likely to expand into that territory.
		if m.r_cultures[o] != m.r_cultures[v] {
			penalty += 0.25
		}
		if m.r_cultures[u] != m.r_cultures[v] {
			penalty += 0.75
		}
		return penalty
	}
}

func (m *Map) getTerritoryBiomeWeightFunc() func(o, u, v int) float64 {
	biomeFunc := m.getRWhittakerModBiomeFunc()
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

func (m *Map) getTerritoryWeightFunc() func(o, u, v int) float64 {
	// Get maxFlux and maxElev for normalizing.
	_, maxFlux := minMax(m.r_flux)
	_, maxElev := minMax(m.r_elevation)

	return func(o, u, v int) float64 {
		// Don't cross from water to land and vice versa,
		// don't do anything below or at sea level.
		if (m.r_elevation[u] > 0) != (m.r_elevation[v] > 0) || m.r_elevation[v] <= 0 {
			return -1
		}

		// Calculate horizontal distance.
		ulat := m.r_latLon[u][0]
		ulon := m.r_latLon[u][1]
		vlat := m.r_latLon[v][0]
		vlon := m.r_latLon[v][1]
		horiz := haversine(ulat, ulon, vlat, vlon) / (2 * math.Pi)

		// TODO: Maybe add a small penalty based on distance from the capital?
		// oLat := m.r_latLon[o][0]
		// oLon := m.r_latLon[o][1]
		// originDist := haversine(vlat, vlon, oLat, oLon) / (2 * math.Pi)

		// Calculate vertical distance.
		vert := (m.r_elevation[v] - m.r_elevation[u]) / maxElev
		if vert > 0 {
			vert /= 10
		}
		diff := 1 + 0.25*math.Pow(vert/horiz, 2)
		diff += 100 * math.Sqrt(m.r_flux[u]/maxFlux)
		if m.r_elevation[u] <= 0 {
			diff = 100
		}
		return horiz * diff
	}
}

// NOTE: The weight function takes three parameters:
// o: The origin/seed region
// u: The region we expand from
// v: The region we expand to
func (m *Map) rPlaceNTerritoriesCustom(seedPoints []int, weight func(o, u, v int) float64) []int {
	var queue territoryQueue
	heap.Init(&queue)

	// 'terr' will hold a mapping of region to territory.
	// The territory ID is the region number of the capital city.
	terr := initRegionSlice(m.mesh.numRegions)
	for i := 0; i < len(seedPoints); i++ {
		terr[seedPoints[i]] = seedPoints[i]
		for _, v := range m.rNeighbors(seedPoints[i]) {
			newdist := weight(seedPoints[i], seedPoints[i], v)
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueRegionEntry{
				score: newdist,
				city:  seedPoints[i],
				vx:    v,
			})
		}
	}

	// Extend territories until the queue is empty.
	for queue.Len() > 0 {
		u := heap.Pop(&queue).(*queueRegionEntry)
		if terr[u.vx] >= 0 {
			continue
		}
		terr[u.vx] = u.city
		for _, v := range m.rNeighbors(u.vx) {
			if terr[v] >= 0 {
				continue
			}
			newdist := weight(u.city, u.vx, v)
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueRegionEntry{
				score: u.score + newdist,
				city:  u.city,
				vx:    v,
			})
		}
	}
	return terr
}

type queueRegionEntry struct {
	index int // The index of the item in the heap.
	score float64
	city  int
	vx    int
}

// territoryQueue implements heap.Interface and holds Items.
type territoryQueue []*queueRegionEntry

func (pq territoryQueue) Len() int { return len(pq) }

func (pq territoryQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	// return pq[i].score > pq[j].score // 3, 2, 1
	return pq[i].score < pq[j].score // 1, 2, 3
}

func (pq territoryQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *territoryQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*queueRegionEntry)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *territoryQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (m *Map) rRelaxTerritories(terr []int, n int) {
	for i := 0; i < n; i++ {
		// TODO: Make sure that we can put some type of constraints on
		// how much a territory can move.
		for r, t := range terr {
			if t < 0 {
				continue
			}
			var nbCountOtherTerr, nbCountSameTerr int
			otherTerr := -1
			for _, v := range m.rNeighbors(r) {
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

// getRTerritoryNeighbors returns a list of territories neighboring the
// territory with the ID 'r' based on the provided slice of len NumRegions
// which maps the index (region id) to their respective territory ID.
func (m *Map) getRTerritoryNeighbors(r int, r_terr []int) []int {
	var res []int
	seenTerritories := make(map[int]bool)
	for i, rg := range r_terr {
		if rg != r {
			continue
		}
		for _, nb := range m.rNeighbors(i) {
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
