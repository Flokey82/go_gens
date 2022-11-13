package genworldvoronoi

import (
	"log"
)

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

func (m *Map) GetCityStates() []*CityState {
	// TODO: Deduplicate with GetEmpires.
	var res []*CityState
	for i := 0; i < m.NumCityStates; i++ {
		c := &CityState{
			ID:      m.cities_r[i].ID,
			Capital: m.cities_r[i],
		}

		// Loop through all cities and gather all that
		// are within the current city state.
		for _, ct := range m.cities_r {
			if m.r_city[ct.ID] == c.ID {
				c.Cities = append(c.Cities, ct)
			}
		}

		// Collect all regions that are part of the
		// current territory.
		for r, terr := range m.r_city {
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

func (m *Map) getCityStateNeighbors(c *CityState) []int {
	return m.getRTerritoryNeighbors(c.ID, m.r_city)
}
