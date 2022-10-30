// Package genworldvoronoi is a port of redblobgames' amazing planet generator.
// See: https://www.redblobgames.com/x/1843-planet-generation
// And: https://github.com/redblobgames/1843-planet-generation
package genworldvoronoi

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/Flokey82/go_gens/vectors"
	opensimplex "github.com/ojrac/opensimplex-go"
)

// ugh globals, sorry
type Map struct {
	BaseObject
	t_flow          []float64      // Triangle flow intensity (rainfall)
	t_downflow_s    []int          // Triangle mapping to side through which water flows downhill.
	order_t         []int          // Triangles in uphill order of elevation.
	s_flow          []float64      // Flow intensity through sides
	r_windvec       []Vertex       // Point / region wind vector
	r_windvec_local []Vertex       // Point / region wind vector (local)
	r_plate         []int          // Point / region to plate mapping
	PlateVectors    []vectors.Vec3 // Plate tectonics / movement vectors
	PlateIsOcean    map[int]bool   // Plate was chosen to be an ocean plate
	plate_r         []int          // Plate seed points / regions
	r_territory     []int          // (political) Point / region mapping to territory / empire
	cities_r        []*City        // (political) City seed points / regions
	r_res_metals    []byte         // (resources) Metal ores
	r_res_gems      []byte         // (resources) Gemstones
	r_res_stone     []byte         // (resources) Different types of stones or minerals
	NumPlates       int            // Number of generated plates
	NumPoints       int            // Number of generated points / regions
	NumCities       int            // Number of generated cities (regions)
	NumMiningTowns  int
	// NumFarmingTowns int
	NumTerritories int           // Number of generated territories
	QuadGeom       *QuadGeometry // Quad geometry generated from the mesh (?)
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
			noise:            opensimplex.NewNormalized(seed),
			mesh:             result.mesh,
		},
		t_downflow_s:    make([]int, mesh.numTriangles),
		order_t:         make([]int, mesh.numTriangles),
		t_flow:          make([]float64, mesh.numTriangles),
		s_flow:          make([]float64, mesh.numSides),
		r_windvec:       make([]Vertex, mesh.numRegions),
		r_windvec_local: make([]Vertex, mesh.numRegions),
		r_res_metals:    make([]byte, mesh.numRegions),
		r_res_gems:      make([]byte, mesh.numRegions),
		r_res_stone:     make([]byte, mesh.numRegions),
		NumPlates:       numPlates,
		NumPoints:       numPoints,
		NumTerritories:  10,
		NumCities:       150,
		NumMiningTowns:  60,
		//NumFarmingTowns: 60,
		QuadGeom: NewQuadGeometry(),
	}
	m.QuadGeom.setMesh(mesh)
	m.generateTriangleCenters()
	m.generateMap()
	return m, nil
}

func (m *Map) generateMap() {
	// Generate tectonic plates.
	start := time.Now()
	m.generatePlates()
	m.assignOceanPlates()
	log.Println("Done plates in ", time.Since(start).String())

	// Calculate elevation.
	start = time.Now()
	m.assignRegionElevation()
	log.Println("Done elevation in ", time.Since(start).String())

	// Identify continents / landmasses.
	start = time.Now()
	m.identifyLandmasses()
	log.Println("Done identify landmasses in ", time.Since(start).String())

	// Assign rainfall, moisture.
	start = time.Now()
	m.assignRainfallBasic()
	// m.assignRainfall(1, moistTransferIndirect, moistOrderWind)
	// m.assignFlux()
	log.Println("Done rainfall in ", time.Since(start).String())

	// Hydrology (based on regions) - EXPERIMENTAL
	start = time.Now()
	// m.assignHydrologyWithFlooding()
	m.assignHydrology()
	// m.getRivers(9000.1)
	// m.r_elevation = m.rErode(0.05)
	log.Println("Done hydrology in ", time.Since(start).String())

	// Place resources
	m.placeResources()

	// Place cities and territories in regions.
	start = time.Now()
	m.rPlaceNCities(m.NumCities, TownTypeDefault)
	m.rPlaceNCities(m.NumMiningTowns, TownTypeMining)
	// m.rPlaceNCities(m.NumFarmingTowns, TownTypeFarming)
	log.Println("Done cities in ", time.Since(start).String())

	start = time.Now()
	m.rPlaceNTerritories(m.NumTerritories)
	log.Println("Done territories in ", time.Since(start).String())

	// Once we have established the territories, we can add trade towns
	// (we need the territories for the trade routes).
	// We should probably establish the trade routes now, so we ensure
	// that the trade towns will still be placed on the nexus points
	// where trade routes meet.
	// start = time.Now()
	// m.rPlaceNCities(30, TownTypeTrading)
	// log.Println("Done trade cities in ", time.Since(start).String())

	// Hydrology (based on triangles)
	// Amit's hydrology code.
	start = time.Now()
	m.assignTriangleValues()
	// m.assignDownflow()
	// m.assignFlow()
	log.Println("Done triangles in ", time.Since(start).String())

	// Quad geometry update.
	// This is really only useful for rendering the map but we don't
	// really use this right now.
	start = time.Now()
	m.QuadGeom.setMap(m.mesh, m)
	log.Println("Done quadgeom in ", time.Since(start).String())
}

// pickRandomRegions picks n random points/regions from the given mesh.
func (m *Map) pickRandomRegions(mesh *TriangleMesh, n int) []int {
	// Reset the random number generator.
	m.resetRand()

	// Pick n random regions.
	chosen_r := make(map[int]bool) // Equivalent of JS new Set()
	for len(chosen_r) < n && len(chosen_r) < mesh.numRegions {
		chosen_r[m.rand.Intn(mesh.numRegions)] = true
	}

	// Convert map back to a slice (yikes).
	//
	// TODO: Do something more clever and efficient than a map that
	// we convert back to a map anyway.
	return convToArray(chosen_r)
}

// assignDistanceField calculates the distance from any point in seeds_r to all other points, but
// don't go past any point in stop_r.
func (m *Map) assignDistanceField(seeds_r []int, stop_r map[int]bool) []float64 {
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
func (m *Map) assignDistanceFieldWithIntensity(seeds_r []int, stop_r map[int]bool, compression map[int]float64) []float64 {
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

// Initialize the noise amplitudes for use in our heightmap.
var amplitudes []float64

func init() {
	const persistence = 2.0 / 3.0
	amplitudes = make([]float64, 5)
	for i := range amplitudes {
		amplitudes[i] = math.Pow(persistence, float64(i))
	}
}

// fbm_noise returns a noise value for the given xyz coordinate.
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
