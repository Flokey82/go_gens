package genworldvoronoi

import (
	"container/heap"
	"log"
	"time"

	"github.com/Flokey82/go_gens/genbiome"
)

type Civ struct {
	*Geo
	*History
	RegionToEmpire    []int       // (political) Point / region mapping to territory / empire
	RegionToCityState []int       // (political) Point / region mapping to city / city state
	Cities            []*City     // (political) City seed points / regions
	RegionToCulture   []int       // (cultural) Point / region mapping to culture
	Cultures          []*Culture  // (cultural) Culture seed points / regions
	Religions         []*Religion // (cultural) Religion seed points / regions
	Settled           []int64     // (cultural) Time of settlement per region
	// SettledBySpecies []int // (cultural) Which species settled the region first
	NumCities          int // Number of generated cities (regions)
	NumCityStates      int // Number of generated city states
	NumMiningTowns     int // Number of generated mining towns
	NumMiningGemsTowns int // Number of generated (gem) mining towns
	NumQuarryTowns     int // Number of generated quarry towns
	NumFarmingTowns    int // Number of generated farming towns
	NumDesertOasis     int // Number of generated desert oases
	NumEmpires         int // Number of generated territories
	NumCultures        int // (Min) Number of generated cultures
	NameGen            *NameGenerators
}

func NewCiv(geo *Geo) *Civ {
	return &Civ{
		Geo:                geo,
		History:            NewHistory(geo.Calendar),
		RegionToEmpire:     initRegionSlice(geo.mesh.numRegions),
		RegionToCityState:  initRegionSlice(geo.mesh.numRegions),
		RegionToCulture:    initRegionSlice(geo.mesh.numRegions),
		Settled:            initTimeSlice(geo.mesh.numRegions),
		NumEmpires:         10,
		NumCities:          150,
		NumCityStates:      150,
		NumMiningTowns:     60,
		NumMiningGemsTowns: 60,
		NumQuarryTowns:     60,
		NumFarmingTowns:    60,
		NumDesertOasis:     10,
		NumCultures:        30,
		NameGen:            NewNameGenerators(geo.Seed),
	}
}

func (m *Civ) generateCivilization() {
	enableCityAging := false

	// TODO: The generation should happen somewhat like this...
	// 0. Calculate time of settlement per region through flood fill.
	// This will allow us to determine the founding date of the cities and
	// settlements.
	m.generateTimeOfSettlement()

	// 1. Generate (species and) cultures.
	// 2. Spread cultures.
	// 3. Generate settlements.
	// 4. Grow settlements.
	// 5. Create organized religions.
	// 6. Spread religions.
	// 7. Select capital cities.
	// 8. Generate city states.
	// 9. Generate empires.

	// Place cultures.
	start := time.Now()
	m.PlaceNCultures(m.NumCultures)
	log.Println("Done cultures in ", time.Since(start).String())

	// Place folk religions.

	// Place cities and territories in regions.
	// TODO: Smaller towns should be found in the vicinity of larger cities.
	start = time.Now()
	m.PlaceNCities(m.NumCities, TownTypeDefault)
	m.PlaceNCities(m.NumFarmingTowns, TownTypeFarming)
	m.PlaceNCities(m.NumDesertOasis, TownTypeDesertOasis)
	m.PlaceNCities(m.NumMiningTowns, TownTypeMining)
	m.PlaceNCities(m.NumMiningGemsTowns, TownTypeMiningGems)
	m.PlaceNCities(m.NumQuarryTowns, TownTypeQuarry)
	log.Println("Done cities in ", time.Since(start).String())

	start = time.Now()
	m.rPlaceNCityStates(m.NumCityStates)
	log.Println("Done city states in ", time.Since(start).String())

	start = time.Now()
	m.regPlaceNEmpires(m.NumEmpires)
	log.Println("Done empires in ", time.Since(start).String())

	// Once we have established the territories, we can add trade towns
	// (we need the territories for the trade routes).
	// We should probably establish the trade routes now, so we ensure
	// that the trade towns will still be placed on the nexus points
	// where trade routes meet.
	// start = time.Now()
	// m.rPlaceNCities(30, TownTypeTrading)
	// log.Println("Done trade cities in ", time.Since(start).String())

	//m.GetEmpires()
	_, maxSettled := minMax64(m.Settled)
	m.Geo.Calendar.SetYear(maxSettled)

	start = time.Now()
	m.calculateAgriculturalPotential(m.Cities)
	log.Println("Done calculating agricultural potential in ", time.Since(start).String())

	start = time.Now()
	m.calculateAttractiveness(m.Cities)
	log.Println("Done calculating attractiveness in ", time.Since(start).String())

	start = time.Now()
	m.calculateResourcePotential(m.Cities)
	log.Println("Done calculating resource potential in ", time.Since(start).String())

	start = time.Now()
	m.calculateEconomicPotential()
	log.Println("Done calculating economic potential in ", time.Since(start).String())

	if enableCityAging {
		// HACK: Age city populations.
		start = time.Now()
		// TODO: Instead we should spawn the cities from the capitals.
		// Also, the theoretical population should be based on the
		// economic potential of the region, the type of settlement,
		// and the time of settlement.
		m.Geo.Calendar.SetYear(0)
		knownCities := len(m.Cities)
		for year := 0; year < int(maxSettled); year++ {
			// Age cities for a year.
			for _, c := range m.getExistingCities() {
				m.tickCityDays(c, 365)
			}

			// Update attractiveness, agricultural potential, and resource potential
			// for new cities.
			if len(m.Cities) > knownCities {
				// TODO: Only update new regions until we have climate change?
				m.calculateAttractiveness(m.Cities[knownCities:])
				m.calculateAgriculturalPotential(m.Cities[knownCities:])
				m.calculateResourcePotential(m.Cities[knownCities:])
				knownCities = len(m.Cities)
			}

			// Recalculate economic potential.
			m.calculateEconomicPotential()
			log.Printf("Aged cities to %d\n", year)

			// Advance year.
			m.Geo.Calendar.TickYear()
		}
		log.Println("Done aging cities in ", time.Since(start).String())
	}
}

