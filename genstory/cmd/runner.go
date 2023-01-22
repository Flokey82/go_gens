package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genstory"
)

func main() {
	log.Println(genstory.NewWorld(time.Now().UnixNano()))
}
