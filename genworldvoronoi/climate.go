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
		height = 0
	}
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
	r_distance_sea := m.assignDistanceField(sea_r, make(map[int]bool))
	for r := 0; r < m.mesh.numRegions; r++ {
		m.r_moisture[r] = 1 - math.Min(float64(r_distance_sea[r]), 100)/float64(100)
	}
}

// getRTemperature returns the average yearly temperature of the given region at the surface.
func (m *Map) getRTemperature(r int, maxElev float64) float64 {
	// TODO: Fix maxElev caching!!!
	return getMeanAnnualTemp(m.r_latLon[r][0]) - getTempFalloffFromAltitude(8850*m.r_elevation[r]/maxElev)
}

// getRTemperature returns the average yearly temperature of the given triangle at the surface.
func (m *Map) getTTemperature(t int) float64 {
	_, maxH := minMax(m.t_elevation) // TODO: Cache somewhere?
	return getMeanAnnualTemp(m.t_latLon[t][0]) - getTempFalloffFromAltitude(8850*m.t_elevation[t]/maxH)
}

// assignWindVectors constructs faux global wind cells reminiscent of a simplified earth model.
// NOTE: This function includes an experimental part that calculates local winds that are influenced
// by the topography / elevation changes. Please note that the code for local winds is incomplete.
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
			change := 45 * (rLatAbs - 30) / 30
			if rLat > 0 {
				degree = 90 - change // Northern hemisphere.
			} else {
				degree = 270 + change // Southern hemisphere.
			}
		} else {
			// NOTE: This is buggy or at least "not nice".
			// +90° ... +60°, -60° ... -90° -> Primitive Hadley Cell.
			// In a polar cell, we turn the wind vector until we are exactly parallel with the equator once we reach 60° Lat.
			change := 90 * (30 - (rLatAbs - 60)) / 30
			if rLat > 0 {
				degree = 180 + change // Northern hemisphere.
			} else {
				degree = 180 - change // Southern hemisphere.
			}
		}
		rad := degToRad(degree)
		r_windvec[i] = Vertex{math.Cos(rad), math.Sin(rad)}
		// log.Println(rLat, degree, r_windvec[i])
	}

	// TODO: Add wind vectors based on local temperature gradients.
	// TODO: Add wind deflection based on altitude changes.
	r_windvec_local := make([]Vertex, m.mesh.numRegions)
	_, maxElev := minMax(m.r_elevation)

	// Determine all sea regions.
	var sea_r []int
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_elevation[r] <= 0 {
			sea_r = append(sea_r, r)
		}
	}
	r_distance_sea := m.assignDistanceField(sea_r, make(map[int]bool))
	for r := range r_windvec_local {
		lat := m.r_latLon[r][0]
		lon := m.r_latLon[r][1]
		temp_r := getMeanAnnualTemp(lat) - getTempFalloffFromAltitude(8850*m.r_elevation[r]/maxElev)
		// if m.r_elevation[r] < 0 {
		// TODO: Use actual distance from ocean to calculate temperature falloff.
		temp_r -= 1 / (r_distance_sea[r] + 1)
		// }
		// Get temperature for r.
		var v vectors.Vec2
		for _, nb := range m.rNeighbors(r) {
			nbLat := m.r_latLon[nb][0]
			nbLon := m.r_latLon[nb][1]
			temp_nb := getMeanAnnualTemp(nbLat) - getTempFalloffFromAltitude(8850*m.r_elevation[nb]/maxElev)
			// if m.r_elevation[nb] < 0 {
			// TODO: Use actual distance from ocean to calculate temperature falloff.
			temp_nb -= 1 / (r_distance_sea[nb] + 1)
			// }
			ve := calcVecFromLatLong(lat, lon, nbLat, nbLon)
			v = v.Add(vectors.Normalize(vectors.NewVec2(ve[0], ve[1])).Mul(temp_nb - temp_r))
		}
		v = vectors.Normalize(v)
		r_windvec_local[r] = Vertex{v.X, v.Y}
	}
	for r := range r_windvec_local {
		rVec := r_windvec[r]
		// Get XYZ Position of r.
		rXYZ := convToVec3(m.r_xyz[r*3 : r*3+3])
		// Convert to polar coordinates.
		rLat := m.r_latLon[r][0]
		rLon := m.r_latLon[r][1]
		h := m.r_elevation[r]
		// Add wind vector to neighbor lat/lon to get the "wind vector lat long" or something like that..
		rwXYZ := convToVec3(latLonToCartesian(rLat+r_windvec[r][0], rLon+r_windvec[r][1])).Normalize()
		var v vectors.Vec2
		v.X = rVec[0]
		v.Y = rVec[1]
		v = vectors.Normalize(v) //v.Mul(h / maxElev)
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
				v = v.Add(vectors.Normalize(vectors.NewVec2(ve[0], ve[1])).Mul(((hnb - h) * dotV / maxElev)))
			}
		}
		v = vectors.Normalize(v)
		r_windvec_local[r] = Vertex{v.X, v.Y}
	}
	// Average wind vectors using neighbor vectors.
	interpolationSteps := 5

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
	moistTransferWalk     = 2
	moistTransferDirect2  = 3
	moistTransferWalk2    = 4
)

