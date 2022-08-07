package simvillage_tiles

import (
	"math/rand"
)

const (
	TileGrass = 243
)

func genLayers(width, height int) [][]int {
	layers := make([][]int, 3)
	// Generate 3 layers: Ground, decoration, and solid objects.
	for i := 0; i < len(layers); i++ {
		layers[i] = make([]int, width*height)
	}

	fillLayer(layers[0], TileGrass)

	// Fill the ground layer.
	return layers
}

func initLayers(layers [][]int) {
	for i := 0; i < len(layers); i++ {
		if i == 0 {
			fillLayer(layers[0], TileGrass)
		} else {
			fillLayer(layers[i], 0)
		}
	}
}

func fillLayer(layer []int, v int) {
	for i := 0; i < len(layer); i++ {
		layer[i] = v
	}
}

func defaultLayers() [][]int {
	return [][]int{
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
			0, 196, 197, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 221, 222, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 286, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 221, 222, 0, 0, 0, 0, 0, 0,
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
	}
}

var (
	flowerTiles = []int{301, 302, 303, 304}                          // Tiles that are flowers
	grassTiles  = []int{243, 243, 243, 243, 243, 243, 243, 218, 244} // Tiles that are grass
)

func genChunk(x, y, width, height int) [][]int {
	// Create a new, empty chunk.
	chunk := make([][]int, 3)
	for i := range chunk {
		chunk[i] = make([]int, width*height)
	}

	// Initialize the random number generator with a unique seed for the given coordinates.
	r := rand.New(rand.NewSource(cash(x, y)))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Fill the ground layer with grass tiles.
			// NOTE: We've repeated the 243 tile a few times in the grassTiles slice
			// to skew the probability in favor of the plain grass tile.
			chunk[0][x+y*width] = grassTiles[r.Intn(len(grassTiles))]

			// Randomly add some flowers on the ground object layer.
			if r.Intn(100) < 2 {
				chunk[1][x+y*width] = flowerTiles[r.Intn(len(flowerTiles))]
			}
		}
	}

	// drawObject draws the given object at (dx, dy).
	drawObject := func(h drawable, dx, dy int) {
		for x := 0; x < h.width; x++ {
			for y := 0; y < h.height; y++ {
				for lIdx, layer := range h.layers {
					if layer[x+y*h.width] != 0 {
						chunk[lIdx][(x+dx)+(y+dy)*width] = layer[x+y*h.width]
					}
				}
			}
		}
	}

	// Randomly put some houses and other objects on the chunk.
	// NOTE: We currently have hardcoded the offset for the houses to (1, 1).
	switch r.Intn(10) {
	case 0:
		drawObject(house1, 1, 1)
	case 1:
		drawObject(house2, 1, 1)
	case 2:
		drawObject(hedge, 1, 1)
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

// drawable is a collection of tiles that can be drawn on a chunk.
// TODO: Use named layers instead of a fixed number of expected layers.
type drawable struct {
	height int     // Height of the object.
	width  int     // Width of the object.
	layers [][]int // Tiles on each layer
}

var house1 = drawable{
	height: 4,
	width:  4,
	layers: [][]int{
		{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		},
		{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		},
		{
			58, 59, 60, 61,
			83, 84, 85, 86,
			108, 109, 110, 111,
			133, 134, 135, 136,
		},
	},
}
var house2 = drawable{
	height: 4,
	width:  6,
	layers: [][]int{
		{
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
		},
		{
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
		},
		{
			63, 64, 65, 66, 67, 68,
			88, 89, 90, 91, 92, 93,
			113, 114, 115, 116, 117, 118,
			138, 139, 140, 141, 142, 143,
		},
	},
}
var hedge = drawable{
	height: 5,
	width:  9,
	layers: [][]int{
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 202, 203, 204, 205, 206, 207, 208, 0,
			0, 227, 228, 229, 230, 231, 232, 233, 0,
			0, 0, 0, 0, 255, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 280, 0, 0, 0, 0,
		},
		{
			176, 177, 178, 179, 180, 181, 182, 183, 184,
			201, 0, 0, 0, 0, 0, 0, 0, 209,
			226, 0, 0, 0, 0, 0, 0, 0, 234,
			251, 252, 253, 254, 0, 256, 257, 258, 259,
			276, 277, 278, 279, 0, 281, 282, 283, 284,
		},
	},
}
