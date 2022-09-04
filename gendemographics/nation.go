package gendemographics

import (
	"fmt"
	"log"
)

// Nation represents a nation.
type Nation struct {
	Size        int           // geographic size in square miles
	Density     int           // population density in people per square mile
	Settlements []*Settlement // settlements within the nation
}

// Various constants.
const (
	PopulationDensityLow        = 20  // population per sqmile
	PopulationDensityMedium     = 40  // population per sqmile
	PopulationPerSqmAgriculture = 120 // 1 sqmile can feed 120 people
)

// NewNation returns a new nation with the given square milage and population density.
func NewNation(size, density int) *Nation {
	return &Nation{
		Size:    size,
		Density: density,
	}
}

// Log prints information on the nation to the console.
func (n *Nation) Log() {
	log.Println(fmt.Sprintf("Population: %d", n.Population()))
	for _, s := range n.Settlements {
		s.Log()
	}
}

// Population returns the population given the population density.
func (n *Nation) Population() int {
	return n.Size * n.Density
}

// Agriculture returns the number of square miles required to feed the population.
func (n *Nation) Agriculture() int {
	return n.Population() / PopulationPerSqmAgriculture
}
