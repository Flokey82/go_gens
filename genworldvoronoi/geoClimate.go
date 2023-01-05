package genworldvoronoi

import (
	"container/list"
	"image/color"
	"log"
	"math"
	"sort"

	"github.com/Flokey82/go_gens/gameconstants"
	"github.com/Flokey82/go_gens/genbiome"
	"github.com/Flokey82/go_gens/vectors"
)

// getRegWhittakerModBiomeFunc returns a function that returns the Whittaker biome
// for a given region.
func (m *Geo) getRegWhittakerModBiomeFunc() func(r int) int {
	_, maxElev := minMax(m.Elevation)
	_, maxMois := minMax(m.Moisture)
	return func(r int) int {
		valElev := m.Elevation[r] / maxElev
		valMois := m.Moisture[r] / maxMois
		regLat := m.LatLon[r][0]
		return getWhittakerModBiome(regLat, valElev, valMois)
	}
}

func getWhittakerModBiome(latitude, elevation, moisture float64) int {
	return genbiome.GetWhittakerModBiome(int(getMeanAnnualTemp(latitude)-getTempFalloffFromAltitude(maxAltitudeFactor*elevation)), int(moisture*maxPrecipitation))
}

func getWhittakerModBiomeColor(latitude, elevation, moisture, intensity float64) color.NRGBA {
	return genbiome.GetWhittakerModBiomeColor(int(getMeanAnnualTemp(latitude)-getTempFalloffFromAltitude(maxAltitudeFactor*elevation)), int(moisture*maxPrecipitation), intensity)
}

// getTempFalloffFromAltitude returns the temperature falloff at a given altitude in meters
// above sea level. (approx. 9.8 °C per 1000 m)
// NOTE: This is definitely not correct :)
// Source: https://www.quora.com/At-what-rate-does-temperature-drop-with-altitude
func getTempFalloffFromAltitude(height float64) float64 {
	if height < 0 {
		return 0.0
	}
	return gameconstants.EarthElevationTemperatureFalloff * height
}

const (
	minTemp          = genbiome.MinTemperatureC
	maxTemp          = genbiome.MaxTemperatureC
	rangeTemp        = maxTemp - minTemp
	maxPrecipitation = genbiome.MaxPrecipitationDM // 450cm
)

// getMeanAnnualTemp returns the temperature at a given latitude within the range of
// -15 °C to +30°C because that's the range in which the Whittaker biomes are defined.
// For this I assume that light hits the globe exactly from a 90° angle with respect
// to the planitary axis.
// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals (facing ratio)
// See: http://www-das.uwyo.edu/~geerts/cwx/notes/chap16/geo_clim.html
// NOTE: -35 °C to +31 °C would be ideally the temp gradient (according to real-life data), but we don't have (yet) any biomes defined for this.
func getMeanAnnualTemp(lat float64) float64 {
	return (math.Sin(degToRad(90-math.Abs(lat))))*rangeTemp + minTemp
}

const maxAltitudeFactor = gameconstants.EarthMaxElevation // How tall is the tallest mountain with an elevation of 1.0?

// getRegTemperature returns the average yearly temperature of the given region at the surface.
func (m *Geo) getRegTemperature(r int, maxElev float64) float64 {
	// TODO: Fix maxElev caching!!!
	return getMeanAnnualTemp(m.LatLon[r][0]) - getTempFalloffFromAltitude(maxAltitudeFactor*m.Elevation[r]/maxElev)
}

// getRTemperature returns the average yearly temperature of the given triangle at the surface.
func (m *Geo) getTriTemperature(t int, maxElev float64) float64 {
	// TODO: Fix maxElev caching!!!
	return getMeanAnnualTemp(m.triLatLon[t][0]) - getTempFalloffFromAltitude(maxAltitudeFactor*m.triElevation[t]/maxElev)
}

