package simvillage_tiles

// chunkSourceFunc returns map chunks for our current world.
type chunkSourceFunc func(x, y int) *MapChunk

// MapCache provides a caching layer for serving MapChunks.
type MapCache struct {
	chunkCache [3][3]*MapChunk // chunkCache[x][y] is the chunk at x-1, y-1.
	curChunkXY [2]int          // curChunkXY[0] is the x position of the current chunk.
	World                      // fetchChunk is a function to fetch a chunk from the source.
}

func newMapCache() *MapCache {
	return &MapCache{
		World: newDefaultWorld(),
	}
}

// setNewWorld sets the given World as the source for chunks.
func (g *MapCache) setNewWorld(f World) {
	// Set the new chunk source.
	g.World = f

	// Reset the cache.
	g.refreshCache(g.curChunkXY)
}

// validCacheIdx returns true if the indices are within the bounds of x[0..2], y[0..2].
func validCacheIdx(x, y int) bool {
	return x >= 0 && x < 3 && y >= 0 && y < 3
}

func (g *MapCache) refreshCache(pChunk [2]int) {
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
				chunkCache[x][y] = g.FetchChunk(pChunk[0]+x-1, pChunk[1]+y-1)
			}
		}
	}
	g.chunkCache = chunkCache
	g.curChunkXY = pChunk
}

// getChunk returns the MapChunk at the given chunk position either from cache or freshly from
// the source (right now it is directly from the random number generator.
func (g *MapCache) getChunk(x, y int) *MapChunk {
	// TODO: Add (g *Game) isInCache(x, y) bool
	if cx, cy := g.curChunkXY[0]-x+1, g.curChunkXY[1]-y+1; validCacheIdx(cx, cy) && g.chunkCache[cx][cy] != nil {
		return g.chunkCache[cx][cy]
	}
	return g.FetchChunk(x, y)
}
