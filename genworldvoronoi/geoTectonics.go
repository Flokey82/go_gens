package genworldvoronoi

import (
	"container/list"
	"math"
	"sort"

	"github.com/Flokey82/go_gens/utils"
	"github.com/Flokey82/go_gens/vectors"
)

// generatePlates generates a number of plate seed points and starts growing the plates
// starting from those seeds in a random order.
func (m *Geo) generatePlates() {
	m.resetRand()
	mesh := m.mesh
	regPlate := make([]int, mesh.numRegions)
	for i := range regPlate {
		regPlate[i] = -1
	}

	// Pick random regions as seed points for plate generation.
	plateRegs := m.pickRandomRegions(utils.Min(m.NumPlates, m.NumPoints))

	var queue []int
	for _, r := range plateRegs {
		queue = append(queue, r)
		regPlate[r] = r
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
	outReg := make([]int, 0, 6)
	for queueOut := 0; queueOut < len(queue); queueOut++ {
		pos := queueOut + m.rand.Intn(len(queue)-queueOut)
		currentReg := queue[pos]
		queue[pos] = queue[queueOut]
		outReg = mesh.r_circulate_r(outReg, currentReg)
		for _, nbReg := range outReg {
			if regPlate[nbReg] == -1 {
				regPlate[nbReg] = regPlate[currentReg]
				queue = append(queue, nbReg)
			}
		}
	}

	// Assign a random movement vector for each plate
	regXYZ := m.XYZ
	plateVectors := make([]vectors.Vec3, mesh.numRegions)
	for _, centerReg := range plateRegs {
		nbReg := mesh.r_circulate_r(outReg, centerReg)[0]
		p0 := convToVec3(regXYZ[3*centerReg : 3*centerReg+3])
		p1 := convToVec3(regXYZ[3*nbReg : 3*nbReg+3])
		plateVectors[centerReg] = vectors.Sub3(p1, p0).Normalize()
	}

	m.PlateRegs = plateRegs
	m.RegionToPlate = regPlate
	m.PlateToVector = plateVectors
}

// assignOceanPlates randomly assigns approx. 50% of the plates as ocean plates.
func (m *Geo) assignOceanPlates() {
	m.resetRand()
	m.PlateIsOcean = make(map[int]bool)
	for _, r := range m.PlateRegs {
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
func (m *Geo) findCollisions() ([]int, []int, []int, map[int]float64) {
	// Use either the largest or smallest compression value.
	useLargestCompression := true

	plateIsOcean := m.PlateIsOcean
	regPlate := m.RegionToPlate
	plateVectors := m.PlateToVector
	numRegions := m.mesh.numRegions
	compressionReg := make(map[int]float64)

	// Initialize the compression measure to either the largest or smallest
	// possible float64 value.
	inf := math.Inf(1)
	if useLargestCompression {
		inf = math.Inf(-1)
	}

	const deltaTime = 1e-11 // simulate movement

	// For each region, I want to know how much it's being compressed
	// into an adjacent region. The "compression" is the change in
	// distance as the two regions move. I'm looking for the adjacent
	// region from a different plate that pushes most into this one
	var mountainRegs, coastlineRegs, oceanRegs []int
	rOut := make([]int, 0, 6)
	var bestReg int
	var bestCompression float64
	for currentReg := 0; currentReg < numRegions; currentReg++ {
		bestCompression = inf
		bestReg = -1
		rOut = m.mesh.r_circulate_r(rOut, currentReg)
		for _, nbReg := range rOut {
			if regPlate[currentReg] != regPlate[nbReg] {
				// sometimes I regret storing xyz in a compact array...
				currentPos := convToVec3(m.XYZ[3*currentReg : 3*currentReg+3])
				neighborPos := convToVec3(m.XYZ[3*nbReg : 3*nbReg+3])

				// simulate movement for deltaTime seconds
				distanceBefore := vectors.Dist3(currentPos, neighborPos)

				plateVec := plateVectors[regPlate[currentReg]].Mul(deltaTime)
				a := vectors.Add3(currentPos, plateVec)

				plateVecNeighbor := plateVectors[regPlate[nbReg]].Mul(deltaTime)
				b := vectors.Add3(neighborPos, plateVecNeighbor)

				distanceAfter := vectors.Dist3(a, b)

				// how much closer did these regions get to each other?
				compression := distanceBefore - distanceAfter

				// Sum up the compression for this region.
				// NOTE: Note sure if this actually makes sense.
				compressionReg[nbReg] += compression

				// keep track of the adjacent region that gets closest.
				// NOTE: changed from compression < bestCompression
				if (compression > bestCompression) == useLargestCompression {
					bestReg = nbReg
					bestCompression = compression
				}
			}
		}

		// Check if we have a ocean region.
		if m.RegionToPlate[currentReg] == currentReg && m.PlateIsOcean[currentReg] {
			oceanRegs = append(oceanRegs, currentReg)
		}

		// Check if we have a collision candidate.
		if bestReg == -1 {
			continue
		}

		compressionReg[currentReg] = bestCompression

		// at this point, bestCompression tells us how much closer
		// we are getting to the region that's pushing into us the most.
		collided := bestCompression > collisionThreshold*deltaTime

		enablePlateCheck := true
		if enablePlateCheck {
			currentPlate := m.RegionToPlate[currentReg]
			bestPlate := m.RegionToPlate[bestReg]
			if plateIsOcean[currentPlate] && plateIsOcean[bestPlate] {
				// If both plates are ocean plates and they collide, a coastline is produced.
				if collided {
					coastlineRegs = append(coastlineRegs, currentReg)
				}
			} else if !plateIsOcean[currentPlate] && !plateIsOcean[bestPlate] {
				// If both plates are non-ocean plates and they collide, mountains are formed.
				if collided {
					mountainRegs = append(mountainRegs, currentReg)
				}
			} else {
				// If the plates are of different types, a collision results in a mountain and
				// drifting apart results in a coastline being defined.
				if collided {
					// If one plate is ocean, mountains only fold up on the non-ocean plate.
					if !plateIsOcean[currentPlate] {
						mountainRegs = append(mountainRegs, currentReg)
					}
				} else {
					// This is incorrect, since can't be certain that we are drifting apart
					// without checking if we have actually a negative compression.
					// I leave this in here, because it just looks cool.
					coastlineRegs = append(coastlineRegs, currentReg)
				}
			}
		} else {
			// If both plates collide, mountains are formed.
			if collided {
				mountainRegs = append(mountainRegs, currentReg)
			}
		}
	}
	return mountainRegs, coastlineRegs, oceanRegs, compressionReg
}

// propagateCompression propagates the compression values from the seed regions
// to all other regions.
func (m *BaseObject) propagateCompression(compression map[int]float64) []float64 {
	// Get the min and max compression value so that we can
	// normalize the compression value, also we need to copy
	// the compression values into a slice so that we can
	// modify them and queue them up.
	cmp := make([]float64, m.mesh.numRegions)
	var cmpSeeds []int
	for r, comp := range compression {
		cmp[r] = comp
		cmpSeeds = append(cmpSeeds, r)
	}
	sort.Ints(cmpSeeds)

	// Queue up the seed regions, shuffle them so that we don't
	// always start with the same regions.
	queue := list.New()
	for _, r := range cmpSeeds {
		queue.PushBack(r)
	}

	// Normalize the compression values.
	minComp, maxComp := minMax(cmp)
	for r := range cmp {
		if cmp[r] > 0 {
			cmp[r] /= maxComp
		} else if cmp[r] < 0 {
			cmp[r] /= math.Abs(minComp)
		}
	}

	// Propagate the compression values.
	outRegs := make([]int, 0, 6)
	for queue.Len() > 0 {
		currentReg := queue.Remove(queue.Front()).(int)
		currentComp := cmp[currentReg]
		for _, nbReg := range m.mesh.r_circulate_r(outRegs, currentReg) {
			// The compression value diminishes over distance.
			// This should be using the inverse square law, but
			// we use a linear function instead.
			distToNb := 1 + m.GetDistance(currentReg, nbReg)
			nbComp := currentComp / distToNb
			if cmp[nbReg] == 0 {
				cmp[nbReg] = nbComp
				queue.PushBack(nbReg)
			} else {
				// Average the compression values.
				// NOTE: I know this is not great, but it works.
				// ... otherwise we get real bad artifacts.
				cmp[nbReg] = (cmp[nbReg] + nbComp) / 2
			}
		}
	}

	// Normalize the compression values.
	minComp, maxComp = minMax(cmp)
	for r := range cmp {
		if cmp[r] > 0 {
			cmp[r] /= maxComp
		} else if cmp[r] < 0 {
			cmp[r] /= math.Abs(minComp)
		}
		// Apply a square falloff to the compression values.
		cmp[r] *= math.Abs(cmp[r])
	}
	return cmp
}

// assignRegionElevation finds collisions between plate regions and assigns
// elevation for each point on the sphere accordingly, which will result in
// mountains, coastlines, etc.
// To ensure variation, opensimplex noise is used to break up any uniformity.
func (m *Geo) assignRegionElevation() {
	// TODO: Use collision values to determine intensity of generated landscape features.
	mountainRegs, coastlineRegs, oceanRegs, compressionReg := m.findCollisions()

	// Sort mountains by compression.
	sort.Slice(mountainRegs, func(i, j int) bool {
		return compressionReg[mountainRegs[i]] > compressionReg[mountainRegs[j]]
	})

	// Take note of all mountains.
	// Since they are sorted by compression, we can use the first m.NumVolcanoes
	// as volcanoes.
	var gotVolcanoes int
	for _, r := range mountainRegs {
		m.RegionIsMountain[r] = true
		if gotVolcanoes < m.NumVolcanoes {
			m.RegionIsVolcano[r] = true
			gotVolcanoes++
		}
	}

	// Take note of the compression of each region.
	m.RegionCompression = compressionReg

	// Distance field generation.
	// I do not quite know how that works, but it is based on:
	// See: https://www.redblobgames.com/x/1728-elevation-control/
	stopReg := make(map[int]bool)
	for _, r := range mountainRegs {
		stopReg[r] = true
	}
	for _, r := range coastlineRegs {
		stopReg[r] = true
	}
	for _, r := range oceanRegs {
		stopReg[r] = true
	}

	// Calculate distance fields.
	// Graph distance from mountains (stops at ocean regions).
	rDistanceA := m.assignDistanceField(mountainRegs, convToMap(oceanRegs))
	// Graph distance from ocean (stops at coastline regions).
	rDistanceB := m.assignDistanceField(oceanRegs, convToMap(coastlineRegs))
	// Graph distance from coastline (stops at all other regions).
	rDistanceC := m.assignDistanceField(coastlineRegs, stopReg)

	// Propagate the compression values.
	compPerReg := m.propagateCompression(compressionReg)

	// This code below calculates the height of a given region based on a linear
	// interpolation of the three distance values above.
	//
	// Ideally, we would use some form of noise using the distance to a mountain / faultline
	// to generate a more natural looking landscape with mountain ridges resulting from the
	// folding of the plates.
	//
	// Since we want a "wave" like appearance, we could use one dimensional noise based on the
	// distance to the faultline with some variation for a more natural look.
	const epsilon = 1e-7
	r_xyz := m.XYZ

	// Exponent for interpolation.
	// n = 1 is a linear interpolation
	// n = 2 is a square interpolation
	// n = 0.5 is a square root interpolation
	na := 1.0 / 1.0
	nb := 1.0 / 1.0
	nc := 1.0 / 1.0
	for r := 0; r < m.mesh.numRegions; r++ {
		a := math.Pow(rDistanceA[r], na) + epsilon // Distance from mountains
		b := math.Pow(rDistanceB[r], nb) + epsilon // Distance from oceans
		c := math.Pow(rDistanceC[r], nc) + epsilon // Distance from coastline
		if m.PlateIsOcean[m.RegionToPlate[r]] {
			// Ocean plates are slightly lower than other plates.
			m.Elevation[r] = -0.1
		}
		if math.IsInf(rDistanceA[r], 0) && math.IsInf(rDistanceB[r], 0) {
			// If the distance from mountains and oceans is unset (infinity),
			// we increase the elevation by 0.1 since we wouldn't be able to
			// calculate the harmonic mean.
			m.Elevation[r] += 0.1
		} else {
			// The height is calculated as weighted harmonic mean of the
			// three distance values.
			f := (1/a - 1/b) / (1/a + 1/b + 1/c)

			// Average with plate compression to get some
			// variation in the landscape.
			f = (f + compPerReg[r]) * 0.5

			// Apply a square falloff to the elevaltion values.
			// f *= math.Abs(f)
			m.Elevation[r] += f
		}
		m.Elevation[r] += m.noise.Eval3(r_xyz[3*r], r_xyz[3*r+1], r_xyz[3*r+2])*2 - 1 // Noise from -1.0 to 1.0
	}
}
