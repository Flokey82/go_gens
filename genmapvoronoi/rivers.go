package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/vmesh"
	"github.com/pzsz/voronoi"
)

func getRiverPaths(h *vmesh.Heightmap, limit float64) [][]voronoi.Vertex {
	dh := h.Downhill()
	flux := getFlux(h)

	var links [][2]voronoi.Vertex
	var above int
	for i := 0; i < h.Len(); i++ {
		if h.Values[i] > 0 {
			above++
		}
	}
	limit *= float64(above) / float64(h.Len())
	for i := 0; i < len(dh); i++ {
		//if isnearedge(h.Mesh, i){
		//    continue
		//}
		if flux.Values[i] > limit && h.Values[i] > 0 && dh[i] >= 0 {
			up := h.Vertices[i]
			down := h.Vertices[dh[i]]
			if h.Values[dh[i]] > 0 {
				links = append(links, [2]voronoi.Vertex{up, down})
			} else {
				links = append(links, [2]voronoi.Vertex{up, {X: (up.X + down.X) / 2, Y: (up.Y + down.Y) / 2}})
			}
		}
	}

	mergedSegs := mergeSegments(links)
	for i := range mergedSegs {
		mergedSegs[i] = relaxPath(mergedSegs[i])
	}
	// return mergeSegments(links).map(relaxPath)
	return mergedSegs
}

func getRivers(h *vmesh.Heightmap, limit float64) []int {
	rivers := make([]int, h.Len())

	// Set up defaults.
	for i := range rivers {
		rivers[i] = -1 // -1 means no river
	}

	dh := h.Downhill()
	flux := getFlux(h)

	var links [][2]int
	var above int
	for i := 0; i < h.Len(); i++ {
		if h.Values[i] > 0 {
			above++
		}
	}
	limit *= float64(above) / float64(h.Len())
	for i := 0; i < len(dh); i++ {
		//if isnearedge(h.Mesh, i){
		//    continue
		//}
		if flux.Values[i] > limit && h.Values[i] > 0 && dh[i] >= 0 {
			up := i
			down := dh[i]
			links = append(links, [2]int{up, down})
		}
	}

	mergedSegs := mergeIndexSegments(links)
	for i := range mergedSegs {
		for _, idx := range mergedSegs[i] {
			rivers[idx] = i
		}
	}
	return rivers
}

func mergeIndexSegments(segs [][2]int) [][]int {
	adj := make(map[int][]int)
	for i := 0; i < len(segs); i++ {
		seg := segs[i]
		a0 := adj[seg[0]]
		a1 := adj[seg[1]]
		a0 = append(a0, seg[1])
		a1 = append(a1, seg[0])
		adj[seg[0]] = a0
		adj[seg[1]] = a1
	}
	done := make(map[int]bool)
	var paths [][]int
	var path []int
	for {
		if path == nil {
			for i := 0; i < len(segs); i++ {
				if done[i] {
					continue
				}
				done[i] = true
				path = []int{segs[i][0], segs[i][1]}
				break
			}
			if path == nil {
				break
			}
		}
		var changed bool
		for i := 0; i < len(segs); i++ {
			if done[i] {
				continue
			}
			if len(adj[path[0]]) == 2 && segs[i][0] == path[0] {
				path = unshiftIndexPath(path, segs[i][1])
			} else if len(adj[path[0]]) == 2 && segs[i][1] == path[0] {
				path = unshiftIndexPath(path, segs[i][0])
			} else if len(adj[path[len(path)-1]]) == 2 && segs[i][0] == path[len(path)-1] {
				path = append(path, segs[i][1])
			} else if len(adj[path[len(path)-1]]) == 2 && segs[i][1] == path[len(path)-1] {
				path = append(path, segs[i][0])
			} else {
				continue
			}
			done[i] = true
			changed = true
			break
		}
		if !changed {
			paths = append(paths, path)
			path = nil
		}
	}
	return paths
}

func unshiftIndexPath(path []int, p int) []int {
	return append([]int{p}, path...)
}

func mergeSegments(segs [][2]voronoi.Vertex) [][]voronoi.Vertex {
	adj := make(map[voronoi.Vertex][]voronoi.Vertex)
	for i := 0; i < len(segs); i++ {
		seg := segs[i]
		a0 := adj[seg[0]]
		a1 := adj[seg[1]]
		a0 = append(a0, seg[1])
		a1 = append(a1, seg[0])
		adj[seg[0]] = a0
		adj[seg[1]] = a1
	}
	done := make(map[int]bool)
	var paths [][]voronoi.Vertex
	var path []voronoi.Vertex
	for {
		if path == nil {
			for i := 0; i < len(segs); i++ {
				if done[i] {
					continue
				}
				done[i] = true
				path = []voronoi.Vertex{segs[i][0], segs[i][1]}
				break
			}
			if path == nil {
				break
			}
		}
		var changed bool
		for i := 0; i < len(segs); i++ {
			if done[i] {
				continue
			}
			if len(adj[path[0]]) == 2 && segs[i][0] == path[0] {
				path = unshiftPath(path, segs[i][1])
			} else if len(adj[path[0]]) == 2 && segs[i][1] == path[0] {
				path = unshiftPath(path, segs[i][0])
			} else if len(adj[path[len(path)-1]]) == 2 && segs[i][0] == path[len(path)-1] {
				path = append(path, segs[i][1])
			} else if len(adj[path[len(path)-1]]) == 2 && segs[i][1] == path[len(path)-1] {
				path = append(path, segs[i][0])
			} else {
				continue
			}
			done[i] = true
			changed = true
			break
		}
		if !changed {
			paths = append(paths, path)
			path = nil
		}
	}
	return paths
}

func unshiftPath(path []voronoi.Vertex, p voronoi.Vertex) []voronoi.Vertex {
	return append([]voronoi.Vertex{p}, path...)
}

func relaxPath(path []voronoi.Vertex) []voronoi.Vertex {
	newpath := []voronoi.Vertex{path[0]}
	for i := 1; i < len(path)-1; i++ {
		newpt := voronoi.Vertex{
			0.25*path[i-1].X + 0.5*path[i].X + 0.25*path[i+1].X,
			0.25*path[i-1].Y + 0.5*path[i].Y + 0.25*path[i+1].Y,
		}
		newpath = append(newpath, newpt)
	}
	newpath = append(newpath, path[len(path)-1])
	return newpath
}
