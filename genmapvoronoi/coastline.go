package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/vmesh"
	"github.com/pzsz/voronoi"
)

func contour(h *vmesh.Heightmap, level float64) [][]voronoi.Vertex {
	var edges [][2]voronoi.Vertex
	for i := 0; i < len(h.Edges); i++ {
		e := h.Edges[i]
		if e.Right == nil {
			continue
		}
		//if isnearedge(h.Mesh, e.IdxA) || isnearedge(h.Mesh, e.IdxB){
		// continue
		//}
		if (h.Values[e.IdxA] > level && h.Values[e.IdxB] <= level) ||
			(h.Values[e.IdxB] > level && h.Values[e.IdxA] <= level) {
			edges = append(edges, [2]voronoi.Vertex{e.Left.Site, e.Right.Site})
		}
	}

	return mergeSegments(edges)
}
