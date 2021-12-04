package gamecs

import (
	"fmt"
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/vectors"
	"log"
	"time"
)

type CAiScheduler struct {
	*aistate.StateMachine

	as  *CAiState
	ap  *CAiPath
	ape *CAiPerception
}

func newCAiScheduler() *CAiScheduler {
	return &CAiScheduler{
		StateMachine: aistate.New(),
	}
}

func (c *CAiScheduler) init(ap *CAiPath, ape *CAiPerception, as *CAiState) {
	// Set up the two states we decide on if we are being threatened.
	sFlee := NewStateFlee(ap, ape)
	sAttack := NewStateAttack(ap, ape)

	// Allow the transition to return one of multiple different transitions.
	c.AddAnySelector(func() aistate.State {
		// Randomly switch between attacking and fleeing.
		// Ultimately we want to decide based on personality or our chances to win.
		if time.Now().Unix()%2 != 0 {
			return sFlee
		}
		return sAttack
	}, func() bool {
		// Check if there are predators around.
		return as.states[sThreatened]
	})

	// This is the default state in which we determine a random point as target.
	sFind := NewStateFind(ap, ape)
	c.AddAnyTransition(sFind, func() bool {
		// Check if there are predators around... if none are around
		// we can go and find a new random spot to move towards.
		return !as.states[sThreatened]
	})

	// Set our initial state.
	c.SetState(sFind)

	c.as = as
	c.ap = ap
	c.ape = ape
}

func (c *CAiScheduler) Update(m *CMovable, delta float64) {
	c.Tick(uint64(delta * 100))
}

const (
	StateTypeFind   aistate.StateType = 0
	StateTypeFlee   aistate.StateType = 1
	StateTypeAttack aistate.StateType = 2
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
