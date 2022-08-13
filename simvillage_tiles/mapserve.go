package simvillage_tiles

// MapServe imitates a database serving map chunks.
type MapServe struct {
	layers *MapChunk // chunk at 0, 0 (this is statically generated)
}

func newMapServe() *MapServe {
	return &MapServe{
		layers: defaultChunk(),
	}
}

// fetchChunk returns the un-cached MapChunk from the generator.
func (g *MapServe) fetchChunk(x, y int) *MapChunk {
	if x != 0 || y != 0 {
		// Generate the chunk at the given position with the given dimensions (in number of tiles).
		return genChunk(x, y, screenWidth/tileSize, screenHeight/tileSize)
	}
	return g.layers // Position 0, 0 is special.
}
