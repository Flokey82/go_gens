package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/vmesh"
	//"github.com/pzsz/voronoi"
	//"log"
	"math"
)

func getSteepness(h *vmesh.Heightmap) *vmesh.Heightmap {
	m := h.Mesh
	steeps := vmesh.NewHeightmap(h.Mesh)
	dh := h.Downhill()
	for r, d := range dh {
		if d < 0 {
			continue // Skip all sinks.
		}
		// Height difference.
		hDiff := h.Values[r] - h.Values[d]

		// Distance.
		dist := distPoints(m.Vertices[r].X, m.Vertices[r].Y, m.Vertices[d].X, m.Vertices[d].Y)

		// Essentially the angle (0°-90°) in the range of 0.0 to 1.0.
		steeps.Values[r] = math.Atan(hDiff/dist) * 2 / math.Pi
	}
	return steeps
}

func erosionRate2(h *vmesh.Heightmap) *vmesh.Heightmap {
	const maxErosionDistance = 3 * 0.03
	const minExp = 1.0
	const varExp = 3.0

	m := h.Mesh
	steeps := getSteepness(h)
	toE := vmesh.NewHeightmap(m)
	flux := getFlux(h)
	for r, fl := range flux.Values {
		// Exponent (e):
		//
		// The exponent determines the shape of the eroded ravine.
		//
		// y = pow(|x|, e)
		//         _  _
		// e = 1:   \/
		//         _   _
		// e = 2:   \_/
		//         _   _
		// e = 4+:  |_|
		exponent := minExp + varExp*(1.0-steeps.Values[r])

		// The amount of flux determines how wide the resulting eroded path is.
		// Since rivers slow down the wider they are, they start to be less erosive with higher flux.
		// Calculate maximum erosion distance based on the water flow intensity aka flux.
		maxDist := fl * maxErosionDistance // Flux times 3 times the current average neighbor distance (at 16384 points)

		// Assign the erosion intensity.
		if toE.Values[r] < maxDist {
			toE.Values[r] = maxDist
			// TODO: sharp drops should carve with higher intensity.
			// toE.Values[r] *= (1 + steeps.Values[r]) / 2
		}
		rVertex := m.Vertices[r]
		seen := make(map[int]bool)
		var doErode func(reg int)
		doErode = func(reg int) {
			if seen[reg] {
				return
			}
			seen[reg] = true
			for _, nb := range h.Neighbours(reg) {
				// Calculate distance.
				dVertex := m.Vertices[nb]
				dist := distPoints(rVertex.X, rVertex.Y, dVertex.X, dVertex.Y)
				if dist > maxDist {
					continue // Skip everything that is too far away.
				}
				distRes := dist / maxDist
				erode := math.Pow(distRes, exponent)
				toErode := maxDist - erode
				if toE.Values[nb] < toErode {
					toE.Values[nb] = toErode
				}
				doErode(nb)
			}
		}
		doErode(r)
	}
	return toE
}

/*
func erosionRate2(h *vmesh.Heightmap) *vmesh.Heightmap {
	const maxErosionDistance = 3 * 0.03
	const minExp = 1.0
	const varExp = 3.0

	m := h.Mesh
	steeps := getSteepness(h)
	toE := vmesh.NewHeightmap(m)
	dh := h.Downhill()
	flux := getFlux(h)
	for r, fl := range flux.Values {
		// Exponent (e):
		//
		// The exponent determines the shape of the eroded ravine.
		//
		// y = pow(|x|, e)
		//         _  _
		// e = 1:   \/
		//         _   _
		// e = 2:   \_/
		//         _   _
		// e = 4+:  |_|
		exponent := minExp + varExp*(1.0-steeps.Values[r])

		// The amount of flux determines how wide the resulting eroded path is.
		// Since rivers slow down the wider they are, they start to be less erosive with higher flux.
		// Calculate maximum erosion distance based on the water flow intensity aka flux.
		maxDist := fl * maxErosionDistance // Flux times 3 times the current average neighbor distance (at 16384 points)

		// Assign the erosion intensity.
		if toE.Values[r] < maxDist {
			toE.Values[r] = maxDist
			// TODO: sharp drops should carve with higher intensity.
			// toE.Values[r] *= (1 + steeps.Values[r]) / 2
		}
		rVertex := m.Vertices[r]

		// MEANDERING
		rDh := dh[r]
		var dhVec [2]float64
		// Calculate downhill vector if there is one.
		if rDh >= 0 {
			dhVec = getVecOf2Vxs(rVertex, m.Vertices[rDh])
		}

		seen := make(map[int]bool)
		var doErode func(reg int)
		doErode = func(reg int) {
			if seen[reg] {
				return
			}
			seen[reg] = true
			for _, nb := range h.Neighbours(reg) {
				// Calculate distance.
				dVertex := m.Vertices[nb]

				mul := 1.0
				// MEANDERING
				if nb != rDh && nb != r {
					// If nb is in direction of vec r->rDh
					// extend erosion distance.
					// Calculate downhill neighbor vector.
					dnVec := getVecOf2Vxs(rVertex, dVertex)

					// Calculate how much the vectors differ.
					v := dot2Vx(dhVec, dnVec)
					if v > 0 {
						log.Println(v)
						mul = 1.8
					} else {
						mul = 0.7
					}
				}

				dist := distPoints(rVertex.X, rVertex.Y, dVertex.X, dVertex.Y) / mul
				log.Println("res1", dist)
				if dist > maxDist {
					continue // Skip everything that is too far away.
				}
				distRes := dist / maxDist
				log.Println("nores", distRes)
				erode := math.Pow(distRes, exponent)
				toErode := maxDist - erode
				if toE.Values[nb] < toErode {
					toE.Values[nb] = toErode
				}
				doErode(nb)
			}
		}
		doErode(r)
	}
	return toE
}

func getVecOf2Vxs(rVertex, dVertex voronoi.Vertex) [2]float64 {
	dx := dVertex.X - rVertex.X
	dy := dVertex.Y - rVertex.Y
	c := math.Sqrt(dx*dx + dy*dy)
	return [2]float64{dx / c, dy / c}
}

func getVecOf2Points(i, j int, m *vmesh.Mesh) [2]float64 {
	iVx := m.Vertices[i]
	jVx := m.Vertices[j]
	x1 := iVx.X - jVx.X
	y1 := iVx.Y - jVx.Y
	eee1 := math.Sqrt(x1*x1 + y1*y1)
	return [2]float64{x1 / eee1, y1 / eee1}
}

func dot2(x1, y1, x2, y2 float64) float64 {
	return x1*x2 + y1*y2
}

func dot2Vx(v1, v2 [2]float64) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}*/
