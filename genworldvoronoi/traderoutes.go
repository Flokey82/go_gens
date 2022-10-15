package genworldvoronoi

import (
	"log"
	"math"
	"sort"

	"github.com/chsc/astar"
)

type TradeNode struct {
	r         *Map
	getNode   func(int) *TradeNode
	index     int       // node index / region number
	used      int       // number of times this node was used for a trade route
	steepness []float64 // cached steepness of all regiones
}

func (n *TradeNode) SetUsed() {
	n.used++
}

func (n *TradeNode) NumNeighbours() int {
	return len(n.r.rNeighbors(n.index))
}

func (n *TradeNode) Neighbour(i int) astar.Node {
	// TODO: Fix this... this is highly inefficient.
	return n.getNode(n.r.rNeighbors(n.index)[i])
}

func (n *TradeNode) Cost(i int) float32 {
	// Discourage underwater paths.
	if n.r.r_elevation[n.index] <= 0 {
		return 999.00
	}
	// TODO: Fix this... this is highly inefficient.
	nIdx := n.r.rNeighbors(n.index)[i]
	if n.r.r_elevation[nIdx] <= 0 {
		return 999.00
	}

	cost := float32(1.0)

	// Altitude changes come with a cost.
	cost += float32(math.Abs(n.r.r_elevation[n.index]-n.r.r_elevation[nIdx])) * 10
	if n.used > 0 {
		cost *= 0.75 / float32(n.used)
	}

	// The steeper the terrain, the more expensive.
	cost *= 1.0 + float32(n.steepness[nIdx]*n.steepness[nIdx])

	// Heavily incentivize re-using existing roads.
	if nUsed := n.Neighbour(i).(*TradeNode).used; nUsed > 0 {
		cost *= 0.25 / float32(nUsed)
	}

	// Bonus if along coast.
	for _, nbnb := range n.r.rNeighbors(nIdx) {
		if n.r.r_elevation[nbnb] <= 0 {
			cost *= 0.65
			break
		}
	}

	// Cost of crossing rivers.
	if (n.r.r_flux[n.index] > n.r.r_rainfall[n.index]) != (n.r.r_flux[nIdx] > n.r.r_rainfall[nIdx]) {
		cost *= 1.4
	}

	// Bonus if along rivers.
	if n.r.r_flux[n.index] > n.r.r_rainfall[n.index] && n.r.r_flux[nIdx] > n.r.r_rainfall[nIdx] {
		cost *= 0.8
	}

	// Penalty for crossing into a new territory
	if n.r.r_territory[n.index] != n.r.r_territory[nIdx] {
		cost += 1.2
	}

	/*
		if n.r.rivers[n.index] < 0 && n.r.rivers[nIdx] >= 0 {
			cost *= 1.4
		}
		if n.r.rivers[n.index] >= 0 && n.r.rivers[n.index] >= 0 {
			cost *= 0.8
		}
		if n.r.terr[n.index] != n.r.terr[nIdx] {
			cost += 1.2
		}*/
	return cost
}

func estimateFunction(start, end astar.Node) float32 {
	return float32(start.(*TradeNode).r.getRDistance(start.(*TradeNode).index, end.(*TradeNode).index))
}

func getTradeRoutes(r *Map) ([][]int, [][]int) {
	nodeCache := make(map[int]*TradeNode)
	steepness := r.getRSteepness()

	// linking will store which cities are linked through a trade route crossing
	// the given region.
	linking := make([][]int, r.mesh.numRegions)
	var getNode func(i int) *TradeNode
	getNode = func(i int) *TradeNode {
		// Make sure we re-use pre-existing nodes.
		n, ok := nodeCache[i]
		if ok {
			return n
		}

		// If we have no cached node for this index,
		// create a new one.
		n = &TradeNode{
			steepness: steepness,
			r:         r,
			index:     i,
			getNode:   getNode,
		}
		nodeCache[i] = n
		return n
	}
	cities := r.cities_r
	pather := astar.New(estimateFunction)
	visited := make(map[[2]int]bool)

	// Paths contains a list of all trade routes represented through
	// a list of connected regions.
	//
	// Note that we still double up if two trade routes happen to
	// share a common section leading up to a city.
	var paths [][]int

	// TODO: Pair up by import/export of goods and taxes to the capital.
	sortCityIdx := make([]int, len(cities))
	for i := range sortCityIdx {
		sortCityIdx[i] = i
	}
	for i, start := range cities {
		// Sort by distance to start as we try to connect the closest towns first.
		sort.Slice(sortCityIdx, func(j, k int) bool {
			return r.getRDistance(start, cities[sortCityIdx[j]]) < r.getRDistance(start, cities[sortCityIdx[k]])
		})
		for _, j := range sortCityIdx {
			end := cities[j]
			// We don't want to link a city to itself and we try to avoid double
			// links (a->b and b->a) as well as we try to only connect towns within
			// the same territory.
			if i == j || visited[[2]int{start, end}] || visited[[2]int{end, start}] || r.r_territory[start] != r.r_territory[end] { //  || math.Abs(float64(i-j)) > float64(5)
				continue
			}
			// Make sure we note that we have visited this city pair.
			visited[[2]int{start, end}] = true

			// Attempt to find a path between the two cities.
			if !pather.Find(getNode(start), getNode(end)) {
				continue
			}
			var newPath []int
			for _, n := range pather.Path() {
				// Mark the node as used.
				n.(*TradeNode).SetUsed()
				nIdx := n.(*TradeNode).index

				// Check if the cities are already in our list for
				// the given region (aka "node index").
				if !isInIntList(linking[nIdx], start) {
					linking[nIdx] = append(linking[nIdx], start)
				}
				if !isInIntList(linking[nIdx], end) {
					linking[nIdx] = append(linking[nIdx], end)
				}

				// Append the region to the path.
				newPath = append(newPath, nIdx)
			}
			paths = append(paths, newPath)
		}
	}
	return paths, linking
}

func isInIntList(l []int, i int) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}
