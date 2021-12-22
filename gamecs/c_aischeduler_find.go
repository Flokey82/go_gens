package gamecs

import (
	"fmt"
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/vectors"
	"log"
)

const StateTypeFind aistate.StateType = 0

// StateFind
type StateFind struct {
	ai *CAi
}

func NewStateFind(ai *CAi) *StateFind {
	return &StateFind{ai: ai}
}

func (s *StateFind) Type() aistate.StateType {
	return StateTypeFind
}

func (s *StateFind) Tick(delta uint64) {
	// Move towards resource, pick up item, etc ...
	if s.ai.CAiPath.active {
		return
	}
	// Select a random point within 128 meters.
	// It would make sense to do something more sensible... looking for resources
	// or whatever.
	s.ai.SetTarget(vectors.RandomVec2(128.0))
	log.Println(fmt.Sprintf("new Target %.2f, %.2f", s.ai.Target.X, s.ai.Target.Y))
}

func (s *StateFind) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateFind) OnExit() {
	log.Printf("leaving state %d", s.Type())
}
