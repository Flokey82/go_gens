package gamecs

import (
	"fmt"
	"log"

	"github.com/Flokey82/aistate"
	"github.com/Flokey82/aitree"
	"github.com/Flokey82/go_gens/vectors"
)

const StateTypeFlee aistate.StateType = 1

// StateFlee
type StateFlee struct {
	ai *CompAi
}

func NewStateFlee(ai *CompAi) *StateFlee {
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
	ai     *CompAi
	ait    *aitree.Tree
	target *Agent
}

func NewStateAttack(ai *CompAi) *StateAttack {
	s := &StateAttack{
		ai:  ai,
		ait: aitree.New(),
	}

	// New action sequence.
	fci := aitree.NewSequence("chase and eliminate threat")
	s.ait.Root = fci

	// Move towards the target.
	am := newActionMoveTo(ai, s.needTarget, func() vectors.Vec2 {
		return s.target.Pos
	})
	fci.Append(am)

	// Attack the target.
	at := newActionAttack(ai, func() *Agent {
		return s.target
	})
	fci.Append(at)
	return s
}

func (s *StateAttack) Type() aistate.StateType {
	return StateTypeAttack
}
func (s *StateAttack) needTarget() bool {
	return !s.foundTarget()
}

func (s *StateAttack) foundTarget() bool {
	return s.target != nil && s.ai.CanSeeEntity(s.target)
}

// findTarget attempts to find a target to attack.
func (s *StateAttack) findTarget() {
	// Check if we have already found a target.
	if s.foundTarget() {
		// Check if the target is dead.
		// If so, we give up pursuing the target.
		if s.target.Dead() {
			s.giveUpTarget()
		}
		return
	}

	// Set our target to the current position of the first entity that we have perceived.
	// (the first is the closest).
	// Ideally we would choose our target based on distance, threat level, etc.
	if len(s.ai.Entities) > 0 {
		for _, e := range s.ai.Entities {
			if e.Dead() {
				continue
			}
			s.target = e
			s.ai.running = true // Run to intercept.
			break
		}
	} else {
		s.giveUpTarget() // No target in sight, give up.
	}
}

// giveUpTarget gives up on the current target.
func (s *StateAttack) giveUpTarget() {
	s.target = nil
	// TODO: Unset aipath target.
	s.ai.running = false // No need to run anymore.
}

func (s *StateAttack) Tick(delta uint64) {
	if s.ait.Tick() == aitree.StateFailure {
		log.Println(fmt.Sprintf("%d: StateAttack failed!!", s.ai.id))
	} else if s.target != nil {
		log.Println(fmt.Sprintf("chasing Target %.2f, %.2f", s.ai.CAiPath.Target.X, s.ai.CAiPath.Target.Y))
	}
}

func (s *StateAttack) OnEnter() {
	log.Printf("entering state %d", s.Type())
	s.findTarget()
}

func (s *StateAttack) OnExit() {
	log.Printf("leaving state %d", s.Type())
	s.giveUpTarget()
}
