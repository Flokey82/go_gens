package genworldvoronoi

// getBorders returns the borders of each territory as list of triangle center points.
func (m *Map) getBorders() [][]int {
	var edges [][2]int
	seen := make(map[[2]int]bool)
	for i := 0; i < len(m.mesh.Halfedges); i++ {
		idxA := m.mesh.s_begin_r(i)
		idxB := m.mesh.s_end_r(i)
		var vx [2]int
		if idxA > idxB {
			vx[0] = m.mesh.s_outer_t(i)
			vx[1] = m.mesh.s_inner_t(i)
		} else {
			vx[0] = m.mesh.s_inner_t(i)
			vx[1] = m.mesh.s_outer_t(i)
		}
		if seen[vx] {
			continue
		}
		seen[vx] = true
		if m.r_elevation[idxA] < 0 || m.r_elevation[idxB] < 0 ||
			(m.r_territory[idxA] < 0 && m.r_territory[idxB] < 0) {
			continue
		}
		if m.r_territory[idxA] != m.r_territory[idxB] {
			edges = append(edges, vx)
		}
	}
	return mergeIndexSegments(edges)
}

// getLakeBorders returns the borders of each lake (regions with same drainage region) as list of triangle center points.
func (m *Map) getLakeBorders() [][]int {
	var edges [][2]int
	seen := make(map[[2]int]bool)
	for i := 0; i < len(m.mesh.Halfedges); i++ {
		idxA := m.mesh.s_begin_r(i)
		idxB := m.mesh.s_end_r(i)
		var vx [2]int
		if idxA > idxB {
			vx[0] = m.mesh.s_outer_t(i)
			vx[1] = m.mesh.s_inner_t(i)
		} else {
			vx[0] = m.mesh.s_inner_t(i)
			vx[1] = m.mesh.s_outer_t(i)
		}
		if seen[vx] {
			continue
		}
		seen[vx] = true
		if m.r_elevation[idxA] < 0 || m.r_elevation[idxB] < 0 ||
			(m.r_drainage[idxA] < 0 && m.r_drainage[idxB] < 0) {
			continue
		}
		if m.r_drainage[idxA] != m.r_drainage[idxB] {
			edges = append(edges, vx)
		}
	}
	return mergeIndexSegments(edges)
}

func (m *Map) contour() [][]int {
	var edges [][2]int
	seen := make(map[[2]int]bool)
	for i := 0; i < len(m.mesh.Halfedges); i++ {
		idxA := m.mesh.s_begin_r(i)
		idxB := m.mesh.s_end_r(i)
		var vx [2]int
		if idxA > idxB {
			vx[0] = m.mesh.s_outer_t(i)
			vx[1] = m.mesh.s_inner_t(i)
		} else {
			vx[0] = m.mesh.s_inner_t(i)
			vx[1] = m.mesh.s_outer_t(i)
		}
		if seen[vx] {
			continue
		}

		if m.r_elevation[idxA] >= 0 && m.r_elevation[idxB] < 0 || m.r_elevation[idxA] < 0 && m.r_elevation[idxB] >= 0 {
			edges = append(edges, vx)
		}
	}

	return mergeIndexSegments(edges)
}
