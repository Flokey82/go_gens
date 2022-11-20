package genworldvoronoi

// getBorders returns the borders of each territory as list of triangle center points.
func (m *Map) getBorders() [][]int {
	return m.getCustomBorders(m.RegionToEmpire)
}

// getLakeBorders returns the borders of each lake (regions with same drainage region) as list of triangle center points.
func (m *Map) getLakeBorders() [][]int {
	return m.getCustomBorders(m.Drainage)
}

// getCustomBorders returns the borders/contours of all region in the supplied slice that have the same value.
func (m *Map) getCustomBorders(regionToID []int) [][]int {
	return m.GetCustomContour(func(idxA, idxB int) bool {
		if m.Elevation[idxA] < 0 || m.Elevation[idxB] < 0 ||
			(regionToID[idxA] < 0 && regionToID[idxB] < 0) {
			return false
		}
		return regionToID[idxA] != regionToID[idxB]
	})
}

// getLandmassBorders returns the borders of each landmass (neighboring regions above sea level)
// as list of triangle center points.
func (m *Map) getLandmassBorders() [][]int {
	return m.GetCustomContour(func(idxA, idxB int) bool {
		return m.Elevation[idxA] >= 0 && m.Elevation[idxB] < 0 || m.Elevation[idxA] < 0 && m.Elevation[idxB] >= 0
	})
}
