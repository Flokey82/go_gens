package genworldvoronoi

import (
	"container/heap"
	"log"
	"sort"
	"time"
)

func (m *BaseObject) getRivers(limit float64) [][]int {
	links := m.getRiverSegments(limit)
	log.Println("start merge")
	start := time.Now()
	defer func() {
		log.Println("Done river segments in ", time.Since(start).String())
	}()
	return mergeIndexSegments(links)
}

func (m *BaseObject) getRiverSegments(limit float64) [][2]int {
	dh := m.r_downhill
	flux := m.r_flux
	var links [][2]int
	for r := 0; r < m.mesh.numRegions; r++ {
		if flux[r] > m.r_rainfall[r] && m.r_elevation[r] >= 0 && dh[r] >= 0 && flux[dh[r]] > m.r_rainfall[dh[r]] {
			// up := r, down := dh[r]
			links = append(links, [2]int{r, dh[r]})
		}
	}
	return links
}

func (m *BaseObject) getWaterBodies() []int {
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

func (m *BaseObject) getWaterBodySizes() map[int]int {
	wbSize := make(map[int]int)
	for _, wb := range m.r_waterbodies {
		if wb >= 0 {
			wbSize[wb]++
		}
	}
	return wbSize
}

func (m *BaseObject) getLakeSizes() map[int]int {
	lakeSize := make(map[int]int)
	for _, drain := range m.r_drainage {
		if drain != -1 {
			lakeSize[drain]++
		}
	}
	return lakeSize
}

func (m *Map) getRivers2(limit float64) [][]int {
	dh := m.r_downhill
	flux := m.r_flux
	_, maxFlux := minMax(flux)
	var links [][2]int
	for r := 0; r < m.mesh.numRegions; r++ {
		if flux[r]/maxFlux > limit && m.r_elevation[r] > 0 && dh[r] >= 0 && flux[dh[r]]/maxFlux > limit {
			// up := r, down := dh[r]
			links = append(links, [2]int{r, dh[r]})
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
		adj[seg[0]] = append(adj[seg[0]], seg[1])
		adj[seg[1]] = append(adj[seg[1]], seg[0])
	}
	var maxSegIdx int
	for s := range adj {
		if s > maxSegIdx {
			maxSegIdx = s
		}
	}
	done := make([]bool, maxSegIdx)
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

			/*
				// Original code.
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
			*/
			if len(adj[path[0]]) == 2 {
				if segs[i][0] == path[0] {
					path = unshiftIndexPath(path, segs[i][1])
					done[i] = true
					changed = true
					break
				}
				if segs[i][1] == path[0] {
					path = unshiftIndexPath(path, segs[i][0])
					done[i] = true
					changed = true
					break
				}
			}
			if len(adj[path[len(path)-1]]) == 2 {
				if segs[i][0] == path[len(path)-1] {
					path = append(path, segs[i][1])
					done[i] = true
					changed = true
					break
				}
				if segs[i][1] == path[len(path)-1] {
					path = append(path, segs[i][0])
					done[i] = true
					changed = true
					break
				}
			}
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
	res := make([]int, len(path)+1)
	res[0] = p
	copy(res[1:], path)
	return res
}

const (
	FloodVariant1 = 0
	FloodVariant2 = 1
)

func (m *Map) fillSinks() []float64 {
	var sea_r []int
	sort_r := make([]int, m.mesh.numRegions)
	is_sea := make([]bool, m.mesh.numRegions)
	for r, h := range m.r_elevation {
		if h < 0 {
			sea_r = append(sea_r, r)
			is_sea[r] = true
		}
		sort_r[r] = r
	}
	epsilon := 1e-13
	log.Println("fill sink epsilon is too big!")
	r_distance_c := m.assignDistanceField(sea_r, make(map[int]bool))
	sort.Slice(sort_r, func(i, j int) bool {
		if r_distance_c[sort_r[i]] == r_distance_c[sort_r[j]] {
			return m.r_elevation[sort_r[i]] < m.r_elevation[sort_r[j]]
		}
		return r_distance_c[sort_r[i]] < r_distance_c[sort_r[j]]
	})

	for {
		var changed bool
		r_sinks := m.BaseObject.getSinks(true, false)
		log.Println("sinks ", len(r_sinks))
		sort.Slice(r_sinks, func(i, j int) bool {
			return r_distance_c[r_sinks[i]] < r_distance_c[r_sinks[j]]
		})
		for _, r := range r_sinks {
			if m.r_flux[r] == 0 || is_sea[r] {
				continue
			}
			closest := r_distance_c[r]
			rn := -1
			for _, nb := range m.rNeighbors(r) {
				if closest > r_distance_c[nb] {
					closest = r_distance_c[nb]
					rn = nb
				}
			}
			if rn >= 0 {
				if m.r_elevation[rn] > m.r_elevation[r] {
					delta := m.r_elevation[rn] - m.r_elevation[r]
					m.r_elevation[r] += epsilon
					m.r_elevation[rn] -= delta
					changed = true
				} else if m.r_elevation[rn] == m.r_elevation[r] {
					m.r_elevation[rn] -= epsilon
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	for false {
		var changed bool
		for _, r := range sort_r {
			if is_sea[r] {
				continue
			}

			for _, nb := range m.rNeighbors(r) {
				if r_distance_c[r] < r_distance_c[nb] {
					if m.r_elevation[r] > m.r_elevation[nb] {
						delta := m.r_elevation[r] - m.r_elevation[nb]
						m.r_elevation[nb] += epsilon // delta / 2
						m.r_elevation[r] -= delta    // (delta / 2) + epsilon
						changed = true
					} else if m.r_elevation[nb] == m.r_elevation[r] {
						m.r_elevation[r] -= epsilon
						changed = true
					}
				} else if r_distance_c[r] > r_distance_c[nb] {
					if m.r_elevation[nb] > m.r_elevation[r] {
						delta := m.r_elevation[nb] - m.r_elevation[r]
						m.r_elevation[r] += epsilon // delta / 2
						m.r_elevation[nb] -= delta  // (delta / 2) + epsilon
						changed = true
					} else if m.r_elevation[nb] == m.r_elevation[r] {
						m.r_elevation[nb] -= epsilon
						changed = true
					}
				}
			}
		}
		if !changed {
			break
		}
	}
	for false {
		var changed bool
		r_sinks := m.BaseObject.getSinks(true, false)
		log.Println("sinks ", len(r_sinks))
		sort.Slice(r_sinks, func(i, j int) bool {
			return r_distance_c[r_sinks[i]] < r_distance_c[r_sinks[j]]
		})
		for _, r := range r_sinks {
			if m.r_flux[r] == 0 {
				continue
			}
			for _, nb := range m.rNeighbors(r) {
				if r_distance_c[r] > r_distance_c[nb] {
					if m.r_elevation[nb] > m.r_elevation[r] {
						delta := m.r_elevation[nb] - m.r_elevation[r]
						m.r_elevation[r] += (delta / 2) //+ epsilon
						m.r_elevation[nb] -= (delta / 2) + epsilon
						changed = true
					} else if m.r_elevation[nb] == m.r_elevation[r] {
						m.r_elevation[nb] -= epsilon
						changed = true
					}
				}
			}
			/*
				closest := r_distance_c[r]
				rn := -1
				for _, nb := range m.rNeighbors(r) {
					if closest > r_distance_c[nb] {
						closest = r_distance_c[nb]
						rn = nb
					}
				}
				if rn >= 0 {
					if m.r_elevation[rn] > m.r_elevation[r] {
						delta := m.r_elevation[rn] - m.r_elevation[r]
						m.r_elevation[r] += (delta / 2) //+ epsilon
						m.r_elevation[rn] -= (delta / 2) + epsilon
						changed = true
					} else if m.r_elevation[rn] == m.r_elevation[r] {
						m.r_elevation[rn] -= epsilon
						changed = true
					}
				}*/
		}
		if !changed {
			return m.r_elevation
		}
	}
	return m.r_elevation
}

// assignHydrology will calculate river systems and fill sinks instead of trying to generate
// water pools.
func (m *Map) assignHydrology() {
	maxAttempts := 5
	erosionAmount := 0.00001

	// Try to flood all sinks.
	var attempts int
	m.BaseObject.assignDownhill(true)
	m.assignFlux(false)
	for {
		// Abort if we have no more sinks or ran out of attempts.
		if attempts > maxAttempts {
			m.r_elevation = m.fillSinks()
			// Regenerate downhill.
			m.BaseObject.assignDownhill(true)

			// Regenerate flux.
			m.assignFlux(true)

			// TODO: Diffuse flux and pool.
			m.assignRainfallBasic()
			// TODO: Fill remaining sinks and re-generate downhill and flux.
			break
		}
		attempts++
		// Reset drains.
		for i := range m.r_drainage {
			m.r_drainage[i] = -1
		}

		// Reset pools.
		for i := range m.r_pool {
			m.r_pool[i] = 0
		}
		m.r_elevation = m.fillSinks()

		// TODO: Diffuse flux and pool.
		m.assignRainfallBasic()

		// Regenerate downhill.
		m.BaseObject.assignDownhill(true)

		// Regenerate flux.
		m.assignFlux(false)
		// Erode a little.
		m.r_elevation = m.rErode(erosionAmount) // NOTE: Flux would change as downhill values would change.
	}
}

// assignHydrologyWithFlooding will calculate river systems and water pools.
func (m *Map) assignHydrologyWithFlooding() {
	maxAttempts := 20
	floodVariant := FloodVariant2
	skipSinksBelowSea := true

	// Reset drains.
	for i := range m.r_drainage {
		m.r_drainage[i] = -1
	}

	// Try to flood all sinks.
	var attempts int
	for {
		// Identify sinks above sea level.
		r_sinks := m.BaseObject.getSinks(skipSinksBelowSea, true)

		// Abort if we have no more sinks or ran out of attempts.
		if len(r_sinks) == 0 || attempts > maxAttempts {
			m.r_elevation = m.fillSinks()
			// Regenerate downhill.
			m.BaseObject.assignDownhill(true)

			// Regenerate flux.
			m.assignFlux(true)

			// TODO: Diffuse flux and pool.
			m.assignRainfall(4, moistTransferDirect2)

			log.Println("ran out of attempts", len(r_sinks))
			// TODO: Fill remaining sinks and re-generate downhill and flux.
			break
		}
		attempts++

		// Now that we want to calculate the max. flux that accumulates in sinks,
		// we will have to disregard any drainage and waterpools.

		// Reset drains.
		for i := range m.r_drainage {
			m.r_drainage[i] = -1
		}

		// Reset pools.
		for i := range m.r_pool {
			m.r_pool[i] = 0
		}

		// Erode a little.
		// m.r_elevation = m.rErode(0.001) // NOTE: Flux would change as downhill values would change.

		// Regenerate downhill and do not skip below sea level.
		m.BaseObject.assignDownhill(false)

		// Regenerate flux.
		m.assignFlux(true)

		// Identify sinks above sea level.
		r_sinks = m.BaseObject.getSinks(false, false)

		// Start from lowest sink.
		sort.Slice(r_sinks, func(i, j int) bool {
			return m.r_elevation[r_sinks[i]] < m.r_elevation[r_sinks[j]]
		})

		// Flood sink up to lowest neighbor + epsilon.
		for _, r := range r_sinks {
			//if m.r_flux[r] < m.r_rainfall[r] {
			//	continue
			//}
			switch floodVariant {
			case FloodVariant1:
				m.floodV1(r, m.r_flux[r])
			case FloodVariant2:
				m.floodV2(r, m.r_flux[r])
			}
		}

		// TODO: Diffuse flux and pool.
		m.assignRainfall(1, moistTransferDirect2)
	}

	// TODO: Triangle downhill.
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

const (
	FluxVolVariantBasic           = 0
	FluxVolVariantBasicWithDrains = 1
	FluxVolVariantWalk1           = 2
	FluxVolVariantWalk2           = 3
)

// assignFlux will populate r_flux by summing up the rainfall for each region from highest to
// lowest using r_downhill to reconstruct the downhill path that water would follow.
// NOTE: This is based on mewo2's terrain generation code
// See: https://github.com/mewo2/terrain
func (m *Map) assignFlux(skipBelowSea bool) {
	m.r_flux = m.getFlux(skipBelowSea)
}

func (m *Map) getFlux(skipBelowSea bool) []float64 {
	// Determines which flux calculation algorithm we use.
	variant := FluxVolVariantBasic

	// Initialize flux values with r_rainfall.
	r_flux := make([]float64, m.mesh.numRegions)
	for i := 0; i < m.mesh.numRegions; i++ {
		if m.r_elevation[i] >= 0 || !skipBelowSea {
			r_flux[i] = m.r_rainfall[i]
		}
	}

	switch variant {
	case FluxVolVariantBasic:
		idxs := make([]int, len(r_flux))
		for i := range r_flux {
			idxs[i] = i
		}
		sort.Slice(idxs, func(a, b int) bool {
			return m.r_elevation[idxs[b]]-m.r_elevation[idxs[a]] < 0
		})

		for _, j := range idxs {
			// Do not copy flux if we are below sea level.
			if (m.r_elevation[j] < 0 && skipBelowSea) || m.r_downhill[j] < 0 {
				continue
			}
			r_flux[m.r_downhill[j]] += r_flux[j]
		}
	case FluxVolVariantBasicWithDrains:
		// Basic variant copying the flux to the downhill neighbor or the drainage.
		// Initialize map for identifying drains and populate initial state of sorted index.
		drains := make(map[int]bool)
		idxs := make([]int, m.mesh.numRegions)
		for i := range idxs {
			if m.r_drainage[i] >= 0 {
				drains[m.r_drainage[i]] = true
				// r_flux[m.r_drainage[i]] += m.r_rainfall[i]
			}
			idxs[i] = i
		}

		// Sort index array.
		sort.Slice(idxs, func(a, b int) bool {
			if (m.r_elevation[idxs[b]] + m.r_pool[idxs[b]]) == (m.r_elevation[idxs[a]] + m.r_pool[idxs[a]]) {
				return drains[idxs[a]]
			}
			return (m.r_elevation[idxs[b]] + m.r_pool[idxs[b]]) < (m.r_elevation[idxs[a]] + m.r_pool[idxs[a]])
		})

		// Copy flux to known drainage point or next lowest neighbor.
		for _, j := range idxs {
			// Do not copy flux if we are below sea level.
			if m.r_elevation[j] < 0 && skipBelowSea {
				continue
			}

			// Check if we are entering a pool that drains somewhere else.
			if m.r_drainage[j] >= 0 {
				r_flux[m.r_drainage[j]] += r_flux[j]
			} else if m.r_downhill[j] >= 0 {
				r_flux[m.r_downhill[j]] += r_flux[j]
			}
		}
	case FluxVolVariantWalk1:
		r_flux_tmp := make([]float64, m.mesh.numRegions)
		for j, fl := range r_flux {
			seen := make(map[int]bool)
			drain := m.r_drainage[j]
			for drain != -1 {
				if m.r_elevation[drain] < 0 && skipBelowSea {
					break
				}
				r_flux_tmp[drain] += fl
				if m.r_drainage[drain] >= 0 && !seen[drain] {
					drain = m.r_drainage[drain]
				} else if m.r_downhill[drain] >= 0 {
					drain = m.r_downhill[drain]
				} else {
					drain = -1
				}
				seen[drain] = true
			}
		}

		// Copy the flux to the resulting flux map.
		for r, fl := range r_flux_tmp {
			r_flux[r] += fl
		}
	case FluxVolVariantWalk2:
		r_flux_tmp := make([]float64, m.mesh.numRegions)
		for j, fl := range r_flux {
			seen := make(map[int]bool)
			r := j
			var chain []int
			for {
				seen[r] = true
				if m.r_drainage[r] >= 0 {
					r = m.r_drainage[r]
				} else {
					r = m.r_downhill[r]
				}
				if r < 0 || m.r_elevation[r] < 0 && skipBelowSea {
					break
				}
				if seen[r] {
					break // Abort to avoid circular references.
				}
				chain = append(chain, r)
				r_flux_tmp[r] += fl
			}
			//r_flux[m.r_drainage[j]] += r_flux[j]
		}

		// Copy the flux to the resulting flux map.
		for r, fl := range r_flux_tmp {
			r_flux[r] += fl
		}
	}
	return r_flux
}

func (m *Map) floodV1(r int, dVol float64) {
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
			// d.sediment *= 0.1
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
	var drain, drainedFrom int
	drainedFrom = -1

	useDrain := true // Use drainage point instead of region draining into drainage point.

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
				if drainfound { // && drainedFrom == -1
					newDrain := -1
					if useDrain {
						newDrain = drain
					} else {
						newDrain = i
					}
					if drainedFrom == -1 || m.r_elevation[newDrain]+m.r_pool[newDrain] < m.r_elevation[drainedFrom]+m.r_pool[drainedFrom] {
						drainedFrom = newDrain
					}
				}
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
			if s != drainedFrom {
				m.r_drainage[s] = drainedFrom // WROOOOONG?????
			}
		}
		delete(boundary, minboundFirst)
		tried[minboundFirst] = false
		r = minboundFirst
	}

	if drainfound {
		if true {
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

			nbs := m.rNeighbors(drain)
			sort.Slice(nbs, func(si, sj int) bool {
				return m.r_elevation[nbs[si]]+m.r_pool[nbs[si]] < m.r_elevation[nbs[sj]]+m.r_pool[nbs[sj]]
			})

			// Fill Neighbors
			for _, neighbor_r := range nbs {
				lowbound(neighbor_r)
			}
			/*
				// Fill Neighbors
				for _, nbs := range m.rNeighbors(drain) {
					lowbound(nbs)
					//for _, nbs2 := range m.rNeighbors(nbs) { // ??????
					//	lowbound(nbs2)
					//}
				}
			*/
		}
		// Water-Level to Plane-Height
		for _, s := range floodset {
			// volume += ((plane > h[ind])?(h[ind] + p[ind] - plane):p[ind])/volumeFactor;
			if plane > m.r_elevation[s] {
				m.r_pool[s] = plane - m.r_elevation[s]
				if s != drainedFrom {
					m.r_drainage[s] = drainedFrom
				}
			} else {
				m.r_pool[s] = 0.0
				m.r_drainage[s] = -1
			}
		}

		for bfirst := range boundary {
			// volume += ((plane > h[ind])?(h[ind] + p[ind] - plane):p[ind])/volumeFactor;
			if plane > m.r_elevation[bfirst] {
				m.r_pool[bfirst] = plane - m.r_elevation[bfirst]
				if bfirst != drainedFrom {
					m.r_drainage[bfirst] = drainedFrom
				}
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
	halfedges := m.mesh.Halfedges
	for i := len(order_t) - 1; i >= 0; i-- {
		// TODO: Describe what's going on here.
		tributary_t := order_t[i]
		flow_s := t_downflow_s[tributary_t]
		if flow_s >= 0 {
			trunk_t := (halfedges[flow_s] / 3)
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
