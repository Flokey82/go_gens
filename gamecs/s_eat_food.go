package gamecs

import (
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/aitree"
	"log"
)

const StateTypeMunch aistate.StateType = 3

// StateEatFood will be active if the agent is hungry and we have food.
type StateEatFood struct {
	ai  *CAi
	ait *aitree.Tree
}

func NewStateEatFood(ai *CAi) *StateEatFood {
	s := &StateEatFood{
		ai:  ai,
		ait: aitree.New(),
	}
	s.ait.Root = newActionConsumeItem(ai, func() *Item {
		a := ai.w.mgr.GetEntityFromID(ai.id)
		return a.CInventory.Find("food")
	})
	return s
}

func (s *StateEatFood) Type() aistate.StateType {
	return StateTypeMunch
}

func (s *StateEatFood) Tick(delta uint64) {
	if s.ait.Tick() == aitree.StateFailure {
		log.Println("Munch failed!!")
	}
}

func (s *StateEatFood) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateEatFood) OnExit() {
	// TODO: Reset tree.
	log.Printf("leaving state %d", s.Type())
}
