package genworldvoronoi

import (
	"log"
	"time"
)

type Civ struct {
	*Geo
	RegionToTerritory []int      // (political) Point / region mapping to territory / empire
	RegionToCityState []int      // (political) Point / region mapping to city / city state
	Cities            []*City    // (political) City seed points / regions
	RegionToCulture   []int      // (cultural) Point / region mapping to culture
	Cultures          []*Culture // (cultural) Culture seed points / regions
	NumCities         int        // Number of generated cities (regions)
	NumCityStates     int        // Number of generated city states
	NumMiningTowns    int        // Number of generated mining towns
	NumFarmingTowns   int        // Number of generated farming towns
	NumDesertOasis    int        // Number of generated desert oases
	NumTerritories    int        // Number of generated territories
	NumCultures       int        // (Min) Number of generated cultures
}

func NewCiv(geo *Geo) *Civ {
	return &Civ{
		Geo:               geo,
		RegionToTerritory: initRegionSlice(geo.mesh.numRegions),
		NumTerritories:    10,
		NumCities:         150,
		NumCityStates:     150,
		NumMiningTowns:    60,
		NumFarmingTowns:   60,
		NumDesertOasis:    10,
		NumCultures:       30,
	}
}

func (m *Civ) generateCivilization() {
	// Place and grow cultures.
	// Place nature religions.
	// ? Create sub branches of religions.
	// Place cities.
	// Grow city states.
	// Create organized religions.
	// Create empires.

	// Place cultures.
	start := time.Now()
	m.PlaceNCultures(m.NumCultures)
	log.Println("Done cultures in ", time.Since(start).String())

	// Place cities and territories in regions.
	start = time.Now()
	m.PlaceNCities(m.NumCities, TownTypeDefault)
	m.PlaceNCities(m.NumMiningTowns, TownTypeMining)
	m.PlaceNCities(m.NumFarmingTowns, TownTypeFarming)
	m.PlaceNCities(m.NumDesertOasis, TownTypeDesertOasis)
	log.Println("Done cities in ", time.Since(start).String())

	start = time.Now()
	m.rPlaceNTerritories(m.NumTerritories)
	log.Println("Done territories in ", time.Since(start).String())

	start = time.Now()
	m.rPlaceNCityStates(m.NumCityStates)
	log.Println("Done city states in ", time.Since(start).String())

	// Once we have established the territories, we can add trade towns
	// (we need the territories for the trade routes).
	// We should probably establish the trade routes now, so we ensure
	// that the trade towns will still be placed on the nexus points
	// where trade routes meet.
	// start = time.Now()
	// m.rPlaceNCities(30, TownTypeTrading)
	// log.Println("Done trade cities in ", time.Since(start).String())

	//m.GetEmpires()
	m.GetEmpires2()
}
