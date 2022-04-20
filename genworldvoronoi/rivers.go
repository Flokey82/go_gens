package genworldvoronoi

import (
	"container/heap"
	"log"
	"sort"
)

func (m *Map) getRivers(limit float64) [][]int {
	dh := m.r_downhill
	flux := m.r_flux

	var links [][2]int
	var above int
	for _, r := range m.r_elevation {
		if r > 0 {
			above++
		}
	}
	limit *= float64(above) / float64(m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		//if isnearedge(h.Mesh, i){
		//    continue
		//}
		if flux[r] > limit && m.r_elevation[r] > 0 && dh[r] >= 0 {
			up := r
			down := dh[r]
			links = append(links, [2]int{up, down})
		}
	}
	log.Println("start merge")

	mergedSegs := mergeIndexSegments(links)
	/*
		// Set up defaults.
		rivers := make([]int, m.mesh.numRegions)
		for i := range rivers {
			rivers[i] = -1 // -1 means no river
		}
		for i := range mergedSegs {
			for _, idx := range mergedSegs[i] {
				rivers[idx] = i
			}
		}
		return rivers
	*/
	return mergedSegs
}

func mergeIndexSegments(segs [][2]int) [][]int {
	log.Println("start adj")
	adj := make(map[int][]int)
	for i := 0; i < len(segs); i++ {
		seg := segs[i]
		a0 := adj[seg[0]]
		a1 := adj[seg[1]]
		a0 = append(a0, seg[1])
		a1 = append(a1, seg[0])
		adj[seg[0]] = a0
		adj[seg[1]] = a1
	}
	done := make(map[int]bool)
	var paths [][]int
	var path []int
	log.Println("start paths")
	for {
		if path == nil {
			for i := 0; i < len(segs); i++ {
				if done[i] {
					continue
				}
				done[i] = true
				path = []int{segs[i][0], segs[i][1]}
				break
			}
			if path == nil {
				break
			}
		}
		var changed bool
		for i := 0; i < len(segs); i++ {
			if done[i] {
				continue
			}
			if len(adj[path[0]]) == 2 && segs[i][0] == path[0] {
				path = unshiftIndexPath(path, segs[i][1])
			} else if len(adj[path[0]]) == 2 && segs[i][1] == path[0] {
				path = unshiftIndexPath(path, segs[i][0])
			} else if len(adj[path[len(path)-1]]) == 2 && segs[i][0] == path[len(path)-1] {
				path = append(path, segs[i][1])
			} else if len(adj[path[len(path)-1]]) == 2 && segs[i][1] == path[len(path)-1] {
				path = append(path, segs[i][0])
			} else {
				continue
			}
			done[i] = true
			changed = true
			break
		}
		if !changed {
			log.Println("done paths", len(paths), "pathlen", len(path))
			paths = append(paths, path)
			path = nil
		}
	}
	return paths
}

func unshiftIndexPath(path []int, p int) []int {
	return append([]int{p}, path...)
}

