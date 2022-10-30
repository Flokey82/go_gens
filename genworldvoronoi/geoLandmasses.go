package genworldvoronoi

import (
	"container/list"
	"log"
)

// identifyLandmasses returns a mapping from region to landmass ID.
// A landmass is a connected number of regions above sealevel.
func (m *Map) identifyLandmasses() []int {
	// NOTE: this is still in need of refinement.
	landMasses := make([]int, m.mesh.numRegions)
	for r := range landMasses {
		landMasses[r] = -1
	}
	var landID int
	var landSizes []int
	for r, h := range m.r_elevation {
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
			if landMasses[r] != -1 || m.r_elevation[r] < 0 {
				return
			}
			landMasses[r] = landID // Assign current landID to the region.
			currentLandSize++      // Increase size of known landmass.
			for _, nb := range m.rNeighbors(r) {
				// Skip if the neighbor region has already been allocated
				// or is below sealevel.
				if landMasses[nb] != -1 || m.r_elevation[nb] < 0 {
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
	log.Println(landSizes)
	return landMasses
}
