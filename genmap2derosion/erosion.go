package genmap2derosion

import (
	"fmt"
	"log"
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

	// Do a series of iterations!
	now := time.Now()
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
// NOTE: This is untested and probably not working very well due to the climate
// simulation and the probably low-quality precipitation data it produces.
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
