package gamerogueish

import (
	"fmt"
	"math/rand"
)

func (g *Game) decideAction(e *Entity) {
	if e.IsDead() {
		return
	}

	// If we are next to the player, attack.
	if e.X == g.player.X && e.Y == g.player.Y {
		g.AddMessage(fmt.Sprintf("%s attacks %s", e.Name, g.player.Name))
		e.Attack(g, g.player)
		return
	}

	// If we can see the player, move towards them.
	if g.FOV.IsInRadius(e.X, e.Y, g.player.X, g.player.Y) {
		moveTowards(g.World, e, g.player)
		return
	}
	// Otherwise, wander.
	wander(g.World, e)
}

func wander(w *World, e *Entity) {
	// Decide on a random direction, see if we can enter the tile.
	// 10 attempts.
	for i := 0; i < 10; i++ {
		dx := rand.Intn(3) - 1
		dy := rand.Intn(3) - 1
		if w.CanMoveTo(e.X+dx, e.Y+dy) {
			e.X += dx
			e.Y += dy
			return
		}
	}
}

func moveTowards(w *World, e, target *Entity) {
	// Calculate the direction to the target.
	dx := target.X - e.X
	dy := target.Y - e.Y
	// If we are already next to the target, we are done.
	if dx*dx+dy*dy == 1 {
		return
	}
	// Normalize the direction.
	if dx > 0 {
		dx = 1
	} else if dx < 0 {
		dx = -1
	}
	if dy > 0 {
		dy = 1
	} else if dy < 0 {
		dy = -1
	}
	if w.CanMoveTo(e.X+dx, e.Y+dy) {
		e.X += dx
		e.Y += dy
	}
	// TODO: Pathfinding.
}
