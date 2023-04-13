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
	crossPath := gengeometry.LShape{
		Width:     1,
		Length:    1,
		WingWidth: 0.35,
	}

	/*crossPath := gengeometry.CircleShape{
		Radius: 0.5,
		Steps:  100,
	}*/

	path := crossPath.GetPath()
	mesh, err := gengeometry.ExtrudePath(path, 0.2)
	if err != nil {
		log.Fatal(err)
	}
	roofMesh, err := gengeometry.TaperPath(path, 0.2)
	if err != nil {
		log.Println(err)
	} else {
		mesh.AddMesh(roofMesh, vectors.NewVec3(0, 0, 0.2))
	}

	// Add a wing to the house.
	wingPath := gengeometry.LShape{
		Width:     0.8,
		Length:    0.8,
		WingWidth: 0.3,
	}
	wingMesh, err := gengeometry.ExtrudePath(wingPath.GetPath(), 0.2)
	if err != nil {
		log.Println(err)
	}

	// Add a roof to the wing.
	roofMesh, err = gengeometry.TaperPath(wingPath.GetPath(), 0.1)
	if err != nil {
		log.Println(err)
	} else {
		wingMesh.AddMesh(roofMesh, vectors.NewVec3(0, 0, 0.2))
	}

	// Add the wing to the house.
	mesh.AddMesh(wingMesh, vectors.NewVec3(1, 0, 0.0))

	// Save the mesh to a file.
	mesh.ExportToObj("test.obj")

	// Generate a squircle.
	squirclePath := gengeometry.SquircleShape{
		Radius: 0.25,
		Steps:  12,
		Width:  2,
		Length: 2,
	}

	// Save the squircle to a file.
	gengeometry.SavePathAsPNG(squirclePath.GetPath(), "testpath.png", 100)

	generateProcBuilding()
}

func generateProcBuilding() {
	// Generate a path for a building.
	procCrossPath := gengeometry.PlusShape{
		Width:     1,
		Length:    1,
		WingWidth: 0.35,
	}

	// Generate a mesh from a path.
	path := procCrossPath.GetPath()
	mesh, err := gengeometry.ExtrudePath(path, 0.2)
	if err != nil {
		log.Fatal(err)
	}
	roofMesh, err := gengeometry.TaperPath(path, 0.2)
	if err != nil {
		log.Println(err)
	} else {
		mesh.AddMesh(roofMesh, vectors.NewVec3(0, 0, 0.2))
	}

	// Now iterate over the sides and add a cylinder to each side.
	for i, s := range gengeometry.GetPathSides(path) {
		corner := gengeometry.RectangleShape{
			Width:  0.1,
			Length: 0.1,
		}
		heightCorner := 0.3
		cornerPath := corner.GetPath()
		cornerMesh, err := gengeometry.ExtrudePath(cornerPath, heightCorner)
		if err != nil {
			log.Println(err)
		} else {
			// Add the corner to the mesh.
			mesh.AddMesh(cornerMesh, vectors.NewVec3(s.Start.X-corner.Width/2, s.Start.Y-corner.Length/2, 0))
		}

		// Add a roof to the corner.
		roofMesh, err := gengeometry.TaperPath(cornerPath, 0.025)
		if err != nil {
			log.Println(err)
		} else {
			mesh.AddMesh(roofMesh, vectors.NewVec3(s.Start.X-corner.Width/2, s.Start.Y-corner.Length/2, heightCorner))
		}

		if i%3 == 0 {
			continue
		}

		// Add support struts in the middle of each side.
		strutPath := gengeometry.RectangleShape{
			Width:  0.025,
			Length: 0.025,
		}

		midPoint := s.Start.Add(s.End).Mul(0.5)

		// We add two struts.
		midPointA := s.Start.Add(midPoint).Mul(0.5)
		midPointB := midPoint.Add(s.End).Mul(0.5)

		for _, midP := range []vectors.Vec2{midPointA, midPointB, midPoint} {
			strutMesh, err := gengeometry.ExtrudePath(strutPath.GetPath(), 0.15)
			if err != nil {
				log.Println(err)
			} else {
				// Add the strut to the mesh.
				mesh.AddMesh(strutMesh, vectors.NewVec3(midP.X-strutPath.Width/2, midP.Y-strutPath.Length/2, 0))
			}

			// Add a roof to the strut.
			roofMesh, err = gengeometry.TaperPath(strutPath.GetPath(), 0.01)
			if err != nil {
				log.Println(err)
			} else {
				mesh.AddMesh(roofMesh, vectors.NewVec3(midP.X-strutPath.Width/2, midP.Y-strutPath.Length/2, 0.15))
			}
		}
	}

	// Save the mesh to a file.
	mesh.ExportToObj("test_2.obj")
}
