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

const StateTypeAttack aistate.StateType = 2

// StateAttack
type StateAttack struct {
	ai     *CAi
	target *Agent
}

func NewStateAttack(ai *CAi) *StateAttack {
	return &StateAttack{ai: ai}
}

func (s *StateAttack) Type() aistate.StateType {
	return StateTypeAttack
}

func (s *StateAttack) findTarget() {
	if s.target != nil && s.ai.CanSeeEntity(s.target) {
		s.ai.SetTarget(s.target.Pos)
		return
	}
	// Set our target we move to to the current position of the first entity that we have perceived.
	// (NOT the closest, the first)
	// Ideally we would choose our target based on distance, threat level, etc.
	s.target = s.ai.Entities[0]
	s.ai.SetTarget(s.target.Pos)
	s.ai.running = true // Run to intercept.
}

func (s *StateAttack) giveUpTarget() {
	s.target = nil
	// TODO: Unset aipath target.
	s.ai.running = false // No need to run anymore.
}

func (s *StateAttack) Tick(delta uint64) {
	if len(s.ai.CAiPerception.Entities) == 0 {
		return
	}

	// TODO: Add behavior tree or something...
	// - chase target
	// - if reached, attack target
	s.findTarget()

	log.Println(fmt.Sprintf("chasing Target %.2f, %.2f", s.ai.CAiPath.Target.X, s.ai.CAiPath.Target.Y))
}

func (s *StateAttack) OnEnter() {
	log.Printf("entering state %d", s.Type())
	s.findTarget()
}

func (s *StateAttack) OnExit() {
	log.Printf("leaving state %d", s.Type())
	s.giveUpTarget()
}
