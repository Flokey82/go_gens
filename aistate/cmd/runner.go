package main

import (
	"github.com/Flokey82/go_gens/aistate"
	"log"
)

func main() {
	// Set up a new state machine.
	s := aistate.New()

	// Set up new states.
	sFindRes := NewState(StateTypeFindResource)
	sMoveToRes := NewState(StateTypeMoveToResource)
	sCollectRes := NewState(StateTypeCollectResource)
	sFlee := NewState(StateTypeFlee)
	sFight := NewState(StateTypeFight)

	aggressive := true

	s.AddAnySelector(func() aistate.State {
		if aggressive {
			return sFight
		}
		return sFlee
	}, func() bool {
		// Check if there are predators around.
		return false
	})

	s.AddTransition(sFlee, sFindRes, func() bool {
		// Check if we are safe again.
		return true
	})

	s.AddTransition(sFindRes, sMoveToRes, func() bool {
		// Check if we have found a resource.
		return true
	})

	s.AddTransition(sMoveToRes, sCollectRes, func() bool {
		// Check if we have reached the resource.
		return true
	})

	s.AddTransition(sCollectRes, sFindRes, func() bool {
		// Check if we have successfully collected the resource.
		return true
	})

	// Set our initial state.
	s.SetState(sFindRes)

	// Run the simulation for a while.
	for i := 0; i < 100; i++ {
		s.Tick(10)
	}
}

const (
	StateTypeFindResource    aistate.StateType = 0
	StateTypeMoveToResource  aistate.StateType = 1
	StateTypeCollectResource aistate.StateType = 2
	StateTypeFlee            aistate.StateType = 3
	StateTypeFight           aistate.StateType = 4
)

// State is a fake implementation of a generic state.
type State struct {
	t aistate.StateType
}

func NewState(t aistate.StateType) *State {
	return &State{t: t}
}

func (s *State) Type() aistate.StateType {
	return s.t
}

func (s *State) Tick(delta uint64) {
	// Move towards resource, pick up item, etc ...
}

func (s *State) OnEnter() {
	log.Printf("entering state %d", s.t)
}

func (s *State) OnExit() {
	log.Printf("leaving state %d", s.t)
}
