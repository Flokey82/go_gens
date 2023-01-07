package genworldvoronoi

import (
	"container/heap"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/fogleman/delaunay"
)

type BaseObject struct {
	XYZ               []float64       // Point / region xyz coordinates
	LatLon            [][2]float64    // Point / region latitude and longitude
	Elevation         []float64       // Point / region elevation
	Moisture          []float64       // Point / region moisture
	Rainfall          []float64       // Point / region rainfall
	Flux              []float64       // Point / region hydrology: throughflow of rainfall
	Waterpool         []float64       // Point / region hydrology: water pool depth
	Downhill          []int           // Point / region mapping to its lowest neighbor
	Drainage          []int           // Point / region mapping of pool to its drainage region
	Waterbodies       []int           // Point / region mapping of pool to waterbody ID
	WaterbodySize     map[int]int     // Waterbody ID to size mapping
	BiomeRegions      []int           // Point / region mapping of regions with the same biome
	BiomeRegionSize   map[int]int     // Biome region ID to size mapping
	Landmasses        []int           // Point / region mapping of regions that are part of the same landmass
	LandmassSize      map[int]int     // Landmass ID to size mapping
	LakeSize          map[int]int     // Lake ID to size mapping
	RegionIsMountain  map[int]bool    // Point / region is a mountain
	RegionIsVolcano   map[int]bool    // Point / region is a volcano
	RegionCompression map[int]float64 // Point / region compression factor
	triMoisture       []float64       // Triangle moisture
	triElevation      []float64       // Triangle elevation
	tXYZ              []float64       // Triangle xyz coordinates
	triPool           []float64       // Triangle water pool depth
	triLatLon         [][2]float64    // Triangle latitude and longitude
	triFlow           []float64       // Triangle flow intensity (rainfall)
	triDownflowSide   []int           // Triangle mapping to side through which water flows downhill.
	orderTri          []int           // Triangles in uphill order of elevation.
	sideFlow          []float64       // Flow intensity through sides
	Seed              int64           // Seed for random number generators
	rand              *rand.Rand      // Rand initialized with above seed
	noise             *Noise          // Opensimplex noise initialized with above seed
	mesh              *TriangleMesh   // Triangle mesh containing the sphere information
}

func newBaseObject(seed int64, sphere *SphereMesh) *BaseObject {
	mesh := sphere.mesh
	return &BaseObject{
		XYZ:               sphere.xyz,
		LatLon:            sphere.latLon,
		Elevation:         make([]float64, mesh.numRegions),
		Moisture:          make([]float64, mesh.numRegions),
		Flux:              make([]float64, mesh.numRegions),
		Waterpool:         make([]float64, mesh.numRegions),
		Rainfall:          make([]float64, mesh.numRegions),
		Downhill:          make([]int, mesh.numRegions),
		Drainage:          make([]int, mesh.numRegions),
		triPool:           make([]float64, mesh.numTriangles),
		triElevation:      make([]float64, mesh.numTriangles),
		triMoisture:       make([]float64, mesh.numTriangles),
		triDownflowSide:   make([]int, mesh.numTriangles),
		orderTri:          make([]int, mesh.numTriangles),
		triFlow:           make([]float64, mesh.numTriangles),
		sideFlow:          make([]float64, mesh.numSides),
		Waterbodies:       make([]int, mesh.numRegions),
		WaterbodySize:     make(map[int]int),
		BiomeRegions:      make([]int, mesh.numRegions),
		BiomeRegionSize:   make(map[int]int),
		Landmasses:        make([]int, mesh.numRegions),
		LandmassSize:      make(map[int]int),
		LakeSize:          make(map[int]int),
		RegionIsMountain:  make(map[int]bool),
		RegionIsVolcano:   make(map[int]bool),
		RegionCompression: make(map[int]float64),
		Seed:              seed,
		rand:              rand.New(rand.NewSource(seed)),
		noise:             NewNoise(6, 2.0/3.0, seed),
		mesh:              sphere.mesh,
	}
}

// resetRand resets the random number generator to its initial state.
func (m *BaseObject) resetRand() {
	m.rand.Seed(m.Seed)
}

