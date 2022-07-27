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

func (w *World) erode(drops int) {
	// Track the Movement of all Particles
	sx := int(w.params.Size.X)
	sy := int(w.params.Size.Y)
	track := make([]int, sx*sy)
	for j := 0; j < drops; j++ {
		// Spawn New Particle
		drop := NewDrop(vectors.NewVec2(
			float64(w.r.Intn(sx)), // Random X pos.
			float64(w.r.Intn(sy)), // Random Y pos.
		))

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
	for i, t := range track {
		if t > 0 {
			w.waterpath[i] = lrateInv*w.waterpath[i] + lrate
		} else {
			w.waterpath[i] = 0.25 * lrateInv * w.waterpath[i]
		}
	}
	if w.params.StoreGIFFrames {
		w.storeGifFrame(w.heightmap, w.waterpath, w.waterpool)
	}
}

func (w *World) erodeRain(cycles int, rmap []float64) {
	// Reset drains.
	for i := range w.waterdrains {
		w.waterdrains[i] = -1
	}

	lrate := float64(0.01)
	lrateInv := (1.0 - lrate)

	now := time.Now()
	sx := int(w.params.Size.X)
	sy := int(w.params.Size.Y)
	// Do a series of iterations!
	for i := 0; i < cycles; i++ {
		fmt.Println(fmt.Sprintf("Erode... (Cycle %d/%d)", i, cycles))

		// Track the Movement of all Particles
		track := make([]int, sx*sy)
		for j := 0; j < len(rmap); j++ {
			if rmap[j] <= minVol {
				continue
			}
			// Spawn New Particle
			drop := NewDrop(vectors.NewVec2(
				float64(j/sy),
				float64(j%sy),
			))
			drop.volume = rmap[j]

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
		for i, t := range track {
			if t > 0 {
				w.waterpath[i] = lrateInv*w.waterpath[i] + lrate
			} else {
				w.waterpath[i] = 0.25 * lrateInv * w.waterpath[i]
			}
		}
		w.storeGifFrame(w.heightmap, w.waterpath, w.waterpool)
	}
	log.Println(time.Since(now))
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
func NewDrop(pos vectors.Vec2) Drop {
	return Drop{
		volume: 1.0, // This will vary in time
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
	dt             = 1.2
	density        = 1.0 // This gives varying amounts of inertia and stuff...
	evapRate       = 0.001
	depositionRate = 0.1
	minVol         = 0.01
	friction       = 0.1
	volumeFactor   = 100.0 // "Water Deposition Rate"
)

func (d *Drop) descend(w *World, track []int) {
	dim := w.params.Size
	var acc vectors.Vec2
	var ind, nind int64
	var effR, effD, effF, muAcc float64
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

		// Lower Friction, Lower Evaporation in Streams
		// makes particles prefer established streams -> "curvy".
		effF = friction * (1.0 - 0.5*w.waterpath[ind])
		effR = evapRate * (1.0 - 0.2*w.waterpath[ind])

		// Newtonian Mechanics
		muAcc = 1 / (d.volume * density)
		acc = vectors.NewVec2(n.X*muAcc, n.Z*muAcc)
		d.speed.AddToThis(acc.Mul(dt))
		d.pos.AddToThis(d.speed.Mul(dt))
		d.speed.MulWithThis(1.0 - dt*effF)

		// New Position
		nind = int64(d.pos.X)*dim.Y + int64(d.pos.Y)

		// Out-Of-Bounds
		if nind >= int64(len(w.waterpool)) || nind < 0 {
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
		bla := dt * cdiff * effD
		d.sediment += bla                  // * effT
		w.heightmap[ind] -= d.volume * bla // * effT
		w.sediment[ind] -= d.volume * bla  // * effT

		// Evaporate (Mass Conservative)
		d.sediment /= (1.0 - dt*effR)
		d.volume *= (1.0 - dt*effR)
	}
}

func (d *Drop) flood(w *World) {
	dim := w.params.Size
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
		// Out of Bounds or position has been tried.
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

	// Iterate
	for d.volume > minVol && fail != 0 {
		set = set[:0]
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
			d.pos.X = float64(int64(drain) / dim.Y)
			d.pos.Y = float64(int64(drain) % dim.Y)

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
		var totalVol float64
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
