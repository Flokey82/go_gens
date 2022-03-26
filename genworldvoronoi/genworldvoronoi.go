// Package genworldvoronoi is a port of redblobgames' amazing planet generator.
// See: https://www.redblobgames.com/x/1843-planet-generation
// And: https://github.com/redblobgames/1843-planet-generation
package genworldvoronoi

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
	opensimplex "github.com/ojrac/opensimplex-go"
)

// ugh globals, sorry
type Map struct {
	t_xyz          []float64
	t_moisture     []float64
	t_elevation    []float64
	t_flow         []float64
	t_downflow_s   []int
	r_xyz          []float64
	r_elevation    []float64
	r_moisture     []float64
	r_plate        []int
	s_flow         []float64
	order_t        []int
	plate_is_ocean map[int]bool
	plate_r        map[int]bool
	mesh           *TriangleMesh
	seed           int64
	rand           *rand.Rand
	noise          opensimplex.Noise
	plate_vec      []vectors.Vec3
	P              int
	N              int
	Q              *QuadGeometry
}

func NewMap(seed int64, P, N int, jitter float64) *Map {
	result := MakeSphere(seed, N, jitter)
	mesh := result.mesh

	m := &Map{
		plate_is_ocean: make(map[int]bool),
		plate_r:        make(map[int]bool),
		r_xyz:          result.r_xyz,
		r_elevation:    make([]float64, mesh.numRegions),
		t_elevation:    make([]float64, mesh.numTriangles),
		r_moisture:     make([]float64, mesh.numRegions),
		t_moisture:     make([]float64, mesh.numTriangles),
		t_downflow_s:   make([]int, mesh.numTriangles),
		order_t:        make([]int, mesh.numTriangles),
		t_flow:         make([]float64, mesh.numTriangles),
		s_flow:         make([]float64, mesh.numSides),
		mesh:           result.mesh,
		seed:           seed,
		rand:           rand.New(rand.NewSource(seed)),
		noise:          opensimplex.New(seed),
		P:              P,
		N:              N,
		Q:              NewQuadGeometry(),
	}
	m.Q.setMesh(mesh)
	m.t_xyz = m.generateTriangleCenters()
	m.generateMap()
	return m
}

func (m *Map) resetRand() {
	m.rand.Seed(m.seed)
}

func (m *Map) ExportOBJ(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	//xy := stereographicProjection(m.r_xyz)
	//for i := 0; i < len(xy); i += 2 {
	//	w.WriteString(fmt.Sprintf("v %f %f %f \n", xy[i], xy[i+1], 2.0)) //
	//}
	for i := 0; i < len(m.r_xyz); i += 3 {
		ve := convToVec3(m.r_xyz[i:]).Mul(1.0 + 0.01*m.r_elevation[i/3])
		w.WriteString(fmt.Sprintf("v %f %f %f \n", ve.X, ve.Y, ve.Z)) //
		//w.WriteString(fmt.Sprintf("v %f %f %f \n", m.r_xyz[i], m.r_xyz[i+1], m.r_xyz[i+2])) //
	}
	w.Flush()
	for i := 0; i < len(m.mesh.Triangles); i += 3 {
		w.WriteString(fmt.Sprintf("f %d %d %d \n", m.mesh.Triangles[i]+1, m.mesh.Triangles[i+1]+1, m.mesh.Triangles[i+2]+1))
		w.Flush()
	}
	w.Flush()
	return nil
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
	m.Q.setMap(m.mesh, m)
}

// Plates

