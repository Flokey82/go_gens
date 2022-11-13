package genworldvoronoi

import (
	"log"
)

// Empire contains information about a territory with the given ID.
// TODO: Maybe drop the regions since we can get that info
// relatively cheaply.
type Empire struct {
	ID       int     // Region where the empire originates (capital)
	Name     string  // Name of the empire
	Emperor  string  // Name of the ruler
	Capital  *City   // Capital city
	Cities   []*City // Cities within the territory
	Language *Language

	// TODO: DO NOT CACHE THIS!
	Regions []int // Regions that are part of the empire
	*Stats
}

func (e *Empire) Log() {
	log.Printf("The Empire of %s: %d cities, %d regions, capital: %s", e.Name, len(e.Cities), len(e.Regions), e.Capital.Name)
	log.Printf("Emperor: %s", e.Emperor)
	e.Stats.Log()
}

func (m *Map) GetEmpires() []*Empire {
	// TODO: Deduplicate with GetCityStates.
	var res []*Empire
	for i := 0; i < m.NumTerritories; i++ {
		capital := m.cities_r[i]
		var lang *Language
		if c := m.getRCulture(capital.ID); c != nil && c.Language != nil {
			lang = c.Language
		} else {
			lang = GenLanguage(m.seed + int64(i))
		}
		e := &Empire{
			ID:       capital.ID,
			Name:     lang.MakeName(),
			Emperor:  lang.MakeFirstName() + " " + lang.MakeLastName(),
			Capital:  capital,
			Language: lang,
		}

		// Loop through all cities and gather all that
		// are within the current territory.
		for _, c := range m.cities_r {
			if m.r_territory[c.ID] == e.ID {
				// TODO: Name cities based on local culture?
				c.Name = e.Language.MakeCityName()
				e.Cities = append(e.Cities, c)
			}
		}

		// Collect all regions that are part of the
		// current territory.
		for r, terr := range m.r_territory {
			if terr == e.ID {
				e.Regions = append(e.Regions, r)
			}
		}
		e.Stats = m.getStats(e.Regions)
		e.Log()
		res = append(res, e)
	}
	return res
}
