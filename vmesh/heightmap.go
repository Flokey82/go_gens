package vmesh

type Heightmap struct {
	*Mesh
	Values        []float64 // elevation data
	//pool     []float64 // water pool data
	//downhill []int     // map of vertex index to next-lower vertex. WARNING THIS IS NOT NOT AUTO UPDATED.
}

func (h *Heightmap) Len() int {
	return len(h.Values)
}

func NewHeightmap(m *Mesh) *Heightmap {
	return &Heightmap{
		Mesh: m,
		Values:    make([]float64, len(m.Vertices)),
	}
}

func (h *Heightmap) Add(hms ...*Heightmap) {
	n := h.Len()
	for i := 0; i < n; i++ {
		for j := 0; j < len(hms); j++ {
			h.Values[i] += hms[j].Values[i]
		}
	}
}

func (h *Heightmap) MinMax() (min, max float64) {
	min = h.Values[0]
	max = h.Values[0]
	for _, v := range h.Values {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return min, max
}

func (h *Heightmap) TriSlope(i int) [2]float64 {
	nbs := h.Neighbours(i)
	if len(nbs) != 3 {
		return [2]float64{0, 0}
	}

	p0 := h.Mesh.Vertices[nbs[0]]
	p1 := h.Mesh.Vertices[nbs[1]]
	p2 := h.Mesh.Vertices[nbs[2]]

	x1 := p1.X - p0.X
	x2 := p2.X - p0.X
	y1 := p1.Y - p0.Y
	y2 := p2.Y - p0.Y

	det := x1*y2 - x2*y1
	h1 := h.Values[nbs[1]] - h.Values[nbs[0]]
	h2 := h.Values[nbs[2]] - h.Values[nbs[0]]

	return [2]float64{
		(y2*h1 - y1*h2) / det,
		(-x2*h1 + x1*h2) / det,
	}
}

		/*
func (h *Heightmap) Slope(i int) float64 {
	s := h.TriSlope(i)
	return math.Sqrt(s[0]*s[0] + s[1]*s[1])
}

func (h *Heightmap) Slopes() *Heightmap {
	//dh := downhill(h)
	slope := NewHeightmap(h.Mesh)
	for i := 0; i < h.Len(); i++ {
		slope.Values[i] = h.Slope(i)
		continue
	}
	return slope
}*/

const (
	VertexOutOfBounds = -2
	VertexSink = -1
)

func (h *Heightmap) Downhill() []int {
	//if h.downhill != nil {
	//	return h.downhill
	//}
	downfrom := func(i int) int {
		if h.IsEdge(i) {
			return VertexOutOfBounds
		}
		best := VertexSink
		besth := h.Values[i]
		nbs := h.Neighbours(i)
		for j := 0; j < len(nbs); j++ {
			if h.Values[nbs[j]] < besth {
				besth = h.Values[nbs[j]]
				best = nbs[j]
			}
		}
		return best
	}
	downs := make([]int, h.Len())
	for i := 0; i < h.Len(); i++ {
		downs[i] = downfrom(i)
	}
	//h.downhill = downs
	return downs
}

func (h *Heightmap) MapF(f func(float64) float64) *Heightmap {
	newh := NewHeightmap(h.Mesh)
	for i, hg := range h.Values {
		newh.Values[i] = f(hg)
	}
	return newh
}