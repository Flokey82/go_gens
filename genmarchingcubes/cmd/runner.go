package main

import (
	"log"

	"github.com/Flokey82/go_gens/genmarchingcubes"
)

func main() {
	const (
		h, w, d = 2, 2, 2
	)
	values := []float64{
		1, 1, 1, 0, 0, 1, 1, 1,
	}
	tris := genmarchingcubes.MarchingCubesGrid(h, w, d, values, 0.5)
	log.Println(tris)
}
