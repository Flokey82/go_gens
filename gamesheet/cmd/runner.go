package main

import (
	"fmt"
	"math/rand"
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
		switch rand.Intn(1000) {
		case 0:
			c.TakeDamage(rand.Intn(10))
		case 1:
			c.TakeAction(rand.Intn(10))
		case 2:
			c.AddExperience(uint16(rand.Intn(int(c.NextLevelXP()))))
		case 3:
			c.Heal(rand.Intn(10))
		case 4:
			//c.RestoreAP(rand.Intn(10))
		case 5:
			c.SetState(gamesheet.StateAwake)
		case 6:
			c.SetState(gamesheet.StateAsleep)
		}
	})

	gl.Start()

	defer gl.Stop()

	for {
	}
}
