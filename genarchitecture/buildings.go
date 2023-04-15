package genarchitecture

import (
	"log"
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/gengeometry"
	"github.com/Flokey82/go_gens/vectors"
)

func GenerateSampleCathedral() {
	// Generate a path for a building.
	procCrossPath := gengeometry.PlusShape{
		Width:     1,
		Length:    1,
		WingWidth: 0.35,
	}

	// Generate a mesh from a path.
	path := procCrossPath.GetPath()
	mesh, err := gengeometry.ExtrudePath(path, 0.2)
	if err != nil {
		log.Fatal(err)
	}
	roofMesh, err := gengeometry.TaperPath(path, 0.2)
	if err != nil {
		log.Println(err)
	} else {
		mesh.AddMesh(roofMesh, vectors.NewVec3(0, 0, 0.2))
	}

	// Now iterate over the sides and add a cylinder to each side.
	for i, s := range gengeometry.GetPathSides(path) {
		corner := gengeometry.RectangleShape{
			Width:  0.1,
			Length: 0.1,
		}
		heightCorner := 0.3
		cornerPath := corner.GetPath()
		cornerMesh, err := gengeometry.ExtrudePath(cornerPath, heightCorner)
		if err != nil {
			log.Println(err)
		} else {
			// Add the corner to the mesh.
			mesh.AddMesh(cornerMesh, vectors.NewVec3(s.Start.X-corner.Width/2, s.Start.Y-corner.Length/2, 0))
		}

		// Add a roof to the corner.
		roofMesh, err := gengeometry.TaperPath(cornerPath, 0.025)
		if err != nil {
			log.Println(err)
		} else {
			mesh.AddMesh(roofMesh, vectors.NewVec3(s.Start.X-corner.Width/2, s.Start.Y-corner.Length/2, heightCorner))
		}

		if i%3 == 0 {
			continue
		}

		// Add support struts in the middle of each side.
		strutPath := gengeometry.RectangleShape{
			Width:  0.025,
			Length: 0.025,
		}

		midPoint := s.Start.Add(s.End).Mul(0.5)

		// We add two struts.
		midPointA := s.Start.Add(midPoint).Mul(0.5)
		midPointB := midPoint.Add(s.End).Mul(0.5)

		for _, midP := range []vectors.Vec2{midPointA, midPointB, midPoint} {
			strutMesh, err := gengeometry.ExtrudePath(strutPath.GetPath(), 0.15)
			if err != nil {
				log.Println(err)
			} else {
				// Add the strut to the mesh.
				mesh.AddMesh(strutMesh, vectors.NewVec3(midP.X-strutPath.Width/2, midP.Y-strutPath.Length/2, 0))
			}

			// Add a roof to the strut.
			roofMesh, err = gengeometry.TaperPath(strutPath.GetPath(), 0.01)
			if err != nil {
				log.Println(err)
			} else {
				mesh.AddMesh(roofMesh, vectors.NewVec3(midP.X-strutPath.Width/2, midP.Y-strutPath.Length/2, 0.15))
			}
		}
	}

	// Save the mesh to a file.
	mesh.ExportToObj("test_2.obj")
}

const (
	ShapeTypeRect = iota
	ShapeTypePlus
	ShapeTypeCircle
	ShapeTypeTriangle
	ShapeTypeH
)

type Rule struct {
	Name         string   // Name of the rule.
	Shape        int      // The shape of the base.
	Roof         bool     // Whether the roof should be generated.
	Reorient     bool     // Whether the shape should be reoriented.
	Width        float64  // Width of the base.
	Length       float64  // Length of the base.
	Height       float64  // Height of the shape.
	Elevation    float64  // How much the roof is elevated from the base.
	Taper        float64  // How much the mesh is tapered.
	RulesSides   []string // Rules to apply to the sides.
	RulesCorners []string // Rules to apply to the corners.
}

func (r *Rule) GetShape() gengeometry.Shape {
	switch r.Shape {
	case ShapeTypeRect:
		return gengeometry.RectangleShape{
			Width:  r.Width,
			Length: r.Length,
		}
	case ShapeTypePlus:
		return gengeometry.PlusShape{
			Width:     r.Width,
			Length:    r.Length,
			WingWidth: 0.35,
		}
	case ShapeTypeCircle:
		return gengeometry.CircleShape{
			Radius: r.Width,
		}
	case ShapeTypeTriangle:
		return gengeometry.TriangleShape{
			Width:  r.Width,
			Length: r.Length,
		}
	case ShapeTypeH:
		return gengeometry.HShape{
			Width:     r.Width,
			Length:    r.Length,
			WingWidth: 0.35,
		}
	}
	return nil
}

