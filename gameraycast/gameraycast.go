package gameraycast

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Map represents the map.
type Map struct {
	X     int   // Map Width
	Y     int   // Map Height
	Array []int // Map Array
}

// NewMap creates a new Map.
func NewMap() *Map {
	return &Map{
		X: 8,
		Y: 8,
		Array: []int{
			1, 1, 1, 1, 1, 1, 1, 1,
			1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 2, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1, 1,
		},
	}
}

// Game represents the game.
type Game struct {
	windowHeight int     // Window Height
	windowWidth  int     // Window Width
	px           float64 // Player X Coordinate
	py           float64 // Player Y Coordinate
	pdx          float64 // Player delta X
	pdy          float64 // Player delta Y
	pa           float64 // Player Angle
	ps           float64 // Player Speed
	pFov         float64 // Player Field of View
	rayPrecision int     // Number of rays to cast
	Scale        int     // Map Scale (how many pixels per tile)
	drawRays     bool    // Draw Rays
	drawMap      bool    // Draw Map
	drawTextures bool    // Draw Texture
	useAltCast   bool    // Use alternative ray casting
	*Map                 // Map
}

// NewGame creates a new Game
func NewGame(m *Map) *Game {
	g := &Game{
		windowHeight: 500,
		windowWidth:  720,
		pdx:          1,   // Player delta X
		pdy:          1,   // Player delta Y
		pa:           0,   // Player Angle
		ps:           2,   // Player Speed
		pFov:         60,  // Player Field of View
		rayPrecision: 120, // Number of rays to cast
		Scale:        64,
		drawRays:     false,
		drawMap:      false,
		drawTextures: true,
		useAltCast:   true,
		Map:          m,
	}

	// Find a valid start position for the player.
	var px, py float64
	for i := 0; i < len(m.Array); i++ {
		if m.Array[i] == 0 {
			px = float64(i%m.X) * float64(g.Scale)
			py = float64(i/m.X) * float64(g.Scale)
			break
		}
	}
	g.px = px // Player X Coordinate
	g.py = py // Player Y Coordinate

	return g
}

