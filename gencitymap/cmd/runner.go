package main

import (
	"log"

	"github.com/Flokey82/go_gens/gencitymap"
)

func main() {
	// Create a rules based map.
	m := gencitymap.NewMap(123, gencitymap.DefaultMapConfig)
	m.Generate()
	for i := 0; i < 1940; i++ {
		m.Step()
	}

	// Create a png image.
	if err := m.ExportToPNG("test_rules.png"); err != nil {
		log.Fatal(err)
	}

	// Create a tensor field based map.
	gen, err := gencitymap.TensorTest()
	if err != nil {
		log.Fatal(err)
	}

	// Create a png image.
	if err := gen.ExportToPNG("test_tensor.png"); err != nil {
		log.Fatal(err)
	}

	// Create an svg image.
	if err := gen.ExportToSVG("test_tensor.svg"); err != nil {
		log.Fatal(err)
	}
}
