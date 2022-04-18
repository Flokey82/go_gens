// Package genworldvoronoi is a port of redblobgames' amazing planet generator.
// See: https://www.redblobgames.com/x/1843-planet-generation
// And: https://github.com/redblobgames/1843-planet-generation
package genworldvoronoi

import (
	"container/heap"
	//"log"
	"math"
	"math/rand"
	"sort"

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
	r_xyz        []float64 // Point / region xyz coordinates
	r_elevation  []float64 // Point / region elevation
	r_moisture   []float64 // Point / region moisture
	r_rainfall   []float64
	r_windvec    []Vertex
	r_plate      []int             // Region to plate mapping
	s_flow       []float64         // Flow intensity through sides
	order_t      []int             // Uphill order of triangles (??)
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
	//m.assignRegionMoisture()
	for i := 0; i < 24; i++ {
		m.assignRainfall()
		// m.assignFlux()
	}
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
		bestCompression = -1.0 // NOTE: Was Infinity
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
		current_plate := m.r_plate[current_r]
		best_plate := m.r_plate[best_r]
		if plateIsOcean[current_plate] && plateIsOcean[best_plate] {
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
				// ocean_r = append(ocean_r, current_r)
			}
		} else if !plateIsOcean[current_plate] && !plateIsOcean[best_plate] {
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

// assignRegionElevation finds collisions between plate regions and assigns
// elevation for each point on the sphere accordingly, which will result in
// mountains, coastlines, etc.
// To ensure variation, opensimplex noise is used to break up any uniformity.
func (m *Map) assignRegionElevation() {
	const epsilon = 1e-3
	// TODO: Use collision values to determine intensity of generated landscape features.
	mountain_r, coastline_r, ocean_r, _ := m.findCollisions()
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
		if m.PlateIsOcean[m.r_plate[r]] {
			m.r_elevation[r] = -0.1
		}
		if r_distance_a[r] == -1 && r_distance_b[r] == -1 { // if a == Infinity && b == Infinity {
			m.r_elevation[r] = 0.1
		} else {
			m.r_elevation[r] += (1/a - 1/b) / (1/a + 1/b + 1/c)
			// m.r_elevation[r] *= (compression_r[r] - minComp) / (maxComp - minComp)
		}
		m.r_elevation[r] += m.fbm_noise(r_xyz[3*r], r_xyz[3*r+1], r_xyz[3*r+2])
	}
}

// assignRegionMoisture assigns moisture to each region based on the plate
// it is on. Ideally we would calculate this using prevailing winds and distance
// from the ocean and whatnot.
func (m *Map) assignRegionMoisture() {
	// TODO: assign region moisture in a better way!
	// for r := 0; r < m.mesh.numRegions; r++ {
	//	m.r_moisture[r] = float64(m.r_plate[r]%10) / 10.0
	// }
	var sea_r []int
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_elevation[r] < 0 {
			sea_r = append(sea_r, r)
		}
	}

	// Assign basic moisture per region based on distance to sea.
	// NOTE: This is currently overridden in assignRainfall().
	r_distance_d := m.assignDistanceField(sea_r, make(map[int]bool))
	for r := 0; r < m.mesh.numRegions; r++ {
		m.r_moisture[r] = 1 - math.Min(float64(r_distance_d[r]), 100)/float64(100)
	}
}

