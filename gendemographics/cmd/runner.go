package main

import (
	"github.com/Flokey82/go_gens/gendemographics"
)

func main() {
	c := gendemographics.New()
	n := c.NewNation(41000, gendemographics.PopulationDensityMedium)
	n.Log()
	//log.Println(gendemographics.GenSettlementSizes(100000))
}
