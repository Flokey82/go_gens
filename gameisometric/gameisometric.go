// Copyright 2021 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gameisometric

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game is an isometric demo game.
type Game struct {
	currentLevel *Level
	width        int
	height       int
	camX         float64
	camY         float64
	camScale     float64
	camScaleTo   float64
	mousePanX    int
	mousePanY    int
	offscreen    *ebiten.Image
}

// NewGame returns a new isometric demo Game.
func NewGame() (*Game, error) {
	l, err := NewLevel()
	if err != nil {
		return nil, fmt.Errorf("failed to create new level: %s", err)
	}

	return &Game{
		currentLevel: l,
		camScale:     1,
		camScaleTo:   1,
		mousePanX:    math.MinInt32,
		mousePanY:    math.MinInt32,
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
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camX -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camX += pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camY -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camY += pan
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
	worldWidth := float64(g.currentLevel.width * g.currentLevel.tileSize / 2)
	worldHeight := float64(g.currentLevel.height * g.currentLevel.tileSize / 2)
	if g.camX < -worldWidth {
		g.camX = -worldWidth
	} else if g.camX > worldWidth {
		g.camX = worldWidth
	}
	if g.camY < -worldHeight {
		g.camY = -worldHeight
	} else if g.camY > 0 {
		g.camY = 0
	}

	// Randomize level.
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		l, err := NewLevel()
		if err != nil {
			return fmt.Errorf("failed to create new level: %s", err)
		}

		g.currentLevel = l
	}

	// Move level up or down.
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.currentLevel.LevelDown()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.currentLevel.LevelUp()
	}

	// If we click, print the tile we clicked on.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := getTileXY(g)
		fmt.Printf("Clicked on tile %d,%d\n", x, y)
		// TODO: Add pathfinding from player to tile.
	}

	return nil
}

// Draw draws the Game on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Render level.
	g.renderLevel(screen)

	// Print game info.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC FV R\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f\nLVL  %d", ebiten.ActualFPS(), ebiten.ActualTPS(), g.camScale, g.camX, g.camY, g.currentLevel.currentLevel))
}

// Layout is called when the Game's layout changes.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width, g.height = outsideWidth, outsideHeight
	return g.width, g.height
}

// cartesianToIso transforms cartesian coordinates into isometric coordinates.
func (g *Game) cartesianToIso(x, y float64) (float64, float64) {
	tileSize := g.currentLevel.tileSize
	ix := (x - y) * float64(tileSize/2)
	iy := (x + y) * float64(tileSize/4)
	return ix, iy
}

// This function might be useful for those who want to modify this example.

// isoToCartesian transforms isometric coordinates into cartesian coordinates.
func (g *Game) isoToCartesian(x, y float64) (float64, float64) {
	tileSize := g.currentLevel.tileSize
	cx := (x/float64(tileSize/2) + y/float64(tileSize/4)) / 2
	cy := (y/float64(tileSize/4) - (x / float64(tileSize/2))) / 2
	return cx, cy
}

// renderLevel draws the current Level on the screen.
func (g *Game) renderLevel(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	padding := float64(g.currentLevel.tileSize) * g.camScale
	cx, cy := float64(g.width/2), float64(g.height/2)

	scaleLater := g.camScale > 1
	target := screen
	scale := g.camScale

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

	for y := 0; y < g.currentLevel.height; y++ {
		for x := 0; x < g.currentLevel.width; x++ {
			xi, yi := g.cartesianToIso(float64(x), float64(y))

			// Skip drawing tiles that are out of the screen.
			drawX, drawY := ((xi-g.camX)*g.camScale)+cx, ((yi+g.camY)*g.camScale)+cy
			if drawX+padding < 0 || drawY+padding < 0 || drawX > float64(g.width) || drawY > float64(g.height) {
				continue
			}

			t := g.currentLevel.tiles[y][x]
			if t == nil {
				continue // No tile at this position.
			}

			op.GeoM.Reset()
			// Move to current isometric position.
			op.GeoM.Translate(xi, yi)
			// Translate camera position.
			op.GeoM.Translate(-g.camX, g.camY)
			// Zoom.
			op.GeoM.Scale(scale, scale)
			// Center.
			op.GeoM.Translate(cx, cy)

			t.Draw(target, op)
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

func getTileXY(g *Game) (int, int) {
	//Get the cursor position
	mx, my := ebiten.CursorPosition()
	winWidth, winHeight := ebiten.WindowSize()
	//Offset for center
	fmx := float64(mx) - float64(winWidth)/2.0
	fmy := float64(my) - float64(winHeight)/2.0

	x, y := fmx+g.camX, fmy-g.camY

	tileSize := g.currentLevel.tileSize

	//Do a half tile mouse shift because of our perspective
	x -= .5 * float64(tileSize)
	y -= .5 * float64(tileSize)
	//Convert isometric
	imx, imy := g.isoToCartesian(x, y)

	return int(math.Ceil(float64(imx))), int(math.Ceil(float64(imy)))
}
