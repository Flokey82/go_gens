package simvillage_tiles

// World is an interface providing chunks at the given coordinates.
// TODO: Add a function that returns the active creatures in this world.
// TODO: Add a function that retrieves a specific tile.
// This would allow us to get information on doors to other maps etc.
type World interface {
	FetchChunk(x, y int) *MapChunk // returns a specific chunk
	TileSet() *TileSet             // returns the tileset for this world
}

// DefaultWorld is the default chunk source.
type DefaultWorld struct {
	layers *MapChunk // chunk at 0, 0 (this is statically generated)
}

// newDefaultWorld returns the default procedural chunk generator.
func newDefaultWorld() *DefaultWorld {
	return &DefaultWorld{
		layers: defaultChunk(),
	}
}

// FetchChunk is a function to fetch a chunk from the world.
func (w *DefaultWorld) FetchChunk(x, y int) *MapChunk {
	if x != 0 || y != 0 {
		// Generate the chunk at the given position with the given dimensions (in number of tiles).
		return genChunk(x, y, screenWidth/tileSize, screenHeight/tileSize)
	}
	return w.layers // Position 0, 0 is special.
}

// TileSet returns the tileset used for this world.
func (w *DefaultWorld) TileSet() *TileSet {
	return tilesDefaultSet
}

// FakeIndoorWorld is a fake indoor map.
type FakeIndoorWorld struct {
}

func (w *FakeIndoorWorld) FetchChunk(x, y int) *MapChunk {
	// Initialize the random number generator with a unique seed for the given coordinates.
	r := newRandForChunk(x, y)

	width := screenWidth / tileSize
	height := screenHeight / tileSize

	// Create a new, empty chunk.
	chunk := newMapChunk(width, height)
	chunk.Ground.fillRandom(dungeonFloorTiles, r)
	return chunk
}

// TileSet returns the tileset used for this world.
func (w *FakeIndoorWorld) TileSet() *TileSet {
	return tilesDungeonSet
}

var dungeonFloorTiles = []int{48, 49, 50, 71, 72, 73, 94, 95, 96} // Tiles that are stone floor (in the dungeon)

// canEnter returns whether the player can enter the tile at (x, y) in the chunk (cX, cY).
// TODO: Improve collision detection.
func canEnter(w World, cX, cY, newX, newY int) bool {
	x, y := getTileXYFromPos(newX, newY)
	if x < 0 || x >= screenWidth/tileSize || y < 0 || y >= screenHeight/tileSize {
		return false
	}
	// TODO: Allow multiple layers to be checked for collision.
	layers := w.FetchChunk(cX, cY)
	return layers.Structures.getTile(x, y) == 0
}
