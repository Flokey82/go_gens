package genworldvoronoi

import (
	"fmt"
	"log"
	"math"

	"github.com/Flokey82/go_gens/genbiome"
)

func (m *Civ) calculateEconomicPotential() {
	// NOTE: This is unfinished right now and a WIP.

	// Calculate the analog of distance between regions by taking the surface
	// of a sphere with radius 1 and dividing it by the number of regions.
	// The square root will work as a somewhat sensible approximation of
	// distance.
	// distRegion := math.Sqrt(4*math.Pi/float64(m.mesh.numRegions))

	// Get the stop regions, which are the cities
	// and calculate the radius in which we can find resources.
	var resourceRadius []float64
	stopRegions := make(map[int]bool)
	for _, c := range m.Cities {
		stopRegions[c.ID] = true
		// The base radius is dependent on the population.
		// The minimum radius is 1.0 and increases with the square
		// root of the population.
		radius := 1.0 + math.Sqrt(float64(c.Population))/5
		resourceRadius = append(resourceRadius, radius)
	}

	// Per resource, we calculate the distance field originating from the
	// cities.
	economicPotential := make([]float64, len(m.Cities))

	calcResourceValues := func(resourceType, resourceMax int) {
		for res := 0; res < resourceMax; res++ {
			// The resource ID also doubles as resource value.
			resourceID := 1 << res

			// Get all regions that contain the resource.
			regions := m.getRegionsWithResource(byte(resourceID), resourceType)
			dist := m.assignDistanceField(regions, stopRegions)

			// Now loop through all cities and check if the distance field
			// indicates that we can find the resource in the radius.
			for i, c := range m.Cities {
				radius := resourceRadius[i]
				// TODO: Make sure we take distance into account.
				if dist[c.ID] <= radius {
					economicPotential[i] += float64(resourceID)
				}
			}
		}
	}

	calcResourceValues(ResourceTypeMetal, ResMaxMetals)
	calcResourceValues(ResourceTypeGem, ResMaxGems)
	calcResourceValues(ResourceTypeStone, ResMaxStones)

	// Now get the agricultural potential of all regions.
	fAr := m.getFitnessArableLand()
	agriculturePerRegion := make([]float64, m.mesh.numRegions)
	for i := range agriculturePerRegion {
		agriculturePerRegion[i] = fAr(i)
	}

	// Get the max value of the agricultural potential.
	_, maxAgr := minMax(agriculturePerRegion)
	for i, c := range m.Cities {
		// Check if we have a positive agricultural potential
		// and add the normalized value to the economic potential.
		if agrPotential := agriculturePerRegion[c.ID]; agrPotential > 0 {
			economicPotential[i] += agrPotential / maxAgr
		}
	}

	// Now we go through all the cities, and see if they might be able to
	// trade with each other. This way they can profit from each other's
	// resources.
	// In the future we make this dependent on geographic features, where
	// mountains or the sea might be a barrier.
	tradePotential := make([]float64, len(m.Cities))
	for i, c := range m.Cities {
		// Calculate the distance field of all cities to the current city.
		dist := m.assignDistanceField([]int{c.ID}, stopRegions)

		// Loop through all cities and check if we can trade with them.
		for j, c2 := range m.Cities {
			// We don't trade with ourselves.
			if i == j {
				continue
			}
			// The trade radius is the sum of the square of the two cities' radius.
			radius := resourceRadius[i]*(1+economicPotential[i]) + resourceRadius[j]*(1+economicPotential[j])

			// If the distance is within the radius, we can trade.
			// the closer we are, the more economic potential we have.
			if dist[c2.ID] <= radius {
				tradePotential[i] = economicPotential[j] * (1 - dist[c2.ID]/radius)
			}
		}
	}

	// Now we add the trade potential to the economic potential.
	for i := range m.Cities {
		economicPotential[i] += tradePotential[i]
	}

	// Log the economic potential of the cities.
	for i, c := range m.Cities {
		log.Printf("City %s has economic potential %f", c.Name, economicPotential[i])
	}
}

