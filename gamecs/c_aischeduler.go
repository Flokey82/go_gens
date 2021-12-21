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
	as := ai.CAiState

	// Set up the two states we decide on if we are being threatened.
	sFlee := NewStateFlee(ai)
	sAttack := NewStateAttack(ai)

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
	sFind := NewStateFind(ai)
	c.AddAnyTransition(sFind, func() bool {
		// Check if there are predators around... if none are around
		// we can go and find a new random spot to move towards.
		return !as.states[sThreatened] && !as.states[sHungry]
	})

	// If we get hungry....
	sMunch := NewStateMunch(ai)
	c.AddAnyTransition(sMunch, func() bool {
		// Check if there are predators around and if we're hungry.
		return !as.states[sThreatened] && as.states[sHungry]
	})

	// Set our initial state.
	c.SetState(sFind)
}

func (c *CAiScheduler) Update(delta float64) {
	c.Tick(uint64(delta * 100))
}
