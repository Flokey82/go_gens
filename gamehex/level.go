package gamehex

import (
	"fmt"
	"image/color"
	"math"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/ojrac/opensimplex-go"
)

type Level struct {
	Width     int
	Height    int
	hexRadius int
	Tiles     []Tile
}

func NewLevel(width, height int) (*Level, error) {
	l := &Level{
		Width:     width,
		Height:    height,
		hexRadius: 32,
		Tiles:     make([]Tile, width*height),
	}

	noise := opensimplex.New(0)
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			// Calculate the noise value for the current tile.
			tx, ty := l.HexTileXYToPixelPos(x, y)
			l.Tiles[y*l.Width+x] = Tile(noise.Eval2(float64(tx), float64(ty)) * 255)
		}
	}

	return l, nil
}

// HexTilePos returns the center coordinates in pixel for the given hex tile in the hex grid.
// NOTE: This is the flat top version of the hex grid.
func (l *Level) HexTileXYToPixelPos(x, y int) (int, int) {
	// Calculate width and height of a hexagon.
	w := l.hexRadius * 2           // width of a hexagon
	h := math.Sqrt(3) * float64(w) // height of a hexagon

	// Calculate x and y position.
	// We start at the top left corner of the grid.
	xPos := float64(x)*float64(w)*0.75 + float64(l.hexRadius)
	yPos := float64(y)*float64(h)/2 + float64(l.hexRadius)

	// If the y position is uneven, we move the x position 0.75 of the width to the right.
	if x%2 != 0 {
		yPos += float64(h) / 4
	}

	return int(xPos), int(yPos)
}

// TileAtPos returns which tile the given position is in.
// NOTE: This is the flat top version of the hex grid.
func (l *Level) TileAtPos(px, py int) (int, int) {
	// Given the position of the mouse, we calculate the tile it is in.
	// Get the column of tiles that we are probably in.
	x := int(math.Floor(float64(px) / float64(l.hexRadius) / 1.5))
	// Get the row of tiles that we are probably in.
	y := int(math.Floor(float64(py) / float64(l.hexRadius) / math.Sqrt(3)))

	// Calculate the center of the tile.
	xCenter, yCenter := l.HexTileXYToPixelPos(x, y)

	// Calculate the distance from the center of the tile to the mouse position.
	dx := float64(px) - float64(xCenter)
	dy := float64(py) - float64(yCenter)

	// If the distance is smaller than the radius of the hexagon, we are in the tile.
	if math.Sqrt(dx*dx+dy*dy) < float64(l.hexRadius) {
		return x, y
	}

	// If we are not in the tile, we check if we are in any of the surrounding tiles.
	for _, p := range l.getNeighbors(x, y) {
		// Calculate the center of the tile.
		xCenter, yCenter := l.HexTileXYToPixelPos(p[0], p[1])

		// Calculate the distance from the center of the tile to the mouse position.
		dx := float64(px) - float64(xCenter)
		dy := float64(py) - float64(yCenter)

		// If the distance is smaller than the radius of the hexagon, we are in the tile.
		if math.Sqrt(dx*dx+dy*dy) < float64(l.hexRadius) {
			return p[0], p[1]
		}
	}

	// If we are not in any of the surrounding tiles, we are not in any tile.
	return -1, -1
}

func (l *Level) getNeighbors(x, y int) [][2]int {
	// Get the neighbors of the given hex tile, depending on the column.
	// On even columns, we need to use different neighbors than on odd columns.
	var nbs [][2]int
	if x%2 == 0 { // even column
		nbs = [][2]int{
			{x + 1, y},
			{x - 1, y},
			{x, y + 1},
			{x, y - 1},
			{x + 1, y - 1},
			{x - 1, y - 1},
		}
	} else { // odd column
		nbs = [][2]int{
			{x + 1, y},
			{x - 1, y},
			{x, y + 1},
			{x, y - 1},
			{x + 1, y + 1},
			{x - 1, y + 1},
		}
	}

	// Add the neighbors that are within the bounds of the level.
	var res [][2]int
	for _, p := range nbs {
		if x, y := p[0], p[1]; x >= 0 && x < l.Width && y >= 0 && y < l.Height {
			res = append(res, [2]int{x, y})
		}
	}
	return res
}

func (l *Level) drawHex(background *ebiten.Image, xCenter, yCenter, scale float64, x, y int, c color.Color) {
	// TODO: Allow drawing operations for translation, rotation and scaling.
	// Draw the hexagon in flat top version.
	//
	//  Flat top version:
	//            ____
	//          /      \
	//    ____ /        \ ____
	//  /      \        /      \
	// /        \ ____ /        \
	// \        /      \        /
	//  \ ____ /        \ ____ /
	hr := float64(l.hexRadius) * scale * 0.99
	points := []vectors.Vec2{
		{X: xCenter + hr, Y: yCenter},
		{X: xCenter + hr/2, Y: yCenter + math.Sqrt(3)/2*hr},
		{X: xCenter - hr/2, Y: yCenter + math.Sqrt(3)/2*hr},
		{X: xCenter - hr, Y: yCenter},
		{X: xCenter - hr/2, Y: yCenter - math.Sqrt(3)/2*hr},
		{X: xCenter + hr/2, Y: yCenter - math.Sqrt(3)/2*hr},
	}

	// Calculate gray value from tile value (minimum is 50% grey, maximum is white)
	colGray := uint8((int(l.Tiles[y*l.Width+x]) + 255) / 2)
	c = genColor(c, float64(colGray)/255)

	// Draw the hexagon.
	for i := 0; i < len(points); i++ {
		ebitenutil.DrawLine(background, points[i].X, points[i].Y, points[(i+1)%len(points)].X, points[(i+1)%len(points)].Y, c)
	}

	// Draw the tile coordinates.
	ebitenutil.DebugPrintAt(background, fmt.Sprintf("%d,%d", x, y), int(xCenter)-l.hexRadius/2, int(yCenter))
}

func genColor(baseColor color.Color, intensity float64) color.Color {
	r, g, b, a := baseColor.RGBA()
	return color.RGBA{
		R: uint8(float64(r) * intensity),
		G: uint8(float64(g) * intensity),
		B: uint8(float64(b) * intensity),
		A: uint8(a),
	}
}

type Tile byte

// TileType represents the type of a tile.
type TileType int

const (
	TileTypeGrass TileType = iota
	TileTypeDirt
	TileTypeWater
	TileTypeSnow
	TileTypeTrees
)

func (t Tile) Type() TileType {
	switch {
	case t < 20:
		return TileTypeWater
	case t < 40:
		return TileTypeDirt
	case t < 80:
		return TileTypeGrass
	default:
		return TileTypeSnow
	}
}

func (t Tile) HasTrees() bool {
	return t > 40 && t < 80
}
