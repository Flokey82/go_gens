package main

import (
	"github.com/Flokey82/go_gens/genworldvoronoi"
)

func main() {
	sp := genworldvoronoi.NewMap(1234, 25, 200000, 0.0)
	sp.ExportPng("test.png")
	sp.ExportOBJ("test.obj")
	//sp.ExportSVG("test.svg")
}
