package genworldvoronoi

import (
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
	for t := 0; t < numTriangles; t++ {
		s0 := 3 * t
		r1 := m.mesh.s_begin_r(s0)
		r2 := m.mesh.s_begin_r(s0 + 1)
		r3 := m.mesh.s_begin_r(s0 + 2)
		t_pool[t] = (1.0 / 3.0) * (r_pool[r1] + r_pool[r2] + r_pool[r3])
		t_elevation[t] = (1.0 / 3.0) * (r_elevation[r1] + r_elevation[r2] + r_elevation[r3])
		t_moisture[t] = (1.0 / 3.0) * (r_moisture[r1] + r_moisture[r2] + r_moisture[r3])
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
		lowest_elevation := 999.0
		rElev := m.r_elevation[r]
		if usePool {
			rElev += m.r_pool[r]
		}
		for _, neighbor_r := range m.rNeighbors(r) {
			elev := m.r_elevation[neighbor_r]
			if usePool {
				elev += m.r_pool[neighbor_r]
			}
			if elev < lowest_elevation && elev < rElev {
				lowest_elevation = elev
				lowest_r = neighbor_r
			}
		}
		r_downhill[r] = lowest_r
	}
	return r_downhill
}

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

func (m *BaseObject) resetRand() {
	m.rand.Seed(m.seed)
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

type interpolated struct {
	num_r int
	BaseObject
}
