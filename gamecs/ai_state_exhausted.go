package gamecs

import (
	"fmt"
	"log"

	"github.com/Flokey82/aistate"
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
	log.Println(fmt.Sprintf("%d: arrived home!!", s.ai.id))
	s.ai.Sleep()
}

func (s *StateRest) OnEnter() {
	log.Printf("entering state %d", s.Type())
	// TODO: Recall home location.
	// - Set navigation target
	// - On arrival: rest, reset exhaustion
	s.ai.SetTarget(s.ai.GetPosition("home"))
}

func (s *StateRest) OnExit() {
	log.Printf("leaving state %d", s.Type())
}
