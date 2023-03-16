package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/Flokey82/go_gens/gengeometry"
	"github.com/mazznoer/colorgrad"
)

func main() {
	// Set up a normal square as polygon.
	poly := gengeometry.Polygon{
		Points: []gengeometry.Point{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 1, Y: 0},
		},
	}

	// Set up a U shape as polygon.
	poly = gengeometry.Polygon{
		Points: []gengeometry.Point{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 1, Y: 0},
			{X: 0.8, Y: 0},
			{X: 0.8, Y: 0.8},
			{X: 0.2, Y: 0.8},
			{X: 0.2, Y: 0},
		},
	}

	// Set up a line that cuts the polygon in half.
	line := gengeometry.Line{
		Start: gengeometry.Point{X: -1, Y: 0.5},
		End:   gengeometry.Point{X: 2, Y: 0.5},
	}

	// Set up a diagonal line that cuts the polygon in half.
	line = gengeometry.Line{
		Start: gengeometry.Point{X: -0.9, Y: -1.1},
		End:   gengeometry.Point{X: 1.2, Y: 1.1},
	}

	// Split into multiple polygons.
	polys := poly.Split(line)

	// draw the polygons to an image.
	f, err := os.Create("test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Get a new color palette.
	grad := colorgrad.Rainbow()
	cols := grad.Colors(uint(len(polys) + 1))

	img := image.NewRGBA(image.Rect(0, 0, 1001, 1001))
	poly.DrawToImage(img, color.RGBA{255, 255, 255, 255}, 1000)
	for i, p := range polys {
		cr, cg, cb, ca := cols[i].RGBA()
		// Generate a color for each polygon.
		col := color.RGBA{
			R: uint8(cr >> 8),
			G: uint8(cg >> 8),
			B: uint8(cb >> 8),
			A: uint8(ca >> 8),
		}
		log.Printf("Poly %d:", i)
		p.Log()
		p.DrawToImage(img, col, 900/float64(i+1))
	}
	line.DrawToImage(img, color.RGBA{0, 0, 255, 255}, 1000)

	png.Encode(f, img)
}
