package gamecs

import (
	"github.com/Flokey82/go_gens/aistate"
	"github.com/Flokey82/go_gens/aitree"
	"github.com/Flokey82/go_gens/vectors"
	"log"
)

const StateTypeFind aistate.StateType = 0

// StateFind is active when an agent is hungry and we don't have any food.
// In this state an agent will wander around until they find a food source.
// The agent will approach the food source and 'pick it up', which will reset
// their hunger and remove the food source from the world map.

type StateFindFood struct {
	ai  *CAi
	ait *aitree.Tree
	it  *Item
}

func NewStateFindFood(ai *CAi) *StateFindFood {
	s := &StateFindFood{
		ai:  ai,
		ait: aitree.New(),
	}

	// TODO: When inventory is full, head home!
	fci := aitree.NewSequence("find and pick up item")
	s.ait.Root = fci
	aw := newActionWander(ai, func() bool {
		if s.foundItem() {
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

	ac := newActionPickUpItem(ai, func() *Item {
		return s.it
	})
	fci.Append(ac)
	return s
}
func (s *StateFindFood) needItem() bool {
	return !s.foundItem()
}

func (s *StateFindFood) foundItem() bool {
	return s.it != nil && s.ai.CanSee(s.it)
}

func (s *StateFindFood) Type() aistate.StateType {
	return StateTypeFind
}

func (s *StateFindFood) Tick(delta uint64) {
	if s.ait.Tick() == aitree.StateFailure {
		log.Println("Find failed!!")
	}
}

func (s *StateFindFood) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateFindFood) OnExit() {
	log.Printf("leaving state %d", s.Type())
	// TODO: Reset tree.
}
