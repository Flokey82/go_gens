package genworldvoronoi

type QuadGeometry struct {
	indexes []int     // indices for indexed drawing mode
	xyz     []float64 // position in 3D-space
	tm      []float64 // temperature, moisture
}

func NewQuadGeometry(mesh *TriangleMesh) *QuadGeometry {
	qg := &QuadGeometry{}
	qg.setMesh(mesh)
	return qg
}

func (qg *QuadGeometry) setMesh(mesh *TriangleMesh) {
	numSides := mesh.numSides
	numRegions := mesh.numRegions
	numTriangles := mesh.numTriangles
	qg.indexes = make([]int, 3*numSides)
	qg.xyz = make([]float64, 3*(numRegions+numTriangles))
	qg.tm = make([]float64, 2*(numRegions+numTriangles))
}

func (qg *QuadGeometry) setMap(mesh *TriangleMesh, m *Geo) {
	const v = 0.95
	xyz := qg.xyz
	tm := qg.tm
	idxs := qg.indexes

	regXYZ := m.XYZ
	triXYZ := m.tXYZ
	xyz = append(regXYZ[:], triXYZ...)
	// TODO: multiply all the r, t points by the elevation, taking V into account

	p := 0
	regElevation := m.Elevation
	regMoisture := m.Moisture
	numRegions := mesh.numRegions
	for r := 0; r < numRegions; r++ {
		tm[p] = regElevation[r]
		tm[p+1] = regMoisture[r]
		p += 2
	}

	triElevation := m.triElevation
	triMoisture := m.triMoisture
	numTriangles := mesh.numTriangles
	for t := 0; t < numTriangles; t++ {
		tm[p] = triElevation[t]
		tm[p+1] = triMoisture[t]
		p += 2
	}

	i := 0
	countValley := 0
	countRidge := 0
	sideFlow := m.sideFlow
	for side := 0; side < mesh.numSides; side++ {
		oppositeSide := mesh.s_opposite_s(side)
		r1 := mesh.s_begin_r(side)
		r2 := mesh.s_begin_r(oppositeSide)
		t1 := mesh.s_inner_t(side)
		t2 := mesh.s_inner_t(oppositeSide)

		// Each quadrilateral is turned into two triangles, so each
		// half-edge gets turned into one. There are two ways to fold
		// a quadrilateral. This is usually a nuisance but in this
		// case it's a feature. See the explanation here
		// https://www.redblobgames.com/x/1725-procedural-elevation/#rendering
		coast := regElevation[r1] < 0.0 || regElevation[r2] < 0.0
		if coast || sideFlow[side] > 0 || sideFlow[oppositeSide] > 0 {
			// It's a coastal or river edge, forming a valley
			idxs[i] = r1
			idxs[i+1] = numRegions + t2
			idxs[i+2] = numRegions + t1
			i += 3
			countValley++
		} else {
			// It's a ridge
			idxs[i] = r1
			idxs[i+1] = r2
			idxs[i+2] = numRegions + t1
			i += 3
			countRidge++
		}
	}
	qg.indexes = idxs
	qg.xyz = xyz
	qg.tm = tm

	//console.log('ridge=', count_ridge, ', valley=', count_valley);
}
