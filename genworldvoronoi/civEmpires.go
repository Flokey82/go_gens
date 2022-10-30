package genworldvoronoi

import (
	"log"

	"github.com/Flokey82/go_gens/gameconstants"
)

// Empire contains information about a territory with the given ID.
// TODO: Maybe drop the regions since we can get that info
// relatively cheaply.
type Empire struct {
	ID        int // ID of the territory
	Name      string
	Emperor   string
	Capital   *City   // Capital city
	Cities    []*City // Cities within the territory
	Regions   []int   // Regions that are part of the empire
	Language  *Language
	ResMetal  [ResMaxMetals]int
	ResGems   [ResMaxGems]int
	TotalArea float64
}

func (e *Empire) Log() {
	log.Printf("The Empire of %s: %d cities, %d regions, capital: %s", e.Name, len(e.Cities), len(e.Regions), e.Capital.Name)
	log.Printf("Total Area: %.2f km2", e.TotalArea*gameconstants.EarthSurface/gameconstants.SphereSurface)
	log.Printf("Emperor: %s", e.Emperor)
	for i := 0; i < ResMaxMetals; i++ {
		log.Printf("Metal %s: %d", metalToString(i), e.ResMetal[i])
	}
	for i := 0; i < ResMaxGems; i++ {
		log.Printf("Gem %s: %d", gemToString(i), e.ResGems[i])
	}
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
				e.TotalArea += m.getRegionArea(r)
				for i := 0; i < ResMaxMetals; i++ {
					if m.r_res_metals[r]&(1<<i) != 0 {
						e.ResMetal[i]++
					}
				}
				for i := 0; i < ResMaxGems; i++ {
					if m.r_res_gems[r]&(1<<i) != 0 {
						e.ResGems[i]++
					}
				}
				e.Regions = append(e.Regions, r)
			}
		}
		e.Log()
		res = append(res, e)
	}
	return res
}