// pickRandomRegions picks n random points/regions from the given mesh.
func (m *BaseObject) pickRandomRegions(n int) []int {
	// Reset the random number generator.
	m.resetRand()

	// Pick n random regions.
	res := make([]int, 0, n)
	for len(res) < n && len(res) < m.mesh.numRegions {
		res = append(res, m.rand.Intn(m.mesh.numRegions))
	}
	sort.Ints(res)
	return res
}

// generateTriCenters iterates through all triangles and generates the centroids for each.
func (m *BaseObject) generateTriCenters() {
	var tXYZ []float64
	var tLatLon [][2]float64
	for t := 0; t < m.mesh.numTriangles; t++ {
		a := m.mesh.s_begin_r(3 * t)
		b := m.mesh.s_begin_r(3*t + 1)
		c := m.mesh.s_begin_r(3*t + 2)
		v3 := getCentroidOfTriangle(
			m.XYZ[3*a:3*a+3],
			m.XYZ[3*b:3*b+3],
			m.XYZ[3*c:3*c+3])
		tXYZ = append(tXYZ, v3.X, v3.Y, v3.Z)
		nla, nlo := latLonFromVec3(v3, 1.0)
		tLatLon = append(tLatLon, [2]float64{nla, nlo})

	}
	m.triLatLon = tLatLon
	m.tXYZ = tXYZ
}

// assignTriValues averages out the values of the mesh points / regions and assigns them
// to the triangles of the mesh (or the triangle centroid).
func (m *BaseObject) assignTriValues() {
	rElevation := m.Elevation
	rMoisture := m.Moisture
	rPool := m.Waterpool
	tElevation := m.triElevation
	tMoisture := m.triMoisture
	tPool := m.triPool
	numTriangles := m.mesh.numTriangles

	const tFraction = 1.0 / 3.0
	for t := 0; t < numTriangles; t++ {
		s0 := 3 * t
		r1 := m.mesh.s_begin_r(s0)
		r2 := m.mesh.s_begin_r(s0 + 1)
		r3 := m.mesh.s_begin_r(s0 + 2)
		tPool[t] = tFraction * (rPool[r1] + rPool[r2] + rPool[r3])
		tElevation[t] = tFraction * (rElevation[r1] + rElevation[r2] + rElevation[r3])
		tMoisture[t] = tFraction * (rMoisture[r1] + rMoisture[r2] + rMoisture[r3])
	}

	// This averages out rainfall to calculate moisture for triangles.
	// Note that this overrides the t_moisture calculated by averaging out r_moisture above.
	for t := 0; t < numTriangles; t++ {
		var moisture float64
		for i := 0; i < 3; i++ {
			s := 3*t + i
			r := m.mesh.s_begin_r(s)
			moisture += m.Rainfall[r] / 3
		}
		tMoisture[t] = moisture
	}
	m.triElevation = tElevation
	m.triPool = tPool
	m.triMoisture = tMoisture
}

// assignDownhill will populate r_downhill with a mapping of region to lowest neighbor region.
// NOTE: This is based on mewo2's terrain generation code
// See: https://github.com/mewo2/terrain
func (m *BaseObject) assignDownhill(usePool bool) {
	m.Downhill = m.GetDownhill(usePool)
}

// GetDownhill will return a mapping of region to lowest neighbor region.
//
// If usePool is true, then the lowest neighbor will be calculated using
// the water pool depth plus the elevation of the region.
func (m *BaseObject) GetDownhill(usePool bool) []int {
	// Here we will map each region to the lowest neighbor.
	rDownhill := make([]int, m.mesh.numRegions)
	outReg := make([]int, 0, 8)
	for r := range rDownhill {
		lowestRegion := -1
		lowestElevation := m.Elevation[r]
		if usePool {
			lowestElevation += m.Waterpool[r]
		}
		for _, nbReg := range m.mesh.r_circulate_r(outReg, r) {
			elev := m.Elevation[nbReg]
			if usePool {
				elev += m.Waterpool[nbReg]
			}
			if elev < lowestElevation {
				lowestElevation = elev
				lowestRegion = nbReg
			}
		}
		rDownhill[r] = lowestRegion
	}
	return rDownhill
}

