package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/genheightmap"
	"github.com/Flokey82/go_gens/vectors"
	"github.com/Flokey82/go_gens/vmesh"
	"github.com/pzsz/voronoi"
	"math"
	"math/rand"
)

// Mesh-based heightmap generation helpers.

func MeshSlope(m *vmesh.Mesh, direction vectors.Vec2) *vmesh.Heightmap {
	return m.ApplyGen(genheightmap.GenSlope(direction))
}

func MeshCone(m *vmesh.Mesh, slope float64) *vmesh.Heightmap {
	return m.ApplyGen(genheightmap.GenCone(slope))
}

func MeshVolCone(m *vmesh.Mesh, slope float64) *vmesh.Heightmap {
	return m.ApplyGen(genheightmap.GenVolCone(slope))
}

func MeshMountains(m *vmesh.Mesh, n int, r float64) *vmesh.Heightmap {
	return m.ApplyGen(genheightmap.GenMountains(m.Extent.Width, m.Extent.Height, n, r))
}

func MeshNoise(m *vmesh.Mesh, slope float64) *vmesh.Heightmap {
	return m.ApplyGen(genheightmap.GenNoise(123456, slope))
}

func MeshRidges(m *vmesh.Mesh, direction vectors.Vec2) *vmesh.Heightmap {
	newvals := vmesh.NewHeightmap(m)
	start := rand.Intn(len(newvals.Values))

	childRidgeDist := 5
	childRidgeChanceFraction := 16 // one in n
	childRidgeHeightFactor := 0.8
	childRidgeLifespanFactor := 0.5
	randomWalkChanceFraction := 2
	defaultRidgeLifespan := 400
	defaultRidgeHeight := 2.0

	var drawRidge func(start, lifespan int, maxHeight float64)
	drawRidge = func(start, lifespan int, maxHeight float64) {
		// TODO: With increasing lifespan sine height.
		var length int
		end := vectors.Vec2{
			X: m.Vertices[start].X + direction.X,
			Y: m.Vertices[start].Y + direction.Y,
		}
		for i := start; length < lifespan; length++ {

			newvals.Values[i] = maxHeight * float64(rand.Intn(10)) / 10
			for _, nb := range newvals.Neighbours(i) {
				if distPoints(m.Vertices[nb].X, m.Vertices[nb].Y, end.X, end.Y) < distPoints(m.Vertices[i].X, m.Vertices[i].Y, end.X, end.Y) {
					i = nb
				}
				if rand.Intn(randomWalkChanceFraction) == 0 {
					break
				}
				if rand.Intn(childRidgeChanceFraction) == 0 {
					br := i
					for p := 0; p < childRidgeDist; p++ {
						for _, nb := range newvals.Neighbours(br) {
							if distPoints(m.Vertices[br].X, m.Vertices[br].Y, m.Vertices[i].X, m.Vertices[i].Y) < distPoints(m.Vertices[i].X, m.Vertices[i].Y, m.Vertices[nb].X, m.Vertices[nb].Y) {
								br = nb
							}
						}
					}
					drawRidge(br, int(float64(lifespan)*childRidgeLifespanFactor), maxHeight*childRidgeHeightFactor)
				}
			}
		}
	}
	drawRidge(start, defaultRidgeLifespan, defaultRidgeHeight)
	return newvals
}

func distPoints(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

func MeshHills(m *vmesh.Mesh, n int, r float64) *vmesh.Heightmap {
	var mounts []voronoi.Vertex
	for i := 0; i < n; i++ {
		op := voronoi.Vertex{m.Extent.Width * (rand.Float64() - 0.5), m.Extent.Height * (rand.Float64() - 0.5)}
		nh := rand.Intn(4) + 1
		for j := 0; j < nh; j++ {
			mounts = append(mounts, voronoi.Vertex{op.X + (rand.Float64()-0.5)*r, op.Y + (rand.Float64()-0.5)*r})
		}
	}
	r = r / 20
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
	f := genheightmap.ModRelax(h.Neighbours, func(idx int) float64 {
		return h.Values[idx]
	})
	for i, v := range h.Values {
		newh.Values[i] = f(i, v) //d3.mean(nbs.map(function (j) {return h[j]}));
	}
	return newh
}

func HeightNormalize(h *vmesh.Heightmap) *vmesh.Heightmap {
	min, max := h.MinMax()
	return h.MapF(genheightmap.ModNormalize(min, max))
}

func HeightPeaky(h *vmesh.Heightmap) *vmesh.Heightmap {
	return HeightNormalize(h).MapF(genheightmap.ModPeaky())
}

func HeightSetSeaLevel(h *vmesh.Heightmap, q float64) *vmesh.Heightmap {
	//newh := vmesh.NewHeightmap(h.Mesh)
	//delta := min + (max-min)*0.1
	//delta := quantile(h, q)
	//for i := 0; i < h.Len(); i++ {
	//	newh.Values[i] = h.Values[i] - delta
	//}
	//return newh
	min, max := h.MinMax()
	return h.MapF(genheightmap.ModSeaLevel(min, max, q))
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
