package genworldvoronoi

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/gameconstants"
	"github.com/Flokey82/go_gens/genbiome"
)

func (m *Civ) getExistingCities() []*City {
	var cities []*City
	for _, c := range m.Cities {
		if c.Founded <= m.History.GetYear() {
			cities = append(cities, c)
		}
	}
	return cities
}

func (m *Civ) calculateEconomicPotential() {
	// TODO: Cities should have several values
	// Some are static, some are dynamic.
	//
	// Static:
	//
	// Static values are based on the region and are not affected by
	// the population.
	//
	// - Local resources
	//   metals, food, etc.
	// - Arable land score
	//   how much land is arable
	// - Climate
	//   how attractive is the climate for settlement
	// - Access to water
	//
	// Dynamic:
	//
	// This is based on the population, which directly impacts the
	// maximum distance from which resources can be gathered,
	// and the number of cities we can trade with.
	//
	// - Trade with nearby cities
	// - Nearby resources
	//
	// Other interesting values to consider:
	//   culture, is capital, etc.

	// We only consider cities that are founded prior or in the current year.
	cities := m.getExistingCities()

	// Calculate the analog of distance between regions by taking the surface
	// of a sphere with radius 1 and dividing it by the number of regions.
	// The square root will work as a somewhat sensible approximation of
	// distance.
	distRegion := math.Sqrt(4 * math.Pi / float64(m.mesh.numRegions))

	// Calculate the base radius in which we can find trade partners.
	var tradeRadius []float64
	for _, c := range cities {
		// The base radius is dependent on the population.
		// ... allow for at least two regions distance.
		radius := c.radius() + 2*distRegion
		tradeRadius = append(tradeRadius, radius)
	}

	economicPotential := make([]float64, len(cities))
	for i, c := range cities {
		economicPotential[i] = c.Resources + c.Agriculture
	}

	// Now we go through all the cities, and see if they might be able to
	// trade with each other. This way they can profit from each other's
	// resources.
	//
	// In the future we make this dependent on geographic features, where
	// mountains or the sea might be a barrier.
	//
	// TODO: This should in particular also take in account what kind of
	// resources are available and which are needed, so we would trade
	// only if we have benefits from it. This would also mean that far
	// away mining towns might profit from trade.
	tradePotential := make([]float64, len(cities))
	for i, c := range cities {
		// Calculate the distance field of all cities to the current city.
		if c.Population == 0 {
			continue
		}

		// Loop through all cities and check if we can trade with them.
		for j, c2 := range cities {
			// We don't trade with ourselves.
			if i == j || c2.Population == 0 {
				continue
			}
			// The trade radius is the sum of the two cities' radius times their economic potential.
			radius := tradeRadius[i]*(1+economicPotential[i]) + tradeRadius[j]*(1+economicPotential[j])

			// If the distance is within the radius, we can trade.
			// However, if the other city has a higher economic potential,
			// we profit less from the trade.
			// TODO: Switch this to population size?
			// the closer we are, the more economic potential we have (up to 20%).
			dist := m.GetDistance(c.ID, c2.ID)
			if dist <= radius {
				if economicPotential[j] > economicPotential[i] {
					// We don't profit as much from the trade (up to 15%).
					tradePotential[i] += economicPotential[i] * (1 - dist/radius) * 0.15
				} else {
					// We profit more from the trade (up to 20%).
					tradePotential[i] += economicPotential[j] * (1 - dist/radius) * 0.2
				}
			}
		}
	}

	// DEBUG: Count the number of cities in range.
	for i, c := range cities {
		// Loop through all cities and check if we can trade with them.
		var count int
		for j, c2 := range cities {
			if i == j {
				continue // We don't trade with ourselves.
			}
			dist := m.GetDistance(c.ID, c2.ID)
			if dist <= tradeRadius[i] {
				count++
			}
		}
		c.TradePartners = count
	}

	// Now normalize trade potential.
	_, maxTrade := minMax(tradePotential)
	if maxTrade > 0 {
		for i := range cities {
			tradePotential[i] /= maxTrade
		}
	}

	// Assign the economic potential.
	for i, c := range cities {
		c.EconomicPotential = economicPotential[i] + tradePotential[i]
		c.Trade = tradePotential[i]
	}
}

func (m *Civ) calculateAttractiveness(cities []*City) {
	// Calculate the attractiveness of the supplied cities.
	attrFunc := m.getAttractivenessFunc()
	for _, c := range cities {
		c.Attractiveness = attrFunc(c.ID)
	}
}

