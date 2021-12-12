package gamecs

import (
	"fmt"
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/vectors"
	"log"
)

const (
	StateTypeFind   aistate.StateType = 0
	StateTypeFlee   aistate.StateType = 1
	StateTypeAttack aistate.StateType = 2
	StateTypeMunch  aistate.StateType = 3
)

// StateFlee
type StateFlee struct {
	ap  *CAiPath
	ape *CAiPerception
}

func NewStateFlee(ap *CAiPath, ape *CAiPerception) *StateFlee {
	return &StateFlee{ap: ap, ape: ape}
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
	s.ap.running = true // Run away!
	// Select a random point to run towards.
	// Ideally we'd choose target location that would lead us away from the threat.
	s.ap.SetTarget(vectors.RandomVec2(128.0))
	log.Println(fmt.Sprintf("fleeing to Target %.2f, %.2f", s.ap.Target.X, s.ap.Target.Y))
}

func (s *StateFlee) OnExit() {
	log.Printf("leaving state %d", s.Type())
	s.ap.running = false // We're safe, no need to run anymore.
}

// StateFind
type StateFind struct {
	ap  *CAiPath
	ape *CAiPerception
}

func NewStateFind(ap *CAiPath, ape *CAiPerception) *StateFind {
	return &StateFind{ap: ap, ape: ape}
}

func (s *StateFind) Type() aistate.StateType {
	return StateTypeFind
}

func (s *StateFind) Tick(delta uint64) {
	// Move towards resource, pick up item, etc ...
	if s.ap.active {
		return
	}
	// Select a random point within 128 meters.
	// It would make sense to do something more sensible... looking for resources
	// or whatever.
	s.ap.SetTarget(vectors.RandomVec2(128.0))
	log.Println(fmt.Sprintf("new Target %.2f, %.2f", s.ap.Target.X, s.ap.Target.Y))
}

func (s *StateFind) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateFind) OnExit() {
	log.Printf("leaving state %d", s.Type())
}

// StateAttack
type StateAttack struct {
	ap  *CAiPath
	ape *CAiPerception
}

func NewStateAttack(ap *CAiPath, ape *CAiPerception) *StateAttack {
	return &StateAttack{ap: ap, ape: ape}
}

func (s *StateAttack) Type() aistate.StateType {
	return StateTypeAttack
}

func (s *StateAttack) Tick(delta uint64) {
	if len(s.ape.Entities) == 0 || s.ap.active {
		return
	}

	// TODO: Add behavior tree or something...
	// - chase target
	// - if reached, attack target

	// For now we just set the target point we move to to the current
	// position of the first entity that we percieved (NOT the closest, the first).
	s.ap.SetTarget(s.ape.Entities[0].Pos)
	log.Println(fmt.Sprintf("chasing Target %.2f, %.2f", s.ap.Target.X, s.ap.Target.Y))
}

func (s *StateAttack) OnEnter() {
	log.Printf("entering state %d", s.Type())
	if len(s.ape.Entities) == 0 {
		return // There is nothing to attack?!
	}
	// Set our target we move to to the current position of the first entity that we have perceived.
	// (NOT the closest, the first)
	// Ideally we would choose our target based on distance, threat level, etc.
	s.ap.SetTarget(s.ape.Entities[0].Pos)
	s.ap.running = true // Run to intercept.
}

func (s *StateAttack) OnExit() {
	log.Printf("leaving state %d", s.Type())
	s.ap.running = false // No need to run anymore.
}

// StateMunch
type StateMunch struct {
	ap  *CAiPath
	ape *CAiPerception
	as  *CAiState
}

func NewStateMunch(ap *CAiPath, ape *CAiPerception, as *CAiState) *StateMunch {
	return &StateMunch{ap: ap, ape: ape, as: as}
}

func (s *StateMunch) Type() aistate.StateType {
	return StateTypeMunch
}

func (s *StateMunch) Tick(delta uint64) {
	if len(s.ape.Items) == 0 || s.ap.active {
		if !s.ap.active {
			s.ap.SetTarget(vectors.RandomVec2(128.0))
		}
		return
	}

	s.as.Eat()
	// TODO: Message that we're munching, so we'd need to reset hunger.
	log.Println(fmt.Sprintf("ate %.2f, %.2f", s.ap.Target.X, s.ap.Target.Y))
}

func (s *StateMunch) OnEnter() {
	log.Printf("entering state %d", s.Type())
	if len(s.ape.Items) == 0 {
		// There is nothing to eat.
		// Select a random point within 128 meters.
		s.ap.SetTarget(vectors.RandomVec2(128.0))
		return
	}
	// Move towards an item.
	s.ap.SetTarget(s.ape.Items[0].Pos)
}

func (s *StateMunch) OnExit() {
	log.Printf("leaving state %d", s.Type())
}
