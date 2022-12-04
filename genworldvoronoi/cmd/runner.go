package main

import (
	"log"

	"github.com/Flokey82/go_gens/genworldvoronoi"
)

func main() {
	sp, err := genworldvoronoi.NewMap(1234, 25, 400000, 0.0)
	if err != nil {
		log.Fatal(err)
	}
	sp.GetEmpires()
	exportPNG := true
	exportOBJ := true
	exportSVG := true
	exportWebp := true
	if exportPNG {
		sp.ExportPng("test.png")
	}
	if exportOBJ {
		sp.ExportOBJ("test.obj")
	}
	if exportSVG {
		sp.ExportSVG("test.svg")
	}
	if exportWebp {
		sp.ExportWebp("test.webp")
	}
}
