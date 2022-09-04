package gendemographics

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

// Range represents a range of values.
type Range [2]int

// InRange returns true if the given value is within the range.
func (r Range) InRange(value int) bool {
	return value >= r.Min() && value <= r.Max()
}

// Min returns the minimum value of the range.
func (r Range) Min() int {
	return r[0]
}

// Max returns the maximum value of the range.
func (r Range) Max() int {
	return r[1]
}

// Sizes of settlement types.
var (
	PopRangeCity    = Range{8000, 12000}
	PopRangeTown    = Range{1000, 8000}
	PopRangeVillage = Range{100, 1000}
	PopRangeHamlet  = Range{10, 100}
)

// GenSettlementPopulations generates a number of cities, towns, settlements represented as population counts.
// NOTE: This is just a placeholder function and should be reworked.
func GenSettlementPopulations(population int) []int {
	var res []int

	// TODO: Fix distribution.
	// Determine the population of the largest city in the kingdom.
	// This is equal to (P × M), where P is equal to the square root
	// of the country’s population, and M is equal to a random roll
	// of 2d4+10 (the average roll is 15).
	//
	// NOTE: I decided to double that to steer away from excessive
	// numbers living in individual dwellings, which would be
	// unrealistic.
	capitalPopulation := int(36 * math.Sqrt(float64(population)))
	res = append(res, capitalPopulation)
	population -= capitalPopulation

	// The second-ranking city will be from 20-80% the size of the largest.
	// To randomly determine this, roll 2d4 times 10% (the average result is 50%)
	//
	// NOTE: I decided to use 10% shrinking instead to prevent a too excessive
	// rural population that aren't part of a settlement.
	secondaryCityPopulation := int(float64(capitalPopulation) * (1.0 - 0.1)) // 90%
	res = append(res, secondaryCityPopulation)
	population -= secondaryCityPopulation

	// Each remaining city will be from 10% to 40% smaller than the previous
	// one (2d4 times 5% – the average result is 25%); continue listing cities
	// for as long as the results maintain a city-scaled population (8,000 or more).
	//
	// Cities tend to be from 8,000-12,000 people. A typical large kingdom will have
	// only a few cities in this population range. Centers of scholarly pursuits
	// (the Universities) tend to be in cities of this size, with only the rare exception
	// thriving in a Big City
	//
	// NOTE: This code doesn't really have an upper limit for city size... so if the
	// capital is huge, we might very well exceed the size of 12,000 people.
	prevPopulation := secondaryCityPopulation
	minPopCity := PopRangeCity.Min()
	for {
		// We use a randomized fraction of 8% shrinking.
		pop := int(float64(prevPopulation) * (1.0 - (0.08 * rand.Float64())))
		if pop < minPopCity {
			break
		}
		res = append(res, pop)
		population -= pop
		prevPopulation = pop
	}

	// To determine the number of towns, start with the number of cities,
	// and multiply it by a roll of 2d8 (the average result is 9).
	//
	// Towns range in population from 1,000-8,000 people, with typical values
	// somewhere around 2,500. Culturally, these are the equivalent to the
	// smaller American cities that line the Interstates.
	//
	// Cities and towns tend to have walls only if they’re politically
	// important and/or frequently threatened.
	numTowns := len(res) * 9
	minPopTown := PopRangeTown.Min()
	for remTowns := numTowns; remTowns > 0; remTowns-- {
		// We use a randomized fraction of 3% shrinking.
		pop := int(float64(prevPopulation) * (1.0 - (0.03 * rand.Float64())))
		if pop < minPopTown {
			break
		}
		res = append(res, pop)
		population -= pop
		prevPopulation = pop

	}

	// Now we'd need to determine the number of villages and hamlets.
	// Villages range in size from 100 to 1,000 people, and hamlets
	// range from 10 to 100 people.
	//
	// Villages and hamlets are the smallest settlements that have
	// a name and a population. They are typically located in the
	// countryside, and are the homes of farmers and other rural
	// workers.
	//
	// NOTE: We don't really have a good way to determine the number
	// of villages and hamlets yet. We could use the population of
	// the kingdom and the number of towns and cities to determine
	// the number of villages and hamlets. But we'd need to know the
	// average ratio of number of cities to number of villages and
	// hamlets.
	minPopHamlet := PopRangeHamlet.Min()
	for population > minPopHamlet {
		// We use a randomized fraction of 1% shrinking right now.
		pop := int(float64(prevPopulation) * (1.0 - (0.01 * rand.Float64())))
		if pop < minPopHamlet {
			break
		}
		res = append(res, pop)
		population -= pop
		prevPopulation = pop
	}

	// The rest would be in small / tiny settlements, individual dwellings, etc.
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

// Log prints information on the settlement to the console.
func (s *Settlement) Log() {
	log.Println(fmt.Sprintf("Settlement (%d) %.2f sqm farmland, %d farms", s.Population, s.Farmland(), CalcNumberFarms(s.Population)))
	log.Println(GenBusinesses(s.Population))
}

// Farmland returns the required farm land to feed the population.
func (s *Settlement) Farmland() float64 {
	return float64(s.Population) / acrePerSqMile
}

// TODO: Settlement types
// - Monastary
// - Burgh

// Some conversion constants.
const (
	acrePerSqMile    = 640 // Number of acres in a square mile
	acrePerHousehold = 30  // Number of acres assigned usually to a household for farming
)

// CalcNumberFarms calculates the number of farms required to feed the given population.
func CalcNumberFarms(population int) int {
	return int(math.Ceil(float64(population) / acrePerHousehold))
}
