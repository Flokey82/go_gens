package gengeometry

import (
	"fmt"
	"log"
	"os"

	"github.com/Flokey82/go_gens/vectors"
)

type Mesh struct {
	Vertices  []vectors.Vec3 // Contains the vertices of the mesh.
	Triangles []int          // Contains the indices of the triangle vertices.
}

func (m *Mesh) ExportToObj(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Write the vertices.
	for _, v := range m.Vertices {
		fmt.Fprintf(f, "v %f %f %f\n", v.X, v.Z, v.Y)
	}

	fmt.Fprintln(f, "")
	// Write the triangles.
	for i := 0; i < len(m.Triangles); i += 3 {
		fmt.Fprintf(f, "f %d %d %d\n", m.Triangles[i]+1, m.Triangles[i+1]+1, m.Triangles[i+2]+1)
	}

	fmt.Fprintln(f, "")

	f.Close()
}

func (m *Mesh) AddMesh(mesh Mesh, verticalOffset float64) {
	lenVerts := len(m.Vertices)
	// Add the vertices.
	for _, v := range mesh.Vertices {
		m.Vertices = append(m.Vertices, vectors.Vec3{X: v.X, Y: v.Y, Z: v.Z + verticalOffset})
	}

	// Add the triangles.
	for _, t := range mesh.Triangles {
		m.Triangles = append(m.Triangles, t+lenVerts)
	}
}

// ExtrudePath extrudes a path to a 3D shape.
func ExtrudePath(path []vectors.Vec2, height float64) Mesh {
	// For every point in the path, create two vertices.
	// The first vertex is the point itself, the second vertex is the point
	// with the height added.
	// The triangles are created by connecting the vertices in the order
	// they are stored in the indices array.

	// Create the vertices.
	vertices := make([]vectors.Vec3, len(path)*2)
	for i, point := range path {
		vertices[i] = vectors.Vec3{X: point.X, Y: point.Y, Z: 0}
		vertices[i+len(path)] = vectors.Vec3{X: point.X, Y: point.Y, Z: height}
	}

	// Create the triangles.
	triangles := make([]int, len(path)*6)
	for i := 0; i < len(path); i++ {
		// Create two triangles for every point in the path.
		triangles[i*6] = i
		triangles[i*6+1] = i + len(path)
		triangles[i*6+2] = (i + 1) % len(path)

		triangles[i*6+3] = (i + 1) % len(path)
		triangles[i*6+4] = i + len(path)
		triangles[i*6+5] = (i+1)%len(path) + len(path)
	}

	// Add triangles for the bottom and top.
	triBottomTop := make([]int, len(path)*3*2)
	for i := 0; i < len(path); i++ {
		// Create two triangles for every point in the path
		// for the bottom and top.
		triBottomTop[i*3] = i
		triBottomTop[i*3+1] = (i + 1) % len(path)
		triBottomTop[i*3+2] = len(path) - 1 - i

		triBottomTop[i*3+len(path)*3] = i + len(path)
		triBottomTop[i*3+1+len(path)*3] = len(path) - 1 - i + len(path)
		triBottomTop[i*3+2+len(path)*3] = (i+1)%len(path) + len(path)
	}

	return Mesh{
		Vertices:  vertices,
		Triangles: append(triangles, triBottomTop...),
	}
}

// TaperPath tapers a path to a 3D shape to generate a roof.
// NOTE: This is WIP since it doesn't get the angle right.
func TaperPath(path []vectors.Vec2, height float64) Mesh {
	// For every point in the path, create two vertices.
	// The first vertex is the point itself, the second vertex is the point
	// with the height added and translated towards the shrunk center.
	// The triangles are created by connecting the vertices in the order
	// they are stored in the indices array.

	// Create the vertices.
	vertices := make([]vectors.Vec3, len(path)*2)
	shrunk := ShrinkPath(path, 0.8)
	shrunk = ShrinkPath(shrunk, 0.6)
	for i, point := range path {
		vertices[i] = vectors.Vec3{X: point.X, Y: point.Y, Z: 0}
		vertices[i+len(path)] = vectors.Vec3{X: shrunk[i].X, Y: shrunk[i].Y, Z: height}
	}

	// Create the triangles.
	triangles := make([]int, len(path)*6)
	for i := 0; i < len(path); i++ {
		// Create two triangles for every point in the path.
		triangles[i*6] = i
		triangles[i*6+1] = i + len(path)
		triangles[i*6+2] = (i + 1) % len(path)

		triangles[i*6+3] = (i + 1) % len(path)
		triangles[i*6+4] = i + len(path)
		triangles[i*6+5] = (i+1)%len(path) + len(path)
	}

	// Add triangles for the bottom and top.
	triBottomTop := make([]int, len(path)*3*2)
	for i := 0; i < len(path); i++ {
		// Create two triangles for every point in the path
		// for the bottom and top.
		triBottomTop[i*3] = i
		triBottomTop[i*3+1] = (i + 1) % len(path)
		triBottomTop[i*3+2] = len(path) - 1 - i

		triBottomTop[i*3+len(path)*3] = i + len(path)
		triBottomTop[i*3+1+len(path)*3] = len(path) - 1 - i + len(path)
		triBottomTop[i*3+2+len(path)*3] = (i+1)%len(path) + len(path)
	}

	return Mesh{
		Vertices:  vertices,
		Triangles: append(triangles, triBottomTop...),
	}
}
