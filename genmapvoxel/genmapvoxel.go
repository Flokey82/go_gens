package genmapvoxel

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Flokey82/go_gens/vectors"

	opensimplex "github.com/ojrac/opensimplex-go"
)

// World represents a voxel world.
type World struct {
	dimX   int64
	dimY   int64
	dimZ   int64
	Voxels [][][]bool    // The voxel grid.
	Values [][][]float64 // Voxels values (full voxel height = 1.0).
}

// New returns a new voxel world.
func New(dimX, dimY, dimZ, seed int64) *World {
	// Generate the world.
	w := &World{
		dimX: dimX,
		dimY: dimY,
		dimZ: dimZ,
	}

	// Initialize the voxel grid
	w.initGrid()

	// Now initialize the opensimplex noise generator and generate the world.
	noise := NewNoise(seed)
	for x := int64(0); x < w.dimX; x++ {
		for y := int64(0); y < w.dimY; y++ {
			// Get the Z dimension by using the noise generator.
			noiseVal := noise.Eval2(float64(x)/float64(w.dimX), float64(y)/float64(w.dimY))
			zVal := ((noiseVal + 1) / 2) * float64(w.dimZ)

			// Round up zVal to the next higher integer.
			zMax := int64(zVal + 0.99999)
			for z := int64(0); z < zMax; z++ {
				w.Voxels[x][y][z] = true              // Set the voxel.
				w.Values[x][y][z] = zVal - float64(z) // Set the voxel value.
			}
		}
	}
	return w
}

// initGrid initializes the voxel grid.
func (w *World) initGrid() {
	w.Voxels = make([][][]bool, w.dimX)
	for x := range w.Voxels {
		w.Voxels[x] = make([][]bool, w.dimY)
		for y := range w.Voxels[x] {
			w.Voxels[x][y] = make([]bool, w.dimZ)
		}
	}

	w.Values = make([][][]float64, w.dimX)
	for x := range w.Values {
		w.Values[x] = make([][]float64, w.dimY)
		for y := range w.Values[x] {
			w.Values[x][y] = make([]float64, w.dimZ)
		}
	}
}

// ExportOBJ exports the world to an OBJ file.
func (w *World) ExportOBJ(filename string, smooth bool) error {
	// Open/create the destination file.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Initialize a new bufio writer.
	wr := bufio.NewWriter(f)

	// Make sure to flush the buffer on exit.
	defer wr.Flush()

	var vertices []vectors.Vec3
	var faces [][]int

	// Iterate over voxel grid and determine visible faces.
	for x := int64(0); x < w.dimX; x++ {
		for y := int64(0); y < w.dimY; y++ {
			for z := int64(0); z < w.dimZ; z++ {
				if !w.Voxels[x][y][z] {
					continue
				}
				// Check which faces are visible and encode them in the index.
				faceIndex := w.getEncodedIndex(x, y, z)

				// If we should smooth the terrain by using the float values, do so.
				value := 1.0
				if smooth {
					value = w.Values[x][y][z]
				}

				// Get the faces to render for this voxel.
				fcs := getFaces(faceIndex, value)

				// Add the vertices to the list.
				for _, f := range fcs {
					var faceIndices []int
					// Translate the faces to the voxel position.
					for _, v := range f.Translate(float64(x), float64(y), float64(z)) {
						vertices = append(vertices, v)
						faceIndices = append(faceIndices, len(vertices))
					}
					// Append the face vertex indices to the list of faces.
					faces = append(faces, faceIndices)
				}
			}
		}
	}

	// Write all the vertices to the file.
	for _, v := range vertices {
		// NOTE: I switched Y and Z since importing into Blender would have the Y axis as the up axis.
		if _, err := wr.WriteString(fmt.Sprintf("v %f %f %f \n", v.X, v.Z, v.Y)); err != nil {
			return err
		}
	}

	// Write all the faces to the file.
	for _, f := range faces {
		if _, err := wr.WriteString(fmt.Sprintf("f %d %d %d %d \n", f[0], f[1], f[2], f[3])); err != nil {
			return err
		}
	}

	return nil
}

