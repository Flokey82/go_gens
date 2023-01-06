package genworldvoronoi

import (
	"math"
)

// Erode erodes all region by the given amount and returns the resulting heightmap.
//
// NOTE: This is based on mewo2's erosion code but limits the eroded height
// to a fraction of the height difference to the downhill neighbor, which
// prevents extreme sinks from forming.
//
// See: https://github.com/mewo2/terrain
func (m *Geo) Erode(amount float64) []float64 {
	// Select the erosion method.
	useAlternative := true

	// Get downhill height diffs so we can ensure that we do not erode
	// any more than that, which would produce sinks (which we try to avoid).
	dhDiff := make([]float64, m.mesh.numRegions)
	for r, dhr := range m.GetDownhill(false) {
		// Skip all sinks which have a downhill value of -1
		if dhr < 0 {
			// NOTE: Sinks have no height diff, but in theory we could give it a
			// negative height diff to fill the sinks during the erosion steps?
			continue
		}
		dhDiff[r] = m.Elevation[r] - m.Elevation[dhr]
	}

	// This will hold our new heightmap.
	newh := make([]float64, m.mesh.numRegions)

	// Get the erosion rate for all regions.
	var er []float64
	if useAlternative {
		er = m.GetErosionRate2()
	} else {
		er = m.GetErosionRate()
	}

	// Get the maximum erosion rate, so we can normalize the erosion values.
	_, maxr := minMax(er)

	// Calculate the new heightmap by applying the erosion rates we have calculated.
	for r, e := range er {
		// We can at most erode amount*dhDiff[r].
		newh[r] = m.Elevation[r] - amount*dhDiff[r]*(e/maxr)
	}
	return newh
}

// GetErosionRate returns the erosion rate per region.
// NOTE: This is based on mewo2's erosion code
// See: https://github.com/mewo2/terrain
func (m *Geo) GetErosionRate() []float64 {
	// TODO: Change to distance so that it is independent of resolution.
	const nbErosionFactor = 0.125

	erodeNeighbors := true         // Traverse neighbors to erode them
	erodeNeighborsRiver := 1       // How deep should we traverse the neighbor-graph? (rivers)
	erodeNeighborsGlacier := 1     // How deep should we traverse the neighbor-graph? (glaciers)
	erodeOnlyAboveSealevel := true // Should we skip erosion below sea level?

	// Get the flux values for all regions.
	flux := m.getFlux(erodeOnlyAboveSealevel)

	// Get max flux so we can normalize the flux values to 0.0 ... 1.0.
	_, maxFlux := minMax(flux)
	if maxFlux == 0 {
		maxFlux = 1
	}

	// Get the slope values for all regions.
	slope := m.GetSlope()

	// This will hold the erosion values for each region.
	newh := make([]float64, m.mesh.numRegions)

	// Get the max height value so we can normalize the elevation values.
	_, maxH := minMax(m.Elevation) // TODO: Cache somewhere?

	// erodeRegion sets the erosion rate for the given region and
	// traverses the neighbor graph up to the remaining depth (rem).
	var erodeRegion func(r, rem int, toErode float64)

	erodeRegion = func(r, rem int, toErode float64) {
		// If we have erosion below sea level, skip this region.
		if erodeOnlyAboveSealevel && m.Elevation[r] < 0 {
			return
		}

		// If the given erosion rate is higher than the currently
		// assigned one, override the value.
		//
		// TODO: If we have a higher erosion rate in newh[r], should we use
		// newh[r] instead of toErode to continue?
		if toErode > newh[r] {
			newh[r] = toErode
		}

		// Check if we have exhausted the max traversal depth.
		rem--
		if rem < 0 || !erodeNeighbors {
			return
		}

		// Additionally visit all neighbors and erode them by a certain
		// fraction of the given erosion value.
		toErode *= nbErosionFactor
		for _, nb := range m.GetRegNeighbors(r) {
			erodeRegion(nb, rem, toErode)
		}
	}

	// Traverse all regions and calculate the erosion rate.
	for r, rSlope := range slope {
		// NOTE: This was directly taken from mewo2's code.
		//
		// Honestly, I have no idea how this works, especially how the river
		// creep is supposed to work or where these constants come from.
		river := math.Sqrt(flux[r]/maxFlux) * rSlope
		creep := rSlope * rSlope
		total := 1000*river + creep
		if total > 200 {
			total = 200
		}

		// Determine how "deep" we traverse the neighbor-graph based on
		// the type of erosion we want to perform.
		erodeNbs := erodeNeighborsRiver

		// If we have avg. temp. below 0, we need to imitate glacial erosion, which
		// is carving "wider" valleys than hydraulic erosion.
		if m.getRegTemperature(r, maxH) < 0 {
			erodeNbs = erodeNeighborsGlacier
		}

		// Traverse the neighbor-graph up to the determined depth and calculate
		// the remaining erosion affecting them given their distance to r.
		erodeRegion(r, erodeNbs, total)
	}
	return newh
}

