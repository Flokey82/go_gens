package main

import (
	"github.com/Flokey82/go_gens/simvillage_simple"
)

func main() {
	v := simvillage_simple.New()
	for i := 0; i < 3000; i++ {
		v.Tick()
	}
}
