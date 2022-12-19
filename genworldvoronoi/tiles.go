package genworldvoronoi

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/llgcode/draw2d/draw2dimg"
	geojson "github.com/paulmach/go.geojson"
)

// GetTile returns the image of the tile at the given coordinates and zoom level.
func (m *Map) GetTile(x, y, zoom int) image.Image {
	// Skip drawing rivers for now.
	drawRivers := false

	// Wrap the tile coordinates.
	x, y = wrapTileCoordinates(x, y, zoom)

	// log.Println("tile", x, y, "zoom", zoom)

	tbb := tileBoundingBox(x, y, zoom)
	la1, lo1, la2, lo2 := tbb.ToLatLon()
	latLonMargin := 20.0 / float64(zoom)

	// Since our mercator conversion gives us absolute pixel coordinates, we need to
	// remove the offset of the tile we are rendering from the path coordinates.
	dx, _ := latLonToPixels(la1, lo1, zoom)
	_, dy2 := latLonToPixels(la2, lo2, zoom)

	// Create a new image to draw the tile on.
	dest := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	gc := draw2dimg.NewGraphicContext(dest)

	// Calculate the min and max elevation and moisture values.
	min, max := minMax(m.Elevation)
	_, maxMois := minMax(m.Moisture)
	for i := 0; i < m.mesh.numRegions; i++ {
		rLat := m.LatLon[i][0]
		rLon := m.LatLon[i][1]

		// Check if we are within the tile with a small margin,
		// taking into account that we might have wrapped around the world.
		if rLat < la1-latLonMargin || rLat > la2+latLonMargin || rLon < lo1-latLonMargin || rLon > lo2+latLonMargin {
			// Check if the tile and the region we are looking at is adjecent to +/- 180 degrees and
			// NOTE: This could be improved by checking if one of the corners of the region is within the tile.
			if lo1 > -175 && lo2 < 175 || rLon < 175 && rLon > -175 {
				continue
			}
		}

		// Draw the path that outlines the region.
		var path [][2]float64
		for _, j := range m.mesh.r_circulate_t(nil, i) {
			tLat := m.t_latLon[j][0]
			tLon := m.t_latLon[j][1]

			// Check if we have wrapped around the world.
			if tLon-rLon > 120 {
				tLon -= 360
			} else if tLon-rLon < -120 {
				tLon += 360
			}

			// Calculate the coordinates of the path point.
			x, y := latLonToPixels(tLat, tLon, zoom)
			path = append(path, [2]float64{(x - dx), (y - dy2)})
		}

		// Now check if the region we are looking at has wrapped around the world /
		// +- 180 degrees. If so, we need to adjust the points in the path.
		if lo1 < -175 && rLon > 175 {
			for i := range path {
				path[i][0] -= float64(sizeFromZoom(zoom))
			}
		} else if lo2 > 175 && rLon < -175 {
			for i := range path {
				path[i][0] += float64(sizeFromZoom(zoom))
			}
		}

		// Calculate the color of the region.
		elev := m.Elevation[i]
		val := (elev - min) / (max - min)
		var col color.NRGBA
		if elev <= 0 {
			col = genBlue(val)
		} else {
			valElev := elev / max
			valMois := m.Moisture[i] / maxMois
			col = getWhittakerModBiomeColor(rLat, valElev, valMois, val)
		}

		// If the path is empty, we can skip it.
		if len(path) == 0 {
			continue
		}

		// Draw the path.
		gc.SetFillColor(col)
		gc.SetStrokeColor(col)
		gc.SetLineWidth(0)
		gc.BeginPath()
		gc.MoveTo(path[0][0], path[0][1])
		for _, p := range path[1:] {
			gc.LineTo(p[0], p[1])
		}
		gc.Close()
		gc.FillStroke()
	}

	// Now we do something completely inefficient and
	// fetch all the rivers and filter them by the tile.
	// We should filter this stuff before we generate the rivers.
	if drawRivers {
		rivers := m.getRivers(0.0001)
		_, maxFlux := minMax(m.Flux)

		for _, river := range rivers {
			gc.SetStrokeColor(color.NRGBA{0, 0, 255, 255})
			gc.SetLineWidth(1)
			gc.BeginPath()
			// Move to the first point.
			rLat, rLon := m.LatLon[river[0]][0], m.LatLon[river[0]][1]
			x, y := latLonToPixels(rLat, rLon, zoom)
			gc.MoveTo(x-dx, y-dy2)
			for _, p := range river[1:] {
				gc.SetLineWidth(1 + 2*(m.Flux[p]/maxFlux))

				// Set the line width based on the flux of the river.
				rLat, rLon = m.LatLon[p][0], m.LatLon[p][1]
				x, y := latLonToPixels(rLat, rLon, zoom)
				gc.LineTo(x-dx, y-dy2)
			}
			gc.Stroke()
		}
	}
	return dest
}

func wrapTileCoordinates(x, y, zoom int) (int, int) {
	// Wrap the tile coordinates.
	x = x % (1 << uint(zoom))
	if x < 0 {
		x += 1 << uint(zoom)
	}
	y = y % (1 << uint(zoom))
	if y < 0 {
		y += 1 << uint(zoom)
	}
	return x, y
}