// pickRandomRegions picks n random points/regions.
func (m *Map) pickRandomRegions(mesh *TriangleMesh, n int) map[int]bool {
	m.resetRand()
	chosen_r := make(map[int]bool) // new Set()
	for len(chosen_r) < n && len(chosen_r) < mesh.numRegions {
		chosen_r[m.rand.Intn(mesh.numRegions)] = true
	}
	return chosen_r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *Map) generatePlates() {
	m.resetRand()
	mesh := m.mesh
	r_plate := make([]int, mesh.numRegions)
	for i := range r_plate {
		r_plate[i] = -1
	}

	// Pick random regions as seed points for plate generation.
	plate_r := m.pickRandomRegions(mesh, min(m.P, m.N))

	var queue []int
	for _, r := range convToArray(plate_r) {
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
	plate_vec := make([]vectors.Vec3, mesh.numRegions)
	for center_r := range plate_r {
		neighbor_r := mesh.r_circulate_r(nil, center_r)[0]
		p0 := convToVec3(r_xyz[3*center_r : 3*center_r+3])
		p1 := convToVec3(r_xyz[3*neighbor_r : 3*neighbor_r+3])
		plate_vec[center_r] = vectors.Sub3(p1, p0).Normalize()
	}

	m.plate_r = plate_r
	m.r_plate = r_plate
	m.plate_vec = plate_vec
}

func (m *Map) assignOceanPlates() {
	m.resetRand()
	m.plate_is_ocean = make(map[int]bool)
	for _, r := range convToArray(m.plate_r) {
		if m.rand.Intn(10) < 5 {
			m.plate_is_ocean[r] = true
			// TODO: either make tiny plates non-ocean, or make sure tiny plates don't create seeds for rivers
		}
	}
}

// Calculate the collision measure, which is the amount
// that any neighbor's plate vector is pushing against
// the current plate vector.
const collisionThreshold = 0.75

func (m *Map) findCollisions() ([]int, []int, []int) {
	r_xyz := m.r_xyz
	plate_is_ocean := m.plate_is_ocean
	r_plate := m.r_plate
	plate_vec := m.plate_vec

	const deltaTime = 1e-2 // simulate movement
	numRegions := m.mesh.numRegions
	var mountain_r, coastline_r, ocean_r, r_out []int
	// For each region, I want to know how much it's being compressed
	// into an adjacent region. The "compression" is the change in
	// distance as the two regions move. I'm looking for the adjacent
	// region from a different plate that pushes most into this one
	for current_r := 0; current_r < numRegions; current_r++ {
		bestCompression := Infinity
		best_r := -1
		r_out = m.mesh.r_circulate_r(r_out, current_r)
		for _, neighbor_r := range r_out {
			if r_plate[current_r] != r_plate[neighbor_r] {
				// sometimes I regret storing xyz in a compact array...
				current_pos := r_xyz[3*current_r : 3*current_r+3]
				neighbor_pos := r_xyz[3*neighbor_r : 3*neighbor_r+3]
				// simulate movement for deltaTime seconds
				distanceBefore := vectors.Dist3(convToVec3(current_pos), convToVec3(neighbor_pos))
				a := vectors.Add3(convToVec3(current_pos), plate_vec[r_plate[current_r]].Mul(deltaTime))
				b := vectors.Add3(convToVec3(neighbor_pos), plate_vec[r_plate[neighbor_r]].Mul(deltaTime))
				distanceAfter := vectors.Dist3(a, b)
				// how much closer did these regions get to each other?
				compression := distanceBefore - distanceAfter
				// keep track of the adjacent region that gets closest.
				if compression < bestCompression {
					best_r = neighbor_r
					bestCompression = compression
				}
			}
		}
		if best_r != -1 {
			// at this point, bestCompression tells us how much closer
			// we are getting to the region that's pushing into us the most.
			collided := bestCompression > collisionThreshold*deltaTime
			if plate_is_ocean[current_r] && plate_is_ocean[best_r] {
				if collided {
					coastline_r = append(coastline_r, current_r)
				} else {
					ocean_r = append(ocean_r, current_r)
				}
			} else if !plate_is_ocean[current_r] && !plate_is_ocean[best_r] {
				if collided {
					mountain_r = append(mountain_r, current_r)
				}
			} else {
				if collided {
					mountain_r = append(mountain_r, current_r)
				} else {
					coastline_r = append(coastline_r, current_r)
				}
			}
		}
	}
	return mountain_r, coastline_r, ocean_r
}

// Calculate the centroid and push it onto an array.
func pushCentroidOfTriangle(out []float64, ax, ay, az, bx, by, bz, cx, cy, cz float64) []float64 {
	// TODO: renormalize to radius 1
	out = append(out, (ax+bx+cx)/3, (ay+by+cy)/3, (az+bz+cz)/3)
	return out
}

func (m *Map) generateTriangleCenters() []float64 {
	var t_xyz []float64
	for t := 0; t < m.mesh.numTriangles; t++ {
		a := m.mesh.s_begin_r(3 * t)
		b := m.mesh.s_begin_r(3*t + 1)
		c := m.mesh.s_begin_r(3*t + 2)
		t_xyz = pushCentroidOfTriangle(t_xyz,
			m.r_xyz[3*a], m.r_xyz[3*a+1], m.r_xyz[3*a+2],
			m.r_xyz[3*b], m.r_xyz[3*b+1], m.r_xyz[3*b+2],
			m.r_xyz[3*c], m.r_xyz[3*c+1], m.r_xyz[3*c+2])
	}
	return t_xyz
}

func (m *Map) assignRegionElevation() {
	const epsilon = 1e-3
	mountain_r, coastline_r, ocean_r := m.findCollisions()
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_plate[r] == r {
			if m.plate_is_ocean[r] {
				ocean_r = append(ocean_r, r)
			} else {
				coastline_r = append(coastline_r, r)
			}
		}
	}

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

	r_xyz := m.r_xyz
	for r := 0; r < m.mesh.numRegions; r++ {
		a := r_distance_a[r] + epsilon
		b := r_distance_b[r] + epsilon
		c := r_distance_c[r] + epsilon
		if a == Infinity && b == Infinity {
			m.r_elevation[r] = 0.1
		} else {
			m.r_elevation[r] = (1/a - 1/b) / (1/a + 1/b + 1/c)
		}
		m.r_elevation[r] += m.fbm_noise(r_xyz[3*r], r_xyz[3*r+1], r_xyz[3*r+2])
	}
}

