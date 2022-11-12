package genworldvoronoi

import (
	"log"
)

type CityState struct {
	ID      int // ID of the city.
	Capital *City
	Cities  []*City
	Regions []int
	*Stats
}

func (c *CityState) Log() {
	log.Printf("The city state of %s: %d cities, %d regions", c.Capital.Name, len(c.Cities), len(c.Regions))
	c.Stats.Log()
}

func (m *Map) GetCityStates() []*CityState {
	var res []*CityState
	for i := 0; i < m.NumCityStates; i++ {
		c := &CityState{
			ID:      m.cities_r[i].R,
			Capital: m.cities_r[i],
		}

		// Loop through all cities and gather all that
		// are within the current city state.
		for _, ct := range m.cities_r {
			if m.r_city[ct.R] == c.ID {
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