// getGlobalWindVector returns a vector for the global wind at the given latitude.
// NOTE: This is based on the trade winds on... well, earth.
// See: https://en.wikipedia.org/wiki/Trade_winds
func getGlobalWindVector(lat float64) Vertex {
	// Based on latitude, we calculate the wind vector angle.
	var degree float64
	if latAbs := math.Abs(lat); latAbs >= 0 && latAbs <= 30 {
		// +30° ... 0°, 0° ... -30° -> Primitive Hadley Cell.
		// In a Hadley cell, we turn the wind vector until we are exactly parallel with the equator once we reach 0° Lat.
		// TODO: This is probably not perfectly parallel at the equator.
		change := 90 * latAbs / 30
		if lat > 0 {
			degree = 180 + change // Northern hemisphere.
		} else {
			degree = 180 - change // Southern hemisphere.
		}
	} else if latAbs > 30 && latAbs <= 60 {
		// +60° ... +30°, -30° ... -60° -> Primitive Mid Latitude Cell.
		// In a mid latitude cell, we turn the wind vector until we are exactly parallel with the 60° Lat.
		// TODO: This is probably not a full 90° turn. Fix this
		change := 90 * (latAbs - 30) / 30
		if lat > 0 {
			degree = 90 - change // Northern hemisphere.
		} else {
			degree = 270 + change // Southern hemisphere.
		}
	} else {
		// NOTE: This is buggy or at least "not nice".
		// +90° ... +60°, -60° ... -90° -> Primitive Hadley Cell.
		// In a polar cell, we turn the wind vector until we are exactly parallel with the equator once we reach 60° Lat.
		change := 90 * (latAbs - 60) / 30
		if lat > 0 {
			degree = 180 + change // Northern hemisphere.
		} else {
			degree = 180 - change // Southern hemisphere.
		}
	}
	rad := degToRad(degree)
	return Vertex{math.Cos(rad), math.Sin(rad)}
}