func (m *Civ) TickCity(c *City) {
	m.resetRand()

	// TODO: Recalculate the economic power of the city.
	// This depends on the population, the local resources in the region,
	// distance to other cities, etc.
	// The population determines the radius in which we can find resources.
	// If the resources are contested, the economic value of the resources
	// will be lower.
	// I think it would make sense to calculate the economic power per region
	// for the entire map first.

	// TODO: Move this to a separate function since we need to tick
	// population growth, production, consumption, etc. independently.

	// Grow city population.
	// There is a 0.19% growth rate per year.
	// So the probability of growth is 0.19% * 1 day / 365 days.
	// TODO: The growth rate should also be based on the relative wealth
	// of the city and the current population. Fix this!
	factor := float64(c.Population) * 0.19 / 365
	if factor >= 1 {
		c.Population += int(math.Ceil(factor * m.rand.Float64()))
	} else if m.rand.Float64() < factor {
		c.Population++
	}

	// TODO: If the "sustainability" of the city is lower than needed to
	// sustain the population, the population will decrease.
	// This can be caused by war, disease, disaster, etc...
	// Some people might migrate to other cities that have a higher
	// sustainability.
}

// City represents a city in the world.
type City struct {
	ID         int       // Region where the city is located
	Name       string    // Name of the city
	Type       TownType  // Type of city
	Score      float64   // Score of the fitness function
	Population int       // Population of the city
	Culture    *Culture  // Culture of the city region
	Language   *Language // Language of the city
}

// String returns a string representation of the city.
func (c *City) String() string {
	return fmt.Sprintf("%s (%d)", c.Name, c.Population)
}

// PlaceNCities places n cities with the highest fitness scores.
func (m *Civ) PlaceNCities(n int, cType TownType) {
	// The fitness function, returning a score from 0.0 to 1.0 for a given region.
	// Select the fitness function based on the city type.
	scoreFunc := cType.GetFitnessFunction(m)

	// The distance seed point function, returning seed points/regions that we
	// want to be far away from.
	// For now we just maximize the distance to cities of the same type.
	distSeedFunc := cType.GetDistanceSeedFunc(m)

	// Place n cities of the given type.
	for i := 0; i < n; i++ {
		c := m.PlaceCity(cType, scoreFunc, distSeedFunc)
		log.Printf("placing %s city %d: %s", cType, i, c.String())
	}
}

// PlaceCity places another city at the region with the highest fitness score.
func (m *Civ) PlaceCity(cType TownType, scoreFunc func(int) float64, distSeedFunc func() []int) *City {
	// Pick the region with the highest fitness score.
	var newcity int
	lastMax := math.Inf(-1)
	for i, val := range m.CalcCityScore(scoreFunc, distSeedFunc) {
		if val > lastMax {
			newcity = i
			lastMax = val
		}
	}

	// Get base population from city type.
	// TODO: Calculate population based on suitability for habitation.
	basePop := cType.FoundingPopulation()
	basePop += 2 * m.rand.Intn(basePop) / (len(m.Cities) + 1)
	c := &City{
		ID:         newcity,
		Score:      lastMax,
		Population: basePop,
		Type:       cType,
		Culture:    m.GetCulture(newcity),
	}

	// If there is no known culture, generate a new one.
	if c.Culture == nil {
		c.Culture = m.PlaceCultureAt(newcity) // TODO: Grow this culture.
	}

	// Use the local language to generate a new city name.
	c.Language = c.Culture.Language
	c.Name = c.Language.MakeCityName()
	m.Cities = append(m.Cities, c)
	return c
}

