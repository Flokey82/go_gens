# Simple Game Loop in Golang
This is a very, very basic game loop... nothing fancy about it.

This implementation is a bare bones variant inspired by kutase's excellent [go-gameloop](https://github.com/kutase/go-gameloop).

## Example

```golang
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

```