func wrapLatLon(la, lo float64) (float64, float64) {
	// Wrap the lat lon coordinates.
	la = math.Mod(la, 180)
	if la < -90 {
		la += 180
	} else if la > 90 {
		la -= 180
	}
	lo = math.Mod(lo, 360)
	if lo < -180 {
		lo += 360
	} else if lo > 180 {
		lo -= 360
	}
	return la, lo
}

func limitLatLon(la, lo float64) (float64, float64) {
	// Limit the lat lon coordinates.
	if la < -90 {
		la = -90
	} else if la > 90 {
		la = 90
	}
	if lo < -180 {
		lo = -180
	} else if lo > 180 {
		lo = 180
	}
	return la, lo
}

type latLonBounds struct {
	la1, lo1, la2, lo2 float64
}

// InBounds checks if the given lat lon coordinates are within the bounds.
// NOTE: We need to take in account that the bounds might have wrapped around the world.
func (b latLonBounds) InBounds(la, lo float64) bool {
	if b.la1 < b.la2 {
		// We wrapped around north or south.
		if la > b.la1 && la < b.la2 {
			return false
		}
	} else {
		if la > b.la1 || la < b.la2 {
			return false
		}
	}
	if b.lo1 > b.lo2 {
		// We wrapped around east or west.
		if lo < b.lo1 && lo > b.lo2 {
			return false
		}
	} else {
		if lo < b.lo1 || lo > b.lo2 {
			return false
		}
	}
	return true
}

func wrapLatitude(la float64) float64 {
	// Wrap the latitude.
	la = math.Mod(la, 180)
	if la < -90 {
		la += 180
	} else if la > 90 {
		la -= 180
	}
	return la
}

func wrapLongitude(lo float64) float64 {
	// Wrap the longitude.
	lo = math.Mod(lo, 360)
	if lo < -180 {
		lo += 360
	} else if lo > 180 {
		lo -= 360
	}
	return lo
}

func limitLatitude(la float64) float64 {
	// Limit the latitude.
	if la < -90 {
		la = -90
	} else if la > 90 {
		la = 90
	}
	return la
}

func limitLongitude(lo float64) float64 {
	// Limit the longitude.
	if lo < -180 {
		lo = -180
	} else if lo > 180 {
		lo = 180
	}
	return lo
}

// GetGeoJSONCities returns all cities as GeoJSON within the given bounds and zoom level.
func (m *Map) GetGeoJSONCities(la1, lo1, la2, lo2 float64, zoom int) []byte {
	geoJSON := geojson.NewFeatureCollection()

	// Fix the bounds if la1, lo1, la2, lo2 are not in the correct order.
	//if la1 < la2 {
	//	la1, la2 = la2, la1
	//}
	//if lo1 > lo2 {
	//	lo1, lo2 = lo2, lo1
	//}
	// Limit the lat lon coordinates.
	//la1, lo1 = limitLatLon(la1, lo1)
	//la2, lo2 = limitLatLon(la2, lo2)

	// Wrap the latitude only if we see less than 180 degrees, otherwise just limit it.
	if math.Abs(la1-la2) < 180 {
		la1 = wrapLatitude(la1)
		la2 = wrapLatitude(la2)
	} else {
		la1 = limitLatitude(la1)
		la2 = limitLatitude(la2)
	}
	// Wrap the longitude only if we see less than 360 degrees.
	if math.Abs(lo1-lo2) < 360 {
		lo1 = wrapLongitude(lo1)
		lo2 = wrapLongitude(lo2)
	} else {
		lo1 = limitLongitude(lo1)
		lo2 = limitLongitude(lo2)
	}
	lbb := latLonBounds{la1, lo1, la2, lo2}
	log.Println(la1, lo1, la2, lo2)

	// Loop through all the cities and check if they are within the tile.
	for _, c := range m.Cities {
		cLat := m.LatLon[c.ID][0]
		cLon := m.LatLon[c.ID][1]

		// Check if we are within the tile with a small margin.
		if !lbb.InBounds(cLat, cLon) {
			continue
		}

		// Add the city to the GeoJSON as a feature.
		f := geojson.NewPointFeature([]float64{cLon, cLat})
		f.SetProperty("id", c.ID)
		f.SetProperty("name", c.Name)
		f.SetProperty("type", c.Type)
		f.SetProperty("culture", c.Culture.Name)
		geoJSON.AddFeature(f)
	}

	log.Println("%d out of %d cities in tile", len(geoJSON.Features), len(m.Cities))

	// Now encode the GeoJSON.
	geoJSONBytes, err := geoJSON.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return geoJSONBytes
}

// GetGeoJSONBorders returns all borders as GeoJSON within the given bounds and zoom level.
func (m *Map) GetGeoJSONBorders(la1, lo1, la2, lo2 float64, zoom int) []byte {
	geoJSON := geojson.NewFeatureCollection()

	// Get all borders and add them to the GeoJSON.
	// Right now we ignore the bounds and zoom level.
	for i, border := range m.getBorders() {
		// Now get the coordinates for each point of the border.
		var borderLatLons [][]float64
		for _, p := range border {
			// Get the lat lon coordinates of the point.
			la := m.t_latLon[p][0]
			lo := m.t_latLon[p][1]
			borderLatLons = append(borderLatLons, []float64{lo, la})
		}
		// Add the border to the GeoJSON as a feature.
		f := geojson.NewLineStringFeature(borderLatLons)
		f.ID = i
		geoJSON.AddFeature(f)
	}
	// Now encode the GeoJSON.
	geoJSONBytes, err := geoJSON.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return geoJSONBytes
}