func (m *Civ) Tick() {
	// Update cities.
	// 1: Update production.
	// 2: Update consumption.
	// 3: Update trade.
	// 4: Update population. (births, deaths, migration)
	// 5: Found new settlements?
	// for _, c := range m.Cities {
	//	m.TickCity(c)
	// }

	// TODO:
	// Update city states.
	// 1: Update wealth / taxation.
	// 2: Update trade.
	// 3: Update politics.
	// (Alliances, wars, taxes, laws, etc.)
	// 4: Update population sentiment.
	// Update empires.
	// (Similar as city states.)
	// Update cultures.
	// 1: Expansion, assimilation, etc.
	// 2: Update culture sentiments.
	// Update religions.
	// (same as cultures)

	// NOTE: In theory we can partially de-duplicate code relating
	// to city states and empires, since they function similarly.
	// We can also de-duplicate cultures and religions.
}

// getRegName attempts to generate a name for the given region.
func (m *Civ) getRegName(r int) string {
	switch m.getRegWhittakerModBiomeFunc()(r) {
	case genbiome.WhittakerModBiomeBorealForestTaiga,
		genbiome.WhittakerModBiomeTemperateRainforest,
		genbiome.WhittakerModBiomeTemperateSeasonalForest,
		genbiome.WhittakerModBiomeTropicalRainforest,
		genbiome.WhittakerModBiomeTropicalSeasonalForest:
		return m.NameGen.Forest.Generate(int64(r), r%2 == 0)
	case genbiome.WhittakerModBiomeHotSwamp,
		genbiome.WhittakerModBiomeWetlands:
		return m.NameGen.Swamp.Generate(int64(r), r%2 == 0)
	}
	return ""
}

