package vmesh

import (
	"github.com/fogleman/delaunay"
	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
	"math"
	"math/rand"
	"sort"
)

type Extent struct {
	Width  float64
	Height float64
}

func (e *Extent) BBox() voronoi.BBox {
	w := e.Width / 2
	h := e.Height / 2
	return voronoi.NewBBox(-w, w, -h, h)
}

var defaultExtent *Extent = &Extent{
	Width:  1.0,
	Height: 1.0,
}

type Edge struct {
	IdxA  int
	IdxB  int
	Left  *voronoi.Cell
	Right *voronoi.Cell
}

type Mesh struct {
	Points      []voronoi.Vertex        // Points used to generate the voronoi diagram
	Voronoi     *voronoi.Diagram        // Voronoi diagram
	Vertices    []voronoi.Vertex        // Vertex points sourrounding a cell
	AdjacentVxs map[int][]int           // Vertex index to adjacent vertices
	VertexTris  map[int][]*voronoi.Cell // Vertex index to bordering cells (triangles)
	Edges       []Edge                  // Edges in voronoi diagram
	Extent      *Extent
}

func (m *Mesh) Triangulate() (*delaunay.Triangulation, error) {
	var pts []delaunay.Point
	for _, p := range m.Vertices {
		pts = append(pts, delaunay.Point{X: p.X, Y: p.Y})
	}
	return delaunay.Triangulate(pts)
}

func (m *Mesh) Distance(i, j int) float64 {
	p := m.Vertices[i]
	q := m.Vertices[j]
	return math.Sqrt((p.X-q.X)*(p.X-q.X) + (p.Y-q.Y)*(p.Y-q.Y))
}

func (m *Mesh) IsEdge(i int) bool {
	return len(m.AdjacentVxs[i]) < 3
}

func (m *Mesh) IsNearEdge(i int) bool {
	x := m.Vertices[i].X
	y := m.Vertices[i].Y
	w := m.Extent.Width
	h := m.Extent.Height
	return x < -0.49*w || x > 0.49*w || y < -0.49*h || y > 0.49*h
}

func (m *Mesh) Neighbours(vxIdx int) []int {
	return m.AdjacentVxs[vxIdx]
}

func (m *Mesh) MapF(f func(v voronoi.Vertex) float64) *Heightmap {
	h := NewHeightmap(m)
	for i := range m.Vertices {
		h.Values[i] = f(m.Vertices[i])
	}
	return h
}

func GenerateGoodMesh(n int, extent *Extent) *Mesh {
	if extent == nil {
		extent = defaultExtent
	}
	return MakeMesh(generateGoodPoints(n, extent), extent)
}

func MakeMesh(pts []voronoi.Vertex, extent *Extent) *Mesh {
	if extent == nil {
		extent = defaultExtent
	}
	vor := voronoiGen(pts, extent)
	var vxs []voronoi.Vertex
	var edges []Edge
	vxids := make(map[voronoi.Vertex]int)
	adj := make(map[int][]int)
	tris := make(map[int][]*voronoi.Cell)
	for i := 0; i < len(vor.Edges); i++ {
		e := vor.Edges[i]
		if e == nil {
			continue
		}
		e0, e0ok := vxids[e.Va.Vertex]
		e1, e1ok := vxids[e.Vb.Vertex]
		if !e0ok {
			e0 = len(vxs)
			vxids[e.Va.Vertex] = e0
			vxs = append(vxs, e.Va.Vertex)
		}
		if !e1ok {
			e1 = len(vxs)
			vxids[e.Vb.Vertex] = e1
			vxs = append(vxs, e.Vb.Vertex)
		}
		adj[e0] = append(adj[e0], e1)
		adj[e1] = append(adj[e1], e0)
		edges = append(edges, Edge{
			IdxA:  e0,
			IdxB:  e1,
			Left:  e.LeftCell,
			Right: e.RightCell,
		})
		if !isInList(tris[e0], e.LeftCell) {
			tris[e0] = append(tris[e0], e.LeftCell)
		}
		if e.RightCell != nil && !isInList(tris[e0], e.RightCell) {
			tris[e0] = append(tris[e0], e.RightCell)
		}
		if !isInList(tris[e1], e.LeftCell) {
			tris[e1] = append(tris[e1], e.LeftCell)
		}
		if e.RightCell != nil && !isInList(tris[e1], e.RightCell) {
			tris[e1] = append(tris[e1], e.RightCell)
		}
	}

	return &Mesh{
		Points:      pts,
		Voronoi:     vor,
		Vertices:    vxs,
		AdjacentVxs: adj,
		VertexTris:  tris,
		Edges:       edges,
		Extent:      extent,
	}
}

func generateGoodPoints(n int, extent *Extent) []voronoi.Vertex {
	if extent == nil {
		extent = defaultExtent
	}

	bbox := extent.BBox()
	pts := utils.RandomSites(bbox, n)
	sort.Slice(pts, func(a, b int) bool {
		return (pts[a].X - pts[b].X) > 0
	})

	// Compute voronoi diagram.
	d := voronoi.ComputeDiagram(pts, bbox, true)

	// Relax using Lloyd's algorithm
	for i := 0; i < 16; i++ {
		pts = utils.LloydRelaxation(d.Cells)
		d = voronoi.ComputeDiagram(pts, bbox, true)
	}
	return pts
}

func generateCellPoints(c *voronoi.Cell, n int) []voronoi.Vertex {
	bbox := getCellBB(c)
	sites := make([]voronoi.Vertex, n)
	w := bbox.Xr - bbox.Xl
	h := bbox.Yb - bbox.Yt
	for j := 0; j < n; j++ {
		for {
			site := voronoi.Vertex{
				X: rand.Float64()*w + bbox.Xl,
				Y: rand.Float64()*h + bbox.Yt,
			}
			if utils.InsideCell(c, site) {
				sites[j] = site
				break
			}
		}
	}
	return sites
}

func getCellBB(c *voronoi.Cell) voronoi.BBox {
	var minX, minY, maxX, maxY float64
	setMinMax := func(v voronoi.Vertex) {
		if v.X < minX {
			minX = v.X
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}
	for _, edge := range c.Halfedges {
		setMinMax(edge.Edge.Va.Vertex)
		setMinMax(edge.Edge.Vb.Vertex)
	}
	return voronoi.NewBBox(minX, maxX, minY, maxY)
}

func voronoiGen(pts []voronoi.Vertex, extent *Extent) *voronoi.Diagram {
	if extent == nil {
		extent = defaultExtent
	}
	return voronoi.ComputeDiagram(pts, extent.BBox(), true)
}

func isInList(l []*voronoi.Cell, c *voronoi.Cell) bool {
	for _, ci := range l {
		if c == ci {
			return true
		}
	}
	return false
}