// Run starts the game.
func (g *Game) Run() {
	ebiten.SetWindowSize(g.windowWidth, g.windowHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// Update updates the game state (keyboard input).
func (g *Game) Update() error {
	g.KeyboardHandler()
	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawMap {
		g.DrawMap(screen)
		ebitenutil.DrawLine(screen, g.px, g.py, g.px+g.pdx*5, g.py+g.pdy*5, color.White)
	}
	g.CastRays(screen)
}

// Layout returns the screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// CastRays casts rays from the player to the map.
func (g *Game) CastRays(screen *ebiten.Image) {
	if g.useAltCast {
		g.rayCastingB(screen)
	} else {
		g.rayCastingA(screen)
	}
}

const degToRad = math.Pi / 180.0

// rayCastingA casts rays from the player to the map.
func (g *Game) rayCastingA(screen *ebiten.Image) {
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

	fovHalf := g.pFov / 2.0
	fovStart := degToRad * fovHalf
	fovStep := degToRad * g.pFov / float64(g.rayPrecision)

	// Keep angle between 0 and 2PI.
	ra = g.pa - fovStart // start at half fov to the left of the player
	if ra < 0 {
		ra += 2 * math.Pi
	} else if ra > 2*math.Pi {
		ra -= 2 * math.Pi
	}

	// Check if we are looking at a door.
	doorInSight := -1

	for r = 0; r < g.rayPrecision; r++ {
		// Check horizontal lines
		dof = 0
		disH := 1000.0 // Distance to horizontal wall
		hx := g.px
		hy := g.py
		aTan := -1 / math.Tan(ra)

		if ra > math.Pi { // Looking up
			ry = float64(((int(g.py) >> 6) << 6)) - 0.0001
			rx = (g.py-ry)*aTan + g.px
			yo = float64(-1 * g.Scale)
			xo = -1 * yo * aTan
		} else if ra < math.Pi { // Looking down
			ry = float64(((int(g.py) >> 6) << 6)) + float64(g.Scale)
			rx = (g.py-ry)*aTan + g.px
			yo = float64(g.Scale)
			xo = -1 * yo * aTan
		} else if ra == 0 || ra == math.Pi { // Looking straight left or right
			rx = g.px
			ry = g.py
			dof = 8
		}
		for dof < 8 {
			mx = (int(rx) >> 6)
			my = (int(ry) >> 6)
			mp = my*g.X + mx
			if mp > 0 && mp < g.X*g.Y && g.Array[mp] > 0 {
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
		disV := 1000.0 // Distance to vertical wall
		vx := g.px
		vy := g.py
		nTan := -math.Tan(ra)

		if ra > math.Pi/2 && ra < 3*math.Pi/2 { // Looking left
			rx = float64(((int(g.px) >> 6) << 6)) - 0.0001
			ry = (g.px-rx)*nTan + g.py
			xo = float64(-1 * g.Scale)
			yo = -1 * xo * nTan
		} else if ra < math.Pi/2 || ra > 3*math.Pi/2 { // Looking right
			rx = float64(((int(g.px) >> 6) << 6)) + float64(g.Scale)
			ry = (g.px-rx)*nTan + g.py
			xo = float64(g.Scale)
			yo = -1 * xo * nTan
		} else if ra == math.Pi/2 || ra == 3*math.Pi/2 { // Looking up or down
			rx = g.px
			ry = g.py
			dof = 8
		}

		for dof < 8 {
			mx = (int(rx) >> 6)
			my = (int(ry) >> 6)
			mp = my*g.X + mx
			if mp > 0 && mp < g.X*g.Y && g.Array[mp] > 0 {
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
		// Keep angle between 0 and 2PI.
		if ca < 0 {
			ca += 2 * math.Pi
		} else if ca > 2*math.Pi {
			ca -= 2 * math.Pi
		}

		disT = disT * math.Cos(ca)
		lineH := float64(g.Scale*g.windowHeight) / disT
		if lineH > float64(g.windowHeight) {
			lineH = float64(g.windowHeight)
		}

		// Find the type of tile we hit with the ray.
		var tileType int
		{
			mx = (int(rx) >> 6)
			my = (int(ry) >> 6)
			mp = my*g.X + mx
			if mp > 0 && mp < g.X*g.Y {
				tileType = g.Array[mp]
			}
		}

		// Set the color of the wall.
		col := getWallColor(tileType)
		if !isVertical {
			col = darkenColor(col, 0.7) // Horizontal walls are darker.
		}

		// TODO: Deduplicate this code!
		// This is the slice of wall that we are going to draw.
		xPos := float64(r * g.windowWidth / g.rayPrecision)
		yPos := float64(g.windowHeight/2) - lineH/2
		width := float64(g.windowWidth / g.rayPrecision)
		height := lineH

		// Draw the floor, which is a rectangle that goes from the bottom of the
		// screen to the bottom of the wall.
		ebitenutil.DrawRect(screen, xPos, yPos+height, width, float64(g.windowHeight), floorColor)

		if !g.drawTextures {
			ebitenutil.DrawRect(screen, xPos, yPos, width, height, col)
		} else {
			// Get texture
			tex := getWallTexture(tileType)

			// Calcule texture position
			texturePositionX := (int(float64(tex.Width())*(rx+ry)/float64(g.Scale)) % tex.Width())

			yIncrementer := (height * 2) / float64(tex.Height())
			yD := float64(g.windowHeight)/2 - height

			for i := 0; i < tex.Height(); i++ {
				ebitenutil.DrawRect(screen, xPos, yD, width, yIncrementer, tex.At(texturePositionX, i))
				yD += yIncrementer
			}
		}

		// TODO: Formalize opening doors and add an animation of the door opening.
		if r == g.rayPrecision/2 && tileType == 2 {
			doorInSight = int(ry/float64(g.Scale))*g.X + int(rx/float64(g.Scale))
		}

		// Keep angle between 0 and 2PI.
		ra += fovStep // render angle += fovStep
		if ra < 0 {
			ra += 2 * math.Pi
		} else if ra > 2*math.Pi {
			ra -= 2 * math.Pi
		}
	}

	// TODO: Move this to a separate function.
	if doorInSight >= 0 {
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			g.Array[doorInSight] = 0
		} else {
			// Print that we can open the door if we press E.
			ebitenutil.DebugPrint(screen, "Press E to open the door")
		}
	}
}

// Raycasting logic
// Alternative algorithm from:
// https://github.com/vinibiavatti1/RayCastingTutorial/blob/master/basic/raycasting.js
func (g *Game) rayCastingB(screen *ebiten.Image) {
	// Calculate the angle of the first ray
	fovHalf := g.pFov / 2.0
	fovStart := degToRad * fovHalf

	// Keep angle between 0 and 2PI
	ra := g.pa - fovStart
	if ra < 0 {
		ra += 2 * math.Pi
	} else if ra > 2*math.Pi {
		ra -= 2 * math.Pi
	}

	// Check if we are looking at a door.
	doorInSight := -1

	fovStep := degToRad * float64(g.pFov) / float64(g.rayPrecision)
	for r := 0; r < g.rayPrecision; r++ {
		// Ray data
		rx := g.px
		ry := g.py

		// Ray path incrementers
		rayCos := math.Cos(ra) / float64(g.rayPrecision)
		raySin := math.Sin(ra) / float64(g.rayPrecision)

		// Wall finder
		var wall int
		for wall == 0 {
			rx += rayCos
			ry += raySin
			wall = g.Array[int(ry/float64(g.Scale))*g.X+int(rx/float64(g.Scale))]
		}

		// Pythagoras theorem
		disH := math.Sqrt(math.Pow(g.px-rx, 2) + math.Pow(g.py-ry, 2))

		// Fish eye correction
		disH = disH * math.Cos(ra-g.pa)

		// Wall height
		lineH := float64(g.Scale) * float64(g.windowHeight/2) / disH

		// Set the color of the wall.
		col := getWallColor(wall)
		if wall == 1 {
			col = darkenColor(col, 0.7) // Horizontal walls are darker.
		}

		// TODO: Deduplicate this code!
		// This is the slice of wall that we are going to draw.
		xPos := float64(r * g.windowWidth / g.rayPrecision)
		yPos := float64(g.windowHeight/2) - lineH/2
		width := float64(g.windowWidth / g.rayPrecision)
		height := lineH

		// Draw the floor, which is a rectangle that goes from the bottom of the
		// screen to the bottom of the wall.
		ebitenutil.DrawRect(screen, xPos, yPos+height, width, float64(g.windowHeight), floorColor)

		if !g.drawTextures {
			ebitenutil.DrawRect(screen, xPos, yPos, width, height, col)
		} else {
			// Get texture
			tex := getWallTexture(wall)

			// Calcule texture position
			texturePositionX := (int(float64(tex.Width())*(rx+ry)/float64(g.Scale)) % tex.Width())

			yIncrementer := (height * 2) / float64(tex.Height())
			yD := float64(g.windowHeight)/2 - height

			for i := 0; i < tex.Height(); i++ {
				ebitenutil.DrawRect(screen, xPos, yD, width, yIncrementer, tex.At(texturePositionX, i))
				yD += yIncrementer
			}
		}

		// TODO: Formalize opening doors and add an animation of the door opening.
		if r == g.rayPrecision/2 && wall == 2 {
			doorInSight = int(ry/float64(g.Scale))*g.X + int(rx/float64(g.Scale))
		}

		// Keep angle between 0 and 2PI
		ra += fovStep // render angle += fovStep
		if ra < 0 {
			ra += 2 * math.Pi
		} else if ra > 2*math.Pi {
			ra -= 2 * math.Pi
		}
	}

	// TODO: Move this to a separate function.
	if doorInSight >= 0 {
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			g.Array[doorInSight] = 0
		} else {
			// Print that we can open the door if we press E.
			ebitenutil.DebugPrint(screen, "Press E to open the door")
		}
	}
}

func Dist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}
