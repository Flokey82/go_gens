package main

import (
	"frohrweck/simvillage"
)

func main() {
	game := simvillage.NewInstance()
	for {
		game.TickDay()
	}
}
