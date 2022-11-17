package genworldvoronoi

import (
	"log"
	"math"

	"github.com/Flokey82/go_gens/genbiome"
)

const (
	TownTypeDefault     = "town"
	TownTypeTrading     = "trading"
	TownTypeMining      = "mining"
	TownTypeFarming     = "agricultural"
	TownTypeDesertOasis = "desert oasis"
)

type City struct {
	ID       int     // Region where the city is located
	Name     string  // TODO: Fill locally.
	Type     string  // Type of city
	Score    float64 // Score of the fitness function
	Culture  *Culture
	Language *Language
}

// rPlaceNCities places n cities with the highest fitness scores.
func (m *Map) rPlaceNCities(n int, cType string) {
	// The fitness function, returning a score from
	// 0.0 to 1.0 for a given region.
	var sf func(int) float64

	// The distance seed point function, returning
	// seed points/regions that we want to be far
	// away from.
	var dsf func() []int

	// Select the fitness function based on the
	// city type.
	switch cType {
	case TownTypeDefault:
		fa := m.getFitnessClimate()
		fb := m.getFitnessCityDefault()
		sf = func(r int) float64 {
			return fa(r) * fb(r)
		}
	case TownTypeTrading:
		sf = m.getFitnessTradingTowns()
	case TownTypeMining:
		sf = m.getFitnessSteepMountains()
	case TownTypeFarming:
		sf = m.getFitnessArableLand()
	case TownTypeDesertOasis:
		// TODO: Improve this fitness function.
		// Right now the oasis are placed at the very edges of
		// deserts, as there is the "best" climate.
		// However, we want them to be trade hubs for desert
		// crossings... so we'll need to place them in the middle
		// of deserts instead.
		fa := m.getFitnessClimate()
		bf := m.getRWhittakerModBiomeFunc()
		sf = func(r int) float64 {
			biome := bf(r)
			if biome == genbiome.WhittakerModBiomeColdDesert ||
				biome == genbiome.WhittakerModBiomeSubtropicalDesert {
				return fa(r)
			}
			return 0
		}
	default:
		return
	}

	// For now we just maximize the distance to cities of the same type.
	dsf = func() []int {
		var cities []int
		for _, c := range m.cities_r {
			if c.Type == cType {
				cities = append(cities, c.ID)
			}
		}
		return cities
	}

	// Place n cities of the given type.
	for i := 0; i < n; i++ {
		log.Println("placing "+cType+" city", i)
		m.rPlaceCity(cType, sf, dsf)
	}
}

// rPlaceCity places another city at the region with the highest fitness score.
func (m *Map) rPlaceCity(cType string, sf func(int) float64, distSeedFunc func() []int) {
	var newcity int
	lastMax := math.Inf(-1)
	for i, val := range m.rCityScore(sf, distSeedFunc) {
		if val > lastMax {
			newcity = i
			lastMax = val
		}
	}
	// TODO: Name city using local culture.
	// If there is no local culture, generate one?
	c := &City{
		ID:      newcity,
		Score:   lastMax,
		Type:    cType,
		Culture: m.getRCulture(newcity),
	}

	// If there is no known culture, generate a new one.
	if c.Culture == nil {
		// TODO: Deduplicate with civCultures.go
		// TODO: Grow this culture.
		newCult := m.newCulture(newcity, m.getRCultureTypeFunc())
		newCult.Regions = []int{newcity}
		newCult.Stats = m.getStats(newCult.Regions)
		m.cultures_r = append(m.cultures_r, newCult)
		m.r_cultures[newcity] = newcity
		c.Culture = newCult
	}

	// Use the local language to generate a new city name.
	c.Language = c.Culture.Language
	c.Name = c.Language.MakeCityName()
	m.cities_r = append(m.cities_r, c)
}

