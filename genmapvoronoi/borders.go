package genmapvoronoi

import (
	"github.com/pzsz/voronoi"
)

// getBorders returns the outlines of political territories by identifying neighboring
// territories and returning the voronoi cell sites of all neighbors that do not have
// the same territory ID.
//
// TODO: Dedup this with getCityBorders.
func getBorders(render *Terrain) [][]voronoi.Vertex {
	terr := render.terr
	h := render.h
	var edges [][2]voronoi.Vertex
	for i := 0; i < len(h.Edges); i++ {
		e := h.Edges[i]
		if e.Right == nil {
			continue
		}
		//if isnearedge(h.Mesh, e.IdxA) || isnearedge(h.Mesh, e.IdxB){
		// continue
		//}
		if h.Values[e.IdxA] < 0 || h.Values[e.IdxB] < 0 {
			continue
		}
		if terr[e.IdxA] != terr[e.IdxB] {
			edges = append(edges, [2]voronoi.Vertex{e.Left.Site, e.Right.Site})
		}
	}

	// Merge all segments into continuous outlines of political territories.
	mergedSegs := mergeSegments(edges)
	for i := range mergedSegs {
		// Relax the outlines.
		mergedSegs[i] = relaxPath(mergedSegs[i])
	}
	return mergedSegs
}

// getCityBorders is virtually identical to getBorders but instead of the national
// territories, we return the outlines of the cities.
//
// TODO: Dedup this with getBorders.
func getCityBorders(render *Terrain) [][]voronoi.Vertex {
	terr := render.cityTerritories
	h := render.h
	var edges [][2]voronoi.Vertex
	for i := 0; i < len(h.Edges); i++ {
		e := h.Edges[i]
		if e.Right == nil {
			continue
		}
		//if isnearedge(h.Mesh, e.IdxA) || isnearedge(h.Mesh, e.IdxB){
		// continue
		//}
		if h.Values[e.IdxA] < 0 || h.Values[e.IdxB] < 0 {
			continue
		}
		if terr[e.IdxA] != terr[e.IdxB] {
			edges = append(edges, [2]voronoi.Vertex{e.Left.Site, e.Right.Site})
		}
	}

	// Merge all segments into continuous outlines of city territories.
	mergedSegs := mergeSegments(edges)
	for i := range mergedSegs {
		// Relax the outlines.
		mergedSegs[i] = relaxPath(mergedSegs[i])
	}
	return mergedSegs
}
