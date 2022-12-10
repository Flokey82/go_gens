package main

import (
	"log"
	"math/rand"

	"github.com/Flokey82/go_gens/simmemory"
)

func main() {
	// Create 20 dwarves.
	var dwarves []*Dwarf
	for i := 0; i < 20; i++ {
		dwarves = append(dwarves, newDwarf())
	}

	// Tick for 1000 days.
	for i := 0; i < 1000; i++ {
		for _, d := range dwarves {
			// Add a new random thought.
			d.Memory.AddThought(simmemory.Thought(rand.Intn(int(simmemory.ThoughtLast))))
			d.Tick()
		}
	}

	// Print the thoughts of the dwarves.
	for i, d := range dwarves {
		log.Printf("Dwarf %d:", i)
		d.Log()
	}
}

type Dwarf struct {
	*simmemory.Memory
}

func newDwarf() *Dwarf {
	return &Dwarf{
		Memory: simmemory.NewMemory(),
	}
}