// rCityScore calculates the fitness value for settlements for all regions.
// distSeedFunc returns a number of regions from which we maximize the distance when
// calculating the fitness score.
func (m *Map) rCityScore(sf func(int) float64, distSeedFunc func() []int) []float64 {
	// TODO: Create different fitness functions for different types of settlement.
	//   - Capital
	//   - Cities / Settlements
	//     ) Proximity to capital!
	//   - Agricultural
	//   - Mining
	//   - ...

	score := make([]float64, m.mesh.numRegions)

	// Get distance to other cities returned by the distSeedFunc.
	r_distance_c := m.assignDistanceField(distSeedFunc(), make(map[int]bool))

	// Get the max distance for normalizing the distance.
	_, maxDistC := minMax(r_distance_c)

	// Calculate the fitness score for each region
	for i := 0; i < m.mesh.numRegions; i++ {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.r_elevation[i] <= 0 || m.r_pool[i] > 0 {
			score[i] = -1.0
			continue
		}
		score[i] = sf(i)

		// Penalty for proximity / bonus for higher distance to other cities.
		//
		// We multiply the score by the distance to other cities, amplifying
		// positive scores.
		//
		// NOTE: Originally this was done with some constant values, which might be better
		// since we are here dependent on the current score we have assigned and cannot
		// recover an initially bad score caused by a low water flux.

		// TODO: Allow for distance to other cities / city types to be a parameter.
		if math.IsInf(r_distance_c[i], 0) {
			continue
		}
		dist := (r_distance_c[i] / maxDistC)
		score[i] *= dist // originally: -= 0.02 / (float64(r_distance_c[i]) + 1e-9)
	}
	return score
}

func (m *Map) getFitnessTradingTowns() func(int) float64 {
	// TODO: Fix this.
	// I think this function should avoid the penalty wrt.
	// proximity to towns of other types.
	_, connecting := m.getTradeRoutes()
	return func(r int) float64 {
		return float64(len(connecting[r]))
	}
}

// getFitnessClimate returns a fitness function that returns high
// scores for regions with high rainfall and high temperatures.
func (m *Map) getFitnessClimate() func(int) float64 {
	_, maxRain := minMax(m.r_rainfall)
	_, maxElev := minMax(m.r_elevation)

	return func(r int) float64 {
		rTemp := m.getRTemperature(r, maxElev)
		if rTemp < 0 {
			return 0.1
		}
		return 0.1 + 0.9*(m.r_rainfall[r]/maxRain)*math.Sqrt(rTemp/maxTemp)
	}
}

func (m *Map) getFitnessCityDefault() func(int) float64 {
	_, maxFlux := minMax(m.r_flux)
	steepness := m.getRSteepness()

	return func(i int) float64 {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.r_elevation[i] <= 0 || m.r_pool[i] > 0 {
			return -1.0
		}

		// Visit all neighbors and modify the score based on their properties.
		var hasWaterBodyBonus bool
		nbs := m.rNeighbors(i)

		// Initialize fitness score with the normalized flux value.
		// This will favor placing cities along (and at the end of)
		// large rivers.
		score := math.Sqrt(m.r_flux[i] / maxFlux)
		for _, nb := range nbs {
			// Add bonus if near ocean or lake.
			if m.isRBelowOrAtSeaLevelOrPool(nb) {
				// We only apply this bonus once.
				if hasWaterBodyBonus {
					continue
				}
				// If a neighbor is below (or at) sea level, or a lake,
				// we increase the fitness value and reduce it by a fraction,
				// depending on the size of the lake or ocean it is part of.
				//
				// TODO: Improve this.

				// If nb is part of a waterbody (ocean) or lake, we reduce the score by a constant factor.
				// The larger the waterbody/lake, the smaller the penalty, which will favor larger waterbodies.
				if wbSize := m.getRLakeOrWaterBodySize(nb); wbSize > 0 {
					hasWaterBodyBonus = true
					score += 0.55 * (1 - 1/(float64(wbSize)+1e-9))
				}
			} else {
				// If the sourrounding terrain is flat, we get a bonus.
				stp := steepness[nb]
				score += 0.5 * (1.0 - stp*stp) / float64(len(nbs))
			}

			// TODO:
			// - Consider biome
			// - Consider sediment/fertility of land.
			// - Add bonus for mountain proximity (mines, resources)
		}

		// The steeper the terrain, the less likely it is to be settled.
		// TODO: Bonus for trade routes.
		stp := steepness[i]
		score *= 1.0 - (stp * stp)
		return score
	}
}
