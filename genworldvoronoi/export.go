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
)

// sizeFromZoom returns the expected size of the world for the mercato projection used below.
func sizeFromZoom(zoom int) int {
	return int(math.Pow(2.0, float64(zoom)) * 256.0)
}

func genBlue(intensity float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(intensity * 255),
		G: uint8(intensity * 255),
		B: 255,
		A: 255,
	}
}

// ExportSVG exports the terrain as SVG to the given path.
// NOTE: This produces broken somewhat incomplete output due to the wraparound of the mesh.
func (m *Map) ExportSVG(path string) error {
	zoom := 2
	filterPathDist := 20.0
	size := sizeFromZoom(zoom)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	svg := svgo.New(f)
	svg.Start(size, size)

	min, max := minMax(m.t_elevation)
	minMois, maxMois := minMax(m.t_moisture)
	for i := 0; i < len(m.mesh.Triangles); i += 3 {
		tmpLine := ""

		// Hacky way to filter paths/triangles that wrap around the entire SVG.
		triLat, triLon := latLonFromVec3(convToVec3(m.t_xyz[i:i+3]).Normalize(), 1.0)
		triX, triY := mercator.LatLonToPixels(triLat, -1*(180-triLon), zoom)
		var skip bool
		for _, j := range m.mesh.Triangles[i : i+3] {
			lat, lon := latLonFromVec3(convToVec3(m.r_xyz[j*3:(j*3)+3]).Normalize(), 1.0)
			x, y := mercator.LatLonToPixels(lat, -1*(180-lon), zoom)
			if dist2([2]float64{x + float64(size)/2, y}, [2]float64{triX + float64(size)/2, triY}) > filterPathDist {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		var path [][2]float64
		for _, j := range m.mesh.Triangles[i : i+3] {
			lat, lon := latLonFromVec3(convToVec3(m.r_xyz[j*3:(j*3)+3]).Normalize(), 1.0)
			x, y := mercator.LatLonToPixels(lat, -1*(180-lon), zoom)
			path = append(path, [2]float64{x + float64(size)/2, y})
		}
		elev := m.t_elevation[i/3]
		val := (elev - min) / (max - min)
		var col color.NRGBA
		if elev < 0 {
			col = genBlue(val)
		} else {
			valElev := elev / max
			// Hacky: Modify elevation based on latitude to compensate for colder weather at the poles and warmer weather at the equator.
			// valElev := math.Max(math.Min((elev/max)+(math.Sqrt(math.Abs(triLat)/90.0)-0.5), max), 0)
			valMois := (m.t_moisture[i/3] - minMois) / (maxMois - minMois)
			valMois = m.t_moisture[i/3] / maxMois
			// col = GetRedblobBiomeColor(int(valElev*4)+1, int(valMois*6)+1, val)
			col = GetWhittakerModBiomeColor(int(getMeanAnnualTemp(triLat)-getTempFalloffFromAltitude(8850*valElev)), int(valMois*45), val)
		}

		svg.Path(svgGenD(path), fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.B)+tmpLine)
	}

	// Rivers
	drawRivers := false
	if drawRivers {
		for i := 0; i < m.mesh.numSides; i++ {
			if m.s_flow[i] > 1 {
				inner_t := m.mesh.s_inner_t(i)
				outer_t := m.mesh.s_outer_t(i)
				if m.t_elevation[inner_t] < 0 && m.t_elevation[outer_t] < 0 {
					continue
				}
				lat, lon := latLonFromVec3(convToVec3(m.t_xyz[inner_t*3:(inner_t*3)+3]).Normalize(), 1.0)
				x1, y1 := mercator.LatLonToPixels(lat, -1*(180-lon), zoom)
				lat, lon = latLonFromVec3(convToVec3(m.t_xyz[outer_t*3:(outer_t*3)+3]).Normalize(), 1.0)
				x2, y2 := mercator.LatLonToPixels(lat, -1*(180-lon), zoom)
				if math.Abs(x1-x2) > float64(size)/2 || math.Abs(y1-y2) > float64(size)/2 {
					continue
				}
				svg.Line(int(x1+float64(size)/2), int(y1), int(x2+float64(size)/2), int(y2), "stroke=\"blue\" stroke-width=\"1\"")
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
	zoom := 1
	size := sizeFromZoom(zoom)
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	min, max := minMax(m.r_elevation)
	minMois, maxMois := minMax(m.r_rainfall)
	for r := 0; r < m.mesh.numRegions; r++ {
		lat, lon := latLonFromVec3(convToVec3(m.r_xyz[r*3:(r*3)+3]).Normalize(), 1.0)
		//log.Println(lat, lon)
		x, y := mercator.LatLonToPixels(lat, -1*(180-lon), zoom)
		val := (m.r_elevation[r] - min) / (max - min)
		var col color.NRGBA
		if elev := m.r_elevation[r]; elev < 0 {
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
			//col = GetWhittakerModBiomeColor(int(getMeanAnnualTemp(lat)-getTempFalloffFromAltitude(8850*valElev)), int(valMois*45), val)
		}
		img.Set(int(x)+size/2, int(y), col)
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
		ve := convToVec3(m.r_xyz[i:]).Mul(1.0 + 0.01*m.r_elevation[i/3])
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
