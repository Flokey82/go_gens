package genworldvoronoi

import (
	"github.com/fogleman/delaunay"
	"math"
)

type Vertex [2]float64

type TriangleMesh struct {
	RVertex            []Vertex
	RInS               []int
	TVertex            []Vertex
	Triangles          []int
	Halfedges          []int
	numBoundaryRegions int
	numSolidSides      int
	numSides           int
	numRegions         int
	numSolidRegions    int
	numTriangles       int
	numSolidTriangles  int
}

// NewTriangleMesh takes partial mesh information and fills in the rest; the
// partial information is generated in create.js or in fromDelaunator.
func NewTriangleMesh(numBoundaryRegions, numSolidSides int, vxs []Vertex, t []int, h []int) *TriangleMesh {
	m := &TriangleMesh{
		numBoundaryRegions: numBoundaryRegions,
		numSolidSides:      numSolidSides,
		RVertex:            vxs,
		Triangles:          t,
		Halfedges:          h,
	}
	m._update()
	return m
}

// Update internal data structures from Delaunator
func (this *TriangleMesh) update(vxs []Vertex, d *delaunay.Triangulation) {
	this.RVertex = vxs
	this.Triangles = d.Triangles
	this.Halfedges = d.Halfedges
	this._update()
}

func (this *TriangleMesh) s_ghost(s int) bool {
	return s >= this.numSolidSides
}

func s_to_t(s int) int {
	return (s / 3)
}
func s_prev_s(s int) int {
	if s%3 == 0 {
		return s + 2
	}
	return s - 1
}
func s_next_s(s int) int {
	if s%3 == 2 {
		return s - 2
	}
	return s + 1
}

/**
 * Update internal data structures to match the input mesh.
 *
 * Use if you have updated the triangles/halfedges with Delaunator
 * and want the dual mesh to match the updated data. Note that
 * this DOES not update boundary regions or ghost elements.
 */
func (this *TriangleMesh) _update() {
	tri := this.Triangles
	hedges := this.Halfedges
	rvtx := this.RVertex
	tvtx := this.TVertex

	this.numSides = len(tri)
	this.numRegions = len(rvtx)
	this.numSolidRegions = this.numRegions - 1 // TODO: only if there are ghosts
	this.numTriangles = this.numSides / 3
	this.numSolidTriangles = this.numSolidSides / 3

	if len(this.TVertex) < this.numTriangles {
		// Extend this array to be big enough
		numOldTriangles := len(tvtx)
		numNewTriangles := this.numTriangles - numOldTriangles
		tvtx = append(tvtx, make([]Vertex, numNewTriangles)...)
		for t := numOldTriangles; t < this.numTriangles; t++ {
			tvtx[t] = Vertex{0, 0}
		}
		this.TVertex = tvtx
	}

	// Construct an index for finding sides connected to a region
	this.RInS = make([]int, this.numRegions)
	for s := 0; s < len(tri); s++ {
		endpoint := tri[s_next_s(s)]
		if this.RInS[endpoint] == 0 || hedges[s] == -1 {
			this.RInS[endpoint] = s
		}
	}

	// Construct triangle coordinates
	for s := 0; s < len(tri); s += 3 {
		t := s / 3
		a := rvtx[tri[s]]
		b := rvtx[tri[s+1]]
		c := rvtx[tri[s+2]]
		if this.s_ghost(s) {
			// ghost triangle center is just outside the unpaired side
			dx := b[0] - a[0]
			dy := b[1] - a[1]
			scale := 10 / math.Sqrt(dx*dx+dy*dy) // go 10units away from side
			tvtx[t][0] = 0.5*(a[0]+b[0]) + dy*scale
			tvtx[t][1] = 0.5*(a[1]+b[1]) - dx*scale
		} else {
			// solid triangle center is at the centroid
			tvtx[t][0] = (a[0] + b[0] + c[0]) / 3
			tvtx[t][1] = (a[1] + b[1] + c[1]) / 3
		}
	}
	this.TVertex = tvtx
}

func (this *TriangleMesh) r_circulate_r(out_r []int, r int) []int {
	s0 := this.RInS[r]
	incoming := s0
	out_r = out_r[:0]
	for {
		out_r = append(out_r, this.s_begin_r(incoming))
		outgoing := s_next_s(incoming)
		incoming = this.Halfedges[outgoing]
		if incoming == -1 || incoming == s0 {
			break
		}
	}
	return out_r
}

func (this *TriangleMesh) r_circulate_t(out_t []int, r int) []int {
	s0 := this.RInS[r]
	incoming := s0
	out_t = out_t[:0]
	for {
		out_t = append(out_t, s_to_t(incoming))
		outgoing := s_next_s(incoming)
		incoming = this.Halfedges[outgoing]
		if incoming == -1 || incoming == s0 {
			break
		}
	}
	return out_t
}

func (this *TriangleMesh) t_circulate_s(out_s []int, t int) []int {
	out_s = make([]int, 3)
	for i := 0; i < 3; i++ {
		out_s[i] = 3*t + i
	}
	return out_s
}

func (this *TriangleMesh) t_circulate_r(out_r []int, t int) []int {
	out_r = make([]int, 3)
	for i := 0; i < 3; i++ {
		out_r[i] = this.Triangles[3*t+i]
	}
	return out_r
}

func (this *TriangleMesh) r_x(r int) float64 { return this.RVertex[r][0] }
func (this *TriangleMesh) r_y(r int) float64 { return this.RVertex[r][1] }
func (this *TriangleMesh) t_x(r int) float64 { return this.TVertex[r][0] }
func (this *TriangleMesh) t_y(r int) float64 { return this.TVertex[r][1] }
func (this *TriangleMesh) s_end_r(s int) int {
	return this.Triangles[s_next_s(s)]
}

func (this *TriangleMesh) s_begin_r(s int) int {
	return this.Triangles[s]
}

func (this *TriangleMesh) s_opposite_s(s int) int {
	return this.Halfedges[s]
}

func (this *TriangleMesh) s_inner_t(s int) int {
	return s_to_t(s)
}

func (this *TriangleMesh) s_outer_t(s int) int {
	return s_to_t(this.Halfedges[s])
}

func (this *TriangleMesh) ghost_r() int          { return this.numRegions - 1 }
func (this *TriangleMesh) r_ghost(r int) bool    { return r == this.numRegions-1 }
func (this *TriangleMesh) t_ghost(t int) bool    { return this.s_ghost(3 * t) }
func (this *TriangleMesh) s_boundary(s int) bool { return this.s_ghost(s) && (s%3 == 0) }
func (this *TriangleMesh) r_boundary(r int) bool { return r < this.numBoundaryRegions }