func (m *Civ) calculateAgriculturalPotential(cities []*City) {
	// Calculate the agricultural potential of the supplied cities.
	fitnessArableFunc := m.getFitnessArableLand()
	for _, c := range cities {
		if agrPotential := fitnessArableFunc(c.ID); agrPotential > 0 {
			c.Agriculture = agrPotential
		}
	}
}

func (m *Civ) calculateResourcePotential(cities []*City) {
	// Now get the resource potential of all cities.
	calcResourceValues := func(res []byte) {
		for _, c := range cities {
			// Sum up the normalized resource values.
			c.Resources += float64(sumResources(res[c.ID])) / 36 // 36 is the maximum value.
		}
	}

	// Reset the resource potential.
	for _, c := range cities {
		c.Resources = 0
	}

	// Calculate the resource potential for each resource.
	calcResourceValues(m.Metals)
	calcResourceValues(m.Gems)
	calcResourceValues(m.Stones)
	calcResourceValues(m.Wood)
	calcResourceValues(m.Various)
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

type disaster struct {
	Name           string
	PopulationLoss float64
}

var (
	disStorm      = disaster{"Storm", 0.01}
	disFire       = disaster{"Fire", 0.02}
	disRockslide  = disaster{"Rockslide", 0.03}
	disCaveIn     = disaster{"Cave In", 0.05}
	disWildfire   = disaster{"Wildfire", 0.07}
	disDrought    = disaster{"Drought", 0.1}
	disFamine     = disaster{"Famine", 0.15}
	disDisease    = disaster{"Disease", 0.25}
	disEarthquake = disaster{"Earthquake", 0.3}
	disFlood      = disaster{"Flood", 0.35}
	disVolcano    = disaster{"Volcanic Eruption", 0.6}
	disPlague     = disaster{"Plague", 0.8}
)

var disasters = []disaster{
	disStorm,
	disFire,
	disRockslide,
	disCaveIn,
	disWildfire,
	disDrought,
	disFamine,
	disDisease,
	disEarthquake,
	disFlood,
	disVolcano,
	disPlague,
}

var sumDisasterProbability float64

func init() {
	// The probabilities are the inverse of the population loss.
	for _, d := range disasters {
		sumDisasterProbability += 1 - d.PopulationLoss
	}
}

func (m *Civ) tickCityDays(c *City, days int) {
	// Check if the city is abandoned.
	if c.Population <= 0 {
		if c.Population < 0 {
			c.Population = 0
		}
		return
	}

	// Enable / disable migration of population when a disaster strikes.
	enableDisasterMigration := true

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
	if m.rand.Intn(100*356) < days {
		// Pick a random disaster given their respective probabilities.
		r := rand.Float64() * sumDisasterProbability
		sumProb := 0.0
		var dis disaster
		for _, d := range disasters {
			sumProb += 1 - d.PopulationLoss
			if r < sumProb {
				dis = d
				break
			}
		}

		if dis.Name == "" {
			log.Fatalf("No disaster was chosen")
		}

		// Calculate the population loss.
		popLoss := dis.PopulationLoss * (2 + rand.Float64()) / 3
		dead := int(math.Ceil(float64(c.Population) * popLoss))
		year := m.Geo.Calendar.GetYear()

		// Add an event to the calendar.
		m.AddEvent(dis.Name, fmt.Sprintf("%d people died", dead), c.Ref())

		// Reduce the population.
		c.Population -= dead
		if c.Population <= 0 {
			c.Population = 0
			return
		}

		// Log the disaster, what type, how many people died and where.
		log.Printf("Year %d: %s, %d people died in %s", year, dis.Name, dead, c.Name)

		// Since there was a disaster, depending on the number of people that
		// died, some people might leave the city.
		//
		// If there is sickness, war, famine, drought, etc, the population might
		// migrate to other cities that are more prosperous or a new settlement
		// might be founded nearby.
		//
		// The bigger the population loss, the more likely it is that people
		// will leave the city.
		if enableDisasterMigration && rand.Float64() < popLoss {
			// Up to 'popLoss' of the population might leave the city.
			leave := int(float64(c.Population) * (popLoss * rand.Float64()))
			m.relocateFromCity(c, leave)
		}
	}

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
		log.Printf("Attractiveness: %.2f, Economic Potential: %.2f, Agriculture: %.2f", c.Attractiveness, c.EconomicPotential, c.Agriculture)
		// The excess population can migrate to other cities or
		// a new settlement might be founded nearby.
		// Since we don't want to constantly migrate people, we just
		// move a larger portion of entire population, so that we drop
		// way below the limit, giving us a chance to grow again for
		// a while.
		excessPopulation := c.Population - maxPop

		// Move 10% of the population or 1.2 times the excess population,
		// whichever is larger.
		excessPopulation = maxInt(excessPopulation*12/10, c.Population/10)

		// Make sure we don't move more than the entire population.
		excessPopulation = minInt(excessPopulation, c.Population)

		m.relocateFromCity(c, excessPopulation)
	}

	// TODO: If a city reaches a certain size it might transition from an
	// agricultural town to a city with a more diverse economy. A mining town
	// might, if there is enough resources, transition to an industrial city.

	// TODO: Update the economic potential of the city if the population
	// has changed.
	// m.calculateEconomicPotential()

	// Update the peak population.
	// TODO: Maybe store the year when the peak population was reached?
	if c.Population > c.MaxPopulation {
		c.MaxPopulation = c.Population
	}
}

