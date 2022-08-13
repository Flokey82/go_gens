package simvillage_tiles

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Creature represents a moving entity in the game.
type Creature struct {
	g         *Game
	pos       [2]int // current creature position relative to the chunk center
	tile      [2]int // current tile coordinates relative to the top left corner of the chunk
	chunk     [2]int // current global chunk coordinates
	tileIdx   int    // cached tile index
	looksLeft bool   // creature looks / moves to the left
	count     int    // for keeping track of the animation frames
}

// NewCreature returns a new creature with the given position.
func NewCreature(g *Game, pos [2]int) *Creature {
	c := &Creature{
		g:         g,
		pos:       pos,
		looksLeft: false,
	}
	c.updateTile() // update the current cached tile position
	return c
}

// updateTile calculates the current tile the creature is on and the
// index of that tile for quicker matching during rendering.
func (c *Creature) updateTile() {
	c.tile[0], c.tile[1] = getTileXYFromPos(c.pos[0], c.pos[1])
	const xCount = screenWidth / tileSize
	c.tileIdx = c.tile[1]*xCount + c.tile[0]
}

// move attempts to move the creature by the given xy delta.
func (c *Creature) move(delta [2]int) {
	if delta[0] == 0 && delta[1] == 0 {
		return
	}
	// Do we want to move left or right?
	c.looksLeft = delta[0] < 0

	newPosX := c.pos[0] + delta[0]
	newPosY := c.pos[1] + delta[1]

	// If we leave the current chunk, update where we are.
	newChunkDeltaX := newPosX / (screenWidth / 2)
	newChunkDeltaY := newPosY / (screenHeight / 2)

	// Can this be done better?
	newPosX -= (screenWidth - 1) * newChunkDeltaX
	newPosY -= (screenHeight - 1) * newChunkDeltaY
	newChunkX := c.chunk[0] + newChunkDeltaX
	newChunkY := c.chunk[1] + newChunkDeltaY

	// Now check if we can go where we'd end up.
	// TODO: Check if we can enter newPosX, newPosY in the new chunk.
	// NOTE: This doesn't check the chunk yet, but since all chunks are the same it sorta works.
	if c.g.canEnter(newChunkX, newChunkY, newPosX, newPosY) {
		// TODO: If only one of both directions would be an illegal move,
		// should we just move along one axis?
		c.pos[0] = newPosX
		c.pos[1] = newPosY
		c.chunk[0] = newChunkX
		c.chunk[1] = newChunkY
		c.updateTile()
	}
	c.count++
}

// Draw renders the creature on the screen with the given offset.
func (c *Creature) Draw(screen *ebiten.Image, cOffsX, cOffsY float64) {
	op := &ebiten.DrawImageOptions{}

	// Center the frame of the player.
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)

	// If the creature looks left, mirror horizontally.
	if c.looksLeft {
		op.GeoM.Scale(-1, 1)
	}

	// Center the player on the screen.
	op.GeoM.Translate(screenWidth/2+float64(c.pos[0])+cOffsX, screenHeight/2+float64(c.pos[1])+cOffsY)

	// Get the animation frame.
	i := (c.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

	// Wrap around the counter so we don't overflow.
	c.count %= frameCount * 5
}

func getTileXYFromPos(x, y int) (int, int) {
	return (x + screenWidth/2) / tileSize, (y + screenHeight/2) / tileSize
}

// getXY returns the x, y tile position of the creature in the tile map.
func (c *Creature) getXY() (int, int) {
	return getTileXYFromPos(c.pos[0], c.pos[1])
}
