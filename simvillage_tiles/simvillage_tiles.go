// Package simvillage_tiles is a tile based village simulation based on:
// https://ebiten.org/examples/tiles.html
// https://ebiten.org/examples/animation.html
package simvillage_tiles

import (
	"bytes"
	"fmt"
	"image"
	"log"

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
	layers          [][]int
	curChunk        [2]int
	playerPos       [2]int
	playerLooksLeft bool
	count           int
}

func NewGame() *Game {
	return &Game{
		layers: defaultLayers(),
	}
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

	// Check if we are moving in any direction.
	if keyPressed := (posDelta[0] != 0 || posDelta[1] != 0); !keyPressed {
		return nil // Nothing to do.
	}

	// Do we want to move left or right?
	g.playerLooksLeft = posDelta[0] < 0

	newPosX := g.playerPos[0] + posDelta[0]
	newPosY := g.playerPos[1] + posDelta[1]
	var newChunkDelta [2]int

	// TODO: If we leave the current chunk, update where we are.
	newChunkDelta[0] = newPosX / (screenWidth / 2)
	newPosX -= (screenWidth - 1) * newChunkDelta[0] // Can this be done better?
	newChunkDelta[1] = newPosY / (screenHeight / 2)
	newPosY -= (screenHeight - 1) * newChunkDelta[1] // Can this be done better?

	// Now check if we can go where we'd end up.
	// TODO: Check if we can enter newPosX, newPosY in the new chunk.
	// NOTE: This doesn't check the chunk yet, but since all chunks are the same it sorta works.
	if g.canEnter(newPosX, newPosY) {
		// TODO: If only one of both directions would be an illegal move,
		// should we just move along one axis?
		g.playerPos[0] = newPosX
		g.playerPos[1] = newPosY
		g.curChunk[0] += newChunkDelta[0]
		g.curChunk[1] += newChunkDelta[1]
	}
	g.count++
	return nil
}

// canEnter returns whether the player can enter the tile at (x, y).
func (g *Game) canEnter(newX, newY int) bool {
	x, y := getTileXYFromPlayerPos(newX, newY)
	if x < 0 || x >= screenWidth/tileSize || y < 0 || y >= screenHeight/tileSize {
		return false
	}
	return g.layers[2][y*screenWidth/tileSize+x] == 0
}

func getTileXYFromPlayerPos(x, y int) (int, int) {
	return (x + screenWidth/2) / tileSize, (y + screenHeight/2) / tileSize
}

// getPlayerXY returns the x, y tile position of the player in the tile map.
func (g *Game) getPlayerXY() (int, int) {
	return getTileXYFromPlayerPos(g.playerPos[0], g.playerPos[1])
}

// getViewportXY returns the x, y tile position of the top left corner of the viewport in the tile map.
func (g *Game) getViewportXY() (int, int) {
	return g.playerPos[0] / tileSize, g.playerPos[1] / tileSize
}

func (g *Game) getChunk(x, y int) [][]int {
	return g.layers // Note: Right now we're faking loading in chunks.
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draws the player on the screen.
	drawPlayer := func(screen *ebiten.Image) {
		op := &ebiten.DrawImageOptions{}
		// Center the frame of the player.
		op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		if g.playerLooksLeft {
			op.GeoM.Scale(-1, 1) // Mirror horizontally.
		}
		// Center the player on the screen.
		op.GeoM.Translate(screenWidth/2, screenHeight/2)

		// Get the animation frame.
		i := (g.count / 5) % frameCount
		sx, sy := frameOX+i*frameWidth, frameOY
		screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

		// Wrap around the counter so we don't overflow.
		g.count = g.count % (frameCount * 5)
	}

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	const xCount = screenWidth / tileSize
	w, _ := tilesImage.Size()
	tileXCount := w / tileSize

	// Get the player location tile index.
	pX, pY := g.getPlayerXY()
	pIdx := pY*xCount + pX

	// Draws the given layer on the screen.
	drawLayer := func(l []int, drawP bool, screen *ebiten.Image, cOffsX, cOffsY float64) {
		// TODO: Draw only the visible tiles and fetch them according to
		// the player's position and the viewport.
		for i, t := range l {
			// Check if we should draw the player on the current layer.
			// If the player is on the current tile, draw the player.
			if drawP && pIdx == i {
				drawPlayer(screen)
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xCount)*tileSize)+cOffsX, float64((i/xCount)*tileSize)+cOffsY)

			// Get the right tile sprite.
			sx := (t % tileXCount) * tileSize
			sy := (t / tileXCount) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

	// Iterate through the layers and draw them.
	drawChunk := func(x, y int) {
		layers := g.getChunk(x+g.curChunk[0], y+g.curChunk[1])
		// Calculate the offset of the current chunk relative
		// to the player position / the center of the screen.
		// TODO: Also tidy up the chunk offset transformation.
		cOffsX := float64(x*screenWidth) - float64(g.playerPos[0])
		cOffsY := float64(y*screenHeight) - float64(g.playerPos[1])
		for i, l := range layers {
			drawLayer(l, i == len(layers)-1 && x == 0 && y == 0, screen, cOffsX, cOffsY)
		}
	}

	// Render current and sourrounding chunks.
	for _, dir := range directions {
		drawChunk(dir[0], dir[1])
	}

	// Draw some information on the screen.
	px, py := g.getPlayerXY()              // Current player tile
	vx, vy := g.getViewportXY()            // Current viewport tile
	cx, cy := g.curChunk[0], g.curChunk[1] // Current chunk
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (T %d, %d V %d, %d C %d, %d)", ebiten.ActualTPS(), px, py, vx, vy, cx, cy))
}

var directions [][2]int = [][2]int{
	{1, 0},   // East
	{0, 1},   // South
	{1, 1},   // SouthEast
	{0, -1},  // North
	{1, -1},  // NorthEast
	{-1, -1}, // NorthWest
	{-1, 0},  // West
	{-1, 1},  // SouthWest
	{0, 0},   // Center
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
