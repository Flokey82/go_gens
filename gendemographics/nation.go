package gendemographics

import (
	"fmt"
	"log"
)

// Nation represents a nation.
type Nation struct {
	Size        int
	Density     int
	Settlements []*Settlement
}

const (
	DensityLow         = 20 // population per sqmile
	DensityMedium      = 40
	defaultAgriculture = 120 // 1 sqmile can feed 120 people
)

// NewNation returns a new nation with the given square milage and population density.
func NewNation(size, density int) *Nation {
	return &Nation{
		Size:    size,
		Density: density,
	}
}

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
	return n.Population() / defaultAgriculture
}
