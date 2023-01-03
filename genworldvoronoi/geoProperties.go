package genworldvoronoi

type RegProperty struct {
	ID                  int
	Elevation           float64 // 0.0-1.0
	Steepness           float64 // 0.0-1.0
	Biome               int     // biome of the city
	DistanceToCoast     float64 // graph distance to the nearest coast
	DistanceToMountain  float64 // graph distance to the nearest mountain
	DistanceToRiver     float64 // graph distance to the nearest river
	DistanceToVolcano   float64 // graph distance to the nearest volcano
	DistanceToFaultline float64 // graph distance to the nearest faultline
	Temperature         float64 // in °C
	Rainfall            float64 // in dm
	DangerRockslide     float64 // 0.0-1.0
	DangerEarthquake    float64 // 0.0-1.0
	DangerVolcano       float64 // 0.0-1.0
	DangerFlood         float64 // 0.0-1.0
	IsValley            bool    // true if the city is in a valley
}

// getRegPropertyFunc returns a function that returns the properties of a region.
// NOTE: This is probably a very greedy function.
func (m *Geo) getRegPropertyFunc() func(int) RegProperty {
	earthquakeChance := m.getEarthquakeChance()
	floodChance := m.getFloodChance()
	volcanoEruptionChance := m.getVolcanoEruptionChance()
	rockSlideAvalancheChance := m.getRockSlideAvalancheChance()
	steepness := m.GetSteepness()
	inlandValleyFunc := m.getFitnessInlandValleys()
	biomeFunc := m.getRegWhittakerModBiomeFunc()
	_, maxElev := minMax(m.Elevation)
	var oceanRegs, mountainRegs, volcanoRegs, riverRegs []int
	stopMountain := make(map[int]bool)
	stopOcean := make(map[int]bool)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.Elevation[r] < 0 {
			oceanRegs = append(oceanRegs, r)
			stopOcean[r] = true
		}
		if m.RegionIsMountain[r] {
			mountainRegs = append(mountainRegs, r)
			stopMountain[r] = true
		}
		if m.RegionIsVolcano[r] {
			volcanoRegs = append(volcanoRegs, r)
		}
		if m.isRegBigRiver(r) {
			riverRegs = append(riverRegs, r)
		}
	}
	distOcean := m.assignDistanceField(oceanRegs, stopMountain)
	distMountain := m.assignDistanceField(mountainRegs, stopOcean)
	distVolcano := m.assignDistanceField(volcanoRegs, stopOcean)
	distRiver := m.assignDistanceField(riverRegs, stopOcean)
	return func(id int) RegProperty {
		return RegProperty{
			ID:                  id,
			Elevation:           m.Elevation[id],
			Steepness:           steepness[id],
			Biome:               biomeFunc(id),
			DistanceToCoast:     distOcean[id],
			DistanceToMountain:  distMountain[id],
			DistanceToRiver:     distRiver[id],
			DistanceToVolcano:   distVolcano[id],
			DistanceToFaultline: m.RegionCompression[id],
			Temperature:         m.getRegTemperature(id, maxElev),
			Rainfall:            m.Rainfall[id],
			DangerRockslide:     rockSlideAvalancheChance[id],
			DangerEarthquake:    earthquakeChance[id],
			DangerVolcano:       volcanoEruptionChance[id],
			DangerFlood:         floodChance[id],
			IsValley:            inlandValleyFunc(id) > 0.8,
		}
	}
}
