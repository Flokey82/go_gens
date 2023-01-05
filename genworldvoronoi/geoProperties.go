package genworldvoronoi

type RegProperty struct {
	ID                  int
	Elevation           float64 // 0.0-1.0
	Steepness           float64 // 0.0-1.0
	Biome               int     // biome of the region
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
	IsValley            bool    // true if the region is a valley
	OnIsland            bool    // true if the region is on an island
}

// getRegPropertyFunc returns a function that returns the properties of a region.
// NOTE: This is probably a very greedy function.
func (m *Geo) getRegPropertyFunc() func(int) RegProperty {
	// TODO: Add chance of tropical storms, wildfires, etc.
	earthquakeChance := m.getEarthquakeChance()
	floodChance := m.getFloodChance()
	volcanoEruptionChance := m.getVolcanoEruptionChance()
	rockSlideAvalancheChance := m.getRockSlideAvalancheChance()
	steepness := m.GetSteepness()
	inlandValleyFunc := m.getFitnessInlandValleys()
	biomeFunc := m.getRegWhittakerModBiomeFunc()
	_, maxElev := minMax(m.Elevation)
	var oceanRegs, mountainRegs, volcanoRegs, riverRegs, faultlineRegs []int
	stopMountain := make(map[int]bool)
	stopOcean := make(map[int]bool)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.Elevation[r] <= 0 {
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
		if m.RegionCompression[r] != 0 {
			faultlineRegs = append(faultlineRegs, r)
		}
	}
	distOcean := m.assignDistanceField(oceanRegs, stopMountain)
	distMountain := m.assignDistanceField(mountainRegs, stopOcean)
	distVolcano := m.assignDistanceField(volcanoRegs, stopOcean)
	distRiver := m.assignDistanceField(riverRegs, stopOcean)
	distFaultline := m.assignDistanceField(faultlineRegs, stopOcean)
	return func(id int) RegProperty {
		// Make sure that we do not have more than 2 neighbours that has a lower elevation.
		// ... because a valley should be surrounded by mountains.
		isValley := inlandValleyFunc(id) > 0.8
		if isValley {
			var count int
			for _, n := range m.GetRegNeighbors(id) {
				if m.Elevation[n] > m.Elevation[id] {
					continue
				}
				count++
				if count > 2 {
					isValley = false
					break
				}
			}
		}

		return RegProperty{
			ID:                  id,
			Elevation:           m.Elevation[id],
			Steepness:           steepness[id],
			Biome:               biomeFunc(id),
			DistanceToCoast:     distOcean[id],
			DistanceToMountain:  distMountain[id],
			DistanceToRiver:     distRiver[id],
			DistanceToVolcano:   distVolcano[id],
			DistanceToFaultline: distFaultline[id],
			Temperature:         m.getRegTemperature(id, maxElev),
			Rainfall:            m.Rainfall[id],
			DangerRockslide:     rockSlideAvalancheChance[id],
			DangerEarthquake:    earthquakeChance[id],
			DangerVolcano:       volcanoEruptionChance[id],
			DangerFlood:         floodChance[id],
			IsValley:            isValley,
			OnIsland:            m.LandmassSize[m.Landmasses[id]] < 15, // TODO: This should use actual geographical area.
		}
	}
}