// relocateFromCity moves a portion of the population from the city to
// another city or a new settlement.
// TODO: Distribute more evenly if a large group of people are moving.
func (m *Civ) relocateFromCity(c *City, population int) {
	// If no one is migrating, then there is nothing to do.
	if population <= 0 {
		return
	}

	// Move out the migrating population.
	if c.Population < population {
		population = c.Population
	}
	c.Population -= population

	// Add an event to the calendar.
	m.AddEvent("Migration", fmt.Sprintf("%d people left", population), c.Ref())

	// Calculate the analog of distance between regions by taking the surface
	// of a sphere with radius 1 and dividing it by the number of regions.
	// The square root will work as a somewhat sensible approximation of
	// distance.
	distRegion := math.Sqrt(4 * math.Pi / float64(m.mesh.numRegions))

	// Per distRegion, the chance of death is 2%.
	const chanceDeath = 0.02
	calcChanceDeath := func(dist float64) float64 {
		return 1 - math.Pow(1-chanceDeath, dist/distRegion)
	}

	cities := m.getExistingCities()

	// Sort the cities by increasing distance from the city.
	sort.Slice(cities, func(i, j int) bool {
		return m.GetDistance(c.ID, cities[i].ID) < m.GetDistance(c.ID, cities[j].ID)
	})

	// Check if any of the n closest cities have enough space.
	numClosestCities := 10

	// The closest city is the city itself, so skip it.
	for _, city := range cities[1:minInt(len(cities), numClosestCities+1)] {
		maxPop := city.MaxPopulationLimit()
		popCapacity := maxPop - city.Population

		// If there is capacity, a portion of the population might move there.
		if popCapacity > 0 {
			// Now pick a fraction of the population that will move to the city,
			// with the largest fraction going to the closest city.
			numMigrants := minInt(population, popCapacity/2)

			// Make sure we don't increase the population by more than 20%,
			// except if the city is abandoned.
			if city.Population > 0 {
				numMigrants = minInt(numMigrants, city.Population/5)
			}

			// Depending on the distance, some of the population might
			// die on the way.
			dist := m.GetDistance(c.ID, city.ID)
			dead := int(math.Ceil(calcChanceDeath(dist) * float64(numMigrants)))
			survived := numMigrants - dead

			// If any survived, move them to the city.
			if survived > 0 {
				// If the city is abandoned, set the economic potential to 1 temporarily.
				if city.Population == 0 {
					city.EconomicPotential = 1
				}
				// Move the population to the closest city.
				city.Population += survived
				if city.Population > city.MaxPopulation {
					city.MaxPopulation = city.Population
				}
				// TODO: Update the economic potential of the city.
				m.AddEvent("Migration", fmt.Sprintf("%d people arrived", population), city.Ref())
			}
			log.Printf("%d people moved from %s to %s, %d died on the way", numMigrants, c.Name, city.Name, dead)

			// Subtract the number of people that moved from the total
			// population that is migrating.
			population -= numMigrants
			if population <= 0 {
				return
			}
		}
	}

	// Make sure we don't place a new settlement in a region that is already
	// occupied.
	occupied := make(map[int]bool)
	for _, city := range m.Cities {
		occupied[city.ID] = true
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
		if depth >= maxDepth {
			return
		}
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
			traverseNeighbors(nb, depth+1)
		}
	}
	traverseNeighbors(c.ID, 0)

	// If we found a suitable region, create a new city there.
	if bestReg != -1 {
		// Depending on the distance, some of the population might die on the way.
		dist := m.GetDistance(c.ID, bestReg)
		dead := int(math.Ceil(calcChanceDeath(dist) * float64(population)))
		survived := population - dead

		// Check if any survived and founded a new city.
		if survived > 0 {
			city := m.placeCityAt(bestReg, m.getRegCityType(bestReg), survived, bestScore)
			city.Founded = m.History.GetYear() + 1 // The city is founded next year.
			city.EconomicPotential = 1             // Set the economic potential to 1 temporarily.
			city.Attractiveness = bestScore

			// TODO: Set the economic potential and attractiveness of the new city.
			log.Printf("%d people moved from %s and founded %s, %d died on the way", population, c.Name, city.Name, dead)
		}
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
	Founded           int64     // Year when the city was founded
	EconomicPotential float64   // Economic potential of the city (DYNAMIC)
	Trade             float64   // Trade value of the city (DYNAMIC)
	Resources         float64   // Resources value of the city (PARTLY DYNAMIC)
	Agriculture       float64   // Agriculture value of the city (STATIC)
	Attractiveness    float64   // Attractiveness of the city (STATIC)
	TradePartners     int       // Number of cities within trade range
}