func (m *Map) assignWindVectors() {
	r_windvec := make([]Vertex, m.mesh.numRegions)

	// Calculate dummy vector that will be the default for any region we don't have a proper logic for yet.
	// windAngleRad := degToRad(biomesParam.wind_angle_deg)
	// dummyVec := Vertex{math.Cos(windAngleRad), math.Sin(windAngleRad)}
	for i := range r_windvec {
		// Determine latitude of current region.
		rXYZ := convToVec3(m.r_xyz[i*3 : i*3+3])
		rLat, _ := latLonFromVec3(rXYZ.Normalize(), 1.0)
		// Based on latitude, we calculate the wind vector angle.
		var degree float64
		if rLatAbs := math.Abs(rLat); rLatAbs >= 0 && rLatAbs <= 30 {
			// +30° ... 0°, 0° ... -30° -> Primitive Hadley Cell.
			// In a Hadley cell, we turn the wind vector until we are exactly parallel with the equator once we reach 0° Lat.
			// TODO: This is probably not perfectly parallel at the equator.
			if rLat > 0 {
				degree = 180 + 90*(30-rLatAbs)/30 // Northern hemisphere.
			} else {
				degree = 360 - 90*(30-rLatAbs)/30 // Southern hemisphere.
			}
		} else if rLatAbs > 30 && rLatAbs <= 60 {
			// +60° ... +30°, -30° ... -60° -> Primitive Mid Latitude Cell.
			// In a mid latitude cell, we turn the wind vector until we are exactly parallel with the 60° Lat.
			// TODO: This is probably not a full 90° turn. Fix this
			if rLat > 0 {
				degree = 0 + 90*(rLatAbs-30)/30 // Northern hemisphere.
			} else {
				degree = 180 - 90*(rLatAbs-30)/30 // Southern hemisphere.
			}
		} else {
			// NOTE: This is buggy or at least "not nice".
			// +90° ... +60°, -60° ... -90° -> Primitive Hadley Cell.
			// In a polar cell, we turn the wind vector until we are exactly parallel with the equator once we reach 60° Lat.
			if rLat > 0 {
				degree = 180 + 90*(90-rLatAbs)/30 // Northern hemisphere.
			} else {
				degree = 360 - 90*(90-rLatAbs)/30 // Southern hemisphere.
			}
		}
		rad := degToRad(degree)
		r_windvec[i] = Vertex{math.Cos(rad), math.Sin(rad)}
	}

	// TODO: Add wind vectors based on local temperature gradients.
	/*
		r_windvec_local := make([]Vertex, m.mesh.numRegions)
		_, maxElev := minMax(m.r_elevation)
		for r := range r_windvec_local {
			s0 := m.mesh.RInS[r]
			incoming := s0
			lat, _ := latLonFromVec3(convToVec3(m.r_xyz[r*3:(r*3)+3]).Normalize(), 1.0)
			// Get temperature for r.
			temp_r := getMeanAnnualTemp(lat) - getTempFalloffFromAltitude(8850*m.r_elevation[r]/maxElev)
			var count int
			for {
				neighbor_r := m.mesh.s_begin_r(incoming)
				// Get temperature for neighbor_r.
				// if temp_neighbor_r > temp_r
				// 	get direction of gradient.
				// sum up using diff to determine average vector
				count++
				outgoing := s_next_s(incoming)
				incoming = m.mesh.Halfedges[outgoing]
				if incoming == s0 {
					break
				}
			}
		}
	*/

	// Average wind vectors using neighbor vectors.
	interpolationSteps := 5
	for i := 0; i < interpolationSteps; i++ {
		r_windvec_interpol := make([]Vertex, m.mesh.numRegions)
		for r := range r_windvec_interpol {
			s0 := m.mesh.RInS[r]
			incoming := s0
			resVec := Vertex{
				r_windvec[r][0],
				r_windvec[r][1],
			}
			var count int
			for {
				neighbor_r := m.mesh.s_begin_r(incoming)
				resVec[0] += r_windvec[neighbor_r][0]
				resVec[1] += r_windvec[neighbor_r][1]
				count++
				outgoing := s_next_s(incoming)
				incoming = m.mesh.Halfedges[outgoing]
				if incoming == s0 {
					break
				}
			}
			resVec[0] /= float64(count + 1)
			resVec[1] /= float64(count + 1)
			r_windvec_interpol[r] = resVec
		}
		r_windvec = r_windvec_interpol
	}

	m.r_windvec = r_windvec
}

