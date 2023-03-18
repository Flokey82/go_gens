package main

import (
	"github.com/Flokey82/go_gens/gencitymap"
)

func main() {
	m := gencitymap.NewMap(123)
	m.Generate()
	for i := 0; i < 1940; i++ {
		m.Step()
	}

	m.ExportToPNG("test.png")
}
