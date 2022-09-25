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
	ai  *CompAi
	ait *aitree.Tree
	it  *Item
}

func NewStateFindFood(ai *CompAi) *StateFindFood {
	s := &StateFindFood{
		ai:  ai,
		ait: aitree.New(),
	}

	// New action sequence.
	fci := aitree.NewSequence("find and pick up item")
	s.ait.Root = fci

	// Wander around until we find an item.
	aw := newActionWander(ai, func() bool {
		if s.foundItem() {
			return true // We have already found an item, so stop wandering.
		}
		// We can see an item, so stop wandering and set the item
		// as our target.
		if len(s.ai.CAiPerception.Items) > 0 {
			s.it = s.ai.CAiPerception.Items[0]
			return true
		}

		// We can't see an item, so keep wandering.
		return false
	})
	fci.Append(aw)

	// Move to the item.
	am := newActionMoveTo(ai, s.needItem, func() vectors.Vec2 {
		return s.it.Pos
	})
	fci.Append(am)

	// Pick up the item.
	ac := newActionPickUpItem(ai, func() *Item {
		return s.it
	})
	fci.Append(ac)
	return s
}

// needItem returns true if we haven't found an item yet.
func (s *StateFindFood) needItem() bool {
	return !s.foundItem()
}

// foundItem returns true if we have an item set
// and it is still in the world (i.e. we can see it).
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
	ai  *CompAi
	ait *aitree.Tree
}

func NewStateEatFood(ai *CompAi) *StateEatFood {
	s := &StateEatFood{
		ai:  ai,
		ait: aitree.New(),
	}

	// Consume food.
	s.ait.Root = newActionConsumeItem(ai, func() *Item {
		// Get our agent from the entity manager.
		a := ai.w.mgr.GetEntityFromID(ai.id)

		// Try to find an item tagged as food.
		return a.CompInventory.Find("food")
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
	ai  *CompAi
	ait *aitree.Tree
}

func NewStateStoreFood(ai *CompAi) *StateStoreFood {
	s := &StateStoreFood{
		ai:  ai,
		ait: aitree.New(),
	}

	// New action sequence.
	ghst := aitree.NewSequence("go home and store stuff")
	s.ait.Root = ghst

	// Go home.
	ghst.Append(newActionMoveTo(ai, func() bool {
		return false
	}, func() vectors.Vec2 {
		return s.ai.GetPosition("home")
	}))

	// Store stuff.
	ghst.Append(newActionTransferItems(ai, func() *CompInventory {
		return s.ai.GetLocation("home").CompInventory
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
