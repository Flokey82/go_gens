package genworldvoronoi

import (
	"container/list"
	"log"
)

func (m *Geo) assignLandmasses() {
	m.Landmasses = m.IdentifyLandmasses()

	lmSize := make(map[int]int)
	for _, lm := range m.Landmasses {
		if lm >= 0 {
			lmSize[lm]++ // Only count regions that are set to a valid ID.
		}
	}
	m.LandmassSize = lmSize
}

// IdentifyLandmasses returns a mapping from region to landmass ID.
// A landmass is a connected number of regions above sealevel.
func (m *Geo) IdentifyLandmasses() []int {
	// NOTE: this is still in need of refinement.
	landMasses := initRegionSlice(m.mesh.numRegions)
	for r := range landMasses {
		// Skip everything that is ocean.
		if m.Elevation[r] <= 0 {
			landMasses[r] = -2
		}
	}

	var landID int
	var landSizes []int
	outReg := make([]int, 0, 8)
	for r, lmID := range landMasses {
		// Skip if the current region has already been allocated
		// or is below sealevel.
		if lmID != -1 {
			continue
		}
		var currentLandSize int
		queue := list.New()
		enqueue := func(r int) {
			// Skip if the current region has already been allocated
			// or is below sealevel.
			if landMasses[r] != -1 {
				return
			}
			landMasses[r] = landID // Assign current landID to the region.
			currentLandSize++      // Increase size of known landmass.
			for _, nb := range m.mesh.r_circulate_r(outReg, r) {
				// Skip if the neighbor region has already been allocated
				// or is below sealevel.
				if landMasses[nb] != -1 {
					continue
				}
				queue.PushBack(nb)
			}
		}
		// Start queue with current region.
		queue.PushBack(r)

		// Process each queue entry until we run out of
		// regions that belong to this landmass.
		for queue.Len() > 0 {
			e := queue.Front()
			enqueue(e.Value.(int))
			queue.Remove(e)
		}

		// Once done, append the current size to the list of landmass-
		// sizes and increment the current landID.
		landSizes = append(landSizes, currentLandSize)
		landID++
	}
	log.Println("number of landmasses", landID)
	// log.Println(landSizes)
	return landMasses
}
