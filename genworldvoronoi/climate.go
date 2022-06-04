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
	return 0.98 * height / 100
}

// getMeanAnnualTemp returns the temperature at a given latitude within the range of
// -35 °C to +31 °C.
// For this I assume that light hits the globe exactly from a 90° angle with respect
// to the planitary axis.
// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/shading-normals (facing ratio)
// See: http://www-das.uwyo.edu/~geerts/cwx/notes/chap16/geo_clim.html
func getMeanAnnualTemp(lat float64) float64 {
	return (math.Sin(degToRad(90-math.Abs(lat))))*66 - 35
}

// assignRegionMoisture assigns moisture to each region based on the plate
// it is on. Ideally we would calculate this using prevailing winds and distance
// from the ocean and whatnot.
// NOTE: This is currently unused.
func (m *Map) assignRegionMoisture() {
	// TODO: assign region moisture in a better way!
	// for r := 0; r < m.mesh.numRegions; r++ {
	//	m.r_moisture[r] = float64(m.r_plate[r]%10) / 10.0
	// }
	var sea_r []int
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_elevation[r] < 0 {
			sea_r = append(sea_r, r)
		}
	}

	// Assign basic moisture per region based on distance to sea.
	// NOTE: This is currently overridden in assignRainfall().
	r_distance_d := m.assignDistanceField(sea_r, make(map[int]bool))
	for r := 0; r < m.mesh.numRegions; r++ {
		m.r_moisture[r] = 1 - math.Min(float64(r_distance_d[r]), 100)/float64(100)
	}
}

func (m *Map) assignWindVectors() {
	r_windvec := make([]Vertex, m.mesh.numRegions)
	for i := range r_windvec {
		// Determine latitude of current region.
		rLat := m.r_latLon[i][0]
		// Based on latitude, we calculate the wind vector angle.
		var degree float64
		if rLatAbs := math.Abs(rLat); rLatAbs >= 0 && rLatAbs <= 30 {
			// +30° ... 0°, 0° ... -30° -> Primitive Hadley Cell.
			// In a Hadley cell, we turn the wind vector until we are exactly parallel with the equator once we reach 0° Lat.
			// TODO: This is probably not perfectly parallel at the equator.
			if rLat > 0 {
				degree = 180 + 90*(30-rLatAbs)/30 // Northern hemisphere.
			} else {
				degree = 360 - 90*(30-rLatAbs)/30 // Southern hemisphere.
			}
		} else if rLatAbs > 30 && rLatAbs <= 60 {
			// +60° ... +30°, -30° ... -60° -> Primitive Mid Latitude Cell.
			// In a mid latitude cell, we turn the wind vector until we are exactly parallel with the 60° Lat.
			// TODO: This is probably not a full 90° turn. Fix this
			if rLat > 0 {
				degree = 0 + 90*(rLatAbs-30)/30 // Northern hemisphere.
			} else {
				degree = 180 - 90*(rLatAbs-30)/30 // Southern hemisphere.
			}
		} else {
			// NOTE: This is buggy or at least "not nice".
			// +90° ... +60°, -60° ... -90° -> Primitive Hadley Cell.
			// In a polar cell, we turn the wind vector until we are exactly parallel with the equator once we reach 60° Lat.
			if rLat > 0 {
				degree = 180 + 90*(90-rLatAbs)/30 // Northern hemisphere.
			} else {
				degree = 360 - 90*(90-rLatAbs)/30 // Southern hemisphere.
			}
		}
		rad := degToRad(degree)
		r_windvec[i] = Vertex{math.Cos(rad), math.Sin(rad)}
	}

	// TODO: Add wind vectors based on local temperature gradients.
	/*
		r_windvec_local := make([]Vertex, m.mesh.numRegions)
		_, maxElev := minMax(m.r_elevation)
		for r := range r_windvec_local {
			s0 := m.mesh.RInS[r]
			incoming := s0
			lat, _ := latLonFromVec3(convToVec3(m.r_xyz[r*3:(r*3)+3]).Normalize(), 1.0)
			// Get temperature for r.
			temp_r := getMeanAnnualTemp(lat) - getTempFalloffFromAltitude(8850*m.r_elevation[r]/maxElev)
			var count int
			for {
				neighbor_r := m.mesh.s_begin_r(incoming)
				// Get temperature for neighbor_r.
				// if temp_neighbor_r > temp_r
				// 	get direction of gradient.
				// sum up using diff to determine average vector
				count++
				outgoing := s_next_s(incoming)
				incoming = m.mesh.Halfedges[outgoing]
				if incoming == s0 {
					break
				}
			}
		}
	*/

	// Average wind vectors using neighbor vectors.
	interpolationSteps := 5
	for i := 0; i < interpolationSteps; i++ {
		r_windvec_interpol := make([]Vertex, m.mesh.numRegions)
		for r := range r_windvec_interpol {
			resVec := Vertex{
				r_windvec[r][0],
				r_windvec[r][1],
			}
			var count int
			for _, neighbor_r := range m.rNeighbors(r) {
				resVec[0] += r_windvec[neighbor_r][0]
				resVec[1] += r_windvec[neighbor_r][1]
				count++
			}
			resVec[0] /= float64(count + 1)
			resVec[1] /= float64(count + 1)
			r_windvec_interpol[r] = resVec
		}
		r_windvec = r_windvec_interpol
	}

	m.r_windvec = r_windvec
}

