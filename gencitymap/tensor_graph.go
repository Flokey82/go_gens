package gencitymap

import (
	"log"
	"math"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

type Intersection struct {
	point    vectors.Vec2
	segments []Segment
}

// Node located along any intersection or point along the simplified road polylines.
type Node struct {
	value     vectors.Vec2
	segments  []vectors.Segment
	neighbors []*Node
	adj       []*Node
}

func (n *Node) addSegment(s vectors.Segment) {
	n.segments = append(n.segments, s)
}

func (n *Node) addNeighbor(node *Node) {
	if node != n {
		n.neighbors = append(n.neighbors, node)
		node.neighbors = append(node.neighbors, n)
	}
}

func (n *Node) removeNeighbor(node *Node) {
	for i, neighbor := range n.neighbors {
		if neighbor == node {
			n.neighbors = append(n.neighbors[:i], n.neighbors[i+1:]...)
			break
		}
	}

	for i, neighbor := range node.neighbors {
		if neighbor == n {
			node.neighbors = append(node.neighbors[:i], node.neighbors[i+1:]...)
			break
		}
	}
}

type Graph struct {
	nodes         []*Node
	intersections []vectors.Vec2
}

// NewGraph creates a graph from a set of streamlines.
// Finds all intersections, and creates a list of Nodes.
func NewGraph(streamlines [][]vectors.Vec2, dstep float64, deleteDangling bool) *Graph {
	// Find all intersections
	log.Println("Finding intersections")
	intersections := findAllIntersections(streamlinesToSegment(streamlines))
	// boo := bush(streamlinesToSegment(streamlines), nil)
	// intersections := boo.Run()
	log.Println("Found intersections:", len(intersections))

	quadtree := newFakeQuadtree()
	nodeAddRadius := 0.01

	log.Println("Adding nodes to quadtree")

	// Add all segment start and endpoints
	for idx := range streamlines {
		streamline := streamlines[idx]
		for i := 0; i < len(streamline); i++ {
			node := &Node{value: streamline[i]}
			if i > 0 {
				node.addSegment(vectorsToSegment(streamline[i-1], streamline[i]))
			}
			if i < len(streamline)-1 {
				node.addSegment(vectorsToSegment(streamline[i], streamline[i+1]))
			}
			fuzzyAddToQuadtree(quadtree, node, nodeAddRadius)
		}
	}
	log.Println("Intersections:", len(intersections))

	// Add all intersections
	for _, intersection := range intersections {
		node := &Node{value: vectors.Vec2{X: intersection.Point.X, Y: intersection.Point.Y}}
		for _, s := range intersection.Segments {
			node.addSegment(s)
		}
		fuzzyAddToQuadtree(quadtree, node, nodeAddRadius)
	}
	log.Println("Done adding nodes to quadtree")

	// For each simplified streamline, build list of nodes in order along streamline
	for _, streamline := range streamlines {
		for i := 0; i < len(streamline)-1; i++ {
			nodesAlongSegment := getNodesAlongSegment(vectorsToSegment(streamline[i], streamline[i+1]), quadtree, nodeAddRadius, dstep)
			if len(nodesAlongSegment) > 1 {
				for j := 0; j < len(nodesAlongSegment)-1; j++ {
					nodesAlongSegment[j].addNeighbor(nodesAlongSegment[j+1])
				}
			} else {
				log.Println("Error Graph.js: segment with less than 2 nodes")
			}
		}
	}

	// Remove dangling nodes
	for _, node := range quadtree.All() {
		if !deleteDangling || !deleteDanglingNodes(node, quadtree) {
			node.adj = make([]*Node, len(node.neighbors))
			copy(node.adj, node.neighbors)
		}
	}
	log.Println("Done removing dangling nodes")

	// Create a list of nodes
	var nodes []*Node
	for _, node := range quadtree.All() {
		if !deleteDangling || len(node.neighbors) > 1 {
			nodes = append(nodes, node)
		}
	}
	log.Println("Done creating list of nodes")

	// Create a list of intersections
	var intersections2 []vectors.Vec2
	for i := range intersections {
		intersections2 = append(intersections2, intersections[i].Point)
	}
	log.Println("Done creating list of intersections")

	return &Graph{
		nodes:         nodes,
		intersections: intersections2,
	}
}

// deleteDanglingNodes removes nodes that are not connected to any other nodes.
// Remove dangling edges from graph to facilitate polygon finding.
func deleteDanglingNodes(n *Node, quadtree QuadTree) bool {
	if len(n.neighbors) == 1 {
		quadtree.Remove(n)
		for _, neighbor := range n.neighbors {
			neighbor.removeNeighbor(n)
			deleteDanglingNodes(neighbor, quadtree)
		}
		return true
	}
	return false
}

// getNodesAlongSegment returns all nodes along a segment.
// Given a segment, step along segment and find all nodes along it.
func getNodesAlongSegment(segment vectors.Segment, quadtree QuadTree, radius float64, step float64) []*Node {
	// Walk dstep along each streamline, adding nodes within dstep/2
	// and connected to this streamline (fuzzy - nodeAddRadius) to list, removing from
	// quadtree and adding them all back at the end
	var foundNodes, nodesAlongSegment []*Node

	start := vectors.Vec2{X: segment.Start.X, Y: segment.Start.Y}
	end := vectors.Vec2{X: segment.End.X, Y: segment.End.Y}

	differenceVector := end.Sub(start)
	step = math.Min(step, differenceVector.Len()/2) // Min of 2 step along vector
	steps := int(math.Ceil(differenceVector.Len() / step))
	// differenceVectorLength := differenceVector.Len()

	for i := 0; i <= steps; i++ {
		currentPoint := start.Add(differenceVector.Mul(float64(i) / float64(steps)))

		// Order nodes, not by 'closeness' but by dot product
		var nodesToAdd []*Node
		closestNode := quadtree.Find(currentPoint, radius+step/2)

		for closestNode != nil {
			quadtree.Remove(closestNode)
			foundNodes = append(foundNodes, closestNode)

			var nodeOnSegment bool
			for _, s := range closestNode.segments {
				if fuzzySegmentsEqual(s, segment, 0.0001) {
					nodeOnSegment = true
					break
				}
			}
			if nodeOnSegment {
				nodesToAdd = append(nodesToAdd, closestNode)
			}
			closestNode = quadtree.Find(currentPoint, radius+step/2)
		}

		sort.Slice(nodesToAdd, func(i, j int) bool {
			return dotProductToSegment(nodesToAdd[i], start, differenceVector) < dotProductToSegment(nodesToAdd[j], start, differenceVector)
		})
		nodesAlongSegment = append(nodesAlongSegment, nodesToAdd...)
	}

	quadtree.Add(foundNodes...)
	return nodesAlongSegment
}

func dotProductToSegment(node *Node, start, differenceVector vectors.Vec2) float64 {
	return differenceVector.Dot(node.value.Sub(start))
}

func fuzzySegmentsEqual(s1 vectors.Segment, s2 vectors.Segment, tolerance float64) bool {
	// From
	if s1.Start.X-s2.Start.X > tolerance {
		return false
	}

	if s1.Start.Y-s2.Start.Y > tolerance {
		return false
	}

	// To
	if s1.End.X-s2.End.X > tolerance {
		return false
	}

	if s1.End.Y-s2.End.Y > tolerance {
		return false
	}

	return true
}

func (g *Graph) dotProductToSegment(node *Node, start vectors.Vec2, differenceVector vectors.Vec2) float64 {
	dotVector := node.value.Sub(start)
	return differenceVector.Dot(dotVector)
}

func fuzzyAddToQuadtree(quadtree QuadTree, node *Node, radius float64) {
	// Only add if there isn't a node within radius
	// Remember to check for double radius when querying tree, or point might be missed
	existingNode := quadtree.Find(node.value, radius)
	if existingNode == nil {
		quadtree.Add(node)
	} else {
		for _, neighbor := range node.neighbors {
			existingNode.addNeighbor(neighbor)
		}

		for _, segment := range node.segments {
			existingNode.addSegment(segment)
		}
	}
}

func streamlinesToSegment(streamlines [][]vectors.Vec2) []vectors.Segment {
	var out []vectors.Segment
	for _, s := range streamlines {
		for i := 0; i < len(s)-1; i++ {
			out = append(out, vectorsToSegment(s[i], s[i+1]))
		}
	}
	return out
}

func vectorsToSegment(v1, v2 vectors.Vec2) vectors.Segment {
	return vectors.Segment{
		Start: v1,
		End:   v2,
	}
}

type intersection struct {
	Point    vectors.Vec2
	Segments []vectors.Segment
}

func findAllIntersections(streamlines []vectors.Segment) []intersection {
	// Find all intersections
	var intersections []intersection
	for i := 0; i < len(streamlines); i++ {
		s1 := streamlines[i]
		for j := i + 1; j < len(streamlines); j++ {
			s2 := streamlines[j]
			if ok, p := s1.Intersects(s2); ok {
				intersections = append(intersections, intersection{Point: p, Segments: []vectors.Segment{s1, s2}})
			}
		}
	}
	return intersections
}

type QuadTree interface {
	Add(node ...*Node)
	Remove(node *Node)
	Find(point vectors.Vec2, radius float64) *Node
	Search(point vectors.Vec2, radius float64) []*Node
	All() []*Node
}

type fakeQuadtree struct {
	nodes []*Node
}

func newFakeQuadtree() *fakeQuadtree {
	return &fakeQuadtree{
		nodes: make([]*Node, 0),
	}
}

func (f *fakeQuadtree) Add(node ...*Node) {
	f.nodes = append(f.nodes, node...)
}

func (f *fakeQuadtree) Remove(node *Node) {
	for i, n := range f.nodes {
		if n == node {
			f.nodes = append(f.nodes[:i], f.nodes[i+1:]...)
			return
		}
	}
}

func (f *fakeQuadtree) Find(point vectors.Vec2, radius float64) *Node {
	for _, n := range f.nodes {
		if n.value.DistanceTo(point) < radius {
			return n
		}
	}
	return nil
}

func (f *fakeQuadtree) Search(point vectors.Vec2, radius float64) []*Node {
	var out []*Node
	for _, n := range f.nodes {
		if n.value.DistanceTo(point) < radius {
			out = append(out, n)
		}
	}
	return out
}

func (f *fakeQuadtree) All() []*Node {
	return f.nodes
}
