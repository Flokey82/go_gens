package main

import (
	"log"

	"github.com/Flokey82/go_gens/gamehex"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	ebiten.SetWindowTitle("Hexagonal (Ebitengine Demo)")
	ebiten.SetWindowSize(500, 600)
	ebiten.SetWindowResizable(true)

	g, err := gamehex.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
