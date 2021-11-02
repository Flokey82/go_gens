// Package aitree implements a very basic behavior tree.
package aitree

// State represents the outcome of an operation
type State int

// Valid outcomes.
const (
	StateRunning State = iota
	StateFailure
	StateSuccess
)

// Node represents a node within the tree.
type Node interface {
	Tick() State
}

// Tree implements a simple behavior tree.
type Tree struct {
	Root Node
}

// New returns a new tree.
func New() *Tree {
	return &Tree{}
}

// Tick will advance the tree by a single tick.
func (t *Tree) Tick() State {
	if t.Root == nil {
		return StateFailure
	}
	return t.Root.Tick()
}
