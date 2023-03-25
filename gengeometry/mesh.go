package gengeometry

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Flokey82/go_gens/vectors"
)

// Mesh represents a 3d mesh that can be exported to a .obj file.
type Mesh struct {
	Vertices  []vectors.Vec3 // Contains the vertices of the mesh.
	Triangles []int          // Contains the indices of the triangle vertices.
}

// ExportToObj exports the mesh to a .obj file.
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

// AddMesh adds a mesh to the current mesh (at a given vertical offset).
func (m *Mesh) AddMesh(mesh *Mesh, position vectors.Vec3) {
	lenVerts := len(m.Vertices)
	// Add the vertices.
	for _, v := range mesh.Vertices {
		m.Vertices = append(m.Vertices, vectors.Vec3{
			X: v.X + position.X,
			Y: v.Y + position.Y,
			Z: v.Z + position.Z,
		})
	}

	// Add the triangles.
	for _, t := range mesh.Triangles {
		m.Triangles = append(m.Triangles, t+lenVerts)
	}
}

// ExtrudePath extrudes a path to a 3D shape.
func ExtrudePath(path []vectors.Vec2, height float64) (*Mesh, error) {
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
	bottom, err := Triangulate(path)
	if err != nil {
		return nil, err
	}

	// Add the offset to the top indices.
	triBottomTop := bottom
	for _, t := range bottom {
		triBottomTop = append(triBottomTop, t+len(path))
	}

	return &Mesh{
		Vertices:  vertices,
		Triangles: append(triangles, triBottomTop...),
	}, nil
}

// TaperPath tapers a path to a 3D shape to generate a roof.
// NOTE: This is WIP since it doesn't get the angle right.
func TaperPath(path []vectors.Vec2, height float64) (*Mesh, error) {
	// For every point in the path, create two vertices.
	// The first vertex is the point itself, the second vertex is the point
	// with the height added and translated towards the shrunk center.
	// The triangles are created by connecting the vertices in the order
	// they are stored in the indices array.

	// Create the vertices.
	vertices := make([]vectors.Vec3, len(path)*2)
	shrunk := StraightSkeleton(path, height, -height*0.8)
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
	bottom, err := Triangulate(path)
	if err != nil {
		return nil, err
	}
	top, err := Triangulate(shrunk)
	if err != nil {
		return nil, err
	}

	// Add the offset to the top indices.
	triBottomTop := bottom
	for _, t := range top {
		triBottomTop = append(triBottomTop, t+len(path))
	}

	return &Mesh{
		Vertices:  vertices,
		Triangles: append(triangles, triBottomTop...),
	}, nil
}

// Triangulate triangulates a polygon using the ear clipping algorithm.
// It returns the indices of each vertex of each triangle in pairs of 3.
func Triangulate(polygon []vectors.Vec2) ([]int, error) {
	// Create a copy of the polygon.
	poly := make([]vectors.Vec2, len(polygon))
	copy(poly, polygon)
	polyIndex := make([]int, len(polygon))
	for i := range polyIndex {
		polyIndex[i] = i
	}

	// Check if the polygon is clockwise or counter-clockwise.
	// If it is clockwise, reverse the order of the points.
	if !isPolyClockwise(poly) {
		for i := len(poly)/2 - 1; i >= 0; i-- {
			opp := len(poly) - 1 - i
			poly[i], poly[opp] = poly[opp], poly[i]
			polyIndex[i], polyIndex[opp] = polyIndex[opp], polyIndex[i]
		}
	}

	// Create the triangles.
	triangles := make([]int, 0)

	// While there are more than 3 points in the polygon, create triangles.
	for len(poly) > 3 {
		// Find the ear.
		ear := findEar(poly)
		if ear == -1 {
			return nil, errors.New("no ear found")
		}

		// Get the triangle formed by ear and the surrounding points.
		t1, t2, t3 := getTriPointIndexes(len(poly), ear)
		triangles = append(triangles, polyIndex[t1], polyIndex[t2], polyIndex[t3])

		// Remove the ear.
		poly = append(poly[:ear], poly[ear+1:]...)
		polyIndex = append(polyIndex[:ear], polyIndex[ear+1:]...)
	}

	// Create the last triangle.
	triangles = append(triangles, polyIndex[0], polyIndex[1], polyIndex[2])

	return triangles, nil
}

func getTriPointIndexes(polyLen, i int) (int, int, int) {
	return (i - 1 + polyLen) % polyLen, i, (i + 1) % polyLen
}

// findEar finds the ear in a polygon.
func findEar(polygon []vectors.Vec2) int {
	// For every point in the polygon, check if it is an ear.
	for i := 0; i < len(polygon); i++ {
		// Check if the point is an ear.
		if isEar(polygon, i) {
			return i
		}
	}

	// No ear found.
	return -1
}