// assignDownhill will populate r_downhill with a mapping of region to lowest neighbor region.
// NOTE: This is based on mewo2's terrain generation code
// See: https://github.com/mewo2/terrain
func (m *Map) assignDownhill() {
	// Here we will map each region to the lowest neighbor.
	r_downhill := make([]int, m.mesh.numRegions)
	for r := range r_downhill {
		lowest_r := -1
		lowest_elevation := 999.0
		for _, neighbor_r := range m.rNeighbors(r) {
			if elev := m.r_elevation[neighbor_r]; elev < lowest_elevation && elev < m.r_elevation[r] {
				lowest_elevation = elev
				lowest_r = neighbor_r
			}
		}
		r_downhill[r] = lowest_r
	}

	/*
		// TODO: Flood sinks.
		// Sort regions by elevation, lowest -> highest.
		idxs := make([]int, m.mesh.numRegions)
		for i := 0; i < m.mesh.numRegions; i++ {
			idxs[i] = i
		}
		sort.Slice(idxs, func(a, b int) bool {
			return m.r_elevation[idxs[b]]-m.r_elevation[idxs[a]] > 0
		})
		// Iterate through sorted regions and identify pits.
		for i := 0; i < m.mesh.numRegions; i++ {
			j := idxs[i]
			if r_downhill[j] == -1 && m.r_elevation[j] > 0 {
				log.Println("pit", j, m.r_elevation[j])
				visited := make(map[int]bool)

				newDH := -1
				newH := m.r_elevation[j]

				// Visit all neighbors until we find a downhill value that is not within the set of visited regions.
				var visit func(r int) bool
				visit = func(r int) bool {
					if visited[r] {
						return true
					}
					visited[r] = true
					if dh := r_downhill[r]; dh >= 0 && !visited[dh] && m.r_elevation[dh] < newH {
						newDH = r_downhill[r]
						newH = m.r_elevation[dh]
						return false
					}
					for _, nb_r := range m.rNeighbors(r) {
						if !visit(nb_r) {
							return false
						}
					}
					return true
				}
				visit(j)
				r_downhill[j] = newDH
				if r_downhill[j] == -1 {
					log.Println("still pit", j, m.r_elevation[j])
				}

			}
		}*/
	m.r_downhill = r_downhill

	// TODO: Triangle downhill.
}

// assignFlux will populate r_flux by summing up the rainfall for each region from highest to
// lowest using r_downhill to reconstruct the downhill path that water would follow.
// NOTE: This is based on mewo2's terrain generation code
// See: https://github.com/mewo2/terrain
func (m *Map) assignFlux() {
	// Initialize flux values with r_rainfall.
	r_flux := make([]float64, m.mesh.numRegions)
	idxs := make([]int, m.mesh.numRegions)
	for i := 0; i < m.mesh.numRegions; i++ {
		idxs[i] = i
		r_flux[i] = m.r_rainfall[i]
	}
	sort.Slice(idxs, func(a, b int) bool {
		return m.r_elevation[idxs[b]]-m.r_elevation[idxs[a]] < 0
	})
	for i := 0; i < m.mesh.numRegions; i++ {
		j := idxs[i]
		if m.r_downhill[j] >= 0 {
			r_flux[m.r_downhill[j]] += r_flux[j]
		}
	}
	m.r_flux = r_flux
}

// Rivers - from mapgen4

// assignTriangleValues averages out the values of the mesh points / regions and assigns them
// to the triangles of the mesh (or the triangle centroid).
func (m *Map) assignTriangleValues() {
	r_elevation := m.r_elevation
	r_moisture := m.r_moisture
	t_elevation := m.t_elevation
	t_moisture := m.t_moisture
	numTriangles := m.mesh.numTriangles
	for t := 0; t < numTriangles; t++ {
		s0 := 3 * t
		r1 := m.mesh.s_begin_r(s0)
		r2 := m.mesh.s_begin_r(s0 + 1)
		r3 := m.mesh.s_begin_r(s0 + 2)
		t_elevation[t] = (1.0 / 3.0) * (r_elevation[r1] + r_elevation[r2] + r_elevation[r3])
		t_moisture[t] = (1.0 / 3.0) * (r_moisture[r1] + r_moisture[r2] + r_moisture[r3])
	}

	// This averages out rainfall to calculate moisture for triangles.
	// Note that this overrides the t_moisture calculated by averaging out r_moisture above.
	for t := 0; t < numTriangles; t++ {
		var moisture float64
		for i := 0; i < 3; i++ {
			s := 3*t + i
			r := m.mesh.s_begin_r(s)
			moisture += m.r_rainfall[r] / 3
		}
		t_moisture[t] = moisture
	}
	m.t_elevation = t_elevation
	m.t_moisture = t_moisture
}

