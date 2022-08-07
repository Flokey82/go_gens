package simvillage_tiles

import "math/rand"

// Dimensions represents the dimensions of a map.
type Dimensions struct {
	Width  int
	Height int
}

// indexToXY returns the x and y position of the given index.
func (m *Dimensions) indexToXY(index int) (x, y int) {
	return index % m.Width, index / m.Width
}

// xyToIndex returns the index of the given position.
func (m *Dimensions) xyToIndex(x, y int) int {
	return x + y*m.Width
}

// MapChunk represents a chunk of the map.
type MapChunk struct {
	Dimensions
	Ground        *Layer // Ground or terrain.
	GroundOverlay *Layer // Overlays like carpets, scratchmarks, etc.
	Objects       *Layer // Objects like stones, flowers, etc.
	Structures    *Layer // Structures like walls, gates, fences, etc.
	Roof          *Layer // Roof or ceiling.
}

// newMapChunk returns a new map chunk with the given width and height.
func newMapChunk(width, height int) *MapChunk {
	return &MapChunk{
		Dimensions: Dimensions{
			Width:  width,
			Height: height,
		},
		Ground:        newLayer(width, height),
		GroundOverlay: newLayer(width, height),
		Objects:       newLayer(width, height),
		Structures:    newLayer(width, height),
		Roof:          newLayer(width, height),
	}
}

// Layer represents a layer on the map.
// Note: This code is in part inspired by cxong's fantastic map generator
// https://github.com/cxong/gomapgen
type Layer struct {
	Dimensions
	Tiles []int
}

// newLayer returns a new layer with the given width and height.
func newLayer(width, height int) *Layer {
	return &Layer{
		Dimensions: Dimensions{
			Width:  width,
			Height: height,
		},
		Tiles: make([]int, width*height),
	}
}

// isValid returns true if the given position is valid.
func (l Layer) isValid(x, y int) bool {
	return x >= 0 && x < l.Width && y >= 0 && y < l.Height
}

// getTile returns the tile at the given position (if the position is valid).
func (l Layer) getTile(x, y int) int {
	if l.isValid(x, y) {
		return l.Tiles[x+y*l.Width]
	}
	return 0
}

// setTile sets the tile at the given position (if the position is valid).
func (l *Layer) setTile(x, y, tile int) {
	if l.isValid(x, y) {
		l.Tiles[x+y*l.Width] = tile
	}
}

// fill fills the map with a single tile type.
func (l *Layer) fill(tile int) {
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			l.setTile(x, y, tile)
		}
	}
}

// fillRandom fills the map with a random selction of given tiles.
func (l *Layer) fillRandom(tiles []int) {
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			l.setTile(x, y, tiles[rand.Intn(len(tiles))]) // TODO: supply random number generator
		}
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

// rect represents a rectangle.
type rect struct {
	x, y int // Top left corner
	w    int // Width
	h    int // Height
}
