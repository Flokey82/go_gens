package genmapvoronoi

import (
	"github.com/pzsz/voronoi"
)

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
	mergedSegs := mergeSegments(edges)
	for i := range mergedSegs {
		mergedSegs[i] = relaxPath(mergedSegs[i])
	}
	//return mergeSegments(edges).map(relaxPath);
	return mergedSegs
}

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
	mergedSegs := mergeSegments(edges)
	for i := range mergedSegs {
		mergedSegs[i] = relaxPath(mergedSegs[i])
	}
	//return mergeSegments(edges).map(relaxPath);
	return mergedSegs
}
