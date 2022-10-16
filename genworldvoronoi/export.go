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
	for i, ll := range m.r_latLon {
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
	drawBorders := false
	drawLakeBorders := false
	drawBelow := false
	drawContour := true
	drawWindDir := false
	drawPlateCompression := false
	drawAltitude := false
	drawTemperature := false
	drawLatitudeDots := false
	drawCityscore := false
	drawRegionTerrain := true
	drawTradeRoutes := true
	drawResources := false

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
			"path.border{ stroke: red;}\n"+
			"path.terrain{ stroke: none;}\n"+
			"path.traderoute{ stroke: lime; stroke-width: 0.5;}\n")

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
		//fitScore := m.resourceFitness()
		//_, maxFit := minMax(fitScore)
		min, max := minMax(m.r_elevation)
		_, maxMois := minMax(m.r_moisture)
		for i := 0; i < em.mesh.numRegions; i++ {
			rLat := em.r_latLon[i][0]
			rLon := em.r_latLon[i][1]
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
			elev := em.r_elevation[i]
			val := (elev - min) / (max - min)
			//val = cityScore[i] / maxS
			//val = fitScore[i] / maxFit
			var col color.NRGBA
			if elev <= 0 {
				col = genBlue(val)
			} else {
				valElev := elev / max
				valMois := em.r_moisture[i] / maxMois
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
				if em.r_pool[j] > 0 {
					poolCount++
				}
				x, y := latLonToPixels(em.r_latLon[j][0], em.r_latLon[j][1], zoom)
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
				x, y := latLonToPixels(em.r_latLon[j][0], em.r_latLon[j][1], zoom)
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

	// drawPath draws a bunch of paths with the given style attributes.
	drawPath := func(paths [][]int, useTriangles bool, style ...string) {
		latLon := m.r_latLon
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
		drawPath(m.getBorders(), true, "class=\"border\"")
	}

	if drawLakeBorders {
		drawPath(m.getLakeBorders(), true, "class=\"lake\"")
	}

	if drawContour {
		drawPath(m.contour(), true, "class=\"contour\"")
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
		for r, rdh := range m.r_downhill {
			if rdh < 0 && m.r_drainage[r] < 0 && m.r_elevation[r] > 0 {
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 2, "fill: rgb(0, 255, 0)")
			}
		}
	}

	if drawWindOrder {
		wind_sort, ord := m.getWindSortOrder()
		minFlux, maxFlux := minMax(wind_sort)
		for _, r := range ord {
			rdh := wind_sort[r]
			col := genGreen((rdh - minFlux) / (maxFlux - minFlux))
			drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
		}
	}

	if drawWindDir {
		windAng := make([]float64, m.mesh.numRegions)
		for i, vec := range m.r_windvec {
			windAng[i] = math.Atan2(vec[0], vec[1])
		}
		minFlux, maxFlux := minMax(windAng)
		for r, rdh := range windAng {
			col := genGreen((rdh - minFlux) / (maxFlux - minFlux))
			drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
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
			drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 2, "fill: rgb(255, 128, 128)")
		}
		for _, r := range coastline_r {
			drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 2, "fill: rgb(128, 255, 128)")
		}
		for _, r := range ocean_r {
			drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 2, "fill: rgb(128, 128, 255)")
		}
		for r := 0; r < m.mesh.numSides; r++ {
			if compression_r[r] != 0 {
				col := genGreen((compression_r[r] - minComp) / (maxComp - minComp))
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawFlux {
		minFlux, maxFlux := minMax(m.r_flux)
		for r, rdh := range m.r_flux {
			if rdh > 0 {
				col := genGreen((rdh - minFlux) / (maxFlux - minFlux))
				col = genGreen(rdh / maxFlux)
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawHumidity {
		minHumid, maxHumid := minMax(m.r_moisture)
		for r, rdh := range m.r_moisture {
			if rdh > 0 {
				col := genGreen((rdh - minHumid) / (maxHumid - minHumid))
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawRainfall {
		minRain, maxRain := minMax(m.r_rainfall)
		for r, rdh := range m.r_rainfall {
			if rdh > 0 {
				col := genGreen((rdh - minRain) / (maxRain - minRain))
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawErosion {
		er := m.rErosionRate()
		minFlux, maxFlux := minMax(er)
		for r, rdh := range m.r_flux {
			if rdh > 0 {
				col := genBlue((rdh - minFlux) / (maxFlux - minFlux))
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawErosion2 {
		er := m.getRErosion2()
		minFlux, maxFlux := minMax(er)
		for r, rdh := range m.r_flux {
			if rdh > 0 {
				col := genBlue((rdh - minFlux) / (maxFlux - minFlux))
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawAltitude {
		er := m.r_elevation
		minHeight, maxHeight := minMax(er)
		minHeight = 0
		for r, rdh := range m.r_elevation {
			if rdh > 0 && r%2 == 0 {
				col := genBlue((rdh - minHeight) / (maxHeight - minHeight))
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawTemperature {
		er := m.r_elevation
		_, maxHeight := minMax(er)
		for r, rdh := range m.r_elevation {
			if rdh > 0 && r%2 == 0 {
				t := m.getRTemperature(r, maxHeight)
				col := genBlue((t - minTemp) / (maxTemp - minTemp))
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
			}
		}
	}

	if drawBelow {
		for r, pVal := range m.r_elevation {
			if pVal <= 0 {
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 2, "fill: rgb(0, 0, 255)")
			}
		}
	}

	// Water pools
	if drawPools {
		for r, pVal := range m.r_pool {
			if pVal > 0 {
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 2, "fill: rgb(0, 0, 255)")
			}
		}
	}

	if drawDrains {
		drains := make(map[int]bool)
		for r, drain := range m.r_drainage {
			if drain >= 0 {
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, "fill: rgb(255, 0, 255)")
			}
			if drain != -1 {
				drains[drain] = true
			}
		}
		for r := range drains {
			drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, "fill: rgb(255, 255, 0)")
		}
	}

	// Cities
	if drawCities {
		for i, r := range m.cities_r {
			radius := 3
			// Capital cities are bigger!
			if i < m.NumTerritories {
				radius = 4
			}
			col := "fill: rgb(255, 0, 0)"
			switch r.Type {
			case TownTypeDefault:
			case TownTypeMining:
				col = "fill: rgb(255, 255, 0)"
				radius = 2
			case TownTypeFarming:
				col = "fill: rgb(55, 255, 0)"
				radius = 1
			}
			drawCircle(m.r_latLon[r.R][0], m.r_latLon[r.R][1], radius, col)
		}
	}

	if drawCityscore {
		scores := m.rCityScore(m.getFitnessCityDefault())
		minScore, maxScore := minMax(scores)
		for r, score := range scores {
			col := genBlue((score - minScore) / (maxScore - minScore))
			drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], 1, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))
		}
	}

	if drawResources {
		// NOTE: This sucks right now.
		res := m.res_metals_r
		for r, t := range res {
			radius := 1
			// Capital cities are bigger!
			if t&ResMetPlatinum > 0 {
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], radius, "fill: rgb(255, 127, 255)")
			}
		}
		for r, t := range res {
			radius := 1
			// Capital cities are bigger!
			if t&ResMetGold > 0 {
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], radius, "fill: rgb(255, 255, 0)")
			}
		}
		for r, t := range res {
			radius := 1
			// Capital cities are bigger!
			if t&ResMetIron > 0 {
				drawCircle(m.r_latLon[r][0], m.r_latLon[r][1], radius, "fill: rgb(200, 127, 0)")
			}
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

func (m *Map) ExportPng(name string) {
	grad := colorgrad.Rainbow()
	cols := grad.Colors(uint(m.NumTerritories))
	terrToCol := make(map[int]int)
	for i, terr := range m.cities_r[:m.NumTerritories] {
		terrToCol[terr.R] = i
	}
	zoom := 1
	size := sizeFromZoom(zoom)
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	min, max := minMax(m.r_elevation)
	_, maxMois := minMax(m.r_rainfall)
	for r := 0; r < m.mesh.numRegions; r++ {
		lat := m.r_latLon[r][0]
		lon := m.r_latLon[r][1]
		// log.Println(lat, lon)
		x, y := latLonToPixels(lat, lon, zoom)
		val := (m.r_elevation[r] - min) / (max - min)
		var col color.NRGBA
		if elev := m.r_elevation[r]; elev <= 0 || m.r_pool[r] > 0 || m.r_flux[r] > 1000 {
			col = genBlue(val)
		} else {
			valElev := elev / max
			// Hacky: Modify elevation based on latitude to compensate for colder weather at the poles and warmer weather at the equator.
			// valElev := math.Max(math.Min((elev/max)+(math.Sqrt(math.Abs(lat)/90.0)-0.5), max), 0)
			valMois := m.r_rainfall[r] / maxMois
			if m.r_territory[r] == 0 {
				col = genbiome.GetWhittakerModBiomeColor(int(getMeanAnnualTemp(lat)-getTempFalloffFromAltitude(maxAltitudeFactor*valElev)), int(valMois*maxPrecipitation), val)
			} else {
				cr, cg, cb, _ := cols[terrToCol[m.r_territory[r]]].RGBA()
				col.R = uint8(float64(255) * float64(cr) / float64(0xffff))
				col.G = uint8(float64(255) * float64(cg) / float64(0xffff))
				col.B = uint8(float64(255) * float64(cb) / float64(0xffff))
				col.A = 255
			}
			// col = GetWhittakerModBiomeColor(int(getMeanAnnualTemp(lat)-getTempFalloffFromAltitude(8850*valElev)), int(valMois*45), val)
		}
		img.Set(int(x), int(y), col)
	}

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
	for i := 0; i < len(m.r_xyz); i += 3 {
		ve := convToVec3(m.r_xyz[i:]).Mul(1.0 + 0.01*(m.r_elevation[i/3]+m.r_pool[i/3]))
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
				w.WriteString(fmt.Sprintf("l %d %d \n", (len(m.r_xyz)/3)+inner_t+1, (len(m.r_xyz)/3)+outer_t+1))
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
			if m.r_plate[begin_r] != m.r_plate[end_r] {
				inner_t := m.mesh.s_inner_t(s)
				outer_t := m.mesh.s_outer_t(s)
				w.WriteString(fmt.Sprintf("l %d %d \n", (len(m.r_xyz)/3)+inner_t+1, (len(m.r_xyz)/3)+outer_t+1))
			}
			w.Flush()
		}
		w.Flush()
	}
	return nil
}
