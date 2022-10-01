package main

import (
	"github.com/Flokey82/go_gens/gamerogueish"
)

func main() {
	g, err := gamerogueish.NewGame()
	if err != nil {
		panic(err)
	}
	g.Start()
}
