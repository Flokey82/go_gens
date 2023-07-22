package genmap2derosion

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

type Drop struct {
	index    int64        // Current position (expressed as index into the heightmap)
	pos      vectors.Vec2 // Current position (expressed as x,y)
	speed    vectors.Vec2 // Speed vector
	volume   float64      // Current water volume of the drop
	sediment float64      // Sediment suspended in the drop
}

// NewDrop returns a new particle at the given position.
func NewDrop(pos vectors.Vec2) Drop {
	return Drop{
		volume: 1.0,
		pos:    pos,
	}
}

// NewDropWithVolume returns a new particle at the given position with the given volume.
func NewDropWithVolume(pos vectors.Vec2, v float64) Drop {
	return Drop{
		pos:    pos,
		volume: v,
	}
}

// Constants for hydrology and erosion.
// TODO: Move this somewhere else.
const (
	dt             = 1.2   // Delta T / time factor
	density        = 1.0   // This gives varying amounts of inertia and stuff...
	evapRate       = 0.001 // The rate of evaporation
	depositionRate = 0.01  // The rate of sediment deposition
	minVol         = 0.01  // Minimum water volume
	friction       = 0.1   // Friction coefficient influencing the speed of the drop
	volumeFactor   = 100.0 // Factor of water volume to height. A volume of 100 equals a height of 1.
	useSedimentMap = true  // Whether to use sediment or not.
)

func (d *Drop) descend(w *World, track []int) {
	dim := w.params.Size
	var acc vectors.Vec2
	var ind, nind int64
	var effR, effD, effF, dropMass float64
	for d.volume > minVol {
		// Initial Position
		ind = int64(d.pos.X)*dim.Y + int64(d.pos.Y)

		// Add to Path
		track[ind]++

		n := w.surfaceNormal(ind)

		// Effective Parameter Set
		// Higher toughness means less erosion.
		// effT := d.depositionRate * math.Max(0.0, 1.0-w.sedimentToughness[ind])

		// Higher plant density means less erosion.
		effD = depositionRate // * math.Max(0.0, 1.0-w.plantdensity[ind])

		// Lower Friction, Lower Evaporation in streams
		// makes particles prefer established streams -> "curvy".
		effF = friction * (1.0 - 0.5*w.waterpath[ind])
		effR = evapRate * (1.0 - 0.2*w.waterpath[ind])

		// Newtonian Mechanics

		// Calculate the mass of the drop.
		dropMass = d.volume * density

		// Calculate the acceleration vector based on the normal vector and drop mass.
		acc = vectors.NewVec2(n.X/dropMass, n.Z/dropMass)

		// Calculate the new speed vector by adding the accelleration vector multiplied
		// with the time factor to the current speed.
		d.speed.AddToThis(acc.Mul(dt))

		// Calculate the new position by adding the speed vector multiplied with the
		// time factor to the current position.
		d.pos.AddToThis(d.speed.Mul(dt))

		// Now slow down the drop based on the friction the drop has experienced during
		// moving for the last time step.
		d.speed.MulWithThis(1.0 - dt*effF)

		// Calculate the new heightmap index from the new position.
		nind = int64(d.pos.X)*dim.Y + int64(d.pos.Y)

		// Check if we are out of bounds (if the drop has left the map).
		if nind >= int64(len(w.waterpool)) || nind < 0 {
			//fmt.Println(fmt.Sprintf("break %f %f", d.pos.x, d.pos.y))
			// Nothing to do with the drop anymore, so set the volume to 0.
			d.volume = 0.0
			break
		}

		// Check if a particle / drop enters pool or is not accelerated.
		if w.waterpool[nind] > 0.0 || (w.waterpath[nind] > 0.3 && acc.Len() < 0.01) {
			// fmt.Println("not accelerated or entering pool")
			// TODO: We should deposit all of our sediment here.
			// w.sediment[nind] += d.sediment
			// d.sediment = 0.0
			break
		}

		// Mass-Transfer (in MASS)
		// Calculate the equilibrium concentration at the current speed.
		c_eq := math.Max(0.0, d.speed.Len()*((w.height(ind))-(w.height(nind)))) // ind?

		// Calculate the difference between what is and what can be suspended in the drop.
		//
		// NOTE: Investigate what is going on here as I'd assume that the drop sediment
		// capacity would be dependent on the water volume?
		eqCDiff := c_eq - d.sediment

		// Calculate how much sediment is either eroded or deposited.
		sedimentDiff := dt * eqCDiff * effD

		// Increase or decrease the sediment in the drop.
		d.sediment += sedimentDiff // * effT

		// Calculate the sediment amount we need to remove from the heightmap.
		sedimentAmount := d.volume * sedimentDiff // * effT

		// Check if we are using the sediment map or not.
		if useSedimentMap {
			// Check if we are eroding or depositing sediment.
			if sedimentAmount > 0 { // Eroding
				if sedimentAmount >= w.sediment[ind] {
					sedimentAmount -= w.sediment[ind]
					w.sediment[ind] = 0.0
				} else {
					w.sediment[ind] -= sedimentAmount
					sedimentAmount = 0.0
				}
				// If we still have some sediment to erode, we will
				// erode the heightmap at the current index.
				if sedimentAmount > 0 {
					w.heightmap[ind] -= sedimentAmount
				}
			} else if sedimentAmount < 0 { // Depositing
				// If we are depositing sediment, we will add the sediment
				// amount to the sediment map.
				w.sediment[ind] -= sedimentAmount
			}
		} else {
			// Remove the calculated sediment amount from the heightmap.
			w.heightmap[ind] -= sedimentAmount
		}

		// Evaporate (Mass Conservative)
		//
		// NOTE: One would think that evaporation would cause some
		// sediment to be deposited, but this is not the case here
		// in this implementation.
		d.sediment /= (1.0 - dt*effR)
		d.volume *= (1.0 - dt*effR)
	}
}

