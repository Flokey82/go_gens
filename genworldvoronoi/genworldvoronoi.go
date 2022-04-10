// Package genworldvoronoi is a port of redblobgames' amazing planet generator.
// See: https://www.redblobgames.com/x/1843-planet-generation
// And: https://github.com/redblobgames/1843-planet-generation
package genworldvoronoi

import (
	"container/heap"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
	opensimplex "github.com/ojrac/opensimplex-go"
)

// ugh globals, sorry
type Map struct {
	t_xyz        []float64         // Triangle xyz coordinates
	t_moisture   []float64         // Triangle moisture
	t_elevation  []float64         // Triangle elevation
	t_flow       []float64         // Flow intensity through triangles
	t_downflow_s []int             // Triangle mapping to side through which water flows downhill.
	r_xyz        []float64         // Point / region xyz coordinates
	r_elevation  []float64         // Point / region elevation
	r_moisture   []float64         // Point / region moisture
	r_plate      []int             // Region to plate mapping
	s_flow       []float64         // Flow intensity through sides
	order_t      []int             // Uphill order of triangles (??)
	PlateVectors []vectors.Vec3    // Plate tectonics / movement vectors
	PlateIsOcean map[int]bool      // Plate was chosed to be an ocean plate
	plate_r      []int             // Plate seed points / regions
	mesh         *TriangleMesh     // Triangle mesh containing the sphere information
	seed         int64             // Seed for random number generators
	rand         *rand.Rand        // Rand initialized with above seed
	noise        opensimplex.Noise // Opensimplex noise initialized with above seed
	NumPlates    int               // Number of generated plates
	NumPoints    int               // Number of generated points / regions
	QuadGeom     *QuadGeometry     // Quad geometry generated from the mesh (?)
}

func NewMap(seed int64, numPlates, numPoints int, jitter float64) (*Map, error) {
	result, err := MakeSphere(seed, numPoints, jitter)
	if err != nil {
		return nil, err
	}
	mesh := result.mesh

	m := &Map{
		PlateIsOcean: make(map[int]bool),
		r_xyz:        result.r_xyz,
		r_elevation:  make([]float64, mesh.numRegions),
		t_elevation:  make([]float64, mesh.numTriangles),
		r_moisture:   make([]float64, mesh.numRegions),
		t_moisture:   make([]float64, mesh.numTriangles),
		t_downflow_s: make([]int, mesh.numTriangles),
		order_t:      make([]int, mesh.numTriangles),
		t_flow:       make([]float64, mesh.numTriangles),
		s_flow:       make([]float64, mesh.numSides),
		mesh:         result.mesh,
		seed:         seed,
		rand:         rand.New(rand.NewSource(seed)),
		noise:        opensimplex.New(seed),
		NumPlates:    numPlates,
		NumPoints:    numPoints,
		QuadGeom:     NewQuadGeometry(),
	}
	m.QuadGeom.setMesh(mesh)
	m.t_xyz = m.generateTriangleCenters()
	m.generateMap()
	return m, nil
}

func (m *Map) resetRand() {
	m.rand.Seed(m.seed)
}

func (m *Map) generateMap() {
	// Plates.
	m.generatePlates()
	m.assignOceanPlates()

	// Elevation.
	m.assignRegionElevation()

	// River / moisture.
	m.assignRegionMoisture()
	m.assignTriangleValues()
	m.assignDownflow()
	m.assignFlow()

	// Quad geometry updete.
	m.QuadGeom.setMap(m.mesh, m)
}

// Plates

// pickRandomRegions picks n random points/regions from the given mesh.
func (m *Map) pickRandomRegions(mesh *TriangleMesh, n int) []int {
	m.resetRand()
	chosen_r := make(map[int]bool) // new Set()
	for len(chosen_r) < n && len(chosen_r) < mesh.numRegions {
		chosen_r[m.rand.Intn(mesh.numRegions)] = true
	}
	return convToArray(chosen_r)
}

