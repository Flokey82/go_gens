package main

import (
	"log"

	"github.com/Flokey82/go_gens/genflora"
)

func main() {
	// Generate a number of random plant lineages.
	var lineages []*genflora.PlantLineage
	for i := 0; i < 1000; i++ {
		lin := genflora.GenRandomLineage()
		log.Println(lin.String())
		lineages = append(lineages, lin)
	}
}
