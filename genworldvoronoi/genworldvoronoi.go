// Package genworldvoronoi is a port of redblobgames' amazing planet generator.
// See: https://www.redblobgames.com/x/1843-planet-generation
// And: https://github.com/redblobgames/1843-planet-generation
package genworldvoronoi

import (
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
	opensimplex "github.com/ojrac/opensimplex-go"
)

// ugh globals, sorry
type Map struct {
	t_xyz        []float64 // Triangle xyz coordinates
	t_moisture   []float64 // Triangle moisture
	t_elevation  []float64 // Triangle elevation
	t_flow       []float64 // Flow intensity through triangles
	t_downflow_s []int     // Triangle mapping to side through which water flows downhill.
	order_t      []int     // Triangles in uphill order of elevation.
	s_flow       []float64 // Flow intensity through sides
	r_xyz        []float64 // Point / region xyz coordinates
	r_elevation  []float64 // Point / region elevation
	r_moisture   []float64 // Point / region moisture
	r_rainfall   []float64
	r_windvec    []Vertex
	r_flux       []float64         // Hydrology (Variant B): Throughflow of rainfall.
	r_downhill   []int             // Mapping of region to its lowest neighbor.
	r_plate      []int             // Region to plate mapping
	PlateVectors []vectors.Vec3    // Plate tectonics / movement vectors
	PlateIsOcean map[int]bool      // Plate was chosen to be an ocean plate
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
		r_flux:       make([]float64, mesh.numRegions),
		r_downhill:   make([]int, mesh.numRegions),
		t_moisture:   make([]float64, mesh.numTriangles),
		t_downflow_s: make([]int, mesh.numTriangles),
		order_t:      make([]int, mesh.numTriangles),
		t_flow:       make([]float64, mesh.numTriangles),
		s_flow:       make([]float64, mesh.numSides),
		r_rainfall:   make([]float64, mesh.numRegions),
		r_windvec:    make([]Vertex, mesh.numRegions),
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
	// m.assignRegionMoisture()
	for i := 0; i < 12; i++ {
		m.assignRainfall()
		// m.assignFlux()
	}

	// Hydrology (based on regions) - EXPERIMENTAL
	m.assignDownhill()
	m.assignFlux()
	// m.getRivers(9000.1)
	m.r_elevation = m.rErode(0.05)

	// Hydrology (based on triangles)
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

// pushCentroidOfTriangle calculates the centroid of a given triange and appends it to the provided slice.
func pushCentroidOfTriangle(out, a, b, c []float64) []float64 {
	// TODO: renormalize to radius 1.
	// v3 := vectors.Vec3{
	//	X: (a[0]+b[0]+c[0]) / 3,
	//	Y: (a[1]+b[1]+c[1]) / 3,
	//	Z: (a[2]+b[2]+c[2]) / 3,
	// }.Normalize()
	// out = append(out, v3.X, v3.Y, v3.Z)
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

// const Infinity = 1.0

// assignDistanceField calculates the distance from any point in seeds_r to all other points, but
// don't go past any point in stop_r.
func (m *Map) assignDistanceField(seeds_r []int, stop_r map[int]bool) []int64 {
	m.resetRand()
	mesh := m.mesh
	numRegions := mesh.numRegions
	r_distance := make([]int64, numRegions)
	for i := range r_distance {
		r_distance[i] = -1 // was: Infinity
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
		for _, neighbor_r := range mesh.r_circulate_r(out_r, current_r) {
			if r_distance[neighbor_r] == -1 && !stop_r[neighbor_r] {
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

func (m *Map) rNeighbors(r int) []int {
	return m.mesh.r_circulate_r(nil, r)
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
