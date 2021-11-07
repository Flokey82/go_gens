package genmap2d

// genHeightMap generates the heightmap using opensimplex noise.
func (m *Map) genHeightMap() {
	w, h := m.Width, m.Height
	hm := initCells(w, h)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			va := m.Noise.Eval2(2*float64(x)/float64(w), 2*float64(y)/float64(h))
			va += m.Noise.Eval2(float64(x)/float64(w), float64(y)/float64(h)) * 0.1
			va += m.Noise.Eval2(3*float64(x)/float64(w), 3*float64(y)/float64(h)) * 0.2
			va += m.Noise.Eval2(10*float64(x)/float64(w), 10*float64(y)/float64(h)) * 0.05
			hm[x][y] = byte(((va / 1.35) * 128) + 128)
		}
	}
	m.HeightMap = hm
}
