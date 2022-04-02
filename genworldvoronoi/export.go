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

// dist2 returns the eucledian distance between two points.
func dist2(a, b [2]float64) float64 {
	xDiff := a[0] - b[0]
	yDiff := a[1] - b[1]
	return float64(math.Sqrt(float64(xDiff*xDiff + yDiff*yDiff)))
}

// ExportSVG exports the terrain as SVG to the given path.
// NOTE: This produces broken somewhat incomplete output due to the wraparound of the mesh.
func (m *Map) ExportSVG(path string) error {
	zoom := 2
	filterPathDist := 20.0
	size := int(math.Pow(2.0, float64(zoom)) * 256.0)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	svg := svgo.New(f)
	svg.Start(size, size)

	min, max := MinMax(m.t_elevation)
	minMois, maxMois := MinMax(m.t_moisture)
	for i := 0; i < len(m.mesh.Triangles); i += 3 {
		tmpLine := ""

		// Hacky way to filter paths/triangles that wrap around the entire SVG.
		triLat, triLon := latLonFromVec3(convToVec3(m.t_xyz[i:i+3]), 1.0)
		triX, triY := mercator.LatLonToPixels(triLat, triLon, zoom)
		var skip bool
		for _, j := range m.mesh.Triangles[i : i+3] {
			lat, lon := latLonFromVec3(convToVec3(m.r_xyz[j*3:(j*3)+3]), 1.0)
			x, y := mercator.LatLonToPixels(lat, lon, zoom)
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
			lat, lon := latLonFromVec3(convToVec3(m.r_xyz[j*3:(j*3)+3]), 1.0)
			x, y := mercator.LatLonToPixels(lat, lon, zoom)
			path = append(path, [2]float64{x + float64(size)/2, y})
		}
		elev := m.t_elevation[i/3]
		val := (elev - min) / (max - min)
		colVal := uint8(val * 255)
		col := color.NRGBA{
			R: colVal,
			G: colVal,
			B: colVal,
			A: 255,
		}
		if elev < 0 {
			col.B = 255
		} else {
			// Hacky: Modify elevation based on latitude to compensate for colder weather at the poles and warmer weather at the equator.
			valElev := math.Max(math.Min((elev/max)+(math.Sqrt(math.Abs(triLat)/90.0)-0.5), max), 0)
			elevClass := int(valElev*4) + 1
			mois := m.t_moisture[i/3]
			valMois := (mois - minMois) / (maxMois - minMois)
			moisClass := int(valMois*6) + 1
			bio := getBiome(elevClass, moisClass)
			bioCol := biomeColor[bio]
			col.R = uint8(val * float64(bioCol.R))
			col.G = uint8(val * float64(bioCol.G))
			col.B = uint8(val * float64(bioCol.B))
		}

		svg.Path(svgGenD(path), fmt.Sprintf("fill: rgb(%d, %d, %d)", col.R, col.G, col.B)+tmpLine)
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
	size := int(math.Pow(2.0, float64(zoom)) * 256.0)
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	min, max := MinMax(m.r_elevation)
	// TODO: assign region moisture in a better way!
	for r := 0; r < m.mesh.numRegions; r++ {
		lat, lon := latLonFromVec3(convToVec3(m.r_xyz[r*3:(r*3)+3]), 1.0)
		x, y := mercator.LatLonToPixels(lat, lon, zoom)
		val := (m.r_elevation[r] - min) / (max - min)
		colVal := uint8(val * 255)
		col := color.NRGBA{
			R: colVal,
			G: colVal,
			B: colVal,
			A: 255,
		}
		if m.r_elevation[r] < 0 {
			col.B = 255
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
	drawPlates := true
	drawRivers := false
	//xy := stereographicProjection(m.r_xyz)
	//for i := 0; i < len(xy); i += 2 {
	//	w.WriteString(fmt.Sprintf("v %f %f %f \n", xy[i], xy[i+1], 2.0)) //
	//}

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

const (
	BioSnow = iota
	BioTundra
	BioBare
	BioScorched
	BioTaiga
	BioShrubland
	BioTemperateDesert
	BioTemperateRainForest
	BioTemperateDeciduousForest
	BioTropicalRainForest
	BioTropicalSeasonalForest
	BioGrassland
	BioSubtropicalDesert
)

func getBiome(height, moisture int) int {
	// see: http://www-cs-students.stanford.edu/~amitp/game-programming/polygon-map-generation/#biomes
	switch height {
	case 1:
		if moisture > 4 {
			return BioTropicalRainForest
		}
		if moisture > 2 {
			return BioTropicalSeasonalForest
		}
		if moisture == 2 {
			return BioGrassland
		}
		return BioSubtropicalDesert
	case 2:
		if moisture == 6 {
			return BioTemperateRainForest
		}
		if moisture > 3 {
			return BioTemperateDeciduousForest
		}
		if moisture > 1 {
			return BioGrassland
		}
		return BioTemperateDesert
	case 3:
		if moisture > 4 {
			return BioTaiga
		}
		if moisture > 2 {
			return BioShrubland
		}
		return BioTemperateDesert
	case 4:
		switch moisture {
		case 3:
			return BioTundra
		case 2:
			return BioBare
		case 1:
			return BioScorched
		default:
			return BioSnow
		}
	}
	return BioSnow
}

var biomeColor = map[int]color.NRGBA{
	BioSnow:                     {0xFF, 0xFF, 0xFF, 0},
	BioTundra:                   {0xDD, 0xDD, 0xBB, 0},
	BioBare:                     {0xBB, 0xBB, 0xBB, 0},
	BioScorched:                 {0x99, 0x99, 0x99, 0},
	BioTaiga:                    {0xCC, 0xD4, 0xBB, 0},
	BioShrubland:                {0xC4, 0xCC, 0xBB, 0},
	BioTemperateDesert:          {0xE4, 0xE8, 0xCA, 0},
	BioTemperateRainForest:      {0xA4, 0xC4, 0xA8, 0},
	BioTemperateDeciduousForest: {0xB4, 0xC9, 0xA9, 0},
	BioTropicalRainForest:       {0x9C, 0xBB, 0xA9, 0},
	BioTropicalSeasonalForest:   {0xA9, 0xCC, 0xA4, 0},
	BioGrassland:                {0xC4, 0xD4, 0xAA, 0},
	BioSubtropicalDesert:        {0xE9, 0xDD, 0xC7, 0},
}
