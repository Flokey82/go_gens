package genworldvoronoi

import (
	"math"
)

// rErode erodes all region by the given amount.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) rErode(amount float64) []float64 {
	er := m.rErosionRate()
	newh := make([]float64, m.mesh.numRegions)
	_, maxr := minMax(er)
	for i := 0; i < m.mesh.numRegions; i++ {
		newh[i] = m.r_elevation[i] - amount*(er[i]/maxr)
	}
	return newh
}

// rErosionRate returns the erosion rate per region.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) rErosionRate() []float64 {
	const nbErosionFactor = 0.015
	m.assignFlux()
	flux := m.r_flux
	_, maxFlux := minMax(m.r_flux)
	slope := m.getRSlope()
	newh := make([]float64, m.mesh.numRegions)
	for i := 0; i < m.mesh.numRegions; i++ {
		river := math.Sqrt((flux[i]/maxFlux)+m.r_pool[i]) * slope[i]
		creep := slope[i] * slope[i]
		total := 1000*river + creep
		if total > 200 {
			total = 200
		}

		// Additionally erode all neighbors by a certain fraction.
		nbs := m.rNeighbors(i)
		for _, nb := range nbs {
			newh[nb] += total * nbErosionFactor
		}
		newh[i] += total
	}
	return newh
}

// getRSlope returns the region slope by averaging the slopes of the triangles
// around a given region.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) getRSlope() []float64 {
	slope := make([]float64, m.mesh.numRegions)
	for i := 0; i < m.mesh.numRegions; i++ {
		s := m.rSlope(i)
		slope[i] = math.Sqrt(s[0]*s[0] + s[1]*s[1])
	}
	return slope
}

func (m *Map) rSlope(i int) [2]float64 {
	var res [2]float64
	var count int
	for _, t := range m.mesh.r_circulate_t(nil, i) {
		slope := m.rTriSlope(m.mesh.t_circulate_r(nil, t))
		res[0] += slope[0]
		res[1] += slope[1]
		count++
	}

	res[0] /= float64(count)
	res[1] /= float64(count)
	return res
}

// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Map) rTriSlope(nbs []int) [2]float64 {
	if len(nbs) != 3 {
		return [2]float64{0, 0}
	}

	p0 := convToVec3(m.r_xyz[nbs[0]*3:])
	p1 := convToVec3(m.r_xyz[nbs[1]*3:])
	p2 := convToVec3(m.r_xyz[nbs[2]*3:])

	x1 := p1.X - p0.X
	x2 := p2.X - p0.X
	y1 := p1.Y - p0.Y
	y2 := p2.Y - p0.Y

	det := x1*y2 - x2*y1
	h1 := m.r_elevation[nbs[1]] - m.r_elevation[nbs[0]]
	h2 := m.r_elevation[nbs[2]] - m.r_elevation[nbs[0]]

	return [2]float64{
		(y2*h1 - y1*h2) / det,
		(-x2*h1 + x1*h2) / det,
	}
}