type biomesParams struct {
	raininess   float64 // 0, 2
	rain_shadow float64 // 0.1, 2
	evaporation float64 // 0, 1
}

const (
	moistTransferDirect   = 0
	moistTransferIndirect = 1
	moistTransferWalk     = 2
	moistTransferDirect2  = 3
)

func (m *Map) assignRainfall(numSteps int) {
	biomesParam := biomesParams{
		raininess:   0.9,
		rain_shadow: 0.9,
		evaporation: 0.9,
	}

	// 1. Initialize
	// 1.1. Determine all sea regions.
	var sea_r []int
	is_sea := make(map[int]bool)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_elevation[r] <= 0 {
			is_sea[r] = true
			sea_r = append(sea_r, r)
		}
	}

	// 1.2. Sort all regions by distance to ocean. Lowest to highest.
	dist_order_r := make([]int, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		dist_order_r[r] = r
	}
	r_distance_d := m.assignDistanceField(sea_r, make(map[int]bool))
	sort.Slice(dist_order_r, func(a, b int) bool {
		if r_distance_d[dist_order_r[a]] == r_distance_d[dist_order_r[b]] {
			return m.r_elevation[dist_order_r[a]] < m.r_elevation[dist_order_r[b]]
		}
		return r_distance_d[dist_order_r[a]] < r_distance_d[dist_order_r[b]]
	})

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

	transferMode := moistTransferDirect2
	for step := 0; step < numSteps; step++ {
		log.Println(step)
		// Evaporation.

		// 2. Assign initial moisture of 1.0 to all regions below or at sea level or replenish
		// moisture through evaporation if our moisture is below 0.
		for _, r := range sea_r {
			if m.r_moisture[r] < 1.0 {
				m.r_moisture[r] = 1.0
			}
			//m.r_rainfall[r] += biomesParam.raininess * m.r_moisture[r]
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
		case moistTransferDirect2:
			sort.Slice(dist_order_r, func(a, b int) bool {
				if m.r_moisture[dist_order_r[a]] == m.r_moisture[dist_order_r[b]] {
					return r_distance_d[dist_order_r[a]] < r_distance_d[dist_order_r[b]]
				}
				return m.r_moisture[dist_order_r[a]] > m.r_moisture[dist_order_r[b]]
			})
			r_moisture_interpol := make([]float64, m.mesh.numRegions)
			// 3.1.A For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r.
			//       This will give us the amount of moisture we transfer to the neighbor region.
			// NOTE: This variant copies moisture from the current region to the neighbors that are in wind direction.
			//       Additionally we use a temporary array to store the transferred moisture and only merge the moisture
			//       once all regions are copied.
			for _, r := range dist_order_r {
				// Get XYZ Position of r.
				rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
				// Convert to polar coordinates.
				rLat := m.r_latLon[r][0]
				rLon := m.r_latLon[r][1]

				// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
				rwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[r][0], rLon+r_windvec[r][1])).Normalize()
				for _, neighbor_r := range m.rNeighbors(r) {
					//if is_sea[neighbor_r] {
					//	continue
					//}
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
						moistOut := m.r_moisture[r] * dotV
						r_moisture_interpol[neighbor_r] += moistOut
					}
				}
			}
			// Now merge the additional moisture with the base moisture and let it rain accordingly.
			for _, r := range dist_order_r {
				// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
				humidity := m.r_moisture[r] + r_moisture_interpol[r]
				rainfall := m.r_rainfall[r] // + biomesParam.raininess*m.r_moisture[r]*dotV
				orographicRainfall := calcRainfall(r, humidity)
				if orographicRainfall > 0.0 {
					rainfall = biomesParam.raininess * orographicRainfall
					humidity -= orographicRainfall
				}
				// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
				// WARNING: The humidity calculation is off.
				// humidity = math.Min(humidity, 1.0)
				// rainfall = math.Min(rainfall, 1.0)
				m.r_rainfall[r] = rainfall
				m.r_moisture[r] = humidity
			}
		case moistTransferDirect:
			// 3.1.B For each region, calculate dot product of Vec r -> r_neighbor and wind vector of r.
			//       This will give us the amount of moisture we transfer to the neighbor region.
			// NOTE: This variant copies moisture from the current region to the neighbors that are in wind direction.
			for _, r := range dist_order_r {
				count := 0
				// Get XYZ Position of r.
				rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
				// Convert to polar coordinates.
				rLat := m.r_latLon[r][0]
				rLon := m.r_latLon[r][1]

				// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
				rwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[r][0], rLon+r_windvec[r][1])).Normalize()
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
			for _, r := range dist_order_r {
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
					rnwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[neighbor_r][0], rLon+r_windvec[neighbor_r][1])).Normalize()

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
				// if m.mesh.r_boundary(r) {
				//	 humidity = 1.0
				// }
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
		case moistTransferWalk:
			seen := make([]bool, m.mesh.numRegions)
			var drill func(r int, transfer float64)
			drill = func(r int, transfer float64) {
				log.Println(transfer)
				if seen[r] || transfer < 0.00001 {
					return
				}
				seen[r] = true
				// Get XYZ Position of r.
				rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
				// Convert to polar coordinates.
				rLat := m.r_latLon[r][0]
				rLon := m.r_latLon[r][1]

				// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
				rwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[r][0], rLon+r_windvec[r][1])).Normalize()
				for _, neighbor_r := range m.rNeighbors(r) {
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
					if dotV > 0 && !seen[neighbor_r] && m.r_elevation[neighbor_r] > 0 { //
						// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
						m.r_moisture[neighbor_r] += transfer * dotV
						humidity := m.r_moisture[neighbor_r]
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
						drill(neighbor_r, transfer*dotV)
					}
				}
			}
			for _, r := range dist_order_r {
				if dist_order_r[r] < 1 || m.r_moisture[r] > 0 {
					drill(r, m.r_moisture[r])
					for i := range seen {
						seen[i] = false
					}
				}
				/*
					humidity := m.r_moisture[r]
					rainfall := m.r_rainfall[r] // + biomesParam.raininess*m.r_moisture[r]*dotV
					orographicRainfall := calcRainfall(r, humidity)
					if orographicRainfall > 0.0 {
						rainfall += biomesParam.raininess * orographicRainfall
						humidity -= orographicRainfall
					}
					// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
					// WARNING: The humidity calculation is off.
					// humidity = math.Min(humidity, 1.0)
					// rainfall = math.Min(rainfall, 1.0)
					m.r_rainfall[r] = rainfall
					m.r_moisture[r] = humidity*/
			}
			log.Println(m.r_rainfall)
			/*
				for _, r := range dist_order_r {
					humidity := m.r_moisture[r]
					rainfall := m.r_rainfall[r] // + biomesParam.raininess*m.r_moisture[r]*dotV
					orographicRainfall := calcRainfall(r, humidity)
					if orographicRainfall > 0.0 {
						rainfall += biomesParam.raininess * orographicRainfall
						humidity -= orographicRainfall
					}
					// TODO: Calculate max humidity at current altitude, temperature, rain off the rest.
					// WARNING: The humidity calculation is off.
					// humidity = math.Min(humidity, 1.0)
					// rainfall = math.Min(rainfall, 1.0)
					m.r_rainfall[r] = rainfall
					m.r_moisture[r] = humidity
				}*/
		}

		// 4. Average moisture and rainfall.
		//m.interpolateRainfallMoisture(1)
	}
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
