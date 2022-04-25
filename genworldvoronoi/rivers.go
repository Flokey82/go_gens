package genworldvoronoi

import (
	"container/heap"
	"log"
	"sort"
)

func (m *Map) getRivers(limit float64) [][]int {
	dh := m.r_downhill
	flux := m.r_flux
	_, maxFlux := minMax(flux)
	var links [][2]int
	for r := 0; r < m.mesh.numRegions; r++ {
		if flux[r]/maxFlux > limit && m.r_elevation[r] > 0 && dh[r] >= 0 && flux[dh[r]]/maxFlux > limit {
			up := r
			down := dh[r]
			links = append(links, [2]int{up, down})
		}
	}
	log.Println("start merge")
	return mergeIndexSegments(links)
}

// getRiverIndices returns a mapping from regions to river ID.
func (m *Map) getRiverIndices(limit float64) []int {
	// Set up defaults.
	rivers := make([]int, m.mesh.numRegions)
	for i := range rivers {
		rivers[i] = -1 // -1 means no river
	}
	for i, riv := range m.getRivers(limit) {
		for _, idx := range riv {
			rivers[idx] = i
		}
	}
	return rivers
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
			//log.Println("done paths", len(paths), "pathlen", len(path))
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

	// Identify sinks above sea level.
	var attempts int
	maxAttempts := 20

	// Reset drains.
	for i := range m.r_drainage {
		m.r_drainage[i] = -1
	}

	// Try to flood all sinks.
	for {
		var r_sinks []int
		for r := range r_downhill {
			lowest_r := -1
			lowest_elevation := 999.0
			for _, neighbor_r := range m.rNeighbors(r) {
				if elev := m.r_elevation[neighbor_r] + m.r_pool[neighbor_r]; elev < lowest_elevation && elev < m.r_elevation[r]+m.r_pool[r] {
					lowest_elevation = elev
					lowest_r = neighbor_r
				}
			}
			r_downhill[r] = lowest_r
			if lowest_r == -1 && m.r_elevation[r] > 0 { // && m.r_drainage[r] < 0
				r_sinks = append(r_sinks, r)
			}
		}
		m.assignFlux()

		// TODO: Regenerate flux.
		if len(r_sinks) == 0 || attempts > maxAttempts {
			log.Println("ran out of attempts", len(r_sinks))
			// TODO: Fill remaining sinks and re-generate downhill and flux.
			break
		}
		attempts++
		// Flood sink up to lowest neighbor + epsilon.
		for _, r := range r_sinks {
			if m.r_flux[r] < m.r_rainfall[r] {
				continue
			}
			m.floodV2(r, m.r_flux[r])
		}
	}
	m.r_downhill = r_downhill

	// TODO: Triangle downhill.
	log.Println(m.r_drainage)
	log.Println(m.r_pool)

	// TODO: Make note of lakes.
	//   - Sum up regions r_pool[r] > 0
	//   - Note lake sizes (for city placement)
	m.r_lake_size = m.getLakeSizes()
	// TODO: Make note of oceans.
	//   - Note ocean sizes (and small waterbodies below sea level)
	m.r_waterbodies = m.getWaterBodies()
	m.r_waterbody_size = m.getWaterBodySizes()
	// TODO: Make note of rivers.
}

func (m *Map) getWaterBodies() []int {
	done := make([]int, m.mesh.numRegions)
	for i := range done {
		if m.r_elevation[i] > 0 {
			done[i] = -2
		} else {
			done[i] = -1
		}
	}
	for r := range done {
		if done[r] != -1 {
			continue
		}
		done[r] = r
		var diveDeeper func(rd int)
		diveDeeper = func(rd int) {
			for _, nbs := range m.rNeighbors(rd) {
				if m.r_elevation[nbs] > 0 || done[nbs] != -1 {
					continue
				}
				done[nbs] = r
				diveDeeper(nbs)
			}
		}
		diveDeeper(r)
	}
	return done
}

func (m *Map) getWaterBodySizes() map[int]int {
	wbSize := make(map[int]int)
	for _, wb := range m.r_waterbodies {
		if wb >= 0 {
			wbSize[wb]++
		}
	}
	return wbSize
}

func (m *Map) getLakeSizes() map[int]int {
	lakeSize := make(map[int]int)
	for _, drain := range m.r_drainage {
		if drain != -1 {
			lakeSize[drain]++
		}
	}
	return lakeSize
}

