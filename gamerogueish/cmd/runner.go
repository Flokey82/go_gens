package main

import (
	"github.com/Flokey82/go_gens/gamerogueish"
)

func main() {
	g, err := gamerogueish.NewGame(gamerogueish.GenWorldSimpleDungeon, 100, 100, -1)
	if err != nil {
		panic(err)
	}
	g.Start()
}