func (m *Civ) generateTimeOfSettlement() {
	// First we pick a "suitable" region where the cradle of civilization
	// will be located.
	// There are some theories where, if we put the origin of civilization
	// in a less suitable region, we will expand to more suitable regions.
	// See: https://forhinhexes.blogspot.com/2019/08/history-xvii-cradle-of-civilizations.html?m=1
	// I feel like this will only work for migration to the most suitable
	// regions, but we know that people will also migrate to less suitable
	// regions, if they have to, or if they are forced to, or if they
	// are just too stubborn to give up.

	// Since we only have one species for now (humans), we will just start
	// with a 'steppe' region, and then expand from there incrementally.

	// How long it takes for the civilization to expand to a region is
	// determined by the characteristics of the region and if there are
	// more suitable regions nearby. So we will use a priority queue
	// to determine the next region to expand to.

	var queue ascPriorityQueue
	heap.Init(&queue)

	// 'settleTime' is the time when a region was settled.
	settleTime := initTimeSlice(m.mesh.numRegions)

	// Now we pick a suitable region to start with (steppe/grassland).
	// We will use the climate fitness function and filter by biome.
	bestRegion := -1
	bestFitness := 0.0
	fa := m.getFitnessClimate()
	bf := m.getRegWhittakerModBiomeFunc()
	for r := 0; r < m.mesh.numRegions; r++ {
		if bf(r) == genbiome.WhittakerModBiomeTemperateGrassland {
			fitness := fa(r)
			if fitness > bestFitness {
				bestFitness = fitness
				bestRegion = r
			}
		}
	}
	if bestRegion == -1 {
		panic("no suitable region found")
	}

	// We will start with a settlement time of 0.
	settleTime[bestRegion] = 0

	// terrainWeight returns high scores for difficult terrain.
	terrainWeight := m.getTerritoryWeightFunc()

	// terrainArable returns high scores if the terrain is arable.
	//terrainArable := m.getFitnessArableLand()

	// TODO: The duration that it takes to settle a region should
	// depend on how many regions there are in total (the size of
	// the regions).
	weight := func(o, u, v int) float64 {
		// Terrain weight.
		// TODO: We should use a slightly different weight function
		// that doesn't treat up- and downhill differently.
		// Also, the penalty should be way higher for "impassable"
		// terrain.
		terrWeight := terrainWeight(bestRegion, u, v)

		// If the terrain weight is positive (or zero), the destination region is land.
		var timeReqired float64
		if terrWeight >= 0 {
			// Settlement on land takes a fraction of 2000 years per (unit) region.
			// 'terrWeight' already takes the actual distance between the regions
			// into account.
			return float64(settleTime[u]) + 2000*terrWeight // * (1-terrainArable(v))
		}

		// If the terrain weight is negative, the source- and/or destination region is ocean.
		// This means, we need boats to get there, which will require more time.
		if m.Elevation[v] > 0 {
			// If we were at sea and arrive at land, we only need a year to disembark.
			timeReqired = 1
		} else if (m.Elevation[v] <= 0) && (m.Elevation[u] <= 0) {
			// Once we are traveling at sea, we travel at a speed of 20 years
			// per (unit) region.
			timeReqired = 20
		} else {
			// We were on land, but the destination is at sea,
			// it takes us 200 years to build a boat.
			timeReqired = 200
		}

		// Calculate the actual distance between the two regions,
		// so we are independent of the mesh resolution.
		actualDist := m.GetDistance(u, v)
		return float64(settleTime[u]) + timeReqired*actualDist
	}

	// Now add the region neighbors to the queue.
	for _, n := range m.GetRegNeighbors(bestRegion) {
		heap.Push(&queue, &queueEntry{
			origin:      bestRegion,
			score:       weight(bestRegion, bestRegion, n),
			destination: n,
		})
	}

	// Expand settlements until we have settled all regions.
	for queue.Len() > 0 {
		u := heap.Pop(&queue).(*queueEntry)

		// Check if the region has already been settled.
		if settleTime[u.destination] >= 0 {
			continue
		}

		// The higher the score, the more difficult it is to settle there,
		// and the longer it took to settle there.
		settleTime[u.destination] = int64(u.score)
		for _, v := range m.GetRegNeighbors(u.destination) {
			// Check if the region has already been settled.
			if settleTime[v] >= 0 {
				continue
			}
			newdist := weight(u.origin, u.destination, v)
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueEntry{
				score:       newdist,
				origin:      u.destination,
				destination: v,
			})
		}
	}

	// TODO: For crossing the ocean, we need to wait for boats to be invented.
	m.Settled = settleTime
}
