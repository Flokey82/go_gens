package genmapvoronoi

import (
	"math"

	"github.com/Flokey82/go_gens/vmesh"
)

// getSteepness returns the steepness per region.
//
// NOTE: We define steepness as the angle to a region from its downhill neighbor
// expressed as a value between 0.0 to 1.0 (representing an angle from 0° to 90°).
func getSteepness(h *vmesh.Heightmap) *vmesh.Heightmap {
	m := h.Mesh

	// This will collect the steepness for each region.
	steeps := vmesh.NewHeightmap(h.Mesh)

	// Get the downhill neighbors for all regions (ignoring water pools for now).
	dh := h.Downhill()
	for r, d := range dh {
		if d < 0 {
			continue // Skip all sinks.
		}
		// In order to calculate the steepness value, we get the great arc distance
		// of each region and its downhill neighbor, as well as the elevation change.
		//
		//     __r            r
		//     | |\            \
		//     | | \            \
		// height|  \            \
		//     | |   \            \
		//     |_|____\dh[r]   ____\dh[r] <- we want to calculate this angle
		//       |dist|
		//
		// We calculate the angle (in radians) as follows:
		// angle = atan(height/dist)
		//
		// Finally, to get the steepness in a range of 0.0 ... 1.0:
		// steepness = angle * 2 / Pi

		// Calculate height difference between r and dh[r].
		hDiff := h.Values[r] - h.Values[d]

		// Distance between r and dh[r].
		dist := distPoints(m.Vertices[r].X, m.Vertices[r].Y, m.Vertices[d].X, m.Vertices[d].Y)

		// Calculate the the angle (0°-90°) expressed as range from 0.0 to 1.0.
		steeps.Values[r] = math.Atan(hDiff/dist) * 2 / math.Pi
	}
	return steeps
}

// getRErosion2 is an alternative erosion calculation which takes in account
// the steepness and flux of each region to determine the shape of eroded
// riverbeds and valleys.
func erosionRate2(h *vmesh.Heightmap) *vmesh.Heightmap {
	const (
		// HACK: About 3 times the current average neighbor distance (at 16384 points)
		maxErosionDistance = 3 * 0.03
		minExp             = 1.0 // Minimum exponent (see below)
		varExp             = 3.0 // Variable exponent (see further below)
	)
	m := h.Mesh
	steeps := getSteepness(h)
	dh := h.Downhill()
	toE := vmesh.NewHeightmap(m)
	flux := getFlux(h)
	for r, fl := range flux.Values {
		// The steeper the incline, the more V-shaped the riverbed becomes,
		// while a river on a low incline will carve a more U shaped riverbed.
		//
		// y = pow(|x|, e)
		//
		// The exponent (e) determines the shape of the eroded ravine, or the
		// remaining height (y) at a certain distance (x) from the river center.
		//         _  _
		// e = 1:   \/
		//         _   _
		// e = 2:   \_/
		//         _   _
		// e = 4+:  |_|
		//
		// Therefore we calculate the exponent (e) for the function shaping
		// the riverbed (shape of erosion pattern) based on the steepness of the region.
		exponent := minExp + varExp*(1.0-steeps.Values[r])

		// The amount of flux determines how wide the resulting eroded path is.
		//
		// NOTE: We assume that a higher flux will mean that theoretically the resulting river
		// will become wider due to the additional required space for transporting the larger
		// water volume.
		//
		// NOTE: Since rivers slow down the wider they are, they start to be less erosive with higher flux.
		//
		// WARNING: The maxDist value currently represents both the max width of the riverbed as well as the
		// maximum erosion value at the river center (r). This is not really ideal, but it works.

		// Calculate maximum erosion distance based on the water flow intensity aka flux.
		maxDist := fl * maxErosionDistance // Flux times 3 times the current average neighbor distance (at 16384 points)

		// Assign the erosion intensity for (r) if it is larger
		// than the current value set for it.
		if toE.Values[r] < maxDist {
			toE.Values[r] = maxDist
			// TODO: sharp drops should carve with higher intensity.
			// toE.Values[r] *= (1 + steeps.Values[r]) / 2
		}
		rVertex := m.Vertices[r]

		// If we have a downhill neighbor, get the vertex for
		// calculating the river segment distance of each neighbor.
		dVertex := rVertex
		if rdh := dh[r]; rdh >= 0 {
			dVertex = m.Vertices[rdh]
		}

		// seen will keep track of all regions whose neighbors doErode()
		// has already visited (to prevent infinite recursion).
		seen := make(map[int]bool)

		// doErode() is the erosion function which will determine how
		// much a region will be eroded depending on the variables
		// we calculated above and the distance of the neighbor
		// region (reg) to the river center region (r).
		var doErode func(reg int)
		doErode = func(reg int) {
			// Skip visiting all neighbor regions we have already visited.
			if seen[reg] {
				return
			}
			// Make sure we note that we (will have) visited all neighbors of reg.
			seen[reg] = true

			// Visit all neighbors of reg and calculate the erosion rate for each.
			for _, nb := range h.Neighbours(reg) {
				nbVertex := m.Vertices[nb]

				// Calculate distance of the neighbor to the river segment (r->rdh).
				dist := distToSegment(rVertex.X, rVertex.Y, dVertex.X, dVertex.Y, nbVertex.X, nbVertex.Y)
				if dist > maxDist {
					continue // Skip everything that is too far away.
				}

				// Calculate distance as value from 0.0 to 1.0, representing
				// the range of 0.0 to maxDist.
				distRes := dist / maxDist

				// Now apply the formula that will shape the riverbed we
				// have described above.
				erode := math.Pow(distRes, exponent)

				// NOTE: Since the formula gives us a higher value with increasing
				// distance, we use the value to REDUCE the maximum erosion value
				// (maxDist). The further we are away from (r), the weaker the
				// erosion.
				toErode := maxDist - erode

				// If we have already an erosion value for (nb), we make sure
				// to use the higher erosion value.
				if toE.Values[nb] < toErode {
					toE.Values[nb] = toErode
				}

				// Visit all neighbors of nb and calculate their erosion values.
				doErode(nb)
			}
		}

		// Call the recursive erosion function for the current region (r).
		// This will calculate the erosion values for all neighbors of (r).
		doErode(r)
	}
	return toE
}