func (w *World) getEncodedIndex(x, y, z int64) byte {
	// Check which faces are visible and encode them in the index.
	var faceIndex byte
	if x > 0 && w.Voxels[x-1][y][z] {
		// Face 0: west.
		faceIndex |= 1 << 0
	}
	if x < w.dimX-1 && w.Voxels[x+1][y][z] {
		// Face 1: east.
		faceIndex |= 1 << 1
	}
	if y > 0 && w.Voxels[x][y-1][z] {
		// Face 2: north.
		faceIndex |= 1 << 2
	}
	if y < w.dimY-1 && w.Voxels[x][y+1][z] {
		// Face 3: south.
		faceIndex |= 1 << 3
	}
	if z < w.dimZ-1 && w.Voxels[x][y][z+1] {
		// Face 4: top.
		faceIndex |= 1 << 4
	}
	if z > 0 && w.Voxels[x][y][z-1] {
		// Face 5: bottom.
		faceIndex |= 1 << 5
	}
	return faceIndex
}

type Side [4]vectors.Vec3

// Translate translates the side by the given offset.
func (s Side) Translate(x, y, z float64) Side {
	for i := range s {
		s[i].X += x
		s[i].Y += y
		s[i].Z += z
	}
	return s
}

// Shrink the side vertically by the given factor.
func (s Side) Shrink(factor float64) Side {
	for i := range s {
		if s[i].Z != 0 {
			s[i].Z = ((s[i].Z + 0.5) * factor) - 0.5
		}
	}
	return s
}

var baseSideEastWest = Side{
	{X: 0, Y: -0.5, Z: -0.5},
	{X: 0, Y: 0.5, Z: -0.5},
	{X: 0, Y: 0.5, Z: 0.5},
	{X: 0, Y: -0.5, Z: 0.5},
}

var baseSideNorthSouth = Side{
	{X: -0.5, Y: 0, Z: -0.5},
	{X: 0.5, Y: 0, Z: -0.5},
	{X: 0.5, Y: 0, Z: 0.5},
	{X: -0.5, Y: 0, Z: 0.5},
}

var baseSideTopBottom = Side{
	{X: -0.5, Y: -0.5, Z: 0},
	{X: 0.5, Y: -0.5, Z: 0},
	{X: 0.5, Y: 0.5, Z: 0},
	{X: -0.5, Y: 0.5, Z: 0},
}

// getFaces returns the faces to render given the encoded faceIndex.
//
// TODO: Allow shrinking the faces for voxels with a value < 1.0 to
// generate "shorter" cubes.
func getFaces(faceIndex byte, height float64) []Side {
	// If the index is 0, then all faces are invisible.
	if faceIndex == 0 {
		return nil
	}

	// Get the face vertices to render for this voxel.
	var sides []Side
	if faceIndex&1<<0 != 0 {
		// Face 0: west.
		sides = append(sides, baseSideEastWest.Translate(-0.5, 0, 0))
	}
	if faceIndex&1<<1 != 0 {
		// Face 1: east.
		sides = append(sides, baseSideEastWest.Translate(0.5, 0, 0))
	}
	if faceIndex&1<<2 != 0 {
		// Face 2: north.
		sides = append(sides, baseSideNorthSouth.Translate(0, -0.5, 0))
	}
	if faceIndex&1<<3 != 0 {
		// Face 3: south.
		sides = append(sides, baseSideNorthSouth.Translate(0, 0.5, 0))
	}
	if faceIndex&1<<4 != 0 {
		// Face 4: top.
		sides = append(sides, baseSideTopBottom.Translate(0, 0, 0.5))
	}
	if faceIndex&1<<5 != 0 {
		// Face 5: bottom.
		sides = append(sides, baseSideTopBottom.Translate(0, 0, -0.5))
	}

	// Check if we need to shrink the sides.
	if height == 1.0 {
		return sides
	}

	// Now shrink the sides vertically by the given factor.
	for i := range sides {
		sides[i] = sides[i].Shrink(height)
	}
	return sides
}

func NewNoise(seed int64) opensimplex.Noise {
	return opensimplex.New(seed)
}
