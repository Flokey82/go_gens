// Package genmarchingsquares implements the marching squares algorithm.
// See: https://en.wikipedia.org/wiki/Marching_squares
package genmarchingsquares

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
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
//
//	 ____nw 1 << 3
//	| ___ne 1 << 2
//	|| __se 1 << 1
//	||| _sw 1 << 0
//	||||
//
// b1111
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
func ExportToPNG(squares [][]byte, dimX, dimY, tileSize int, filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, dimX*tileSize, dimY*tileSize))
	gc := draw2dimg.NewGraphicContext(img)

	// Set some properties
	gc.SetLineWidth(5)
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	for x := 0; x < dimX; x++ {
		for y := 0; y < dimY; y++ {
			drawTile(gc, tileSize, x, y, squares[x][y])
		}
	}
	gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	for x := 0; x < dimX; x++ {
		for y := 0; y < dimY; y++ {
			drawCode(gc, x, y, tileSize, squares[x][y])
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

var (
	// Vector points for one active pixel (configuration 1)
	baseOffsetOneTile_1 = [][2]float64{
		{0, 0.5},
		{0.5, 1},
		{0, 1},
	}

	// Vector points for two active pixels / half tile (configuration 3)
	baseOffsetHalfTile_3 = [][2]float64{
		{0, 0.5},
		{1, 0.5},
		{1, 1},
		{0, 1},
	}

	// Vector points for two active pixels / diagonal tile (configuration 5)
	baseOffsetDiagonalTile_5 = [][2]float64{
		{0.5, 0},
		{1, 0},
		{1, 0.5},
		{0.5, 1},
		{0, 1},
		{0, 0.5},
	}

	// Vector points for three active pixels / three quarter tile (configuration 7)
	baseOffsetThreeQuarterTile_7 = [][2]float64{
		{0.5, 0},
		{1, 0},
		{1, 1},
		{0, 1},
		{0, 0.5},
	}

	// Vector points for four active pixels / full tile (configuration 15)
	baseOffsetFull = [][2]float64{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}
)

// drawTile draws the given encoded tile at the x/y grid coordinate.
func drawTile(gc *draw2dimg.GraphicContext, tileSize, tileX, tileY int, encTile byte) {
	offsX := float64(tileX * tileSize)
	offsY := float64(tileY * tileSize)

	// offsetPoints rotates and translates the points to the tile location.
	offsetPoints := func(baseOffset [][2]float64, angle float64) [][2]float64 {
		var points [][2]float64
		for _, offset := range rotatePoints(baseOffset, angle) {
			points = append(points, [2]float64{
				offsX + offset[0]*float64(tileSize),
				offsY + offset[1]*float64(tileSize),
			})
		}
		return points
	}

	var (
		baseOffset [][2]float64
		angle      float64
	)
	switch encTile {
	case 0:
		// No points, nothing to do.
		return
	case 1: // Draw triangle, single point (configuration 1).
		baseOffset = baseOffsetOneTile_1
		angle = 0
	case 2: // Draw triangle, single point (configuration 2).
		baseOffset = baseOffsetOneTile_1
		angle = -90
	case 4: // Draw triangle, single point (configuration 4).
		baseOffset = baseOffsetOneTile_1
		angle = -180
	case 8: // Draw triangle, single point (configuration 8).
		baseOffset = baseOffsetOneTile_1
		angle = -270
	case 3: // Draw half tile (configuration 3).
		baseOffset = baseOffsetHalfTile_3
		angle = 0
	case 6: // Draw half tile (configuration 6).
		baseOffset = baseOffsetHalfTile_3
		angle = -90
	case 9: // Draw half tile (configuration 9).
		baseOffset = baseOffsetHalfTile_3
		angle = 90
	case 12: // Draw half tile (configuration 12).
		baseOffset = baseOffsetHalfTile_3
		angle = 180
	case 5: // Draw diagonal (configuration 5).
		baseOffset = baseOffsetDiagonalTile_5
		angle = 0
	case 10: // Draw diagonal (configuration 10).
		baseOffset = baseOffsetDiagonalTile_5
		angle = 90
	case 7: // Draw tile minus triangle, 3 points (configuration 7).
		baseOffset = baseOffsetThreeQuarterTile_7
		angle = 0
	case 11: // Draw tile minus triangle, 3 points (configuration 11).
		baseOffset = baseOffsetThreeQuarterTile_7
		angle = 90
	case 13: // Draw tile minus triangle, 3 points (configuration 13).
		baseOffset = baseOffsetThreeQuarterTile_7
		angle = 180
	case 14: // Draw tile minus triangle, 3 points (configuration 14).
		baseOffset = baseOffsetThreeQuarterTile_7
		angle = 270
	case 15: // Full tile (configuration 15)
		baseOffset = baseOffsetFull
		angle = 0
	}
	drawPolygon(gc, offsetPoints(baseOffset, angle))
}

// drawCode draws circles at the tile corners which are encoded in the tile.
// A red circle indicates a set bit.
func drawCode(gc *draw2dimg.GraphicContext, tileX, tileY, tileSize int, enctile byte) {
	offsX := float64(tileX * tileSize)
	offsY := float64(tileY * tileSize)
	radius := 10.0
	if enctile&(1<<0) != 0 {
		draw2dkit.Circle(gc, offsX, offsY+float64(tileSize), radius)
	}
	if enctile&(1<<1) != 0 {
		draw2dkit.Circle(gc, offsX+float64(tileSize), offsY+float64(tileSize), radius)
	}
	if enctile&(1<<2) != 0 {
		draw2dkit.Circle(gc, offsX+float64(tileSize), offsY, radius)
	}
	if enctile&(1<<3) != 0 {
		draw2dkit.Circle(gc, offsX, offsY, radius)
	}
	gc.Fill()
}

// drawPolygon draws a polygon from the given points.
func drawPolygon(gc *draw2dimg.GraphicContext, points [][2]float64) {
	gc.MoveTo(points[0][0], points[0][1])
	for _, p := range points[1:] {
		gc.LineTo(p[0], p[1])
	}
	gc.Close()
	gc.Fill()
}

// RotatePoints rotates a number of given polygon points by a specified angle.
func rotatePoints(points [][2]float64, angleDeg float64) [][2]float64 {
	angleRad := angleDeg * math.Pi / 180 // Convert angle to rad.
	var res [][2]float64
	for _, srcPt := range points {
		// Rotate point around center of tile (0.5, 0.5).
		res = append(res, rotatePoint(0.5, 0.5, angleRad, srcPt))
	}
	return res
}

// cx, cy defines the point around which we rotate.
// Based on: https://stackoverflow.com/questions/2259476/rotating-a-point-about-another-point-2d
func rotatePoint(cx, cy, angleRad float64, p [2]float64) [2]float64 {
	s := math.Sin(angleRad)
	c := math.Cos(angleRad)

	// Translate point back to origin.
	p[0] -= cx
	p[1] -= cy

	// Rotate point.
	xnew := p[0]*c - p[1]*s
	ynew := p[0]*s + p[1]*c

	// Translate point back and return.
	return [2]float64{xnew + cx, ynew + cy}
}
