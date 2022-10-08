package genmap2d

import (
	"image/color"
)

// Various tile IDs.
const (
	TileIDGrass byte = iota
	TileIDWater
	TileIDTree
	TileIDSand
	TileIDMountain
	TileIDSnow
	TileIDVillage
	TileIDMax
)

// TileFromHeight returns the tile ID for a given height.
func (m *Map) TileFromHeight(h int) byte {
	if h <= 0 {
		return TileIDWater
	}
	if h <= 2 {
		return TileIDSand
	}
	if h >= 90 {
		return TileIDSnow
	}
	if h >= 70 {
		return TileIDMountain
	}
	if h >= 30 && h <= 60 {
		return TileIDTree
	}
	return TileIDGrass
}

// TileColor returns the color of a tile based on the given tile ID.
func (m *Map) TileColor(tID byte) color.Color {
	switch tID {
	case TileIDGrass:
		return color.RGBA{0x34, 0x8C, 0x31, 0xff}
	case TileIDWater:
		return color.RGBA{0x00, 0x75, 0x77, 0xff}
	case TileIDTree:
		return color.RGBA{0x42, 0x69, 0x2f, 0xff}
	case TileIDSand:
		return color.RGBA{0xc2, 0xb2, 0x80, 0xff}
	case TileIDMountain:
		return color.RGBA{0x91, 0x8E, 0x85, 0xff}
	case TileIDSnow:
		return color.RGBA{0xFF, 0xFF, 0xFF, 0xff}
	case TileIDVillage:
		return color.RGBA{0xFF, 0x00, 0x00, 0xff}
	default:
		return color.RGBA{0x00, 0x00, 0x00, 0xff}
	}
}

// setup assigns all tiles based on the height of a given point.
func (m *Map) setup() {
	m.run(func(x, y int) byte {
		return m.TileFromHeight(int(m.Elevation[m.GetIndex(x, y)]) - 120)
	})
}

// run executes function 'f' for every point and assignes the returned tile ID.
func (m *Map) run(f func(x, y int) byte) {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.Cells[m.GetIndex(x, y)] = f(x, y)
		}
	}
}

// tilesInRadius returns all tile IDs within a given radious around a given point.
func (m *Map) tilesInRadius(x, y, r int) []byte {
	var res []byte
	for cx := x - r; cx < x+r; cx++ {
		if cx < 0 || cx >= m.Width {
			continue
		}
		for cy := y - r; cy < y+r; cy++ {
			if cy < 0 || cy >= m.Height {
				continue
			}
			if dist(cx, cy, x, y) > r {
				// Not in circle.
				continue
			}
			res = append(res, m.Cells[m.GetIndex(cx, cy)])
		}
	}
	return res
}
