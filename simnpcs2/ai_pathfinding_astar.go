package simnpcs2

import (
	"errors"
	"fmt"
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

func (w *World) Pathfind(start, end *vectors.Vec2) ([]int, error) {
	// Convert the start and end vectors to cell indices.
	startIdx := w.PosToCellIdx(start)
	endIdx := w.PosToCellIdx(end)

	// Check if the start and end indices are valid.
	if err := w.CheckIdxReachable(startIdx); err != nil {
		return nil, err
	}
	if err := w.CheckIdxReachable(endIdx); err != nil {
		return nil, err
	}

	// Get all obstacles.
	// TODO: Maybe cache this?
	var obstacles []Node
	for i, blocked := range w.Cells {
		if blocked {
			obstacles = append(obstacles, Node{
				X: i % w.Width,
				Y: i / w.Width,
			})
		}
	}

	// Create a new A* instance.
	astar, err := newAStar(Config{
		GridWidth:    w.Width,
		GridHeight:   w.Height,
		InvalidNodes: obstacles,
	})
	if err != nil {
		return nil, err
	}

	// Calculate the path.
	startNode := Node{X: startIdx % w.Width, Y: startIdx / w.Width}
	endNode := Node{X: endIdx % w.Width, Y: endIdx / w.Width}
	path, err := astar.FindPath(startNode, endNode)
	if err != nil {
		return nil, err
	}

	// Convert the path to cell indices (and reverse it).
	pathIdx := make([]int, len(path))
	for i := range path {
		node := path[len(path)-1-i]
		pathIdx[i] = node.Y*w.Width + node.X
	}

	return pathIdx, nil
}

// Config holds important settings
// to perform the calculation
//
// GridWidth and GridHeight are required and represents
// the size of the grid
//
// InvalidNodes can be used to add not accessible nodes like obstacles etc.
// WeightedNodes can be used to add nodes to be avoided like mud or mountains
type Config struct {
	GridWidth, GridHeight int
	InvalidNodes          []Node
	WeightedNodes         []Node
}

type astar struct {
	config               Config
	openList, closedList List
	startNode, endNode   Node
}

var ErrGridSize = errors.New("grid size must be at least 2x2")

// newAStar creates a new astar instance
// NOTE: This is a temporary drop from this fantastic package:
// https://github.com/jpierer/astar
func newAStar(config Config) (*astar, error) {
	if config.GridWidth < 2 || config.GridHeight < 2 {
		return nil, ErrGridSize
	}
	a := &astar{config: config}
	return a.init(), nil
}

// init initialised needed properties
// internal function
func (a *astar) init() *astar {
	// add invalidNodes directly to the closedList
	a.closedList.Add(a.config.InvalidNodes...)
	return a
}

// H caluclates the absolute distance between
// nodeA and nodeB calculates by the manhattan distance
func (a *astar) H(nodeA Node, nodeB Node) int {
	absX := math.Abs(float64(nodeA.X - nodeB.X))
	absY := math.Abs(float64(nodeA.Y - nodeB.Y))
	return int(absX + absY)
}

// GetNeighborNodes calculates the next neighbors of the given node
// if a neighbor node is not accessible the node will be ignored
func (a *astar) GetNeighborNodes(node Node) []Node {
	var neighborNodes []Node

	upNode := Node{X: node.X, Y: node.Y + 1, parent: &node}
	if a.isAccessible(upNode) {
		neighborNodes = append(neighborNodes, upNode)
	}

	downNode := Node{X: node.X, Y: node.Y - 1, parent: &node}
	if a.isAccessible(downNode) {
		neighborNodes = append(neighborNodes, downNode)
	}

	leftNode := Node{X: node.X - 1, Y: node.Y, parent: &node}
	if a.isAccessible(leftNode) {
		neighborNodes = append(neighborNodes, leftNode)
	}

	rightNode := Node{X: node.X + 1, Y: node.Y, parent: &node}
	if a.isAccessible(rightNode) {
		neighborNodes = append(neighborNodes, rightNode)
	}

	return neighborNodes
}

// isAccessible checks if the node is reachable in the grid
// and is not in the invalidNodes slice
func (a *astar) isAccessible(node Node) bool {

	// if node is out of bound
	if node.X < 0 || node.Y < 0 || node.X > a.config.GridWidth-1 || node.Y > a.config.GridHeight-1 {
		return false
	}

	// check if the node is in the closedList
	// the predefined invalidNodes are also in this list
	if a.closedList.Contains(node) {
		return false
	}

	return true
}

// IsEndNode checks if the given node has
// equal node coordinates with the end node
func (a *astar) IsEndNode(checkNode, endNode Node) bool {
	return checkNode.X == endNode.X && checkNode.Y == endNode.Y
}

var (
	ErrNoPathFound = errors.New("no path found")
	ErrNoNodeFound = errors.New("no node found")
)

// FindPath starts the a* algorithm for the given start and end node
// The return value will be the fastest way represented as a nodes slice
//
// If no path was found it returns nil and an error
func (a *astar) FindPath(startNode, endNode Node) ([]Node, error) {
	a.startNode = startNode
	a.endNode = endNode

	defer func() {
		a.openList.Clear()
		a.closedList.Clear()
	}()

	a.openList.Add(startNode)

	for !a.openList.IsEmpty() {

		currentNode, err := a.openList.GetMinFNode()
		if err != nil {
			return nil, fmt.Errorf("cannot get minF node %v", err)
		}

		a.openList.Remove(currentNode)
		a.closedList.Add(currentNode)

		// we found the path
		if a.IsEndNode(currentNode, endNode) {
			return a.getNodePath(currentNode), nil
		}

		neighbors := a.GetNeighborNodes(currentNode)
		for _, neighbor := range neighbors {
			if a.closedList.Contains(neighbor) {
				continue
			}
			a.calculateNode(&neighbor)
			if !a.openList.Contains(neighbor) {
				a.openList.Add(neighbor)
			}
		}

	}

	return nil, ErrNoPathFound
}

// calculateNode calculates the F, G and H value for the given node
func (a *astar) calculateNode(node *Node) {
	node.g++

	// check for special node weighting
	for _, wNode := range a.config.WeightedNodes {
		if node.X == wNode.X && node.Y == wNode.Y {
			node.g = node.g + wNode.Weighting
			break
		}
	}

	node.h = a.H(*node, a.endNode)
	node.f = node.g + node.h
}

// getNodePath returns the chain of parent nodes
// the given node will be still included in the nodes slice
func (a *astar) getNodePath(currentNode Node) []Node {
	var nodePath []Node
	nodePath = append(nodePath, currentNode)
	for {
		if currentNode.parent == nil {
			break
		}

		parentNode := *currentNode.parent

		// if the end of node chain
		if parentNode.parent == nil {
			break
		}

		nodePath = append(nodePath, parentNode)
		currentNode = parentNode
	}
	return nodePath
}

// Node represents a simple node
// X and Y represents the nodes coordinates on the grid
//
// IMPORTANT: The grid coordinates starts on the "bottom left" -> X:0 / Y:0
//
// With the Weighting value you can set the nodes heavy grade
// so a node with mud or water are heavier as gras or street
type Node struct {
	f, g, h   int
	X, Y      int
	Weighting int
	parent    *Node
}

// String returns formatted values of the node
func (n Node) String() string {
	return fmt.Sprintf("Node [X:%d Y:%d F:%d G:%d H:%d]", n.X, n.Y, n.f, n.g, n.h)
}

// List represents a list of nodes
type List struct {
	nodes []Node
}

// NewList creates a new list
func NewList() *List {
	return &List{}
}

// Add one or more nodes to the list
func (l *List) Add(nodes ...Node) {
	l.nodes = append(l.nodes, nodes...)
}

// All returns the full list of nodes
func (l *List) All() []Node {
	return l.nodes
}

// Remove a node from the list
// if the node is not found we do nothing
func (l *List) Remove(removeNode Node) {
	index := l.GetIndex(removeNode)
	if index >= 0 {
		l.nodes = append(l.nodes[:index], l.nodes[index+1:]...)
	}
}

// GetIndex returns the index of the node in the list
// if the node is not found the return value is -1
func (l *List) GetIndex(searchNode Node) int {
	for index, node := range l.nodes {
		if node.X == searchNode.X && node.Y == searchNode.Y {
			return index
		}
	}
	return -1
}

// Contains check if a node is in the list
func (l *List) Contains(searchNode Node) bool {
	return l.GetIndex(searchNode) >= 0
}

// IsEmpty returns if the nodes list has nodes or not
func (l *List) IsEmpty() bool {
	return len(l.nodes) == 0
}

// Clear removes all nodes from the list
func (l *List) Clear() {
	l.nodes = []Node{}
}

// GetIndexOfMinF returns the index of the nodes list
// with the smallest node.F value
//
// if no node is found it returns -1
func (l *List) GetIndexOfMinF() int {
	lastNode := Node{}
	lastNodeIndex := -1
	for index, node := range l.nodes {
		if lastNodeIndex == -1 || node.f < lastNode.f {
			lastNode = node
			lastNodeIndex = index
		}
	}
	return lastNodeIndex
}

// GetMinFNode returns the node with the smallest node.F value
func (l *List) GetMinFNode() (Node, error) {
	minFIndex := l.GetIndexOfMinF()
	if minFIndex == -1 {
		return Node{}, ErrNoNodeFound
	}
	return l.nodes[minFIndex], nil
}
