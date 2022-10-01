package main

import (
	"github.com/Flokey82/go_gens/gamerogueish"
)

func main() {
	g, err := gamerogueish.NewGame(gamerogueish.GenSimpleDungeon)
	if err != nil {
		panic(err)
	}
	g.Start()
}
