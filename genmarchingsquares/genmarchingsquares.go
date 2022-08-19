// Package genmarchingsquares implements the marching squares algorithm.
// See: https://en.wikipedia.org/wiki/Marching_squares
package genmarchingsquares

// MarchSquares returns a grid of tiles encoded as 4 bit values that are generated from the given pixel grid.
func MarchSquares(pixels [][]bool, dimX, dimY int) [][]byte {
	squares := make([][]byte, dimX-1)
	for i := range squares {
		squares[i] = make([]byte, dimY-1)
	}

	// Walk the pixel grid and generate the tiles.
	//
	//  _____ pixel grid
	// | | | |
	// v v v v
	// +-+-+-+
	// | | | | <-
	// +-+-+-+   | tile grid
	// | | | | <-
	// +-+-+-+
	for x := 0; x < dimX-1; x++ {
		for y := 0; y < dimY-1; y++ {
			squares[x][y] = encodeTile(
				pixels[x][y],     // nw -> 8
				pixels[x+1][y],   // ne -> 4
				pixels[x+1][y+1], // se -> 2
				pixels[x][y+1],   // sw -> 1
			)
		}
	}
	return squares
}

// Encoded tile:
//
// 8-4    nw-ne
// | | <- |   |
// 1-2    sw-se
func encodeTile(nw, ne, se, sw bool) byte {
	var val byte
	if nw {
		val |= 1 << 3
	}
	if ne {
		val |= 1 << 2
	}
	if se {
		val |= 1 << 1
	}
	if sw {
		val |= 1 << 0
	}
	return val
}
