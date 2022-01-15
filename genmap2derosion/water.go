package genmap2derosion

import (
	"fmt"
	"github.com/Flokey82/go_gens/vectors"
	//"log"
	"math"
)

// Hydraulic erosion
func (w *World) erode(cycles, drops int) {
	// Reset drains.
	for i := range w.waterdrains {
		w.waterdrains[i] = -1
	}

	// Do a series of iterations!
	for i := 0; i < cycles; i++ {
		fmt.Println(fmt.Sprintf("Erode... (Cycle %d/%d)", i, cycles))

		// Track the Movement of all Particles
		track := make([]int, w.dim.X*w.dim.Y)
		for j := 0; j < drops; j++ {
			// Spawn New Particle
			drop := NewDrop(vectors.NewVec2(
				float64((w.r.Int())%int(w.dim.X)), // Random X pos.
				float64((w.r.Int())%int(w.dim.Y)), // Random Y pos.
			), w.fdim)

			spill := 5
			for drop.volume > minVol && spill != 0 {
				drop.descend(w, track)

				if drop.volume > minVol {
					drop.flood(w)
				}
				spill--
			}
		}

		// Update Path
		lrate := float64(0.01)
		lrateInv := (1.0 - lrate)
		for i := 0; i < len(track); i++ {
			if track[i] > 0 {
				w.waterpath[i] = lrateInv*w.waterpath[i] + lrate
			} else {
				w.waterpath[i] = 0.0 // 0.5 * lrateInv * w.waterpath[i]
			}
		}
	}
}

type Drop struct {
	// Properties
	index    int64
	pos      vectors.Vec2
	speed    vectors.Vec2
	volume   float64 // This will vary in time
	sediment float64 // Sediment concentration
}

// NewDrop returns a new particle at the given position.
func NewDrop(pos, dim vectors.Vec2) Drop {
	var d Drop
	d.init()
	d.pos = pos
	d.index = int64(pos.X*dim.Y + pos.Y)
	return d
}

// NewDropWithVolume returns a new particle at the given position with the given volume.
func NewDropWithVolume(p, dim vectors.Vec2, v float64) Drop {
	var d Drop
	d.init()
	d.pos = p
	d.index = int64(p.X*dim.X + p.X)
	d.volume = v
	return d
}

const (
	// Parameters (move this somewhere else)
	dt             = 1.2
	density        = 1.0 // This gives varying amounts of inertia and stuff...
	evapRate       = 0.001
	depositionRate = 0.1
	minVol         = 0.01
	friction       = 0.1
	volumeFactor   = 100.0 // "Water Deposition Rate"
)

func (d *Drop) init() {
	d.speed = vectors.NewVec2(0.0, 0.0)
	d.volume = 1.0   // This will vary in time
	d.sediment = 0.0 // Sediment concentration
}

func (d *Drop) descend(w *World, track []int) {
	dim := w.dim
	var ipos vectors.IVec2

	for d.volume > minVol {
		// Initial Position
		ipos.X = int64(d.pos.X)
		ipos.Y = int64(d.pos.Y)
		ind := ipos.X*dim.Y + ipos.Y

		// Add to Path
		track[ind]++

		n := w.surfaceNormal(ind)

		// Effective Parameter Set
		// Higher toughness means less erosion.
		// effT := d.depositionRate * math.Max(0.0, 1.0-w.sedimentToughness[ind])

		// Higher plant density means less erosion.
		effD := depositionRate // * math.Max(0.0, 1.0-w.plantdensity[ind])

		// Lower Friction, Lower Evaporation in Streams
		// makes particles prefer established streams -> "curvy".
		effF := friction * (1.0 - 0.5*w.waterpath[ind])
		effR := evapRate * (1.0 - 0.2*w.waterpath[ind])

		// Newtonian Mechanics
		muAcc := 1 / (d.volume * density)
		acc := vectors.NewVec2(n.X*muAcc, n.Z*muAcc)
		d.speed.AddToThis(acc.Mul(dt))
		d.pos.AddToThis(d.speed.Mul(dt))
		d.speed.MulWithThis(1.0 - dt*effF)

		// New Position
		nind := int64(d.pos.X)*dim.Y + int64(d.pos.Y)

		// Out-Of-Bounds
		if nind >= int64(len(w.waterpool)) || nind < 0 || d.pos.X < 0 || d.pos.Y < 0 || int64(d.pos.X) >= int64(dim.X) || int64(d.pos.Y) >= int64(dim.Y) {
			//fmt.Println(fmt.Sprintf("break %f %f", d.pos.x, d.pos.y))
			d.volume = 0.0
			break
		}

		// Particle enters pool or is not accelerated.
		if w.waterpool[nind] > 0.0 || (w.waterpath[nind] > 0.3 && acc.Len() < 0.01) {
			// fmt.Println("not accelerated or entering pool")
			break
		}

		// Mass-Transfer (in MASS)
		c_eq := math.Max(0.0, d.speed.Len()*(w.heightmap[ind]-w.heightmap[nind])) // ind?
		cdiff := c_eq - d.sediment
		d.sediment += dt * cdiff * effD                  // * effT
		w.heightmap[ind] -= d.volume * dt * cdiff * effD // * effT
		w.sediment[ind] -= d.volume * dt * cdiff * effD  // * effT

		// Evaporate (Mass Conservative)
		d.sediment /= (1.0 - dt*effR)
		d.volume *= (1.0 - dt*effR)
	}
}