// assignDownflow starts with triangles that are considered "ocean" and works its way
// uphill to build a graph of child/parents that will allow us later to determine water
// flux and whatnot.
//
// NOTE: This is the original code that Amit uses in his procedural planets project.
// He uses triangle centroids for his river generation, where I prefer to use the regions
// directly.
func (m *BaseObject) assignDownflow() {
	// Use a priority queue, starting with the ocean triangles and
	// moving upwards using elevation as the priority, to visit all
	// the land triangles.
	queue := make(ascPriorityQueue, 0)
	numTriangles := m.mesh.numTriangles
	queueIn := 0
	for i := range m.triDownflowSide {
		m.triDownflowSide[i] = -999
	}
	heap.Init(&queue)

	// Part 1: ocean triangles get downslope assigned to the lowest neighbor.
	for t := 0; t < numTriangles; t++ {
		if m.triElevation[t] < 0 {
			bestSide := -1
			bestElevation := m.triElevation[t]
			for j := 0; j < 3; j++ {
				side := 3*t + j
				elevation := m.triElevation[m.mesh.s_outer_t(side)]
				if elevation < bestElevation {
					bestSide = side
					bestElevation = elevation
				}
			}
			m.orderTri[queueIn] = t
			queueIn++
			m.triDownflowSide[t] = bestSide
			heap.Push(&queue, &queueEntry{
				destination: t,
				score:       m.triElevation[t],
				index:       t,
			})
		}
	}

	// Part 2: land triangles get visited in elevation priority.
	for queueOut := 0; queueOut < numTriangles; queueOut++ {
		current_t := heap.Pop(&queue).(*queueEntry).destination
		for j := 0; j < 3; j++ {
			s := 3*current_t + j
			neighbor_t := m.mesh.s_outer_t(s) // uphill from current_t
			if m.triDownflowSide[neighbor_t] == -999 && m.triElevation[neighbor_t] >= 0.0 {
				m.triDownflowSide[neighbor_t] = m.mesh.s_opposite_s(s)
				m.orderTri[queueIn] = neighbor_t
				queueIn++
				heap.Push(&queue, &queueEntry{
					destination: neighbor_t,
					score:       m.triElevation[neighbor_t],
				})
			}
		}
	}
}

// GetDistance calculate the distance between two regions using
// the lat long and haversine.
func (m *BaseObject) GetDistance(r1, r2 int) float64 {
	la1, lo1 := m.LatLon[r1][0], m.LatLon[r1][1]
	la2, lo2 := m.LatLon[r2][0], m.LatLon[r2][1]
	return haversine(la1, lo1, la2, lo2)
}

// GetRegNeighbors returns the neighbor regions of a region.
func (m *BaseObject) GetRegNeighbors(r int) []int {
	return m.mesh.r_circulate_r(nil, r)
}

func (m *BaseObject) getLowestRegNeighbor(r int) int {
	lowestReg := -1
	lowestElev := 999.0
	rElev := m.Elevation[r]
	for _, nbReg := range m.GetRegNeighbors(r) {
		elev := m.Elevation[nbReg]
		if elev < lowestElev && elev < rElev {
			lowestElev = elev
			lowestReg = nbReg
		}
	}
	return lowestReg
}

// TestAreas essentially sums up the surface area of all the regions
// and prints the total.. which shows that we're pretty close to the
// surface area of a unit sphere. :) Yay!
func (m *BaseObject) TestAreas() {
	var tot float64
	for i := 0; i < m.mesh.numRegions; i++ {
		a := m.GetRegArea(i)
		tot += a
		log.Println(a)
	}
	log.Println(tot)
}

// GetRegArea returns the surface area of a region on a unit sphere.
func (m *BaseObject) GetRegArea(r int) float64 {
	regLatLon := m.LatLon[r]
	tris := m.mesh.r_circulate_t(make([]int, 0, 6), r)
	dists := make([]float64, len(tris))
	for i, tri := range tris {
		dLatLon := m.triLatLon[tri]
		dists[i] = haversine(regLatLon[0], regLatLon[1], dLatLon[0], dLatLon[1])
	}
	var area float64
	for ti0, t0 := range tris {
		ti1 := (ti0 + 1) % len(tris)
		t1 := tris[ti1]
		t0LatLon := m.triLatLon[t0]
		t1LatLon := m.triLatLon[t1]
		a := dists[ti0]
		b := dists[ti1]
		c := haversine(t0LatLon[0], t0LatLon[1], t1LatLon[0], t1LatLon[1])
		area += heronsTriArea(a, b, c)
	}
	return area
}

