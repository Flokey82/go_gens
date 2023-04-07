package gengeometry

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/llgcode/draw2d/draw2dimg"
)

// BlockyShape represents a shape that is made up of blocks / tiles.
type BlockyShape struct {
	Width, Length int    // Width and length of the shape in tiles
	Set           []bool // Set of tiles that are set
}

func NewRandShape(width, length int) BlockyShape {
	g := NewBlockyShape(width, length)
	for i := range g.Set {
		g.Set[i] = rand.Float64() < 0.5
	}
	return g
}

// NewBlockyShape creates a new BlockyShape with the given width and length.
func NewBlockyShape(width, length int) BlockyShape {
	return BlockyShape{
		Width:  width,
		Length: length,
		Set:    make([]bool, width*length),
	}
}

// Print prints the shape to the console
func (g BlockyShape) Print() {
	for i, v := range g.Set {
		if i%g.Width == 0 {
			fmt.Println()
		}
		if v {
			fmt.Print("X")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

// RenderToImage renders the shape to a supplied image with the given color.
func (g BlockyShape) RenderToImage(img *image.RGBA, color color.RGBA) {
	// New context
	ctx := draw2dimg.NewGraphicContext(img)
	ctx.SetStrokeColor(color)
	ctx.SetFillColor(color)
	ctx.SetLineWidth(1)

	for _, path := range g.GetPaths() {
		ctx.BeginPath()
		for i, v := range path {
			if i == 0 {
				ctx.MoveTo(v.X*100, v.Y*100)
			} else {
				ctx.LineTo(v.X*100, v.Y*100)
			}
		}
		ctx.Close()
		ctx.FillStroke()
	}
	ctx.Close()
}

// GetPaths returns all closed paths that make up the shape.
// TODO: This is a very naive implementation that just adds each square as a path.
// It should be improved to merge squares that are next to each other.
func (g BlockyShape) GetPaths() [][]vectors.Vec2 {
	var shapes [][]int
	visited := make([]bool, len(g.Set))
	var walkConnected func(i int) []int
	walkConnected = func(i int) []int {
		if visited[i] || !g.Set[i] {
			return nil
		}
		visited[i] = true

		// Add all neighbors to the shape and recurse
		var shape []int
		shape = append(shape, i)
		for _, nb := range g.GetNeighbors(i) {
			shape = append(shape, walkConnected(nb)...)
		}
		return shape
	}

	// Walk through all tiles and find all shapes
	for i, v := range g.Set {
		if !v || visited[i] {
			continue
		}
		shapes = append(shapes, walkConnected(i))
	}

	// Convert the shapes to paths
	var paths [][]vectors.Vec2
	for _, shape := range shapes {
		// Just add each square as a path
		// TODO: Merge squares that are next to each other
		for _, idx := range shape {
			x := idx % g.Width
			y := idx / g.Width
			paths = append(paths, []vectors.Vec2{
				{X: float64(x), Y: float64(y)},
				{X: float64(x), Y: float64(y + 1)},
				{X: float64(x + 1), Y: float64(y + 1)},
				{X: float64(x + 1), Y: float64(y)},
			})
		}
	}

	return paths
}

// ConnectionPoints returns the indices of all tiles that are connection points.
// A connection point is a tile that has only one neighbor.
func (g BlockyShape) ConnectionPoints() []int {
	// All tiles that have only one neighbour are connection points
	var connectionPoints []int
	for i, v := range g.Set {
		if !v {
			continue
		}

		// Check if the tile has only one neighbor
		setNeighbors := g.GetNeighbors(i)
		if len(setNeighbors) == 1 || len(setNeighbors) == 2 {
			connectionPoints = append(connectionPoints, i)
		}
	}
	return connectionPoints
}

// GetIsNeighborset returns a boolean array that indicates if the tile at the
// given index has a neighbor in the given direction.
//
// The directions are:
// 0: Top
// 1: Bottom
// 2: Left
// 3: Right
func (g *BlockyShape) GetIsNeighborset(idx int) [4]bool {
	var isNeighborset [4]bool
	// Look up all neighbors and make sure we stay within the bounds of the
	// shape
	neighbors := []int{
		idx - g.Width, // Top
		idx + g.Width, // Bottom
		idx - 1,       // Left
		idx + 1,       // Right
	}
	for i, n := range neighbors {
		if n < 0 || n >= len(g.Set) {
			continue
		}
		if g.Set[n] {
			isNeighborset[i] = true
		}
	}
	return isNeighborset
}

// GetNeighbors returns the indices of all set neighbors of the given index.
func (g BlockyShape) GetNeighbors(idx int) []int {
	// Look up all neighbors and make sure we stay within the bounds of the
	// shape.
	neighbors := []int{
		idx - g.Width, // Top
		idx + g.Width, // Bottom
		idx - 1,       // Left
		idx + 1,       // Right
	}
	var setNeighbors []int
	for _, n := range neighbors {
		if n < 0 || n >= len(g.Set) {
			continue
		}
		if g.Set[n] {
			setNeighbors = append(setNeighbors, n)
		}
	}
	return setNeighbors
}

// GetConnectionDirection returns the direction of the connection of the given index.
func (g BlockyShape) GetConnectionDirection(idx int) int {
	// Get the direction of the connection
	neighbors := g.GetNeighbors(idx)
	if len(neighbors) != 1 {
		return -1
	}
	return neighbors[0] - idx // +1 for right, -1 for left, +width for down, -width for up
}

const (
	RotateRight = iota
	RotateLeft
)

func (g BlockyShape) GetRotated(rotateLeft bool, steps int) BlockyShape {
	// Truncate the number of steps to 0-3
	// 0 steps: no rotation
	// 1 step: rotate 90 degrees
	// 2 steps: rotate 180 degrees
	// 3 steps: rotate 270 degrees
	steps = steps % 4

	// Rotate the shape
	var rotated BlockyShape

	// Figure out the new width and length
	if steps%2 == 0 {
		rotated = NewBlockyShape(g.Width, g.Length)
	} else {
		rotated = NewBlockyShape(g.Length, g.Width)
	}

	// Rotate the shape.
	for i, v := range g.Set {
		if !v {
			continue
		}

		// Get the current x and y coordinates.
		x := i % g.Width
		y := i / g.Width

		// Depending on the rotation direction and number of steps, we need to
		// rotate the index accordingly.
		var newX, newY int
		if steps == 0 {
			newX = x
			newY = y
		} else if steps == 2 {
			newX = g.Width - x - 1
			newY = g.Length - y - 1
		} else if rotateLeft == (steps == 1) {
			// steps == 1 || steps == 3
			newX = g.Length - y - 1
			newY = x
		} else {
			// steps == 1 || steps == 3
			newX = y
			newY = g.Width - x - 1
		}
		rotated.Set[newY*rotated.Width+newX] = true
	}
	return rotated
}

// O-Shape:
//  ______
// |  __  |
// | |__| |
// |______|
//
// U-Shape:
//  _    _
// | |__| |
// |______|
//
// L-Shape:
//  _
// | |____
// |______|
//
// J-Shape:
//    ____
//   |__  |
//  ____| |
// |______|
//
// T-Shape:
//  ______
// |_    _|
//   |__|
//
// Plus-Shape:
//    __
//  _|  |_
// |_    _|
//   |__|
//
// Rectangle:
//  ______
// |______|
//

var (
	BlockyOShape = BlockyShape{
		Width:  3,
		Length: 3,
		Set: []bool{
			true, true, true,
			true, false, true,
			true, true, true,
		},
	}

	BlockyUShape = BlockyShape{
		Width:  3,
		Length: 3,
		Set: []bool{
			true, false, true,
			true, false, true,
			true, true, true,
		},
	}

	BlockyLShape = BlockyShape{
		Width:  3,
		Length: 3,
		Set: []bool{
			true, false, false,
			true, false, false,
			true, true, true,
		},
	}

	BlockyJShape = BlockyShape{
		Width:  3,
		Length: 3,
		Set: []bool{
			false, true, true,
			false, false, true,
			true, true, true,
		},
	}

	BlockyTShape = BlockyShape{
		Width:  3,
		Length: 3,
		Set: []bool{
			true, true, true,
			false, true, false,
			false, true, false,
		},
	}

	BlockyPlusShape = BlockyShape{
		Width:  3,
		Length: 3,
		Set: []bool{
			false, true, false,
			true, true, true,
			false, true, false,
		},
	}
)

var BlockyShapes = []BlockyShape{
	BlockyOShape,
	BlockyUShape,
	BlockyLShape,
	BlockyJShape,
	BlockyTShape,
	BlockyPlusShape,
}

// GetBresenhamLine returns a list of points that are on the line between pos1 and pos2.
// The line is drawn using the Bresenham algorithm.
// See: http://www.roguebasin.com/index.php/Bresenham%27s_Line_Algorithm
func GetBresenhamLine(pos1, pos2 vectors.IVec2) (points []vectors.IVec2) {
	x1, y1 := pos1.X, pos1.Y
	x2, y2 := pos2.X, pos2.Y

	isSteep := abs(y2-y1) > abs(x2-x1)
	if isSteep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}

	reversed := false
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		reversed = true
	}

	deltaX := x2 - x1
	deltaY := abs(y2 - y1)
	err := deltaX / 2
	y := y1
	var ystep int64

	if y1 < y2 {
		ystep = 1
	} else {
		ystep = -1
	}

	for x := x1; x < x2+1; x++ {
		if isSteep {
			points = append(points, vectors.IVec2{X: y, Y: x})
		} else {
			points = append(points, vectors.IVec2{X: x, Y: y})
		}
		err -= deltaY
		if err < 0 {
			y += ystep
			err += deltaX
		}
	}

	if reversed {
		//Reverse the slice
		for i, j := 0, len(points)-1; i < j; i, j = i+1, j-1 {
			points[i], points[j] = points[j], points[i]
		}
	}

	return
}

func abs(x int64) int64 {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0
	}
	return x
}