type biomesParams struct {
	raininess   float64 // 0, 2
	rain_shadow float64 // 0.1, 2
	evaporation float64 // 0, 1
}

func (m *Map) assignRainfall() {
	biomesParam := biomesParams{
		raininess:   0.9,
		rain_shadow: 0.5,
		evaporation: 0.9,
	}

	// 1. Assign wind vector for every region
	m.assignWindVectors()
	r_windvec := m.r_windvec

	// 2. Assign initial moisture of 1.0 to all regions below or at sea level
	var sea_r []int
	dist_order_r := make([]int, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		dist_order_r[r] = r
		if m.r_elevation[r] <= 0 {
			sea_r = append(sea_r, r)
			m.r_moisture[r] = 1.0
			m.r_rainfall[r] += biomesParam.raininess * m.r_moisture[r]
		}
	}

	// 3. Sort all regions by distance to ocean. Lowest to highest.
	r_distance_d := m.assignDistanceField(sea_r, make(map[int]bool))
	sort.Slice(dist_order_r, func(a, b int) bool {
		if r_distance_d[dist_order_r[a]] == r_distance_d[dist_order_r[b]] {
			return m.r_elevation[dist_order_r[a]] < m.r_elevation[dist_order_r[b]]
		}
		return r_distance_d[dist_order_r[a]] < r_distance_d[dist_order_r[b]]
	})

	moistureTransferSourceToDest := true
	if moistureTransferSourceToDest {
		// 4. For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r.
		//    This will give us the amount of moisture we transfer to the neighbor region.
		_, maxH := minMax(m.r_elevation)
		for _, r := range dist_order_r {
			count := 0
			s0 := m.mesh.RInS[r]
			incoming := s0
			// Get XYZ Position of r.
			rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
			// Convert to polar coordinates.
			rLat, rLon := latLonFromVec3(rXYZ.Normalize(), 1.0)

			// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
			rwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[r][0], rLon+r_windvec[r][1])).Normalize()
			for {
				neighbor_r := m.mesh.s_begin_r(incoming)
				// Calculate dot product of wind vector to vector r -> neighbor_r.
				// Get XYZ Position of r_neighbor.
				rnXYZ := convToVec3(m.r_xyz[neighbor_r*3 : neighbor_r*3+3])

				// Calculate Vector between r and neighbor_r.
				va := vectors.Sub3(rnXYZ, rXYZ).Normalize()

				// Calculate Vector between r and wind_r.
				vb := vectors.Sub3(rwXYZ, rXYZ).Normalize()

				// Calculate dot product between va and vb.
				// This will give us how much the current region lies within the wind direction of the
				// current neighbor.
				// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
				dotV := vectors.Dot3(va, vb)
				if dotV > 0 {
					// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
					count++
					var humidity, rainfall float64
					humidity = m.r_moisture[neighbor_r] + m.r_moisture[r]*dotV
					rainfall = m.r_rainfall[neighbor_r] + biomesParam.raininess*m.r_moisture[r]*dotV
					heightVal := 1 - (m.r_elevation[neighbor_r] / maxH)
					if humidity > heightVal {
						orographicRainfall := biomesParam.rain_shadow * (humidity - heightVal)
						rainfall += biomesParam.raininess * orographicRainfall
						humidity -= orographicRainfall
					}
					// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
					// WARNING: The humidity calculation is off.
					humidity = math.Min(humidity, 1.0)
					// rainfall = math.Min(rainfall, 1.0)
					m.r_rainfall[neighbor_r] = rainfall
					m.r_moisture[neighbor_r] = humidity
				}

				outgoing := s_next_s(incoming)
				incoming = m.mesh.Halfedges[outgoing]
				if incoming == s0 {
					break
				}
			}
		}
	} else {
		// 4. For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r_neighbor.
		//    This will give us the amount of moisture we transfer from the neighbor region.
		for _, r := range dist_order_r {
			count := 0
			sum := 0.0
			s0 := m.mesh.RInS[r]
			incoming := s0
			// Get XYZ Position of r.
			rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
			for {
				neighbor_r := m.mesh.s_begin_r(incoming)
				// Calculate dot product of wind vector to vector r -> neighbor_r.
				// Get XYZ Position of r_neighbor.
				rnXYZ := convToVec3(m.r_xyz[neighbor_r*3 : neighbor_r*3+3])

				// Convert to polar coordinates.
				rLat, rLon := latLonFromVec3(rnXYZ.Normalize(), 1.0)

				// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
				rnwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[neighbor_r][0], rLon+r_windvec[neighbor_r][1])).Normalize()

				// Calculate Vector between r and neighbor_r.
				va := vectors.Sub3(rXYZ, rnXYZ).Normalize()

				// Calculate Vector between neightbor_r and wind_neighbor_r.
				vb := vectors.Sub3(rnwXYZ, rnXYZ).Normalize()

				// Calculate dot product between va and vb.
				// This will give us how much the current region lies within the wind direction of the
				// current neighbor.
				// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
				dotV := vectors.Dot3(va, vb)
				if dotV > 0 {
					// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
					count++
					sum += m.r_moisture[neighbor_r] * dotV
				}
				outgoing := s_next_s(incoming)
				incoming = m.mesh.Halfedges[outgoing]
				if incoming == s0 {
					break
				}
			}

			var humidity, rainfall float64
			humidity = m.r_moisture[r]
			if count > 0 {
				// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
				// WARNING: The humidity calculation is off.
				humidity = math.Min(humidity+sum, 1.0) // / float64(count)
				rainfall = math.Min(rainfall+biomesParam.raininess*sum, 1.0)
			}
			// if m.mesh.r_boundary(r) {
			//	 humidity = 1.0
			// }
			if m.r_elevation[r] <= 0.0 {
				// evaporation := biomesParam.evaporation * (-m.r_elevation[r])
				// humidity = evaporation
				humidity = m.r_moisture[r]
			}
			if humidity > 1.0-m.r_elevation[r] {
				orographicRainfall := biomesParam.rain_shadow * (humidity - (1.0 - m.r_elevation[r]))
				rainfall += biomesParam.raininess * orographicRainfall
				humidity -= orographicRainfall
			}
			m.r_rainfall[r] = rainfall
			m.r_moisture[r] = humidity

		}
	}

	// Average moisture and rainfall.
	interpolationSteps := 1
	for i := 0; i < interpolationSteps; i++ {
		r_moisture_interpol := make([]float64, m.mesh.numRegions)
		r_rainfall_interpol := make([]float64, m.mesh.numRegions)
		for r := range r_moisture_interpol {
			s0 := m.mesh.RInS[r]
			incoming := s0
			rMoist := m.r_moisture[r]
			rRain := m.r_rainfall[r]
			var count int
			for {
				neighbor_r := m.mesh.s_begin_r(incoming)
				rMoist += m.r_moisture[neighbor_r]
				rRain += m.r_rainfall[neighbor_r]
				count++
				outgoing := s_next_s(incoming)
				incoming = m.mesh.Halfedges[outgoing]
				if incoming == s0 {
					break
				}
			}
			rMoist /= float64(count + 1)
			r_moisture_interpol[r] = rMoist
			rRain /= float64(count + 1)
			r_rainfall_interpol[r] = rRain
		}
		m.r_moisture = r_moisture_interpol
		m.r_rainfall = r_rainfall_interpol
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
			t_flow[t] = 1 / float64(numTriangles) * t_moisture[t] //0.5 * t_moisture[t] * t_moisture[t]
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
			trunk_t := (_halfedges[flow_s] / 3)
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
		out_r = mesh.r_circulate_r(out_r, current_r)
		for _, neighbor_r := range out_r {
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
	min, max := hm[0], hm[0]
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
