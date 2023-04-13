package genarchitecture

import (
	"log"

	"github.com/Flokey82/go_gens/gengeometry"
	"github.com/Flokey82/go_gens/vectors"
)

func GenerateSampleCathedral() {
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