/*
// Garbage code that attempts to get the erosion to carve "meandering" paths.
func erosionRate2(h *vmesh.Heightmap) *vmesh.Heightmap {
	const maxErosionDistance = 3 * 0.03
	const minExp = 1.0
	const varExp = 3.0

	m := h.Mesh
	steeps := getSteepness(h)
	toE := vmesh.NewHeightmap(m)
	dh := h.Downhill()
	flux := getFlux(h)
	for r, fl := range flux.Values {
		// Exponent (e):
		//
		// The exponent determines the shape of the eroded ravine.
		//
		// y = pow(|x|, e)
		//         _  _
		// e = 1:   \/
		//         _   _
		// e = 2:   \_/
		//         _   _
		// e = 4+:  |_|
		exponent := minExp + varExp*(1.0-steeps.Values[r])

		// The amount of flux determines how wide the resulting eroded path is.
		// Since rivers slow down the wider they are, they start to be less erosive with higher flux.
		// Calculate maximum erosion distance based on the water flow intensity aka flux.
		maxDist := fl * maxErosionDistance // Flux times 3 times the current average neighbor distance (at 16384 points)

		// Assign the erosion intensity.
		if toE.Values[r] < maxDist {
			toE.Values[r] = maxDist
			// TODO: sharp drops should carve with higher intensity.
			// toE.Values[r] *= (1 + steeps.Values[r]) / 2
		}
		rVertex := m.Vertices[r]

		// MEANDERING
		rDh := dh[r]
		var dhVec [2]float64
		// Calculate downhill vector if there is one.
		if rDh >= 0 {
			dhVec = getVecOf2Vxs(rVertex, m.Vertices[rDh])
		}

		seen := make(map[int]bool)
		var doErode func(reg int)
		doErode = func(reg int) {
			if seen[reg] {
				return
			}
			seen[reg] = true
			for _, nb := range h.Neighbours(reg) {
				// Calculate distance.
				dVertex := m.Vertices[nb]

				mul := 1.0
				// MEANDERING
				if nb != rDh && nb != r {
					// If nb is in direction of vec r->rDh
					// extend erosion distance.
					// Calculate downhill neighbor vector.
					dnVec := getVecOf2Vxs(rVertex, dVertex)

					// Calculate how much the vectors differ.
					v := dot2Vx(dhVec, dnVec)
					if v > 0 {
						log.Println(v)
						mul = 1.8
					} else {
						mul = 0.7
					}
				}

				dist := distPoints(rVertex.X, rVertex.Y, dVertex.X, dVertex.Y) / mul
				log.Println("res1", dist)
				if dist > maxDist {
					continue // Skip everything that is too far away.
				}
				distRes := dist / maxDist
				log.Println("nores", distRes)
				erode := math.Pow(distRes, exponent)
				toErode := maxDist - erode
				if toE.Values[nb] < toErode {
					toE.Values[nb] = toErode
				}
				doErode(nb)
			}
		}
		doErode(r)
	}
	return toE
}

func getVecOf2Vxs(rVertex, dVertex voronoi.Vertex) [2]float64 {
	dx := dVertex.X - rVertex.X
	dy := dVertex.Y - rVertex.Y
	c := math.Sqrt(dx*dx + dy*dy)
	return [2]float64{dx / c, dy / c}
}

func getVecOf2Points(i, j int, m *vmesh.Mesh) [2]float64 {
	iVx := m.Vertices[i]
	jVx := m.Vertices[j]
	x1 := iVx.X - jVx.X
	y1 := iVx.Y - jVx.Y
	eee1 := math.Sqrt(x1*x1 + y1*y1)
	return [2]float64{x1 / eee1, y1 / eee1}
}

func dot2(x1, y1, x2, y2 float64) float64 {
	return x1*x2 + y1*y2
}

func dot2Vx(v1, v2 [2]float64) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}*/
