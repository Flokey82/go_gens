package gameraycast

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) DrawMap(screen *ebiten.Image) {
	var x, y int
	var xo, yo int

	for y = 0; y < g.Y; y++ {
		for x = 0; x < g.X; x++ {
			var tileColor = color.Gray16{0xffff}
			if g.Array[y*g.X+x] == 1 {
				tileColor = color.Gray16{0xffff}
			} else {
				tileColor = color.Gray16{0}
			}

			xo = x * g.Scale
			yo = y * g.Scale

			ebitenutil.DrawRect(screen, float64(xo+1), float64(yo+1), float64(g.Scale-1), float64(g.Scale-1), tileColor)
		}
	}
}
