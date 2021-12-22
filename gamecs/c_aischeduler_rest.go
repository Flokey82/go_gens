package gamecs

import (
	"github.com/Flokey82/go_gens/aistate"
	"log"
)

const StateTypeRest aistate.StateType = 4

// StateRest will be active once an agent is exhausted.
// In this state, the agent will attempt to return home to rest.
type StateRest struct {
	ai *CAi
}

func NewStateRest(ai *CAi) *StateRest {
	return &StateRest{
		ai: ai,
	}
}

func (s *StateRest) Type() aistate.StateType {
	return StateTypeRest
}

func (s *StateRest) Tick(delta uint64) {
	if s.ai.CAiPath.active {
		return
	}
	log.Println("arrived home!")
	s.ai.Sleep()
}

func (s *StateRest) OnEnter() {
	log.Printf("entering state %d", s.Type())
	// TODO: Recall home location.
	// - Set navigation target
	// - On arrival: rest, reset exhaustion
	loc := s.ai.GetLocation("home")
	s.ai.SetTarget(loc)
}

func (s *StateRest) OnExit() {
	log.Printf("leaving state %d", s.Type())
}
