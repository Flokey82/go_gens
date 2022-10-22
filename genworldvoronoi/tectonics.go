package genworldvoronoi

import (
	"math"

	"github.com/Flokey82/go_gens/utils"
	"github.com/Flokey82/go_gens/vectors"
)

// generatePlates generates a number of plate seed points and starts growing the plates
// starting from those seeds in a random order.
func (m *Map) generatePlates() {
	m.resetRand()
	mesh := m.mesh
	r_plate := make([]int, mesh.numRegions)
	for i := range r_plate {
		r_plate[i] = -1
	}

	// Pick random regions as seed points for plate generation.
	plate_r := m.pickRandomRegions(mesh, utils.Min(m.NumPlates, m.NumPoints))

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
	nInf := math.Inf(-1)

	const deltaTime = 1e-7 // simulate movement

	// For each region, I want to know how much it's being compressed
	// into an adjacent region. The "compression" is the change in
	// distance as the two regions move. I'm looking for the adjacent
	// region from a different plate that pushes most into this one
	var mountain_r, coastline_r, ocean_r, r_out []int
	var best_r int
	var bestCompression float64
	for current_r := 0; current_r < numRegions; current_r++ {
		bestCompression = nInf // NOTE: Was Infinity
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

		enablePlateCheck := true
		if enablePlateCheck {
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
				} else {
					// coastline_r = append(coastline_r, current_r)
				}
			} else {
				// If the plates are of different types, a collision results in a mountain and
				// drifting apart results in a coastline being defined.
				if collided {
					// If one plate is ocean, mountains only fold up on the non-ocean plate.
					if !plateIsOcean[current_plate] {
						mountain_r = append(mountain_r, current_r)
					}
				} else {
					coastline_r = append(coastline_r, current_r)
				}
			}
		} else {
			// If both plates collide, mountains are formed.
			if collided {
				mountain_r = append(mountain_r, current_r)
			}
		}
	}
	return mountain_r, coastline_r, ocean_r, compression_r
}

const (
	RTypeNone = iota
	RTypeMountain
	RTypeCoastline
	RTypeOcean
)

// assignRegionElevation finds collisions between plate regions and assigns
// elevation for each point on the sphere accordingly, which will result in
// mountains, coastlines, etc.
// To ensure variation, opensimplex noise is used to break up any uniformity.
func (m *Map) assignRegionElevation() {
	useDistanceFieldWithCompression := true

	// TODO: Use collision values to determine intensity of generated landscape features.
	mountain_r, coastline_r, ocean_r, compression_r := m.findCollisions()
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_plate[r] == r && m.PlateIsOcean[r] {
			ocean_r = append(ocean_r, r)
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

	var r_distance_a, r_distance_b, r_distance_c []float64
	if useDistanceFieldWithCompression {
		// Calculate distance fields using the compression values of each region.

		// Graph distance from mountains (stops at ocean regions).
		r_distance_a = m.assignDistanceFieldWithIntensity(mountain_r, convToMap(ocean_r), compression_r)
		// Graph distance from ocean (stops at coastline regions).
		r_distance_b = m.assignDistanceFieldWithIntensity(ocean_r, convToMap(coastline_r), compression_r)
		// Graph distance from coastline (stops at all other regions).
		r_distance_c = m.assignDistanceFieldWithIntensity(coastline_r, stop_r, compression_r)
	} else {
		// Calculate distance fields.

		// Graph distance from mountains (stops at ocean regions).
		r_distance_a = m.assignDistanceField(mountain_r, convToMap(ocean_r))
		// Graph distance from ocean (stops at coastline regions).
		r_distance_b = m.assignDistanceField(ocean_r, convToMap(coastline_r))
		// Graph distance from coastline (stops at all other regions).
		r_distance_c = m.assignDistanceField(coastline_r, stop_r)
	}

	// This code below calculates the height of a given region based on a linear
	// interpolation of the three distance values above.
	//
	// Ideally, we would use some form of noise using the distance to a mountain / faultline
	// to generate a more natural looking landscape with mountain ridges resulting from the
	// folding of the plates.
	//
	// Since we want a "wave" like appearance, we could use one dimensional noise based on the
	// distance to the faultline with some variation for a more natural look.
	const epsilon = 1e-3
	r_xyz := m.r_xyz
	for r := 0; r < m.mesh.numRegions; r++ {
		a := r_distance_a[r] + epsilon // Distance from mountains
		b := r_distance_b[r] + epsilon // Distance from oceans
		c := r_distance_c[r] + epsilon // Distance from coastline
		if m.PlateIsOcean[m.r_plate[r]] {
			// Ocean plates are slightly lower than other plates.
			m.r_elevation[r] = -0.1
		}
		if math.IsInf(r_distance_a[r], 0) && math.IsInf(r_distance_b[r], 0) {
			// If the distance from mountains and oceans is unset (infinity),
			// we increase the elevation by 0.1 since we wouldn't be able to
			// calculate the harmonic mean.
			m.r_elevation[r] += 0.1
		} else {
			// The height is calculated as weighted harmonic mean of the
			// three distance values.
			f := (1/a - 1/b) / (1/a + 1/b + 1/c)
			m.r_elevation[r] += f
		}
		m.r_elevation[r] += m.fbm_noise(r_xyz[3*r], r_xyz[3*r+1], r_xyz[3*r+2])*2 - 1 // Noise from -1.0 to 1.0
	}
}
