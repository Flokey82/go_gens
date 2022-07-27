package gamecs

import (
	"fmt"
	"log"

	"github.com/Flokey82/aistate"
	"github.com/Flokey82/aitree"
	"github.com/Flokey82/go_gens/vectors"
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
		s.it = nil
		log.Println(fmt.Sprintf("%d: Find failed!!", s.ai.id))
	}
}

func (s *StateFindFood) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateFindFood) OnExit() {
	log.Printf("leaving state %d", s.Type())
	// TODO: Reset tree.
	s.it = nil
}

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
		log.Println(fmt.Sprintf("%d: Munch failed!!", s.ai.id))
	}
}

func (s *StateEatFood) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateEatFood) OnExit() {
	// TODO: Reset tree.
	log.Printf("leaving state %d", s.Type())
}

const StateTypeStore aistate.StateType = 5

type StateStoreFood struct {
	ai  *CAi
	ait *aitree.Tree
}

func NewStateStoreFood(ai *CAi) *StateStoreFood {
	s := &StateStoreFood{
		ai:  ai,
		ait: aitree.New(),
	}
	ghst := aitree.NewSequence("go home and store stuff")
	s.ait.Root = ghst

	ghst.Append(newActionMoveTo(ai, func() bool {
		return false
	}, func() vectors.Vec2 {
		return s.ai.GetPosition("home")
	}))
	ghst.Append(newActionTransferItems(ai, func() *CInventory {
		return s.ai.GetLocation("home").CInventory
	}))
	return s
}

func (s *StateStoreFood) Type() aistate.StateType {
	return StateTypeStore
}

func (s *StateStoreFood) Tick(delta uint64) {
	if s.ait.Tick() == aitree.StateFailure {
		log.Println(fmt.Sprintf("%d: Store failed!!", s.ai.id))
	}
}

func (s *StateStoreFood) OnEnter() {
	log.Printf("entering state %d", s.Type())
}

func (s *StateStoreFood) OnExit() {
	// TODO: Reset tree.
	log.Printf("leaving state %d", s.Type())
}
