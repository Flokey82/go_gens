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
	dx, dy := findPath(w, e, target)
	if w.CanMoveTo(e.X+dx, e.Y+dy) {
		e.X += dx
		e.Y += dy
	}
	// TODO: Cache the path and adjust the monster speed, so that we have a chance
	// to escape.
}

// findPath finds a path from e to target and returns dx, dy to move in
// this turn.
func findPath(w *World, e, target *Entity) (int, int) {
	// New A* pathfinder.
	p := NewPathfinder(w, e.X, e.Y, target.X, target.Y)
	// Find the path.
	path := p.FindPath()
	// If there is no path, return 0, 0.
	// Make sure we don't return the exact same position.
	if len(path) <= 1 {
		return 0, 0
	}
	// Return the next step in the path.
	return path[1].X - e.X, path[1].Y - e.Y
}

// Node represents a node in the pathfinder.
type Node struct {
	Parent *Node
	X, Y   int
}

// Pathfinder represents a pathfinder.
type Pathfinder struct {
	world *World
	start *Node
	end   *Node
}

// NewPathfinder creates a new pathfinder.
func NewPathfinder(w *World, x, y, tx, ty int) *Pathfinder {
	return &Pathfinder{
		world: w,
		start: &Node{X: x, Y: y},
		end:   &Node{X: tx, Y: ty},
	}
}

// FindPath finds a path from the start to the end.
// NOTE: Poor man's A*.
// TODO: Use heuristics to prioritize the search to nodes that are closer to the
// end node.
func (p *Pathfinder) FindPath() []*Node {
	var closed []*Node
	open := []*Node{p.start}
	// While there are still nodes to check.
	for len(open) > 0 {
		// Get the first node from the open list.
		// NOTE: This is pretty inefficient to do it this way due to the
		// constant re-allocation. In theory we could allocate a slice with a
		// reasonable capacity and use an index to keep track of the position
		// and once we've exhausted the capacity, we can copy the last node to
		// the first position and set the index to 0 and truncate the slice.
		n := open[0]
		open = open[1:]
		// If this is the end node, we are done.
		if n.X == p.end.X && n.Y == p.end.Y {
			return p.reconstructPath(n)
		}
		// Add the node to the closed list.
		closed = append(closed, n)
		// Get the neighbors.
		neighbors := p.getNeighbors(n)
		// For each neighbor.
		for _, neighbor := range neighbors {
			// If the neighbor is in the closed list, skip it.
			if isInList(neighbor, closed) {
				continue
			}
			// If the neighbor is not in the open list, add it.
			if !isInList(neighbor, open) {
				open = append(open, neighbor)
			}
		}
	}
	// No path found.
	return nil
}

// reconstructPath reconstructs the path from the end node.
func (p *Pathfinder) reconstructPath(n *Node) []*Node {
	path := []*Node{n}
	for n.Parent != nil {
		n = n.Parent
		path = append([]*Node{n}, path...)
	}
	return path
}

// getNeighbors returns the neighbors of the node.
// TODO: Re-use a pre-allocated slice for the neighbors.
func (p *Pathfinder) getNeighbors(n *Node) []*Node {
	var neighbors []*Node
	// Check each direction.
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			// Skip the center.
			if dx == 0 && dy == 0 {
				continue
			}
			// If we can move to the tile, add it as a neighbor.
			if p.world.CanMoveTo(n.X+dx, n.Y+dy) {
				neighbors = append(neighbors, &Node{Parent: n, X: n.X + dx, Y: n.Y + dy})
			}
		}
	}
	return neighbors
}

// isInList returns true if the node is in the list.
func isInList(n *Node, list []*Node) bool {
	for _, node := range list {
		if node.X == n.X && node.Y == n.Y {
			return true
		}
	}
	return false
}
