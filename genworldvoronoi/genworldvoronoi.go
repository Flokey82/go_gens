// Package genworldvoronoi is a port of redblobgames' amazing planet generator.
// See: https://www.redblobgames.com/x/1843-planet-generation
// And: https://github.com/redblobgames/1843-planet-generation
package genworldvoronoi

import (
	"log"
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
	r_city          []int          // (political) Point / region mapping to city / city state
	cities_r        []*City        // (political) City seed points / regions
	r_cultures      []int          // (cultural) Point / region mapping to culture
	cultures_r      []*Culture     // (cultural) Culture seed points / regions
	r_res_metals    []byte         // (resources) Metal ores
	r_res_gems      []byte         // (resources) Gemstones
	r_res_stone     []byte         // (resources) Different types of stones or minerals
	NumPlates       int            // Number of generated plates
	NumVolcanoes    int            // Number of generated volcanoes
	NumPoints       int            // Number of generated points / regions
	NumCities       int            // Number of generated cities (regions)
	NumCityStates   int            // Number of generated city states
	NumMiningTowns  int            // Number of generated mining towns
	NumFarmingTowns int            // Number of generated farming towns
	NumDesertOasis  int            // Number of generated desert oases
	NumTerritories  int            // Number of generated territories
	NumCultures     int            // (Min) Number of generated cultures
	QuadGeom        *QuadGeometry  // Quad geometry generated from the mesh (?)
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
		r_territory:     initRegionSlice(mesh.numRegions),
		NumPlates:       numPlates,
		NumVolcanoes:    numPlates, // TODO: Allow independent configuration.
		NumPoints:       numPoints,
		NumTerritories:  10,
		NumCities:       150,
		NumCityStates:   150,
		NumMiningTowns:  60,
		NumFarmingTowns: 60,
		NumDesertOasis:  10,
		NumCultures:     30,
		QuadGeom:        NewQuadGeometry(),
	}
	m.QuadGeom.setMesh(mesh)
	m.generateTriangleCenters()
	m.generateMap()
	return m, nil
}

func (m *Map) generateMap() {
	// Build geography.

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

	// CIVILIZATION STUFF:
	// Place and grow cultures.
	// Place nature religions.
	// ? Create sub branches of religions.
	// Place cities.
	// Grow city states.
	// Create organized religions.
	// Create empires.

	// Place cultures.
	start = time.Now()
	m.rPlaceNCultures(m.NumCultures)
	log.Println("Done cultures in ", time.Since(start).String())

	// Place cities and territories in regions.
	start = time.Now()
	m.rPlaceNCities(m.NumCities, TownTypeDefault)
	m.rPlaceNCities(m.NumMiningTowns, TownTypeMining)
	m.rPlaceNCities(m.NumFarmingTowns, TownTypeFarming)
	m.rPlaceNCities(m.NumDesertOasis, TownTypeDesertOasis)
	log.Println("Done cities in ", time.Since(start).String())

	start = time.Now()
	m.rPlaceNTerritories(m.NumTerritories)
	log.Println("Done territories in ", time.Since(start).String())

	start = time.Now()
	m.rPlaceNCityStates(m.NumCityStates)
	log.Println("Done city states in ", time.Since(start).String())

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
	//m.GetEmpires()
	m.GetEmpires2()
}
