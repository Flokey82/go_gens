package gamecs

import (
	"log"

	"github.com/Flokey82/aistate"
)

const StateTypeIdle aistate.StateType = 7

type StateIdle struct {
	ai *CompAi
}

func NewStateIdle(ai *CompAi) *StateIdle {
	return &StateIdle{ai: ai}
}

func (s *StateIdle) Type() aistate.StateType {
	return StateTypeIdle
}

func (s *StateIdle) Tick(delta uint64) {
}

func (s *StateIdle) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateIdle) OnExit() {
	log.Printf("leaving state %d", s.Type())
}
