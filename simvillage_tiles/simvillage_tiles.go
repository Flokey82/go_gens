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
	"github.com/hajimehoshi/ebiten/inpututil"

	_ "embed"
)

//go:embed tiles/dungeon_tiles.png
var dungeon_png []byte

//go:embed tiles/roguelikeitems.png
var items_png []byte

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
	tilesDefaultSet *TileSet
	tilesDungeonSet *TileSet
	tilesItemsSet   *TileSet
	runnerImage     *ebiten.Image
)

func init() {
	var err error

	// TODO: Move this to the world.go and the creatures.go files.
	tilesDefaultSet, err = NewTileSet(images.Tiles_png)
	if err != nil {
		log.Fatal(err)
	}

	tilesDungeonSet, err = NewTileSet(dungeon_png)
	if err != nil {
		log.Fatal(err)
	}

	tilesItemsSet, err = NewTileSet(items_png)
	if err != nil {
		log.Fatal(err)
	}

	// Decode an image from the image file's byte slice.
	// TODO: Move to TileSet.
	imgRunner, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(imgRunner)
}

// Game represents the main game struct.
// TODO: Move creatures to the individual world structs.
// ... or add a pointer to each creature that points to
// its world, so all creatures everywhere are still updated.
type Game struct {
	*MapCache
	player    *Creature   // player
	creatures []*Creature // NPCs (and player)
	indoors   bool        // for toggling between worlds (hacky)
	tileSet   *TileSet    // Tile set of the current world
	dWorld    World       // Default world
	iWorld    World       // Fake indoor world
}

func NewGame() *Game {
	g := &Game{
		tileSet: tilesDefaultSet,
		dWorld:  newDefaultWorld(),
		iWorld:  newFakeIndoorWorld(),
	}
	g.MapCache = newMapCache(g.dWorld)
	g.player = NewCreature(g, g.dWorld, [2]int{0, 0})
	g.refreshCache([2]int{0, 0})

	// Add the player to the creature index.
	g.addCreature(g.player)

	// Add some creatures.
	g.addCreature(NewCreature(g, g.dWorld, [2]int{2, 2}))
	g.addCreature(NewCreature(g, g.dWorld, [2]int{20, 2}))
	g.addCreature(NewCreature(g, g.iWorld, [2]int{2, 20}))
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

	// Transport player to another world.
	// TODO:
	// - Set player position in the new world
	//   and remember player position in the old world.
	// - Check if we are at a door.
	//   Each door would have a destination world and coordinate.
	//   The destination world would have a door that leads back
	//   to the entrance.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.changeWorld()
	}

	// If the currently cached center chunk does not match the
	// player chunk position, we need to refresh the render cache.
	if g.player.chunk != g.curChunkXY {
		g.refreshCache(g.player.chunk)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) && g.isPlayerOnTrigger() {
		trg := g.getPlayerTrigger()
		g.teleportPlayer(trg.Destination.X, trg.Destination.Y, trg.Destination.Chunk)
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

// changeWorld toggles between indoor and outdoor map.
// NOTE: This is just for testing purposes.
func (g *Game) changeWorld() {
	if g.indoors = !g.indoors; g.indoors {
		g.setWorld(g.iWorld)
	} else {
		g.setWorld(g.dWorld)
	}
}

// setWorld sets the currently rendered world.
func (g *Game) setWorld(w World) {
	g.MapCache.setNewWorld(w) // update the map cache
	g.tileSet = w.TileSet()   // update the render tile set
	g.player.w = w            // update the player's world
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	const xCount = screenWidth / tileSize

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
				screen.DrawImage(g.tileSet.TileImage(t), op)
			}
		}

		drawItems := func(items []*Item) {
			for _, it := range items {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(it.Position[0]*tileSize)+cOffsX, float64(it.Position[1]*tileSize)+cOffsY)

				// Get the right tile sprite.
				screen.DrawImage(tilesItemsSet.TileImage(it.Tile), op)
			}
		}

		// TODO: Once we have moved to a more generic layer naming system,
		// move to iterating through layers in the set order again.
		drawLayer(layers.Ground, false)
		drawLayer(layers.GroundOverlay, false)
		drawLayer(layers.Objects, false)
		drawItems(layers.Items)
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
	px, py := g.player.getTileXY()                 // Current player tile
	fx, fy := g.player.facingTile()                // The tile the player is facing at
	cx, cy := g.player.chunk[0], g.player.chunk[1] // Current chunk

	var triggerStr string
	// Get current chunk and compare trigger x, y to player x, y.
	// If it is a match, draw the trigger text.
	// TODO: Find a better way to determine if a trigger is set.
	if g.isPlayerOnTrigger() {
		// Draw the trigger text.
		triggerStr = "\nTRIGGER (door) found!"
	}

	// Draw ticks per second (TPS), current tile (T), viewport center tile (V), and current chunk (C).
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f (T %d, %d C %d, %d F %d, %d)%s", ebiten.ActualTPS(), px, py, cx, cy, fx, fy, triggerStr))
}

func (g *Game) isPlayerOnTrigger() bool {
	// Get current chunk and compare trigger x, y to player x, y.
	// If it is a match, draw the trigger text.
	px, py := g.player.getTileXY()                 // Current player tile
	cx, cy := g.player.chunk[0], g.player.chunk[1] // Current chunk

	for _, tp := range g.getChunk(cx, cy).Triggers {
		if tp.Position == [2]int{px, py} {
			return true
		}
	}
	return false
}

func (g *Game) getPlayerTrigger() Trigger {
	// Get current chunk and compare trigger x, y to player x, y.
	// If it is a match, draw the trigger text.
	px, py := g.player.getTileXY()                 // Current player tile
	cx, cy := g.player.chunk[0], g.player.chunk[1] // Current chunk

	for _, tp := range g.getChunk(cx, cy).Triggers {
		if tp.Position == [2]int{px, py} {
			return tp
		}
	}
	return Trigger{}
}

func (g *Game) teleportPlayer(x, y int, chunk [2]int) {
	g.player.pos = [2]int{x, y}
	g.player.chunk = chunk
	g.refreshCache(g.player.chunk)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// addCreature adds a creature to the game.
func (g *Game) addCreature(c *Creature) {
	g.creatures = append(g.creatures, c)
}

// getCreatures returns the creatures at the given tile in the given chunk in
// the current world (of the player).
func (g *Game) getCreatures(tileIdx int, chunk [2]int) []*Creature {
	var creatures []*Creature
	for _, c := range g.creatures {
		if c.w == c.g.player.w && c.tileIdx == tileIdx && c.chunk == chunk {
			creatures = append(creatures, c)
		}
	}
	return creatures
}
