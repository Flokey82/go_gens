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
			var val byte
			if pixels[x][y] {
				val |= 1 << 0
			}
			if pixels[x+1][y] {
				val |= 1 << 1
			}
			if pixels[x+1][y+1] {
				val |= 1 << 2
			}
			if pixels[x][y+1] {
				val |= 1 << 3
			}
			squares[x][y] = val
		}
	}
	return squares
}
