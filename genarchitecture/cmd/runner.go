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

	// Set up the sample rules.
	rc := genarchitecture.NewRuleCollection()
	for _, r := range genarchitecture.SampleRules {
		rc.AddRule(r)
	}
	rc.Root = genarchitecture.SampleRules[0]

	// Run the rules.
	mesh := rc.Run()
	mesh.ExportToObj("test_3.obj")
}
