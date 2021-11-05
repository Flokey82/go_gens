// Package gencellular implements cellular automata (for now Conway's Game of Life).
package gencellular

import (
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"
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

// Culture was a very smart way to call whatever holds the cells.
type Culture struct {
	Cells      [2][][]bool       // Cell buffers.
	Height     int               // Number of cells over x.
	Width      int               // Number of cells over y.
	Generation int               // Number of ticks.
	Init       SeedFunc          // Function used to initialize the culture.
	images     []*image.Paletted // Generated frame used to construct the GIF.
	palette    []color.Color     // Default color palette.
	delays     []int             // Delay for each individual frame (0 for now).
}

// New returns a new cell culture with the given height and width.
func New(height, width int) *Culture {
	return NewCustom(height, width, SeedDefault)
}

// NewCustom allows for a custom seed init function.
func NewCustom(height, width int, f SeedFunc) *Culture {
	c := &Culture{
		Height: height,
		Width:  width,
		Init:   f,
		palette: []color.Color{
			color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
			color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
			color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
			color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
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
	c.images = nil
	c.delays = nil
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
			neighbors := c.countNeighbours(current, x, y)
			if current[x][y] {
				next[x][y] = neighbors == 2 || neighbors == 3
			} else {
				next[x][y] = neighbors == 3
			}
		}
	}
	c.storeGifFrame(current)
	c.Generation++
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

func (c *Culture) storeGifFrame(cells [][]bool) {
	img := image.NewPaletted(image.Rect(0, 0, c.Width, c.Height), c.palette)
	c.images = append(c.images, img)
	c.delays = append(c.delays, 0)
	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			if cells[x][y] {
				img.Set(x, y, color.RGBA{0xFF, 0x00, 0x00, 255})
			} else {
				img.Set(x, y, color.RGBA{0x00, 0x00, 0x00, 255})
			}
		}
	}
}

// Export all frames to a GIF under the given path.
func (c *Culture) ExportGif(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	if err := gif.EncodeAll(f, &gif.GIF{
		Image: c.images,
		Delay: c.delays,
	}); err != nil {
		return err
	}
	return f.Close()
}
