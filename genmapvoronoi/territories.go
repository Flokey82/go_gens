package genmapvoronoi

import (
	"container/heap"
	"math"
)

type queueEntry struct {
	index int // The index of the item in the heap.
	score float64
	city  int
	vx    int
}

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*queueEntry

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].score > pq[j].score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*queueEntry)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func getTerritories(render *Terrain) []int {
	h := render.h
	cities := render.cities
	n := render.params.nterrs
	if n > len(render.cities) {
		n = len(render.cities)
	}
	flux := getFlux(h)
	terr := make([]int, h.Len())
	var queue PriorityQueue
	heap.Init(&queue)
	weight := func(u, v int) float64 {
		horiz := h.Distance(u, v)
		vert := h.Values[v] - h.Values[u]
		if vert > 0 {
			vert /= 10
		}
		diff := 1 + 0.25*math.Pow(vert/horiz, 2)
		diff += 100 * math.Sqrt(flux.Values[u])
		if h.Values[u] <= 0 {
			diff = 100
		}
		if (h.Values[u] > 0) != (h.Values[v] > 0) {
			return 1000
		}
		return horiz * diff
	}
	for i := 0; i < n; i++ {
		terr[cities[i]] = cities[i]
		nbs := h.Neighbours(cities[i])
		for j := 0; j < len(nbs); j++ {
			heap.Push(&queue, &queueEntry{
				score: weight(cities[i], nbs[j]),
				city:  cities[i],
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
		nbs := h.Neighbours(u.vx)
		for i := 0; i < len(nbs); i++ {
			v := nbs[i]
			if terr[v] != 0 {
				continue
			}
			newdist := weight(u.vx, v)
			heap.Push(&queue, &queueEntry{
				score: u.score + newdist,
				city:  u.city,
				vx:    v,
			})
		}
	}
	return terr
}

func getCityTerritories(render *Terrain, terr []int) []int {
	h := render.h
	cities := render.cities
	n := len(render.cities)
	rankCity := make(map[int]int)
	// TODO: Change to size by population.
	for i, city := range cities {
		rankCity[city] = len(cities) - i
	}
	cellMult := 5
	numCells := make(map[int]int)
	//scores := cityScore(h, cities)
	flux := getFlux(h)
	cterr := make([]int, h.Len())
	var queue PriorityQueue
	heap.Init(&queue)
	weight := func(u, v int) float64 {
		horiz := h.Distance(u, v)
		if horiz > 0 {
			horiz /= 15
		}
		vert := h.Values[v] - h.Values[u]
		if vert > 0 {
			vert /= 2
		}
		diff := 1 + 0.25*math.Pow(vert/horiz, 2)
		diff += 20 * math.Sqrt(flux.Values[u])
		if h.Values[u] <= 0 {
			diff = 100
		}
		if (h.Values[u] > 0) != (h.Values[v] > 0) {
			return 1000
		}
		return horiz * diff
	}
	for i := 0; i < n; i++ {
		cterr[cities[i]] = cities[i]
		nbs := h.Neighbours(cities[i])
		for j := 0; j < len(nbs); j++ {
			heap.Push(&queue, &queueEntry{
				score: weight(cities[i], nbs[j]),
				city:  cities[i],
				vx:    nbs[j],
			})
		}
	}
	for queue.Len() > 0 {
		u := queue.Pop().(*queueEntry)
		// Check if territory is claimed.
		if cterr[u.vx] != 0 {
			continue
		}
		// Check if city territory is the same as the potential extended territory.
		if terr[u.city] != terr[u.vx] {
			continue
		}
		// Check if we have enough territory for the city taken in account its rank.
		if numCells[u.city] > rankCity[u.city]*cellMult {
			continue
		}
		numCells[u.city]++
		cterr[u.vx] = u.city
		nbs := h.Neighbours(u.vx)
		for i := 0; i < len(nbs); i++ {
			v := nbs[i]
			if cterr[v] != 0 {
				continue
			}
			// Check if city territory is the same as the potential extended territory.
			if terr[u.city] != terr[v] {
				continue
			}
			newdist := weight(u.city, v)
			heap.Push(&queue, &queueEntry{
				score: u.score + newdist,
				city:  u.city,
				vx:    v,
			})
		}
	}
	return cterr
}