// Ref returns the object reference of the city.
func (c *City) Ref() ObjectReference {
	return ObjectReference{
		ID:   c.ID,
		Type: ObjectTypeCity,
	}
}

func (c *City) radius() float64 {
	// In kilometers.
	if c.Population <= 0 {
		return 0
	}
	return 100 * math.Sqrt(float64(c.Population)/math.Pi) / gameconstants.EarthCircumference
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
		Founded:       m.Settled[r] + m.rand.Int63n(100),
	}

	// TODO: Set agricultural potential and resources based on the region.

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

// getRegCityType returns the optimal type of city for a given region.
func (m *Civ) getRegCityType(r int) TownType {
	// If we have a lot of metals, gems, etc. we have a mining town.
	if m.Metals[r] > 0 || m.Gems[r] > 0 {
		return TownTypeMining
	}

	// If we have stone, we have a quarry.
	if m.Stones[r] > 0 {
		return TownTypeQuarry
	}

	// TODO: Cache this somehow.
	if m.getFitnessArableLand()(r) > 0.5 {
		return TownTypeFarming
	}
	// TODO: Add more types of cities.
	return TownTypeDefault
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
	TownTypeMiningGems  TownType = "mining (gems)"
	TownTypeQuarry      TownType = "quarry"
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
	case TownTypeQuarry, TownTypeMining, TownTypeMiningGems:
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

func (m *Civ) getFitnessProximityToCities(except ...TownType) func(int) float64 {
	var cities []int
	exceptMap := make(map[TownType]bool)
	for _, t := range except {
		exceptMap[t] = true
	}
	for _, c := range m.Cities {
		if !exceptMap[c.Type] {
			cities = append(cities, c.ID)
		}
	}
	distCities := m.assignDistanceField(cities, make(map[int]bool))
	_, maxDist := minMax(distCities)
	if maxDist == 0 {
		maxDist = 1
	}
	return func(r int) float64 {
		if distCities[r] == 0 {
			return 0
		}
		return 1 - float64(distCities[r])/maxDist
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
	case TownTypeQuarry:
		fa := m.getFitnessSteepMountains()
		fb := m.getFitnessClimate()
		fc := m.getFitnessProximityToWater()
		fd := m.getFitnessProximityToCities(TownTypeMining, TownTypeMiningGems, TownTypeQuarry)
		return func(r int) float64 {
			if m.Stones[r] == 0 {
				return -1.0
			}
			return fd(r) * (fa(r)*fb(r) + fc(r)) / 2
		}
	case TownTypeMining:
		fa := m.getFitnessSteepMountains()
		fb := m.getFitnessClimate()
		fc := m.getFitnessProximityToWater()
		fd := m.getFitnessProximityToCities(TownTypeMining, TownTypeMiningGems, TownTypeQuarry)
		return func(r int) float64 {
			if m.Metals[r] == 0 {
				return -1.0
			}
			return fd(r) * (fa(r)*fb(r) + fc(r)) / 2
		}
	case TownTypeMiningGems:
		fa := m.getFitnessSteepMountains()
		fb := m.getFitnessClimate()
		fc := m.getFitnessProximityToWater()
		fd := m.getFitnessProximityToCities(TownTypeMining, TownTypeMiningGems, TownTypeQuarry)
		return func(r int) float64 {
			if m.Gems[r] == 0 {
				return -1.0
			}
			return fd(r) * (fa(r)*fb(r) + fc(r)) / 2
		}
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
