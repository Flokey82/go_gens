package genmapvoronoi

import (
	"fmt"
	"github.com/pzsz/voronoi"
	"log"
	"math"
	"math/rand"
	"os"

	svgo "github.com/ajstarks/svgo"
)

// ExportOBJ returns a Wavefront OBJ file representing the heightmap.
func (r *Terrain) ExportOBJ(path string) error {
	return r.h.ExportOBJ(path)
}

// ExportSVG exports the terrain as SVG to the given path.
func (r *Terrain) ExportSVG(path string) error {
	params := r.params
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	width := 3500
	height := int(float64(width) * params.Extent.Height / params.Extent.Width)
	svg := svgo.New(f)
	svg.Start(width, height)

	svgVisualizeHeight(svg, r, width, height)
	svgDrawPaths(svg, r.riverPaths, "stroke=\"blue\" fill=\"none\" stroke-width=\"2\"", width, height)
	svgDrawPaths(svg, r.coasts, "stroke=\"black\" fill=\"none\" stroke-width=\"3\"", width, height)
	svgDrawPaths(svg, r.borders, "stroke=\"red\" fill=\"none\" stroke-width=\"2\"", width, height)
	svgDrawPaths(svg, r.cityBorders, "stroke=\"purple\" fill=\"none\"", width, height)
	svgVisualizeSlopes(svg, r, width, height)
	svgVisualizeCities(svg, r, width, height)
	svgVisualizeRidges(svg, r, width, height)
	//svgDrawLabels(svg, render)

	svg.End()
	return nil
}

func svgVisualizeHeight(svg *svgo.SVG, r *Terrain, width, height int) {
	h := r.h
	h = erosionRate2(h)
	//h = erosionRate3(r.bd,h)
	//h = getFlux(h)
	//h = calcWind(r.bd)
	log.Println(h.Values)
	min, max := h.MinMax()
	for i := range h.VertexTris {
		tmpLine := ""
		var path []voronoi.Vertex
		for j := range h.VertexTris[i] {
			path = append(path, h.VertexTris[i][j].Site)
		}
		if h.Values[i] <= 0 {
			rr := int(math.Abs(((min - h.Values[i]) / min)) * 68)
			rg := int(math.Abs(((min - h.Values[i]) / min)) * 68)
			rb := int(math.Abs(((min - h.Values[i]) / min)) * 255)
			svg.Path(svgGenD(path, width, height), fmt.Sprintf("fill: rgb(%d, %d, %d)", rr, rg, rb)+tmpLine)
		} else {
			rr := int((h.Values[i] / max) * 255)
			rg := int((h.Values[i] / max) * 255)
			if r.sediment.Values[i] > 0 {
				rg = 255
			}
			svg.Path(svgGenD(path, width, height), fmt.Sprintf("fill: rgb(%d, %d, %d)", rr, rg, rr)+tmpLine)
		}
	}
}

func svgDrawPaths(svg *svgo.SVG, paths [][]voronoi.Vertex, attr string, width, height int) {
	for _, path := range paths {
		svg.Path(svgGenD(path, width, height), attr)
	}
}

func svgGenD(path []voronoi.Vertex, width, height int) string {
	var str string

	for i, p := range path {
		if i == 0 {
			str += fmt.Sprintf("M %f,%f", (p.X+0.5)*float64(width), (p.Y+0.5)*float64(height))
			continue
		}
		str += fmt.Sprintf("L %f,%f", (p.X+0.5)*float64(width), (p.Y+0.5)*float64(height))
	}
	return str
}

func svgVisualizeCities(svg *svgo.SVG, render *Terrain, width, height int) {
	cities := render.cities
	h := render.h
	n := render.params.NumTerritories

	for i, city := range cities {
		r := 4
		if i < n {
			r = 10
		}
		svg.Circle(int((h.Vertices[city].X+0.5)*float64(width)), int((h.Vertices[city].Y+0.5)*float64(height)), r, "fill=\"red\"")
	}
}

