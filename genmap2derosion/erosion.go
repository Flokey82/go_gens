package genmap2derosion

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Flokey82/go_gens/vectors"
)

var (
	lrate    = float64(0.01)
	lrateInv = (1.0 - lrate)
)

// Hydraulic erosion
func (w *World) doErosion(cycles, drops int) {
	// Reset drains.
	for i := range w.waterdrains {
		w.waterdrains[i] = -1
	}

	now := time.Now()
	// Do a series of iterations!
	for i := 0; i < cycles; i++ {
		fmt.Println(fmt.Sprintf("Erode... (Cycle %d/%d)", i, cycles))
		w.erode(drops)
	}
	log.Println(time.Since(now))
}

// erode performs one iteration of erosion with the given number of drops.
func (w *World) erode(drops int) {
	sx := int(w.params.Size.X)
	sy := int(w.params.Size.Y)

	// We keep track of the movement of all drops to determine flow patterns
	// and river paths.
	track := make([]int, sx*sy)
	for j := 0; j < drops; j++ {
		// Spawn new particle at a random position.
		drop := NewDrop(vectors.NewVec2(
			float64(w.r.Intn(sx)),
			float64(w.r.Intn(sy)),
		))

		// Spill limits the number of times we can perform a flood and/or
		// attempt to move the drop downhll.
		spill := 5
		for drop.volume > minVol && spill != 0 {
			// Move the drop downhill and keep track of the path
			// that it takes.
			drop.descend(w, track)

			// If we still have a sizable water volume left after
			// moving the drop, perform a flood(fill) at the current
			// position where the drop came to rest.
			if drop.volume > minVol {
				drop.flood(w)
			}
			spill--
		}
	}

	// Update the waterpath by checking if we recorded drops
	// passing through any given location.
	for i, t := range track {
		if t > 0 {
			// We had some drops come through, so we refesh the value
			// at the given position by adding the lrate constant.
			//
			// NOTE: This doesn't seem quite right? Shouldn't a higher
			// drop count mean that we increase the path value proportionally
			// to the throughflow?
			w.waterpath[i] = lrateInv*w.waterpath[i] + lrate
		} else {
			// Fade a path value if we didn't have any drops passing through
			// during this pass.
			w.waterpath[i] = 0.25 * lrateInv * w.waterpath[i]
		}
	}

	// If we should store GIF frames, do so.
	if w.params.StoreGIFFrames {
		w.storeGifFrame(w.heightmap, w.waterpath, w.waterpool)
	}
}

// erodeRain is an experimental variation of the "erode" function which initializes
// drops based on the precipitation values we have calculated for each location.
//
// NOTE: This is untested and probably not working very well.
func (w *World) erodeRain(cycles int, rmap []float64) {
	// Reset all recorded drains.
	for i := range w.waterdrains {
		w.waterdrains[i] = -1
	}

	lrate := float64(0.01)
	lrateInv := (1.0 - lrate)

	now := time.Now()
	sx := int(w.params.Size.X)
	sy := int(w.params.Size.Y)

	// Perform a number of erosion cycles.
	for i := 0; i < cycles; i++ {
		fmt.Println(fmt.Sprintf("Erode... (Cycle %d/%d)", i, cycles))

		// Track the Movement of all Particles
		track := make([]int, sx*sy)

		// Instead of spawning a given number of drops at random locations,
		// we spawn drops at all locations where rain falls.
		for j := 0; j < len(rmap); j++ {
			// If the precipitation value is below the minimum,
			// skip this location.
			if rmap[j] <= minVol {
				continue
			}
			// Spawn a new particle at the given location.
			drop := NewDrop(vectors.NewVec2(
				float64(j/sy),
				float64(j%sy),
			))

			// Set the drop volume to the precipitation value.
			drop.volume = rmap[j]

			// Spill limits the number of times we can perform a flood and/or
			// attempt to move the drop downhll.
			spill := 5

			// As long as we have still enough water to move the drop,
			// keep moving it.
			for drop.volume > minVol && spill != 0 {
				// Move the drop downhill and keep track of the path
				// that it takes.
				drop.descend(w, track)

				// If we still have a sizable water volume left after
				// moving the drop, perform a flood(fill) at the current
				// position where the drop came to rest.
				if drop.volume > minVol {
					drop.flood(w)
				}
				spill--
			}
		}

		// Update the waterpath by checking if we recorded drops
		// passing through any given location.
		for i, t := range track {
			if t > 0 {
				// We had some drops come through, so we refesh the value
				// at the given position by adding the lrate constant.
				//
				// NOTE: This doesn't seem quite right? Shouldn't a higher
				// drop count mean that we increase the path value proportionally
				// to the throughflow?
				w.waterpath[i] = lrateInv*w.waterpath[i] + lrate
			} else {
				// Fade a path value if we didn't have any drops passing through
				// during this pass.
				w.waterpath[i] = 0.25 * lrateInv * w.waterpath[i]
			}
		}

		// If we should store GIF frames, do so.
		w.storeGifFrame(w.heightmap, w.waterpath, w.waterpool)
	}
	log.Println(time.Since(now))
}

