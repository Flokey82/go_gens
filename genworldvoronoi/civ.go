package genworldvoronoi

import (
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
	NumCities         int         // Number of generated cities (regions)
	NumCityStates     int         // Number of generated city states
	NumMiningTowns    int         // Number of generated mining towns
	NumFarmingTowns   int         // Number of generated farming towns
	NumDesertOasis    int         // Number of generated desert oases
	NumEmpires        int         // Number of generated territories
	NumCultures       int         // (Min) Number of generated cultures
	NameGen           *NameGenerators
}

func NewCiv(geo *Geo) *Civ {
	return &Civ{
		Geo:               geo,
		RegionToEmpire:    initRegionSlice(geo.mesh.numRegions),
		RegionToCityState: initRegionSlice(geo.mesh.numRegions),
		RegionToCulture:   initRegionSlice(geo.mesh.numRegions),
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
	m.rPlaceNEmpires(m.NumEmpires)
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
	// 1: Update city states.
	// 1.1: Update wealth / taxation.
	// 1.2: Update trade.
	// 1.3: Update politics.
	// (Alliances, wars, taxes, laws, etc.)
	// 1.4: Update population sentiment.
	// 2: Update empires.
	// (Similar as city states.)
	// 3: Update cultures.
	// 3.1: Expansion, assimilation, etc.
	// 3.2: Update culture sentiments.
	// 4: Update religions.
	// (same as cultures)

	// NOTE: In theory we can partially de-duplicate code relating
	// to city states and empires, since they function similarly.
	// We can also de-duplicate cultures and religions.
}

func (m *Civ) getRegionName(r int) string {
	switch m.getRWhittakerModBiomeFunc()(r) {
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
