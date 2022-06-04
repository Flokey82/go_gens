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

	svgo "github.com/ajstarks/svgo"
	"github.com/davvo/mercator"
	"github.com/mazznoer/colorgrad"
)

const tileSize = 256

// sizeFromZoom returns the expected size of the world for the mercato projection used below.
func sizeFromZoom(zoom int) int {
	return int(math.Pow(2.0, float64(zoom)) * float64(tileSize))
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
	log.Println(lat1, lon1, lat2, lon2)
	/*if lat1 < 0 {
		lat1 += 180 // This sucks.
	}
	if lat2 < 0 {
		lat2 += 180 // This sucks.
	}
	if lon1 < 0 {
		lon1 += 360 // This sucks.
	}
	if lon2 < 0 {
		lon2 += 360 // This sucks.
	}*/
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
		l0 := ll[0]
		//if l0 < 0 {
		//	l0 += 180 // This sucks.
		//}
		if l0 < lat1 || l0 >= lat2 {
			continue
		}
		l1 := ll[1]
		//if l1 < 0 {
		//	l1 += 360 // This sucks.
		//}
		if l1 < lon1 || l1 >= lon2 {
			continue
		}
		r.r = append(r.r, i)
	}
	for i, ll := range m.t_latLon {
		l0 := ll[0]
		//if l0 < 0 {
		//	l0 += 180 // This sucks.
		//}
		if l0 < lat1 || l0 >= lat2 {
			continue
		}
		l1 := ll[1]
		//if l1 < 0 {
		//	l1 += 360 // This sucks.
		//}
		if l1 < lon1 || l1 >= lon2 {
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
	drawHumidity := false
	drawRainfall := false
	drawBorders := true
	drawLakeBorders := true
	drawBelow := false

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

	em := m
	// Hack to test tile fetching
	// 113.48673955688815 180 139.02010193037987 225
	// 139.02010193037987 180 0 225
	//tbb := tileBoundingBox(0, 0, 0)
	//la1, lo1, la2, lo2 := tbb.ToLatLon()
	//re := m.getBB(la1, lo1, la2, lo2)
	//em, err := m.interpolate(re.r)
	//if err != nil {
	//	panic(err)
	//}
	// end hack
	min, max := minMax(m.t_elevation)
	minMois, maxMois := minMax(m.t_moisture)
	for i := 0; i < len(em.mesh.Triangles); i += 3 {
		tmpLine := ""

		// Hacky way to filter paths/triangles that wrap around the entire SVG.
		triLat := em.t_latLon[i/3][0]
		triLon := em.t_latLon[i/3][1]
		triX, triY := mercator.LatLonToPixels(triLat, triLon, zoom)
		var skip bool
		var poolCount int
		for _, j := range em.mesh.Triangles[i : i+3] {
			if em.r_pool[j] > 0 {
				poolCount++
			}
			x, y := mercator.LatLonToPixels(em.r_latLon[j][0], em.r_latLon[j][1], zoom)
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
			x, y := mercator.LatLonToPixels(em.r_latLon[j][0], em.r_latLon[j][1], zoom)
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
			valMois := (em.t_moisture[i/3] - minMois) / (maxMois - minMois)
			valMois = em.t_moisture[i/3] / maxMois
			// col = GetRedblobBiomeColor(int(valElev*4)+1, int(valMois*6)+1, val)
			col = GetWhittakerModBiomeColor(int(getMeanAnnualTemp(triLat)-getTempFalloffFromAltitude(8850*valElev)), int(valMois*45), val)
		}

		svg.Path(svgGenD(path), fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.B)+tmpLine)
	}

	if drawBorders {
		for _, border := range m.getBorders() {
			var path [][2]float64
			for _, borderSeg := range border {
				x, y := mercator.LatLonToPixels(m.t_latLon[borderSeg][0], m.t_latLon[borderSeg][1], zoom)
				path = append(path, [2]float64{x, y})
			}
			svg.Path(svgGenD(path), "stroke=\"red\" fill=\"none\" stroke-width=\"0.5\"")
		}
	}

	if drawLakeBorders {
		for _, border := range m.getLakeBorders() {
			var path [][2]float64
			for _, borderSeg := range border {
				x, y := mercator.LatLonToPixels(m.t_latLon[borderSeg][0], m.t_latLon[borderSeg][1], zoom)
				path = append(path, [2]float64{x, y})
			}
			svg.Path(svgGenD(path), "stroke=\"blue\"", "fill=\"blue\"", "fill-opacity=\"0.5\"", "stroke-width=\"0.5\"")
		}
	}

	// Rivers
	if drawRiversA {
		for _, riv := range m.getRivers(0.005) {
			var path [][2]float64
			for _, rivseg := range riv {
				x, y := mercator.LatLonToPixels(m.r_latLon[rivseg][0], m.r_latLon[rivseg][1], zoom)
				path = append(path, [2]float64{x, y})
			}
			svg.Path(svgGenD(path), "stroke=\"blue\" fill=\"none\" stroke-width=\"0.5\"")
		}
	}

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
			x1, y1 := mercator.LatLonToPixels(m.t_latLon[inner_t][0], m.t_latLon[inner_t][1], zoom)
			x2, y2 := mercator.LatLonToPixels(m.t_latLon[outer_t][0], m.t_latLon[outer_t][1], zoom)
			if math.Abs(x1-x2) > float64(size)/2 || math.Abs(y1-y2) > float64(size)/2 {
				continue
			}
			svg.Line(int(x1), int(y1), int(x2), int(y2), "stroke=\"blue\" stroke-width=\"1\"")
		}
	}

	// Sinks
	if drawSinks {
		for r, rdh := range m.r_downhill {
			if rdh < 0 && m.r_drainage[r] < 0 && m.r_elevation[r] > 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				r := 2
				svg.Circle(int(x), int(y), r, "fill: rgb(0, 255, 0)")

			}
		}
	}

	if drawFlux {
		minFlux, maxFlux := minMax(m.r_flux)
		for r, rdh := range m.r_flux {
			if rdh > 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				r := 1
				col := genGreen((rdh - minFlux) / (maxFlux - minFlux))
				col = genGreen(rdh / maxFlux)
				svg.Circle(int(x), int(y), r, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))
			}
		}
	}

	if drawHumidity {
		minHumid, maxHumid := minMax(m.r_moisture)
		for r, rdh := range m.r_moisture {
			if rdh > 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				r := 1
				col := genGreen((rdh - minHumid) / (maxHumid - minHumid))
				svg.Circle(int(x), int(y), r, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))

			}
		}
	}

	if drawRainfall {
		minRain, maxRain := minMax(m.r_rainfall)
		for r, rdh := range m.r_rainfall {
			if rdh > 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				r := 1
				col := genGreen((rdh - minRain) / (maxRain - minRain))
				svg.Circle(int(x), int(y), r, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.R, col.R))

			}
		}
	}

	if drawErosion {
		er := m.rErosionRate()
		minFlux, maxFlux := minMax(er)
		for r, rdh := range m.r_flux {
			if rdh > 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				r := 1
				col := genBlue((rdh - minFlux) / (maxFlux - minFlux))
				svg.Circle(int(x), int(y), r, fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.G))

			}
		}
	}
	if drawBelow {
		for r, pVal := range m.r_elevation {
			if pVal <= 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				svg.Circle(int(x), int(y), 2, "fill: rgb(0, 0, 255)")
			}
		}
	}

	// Water pools
	if drawPools {
		for r, pVal := range m.r_pool {
			if pVal > 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				svg.Circle(int(x), int(y), 2, "fill: rgb(0, 0, 255)")
			}
		}
	}
	if drawDrains {
		drains := make(map[int]bool)
		for r, drain := range m.r_drainage {
			if drain >= 0 {
				x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
				r := 1
				svg.Circle(int(x), int(y), r, "fill: rgb(255, 0, 255)")
			}
			if drain != -1 {
				drains[drain] = true
			}
		}
		for r := range drains {
			x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
			r := 1
			svg.Circle(int(x), int(y), r, "fill: rgb(255, 255, 0)")
		}
	}

	// Cities
	if drawCities {
		for i, r := range m.cities_r {
			x, y := mercator.LatLonToPixels(m.r_latLon[r][0], m.r_latLon[r][1], zoom)
			r := 2
			// Capital cities are bigger!
			if i < m.NumTerritories {
				r = 4
			}
			svg.Circle(int(x), int(y), r, "fill: rgb(255, 0, 0)")
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
		terrToCol[terr] = i
	}
	zoom := 1
	size := sizeFromZoom(zoom)
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	min, max := minMax(m.r_elevation)
	minMois, maxMois := minMax(m.r_rainfall)
	for r := 0; r < m.mesh.numRegions; r++ {
		lat := m.r_latLon[r][0]
		lon := m.r_latLon[r][1]
		//log.Println(lat, lon)
		x, y := mercator.LatLonToPixels(lat, lon, zoom)
		val := (m.r_elevation[r] - min) / (max - min)
		var col color.NRGBA
		if elev := m.r_elevation[r]; elev <= 0 || m.r_pool[r] > 0 || m.r_flux[r] > 1000 {
			col = genBlue(val)
		} else {
			valElev := elev / max
			// Hacky: Modify elevation based on latitude to compensate for colder weather at the poles and warmer weather at the equator.
			// valElev := math.Max(math.Min((elev/max)+(math.Sqrt(math.Abs(lat)/90.0)-0.5), max), 0)
			valMois := (m.r_rainfall[r] - minMois) / (maxMois - minMois)
			valMois = m.r_rainfall[r] / maxMois
			col = GetRedblobBiomeColor(int(valElev*4)+1, int(valMois*6)+1, val)
			col.R = uint8(255 * valMois)
			col.G = 0
			col.B = uint8(255 * (1 - valMois))
			if m.r_territory[r] == 0 {
				col = GetWhittakerModBiomeColor(int(getMeanAnnualTemp(lat)-getTempFalloffFromAltitude(8850*valElev)), int(valMois*45), val)
			} else {
				cr, cg, cb, _ := cols[terrToCol[m.r_territory[r]]].RGBA()
				col.R = uint8(float64(255) * float64(cr) / float64(0xffff))
				col.G = uint8(float64(255) * float64(cg) / float64(0xffff))
				col.B = uint8(float64(255) * float64(cb) / float64(0xffff))
			}
			//col = GetWhittakerModBiomeColor(int(getMeanAnnualTemp(lat)-getTempFalloffFromAltitude(8850*valElev)), int(valMois*45), val)
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
