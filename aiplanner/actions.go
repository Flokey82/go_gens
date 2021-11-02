package aiplanner

import (
	"errors"
)

// ErrActionInvalid indicates that the current action can not be executed with the given world state.
var ErrActionInvalid = errors.New("action invalid")

// Action is the interface all actions have to implement
type Action interface {
	String() string
	Cost() int
	CanRunIf(Planner, WorldState) bool
	Simulate(p Planner, worldState WorldFork) (WorldFork, error)
}

// ensure DefaultAction implements Action
var _ Action = (*DefaultAction)(nil)

// DefaultAction is a basic action implementation.
type DefaultAction struct {
	name       string
	cost       int
	conditions PlannerState
	effects    PlannerState
}

// NewAction returns a new default action.
func NewAction(name string, cost int, conditions, effects PlannerState) *DefaultAction {
	return &DefaultAction{
		name:       name,
		cost:       cost,
		conditions: conditions,
		effects:    effects,
	}
}

// String returns the name of the action.
func (a *DefaultAction) String() string {
	return a.name
}

// Cost returns the cost of the action.
func (a *DefaultAction) Cost() int {
	return a.cost
}

// CanRunIf returns true if 'agent' can perform action 'a' given the hypothetical world state 'world'.
func (a *DefaultAction) CanRunIf(agent Planner, world WorldState) bool {
	conditionsMet := world.Contains(a.conditions)
	effectsAchieved := world.Contains(a.effects)
	return conditionsMet && !effectsAchieved
}

// Simulate runs the action and returns a new hypothetical world state.
func (a *DefaultAction) Simulate(agent Planner, worldState WorldFork) (WorldFork, error) {
	if !a.CanRunIf(agent, worldState) {
		return nil, ErrActionInvalid
	}

	// Fork the world state and apply the effects of the current action.
	newState := worldState.Fork()
	newState.Update(a.effects)
	return newState, nil
}
