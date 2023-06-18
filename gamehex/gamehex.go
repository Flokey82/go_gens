// Package gamehex implements an example hexagonal game board and uses ebiten for rendering.
package gamehex

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Game struct {
	currentLevel *Level
	width        int
	height       int
	offscreen    *ebiten.Image
	clickedTile  [2]int
	camX         float64
	camY         float64
	camScale     float64
	camScaleTo   float64
	mousePanX    int
	mousePanY    int
	sprites      *SpriteSheet
	pathFinder
}

// NewGame returns a new isometric demo Game.
func NewGame() (*Game, error) {
	l, err := NewLevel(10, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to create new level: %s", err)
	}
	s, err := LoadSpriteSheet(l.hexRadius)
	if err != nil {
		return nil, fmt.Errorf("failed to load spritesheet: %s", err)
	}

	return &Game{
		currentLevel: l,
		clickedTile:  [2]int{-1, -1},
		camScale:     1,
		camScaleTo:   1,
		mousePanX:    math.MinInt32,
		mousePanY:    math.MinInt32,
		sprites:      s,
		pathFinder:   newPathFinder(),
	}, nil
}

// Update reads current user input and updates the Game state.
func (g *Game) Update() error {
	// Update target zoom level.
	var scrollY float64
	if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		scrollY = -0.25
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		scrollY = 0.25
	} else {
		_, scrollY = ebiten.Wheel()
		if scrollY < -1 {
			scrollY = -1
		} else if scrollY > 1 {
			scrollY = 1
		}
	}
	g.camScaleTo += scrollY * (g.camScaleTo / 7)

	// Clamp target zoom level.
	if g.camScaleTo < 0.01 {
		g.camScaleTo = 0.01
	} else if g.camScaleTo > 100 {
		g.camScaleTo = 100
	}

	// Smooth zoom transition.
	div := 10.0
	if g.camScaleTo > g.camScale {
		g.camScale += (g.camScaleTo - g.camScale) / div
	} else if g.camScaleTo < g.camScale {
		g.camScale -= (g.camScale - g.camScaleTo) / div
	}

	// Pan camera via keyboard.
	pan := 7.0 / g.camScale
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camX -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camX += pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camY -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camY += pan
	}

	// Select tiles via keyboard.
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.clickedTile[0]--
		if g.clickedTile[0] < 0 {
			g.clickedTile[0] = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.clickedTile[0]++
		if g.clickedTile[0] >= g.currentLevel.Width {
			g.clickedTile[0] = g.currentLevel.Width - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.clickedTile[1]++
		if g.clickedTile[1] >= g.currentLevel.Height {
			g.clickedTile[1] = g.currentLevel.Height - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.clickedTile[1]--
		if g.clickedTile[1] < 0 {
			g.clickedTile[1] = 0
		}
	}

	// Pan camera via mouse.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if g.mousePanX == math.MinInt32 && g.mousePanY == math.MinInt32 {
			g.mousePanX, g.mousePanY = ebiten.CursorPosition()
		} else {
			x, y := ebiten.CursorPosition()
			dx, dy := float64(g.mousePanX-x)*(pan/100), float64(g.mousePanY-y)*(pan/100)
			g.camX, g.camY = g.camX-dx, g.camY+dy
		}
	} else if g.mousePanX != math.MinInt32 || g.mousePanY != math.MinInt32 {
		g.mousePanX, g.mousePanY = math.MinInt32, math.MinInt32
	}

	// Clamp camera position.
	// TODO: Fix actual width and height. Clamping is currently inaccurate.
	worldWidth := float64(g.currentLevel.Width*g.currentLevel.hexRadius) * 1.5
	worldHeight := float64(g.currentLevel.Height*g.currentLevel.hexRadius) * math.Sqrt(3)
	if g.camX < 0 {
		g.camX = 0
	} else if g.camX > worldWidth {
		g.camX = worldWidth
	}
	if g.camY < -worldHeight {
		g.camY = -worldHeight
	} else if g.camY > 0 {
		g.camY = 0
	}

	// Handle pathfinding input.
	g.handlePathInput()

	// If we have a mouse click, we calculate the tile we clicked on and store it
	// for rendering later.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.clickedTile = g.getTileAtCursor()
	}
	return nil
}

func (g *Game) getTileAtCursor() [2]int {
	x, y := ebiten.CursorPosition()
	x, y = int(float64(x)/g.camScale), int(float64(y)/g.camScale)
	x, y = x+int(g.camX), y-int(g.camY)
	cX, cY := g.currentLevel.TileAtPos(x, y)
	return [2]int{cX, cY}
}

