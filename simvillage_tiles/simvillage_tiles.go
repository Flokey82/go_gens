// Package simvillage_tiles is a tile based village simulation based on:
// https://ebiten.org/examples/tiles.html
// https://ebiten.org/examples/animation.html
package simvillage_tiles

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
)

const (
	screenWidth  = 240
	screenHeight = 240
	tileSize     = 16
	frameOX      = 0
	frameOY      = 32
	frameWidth   = 32
	frameHeight  = 32
	frameCount   = 8
)

var (
	tilesImage  *ebiten.Image
	runnerImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)

	// Decode an image from the image file's byte slice.
	imgRunner, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(imgRunner)
}

type Game struct {
	layers    [][]int
	player    *Creature
	creatures []*Creature
}

func NewGame() *Game {
	g := &Game{
		layers: defaultLayers(),
	}
	g.player = NewCreature(g, [2]int{0, 0})

	// Add the player to the creature index.
	g.addCreature(g.player)

	// Add some creatures.
	g.addCreature(NewCreature(g, [2]int{2, 2}))
	g.addCreature(NewCreature(g, [2]int{20, 2}))
	g.addCreature(NewCreature(g, [2]int{2, 20}))
	return g
}

func (g *Game) Update() error {
	// Handle input.
	var posDelta [2]int
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		posDelta[0]--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		posDelta[0]++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		posDelta[1]--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		posDelta[1]++
	}
	g.player.move(posDelta)

	// Handle "AI".
	// NOTE: This just makes the creatures move randomly.
	for _, c := range g.creatures {
		if c == g.player {
			continue // Skip the player.
		}
		c.move([2]int{rand.Intn(3) - 1, rand.Intn(3) - 1})
	}
	return nil
}

// canEnter returns whether the player can enter the tile at (x, y) in the chunk (cX, cY).
func (g *Game) canEnter(cX, cY, newX, newY int) bool {
	x, y := getTileXYFromPos(newX, newY)
	if x < 0 || x >= screenWidth/tileSize || y < 0 || y >= screenHeight/tileSize {
		return false
	}
	layers := g.getChunk(cX, cY)
	return layers[2][y*screenWidth/tileSize+x] == 0
}

// getViewportXY returns the x, y tile position of the top left corner of the viewport in the tile map.
func (g *Game) getViewportXY() (int, int) {
	return g.player.pos[0] / tileSize, g.player.pos[1] / tileSize
}

// getChunk returns the layers at the given chunk position.
// Note: Right now we're generating chunks on the fly... We should find a way to cache them.
func (g *Game) getChunk(x, y int) [][]int {
	if x != 0 || y != 0 {
		return genChunk(x, y, screenWidth/tileSize, screenHeight/tileSize)
	}
	return g.layers
}

// addCreature adds a creature to the game.
func (g *Game) addCreature(c *Creature) {
	g.creatures = append(g.creatures, c)
}

// getCreatures returns the creatures at the given tile in the given chunk.
func (g *Game) getCreatures(tileIdx int, chunk [2]int) []*Creature {
	var creatures []*Creature
	for _, c := range g.creatures {
		if c.tileIdx == tileIdx && c.chunk == chunk {
			creatures = append(creatures, c)
		}
	}
	return creatures
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	const xCount = screenWidth / tileSize
	w, _ := tilesImage.Size()
	tileXCount := w / tileSize

	// Iterate through the layers and draw them.
	var cxy [2]int

	// Viewport chunk.
	vpChunk := g.player.chunk

	// Viewport position.
	vpPos := g.player.pos

	drawChunk := func(x, y int) {
		layers := g.getChunk(x+vpChunk[0], y+vpChunk[1])
		// Calculate the offset of the current chunk relative
		// to the player position / the center of the screen.
		// TODO: Also tidy up the chunk offset transformation.
		cOffsX := float64(x*screenWidth) - float64(vpPos[0])
		cOffsY := float64(y*screenHeight) - float64(vpPos[1])
		cxy[0] = x + vpChunk[0]
		cxy[1] = y + vpChunk[1]
		for i, l := range layers {
			drawP := i == len(layers)-1
			// TODO: Draw only the visible tiles and fetch them according to
			// the player's position and the viewport.
			for i, t := range l {
				// Draws the given layer on the screen.
				// Check if we should draw any creatures on the current layer and tile.
				// TODO: Optimize this.
				if drawP {
					for _, cr := range g.getCreatures(i, cxy) {
						cr.Draw(screen, cOffsX, cOffsY)
					}
				}

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64((i%xCount)*tileSize)+cOffsX, float64((i/xCount)*tileSize)+cOffsY)

				// Get the right tile sprite.
				sx := (t % tileXCount) * tileSize
				sy := (t / tileXCount) * tileSize
				screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
			}
		}
	}

	// Render current and sourrounding chunks.
	for _, dir := range [][2]int{
		{1, 0},   // East
		{0, 1},   // South
		{1, 1},   // SouthEast
		{0, -1},  // North
		{1, -1},  // NorthEast
		{-1, -1}, // NorthWest
		{-1, 0},  // West
		{-1, 1},  // SouthWest
		{0, 0},   // Center
	} {
		drawChunk(dir[0], dir[1])
	}

	// Draw some information on the screen.
	px, py := g.player.getXY()                     // Current player tile
	vx, vy := g.getViewportXY()                    // Current viewport tile
	cx, cy := g.player.chunk[0], g.player.chunk[1] // Current chunk
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (T %d, %d V %d, %d C %d, %d)", ebiten.ActualTPS(), px, py, vx, vy, cx, cy))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
