package genworldvoronoi

import (
	"log"
)

func (m *Civ) rPlaceNCityStates(n int) {
	m.resetRand()
	// Territories are based on cities acting as their capital.
	// Since the algorithm places the cities with the highes scores
	// first, we use the top 'n' cities as the capitals for the
	// territories.
	var seedCities []int
	for i, c := range m.Cities {
		if i >= n {
			break
		}
		seedCities = append(seedCities, c.ID)
	}
	weight := m.getTerritoryWeightFunc()
	biomeWeight := m.getTerritoryBiomeWeightFunc()
	cultureWeight := m.getTerritoryCultureWeightFunc()

	m.RegionToCityState = m.regPlaceNTerritoriesCustom(seedCities, func(o, u, v int) float64 {
		// TODO: Make sure we take in account expansionism, wealth, score, and culture.
		w := weight(o, u, v)
		if w < 0 {
			return -1
		}
		b := biomeWeight(o, u, v)
		if b < 0 {
			return -1
		}
		c := cultureWeight(o, u, v)
		if c < 0 {
			return -1
		}
		return (w + b + c) / 3
	})

	// Before relaxing the territories, we'd need to ensure that we only
	// relax without changing the borders of the empire...
	// So we'd only re-assign IDs that belong to the same territory.
	// m.rRelaxTerritories(m.r_city, 5)
}

// CityState represents a territory governed by a single city.
type CityState struct {
	ID      int     // Region where the city state originates
	Capital *City   // Capital city
	Cities  []*City // Cities within the city state

	// TODO: DO NOT CACHE THIS!
	Regions []int
	*Stats
}

func (c *CityState) Log() {
	log.Printf("The city state of %s: %d cities, %d regions", c.Capital.Name, len(c.Cities), len(c.Regions))
	c.Stats.Log()
}

func (m *Civ) GetCityStates() []*CityState {
	// TODO: Deduplicate with GetEmpires.
	var res []*CityState
	for i := 0; i < m.NumCityStates; i++ {
		c := &CityState{
			ID:      m.Cities[i].ID,
			Capital: m.Cities[i],
		}

		// Loop through all cities and gather all that
		// are within the current city state.
		for _, ct := range m.Cities {
			if m.RegionToCityState[ct.ID] == c.ID {
				c.Cities = append(c.Cities, ct)
			}
		}

		// Collect all regions that are part of the
		// current territory.
		for r, terr := range m.RegionToCityState {
			if terr == c.ID {
				c.Regions = append(c.Regions, r)
			}
		}
		c.Stats = m.getStats(c.Regions)
		c.Log()
		res = append(res, c)
	}
	return res
}

// getCityStateNeighbors returns all city states that are neighbors of the
// given city state.
func (m *Civ) getCityStateNeighbors(c *CityState) []int {
	return m.getTerritoryNeighbors(c.ID, m.RegionToCityState)
}