// assignWindVectors constructs faux global wind cells reminiscent of a simplified earth model.
// NOTE: This function includes an experimental part that calculates local winds that are influenced
// by the topography / elevation changes. Please note that the code for local winds is incomplete.
func (m *Geo) assignWindVectors() {
	// Based on latitude of each region, we calculate the wind vector.
	regWindVec := make([]Vertex, m.mesh.numRegions)
	for i := range regWindVec {
		regWindVec[i] = getGlobalWindVector(m.LatLon[i][0])
	}

	// Local wind vectors.
	regWindVecLocal := make([]Vertex, m.mesh.numRegions)
	_, maxElev := minMax(m.Elevation)

	// Experimental: Local wind vectors based on temperature gradients.
	useTempGradient := false

	// NOTE: This is currently overridden by the altitude changes below.
	if useTempGradient {
		// Add local wind vectors based on local temperature gradients.
		//
		// NOTE: You won't be happy about the results of the temp gradient anyway
		// since everything becomes quite "patchy".
		//
		// In plain English: This is garbage code :(
		// I suspect that the wind is deflected too much by minimal temperature changes
		// and I am too lazy to really look into it.

		// Determine all sea regions.
		var seaRegs []int
		for r := 0; r < m.mesh.numRegions; r++ {
			if m.Elevation[r] <= 0 {
				seaRegs = append(seaRegs, r)
			}
		}
		regDistanceSea := m.assignDistanceField(seaRegs, make(map[int]bool))
		for r := range regWindVecLocal {
			regVec := regWindVec[r]
			lat := m.LatLon[r][0]
			lon := m.LatLon[r][1]
			tempReg := getMeanAnnualTemp(lat) - getTempFalloffFromAltitude(8850*m.Elevation[r]/maxElev)
			if m.Elevation[r] < 0 {
				// TODO: Use actual distance from ocean to calculate temperature falloff.
				tempReg -= 1 / (regDistanceSea[r] + 1)
			}
			// Get temperature for r.
			v := vectors.Normalize(vectors.Vec2{
				X: regVec[0],
				Y: regVec[1],
			})
			for _, nb := range m.GetRegNeighbors(r) {
				nbLat := m.LatLon[nb][0]
				nbLon := m.LatLon[nb][1]
				tempNb := getMeanAnnualTemp(nbLat) - getTempFalloffFromAltitude(8850*m.Elevation[nb]/maxElev)
				if m.Elevation[nb] < 0 {
					// TODO: Use actual distance from ocean to calculate temperature falloff.
					tempNb -= 1 / (regDistanceSea[nb] + 1)
				}
				ve := calcVecFromLatLong(lat, lon, nbLat, nbLon)
				v = v.Add(vectors.Normalize(vectors.NewVec2(ve[0], ve[1])).Mul(tempNb - tempReg))
			}
			v = vectors.Normalize(v)
			regWindVecLocal[r] = Vertex{v.X, v.Y}
		}
	} else {
		// Add wind deflection based on altitude changes.
		for r := range regWindVecLocal {
			regVec := regWindVec[r]
			// Get XYZ Position of r.
			regXYZ := convToVec3(m.XYZ[r*3 : r*3+3])
			// Convert to polar coordinates.
			regLat := m.LatLon[r][0]
			regLon := m.LatLon[r][1]
			h := m.Elevation[r]
			if h < 0 {
				h = 0
			}
			// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
			rwXYZ := convToVec3(latLonToCartesian(regLat+regWindVec[r][1], regLon+regWindVec[r][0])).Normalize()
			v := vectors.Normalize(vectors.Vec2{
				X: regVec[0],
				Y: regVec[1],
			}) // v.Mul(h / maxElev)
			vw := calcVecFromLatLong(regLat, regLon, regLat+regWindVec[r][1], regLon+regWindVec[r][0])
			v0 := vectors.Normalize(vectors.Vec2{
				X: vw[0],
				Y: vw[1],
			})

			// Calculate Vector between r and wind_r.
			vb := vectors.Sub3(rwXYZ, regXYZ).Normalize()

			for _, nbReg := range m.GetRegNeighbors(r) {
				// if is_sea[neighbor_r] {
				//	continue
				// }
				// Calculate dot product of wind vector to vector r -> neighbor_r.
				// Get XYZ Position of r_neighbor.
				rnXYZ := convToVec3(m.XYZ[nbReg*3 : nbReg*3+3])

				// Calculate Vector between r and neighbor_r.
				va := vectors.Sub3(rnXYZ, regXYZ).Normalize()

				// Calculate dot product between va and vb.
				// This will give us how much the current region lies within the wind direction of the
				// current neighbor.
				// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
				dotV := vectors.Dot3(va, vb)
				hnb := m.Elevation[nbReg]
				if hnb < 0 {
					hnb = 0
				}
				if dotV > 0 {
					nbLat := m.LatLon[nbReg][0]
					nbLon := m.LatLon[nbReg][1]
					ve := calcVecFromLatLong(regLat, regLon, nbLat, nbLon)
					vx := vectors.Normalize(v0.Sub(vectors.Normalize(vectors.Vec2{
						X: ve[0],
						Y: ve[1],
					})))
					// The higher the dot product (the more direct the neighbor is in wind direction), the higher
					// the influence of an elevation change. So a steep mountain ahead will slow the wind down.
					// If a steep mountain is to the left, the wind vector will be pushed to the right.
					v = v.Add(vx.Mul(dotV * (hnb - h) / maxElev))
				}
			}
			v = vectors.Normalize(v)
			regWindVecLocal[r] = Vertex{v.X, v.Y}
		}
	}
	// Average wind vectors using neighbor vectors.
	interpolationSteps := 0

	m.RegionToWindVecLocal = m.interpolateWindVecs(regWindVecLocal, interpolationSteps)
	m.RegionToWindVec = m.interpolateWindVecs(regWindVec, interpolationSteps)
}