func (m *Map) assignRegionMoisture() {
	// TODO: assign region moisture in a better way!
	for r := 0; r < m.mesh.numRegions; r++ {
		m.r_moisture[r] = float64(m.r_plate[r]%10) / 10.0
	}
}

// Rivers - from mapgen4
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
			_queue.Push(&Item{ID: t, Value: m.t_elevation[t], Index: t})
		}
	}

	// Part 2: land triangles get visited in elevation priority.
	for queue_out := 0; queue_out < numTriangles; queue_out++ {
		current_t := _queue.Pop().(*Item).ID
		for j := 0; j < 3; j++ {
			s := 3*current_t + j
			neighbor_t := m.mesh.s_outer_t(s) // uphill from current_t
			if m.t_downflow_s[neighbor_t] == -999 && m.t_elevation[neighbor_t] >= 0.0 {
				m.t_downflow_s[neighbor_t] = m.mesh.s_opposite_s(s)
				m.order_t[queue_in] = neighbor_t
				queue_in++
				_queue.Push(&Item{ID: neighbor_t, Value: m.t_elevation[neighbor_t]})
				heap.Init(&_queue)
			}
		}
	}
}

func (m *Map) assignFlow() {
	s_flow := m.s_flow
	for i := range s_flow {
		s_flow[i] = 0
	}

	t_flow := m.t_flow
	t_elevation := m.t_elevation
	t_moisture := m.t_moisture
	numTriangles := m.mesh.numTriangles
	for t := 0; t < numTriangles; t++ {
		if t_elevation[t] >= 0.0 {
			t_flow[t] = 0.5 * t_moisture[t] * t_moisture[t]
		} else {
			t_flow[t] = 0
		}
	}

	order_t := m.order_t
	t_downflow_s := m.t_downflow_s
	_halfedges := m.mesh.Halfedges
	for i := len(order_t) - 1; i >= 0; i-- {
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

// Distance from any point in seeds_r to all other points, but
// don't go past any point in stop_r.
func (m *Map) assignDistanceField(seeds_r []int, stop_r map[int]bool) []float64 {
	m.resetRand()
	mesh := m.mesh
	numRegions := mesh.numRegions
	r_distance := make([]float64, numRegions)
	for i := range r_distance {
		r_distance[i] = Infinity
	}

	var queue []int
	for _, r := range seeds_r {
		queue = append(queue, r) //.push(r)
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

				queue = append(queue, neighbor_r) //queue.push(neighbor_r)
			}
		}
	}
	// TODO: possible enhancement: keep track of which seed is closest
	// to this point, so that we can assign variable mountain/ocean
	// elevation to each seed instead of them always being +1/-1
	return r_distance
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

func convToVec3(a []float64) vectors.Vec3 {
	return vectors.Vec3{a[0], a[1], a[2]}
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
