package genworldvoronoi

import (
	"log"
)

// Empire contains information about a territory with the given ID.
// TODO: Maybe drop the regions since we can get that info
// relatively cheaply.
type Empire struct {
	ID       int // ID of the territory
	Name     string
	Emperor  string
	Capital  *City   // Capital city
	Cities   []*City // Cities within the territory
	Regions  []int   // Regions that are part of the empire
	Language *Language
	*Stats
}

func (e *Empire) Log() {
	log.Printf("The Empire of %s: %d cities, %d regions, capital: %s", e.Name, len(e.Cities), len(e.Regions), e.Capital.Name)
	log.Printf("Emperor: %s", e.Emperor)
	e.Stats.Log()
}

func (m *Map) GetEmpires() []*Empire {
	var res []*Empire
	for i := 0; i < m.NumTerritories; i++ {
		lang := GenLanguage(m.seed + int64(i))
		e := &Empire{
			ID:       m.cities_r[i].R,
			Name:     lang.MakeName(),
			Emperor:  lang.MakeFirstName() + " " + lang.MakeLastName(),
			Capital:  m.cities_r[i],
			Language: lang,
		}

		// TODO: Name empire, name cities.

		// Loop through all cities and gather all that
		// are within the current territory.
		for _, c := range m.cities_r {
			if m.r_territory[c.R] == e.ID {
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