// interpolateWindVecs interpolates the given wind vectors at their respective regions by
// mixing them with the wind vectors of their neighbor regions.
func (m *Geo) interpolateWindVecs(in []Vertex, steps int) []Vertex {
	// Average wind vectors using neighbor vectors.
	for i := 0; i < steps; i++ {
		regWindVecInterpolated := make([]Vertex, m.mesh.numRegions)
		for r := range regWindVecInterpolated {
			resVec := Vertex{
				in[r][0],
				in[r][1],
			}
			var count int
			for _, nbReg := range m.GetRegNeighbors(r) {
				resVec[0] += in[nbReg][0]
				resVec[1] += in[nbReg][1]
				count++
			}
			resVec[0] /= float64(count + 1)
			resVec[1] /= float64(count + 1)
			regWindVecInterpolated[r] = resVec
		}
		in = regWindVecInterpolated
	}
	return in
}

type biomesParams struct {
	raininess   float64 // 0, 2
	rainShadow  float64 // 0.1, 2
	evaporation float64 // 0, 1
}

const (
	moistTransferDirect   = 0
	moistTransferIndirect = 1
	moistOrderWind        = 0
	moistOrderOther       = 1
)

// assignRainfall is an overengineered logic that is supposed to calculate the transfer
// of moisture across the globe based on global winds using distinct approaches.
// Unfortunately, this is highly bugged and not as useful as the simpler version
// 'assignRainfallBasic'.
func (m *Geo) assignRainfall(numSteps, transferMode, sortOrder int) {
	biomesParam := biomesParams{
		raininess:   0.9,
		rainShadow:  0.9,
		evaporation: 0.9,
	}

	// 1. Initialize
	// 1.1. Determine all sea regions.
	var seaRegs, landRegs []int
	isSea := make([]bool, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.Elevation[r] < 0 {
			isSea[r] = true
			seaRegs = append(seaRegs, r)
		} else {
			landRegs = append(landRegs, r)
		}
	}

	var sortOrderRegs []int
	if sortOrder == moistOrderOther {
		// 1.2. Sort all regions by distance to ocean. Lowest to highest.
		distOrderRegs := make([]int, m.mesh.numRegions)
		for r := 0; r < m.mesh.numRegions; r++ {
			distOrderRegs[r] = r
		}
		regDistanceSea := m.assignDistanceField(seaRegs, make(map[int]bool))
		sort.Slice(distOrderRegs, func(a, b int) bool {
			if regDistanceSea[distOrderRegs[a]] == regDistanceSea[distOrderRegs[b]] {
				return m.Elevation[distOrderRegs[a]] < m.Elevation[distOrderRegs[b]]
			}
			return regDistanceSea[distOrderRegs[a]] < regDistanceSea[distOrderRegs[b]]
		})
		sortOrderRegs = distOrderRegs
	} else {
		// 1.2. Sort the indices in wind-order so we can ensure that we push the moisture
		// in their logical sequence across the globe.
		_, sortOrderRegs = m.getWindSortOrder() // Works reasonably well.
	}

	// 1.3. Assign wind vector for every region
	m.assignWindVectors()
	regWindVec := m.RegionToWindVec
	_, maxH := minMax(m.Elevation)

	calcRainfall := func(r int, humidity float64) float64 {
		regElev := m.Elevation[r]
		if regElev < 0 {
			regElev = 0 // Set to sea-level
		}
		heightVal := 1 - (regElev / maxH)
		if humidity > heightVal {
			return biomesParam.rainShadow * (humidity - heightVal)
		}
		return 0
	}

	for step := 0; step < numSteps; step++ {
		log.Println(step)
		// Evaporation.

		// 2. Assign initial moisture of 1.0 to all regions below or at sea level or replenish
		// moisture through evaporation if our moisture is below 0.
		for _, r := range seaRegs {
			if m.Moisture[r] < 1.0 {
				m.Moisture[r] = 1.0
			}
			// m.r_rainfall[r] += biomesParam.raininess * m.r_moisture[r]
		}

		// Rivers should experience some evaporation.
		for r, fluxval := range m.Flux {
			if m.Moisture[r] < fluxval && m.Moisture[r] < 1.0 {
				m.Moisture[r] = 1.0 // TODO: Should depend on available water.
			}
		}

		// Water pools should experience some evaporation.
		for r, poolval := range m.Waterpool {
			if poolval > 0 && m.Moisture[r] < 1.0 {
				m.Moisture[r] = 1.0 // TODO: Should depend on available water.
			}
		}
		// m.interpolateRainfallMoisture(1)

		// 3. Transfer moisture based on wind vectors.
		switch transferMode {
		case moistTransferDirect:
			// 3.1.B For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r.
			//       This will give us the amount of moisture we transfer to the neighbor region.
			// NOTE: This variant copies moisture from the current region to the neighbors that are in wind direction.
			for _, r := range sortOrderRegs {
				count := 0
				// Get XYZ Position of r.
				regXYZ := convToVec3(m.XYZ[r*3 : r*3+3])
				// Convert to polar coordinates.
				regLat := m.LatLon[r][0]
				regLon := m.LatLon[r][1]

				// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
				regToWindVec3 := convToVec3(latLonToCartesian(regLat+regWindVec[r][1], regLon+regWindVec[r][0])).Normalize()
				for _, nbReg := range m.GetRegNeighbors(r) {
					if isSea[nbReg] {
						continue
					}
					// Calculate dot product of wind vector to vector r -> neighbor_r.
					// Get XYZ Position of r_neighbor.
					regToNbVec3 := convToVec3(m.XYZ[nbReg*3 : nbReg*3+3])

					// Calculate Vector between r and neighbor_r.
					va := vectors.Sub3(regToNbVec3, regXYZ).Normalize()

					// Calculate Vector between r and wind_r.
					vb := vectors.Sub3(regToWindVec3, regXYZ).Normalize()

					// Calculate dot product between va and vb.
					// This will give us how much the current region lies within the wind direction of the
					// current neighbor.
					// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
					dotV := vectors.Dot3(va, vb)
					if dotV > 0 {
						// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
						count++
						humidity := m.Moisture[nbReg] + m.Moisture[r]*dotV
						rainfall := m.Rainfall[nbReg] // + biomesParam.raininess*m.r_moisture[r]*dotV
						orographicRainfall := calcRainfall(nbReg, humidity)
						if orographicRainfall > 0.0 {
							rainfall += biomesParam.raininess * orographicRainfall
							humidity -= orographicRainfall
						}
						// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
						// WARNING: The humidity calculation is off.
						// humidity = math.Min(humidity, 1.0)
						// rainfall = math.Min(rainfall, 1.0)
						m.Rainfall[nbReg] = rainfall
						m.Moisture[nbReg] = humidity
					}
				}
			}
		case moistTransferIndirect:
			// 3.2. For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r_neighbor.
			//    This will give us the amount of moisture we transfer from the neighbor region.
			// NOTE: This variant copies moisture to the current region from the neighbors depending on their wind direction.
			for _, r := range sortOrderRegs {
				count := 0
				sum := 0.0
				// Get XYZ Position of r as vector3
				regVec3 := convToVec3(m.XYZ[r*3 : r*3+3])
				for _, nbReg := range m.GetRegNeighbors(r) {
					// Calculate dot product of wind vector to vector r -> neighbor_r.
					// Get XYZ Position of r_neighbor.
					regToNbVec3 := convToVec3(m.XYZ[nbReg*3 : nbReg*3+3])

					// Convert to polar coordinates.
					rLat := m.LatLon[nbReg][0]
					rLon := m.LatLon[nbReg][1]

					// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
					nbToWindVec3 := convToVec3(latLonToCartesian(rLat+regWindVec[nbReg][1], rLon+regWindVec[nbReg][0])).Normalize()

					// Calculate Vector between r and neighbor_r.
					va := vectors.Sub3(regVec3, regToNbVec3).Normalize()

					// Calculate Vector between neightbor_r and wind_neighbor_r.
					vb := vectors.Sub3(nbToWindVec3, regToNbVec3).Normalize()

					// Calculate dot product between va and vb.
					// This will give us how much the current region lies within the wind direction of the
					// current neighbor.
					// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
					dotV := vectors.Dot3(va, vb)
					if dotV > 0 {
						// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
						count++
						sum += m.Moisture[nbReg] * dotV
					}
				}

				var humidity, rainfall float64
				humidity = m.Moisture[r]
				if count > 0 {
					// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
					// WARNING: The humidity calculation is off.
					humidity = math.Min(humidity+sum, 1.0) // / float64(count)
					rainfall = math.Min(rainfall+biomesParam.raininess*sum, 1.0)
				}
				if m.Elevation[r] <= 0.0 {
					// evaporation := biomesParam.evaporation * (-m.r_elevation[r])
					// humidity = evaporation
					humidity = m.Moisture[r]
				}
				orographicRainfall := calcRainfall(r, humidity)
				if orographicRainfall > 0.0 {
					rainfall += biomesParam.raininess * orographicRainfall
					humidity -= orographicRainfall
				}
				m.Rainfall[r] = rainfall
				m.Moisture[r] = humidity
			}
		}

		// 4. Average moisture and rainfall.
		// m.interpolateRainfallMoisture(1)
	}
}

