package genarchitecture

import (
	"log"
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
)

type Rule struct {
	Name         string
	Shape        int
	Roof         bool
	Width        float64
	Length       float64
	RulesSides   []string
	RulesCorners []string
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
	}
	return nil
}

var SampleRules = []*Rule{
	{
		Name:         "base",
		Shape:        ShapeTypePlus,
		Width:        2,
		Length:       2,
		Roof:         true,
		RulesSides:   []string{"wing"},
		RulesCorners: []string{"corner"},
	},
	{
		Name:         "wing",
		Shape:        ShapeTypeRect,
		Width:        0.35,
		Length:       0.35,
		Roof:         true,
		RulesSides:   []string{"strut"},
		RulesCorners: []string{"strut"},
	},
	{
		Name:         "corner",
		Shape:        ShapeTypeRect,
		Width:        0.1,
		Length:       0.1,
		Roof:         true,
		RulesSides:   []string{"strut"},
		RulesCorners: []string{"strut"},
	},
	{
		Name:         "strut",
		Shape:        ShapeTypeRect,
		Width:        0.025,
		Length:       0.025,
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

		path := rule.GetShape().GetPath()

		entry.CurrentPos.X -= rule.Width / 2
		entry.CurrentPos.Y -= rule.Length / 2

		// Extrude the path.
		extrudedMesh, err := gengeometry.ExtrudePath(path, 0.2)
		if err != nil {
			log.Println(err)
		} else {
			mesh.AddMesh(extrudedMesh, entry.CurrentPos)
		}

		if rule.Roof {
			// Add a roof to the path.
			roofMesh, err := gengeometry.TaperPath(path, 0.05)
			if err != nil {
				log.Println(err)
			} else {
				mesh.AddMesh(roofMesh, vectors.NewVec3(entry.CurrentPos.X, entry.CurrentPos.Y, 0.2))
			}
		}
		// If there is a rule for the sides, add them to the stack.
		if len(rule.RulesSides) > 0 {
			sideRule := rule.RulesSides[rand.Intn(len(rule.RulesSides))]
			for _, s := range gengeometry.GetPathSides(path) {
				midPoint := s.Start.Add(s.End).Mul(0.5)
				// Add the side to the stack.
				stack = append(stack, stackEntry{
					CurrentPos:  vectors.NewVec3(entry.CurrentPos.X+midPoint.X, entry.CurrentPos.Y+midPoint.Y, entry.CurrentPos.Z),
					RuleToApply: sideRule,
				})
			}
		}

		// If there is a rule for the corners, add them to the stack.
		if len(rule.RulesCorners) > 0 {
			cornerRule := rule.RulesCorners[rand.Intn(len(rule.RulesCorners))]
			for _, c := range path {
				// Add the corner to the stack.
				stack = append(stack, stackEntry{
					CurrentPos:  vectors.NewVec3(entry.CurrentPos.X+c.X, entry.CurrentPos.Y+c.Y, entry.CurrentPos.Z),
					RuleToApply: cornerRule,
				})
			}
		}
	}

	return mesh
}