// generatePlates generates a number of plate seed points and starts growing the plates
// starting from those seeds in a ranom order.
func (m *Map) generatePlates() {
	m.resetRand()
	mesh := m.mesh
	r_plate := make([]int, mesh.numRegions)
	for i := range r_plate {
		r_plate[i] = -1
	}

	// Pick random regions as seed points for plate generation.
	plate_r := m.pickRandomRegions(mesh, min(m.NumPlates, m.NumPoints))

	var queue []int
	for _, r := range plate_r {
		queue = append(queue, r)
		r_plate[r] = r
	}
	// In Breadth First Search (BFS) the queue will be all elements in
	// queue[queue_out ... queue.length-1]. Pushing onto the queue
	// adds an element to the end, increasing queue.length. Popping
	// from the queue removes an element from the beginning by
	// increasing queue_out.

	// To add variety, use a random search instead of a breadth first
	// search. The frontier of elements to be expanded is still
	// queue[queue_out ... queue.length-1], but pick a random element
	// to pop instead of the earliest one. Do this by swapping
	// queue[pos] and queue[queue_out].
	var out_r []int
	for queue_out := 0; queue_out < len(queue); queue_out++ {
		pos := queue_out + m.rand.Intn(len(queue)-queue_out)
		current_r := queue[pos]
		queue[pos] = queue[queue_out]
		out_r = mesh.r_circulate_r(out_r, current_r)
		for _, neighbor_r := range out_r {
			if r_plate[neighbor_r] == -1 {
				r_plate[neighbor_r] = r_plate[current_r]
				queue = append(queue, neighbor_r)
			}
		}
	}

	// Assign a random movement vector for each plate
	r_xyz := m.r_xyz
	plateVectors := make([]vectors.Vec3, mesh.numRegions)
	for _, center_r := range plate_r {
		neighbor_r := mesh.r_circulate_r(nil, center_r)[0]
		p0 := convToVec3(r_xyz[3*center_r : 3*center_r+3])
		p1 := convToVec3(r_xyz[3*neighbor_r : 3*neighbor_r+3])
		plateVectors[center_r] = vectors.Sub3(p1, p0).Normalize()
	}

	m.plate_r = plate_r
	m.r_plate = r_plate
	m.PlateVectors = plateVectors
}

// assignOceanPlates randomly assigns approx. 50% of the plates as ocean plates.
func (m *Map) assignOceanPlates() {
	m.resetRand()
	m.PlateIsOcean = make(map[int]bool)
	for _, r := range m.plate_r {
		if m.rand.Intn(10) < 5 {
			m.PlateIsOcean[r] = true
			// TODO: either make tiny plates non-ocean, or make sure tiny plates don't create seeds for rivers
		}
	}
	log.Println(m.PlateIsOcean)
	pm := make(map[int]int)
	for _, r := range m.r_plate {
		pm[r]++
	}
	log.Println(pm)
}

// Calculate the collision measure, which is the amount
// that any neighbor's plate vector is pushing against
// the current plate vector.
const collisionThreshold = 0.75

// findCollisions iterates through all regions and finds the regions whose neighbor points
// belong to a different plate. This subset of points is than moved using their respective
// tectonic plate vector and if they approach each other to an extent where they exceed the
// collision threshold, a collision is noted. Depending on the type of plates involved in a
// collision, they produce certain effects like forming a coastline, mountains, etc.
//
// FIXME: The smaller the distance of the cells, the more likely a plate moves past the neighbor plate.
// This causes all kinds of issues.
func (m *Map) findCollisions() ([]int, []int, []int, map[int]float64) {
	plateIsOcean := m.PlateIsOcean
	r_plate := m.r_plate
	plateVectors := m.PlateVectors
	numRegions := m.mesh.numRegions
	compression_r := make(map[int]float64)

	const deltaTime = 1e-7 // simulate movement

	// For each region, I want to know how much it's being compressed
	// into an adjacent region. The "compression" is the change in
	// distance as the two regions move. I'm looking for the adjacent
	// region from a different plate that pushes most into this one
	var mountain_r, coastline_r, ocean_r, r_out []int
	var best_r int
	var bestCompression float64
	for current_r := 0; current_r < numRegions; current_r++ {
		bestCompression = 0.0 // NOTE: Was Infinity
		best_r = -1
		r_out = m.mesh.r_circulate_r(r_out, current_r)
		for _, neighbor_r := range r_out {
			if r_plate[current_r] != r_plate[neighbor_r] {
				// sometimes I regret storing xyz in a compact array...
				current_pos := convToVec3(m.r_xyz[3*current_r : 3*current_r+3])
				neighbor_pos := convToVec3(m.r_xyz[3*neighbor_r : 3*neighbor_r+3])

				// simulate movement for deltaTime seconds
				distanceBefore := vectors.Dist3(current_pos, neighbor_pos)

				plateVec := plateVectors[r_plate[current_r]].Mul(deltaTime)
				a := vectors.Add3(current_pos, plateVec)

				plateVecNeighbor := plateVectors[r_plate[neighbor_r]].Mul(deltaTime)
				b := vectors.Add3(neighbor_pos, plateVecNeighbor)

				distanceAfter := vectors.Dist3(a, b)

				// how much closer did these regions get to each other?
				compression := distanceBefore - distanceAfter

				// keep track of the adjacent region that gets closest.
				if compression > bestCompression { // NOTE: changed from compression < bestCompression
					best_r = neighbor_r
					bestCompression = compression
				}
			}
		}
		// Check if we have a collision candidate.
		if best_r == -1 {
			continue
		}
		compression_r[best_r] += bestCompression

		// at this point, bestCompression tells us how much closer
		// we are getting to the region that's pushing into us the most.
		collided := bestCompression > collisionThreshold*deltaTime
		if plateIsOcean[current_r] && plateIsOcean[best_r] {
			// If both plates are ocean plates and they collide, a coastline is produced,
			// while if they "drift apart" (which is not quite correct in our code, since
			// drifting apart can already be a collision below the threshold), we mark it
			// as "ocean" representing a rift.
			if collided {
				coastline_r = append(coastline_r, current_r)
			} else {
				// In theory, this is not 100% correct, as plates that drift apart result
				// at times in volcanic islands that are formed from escaping magma.
				// See: https://www.icelandontheweb.com/articles-on-iceland/nature/geology/tectonic-plates
				ocean_r = append(ocean_r, current_r)
			}
		} else if !plateIsOcean[current_r] && !plateIsOcean[best_r] {
			// If both plates are non-ocean plates and they collide, mountains are formed.
			if collided {
				mountain_r = append(mountain_r, current_r)
			}
		} else {
			// If the plates are of different types, a collision results in a mountain and
			// drifting apart results in a coastline being defined.
			if collided {
				mountain_r = append(mountain_r, current_r)
			} else {
				coastline_r = append(coastline_r, current_r)
			}
		}
	}
	return mountain_r, coastline_r, ocean_r, compression_r
}

