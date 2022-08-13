package simvillage_tiles

import (
	"math/rand"
)

// Layer represents a layer on the map.
// TODO: Provide a (bool) Collidable property which will determine if the
// layer will be checked in the collision detection.
// Note: This code is in part inspired by cxong's fantastic map generator
// https://github.com/cxong/gomapgen
type Layer struct {
	Dimensions
	Tiles []int
}

// newLayer returns a new layer with the given width and height.
func newLayer(width, height int) *Layer {
	return &Layer{
		Dimensions: NewDimensions(width, height),
		Tiles:      make([]int, width*height),
	}
}

// isValid returns true if the given position is valid.
func (l Layer) isValid(x, y int) bool {
	return x >= 0 && x < l.Width && y >= 0 && y < l.Height
}

// getTile returns the tile at the given position (if the position is valid).
func (l Layer) getTile(x, y int) int {
	if l.isValid(x, y) {
		return l.Tiles[l.xyToIndex(x, y)]
	}
	return 0
}

// setTile sets the tile at the given position (if the position is valid).
func (l *Layer) setTile(x, y, tile int) {
	if l.isValid(x, y) {
		l.Tiles[l.xyToIndex(x, y)] = tile
	}
}

// fill fills the map with a single tile type.
func (l *Layer) fill(tile int) {
	for i := range l.Tiles {
		l.Tiles[i] = tile
	}
}

// fillRandom fills the map with a random selction of given tiles.
func (l *Layer) fillRandom(tiles []int, r *rand.Rand) {
	for i := range l.Tiles {
		l.Tiles[i] = tiles[r.Intn(len(tiles))]
	}
}

// drawRectangle draws a rectangle, optional filled.
func (l *Layer) drawRectangle(r rect, tile int, filled bool) {
	for y := r.y; y < r.y+r.h; y++ {
		for x := r.x; x < r.x+r.w; x++ {
			if filled || x == r.x || y == r.y || x == r.x+r.w-1 || y == r.y+r.h-1 {
				l.setTile(x, y, tile)
			}
		}
	}
}

// Dimensions represents the dimensions of a map.
type Dimensions struct {
	Width  int
	Height int
}

// NewDimensions returns a new Dimensions struct.
func NewDimensions(width, height int) Dimensions {
	return Dimensions{
		Width:  width,
		Height: height,
	}
}

// indexToXY returns the x and y position of the given index.
func (m *Dimensions) indexToXY(index int) (x, y int) {
	return index % m.Width, index / m.Width
}

// xyToIndex returns the index of the given position.
func (m *Dimensions) xyToIndex(x, y int) int {
	return x + y*m.Width
}

// rect represents a rectangle.
type rect struct {
	x, y int // Top left corner
	w    int // Width
	h    int // Height
}