func (m *Geo) getWindSortOrder() ([]float64, []int) {
	m.assignWindVectors()
	useAlternativeWindSort := true

	// TODO: Add bool parameter to switch between local winds and global winds.
	// regWindVec := m.regWindVec
	regWindVec := m.RegionToWindVecLocal
	windOrderRegs := make([]int, m.mesh.numRegions)
	regWindSort := make([]float64, m.mesh.numRegions)
	// Sort all regions by latitude and longitude and their wind vector.
	// This will give us a logical order in which we can push the moisture across the globe.
	if useAlternativeWindSort {
		for r := 0; r < m.mesh.numRegions; r++ {
			// Get XYZ Position of r as vector3
			regVec3 := convToVec3(m.XYZ[r*3 : r*3+3])
			// Get XYZ Position of r_neighbor.
			regToWindVec3 := convToVec3(latLonToCartesian(m.LatLon[r][0]+regWindVec[r][1], m.LatLon[r][1]+regWindVec[r][0])).Normalize()
			// Calculate Vector between r and neighbor_r.
			va := vectors.Sub3(regVec3, regToWindVec3).Normalize()
			// Calculate dot product between va and vb.
			// This will give us how much the current region lies within the wind direction of the
			// current neighbor.
			// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
			dotV := vectors.Dot3(va, regToWindVec3)
			regWindSort[r] = dotV
			windOrderRegs[r] = r
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(regWindSort)))
		sort.Sort(sort.Reverse(sort.IntSlice(windOrderRegs)))
	} else {
		for r := 0; r < m.mesh.numRegions; r++ {
			windOrderRegs[r] = r
			// TODO: modify the sort order by ensuring longitude wraps around...??
			lat := (m.LatLon[r][0]) * regWindVec[r][1] / math.Abs(regWindVec[r][1]) // radToDeg(r_windvec[r][1])
			lon := (m.LatLon[r][1]) * regWindVec[r][0] / math.Abs(regWindVec[r][0]) // radToDeg(r_windvec[r][0])
			regWindSort[r] = (lat + lon)
		}

		// Sort the indices in wind-order so we can ensure that we push the moisture
		// in their logical sequence across the globe.
		sort.Slice(windOrderRegs, func(a, b int) bool {
			return regWindSort[windOrderRegs[a]] < regWindSort[windOrderRegs[b]]
		})
	}
	return regWindSort, windOrderRegs
}

