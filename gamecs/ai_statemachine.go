package gamecs

import (
	"github.com/Flokey82/aistate"
)

type CAiScheduler struct {
	*aistate.StateMachine
}

func newCAiScheduler() *CAiScheduler {
	return &CAiScheduler{
		StateMachine: aistate.New(),
	}
}

// init initializes the state machine that controls agent behavior.
//
// NOTE: This is super-messy and would need to be re-written from
// scratch. There is no real rhyme or reason behind many of the
// transitions and there might be dead ends. I'll probably start
// over again in a new project at some point.
func (c *CAiScheduler) init(ai *CompAi) {
	// TODO: Add 'any' transitions in order of priority.
	// Set up the two states we decide on if we are being threatened.
	sFlee := NewStateFlee(ai)     // Flee from predator
	sAttack := NewStateAttack(ai) // Attack

	// Allow the transition to return one of multiple different transitions.
	c.AddAnySelector(func() aistate.State {
		// Ultimately we want to decide based on personality or our chances to win.
		// TODO: Check if we have enough action points to attack.
		if ai.Conflict() && !ai.CAiStatus.states[sInjured] {
			return sAttack
		}
		return sFlee
	}, ai.CAiStatus.IsFunc(sThreatened))

	// Empty our inventory if it is full.
	sStore := NewStateStoreFood(ai)
	c.AddAnyTransition(sStore, func() bool {
		return ai.w.mgr.GetEntityFromID(ai.id).CompInventory.IsFull()
	})

	// If we get hungry....
	sFind := NewStateFindFood(ai) // Find food
	sMunch := NewStateEatFood(ai) // Eat food
	c.AddAnySelector(func() aistate.State {
		if ai.CAiStatus.HasFood() {
			return sMunch // If we have food, we can go munch.
		}
		return sFind // ... otherwise we have to find food first.
	}, ai.CAiStatus.IsFunc(sHungry))

	// If we get sleepy, get some rest.
	sRest := NewStateRest(ai)
	c.AddAnyTransition(sRest, ai.CAiStatus.IsFunc(sExhausted))

	// Make sure we always have some food in our pocket.
	c.AddAnyTransition(sFind, ai.CAiStatus.NoFood)

	// Add selector to exit attack state.
	c.AddSelector(sAttack, func() aistate.State {
		if c.Previous != nil {
			return c.Previous
		}
		return sFind
	}, ai.CAiStatus.IsNotFunc(sThreatened))

	// Add selector to exit flee state.
	c.AddSelector(sFlee, func() aistate.State {
		if c.Previous != nil {
			return c.Previous
		}
		return sFind
	}, ai.CAiStatus.IsNotFunc(sThreatened))

	// Set our initial state.
	c.SetState(sFind)
}

func (c *CAiScheduler) Update(delta float64) {
	c.Tick(uint64(delta * 100))
}
