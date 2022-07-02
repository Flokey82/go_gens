package genworldvoronoi

import (
	"container/heap"
	"container/list"
	"log"
	"math"
	//"github.com/Flokey82/go_gens/vectors"
)

// identifyLandmasses returns a mapping from region to landmass ID.
// A landmass is a connected number of regions above sealevel.
func (m *Map) identifyLandmasses() []int {
	// NOTE: this is still in need of refinement.
	landMasses := make([]int, m.mesh.numRegions)
	for r := range landMasses {
		landMasses[r] = -1
	}
	var landID int
	var landSizes []int
	for r, h := range m.r_elevation {
		// Skip if the current region has already been allocated
		// or is below sealevel.
		if landMasses[r] != -1 || h < 0 {
			continue
		}
		var currentLandSize int
		queue := list.New()
		enqueue := func(r int) {
			// Skip if the current region has already been allocated
			// or is below sealevel.
			if landMasses[r] != -1 || m.r_elevation[r] < 0 {
				return
			}
			landMasses[r] = landID // Assign current landID to the region.
			currentLandSize++      // Increase size of known landmass.
			for _, nb := range m.rNeighbors(r) {
				// Skip if the neighbor region has already been allocated
				// or is below sealevel.
				if landMasses[nb] != -1 || m.r_elevation[nb] < 0 {
					continue
				}
				queue.PushBack(nb)
			}
		}
		// Start queue with current region.
		queue.PushBack(r)

		// Process each queue entry until we run out of
		// regions that belong to this landmass.
		for queue.Len() > 0 {
			e := queue.Front()
			enqueue(e.Value.(int))
			queue.Remove(e)
		}

		// Once done, append the current size to the list of landmass-
		// sizes and increment the current landID.
		landSizes = append(landSizes, currentLandSize)
		landID++
	}
	log.Println("number of landmasses", landID)
	log.Println(landSizes)
	return landMasses
}

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
	// return pq[i].score > pq[j].score // 3, 2, 1
	return pq[i].score < pq[j].score // 1, 2, 3
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
		ulat := m.r_latLon[u][0]
		ulon := m.r_latLon[u][1]

		vlat := m.r_latLon[v][0]
		vlon := m.r_latLon[v][1]

		horiz := haversine(ulat, ulon, vlat, vlon) / (2 * math.Pi)
		vert := (m.r_elevation[v] - m.r_elevation[u]) / maxElev
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
		for _, v := range m.rNeighbors(m.cities_r[i]) {
			if m.r_elevation[v] <= 0 {
				continue
			}
			heap.Push(&queue, &queueEntry{
				score: weight(m.cities_r[i], v),
				city:  m.cities_r[i],
				vx:    v,
			})
		}
	}
	for queue.Len() > 0 {
		u := heap.Pop(&queue).(*queueEntry)
		if terr[u.vx] != 0 {
			continue
		}
		terr[u.vx] = u.city
		for _, v := range m.rNeighbors(u.vx) {
			if terr[v] != 0 || m.r_elevation[v] < 0 {
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