type Drop struct {
	// Properties
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

const (
	// Parameters (move this somewhere else)
	dt             = 1.2   // Delta T / time factor
	density        = 1.0   // This gives varying amounts of inertia and stuff...
	evapRate       = 0.001 // The rate of evaporation
	depositionRate = 0.1   // The rate of sediment deposition
	minVol         = 0.01  // Minimum water volume
	friction       = 0.1   // Friction coefficient influencing the speed of the drop
	volumeFactor   = 100.0 // "Water Deposition Rate" / water volume multiplier for the flood algorithm
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
			break
		}

		// Mass-Transfer (in MASS)
		// Calculate the equilibrium concentration at the current speed.
		c_eq := math.Max(0.0, d.speed.Len()*(w.heightmap[ind]-w.heightmap[nind])) // ind?

		// Calculate the difference between what is and what can be suspended in the drop.
		//
		// NOTE: Investigate what is going on here as I'd assume that the drop sediment
		// capacity would be dependent on the water volume?
		eqCDiff := c_eq - d.sediment

		// Calculate how much sediment is either eroded or deposited.
		sedimentDiff := dt * eqCDiff * effD

		// Increase or decrease the sediment in the drop.
		d.sediment += sedimentDiff // * effT

		// Remove the calculated sediment amount from the heightmap.
		w.heightmap[ind] -= d.volume * sedimentDiff // * effT

		// Remove the calculated sediment amount from the sedimentmap.
		// This we use to keep track of where erosion has happened and
		// where sediment has been deposited to identify potential
		// fertile land etc.
		w.sediment[ind] -= d.volume * sedimentDiff // * effT

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
	plane := w.heightmap[d.index] + w.waterpool[d.index]
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
		currHeight := w.heightmap[i] + w.waterpool[i]
		if plane < currHeight {
			return
		}

		// Drainage Point
		if initialplane > currHeight {
			if !drainfound || currHeight < w.waterpool[drain]+w.heightmap[drain] {
				// No Drain yet or lower drain.
				drain = int(i)
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
			plane = (1.0-drainage)*initialplane + drainage*(w.heightmap[drain]+w.waterpool[drain])

			// Compute the new height for all flooded cells / indices / locations
			// that we stored in the floodset.
			for _, s := range set {
				if plane > w.heightmap[s] {
					// If the current plane is heigher than the current heightmap value,
					// we make up the difference by adding the difference to the waterpool
					// height.
					w.waterpool[s] = plane - w.heightmap[s]
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
			totalVol += volumeFactor * (plane - (w.heightmap[s] + w.waterpool[s]))
		}

		// Check if we can fill the volume of the sink.
		if totalVol <= d.volume && initialplane < plane {
			// If we can, raise water level to plane height.
			for _, s := range set {
				w.waterpool[s] = plane - w.heightmap[s]
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
			initialplane = plane
		}
		plane += 0.5 * (d.volume - totalVol) / float64(len(set)) / volumeFactor
	}

	// Couldn't place the volume (for some reason), so ignore this drop and set
	// the water volume to 0, which will discard this drop in the erosion function.
	if fail == 0 {
		d.volume = 0.0
	}
}
