package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/vmesh"
	"math"
	"sort"
)

func doErosion(h *vmesh.Heightmap, amount float64, n int) *vmesh.Heightmap {
	h = HeightFillSinks(h)
	for i := 0; i < n; i++ {
		h = erode(h, amount)
		h = HeightFillSinks(h)
	}
	return h
}

func erode(h *vmesh.Heightmap, amount float64) *vmesh.Heightmap {
	er := erosionRate(h)
	newh := vmesh.NewHeightmap(h.Mesh)
	_, maxr := er.MinMax()
	for i := 0; i < h.Len(); i++ {
		newh.Values[i] = h.Values[i] - amount*(er.Values[i]/maxr)
	}
	return newh
}

func erosionRate(h *vmesh.Heightmap) *vmesh.Heightmap {
	flux := getFlux(h)
	slope := getSlope(h)
	newh := vmesh.NewHeightmap(h.Mesh)
	for i := 0; i < h.Len(); i++ {
		river := math.Sqrt(flux.Values[i]) * slope.Values[i]
		creep := slope.Values[i] * slope.Values[i]
		total := 1000*river + creep
		if total > 200 {
			total = 200
		}
		nbs := h.Neighbours(i)
		for _, nb := range nbs {
			newh.Values[nb] += total * 0.25
		}

		newh.Values[i] += total
	}
	return newh
}

func getFlux(h *vmesh.Heightmap) *vmesh.Heightmap {
	dh := h.Downhill()
	idxs := make([]int, h.Len())
	flux := vmesh.NewHeightmap(h.Mesh)
	for i := 0; i < len(h.Values); i++ {
		idxs[i] = i
		flux.Values[i] = 1.0 / float64(h.Len())
	}
	sort.Slice(idxs, func(a, b int) bool {
		return h.Values[idxs[b]]-h.Values[idxs[a]] < 0
	})

	for i := 0; i < h.Len(); i++ {
		j := idxs[i]
		if dh[j] >= 0 {
			flux.Values[dh[j]] += flux.Values[j]
		}
	}
	return flux
}

func getSlope(h *vmesh.Heightmap) *vmesh.Heightmap {
	slope := vmesh.NewHeightmap(h.Mesh)
	for i := 0; i < h.Len(); i++ {
		s := h.TriSlope(i)
		slope.Values[i] = math.Sqrt(s[0]*s[0] + s[1]*s[1])
	}
	return slope
}

func getSinks(h *vmesh.Heightmap) []int {
	dh := h.Downhill()
	sinks := make([]int, h.Len())
	flux := getFlux(h)
	var skip bool
	var sumFlux float64
	for i := 0; i < len(dh); i++ {
		skip = false
		sumFlux = 0.0
		nbs := h.Neighbours(i)
		for _, nb := range nbs {
			if h.Values[nb] < h.Values[i] {
				skip = true
				break
			}
			sumFlux += flux.Values[nb]
		}
		if skip {
			continue
		}
		if sumFlux > 0.001 {
			for _, nb := range nbs {
				sinks[nb] = i
			}
			sinks[i] = i
		}

	}
	return sinks
}

func findSinks(h *vmesh.Heightmap) []int {
	dh := h.Downhill()
	sinks := make([]int, h.Len())
	for i := 0; i < len(dh); i++ {
		node := i
		for {
			if h.IsEdge(node) {
				sinks[i] = vmesh.VertexOutOfBounds
				break
			}
			if dh[node] == vmesh.VertexSink {
				sinks[i] = node
				break
			}
			node = dh[node]
		}
	}
	return sinks
}

/*
func quantile(h*Heightmap, q float64) {
    var sortedh = [];
    for (var i = 0; i < h.length; i++) {
        sortedh[i] = h[i];
    }
    sortedh.sort(d3.ascending);
    return d3.quantile(sortedh, q);
}*/