// pushCentroidOfTriangle calculates the centroid of a given triange and appends it to the provided slice.
func pushCentroidOfTriangle(out, a, b, c []float64) []float64 {
	// TODO: renormalize to radius 1
	// v3 := vectors.Vec3{
	//	X: (a[0]+b[0]+c[0]) / 3,
	//	Y: (a[1]+b[1]+c[1]) / 3,
	//	Z: (a[2]+b[2]+c[2]) / 3,
	// }.Normalize()
	//out = append(out, v3.X, v3.Y, v3.Z)
	return append(out, (a[0]+b[0]+c[0])/3, (a[1]+b[1]+c[1])/3, (a[2]+b[2]+c[2])/3)
}

// generateTriangleCenters iterates through all triangles and generates the centroids for each.
func (m *Map) generateTriangleCenters() []float64 {
	var t_xyz []float64
	for t := 0; t < m.mesh.numTriangles; t++ {
		a := m.mesh.s_begin_r(3 * t)
		b := m.mesh.s_begin_r(3*t + 1)
		c := m.mesh.s_begin_r(3*t + 2)
		t_xyz = pushCentroidOfTriangle(t_xyz,
			m.r_xyz[3*a:3*a+3],
			m.r_xyz[3*b:3*b+3],
			m.r_xyz[3*c:3*c+3])
	}
	return t_xyz
}

// assignRegionElevation finds collisions between plate regions and assigns
// elevation for each point on the sphere accordingly, which will result in
// mountains, coastlines, etc.
// To ensure variation, opensimplex noise is used to break up any uniformity.
func (m *Map) assignRegionElevation() {
	const epsilon = 1e-3
	// TODO: Use collision values to determine intensity of generated landscape
	// features.
	mountain_r, coastline_r, ocean_r, _ := m.findCollisions()
	log.Println(mountain_r)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_plate[r] == r {
			if m.PlateIsOcean[r] {
				ocean_r = append(ocean_r, r)
			} else {
				coastline_r = append(coastline_r, r)
			}
		}
	}

	// Distance field generation.
	// I do not quite know how that works, but it is based on:
	// See: https://www.redblobgames.com/x/1728-elevation-control/
	stop_r := make(map[int]bool)
	for _, r := range mountain_r {
		stop_r[r] = true
	}
	for _, r := range coastline_r {
		stop_r[r] = true
	}
	for _, r := range ocean_r {
		stop_r[r] = true
	}

	r_distance_a := m.assignDistanceField(mountain_r, convToMap(ocean_r))
	r_distance_b := m.assignDistanceField(ocean_r, convToMap(coastline_r))
	r_distance_c := m.assignDistanceField(coastline_r, stop_r)

	// Get min/max compression.
	// var compVals []float64
	// for _, v := range compression_r {
	//   compVals = append(compVals, v)
	// }
	// minComp, maxComp := minMax(compVals)

	r_xyz := m.r_xyz
	for r := 0; r < m.mesh.numRegions; r++ {
		a := float64(r_distance_a[r]) + epsilon
		b := float64(r_distance_b[r]) + epsilon
		c := float64(r_distance_c[r]) + epsilon
		if a == Infinity && b == Infinity {
			m.r_elevation[r] = 0.1
		} else {
			m.r_elevation[r] = (1/a - 1/b) / (1/a + 1/b + 1/c)
			// m.r_elevation[r] *= ((compression_r[r] - minComp) / (maxComp - minComp))
		}
		m.r_elevation[r] += m.fbm_noise(r_xyz[3*r], r_xyz[3*r+1], r_xyz[3*r+2])
	}
}

