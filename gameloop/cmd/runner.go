package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/gameloop"
)

func main() {
	gameloop.New(time.Second, func(delta float64) {
		log.Println("loop!")
	})
}