func svgVisualizeRidges(svg *svgo.SVG, render *Terrain, width, height int) {
	h := render.h
	//visited := make(map[int]bool)
	var ridges [][]voronoi.Vertex
	for i := 0; i < h.Len(); i++ {
		nbs := h.Neighbours(i)
		hei := h.Values[i]
		var nbRidge int
		var lowCount, hiCount int
		for _, nb := range nbs {
			if h.Values[nb] < hei*0.95 {
				lowCount++
				continue
			}
			if math.Abs(h.Values[nb]-hei) < hei*0.15 {
				hiCount++
				nbRidge = nb
			}
		}
		if lowCount >= 1 && hiCount >= 1 {
			vxxs := []voronoi.Vertex{h.Vertices[i], h.Vertices[nbRidge]}
			ridges = append(ridges, vxxs)
		}
	}
	svgDrawPaths(svg, ridges, "stroke=\"gray\" fill=\"none\" stroke-width=\"0.5\"", width, height)
}

/*
func visualizeRidges(svg *svgo.SVG, render *Terrain, width, height int) {
	h := render.h

	//visited := make(map[int]bool)
	var ridges [][]voronoi.Vertex
	for i := 0; i < len(h.h); i++ {
		nbs := h.Mesh.Neighbours(i)
		hei := h.h[i]
		var nbRidge int
		var lowCount, hiCount int
		for _, nb := range nbs {
			if h.h[nb] < hei*0.95 {
				lowCount++
				continue
			}
			if math.Abs(h.h[nb]-hei) < hei*0.15 {
				hiCount++
				nbRidge = nb
			}
		}
		if (lowCount >= 1 && hiCount >= 1) && nbRidge != 0 {
			ridges = append(ridges, []voronoi.Vertex{h.Mesh.Vertices[i], h.Mesh.Vertices[nbRidge]})
		}
	}
	drawPaths(svg, ridges, "stroke=\"gray\" fill=\"none\" stroke-width=\"0.5\"", width, height)
}*/

func svgVisualizeSlopes(svg *svgo.SVG, render *Terrain, width, height int) {
	h := render.h
	var sunStrokes, shadeStrokes [][]voronoi.Vertex
	r := 0.25 / math.Sqrt(float64(h.Len()))
	for i := 0; i < h.Len(); i++ {
		if h.Values[i] <= 0 || render.rivers[i] >= 0 {
			continue
		}
		//if isnearedge(h.Mesh, i){
		//	continue
		//}
		nbs := h.Neighbours(i)
		nbs = append(nbs, i)
		var s, s2 float64
		for j := 0; j < len(nbs); j++ {
			slopes := h.TriSlope(nbs[j])
			s += slopes[0] / 10
			s2 += slopes[1]
		}
		s /= float64(len(nbs))
		s2 /= float64(len(nbs))
		if math.Abs(s) < runif(0.01, 0.4) {
			continue
		}
		l := r * runif(1, 2) * (1 - 0.2*math.Pow(math.Atan(s), 2)) * math.Exp(s2/100)
		x := h.Vertices[i].X
		y := h.Vertices[i].Y

		// Shadow experiment
		ts := h.TriSlope(i)

		if math.Abs(l*s) > 2*r {
			n := math.Floor(math.Abs(l * s / r))
			l /= n
			if n > 4 { // was 4
				n = 4
			}
			for j := 0; j < int(n); j++ {
				u := rand.Float64() * r
				v := rand.Float64() * r

				// Shadow experiment
				if ts[0] > 0 { // ts[1] <= 0 &&
					sunStrokes = append(sunStrokes, []voronoi.Vertex{{x + u - l, y + v + l*s}, {x + u + l, y + v - l*s}})
				} else {
					shadeStrokes = append(shadeStrokes, []voronoi.Vertex{{x + u - l, y + v + l*s}, {x + u + l, y + v - l*s}})
				}
				//strokes = append(strokes, []voronoi.Vertex{voronoi.Vertex{x + u - l, y + v + l*s}, voronoi.Vertex{x + u + l, y + v - l*s}})
			}
		} else {

			// Shadow experiment
			if ts[0] > 0 { // ts[1] <= 0 &&
				sunStrokes = append(sunStrokes, []voronoi.Vertex{{x - l, y + l*s}, {x + l, y - l*s}})
			} else {
				shadeStrokes = append(shadeStrokes, []voronoi.Vertex{{x - l, y + l*s}, {x + l, y - l*s}})
			}
			//strokes = append(strokes, []voronoi.Vertex{voronoi.Vertex{x - l, y + l*s}, voronoi.Vertex{x + l, y - l*s}})
		}
	}
	svgDrawPaths(svg, shadeStrokes, "stroke=\"black\" fill=\"none\" stroke-width=\"1\"", width, height)
	svgDrawPaths(svg, sunStrokes, "stroke=\"gray\" fill=\"none\" stroke-width=\"0.5\"", width, height)
}
