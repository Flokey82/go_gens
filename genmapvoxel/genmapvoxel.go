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
	voxels [][][]bool
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
			zMax := int64(((noiseVal + 1) / 2) * float64(w.dimZ))
			for z := int64(0); z < zMax; z++ {
				w.voxels[x][y][z] = true // Set the voxel.
			}
		}
	}
	return w
}

// initGrid initializes the voxel grid.
func (w *World) initGrid() {
	w.voxels = make([][][]bool, w.dimX)
	for x := range w.voxels {
		w.voxels[x] = make([][]bool, w.dimY)
		for y := range w.voxels[x] {
			w.voxels[x][y] = make([]bool, w.dimZ)
		}
	}
}

// ExportOBJ exports the world to an OBJ file.
func (w *World) ExportOBJ(filename string) error {
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
				if !w.voxels[x][y][z] {
					continue
				}
				// Check which faces are visible and encode them in the index.
				faceIndex := w.getEncodedIndex(x, y, z)

				// Get the faces to render for this voxel.
				fcs := getFaces(faceIndex)

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
	if x > 0 && w.voxels[x-1][y][z] {
		// Face 0: west.
		faceIndex |= 1 << 0
	}
	if x < w.dimX-1 && w.voxels[x+1][y][z] {
		// Face 1: east.
		faceIndex |= 1 << 1
	}
	if y > 0 && w.voxels[x][y-1][z] {
		// Face 2: north.
		faceIndex |= 1 << 2
	}
	if y < w.dimY-1 && w.voxels[x][y+1][z] {
		// Face 3: south.
		faceIndex |= 1 << 3
	}
	if z < w.dimZ-1 && w.voxels[x][y][z+1] {
		// Face 4: top.
		faceIndex |= 1 << 4
	}
	if z > 0 && w.voxels[x][y][z-1] {
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
func getFaces(faceIndex byte) []Side {
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
	return sides
}

func NewNoise(seed int64) opensimplex.Noise {
	return opensimplex.New(seed)
}
