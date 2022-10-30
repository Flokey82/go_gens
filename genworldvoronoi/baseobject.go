package genworldvoronoi

import (
	"log"
	"math"
	"math/rand"

	opensimplex "github.com/ojrac/opensimplex-go"

	"github.com/fogleman/delaunay"
)

type BaseObject struct {
	r_xyz            []float64    // Point / region xyz coordinates
	r_latLon         [][2]float64 // Point / region latitude and longitude
	r_elevation      []float64    // Point / region elevation
	r_moisture       []float64    // Point / region moisture
	r_rainfall       []float64    // Point / region rainfall
	r_flux           []float64    // Point / region hydrology: throughflow of rainfall
	r_pool           []float64
	r_downhill       []int // Point / region mapping to its lowest neighbor
	r_drainage       []int
	r_waterbodies    []int
	r_waterbody_size map[int]int
	r_lake_size      map[int]int
	t_moisture       []float64 // Triangle moisture
	t_elevation      []float64 // Triangle elevation
	t_xyz            []float64 // Triangle xyz coordinates
	t_pool           []float64
	t_latLon         [][2]float64
	seed             int64             // Seed for random number generators
	rand             *rand.Rand        // Rand initialized with above seed
	noise            opensimplex.Noise // Opensimplex noise initialized with above seed
	mesh             *TriangleMesh     // Triangle mesh containing the sphere information
}

// resetRand resets the random number generator to its initial state.
func (m *BaseObject) resetRand() {
	m.rand.Seed(m.seed)
}

// pickRandomRegions picks n random points/regions from the given mesh.
func (m *BaseObject) pickRandomRegions(n int) []int {
	// Reset the random number generator.
	m.resetRand()

	// Pick n random regions.
	chosen_r := make(map[int]bool) // Equivalent of JS new Set()
	for len(chosen_r) < n && len(chosen_r) < m.mesh.numRegions {
		chosen_r[m.rand.Intn(m.mesh.numRegions)] = true
	}

	// Convert map back to a slice (yikes).
	//
	// TODO: Do something more clever and efficient than a map that
	// we convert back to a map anyway.
	return convToArray(chosen_r)
}

// generateTriangleCenters iterates through all triangles and generates the centroids for each.
func (m *BaseObject) generateTriangleCenters() {
	var t_xyz []float64
	var t_latLon [][2]float64
	for t := 0; t < m.mesh.numTriangles; t++ {
		a := m.mesh.s_begin_r(3 * t)
		b := m.mesh.s_begin_r(3*t + 1)
		c := m.mesh.s_begin_r(3*t + 2)
		v3 := getCentroidOfTriangle(
			m.r_xyz[3*a:3*a+3],
			m.r_xyz[3*b:3*b+3],
			m.r_xyz[3*c:3*c+3])
		t_xyz = append(t_xyz, v3.X, v3.Y, v3.Z)
		nla, nlo := latLonFromVec3(v3, 1.0)
		t_latLon = append(t_latLon, [2]float64{nla, nlo})

	}
	m.t_latLon = t_latLon
	m.t_xyz = t_xyz
}

