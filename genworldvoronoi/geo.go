package genworldvoronoi

import (
	"log"
	"math/rand"
	"time"

	"github.com/Flokey82/go_gens/vectors"
	opensimplex "github.com/ojrac/opensimplex-go"
)

type Geo struct {
	BaseObject
	t_flow               []float64      // Triangle flow intensity (rainfall)
	t_downflow_s         []int          // Triangle mapping to side through which water flows downhill.
	order_t              []int          // Triangles in uphill order of elevation.
	s_flow               []float64      // Flow intensity through sides
	PlateToVector        []vectors.Vec3 // Plate tectonics / movement vectors
	PlateIsOcean         map[int]bool   // Plate was chosen to be an ocean plate
	PlateRegions         []int          // Plate seed points / regions
	RegionToWindVec      []Vertex       // Point / region wind vector
	RegionToWindVecLocal []Vertex       // Point / region wind vector (local)
	RegionToPlate        []int          // Point / region to plate mapping
	RegionToResMetals    []byte         // (resources) Metal ores
	RegionToResGems      []byte         // (resources) Gemstones
	RegionToResStones    []byte         // (resources) Different types of stones or minerals
	NumPlates            int            // Number of generated plates
	NumVolcanoes         int            // Number of generated volcanoes
	NumPoints            int            // Number of generated points / regions
	QuadGeom             *QuadGeometry  // Quad geometry generated from the mesh (?)
}

func newGeo(seed int64, numPlates, numPoints int, jitter float64) (*Geo, error) {
	result, err := MakeSphere(seed, numPoints, jitter)
	if err != nil {
		return nil, err
	}
	mesh := result.mesh

	m := &Geo{
		PlateIsOcean: make(map[int]bool),
		BaseObject: BaseObject{
			XYZ:           result.r_xyz,
			LatLon:        result.r_latLon,
			Elevation:     make([]float64, mesh.numRegions),
			Moisture:      make([]float64, mesh.numRegions),
			Flux:          make([]float64, mesh.numRegions),
			Waterpool:     make([]float64, mesh.numRegions),
			Rainfall:      make([]float64, mesh.numRegions),
			Downhill:      make([]int, mesh.numRegions),
			Drainage:      make([]int, mesh.numRegions),
			t_pool:        make([]float64, mesh.numTriangles),
			t_elevation:   make([]float64, mesh.numTriangles),
			t_moisture:    make([]float64, mesh.numTriangles),
			Waterbodies:   make([]int, mesh.numRegions),
			WaterbodySize: make(map[int]int),
			LakeSize:      make(map[int]int),
			Seed:          seed,
			rand:          rand.New(rand.NewSource(seed)),
			noise:         opensimplex.NewNormalized(seed),
			mesh:          result.mesh,
		},
		t_downflow_s:         make([]int, mesh.numTriangles),
		order_t:              make([]int, mesh.numTriangles),
		t_flow:               make([]float64, mesh.numTriangles),
		s_flow:               make([]float64, mesh.numSides),
		RegionToWindVec:      make([]Vertex, mesh.numRegions),
		RegionToWindVecLocal: make([]Vertex, mesh.numRegions),
		RegionToResMetals:    make([]byte, mesh.numRegions),
		RegionToResGems:      make([]byte, mesh.numRegions),
		RegionToResStones:    make([]byte, mesh.numRegions),
		NumPlates:            numPlates,
		NumVolcanoes:         numPlates, // TODO: Allow independent configuration.
		NumPoints:            numPoints,
		QuadGeom:             NewQuadGeometry(),
	}
	m.QuadGeom.setMesh(mesh)
	m.generateTriangleCenters()
	return m, nil
}

func (m *Geo) generateGeology() {
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
	m.IdentifyLandmasses()
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

// GetCustomContour returns a countour by tracing the region borders determined
// by the function. The result is a list of sequential triangle center points.
// Each sequence represents a continous border around regions of the same type.
//
// The function returns true if the region borders should be traced / if
// the regions do not belong to the same group of regions.
func (m *Geo) GetCustomContour(f func(idxA, idxB int) bool) [][]int {
	var edges [][2]int
	seen := make(map[[2]int]bool)
	for i := 0; i < len(m.mesh.Halfedges); i++ {
		idxA := m.mesh.s_begin_r(i)
		idxB := m.mesh.s_end_r(i)
		var vx [2]int
		if idxA > idxB {
			vx[0] = m.mesh.s_outer_t(i)
			vx[1] = m.mesh.s_inner_t(i)
		} else {
			vx[0] = m.mesh.s_inner_t(i)
			vx[1] = m.mesh.s_outer_t(i)
		}
		if seen[vx] {
			continue
		}
		seen[vx] = true
		if f(idxA, idxB) {
			edges = append(edges, vx)
		}
	}

	return mergeIndexSegments(edges)
}