func (m *Map) flood(r int, dVol float64) {
	const volumeFactor = 100.0 // "Water Deposition Rate"
	const epsilon = 1e-3
	const minVol = 0.01
	const drainage = 0.01
	// TODO: Return remaining volume
	plane := m.r_elevation[r] + m.r_pool[r]
	initialplane := plane

	// Floodset
	set := make([]int, 0, 1024)
	fail := 200
	tried := make([]bool, m.mesh.numRegions)
	var drain int
	var drainfound bool
	var fill func(i int)
	fill = func(i int) {
		// Out of Bounds or position has been tried.
		if i < 0 || tried[i] {
			return
		}
		tried[i] = true

		// Wall / Boundary
		currHeight := m.r_elevation[i] + m.r_pool[i]
		if plane < currHeight {
			return
		}

		// Drainage Point
		if initialplane > currHeight {
			if !drainfound || currHeight < m.r_pool[drain]+m.r_elevation[drain] {
				// No Drain yet or lower drain.
				drain = i
			}

			drainfound = true
			return
		}

		// Part of the Pool
		set = append(set, i)
		nbs := m.rNeighbors(i)
		sort.Slice(nbs, func(si, sj int) bool {
			return m.r_elevation[nbs[si]]+m.r_pool[nbs[si]] < m.r_elevation[nbs[sj]]+m.r_pool[nbs[sj]]
		})
		for _, neighbor_r := range nbs {
			fill(neighbor_r)
		}
	}

	// Iterate
	for dVol > minVol && fail != 0 {
		set = set[:0]
		for i := range tried {
			tried[i] = false
		}
		drain = 0
		drainfound = false

		// Perform Flood
		fill(r)

		// Drainage Point
		if drainfound {
			// Set the New Waterlevel (Slowly)
			plane = (1.0-drainage)*initialplane + drainage*(m.r_elevation[drain]+m.r_pool[drain])

			// Compute the New Height
			for _, s := range set {
				if plane > m.r_elevation[s] {
					m.r_pool[s] = plane - m.r_elevation[s]
					m.r_drainage[s] = drain
				} else {
					m.r_pool[s] = 0.0
					m.r_drainage[s] = -1
				}
			}
			// Remove Sediment
			//d.sediment *= 0.1
			log.Println(r, "found drain!")
			break
		}

		// Get Volume under Plane
		// So we sum up the difference between plane and (height[s]+pool[s]) which
		// gives up the total missing volume required for a full flood.
		var totalVol float64
		for _, s := range set {
			totalVol += volumeFactor * (plane - (m.r_elevation[s] + m.r_pool[s]))
		}
		// log.Println("totalVol", totalVol, "dVol", dVol, "setLen", len(set))
		// We can fill the volume of the sink.
		if totalVol <= dVol && initialplane < plane {
			// Raise water level to plane height.
			for _, s := range set {
				m.r_pool[s] = plane - m.r_elevation[s]
			}

			// Adjust flux Volume
			dVol -= totalVol
			totalVol = 0.0
		} else {
			fail-- // Plane was too high.
		}

		// Adjust Planes
		if plane > initialplane {
			initialplane = plane
		}
		// log.Println("plane before", plane)
		plane += 0.5 * (dVol - totalVol) / float64(len(set)) / volumeFactor
		log.Println(r, "plane after", plane)
	}
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
		if m.r_drainage[i] >= 0 {
			r_flux[m.r_drainage[i]] += r_flux[i]
		} else if m.r_elevation[i] >= 0 {
			r_flux[i] += m.r_rainfall[i]
		}
	}
	sort.Slice(idxs, func(a, b int) bool {
		return (m.r_elevation[idxs[b]] + m.r_pool[idxs[b]]) < (m.r_elevation[idxs[a]] + m.r_pool[idxs[a]])
	})
	for i := 0; i < m.mesh.numRegions; i++ {
		j := idxs[i]
		// Check if we are entering a pool that drains somewhere else.
		if m.r_drainage[j] >= 0 {
			r_flux[m.r_drainage[j]] += r_flux[j]
		} else if m.r_downhill[j] >= 0 {
			r_flux[m.r_downhill[j]] += r_flux[j]
		}
	}
	m.r_flux = r_flux
	log.Println(minMax(r_flux))
}

