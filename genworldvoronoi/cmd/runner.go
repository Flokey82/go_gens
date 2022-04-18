package main

import (
	"github.com/Flokey82/go_gens/genworldvoronoi"
	"log"
)

func main() {
	sp, err := genworldvoronoi.NewMap(1234, 25, 200000, 0.0)
	if err != nil {
		log.Fatal(err)
	}
	sp.ExportPng("test.png")
	sp.ExportOBJ("test.obj")
	sp.ExportSVG("test.svg")
}