func (d *Drop) flood(w *World) {
	dim := w.dim
	// Current Height
	d.index = int64(d.pos.X)*dim.Y + int64(d.pos.Y)
	plane := w.heightmap[d.index] + w.waterpool[d.index]
	initialplane := plane

	// Floodset
	set := make([]int64, 0, 1024)
	fail := 10
	drainage := w.drainage
	size := dim.X * dim.Y
	tried := make([]bool, size)

	var drain int
	var drainfound bool
	var fill func(i int64)
	fill = func(i int64) {

		// Out of Bounds
		if i < 0 || i >= size {
			return
		}

		// Position has been tried
		if tried[i] {
			return
		}
		tried[i] = true

		// Wall / Boundary
		if plane < w.heightmap[i]+w.waterpool[i] {
			return
		}

		// Drainage Point
		if initialplane > w.heightmap[i]+w.waterpool[i] {
			if !drainfound || w.waterpool[i]+w.heightmap[i] < w.waterpool[drain]+w.heightmap[drain] {
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

	// Iterate
	for d.volume > minVol && fail != 0 {
		set = make([]int64, 0, 1024)
		for i := range tried {
			tried[i] = false
		}
		drain = 0
		drainfound = false

		// Perform Flood
		fill(d.index)

		// Drainage Point
		if drainfound {

			// Set the Drop Position and Evaporate
			d.pos = vectors.NewVec2(float64(int64(drain)/dim.Y), float64(int64(drain)%dim.Y))

			// Set the New Waterlevel (Slowly)
			plane = (1.0-drainage)*initialplane + drainage*(w.heightmap[drain]+w.waterpool[drain])

			// Compute the New Height
			for _, s := range set {
				if plane > w.heightmap[s] {
					w.waterpool[s] = plane - w.heightmap[s]
					w.waterdrains[s] = drain
				} else {
					w.waterpool[s] = 0.0
				}
			}
			// Remove Sediment
			d.sediment *= 0.1
			break
		}

		// Get Volume under Plane
		// So we sum up the difference between plane and (height[s]+pool[s]) which
		// gives up the total missing volume required for a full flood.
		totalVol := float64(0.0)
		for _, s := range set {
			totalVol += volumeFactor * (plane - (w.heightmap[s] + w.waterpool[s]))
		}

		// We can fill the volume of the sink.
		if totalVol <= d.volume && initialplane < plane {
			// Raise water level to plane height.
			for _, s := range set {
				w.waterpool[s] = plane - w.heightmap[s]
			}

			// Adjust Drop Volume
			d.volume -= totalVol
			totalVol = 0.0
		} else {
			fail-- // Plane was too high.
			//log.Println("fail")
		}

		// Adjust Planes
		if plane > initialplane {
			initialplane = plane
		}
		plane += 0.5 * (d.volume - totalVol) / float64(len(set)) / volumeFactor
	}

	// Couldn't place the volume (for some reason), so ignore this drop.
	if fail == 0 {
		d.volume = 0.0
	}
}