// assignRainfall is an overengineered logic that is supposed to calculate the transfer
// of moisture across the globe based on global winds using distinct approaches.
// Unfortunately, this is highly bugged and not as useful as the simpler version
// 'assignRainfallBasic'.
func (m *Map) assignRainfall(numSteps, transferMode int) {
	biomesParam := biomesParams{
		raininess:   0.9,
		rain_shadow: 0.9,
		evaporation: 0.9,
	}

	// 1. Initialize
	// 1.1. Determine all sea regions.
	var sea_r, land_r []int
	is_sea := make(map[int]bool)
	is_land := make(map[int]bool)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.r_elevation[r] <= 0 {
			is_sea[r] = true
			sea_r = append(sea_r, r)
		} else {
			is_land[r] = true
			land_r = append(sea_r, r)
		}
	}

	// 1.2. Sort all regions by distance to ocean. Lowest to highest.
	dist_order_r := make([]int, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		dist_order_r[r] = r
	}
	r_distance_sea := m.assignDistanceField(sea_r, is_sea)
	r_distance_land := m.assignDistanceField(land_r, is_land)
	sort.Slice(dist_order_r, func(a, b int) bool {
		if r_distance_sea[dist_order_r[a]] == r_distance_sea[dist_order_r[b]] {
			return m.r_elevation[dist_order_r[a]] < m.r_elevation[dist_order_r[b]]
		}
		return r_distance_sea[dist_order_r[a]] < r_distance_sea[dist_order_r[b]]
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
		case moistTransferDirect2:
			sort.Slice(dist_order_r, func(a, b int) bool {
				if m.r_moisture[dist_order_r[a]] == m.r_moisture[dist_order_r[b]] {
					return r_distance_sea[dist_order_r[a]] < r_distance_sea[dist_order_r[b]]
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
				if seen[r] || transfer < 0.001 {
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
		case moistTransferWalk2:
			seen := make([]bool, m.mesh.numRegions)
			tmp := make([]float64, m.mesh.numRegions)
			var drill func(r int, transfer float64)
			drill = func(r int, transfer float64) {
				//log.Println(transfer)
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
					if dotV > 0 && !seen[neighbor_r] && m.r_elevation[neighbor_r] >= 0 {
						// Only positive dot products mean that we lie within 90°, so 'in wind direction'.
						tmp[neighbor_r] += transfer * dotV
						// nLat := m.r_latLon[neighbor_r][0]
						// nLon := m.r_latLon[neighbor_r][1]
						// dist := haversine(rLat, rLon, nLat, nLon) / (2 * math.Pi)
						drill(neighbor_r, transfer*dotV)
					}
				}
			}
			for _, r := range sea_r {
				if r_distance_land[r] == 0 {
					drill(r, 1)
					for i := range seen {
						seen[i] = false
					}
				}
			}
			for r, h := range tmp {
				m.r_moisture[r] += h
				tmp[r] = 0
			}
			/*
				for _, r := range dist_order_r {
					if r_distance_sea[r] >= 0 && m.r_moisture[r] > 0 {
						drill(r, m.r_moisture[r])
						for i := range seen {
							seen[i] = false
						}
					}
				}
				for r, h := range tmp {
					m.r_moisture[r] += h
					tmp[r] = 0
				}*/
			/*
				for _, r := range dist_order_r {
					if r_distance_sea[r] >= 0 && m.r_moisture[r] > 0 {
						log.Println(r_distance_sea[r])
						drill(r, m.r_moisture[r])
						for i := range seen {
							seen[i] = false
						}
					}
				}
				for r, h := range tmp {
					m.r_moisture[r] += h
					humidity := m.r_moisture[r]
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
				}*/
			for _, r := range dist_order_r {
				humidity := m.r_moisture[r]
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
		}

		// 4. Average moisture and rainfall.
		// m.interpolateRainfallMoisture(1)
	}
}

// assignRainfallBasic visits all regions in down-wind order and copies moisture from its up-wind
// neighbors.
func (m *Map) assignRainfallBasic() {
	biomesParam := biomesParams{
		raininess:   0.9,
		rain_shadow: 0.8,
		evaporation: 0.9,
	}
	humidityFromRiver := 0.25
	humidityFromSea := 1.0
	evaporateRivers := false

	m.assignWindVectors()
	r_windvec := m.r_windvec
	// r_windvec := m.r_windvec_local
	wind_order_r := make([]int, m.mesh.numRegions)
	r_wind_sort := make([]float64, m.mesh.numRegions)
	r_wind_sort_mod := make([]float64, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		wind_order_r[r] = r
		// Here I modified the sort order by ensuring longitude wraps around... Verify??
		lon := (m.r_latLon[r][1] + 180) * r_windvec[r][1]
		if lon < 0 {
			lon += 360 // make sure negative numbers wrap around to 360.0
		}
		r_wind_sort[r] = math.Abs(m.r_latLon[r][0]*r_windvec[r][0]) + lon

		// Make sure we wrap around 360° correctly (with mod to preven values > 360.0).
		lonMod := math.Mod((m.r_latLon[r][1] + 180*r_windvec[r][1]), 360.0)
		if lonMod < 0 {
			lonMod += 360
		}
		r_wind_sort_mod[r] = math.Abs(m.r_latLon[r][0])*r_windvec[r][0] + lonMod
	}

	// Sort the indices in wind-order so we can ensure that we push the moisture
	// in their logical sequence across the globe.
	sort.Slice(wind_order_r, func(a, b int) bool {
		return r_wind_sort[wind_order_r[a]] > r_wind_sort[wind_order_r[b]]
	})

	_, maxH := minMax(m.r_elevation)

	// calcRainfall returns the amount of rain shed given the region and humidity.
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
	for i := 0; i < 1; i++ {
		// 1. Assign initial moisture of 1.0 to all regions below or at sea level or replenish
		// moisture through evaporation if our moisture is below 0.
		for r, h := range m.r_elevation {
			if m.r_moisture[r] < humidityFromSea && h < 0 {
				m.r_moisture[r] = humidityFromSea
			}
			// m.r_rainfall[r] += biomesParam.raininess * m.r_moisture[r]
		}

		// Rivers should experience some evaporation.
		// if evaporateRivers {
		//  for r, fluxval := range m.r_flux {
		// 	 if m.r_moisture[r] < fluxval && m.r_moisture[r] < humidityFromRiver {
		// 	 	 m.r_moisture[r] = humidityFromRiver // TODO: Should depend on available water.
		//	 }
		//  }
		// }

		// Water pools should experience some evaporation.
		for r, poolval := range m.r_pool {
			if m.r_moisture[r] < humidityFromSea && poolval > 0 {
				m.r_moisture[r] = humidityFromSea // TODO: Should depend on available water.
			}
		}

		// Visit regions in wind order and copy the moisture from the neighbor regious that are
		// up-wind.
		for _, r := range wind_order_r {
			var count int
			var sum float64
			// TODO: Here I modified the sort order by ensuring longitude wraps around... Verify??
			for _, neighbor_r := range m.rNeighbors(r) {
				// Check if the neighbor region is up-wind (that the wind blows from neighbor_r to r).
				// We use both the wind-order using modulo (wrap around 360° longitude) and the normal order to
				// catch cases like a neighbor being at a negative order, or just at 360° longitude, blowing across to a neighbor
				// sitting at 0° etc.
				// NOTE: This needs to be thoroughly checkt if this fix works for the border regions.
				if r_wind_sort_mod[neighbor_r] > r_wind_sort_mod[r] || r_wind_sort[neighbor_r] > r_wind_sort[r] {
					count++
					sum += m.r_moisture[neighbor_r]
				}
			}

			var humidity, rainfall float64
			if count > 0 {
				humidity = sum / float64(count)
				rainfall = biomesParam.raininess * humidity
			} else {
				humidity = m.r_moisture[r]
			}
			if m.r_elevation[r] < 0 && humidity < humidityFromSea {
				// evaporation := biomesParam.evaporation * (-m.r_elevation[r])
				// humidity = 1 + evaporation
				humidity = humidityFromSea
			} else if evaporateRivers && m.r_rainfall[r] < m.r_flux[r] && humidity < humidityFromRiver {
				// evaporation := biomesParam.evaporation * m.r_flux[r]
				// humidity += evaporation
				humidity = humidityFromRiver
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

func calcVecFromLatLong(lat1, lon1, lat2, lon2 float64) [2]float64 {
	lat1 *= math.Pi / 180
	lon1 *= math.Pi / 180
	lat2 *= math.Pi / 180
	lon2 *= math.Pi / 180
	return [2]float64{
		math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1), // X
		math.Sin(lon2-lon1) * math.Cos(lat2),                                              // Y
	}
}