var SampleRules = []*Rule{
	{
		Name:         "base",
		Shape:        ShapeTypePlus,
		Width:        2,
		Length:       2,
		Height:       0.2,
		Roof:         true,
		RulesSides:   []string{"wing"},
		RulesCorners: []string{"corner"},
	},
	{
		Name:         "wing",
		Shape:        ShapeTypeRect,
		Width:        0.35,
		Length:       0.35,
		Height:       0.2,
		Roof:         true,
		RulesSides:   []string{"strut"},
		RulesCorners: []string{"strut"},
	},
	{
		Name:         "corner",
		Shape:        ShapeTypeRect,
		Width:        0.1,
		Length:       0.1,
		Height:       0.2,
		Roof:         true,
		RulesSides:   []string{"strut"},
		RulesCorners: []string{"strut"},
	},
	{
		Name:         "strut",
		Shape:        ShapeTypeRect,
		Width:        0.025,
		Length:       0.025,
		Height:       0.2,
		RulesSides:   []string{},
		RulesCorners: []string{},
	},
}

var SampleRules2 = []*Rule{
	{
		Name:       "base",
		Shape:      ShapeTypeRect,
		Width:      2,
		Length:     2,
		Taper:      1,
		Height:     1,
		RulesSides: []string{"side"},
	},
	{
		Name:         "side",
		Shape:        ShapeTypeH,
		Width:        1.51,
		Length:       0.51,
		Height:       0.35,
		Roof:         true,
		Reorient:     true,
		RulesSides:   []string{"strut"},
		RulesCorners: []string{"strut"},
	},
	{
		Name:         "strut",
		Shape:        ShapeTypeRect,
		Width:        0.025,
		Length:       0.025,
		Height:       0.35,
		Reorient:     true,
		RulesSides:   []string{},
		RulesCorners: []string{},
	},
}

type RuleCollection struct {
	Rules map[string]*Rule
	All   []*Rule
	Root  *Rule
}

func NewRuleCollection() *RuleCollection {
	return &RuleCollection{
		Rules: make(map[string]*Rule),
	}
}

func (rc *RuleCollection) AddRule(r *Rule) {
	rc.Rules[r.Name] = r
	rc.All = append(rc.All, r)
}

type stackEntry struct {
	CurrentPos  vectors.Vec3
	CurrentDir  vectors.Vec2
	RuleToApply string
}

func (rc *RuleCollection) Run() *gengeometry.Mesh {
	// Create a stack for all the entries we process.
	stack := []stackEntry{}

	// Create a mesh to store the result.
	mesh := &gengeometry.Mesh{}

	// Add the root rule to the stack.
	stack = append(stack, stackEntry{
		CurrentPos:  vectors.NewVec3(0, 0, 0),
		RuleToApply: rc.Root.Name,
	})

	for i := 0; i < len(stack); i++ {
		log.Println("Stack size:", len(stack))

		// Get the entry, evaluate, and add new entries to the stack.
		entry := stack[i]
		rule := rc.Rules[entry.RuleToApply]
		log.Println("Applying rule:", rule.Name)

		// TODO: Change the orientation of the shape based on where the shape is placed.
		//	   ^
		//   _| |_
		// < _   _ >
		//    | |
		//     v
		path := rule.GetShape().GetPath()
		midPointParent := gengeometry.CenterOfPath(path)

		entry.CurrentPos.X -= rule.Width / 2
		entry.CurrentPos.Y -= rule.Length / 2
		entry.CurrentPos.Z += rule.Elevation

		if rule.Reorient {
			// Calculate angle from vector.
			curAngle := math.Atan2(entry.CurrentDir.Y, entry.CurrentDir.X)
			path = gengeometry.RotatePolygonAroundPoint(path, vectors.NewVec2(rule.Width/2, rule.Length/2), curAngle)
		}

		// Extrude the path.
		var extrudedMesh *gengeometry.Mesh
		var err error
		if rule.Taper > 0 {
			extrudedMesh, err = gengeometry.TaperPath(path, rule.Taper)
		} else {
			extrudedMesh, err = gengeometry.ExtrudePath(path, rule.Height)
		}
		if err != nil {
			log.Println(err)
		} else {
			mesh.AddMesh(extrudedMesh, entry.CurrentPos)
		}

		if rule.Roof {
			// TODO: There is a problem when adding a roof to a tapered shape... it will use
			// the bottom path instead of the top path (which is smaller).
			// This will look like a mushroom.

			// Add a roof to the path.
			roofMesh, err := gengeometry.TaperPath(path, 0.05)
			if err != nil {
				log.Println(err)
			} else {
				mesh.AddMesh(roofMesh, vectors.NewVec3(entry.CurrentPos.X, entry.CurrentPos.Y, rule.Height))
			}
		}

		// If there is a rule for the sides, add them to the stack.
		if len(rule.RulesSides) > 0 {
			sideRule := rule.RulesSides[rand.Intn(len(rule.RulesSides))]
			for _, s := range gengeometry.GetPathSides(path) {
				midPoint := s.Start.Add(s.End).Mul(0.5)
				// dir is the normal vector (e.g. 90° to the side) of the side (pointing outwards).
				dir := vectors.Normalize(s.End.Sub(s.Start)).Perpendicular()

				// Add the side to the stack.
				stack = append(stack, stackEntry{
					CurrentPos:  vectors.NewVec3(entry.CurrentPos.X+midPoint.X, entry.CurrentPos.Y+midPoint.Y, entry.CurrentPos.Z),
					CurrentDir:  dir,
					RuleToApply: sideRule,
				})
			}
		}

		currPos2 := midPointParent

		// If there is a rule for the corners, add them to the stack.
		if len(rule.RulesCorners) > 0 {
			cornerRule := rule.RulesCorners[rand.Intn(len(rule.RulesCorners))]
			for _, c := range path {
				dir := vectors.Normalize(c.Sub(currPos2))
				// Add the corner to the stack.
				stack = append(stack, stackEntry{
					CurrentPos:  vectors.NewVec3(entry.CurrentPos.X+c.X, entry.CurrentPos.Y+c.Y, entry.CurrentPos.Z),
					CurrentDir:  dir,
					RuleToApply: cornerRule,
				})
			}
		}
	}

	return mesh
}

