package genworldvoronoi

import (
	"log"

	"github.com/Flokey82/go_gens/genbiome"
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
	if p.IsValley {
		str += " in a valley"
	} else if p.Steepness > 0.5 {
		if p.Elevation > 0.5 {
			str += " on a mountain"
		} else {
			str += " on a hillside"
		}
	}
	str += " in the " + genbiome.WhittakerModBiomeToString(p.Biome) + ".\n"

	// Generate some flavor text describing the environment based on temerature, humidity, and biome.

	if p.DistanceToVolcano < 3 {
		if p.DistanceToVolcano < 1 {
			str += " It is located on a volcano"
		} else {
			str += " It is located near a volcano"
		}
		if p.DangerVolcano > 0.2 {
			if p.DangerVolcano > 0.5 {
				str += " and is in constant danger of being destroyed by a volcanic eruption"
			} else {
				str += " and faces the threat of a possible volcanic eruption"
			}
		}
		str += ". "
	} else if p.DistanceToMountain < 3 {
		if p.DistanceToMountain < 1 {
			str += " It is located on a mountain"
		} else {
			str += " It is located near a mountain"
		}
		if p.DangerRockslide > 0.2 {
			if p.DangerRockslide > 0.5 {
				str += " and is in constant danger of being destroyed by a rockslide"
			} else {
				str += " and suffers occasional rockslides"
			}
		}
		str += ". "
	}
	if p.DistanceToRiver < 1 {
		str += " It is located near a river"
		if p.DangerFlood > 0.2 {
			if p.DangerFlood > 0.5 {
				str += " and is threatened by frequent floods."
			} else {
				str += " and experiences occasional floods."
			}
		}
		str += ". \n"
	}
	return str + generateFlavorTextForBiome(p.Biome)
}
