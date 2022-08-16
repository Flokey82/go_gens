package genmapvoronoi

import (
	"math"
	"sort"

	"github.com/Flokey82/go_gens/vmesh"
)

// doErosion performs n erosion passes using the given amount.
func doErosion(h *vmesh.Heightmap, amount float64, n int) (eroded *vmesh.Heightmap, sediment *vmesh.Heightmap) {
	eroded = HeightFillSinks(h)
	for i := 0; i < n; i++ {
		eroded = erode(eroded, amount)
		eroded = HeightFillSinks(eroded)
	}
	return eroded, eroded.Diff(h)
}

// erode performs one single pass of erosion over the heightmap.
func erode(h *vmesh.Heightmap, amount float64) *vmesh.Heightmap {
	er := erosionRate2(h) // NOTE: We currently use the alternative erosion calculation.
	newh := vmesh.NewHeightmap(h.Mesh)
	_, maxr := er.MinMax()
	for i := 0; i < h.Len(); i++ {
		newh.Values[i] = h.Values[i] - amount*(er.Values[i]/maxr)
	}
	return newh
}

// erosionRate calculates the erosion rate for each vertex in the heightmap
// based on its flux and slope.
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

// getFlux calculates the water flux for each vertex in the heightmap.
func getFlux(h *vmesh.Heightmap) *vmesh.Heightmap {
	dh := h.Downhill()

	// Initialize flux and the slice we use to sort the vertices.
	idxs := make([]int, h.Len())
	flux := vmesh.NewHeightmap(h.Mesh)
	for i := 0; i < len(h.Values); i++ {
		idxs[i] = i

		// Initialize the flux with 1.0 divided by the number of vertices, which
		// will result in flux being at most 1.0 for any vertex.
		flux.Values[i] = 1.0 / float64(h.Len())
	}

	// Sort the vertices from highest to lowest.
	sort.Slice(idxs, func(a, b int) bool {
		return h.Values[idxs[a]] > h.Values[idxs[b]]
	})

	// Iterate over all vertices from highest to lowest and
	// add its flux value to the flux value of its downhill neighbour.
	for i := 0; i < h.Len(); i++ {
		j := idxs[i]
		// Skip sinks that don't have a downhill neighbour (indicated by a -1).
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

// getSinks attempts to identify sinks by looking for vertices that are lower
// than all their neighbors and have an inflow of water (by summing up the flux of
// all of their neighbors).
//
// NOTE: This seems to be garbage code.
// Not sure why I wrote this, seems pretty dumb. Better use findSinks instead.
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

// findSinks returns a mapping of vertex to sink index.
//
// We traverse each point downhill until we reach a downhill value of -1
// and assign the index of the lowest vertex as a sinkID, while we ignore "sinks"
// that are the edge of the map.
//
// Not sure how well this works.
func findSinks(h *vmesh.Heightmap) []int {
	dh := h.Downhill()
	sinks := make([]int, h.Len())
	for i := 0; i < len(dh); i++ {

		// Traverse vertex i downhill until we reach
		// either the edge of the map, or a sink.
		node := i
		for {
			// If we have reached the edge of the map,
			// abort. (This is not considered a sink,
			// even if it is lower than all its neighbours.)
			if h.IsEdge(node) {
				sinks[i] = vmesh.VertexOutOfBounds
				break
			}

			// If the downhill value indicates that 'node' is
			// a sink, we assign 'node' as an ID to the vertex
			// 'i' (which we started traversing from).
			//
			// This will allow us to know which vertices will
			// drain into what sink, giving us information on
			// which vertices will be part of what "lakes" (sinks).
			if dh[node] == vmesh.VertexSink {
				sinks[i] = node
				break
			}

			// We haven't reached the end of our traversal,
			// so we continue traversing.
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
