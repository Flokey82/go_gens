package main

import (
	"log"

	"github.com/Flokey82/go_gens/genmapvoxel"
	"github.com/Flokey82/go_gens/genmarchingcubes"
)

func main() {
	const (
		h, w, d = 32, 32, 32
	)
	wld := genmapvoxel.New(h, w, d, 12345)
	values := make([]float64, h*w*d)
	for x := range wld.Voxels {
		for y := range wld.Voxels[x] {
			for z, val := range wld.Voxels[x][y] {
				if val {
					values[x+y*w+z*w*h] = 1
				}
			}
		}
	}
	tris := genmarchingcubes.MarchingCubesGrid(h, w, d, values, 0.5)
	if err := genmarchingcubes.ExportToOBJ("tmp.obj", tris); err != nil {
		log.Println(err)
	}
}
