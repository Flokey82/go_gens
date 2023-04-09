package gameraycast

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) DrawMap(screen *ebiten.Image) {
	var x, y int
	var xo, yo int

	for y = 0; y < g.mapy; y++ {
		for x = 0; x < g.mapx; x++ {
			var tileColor = color.Gray16{0xffff}
			if g.mapArray[y*g.mapx+x] == 1 {
				tileColor = color.Gray16{0xffff}
			} else {
				tileColor = color.Gray16{0}
			}

			xo = x * g.mapScale
			yo = y * g.mapScale

			ebitenutil.DrawRect(screen, float64(xo+1), float64(yo+1), float64(g.mapScale-1), float64(g.mapScale-1), tileColor)
		}
	}
}
