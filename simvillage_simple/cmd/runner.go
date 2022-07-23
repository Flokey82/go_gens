package main

import (
	"github.com/Flokey82/go_gens/simvillage_simple"
)

func main() {
	v := simvillage_simple.New()
	v.AddRandomPerson()
	v.AddRandomPerson()
	v.AddRandomPerson()
	v.AddRandomPerson()
	v.AddRandomPerson()
	for i := 0; i < 60000; i++ {
		v.Tick()
		if len(v.People) == 0 {
			break
		}
	}
}
