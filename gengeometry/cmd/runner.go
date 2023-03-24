package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/Flokey82/go_gens/gengeometry"
	"github.com/Flokey82/go_gens/vectors"
	"github.com/mazznoer/colorgrad"
)

func main() {
	// Set up a normal square as polygon.
	poly := gengeometry.Polygon{
		Points: []vectors.Vec2{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
			{X: 1, Y: 0},
		},
	}

	// Set up a U shape as polygon.
	poly = gengeometry.Polygon{
		Points: []vectors.Vec2{
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
	line := vectors.NewSegment(
		vectors.Vec2{X: -1, Y: 0.5},
		vectors.Vec2{X: 2, Y: 0.5},
	)

	// Set up a diagonal line that cuts the polygon in half.
	line = vectors.NewSegment(
		vectors.Vec2{X: -0.9, Y: -1.1},
		vectors.Vec2{X: 1.2, Y: 1.1},
	)

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
		p.DrawToImage(img, col, 900/float64(i+1))
	}
	gengeometry.DrawLine(img, line.Start, line.End, color.RGBA{0, 0, 255, 255}, 1000)

	png.Encode(f, img)

	// Make a house :]
	// Generate a mesh from a path.
	crossPath := gengeometry.HShape{
		Width:     1,
		Length:    1,
		WingWidth: 0.2,
	}
	path := crossPath.GetPath()
	mesh, err := gengeometry.ExtrudePath(path, 0.2)
	if err != nil {
		log.Fatal(err)
	}
	roofMesh, err := gengeometry.TaperPath(path, 0.1)
	if err != nil {
		log.Fatal(err)
	}
	mesh.AddMesh(roofMesh, 0.2)

	mesh.ExportToObj("test.obj")
}
