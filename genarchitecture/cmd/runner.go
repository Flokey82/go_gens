package main

import (
	"log"

	"github.com/Flokey82/go_gens/genarchitecture"
)

func main() {
	st := genarchitecture.GenerateStyle(genarchitecture.Materials)
	st.ExportSvg("test.svg")
	log.Println(st.Description())

	// Generate a sample building.
	genarchitecture.GenerateSampleCathedral()
}
