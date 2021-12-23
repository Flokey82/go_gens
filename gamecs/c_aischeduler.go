package gamecs

import (
	"github.com/Flokey82/go_gens/aistate"
	"time"
)

type CAiScheduler struct {
	*aistate.StateMachine
}

func newCAiScheduler() *CAiScheduler {
	return &CAiScheduler{
		StateMachine: aistate.New(),
	}
}

func (c *CAiScheduler) init(ai *CAi) {
	// Set up the two states we decide on if we are being threatened.
	sFlee := NewStateFlee(ai)     // Flee from predator
	sAttack := NewStateAttack(ai) // Attack

	// Allow the transition to return one of multiple different transitions.
	c.AddAnySelector(func() aistate.State {
		// Randomly switch between attacking and fleeing.
		// Ultimately we want to decide based on personality or our chances to win.
		if time.Now().Unix()%2 != 0 {
			return sFlee
		}
		return sAttack
	}, func() bool {
		return ai.CAiStatus.states[sThreatened]
	})

	// If we get hungry....
	sFind := NewStateFindFood(ai) // Find food
	sMunch := NewStateEatFood(ai) // Eat food
	c.AddAnySelector(func() aistate.State {
		if ai.CAiStatus.HasFood() {
			return sMunch // If we have food, we can go munch.
		}
		return sFind // ... otherwise we have to find food first.
	}, func() bool {
		return ai.CAiStatus.states[sHungry]
	})

	// If we get sleepy....
	sRest := NewStateRest(ai)
	c.AddAnyTransition(sRest, func() bool {
		return ai.CAiStatus.states[sExhausted]
	})

	c.AddAnyTransition(sFind, func() bool {
		// Check if there are predators around... if none are around
		// we can go and find a new random spot to move towards.
		return !ai.CAiStatus.states[sThreatened] && !ai.CAiStatus.states[sHungry]
	})

	// Set our initial state.
	c.SetState(sFind)
}

func (c *CAiScheduler) Update(delta float64) {
	c.Tick(uint64(delta * 100))
}
