package simvillage_tiles

import (
	"math/rand"
)

var (
	flowerTiles = []int{301, 302, 303, 304}                          // Tiles that are flowers
	grassTiles  = []int{243, 243, 243, 243, 243, 243, 243, 218, 244} // Tiles that are grass
)

func genChunk(x, y, width, height int) *MapChunk {
	// Create a new, empty chunk.
	chunk := newMapChunk(width, height)

	// Initialize the random number generator with a unique seed for the given coordinates.
	r := rand.New(rand.NewSource(cash(x, y)))
	// Fill the ground layer with grass tiles.
	// NOTE: We've repeated the 243 tile a few times in the grassTiles slice
	// to skew the probability in favor of the plain grass tile.
	chunk.Ground.fillRandom(grassTiles, r)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Randomly add some flowers on the ground overlay layer.
			if r.Intn(100) < 2 {
				chunk.GroundOverlay.setTile(x, y, flowerTiles[r.Intn(len(flowerTiles))])
			}
		}
	}

	// Randomly put some houses and other objects on the chunk.
	// NOTE: We currently have hardcoded the offset for the houses to (1, 1).
	switch r.Intn(10) {
	case 0:
		chunk.drawObject(house1, 1, 1)
	case 1:
		chunk.drawObject(house2, 1, 1)
	case 2:
		chunk.drawObject(house3, 1, 1)
	case 3:
		chunk.drawObject(hedge, 1, 1)
	}
	return chunk
}

// cash stands for chaos hash :D
// See: https://stackoverflow.com/a/37221804
func cash(x, y int) int64 {
	h := x*374761393 + y*668265263 //all constants are prime
	h = (h ^ (h >> 13)) * 1274126177
	return int64(h ^ (h >> 16))
}

// defaultChunk returns the default MapChunk.
func defaultChunk() *MapChunk {
	m := newMapChunk(screenWidth/tileSize, screenHeight/tileSize)
	m.Ground.Tiles = []int{
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
	}
	m.GroundOverlay.Tiles = []int{
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
		0, 196, 197, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 221, 222, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 286, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 221, 222, 0, 0, 0, 0, 0, 0,
	}
	m.Structures.Tiles = []int{
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
	}
	return m
}