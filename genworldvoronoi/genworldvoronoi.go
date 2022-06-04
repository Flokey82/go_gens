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
	BaseObject
	t_flow         []float64      // Triangle flow intensity (rainfall)
	t_downflow_s   []int          // Triangle mapping to side through which water flows downhill.
	order_t        []int          // Triangles in uphill order of elevation.
	s_flow         []float64      // Flow intensity through sides
	r_windvec      []Vertex       // Point / region wind vector
	r_plate        []int          // Point / region to plate mapping
	r_territory    []int          // Point / region mapping to territory (political)
	PlateVectors   []vectors.Vec3 // Plate tectonics / movement vectors
	PlateIsOcean   map[int]bool   // Plate was chosen to be an ocean plate
	plate_r        []int          // Plate seed points / regions
	cities_r       []int          // City seed points / regions
	NumPlates      int            // Number of generated plates
	NumPoints      int            // Number of generated points / regions
	NumCities      int            // Number of generated cities (regions)
	NumTerritories int            // Number of generated territories
	QuadGeom       *QuadGeometry  // Quad geometry generated from the mesh (?)
}

func NewMap(seed int64, numPlates, numPoints int, jitter float64) (*Map, error) {
	result, err := MakeSphere(seed, numPoints, jitter)
	if err != nil {
		return nil, err
	}
	mesh := result.mesh

	m := &Map{
		PlateIsOcean: make(map[int]bool),
		BaseObject: BaseObject{
			r_xyz:            result.r_xyz,
			r_latLon:         result.r_latLon,
			r_elevation:      make([]float64, mesh.numRegions),
			r_moisture:       make([]float64, mesh.numRegions),
			r_flux:           make([]float64, mesh.numRegions),
			r_pool:           make([]float64, mesh.numRegions),
			r_rainfall:       make([]float64, mesh.numRegions),
			r_downhill:       make([]int, mesh.numRegions),
			r_drainage:       make([]int, mesh.numRegions),
			t_pool:           make([]float64, mesh.numTriangles),
			t_elevation:      make([]float64, mesh.numTriangles),
			t_moisture:       make([]float64, mesh.numTriangles),
			r_waterbodies:    make([]int, mesh.numRegions),
			r_waterbody_size: make(map[int]int),
			r_lake_size:      make(map[int]int),
			seed:             seed,
			rand:             rand.New(rand.NewSource(seed)),
			noise:            opensimplex.New(seed),
			mesh:             result.mesh,
		},
		t_downflow_s:   make([]int, mesh.numTriangles),
		order_t:        make([]int, mesh.numTriangles),
		t_flow:         make([]float64, mesh.numTriangles),
		s_flow:         make([]float64, mesh.numSides),
		r_windvec:      make([]Vertex, mesh.numRegions),
		NumPlates:      numPlates,
		NumPoints:      numPoints,
		NumTerritories: 10,
		NumCities:      50,
		QuadGeom:       NewQuadGeometry(),
	}
	m.QuadGeom.setMesh(mesh)
	m.generateTriangleCenters()
	m.generateMap()
	return m, nil
}

func (m *Map) generateMap() {
	// Plates.
	m.generatePlates()
	m.assignOceanPlates()

	// Elevation.
	m.assignRegionElevation()

	// River / moisture.
	// m.assignRegionMoisture()
	m.assignRainfall(10)
	// m.assignFlux()

	// Hydrology (based on regions) - EXPERIMENTAL
	m.assignDownhill()
	m.assignFlux()
	//m.makeItRain()
	// m.getRivers(9000.1)
	// m.r_elevation = m.rErode(0.05)

	// Place cities and territories in regions.
	m.rPlaceNCities(m.NumCities)
	m.rPlaceNTerritories(m.NumTerritories)

	// Hydrology (based on triangles)
	m.assignTriangleValues()
	// m.assignDownflow()
	// m.assignFlow()

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
func getCentroidOfTriangle(a, b, c []float64) vectors.Vec3 {
	return vectors.Vec3{
		X: (a[0] + b[0] + c[0]) / 3,
		Y: (a[1] + b[1] + c[1]) / 3,
		Z: (a[2] + b[2] + c[2]) / 3,
	}.Normalize()
}

// const Infinity = 1.0

// assignDistanceField calculates the distance from any point in seeds_r to all other points, but
// don't go past any point in stop_r.
func (m *Map) assignDistanceField(seeds_r []int, stop_r map[int]bool) []float64 {
	m.resetRand()
	mesh := m.mesh
	numRegions := mesh.numRegions
	r_distance := make([]float64, numRegions)
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

func (m *Map) assignDistanceField2(seeds_r []int, stop_r map[int]bool, compression map[int]float64) []float64 {
	enableNegativeCompression := true
	enablePositiveCompression := true

	m.resetRand()
	mesh := m.mesh
	numRegions := mesh.numRegions
	r_distance := make([]float64, numRegions)
	for i := range r_distance {
		r_distance[i] = -1 // was: Infinity
	}

	var queue []int
	for _, r := range seeds_r {
		queue = append(queue, r)
		r_distance[r] = 0
	}

	maxComp := 0.0
	minComp := 0.0
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
	for queue_out := 0; queue_out < len(queue); queue_out++ {
		pos := queue_out + m.rand.Intn(len(queue)-queue_out)
		current_r := queue[pos]
		current_comp := compression[current_r]
		current_dist := r_distance[current_r]
		queue[pos] = queue[queue_out]
		for _, neighbor_r := range mesh.r_circulate_r(out_r, current_r) {
			if r_distance[neighbor_r] == -1 && !stop_r[neighbor_r] {
				r_distance[neighbor_r] = current_dist + 1
				if current_comp > 0 && enablePositiveCompression {
					r_distance[neighbor_r] -= current_comp / maxComp
				} else if current_comp < 0 && enableNegativeCompression {
					r_distance[neighbor_r] += current_comp / minComp
				}
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
