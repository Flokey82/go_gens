package genworldvoronoi

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/Flokey82/go_gens/gameconstants"
	"github.com/Flokey82/go_gens/genbiome"
	"github.com/davvo/mercator"
	"github.com/llgcode/draw2d/draw2dimg"

	geojson "github.com/paulmach/go.geojson"
)

// GetTile returns the image of the tile at the given coordinates and zoom level.
func (m *Map) GetTile(x, y, zoom int) image.Image {
	// Skip drawing rivers for now.
	drawRivers := true

	// Skip drawing wind vectors for now.
	drawWindVectors := false

	// Wrap the tile coordinates.
	x, y = wrapTileCoordinates(x, y, zoom)

	// log.Println("tile", x, y, "zoom", zoom)

	tbb := newTileBoundingBox(x, y, zoom)
	la1, lo1, la2, lo2 := tbb.toLatLon()
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
	out_t := make([]int, 0, 6)
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
		for _, j := range m.mesh.r_circulate_t(out_t, i) {
			tLat := m.triLatLon[j][0]
			tLon := m.triLatLon[j][1]

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

	// Draw all the wind vectors on top.
	if drawWindVectors {
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
			// Now draw the wind vector for the region.
			//windVec := m.RegionToWindVec[i]
			windVec := m.RegionToWindVecLocal[i]
			// Calculate the coordinates of the center of the region.
			x, y := latLonToPixels(rLat, rLon, zoom)
			x -= dx
			y -= dy2

			// Calculate the length of the wind vector.
			length := math.Sqrt(windVec[0]*windVec[0] + windVec[1]*windVec[1])

			// Calculate the angle of the wind vector.
			angle := math.Atan2(windVec[1], windVec[0])

			// Calculate the coordinates of the end of the wind vector.
			// Since we are on a computer screen, we need to flip the y-axis.
			x2 := x + math.Cos(angle)*length*50
			y2 := y - math.Sin(angle)*length*50

			// Draw the wind vector.
			gc.SetStrokeColor(color.NRGBA{0, 0, 0, 255})
			gc.SetLineWidth(1)
			gc.BeginPath()
			gc.MoveTo(x, y)
			gc.LineTo(x2, y2)
			gc.Stroke()

			// Draw the arrow head.
			gc.SetStrokeColor(color.NRGBA{0, 0, 0, 255})
			gc.SetLineWidth(1)
			gc.BeginPath()
			gc.MoveTo(x2, y2)
			gc.LineTo(x2-math.Cos(angle+math.Pi/6)*5, y2+math.Sin(angle+math.Pi/6)*5)
			gc.Stroke()

			gc.SetStrokeColor(color.NRGBA{0, 0, 0, 255})
			gc.SetLineWidth(1)
			gc.BeginPath()
			gc.MoveTo(x2, y2)
			gc.LineTo(x2-math.Cos(angle-math.Pi/6)*5, y2+math.Sin(angle-math.Pi/6)*5)
			gc.Stroke()
		}
	}

	// Now we do something completely inefficient and
	// fetch all the rivers and filter them by the tile.
	// We should filter this stuff before we generate the rivers.
	if drawRivers {
		rivers := m.getRiversInLatLonBB(0.001/float64(int(1)<<zoom), la1-latLonMargin, lo1-latLonMargin, la2+latLonMargin, lo2+latLonMargin)
		_, maxFlux := minMax(m.Flux)

		for _, river := range rivers {
			// TODO: Fix wrapping around the world.
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
				x -= dx
				y -= dy2
				gc.LineTo(x, y)
				// TODO: Use steepness to determine the amplitude of meandering.
				// The less steep the river is, the more it meanders.
				// lx, ly := gc.LastPoint()
				// gc.CubicCurveTo((x+2*lx)/3, ly, lx, (ly+2*y)/3, x, y)
			}
			gc.Stroke()
			gc.Close()
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
func (m *Map) GetGeoJSONCities(la1, lo1, la2, lo2 float64, zoom int) ([]byte, error) {
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

	// Get the last settled year.
	_, maxSettled := minMax64(m.Settled)
	distRegion := math.Sqrt(4 * math.Pi / float64(m.mesh.numRegions))

	biomeFunc := m.getRegWhittakerModBiomeFunc()
	_, maxElev := minMax(m.Elevation)
	_, maxMois := minMax(m.Moisture)

	// Loop through all the cities and check if they are within the tile.
	// TODO: Just show the largest cities for lower zoom levels.
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
		f.SetProperty("name", c.Name+fmt.Sprintf("lat: %f, lon: %f", cLat, cLon))
		f.SetProperty("type", c.Type)
		f.SetProperty("culture", c.Culture.Name)
		f.SetProperty("population", c.Population)
		f.SetProperty("popgrowth", c.PopulationGrowthRate())
		f.SetProperty("maxpop", c.MaxPopulation)
		f.SetProperty("maxpoplimit", c.MaxPopulationLimit())
		f.SetProperty("settled", maxSettled-c.Founded)
		temperature := m.getRegTemperature(c.ID, maxElev)
		precip := maxPrecipitation * m.Moisture[c.ID] / maxMois
		elev := maxAltitudeFactor * m.Elevation[c.ID] / maxElev
		f.SetProperty("biome", genbiome.WhittakerModBiomeToString(biomeFunc(c.ID))+
			fmt.Sprintf(" (%.1f°C, %.1fdm, %.1fm)", temperature, precip, elev))
		f.SetProperty("attractiveness", c.Attractiveness)
		f.SetProperty("economic", c.EconomicPotential)
		f.SetProperty("agriculture", c.Agriculture)
		f.SetProperty("trade", c.Trade)
		f.SetProperty("resources", c.Resources)
		f.SetProperty("radius", (c.radius()+2*distRegion)*gameconstants.EarthCircumference/(2*math.Pi))
		f.SetProperty("tradepartners", c.TradePartners)
		var msgs []string
		hist := m.History.GetEvents(c.ID, ObjectTypeCity)

		// Only show the last 10 events.
		numEvents := len(hist)
		if numEvents > 10 {
			numEvents = 10
		}
		for _, event := range hist[len(hist)-numEvents:] {
			msgs = append(msgs, event.String())
		}
		f.SetProperty("history", msgs)

		// Generate the list of local resources.
		var resources []string
		// Metals.
		for i := 0; i < ResMaxMetals; i++ {
			if m.Metals[c.ID]&(1<<i) != 0 {
				resources = append(resources, metalToString(i))
			}
		}
		// Gems.
		for i := 0; i < ResMaxGems; i++ {
			if m.Gems[c.ID]&(1<<i) != 0 {
				resources = append(resources, gemToString(i))
			}
		}
		// Stones.
		for i := 0; i < ResMaxStones; i++ {
			if m.Stones[c.ID]&(1<<i) != 0 {
				resources = append(resources, stoneToString(i))
			}
		}
		// Woods.
		for i := 0; i < ResMaxWoods; i++ {
			if m.Wood[c.ID]&(1<<i) != 0 {
				resources = append(resources, woodToString(i))
			}
		}
		// Various.
		for i := 0; i < ResMaxVarious; i++ {
			if m.Various[c.ID]&(1<<i) != 0 {
				resources = append(resources, variousToString(i))
			}
		}
		f.SetProperty("reslist", resources)

		geoJSON.AddFeature(f)
	}

	log.Printf("%d out of %d cities in tile", len(geoJSON.Features), len(m.Cities))

	// Now encode the GeoJSON.
	geoJSONBytes, err := geoJSON.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return geoJSONBytes, nil
}

// GetGeoJSONBorders returns all borders as GeoJSON within the given bounds and zoom level.
func (m *Map) GetGeoJSONBorders(la1, lo1, la2, lo2 float64, zoom int) ([]byte, error) {
	geoJSON := geojson.NewFeatureCollection()

	// Get all borders and add them to the GeoJSON.
	// Right now we ignore the bounds and zoom level.
	for i, border := range m.getBorders() {
		// Now get the coordinates for each point of the border.
		var borderLatLons [][]float64
		for _, p := range border {
			// Get the lat lon coordinates of the point.
			la := m.triLatLon[p][0]
			lo := m.triLatLon[p][1]
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
		return nil, err
	}
	return geoJSONBytes, nil
}

const tileSize = 256

// sizeFromZoom returns the expected size of the world for the mercato projection used below.
func sizeFromZoom(zoom int) int {
	return int(math.Pow(2.0, float64(zoom)) * float64(tileSize))
}

func latLonToPixels(lat, lon float64, zoom int) (float64, float64) {
	return mercator.LatLonToPixels(-1*lat, lon, zoom)
}

// tileBoundingBox represents a bounding box in pixels for a tile.
type tileBoundingBox struct {
	x1, y1 float64
	x2, y2 float64
	zoom   int
}

// toLatLon returns the lat lon coordinates of the north-west and
// south-east corners of the bounding box.
func (t *tileBoundingBox) toLatLon() (lat1, lon1, lat2, lon2 float64) {
	lat1, lon1 = mercator.PixelsToLatLon(t.x1, t.y1, t.zoom)
	lat2, lon2 = mercator.PixelsToLatLon(t.x2, t.y2, t.zoom)
	return
}

// newTileBoundingBox returns a new tile bounding box for the given tile coordinates
// and zoom level.
func newTileBoundingBox(tx, ty, zoom int) tileBoundingBox {
	return tileBoundingBox{
		x1:   float64(tx * tileSize),
		y1:   float64(ty * tileSize),
		x2:   float64((tx + 1) * tileSize),
		y2:   float64((ty + 1) * tileSize),
		zoom: zoom,
	}
}

// boundingBoxResult contains the results of a bounding box query.
type boundingBoxResult struct {
	Regions   []int // Regions withi the bounding box.
	Triangles []int // Triangles within the bounding box.
}

// getBoundingBoxRegions returns all regions and triangles within the given lat/lon bounding box.
//
// TODO: Add margin in order to also return regions/triangles that are partially
// within the bounding box.
func (m *BaseObject) getBoundingBoxRegions(lat1, lon1, lat2, lon2 float64) *boundingBoxResult {
	r := &boundingBoxResult{}
	// TODO: Add convenience function to check against bounding box.
	for i, ll := range m.LatLon {
		if l0, l1 := ll[0], ll[1]; l0 < lat1 || l0 >= lat2 || l1 < lon1 || l1 >= lon2 {
			continue
		}
		r.Regions = append(r.Regions, i)
	}
	for i, ll := range m.triLatLon {
		if l0, l1 := ll[0], ll[1]; l0 < lat1 || l0 >= lat2 || l1 < lon1 || l1 >= lon2 {
			continue
		}
		r.Triangles = append(r.Triangles, i)
	}
	return r
}