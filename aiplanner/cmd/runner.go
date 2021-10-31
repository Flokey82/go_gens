package main

import (
	"fmt"
	"github.com/Flokey82/go_gens/aiplanner"
)

func main() {
	p := aiplanner.NewPlanner()

	emptyState := aiplanner.PlannerState{}

	fightEnemyHandToHandAction := aiplanner.NewAction(
		"fight enemy hand to hand", 3,
		aiplanner.PlannerState{"enemyInMeelee": true},
		aiplanner.PlannerState{"killedEnemy": true},
	)
	moveNextToEnemyAction := aiplanner.NewAction(
		"move next to enemy", 4, emptyState,
		aiplanner.PlannerState{"enemyInMeelee": true},
	)

	pickupGunAction := aiplanner.NewAction(
		"pick up gun", 1, emptyState,
		aiplanner.PlannerState{"armed": true},
	)
	shootEnemyAction := aiplanner.NewAction(
		"shoot enemy", 1,
		aiplanner.PlannerState{"enemyInSight": true, "armed": true},
		aiplanner.PlannerState{"killedEnemy": true, "threatened": false},
	)
	moveToEnemyAction := aiplanner.NewAction(
		"move to enemy", 3, emptyState,
		aiplanner.PlannerState{"enemyInSight": true},
	)

	findAppleAction := aiplanner.NewAction("find apple", 1,
		aiplanner.PlannerState{"threatened": false},
		aiplanner.PlannerState{"hasFood": true},
	)
	eatAppleAction := aiplanner.NewAction("eat food", 1,
		aiplanner.PlannerState{"threatened": false, "hasFood": true},
		aiplanner.PlannerState{"hungry": false},
	)

	p.Actions = []aiplanner.Action{
		fightEnemyHandToHandAction,
		pickupGunAction,
		shootEnemyAction,
		findAppleAction,
		eatAppleAction,
		moveToEnemyAction,
		moveNextToEnemyAction,
	}

	p.WorldState.Set("threatened", true)
	p.WorldState.Set("hungry", true)
	p.Goals = aiplanner.PlannerState{
		"hungry": false,
	}

	bestPlan, bestScore := p.GetBestPlan()
	fmt.Println(bestPlan, bestScore)
	fmt.Println(p.ValidatePlan(bestPlan))
}
