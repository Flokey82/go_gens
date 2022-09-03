package genworldvoronoi

import (
	"log"
	"math"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

// getTempFalloffFromAltitude returns the temperature falloff at a given altitude in meters
// above sea level. (approx. 9.8 °C per 1000 m)
// NOTE: This is definitely not correct :)
// Source: https://www.quora.com/At-what-rate-does-temperature-drop-with-altitude
func getTempFalloffFromAltitude(height float64) float64 {
	if height < 0 {
		return 0.0
	}
	const falloffFactor = 0.98 / 100
	return falloffFactor * height
}

const (
	minTemp          = -15.0
	maxTemp          = 30.0
	rangeTemp        = maxTemp - minTemp
	maxPrecipitation = 45.0 // 450dm
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

const maxAltitudeFactor = 8850 / 2 // How tall is the tallest mountain with an elevation of 1.0?

// getRTemperature returns the average yearly temperature of the given region at the surface.
func (m *Map) getRTemperature(r int, maxElev float64) float64 {
	// TODO: Fix maxElev caching!!!
	return getMeanAnnualTemp(m.r_latLon[r][0]) - getTempFalloffFromAltitude(maxAltitudeFactor*m.r_elevation[r]/maxElev)
}

// getRTemperature returns the average yearly temperature of the given triangle at the surface.
func (m *Map) getTTemperature(t int, maxElev float64) float64 {
	// TODO: Fix maxElev caching!!!
	return getMeanAnnualTemp(m.t_latLon[t][0]) - getTempFalloffFromAltitude(maxAltitudeFactor*m.t_elevation[t]/maxElev)
}

// getGlobalWindVector returns a vector for the global wind at the given latitude.
// NOTE: This is based on the trade winds on... well, earth.
// See: https://en.wikipedia.org/wiki/Trade_winds
func getGlobalWindVector(rLat float64) Vertex {
	// Based on latitude, we calculate the wind vector angle.
	var degree float64
	if rLatAbs := math.Abs(rLat); rLatAbs >= 0 && rLatAbs <= 30 {
		// +30° ... 0°, 0° ... -30° -> Primitive Hadley Cell.
		// In a Hadley cell, we turn the wind vector until we are exactly parallel with the equator once we reach 0° Lat.
		// TODO: This is probably not perfectly parallel at the equator.
		change := 90 * rLatAbs / 30
		if rLat > 0 {
			degree = 180 + change // Northern hemisphere.
		} else {
			degree = 180 - change // Southern hemisphere.
		}
	} else if rLatAbs > 30 && rLatAbs <= 60 {
		// +60° ... +30°, -30° ... -60° -> Primitive Mid Latitude Cell.
		// In a mid latitude cell, we turn the wind vector until we are exactly parallel with the 60° Lat.
		// TODO: This is probably not a full 90° turn. Fix this
		change := 90 * (rLatAbs - 30) / 30
		if rLat > 0 {
			degree = 90 - change // Northern hemisphere.
		} else {
			degree = 270 + change // Southern hemisphere.
		}
	} else {
		// NOTE: This is buggy or at least "not nice".
		// +90° ... +60°, -60° ... -90° -> Primitive Hadley Cell.
		// In a polar cell, we turn the wind vector until we are exactly parallel with the equator once we reach 60° Lat.
		change := 90 * (rLatAbs - 60) / 30
		if rLat > 0 {
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
func (m *Map) assignWindVectors() {
	r_windvec := make([]Vertex, m.mesh.numRegions)
	// Based on latitude of each region, we calculate the wind vector.
	for i := range r_windvec {
		r_windvec[i] = getGlobalWindVector(m.r_latLon[i][0])
	}

	r_windvec_local := make([]Vertex, m.mesh.numRegions)
	_, maxElev := minMax(m.r_elevation)

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
		var sea_r []int
		for r := 0; r < m.mesh.numRegions; r++ {
			if m.r_elevation[r] <= 0 {
				sea_r = append(sea_r, r)
			}
		}
		r_distance_sea := m.assignDistanceField(sea_r, make(map[int]bool))
		for r := range r_windvec_local {
			rVec := r_windvec[r]
			lat := m.r_latLon[r][0]
			lon := m.r_latLon[r][1]
			temp_r := getMeanAnnualTemp(lat) - getTempFalloffFromAltitude(8850*m.r_elevation[r]/maxElev)
			if m.r_elevation[r] < 0 {
				// TODO: Use actual distance from ocean to calculate temperature falloff.
				temp_r -= 1 / (r_distance_sea[r] + 1)
			}
			// Get temperature for r.
			v := vectors.Normalize(vectors.Vec2{
				X: rVec[0],
				Y: rVec[1],
			})
			for _, nb := range m.rNeighbors(r) {
				nbLat := m.r_latLon[nb][0]
				nbLon := m.r_latLon[nb][1]
				temp_nb := getMeanAnnualTemp(nbLat) - getTempFalloffFromAltitude(8850*m.r_elevation[nb]/maxElev)
				if m.r_elevation[nb] < 0 {
					// TODO: Use actual distance from ocean to calculate temperature falloff.
					temp_nb -= 1 / (r_distance_sea[nb] + 1)
				}
				ve := calcVecFromLatLong(lat, lon, nbLat, nbLon)
				v = v.Add(vectors.Normalize(vectors.NewVec2(ve[0], ve[1])).Mul(temp_nb - temp_r))
			}
			v = vectors.Normalize(v)
			r_windvec_local[r] = Vertex{v.X, v.Y}
		}
	} else {
		// Add wind deflection based on altitude changes.
		for r := range r_windvec_local {
			rVec := r_windvec[r]
			// Get XYZ Position of r.
			rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
			// Convert to polar coordinates.
			rLat := m.r_latLon[r][0]
			rLon := m.r_latLon[r][1]
			h := m.r_elevation[r]
			// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
			rwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[r][1], rLon+r_windvec[r][0])).Normalize()
			v := vectors.Normalize(vectors.Vec2{
				X: rVec[0],
				Y: rVec[1],
			}) //v.Mul(h / maxElev)
			for _, neighbor_r := range m.rNeighbors(r) {
				// if is_sea[neighbor_r] {
				//	continue
				// }
				// Calculate dot product of wind vector to vector r -> neighbor_r.
				// Get XYZ Position of r_neighbor.
				rnXYZ := convToVec3(m.r_xyz[neighbor_r*3 : neighbor_r*3+3])

				// Calculate Vector between r and neighbor_r.
				va := vectors.Sub3(rnXYZ, rXYZ).Normalize()

				// Calculate Vector between r and wind_r.
				vb := vectors.Sub3(rwXYZ, rXYZ).Normalize()

				// Calculate dot product between va and vb.
				// This will give us how much the current region lies within the wind direction of the
				// current neighbor.
				// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
				dotV := vectors.Dot3(va, vb)
				hnb := m.r_elevation[neighbor_r]
				if dotV > 0 && hnb != h && h >= 0 && hnb >= 0 {
					nbLat := m.r_latLon[neighbor_r][0]
					nbLon := m.r_latLon[neighbor_r][1]
					ve := calcVecFromLatLong(rLat, rLon, nbLat, nbLon)
					// The higher the dot product (the more direct the neighbor is in wind direction), the higher
					// the influence of an elevation change. So a steep mountain ahead will slow the wind down.
					// If a steep mountain is too the left, the wind vector will be pushed to the right.
					v = v.Add(vectors.Normalize(vectors.NewVec2(ve[0], ve[1])).Mul((h - hnb) * dotV / maxElev))
				}
			}
			v = vectors.Normalize(v)
			r_windvec_local[r] = Vertex{v.X, v.Y}
		}
	}
	// Average wind vectors using neighbor vectors.
	interpolationSteps := 0

	m.r_windvec_local = m.interpolateWindVecs(r_windvec_local, interpolationSteps)
	m.r_windvec = m.interpolateWindVecs(r_windvec, interpolationSteps)
}

// interpolateWindVecs interpolates the given wind vectors at their respective regions by
// mixing them with the wind vectors of their neighbor regions.
func (m *Map) interpolateWindVecs(in []Vertex, steps int) []Vertex {
	// Average wind vectors using neighbor vectors.
	for i := 0; i < steps; i++ {
		r_windvec_interpol := make([]Vertex, m.mesh.numRegions)
		for r := range r_windvec_interpol {
			resVec := Vertex{
				in[r][0],
				in[r][1],
			}
			var count int
			for _, neighbor_r := range m.rNeighbors(r) {
				resVec[0] += in[neighbor_r][0]
				resVec[1] += in[neighbor_r][1]
				count++
			}
			resVec[0] /= float64(count + 1)
			resVec[1] /= float64(count + 1)
			r_windvec_interpol[r] = resVec
		}
		in = r_windvec_interpol
	}
	return in
}

type biomesParams struct {
	raininess   float64 // 0, 2
	rain_shadow float64 // 0.1, 2
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
func (m *Map) assignRainfall(numSteps, transferMode, sortOrder int) {
	biomesParam := biomesParams{
		raininess:   0.9,
		rain_shadow: 0.9,
		evaporation: 0.9,
	}

	// 1. Initialize
	// 1.1. Determine all sea regions.
	var sea_r, land_r []int
	is_sea := make([]bool, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_elevation[r] < 0 {
			is_sea[r] = true
			sea_r = append(sea_r, r)
		} else {
			land_r = append(land_r, r)
		}
	}

	var sort_order []int
	if sortOrder == moistOrderOther {
		// 1.2. Sort all regions by distance to ocean. Lowest to highest.
		dist_order_r := make([]int, m.mesh.numRegions)
		for r := 0; r < m.mesh.numRegions; r++ {
			dist_order_r[r] = r
		}
		r_distance_sea := m.assignDistanceField(sea_r, make(map[int]bool))
		sort.Slice(dist_order_r, func(a, b int) bool {
			if r_distance_sea[dist_order_r[a]] == r_distance_sea[dist_order_r[b]] {
				return m.r_elevation[dist_order_r[a]] < m.r_elevation[dist_order_r[b]]
			}
			return r_distance_sea[dist_order_r[a]] < r_distance_sea[dist_order_r[b]]
		})
		sort_order = dist_order_r
	} else {
		// 1.2. Sort the indices in wind-order so we can ensure that we push the moisture
		// in their logical sequence across the globe.
		_, sort_order = m.getWindSortOrder() // Works reasonably well.
	}

	// 1.3. Assign wind vector for every region
	m.assignWindVectors()
	r_windvec := m.r_windvec
	_, maxH := minMax(m.r_elevation)

	calcRainfall := func(r int, humidity float64) float64 {
		r_elev := m.r_elevation[r]
		if r_elev < 0 {
			r_elev = 0 // Set to sea-level
		}
		heightVal := 1 - (r_elev / maxH)
		if humidity > heightVal {
			return biomesParam.rain_shadow * (humidity - heightVal)
		}
		return 0
	}

	for step := 0; step < numSteps; step++ {
		log.Println(step)
		// Evaporation.

		// 2. Assign initial moisture of 1.0 to all regions below or at sea level or replenish
		// moisture through evaporation if our moisture is below 0.
		for _, r := range sea_r {
			if m.r_moisture[r] < 1.0 {
				m.r_moisture[r] = 1.0
			}
			// m.r_rainfall[r] += biomesParam.raininess * m.r_moisture[r]
		}

		// Rivers should experience some evaporation.
		for r, fluxval := range m.r_flux {
			if m.r_moisture[r] < fluxval && m.r_moisture[r] < 1.0 {
				m.r_moisture[r] = 1.0 // TODO: Should depend on available water.
			}
		}

		// Water pools should experience some evaporation.
		for r, poolval := range m.r_pool {
			if poolval > 0 && m.r_moisture[r] < 1.0 {
				m.r_moisture[r] = 1.0 // TODO: Should depend on available water.
			}
		}
		// m.interpolateRainfallMoisture(1)

		// 3. Transfer moisture based on wind vectors.
		switch transferMode {
		case moistTransferDirect:
			// 3.1.B For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r.
			//       This will give us the amount of moisture we transfer to the neighbor region.
			// NOTE: This variant copies moisture from the current region to the neighbors that are in wind direction.
			for _, r := range sort_order {
				count := 0
				// Get XYZ Position of r.
				rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
				// Convert to polar coordinates.
				rLat := m.r_latLon[r][0]
				rLon := m.r_latLon[r][1]

				// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
				rwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[r][1], rLon+r_windvec[r][0])).Normalize()
				for _, neighbor_r := range m.rNeighbors(r) {
					if is_sea[neighbor_r] {
						continue
					}
					// Calculate dot product of wind vector to vector r -> neighbor_r.
					// Get XYZ Position of r_neighbor.
					rnXYZ := convToVec3(m.r_xyz[neighbor_r*3 : neighbor_r*3+3])

					// Calculate Vector between r and neighbor_r.
					va := vectors.Sub3(rnXYZ, rXYZ).Normalize()

					// Calculate Vector between r and wind_r.
					vb := vectors.Sub3(rwXYZ, rXYZ).Normalize()

					// Calculate dot product between va and vb.
					// This will give us how much the current region lies within the wind direction of the
					// current neighbor.
					// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
					dotV := vectors.Dot3(va, vb)
					if dotV > 0 {
						// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
						count++
						humidity := m.r_moisture[neighbor_r] + m.r_moisture[r]*dotV
						rainfall := m.r_rainfall[neighbor_r] // + biomesParam.raininess*m.r_moisture[r]*dotV
						orographicRainfall := calcRainfall(neighbor_r, humidity)
						if orographicRainfall > 0.0 {
							rainfall += biomesParam.raininess * orographicRainfall
							humidity -= orographicRainfall
						}
						// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
						// WARNING: The humidity calculation is off.
						// humidity = math.Min(humidity, 1.0)
						// rainfall = math.Min(rainfall, 1.0)
						m.r_rainfall[neighbor_r] = rainfall
						m.r_moisture[neighbor_r] = humidity
					}
				}
			}
		case moistTransferIndirect:
			// 3.2. For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r_neighbor.
			//    This will give us the amount of moisture we transfer from the neighbor region.
			// NOTE: This variant copies moisture to the current region from the neighbors depending on their wind direction.
			for _, r := range sort_order {
				count := 0
				sum := 0.0
				// Get XYZ Position of r.
				rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
				for _, neighbor_r := range m.rNeighbors(r) {
					// Calculate dot product of wind vector to vector r -> neighbor_r.
					// Get XYZ Position of r_neighbor.
					rnXYZ := convToVec3(m.r_xyz[neighbor_r*3 : neighbor_r*3+3])

					// Convert to polar coordinates.
					rLat := m.r_latLon[neighbor_r][0]
					rLon := m.r_latLon[neighbor_r][1]

					// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
					rnwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[neighbor_r][1], rLon+r_windvec[neighbor_r][0])).Normalize()

					// Calculate Vector between r and neighbor_r.
					va := vectors.Sub3(rXYZ, rnXYZ).Normalize()

					// Calculate Vector between neightbor_r and wind_neighbor_r.
					vb := vectors.Sub3(rnwXYZ, rnXYZ).Normalize()

					// Calculate dot product between va and vb.
					// This will give us how much the current region lies within the wind direction of the
					// current neighbor.
					// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals
					dotV := vectors.Dot3(va, vb)
					if dotV > 0 {
						// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
						count++
						sum += m.r_moisture[neighbor_r] * dotV
					}
				}

				var humidity, rainfall float64
				humidity = m.r_moisture[r]
				if count > 0 {
					// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
					// WARNING: The humidity calculation is off.
					humidity = math.Min(humidity+sum, 1.0) // / float64(count)
					rainfall = math.Min(rainfall+biomesParam.raininess*sum, 1.0)
				}
				if m.r_elevation[r] <= 0.0 {
					// evaporation := biomesParam.evaporation * (-m.r_elevation[r])
					// humidity = evaporation
					humidity = m.r_moisture[r]
				}
				orographicRainfall := calcRainfall(r, humidity)
				if orographicRainfall > 0.0 {
					rainfall += biomesParam.raininess * orographicRainfall
					humidity -= orographicRainfall
				}
				m.r_rainfall[r] = rainfall
				m.r_moisture[r] = humidity
			}
		}

		// 4. Average moisture and rainfall.
		// m.interpolateRainfallMoisture(1)
	}
}

