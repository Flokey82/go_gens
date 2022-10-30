package genworldvoronoi

import (
	"log"
	"math"
)

// rErode erodes all region by the given amount.
//
// NOTE: This is based on mewo2's erosion code but limits the eroded height
// to a fraction of the height difference to the downhill neighbor, which
// prevents extreme sinks from forming.
//
// See: https://github.com/mewo2/terrain
func (m *Map) rErode(amount float64) []float64 {
	// Get downhill height diffs so we can ensure that we do not erode
	// any more than that, which would produce sinks (which we try to avoid).
	dhDiff := make([]float64, m.mesh.numRegions)
	for r, dhr := range m.getDownhill(false) {
		// Skip all sinks which have a downhill value of -1
		if dhr < 0 {
			// NOTE: Sinks have no height diff, but in theory we could give it a
			// negative height diff to fill the sinks during the erosion steps?
			continue
		}
		dhDiff[r] = m.r_elevation[r] - m.r_elevation[dhr]
	}

	// This will hold our new heightmap.
	newh := make([]float64, m.mesh.numRegions)

	// Get the erosion rate for all regions.
	er := m.rErosionRate()

	// Get the maximum erosion rate, so we can normalize the erosion values.
	_, maxr := minMax(er)

	// Calculate the new heightmap by applying the erosion rates we have calculated.
	for r, e := range er {
		// We can at most erode amount*dhDiff[r].
		newh[r] = m.r_elevation[r] - amount*dhDiff[r]*(e/maxr)
	}
	return newh
}

// rErosionRate returns the erosion rate per region.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) rErosionRate() []float64 {
	// TODO: Change to distance so that it is independent of resolution.
	const nbErosionFactor = 0.125

	erodeNeighbors := true          // Traverse neighbors to erode them
	erodeNeighborsRiver := 1        // How deep should we traverse the neighbor-graph? (rivers)
	erodeNeighborsGlacier := 3      // How deep should we traverse the neighbor-graph? (glaciers)
	erodeOnlyAboveSealevel := false // Should we skip erosion below sea level?

	// Get the flux values for all regions.
	flux := m.getFlux(erodeOnlyAboveSealevel)

	// Get max flux so we can normalize the flux values to 0.0 ... 1.0.
	_, maxFlux := minMax(flux)

	// Get the slope values for all regions.
	slope := m.getRSlope()

	// This will hold the erosion values for each region.
	newh := make([]float64, m.mesh.numRegions)

	// Get the max height value so we can normalize the elevation values.
	_, maxH := minMax(m.r_elevation) // TODO: Cache somewhere?

	// erodeRegion sets the erosion rate for the given region and
	// traverses the neighbor graph up to the remaining depth (rem).
	var erodeRegion func(r, rem int, toErode float64)

	erodeRegion = func(r, rem int, toErode float64) {
		// If the given erosion rate is higher than the currently
		// assigned one, override the value.
		//
		// TODO: If we have a higher erosion rate in newh[r], should we use
		// newh[r] instead of toErode to continue?
		if toErode > newh[r] {
			newh[r] = toErode
		}

		// Check if we have exhausted the max traversal depth.
		rem--
		if rem < 0 || !erodeNeighbors {
			return
		}

		// Additionally visit all neighbors and erode them by a certain
		// fraction of the given erosion value.
		toErode *= nbErosionFactor
		for _, nb := range m.rNeighbors(r) {
			erodeRegion(nb, rem, toErode)
		}
	}

	log.Println("start erosion")
	for r, rslope := range slope {
		// NOTE: This was directly taken from mewo2's code.
		//
		// Honestly, I have no idea how this works, especially how the river
		// creep is supposed to work or where these constants come from.
		river := math.Sqrt(flux[r]/maxFlux) * rslope
		creep := rslope * rslope
		total := 1000*river + creep
		if total > 200 {
			total = 200
		}

		// Determine how "deep" we traverse the neighbor-graph based on
		// the type of erosion we want to perform.
		erodeNbs := erodeNeighborsRiver

		// If we have avg. temp. below 0, we need to imitate glacial erosion, which
		// is carving "wider" valleys than hydraulic erosion.
		if m.getRTemperature(r, maxH) < 0 {
			erodeNbs = erodeNeighborsGlacier
		}

		// Traverse the neighbor-graph up to the determined depth and calculate
		// the remaining erosion affecting them given their distance to r.
		erodeRegion(r, erodeNbs, total)
	}
	log.Println("stop erosion")
	return newh
}