// assignDownflow starts with triangles that are considered "ocean" and works its way
// uphill to build a graph of child/parents that will allow us later to determine water
// flux and whatnot.
func (m *Map) assignDownflow() {
	// Use a priority queue, starting with the ocean triangles and
	// moving upwards using elevation as the priority, to visit all
	// the land triangles.
	_queue := make(PriorityQueue, 0)
	numTriangles := m.mesh.numTriangles
	queue_in := 0
	for i := range m.t_downflow_s {
		m.t_downflow_s[i] = -999
	}
	heap.Init(&_queue)

	// Part 1: ocean triangles get downslope assigned to the lowest neighbor.
	for t := 0; t < numTriangles; t++ {
		if m.t_elevation[t] < 0 {
			best_s := -1
			best_e := m.t_elevation[t]
			for j := 0; j < 3; j++ {
				s := 3*t + j
				e := m.t_elevation[m.mesh.s_outer_t(s)]
				if e < best_e {
					best_e = e
					best_s = s
				}
			}
			m.order_t[queue_in] = t
			queue_in++
			m.t_downflow_s[t] = best_s
			heap.Push(&_queue, &Item{ID: t, Value: m.t_elevation[t], Index: t})
		}
	}

	// Part 2: land triangles get visited in elevation priority.
	for queue_out := 0; queue_out < numTriangles; queue_out++ {
		current_t := heap.Pop(&_queue).(*Item).ID
		for j := 0; j < 3; j++ {
			s := 3*current_t + j
			neighbor_t := m.mesh.s_outer_t(s) // uphill from current_t
			if m.t_downflow_s[neighbor_t] == -999 && m.t_elevation[neighbor_t] >= 0.0 {
				m.t_downflow_s[neighbor_t] = m.mesh.s_opposite_s(s)
				m.order_t[queue_in] = neighbor_t
				queue_in++
				heap.Push(&_queue, &Item{ID: neighbor_t, Value: m.t_elevation[neighbor_t]})
			}
		}
	}
}

// assignFlow calculates the water flux by traversing the graph generated with
// assignDownflow in reverse order (so, downhill?) and summing up the moisture.
func (m *Map) assignFlow() {
	s_flow := m.s_flow

	// Clear all existing water flux values.
	for i := range s_flow {
		s_flow[i] = 0
	}

	t_flow := m.t_flow
	t_elevation := m.t_elevation
	t_moisture := m.t_moisture

	// Set the flux value for each triangle above sealevel to
	// half of its moisture squared as its initial state.
	numTriangles := m.mesh.numTriangles
	for t := 0; t < numTriangles; t++ {
		if t_elevation[t] >= 0.0 {
			t_flow[t] = 0.5 * t_moisture[t] * t_moisture[t]
		} else {
			t_flow[t] = 0
		}
	}

	// Now traverse the flux graph in reverse order and sum up
	// the moisture of all tributaries while descending.
	order_t := m.order_t
	t_downflow_s := m.t_downflow_s
	_halfedges := m.mesh.Halfedges
	for i := len(order_t) - 1; i >= 0; i-- {
		// TODO: Describe what's going on here.
		tributary_t := order_t[i]
		flow_s := t_downflow_s[tributary_t]
		if flow_s >= 0 {
			trunk_t := (_halfedges[flow_s] / 3)
			t_flow[trunk_t] += t_flow[tributary_t]
			s_flow[flow_s] += t_flow[tributary_t] // TODO: isn't s_flow[flow_s] === t_flow[?]
			if t_elevation[trunk_t] > t_elevation[tributary_t] {
				t_elevation[trunk_t] = t_elevation[tributary_t]
			}
		}
	}
	m.t_flow = t_flow
	m.s_flow = s_flow
	m.t_elevation = t_elevation
}

type Item struct {
	ID    int
	Value float64
	Index int // The index of the item in the heap.
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest based on expiration number as the priority
	// The lower the expiry, the higher the priority
	return pq[i].Value < pq[j].Value
}

// We just implement the pre-defined function in interface of heap.

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
