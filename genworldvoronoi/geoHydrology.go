package genworldvoronoi

import (
	"log"
	"math"
	"sort"
)

const (
	FloodVariant1 = 0
	FloodVariant2 = 1
)

// assignHydrology will calculate river systems and fill sinks instead of trying to generate
// water pools.
func (m *Geo) assignHydrology() {
	maxAttempts := 3
	erosionAmount := 0.01 // Erode 1% of delta-h per pass.

	// HACK: Fill all sinks that are below sea level and a single region
	// below sea level.
Loop:
	for _, r := range m.GetSinks(false, false) {
		// Check if all neighbors are above sea level.
		lowest := math.Inf(0)
		for _, nb := range m.GetRegNeighbors(r) {
			if !m.isRegBelowOrAtSeaLevelOrPool(r) {
				continue Loop
			}
			if m.Elevation[nb] < lowest {
				lowest = m.Elevation[nb]
			}
		}
		m.Elevation[r] = lowest
	}

	// Start off by filling sinks.
	m.Elevation = m.FillSinks()

	// Try to flood all sinks.
	var attempts int
	m.BaseObject.assignDownhill(true)
	m.assignFlux(false)
	for {
		// Abort if we have no more sinks or ran out of attempts.
		if attempts > maxAttempts {
			m.Elevation = m.FillSinks()
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
		for i := range m.Drainage {
			m.Drainage[i] = -1
		}

		// Reset pools.
		for i := range m.Waterpool {
			m.Waterpool[i] = 0
		}
		m.Elevation = m.FillSinks()

		// TODO: Diffuse flux and pool.
		m.assignRainfallBasic()

		// Regenerate downhill.
		m.BaseObject.assignDownhill(true)

		// Regenerate flux.
		m.assignFlux(false)

		// Erode a little.
		m.Elevation = m.Erode(erosionAmount) // NOTE: Flux would change as downhill values would change.
	}

	// TODO: Move this somewhere else.
	m.LakeSize = m.getLakeSizes()
	// TODO: Make note of oceans.
	//   - Note ocean sizes (and small waterbodies below sea level)
	m.assignWaterbodies()
}

// assignHydrologyWithFlooding will calculate river systems and water pools.
func (m *Geo) assignHydrologyWithFlooding() {
	maxAttempts := 20
	floodVariant := FloodVariant2
	skipSinksBelowSea := true

	// Reset drains.
	for i := range m.Drainage {
		m.Drainage[i] = -1
	}

	// Try to flood all sinks.
	var attempts int
	for {
		// Identify sinks above sea level.
		r_sinks := m.BaseObject.GetSinks(skipSinksBelowSea, true)

		// Abort if we have no more sinks or ran out of attempts.
		if len(r_sinks) == 0 || attempts > maxAttempts {
			m.Elevation = m.FillSinks()
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
		for i := range m.Drainage {
			m.Drainage[i] = -1
		}

		// Reset pools.
		for i := range m.Waterpool {
			m.Waterpool[i] = 0
		}

		// Erode a little.
		// m.r_elevation = m.rErode(0.01) // NOTE: Flux would change as downhill values would change.

		// Regenerate downhill and do not skip below sea level.
		m.BaseObject.assignDownhill(false)

		// Regenerate flux.
		m.assignFlux(true)

		// Identify sinks above sea level.
		r_sinks = m.BaseObject.GetSinks(false, false)

		// Start from lowest sink.
		sort.Slice(r_sinks, func(i, j int) bool {
			return m.Elevation[r_sinks[i]] < m.Elevation[r_sinks[j]]
		})

		// Flood sink up to lowest neighbor + epsilon.
		for _, r := range r_sinks {
			//if m.r_flux[r] < m.r_rainfall[r] {
			//	continue
			//}
			switch floodVariant {
			case FloodVariant1:
				m.floodV1(r, m.Flux[r])
			case FloodVariant2:
				m.floodV2(r, m.Flux[r])
			}
		}

		// TODO: Diffuse flux and pool.
		m.assignRainfall(1, moistTransferDirect, moistOrderWind)
	}

	// TODO: Triangle downhill.
	// TODO: Make note of lakes.
	//   - Sum up regions r_pool[r] > 0
	//   - Note lake sizes (for city placement)
	m.LakeSize = m.getLakeSizes()
	// TODO: Make note of rivers.
	m.assignWaterbodies()
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
func (m *Geo) assignFlux(skipBelowSea bool) {
	m.Flux = m.getFlux(skipBelowSea)
}

// getFlux calculates and returns the water flux values for each region.
func (m *Geo) getFlux(skipBelowSea bool) []float64 {
	// Determines which flux calculation algorithm we use.
	variant := FluxVolVariantBasic

	// Initialize flux values with r_rainfall.
	regFlux := make([]float64, m.mesh.numRegions)
	for i := 0; i < m.mesh.numRegions; i++ {
		if m.Elevation[i] >= 0 || !skipBelowSea {
			regFlux[i] = m.Rainfall[i]
		}
	}

	switch variant {
	case FluxVolVariantBasic:
		// This is most basic flux calculation.
		// Sort regions by elevation in descending order.
		idxs := make([]int, len(regFlux))
		for i := range regFlux {
			idxs[i] = i
		}
		sort.Slice(idxs, func(a, b int) bool {
			return m.Elevation[idxs[a]] > m.Elevation[idxs[b]]
		})

		// Highest elevation first.
		for _, r := range idxs {
			// Skip calculation if we are below sea level or there is no downhill
			// neighbor where the water could flow to.
			// NOTE: In this case we allow water to flow to sea level.
			if (m.Elevation[r] < 0 && skipBelowSea) || m.Downhill[r] < 0 {
				continue
			}

			// Add the flux of the region to the downhill neighbor.
			regFlux[m.Downhill[r]] += regFlux[r]
		}
	case FluxVolVariantBasicWithDrains:
		// Basic variant copying the flux to the downhill neighbor or the drainage.
		// Initialize map for identifying drains and populate initial state of sorted index.
		drains := make(map[int]bool)
		idxs := make([]int, m.mesh.numRegions)
		for i := range idxs {
			if m.Drainage[i] >= 0 {
				drains[m.Drainage[i]] = true
				// r_flux[m.r_drainage[i]] += m.r_rainfall[i]
			}
			idxs[i] = i
		}

		// Sort index array.
		sort.Slice(idxs, func(a, b int) bool {
			if (m.Elevation[idxs[b]] + m.Waterpool[idxs[b]]) == (m.Elevation[idxs[a]] + m.Waterpool[idxs[a]]) {
				return drains[idxs[a]]
			}
			return (m.Elevation[idxs[b]] + m.Waterpool[idxs[b]]) < (m.Elevation[idxs[a]] + m.Waterpool[idxs[a]])
		})

		// Copy flux to known drainage point or next lowest neighbor.
		for _, j := range idxs {
			// Do not copy flux if we are below sea level.
			// NOTE: In this case we allow water to flow to sea level.
			if m.Elevation[j] < 0 && skipBelowSea {
				continue
			}

			// Check if we are entering a pool that drains somewhere else.
			if m.Drainage[j] >= 0 {
				// If there is a drainage point set for the current region,
				// which indicates that this region is part of a lake.
				// In this case we copy the flux directly to the region where
				// this region drains into.
				regFlux[m.Drainage[j]] += regFlux[j]
			} else if m.Downhill[j] >= 0 {
				// Add the flux of the region to the downhill neighbor.
				regFlux[m.Downhill[j]] += regFlux[j]
			}
		}
	case FluxVolVariantWalk1:
		// This seems incomplete as it will only calculate the flux
		// if a drainage point is set.
		// I put in a quick fix as I type this, but I didn't test the
		// result, so no guarantees.
		regFluxTmp := make([]float64, m.mesh.numRegions)
		for j, fl := range regFlux {
			seen := make(map[int]bool)
			drain := m.Drainage[j]
			if drain == -1 {
				drain = m.Downhill[j]
			}
			for drain != -1 {
				// NOTE: In this case we allow water to flow to sea level.
				if m.Elevation[drain] < 0 && skipBelowSea {
					break
				}
				regFluxTmp[drain] += fl
				if m.Drainage[drain] >= 0 && !seen[drain] {
					drain = m.Drainage[drain]
				} else if m.Downhill[drain] >= 0 {
					drain = m.Downhill[drain]
				} else {
					drain = -1
				}
				seen[drain] = true
			}
		}

		// Copy the flux to the resulting flux map.
		for r, fl := range regFluxTmp {
			regFlux[r] += fl
		}
	case FluxVolVariantWalk2:
		// This variant will walk downhill for each region until we
		// can't find neither a downhill neighbor nor a drainage point.
		regFluxTmp := make([]float64, m.mesh.numRegions)
		for j, fl := range regFlux {
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
				if m.Drainage[r] >= 0 {
					r = m.Drainage[r] // continue with drainage point
				} else {
					r = m.Downhill[r] // use downhill neighbor
				}

				// If we couldn't find a region to drain into, or if
				// we are below sea level, stop here.
				// NOTE: In this case we allow water to flow to sea level.
				if r < 0 || m.Elevation[r] < 0 && skipBelowSea {
					break
				}
				// Abort if we have already visited r to avoid circular
				// references.
				if seen[r] {
					break
				}
				// chain = append(chain, r)
				regFluxTmp[r] += fl
			}
			// Not sure why this was here.
			// r_flux[m.r_drainage[j]] += r_flux[j]
		}

		// Copy the flux to the resulting flux map.
		for r, fl := range regFluxTmp {
			regFlux[r] += fl
		}
	}
	return regFlux
}

// floodV1 is the first variant of the flood fill algorithm, which finds
// drainage points for sinks and generates lakes.
// Don't ask me how this works in detail as I do not know.
//
// Only thing I know is that it is based on Nick McDonald's old flood fill
// he used in simple_hydrology.
//
// TODO: Return remaining volume
func (m *Geo) floodV1(r int, dVol float64) {
	const (
		volumeFactor = 100.0 // "Water Deposition Rate"
		epsilon      = 1e-3
		minVol       = 0.01
		drainage     = 0.01
	)

	plane := m.Elevation[r] + m.Waterpool[r]
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
		currHeight := m.Elevation[i] + m.Waterpool[i]
		if plane < currHeight {
			return
		}

		// Drainage Point
		if initialplane > currHeight {
			if !drainfound || currHeight < m.Waterpool[drain]+m.Elevation[drain] {
				// No Drain yet or lower drain.
				drain = i
			}

			drainfound = true
			return
		}

		// Part of the Pool
		set = append(set, i)
		nbs := m.GetRegNeighbors(i)

		// Pre-sort neighbors by height (elevation + water pool).
		//
		// NOTE: The regions are sorted in ascending order, so the first
		// region in the list will be the lowest one.
		sort.Slice(nbs, func(si, sj int) bool {
			return m.Elevation[nbs[si]]+m.Waterpool[nbs[si]] < m.Elevation[nbs[sj]]+m.Waterpool[nbs[sj]]
		})

		// Expand floodset by attempting to fill all neighbors.
		for _, nbReg := range nbs {
			fill(nbReg)
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
			plane = (1.0-drainage)*initialplane + drainage*(m.Elevation[drain]+m.Waterpool[drain])

			// Compute the New Height
			for _, s := range set {
				if plane > m.Elevation[s] {
					m.Waterpool[s] = plane - m.Elevation[s]
					m.Drainage[s] = drain
				} else {
					m.Waterpool[s] = 0.0
					m.Drainage[s] = -1
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
			totalVol += volumeFactor * (plane - (m.Elevation[s] + m.Waterpool[s]))
		}
		// log.Println("totalVol", totalVol, "dVol", dVol, "setLen", len(set))
		// We can fill the volume of the sink.
		if totalVol <= dVol && initialplane < plane {
			// Raise water level to plane height.
			for _, s := range set {
				m.Waterpool[s] = plane - m.Elevation[s]
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
func (m *Geo) floodV2(r int, dVol float64) bool {
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
		currHeight := m.Elevation[i] + m.Waterpool[i]
		if plane < currHeight {
			boundary[i] = currHeight
			return true
		}

		// Drainage Point
		if currHeight < plane {
			// No Drain yet
			if !drainfound || currHeight < m.Waterpool[drain]+m.Elevation[drain] {
				drain = i
			}
			drainfound = true
			return false
		}

		// Part of the Pool
		floodset = append(floodset, i)
		nbs := m.GetRegNeighbors(i)
		sort.Slice(nbs, func(si, sj int) bool {
			return m.Elevation[nbs[si]]+m.Waterpool[nbs[si]] < m.Elevation[nbs[sj]]+m.Waterpool[nbs[sj]]
		})
		for _, nbReg := range nbs {
			if !findset(nbReg, plane) {
				if drainfound { // && drainedFrom == -1
					newDrain := -1
					if useDrain {
						newDrain = drain
					} else {
						newDrain = i
					}
					if drainedFrom == -1 || m.Elevation[newDrain]+m.Waterpool[newDrain] < m.Elevation[drainedFrom]+m.Waterpool[drainedFrom] {
						drainedFrom = newDrain
					}
				}
				return false
			}
		}
		return true
	}

	plane := m.Waterpool[r] + m.Elevation[r]
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
			m.Waterpool[s] = plane - m.Elevation[s]
			if s != drainedFrom {
				m.Drainage[s] = drainedFrom // WROOOOONG?????
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
				if i < 0 || m.Waterpool[i] == 0 {
					return
				}
				// Below Drain Height
				if m.Elevation[i]+m.Waterpool[i] < m.Elevation[drain]+m.Waterpool[drain] {
					return
				}
				// Higher than Plane (we want lower)
				if m.Elevation[i]+m.Waterpool[i] >= plane {
					return
				}
				plane = m.Elevation[i] + m.Waterpool[i]
			}

			nbs := m.GetRegNeighbors(drain)
			sort.Slice(nbs, func(si, sj int) bool {
				return m.Elevation[nbs[si]]+m.Waterpool[nbs[si]] < m.Elevation[nbs[sj]]+m.Waterpool[nbs[sj]]
			})

			// Fill Neighbors
			for _, nbReg := range nbs {
				lowbound(nbReg)
				// Fill neighbors of neighbors
				// for _, nbs2 := range m.rNeighbors(nbs) { // ??????
				//	lowbound(nbs2)
				// }
			}
		}

		// Water-Level to Plane-Height
		for _, s := range floodset {
			// volume += ((plane > h[ind])?(h[ind] + p[ind] - plane):p[ind])/volumeFactor;
			if plane > m.Elevation[s] {
				m.Waterpool[s] = plane - m.Elevation[s]
				if s != drainedFrom {
					m.Drainage[s] = drainedFrom
				}
			} else {
				m.Waterpool[s] = 0.0
				m.Drainage[s] = -1
			}
		}

		for bfirst := range boundary {
			// volume += ((plane > h[ind])?(h[ind] + p[ind] - plane):p[ind])/volumeFactor;
			if plane > m.Elevation[bfirst] {
				m.Waterpool[bfirst] = plane - m.Elevation[bfirst]
				if bfirst != drainedFrom {
					m.Drainage[bfirst] = drainedFrom
				}
			} else {
				m.Waterpool[bfirst] = 0.0
				m.Drainage[bfirst] = -1
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

// assignFlow calculates the water flux by traversing the graph generated with
// assignDownflow in reverse order (so, downhill?) and summing up the moisture.
//
// NOTE: This is the original code that Amit uses in his procedural planets project.
// He uses triangle centroids for his river generation, where I prefer to use the regions
// directly.
func (m *BaseObject) assignFlow() {
	sideFlow := m.sideFlow

	// Clear all existing water flux values.
	for i := range sideFlow {
		sideFlow[i] = 0
	}

	triFlow := m.triFlow
	triElevation := m.triElevation
	triMoisture := m.triMoisture

	// Set the flux value for each triangle above sealevel to
	// half of its moisture squared as its initial state.
	numTriangles := m.mesh.numTriangles
	for t := 0; t < numTriangles; t++ {
		if triElevation[t] >= 0.0 {
			triFlow[t] = 0.5 * triMoisture[t] * triMoisture[t]
		} else {
			triFlow[t] = 0
		}
	}

	// Now traverse the flux graph in reverse order and sum up
	// the moisture of all tributaries while descending.
	orderTris := m.orderTri
	triDownflowSide := m.triDownflowSide
	halfedges := m.mesh.Halfedges
	for i := len(orderTris) - 1; i >= 0; i-- {
		// TODO: Describe what's going on here.
		tributaryTri := orderTris[i]
		flowSide := triDownflowSide[tributaryTri]
		if flowSide >= 0 {
			trunkTri := (halfedges[flowSide] / 3)
			triFlow[trunkTri] += triFlow[tributaryTri]
			sideFlow[flowSide] += triFlow[tributaryTri] // TODO: isn't s_flow[flow_s] === t_flow[?]
			if triElevation[trunkTri] > triElevation[tributaryTri] {
				triElevation[trunkTri] = triElevation[tributaryTri]
			}
		}
	}
	m.triFlow = triFlow
	m.sideFlow = sideFlow
	m.triElevation = triElevation
}
