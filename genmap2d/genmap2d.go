// Package genmap2d implements... <drumroll...> a simple 2d map generator!
package genmap2d

import (
	"image"
	"image/png"
	"math/rand"
	"os"

	opensimplex "github.com/ojrac/opensimplex-go"
)

// Map represents a generated map.
type Map struct {
	Width     int               // Width of the map
	Height    int               // Height of the map
	Villages  []*VillageScore   // Generated villages
	Elevation []float64         // Elevation values for each cell
	Cells     []byte            // Cells contains the assigned tile IDs
	Rand      *rand.Rand        // Rand initialized with the provided seed
	Noise     opensimplex.Noise // Noise initialized with the provided seed
}

// New returns a new map with the given dimensions generated using the given seed.
func New(width, height int, seed int64) *Map {
	m := &Map{
		Width:     width,
		Height:    height,
		Elevation: make([]float64, width*height),
		Cells:     make([]byte, width*height),
		Rand:      rand.New(rand.NewSource(seed)),
		Noise:     opensimplex.New(seed),
	}
	m.genHeightMap()
	m.setup()
	return m
}

func (m *Map) GetIndex(x, y int) int {
	return x + y*m.Width
}

func (m *Map) GetCoordinates(idx int) (int, int) {
	return idx % m.Width, idx / m.Width
}

// Export last frame to a PNG under the given path.
func (m *Map) ExportPng(path string) error {
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, m.Width, m.Height))

	// Set each pixel to the appropriate tile color.
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			img.Set(x, y, m.TileColor(m.Cells[m.GetIndex(x, y)]))
		}
	}

	// Create a file at the given path.
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	// Write the png and close the file.
	defer f.Close()
	return png.Encode(f, img)
}