// Flood performs a flood algorithm for the drop on its current position.
func (d *Drop) flood(w *World) {
	dim := w.params.Size

	// Get current surface height (heightmap + waterpool height).
	d.index = int64(d.pos.X)*dim.Y + int64(d.pos.Y)

	// We use 'plane' as the initial height to identify potential drains.
	plane := w.height(d.index) + w.waterpool[d.index]
	initialplane := plane

	// Floodset will keep track of the flooded cells.
	set := make([]int64, 0, 1024)

	// Fail limits the number of attempts to find a drain.
	fail := 10

	drainage := w.drainage
	size := dim.X * dim.Y

	// Tried keeps track of all cells / indices / locations we
	// have visited during the flood fill.
	tried := make([]bool, size)

	var drain int
	var drainfound bool

	// Fill will flood a given cell / index / location and its
	// neighbors until a drain is found.
	var fill func(i int64)
	fill = func(i int64) {
		// BIG_TODO: Document this code.

		// Out of Bounds or cell / index / location has been visited.
		if i < 0 || i >= size || tried[i] {
			return
		}
		tried[i] = true

		// Wall / Boundary
		//
		// Calculate current height of this cell / index / location
		currHeight := w.height(i) + w.waterpool[i]
		if plane < currHeight {
			// If this cell / index / location has a surface height
			// above the target flood plane, it is considered a wall or
			// a boundary of the sink we try to flood up to the height
			// of 'plane' and we return early.
			//
			// currentHeight -.
			//                 \------- <- plane
			//                  \______
			//
			return
		}

		// Drainage Point
		// If this cell / index / location has a surface height that is below the
		// initial plane that we started width, we have found a drainage point.
		if initialplane > currHeight {
			// Only set the drainage point if no drain has been found yet or if the
			// new drainage point is below the surface height of the previous drain.
			if !drainfound || currHeight < w.height(int64(drain))+w.waterpool[drain] {
				drain = int(i) // No Drain yet or lower drain.
			}
			drainfound = true
			return
		}

		// Part of the Pool
		set = append(set, i)
		fill(i + dim.Y) // Fill Neighbors
		fill(i - dim.Y)
		fill(i + 1)
		fill(i - 1)
		fill(i + dim.Y + 1) // Diagonals (Improves Drainage)
		fill(i - dim.Y - 1)
		fill(i + dim.Y - 1)
		fill(i - dim.Y + 1)
	}

	// Iterate until we have exhausted the drop's water volume,
	// we have found a drain, or we have exceeded the fail limit.
	for d.volume > minVol && fail != 0 {
		// Reset the floodset.
		set = set[:0]

		// Reset the array that keeps track of visited cells.
		for i := range tried {
			tried[i] = false
		}

		// Reset drain and drainfound.
		drain = 0
		drainfound = false

		// Perform flood at the current drop index / position / cell.
		fill(d.index)

		// Check if the fill function has found a drain.
		if drainfound {
			// Set the new prop position to where we drain
			// the flooded pool.
			d.pos.X = float64(int64(drain) / dim.Y)
			d.pos.Y = float64(int64(drain) % dim.Y)

			// Set the new waterlevel / surface plane (slowly, using the drainage factor).
			// NOTE: I am not sure why we do this istead of just setting the waterlevel to the
			// surface plane at the drainage point directly.
			plane = (1.0-drainage)*initialplane + drainage*(w.height(int64(drain))+w.waterpool[drain])

			// Compute the new height for all flooded cells / indices / locations
			// that we stored in the floodset.
			for _, s := range set {
				if plane > w.heightmap[s]+w.sediment[s] {
					// If the current plane is heigher than the current heightmap value,
					// we make up the difference by adding the difference to the waterpool
					// height.
					w.waterpool[s] = plane - w.height(s)
					w.waterdrains[s] = drain
				} else {
					// If the current flood plane is below the heightmap,
					// make sure we set the waterpool height to 0.
					w.waterpool[s] = 0.0
				}
			}
			// Remove sediment.
			// NOTE: Why do we do this without depositing it somewhere?
			d.sediment *= 0.1
			break
		}

		// Calculate the total water volume under the flood plane.
		// We sum up the difference between plane and (height[s]+pool[s]) which
		// gives up the total missing volume required for a full flood.
		var totalVol float64
		for _, s := range set {
			totalVol += volumeFactor * (plane - (w.height(s) + w.waterpool[s]))
		}

		// Check if we can fill the volume of the sink.
		if totalVol <= d.volume && initialplane < plane {
			// If we can, raise water level to plane height.
			for _, s := range set {
				w.waterpool[s] = plane - w.height(s)
			}

			// Adjust the drop volume by removing the volume
			// that we have spent filling the sink up to the
			// flood plane.
			d.volume -= totalVol
			totalVol = 0.0
		} else {
			// We weren't able to fill the plane since the
			// water volume of the drow was not enough to
			// fill the sink up to the plane height.
			fail-- // Plane was too high.
		}

		// Adjust the planes.
		// BIG_TODO: Figure out why we do this in this way and how
		// planes actually work and what the logic behind adjusting
		// is.
		if plane > initialplane {
			// Since we have failed to find a drain, we will set the initial
			// plane to the current plane, effectively increasing the new
			// initial surface height.
			initialplane = plane
		}

		// We increase the flood plane height by increasing the height by a
		// factor calculated from the difference between drop volume and total
		// flood volume spread over the number of cells in the floodset.
		volRemaining := d.volume - totalVol

		// Divide the remaining volume by the number of flooded cells.
		volRemainingPerCell := volRemaining / float64(len(set))

		// Once we have calculated the new plane, we will retry to fill the sink
		// given the new flood plane height.
		//
		// NOTE: The potential height per cell is the volume divided by the volume
		// factor.
		//
		// NOTE: By adding half of the remaining volume height during each iteration
		// we are getting closer and closer to the maximum possible plane height
		// given the water volume of the drop.
		//
		// 0.5 -> 0.75 -> 0.875 -> 0.9375 -> 0.96875 -> 0.984375 ...
		plane += 0.5 * volRemainingPerCell / volumeFactor
	}

	// Couldn't place the volume (for some reason), so ignore this drop and set
	// the water volume to 0, which will discard this drop in the erosion function.
	if fail == 0 {
		d.volume = 0.0
	}
}