func (m *Geo) assignRainfallBasic() {
	// NOTE: This still has issues with the wrap around at +/- 180° long
	biomesParam := biomesParams{
		raininess:   0.9,
		rainShadow:  0.9,
		evaporation: 0.9,
	}
	humidityFromRiver := 1.0
	humidityFromSea := 1.0
	humidityFromPool := 1.0
	evaporateRivers := true
	evaporatePools := false

	_, maxFlux := minMax(m.Flux)
	_, maxPool := minMax(m.Waterpool)
	minElev, maxElev := minMax(m.Elevation)
	if minElev == 0 {
		minElev = 1
	}

	// Sort the indices in wind-order so we can ensure that we push the moisture
	// in their logical sequence across the globe.
	_, windOrderRegs := m.getWindSortOrder()
	regWindVec := m.RegionToWindVecLocal

	// calcRainfall returns the amount of rain shed given the region and humidity.
	calcRainfall := func(r int, humidity float64) float64 {
		elev := m.Elevation[r]
		if elev < 0 {
			elev = 0 // Set to sea-level
		}
		heightVal := 1 - (elev / maxElev)
		if humidity > heightVal {
			return biomesParam.rainShadow * (humidity - heightVal)
		}
		return 0
	}

	// Evaporation.
	// 1. Assign initial moisture of 1.0 to all regions below or at sea level or replenish
	// moisture through evaporation if our moisture is below 0.
	for r, h := range m.Elevation {
		if h <= 0 {
			m.Moisture[r] = math.Max(m.Moisture[r], humidityFromSea)
		}
	}

	// Rivers should experience some evaporation.
	if evaporateRivers {
		for r, fluxval := range m.Flux {
			if m.isRegBigRiver(r) {
				evaporation := humidityFromRiver * fluxval / maxFlux
				m.Moisture[r] = math.Max(m.Moisture[r], evaporation)
			}
		}
	}

	// Water pools should experience some evaporation.
	//
	// NOTE: Currently this is not used since flood algorithms are deactivated so
	// the value for water pools is always 0.
	if evaporatePools {
		for r, poolval := range m.Waterpool {
			if poolval > 0 {
				evaporation := humidityFromPool * poolval / maxPool
				m.Moisture[r] = math.Max(m.Moisture[r], evaporation)
			}
		}
	}

	// Visit regions in wind order and copy the moisture from the neighbor regious that are
	// up-wind.
	//
	// NOTE: Since we start and stop at +- 180° long, we need to run the code several times
	// to ensure that moisture is pushed across the longitude wrap-around.
	for i := 0; i < 4; i++ {
		for _, r := range windOrderRegs {
			var humidity float64

			// Calculate humidity.
			for _, nbReg := range m.GetRegNeighbors(r) {
				rL := m.LatLon[r]
				nL := m.LatLon[nbReg]

				// TODO: Check dot product of wind vector (r) and neighbour->r.
				vVec := normal2(calcVecFromLatLong(nL[0], nL[1], nL[0]+regWindVec[nbReg][1], nL[1]+regWindVec[nbReg][0]))
				nVec := normal2(calcVecFromLatLong(nL[0], nL[1], rL[0], rL[1]))
				dotV := dot2(vVec, nVec)

				// Check if the neighbor region is up-wind (that the wind blows from neighbor_r to r) / dotV is positive.
				if dotV > 0.0 {
					humidity += m.Moisture[nbReg] * dotV
				}
			}

			// Evaporation.
			if m.Elevation[r] <= 0 {
				evaporation := biomesParam.evaporation * humidityFromSea * m.Elevation[r] / minElev
				humidity = math.Max(humidity, evaporation)
			} else if evaporateRivers && m.isRegBigRiver(r) {
				evaporation := biomesParam.evaporation * humidityFromRiver * m.Flux[r] / maxFlux
				humidity = math.Max(humidity, evaporation)
			} else if evaporatePools && m.Waterpool[r] > 0 {
				evaporation := biomesParam.evaporation * humidityFromPool * m.Waterpool[r] / maxPool
				humidity = math.Max(humidity, evaporation)
			}

			// Calculate orographic rainfall caused by elevation changes.
			rainfall := biomesParam.raininess * calcRainfall(r, humidity)
			m.Rainfall[r] = rainfall
			m.Moisture[r] = humidity - rainfall
		}
	}
	m.interpolateRainfallMoisture(5)
}

