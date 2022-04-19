package genworldvoronoi

import (
	"log"
)

func (m *Map) getRivers(limit float64) [][]int {
	dh := m.r_downhill
	flux := m.r_flux

	var links [][2]int
	var above int
	for _, r := range m.r_elevation {
		if r > 0 {
			above++
		}
	}
	limit *= float64(above) / float64(m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		//if isnearedge(h.Mesh, i){
		//    continue
		//}
		if flux[r] > limit && m.r_elevation[r] > 0 && dh[r] >= 0 {
			up := r
			down := dh[r]
			links = append(links, [2]int{up, down})
		}
	}
	log.Println("start merge")

	mergedSegs := mergeIndexSegments(links)
	/*
		// Set up defaults.
		rivers := make([]int, m.mesh.numRegions)
		for i := range rivers {
			rivers[i] = -1 // -1 means no river
		}
		for i := range mergedSegs {
			for _, idx := range mergedSegs[i] {
				rivers[idx] = i
			}
		}
		return rivers
	*/
	return mergedSegs
}

func mergeIndexSegments(segs [][2]int) [][]int {
	log.Println("start adj")
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
	log.Println("start paths")
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
			log.Println("done paths", len(paths), "pathlen", len(path))
			paths = append(paths, path)
			path = nil
		}
	}
	return paths
}

func unshiftIndexPath(path []int, p int) []int {
	return append([]int{p}, path...)
}
