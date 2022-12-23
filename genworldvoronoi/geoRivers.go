package genworldvoronoi

import (
	"log"
	"time"
)

// getRivers returns the merged river segments whose flux exceeds the provided limit.
// Each river is represented as a sequence of region indices.
func (m *BaseObject) getRivers(limit float64) [][]int {
	// Get segments that are valid river segments.
	links := m.getRiverSegments(limit)

	// Merge the segments that are connected to each other into logical region sequences.
	log.Println("start merge")
	start := time.Now()
	defer func() {
		log.Println("Done river segments in ", time.Since(start).String())
	}()
	return mergeIndexSegments(links)
}

func (m *BaseObject) getRiversInLatLonBB(limit float64, minLat, minLon, maxLat, maxLon float64) [][]int {
	// Get segments that are valid river segments.
	links := m.getRiverSegments(limit)

	// Merge the segments that are connected to each other into logical region sequences.
	log.Println("start merge")
	start := time.Now()
	defer func() {
		log.Println("Done river segments in ", time.Since(start).String())
	}()
	// Filter out all segments that are not in the bounding box.
	var filtered [][2]int
	for _, link := range links {
		lat1, lon1 := m.LatLon[link[0]][0], m.LatLon[link[0]][1]
		lat2, lon2 := m.LatLon[link[1]][0], m.LatLon[link[1]][1]

		// If both points are outside the bounding box, skip the segment.
		if (lat1 < minLat || lat1 > maxLat || lon1 < minLon || lon1 > maxLon) &&
			(lat2 < minLat || lat2 > maxLat || lon2 < minLon || lon2 > maxLon) {
			continue
		}
		filtered = append(filtered, link)
	}
	return mergeIndexSegments(filtered)
}

// getRiverSegments returns all region / downhill neighbor pairs whose flux values
// exceed the provided limit / threshold.
func (m *BaseObject) getRiverSegments(limit float64) [][2]int {
	// NOTE: Should we re-generate downhill and flux, just in case erosion
	// or other factors might have changed this?

	// Get (cached) downhill neighbors.
	dh := m.Downhill

	// Get (cached) flux values.
	flux := m.Flux

	// Adjust the limit to be a fraction of the max flux.
	// This will save us a lot of cycles when comparing
	// flux values to the limit.
	_, maxFlux := minMax(flux)
	limit *= maxFlux

	// Find all link segments that have a high enough flux value.
	var links [][2]int
	for r := 0; r < m.mesh.numRegions; r++ {
		// Skip all regions that are sinks / have no downhill neighbor or
		// regions below sea level.
		if dh[r] < 0 || m.Elevation[r] < 0 {
			continue
		}

		// Skip all regions with flux values that are equal to the rainfall in the region,
		// which is the minimum flux value / the default state for regions without
		// water influx.
		// NOTE: Rivers need at least one contributor region and would therefore have a flux
		// value that is higher than the rainfall in the region.
		if flux[r] <= m.Rainfall[r] || flux[dh[r]] <= m.Rainfall[dh[r]] {
			continue
		}

		// NOTE: Right now we skip segments if both flux values are
		// below the limit.
		if flux[r] >= limit && flux[dh[r]] >= limit {
			// NOTE: The river segment always flows from seg[0] to seg[1].
			links = append(links, [2]int{r, dh[r]})
		}
	}
	return links
}

// getRiverIndices returns a mapping from regions to river ID.
func (m *Geo) getRiverIndices(limit float64) []int {
	// Set up defaults.
	rivers := make([]int, m.mesh.numRegions)
	for i := range rivers {
		rivers[i] = -1 // -1 means no river
	}
	for i, riv := range m.getRivers(limit) {
		for _, idx := range riv {
			rivers[idx] = i
		}
	}
	return rivers
}