// assignTriangleValues averages out the values of the mesh points / regions and assigns them
// to the triangles of the mesh (or the triangle centroid).
func (m *BaseObject) assignTriangleValues() {
	r_elevation := m.r_elevation
	r_moisture := m.r_moisture
	r_pool := m.r_pool
	t_elevation := m.t_elevation
	t_moisture := m.t_moisture
	t_pool := m.t_pool
	numTriangles := m.mesh.numTriangles

	const tFraction = 1.0 / 3.0
	for t := 0; t < numTriangles; t++ {
		s0 := 3 * t
		r1 := m.mesh.s_begin_r(s0)
		r2 := m.mesh.s_begin_r(s0 + 1)
		r3 := m.mesh.s_begin_r(s0 + 2)
		t_pool[t] = tFraction * (r_pool[r1] + r_pool[r2] + r_pool[r3])
		t_elevation[t] = tFraction * (r_elevation[r1] + r_elevation[r2] + r_elevation[r3])
		t_moisture[t] = tFraction * (r_moisture[r1] + r_moisture[r2] + r_moisture[r3])
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

// assignDownhill will populate r_downhill with a mapping of region to lowest neighbor region.
// NOTE: This is based on mewo2's terrain generation code
// See: https://github.com/mewo2/terrain
func (m *BaseObject) assignDownhill(usePool bool) {
	m.r_downhill = m.getDownhill(usePool)
}

func (m *BaseObject) getDownhill(usePool bool) []int {
	// Here we will map each region to the lowest neighbor.
	r_downhill := make([]int, m.mesh.numRegions)
	for r := range r_downhill {
		lowest_r := -1
		lowest_elevation := m.r_elevation[r]
		if usePool {
			lowest_elevation += m.r_pool[r]
		}
		for _, neighbor_r := range m.rNeighbors(r) {
			elev := m.r_elevation[neighbor_r]
			if usePool {
				elev += m.r_pool[neighbor_r]
			}
			if elev < lowest_elevation {
				lowest_elevation = elev
				lowest_r = neighbor_r
			}
		}
		r_downhill[r] = lowest_r
	}
	return r_downhill
}

// Calculate distance using the lat long and haversine.
func (m *BaseObject) getRDistance(r1, r2 int) float64 {
	la1, lo1 := m.r_latLon[r1][0], m.r_latLon[r1][1]
	la2, lo2 := m.r_latLon[r2][0], m.r_latLon[r2][1]
	return haversine(la1, lo1, la2, lo2)
}

func (m *BaseObject) rNeighbors(r int) []int {
	return m.mesh.r_circulate_r(nil, r)
}

func (m *BaseObject) getLowestNeighbor(r int) int {
	lowest_r := -1
	lowest_elevation := 999.0
	rElev := m.r_elevation[r]
	for _, neighbor_r := range m.rNeighbors(r) {
		elev := m.r_elevation[neighbor_r]
		if elev < lowest_elevation && elev < rElev {
			lowest_elevation = elev
			lowest_r = neighbor_r
		}
	}
	return lowest_r
}

func (m *BaseObject) TestAreas() {
	var tot float64
	for i := 0; i < m.mesh.numRegions; i++ {
		a := m.getRegionArea(i)
		tot += a
		log.Println(a)
	}
	log.Println(tot)
}

// getRegionArea returns the surface area of a region on a unit sphere.
func (m *BaseObject) getRegionArea(r int) float64 {
	rLatLon := m.r_latLon[r]
	ts := m.mesh.r_circulate_t(nil, r)
	dists := make([]float64, len(ts))
	for i, t := range ts {
		dLatLon := m.t_latLon[t]
		dists[i] = haversine(rLatLon[0], rLatLon[1], dLatLon[0], dLatLon[1])
	}
	var area float64
	for ti0, t0 := range ts {
		ti1 := (ti0 + 1) % len(ts)
		t1 := ts[ti1]
		t0LatLon := m.t_latLon[t0]
		t1LatLon := m.t_latLon[t1]
		a := dists[ti0]
		b := dists[ti1]
		c := haversine(t0LatLon[0], t0LatLon[1], t1LatLon[0], t1LatLon[1])
		area += heronsTriArea(a, b, c)
	}
	return area
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

// getRSlope returns the region slope by averaging the slopes of the triangles
// around a given region.
//
// NOTE: This is based on mewo2's erosion code but uses rPolySlope instead of
// rSlope, which determines the slope based on all neighbors.
//
// See: https://github.com/mewo2/terrain
func (m *BaseObject) getRSlope() []float64 {
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

// getRSteepness returns the steepness per region.
//
// NOTE: We define steepness as the angle to a region from its downhill neighbor
// expressed as a value between 0.0 to 1.0 (representing an angle from 0° to 90°).
func (m *BaseObject) getRSteepness() []float64 {
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
func (m *BaseObject) rPolySlope(i int) [2]float64 {
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
func (m *BaseObject) rSlope(i int) [2]float64 {
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
func (m *BaseObject) rTriSlope(nbs []int) [2]float64 {
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

// getSinks returns all regions that do not have a downhill neighbor.
// If 'skipSinksBelowSea' is true, regions below sea level are excluded.
// If 'usePool' is true, water pool data is used to determine if the sink is a lake.
func (m *BaseObject) getSinks(skipSinksBelowSea, usePool bool) []int {
	// Identify sinks above sea level.
	var r_sinks []int
	for r, lowest_r := range m.getDownhill(usePool) {
		if lowest_r == -1 && (!skipSinksBelowSea || m.r_elevation[r] >= 0) { // && m.r_drainage[r] < 0
			r_sinks = append(r_sinks, r)
		}
	}
	return r_sinks
}

// fillSinks is an implementation of the algorithm described in
// https://www.researchgate.net/publication/240407597_A_fast_simple_and_versatile_algorithm_to_fill_the_depressions_of_digital_elevation_models
// and a partial port of the implementation in:
// https://github.com/Rob-Voss/Learninator/blob/master/js/lib/Terrain.js
//
// NOTE: This algorithm produces a too uniform result at the moment, resulting
// in very artificially looking rivers. It lacks some kind of variation like
// noise. It's very fast and less destructive than my other, home-grown algorithm.
// Maybe it's worth to combine the two in some way?
func (m *BaseObject) fillSinks() []float64 {
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

// assignDistanceField calculates the distance from any point in seeds_r to all other points, but
// don't go past any point in stop_r.
func (m *BaseObject) assignDistanceField(seeds_r []int, stop_r map[int]bool) []float64 {
	// Reset the random number generator.
	m.resetRand()

	inf := math.Inf(0)
	mesh := m.mesh
	numRegions := mesh.numRegions

	// Initialize the distance values for all regions to -1.
	r_distance := make([]float64, numRegions)
	for i := range r_distance {
		r_distance[i] = inf
	}

	// Initialize the queue for the breadth first search with
	// the seed regions.
	var queue []int
	for _, r := range seeds_r {
		queue = append(queue, r)
		r_distance[r] = 0
	}

	// Random search adapted from breadth first search.
	var out_r []int

	// TODO: Improve the queue. Currently this is growing unchecked.
	for queue_out := 0; queue_out < len(queue); queue_out++ {
		pos := queue_out + m.rand.Intn(len(queue)-queue_out)
		current_r := queue[pos]
		queue[pos] = queue[queue_out]
		for _, neighbor_r := range mesh.r_circulate_r(out_r, current_r) {
			if !math.IsInf(r_distance[neighbor_r], 0) || stop_r[neighbor_r] {
				continue
			}

			// If the current distance value for neighbor_r is unset (-1)
			// and if neighbor_r is not a "stop region", we set the distance
			// value to the distance value of current_r, incremented by 1.
			r_distance[neighbor_r] = r_distance[current_r] + 1
			queue = append(queue, neighbor_r)
		}
	}

	// TODO: possible enhancement: keep track of which seed is closest
	// to this point, so that we can assign variable mountain/ocean
	// elevation to each seed instead of them always being +1/-1
	return r_distance
}

// assignDistanceFieldWithIntensity is almost identical to assignDistanceField.
// The main difference is that the distance value of each region is reduced by the compression value.
func (m *BaseObject) assignDistanceFieldWithIntensity(seeds_r []int, stop_r map[int]bool, compression map[int]float64) []float64 {
	enableNegativeCompression := true
	enablePositiveCompression := true

	// Reset the random number generator.
	m.resetRand()

	inf := math.Inf(0)
	mesh := m.mesh
	numRegions := mesh.numRegions

	// Initialize the distance values for all regions to -1.
	r_distance := make([]float64, numRegions)
	for i := range r_distance {
		r_distance[i] = inf
	}

	// Initialize the queue for the breadth first search with
	// the seed regions.
	var queue []int
	for _, r := range seeds_r {
		queue = append(queue, r)
		r_distance[r] = 0
	}

	// Get the min and max compression value so that we can
	// normalize the compression value.
	var maxComp, minComp float64
	for _, comp := range compression {
		if comp > maxComp {
			maxComp = comp
		}
		if comp < minComp {
			minComp = comp
		}
	}

	// Random search adapted from breadth first search.
	var out_r []int

	// TODO: Improve the queue. Currently this is growing unchecked.
	for queue_out := 0; queue_out < len(queue); queue_out++ {
		pos := queue_out + m.rand.Intn(len(queue)-queue_out)
		current_r := queue[pos]
		current_comp := compression[current_r]
		current_dist := r_distance[current_r]
		queue[pos] = queue[queue_out]
		for _, neighbor_r := range mesh.r_circulate_r(out_r, current_r) {
			if !math.IsInf(r_distance[neighbor_r], 0) || stop_r[neighbor_r] {
				continue
			}

			// If the current distance value for neighbor_r is unset (-1)
			// and if neighbor_r is not a "stop region", we set the distance
			// value to the distance value of current_r, incremented by 1.
			r_distance[neighbor_r] = current_dist + 1

			// Apply the compression of the current region to the distance
			// value of neighbor_r.
			if current_comp > 0 && enablePositiveCompression {
				// If positive compression is enabled and the compression is... well
				// positive, we subtract the normalized compression value from the
				// distance value for neighbor_r.
				r_distance[neighbor_r] -= current_comp / maxComp
			} else if current_comp < 0 && enableNegativeCompression {
				// If negative compression is enabled and the compression is... well
				// negative, we add the normalized compression value to the distance
				// value for neighbor_r.
				r_distance[neighbor_r] += current_comp / minComp
			}
			// Add neighbor_r to the queue.
			queue = append(queue, neighbor_r)
		}
	}

	// TODO: possible enhancement: keep track of which seed is closest
	// to this point, so that we can assign variable mountain/ocean
	// elevation to each seed instead of them always being +1/-1
	return r_distance
}

type interpolated struct {
	num_r int
	BaseObject
}

// interpolate adds for each neighboring region pair one intermediate,
// interpolated region, increasing the "resolution" for the given regions.
func (m *BaseObject) interpolate(rr []int) (*interpolated, error) {
	// Get all points within bounds.
	var ipl interpolated
	seen := make(map[[2]int]bool)
	for _, r := range rr {
		ipl.num_r++
		rxyz := m.r_xyz[r*3 : (r*3)+3]
		ipl.r_xyz = append(ipl.r_xyz, rxyz...)
		ipl.r_moisture = append(ipl.r_moisture, m.r_moisture[r])
		ipl.r_rainfall = append(ipl.r_rainfall, m.r_rainfall[r])
		ipl.r_flux = append(ipl.r_flux, m.r_flux[r])
		ipl.r_pool = append(ipl.r_pool, m.r_pool[r])
		ipl.r_elevation = append(ipl.r_elevation, m.r_elevation[r])

		// Circulate_r all points and add midpoints.
		for _, rn := range m.mesh.r_circulate_r(nil, r) {
			// Check if we already added a midpoint for this edge.
			var check [2]int
			if r < rn {
				check[0] = r
				check[1] = rn
			} else {
				check[0] = rn
				check[1] = r
			}
			if seen[check] {
				continue
			}
			seen[check] = true

			// Generate midpoint and average values.
			rnxyz := m.r_xyz[rn*3 : (rn*3)+3]
			mid := convToVec3([]float64{
				(rxyz[0] + rnxyz[0]) / 2,
				(rxyz[1] + rnxyz[1]) / 2,
				(rxyz[2] + rnxyz[2]) / 2,
			}).Normalize()
			ipl.r_xyz = append(ipl.r_xyz, mid.X, mid.Y, mid.Z)
			ipl.num_r++

			// Calculate diff and use noise to add variation.
			nvl := (m.noise.Eval3(mid.X, mid.Y, mid.Z) + 1) / 2
			diffElevation := m.r_elevation[rn] - m.r_elevation[r]
			diffMoisture := m.r_moisture[rn] - m.r_moisture[r]
			diffRainfall := m.r_rainfall[rn] - m.r_rainfall[r]
			diffFlux := m.r_flux[rn] - m.r_flux[r]
			diffPool := m.r_pool[rn] - m.r_pool[r]

			// TODO: Add some better variation with the water pool and stuff.
			// TODO: Add flood fill, downhill and flux?

			ipl.r_elevation = append(ipl.r_elevation, m.r_elevation[r]+(diffElevation*nvl))
			ipl.r_moisture = append(ipl.r_moisture, m.r_moisture[r]+(diffMoisture*nvl))
			ipl.r_rainfall = append(ipl.r_rainfall, m.r_rainfall[r]+(diffRainfall*nvl))
			ipl.r_flux = append(ipl.r_flux, m.r_flux[r]+(diffFlux*nvl))
			ipl.r_pool = append(ipl.r_pool, m.r_pool[r]+(diffPool*nvl))
		}
	}
	for r := 0; r < len(ipl.r_xyz); r += 3 {
		// HACKY! Fix this properly!
		nla, nlo := latLonFromVec3(convToVec3(ipl.r_xyz[r:r+3]).Normalize(), 1.0)
		ipl.r_latLon = append(ipl.r_latLon, [2]float64{nla, nlo})
	}
	xy := stereographicProjection(ipl.r_xyz)
	var pts []delaunay.Point
	for i := 0; i < len(xy); i += 2 {
		pts = append(pts, delaunay.Point{X: xy[i], Y: xy[i+1]})
	}

	// Delaunay to rebuild triangles, halfedges etc.
	tri, err := delaunay.Triangulate(pts)
	if err != nil {
		return nil, err
	}
	mesh := NewTriangleMesh(0, len(tri.Triangles), make([]Vertex, ipl.num_r), tri.Triangles, tri.Halfedges)
	ipl.mesh = mesh
	ipl.t_pool = make([]float64, mesh.numTriangles)
	ipl.t_elevation = make([]float64, mesh.numTriangles)
	ipl.t_moisture = make([]float64, mesh.numTriangles)
	ipl.assignDownhill(true)
	ipl.assignTriangleValues()
	ipl.generateTriangleCenters()
	ipl.seed = m.seed
	ipl.rand = rand.New(rand.NewSource(m.seed))
	ipl.noise = m.noise
	return &ipl, nil
}
