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
