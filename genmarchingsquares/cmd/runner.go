package main

import (
	"log"

	"github.com/Flokey82/go_gens/genmarchingsquares"
)

func main() {
	const (
		dimX = 4
		dimY = 4
	)
	pixels := make([][]bool, dimX)
	pixels[0] = []bool{false, true, true, false}
	pixels[1] = []bool{true, false, false, true}
	pixels[2] = []bool{true, false, false, true}
	pixels[3] = []bool{true, true, true, false}
	squares := genmarchingsquares.MarchSquares(pixels, dimX, dimY)
	if err := genmarchingsquares.ExportToPNG(squares, dimX-1, dimY-1, 128, "squares.png"); err != nil {
		log.Println(err)
	}
}