func (m *BaseObject) isRegBelowOrAtSeaLevelOrPool(r int) bool {
	return m.Elevation[r] <= 0 || m.Waterpool[r] > 0
}

func (m *BaseObject) isRegLakeOrWaterBody(r int) bool {
	return m.isRegWaterBody(r) || m.isRegLake(r)
}

func (m *BaseObject) isRegWaterBody(r int) bool {
	return m.Waterbodies[r] >= 0
}

func (m *BaseObject) isRegLake(r int) bool {
	return m.Drainage[r] >= 0 || m.Waterpool[r] > 0
}

func (m *BaseObject) isRegRiver(r int) bool {
	return m.Flux[r] > m.Rainfall[r]
}

func (m *BaseObject) isRegBigRiver(r int) bool {
	return m.Flux[r] > m.Rainfall[r]*2
}

// GetSlope returns the region slope by averaging the slopes of the triangles
// around a given region.
//
// NOTE: This is based on mewo2's erosion code but uses rPolySlope instead of
// rSlope, which determines the slope based on all neighbors.
//
// See: https://github.com/mewo2/terrain
func (m *BaseObject) GetSlope() []float64 {
	slope := make([]float64, m.mesh.numRegions)
	for r, dhReg := range m.GetDownhill(false) {
		// Sinks have no slope, so we skip them.
		if dhReg < 0 {
			continue
		}

		// Get the slope vector.
		// The slope value we want is the length of the vector returned by rPolySlope.
		// NOTE: We use improved poly-slope code, which uses all neighbors for
		// the slope calculation.
		s := m.regPolySlope(r)
		slope[r] = math.Sqrt(s[0]*s[0] + s[1]*s[1])
	}
	return slope
}

