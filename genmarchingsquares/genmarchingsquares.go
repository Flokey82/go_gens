// Package genmarchingsquares implements the marching squares algorithm.
// See: https://en.wikipedia.org/wiki/Marching_squares
package genmarchingsquares

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

// MarchSquares returns a grid of tiles encoded as 4 bit values that are generated from the given pixel grid.
func MarchSquares(pixels [][]bool, dimX, dimY int) [][]byte {
	squares := make([][]byte, dimX-1)
	for i := range squares {
		squares[i] = make([]byte, dimY-1)
	}

	// Walk the pixel grid and generate the tiles.
	//
	//  _____ pixel grid
	// | | | |
	// v v v v
	// +-+-+-+
	// | | | | <-
	// +-+-+-+   | tile grid
	// | | | | <-
	// +-+-+-+
	for x := 0; x < dimX-1; x++ {
		for y := 0; y < dimY-1; y++ {
			squares[x][y] = encodeTile(
				pixels[x][y],     // nw -> 8
				pixels[x+1][y],   // ne -> 4
				pixels[x+1][y+1], // se -> 2
				pixels[x][y+1],   // sw -> 1
			)
		}
	}
	return squares
}

// Encoded tile:
//
// 8-4    nw-ne
// | | <- |   |
// 1-2    sw-se
func encodeTile(nw, ne, se, sw bool) byte {
	var val byte
	if nw {
		val |= 1 << 3
	}
	if ne {
		val |= 1 << 2
	}
	if se {
		val |= 1 << 1
	}
	if sw {
		val |= 1 << 0
	}
	return val
}

// ExportToPNG exports the given encoded tiles to PNG.
func ExportToPNG(squares [][]byte, dimX, dimY, tileSize int, filename string) {
	img := image.NewRGBA(image.Rect(0, 0, dimX*tileSize, dimY*tileSize))
	gc := draw2dimg.NewGraphicContext(img)

	// Set some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(5)
	for x := 0; x < dimX; x++ {
		for y := 0; y < dimY; y++ {
			drawTile(gc, tileSize, x, y, squares[x][y])
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

// drawTile draws the given encoded tile at the x/y grid coordinate.
func drawTile(gc *draw2dimg.GraphicContext, tileSize, tileX, tileY int, encTile byte) {
	offsX := float64(tileX * tileSize)
	offsY := float64(tileY * tileSize)
	switch encTile {
	case 0:
		// No points, nothing to do.
	case 1, 2, 4, 8:
		// Draw triangle, single point.
		var offsets [3][2]float64
		switch encTile {
		case 1:
			offsets = [3][2]float64{
				{0, 0.5},
				{0.5, 1},
				{0, 1},
			}
		case 2:
			offsets = [3][2]float64{
				{0.5, 1},
				{1, 1},
				{1, 0.5},
			}
		case 4:
			offsets = [3][2]float64{
				{1, 0.5},
				{1, 0},
				{0.5, 0},
			}
		case 8:
			offsets = [3][2]float64{
				{0.5, 0},
				{0, 0},
				{0, 0.5},
			}
		}
		var points [][2]float64
		for _, offset := range offsets {
			points = append(points, [2]float64{
				offsX + offset[0]*float64(tileSize),
				offsY + offset[1]*float64(tileSize),
			})
		}
		drawPolygon(gc, points)
	case 3, 6, 9, 12:
		// Draw half tile.
		var offsets [2][2]float64
		switch encTile {
		case 3:
			offsets = [2][2]float64{
				{0, 0.5},
				{1, 1},
			}
		case 6:
			offsets = [2][2]float64{
				{0.5, 0},
				{1, 1},
			}
		case 9:
			offsets = [2][2]float64{
				{0, 0},
				{0.5, 1},
			}
		case 12:
			offsets = [2][2]float64{
				{0, 0},
				{1, 0.5},
			}
		}
		drawRectangle(gc, [2]float64{offsX + offsets[0][0]*float64(tileSize), offsY + offsets[0][1]*float64(tileSize)},
			[2]float64{offsX + offsets[1][0]*float64(tileSize), offsY + offsets[1][1]*float64(tileSize)})
	case 5, 10:
		// Draw diagonal.
		var offsets [6][2]float64
		switch encTile {
		case 5:
			offsets = [6][2]float64{
				{0.5, 0},
				{1, 0},
				{1, 0.5},
				{0.5, 1},
				{0, 1},
				{0, 0.5},
			}
		case 10:
			offsets = [6][2]float64{
				{0, 0.5},
				{0.5, 1},
				{1, 1},
				{1, 0.5},
				{0.5, 0},
				{0, 0},
			}
		}
		var points [][2]float64
		for _, offset := range offsets {
			points = append(points, [2]float64{
				offsX + offset[0]*float64(tileSize),
				offsY + offset[1]*float64(tileSize),
			})
		}
		drawPolygon(gc, points)
	case 7, 11, 13, 14:
		// Draw tile minus triangle, 3 points
		var offsets [5][2]float64
		switch encTile {
		case 7:
			offsets = [5][2]float64{
				{0.5, 0},
				{1, 0},
				{1, 1},
				{0, 1},
				{0, 0.5},
			}
		case 11:
			offsets = [5][2]float64{
				{0, 0},
				{0.5, 0},
				{1, 0.5},
				{1, 1},
				{0, 1},
			}
		case 13:
			offsets = [5][2]float64{
				{0, 0},
				{1, 0},
				{1, 0.5},
				{0.5, 1},
				{0, 1},
			}
		case 14:
			offsets = [5][2]float64{
				{0, 0},
				{1, 0},
				{1, 1},
				{0.5, 1},
				{0, 0.5},
			}
		}
		var points [][2]float64
		for _, offset := range offsets {
			points = append(points, [2]float64{
				offsX + offset[0]*float64(tileSize),
				offsY + offset[1]*float64(tileSize),
			})
		}
		drawPolygon(gc, points)
	case 15:
		// Full tile
		drawRectangle(gc,
			[2]float64{float64(tileX * tileSize), float64(tileY * tileSize)},
			[2]float64{float64(tileX*tileSize) + float64(tileSize), float64(tileY*tileSize) + float64(tileSize)},
		)
	}
}

func drawRectangle(gc *draw2dimg.GraphicContext, p1, p2 [2]float64) {
	gc.MoveTo(p1[0], p1[1])
	gc.LineTo(p2[0], p1[1])
	gc.LineTo(p2[0], p2[1])
	gc.LineTo(p1[0], p2[1])
	gc.Close()
	gc.Fill()
}

func drawPolygon(gc *draw2dimg.GraphicContext, points [][2]float64) {
	gc.MoveTo(points[0][0], points[0][1])
	for _, p := range points[1:] {
		gc.LineTo(p[0], p[1])
	}
	gc.Close()
	gc.Fill()
}
