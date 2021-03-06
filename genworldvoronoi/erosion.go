package genworldvoronoi

import (
	"log"
	"math"
)

// rErode erodes all region by the given amount.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) rErode(amount float64) []float64 {
	// Get downhill height diffs so we can ensure that we do not erode
	// any more than that, which would produce sinks (which we try to avoid).
	dhDiff := make([]float64, m.mesh.numRegions)
	for r, dhr := range m.getDownhill(false) {
		if dhr < 0 {
			// Sinks have no height diff, but in theory we could give it a negative
			// height diff to fill the sinks during the erosion steps?
			continue
		}
		dhDiff[r] = m.r_elevation[r] - m.r_elevation[dhr]
	}

	// Now apply the erosion.
	newh := make([]float64, m.mesh.numRegions)
	er := m.rErosionRate()
	_, maxr := minMax(er)
	for r, e := range er {
		// We can at most erode amount*dhDiff[r].
		newh[r] = m.r_elevation[r] - amount*dhDiff[r]*(e/maxr)
	}
	return newh
}

func (m *Map) rErodeOld(amount float64) []float64 {
	newh := make([]float64, m.mesh.numRegions)
	er := m.rErosionRate()
	_, maxr := minMax(er)
	for r, e := range er {
		newh[r] = m.r_elevation[r] - amount*(e/maxr)
	}
	return newh
}

// rErosionRate returns the erosion rate per region.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) rErosionRate() []float64 {
	// TODO: Change to distance so that it is independent of resolution.
	const nbErosionFactor = 0.125
	erodeNeighbors := true
	erodeNeighborsRiver := 1
	erodeNeighborsGlacier := 3
	erodeOnlyAboveSealevel := false
	flux := m.getFlux(erodeOnlyAboveSealevel)
	_, maxFlux := minMax(flux)
	slope := m.getRSlope()
	newh := make([]float64, m.mesh.numRegions)

	_, maxH := minMax(m.r_elevation) // TODO: Cache somewhere?
	var erodeRegion func(r, rem int, toErode float64)
	erodeRegion = func(r, rem int, toErode float64) {
		if toErode > newh[r] {
			newh[r] = toErode
		}
		rem--
		// Check erosion depth.
		if rem < 0 || !erodeNeighbors {
			return
		}
		// Additionally erode all neighbors by a certain fraction.
		toErode *= nbErosionFactor
		for _, nb := range m.rNeighbors(r) {
			erodeRegion(nb, rem, toErode)
		}
	}
	log.Println("start erosion")
	for r, s := range slope {
		river := math.Sqrt(flux[r]/maxFlux) * s //flux[r] * slope[r] / maxFlux
		creep := s * s
		total := 1000*river + creep
		if total > 200 {
			total = 200
		}
		// If we have avg. temp. below 0, we need to imitate glacial erosion, which
		// is "wider" than hydraulic erosion.
		erodeNbs := erodeNeighborsRiver
		// TODO: Fix maxElev caching!!!
		if m.getRTemperature(r, maxH) < 0 {
			erodeNbs = erodeNeighborsGlacier
		}
		erodeRegion(r, erodeNbs, total)
	}
	log.Println("stop erosion")
	return newh
}

// getRSlope returns the region slope by averaging the slopes of the triangles
// around a given region.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) getRSlope() []float64 {
	usePolySlope := true // Use neighbor regions for calculating slope.
	slope := make([]float64, m.mesh.numRegions)
	for r, dhr := range m.getDownhill(false) {
		// Sinks have no slope.
		if dhr < 0 {
			continue
		}
		var s [2]float64
		if usePolySlope {
			s = m.rPolySlope(r)
		} else {
			s = m.rSlope(r)
		}
		// The slope is the length of the vector returned by
		// rSlope.
		slope[r] = math.Sqrt(s[0]*s[0] + s[1]*s[1])
	}
	return slope
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

