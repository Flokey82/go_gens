package aitree

// Sequence implements a sequence decorator node.
type Sequence struct {
	Name     string
	Children []Node
	Current  int // Current active child node
}

// NewSequence returns a new sequence decorator node.
func NewSequence(name string) *Sequence {
	return &Sequence{
		Name: name,
	}
}

// Append a new child node.
func (node *Sequence) Append(c Node) {
	node.Children = append(node.Children, c)
}

// Tick will attempt to run all children in the sequence and returns success if all complete successfully.
func (node *Sequence) Tick() State {
	if node.Current < len(node.Children) {
		result := node.Children[node.Current].Tick()
		if result == StateRunning {
			return StateRunning
		}
		if result == StateFailure {
			node.Current = 0
			return StateFailure
		}
		node.Current++
		if node.Current < len(node.Children) {
			return StateRunning
		}
		node.Current = 0
		return StateSuccess
	}
	return StateSuccess
}

// Selector implements a selector decorator node.
type Selector struct {
	Name     string
	Children []Node
	Current  int // Current active child node
}

// NewSelector returns a new selector decorator node.
func NewSelector(name string) *Selector {
	return &Selector{
		Name: name,
	}
}

// Append a new child node.
func (node *Selector) Append(c Node) {
	node.Children = append(node.Children, c)
}

// Tick will attempt to run each child sequentially and returns success if one of them succeeds.
func (node *Selector) Tick() State {
	if node.Current < len(node.Children) {
		result := node.Children[node.Current].Tick()
		if result == StateRunning {
			return StateRunning
		}
		if result == StateSuccess {
			node.Current = 0
			return StateSuccess
		}
		node.Current++
		if node.Current < len(node.Children) {
			return StateRunning
		}
		node.Current = 0
		return StateFailure
	}
	return StateFailure
}
