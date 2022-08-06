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
	frameOX      = 0
	frameOY      = 32
	frameWidth   = 32
	frameHeight  = 32
	frameCount   = 8
)

const (
	tileSize = 16
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
	playerPos       [2]int
	playerLooksLeft bool
	count           int
}

func NewGame() *Game {
	g := &Game{
		layers: [][]int{
			{
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
	}
	return g
}

func (g *Game) Update() error {
	var keyPressed bool
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.canEnter(g.playerPos[0]-1, g.playerPos[1]) {
			g.playerPos[0] -= 1
			g.playerLooksLeft = true
			keyPressed = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.canEnter(g.playerPos[0]+1, g.playerPos[1]) {
			g.playerPos[0] += 1
			g.playerLooksLeft = false
			keyPressed = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.canEnter(g.playerPos[0], g.playerPos[1]-1) {
			g.playerPos[1] -= 1
			keyPressed = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.canEnter(g.playerPos[0], g.playerPos[1]+1) {
			g.playerPos[1] += 1
			keyPressed = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		px, py := g.getPlayerXY()
		g.layers[1][py*screenWidth/tileSize+px] = rand.Intn(255)
		keyPressed = true
	}
	if keyPressed {
		g.count++
	}
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

// getPlayerXY returns the x, y position of the player in the tile map.
func (g *Game) getPlayerXY() (int, int) {
	return getTileXYFromPlayerPos(g.playerPos[0], g.playerPos[1])
}

func (g *Game) Draw(screen *ebiten.Image) {
	w, _ := tilesImage.Size()
	tileXCount := w / tileSize

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	const xCount = screenWidth / tileSize

	// Draws the player on the screen.
	drawPlayer := func() {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		if g.playerLooksLeft {
			op.GeoM.Scale(-1, 1)
		}
		op.GeoM.Translate(screenWidth/2+float64(g.playerPos[0]), screenHeight/2+float64(g.playerPos[1]))
		i := (g.count / 5) % frameCount
		sx, sy := frameOX+i*frameWidth, frameOY
		screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	}

	// Draws the given layer on the screen.
	drawLayer := func(l []int, drawP bool) {
		pX, pY := g.getPlayerXY()
		for i, t := range l {
			tX, tY := i%xCount, i/xCount
			if pX == tX && pY == tY && drawP {
				drawPlayer()
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))

			sx := (t % tileXCount) * tileSize
			sy := (t / tileXCount) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

	// Iterate through the layers and draw them.
	for i, l := range g.layers {
		// TODO: Make sure we draw all tiles behind the player first.
		drawLayer(l, i == len(g.layers)-1)
	}

	// Draw some information on the screen.
	px, py := g.getPlayerXY()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (tp %d, %d)", ebiten.ActualTPS(), px, py))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
