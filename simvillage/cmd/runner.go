package main

import (
	"github.com/Flokey82/go_gens/simvillage"
)

func main() {
	game := simvillage.NewInstance()
	for {
		game.TickDay()
	}
}