// Draw draws the Game on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Render level.
	g.renderLevel(screen)

	// Render path info.
	g.drawPathInfo(screen)

	// Print game info.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC P\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f", ebiten.ActualFPS(), ebiten.ActualTPS(), g.camScale, g.camX, g.camY))
}

// Layout is called when the Game's layout changes.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width, g.height = outsideWidth, outsideHeight
	return g.width, g.height
}

// renderLevel draws the current Level on the screen.
func (g *Game) renderLevel(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	padding := float64(g.currentLevel.hexRadius) * g.camScale
	cx, cy := 0.0, 0.0 // float64(g.width/2), float64(g.height/2)

	scaleLater := g.camScale > 1
	scale := g.camScale
	target := screen

	// When zooming in, tiles can have slight bleeding edges.
	// To avoid them, render the result on an offscreen first and then scale it later.
	if scaleLater {
		if g.offscreen != nil {
			w, h := g.offscreen.Size()
			sw, sh := screen.Size()
			if w != sw || h != sh {
				g.offscreen.Dispose()
				g.offscreen = nil
			}
		}
		if g.offscreen == nil {
			g.offscreen = ebiten.NewImage(screen.Size())
		}
		target = g.offscreen
		target.Clear()
		scale = 1
	}

	nbOfSelected := make(map[[2]int]bool)
	if g.clickedTile[0] != -1 && g.clickedTile[1] != -1 {
		for _, n := range g.currentLevel.getNeighbors(g.clickedTile[0], g.clickedTile[1]) {
			nbOfSelected[n] = true
		}
	}

	// Hacky way to check if a tile is part of the path.
	isPath := make(map[[2]int]bool)
	for _, p := range g.path {
		isPath[p] = true
	}

	for y := 0; y < g.currentLevel.Height; y++ {
		for x := 0; x < g.currentLevel.Width; x++ {
			// Tint the tile red if it was clicked.
			var c color.Color
			if g.clickedTile[0] == x && g.clickedTile[1] == y {
				c = color.RGBA{255, 0, 0, 255}
			} else if nbOfSelected[[2]int{x, y}] {
				c = color.RGBA{0, 255, 0, 255}
			} else {
				c = color.White
			}
			xi, yi := g.currentLevel.HexTileXYToPixelPos(x, y)

			// Skip drawing tiles that are out of the screen.
			drawX, drawY := ((float64(xi)-g.camX)*g.camScale)+cx, ((float64(yi)+g.camY)*g.camScale)+cy
			if drawX+padding < 0 || drawY+padding < 0 || drawX > float64(g.width) || drawY > float64(g.height) {
				continue
			}

			op.GeoM.Reset()
			// Move to current isometric position.
			op.GeoM.Translate(float64(xi), float64(yi))
			// Translate camera position.
			op.GeoM.Translate(-g.camX, g.camY)
			// Translate from center of hexagon.
			op.GeoM.Translate(-float64(g.currentLevel.hexRadius), -float64(g.currentLevel.hexRadius)*math.Sqrt(3)/2)
			// Zoom.
			op.GeoM.Scale(scale, scale)
			// Center.
			op.GeoM.Translate(cx, cy)

			// Reset color matrix and tint if needed.
			op.ColorM.Reset()
			// Tint if it is part of the path.
			if isPath[[2]int{x, y}] {
				op.ColorM.Scale(1, 0.5, 0.5, 1)
			}

			// Draw tile.
			val := g.currentLevel.Tiles[y*g.currentLevel.Width+x]

			// Draw terrain.
			switch val.Type() {
			case TileTypeWater:
				target.DrawImage(g.sprites.Water, op)
			case TileTypeDirt:
				target.DrawImage(g.sprites.Dirt, op)
			case TileTypeGrass:
				target.DrawImage(g.sprites.Grass, op)
			default:
				target.DrawImage(g.sprites.Snow, op)
			}

			// Draw trees if present.
			if val.HasTrees() {
				target.DrawImage(g.sprites.Trees, op)
			}

			g.currentLevel.drawHex(target, (float64(xi)-g.camX+cx)*scale, (float64(yi)+g.camY+cy)*scale, scale, x, y, c)
		}
	}

	if scaleLater {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Scale(float64(g.camScale), float64(g.camScale))
		op.GeoM.Translate(cx, cy)
		screen.DrawImage(target, op)
	}
}
