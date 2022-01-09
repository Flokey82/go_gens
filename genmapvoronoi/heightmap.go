package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/vmesh"
	"github.com/pzsz/voronoi"
	"math"
	"math/rand"
)

// Mesh-based heightmap generation helpers.

func MeshSlope(m *vmesh.Mesh, direction [2]float64) *vmesh.Heightmap {
	return m.MapF(func(x voronoi.Vertex) float64 {
		return x.X*direction[0] + x.Y*direction[1]
	})
}

func MeshCone(m *vmesh.Mesh, slope float64) *vmesh.Heightmap {
	return m.MapF(func(x voronoi.Vertex) float64 {
		return math.Pow(x.X*x.X+x.Y*x.Y, 0.5) * slope
	})
}

func MeshVolCone(m *vmesh.Mesh, slope float64) *vmesh.Heightmap {
	return m.MapF(func(x voronoi.Vertex) float64 {
		dist := math.Pow(x.X*x.X+x.Y*x.Y, 0.5)
		if dist < 0.1 {
			return -4 * dist * slope
		}
		return dist * slope
	})
}

func MeshMountains(m *vmesh.Mesh, n int, r float64) *vmesh.Heightmap {
	var mounts []voronoi.Vertex
	for i := 0; i < n; i++ {
		mounts = append(mounts, voronoi.Vertex{m.Extent.Width * (rand.Float64() - 0.5), m.Extent.Height * (rand.Float64() - 0.5)})
	}
	newvals := vmesh.NewHeightmap(m)
	for i := 0; i < len(m.Vertices); i++ {
		p := m.Vertices[i]
		for j := 0; j < n; j++ {
			m := mounts[j]
			newvals.Values[i] += math.Pow(math.Exp(-((p.X-m.X)*(p.X-m.X)+(p.Y-m.Y)*(p.Y-m.Y))/(2*r*r)), 2)
		}
	}
	return newvals
}

// Heightmap modification helpers.

func HeightRelax(h *vmesh.Heightmap) *vmesh.Heightmap {
	newh := vmesh.NewHeightmap(h.Mesh)
	for i := 0; i < h.Len(); i++ {
		nbs := h.Neighbours(i)
		if len(nbs) < 3 {
			newh.Values[i] = 0
			continue
		}
		var heights []float64
		for _, j := range nbs {
			heights = append(heights, h.Values[j])
		}
		newh.Values[i] = CalcMean(heights) //d3.mean(nbs.map(function (j) {return h[j]}));
	}
	return newh
}

func HeightNormalize(h *vmesh.Heightmap) *vmesh.Heightmap {
	lo, hi := h.MinMax()
	return h.MapF(func(x float64) float64 {
		return (x - lo) / (hi - lo)
	})
}

func HeightPeaky(h *vmesh.Heightmap) *vmesh.Heightmap {
	return HeightNormalize(h).MapF(math.Sqrt)
}

func HeightSetSeaLevel(h *vmesh.Heightmap, q float64) *vmesh.Heightmap {
	newh := vmesh.NewHeightmap(h.Mesh)
	min, max := newh.MinMax()
	delta := min + (max-min)*0.1
	//delta := quantile(h, q)
	for i := 0; i < h.Len(); i++ {
		newh.Values[i] = h.Values[i] - delta
	}
	return newh
}

func HeightCleanCoast(h *vmesh.Heightmap, iters int) *vmesh.Heightmap {
	for iter := 0; iter < iters; iter++ {
		changed := 0
		newh := vmesh.NewHeightmap(h.Mesh)
		for i := 0; i < h.Len(); i++ {
			newh.Values[i] = h.Values[i]
			nbs := h.Neighbours(i)
			if h.Values[i] <= 0 || len(nbs) != 3 {
				continue
			}
			count := 0
			best := -999999.0
			for j := 0; j < len(nbs); j++ {
				if h.Values[nbs[j]] > 0 {
					count++
				} else if h.Values[nbs[j]] > best {
					best = h.Values[nbs[j]]
				}
			}
			if count > 1 {
				continue
			}
			newh.Values[i] = best / 2
			changed++
		}
		h = newh
		newh = vmesh.NewHeightmap(h.Mesh)
		for i := 0; i < h.Len(); i++ {
			newh.Values[i] = h.Values[i]
			nbs := h.Neighbours(i)
			if h.Values[i] > 0 || len(nbs) != 3 {
				continue
			}
			count := 0
			best := 999999.0
			for j := 0; j < len(nbs); j++ {
				if h.Values[nbs[j]] <= 0 {
					count++
				} else if h.Values[nbs[j]] < best {
					best = h.Values[nbs[j]]
				}
			}
			if count > 1 {
				continue
			}
			newh.Values[i] = best / 2
			changed++
		}
		h = newh
	}
	return h
}

func HeightFillSinks(h *vmesh.Heightmap) *vmesh.Heightmap { //, epsilon float64
	epsilon := 1e-5
	infinity := 999999.0
	newh := vmesh.NewHeightmap(h.Mesh)
	for i := 0; i < h.Len(); i++ {
		if h.IsNearEdge(i) {
			newh.Values[i] = h.Values[i]
		} else {
			newh.Values[i] = infinity
		}
	}
	for {
		var changed bool
		for i := 0; i < h.Len(); i++ {
			if newh.Values[i] == h.Values[i] {
				continue
			}
			nbs := h.Neighbours(i)
			for j := 0; j < len(nbs); j++ {
				if h.Values[i] >= newh.Values[nbs[j]]+epsilon {
					newh.Values[i] = h.Values[i]
					changed = true
					break
				}
				oh := newh.Values[nbs[j]] + epsilon
				if (newh.Values[i] > oh) && (oh > h.Values[i]) {
					newh.Values[i] = oh
					changed = true
				}
			}
		}
		if !changed {
			return newh
		}
	}
}

func CalcMean(nums []float64) float64 {
	total := 0.0
	for _, v := range nums {
		total += v
	}
	return total / float64(len(nums))
}

func isInList(l []*voronoi.Cell, c *voronoi.Cell) bool {
	for _, ci := range l {
		if c == ci {
			return true
		}
	}
	return false
}

func isInIntList(l []int, c int) bool {
	for _, ci := range l {
		if c == ci {
			return true
		}
	}
	return false
}

func runif(lo, hi float64) float64 {
	return lo + rand.Float64()*(hi-lo)
}

func normalizeVec2(vec [2]float64) [2]float64 {
	var dest [2]float64
	x := vec[0]
	y := vec[1]
	length := (x * x) + (y * y)
	if length > 0 {
		length = 1.0 / math.Sqrt(length)
		dest[0] = vec[0] * length
		dest[1] = vec[1] * length
	}
	return dest
}

func randomVector2(scale float64) [2]float64 {
	return [2]float64{scale * rand.Float64(), scale * rand.Float64()}
}

func crossVector2(v1, v2 [2]float64) float64 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}

func lenVector2(v [2]float64) float64 {
	return float64(math.Hypot(v[0], v[1]))
}
