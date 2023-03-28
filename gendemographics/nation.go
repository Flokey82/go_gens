package gendemographics

import (
	"log"
	"math"
)

// Nation represents a nation.
type Nation struct {
	Size        int           // geographic size in square miles
	Density     int           // population density in people per square mile
	Settlements []*Settlement // settlements within the nation
}

// Various constants.
const (
	PopulationDensityLow        = 20    // population per sqmile
	PopulationDensityMedium     = 40    // population per sqmile
	PopulationPerSqmAgriculture = 120   // 1 sqmile can feed 120 people
	PopulationPerCastle         = 50000 // 1 castle per 50k people
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
	log.Printf("Population: %d\n", n.Population())
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

// Castles returns the number of active castles in the nation.
//
// 75% of all castles will be in the civilized (settled) areas of a kingdom.
// The other 25% will be in the “wilderness,” along borders, etc.
//
// TODO: There should be different types / sizes of fortifications.
func (n *Nation) Castles() int {
	return n.Population() / PopulationPerCastle
}

// RuinedCastles returns the number of abandoned castles in the nation.
//
// From Medieval Demographics Made Easy by S. John Ross:
//
// Ruins, first of all, depend on the age of the region. The following formula
// is only a guide. The frequency of ruins in Europe varied greatly depending
// on military history and remoteness of the area.
//
// To determine the approximate number of ruined fortifications, divide the
// kingdom’s population by five million. Multiply the result by the square root
// of the kingdom’s age. If the kingdom has changed hands a lot, use the total age
// – the number of years that castle-building people have lived there, regardless
// of the Royal Lineage.
//
// TODO: Maybe that should be a result of a simulation. If a fortification doesn't
// make sense anymore, it can be abandoned and fall into disrepair.
// This might also happen if there is a major recession etc.
func (n *Nation) RuinedCastles(age int) int {
	pop := float64(n.Population())
	return int((pop / 5000000) * math.Sqrt(float64(age)))
}
