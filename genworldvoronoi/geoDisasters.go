package genworldvoronoi

import "math"

func (m *Geo) getEarthquakeChance() []float64 {
	// Get distance field from fault lines using the plate compression.
	compression := m.propagateCompression(m.RegionCompression)

	// Now get the chance of earthquake for each region.
	earthquakeChance := make([]float64, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		earthquakeChance[r] = math.Abs(compression[r])
	}
	return earthquakeChance
}

func (m *Geo) getFloodChance() []float64 {
	// Now get the chance of flood for each region.
	floodChance := make([]float64, m.mesh.numRegions)
	_, maxFlux := minMax(m.Flux)
	steepness := m.GetSteepness()
	for r := 0; r < m.mesh.numRegions; r++ {
		// We use the flux of water and the steepness in the region
		// to determine the chance of a flood.
		// NOTE: This should also apply to lakes.
		floodChance[r] = (1 - steepness[r]) * m.Flux[r] / maxFlux
	}

	// Normalize the flood chance.
	_, maxFloodChance := minMax(floodChance)
	for r := 0; r < m.mesh.numRegions; r++ {
		floodChance[r] /= maxFloodChance
	}
	return floodChance
}

func (m *Geo) getVolcanoEruptionChance() []float64 {
	return m.getDownhillDisaster(m.RegionIsVolcano, 0.05)
}

func (m *Geo) getRockSlideAvalancheChance() []float64 {
	return m.getDownhillDisaster(m.RegionIsMountain, 0.1)
}

func (m *Geo) getDownhillDisaster(origins map[int]bool, steepnessLimit float64) []float64 {
	steepness := m.GetSteepness()
	downhill := m.GetDownhill(true)

	// Start at the origin regions and go downhill until the terrain is too
	// flat or we reach the ocean.
	chance := make([]float64, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		if !origins[r] {
			continue
		}

		// Go downhill until the steepness is too low or we reach the ocean.
		rdh := r
		danger := 1.0
		for rdh != -1 && steepness[rdh] > steepnessLimit && m.Elevation[rdh] > 0 {
			// Add the danger of the region to the chance of being affected by a
			// downhill disaster.
			chance[rdh] += danger
			rdh = downhill[rdh]
			danger *= steepness[rdh]
		}
	}
	return chance
}

func (m *Civ) getDisasterFunc() func(r int) []disaster {
	// distRegion := math.Sqrt(4 * math.Pi / float64(m.mesh.numRegions))
	// biomeFunc := m.getRegWhittakerModBiomeFunc()
	_, maxElev := minMax(m.Elevation)
	var volcanoes, mountains, faultlines []int
	isBigRiver := make(map[int]bool)
	isFireDanger := make(map[int]bool)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.RegionIsMountain[r] {
			mountains = append(mountains, r)
		}
		if m.RegionIsVolcano[r] {
			volcanoes = append(volcanoes, r)
		}
		if m.RegionCompression[r] > 0 {
			faultlines = append(faultlines, r)
		}
		if m.isRegBigRiver(r) {
			isBigRiver[r] = true
		}
		// Determine if there is danger of fire by checking if the region is
		// hot and relatively dry while still having vegetation.
		temp := m.getRegTemperature(r, maxElev)
		if temp > 25 && m.Moisture[r] < 0.2 {
			isFireDanger[r] = true
		}
	}

	// Get distance field from volanoes.
	distVolcanoes := m.assignDistanceField(mountains, make(map[int]bool))
	// Get distance field from fault lines.
	distMountains := m.assignDistanceField(mountains, make(map[int]bool))
	// Get distance field from mountains.
	distFaultlines := m.assignDistanceField(faultlines, make(map[int]bool))

	// TODO: Instead, introduce a new property of disasters that determines
	// how likely they are to occur. Then, we can take in account how far
	// away the disaster is from the region.
	return func(regionID int) []disaster {
		// Now get the disasters that might affect the region.
		var ds []disaster
		// Check if the region is close to a volcano.
		if distVolcanoes[regionID] < 3 {
			ds = append(ds, disVolcano)
		}
		// Check if the region is close to a mountain.
		if distMountains[regionID] < 3 {
			ds = append(ds, disRockslide)
		}
		// Check if the region is close to a fault line.
		if distFaultlines[regionID] < 3 {
			ds = append(ds, disEarthquake)
		}
		// Check if the region is at a big river.
		if isBigRiver[regionID] {
			ds = append(ds, disFlood)
		}
		// Check if we have a fire danger.
		if isFireDanger[regionID] {
			ds = append(ds, disWildfire)
		}
		return ds
	}
}
