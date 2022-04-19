package main

import (
	"github.com/Flokey82/go_gens/genworldvoronoi"
	"log"
)

func main() {
	sp, err := genworldvoronoi.NewMap(1234, 25, 400000, 0.0)
	if err != nil {
		log.Fatal(err)
	}
	exportPNG := true
	exportOBJ := true
	exportSVG := true
	if exportPNG {
		sp.ExportPng("test.png")
	}
	if exportOBJ {
		sp.ExportOBJ("test.obj")
	}
	if exportSVG {
		sp.ExportSVG("test.svg")
	}
}
