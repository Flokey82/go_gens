package genworldvoronoi

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/genbiome"
)

func (m *Civ) calculateEconomicPotential() {
	// TODO: Cities should have several values:
	// - Prosperity
	//   trade, industry, agriculture, etc.
	// - Attractiveness
	//   climate, culture, etc.

	// NOTE: This is unfinished right now and a WIP.

	// Calculate the analog of distance between regions by taking the surface
	// of a sphere with radius 1 and dividing it by the number of regions.
	// The square root will work as a somewhat sensible approximation of
	// distance.
	distRegion := math.Sqrt(4 * math.Pi / float64(m.mesh.numRegions))

	// Get the stop regions, which are the cities
	// and calculate the radius in which we can find resources.
	var resourceRadius []float64
	stopRegs := make(map[int]bool)
	for _, c := range m.Cities {
		stopRegs[c.ID] = true
		// The base radius is dependent on the population.
		// The minimum radius is 1.0 and increases with the square
		// root of the population.
		radius := 1.0 + math.Sqrt(float64(c.Population))*distRegion
		resourceRadius = append(resourceRadius, radius)
	}

	// Per resource, we calculate the distance field originating from the
	// cities.
	economicPotential := make([]float64, len(m.Cities))

	calcResourceValues := func(resourceType, resourceMax int) {
		for res := 0; res < resourceMax; res++ {
			// The resource ID also doubles as resource value.
			resourceID := 1 << res

			// Get all regs that contain the resource.
			regs := m.getRegsWithResource(byte(resourceID), resourceType)
			dist := m.assignDistanceField(regs, stopRegs)

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

	// Normalize the economic potential, so we have a value between 0 and 1.
	_, maxEcon := minMax(economicPotential)
	for i := range economicPotential {
		economicPotential[i] /= maxEcon
	}

	// Now get the agricultural potential of all regions.
	fitnessArableFunc := m.getFitnessArableLand()
	agriculturePerRegion := make([]float64, m.mesh.numRegions)
	for i := range agriculturePerRegion {
		agriculturePerRegion[i] = fitnessArableFunc(i)
	}

	// Add the normalized agricultural potential to the economic potential.
	_, maxAgr := minMax(agriculturePerRegion)
	for i, c := range m.Cities {
		// Check if we have a positive agricultural potential
		// and add the normalized value to the economic potential.
		if agrPotential := agriculturePerRegion[c.ID]; agrPotential > 0 {
			economicPotential[i] += agrPotential / maxAgr
		}
	}
	// Now we have a economic potential ranging from 0 to 2.

	// Now we go through all the cities, and see if they might be able to
	// trade with each other. This way they can profit from each other's
	// resources.
	//
	// In the future we make this dependent on geographic features, where
	// mountains or the sea might be a barrier.
	tradePotential := make([]float64, len(m.Cities))
	for i, c := range m.Cities {
		// Calculate the distance field of all cities to the current city.

		//var count int
		// Loop through all cities and check if we can trade with them.
		for j, c2 := range m.Cities {
			// We don't trade with ourselves.
			if i == j {
				continue
			}
			// The trade radius is the sum of the two cities' radius times their economic potential.
			radius := resourceRadius[i]*economicPotential[i] + resourceRadius[j]*economicPotential[j]

			// If the distance is within the radius, we can trade.
			// However, if the other city has a higher economic potential,
			// we profit less from the trade.
			// TODO: Switch this to population size?
			// the closer we are, the more economic potential we have (up to 20%).
			dist := m.GetDistance(c.ID, c2.ID)
			if dist <= radius {
				if economicPotential[j] > economicPotential[i] {
					// We don't profit as much from the trade (up to 10%).
					tradePotential[i] += economicPotential[i] * (1 - dist/radius) * 0.1
				} else {
					// We profit more from the trade (up to 20%).
					tradePotential[i] += economicPotential[j] * (1 - dist/radius) * 0.2
				}
			}
		}
		//log.Printf("City %s (%s) can trade with %d cities.", c.Name, c.Type, count)
	}

	// Now we add the normalized trade potential to the economic potential.
	_, maxTrade := minMax(tradePotential)
	for i := range m.Cities {
		economicPotential[i] += tradePotential[i] / maxTrade
	}
	// Now we have a economic potential ranging from 0 to 3.

	// Assign the economic potential (range 0 to 3).
	for i, c := range m.Cities {
		c.EconomicPotential = economicPotential[i]

		// Log the economic potential (remove later)
		// log.Printf("City %s (%s) has economic potential %f", c.Name, c.Type, c.EconomicPotential)
	}
}

func (m *Civ) calculateAttractiveness() {
	// The attractiveness of a region is dependent on the following factors:
	// - Climate and elevation
	// - Distance to water (ocean, river, lake)
	// - Arable land (self-sufficiency)
	climateFitnessFunc := m.getFitnessClimate()
	arableLandFitnessFunc := m.getFitnessArableLand()
	proximityToWaterFitnessFunc := m.getFitnessProximityToWater()

	// TODO: Use getAttractivenessFunc
	// Calculate the attractiveness of all cities.
	for _, c := range m.Cities {
		// Log how close the city is to water.
		// log.Printf("City %s (%s) has proximity %f to water", c.Name, c.Type, proximityToWaterFitnessFunc(c.ID))
		// Log the arable land.
		// log.Printf("City %s (%s) has %f arable land", c.Name, c.Type, arableLandFitnessFunc(c.ID))
		// Log the climate.
		// log.Printf("City %s (%s) has climate %f", c.Name, c.Type, climateFitnessFunc(c.ID))
		// The attractiveness is the average of the fitness functions.
		c.Attractiveness = (climateFitnessFunc(c.ID) + arableLandFitnessFunc(c.ID) + proximityToWaterFitnessFunc(c.ID)) / 3
	}
}

func (m *Civ) getAttractivenessFunc() func(int) float64 {
	// The attractiveness of a region is dependent on the following factors:
	// - Climate and elevation
	// - Distance to water (ocean, river, lake)
	// - Arable land (self-sufficiency)
	climateFitnessFunc := m.getFitnessClimate()
	arableLandFitnessFunc := m.getFitnessArableLand()
	proximityToWaterFitnessFunc := m.getFitnessProximityToWater()

	return func(regionID int) float64 {
		// The attractiveness is the average of the fitness functions.
		return (climateFitnessFunc(regionID) + arableLandFitnessFunc(regionID) + proximityToWaterFitnessFunc(regionID)) / 3
	}
}

func (m *Civ) tickCityDays(c *City, days int) {
	// Check if the city is abandoned.
	if c.Population <= 0 {
		return
	}

	// There is a chance of some form of disaster.
	// For example, a mining town should have a chance of a cave in,
	// or a farming town should have a chance of a drought.
	// Towns close to volcanoes should have a chance of a volcanic
	// eruption while fault lines should have a chance of an earthquake.
	// If towns are heavily affected, they might be destroyed or abandoned.
	// Also the life expectancy might be low in early history.
	// Wars might also have a negative effect on the population.

	// TODO: With increasing population, the city is be more prone to famine
	// or disease.

	// Check if a random disaster strikes.
	if c.Population > 4000 && rand.Intn(100*356) < days {
		dead := rand.Intn(c.Population / 2)

		// Randomly choose a disaster from the list of strings.
		disasters := []string{
			"Plague",
			"Famine",
			"Drought",
			"Earthquake",
			"Volcanic eruption",
			"War",
			"Fire",
			"Riots",
			"Revolution",
			"Rebellion",
			"Storm",
			"Blizzard",
			"Wildfire",
			"Sinkhole",
		}
		disaster := disasters[rand.Intn(len(disasters))]
		year := m.Geo.Calendar.GetYear()
		m.AddEvent(&Event{
			Year: year,
			Type: disaster,
			Msg:  fmt.Sprintf("%d people died in %s", dead, c.Name),
			ID: ObjectReference{
				Type: ObjectTypeCity,
				ID:   c.ID,
			},
		})
		// Log the disaster, what type, how many people died and where.
		log.Printf("Year %d, Disaster: %s, %d people died in %s", year, disaster, dead, c.Name)
		c.Population -= dead
		if c.Population < 0 {
			c.Population = 0
		}
	}

	// Check if there is still anyone alive.
	if c.Population == 0 {
		return
	}

	// TODO: If there is sickness, war, famine, drought, etc, the population
	// might migrate to other cities that are more prosperous or a new settlement
	// might be founded nearby.

	// TODO: If a city reaches a certain size it might transition from an
	// agricultural town to a city with a more diverse economy. A mining town
	// might, if there is enough resources, transition to an industrial city.

	// In the middle ages, the average population growth was 0.16%-ish per year.
	// See: https://en.wikipedia.org/wiki/Medieval_demography
	//
	// TODO: The population growth should be dependent on the economic power
	// and if there is famine, war, drought, sickness, etc.
	// TODO: Also take in account what size of population the city can sustain.
	factor := float64(c.Population) * c.PopulationGrowthRate() * float64(days) / 365
	if factor >= 1 {
		c.Population += int(math.Ceil(factor))
	} else if m.rand.Float64() < factor {
		c.Population++
	}

	// Calculate the limit of the population based on attractiveness and
	// economic potential.
	maxPop := c.MaxPopulationLimit()
	if c.Population > maxPop {
		log.Println("City population limit reached:", c.Name, c.Population, maxPop)
		excessPopulation := c.Population - maxPop
		// TODO: The excess population should migrate to other cities or
		// a new settlement might be founded nearby.
		// Since we don't want to constantly migrate people, we just
		// move a larger portion of entire population, so that we drop
		// way below the limit, giving us a chance to grow again for
		// a while.

		// Move 10% of the population.
		excessPopulation = maxInt(excessPopulation*12/10, c.Population/20)
		excessPopulation = minInt(excessPopulation, c.Population-1)
		m.relocateFromCity(c, excessPopulation)
		/*
			// Find the best suitable neighbor region.
			bestReg := -1
			bestScore := 0.0

			climateFitnessFunc := m.getFitnessClimate()
			arableLandFitnessFunc := m.getFitnessArableLand()
			proximityToWaterFitnessFunc := m.getFitnessProximityToWater()
			for _, nb := range m.GetRegNeighbors(c.ID) {
				attr := (climateFitnessFunc(c.ID) + arableLandFitnessFunc(c.ID) + proximityToWaterFitnessFunc(c.ID)) / 3
				if attr > bestScore {
					bestScore = attr
					bestReg = nb
				}
			}
			if bestReg != -1 {
				m.placeCityAt(bestReg, TownTypeFarming, excessPopulation, bestScore)
			}*/
		c.Population -= excessPopulation
		m.calculateEconomicPotential()
	}

	// Update the peak population.
	// TODO: Maybe store the year when the peak population was reached?
	if c.Population > c.MaxPopulation {
		c.MaxPopulation = c.Population
	}
}

// relocateFromCity moves a portion of the population from the city to
// another city or a new settlement.
func (m *Civ) relocateFromCity(c *City, population int) {
	// Copy the slice of cities, so that we can modify it.
	cities := make([]*City, 0, len(m.Cities)-1)
	occupied := make(map[int]bool)
	for _, city := range m.Cities {
		occupied[city.ID] = true

		// Remove the origin city from the list of cities.
		if city.ID != c.ID {
			cities = append(cities, city)
		}
	}

	// Sort the cities by distance from the city.
	sort.Slice(cities, func(i, j int) bool {
		return m.GetDistance(c.ID, cities[i].ID) < m.GetDistance(c.ID, cities[j].ID)
	})

	// Check if any of the three closest cities have enough space.
	maxClosestCities := 5
	for i, city := range cities {
		if i >= maxClosestCities {
			break
		}
		maxPop := city.MaxPopulationLimit()
		if city.Population+population*12/10 <= maxPop {
			// Move the population to the closest city.
			city.Population += population
			return
		}
	}

	// Since the closest city doesn't have enough space, we need to
	// create a new settlement.
	attFunc := m.getAttractivenessFunc()

	// Find the best suitable neighbor region.
	bestReg := -1
	bestScore := 0.0
	maxDepth := 10
	seenRegions := make(map[int]bool)
	var traverseNeighbors func(id int, depth int)
	traverseNeighbors = func(id int, depth int) {
		for _, nb := range m.GetRegNeighbors(id) {
			if seenRegions[nb] {
				continue
			}
			seenRegions[nb] = true
			attr := attFunc(nb)
			if attr > bestScore && !occupied[nb] && nb != c.ID {
				bestScore = attr
				bestReg = nb
			}
			if depth < maxDepth {
				traverseNeighbors(nb, depth+1)
			}
		}
	}
	traverseNeighbors(c.ID, 0)
	if bestReg != -1 {
		c := m.placeCityAt(bestReg, TownTypeFarming, population, bestScore)
		c.Founded = m.History.GetYear()
		c.Attractiveness = bestScore
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
	m.tickCityDays(c, 1)

	// TODO: If the "sustainability" of the city is lower than needed to
	// sustain the population, the population will decrease.
	// This can be caused by war, disease, disaster, etc...
	// Some people might migrate to other cities that have a higher
	// sustainability.
}

// City represents a city in the world.
type City struct {
	ID                int       // Region where the city is located
	Name              string    // Name of the city
	Type              TownType  // Type of city
	Score             float64   // Score of the fitness function
	Population        int       // Population of the city
	MaxPopulation     int       // Maximum population of the city
	Culture           *Culture  // Culture of the city region
	Language          *Language // Language of the city
	EconomicPotential float64   // Economic potential of the city
	Attractiveness    float64   // Attractiveness of the city
	Founded           int64     // Year when the city was founded
}

// String returns a string representation of the city.
func (c *City) String() string {
	return fmt.Sprintf("%s (%d)", c.Name, c.Population)
}

// MaxPopulationLimit returns the maximum population sustainable by the city.
func (c *City) MaxPopulationLimit() int {
	return int(2000 * math.Pow((c.EconomicPotential+c.Attractiveness), 2))
}

// PopulationGrowthRate returns the population growth rate per year.
func (c *City) PopulationGrowthRate() float64 {
	return 0.3 * (c.EconomicPotential + c.Attractiveness) / (2 * 100)
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
	occupied := make(map[int]bool)
	for _, c := range m.Cities {
		occupied[c.ID] = true
	}

	var newcity int
	lastMax := math.Inf(-1)
	for i, val := range m.CalcCityScore(scoreFunc, distSeedFunc) {
		if val > lastMax && !occupied[i] {
			newcity = i
			lastMax = val
		}
	}
	if newcity == -1 {
		panic("no suitable location found")
	}
	// Get base population from city type.
	// TODO: Calculate population based on suitability for habitation.
	basePop := cType.FoundingPopulation()
	basePop += 2 * m.rand.Intn(basePop) / (len(m.Cities) + 1)
	return m.placeCityAt(newcity, cType, basePop, lastMax)
}

func (m *Civ) placeCityAt(r int, cType TownType, pop int, score float64) *City {
	c := &City{
		ID:            r,
		Score:         score,
		Population:    pop,
		MaxPopulation: pop,
		Type:          cType,
		Culture:       m.GetCulture(r),
		Founded:       m.Settled[r],
	}

	// If there is no known culture, generate a new one.
	if c.Culture == nil {
		c.Culture = m.PlaceCultureAt(r) // TODO: Grow this culture.
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

	return func(r int) float64 {
		// If we are below (or at) sea level, or we are in a pool of water,
		// assign lowest score and continue.
		if m.Elevation[r] <= 0 || m.Waterpool[r] > 0 {
			return -1.0
		}

		// Visit all neighbors and modify the score based on their properties.
		var hasWaterBodyBonus bool
		nbs := m.GetRegNeighbors(r)

		// Initialize fitness score with the normalized flux value.
		// This will favor placing cities along (and at the end of)
		// large rivers.
		score := math.Sqrt(m.Flux[r] / maxFlux)
		for _, nb := range nbs {
			// Add bonus if near ocean or lake.
			if m.isRegBelowOrAtSeaLevelOrPool(nb) {
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
				if wbSize := m.getRegLakeOrWaterBodySize(nb); wbSize > 0 {
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
		stp := steepness[r]
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
		bf := m.getRegWhittakerModBiomeFunc()
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
