package simvillage_tiles

// World is an interface providing chunks at the given coordinates.
// TODO: Add a function that returns the active creatures in this world.
type World interface {
	FetchChunk(x, y int) *MapChunk
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
