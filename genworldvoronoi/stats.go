package genworldvoronoi

import (
	"log"

	"github.com/Flokey82/go_gens/gameconstants"
	"github.com/Flokey82/go_gens/genbiome"
)

type Stats struct {
	NumRegions int
	ResMetal   [ResMaxMetals]int
	ResGems    [ResMaxGems]int
	ResStones  [ResMaxStones]int
	TotalArea  float64
	Biomes     map[int]int
	Desert     int
	Forest     int
	RainForest int
	Snow       int
	Swamp      int
	Wetlands   int
}

func (m *Map) getStats(rr []int) *Stats {
	st := &Stats{
		NumRegions: len(rr),
		Biomes:     make(map[int]int),
	}
	biomeFunc := m.getRWhittakerModBiomeFunc()
	for _, r := range rr {
		st.TotalArea += m.getRegionArea(r)
		for i := 0; i < ResMaxMetals; i++ {
			if m.r_res_metals[r]&(1<<i) != 0 {
				st.ResMetal[i]++
			}
		}
		for i := 0; i < ResMaxGems; i++ {
			if m.r_res_gems[r]&(1<<i) != 0 {
				st.ResGems[i]++
			}
		}
		for i := 0; i < ResMaxStones; i++ {
			if m.r_res_stone[r]&(1<<i) != 0 {
				st.ResStones[i]++
			}
		}
		b := biomeFunc(r)
		st.Biomes[b]++

		switch b {
		case genbiome.WhittakerModBiomeColdDesert, genbiome.WhittakerBiomeSubtropicalDesert:
			st.Desert++
		case genbiome.WhittakerModBiomeTropicalRainForest,
			genbiome.WhittakerModBiomeTemperateRainForest:
			st.RainForest++
		case genbiome.WhittakerModBiomeTropicalSeasonalForest,
			genbiome.WhittakerModBiomeTemperateSeasonalForest:
			st.Forest++
		case genbiome.WhittakerModBiomeSnow:
			st.Snow++
		case genbiome.WhittakerModBiomeHotSwamp:
			st.Swamp++
		case genbiome.WhittakerModBiomeWetlands:
			st.Wetlands++
		}
	}
	return st
}

func (s *Stats) Log() {
	log.Printf("Total Area: %.2f km2", s.TotalArea*gameconstants.EarthSurface/gameconstants.SphereSurface)
	for i := 0; i < ResMaxMetals; i++ {
		log.Printf("Metal %s: %d", metalToString(i), s.ResMetal[i])
	}
	for i := 0; i < ResMaxGems; i++ {
		log.Printf("Gem %s: %d", gemToString(i), s.ResGems[i])
	}
	for i := 0; i < ResMaxStones; i++ {
		log.Printf("Stone %s: %d", stoneToString(i), s.ResStones[i])
	}
	log.Printf("Desert: %.2f%%", 100*float64(s.Desert)/float64(s.NumRegions))
	log.Printf("RainForest: %.2f%%", 100*float64(s.RainForest)/float64(s.NumRegions))
	log.Printf("Forest: %.2f%%", 100*float64(s.Forest)/float64(s.NumRegions))
	log.Printf("Snow: %.2f%%", 100*float64(s.Snow)/float64(s.NumRegions))
	log.Printf("Swamp: %.2f%%", 100*float64(s.Swamp)/float64(s.NumRegions))
	log.Printf("Wetlands: %.2f%%", 100*float64(s.Wetlands)/float64(s.NumRegions))
}