// GetSteepness returns the steepness of every region to their downhill neighbor.
//
// NOTE: We define steepness as the angle to a region from its downhill neighbor
// expressed as a value between 0.0 to 1.0 (representing an angle from 0° to 90°).
func (m *BaseObject) GetSteepness() []float64 {
	// This will collect the steepness for each region.
	steeps := make([]float64, m.mesh.numRegions)

	// Get the downhill neighbors for all regions (ignoring water pools for now).
	dh := m.GetDownhill(false)
	for r, dhReg := range dh {
		if dhReg < 0 {
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
		hDiff := m.Elevation[r] - m.Elevation[dhReg]

		// Great arc distance between the lat/lon coordinates of r and dh[r].
		regLatLon := m.LatLon[r]
		dhRegLatLon := m.LatLon[dhReg]
		dist := haversine(regLatLon[0], regLatLon[1], dhRegLatLon[0], dhRegLatLon[1])

		// Calculate the the angle (0°-90°) expressed as range from 0.0 to 1.0.
		steeps[r] = math.Atan(hDiff/dist) * 2 / math.Pi
	}
	return steeps
}

// regPolySlope calculates the slope of a region, taking in account all neighbors (which form a polygon).
func (m *BaseObject) regPolySlope(i int) [2]float64 {
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

	// Get the origin vector of the center region.
	// We will rotate the points with this vector until the polygon is facing upwards.
	center := convToVec3(m.XYZ[i*3:]).Normalize()

	// Get the axis of rotation.
	axis := center.Cross(vectors.Up)

	// Calculate the angle of rotation.
	angle := math.Acos(vectors.Up.Dot(center) / (vectors.Up.Len() * center.Len()))

	var normal vectors.Vec3
	nbs := m.GetRegNeighbors(i)
	for j, r := range nbs {
		jNext := nbs[(j+1)%len(nbs)]
		// Get the current and next vertex and scale the vector by the height factor
		// and elevation, then rotate the vector around the axis.
		current := convToVec3(m.XYZ[r*3:]).
			Rotate(axis, angle).
			Mul(1 + 0.1*m.Elevation[r])
		next := convToVec3(m.XYZ[jNext*3:]).
			Rotate(axis, angle).
			Mul(1 + 0.1*m.Elevation[jNext])
		normal.X += (current.Z - next.Z) * (current.Y + next.Y)
		normal.Y += (current.Y - next.Y) * (current.X + next.X)
		normal.Z += (current.X - next.X) * (current.Z + next.Z)
	}
	normal = normal.Normalize()
	return [2]float64{normal.X / -normal.Z, normal.Y / -normal.Z} // TODO: Normalize
}

// regSlope returns the x/y vector for a given region by averaging the
// x/y vectors of the neighbor triangle centers.
func (m *BaseObject) regSlope(i int) [2]float64 {
	var res [2]float64
	var count int

	// Buffer for circulating r and t.
	outTri := make([]int, 0, 6)
	outReg := make([]int, 0, 6)

	// NOTE: This is way less accurate. In theory we'd need
	// to calculate the normal of a polygon.
	// See solution rSlope2.
	for _, t := range m.mesh.r_circulate_t(outTri, i) {
		slope := m.regTriSlope(t, m.mesh.t_circulate_r(outReg, t))
		res[0] += slope[0]
		res[1] += slope[1]
		count++
	}
	res[0] /= float64(count)
	res[1] /= float64(count)
	return res
}

// regTriSlope calculates the slope based on three regions.
//
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
//
// WARNING: This only takes in account 3 neighbors!!
// Our implementation however has at times more than 3!
func (m *BaseObject) regTriSlope(t int, nbs []int) [2]float64 {
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

	// Calculate the normal of the triangle.
	normal := m.regTriNormal(t, nbs)

	// Calculate the baricentric coordinates of the triangle center.

	det := normal.Z // negative Z?
	return [2]float64{
		normal.X / det,
		normal.Y / det,
	}
}

func (m *BaseObject) regTriNormal(t int, nbs []int) vectors.Vec3 {

	// Rotate the points so that the triangle is facing upwards.
	// So we calculate the difference between the center vector and the
	// global up vector.
	// Then we rotate the points by the resulting difference vector.
	// This is done by calculating the cross product of the two vectors.
	// The cross product is the axis of rotation and the length of the
	// cross product is the angle of rotation.

	// Get the origin vector of the triangle center.
	// We will rotate the points with this vector until the triangle is facing upwards.
	center := convToVec3(m.tXYZ[t*3:]).Normalize()

	// Get the axis to rotate the 'center' vector to the global up vector.
	axis := center.Cross(vectors.Up)

	// Calculate the angle of rotation.
	angle := math.Acos(vectors.Up.Dot(center) / (vectors.Up.Len() * center.Len()))

	// Get the three points of the triangle.
	p0 := convToVec3(m.XYZ[nbs[0]*3:])
	p1 := convToVec3(m.XYZ[nbs[1]*3:])
	p2 := convToVec3(m.XYZ[nbs[2]*3:])

	p0 = p0.Rotate(axis, angle).Mul(1 + 0.1*m.Elevation[nbs[0]])
	p1 = p1.Rotate(axis, angle).Mul(1 + 0.1*m.Elevation[nbs[1]])
	p2 = p2.Rotate(axis, angle).Mul(1 + 0.1*m.Elevation[nbs[2]])

	// Calculate the normal.
	return p1.Sub(p0).Cross(p2.Sub(p0)).Normalize()
}

// GetSinks returns all regions that do not have a downhill neighbor.
// If 'skipSinksBelowSea' is true, regions below sea level are excluded.
// If 'usePool' is true, water pool data is used to determine if the sink is a lake.
func (m *BaseObject) GetSinks(skipSinksBelowSea, usePool bool) []int {
	// Identify sinks above sea level.
	var regSinks []int
	for r, lowestReg := range m.GetDownhill(usePool) {
		if lowestReg == -1 && (!skipSinksBelowSea || m.Elevation[r] >= 0) { // && m.r_drainage[r] < 0
			regSinks = append(regSinks, r)
		}
	}
	return regSinks
}

// FillSinks is an implementation of the algorithm described in
// https://www.researchgate.net/publication/240407597_A_fast_simple_and_versatile_algorithm_to_fill_the_depressions_of_digital_elevation_models
// and a partial port of the implementation in:
// https://github.com/Rob-Voss/Learninator/blob/master/js/lib/Terrain.js
//
// NOTE: This algorithm produces a too uniform result at the moment, resulting
// in very artificially looking rivers. It lacks some kind of variation like
// noise. It's very fast and less destructive than my other, home-grown algorithm.
// Maybe it's worth to combine the two in some way?
func (m *BaseObject) FillSinks() []float64 {
	// Reset the RNG.
	m.resetRand()

	inf := math.Inf(0)
	baseEpsilon := 1.0 / (float64(m.mesh.numRegions) * 1000.0)
	newHeight := make([]float64, m.mesh.numRegions)
	for i := range newHeight {
		if m.Elevation[i] <= 0 {
			// Set the elevation at or below sea level to the current
			// elevation.
			newHeight[i] = m.Elevation[i]
		} else {
			// Set the elevation above sea level to infinity.
			newHeight[i] = inf
		}
	}

	// Loop until no more changes are made.
	var epsilon float64
	outReg := make([]int, 0, 8)
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
		for _, r := range m.rand.Perm(len(m.Elevation)) {
			// Skip all regions that have the same elevation as in
			// the current heightmap.
			if newHeight[r] == m.Elevation[r] {
				continue
			}

			// Iterate over all neighbors in a random order.
			nbs := m.mesh.r_circulate_r(outReg, r)
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
				if m.Elevation[r] >= newHeight[nb]+epsilon {
					newHeight[r] = m.Elevation[r]
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
				if newHeight[r] > oh && oh > m.Elevation[r] {
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

// assignDistanceField calculates the distance from any point in seedRegs to all other points, but
// don't go past any point in stopReg.
func (m *BaseObject) assignDistanceField(seedRegs []int, stopReg map[int]bool) []float64 {
	// Reset the random number generator.
	m.resetRand()

	inf := math.Inf(0)
	mesh := m.mesh
	numRegions := mesh.numRegions

	// Initialize the distance values for all regions to +Inf.
	regDistance := make([]float64, numRegions)
	for i := range regDistance {
		regDistance[i] = inf
	}

	// Initialize the queue for the breadth first search with
	// the seed regions.
	queue := make([]int, len(seedRegs))
	for i, r := range seedRegs {
		queue[i] = r
		regDistance[r] = 0
	}

	// Allocate a slice for the output of mesh.r_circulate_r.
	outRegs := make([]int, 0, 6)

	// Random search adapted from breadth first search.
	// TODO: Improve the queue. Currently this is growing unchecked.
	for queueOut := 0; queueOut < len(queue); queueOut++ {
		pos := queueOut + m.rand.Intn(len(queue)-queueOut)
		currentReg := queue[pos]
		queue[pos] = queue[queueOut]
		for _, nbReg := range mesh.r_circulate_r(outRegs, currentReg) {
			if !math.IsInf(regDistance[nbReg], 0) || stopReg[nbReg] {
				continue
			}

			// If the current distance value for neighbor_r is unset (-1)
			// and if neighbor_r is not a "stop region", we set the distance
			// value to the distance value of current_r, incremented by 1.
			regDistance[nbReg] = regDistance[currentReg] + 1
			queue = append(queue, nbReg)
		}

		// If we have consumed over 1000000 elements in the queue,
		// we reset the queue to the remaining elements.
		if queueOut > 10000 {
			n := copy(queue, queue[queueOut:])
			queue = queue[:n]
			queueOut = 0
		}
	}

	// TODO: possible enhancement: keep track of which seed is closest
	// to this point, so that we can assign variable mountain/ocean
	// elevation to each seed instead of them always being +1/-1
	return regDistance
}

type interpolated struct {
	numRegions int
	BaseObject
}

// interpolate adds for each neighboring region pair one intermediate,
// interpolated region, increasing the "resolution" for the given regions.
func (m *BaseObject) interpolate(regions []int) (*interpolated, error) {
	var ipl interpolated
	ipl.Seed = m.Seed
	ipl.rand = rand.New(rand.NewSource(m.Seed))

	// Increase the resolution by one octave.
	ipl.noise = m.noise.PlusOneOctave()

	// Get all points within bounds.
	seen := make(map[[2]int]bool)

	// Carry over mountains, volcanoes and compression.
	regionIsMountain := make(map[int]bool)
	regionIsVolcano := make(map[int]bool)
	regionCompression := make(map[int]float64)
	outRegs := make([]int, 0, 6)

	for _, r := range regions {
		if m.RegionIsMountain[r] {
			regionIsMountain[ipl.numRegions] = true
		}
		if m.RegionIsVolcano[r] {
			regionIsVolcano[ipl.numRegions] = true
		}
		if m.RegionCompression[r] != 0 {
			regionCompression[ipl.numRegions] = m.RegionCompression[r]
		}

		ipl.numRegions++
		rxyz := m.XYZ[r*3 : (r*3)+3]
		ipl.XYZ = append(ipl.XYZ, rxyz...)
		ipl.Moisture = append(ipl.Moisture, m.Moisture[r])
		ipl.Rainfall = append(ipl.Rainfall, m.Rainfall[r])
		ipl.Flux = append(ipl.Flux, m.Flux[r])
		ipl.Waterpool = append(ipl.Waterpool, m.Waterpool[r])
		ipl.Elevation = append(ipl.Elevation, m.Elevation[r])

		// Circulate_r all points and add midpoints.
		for _, nbReg := range m.mesh.r_circulate_r(outRegs, r) {
			// Check if we already added a midpoint for this edge.
			var check [2]int
			if r < nbReg {
				check[0] = r
				check[1] = nbReg
			} else {
				check[0] = nbReg
				check[1] = r
			}
			if seen[check] {
				continue
			}
			seen[check] = true

			// Generate midpoint and average values.
			rnxyz := m.XYZ[nbReg*3 : (nbReg*3)+3]
			mid := convToVec3([]float64{
				(rxyz[0] + rnxyz[0]) / 2,
				(rxyz[1] + rnxyz[1]) / 2,
				(rxyz[2] + rnxyz[2]) / 2,
			}).Normalize()
			ipl.XYZ = append(ipl.XYZ, mid.X, mid.Y, mid.Z)
			ipl.numRegions++

			// Calculate diff and use noise to add variation.
			nvl := (ipl.noise.Eval3(mid.X, mid.Y, mid.Z) + 1) / 2
			diffElevation := m.Elevation[nbReg] - m.Elevation[r]
			diffMoisture := m.Moisture[nbReg] - m.Moisture[r]
			diffRainfall := m.Rainfall[nbReg] - m.Rainfall[r]
			diffFlux := m.Flux[nbReg] - m.Flux[r]
			diffPool := m.Waterpool[nbReg] - m.Waterpool[r]

			// TODO: Add some better variation with the water pool and stuff.
			// TODO: Add flood fill, downhill and flux?
			// TODO: Average compression values?

			ipl.Elevation = append(ipl.Elevation, m.Elevation[r]+(diffElevation*nvl))
			ipl.Moisture = append(ipl.Moisture, m.Moisture[r]+(diffMoisture*nvl))
			ipl.Rainfall = append(ipl.Rainfall, m.Rainfall[r]+(diffRainfall*nvl))
			ipl.Flux = append(ipl.Flux, m.Flux[r]+(diffFlux*nvl))
			ipl.Waterpool = append(ipl.Waterpool, m.Waterpool[r]+(diffPool*nvl))
		}
	}
	for r := 0; r < len(ipl.XYZ); r += 3 {
		// HACKY! Fix this properly!
		nla, nlo := latLonFromVec3(convToVec3(ipl.XYZ[r:r+3]).Normalize(), 1.0)
		ipl.LatLon = append(ipl.LatLon, [2]float64{nla, nlo})
	}
	xy := stereographicProjection(ipl.XYZ)
	var pts []delaunay.Point
	for i := 0; i < len(xy); i += 2 {
		pts = append(pts, delaunay.Point{X: xy[i], Y: xy[i+1]})
	}

	// Delaunay to rebuild triangles, halfedges etc.
	tri, err := delaunay.Triangulate(pts)
	if err != nil {
		return nil, err
	}
	mesh := NewTriangleMesh(0, len(tri.Triangles), make([]Vertex, ipl.numRegions), tri.Triangles, tri.Halfedges)
	ipl.mesh = mesh
	ipl.RegionIsMountain = regionIsMountain
	ipl.RegionIsVolcano = regionIsVolcano
	ipl.RegionCompression = regionCompression
	ipl.triPool = make([]float64, mesh.numTriangles)
	ipl.triElevation = make([]float64, mesh.numTriangles)
	ipl.triMoisture = make([]float64, mesh.numTriangles)
	ipl.triDownflowSide = make([]int, mesh.numTriangles)
	ipl.orderTri = make([]int, mesh.numTriangles)
	ipl.triFlow = make([]float64, mesh.numTriangles)
	ipl.sideFlow = make([]float64, mesh.numSides)
	ipl.assignDownhill(true)
	ipl.assignTriValues()
	ipl.generateTriCenters()
	ipl.assignDownflow()
	ipl.assignFlow()
	return &ipl, nil
}
