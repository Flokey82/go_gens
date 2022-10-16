package genworldvoronoi

import (
	"container/heap"
	"log"
	"math"
	"sort"
	"time"
)

// getRivers returns the merged river segments whose flux exceeds the provided limit.
// Each river is represented as a sequence of region indices.
func (m *BaseObject) getRivers(limit float64) [][]int {
	// Get segments that are valid river segments.
	links := m.getRiverSegments(limit)

	// Merge the segments that are connected to each other into logical region sequences.
	log.Println("start merge")
	start := time.Now()
	defer func() {
		log.Println("Done river segments in ", time.Since(start).String())
	}()
	return mergeIndexSegments(links)
}

// getRivers2 is a modified, simpler version of getRivers and getRiverSegments.
// NOTE: Not sure why this exists, since getRivers works just fine?
func (m *Map) getRivers2(limit float64) [][]int {
	dh := m.r_downhill
	flux := m.r_flux

	// Adjust the limit to be a fraction of the max flux.
	// This will save us a lot of cycles when comparing
	// flux values to the limit.
	_, maxFlux := minMax(flux)
	limit *= maxFlux

	var links [][2]int
	for r := 0; r < m.mesh.numRegions; r++ {
		if flux[r] > limit && m.r_elevation[r] > 0 && dh[r] >= 0 && flux[dh[r]] > limit {
			// up := r, down := dh[r]
			links = append(links, [2]int{r, dh[r]})
		}
	}
	log.Println("start merge")
	return mergeIndexSegments(links)
}

// getRiverSegments returns all region / downhill neighbor pairs whose flux values
// exceed the provided limit / threshold.
func (m *BaseObject) getRiverSegments(limit float64) [][2]int {
	// NOTE: Should we re-generate downhill and flux, just in case erosion
	// or other factors might have changed this?

	// Get (cached) downhill neighbors.
	dh := m.r_downhill

	// Get (cached) flux values.
	flux := m.r_flux

	// Adjust the limit to be a fraction of the max flux.
	// This will save us a lot of cycles when comparing
	// flux values to the limit.
	_, maxFlux := minMax(flux)
	limit *= maxFlux

	// Find all link segments that have a high enough flux value.
	var links [][2]int
	for r := 0; r < m.mesh.numRegions; r++ {
		// Skip all regions that are sinks / have no downhill neighbor or
		// regions below sea level.
		if dh[r] < 0 || m.r_elevation[r] < 0 {
			continue
		}

		// Skip all regions with flux values that are equal to the rainfall in the region,
		// which is the minimum flux value / the default state for regions without
		// water influx.
		// NOTE: Rivers need at least one contributor region and would therefore have a flux
		// value that is higher than the rainfall in the region.
		if flux[r] <= m.r_rainfall[r] || flux[dh[r]] <= m.r_rainfall[dh[r]] {
			continue
		}

		// NOTE: Right now we skip segments if both flux values are
		// below the limit.
		if flux[r] >= limit && flux[dh[r]] >= limit {
			// NOTE: The river segment always flows from seg[0] to seg[1].
			links = append(links, [2]int{r, dh[r]})
		}
	}
	return links
}

// getWaterBodies returns a slice which all regions to enumerated waterbodies/oceans.
//
// NOTE: For regions that are not part of an ocean (elevation above sea level)
// a value of -2 is assigned.
func (m *BaseObject) getWaterBodies() []int {
	// Initialize the waterbody (ocean) mapping.
	done := make([]int, m.mesh.numRegions)
	for i := range done {
		if m.r_elevation[i] > 0 {
			done[i] = -2 // Non-ocean regions above sealevel.
		} else {
			done[i] = -1 // Ocean regions that have not been visited yet.
		}
	}

	for r := range done {
		// Skip regions that have already been visited or that are
		// non-ocean / above sealevel.
		if done[r] != -1 {
			continue
		}
		// Set the region index (r) as the ID for the new waterbody.
		done[r] = r

		// diveDeeper is a recursive function that performs a sort
		// of flood fill, assigning the current waterbody ID to all
		// neighboring regions that are ocean regions.
		//
		// TODO: Maybe use a queue instead... we might exceed Go's
		// stack size calling this recursively regardless of how deep
		// the execution stack might go.
		var diveDeeper func(rd int)
		diveDeeper = func(rd int) {
			for _, nbs := range m.rNeighbors(rd) {
				// If we have reached land or already visited nbs, skip.
				if m.r_elevation[nbs] > 0 || done[nbs] != -1 {
					continue
				}
				// Assign the source region index to nbs.
				done[nbs] = r

				// Visit neighbors of nbs.
				diveDeeper(nbs)
			}
		}

		// Recursively assign the waterbody ID / region index (r)
		// to all suitable neighbor regions and their neighbors,
		// and so on.
		diveDeeper(r)
	}
	return done
}

