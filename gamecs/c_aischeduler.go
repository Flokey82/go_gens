package gamecs

import (
	"fmt"
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/vectors"
	"log"
)

type CAiScheduler struct {
	*aistate.StateMachine
}

func newCAiScheduler() *CAiScheduler {
	return &CAiScheduler{
		StateMachine: aistate.New(),
	}
}

func (c *CAiScheduler) init(ap *CAiPath, ape *CAiPerception) {
	sFlee := NewStateFlee(ap, ape)
	sFind := NewStateFind(ap, ape)

	c.AddAnyTransition(sFlee, func() bool {
		// Check if there are predators around.
		return len(ape.Entities) > 0
	})

	c.AddAnyTransition(sFind, func() bool {
		// Check if there are predators around.
		return len(ape.Entities) == 0
	})

	// Set our initial state.
	c.SetState(sFind)
}

func (c *CAiScheduler) Update(m *CMovable, ap *CAiPath, ape *CAiPerception, delta float64) {
	c.Tick(uint64(delta * 100))
}

const (
	StateTypeFind aistate.StateType = 0
	StateTypeFlee aistate.StateType = 1
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
	// Move towards resource, pick up item, etc ...
}

func (s *StateFlee) OnEnter() {
	log.Printf("entering state %d", s.Type())
	s.ap.running = true
	s.ap.SetTarget(vectors.RandomVec2(64.0))
	log.Println(fmt.Sprintf("fleeing to Target %.2f, %.2f", s.ap.Target.X, s.ap.Target.Y))
}

func (s *StateFlee) OnExit() {
	log.Printf("leaving state %d", s.Type())
	s.ap.running = false
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
	s.ap.SetTarget(vectors.RandomVec2(64.0)) // Random point within 18 meters.
	log.Println(fmt.Sprintf("new Target %.2f, %.2f", s.ap.Target.X, s.ap.Target.Y))
}

func (s *StateFind) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateFind) OnExit() {
	log.Printf("leaving state %d", s.Type())
}
