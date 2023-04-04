package gencitymap

import (
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
)

type PolygonParams struct {
	MaxLength      int
	MinArea        float64
	ShrinkSpacing  float64
	ChanceNoDivide float64
}

// PolygonFinder finds polygons in a graph, used for finding lots and parks.
type PolygonFinder struct {
	Polygons        [][]vectors.Vec2
	ShrunkPolygons  [][]vectors.Vec2
	DividedPolygons [][]vectors.Vec2
	toShrink        [][]vectors.Vec2
	resolveShrink   func()
	toDivide        [][]vectors.Vec2
	resolveDivide   func()
	Nodes           []*Node
	Params          PolygonParams
	TensorField     *TensorField
}

func NewPolygonFinder(nodes []*Node, params PolygonParams, tensorField *TensorField) *PolygonFinder {
	return &PolygonFinder{
		Nodes:       nodes,
		Params:      params,
		TensorField: tensorField,
	}
}

func (p *PolygonFinder) getPolygons() [][]vectors.Vec2 {
	if len(p.DividedPolygons) > 0 {
		return p.DividedPolygons
	}

	if len(p.ShrunkPolygons) > 0 {
		return p.ShrunkPolygons
	}

	return p.Polygons
}

func (p *PolygonFinder) Reset() {
	p.toShrink = nil
	p.toDivide = nil
	p.Polygons = nil
	p.ShrunkPolygons = nil
	p.DividedPolygons = nil
}

func (p *PolygonFinder) Update() bool {
	change := false
	if len(p.toShrink) > 0 {
		resolve := len(p.toShrink) == 1
		if p.stepShrink(p.toShrink[len(p.toShrink)-1]) {
			change = true
		}

		if resolve {
			p.resolveShrink()
		}
	}

	if len(p.toDivide) > 0 {
		resolve := len(p.toDivide) == 1
		if p.stepDivide(p.toDivide[len(p.toDivide)-1]) {
			change = true
		}

		if resolve {
			p.resolveDivide()
		}
	}

	return change
}

// Shrink shrinks the polygons by the given amount.
// Properly shrink polygon so the edges are all the same distance from the road.
func (p *PolygonFinder) Shrink(animate bool) {
	if len(p.Polygons) == 0 {
		p.findPolygons()
	}

	/*if animate {
		if len(p.Polygons) == 0 {
			return
		}

		p.toShrink = p.Polygons
		p.resolveShrink = resolve
	} else {*/
	p.ShrunkPolygons = nil
	for _, poly := range p.Polygons {
		p.stepShrink(poly)
	}
	//resolve()
	//}
}

func (p *PolygonFinder) stepShrink(polygon []vectors.Vec2) bool {
	shrunk := ResizeGeometry(polygon, -p.Params.ShrinkSpacing, true)
	if len(shrunk) > 0 {
		p.ShrunkPolygons = append(p.ShrunkPolygons, shrunk)
		return true
	}

	return false
}

func (p *PolygonFinder) Divide(animate bool) {
	if len(p.Polygons) == 0 {
		p.findPolygons()
	}

	polygons := p.Polygons
	if len(p.ShrunkPolygons) > 0 {
		polygons = p.ShrunkPolygons
	}

	/*if animate {
		if len(polygons) == 0 {
			return
		}

		p.toDivide = polygons
		p.resolveDivide = resolve
	} else {*/
	p.DividedPolygons = nil
	for _, poly := range polygons {
		p.stepDivide(poly)
	}
	//resolve()
	//}
}

func (p *PolygonFinder) stepDivide(polygon []vectors.Vec2) bool {
	// TODO need to filter shrunk polygons using aspect ratio, area
	// this skips the filter in PolygonUtil.subdividePolygon
	if p.Params.ChanceNoDivide > 0 && rand.Float64() < p.Params.ChanceNoDivide {
		p.DividedPolygons = append(p.DividedPolygons, polygon)
		return true
	}

	divided := SubdividePolygon(polygon, p.Params.MinArea)
	if len(divided) > 0 {
		p.DividedPolygons = append(p.DividedPolygons, divided...)
		return true
	}

	return false
}

