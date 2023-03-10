package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genreligion"
)

func main() {
	for i := 0; i < 10; i++ {
		if myth, err := genreligion.NewWorld(time.Now().UnixNano()); err != nil {
			log.Println(i, err)
		} else {
			log.Println(i, myth)
		}
	}
}
