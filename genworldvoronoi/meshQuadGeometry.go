package genworldvoronoi

type QuadGeometry struct {
    I   []int
    xyz []float64
    tm  []float64
}

func NewQuadGeometry() *QuadGeometry {
    /* xyz = position in 3-space;
       tm = temperature, moisture
       I = indices for indexed drawing mode */
    return &QuadGeometry{}
}

func (this *QuadGeometry) setMesh(mesh *TriangleMesh) {
    numSides := mesh.numSides
    numRegions := mesh.numRegions
    numTriangles := mesh.numTriangles
    this.I = make([]int, 3*numSides)
    this.xyz = make([]float64, 3*(numRegions+numTriangles))
    this.tm = make([]float64, 2*(numRegions+numTriangles))
}

func (this *QuadGeometry) setMap(mesh *TriangleMesh, m *Map) {
    const V = 0.95
    xyz := this.xyz
    tm := this.tm
    I := this.I

    r_xyz := m.r_xyz
    t_xyz := m.t_xyz
    xyz = r_xyz[:]              //.set(r_xyz)
    xyz = append(xyz, t_xyz...) //.set(t_xyz, r_xyz.length)
    // TODO: multiply all the r, t points by the elevation, taking V into account

    p := 0
    r_elevation := m.r_elevation
    r_moisture := m.r_moisture
    numRegions := mesh.numRegions
    for r := 0; r < numRegions; r++ {
        tm[p] = r_elevation[r]
        tm[p+1] = r_moisture[r]
        p += 2
    }

    t_elevation := m.t_elevation
    t_moisture := m.t_moisture
    numTriangles := mesh.numTriangles
    for t := 0; t < numTriangles; t++ {
        tm[p] = t_elevation[t]
        tm[p+1] = t_moisture[t]
        p += 2
    }

    i := 0
    count_valley := 0
    count_ridge := 0
    s_flow := m.s_flow
    numSides := mesh.numSides
    for s := 0; s < numSides; s++ {
        opposite_s := mesh.s_opposite_s(s)
        r1 := mesh.s_begin_r(s)
        r2 := mesh.s_begin_r(opposite_s)
        t1 := mesh.s_inner_t(s)
        t2 := mesh.s_inner_t(opposite_s)

        // Each quadrilateral is turned into two triangles, so each
        // half-edge gets turned into one. There are two ways to fold
        // a quadrilateral. This is usually a nuisance but in this
        // case it's a feature. See the explanation here
        // https://www.redblobgames.com/x/1725-procedural-elevation/#rendering
        coast := r_elevation[r1] < 0.0 || r_elevation[r2] < 0.0
        if coast || s_flow[s] > 0 || s_flow[opposite_s] > 0 {
            // It's a coastal or river edge, forming a valley
            I[i] = r1
            I[i+1] = numRegions + t2
            I[i+2] = numRegions + t1
            i += 3
            count_valley++
        } else {
            // It's a ridge
            I[i] = r1
            I[i+1] = r2
            I[i+2] = numRegions + t1
            i += 3
            count_ridge++
        }
    }
    this.I = I
    this.xyz = xyz
    this.tm = tm

    //console.log('ridge=', count_ridge, ', valley=', count_valley);
}