// Graph grammar:
// A node represents something to be re-written or drawn.
// Rules apply to nodes of a type and define if and how they are re-written.
// A rule can return multiple nodes and either replace the original node or add new nodes.
// Also, each node may be re-written only once (if this is explicitly set) or multiple times.

type Node struct {
	ID       string      // The ID of the node.
	Parent   *Node       // The parent node.
	Children []*Node     // The child nodes.
	Data     interface{} // The data associated with the node.
	Replaced bool        // If true, this node will not be drawn, or re-written.
}

func newNode(id string, data interface{}) *Node {
	return &Node{
		ID:   id,
		Data: data,
	}
}

type NodeRule struct {
	ID          string              // The ID of the rule.
	F           func(*Node) []*Node // The function that is applied to the node and produces new nodes.
	ReplaceNode bool                // If true, the node will be replaced by the new nodes.
}

func (nr *NodeRule) Apply(node *Node) []*Node {
	res := nr.F(node)
	for _, n := range res {
		n.Parent = node
	}
	// If the rule replaces the node, set the replaced flag.
	if nr.ReplaceNode {
		node.Replaced = true
	}
	// Append the new nodes as children to the current node.
	node.Children = append(node.Children, res...)
	return res
}

type ShapeData struct {
	Shape       gengeometry.Shape
	CurrentPos  vectors.Vec3 // Center of the shape.
	CurrentDir  vectors.Vec2 // Direction of the shape.
	HeightScale float64      // The height scale of the shape.
	Reorient    bool         // If true, the shape will be re-oriented based on the direction of the node.
}

