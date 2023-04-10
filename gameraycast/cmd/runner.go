package main

import "github.com/Flokey82/go_gens/gameraycast"

func main() {
	g := gameraycast.NewGame(gameraycast.NewMap())
	g.Run()
}