func (m *Map) getRErosion2() []float64 {
	const maxErosionDistance = 3 * 0.006
	const minExp = 1.0
	const varExp = 3.0

	steeps := m.getRSteepness()
	erodeOnlyAboveSealevel := false

	toE := make([]float64, m.mesh.numRegions)
	flux := m.getFlux(erodeOnlyAboveSealevel)
	_, maxFlux := minMax(flux)
	for r, fl := range flux {
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
		exponent := minExp + varExp*(1.0-steeps[r])

		// The amount of flux determines how wide the resulting eroded path is.
		// Since rivers slow down the wider they are, they start to be less erosive with higher flux.
		fluxVal := fl / maxFlux

		// Calculate maximum erosion distance based on the water flow intensity aka flux.
		maxDist := fluxVal * maxErosionDistance

		// Assign the erosion intensity.
		if toE[r] < maxDist {
			toE[r] = maxDist
			// TODO: sharp drops should carve with higher intensity.
		}
		rLatLon := m.r_latLon[r]
		seen := make(map[int]bool)
		var doErode func(reg int)
		doErode = func(reg int) {
			if seen[reg] {
				return
			}
			seen[reg] = true
			for _, nb := range m.rNeighbors(reg) {
				// Calculate great arc distance.
				dLatLon := m.r_latLon[nb]
				dist := haversine(rLatLon[0], rLatLon[1], dLatLon[0], dLatLon[1])
				if dist > maxDist {
					continue // Skip everything that is too far away.
				}
				distRes := dist / maxDist
				erode := math.Pow(distRes, exponent)
				toErode := maxDist - erode
				if toE[nb] < toErode {
					toE[nb] = toErode
				}
				doErode(nb)
			}
		}
		doErode(r)
	}
	return toE
}

func (m *Map) getRSteepness() []float64 {
	steeps := make([]float64, m.mesh.numRegions)
	dh := m.getDownhill(false)
	for r, d := range dh {
		if d < 0 {
			continue // Skip all sinks.
		}
		// Height difference.
		hDiff := m.r_elevation[r] - m.r_elevation[d]

		// Great arc distance.
		rLatLon := m.r_latLon[r]
		dLatLon := m.r_latLon[d]
		dist := haversine(rLatLon[0], rLatLon[1], dLatLon[0], dLatLon[1])

		// Essentially the angle (0??-90??) in the range of 0.0 to 1.0.
		steeps[r] = math.Atan(hDiff/dist) * 2 / math.Pi
	}
	return steeps
}

func (m *Map) rPolySlope(i int) [2]float64 {
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
	/*
		Correct implementation: https://www.khronos.org/opengl/wiki/Calculating_a_Surface_Normal

		Begin Function CalculateSurfaceNormal (Input Polygon) Returns Vector

		   Set Vertex Normal to (0, 0, 0)

		   Begin Cycle for Index in [0, Polygon.vertexNumber)

		      Set Vertex Current to Polygon.verts[Index]
		      Set Vertex Next    to Polygon.verts[(Index plus 1) mod Polygon.vertexNumber]

		      Set Normal.X to Sum of Normal.X and (multiply (Current.Z minus Next.Z) by (Current.Y plus Next.Y))
		      Set Normal.Z to Sum of Normal.Z and (multiply (Current.Y minus Next.Y) by (Current.X plus Next.X))
		      Set Normal.Y to Sum of Normal.Y and (multiply (Current.X minus Next.X) by (Current.Z plus Next.Z))

		   End Cycle

		   Returning Normalize(Normal)

		End Function
	*/
	return [2]float64{normal[0] / -normal[2], normal[1] / -normal[2]} // TODO: Normalize
}

// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) rTriSlope(nbs []int) [2]float64 {
	if len(nbs) != 3 {
		return [2]float64{0, 0}
	}

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

/*
See: https://www.khronos.org/opengl/wiki/Calculating_a_Surface_Normal
Begin Function CalculateSurfaceNormal (Input Triangle) Returns Vector

	Set Vector U to (Triangle.p2 minus Triangle.p1)
	Set Vector V to (Triangle.p3 minus Triangle.p1)

	Set Normal.X to (multiply U.Z by V.Y) minus (multiply U.Y by V.Z)
	Set Normal.Z to (multiply U.Y by V.X) minus (multiply U.X by V.Y)
	Set Normal.Y to (multiply U.X by V.Z) minus (multiply U.Z by V.X)

	Returning Normal

End Function
*/