// getRSlope returns the region slope by averaging the slopes of the triangles
// around a given region.
//
// NOTE: This is based on mewo2's erosion code but uses rPolySlope instead of
// rSlope, which determines the slope based on all neighbors.
//
// See: https://github.com/mewo2/terrain
func (m *Map) getRSlope() []float64 {
	slope := make([]float64, m.mesh.numRegions)
	for r, dhr := range m.getDownhill(false) {
		// Sinks have no slope, so we skip them.
		if dhr < 0 {
			continue
		}

		// Get the slope vector.
		// The slope value we want is the length of the vector returned by rPolySlope.
		// NOTE: We use improved poly-slope code, which uses all neighbors for
		// the slope calculation.
		s := m.rPolySlope(r)
		slope[r] = math.Sqrt(s[0]*s[0] + s[1]*s[1])
	}
	return slope
}

// rErosionRate2 is an alternative erosion calculation which takes in account
// the steepness and flux of each region to determine the shape of eroded
// riverbeds and valleys.
func (m *Map) rErosionRate2() []float64 {
	const (
		// HACK: That's about 3 neighbors away at 400.000 points. This should not be hardcoded.
		maxErosionDistance = 3 * 0.006
		minExp             = 1.0 // Minimum exponent (see below)
		varExp             = 3.0 // Variable exponent (see further below)
	)

	// Get the steepness of each region to its downhill neighbor.
	steeps := m.getRSteepness()

	// This will collect the erosion values for each region.
	toE := make([]float64, m.mesh.numRegions)

	// If true, we do not consider erosion below sea level.
	erodeOnlyAboveSealevel := false

	// Get the water flux for each region.
	flux := m.getFlux(erodeOnlyAboveSealevel)

	// Get the maximum flux value so we can normalize the flux values
	_, maxFlux := minMax(flux)

	for r, fl := range flux {
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
		exponent := minExp + varExp*(1.0-steeps[r])

		// The amount of flux determines how wide the resulting eroded path is.
		//
		// NOTE: We assume that a higher flux will mean that theoretically the resulting river
		// will become wider due to the additional required space for transporting the larger
		// water volume.
		fluxVal := fl / maxFlux

		// Calculate maximum erosion distance based on the water flow intensity aka flux.
		//
		// NOTE: Since rivers slow down the wider they are, they start to be less erosive with higher flux.
		//
		// WARNING: The maxDist value currently represents both the max width of the riverbed as well as the
		// maximum erosion value at the river center (r). This is not really ideal, but it works.
		maxDist := fluxVal * maxErosionDistance

		// Assign the erosion intensity for (r) if it is larger
		// than the current value set for it.
		if toE[r] < maxDist {
			toE[r] = maxDist // TODO: sharp drops should carve with higher intensity.
		}

		// Get lat/lon of the current region to calculate the great arc
		// distance of each neighbor visited by doErode().
		rLatLon := m.r_latLon[r]

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
			for _, nb := range m.rNeighbors(reg) {
				// Calculate great arc distance to the center of the river (r).
				dLatLon := m.r_latLon[nb]
				dist := haversine(rLatLon[0], rLatLon[1], dLatLon[0], dLatLon[1])

				// Skip everything that is too far away.
				if dist > maxDist {
					continue
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
				if toE[nb] < toErode {
					toE[nb] = toErode
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

// getRSteepness returns the steepness per region.
//
// NOTE: We define steepness as the angle to a region from its downhill neighbor
// expressed as a value between 0.0 to 1.0 (representing an angle from 0° to 90°).
func (m *Map) getRSteepness() []float64 {
	// This will collect the steepness for each region.
	steeps := make([]float64, m.mesh.numRegions)

	// Get the downhill neighbors for all regions (ignoring water pools for now).
	dh := m.getDownhill(false)
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
		hDiff := m.r_elevation[r] - m.r_elevation[d]

		// Great arc distance between the lat/lon coordinates of r and dh[r].
		rLatLon := m.r_latLon[r]
		dLatLon := m.r_latLon[d]
		dist := haversine(rLatLon[0], rLatLon[1], dLatLon[0], dLatLon[1])

		// Calculate the the angle (0°-90°) expressed as range from 0.0 to 1.0.
		steeps[r] = math.Atan(hDiff/dist) * 2 / math.Pi
	}
	return steeps
}

// rPolySlope calculates the slope of a region, taking in account all neighbors (which form a polygon).
func (m *Map) rPolySlope(i int) [2]float64 {
	// See: https://www.khronos.org/opengl/wiki/Calculating_a_Surface_Normal
	//
	// Begin Function CalculateSurfaceNormal (Input Polygon) Returns Vector
	//  Set Vertex Normal to (0, 0, 0)
	//
	//  Begin Cycle for Index in [0, Polygon.vertexNumber)
	//    Set Vertex Current to Polygon.verts[Index]
	//    Set Vertex Next    to Polygon.verts[(Index plus 1) mod Polygon.vertexNumber]
	//
	//    Set Normal.X to Sum of Normal.X and (multiply (Current.Z minus Next.Z) by (Current.Y plus Next.Y))
	//    Set Normal.Z to Sum of Normal.Z and (multiply (Current.Y minus Next.Y) by (Current.X plus Next.X))
	//    Set Normal.Y to Sum of Normal.Y and (multiply (Current.X minus Next.X) by (Current.Z plus Next.Z))
	//  End Cycle
	//
	//  Returning Normalize(Normal)
	// End Function

	var normal [3]float64
	nbs := m.rNeighbors(i)
	for j, r := range nbs {
		jNext := nbs[(j+1)%len(nbs)]
		current := convToVec3(m.r_xyz[r*3:])
		next := convToVec3(m.r_xyz[jNext*3:])
		normal[0] += (current.Z - next.Z) * (current.Y + next.Y)
		normal[1] += (current.Y - next.Y) * (current.X + next.X)
		normal[2] += (current.X - next.X) * (current.Z + next.Z)
	}
	return [2]float64{normal[0] / -normal[2], normal[1] / -normal[2]} // TODO: Normalize
}

// rSlope returns the x/y vector for a given region by averaging the
// x/y vectors of the neighbor triangle centers.
func (m *Map) rSlope(i int) [2]float64 {
	var res [2]float64
	var count int

	// NOTE: This is way less accurate. In theory we'd need
	// to calculate the normal of a polygon.
	// See solution rSlope2.
	for _, t := range m.mesh.r_circulate_t(nil, i) {
		slope := m.rTriSlope(m.mesh.t_circulate_r(nil, t))
		res[0] += slope[0]
		res[1] += slope[1]
		count++
	}
	res[0] /= float64(count)
	res[1] /= float64(count)
	return res
}

// rTriSlope calculates the slope based on three regions.
//
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
//
// WARNING: This only takes in account 3 neighbors!!
// Our implementation however has at times more than 3!
func (m *Map) rTriSlope(nbs []int) [2]float64 {
	// Skip if we don't have enough regions.
	if len(nbs) != 3 {
		return [2]float64{0, 0}
	}

	// I assume that this is what this code is based on...?
	//
	// See: https://www.khronos.org/opengl/wiki/Calculating_a_Surface_Normal
	//
	// Begin Function CalculateSurfaceNormal (Input Triangle) Returns Vector
	//
	//	Set Vector U to (Triangle.p2 minus Triangle.p1)
	//	Set Vector V to (Triangle.p3 minus Triangle.p1)
	//
	//	Set Normal.X to (multiply U.Z by V.Y) minus (multiply U.Y by V.Z)
	//	Set Normal.Z to (multiply U.Y by V.X) minus (multiply U.X by V.Y)
	//	Set Normal.Y to (multiply U.X by V.Z) minus (multiply U.Z by V.X)
	//
	//	Returning Normal
	//
	// End Function

	p0 := convToVec3(m.r_xyz[nbs[0]*3:])
	p1 := convToVec3(m.r_xyz[nbs[1]*3:])
	p2 := convToVec3(m.r_xyz[nbs[2]*3:])

	x1 := p1.X - p0.X
	x2 := p2.X - p0.X
	y1 := p1.Y - p0.Y
	y2 := p2.Y - p0.Y
	z1 := m.r_elevation[nbs[1]] - m.r_elevation[nbs[0]]
	z2 := m.r_elevation[nbs[2]] - m.r_elevation[nbs[0]]

	det := x1*y2 - y1*x2 // negative Z?
	return [2]float64{
		(z1*y2 - y1*z2) / det,
		(x1*z2 - z1*x2) / det,
	}
}
