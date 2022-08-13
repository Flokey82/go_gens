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
	*MapServe
	player     *Creature       // player
	creatures  []*Creature     // NPCs (and player)
	chunkCache [3][3]*MapChunk // cached chunks
	curChunkXY [2]int          // location of the cache (x,y of chunkCache[1][1])
}

func NewGame() *Game {
	g := &Game{
		MapServe: newMapServe(),
	}
	g.player = NewCreature(g, [2]int{0, 0})
	g.refreshCache()

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

	// If the currently cached center chunk does not match the
	// player chunk position, we need to refresh the cache.
	if g.player.chunk != g.curChunkXY {
		g.refreshCache()
	}

	// Handle "AI".
	for _, c := range g.creatures {
		// Skip the player, they move on their own.
		if c == g.player {
			continue
		}
		// NOTE: This just makes the creatures move randomly.
		c.move([2]int{rand.Intn(3) - 1, rand.Intn(3) - 1})
	}
	return nil
}

func (g *Game) refreshCache() {
	pChunk := g.player.chunk

	// Calculate chunk delta
	chunkDelta := [2]int{
		pChunk[0] - g.curChunkXY[0],
		pChunk[1] - g.curChunkXY[1],
	}

	// We move chunks in cache according to delta and copy the old chunks
	// to their new position in a new cache.
	//
	// (o): old player location chunk
	// (n): new player location chunk
	// [d]: chunks discarded
	// [f]: chunks fetched
	// (*): chunks copied / re-used
	//
	// The chunk delta in this case is -1, 1 as we are
	// moving to the left by one chunk and down by one.
	//
	// Note that 0, 0 is always considered the top-left
	// corner.
	//
	// Therefore:
	// - moving left is x-1, moving right is x+1
	// - moving up is y-1, moving down is y+1
	//
	//    |-------| old 3x3 cache
	//    [d][d][d]
	// [f]( )(o)[d]
	// [f](n)( )[d]
	// [f][f][f]
	// |-------| new 3x3 cache
	//
	// Example 1:
	//
	// The position x:0, y:2 in the new cache would be
	// x:-1, y:3 in the old cache, given the delta of -1, 1.
	//
	// ... Since this is out of bounds (x[0..2], y[0..2])
	// we will not find this position in the old cache and
	// will have to fetch the chunk freshly.
	//
	// Example 2:
	//
	// The position x:1, y:1 in the new cache would be
	// x:0, y:2 in the old cache, given the delta of -1, 1.
	//
	// ... Since this is within bounds (x[0..2], y[0..2])
	// we can simply copy the chunk from the old cache at
	// x:0, y:2 to the new position x:1, y:1 in the new cache.
	var chunkCache [3][3]*MapChunk

	// Iterate through the new cache slots.
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			// Calculate x,y chunk cache position in "old" cache.
			cdx := x + chunkDelta[0]
			cdy := y + chunkDelta[1]

			// If we are within the bounds of the old cache, we re-use the chunk we have in the
			// old cache, if one is present (e.g. not nil).
			//
			// If the tile in the old cache at cdx,cdy is nil, we have likely not initialized
			// the cache yet and need to fetch the chunk anyway.
			if validCacheIdx(cdx, cdy) && g.chunkCache[cdx][cdy] != nil {
				chunkCache[x][y] = g.chunkCache[cdx][cdy]
			} else {
				// If we are out of bounds of the old cache, we fetch the chunk.
				chunkCache[x][y] = g.fetchChunk(pChunk[0]+x-1, pChunk[1]+y-1)
			}
		}
	}
	g.chunkCache = chunkCache
	g.curChunkXY = g.player.chunk
}

// validCacheIdx returns true if the indices are within the bounds of x[0..2], y[0..2].
func validCacheIdx(x, y int) bool {
	return x >= 0 && x < 3 && y >= 0 && y < 3
}

// canEnter returns whether the player can enter the tile at (x, y) in the chunk (cX, cY).
// TODO: Improve collision detection.
func (g *Game) canEnter(cX, cY, newX, newY int) bool {
	x, y := getTileXYFromPos(newX, newY)
	if x < 0 || x >= screenWidth/tileSize || y < 0 || y >= screenHeight/tileSize {
		return false
	}
	// TODO: Allow multiple layers to be checked for collision.
	layers := g.getChunk(cX, cY)
	return layers.Structures.getTile(x, y) == 0
}

// getViewportXY returns the x, y tile position of the top left corner of the viewport in the tile map.
func (g *Game) getViewportXY() (int, int) {
	return g.player.pos[0] / tileSize, g.player.pos[1] / tileSize
}

// getChunk returns the MapChunk at the given chunk position either from cache or freshly from
// the source (right now it is directly from the random number generator.
func (g *Game) getChunk(x, y int) *MapChunk {
	// TODO: Add (g *Game) isInCache(x, y) bool
	if cx, cy := g.curChunkXY[0]-x+1, g.curChunkXY[1]-y+1; validCacheIdx(cx, cy) && g.chunkCache[cx][cy] != nil {
		return g.chunkCache[cx][cy]
	}
	return g.fetchChunk(x, y)
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
		layers := g.chunkCache[x+1][y+1] //g.getChunk(x+vpChunk[0], y+vpChunk[1])
		// Calculate the offset of the current chunk relative
		// to the player position / the center of the screen.
		// TODO: Also tidy up the chunk offset transformation.
		cOffsX := float64(x*screenWidth) - float64(vpPos[0])
		cOffsY := float64(y*screenHeight) - float64(vpPos[1])
		cxy[0] = x + vpChunk[0]
		cxy[1] = y + vpChunk[1]

		drawLayer := func(l *Layer, drawP bool) {
			// TODO: Draw only the visible tiles and fetch them according to
			// the player's position and the viewport.
			for i, t := range l.Tiles {
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

		// TODO: Once we have moved to a more generic layer naming system,
		// move to iterating through layers in the set order again.
		drawLayer(layers.Ground, false)
		drawLayer(layers.GroundOverlay, false)
		drawLayer(layers.Objects, false)
		drawLayer(layers.Structures, true)
		drawLayer(layers.Roof, false)
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

	// Draw some debug information on the screen.
	g.drawDebugInfo(screen)
}

// drawDebugInfo prints debug information on the screen.
func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	px, py := g.player.getXY()                     // Current player tile
	vx, vy := g.getViewportXY()                    // Current viewport tile
	cx, cy := g.player.chunk[0], g.player.chunk[1] // Current chunk

	// Draw ticks per second (TPS), current tile (T), viewport center tile (V), and current chunk (C).
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (T %d, %d V %d, %d C %d, %d)", ebiten.ActualTPS(), px, py, vx, vy, cx, cy))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
