package genworldvoronoi

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/Flokey82/go_gens/genbiome"
	svgo "github.com/ajstarks/svgo"
	"github.com/davvo/mercator"
	"github.com/mazznoer/colorgrad"
	"github.com/sizeofint/webpanimation"
)

const tileSize = 256

// sizeFromZoom returns the expected size of the world for the mercato projection used below.
func sizeFromZoom(zoom int) int {
	return int(math.Pow(2.0, float64(zoom)) * float64(tileSize))
}

func latLonToPixels(lat, lon float64, zoom int) (float64, float64) {
	return mercator.LatLonToPixels(-1*lat, lon, zoom)
}

type tileBB struct {
	x1   float64
	y1   float64
	x2   float64
	y2   float64
	zoom int
}

func (t *tileBB) ToLatLon() (lat1, lon1, lat2, lon2 float64) {
	lat1, lon1 = mercator.PixelsToLatLon(t.x1, t.y1, t.zoom)
	lat2, lon2 = mercator.PixelsToLatLon(t.x2, t.y2, t.zoom)
	return
}

func tileBoundingBox(tx, ty, zoom int) tileBB {
	return tileBB{
		x1:   float64(tx * tileSize),
		y1:   float64(ty * tileSize),
		x2:   float64((tx + 1) * tileSize),
		y2:   float64((ty + 1) * tileSize),
		zoom: zoom,
	}
}

type QueryResult struct {
	r []int
	t []int
}

func (m *BaseObject) getBB(lat1, lon1, lat2, lon2 float64) *QueryResult {
	r := &QueryResult{}
	// TODO: Add convenience function to check against bounding box.
	for i, ll := range m.LatLon {
		if l0, l1 := ll[0], ll[1]; l0 < lat1 || l0 >= lat2 || l1 < lon1 || l1 >= lon2 {
			continue
		}
		r.r = append(r.r, i)
	}
	for i, ll := range m.t_latLon {
		if l0, l1 := ll[0], ll[1]; l0 < lat1 || l0 >= lat2 || l1 < lon1 || l1 >= lon2 {
			continue
		}
		r.t = append(r.t, i)
	}
	log.Println(lat1, lon1, lat2, lon2)
	return r
}

func genBlue(intensity float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(intensity * 255),
		G: uint8(intensity * 255),
		B: 255,
		A: 255,
	}
}

func genGreen(intensity float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(intensity * 255),
		B: uint8((1 - intensity) * 255),
		G: 255,
		A: 255,
	}
}

