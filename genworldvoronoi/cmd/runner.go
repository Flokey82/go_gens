package main

import (
	"github.com/Flokey82/go_gens/genworldvoronoi"
)

func main() {
	sp := genworldvoronoi.NewMap(1234, 20, 100000, 0.0)
	sp.ExportOBJ("test.obj")
}