func (m *Geo) interpolateRainfallMoisture(interpolationSteps int) {
	for i := 0; i < interpolationSteps; i++ {
		regMoistureInterpol := make([]float64, m.mesh.numRegions)
		regRainfallInterpol := make([]float64, m.mesh.numRegions)
		for r := range regMoistureInterpol {
			rMoist := m.Moisture[r]
			rRain := m.Rainfall[r]
			var count int
			for _, nbReg := range m.GetRegNeighbors(r) {
				// Gravity! Water moves downwards.
				// This is not super-accurate since you'd have to take
				// in account how steep the slope is etc.
				if m.Elevation[r] >= m.Elevation[nbReg] {
					continue
				}
				rMoist += m.Moisture[nbReg]
				rRain += m.Rainfall[nbReg]
				count++
			}
			regMoistureInterpol[r] = rMoist / float64(count+1)
			regRainfallInterpol[r] = rRain / float64(count+1)
		}
		m.Moisture = regMoistureInterpol
		m.Rainfall = regRainfallInterpol
	}
}

func (m *Geo) assignBiomeRegions() {
	// Identify connected regions with the same biome.
	m.BiomeRegions = m.identifyBiomeRegions()

	regSize := make(map[int]int)
	for _, lm := range m.BiomeRegions {
		if lm >= 0 {
			regSize[lm]++ // Only count regions that are set to a valid ID.
		}
	}
	m.BiomeRegionSize = regSize
}