// getWaterBodySizes return a mapping of waterbody IDs to the number of regions
// associated with each waterbody ID.
func (m *BaseObject) getWaterBodySizes() map[int]int {
	wbSize := make(map[int]int)
	for _, wb := range m.r_waterbodies {
		if wb >= 0 {
			wbSize[wb]++ // Only count regions that are set to a valid ID.
		}
	}
	return wbSize
}

// getLakeSizes returns a mapping of drainage region to the number of regions that
// drain to this point, effectively summing up the size of each lake.
func (m *BaseObject) getLakeSizes() map[int]int {
	lakeSize := make(map[int]int)
	for _, drain := range m.r_drainage {
		if drain != -1 {
			lakeSize[drain]++ // Only count regions that have a drainage point assigned.
		}
	}
	return lakeSize
}

func (m *BaseObject) isRBelowOrAtSeaLevelOrPool(r int) bool {
	return m.r_elevation[r] <= 0 || m.r_pool[r] > 0
}

func (m *BaseObject) isRLakeOrWaterBody(r int) bool {
	return m.r_waterbodies[r] >= 0 || m.r_drainage[r] >= 0
}

func (m *BaseObject) isRiver(r int) bool {
	return m.r_flux[r] > m.r_rainfall[r]
}

// getRLakeOrWaterBodySiyze returns the size of the lake or waterbody that the
// provided region is part of.
func (m *BaseObject) getRLakeOrWaterBodySize(r int) int {
	if m.r_waterbodies[r] >= 0 {
		return m.r_waterbody_size[m.r_waterbodies[r]]
	}
	if m.r_drainage[r] >= 0 {
		return m.r_lake_size[m.r_drainage[r]]
	}
	return 0
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

// mergeIndexSegments matches up the ends of the segments (region pairs) and returns
// a slice containing all continuous, connected segments as sequence of connected regions.
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

// fillSinks attempts to remove sinks in the terrain using different approaches.
// NOTE: This is very much WIP, so forgive me for the mess.
func (m *Map) fillSinks() []float64 {
	// This will define which algorithm to use.
	usePlanchonDarboux := true

	// Use alternative fill sink algorithm?
	if usePlanchonDarboux {
		// Return sinks filled with Plançon-Darboux.
		return m.fillSinksPlanchonDarboux()
	}

	epsilon := 1e-13

	// Gather all regions that are below sealevel.
	var sea_r []int
	is_sea := make([]bool, m.mesh.numRegions)
	for r, h := range m.r_elevation {
		if h < 0 {
			sea_r = append(sea_r, r)
			is_sea[r] = true
		}
	}

	// Find all regions that are seeds for mountains by running the
	// continent collision code.
	//
	// NOTE: This is very inefficient to call this each time.
	mountain_r, _, _, _ := m.findCollisions()

	log.Println("Reminder for dev: fill sink epsilon might be too big!")

	// Calculate the distance field that gives us the distance of each region
	// to the sea tiles.
	r_distance_c := m.assignDistanceField(sea_r, convToMap(mountain_r))

	// Loop until we didn't change the terrain anymore.
	for {
		var changed bool

		// Get all sinks (regions without downhill neighbor).
		r_sinks := m.BaseObject.getSinks(true, false)

		log.Println("sinks remaining: ", len(r_sinks))

		// Sort the sinks by the distance to the ocean.
		sort.Slice(r_sinks, func(i, j int) bool {
			// Sort in ascending order.
			// The sinks closest to the ocean come first etc.
			return r_distance_c[r_sinks[i]] < r_distance_c[r_sinks[j]]
		})

		// Iterate over the sinks.
		for _, r := range r_sinks {
			// Skip sinks that are in the ocean or don't have any waterflow.
			// We only fill sinks to avoid issues when generating rivers.
			if m.r_flux[r] == 0 || is_sea[r] {
				continue
			}

			// Get the neighbor that is closest to the ocean with the lowest elevation.
			closest := r_distance_c[r]
			lowest := 9999.0
			rn := -1
			for _, nb := range m.rNeighbors(r) {
				if closest > r_distance_c[nb] && lowest > m.r_elevation[nb] {
					closest = r_distance_c[nb]
					lowest = m.r_elevation[nb]
					rn = nb
				}
			}

			// If we have found a suitable neighbor, we increase the elevation of
			// the sink while reducing the height of the neighbor closer to the ocean.
			// This will "dig a trench" towards the ocean.
			//
			// TODO:
			// - If we didn't find a suitable neighbor, we could still fill the sink a little?
			// - Reducing the neighbor by the elevation delta is a bit brutal, we could do a 50/50 split.
			if rn >= 0 {
				if m.r_elevation[rn] > m.r_elevation[r] {
					// If the elevation of the neighbor is higher than the sinks', we reduce
					// the neighbor's elevation by the elevation delta while increasing the
					// elevation of the sink by epsilon.
					delta := m.r_elevation[rn] - m.r_elevation[r]
					m.r_elevation[r] += epsilon
					m.r_elevation[rn] -= delta
					changed = true
				} else if m.r_elevation[rn] == m.r_elevation[r] {
					// If the elevation happens to be identical, we reduce the elevation of
					// the neighbor by epsilon.
					m.r_elevation[rn] -= epsilon
					changed = true
				}
			}
		}

		// If no change has occured in this pass, we can stop.
		if !changed {
			break
		}
	}

	// This pass is currently deactivated and contains another variant of the above.
	// Yeah... I can't really remember what my intention were here.
	if false {
		sort_r := make([]int, m.mesh.numRegions)
		for r := range m.r_elevation {
			sort_r[r] = r
		}
		sort.Slice(sort_r, func(i, j int) bool {
			// If the distance to the sea is the same, sort by ascending elevation.
			if r_distance_c[sort_r[i]] == r_distance_c[sort_r[j]] {
				return m.r_elevation[sort_r[i]] < m.r_elevation[sort_r[j]]
			}
			// Sort by distance to the sea in ascending order.
			return r_distance_c[sort_r[i]] < r_distance_c[sort_r[j]]
		})
		for {
			var changed bool
			for _, r := range sort_r {
				// Skip all ocean regions.
				if is_sea[r] {
					continue
				}

				// Iterate over all neighbors.
				for _, nb := range m.rNeighbors(r) {
					// If the distance to the sea for r is shorter than the distance for
					// nb, we reduce the elevation of r by the elevation delta and increase
					// the elevation of nb by epsilon.
					a, b := r, nb
					if r_distance_c[r] > r_distance_c[nb] {
						a, b = nb, r // If nb is closer to the sea than nb, we do the oposite to the above.
					}
					if m.r_elevation[a] > m.r_elevation[b] {
						delta := m.r_elevation[a] - m.r_elevation[b]
						m.r_elevation[b] += epsilon // delta / 2
						m.r_elevation[a] -= delta   // (delta / 2) + epsilon
						changed = true
					} else if m.r_elevation[b] == m.r_elevation[a] {
						m.r_elevation[a] -= epsilon
						changed = true
					}

					// NOTE: The old code did not alter elevations if the distance to the ocean was the same?
				}
			}

			// If no change has occured in this pass, we can stop.
			if !changed {
				break
			}
		}
	}
	return m.r_elevation
}

// fillSinksPlanchonDarboux is an implementation of the algorithm described in
// https://www.researchgate.net/publication/240407597_A_fast_simple_and_versatile_algorithm_to_fill_the_depressions_of_digital_elevation_models
// and a partial port of the implementation in:
// https://github.com/Rob-Voss/Learninator/blob/master/js/lib/Terrain.js
//
// NOTE: This algorithm produces a too uniform result at the moment, resulting
// in very artificially looking rivers. It lacks some kind of variation like
// noise. It's very fast and less destructive than my other, home-grown algorithm.
// Maybe it's worth to combine the two in some way?
func (m *Map) fillSinksPlanchonDarboux() []float64 {
	// Reset the RNG.
	m.resetRand()

	inf := math.Inf(0)
	baseEpsilon := 1.0 / (float64(m.mesh.numRegions) * 1000.0)
	newHeight := make([]float64, m.mesh.numRegions)
	for i := range newHeight {
		if m.r_elevation[i] <= 0 {
			// Set the elevation at or below sea level to the current
			// elevation.
			newHeight[i] = m.r_elevation[i]
		} else {
			// Set the elevation above sea level to infinity.
			newHeight[i] = inf
		}
	}

	// Loop until no more changes are made.
	var epsilon float64
	for {
		// Variation.
		//
		// In theory we could use noise or random values to slightly
		// alter epsilon here. It should still work, albeit a bit slower.
		// The idea is to make the algorithm less destructive and more
		// natural looking.
		//
		// NOTE: I've decided to use m.rand.Float64() instead of noise.
		epsilon = baseEpsilon * m.rand.Float64()

		changed := false

		// By shuffling the order in which we parse regions,
		// we ensure a more natural look.
		for _, r := range m.rand.Perm(len(m.r_elevation)) {
			// Skip all regions that have the same elevation as in
			// the current heightmap.
			if newHeight[r] == m.r_elevation[r] {
				continue
			}

			// Iterate over all neighbors in a random order.
			nbs := m.rNeighbors(r)
			for _, i := range m.rand.Perm(len(nbs)) {
				nb := nbs[i]
				// Since we have set all inland regions to infinity,
				// we will only succeed here if the newHeight of the neighbor
				// is either below sea level or if the newHeight has already
				// been set AND if the elevation is higher than the neighbors.
				//
				// This means that we're working our way inland, starting from
				// the coast, comparing each region with the processed / set
				// neighbors (that aren't set to infinity) in the new heightmap
				// until we run out of regions that need change.
				if m.r_elevation[r] >= newHeight[nb]+epsilon {
					newHeight[r] = m.r_elevation[r]
					changed = true
					break
				}

				// If we reach this point, the neighbor in the new heightmap
				// is higher than the current elevation of 'r'.
				// This can mean two things. Either the neighbor is set to infinity
				// or the current elevation might indicate a sink.

				// So we check if the newHeight of r is larger than the
				// newHeight of the neighbor (plus epsilon), which will ensure that
				// the newHeight of neighbor is not set to infinity.
				//
				// Additionally we check if the newHeight of the neighbor
				// is higher than the current height of r, which ensures that if the
				// current elevation indicates a sink, we will fill up the sink to the
				// new neighbor height plus epsilon.
				//
				// TODO: Simplify this comment word salad.
				oh := newHeight[nb] + epsilon
				if newHeight[r] > oh && oh > m.r_elevation[r] {
					newHeight[r] = oh
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	return newHeight
}

// assignHydrology will calculate river systems and fill sinks instead of trying to generate
// water pools.
func (m *Map) assignHydrology() {
	maxAttempts := 3
	erosionAmount := 0.01 // Erode 1% of delta-h per pass.

	// HACK: Fill all sinks that are below sea level and a single region
	// below sea level.
Loop:
	for _, r := range m.getSinks(false, false) {
		// Check if all neighbors are above sea level.
		lowest := math.Inf(0)
		for _, nb := range m.rNeighbors(r) {
			if !m.isRBelowOrAtSeaLevelOrPool(r) {
				continue Loop
			}
			if m.r_elevation[nb] < lowest {
				lowest = m.r_elevation[nb]
			}
		}
		m.r_elevation[r] = lowest
	}

	// Start off by filling sinks.
	m.r_elevation = m.fillSinks()

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

	// TODO: Move this somewhere else.
	m.r_lake_size = m.getLakeSizes()
	// TODO: Make note of oceans.
	//   - Note ocean sizes (and small waterbodies below sea level)
	m.r_waterbodies = m.getWaterBodies()
	m.r_waterbody_size = m.getWaterBodySizes()
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
			m.assignRainfall(4, moistTransferDirect, moistOrderWind)

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
		// m.r_elevation = m.rErode(0.01) // NOTE: Flux would change as downhill values would change.

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
		m.assignRainfall(1, moistTransferDirect, moistOrderWind)
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

// getFlux calculates and returns the water flux values for each region.
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
		// This is most basic flux calculation.
		// Sort regions by elevation in descending order.
		idxs := make([]int, len(r_flux))
		for i := range r_flux {
			idxs[i] = i
		}
		sort.Slice(idxs, func(a, b int) bool {
			return m.r_elevation[idxs[a]] > m.r_elevation[idxs[b]]
		})

		// Highest elevation first.
		for _, r := range idxs {
			// Skip calculation if we are below sea level or there is no downhill
			// neighbor where the water could flow to.
			if (m.r_elevation[r] < 0 && skipBelowSea) || m.r_downhill[r] < 0 {
				continue
			}

			// Add the flux of the region to the downhill neighbor.
			r_flux[m.r_downhill[r]] += r_flux[r]
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
				// If there is a drainage point set for the current region,
				// which indicates that this region is part of a lake.
				// In this case we copy the flux directly to the region where
				// this region drains into.
				r_flux[m.r_drainage[j]] += r_flux[j]
			} else if m.r_downhill[j] >= 0 {
				// Add the flux of the region to the downhill neighbor.
				r_flux[m.r_downhill[j]] += r_flux[j]
			}
		}
	case FluxVolVariantWalk1:
		// This seems incomplete as it will only calculate the flux
		// if a drainage point is set.
		// I put in a quick fix as I type this, but I didn't test the
		// result, so no guarantees.
		r_flux_tmp := make([]float64, m.mesh.numRegions)
		for j, fl := range r_flux {
			seen := make(map[int]bool)
			drain := m.r_drainage[j]
			if drain == -1 {
				drain = m.r_downhill[j]
			}
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
		// This variant will walk downhill for each region until we
		// can't find neither a downhill neighbor nor a drainage point.
		r_flux_tmp := make([]float64, m.mesh.numRegions)
		for j, fl := range r_flux {
			// Seen will keep track of the regions that we have
			// already visited for this region. This will prevent
			// any infinite recursions that might be caused by
			// drains that loop back to themselves.
			//
			// Not ideal, but it is what it is.
			seen := make(map[int]bool)

			// Start at region r.
			r := j

			// Not sure why I kept the chain of visited regions but I
			// assume it is useful for debugging.
			// var chain []int
			for {
				// Note that we've visited the region.
				seen[r] = true

				// Check if we have a drainage point.
				if m.r_drainage[r] >= 0 {
					r = m.r_drainage[r] // continue with drainage point
				} else {
					r = m.r_downhill[r] // use downhill neighbor
				}

				// If we couldn't find a region to drain into, or if
				// we are below sea level, stop here.
				if r < 0 || m.r_elevation[r] < 0 && skipBelowSea {
					break
				}
				// Abort if we have already visited r to avoid circular
				// references.
				if seen[r] {
					break
				}
				// chain = append(chain, r)
				r_flux_tmp[r] += fl
			}
			// Not sure why this was here.
			// r_flux[m.r_drainage[j]] += r_flux[j]
		}

		// Copy the flux to the resulting flux map.
		for r, fl := range r_flux_tmp {
			r_flux[r] += fl
		}
	}
	return r_flux
}

// floodV1 is the first variant of the flood fill algorithm, which finds
// drainage points for sinks and generates lakes.
// Don't ask me how this works in detail as I do not know.
//
// Only thing I know is that it is based on Nick McDonald's old flood fill
// he used in simple_hydrology.
//
// TODO: Return remaining volume
func (m *Map) floodV1(r int, dVol float64) {
	const (
		volumeFactor = 100.0 // "Water Deposition Rate"
		epsilon      = 1e-3
		minVol       = 0.01
		drainage     = 0.01
	)

	plane := m.r_elevation[r] + m.r_pool[r]
	initialplane := plane

	// Floodset contains all regions that are part of a floodplain.
	set := make([]int, 0, 1024)

	// Abort after 200 attempts.
	fail := 200

	// Keep track of the regions we have visitad during a flood fill attempt.
	tried := make([]bool, m.mesh.numRegions)
	var drain int
	var drainfound bool
	var fill func(i int)
	fill = func(i int) {
		// Out of bounds, or region has been visited ("tried") previously.
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

		// Pre-sort neighbors by height (elevation + water pool).
		//
		// NOTE: The regions are sorted in ascending order, so the first
		// region in the list will be the lowest one.
		sort.Slice(nbs, func(si, sj int) bool {
			return m.r_elevation[nbs[si]]+m.r_pool[nbs[si]] < m.r_elevation[nbs[sj]]+m.r_pool[nbs[sj]]
		})

		// Expand floodset by attempting to fill all neighbors.
		for _, neighbor_r := range nbs {
			fill(neighbor_r)
		}
	}

	// Iterate
	for dVol > minVol && fail != 0 {
		set = set[:0]

		// Reset the visited regions.
		for i := range tried {
			tried[i] = false
		}

		// Reset the drain and drainfound flag.
		drain = 0
		drainfound = false

		// Perform flooding of initial region.
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

		// Adjust planes.
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
//
// NOTE: This is the original code that Amit uses in his procedural planets project.
// He uses triangle centroids for his river generation, where I prefer to use the regions
// directly.
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
//
// NOTE: This is the original code that Amit uses in his procedural planets project.
// He uses triangle centroids for his river generation, where I prefer to use the regions
// directly.
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
