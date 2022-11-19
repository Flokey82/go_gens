package genworldvoronoi

// getWaterBodies returns a slice which all regions to enumerated waterbodies/oceans.
//
// NOTE: For regions that are not part of an ocean (elevation above sea level)
// a value of -2 is assigned.
func (m *BaseObject) getWaterBodies() []int {
	// Initialize the waterbody (ocean) mapping.
	done := make([]int, m.mesh.numRegions)
	for i := range done {
		if m.Elevation[i] > 0 {
			done[i] = -2 // Non-ocean regions above sealevel.
		} else {
			done[i] = -1 // Ocean regions that have not been visited yet.
		}
	}

	for r := range done {
		// Skip regions that have already been visited or that are
		// non-ocean / above sealevel.
		if done[r] != -1 {
			continue
		}
		// Set the region index (r) as the ID for the new waterbody.
		done[r] = r

		// diveDeeper is a recursive function that performs a sort
		// of flood fill, assigning the current waterbody ID to all
		// neighboring regions that are ocean regions.
		//
		// TODO: Maybe use a queue instead... we might exceed Go's
		// stack size calling this recursively regardless of how deep
		// the execution stack might go.
		var diveDeeper func(rd int)
		diveDeeper = func(rd int) {
			for _, nbs := range m.GetRegionNeighbors(rd) {
				// If we have reached land or already visited nbs, skip.
				if m.Elevation[nbs] > 0 || done[nbs] != -1 {
					continue
				}
				// Assign the source region index to nbs.
				done[nbs] = r

				// Visit neighbors of nbs.
				diveDeeper(nbs)
			}
		}

		// Recursively assign the waterbody ID / region index (r)
		// to all suitable neighbor regions and their neighbors,
		// and so on.
		diveDeeper(r)
	}
	return done
}

// getWaterBodySizes return a mapping of waterbody IDs to the number of regions
// associated with each waterbody ID.
func (m *BaseObject) getWaterBodySizes() map[int]int {
	wbSize := make(map[int]int)
	for _, wb := range m.Waterbodies {
		if wb >= 0 {
			wbSize[wb]++ // Only count regions that are set to a valid ID.
		}
	}
	return wbSize
}

// getLakeSizes returns a mapping of drainage region to the number of regions that
// drain to this point, effectively summing up the size of each lake.
func (m *BaseObject) getLakeSizes() map[int]int {
	lakeSize := make(map[int]int)
	for _, drain := range m.Drainage {
		if drain != -1 {
			lakeSize[drain]++ // Only count regions that have a drainage point assigned.
		}
	}
	return lakeSize
}

// getRLakeOrWaterBodySize returns the size of the lake or waterbody that the
// provided region is part of.
func (m *BaseObject) getRLakeOrWaterBodySize(r int) int {
	if m.Waterbodies[r] >= 0 {
		return m.WaterbodySize[m.Waterbodies[r]]
	}
	if m.Drainage[r] >= 0 {
		return m.LakeSize[m.Drainage[r]]
	}
	return 0
}
