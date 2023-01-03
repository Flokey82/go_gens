package genworldvoronoi

import (
	"log"
)

func (m *Civ) generateCitiesFlavorText() {
	rpFunc := m.getRegPropertyFunc()

	for _, c := range m.Cities {
		flvTxt := m.generateCityFlavorText(c, rpFunc(c.ID))
		log.Println(c.Name, flvTxt)
	}
}

// generateCityFlavorText generates a flavor text for a city.
func (m *Civ) generateCityFlavorText(c *City, p RegProperty) string {
	str := c.Name + " is a "
	if c.Population == 0 {
		str += "deserted "
		if c.MaxPopulation > 1000 {
			str += "town"
		} else {
			str += "village"
		}
	} else if c.Population < 100 {
		if c.MaxPopulation > 1000 {
			str += "desolate town"
		} else {
			str += "small village"
		}
	} else if c.Population < 1000 {
		str += "small town"
	} else if c.Population < 10000 {
		str += "large town"
	} else {
		str += "large city"
	}
	if p.IsValley && p.DistanceToCoast > 3 {
		str += " in a valley"
	} else if p.Steepness > 0.5 {
		if p.Elevation > 0.5 {
			str += " on a mountain"
		} else if p.DistanceToCoast <= 1 {
			str += " on a coastal cliff"
		} else {
			str += " on a hillside"
		}
	} else if p.DistanceToCoast <= 1 {
		str += " on the coast"
	}
	str += ".\n"

	// Generate some flavor text describing the region.
	str += m.GenerateRegPropertyDescription(p)

	// ... and finally add some flavor text for the biome.
	return str + generateFlavorTextForBiome(int64(c.ID), p.Biome)
}
