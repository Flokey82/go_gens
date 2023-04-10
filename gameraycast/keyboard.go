package gameraycast

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) KeyboardHandler() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.Array[int((int(g.py)>>6))*g.X+int((int(g.px+g.pdx)>>6))] == 0 {
			g.px += g.pdx
		}
		if g.Array[int((int(g.py+g.pdy)>>6))*g.X+int((int(g.px)>>6))] == 0 {
			g.py += g.pdy
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.Array[int((int(g.py)>>6))*g.X+int((int(g.px-g.pdx)>>6))] == 0 {
			g.px -= g.pdx
		}
		if g.Array[int((int(g.py-g.pdy)>>6))*g.X+int((int(g.px)>>6))] == 0 {
			g.py -= g.pdy
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.pa -= 0.05
		if g.pa < 0 {
			g.pa = 2 * math.Pi
		}
		g.pdx = math.Cos(g.pa) * g.ps
		g.pdy = math.Sin(g.pa) * g.ps
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.pa += 0.05
		if g.pa > math.Pi*2 {
			g.pa = 0
		}
		g.pdx = math.Cos(g.pa) * g.ps
		g.pdy = math.Sin(g.pa) * g.ps
	}
}
