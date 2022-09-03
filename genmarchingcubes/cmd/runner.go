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
	for x := range wld.Values {
		for y := range wld.Values[x] {
			for z, val := range wld.Values[x][y] {
				values[x+y*w+z*w*h] = val
			}
		}
	}
	tris := genmarchingcubes.MarchingCubesGrid(h, w, d, values, 0.5)
	if err := genmarchingcubes.ExportToOBJ("tmp.obj", tris); err != nil {
		log.Println(err)
	}
}