// CalcCityScore calculates the fitness value for settlements for all regions.
//
// 'sf': Fitness function for scoring a region.
// 'distSeedFunc': Returns a number of regions from which we maximize the distance.
func (m *Civ) CalcCityScore(sf func(int) float64, distSeedFunc func() []int) []float64 {
	sfCity := func(r int) float64 {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.Elevation[r] <= 0 || m.Waterpool[r] > 0 {
			return -1.0
		}
		return sf(r)
	}

	return m.CalcFitnessScore(sfCity, distSeedFunc)
}

func (m *Civ) getFitnessTradingTowns() func(int) float64 {
	// TODO: Fix this.
	// I think this function should avoid the penalty wrt.
	// proximity to towns of other types.
	_, connecting := m.getTradeRoutes()
	return func(r int) float64 {
		return float64(len(connecting[r]))
	}
}

func (m *Civ) getFitnessCityDefault() func(int) float64 {
	_, maxFlux := minMax(m.Flux)
	steepness := m.GetSteepness()

	return func(i int) float64 {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.Elevation[i] <= 0 || m.Waterpool[i] > 0 {
			return -1.0
		}

		// Visit all neighbors and modify the score based on their properties.
		var hasWaterBodyBonus bool
		nbs := m.GetRegionNeighbors(i)

		// Initialize fitness score with the normalized flux value.
		// This will favor placing cities along (and at the end of)
		// large rivers.
		score := math.Sqrt(m.Flux[i] / maxFlux)
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

// TownType represents the type of a city.
type TownType string

// The different types of cities.
const (
	TownTypeDefault     TownType = "town"
	TownTypeTrading     TownType = "trading"
	TownTypeMining      TownType = "mining"
	TownTypeFarming     TownType = "agricultural"
	TownTypeDesertOasis TownType = "desert oasis"
)

// FoundingPopulation returns the starting population of a city type.
func (t TownType) FoundingPopulation() int {
	switch t {
	case TownTypeDefault:
		return 100
	case TownTypeTrading:
		return 80
	case TownTypeMining:
		return 20
	case TownTypeFarming:
		return 20
	case TownTypeDesertOasis:
		return 20
	default:
		log.Fatalf("unknown city type: %s", t)
	}
	return 0
}

// GetDistanceSeedFunc returns the distance seed function for a city type.
func (t TownType) GetDistanceSeedFunc(m *Civ) func() []int {
	// For now we just maximize the distance to cities of the same type.
	return func() []int {
		var cities []int
		for _, c := range m.Cities {
			if c.Type == t {
				cities = append(cities, c.ID)
			}
		}
		return cities
	}
}

// GetFitnessFunction returns the fitness function for a city type.
func (t TownType) GetFitnessFunction(m *Civ) func(int) float64 {
	// TODO: Create different fitness functions for different types of settlement.
	//   - Capital
	//   - Cities / Settlements
	//     ) Proximity to capital!
	//   - Agricultural
	//   - Mining
	//   - ...
	switch t {
	case TownTypeDefault:
		fa := m.getFitnessClimate()
		fb := m.getFitnessCityDefault()
		return func(r int) float64 {
			return fa(r) * fb(r)
		}
	case TownTypeTrading:
		return m.getFitnessTradingTowns()
	case TownTypeMining:
		return m.getFitnessSteepMountains()
	case TownTypeFarming:
		return m.getFitnessArableLand()
	case TownTypeDesertOasis:
		// TODO: Improve this fitness function.
		// Right now the oasis are placed at the very edges of
		// deserts, as there is the "best" climate.
		// However, we want them to be trade hubs for desert
		// crossings... so we'll need to place them in the middle
		// of deserts instead.
		fa := m.getFitnessClimate()
		bf := m.getRWhittakerModBiomeFunc()
		return func(r int) float64 {
			biome := bf(r)
			if biome == genbiome.WhittakerModBiomeColdDesert ||
				biome == genbiome.WhittakerModBiomeSubtropicalDesert {
				return fa(r)
			}
			return 0
		}
	default:
		log.Fatalf("unknown city type: %s", t)
	}
	return nil
}
