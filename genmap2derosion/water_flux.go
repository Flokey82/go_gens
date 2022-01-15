package genmap2derosion

import (
	"log"
	"sort"
)

func (w *World) getDownhill() [][]int {
	idxs := make([][]int, len(w.heightmap))
	for i, h := range w.heightmap {
		nbs := getNeighbors(i, w.heightmap[:], int(w.dim.Y))
		for _, nb := range nbs {
			if nbh := w.heightmap[nb]; nbh < h {
				idxs[i] = append(idxs[i], nb)
			}
		}
	}
	return idxs
}

func (w *World) getRawFlux() []float64 {
	flux := make([]float64, len(w.heightmap))
	idxs := make([]int, len(w.heightmap))
	for i := 0; i < len(w.heightmap); i++ {
		idxs[i] = i
		flux[i] = 1.0 / float64(len(w.heightmap))
	}
	sort.Slice(idxs, func(a, b int) bool {
		return w.heightmap[idxs[b]]-w.heightmap[idxs[a]] < 0
	})

	dh := w.getDownhill()

	// Store all sinks / dead ends.
	var sinks []int
	for i := range w.heightmap {
		j := idxs[i]
		h := w.heightmap[j]
		fluxVol := flux[j]
		var sumh float64
		nbs := dh[j]
		for _, nb := range nbs {
			sumh += h - w.heightmap[nb]
		}
		if len(nbs) == 0 {
			sinks = append(sinks, j)
		} else {
			for _, nb := range nbs {
				flux[nb] += ((h - w.heightmap[nb]) / sumh) * fluxVol
			}
		}
	}
	w.fluxwaterpool = [worldsize * worldsize]float64{}

	dim := w.dim
	for _, idx := range sinks {
		flu := flux[idx]
		volume := flux[idx] //*float64(len(w.heightmap))
		// Take flux at location and distribute in pool.
		// Fill sink.
		plane := w.heightmap[idx] + w.fluxwaterpool[idx]
		initialplane := plane

		// Floodset
		set := make([]int64, 0, 1024)
		fail := 10
		volumeFactor := 100.0 * float64(len(w.heightmap))
		drainage := w.drainage
		size := len(w.heightmap)
		tried := make([]bool, size)

		var drain int
		var drainfound bool
		var fill func(i int64)
		fill = func(i int64) {
			// Out of Bounds
			if i < 0 || i >= int64(size) {
				return
			}

			// Position has been tried
			if tried[i] {
				return
			}
			tried[i] = true

			// Wall / Boundary
			if plane < w.heightmap[i]+w.fluxwaterpool[i] {
				return
			}

			// Drainage Point
			if initialplane > w.heightmap[i]+w.fluxwaterpool[i] {
				if !drainfound || w.waterpool[i]+w.heightmap[i] < w.fluxwaterpool[drain]+w.heightmap[drain] {
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
		for volume > 0.001 && fail != 0 {

			set = make([]int64, 0, 1024)
			for i := range tried {
				tried[i] = false
			}
			drain = 0
			drainfound = false

			// Perform Flood
			fill(int64(idx))

			// Drainage Point
			if drainfound {
				// Set the New Waterlevel (Slowly)
				plane = (1.0-drainage)*initialplane + drainage*(w.heightmap[drain]+w.fluxwaterpool[drain])

				// Compute the New Height
				for _, s := range set {
					if plane > w.heightmap[s] {
						w.fluxwaterpool[s] = plane - w.heightmap[s]
						flux[s] += flu / float64(len(set))
					} else {
						w.fluxwaterpool[s] = 0.0
					}
				}
				break
			}

			// Get Volume under Plane
			tVol := float64(0.0)
			for _, s := range set {
				tVol += volumeFactor * (plane - (w.heightmap[s] + w.fluxwaterpool[s] + flux[s])) //
			}
			// We can partially fill this volume
			if tVol <= volume && initialplane < plane {
				// Raise water level to plane height
				for _, s := range set {
					w.fluxwaterpool[s] = plane - w.heightmap[s]
				}
				// Adjust Drop Volume
				volume -= tVol
				tVol = 0.0
			} else {
				fail-- // Plane was too high.
			}

			// Adjust Planes
			if plane > initialplane {
				initialplane = plane
			}
			plane += 0.5 * (volume - tVol) / float64(len(set)) / volumeFactor
		}

		// Couldn't place the volume (for some reason)- so ignore this drop.
		if fail == 0 {
			volume = 0.0
		}
	}
	return flux
}

func (w *World) getFlux() []float64 {
	flux := make([]float64, len(w.heightmap))
	idxs := make([]int, len(w.heightmap))
	for i := 0; i < len(w.heightmap); i++ {
		idxs[i] = i
		flux[i] = 1.0 / float64(len(w.heightmap))
	}
	sort.Slice(idxs, func(a, b int) bool {
		return w.heightmap[idxs[b]]-w.heightmap[idxs[a]] < 0
	})

	// Store all sinks / dead ends.
	var sinks []int

	// Find all points that are known to be in a water pool
	// and add their flux value to the drainage point.
	for i := range w.heightmap {
		j := idxs[i]
		fluxVol := flux[j]
		if w.waterdrains[j] < 0 {
			continue
		}
		for w.waterdrains[j] >= 0 {
			j = w.waterdrains[j]
		}
		flux[j] += fluxVol
		flux[j] = 0
	}
	for i := range w.heightmap {
		j := idxs[i]
		h := w.heightmap[j]
		fluxVol := flux[j]
		// NOTE: This can deadlock!
		for w.waterdrains[j] >= 0 {
			j = w.waterdrains[j]
		}
		var sumh float64
		var lower []int
		nbs := getNeighbors(j, w.heightmap[:], int(w.dim.Y))
		for _, nb := range nbs {
			if nbh := w.heightmap[nb]; nbh < h {
				sumh += h - nbh
				lower = append(lower, nb)
			}
		}
		if len(lower) == 0 {
			// Sink.
			sinks = append(sinks, j)
		} else {
			for _, nb := range lower {
				fluxShed := ((h - w.heightmap[nb]) / sumh) * fluxVol
				if w.waterdrains[nb] >= 0 {
					if h < w.heightmap[w.waterdrains[nb]] {
						log.Println("!!!drain > h")
					} else {
						log.Println("drain < h")
					}
					flux[w.waterdrains[nb]] += fluxShed
				} else {
					flux[nb] += fluxShed
				}
			}
		}
	}
	return flux
}
