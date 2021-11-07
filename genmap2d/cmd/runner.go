package main

import (
	"github.com/Flokey82/go_gens/genmap2d"
)

func main() {
	v := genmap2d.New(128, 128, 1234)
	for i := 0; i < 5; i++ {
		v.PlaceVillage()
	}
	v.ExportPng("rgb.png")
}
