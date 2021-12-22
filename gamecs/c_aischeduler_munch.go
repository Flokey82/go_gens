package gamecs

import (
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/aitree"
	"github.com/Flokey82/go_gens/vectors"
	"log"
)

const StateTypeMunch aistate.StateType = 3

// StateMunch
type StateMunch struct {
	ai  *CAi
	ait *aitree.Tree
	it  *Item
}

func NewStateMunch(ai *CAi) *StateMunch {
	s := &StateMunch{
		ai:  ai,
		ait: aitree.New(),
	}

	fci := aitree.NewSequence("find and consume item")
	s.ait.Root = fci
	aw := newActionWander(ai, func() bool {
		if s.hasItem() {
			return true
		}
		if len(s.ai.CAiPerception.Items) > 0 {
			s.it = s.ai.CAiPerception.Items[0]
			return true
		}
		return false
	})
	fci.Append(aw)

	am := newActionMoveTo(ai, s.needItem, func() vectors.Vec2 {
		return s.it.Pos
	})
	fci.Append(am)

	ac := newActionConsumeItem(ai, s.needItem, func() *Item {
		return s.it
	})
	fci.Append(ac)
	return s
}

func (s *StateMunch) needItem() bool {
	return !s.hasItem()
}

func (s *StateMunch) hasItem() bool {
	return s.it != nil && s.ai.CanSee(s.it)
}

func (s *StateMunch) Type() aistate.StateType {
	return StateTypeMunch
}

func (s *StateMunch) Tick(delta uint64) {
	if s.ait.Tick() == aitree.StateFailure {
		log.Println("Munch failed!!")
	}
}

func (s *StateMunch) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateMunch) OnExit() {
	log.Printf("leaving state %d", s.Type())
}
