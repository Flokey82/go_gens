package genworldvoronoi

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/vectors"
)

type Geo struct {
	*Calendar
	*BaseObject
	*Resources                          // Natural resources.
	PlateToVector        []vectors.Vec3 // Plate tectonics / movement vectors
	PlateIsOcean         map[int]bool   // Plate was chosen to be an ocean plate
	PlateRegs            []int          // Plate seed points / regions
	RegionToWindVec      []Vertex       // Point / region wind vector
	RegionToWindVecLocal []Vertex       // Point / region wind vector (local)
	RegionToPlate        []int          // Point / region to plate mapping
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
		Calendar:             NewCalendar(),
		PlateIsOcean:         make(map[int]bool),
		BaseObject:           newBaseObject(seed, result),
		Resources:            newResources(mesh.numRegions),
		RegionToWindVec:      make([]Vertex, mesh.numRegions),
		RegionToWindVecLocal: make([]Vertex, mesh.numRegions),
		NumPlates:            numPlates,
		NumVolcanoes:         10, // TODO: Allow independent configuration.
		NumPoints:            numPoints,
		QuadGeom:             NewQuadGeometry(mesh),
	}
	m.generateTriCenters()
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

	// Now that water is assigned, we can make note of waterbodies.
	// NOTE: Lake sizes are assigned in assignHydrology etc.
	start = time.Now()
	m.assignWaterbodies()
	log.Println("Done waterbodies in ", time.Since(start).String())

	// Place resources
	start = time.Now()
	m.placeResources()
	log.Println("Done placing resources in ", time.Since(start).String())

	// Hydrology (based on triangles)
	// Amit's hydrology code.
	start = time.Now()
	m.assignTriValues()
	// m.assignDownflow()
	// m.assignFlow()
	log.Println("Done triangles in ", time.Since(start).String())

	// Quad geometry update.
	// This is really only useful for rendering the map but we don't
	// really use this right now.
	start = time.Now()
	m.QuadGeom.setMap(m.mesh, m)
	log.Println("Done quadgeom in ", time.Since(start).String())

	// Identify continents / landmasses.
	start = time.Now()
	m.assignLandmasses()
	log.Println("Done identify landmasses in ", time.Since(start).String())

	// Update the biome regions.
	// This will be interesting to determine place names, impact on
	// pathfinding (navigating around difficult terrain), etc.
	start = time.Now()
	m.assignBiomeRegions()
	log.Println("Done biome regions in ", time.Since(start).String())
}

func (m *Geo) Tick() {
	// Advance the calendar.
	m.Calendar.Tick()
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
