// Package gencellular implements cellular automata (for now Conway's Game of Life).
package gencellular

import (
	"math/rand"
)

// SeedFunc is the type of function used to initialize the simulation.
type SeedFunc func(cells [][]bool, w, h int)

// SeedDefault is the default init function used to initialize the simulation.
func SeedDefault(cells [][]bool, w, h int) {
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cells[x][y] = rand.Intn(5) > 1
		}
	}
}

// EvalFunc is the type of function used to evaluate the state of a cell in the next generation.
type EvalFunc func(currState bool, numNeighbors int) bool

// EvalDefault is the default evaluation function used to evaluate the state of a cell in the next generation.
func EvalDefault(currState bool, numNeighbors int) bool {
	return (currState && numNeighbors == 2) || numNeighbors == 3
}

// Culture was a very smart way to call whatever holds the cells.
type Culture struct {
	Cells      [2][][]bool // Cell buffers.
	Height     int         // Number of cells over x.
	Width      int         // Number of cells over y.
	Generation int         // Number of ticks.
	Init       SeedFunc    // Function used to initialize the culture.
	Eval       EvalFunc    // Function used to evaluate the state of a cell in the next generation.
}

// New returns a new cell culture with the given height and width.
func New(height, width int) *Culture {
	return NewCustom(height, width, SeedDefault, EvalDefault)
}

// NewCustom allows for a custom seed init and state eval function.
func NewCustom(height, width int, sf SeedFunc, ef EvalFunc) *Culture {
	c := &Culture{
		Height: height,
		Width:  width,
		Init:   sf,
		Eval:   ef,
	}
	c.Reset()
	return c
}

// Reset the simulation.
func (c *Culture) Reset() {
	c.Cells[0] = initCells(c.Width, c.Height)
	c.Cells[1] = initCells(c.Width, c.Height)
	c.Init(c.Cells[0], c.Width, c.Height)
	c.Generation = 0
}

// initCells initializes a two dimensional slice of bools with the given 'w'idth and 'h'eight.
func initCells(w, h int) [][]bool {
	nc := make([][]bool, w)
	for i := range nc {
		nc[i] = make([]bool, h)
	}
	return nc
}

// Tick advances the culture by one tick.
func (c *Culture) Tick() {
	current := c.Cells[c.Generation%2]
	next := c.Cells[(c.Generation+1)%2]
	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			next[x][y] = c.Eval(current[x][y], c.countNeighbours(current, x, y))
		}
	}
	c.Generation++
}

var directions [][2]int = [][2]int{
	{1, 0},   // East
	{0, 1},   // South
	{1, 1},   // SouthEast
	{0, -1},  // North
	{1, -1},  // NorthEast
	{-1, -1}, // NorthWest
	{-1, 0},  // West
	{-1, 1},  // SouthWest
}

// countNeighbours counts the neighbours of the given cell.
func (c *Culture) countNeighbours(cells [][]bool, x, y int) (sum int) {
	var dx, dy int
	for _, dir := range directions {
		dx = x + dir[0]
		if dx < 0 || dx >= c.Width {
			continue
		}
		dy = y + dir[1]
		if dy < 0 || dy >= c.Height || !cells[dx][dy] {
			continue
		}
		sum++
	}
	return sum
}
