package genworldvoronoi

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/Flokey82/go_gens/gameconstants"
	"github.com/Flokey82/go_gens/genbiome"
	"github.com/Flokey82/go_gens/vectors"
	"github.com/davvo/mercator"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/mazznoer/colorgrad"

	geojson "github.com/paulmach/go.geojson"
)

// GetTile returns the image of the tile at the given coordinates and zoom level.
func (m *Map) GetTile(x, y, zoom, displayMode int, drawWindVectors, drawRivers, drawShadows bool) image.Image {
	var colorFunc func(int) color.Color

	switch displayMode {
	case 13, 14, 15, 16:
		colorGrad := colorgrad.Rainbow()
		terrToColor := make(map[int]int)
		var territory []int
		var terrLen int
		if displayMode == 13 {
			terr := m.Cities[:m.NumCityStates]
			terrLen = len(terr)
			for i, c := range terr {
				terrToColor[c.ID] = i
			}
			territory = m.RegionToCityState
		} else if displayMode == 14 {
			terr := m.Cities[:m.NumEmpires]
			terrLen = len(terr)
			for i, c := range terr {
				terrToColor[c.ID] = i
			}
			territory = m.RegionToEmpire
		} else if displayMode == 15 {
			terr := m.Cultures
			terrLen = len(terr)
			for i, c := range terr {
				terrToColor[c.ID] = i
			}
			territory = m.RegionToCulture
		} else {
			terr := m.Species
			terrLen = len(terr)
			for i, c := range terr {
				terrToColor[c.Origin] = i
			}
			territory = m.SpeciesRegions
		}

		min, max := minMax(m.Elevation)
		_, maxMois := minMax(m.Moisture)
		cols := colorGrad.Colors(uint(terrLen))
		colorFunc = func(i int) color.Color {
			// Calculate the color of the region.
			rLat := m.LatLon[i][0]
			elev := m.Elevation[i]
			val := (elev - min) / (max - min)
			if territory[i] == -1 {
				if elev <= 0 {
					return genBlue(val)
				} else {
					valElev := elev / max
					valMois := m.Moisture[i] / maxMois
					return getWhittakerModBiomeColor(rLat, valElev, valMois, val)
				}
			}
			terrID := terrToColor[territory[i]]
			return genColor(cols[terrID], val)
		}
	default:
		vals := m.Elevation
		if displayMode == 1 {
			vals = m.Moisture
		} else if displayMode == 2 {
			vals = m.Rainfall
		} else if displayMode == 3 {
			vals = m.Flux
		} else if displayMode == 4 {
			vals = m.propagateCompression(m.RegionCompression)
		} else if displayMode == 5 {
			vals = m.getEarthquakeChance()
		} else if displayMode == 6 {
			vals = m.getVolcanoEruptionChance()
		} else if displayMode == 7 {
			vals = m.getRockSlideAvalancheChance()
		} else if displayMode == 8 {
			vals = m.getFloodChance()
		} else if displayMode == 9 {
			vals = m.GetErosionRate()
		} else if displayMode == 10 {
			vals = m.GetErosionRate2()
		} else if displayMode == 11 {
			vals = m.GetSteepness()
		} else if displayMode == 12 {
			vals = m.GetSlope()
		}

		// Calculate the min and max elevation.
		_, max := minMax(m.Elevation)
		_, maxMois := minMax(m.Moisture)
		minVal, maxVal := minMax(vals)
		colorFunc = func(i int) color.Color {
			// Calculate the color of the region.
			rLat := m.LatLon[i][0]
			elev := m.Elevation[i]
			val := (vals[i] - minVal) / (maxVal - minVal)
			var col color.NRGBA
			if elev <= 0 {
				col = genBlue(val)
			} else {
				valElev := elev / max
				valMois := m.Moisture[i] / maxMois
				col = getWhittakerModBiomeColor(rLat, valElev, valMois, val)
			}
			return col
		}
	}

	// Wrap the tile coordinates.
	x, y = wrapTileCoordinates(x, y, zoom)

	tbb := newTileBoundingBox(x, y, zoom)
	la1, lo1, la2, lo2 := tbb.toLatLon()
	latLonMargin := 20 / float64(zoom)

	la1Margin := math.Max(-90, math.Min(90, la1-latLonMargin))
	la2Margin := math.Max(-90, math.Min(90, la2+latLonMargin))
	lo1Margin := lo1 - latLonMargin
	lo2Margin := lo2 + latLonMargin

	// Since our mercator conversion gives us absolute pixel coordinates, we need to
	// remove the offset of the tile we are rendering from the path coordinates.
	dx, _ := latLonToPixels(la1, lo1, zoom)
	_, dy2 := latLonToPixels(la2, lo2, zoom)

	// Create a new image to draw the tile on.
	dest := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	gc := draw2dimg.NewGraphicContext(dest)

	out_t := make([]int, 0, 6)
	gc.SetLineWidth(1)
	for i := 0; i < m.mesh.numRegions; i++ {
		rLat := m.LatLon[i][0]
		rLon := m.LatLon[i][1]

		// Check if we are within the tile with a small margin,
		// taking into account that we might have wrapped around the world.
		if rLat < la1Margin || rLat > la2Margin || rLon < lo1Margin || rLon > lo2Margin {
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
		col := colorFunc(i)

		// If the path is empty, we can skip it.
		if len(path) == 0 {
			continue
		}

		// Draw the path.
		gc.SetStrokeColor(col)
		gc.SetFillColor(col)
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
		// Set the color and line width of the wind vectors.
		gc.SetStrokeColor(color.NRGBA{0, 0, 0, 255})
		gc.SetLineWidth(1)
		for i := 0; i < m.mesh.numRegions; i++ {
			rLat := m.LatLon[i][0]
			rLon := m.LatLon[i][1]

			// Check if we are within the tile with a small margin,
			// taking into account that we might have wrapped around the world.
			if rLat < la1Margin || rLat > la2Margin || rLon < lo1Margin || rLon > lo2Margin {
				// Check if the tile and the region we are looking at is adjecent to +/- 180 degrees and
				// NOTE: This could be improved by checking if one of the corners of the region is within the tile.
				if lo1 > -175 && lo2 < 175 || rLon < 175 && rLon > -175 {
					continue
				}
			}
			// Now draw the wind vector for the region.
			// windVec := m.RegionToWindVec[i]
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
			gc.BeginPath()
			gc.MoveTo(x, y)
			gc.LineTo(x2, y2)
			gc.Stroke()

			// Draw the arrow head.
			gc.BeginPath()
			gc.MoveTo(x2, y2)
			gc.LineTo(x2-math.Cos(angle+math.Pi/6)*5, y2+math.Sin(angle+math.Pi/6)*5)
			gc.Stroke()

			gc.BeginPath()
			gc.MoveTo(x2, y2)
			gc.LineTo(x2-math.Cos(angle-math.Pi/6)*5, y2+math.Sin(angle-math.Pi/6)*5)
			gc.Stroke()
		}
	}

	if drawShadows {
		min, max := minMax(m.triElevation)
		if max == 0 {
			max = 1
		}
		_, maxMois := minMax(m.triMoisture)
		if maxMois == 0 {
			maxMois = 1
		}

		// Set our initial line width.
		gc.SetLineWidth(1)

		// Set the global light direction almost straight up, with a slight offset to the right.
		lightDir := vectors.Vec3{X: 1.0, Y: 1.0, Z: 1.0}.Normalize()
	Loop:
		for i := 0; i < len(m.mesh.Triangles); i += 3 {
			// Hacky way to filter paths/triangles that wrap around the entire SVG.
			triLat := m.triLatLon[i/3][0]
			triLon := m.triLatLon[i/3][1]

			// Check if we are within the tile with a small margin, taking into account that we might have wrapped around the world.
			// Also keep in mind that the latitude and longitude can be negative, so we need to add the margin to the lower bound,
			// and subtract the margin from the upper bound.
			if triLat < la1Margin || triLat > la2Margin || triLon < lo1Margin || triLon > lo2Margin {
				// Check if the tile and the region we are looking at is adjecent to +/- 180 degrees and
				// NOTE: This could be improved by checking if one of the corners of the region is within the tile.
				if lo1 > -175 && lo2 < 175 || triLon < 175 && triLon > -175 {
					continue
				}
			}
			var poolCount int
			for _, j := range m.mesh.Triangles[i : i+3] {
				if m.Waterpool[j] > 0 {
					poolCount++
				}
			}

			// Draw the path that outlines the region.
			var path [][2]float64
			for _, j := range m.mesh.t_circulate_r(out_t, i/3) {
				rLat := m.LatLon[j][0]
				rLon := m.LatLon[j][1]

				// Check if we the region is across the +/- 180 degrees longitude line compared to the triangle.
				// In this case, the longitude is almost 360 degrees off, which means we need to adjust the longitude.
				if rLon-triLon > 110 {
					rLon -= 360
				} else if rLon-triLon < -110 {
					rLon += 360
				}

				// Calculate the coordinates of the path point.
				x, y := latLonToPixels(rLat, rLon, zoom)
				p := [2]float64{(x - dx), (y - dy2)}
				// Check if we are way outside the tile.
				if p[0] < -1000 || p[0] > 1000 || p[1] < -1000 || p[1] > 1000 {
					continue Loop
				}
				path = append(path, p)
			}

			// Now check if the region we are looking at has wrapped around the world /
			// +- 180 degrees. If so, we need to adjust the points in the path.
			if lo1 < -175 && triLon > 175 {
				for i := range path {
					path[i][0] -= float64(sizeFromZoom(zoom))
				}
			} else if lo2 > 175 && triLon < -175 {
				for i := range path {
					path[i][0] += float64(sizeFromZoom(zoom))
				}
			}

			elev := m.triElevation[i/3]
			val := (elev - min) / (max - min)
			var col color.NRGBA
			if elev <= 0 || poolCount > 2 {
				col = genBlue(val)
			} else {
				// Get the slope of the triangle.
				slope := m.regTriNormal(i/3, m.mesh.t_circulate_r(out_t, i/3))

				// Now take the dot product of the slope and our global
				// light direction to get the amount of light on the triangle.
				light := math.Max(0, vectors.Dot3(slope, lightDir))

				// Calculate the brightness of the triangle.
				// For shaded reliefs the contrast should increase by elevation.
				// http://www.reliefshading.com/design/
				brightness := val * (1 - val*(1-light))
				col = getWhittakerModBiomeColor(triLat, elev/max, m.triMoisture[i/3]/maxMois, brightness)
			}

			// If the path is empty, we can skip it.
			if len(path) == 0 {
				continue
			}

			// Draw the path.
			gc.SetStrokeColor(col)
			gc.SetFillColor(col)
			gc.BeginPath()
			gc.MoveTo(path[0][0], path[0][1])
			for _, p := range path[1:] {
				gc.LineTo(p[0], p[1])
			}
			gc.Close()
			gc.FillStroke()
		}
	}

	// Now we do something completely inefficient and
	// fetch all the rivers and filter them by the tile.
	// We should filter this stuff before we generate the rivers.
	if drawRivers {
		rivers := m.getRiversInLatLonBB(0.001/float64(int(1)<<zoom), la1-latLonMargin, lo1-latLonMargin, la2+latLonMargin, lo2+latLonMargin)
		_, maxFlux := minMax(m.Flux)

		// Set our stroke color to a nice river blue.
		gc.SetStrokeColor(color.NRGBA{0, 0, 255, 255})

		for _, river := range rivers {
			// Set the initial line width.
			gc.SetLineWidth(1)
			gc.BeginPath()

			// Move to the first point.
			rLat, rLon := m.LatLon[river[0]][0], m.LatLon[river[0]][1]
			x, y := latLonToPixels(rLat, rLon, zoom)
			gc.MoveTo(x-dx, y-dy2)
			for i, p := range river[1:] {
				// Set the line width based on the flux of the river, averaged with the previous flux.
				gc.SetLineWidth(4 * math.Sqrt((m.Flux[p]+m.Flux[river[i]])/(2*maxFlux)))

				// Set the line width based on the flux of the river.
				rLat, rLon = m.LatLon[p][0], m.LatLon[p][1]

				// Now compare the longitude to the previous longitude.
				// If we have crossed the +- 180 degree boundary, we need to
				// draw to a fake point at the same latitude but on the same side of the world.
				if diff := rLon - m.LatLon[river[i]][1]; math.Abs(diff) > 110 {
					rLonFake := rLon - 360
					if diff < 0 {
						rLonFake = rLon + 360
					}
					// Draw to the fake point.
					x, y := latLonToPixels(rLat, rLonFake, zoom)
					gc.LineTo(x-dx, y-dy2)
					gc.Stroke()

					// Move to the real point and start a new path.
					x, y = latLonToPixels(rLat, rLon, zoom)
					gc.BeginPath()
					gc.MoveTo(x-dx, y-dy2)
				}

				x, y := latLonToPixels(rLat, rLon, zoom)
				x -= dx
				y -= dy2

				// If we are below sea level, interpolate the point with the previous point.
				if m.Elevation[p] < 0 {
					// Draw from the last position to the midpoint.
					// This will cause the river to end at the sea level.
					lx, ly := gc.LastPoint()
					gc.LineTo((x+lx)/2, (y+ly)/2)

					// Move to the new point.
					gc.MoveTo(x, y)
				} else if m.Elevation[river[i]] < 0 {
					// If the previous point was below sea level, interpolate the point with the next point.
					// This will cause the river to start at the sea level.
					lx, ly := gc.LastPoint()
					gc.MoveTo((x+lx)/2, (y+ly)/2)

					// Draw to the new point.
					gc.LineTo(x, y)
				} else {
					gc.LineTo(x, y)
				}
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

	regPropertyFunc := m.getRegPropertyFunc()

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
		f.SetProperty("name", c.Name)
		f.SetProperty("type", c.Type)
		f.SetProperty("culture", fmt.Sprintf("%s (%s)", c.Culture.Name, c.Culture.Type))
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
		f.SetProperty("coordinates", fmt.Sprintf("lat %.2f, lon %.2f", cLat, cLon))
		f.SetProperty("attractiveness", c.Attractiveness)
		f.SetProperty("economic", c.EconomicPotential)
		f.SetProperty("agriculture", c.Agriculture)
		f.SetProperty("trade", c.Trade)
		f.SetProperty("resources", c.Resources)
		f.SetProperty("radius", (c.radius()+2*distRegion)*gameconstants.EarthCircumference/(2*math.Pi))
		f.SetProperty("tradepartners", c.TradePartners)
		f.SetProperty("flavortext", m.generateCityFlavorText(c, regPropertyFunc(c.ID)))
		var sName string
		if m.SpeciesRegions[c.ID] >= 0 {
			var s *Species
			for _, sp := range m.Species {
				if sp.Origin == m.SpeciesRegions[c.ID] {
					s = sp
					break
				}
			}
			if s != nil {
				sName = s.Name
				if sName == "" {
					sName = s.String()
				}
			}
		}
		f.SetProperty("species", sName)
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
func (m *Map) GetGeoJSONBorders(la1, lo1, la2, lo2 float64, zoom, displayMode int) ([]byte, error) {
	geoJSON := geojson.NewFeatureCollection()
	var borders [][]int
	switch displayMode {
	case 1:
		borders = m.getCustomBorders(m.RegionToCityState)
	case 2:
		borders = m.getCustomBorders(m.RegionToCulture)
	case 3:
		borders = m.getCustomBorders(m.RegionToPlate)
	case 4:
		borders = m.getCustomBorders(m.BiomeRegions)
	case 5:
		// Nothing.
	default:
		borders = m.getCustomBorders(m.RegionToEmpire)
	}

	// Get all borders and add them to the GeoJSON.
	// Right now we ignore the bounds and zoom level.
	for i, border := range borders {
		// Now get the coordinates for each point of the border.
		var borderLatLons [][]float64
		for _, p := range border {
			// Get the lat lon coordinates of the point.
			la := m.triLatLon[p][0]
			lo := m.triLatLon[p][1]

			// Check if we have crossed the 180 degree longitude line.
			// If so, we stop here, add the border to the GeoJSON and start a new one.
			if len(borderLatLons) > 0 && math.Abs(borderLatLons[len(borderLatLons)-1][0]-lo) > 180 {
				// Add the border to the GeoJSON as a feature.
				f := geojson.NewLineStringFeature(borderLatLons)
				f.ID = i
				geoJSON.AddFeature(f)

				// Start a new border.
				borderLatLons = nil
			}
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
