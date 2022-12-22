package genworldvoronoi

import (
	"container/heap"
	"log"
	"time"

	"github.com/Flokey82/go_gens/genbiome"
)

type Civ struct {
	*Geo
	RegionToEmpire    []int       // (political) Point / region mapping to territory / empire
	RegionToCityState []int       // (political) Point / region mapping to city / city state
	Cities            []*City     // (political) City seed points / regions
	RegionToCulture   []int       // (cultural) Point / region mapping to culture
	Cultures          []*Culture  // (cultural) Culture seed points / regions
	Religions         []*Religion // (cultural) Religion seed points / regions
	Settled           []int64     // (cultural) Time of settlement per region
	// SettledBySpecies []int // (cultural) Which species settled the region first
	NumCities       int // Number of generated cities (regions)
	NumCityStates   int // Number of generated city states
	NumMiningTowns  int // Number of generated mining towns
	NumFarmingTowns int // Number of generated farming towns
	NumDesertOasis  int // Number of generated desert oases
	NumEmpires      int // Number of generated territories
	NumCultures     int // (Min) Number of generated cultures
	NameGen         *NameGenerators
}

func NewCiv(geo *Geo) *Civ {
	return &Civ{
		Geo:               geo,
		RegionToEmpire:    initRegionSlice(geo.mesh.numRegions),
		RegionToCityState: initRegionSlice(geo.mesh.numRegions),
		RegionToCulture:   initRegionSlice(geo.mesh.numRegions),
		Settled:           initTimeSlice(geo.mesh.numRegions),
		NumEmpires:        10,
		NumCities:         150,
		NumCityStates:     150,
		NumMiningTowns:    60,
		NumFarmingTowns:   60,
		NumDesertOasis:    10,
		NumCultures:       30,
		NameGen:           NewNameGenerators(geo.Seed),
	}
}

func (m *Civ) generateCivilization() {
	// TODO: The generation should happen somewhat like this...
	// 0. Calculate time of settlement per region through flood fill.
	// This will allow us to determine the founding date of the cities and
	// settlements.
	m.generateTimeOfSettlement()
	// 1. Generate (races and) cultures.
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
	m.PlaceNCities(m.NumMiningTowns, TownTypeMining)
	m.PlaceNCities(m.NumFarmingTowns, TownTypeFarming)
	m.PlaceNCities(m.NumDesertOasis, TownTypeDesertOasis)
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

	// HACK: Age city populations.
	// TODO: Instead we should spawn the cities from the capitals.
	// Also, the theoretical population should be based on the
	// economic potential of the region, the type of settlement,
	// and the time of settlement.
	_, maxSettled := minMax64(m.Settled)
	for _, c := range m.Cities {
		// Tick each city for the number of years since it was settled.
		for j := 0; j < int(maxSettled-m.Settled[c.ID]); j++ {
			m.tickCityDays(c, 365)
		}
	}

	//m.GetEmpires()
	m.calculateEconomicPotential()
}

func (m *Civ) Tick() {
	// Update cities.
	// 1: Update production.
	// 2: Update consumption.
	// 3: Update trade.
	// 4: Update population. (births, deaths, migration)
	// 5: Found new settlements?
	for _, c := range m.Cities {
		m.TickCity(c)
	}
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
		genbiome.WhittakerModBiomeTemperateRainForest,
		genbiome.WhittakerModBiomeTemperateSeasonalForest,
		genbiome.WhittakerModBiomeTropicalRainForest,
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

		// If the terrain weight is negative, the region is ocean.
		// This means, we need boats to get there, which will require
		// more time.
		if terrWeight < 0 {
			// Once we are at sea, we travel at a speed of 20 years per
			// region.
			if (m.Elevation[v] <= 0) && (m.Elevation[u] <= 0) {
				return float64(settleTime[u]) + 20
			} else if m.Elevation[v] > 0 {
				// If we arrive at land, we only need a year.
				return float64(settleTime[u]) + 1
			}
			// It takes us 100 years to build a boat.
			return float64(settleTime[u]) + 100
		}
		// The settle time is a fraction of 1000 years.
		return float64(settleTime[u]) + 1000*terrWeight //*(1-terrainArable(v))
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
