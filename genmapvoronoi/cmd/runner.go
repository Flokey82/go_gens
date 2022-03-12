package main

import (
	"github.com/Flokey82/go_gens/genmapvoronoi"
	"log"
)

func main() {
	r := genmapvoronoi.NewTerrain(genmapvoronoi.DefaultParams)
	if err := r.ExportSVG("test.svg"); err != nil {
		log.Fatal(err)
	}

	if err := r.ExportOBJ("tmp.obj"); err != nil {
		log.Fatal(err)
	}
}