// isEar checks if a point in a polygon is an ear.
func isEar(polygon []vectors.Vec2, i int) bool {
	// Check if the point is convex.
	// NOTE: This is a bit iffy... not sure if the order of the points might
	// influence if we are looking at the positive or negative side of the
	// angle.
	if !isConvex(polygon, i) {
		return false
	}

	// Check if any other point is inside the triangle formed by the point.
	if isInside(polygon, i) {
		return false
	}

	// We have to make sure, that the side of the triangle that connects
	// the previous and next point is inside the polygon and not outside.
	// If it is outside, we have to skip this triangle.
	// This is done by checking if the midpoint of the side is inside the
	// polygon.
	t1, _, t3 := getTriPointIndexes(len(polygon), i)

	// Get the midpoint of the side.
	mid := vectors.Vec2{
		X: (polygon[t1].X + polygon[t3].X) / 2,
		Y: (polygon[t1].Y + polygon[t3].Y) / 2,
	}

	// Check if the midpoint of the side is inside the polygon.
	if !isInsidePolygon(mid, polygon) {
		// The midpoint is outside the polygon, so we have to skip this
		// triangle.
		return false
	}

	// The point is an ear.
	return true
}

// isConvex checks if a point in a polygon is convex.
func isConvex(polygon []vectors.Vec2, i int) bool {
	// Get the triangle formed by i and the surrounding points.
	p1Idx, p2Idx, p3Idx := getTriPointIndexes(len(polygon), i)

	p1 := polygon[p1Idx]
	p2 := polygon[p2Idx]
	p3 := polygon[p3Idx]

	// Get the vectors from p2 to p1 and p2 to p3.
	v1 := vectors.Sub2(p1, p2)
	v2 := vectors.Sub2(p3, p2)

	// Get the cross product of the vectors.
	cross := vectors.Cross2(v1, v2)

	// If the cross product is positive, the angle is less than 180 degrees.
	// If the cross product is negative, the angle is greater than 180 degrees.
	// If the cross product is 0, the angle is 180 degrees.
	// We want the angle to be less than 180 degrees, so we check if the
	// cross product is positive.
	return cross > 0
}

// isInside checks if any point in a polygon is inside the triangle
// formed aruond the point with index i.
// TODO: We could replace this code with a check if p1->p3 intersect
// any of the sides of the polygon.
func isInside(polygon []vectors.Vec2, i int) bool {
	// Get the triangle formed by i and the surrounding points.
	p1Idx, p2Idx, p3Idx := getTriPointIndexes(len(polygon), i)

	p1 := polygon[p1Idx]
	p2 := polygon[p2Idx]
	p3 := polygon[p3Idx]

	// Create a triangle from the points.
	triangle := []vectors.Vec2{p1, p2, p3}

	// Check if any of the other points is inside the triangle.
	for idx, point := range polygon {
		// Skip the points that are part of the triangle.
		if idx == p1Idx || idx == p2Idx || idx == p3Idx {
			continue
		}
		if vectors.PointInTriangle(point, triangle) {
			return true
		}
	}

	// No point is inside the triangle.
	return false
}

// isInsidePolygon checks if a point is inside a polygon.
func isInsidePolygon(point vectors.Vec2, polygon []vectors.Vec2) bool {
	// Create a ray from the point in any direction.
	ray := vectors.Ray2{
		Origin: point,
		Dir:    vectors.Vec2{X: 0.1, Y: 0.1},
	}

	var intersections int
	for i := 0; i < len(polygon); i++ {
		segment := vectors.Segment{
			Start: polygon[i],
			End:   polygon[(i+1)%len(polygon)],
		}
		if ray.Intersects(segment) {
			intersections++
		}
	}

	// Check if the ray intersects the polygon an odd number of times.
	// If it does, the point is inside the polygon.
	return intersections%2 == 1
}

func isPolyClockwise(polygon []vectors.Vec2) bool {
	var sum float64
	for i := 0; i < len(polygon); i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%len(polygon)]
		sum += (p2.X - p1.X) * (p2.Y + p1.Y)
	}
	return sum > 0
}

/*
func reversePolygon(polygon []vectors.Vec2) []vectors.Vec2 {
	reversed := make([]vectors.Vec2, len(polygon))
	for i, p := range polygon {
		reversed[len(polygon)-1-i] = p
	}
	return reversed
}
*/

// RotatePolygonAroundPoint rotates a polygon around a point by the given angle.
func RotatePolygonAroundPoint(polygon []vectors.Vec2, point vectors.Vec2, angle float64) []vectors.Vec2 {
	rotated := make([]vectors.Vec2, len(polygon))
	for i, p := range polygon {
		rotated[i] = p.RotateAroundPoint(point, angle)
	}
	return rotated
}

// AngleOfSide returns the absolute angle of the normal of the side of the
// polygon with the given index.
func AngleOfSide(polygon []vectors.Vec2, i int) float64 {
	// Get the points of the side.
	p1Idx, p2Idx := getSidePointIndexes(len(polygon), i)
	p1 := polygon[p1Idx]
	p2 := polygon[p2Idx]

	// Get the normal of the side.
	normal := vectors.Sub2(p2, p1).Normal()

	// Get the angle of the normal.
	return normal.Angle()
}

// getSidePointIndexes returns the indexes of the points of the side of the
// polygon with the given index.
func getSidePointIndexes(polygonLen, i int) (int, int) {
	return i, (i + 1) % polygonLen
}
