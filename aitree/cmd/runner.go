package main

import (
	"github.com/Flokey82/go_gens/aitree"
	"log"
)

func main() {
	// Set up a new state machine.
	s := aitree.New()

	fdp := aitree.NewSequence("find and drink potion")
	s.Root = fdp

	// Construct a sub-tree to search for a potion.
	fp := aitree.NewSelector("find a potion")
	fdp.Append(fp)

	// No potion in inventory.
	hasPotion := newFakeLeaf("potion in inventory", aitree.StateFailure, 1)
	fp.Append(hasPotion)

	// Look around for one.
	searchForPotion := aitree.NewSelector("search for potion")
	fp.Append(searchForPotion)

	searchChest := newFakeLeaf("search in chest", aitree.StateFailure, 10)
	searchForPotion.Append(searchChest)

	searchCloset := newFakeLeaf("search in closet", aitree.StateFailure, 2)
	searchForPotion.Append(searchCloset)

	searchDesk := newFakeLeaf("search in desk", aitree.StateSuccess, 1)
	searchForPotion.Append(searchDesk)

	// Construct a sub-tree for drinking the potion.
	dp := aitree.NewSequence("drink potion")
	fdp.Append(dp)

	// Check expiry date on bottle.
	checkExpiry := newFakeLeaf("checking expiry date", aitree.StateSuccess, 1)
	dp.Append(checkExpiry)

	// Uncork bottle.
	uncorkBottle := newFakeLeaf("uncork bottle", aitree.StateSuccess, 1)
	dp.Append(uncorkBottle)

	// Tastes like raspberry.
	drinkPotion := newFakeLeaf("drink from bottle", aitree.StateSuccess, 1)
	dp.Append(drinkPotion)

	for s.Tick() == aitree.StateRunning {
		// Tick.
	}
}

type FakeLeaf struct {
	Name     string
	Outcome  aitree.State
	Current  int
	Attempts int
}

func newFakeLeaf(name string, outcome aitree.State, attempts int) *FakeLeaf {
	return &FakeLeaf{
		Name:     name,
		Outcome:  outcome,
		Attempts: attempts,
	}
}

func (l *FakeLeaf) Tick() aitree.State {
	log.Println(l.Name)
	for l.Current < l.Attempts {
		log.Println("run " + l.Name)
		l.Current++
		if l.Current < l.Attempts {
			return aitree.StateRunning
		}
		l.Current = 0
		log.Println("done with " + l.Name)
		return l.Outcome
	}
	return l.Outcome
}
