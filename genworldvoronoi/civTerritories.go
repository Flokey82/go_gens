package genworldvoronoi

import (
	"container/heap"
	"math"
)

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

func (m *Map) rPlaceNTerritories(n int) {
	// Territories are based on cities acting as their capital.
	// Since the algorithm places the cities with the highes scores
	// first, we use the top 'n' cities as the capitals for the
	// territories.
	var seedCities []int
	for i, c := range m.cities_r {
		if i >= n {
			break
		}
		seedCities = append(seedCities, c.R)
	}
	weight := m.getTerritoryWeightFunc()
	m.r_territory = m.rPlaceNTerritoriesCustom(seedCities, weight)
}

func (m *Map) rPlaceNCityStates(n int) []int {
	// Territories are based on cities acting as their capital.
	// Since the algorithm places the cities with the highes scores
	// first, we use the top 'n' cities as the capitals for the
	// territories.
	var seedCities []int
	for i, c := range m.cities_r {
		if i >= n {
			break
		}
		seedCities = append(seedCities, c.R)
	}
	weight := m.getTerritoryWeightFunc()

	return m.rPlaceNTerritoriesCustom(seedCities, func(u, v int) float64 {
		if m.r_territory[u] != m.r_territory[v] {
			return -1
		}
		return weight(u, v)
	})
}

func (m *Map) getTerritoryWeightFunc() func(u, v int) float64 {
	// Get maxFlux and maxElev for normalizing.
	_, maxFlux := minMax(m.r_flux)
	_, maxElev := minMax(m.r_elevation)

	return func(u, v int) float64 {
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

func (m *Map) rPlaceNTerritoriesCustom(seedPoints []int, weight func(u, v int) float64) []int {
	var queue territoryQueue
	heap.Init(&queue)

	// 'terr' will hold a mapping of region to territory.
	// The territory ID is the region number of the capital city.
	terr := make([]int, m.mesh.numRegions)
	for r := range terr {
		terr[r] = -1
	}
	for i := 0; i < len(seedPoints); i++ {
		terr[seedPoints[i]] = seedPoints[i]
		for _, v := range m.rNeighbors(seedPoints[i]) {
			newdist := weight(seedPoints[i], v)
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
			newdist := weight(u.vx, v)
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
