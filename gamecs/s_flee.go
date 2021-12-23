package gamecs

import (
	"fmt"
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/vectors"
	"log"
)

const StateTypeFlee aistate.StateType = 1

// StateFlee
type StateFlee struct {
	ai *CAi
}

func NewStateFlee(ai *CAi) *StateFlee {
	return &StateFlee{ai: ai}
}

func (s *StateFlee) Type() aistate.StateType {
	return StateTypeFlee
}

func (s *StateFlee) Tick(delta uint64) {
	// Check if we are being chased!
	// Are we safe?

	// TODO: Return false if the state isn't complete yet.
	// Return true if we're done and safe.
}

func (s *StateFlee) OnEnter() {
	log.Printf("entering state %d", s.Type())
	s.ai.CAiPath.running = true // Run away!
	// Select a random point to run towards.
	// Ideally we'd choose target location that would lead us away from the threat.
	s.ai.SetTarget(vectors.RandomVec2(128.0))
	log.Println(fmt.Sprintf("fleeing to Target %.2f, %.2f", s.ai.CAiPath.Target.X, s.ai.CAiPath.Target.Y))
}

func (s *StateFlee) OnExit() {
	log.Printf("leaving state %d", s.Type())
	s.ai.CAiPath.running = false // We're safe, no need to run anymore.
}
