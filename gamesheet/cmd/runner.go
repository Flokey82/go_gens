package main

import (
	"fmt"
	"time"

	"github.com/Flokey82/gameloop"
	"github.com/Flokey82/go_gens/gamesheet"
)

func main() {
	// Create a new character.
	c := gamesheet.New(100, 100, 0, 180, 180, 180, 180)

	gl := gameloop.New(time.Millisecond, func(delta float64) {
		// Update the character.
		c.Tick(delta)

		fmt.Print("\033[H\033[2J")
		c.Log()
	})

	gl.Start()

	defer gl.Stop()

	for {
	}
}
