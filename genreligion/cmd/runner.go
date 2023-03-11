package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genreligion"
)

func main() {
	for i := 0; i < 10; i++ {
		seed := time.Now().UnixNano()
		gen := genreligion.NewGenerator(seed, nil)
		god, err := gen.GetDeity()
		if err != nil {
			log.Println(i, err)
		} else if myth, err := gen.NewCreation(god); err != nil {
			log.Println(i, err)
		} else {
			log.Println(i, myth)
		}
	}
}