// ExportSVG exports the terrain as SVG to the given path.
// NOTE: This produces broken somewhat incomplete output due to the wraparound of the mesh.
func (m *Map) ExportSVG(path string) error {
	drawRiversA := true
	drawRiversB := false
	drawFlux := false
	drawDrains := false
	drawCities := true
	drawSinks := false
	drawPools := false
	drawErosion := false
	drawErosion2 := false
	drawHumidity := false
	drawWindOrder := false
	drawRainfall := false
	drawBorders := true
	drawLakeBorders := false
	drawBelow := false
	drawLandmassContour := true
	drawWindDir := false
	drawPlateCompression := false
	drawAltitude := false
	drawTemperature := false
	drawLatitudeDots := false
	drawCityscore := false
	drawRegionTerrain := true
	drawTradeRoutes := false
	drawResources := true

	zoom := 3
	filterPathDist := 20.0
	size := sizeFromZoom(zoom)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	svg := svgo.New(f)
	svg.Start(size, size)
	svg.Style("text/css",
		"path { fill: none; stroke-width: 0.5; }\n"+
			"path.contour{ stroke: black;}\n"+
			".river{ stroke: blue;}\n"+
			"path.lake{ fill: blue; stroke: blue; fill-opacity: 0.5;}\n"+
			"path.border{ stroke: red; stroke-width: 1;"+
			"stroke-linecap: butt;}\n"+
			"path.cityborder{ stroke: black; stroke-width: 1;"+
			"stroke-linecap: butt;}\n"+
			"path.terrain{ stroke: none;}\n"+
			"path.traderoute{"+
			"stroke: lime;"+
			"stroke-width: 0.5;}\n"+
			"text{"+
			"font-weight: bold;"+
			"font-family: \"Palatino Linotype\", \"Book Antiqua\", Palatino, serif;"+
			"fill: white;"+
			"stroke: black;"+
			"stroke-width: 2;"+
			"text-anchor: start;"+
			"stroke-linejoin: round;"+
			"paint-order: stroke;}\n"+
			"text.mine{"+
			"font-size: 6px;}\n"+
			"text.city{"+
			"font-size: 8px;}\n"+
			"text.capital{"+
			"font-size: 12px;}\n")
	em := m
	// Hack to test tile fetching
	// 113.48673955688815 180 139.02010193037987 225
	// 139.02010193037987 180 0 225
	// tbb := tileBoundingBox(0, 0, 0)
	// la1, lo1, la2, lo2 := tbb.ToLatLon()
	// re := m.getBB(la1, lo1, la2, lo2)
	// em, err := m.interpolate(re.r)
	// if err != nil {
	//	panic(err)
	// }
	// end hack

	// Use regions instead of triangles to render terrain.
	if drawRegionTerrain {
		//cities_r := m.cities_r
		//m.cities_r = nil
		//cityScore := m.rCityScore()
		//m.cities_r = cities_r
		//_, maxS := minMax(cityScore)
		//fitScore := m.genNoise()
		//_, maxFit := minMax(fitScore)
		//solarRad := m.calcSolarRadiation(172 / 2)
		//minSol, maxSol := minMax(solarRad)
		min, max := minMax(m.Elevation)
		_, maxMois := minMax(m.Moisture)
		for i := 0; i < em.mesh.numRegions; i++ {
			rLat := em.LatLon[i][0]
			rLon := em.LatLon[i][1]
			rX, rY := latLonToPixels(rLat, rLon, zoom)
			var skip bool
			for _, j := range em.mesh.r_circulate_t(nil, i) {
				x, y := latLonToPixels(em.t_latLon[j][0], em.t_latLon[j][1], zoom)
				if dist2([2]float64{x, y}, [2]float64{rX, rY}) > filterPathDist {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
			var path [][2]float64
			for _, j := range em.mesh.r_circulate_t(nil, i) {
				x, y := latLonToPixels(em.t_latLon[j][0], em.t_latLon[j][1], zoom)
				path = append(path, [2]float64{x, y})
			}
			elev := em.Elevation[i]
			val := (elev - min) / (max - min)
			//val = (solarRad[i] - minSol) / (maxSol - minSol)
			//val = cityScore[i] / maxS
			//val = fitScore[i] / maxFit
			var col color.NRGBA
			if elev <= 0 {
				col = genBlue(val)
			} else {
				valElev := elev / max
				valMois := em.Moisture[i] / maxMois
				col = genbiome.GetWhittakerModBiomeColor(int(getMeanAnnualTemp(rLat)-getTempFalloffFromAltitude(maxAltitudeFactor*valElev)), int(valMois*maxPrecipitation), val)
			}
			svg.Path(svgGenD(path), fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.B), "class=\"terrain\"")
		}
	} else {
		min, max := minMax(m.t_elevation)
		_, maxMois := minMax(m.t_moisture)
		for i := 0; i < len(em.mesh.Triangles); i += 3 {
			// Hacky way to filter paths/triangles that wrap around the entire SVG.
			triLat := em.t_latLon[i/3][0]
			triLon := em.t_latLon[i/3][1]
			triX, triY := latLonToPixels(triLat, triLon, zoom)
			var skip bool
			var poolCount int
			for _, j := range em.mesh.Triangles[i : i+3] {
				if em.Waterpool[j] > 0 {
					poolCount++
				}
				x, y := latLonToPixels(em.LatLon[j][0], em.LatLon[j][1], zoom)
				if dist2([2]float64{x, y}, [2]float64{triX, triY}) > filterPathDist {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			var path [][2]float64
			for _, j := range em.mesh.Triangles[i : i+3] {
				x, y := latLonToPixels(em.LatLon[j][0], em.LatLon[j][1], zoom)
				path = append(path, [2]float64{x, y})
			}
			elev := em.t_elevation[i/3]
			val := (elev - min) / (max - min)
			var col color.NRGBA
			if elev <= 0 || poolCount > 2 {
				col = genBlue(val)
			} else {
				valElev := elev / max
				// Hacky: Modify elevation based on latitude to compensate for colder weather at the poles and warmer weather at the equator.
				// valElev := math.Max(math.Min((elev/max)+(math.Sqrt(math.Abs(triLat)/90.0)-0.5), max), 0)
				valMois := em.t_moisture[i/3] / maxMois
				col = genbiome.GetWhittakerModBiomeColor(int(getMeanAnnualTemp(triLat)-getTempFalloffFromAltitude(maxAltitudeFactor*valElev)), int(valMois*maxPrecipitation), val)
			}
			svg.Path(svgGenD(path), fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.B), "class=\"terrain\"")
		}
	}

	// drawCircle draws a circle at the given lat/lon coordinates.
	drawCircle := func(lat, lon float64, r int, color string) {
		x, y := latLonToPixels(lat, lon, zoom)
		svg.Circle(int(x), int(y), r, color)
	}

	drawText := func(lat, lon float64, text string, style ...string) {
		x, y := latLonToPixels(lat, lon, zoom)
		y -= 3
		svg.Text(int(x), int(y), text, style...)
	}

	// drawPath draws a bunch of paths with the given style attributes.
	drawPath := func(paths [][]int, useTriangles bool, style ...string) {
		latLon := m.LatLon
		if useTriangles {
			latLon = m.t_latLon
		}
		for _, border := range paths {
			var path [][2]float64
			for _, borderSeg := range border {
				x, y := latLonToPixels(latLon[borderSeg][0], latLon[borderSeg][1], zoom)

				// This check prevents long lines across the SVG if the path happens to wrap around
				// 180° longitude.
				if len(path) >= 1 && dist2(path[len(path)-1], [2]float64{x, y}) > filterPathDist {
					svg.Path(svgGenD(path), style...)
					path = nil
				}
				path = append(path, [2]float64{x, y})
			}
			svg.Path(svgGenD(path), style...)
		}
	}

	if drawLatitudeDots {
		drawCircle(43.0, -80.0, 4, "fill: rgb(123, 255, 23)")
		drawCircle(-43.0, 80.0, 4, "fill: rgb(123, 255, 23)")
		drawCircle(60.0, 0.0, 4, "fill: rgb(123, 255, 223)")
		drawCircle(30.0, 0.0, 4, "fill: rgb(123, 255, 223)")
		drawCircle(0.0, 0.0, 4, "fill: rgb(123, 255, 223)")
		drawCircle(-30.0, 0.0, 4, "fill: rgb(0, 255, 223)")
		drawCircle(-60.0, 0.0, 4, "fill: rgb(0, 255, 223)")
	}

	if drawBorders {
		log.Println("TODO: Place city states first and grow empires from city states?")
		drawPath(m.getCustomBorders(m.RegionToCityState), true, "class=\"cityborder\"")
		drawPath(m.getBorders(), true, "class=\"border\"")
	}

	if drawLakeBorders {
		drawPath(m.getLakeBorders(), true, "class=\"lake\"")
	}

	if drawLandmassContour {
		drawPath(m.getLandmassBorders(), true, "class=\"contour\"")
	}

	// Rivers (based on regions)
	if drawRiversA {
		drawPath(m.getRivers(0.001), false, "class=\"river\"")

		// Skip frozen regions
		// TODO: Fix maxElev caching!!!
		// if m.getRTemperature(rivseg, maxR) < 0 {
		//	continue
		// }
		//
		// Alternative:
		//
		// valMois := em.r_moisture[rivseg] / maxMois
		// if genbiome.GetWhittakerModBiome(int(m.getRTemperature(rivseg, maxR)), int(valMois*45)) == WhittakerModBiomeSnow {
		// 	continue
		// }
	}

	if drawTradeRoutes {
		paths, _ := m.getTradeRoutes()
		drawPath(paths, false, "class=\"traderoute\"")
	}

	// Rivers (based on triangles)
	if drawRiversB {
		for i := 0; i < m.mesh.numSides; i++ {
			if m.s_flow[i] < 10000 {
				continue
			}
			inner_t := m.mesh.s_inner_t(i)
			outer_t := m.mesh.s_outer_t(i)
			if m.t_elevation[inner_t] < 0 && m.t_elevation[outer_t] < 0 {
				continue
			}
			x1, y1 := latLonToPixels(m.t_latLon[inner_t][0], m.t_latLon[inner_t][1], zoom)
			x2, y2 := latLonToPixels(m.t_latLon[outer_t][0], m.t_latLon[outer_t][1], zoom)
			if math.Abs(x1-x2) > float64(size)/2 || math.Abs(y1-y2) > float64(size)/2 {
				continue
			}
			svg.Line(int(x1), int(y1), int(x2), int(y2), "class=\"river\"")
		}
	}

	// Sinks
	if drawSinks {
		for r, rdh := range m.Downhill {
			if rdh < 0 && m.Drainage[r] < 0 && m.Elevation[r] > 0 {
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 2, "fill: rgb(0, 255, 0)")
			}
		}
	}

	if drawWindOrder {
		wind_sort, ord := m.getWindSortOrder()
		minFlux, maxFlux := minMax(wind_sort)
		for _, r := range ord {
			rdh := wind_sort[r]
			col := genGreen((rdh - minFlux) / (maxFlux - minFlux))
			drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
		}
	}

	if drawWindDir {
		windAng := make([]float64, m.mesh.numRegions)
		for i, vec := range m.RegionToWindVec {
			windAng[i] = math.Atan2(vec[0], vec[1])
		}
		minFlux, maxFlux := minMax(windAng)
		for r, rdh := range windAng {
			col := genGreen((rdh - minFlux) / (maxFlux - minFlux))
			drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
		}
	}

	if drawPlateCompression {
		mountain_r, coastline_r, ocean_r, compression_r := m.findCollisions()
		var minComp, maxComp float64
		for _, comp := range compression_r {
			if comp < minComp {
				minComp = comp
			}
			if comp > maxComp {
				maxComp = comp
			}
		}
		for _, r := range mountain_r {
			drawCircle(m.LatLon[r][0], m.LatLon[r][1], 2, "fill: rgb(255, 128, 128)")
		}
		for _, r := range coastline_r {
			drawCircle(m.LatLon[r][0], m.LatLon[r][1], 2, "fill: rgb(128, 255, 128)")
		}
		for _, r := range ocean_r {
			drawCircle(m.LatLon[r][0], m.LatLon[r][1], 2, "fill: rgb(128, 128, 255)")
		}
		for r := 0; r < m.mesh.numSides; r++ {
			if compression_r[r] != 0 {
				col := genGreen((compression_r[r] - minComp) / (maxComp - minComp))
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawFlux {
		minFlux, maxFlux := minMax(m.Flux)
		for r, rdh := range m.Flux {
			if rdh > 0 {
				col := genGreen((rdh - minFlux) / (maxFlux - minFlux))
				col = genGreen(rdh / maxFlux)
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawHumidity {
		minHumid, maxHumid := minMax(m.Moisture)
		for r, rdh := range m.Moisture {
			if rdh > 0 {
				col := genGreen((rdh - minHumid) / (maxHumid - minHumid))
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawRainfall {
		minRain, maxRain := minMax(m.Rainfall)
		for r, rdh := range m.Rainfall {
			if rdh > 0 {
				col := genGreen((rdh - minRain) / (maxRain - minRain))
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawErosion {
		er := m.GetErosionRate()
		minFlux, maxFlux := minMax(er)
		for r, rdh := range m.Flux {
			if rdh > 0 {
				col := genBlue((rdh - minFlux) / (maxFlux - minFlux))
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawErosion2 {
		er := m.GetErosionRate2()
		minFlux, maxFlux := minMax(er)
		for r, rdh := range m.Flux {
			if rdh > 0 {
				col := genBlue((rdh - minFlux) / (maxFlux - minFlux))
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawAltitude {
		er := m.Elevation
		minHeight, maxHeight := minMax(er)
		minHeight = 0
		for r, rdh := range m.Elevation {
			if rdh > 0 && r%2 == 0 {
				col := genBlue((rdh - minHeight) / (maxHeight - minHeight))
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawTemperature {
		er := m.Elevation
		_, maxHeight := minMax(er)
		for r, rdh := range m.Elevation {
			if rdh > 0 && r%2 == 0 {
				t := m.getRTemperature(r, maxHeight)
				col := genBlue((t - minTemp) / (maxTemp - minTemp))
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawBelow {
		for r, pVal := range m.Elevation {
			if pVal <= 0 {
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 2, "fill: rgb(0, 0, 255)")
			}
		}
	}

	// Water pools
	if drawPools {
		for r, pVal := range m.Waterpool {
			if pVal > 0 {
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 2, "fill: rgb(0, 0, 255)")
			}
		}
	}

	if drawDrains {
		drains := make(map[int]bool)
		for r, drain := range m.Drainage {
			if drain >= 0 {
				drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, "fill: rgb(255, 0, 255)")
			}
			if drain != -1 {
				drains[drain] = true
			}
		}
		for r := range drains {
			drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, "fill: rgb(255, 255, 0)")
		}
	}

	if drawResources {
		grad := colorgrad.Rainbow()
		cols := grad.Colors(uint(ResMaxMetals))

		// NOTE: This sucks right now.
		res := m.Metals
		radius := 1
		count := make([]int, ResMaxMetals)
		for i := 0; i < ResMaxMetals; i++ {
			cr, cg, cb, _ := cols[i].RGBA()
			col := fmt.Sprintf("fill: rgb(%d, %d, %d)", cr/(0xffff/255), cg/(0xffff/255), cb/(0xffff/255))
			for r, t := range res {
				if t&(1<<i) > 0 {
					count[i]++
					drawCircle(m.LatLon[r][0], m.LatLon[r][1], radius, col)
				}
			}
		}
		for i := 0; i < ResMaxMetals; i++ {
			log.Printf("Metal %s: %d", metalToString(i), count[i])
		}
	}

	// Cities
	if drawCities {
		for i, r := range m.Cities {
			radius := 2
			class := "class=\"city\""
			col := "fill: rgb(255, 165, 0)"

			// Capital cities are bigger!
			if i < m.NumEmpires {
				radius = 4
				class = "class=\"capital\""
				col = "fill: rgb(255, 0, 0)"
			}
			switch r.Type {
			case TownTypeDefault:
			case TownTypeMining:
				col = "fill: rgb(255, 255, 0)"
				radius = 2
			case TownTypeFarming:
				col = "fill: rgb(55, 255, 0)"
				radius = 1
			case TownTypeDesertOasis:
				col = "fill: rgb(55, 0, 255)"
				radius = 1
			}
			drawCircle(m.LatLon[r.ID][0], m.LatLon[r.ID][1], radius, col)
			drawText(m.LatLon[r.ID][0], m.LatLon[r.ID][1], r.Name, class)
		}
		// TODO: Move labels to avoid overlap.
	}

	if drawCityscore {
		scores := m.CalcCityScore(m.getFitnessCityDefault(), func() []int { return nil })
		minScore, maxScore := minMax(scores)
		for r, score := range scores {
			col := genBlue((score - minScore) / (maxScore - minScore))
			drawCircle(m.LatLon[r][0], m.LatLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
		}
	}

	svg.End()
	return nil
}

func svgGenD(path [][2]float64) string {
	var str string
	for i, p := range path {
		if i == 0 {
			str += fmt.Sprintf("M %f,%f", p[0], p[1])
			continue
		}
		str += fmt.Sprintf("L %f,%f", p[0], p[1])
	}
	return str
}

func (m *Map) ExportWebp(name string) {
	zoom := 1
	size := sizeFromZoom(zoom)
	webpanim := webpanimation.NewWebpAnimation(size, size, 0)
	webpanim.WebPAnimEncoderOptions.SetKmin(9)
	webpanim.WebPAnimEncoderOptions.SetKmax(17)
	defer webpanim.ReleaseMemory() // don't forget call this or you will have memory leaks
	webpConfig := webpanimation.NewWebpConfig()
	webpConfig.SetLossless(1)

	timeline := 0
	timestep := 50
	for i := 0; i < 366; i++ {
		m.Step()
		if err := webpanim.AddFrame(m.getImage(), timeline, webpConfig); err != nil {
			log.Fatal(err)
		}
		timeline += timestep
	}

	if err := webpanim.AddFrame(nil, timeline, webpConfig); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	// encode animation and write result bytes in buffer
	if err = webpanim.Encode(f); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (m *Map) getImage() image.Image {
	grad := colorgrad.Rainbow()

	terrToCol := make(map[int]int)

	terr := m.Cities[:m.NumCityStates]
	territory := m.RegionToCityState
	//terr := m.Cultures
	//territory := m.RegionToCulture
	for i, c := range terr {
		terrToCol[c.ID] = i
		//log.Printf("%d: %s %f", i, c.Type, c.Expansionism)
	}
	cols := grad.Colors(uint(len(terr)))

	zoom := 1
	size := sizeFromZoom(zoom)
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	min, max := minMax(m.Elevation)
	_, maxMois := minMax(m.Rainfall)
	for r := 0; r < m.mesh.numRegions; r++ {
		lat := m.LatLon[r][0]
		lon := m.LatLon[r][1]
		// log.Println(lat, lon)
		x, y := latLonToPixels(lat, lon, zoom)
		val := (m.Elevation[r] - min) / (max - min)
		var col color.NRGBA
		if elev := m.Elevation[r]; elev <= 0 || m.Waterpool[r] > 0 || m.Flux[r] > 1000 {
			col = genBlue(val)
		} else {
			valElev := elev / max
			// Hacky: Modify elevation based on latitude to compensate for colder weather at the poles and warmer weather at the equator.
			// valElev := math.Max(math.Min((elev/max)+(math.Sqrt(math.Abs(lat)/90.0)-0.5), max), 0)
			valMois := m.Rainfall[r] / maxMois
			if territory[r] != 0 && false {
				cr, cg, cb, _ := cols[terrToCol[territory[r]]].RGBA()
				col.R = uint8(float64(255) * float64(cr) / float64(0xffff))
				col.G = uint8(float64(255) * float64(cg) / float64(0xffff))
				col.B = uint8(float64(255) * float64(cb) / float64(0xffff))
				col.A = 255
			} else {
				temMin, temMax := m.GetMinMaxTemperature(lat)
				temAvg := (temMin + temMax) / 2
				col = genbiome.GetWhittakerModBiomeColor(int(temAvg-getTempFalloffFromAltitude(maxAltitudeFactor*valElev)), int(valMois*maxPrecipitation), val)
			}
			// col = GetWhittakerModBiomeColor(int(getMeanAnnualTemp(lat)-getTempFalloffFromAltitude(8850*valElev)), int(valMois*45), val)
		}
		img.Set(int(x), int(y), col)
	}
	return img
}

func (m *Map) ExportPng(name string) {
	img := m.getImage()

	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (m *Map) ExportOBJ(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	drawPlates := false
	drawRivers := false
	/*
		// This will export the quad geometry.
		// Vertices
		for i := 0; i < len(m.QuadGeom.xyz); i += 3 {
			ve := convToVec3(m.QuadGeom.xyz[i:]).Mul(1.0 + 0.01*m.QuadGeom.tm[(i/3)*2])
			w.WriteString(fmt.Sprintf("v %f %f %f \n", ve.X, ve.Y, ve.Z))
		}

		// Globe
		for i := 0; i < len(m.QuadGeom.I); i += 3 {
			w.WriteString(fmt.Sprintf("f %d %d %d \n", m.QuadGeom.I[i]+1, m.QuadGeom.I[i+1]+1, m.QuadGeom.I[i+2]+1))
			w.Flush()
		}
		w.Flush()
	*/

	// Vertices
	for i := 0; i < len(m.XYZ); i += 3 {
		ve := convToVec3(m.XYZ[i:]).Mul(1.0 + 0.01*(m.Elevation[i/3]+m.Waterpool[i/3]))
		w.WriteString(fmt.Sprintf("v %f %f %f \n", ve.X, ve.Y, ve.Z))
	}

	// Triangle vertices
	if drawPlates || drawRivers {
		for i := 0; i < len(m.t_xyz); i += 3 {
			ve := convToVec3(m.t_xyz[i:]).Mul(1.03 + 0.01*m.t_elevation[i/3])
			w.WriteString(fmt.Sprintf("v %f %f %f \n", ve.X, ve.Y, ve.Z))
		}
		w.Flush()
	}

	// Globe
	for i := 0; i < len(m.mesh.Triangles); i += 3 {
		w.WriteString(fmt.Sprintf("f %d %d %d \n", m.mesh.Triangles[i]+1, m.mesh.Triangles[i+1]+1, m.mesh.Triangles[i+2]+1))
		w.Flush()
	}
	w.Flush()

	// Rivers
	if drawRivers {
		for i := 0; i < m.mesh.numSides; i++ {
			if m.s_flow[i] > 1 {
				inner_t := m.mesh.s_inner_t(i)
				outer_t := m.mesh.s_outer_t(i)
				if m.t_elevation[inner_t] < 0 && m.t_elevation[outer_t] < 0 {
					continue
				}
				w.WriteString(fmt.Sprintf("l %d %d \n", (len(m.XYZ)/3)+inner_t+1, (len(m.XYZ)/3)+outer_t+1))
				w.Flush()
			}
		}
		w.Flush()
	}

	// Plates
	if drawPlates {
		for s := 0; s < m.mesh.numSides; s++ {
			begin_r := m.mesh.s_begin_r(s)
			end_r := m.mesh.s_end_r(s)
			if m.RegionToPlate[begin_r] != m.RegionToPlate[end_r] {
				inner_t := m.mesh.s_inner_t(s)
				outer_t := m.mesh.s_outer_t(s)
				w.WriteString(fmt.Sprintf("l %d %d \n", (len(m.XYZ)/3)+inner_t+1, (len(m.XYZ)/3)+outer_t+1))
			}
			w.Flush()
		}
		w.Flush()
	}
	return nil
}
