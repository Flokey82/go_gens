package gendemographics

import (
	"fmt"
	"log"
	"math"
)

const (
	PopTown    = 1500
	PopVillage = 150
	PopHamlet  = 50
)

// GenSettlementPopulations generates a number of cities, towns, settlements represented as population counts.
// NOTE: This is just a placeholder function and should be reworked.
func GenSettlementPopulations(population int) []int {
	var res []int

	// TODO: Fix distribution.
	capitalPopulation := int(15 * math.Sqrt(float64(population)))
	for population > 100 && capitalPopulation > 100 {
		res = append(res, capitalPopulation)
		population -= capitalPopulation
		capitalPopulation = int(float64(capitalPopulation) * 0.8)
	}

	// The rest would be in small / tiny villages.
	if population > 0 {
		res = append(res, population)
	}
	log.Println(res)
	return res
}

// Settlement represents a settlement.
type Settlement struct {
	Population int
	Businesses map[string]int
}

// NewSettlement returns a new settlement with the given population.
func NewSettlement(population int) *Settlement {
	return &Settlement{
		Population: population,
		Businesses: make(map[string]int),
	}
}

func (s *Settlement) Log() {
	log.Println(fmt.Sprintf("Settlement (%d) %.2f sqm farmland, %d farms", s.Population, s.Farmland(), CalcNumberFarms(s.Population)))
	log.Println(GenBusinesses(s.Population))
}

// Required farm land.
func (s *Settlement) Farmland() float64 {
	return float64(s.Population) / acrePerSqMile
}

// TODO: Settlement types
// - Monastary
// - Burgh

const (
	acrePerSqMile    = 640 // Number of acres in a square mile
	acrePerHousehold = 30  // Number of acres assigned usually to a household for farming
)

func CalcNumberFarms(population int) int {
	return int(math.Ceil(float64(population) / acrePerHousehold))
}