func (p *PolygonFinder) findPolygons() {
	// Node
	// x, y, value (Vector2), adj (list of node refs)
	// Gonna edit adj for now

	// Walk a clockwise path until polygon found or limit reached
	// When we find a polygon, mark all edges as traversed (in particular direction)
	// Each edge separates two polygons
	// If edge already traversed in this direction, this polygon has already been found
	p.ShrunkPolygons = nil
	p.DividedPolygons = nil
	polygons := make([][]vectors.Vec2, 0)

	for _, node := range p.Nodes {
		if len(node.adj) < 2 {
			continue
		}

		for _, nextNode := range node.adj {
			polygon := p.recursiveWalk([]*Node{node, nextNode})
			if polygon != nil && len(polygon) < p.Params.MaxLength {
				p.removePolygonAdjacencies(polygon)
				var polygonPoints []vectors.Vec2
				for _, n := range polygon {
					polygonPoints = append(polygonPoints, n.value)
				}
				polygons = append(polygons, polygonPoints)
			}
		}
	}

	p.Polygons = p.filterPolygonsByWater(polygons)
}

func (p *PolygonFinder) filterPolygonsByWater(polygons [][]vectors.Vec2) [][]vectors.Vec2 {
	out := make([][]vectors.Vec2, 0)
	for _, poly := range polygons {
		averagePoint := AveragePoint(poly)
		if p.TensorField.onLand(averagePoint) && !p.TensorField.inParks(averagePoint) {
			out = append(out, poly)
		}
	}

	return out
}

func (p *PolygonFinder) removePolygonAdjacencies(polygon []*Node) {
	for i := 0; i < len(polygon); i++ {
		node := polygon[i]
		nextNode := polygon[(i+1)%len(polygon)]
		index := indexOfNode(node.adj, nextNode)
		if index >= 0 {
			node.adj = append(node.adj[:index], node.adj[index+1:]...)
		}
	}
}

func indexOfNode(nodes []*Node, node *Node) int {
	for i, n := range nodes {
		if n == node {
			return i
		}
	}
	return -1
}

/*
func (p *PolygonFinder) recursiveWalk(path []*Node) []*Node {
	lastNode := path[len(path)-1]
	secondLastNode := path[len(path)-2]
	nextNode := p.findNextNode(lastNode, secondLastNode)
	if nextNode == nil {
		return nil
	}
	path = append(path, nextNode)
	if nextNode == path[0] {
		return path
	}
	return p.recursiveWalk(path)
}*/

func (p *PolygonFinder) recursiveWalk(visited []*Node) []*Node {
	if len(visited) >= p.Params.MaxLength {
		return nil
	}
	// TODO backtracking to find polygons with dead end roads inside them
	nextNode := p.getRightmostNode(visited[len(visited)-2], visited[len(visited)-1])
	if nextNode == nil {
		return nil // Currently ignores polygons with dead end inside
	}

	visitedIndex := indexOfNode(visited, nextNode)
	if visitedIndex >= 0 {
		return visited[visitedIndex:]
	}
	visited = append(visited, nextNode)
	return p.recursiveWalk(visited)
}

func (p *PolygonFinder) getRightmostNode(nodeFrom, nodeTo *Node) *Node {
	// We want to turn right at every junction
	if len(nodeTo.adj) == 0 {
		return nil
	}

	backwardsDifferenceVector := nodeFrom.value.Sub(nodeTo.value)
	transformAngle := math.Atan2(backwardsDifferenceVector.Y, backwardsDifferenceVector.X)

	var rightmostNode *Node
	smallestTheta := math.Pi * 2

	for _, nextNode := range nodeTo.adj {
		if nextNode != nodeFrom {
			nextVector := nextNode.value.Sub(nodeTo.value)
			nextAngle := math.Atan2(nextVector.Y, nextVector.X) - transformAngle
			if nextAngle < 0 {
				nextAngle += math.Pi * 2
			}

			if nextAngle < smallestTheta {
				smallestTheta = nextAngle
				rightmostNode = nextNode
			}
		}
	}

	return rightmostNode
}
