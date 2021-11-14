package main

import (
	"github.com/Flokey82/go_gens/gameloop"
	"log"
	"time"
)

func main() {
	gameloop.New(time.Second, func(delta float64) {
		log.Println("loop!")
	})
}