// GetErosionRate2 is an alternative erosion calculation which takes in account
// the steepness and flux of each region to determine the shape of eroded
// riverbeds and valleys.
func (m *Geo) GetErosionRate2() []float64 {
	const (
		distRegions = 3.0 // Number of regions to traverse at max flux.
		minExp      = 1.0 // Minimum exponent (see below)
		varExp      = 3.0 // Variable exponent (see further below)
	)

	// For now we set the maximum erosion distance to 'distRegions' times the distance
	// between two regions. This is a bit arbitrary, but it seems to work well.
	distRegion := math.Sqrt(4 * math.Pi / float64(m.mesh.numRegions))
	maxErosionDistance := distRegions * distRegion

	// Get the steepness of each region to its downhill neighbor.
	steeps := m.GetSteepness()

	// This will collect the erosion values for each region.
	toE := make([]float64, m.mesh.numRegions)

	// If true, we do not consider erosion below sea level.
	erodeOnlyAboveSealevel := true

	// Get the water flux for each region.
	flux := m.getFlux(erodeOnlyAboveSealevel)

	// Get the maximum flux value so we can normalize the flux values
	_, maxFlux := minMax(flux)

	for r, fl := range flux {
		// The steeper the incline, the more V-shaped the riverbed becomes,
		// while a river on a low incline will carve a more U shaped riverbed.
		//
		// y = pow(|x|, e)
		//
		// The exponent (e) determines the shape of the eroded ravine, or the
		// remaining height (y) at a certain distance (x) from the river center.
		//         _  _
		// e = 1:   \/
		//         _   _
		// e = 2:   \_/
		//         _   _
		// e = 4+:  |_|
		//
		// Therefore we calculate the exponent (e) for the function shaping
		// the riverbed (shape of erosion pattern) based on the steepness of the region.
		exponent := minExp + varExp*(1.0-steeps[r])

		// The amount of flux determines how wide the resulting eroded path is.
		//
		// NOTE: We assume that a higher flux will mean that theoretically the resulting river
		// will become wider due to the additional required space for transporting the larger
		// water volume.
		fluxVal := math.Sqrt(fl / maxFlux)

		// Calculate maximum erosion distance based on the water flow intensity aka flux.
		//
		// NOTE: Since rivers slow down the wider they are, they start to be less erosive with higher flux.
		//
		// WARNING: The maxDist value currently represents both the max width of the riverbed as well as the
		// maximum erosion value at the river center (r). This is not really ideal, but it works.
		maxDist := fluxVal * maxErosionDistance
		if maxDist == 0 {
			continue
		}

		// Assign the erosion intensity for (r) if it is larger
		// than the current value set for it.
		// TODO: sharp drops should carve with higher intensity.
		toE[r] = math.Max(maxDist, toE[r])

		// Get lat/lon of the current region to calculate the great arc
		// distance of each neighbor visited by doErode().
		rLatLon := m.LatLon[r]

		// If we have a downhill neighbor, use its lat lon coordinates
		// for the arc segment distance.
		dLatLon := rLatLon
		if rdh := m.Downhill[r]; rdh >= 0 {
			dLatLon = m.LatLon[rdh]
		}

		// seen will keep track of all regions whose neighbors doErode()
		// has already visited (to prevent infinite recursion).
		seen := make(map[int]bool)

		// doErode() is the erosion function which will determine how
		// much a region will be eroded depending on the variables
		// we calculated above and the distance of the neighbor
		// region (reg) to the river center region (r).
		var doErode func(reg int)

		doErode = func(reg int) {
			// Skip visiting all neighbor regions we have already visited.
			if seen[reg] {
				return
			}
			// Make sure we note that we (will have) visited all neighbors of reg.
			seen[reg] = true

			// Visit all neighbors of reg and calculate the erosion rate for each.
			for _, nb := range m.GetRegNeighbors(reg) {
				nbLatLon := m.LatLon[nb]

				// Calculate distance to the arc/line segment of the river (r->rdh).
				dist := crossArc(rLatLon[0], rLatLon[1], dLatLon[0], dLatLon[1], nbLatLon[0], nbLatLon[1])

				// Skip everything that is too far away.
				if math.IsNaN(dist) || dist > maxDist {
					continue
				}

				// Calculate distance as value from 0.0 to 1.0, representing
				// the range of 0.0 to maxDist.
				distRes := dist / maxDist

				// Now apply the formula that will shape the riverbed we
				// have described above.
				erode := math.Pow(distRes, exponent)

				// Check if we've exceeded the maximum erosion value distance.
				if erode > maxDist {
					return
				}

				// NOTE: Since the formula gives us a higher value with increasing
				// distance, we use the value to REDUCE the maximum erosion value
				// (maxDist). The further we are away from (r), the weaker the
				// erosion.
				toErode := maxDist - erode

				// If we have already an erosion value for (nb), we make sure
				// to use the higher erosion value.
				toE[nb] = math.Max(toErode, toE[nb])

				// Visit all neighbors of nb and calculate their erosion values.
				doErode(nb)
			}
		}

		// Call the recursive erosion function for the current region (r).
		// This will calculate the erosion values for all neighbors of (r).
		doErode(r)
	}
	return toE
}