// assignRegionMoisture assigns moisture to each region based on the plate
// it is on. Ideally we would calculate this using prevailing winds and distance
// from the ocean and whatnot.
func (m *Map) assignRegionMoisture() {
	// TODO: assign region moisture in a better way!
	for r := 0; r < m.mesh.numRegions; r++ {
		m.r_moisture[r] = float64(m.r_plate[r]%10) / 10.0
	}
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
			trunk_t := (_halfedges[flow_s] / 3) | 0
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

const Infinity = 1.0

// assignDistanceField calculates the distance from any point in seeds_r to all other points, but
// don't go past any point in stop_r.
func (m *Map) assignDistanceField(seeds_r []int, stop_r map[int]bool) []int64 {
	m.resetRand()
	mesh := m.mesh
	numRegions := mesh.numRegions
	r_distance := make([]int64, numRegions)
	for i := range r_distance {
		r_distance[i] = Infinity
	}

	var queue []int
	for _, r := range seeds_r {
		queue = append(queue, r)
		r_distance[r] = 0
	}

	// Random search adapted from breadth first search.
	var out_r []int
	for queue_out := 0; queue_out < len(queue); queue_out++ {
		pos := queue_out + m.rand.Intn(len(queue)-queue_out)
		current_r := queue[pos]
		queue[pos] = queue[queue_out]
		out_r = mesh.r_circulate_r(out_r, current_r)
		for _, neighbor_r := range out_r {
			if r_distance[neighbor_r] == Infinity && !stop_r[neighbor_r] {
				r_distance[neighbor_r] = r_distance[current_r] + 1
				queue = append(queue, neighbor_r)
			}
		}
	}

	// TODO: possible enhancement: keep track of which seed is closest
	// to this point, so that we can assign variable mountain/ocean
	// elevation to each seed instead of them always being +1/-1
	return r_distance
}

const persistence = 2.0 / 3.0

var amplitudes []float64

func init() {
	amplitudes = make([]float64, 5)
	for i := range amplitudes {
		amplitudes[i] = math.Pow(persistence, float64(i))
	}
}

func (m *Map) fbm_noise(nx, ny, nz float64) float64 {
	sum := 0.0
	sumOfAmplitudes := 0.0
	for octave := 0; octave < len(amplitudes); octave++ {
		frequency := 1 << octave
		sum += amplitudes[octave] * m.noise.Eval3(nx*float64(frequency), ny*float64(frequency), nz*float64(frequency))
		sumOfAmplitudes += amplitudes[octave]
	}
	return sum / sumOfAmplitudes
}

// min is the int equivalent of math.Min(a, b).
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// minMax returns the smallest and largest value in hm.
func minMax(hm []float64) (float64, float64) {
	if len(hm) == 0 {
		return 0, 0
	}
	min := hm[0]
	max := hm[0]
	for _, h := range hm {
		if h > max {
			max = h
		}

		if h < min {
			min = h
		}
	}
	return min, max
}

func convToMap(in []int) map[int]bool {
	res := make(map[int]bool)
	for _, v := range in {
		res[v] = true
	}
	return res
}

func convToArray(in map[int]bool) []int {
	var res []int
	for v := range in {
		res = append(res, v)
	}
	sort.Ints(res)
	return res
}

// convToVec3 converts a float slice containing 3 values into a vectors.Vec3.
func convToVec3(xyz []float64) vectors.Vec3 {
	return vectors.Vec3{xyz[0], xyz[1], xyz[2]}
}