// Flooding Algorithm Overhaul:
// Currently, I can only flood at my position as long as we are rising.
// Then I return and let the particle descend. This should only happen if I can't find a closed set to fill.
// So: Rise and fill, removing the volume as we go along.
// Then: If we find a lower point, try to rise and fill from there.
//
// See: https://github.com/weigert/SimpleHydrology/blob/master/source/water.h
func (m *Map) floodV2(r int, dVol float64) bool {
	minVol := 0.001
	if dVol < minVol {
		return false
	}
	volumeFactor := 0.5

	// Either try to find a closed set under this plane, which has a certain volume,
	// or raise the plane till we find the correct closed set height.
	// And only if it can't be found, re-emit the particle.
	tried := make([]bool, m.mesh.numRegions)
	boundary := make(map[int]float64)
	var floodset []int
	var drainfound bool
	var drain int

	// Returns whether the set is closed at given height
	var findset func(i int, plane float64) bool
	findset = func(i int, plane float64) bool {
		// Out of Bounds or position has been tried.
		if i < 0 || tried[i] {
			return true
		}
		tried[i] = true

		// Wall / Boundary
		currHeight := m.r_elevation[i] + m.r_pool[i]
		if plane < currHeight {
			boundary[i] = currHeight
			return true
		}

		// Drainage Point
		if currHeight < plane {
			// No Drain yet
			if !drainfound || currHeight < m.r_pool[drain]+m.r_elevation[drain] {
				drain = i
			}
			drainfound = true
			return false
		}

		// Part of the Pool
		floodset = append(floodset, i)
		nbs := m.rNeighbors(i)
		sort.Slice(nbs, func(si, sj int) bool {
			return m.r_elevation[nbs[si]]+m.r_pool[nbs[si]] < m.r_elevation[nbs[sj]]+m.r_pool[nbs[sj]]
		})
		for _, neighbor_r := range nbs {
			if !findset(neighbor_r, plane) {
				return false
			}
		}
		return true
	}

	plane := m.r_pool[r] + m.r_elevation[r]
	minboundFirst := r
	minboundSecond := plane
	for dVol > minVol && findset(r, plane) {
		// Find the Lowest Element on the Boundary
		minboundFirst = -1
		for bfirst, bsecond := range boundary {
			if bsecond < minboundSecond || minboundFirst == -1 {
				minboundFirst = bfirst
				minboundSecond = bsecond
			}
		}
		// Compute the Height of our Volume over the Set
		vheight := dVol * volumeFactor / float64(len(floodset))

		// Not High Enough: Fill 'er up
		if plane+vheight < minboundSecond {
			plane += vheight
		} else {
			dVol -= (minboundSecond - plane) / volumeFactor * float64(len(floodset))
			plane = minboundSecond
		}

		for _, s := range floodset {
			m.r_pool[s] = plane - m.r_elevation[s]
			m.r_drainage[s] = drain // WROOOOONG?????
		}
		delete(boundary, minboundFirst)
		tried[minboundFirst] = false
		r = minboundFirst
	}

	if drainfound {
		// Search for Exposed Neighbor with Non-Zero Waterlevel
		var lowbound func(i int)
		lowbound = func(i int) {
			// Out-Of-Bounds
			if i < 0 || m.r_pool[i] == 0 {
				return
			}
			// Below Drain Height
			if m.r_elevation[i]+m.r_pool[i] < m.r_elevation[drain]+m.r_pool[drain] {
				return
			}
			// Higher than Plane (we want lower)
			if m.r_elevation[i]+m.r_pool[i] >= plane {
				return
			}
			plane = m.r_elevation[i] + m.r_pool[i]
		}
		// Fill Neighbors
		for _, nbs := range m.rNeighbors(drain) {
			lowbound(nbs)
			//for _, nbs2 := range m.rNeighbors(nbs) { // ??????
			//	lowbound(nbs2)
			//}
		}

		// Water-Level to Plane-Height
		for _, s := range floodset {
			//  volume += ((plane > h[ind])?(h[ind] + p[ind] - plane):p[ind])/volumeFactor;
			if plane > m.r_elevation[s] {
				m.r_pool[s] = plane - m.r_elevation[s]
				m.r_drainage[s] = drain
			} else {
				m.r_pool[s] = 0.0
				m.r_drainage[s] = -1
			}
		}

		for bfirst := range boundary {
			//  volume += ((plane > h[ind])?(h[ind] + p[ind] - plane):p[ind])/volumeFactor;
			if plane > m.r_elevation[bfirst] {
				m.r_pool[bfirst] = plane - m.r_elevation[bfirst]
				m.r_drainage[bfirst] = drain
			} else {
				m.r_pool[bfirst] = 0.0
				m.r_drainage[bfirst] = -1
			}
		}
		// sediment *= oldvolume/volume;
		// sediment /= float64(len(floodset)) //Distribute Sediment in Pool
		r = drain
		return true
	}
	return false
}

// Rivers - from mapgen4

// assignTriangleValues averages out the values of the mesh points / regions and assigns them
// to the triangles of the mesh (or the triangle centroid).
func (m *Map) assignTriangleValues() {
	r_elevation := m.r_elevation
	r_moisture := m.r_moisture
	r_pool := m.r_pool
	t_elevation := m.t_elevation
	t_moisture := m.t_moisture
	t_pool := m.t_pool
	numTriangles := m.mesh.numTriangles
	for t := 0; t < numTriangles; t++ {
		s0 := 3 * t
		r1 := m.mesh.s_begin_r(s0)
		r2 := m.mesh.s_begin_r(s0 + 1)
		r3 := m.mesh.s_begin_r(s0 + 2)
		t_pool[t] = (1.0 / 3.0) * (r_pool[r1] + r_pool[r2] + r_pool[r3])
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
	m.t_pool = t_pool
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
