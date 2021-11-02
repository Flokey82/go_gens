// Package aiplanner provides a very simple GOAP implementation.
//
// TODO:
// - Separate agent goals from planner.
// - Finalize real world state interface with mock realWorld implementation.
// - Clean up and document.
//
// This code has been inspired by:
// https://github.com/jamiecollinson/go-goap
//
// WARNING: Work in progress!
package aiplanner

import "log"

// Planner implements a simple GOAP.
type Planner struct {
	Actions    []Action     // All known actions.
	WorldState *realWorld   // Queryable world state.
	Goals      PlannerState // This should be separate / the agent.
}

// NewPlanner returns a new planner.
func NewPlanner() *Planner {
	return &Planner{
		WorldState: &realWorld{},
		Goals:      make(PlannerState),
	}
}

// possibleActions returns all actions that can run given the hypothetical worldstate.
func (a *Planner) possibleActions(world WorldState) []Action {
	validActions := []Action{}
	for _, action := range a.Actions {
		if action.CanRunIf(*a, world) {
			validActions = append(validActions, action)
		}
	}
	return validActions
}

// ValidatePlan plays through the given plan and checks if it is still valid.
func (a *Planner) ValidatePlan(currentPlan Plan) bool {
	return a.validatePlan(currentPlan, a.WorldState.Fork())
}

// validatePlan plays through the given plan and checks if it is still valid.
func (a *Planner) validatePlan(currentPlan Plan, world WorldFork) bool {
	var err error
	for _, action := range currentPlan {
		if !action.CanRunIf(*a, world) {
			return false
		}
		world, err = action.Simulate(*a, world)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// Check if goals are met in the resulting world state.
		// TODO: Check if we would have remaining actions that are not neccessary for
		// the current goal.
		if world.Contains(a.Goals) { // newAgent.goalsMet()
			return true
		}
	}
	return false
}

// GetPlans returns all possible plans.
func (a *Planner) GetPlans(currentPlan Plan) []Plan {
	return a.getPlans(currentPlan, a.WorldState.Fork())
}

func (a *Planner) getPlans(currentPlan Plan, world WorldFork) []Plan {
	var results []Plan
	for _, action := range a.possibleActions(world) {
		newPlan := append(currentPlan, action)
		// TODO: Discard plan if cost is now greater than the complete plan with the lowest cost?

		newWorld, err := action.Simulate(*a, world)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		// Check if goals are met.
		if newWorld.Contains(a.Goals) {
			results = append(results, newPlan)
		} else {
			results = append(results, a.getPlans(newPlan, newWorld)...)
		}
	}

	return results
}

// GetBestPlan returns the best, cheapest plan.
func (a *Planner) GetBestPlan() (Plan, int) {
	// NOTE: This is a rather greedy approach.
	//
	// In theory we can keep track of the shortest complete plan and skip all
	// plans that are longer (or just as long but incomplete) in the plan
	// generation function.
	plans := a.GetPlans(Plan{})
	var bestPlan Plan
	bestCost := 99999
	for _, plan := range plans {
		if cost := plan.Cost(); cost < bestCost {
			bestPlan = plan
			bestCost = cost
		}
	}
	return bestPlan, bestCost
}

// Plan represents a series of actions.
// TODO:
// - Implement cost caching.
// - Implement removing actions on partial completion?
// - Implement append function (+ cost caching).
type Plan []Action

// Cost returns the sum of the cost of all actions.
func (p *Plan) Cost() int {
	var cost int
	for _, action := range *p {
		cost += action.Cost()
	}
	return cost
}
