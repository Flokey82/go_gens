package genworldvoronoi

import (
	"container/heap"
	"math"
	//"github.com/Flokey82/go_gens/vectors"
)

type queueEntry struct {
	index int // The index of the item in the heap.
	score float64
	city  int
	vx    int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue2 []*queueEntry

func (pq PriorityQueue2) Len() int { return len(pq) }

func (pq PriorityQueue2) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].score > pq[j].score
}

func (pq PriorityQueue2) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue2) Push(x interface{}) {
	n := len(*pq)
	item := x.(*queueEntry)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue2) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (m *Map) rPlaceNTerritories(n int) {
	_, maxFlux := minMax(m.r_flux)
	_, maxElev := minMax(m.r_elevation)
	if n > len(m.cities_r) {
		n = len(m.cities_r)
	}
	terr := make([]int, m.mesh.numRegions)
	var queue PriorityQueue2
	heap.Init(&queue)
	weight := func(u, v int) float64 {
		uVec := convToVec3(m.r_xyz[u*3 : u*3+3])
		ulat, ulon := latLonFromVec3(uVec.Normalize(), 1.0)

		vVec := convToVec3(m.r_xyz[v*3 : v*3+3])
		vlat, vlon := latLonFromVec3(vVec.Normalize(), 1.0)
		horiz := haversine(ulat, ulon, vlat, vlon) / (2 * math.Pi)
		//horiz := vectors.Dist3(uVec, vVec)
		vert := m.r_elevation[v]/maxElev - m.r_elevation[u]/maxElev
		if vert > 0 {
			vert /= 10
		}
		diff := 1 + 0.25*math.Pow(vert/horiz, 2)
		diff += 100 * math.Sqrt(m.r_flux[u]/maxFlux)
		if m.r_elevation[u] <= 0 {
			diff = 100
		}
		if (m.r_elevation[u] > 0) != (m.r_elevation[v] > 0) {
			return -1
		}
		return horiz * diff
	}
	for i := 0; i < n; i++ {
		terr[m.cities_r[i]] = m.cities_r[i]
		nbs := m.rNeighbors(m.cities_r[i])
		for j := 0; j < len(nbs); j++ {
			heap.Push(&queue, &queueEntry{
				score: weight(m.cities_r[i], nbs[j]),
				city:  m.cities_r[i],
				vx:    nbs[j],
			})
		}
	}
	for queue.Len() > 0 {
		u := queue.Pop().(*queueEntry)
		if terr[u.vx] != 0 {
			continue
		}
		terr[u.vx] = u.city
		nbs := m.rNeighbors(u.vx)
		for i := 0; i < len(nbs); i++ {
			v := nbs[i]
			if terr[v] != 0 {
				continue
			}
			newdist := weight(u.vx, v)
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueEntry{
				score: u.score + newdist,
				city:  u.city,
				vx:    v,
			})
		}
	}
	m.r_territory = terr
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// distance between latitudes and longitudes
	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)

	// convert to radians
	lat1 = degToRad(lat1)
	lat2 = degToRad(lat2)

	// apply formulae
	a := math.Pow(math.Sin(dLat/2), 2) + math.Pow(math.Sin(dLon/2), 2)*math.Cos(lat1)*math.Cos(lat2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c
}