func (m *Map) getWindSortOrder() ([]float64, []int) {
	m.assignWindVectors()
	// TODO: Add bool parameter to switch between local winds and global winds.
	//r_windvec := m.r_windvec
	r_windvec := m.r_windvec_local
	wind_order_r := make([]int, m.mesh.numRegions)
	r_wind_sort := make([]float64, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		wind_order_r[r] = r
		// TODO: modify the sort order by ensuring longitude wraps around...??
		lat := (m.r_latLon[r][0]) * r_windvec[r][1] / math.Abs(r_windvec[r][1]) // radToDeg(r_windvec[r][1])
		lon := (m.r_latLon[r][1]) * r_windvec[r][0] / math.Abs(r_windvec[r][0]) // radToDeg(r_windvec[r][0])
		r_wind_sort[r] = (lat + lon)
	}

	// Sort the indices in wind-order so we can ensure that we push the moisture
	// in their logical sequence across the globe.
	sort.Slice(wind_order_r, func(a, b int) bool {
		return r_wind_sort[wind_order_r[a]] < r_wind_sort[wind_order_r[b]]
	})
	return r_wind_sort, wind_order_r
}

func (m *Map) assignRainfallBasic() {
	// NOTE: This still has issues with the wrap around at +/- 180° long
	biomesParam := biomesParams{
		raininess:   0.9,
		rain_shadow: 0.9,
		evaporation: 0.9,
	}
	humidityFromRiver := 0.25
	humidityFromSea := 1.0
	evaporateRivers := false

	// Sort the indices in wind-order so we can ensure that we push the moisture
	// in their logical sequence across the globe.
	_, wind_order_r := m.getWindSortOrder()
	r_windvec := m.r_windvec_local

	// calcRainfall returns the amount of rain shed given the region and humidity.
	_, maxH := minMax(m.r_elevation)
	calcRainfall := func(r int, humidity float64) float64 {
		r_elev := m.r_elevation[r]
		if r_elev < 0 {
			r_elev = 0 // Set to sea-level
		}
		heightVal := 1 - (r_elev / maxH)
		if humidity > heightVal {
			return biomesParam.rain_shadow * (humidity - heightVal)
		}
		return 0
	}

	// Evaporation.
	// 1. Assign initial moisture of 1.0 to all regions below or at sea level or replenish
	// moisture through evaporation if our moisture is below 0.
	for r, h := range m.r_elevation {
		if h < 0 && m.r_moisture[r] < humidityFromSea {
			m.r_moisture[r] = humidityFromSea
		}
	}

	// Rivers should experience some evaporation.
	if evaporateRivers {
		for r, fluxval := range m.r_flux {
			if m.r_moisture[r] < fluxval && m.r_moisture[r] < humidityFromRiver {
				m.r_moisture[r] = humidityFromRiver // TODO: Should depend on available water.
			}
		}
	}

	// Water pools should experience some evaporation.
	for r, poolval := range m.r_pool {
		if m.r_moisture[r] < humidityFromSea && poolval > 0 {
			m.r_moisture[r] = humidityFromSea // TODO: Should depend on available water.
		}
	}

	// Visit regions in wind order and copy the moisture from the neighbor regious that are
	// up-wind.
	for i := 0; i < 4; i++ {
		for _, r := range wind_order_r {
			var humidity float64

			// Calculate humidity.
			for _, neighbor_r := range m.rNeighbors(r) {
				rL := m.r_latLon[r]
				nL := m.r_latLon[neighbor_r]

				// TODO: Check dot product of wind vector (r) and neighbour->r.
				vVec := normal2(calcVecFromLatLong(nL[0], nL[1], nL[0]+r_windvec[neighbor_r][1], nL[1]+r_windvec[neighbor_r][0]))

				//nVec := normal2([2]float64{rL[1] - nL[1], rL[0] - nL[0]})
				nVec := normal2(calcVecFromLatLong(nL[0], nL[1], rL[0], rL[1]))
				dotV := dot2(vVec, nVec)

				// Check if the neighbor region is up-wind (that the wind blows from neighbor_r to r).
				// We use both the wind-order using modulo (wrap around 360° longitude) and the normal order to
				// catch cases like a neighbor being at a negative order, or just at 360° longitude, blowing across to a neighbor
				// sitting at 0° etc.
				// NOTE: This needs to be thoroughly checked if this fix works for the border regions.
				if dotV > 0.0 {
					humidity += m.r_moisture[neighbor_r] * dotV
				}
			}

			// Set base rainfall.
			rainfall := biomesParam.raininess * humidity

			// Evaporation.
			// TODO: Remove the evaporation here if possible and instead
			// rely on the evaporation step above.
			if m.r_elevation[r] < 0 {
				evaporation := biomesParam.evaporation * (-m.r_elevation[r])
				humidity += evaporation
			} else if evaporateRivers {
				evaporation := biomesParam.evaporation * m.r_flux[r]
				humidity += evaporation
			}

			// Calculate orographic rainfall caused by elevation changes.
			orographicRainfall := calcRainfall(r, humidity)
			if orographicRainfall > 0.0 {
				rainfall += biomesParam.raininess * orographicRainfall
				humidity -= orographicRainfall
			}
			m.r_rainfall[r] = rainfall
			m.r_moisture[r] = humidity
		}
	}
	m.interpolateRainfallMoisture(15)
}

func (m *Map) interpolateRainfallMoisture(interpolationSteps int) {
	for i := 0; i < interpolationSteps; i++ {
		r_moisture_interpol := make([]float64, m.mesh.numRegions)
		r_rainfall_interpol := make([]float64, m.mesh.numRegions)
		for r := range r_moisture_interpol {
			rMoist := m.r_moisture[r]
			rRain := m.r_rainfall[r]
			var count int
			for _, neighbor_r := range m.rNeighbors(r) {
				// Gravity! Water moves downwards.
				// This is not super-accurate since you'd have to take
				// in account how steep the slope is etc.
				if m.r_elevation[r] >= m.r_elevation[neighbor_r] {
					continue
				}
				rMoist += m.r_moisture[neighbor_r]
				rRain += m.r_rainfall[neighbor_r]
				count++
			}
			r_moisture_interpol[r] = rMoist / float64(count+1)
			r_rainfall_interpol[r] = rRain / float64(count+1)
		}
		m.r_moisture = r_moisture_interpol
		m.r_rainfall = r_rainfall_interpol
	}
}