// identifyBiomeRegions identifies connected regions with the same biome.
func (m *Geo) identifyBiomeRegions() []int {
	// We use a flood fill algorithm to identify regions with the same biome
	biomeToRegs := initRegionSlice(m.mesh.numRegions)
	// Set all ocean regions to -2.
	for r := range biomeToRegs {
		if m.Elevation[r] <= 0.0 {
			biomeToRegs[r] = -2
		}
	}
	biomeFunc := m.getRegWhittakerModBiomeFunc()

	// Use a queue to implement the flood fill algorithm.
	outRegs := make([]int, 0, 6)
	floodFill := func(id int) {
		queue := list.New()
		// Get the biome that is represented by the region ID.
		biome := biomeFunc(id)
		queue.PushBack(id)

		// The region ID will serve as a representative of the biome.
		biomeToRegs[id] = id

		// Now flood fill all regions with the same biome.
		for queue.Len() > 0 {
			e := queue.Front()
			if e == nil {
				break
			}
			queue.Remove(e)
			nbID := e.Value.(int)

			for _, n := range m.mesh.r_circulate_r(outRegs, nbID) {
				if biomeToRegs[n] == -1 && biomeFunc(n) == biome {
					queue.PushBack(n)

					// The region ID will serve as a representative of the biome.
					biomeToRegs[n] = id
				}
			}
		}
	}

	// Loop through all regions and pick the first region that has not been
	// assigned a biome yet. Then flood fill all regions with the same biome.
	for id := 0; id < m.mesh.numRegions; id++ {
		if biomeToRegs[id] == -1 {
			floodFill(id)
		}
	}
	return biomeToRegs
}
