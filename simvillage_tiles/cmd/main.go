package main

import (
	"log"

	"github.com/Flokey82/go_gens/simvillage_tiles"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

func main() {
	g := simvillage_tiles.NewGame()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Tiles (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
