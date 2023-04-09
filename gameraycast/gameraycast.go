package gameraycast

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	windowHeight int     // Window Height
	windowWidth  int     // Window Width
	px           float64 // Player X Coordinate
	py           float64 // Player Y Coordinate
	pdx          float64 // Player delta X
	pdy          float64 // Player delta Y
	pa           float64 // Player Angle
	ps           float64 // Player Speed
	mapx         int     // Map Width
	mapy         int     // Map Height
	mapScale     int     // Map Unit Size
	mapArray     [64]int // Map Array
	drawRays     bool    // Draw Rays
	drawMap      bool    // Draw Map
}

func newGame() *Game {
	return &Game{
		windowHeight: 500,
		windowWidth:  720,
		px:           300, // Player X Coordinate
		py:           300, // Player Y Coordinate
		pdx:          1,   // Player delta X
		pdy:          1,   // Player delta Y
		pa:           0,   // Player Angle
		ps:           2,   // Player Speed
		mapx:         8,   // Map Width
		mapy:         8,   // Map Height
		mapScale:     64,  // Map Unit Size
		mapArray: [64]int{
			1, 1, 1, 1, 1, 1, 1, 1,
			1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 1, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1, 1,
		},
		drawRays: false,
		drawMap:  false,
	}
}

func (g *Game) Update() error {
	g.KeyboardHandler()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawMap {
		g.DrawMap(screen)
		ebitenutil.DrawLine(screen, g.px, g.py, g.px+g.pdx*5, g.py+g.pdy*5, color.White)
	}
	g.CastRays(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) CastRays(screen *ebiten.Image) {
	var (
		r    int
		mx   int
		my   int
		mp   int
		dof  int
		rx   float64
		ry   float64
		ra   float64
		xo   float64
		yo   float64
		disT float64
	)

	ra = g.pa - 0.00872665*30 // half a degree in radians times 30
	if ra < 0 {
		ra += 2 * math.Pi
	} else if ra > 2*math.Pi {
		ra -= 2 * math.Pi
	}

	for r = 0; r < 120; r++ {
		// Check horizontal lines
		dof = 0
		disH := 1000.0 // Distance to horizontal wall
		hx := g.px
		hy := g.py
		aTan := -1 / math.Tan(ra)

		if ra > math.Pi { // Looking up
			ry = float64(((int(g.py) >> 6) << 6)) - 0.0001
			rx = (g.py-ry)*aTan + g.px
			yo = float64(-1 * g.mapScale)
			xo = -1 * yo * aTan
		} else if ra < math.Pi { // Looking down
			ry = float64(((int(g.py) >> 6) << 6)) + 64.0
			rx = (g.py-ry)*aTan + g.px
			yo = float64(g.mapScale)
			xo = -1 * yo * aTan
		} else if ra == 0 || ra == math.Pi { // Looking straight left or right
			rx = g.px
			ry = g.py
			dof = 8
		}
		for dof < 8 {
			mx = (int(rx) >> 6)
			my = (int(ry) >> 6)
			mp = my*g.mapx + mx
			if mp > 0 && mp < g.mapx*g.mapy && g.mapArray[mp] > 0 {
				disH = Dist(g.px, g.py, rx, ry)
				hx = rx
				hy = ry
				dof = 8
			} else {
				rx = rx + xo
				ry = ry + yo
				dof += 1
			}
		}

		// Check vertical lines
		dof = 0
		var disV float64 = 1000 // Distance to vertical wall
		vx := g.px
		vy := g.py
		nTan := -math.Tan(ra)

		if ra > math.Pi/2 && ra < 3*math.Pi/2 { // Looking left
			rx = float64(((int(g.px) >> 6) << 6)) - 0.0001
			ry = (g.px-rx)*nTan + g.py
			xo = float64(-1 * g.mapScale)
			yo = -1 * xo * nTan
		} else if ra < math.Pi/2 || ra > 3*math.Pi/2 { // Looking right
			rx = float64(((int(g.px) >> 6) << 6)) + 64.0
			ry = (g.px-rx)*nTan + g.py
			xo = float64(g.mapScale)
			yo = -1 * xo * nTan
		} else if ra == math.Pi/2 || ra == 3*math.Pi/2 { // Looking up or down
			rx = g.px
			ry = g.py
			dof = 8
		}

		for dof < 8 {
			mx = (int(rx) >> 6)
			my = (int(ry) >> 6)
			mp = my*g.mapx + mx
			if mp > 0 && mp < g.mapx*g.mapy && g.mapArray[mp] > 0 {
				disV = Dist(g.px, g.py, rx, ry)
				vx = rx
				vy = ry
				dof = 8
			} else {
				rx = rx + xo
				ry = ry + yo
				dof += 1
			}
		}

		var isVertical bool
		if disH < disV {
			rx = hx
			ry = hy
			disT = disH
			isVertical = false
		} else if disV < disH {
			rx = vx
			ry = vy
			disT = disV
			isVertical = true
		}

		if g.drawRays {
			ebitenutil.DrawLine(screen, g.px, g.py, rx, ry, color.RGBA{255, 128, 0, 255})
		}

		ca := g.pa - ra
		if ca < 0 {
			ca += 2 * math.Pi
		} else if ca > 2*math.Pi {
			ca -= 2 * math.Pi
		}

		disT = disT * math.Cos(ca)
		lineH := float64(g.mapScale*g.windowHeight) / disT
		if lineH > float64(g.windowHeight) {
			lineH = float64(g.windowHeight)
		}

		if isVertical {
			ebitenutil.DrawRect(screen, float64(r*g.windowWidth/120), float64(g.windowHeight/2)-lineH/2, float64(g.windowWidth/120), lineH, color.RGBA{255, 0, 0, 255})
		} else {
			ebitenutil.DrawRect(screen, float64(r*g.windowWidth/120), float64(g.windowHeight/2)-lineH/2, float64(g.windowWidth/120), lineH, color.RGBA{204, 0, 0, 255})
		}

		ra += 0.00872665 // half a degree in radians
		if ra < 0 {
			ra += 2 * math.Pi
		} else if ra > 2*math.Pi {
			ra -= 2 * math.Pi
		}
	}
}

func Run() {
	g := newGame()
	ebiten.SetWindowSize(g.windowWidth, g.windowHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func Dist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}
