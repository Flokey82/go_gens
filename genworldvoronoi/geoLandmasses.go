package genworldvoronoi

import (
	"container/list"
	"log"
)

// IdentifyLandmasses returns a mapping from region to landmass ID.
// A landmass is a connected number of regions above sealevel.
func (m *Geo) IdentifyLandmasses() []int {
	// NOTE: this is still in need of refinement.
	landMasses := make([]int, m.mesh.numRegions)
	for r := range landMasses {
		landMasses[r] = -1
	}
	var landID int
	var landSizes []int
	for r, h := range m.Elevation {
		// Skip if the current region has already been allocated
		// or is below sealevel.
		if landMasses[r] != -1 || h < 0 {
			continue
		}
		var currentLandSize int
		queue := list.New()
		enqueue := func(r int) {
			// Skip if the current region has already been allocated
			// or is below sealevel.
			if landMasses[r] != -1 || m.Elevation[r] < 0 {
				return
			}
			landMasses[r] = landID // Assign current landID to the region.
			currentLandSize++      // Increase size of known landmass.
			for _, nb := range m.GetRegNeighbors(r) {
				// Skip if the neighbor region has already been allocated
				// or is below sealevel.
				if landMasses[nb] != -1 || m.Elevation[nb] < 0 {
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

// GetLandmassSizes returns a mapping of landmass ID to size in regions.
func (m *Geo) GetLandmassSizes() map[int]int {
	lmSize := make(map[int]int)
	for _, lm := range m.IdentifyLandmasses() {
		if lm >= 0 {
			lmSize[lm]++ // Only count regions that are set to a valid ID.
		}
	}
	return lmSize
}