// A graph holds the result of a graph grammar.
// WARNING: This is not working properly yet!!!!
// We either run the evaluation until until we reached terminal nodes or we run it for a fixed number of iterations.
func Eval() *Node {
	var stack []*Node

	root := newNode("root", &ShapeData{
		Shape: &gengeometry.HShape{
			Width:     2,
			Length:    2,
			WingWidth: 0.55,
		},
		HeightScale: 0.2,
	})

	// Add the root node to the stack.
	stack = append(stack, root)

	rulesByName := map[string]*NodeRule{
		"root": &NodeRule{
			ID:          "root",
			ReplaceNode: false,
			F: func(node *Node) []*Node {
				dat := node.Data.(*ShapeData)
				path := dat.Shape.GetPath()
				pathMid := gengeometry.CenterOfPath(path)
				dat.CurrentPos.X += pathMid.X
				dat.CurrentPos.Y += pathMid.Y

				var res []*Node
				for _, s := range gengeometry.GetPathSides(path) {
					shapeChild := &gengeometry.RectangleShape{
						Width:  0.25,
						Length: 0.25,
					}
					midPoint := s.Start.Add(s.End).Mul(0.5)
					pos := vectors.NewVec3(dat.CurrentPos.X+midPoint.X, dat.CurrentPos.Y+midPoint.Y, dat.CurrentPos.Z)
					midChild := gengeometry.CenterOfPath(shapeChild.GetPath())
					pos.X -= midChild.X
					pos.Y -= midChild.Y
					// dir is the normal vector (e.g. 90° to the side) of the side (pointing outwards).
					dir := vectors.Normalize(s.End.Sub(s.Start)).Perpendicular()
					res = append(res, newNode("Wing", &ShapeData{
						Shape:       shapeChild,
						CurrentPos:  pos,
						CurrentDir:  dir,
						HeightScale: 2 * dat.HeightScale,
						Reorient:    true,
					}))
				}

				currPos2 := vectors.NewVec2(pathMid.X, pathMid.Y)

				for _, p := range path {
					shapeChild := &gengeometry.RectangleShape{
						Width:  0.15,
						Length: 0.15,
					}
					// Direction of the path point pointing towards the center of the parent shape.
					dir := vectors.Normalize(p.Sub(currPos2))
					pos := vectors.NewVec3(dat.CurrentPos.X+p.X, dat.CurrentPos.Y+p.Y, dat.CurrentPos.Z)
					// Translate the position by the child shape mid point.
					// This will cause the child being centered on the path point.
					midChild := gengeometry.CenterOfPath(shapeChild.GetPath())
					pos.X -= midChild.X
					pos.Y -= midChild.Y
					res = append(res, newNode("Wing", &ShapeData{
						Shape:       shapeChild,
						CurrentPos:  pos,
						CurrentDir:  dir,
						HeightScale: 1 * dat.HeightScale,
						Reorient:    true,
					}))
				}
				return res
			},
		},
		"Wing": &NodeRule{
			ID:          "Wing",
			ReplaceNode: false,
			F: func(node *Node) []*Node {
				dat := node.Data.(*ShapeData)
				path := dat.Shape.GetPath()
				pathMid := gengeometry.CenterOfPath(path)
				dat.CurrentPos.X += pathMid.X
				dat.CurrentPos.Y += pathMid.Y

				var res []*Node
				for _, s := range gengeometry.GetPathSides(path) {
					shapeChild := &gengeometry.RectangleShape{
						Width:  0.05,
						Length: 0.05,
					}
					midPoint := s.Start.Add(s.End).Mul(0.5)
					pos := vectors.NewVec3(dat.CurrentPos.X+midPoint.X, dat.CurrentPos.Y+midPoint.Y, dat.CurrentPos.Z)
					midChild := gengeometry.CenterOfPath(shapeChild.GetPath())
					pos.X -= midChild.X
					pos.Y -= midChild.Y
					// dir is the normal vector (e.g. 90° to the side) of the side (pointing outwards).
					dir := vectors.Normalize(s.End.Sub(s.Start)).Perpendicular()
					res = append(res, newNode("Strut", &ShapeData{
						Shape:       shapeChild,
						CurrentPos:  pos,
						CurrentDir:  dir,
						HeightScale: 1 * dat.HeightScale,
						Reorient:    true,
					}))
				}

				currPos2 := vectors.NewVec2(pathMid.X, pathMid.Y)

				for _, p := range path {
					shapeChild := &gengeometry.RectangleShape{
						Width:  0.05,
						Length: 0.05,
					}
					// Direction of the path point pointing towards the center of the parent shape.
					dir := vectors.Normalize(p.Sub(currPos2))
					pos := vectors.NewVec3(dat.CurrentPos.X+p.X, dat.CurrentPos.Y+p.Y, dat.CurrentPos.Z)
					// Translate the position by the child shape mid point.
					// This will cause the child being centered on the path point.
					midChild := gengeometry.CenterOfPath(shapeChild.GetPath())
					pos.X -= midChild.X
					pos.Y -= midChild.Y
					res = append(res, newNode("Strut", &ShapeData{
						Shape:       shapeChild,
						CurrentPos:  pos,
						CurrentDir:  dir,
						HeightScale: 1 * dat.HeightScale,
						Reorient:    true,
					}))
				}
				return res
			},
		},
	}

	// Run the evaluation until we reached terminal nodes.
	for len(stack) > 0 {
		// Pop the last node from the stack.
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check if the node is already replaced.
		if node.Replaced {
			continue
		}

		// Check if the node has a rule.
		rule, ok := rulesByName[node.ID]
		if !ok {
			continue
		}

		// Apply the rule.
		newNodes := rule.Apply(node)
		log.Println("Applied rule", rule.ID, "to node", node.ID, "resulting in", len(newNodes), "new nodes.")

		// Add the new nodes to the stack.
		stack = append(stack, newNodes...)
	}

	return root
}

func ConvertNodeToMesh(node *Node, mesh *gengeometry.Mesh) {
	dat := node.Data.(*ShapeData)
	path := dat.Shape.GetPath()
	if dat.Reorient {
		// Calculate angle from vector.
		curAngle := math.Atan2(dat.CurrentDir.Y, dat.CurrentDir.X)
		center := gengeometry.CenterOfPath(path)
		path = gengeometry.RotatePolygonAroundPoint(path, center, curAngle)
	}
	if !node.Replaced {
		me, err := gengeometry.ExtrudePath(path, dat.HeightScale)
		if err != nil {
			panic(err)
		} else {
			mesh.AddMesh(me, dat.CurrentPos)
		}
	}

	for _, child := range node.Children {
		ConvertNodeToMesh(child, mesh)
	}
}
