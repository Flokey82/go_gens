package genlsystem

import (
	"bufio"
	"fmt"
	"github.com/Flokey82/go_gens/vectors"
	"image/color"
	"math"
	"os"
)

type Bounds3d struct {
	minX, minY, minZ float64
	maxX, maxY, maxZ float64
}

func (a *Bounds3d) AddPoint(x, y, z float64) {
	if a.minX > x {
		a.minX = x
	}
	if a.maxX < x {
		a.maxX = x
	}
	if a.minY > y {
		a.minY = y
	}
	if a.maxY < y {
		a.maxY = y
	}
	if a.minZ > z {
		a.minZ = z
	}
	if a.maxZ < z {
		a.maxZ = z
	}
}

func (a *Bounds3d) Size() (width, height, depth float64) {
	return a.maxX - a.minX, a.maxY - a.minY, a.maxZ - a.minZ
}

type stack3d struct {
	x, y, z float64
	color   color.Color
	width   float64
	forward vectors.Vec3
	up      vectors.Vec3
	prev    *stack3d
}

type line3d struct {
	x1, y1, x2, y2, z1, z2 float64
	color                  color.Color
	width                  float64
}

// Turtle3d implements a turtle-esque drawing system IN 3D!
// This code is partially inspired by:
// https://github.com/yalue/l_system_3d
// and:
// https://github.com/recp/cglm/blob/master/include/cglm/vec3.h
type Turtle3d struct {
	rules    map[string]func(*Turtle3d)
	cur      *stack3d
	lines    []line3d
	boundary *Bounds3d
}

// NewTurtle3d returns a new Turtle3d struct.
func NewTurtle3d(rules map[string]func(*Turtle3d)) *Turtle3d {
	return &Turtle3d{
		rules: rules,
		cur: &stack3d{
			color:   color.RGBA{0x00, 0x00, 0x00, 0xFF},
			width:   1,
			forward: vectors.NewVec3(1, 0, 0),
			up:      vectors.NewVec3(0, 1, 0),
		},
		lines:    []line3d{},
		boundary: &Bounds3d{},
	}
}

// save position and angle
func (t *Turtle3d) Save() {
	t.cur = &stack3d{
		x:       t.cur.x,
		y:       t.cur.y,
		z:       t.cur.z,
		color:   t.cur.color,
		width:   t.cur.width,
		prev:    t.cur, // stackception
		up:      t.cur.up,
		forward: t.cur.forward,
	}
}

// restore position and angle
func (t *Turtle3d) Restore() {
	if t.cur == nil {
		return
	}

	if p := t.cur.prev; p != nil {
		t.cur.prev = nil // stackalypse
		t.cur = p
	}
}

// Move moves forward without drawing.
func (t *Turtle3d) Move(f float64) {
	x, y, z := t.advanceRotateEtc(f)
	t.cur.x += x
	t.cur.y += y
	t.cur.z += z
	t.boundary.AddPoint(t.cur.x, t.cur.y, t.cur.z)
}

func (t *Turtle3d) advanceRotateEtc(distance float64) (x, y, z float64) {
	change := t.cur.forward.Mul(distance)
	x = change.X
	y = change.Y
	z = change.Z
	return
}

// Draw draws a line forward.
func (t *Turtle3d) Draw(f float64) {
	sx, sy, sz := t.cur.x, t.cur.y, t.cur.z
	x, y, z := t.advanceRotateEtc(f)
	t.cur.x += x
	t.cur.y += y
	t.cur.z += z

	t.boundary.AddPoint(t.cur.x, t.cur.y, t.cur.z)

	t.lines = append(t.lines, line3d{
		x1:    sx,
		y1:    sy,
		z1:    sz,
		x2:    t.cur.x,
		y2:    t.cur.y,
		z2:    t.cur.z,
		color: t.cur.color,
		width: t.cur.width,
	})
}

// GetPosition gets the current position.
func (t *Turtle3d) GetPosition() (x, y, z float64) {
	return t.cur.x, t.cur.y, t.cur.z
}

func degToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

// Rotate is the same as Yaw.
func (t *Turtle3d) Rotate(angle float64) {
	glm_vec3_rotate(&t.cur.forward, degToRadians(angle), t.cur.up)
}

func (t *Turtle3d) Pitch(angle float64) {
	right := vectors.Cross3(t.cur.forward, t.cur.up)
	glm_vec3_rotate(&t.cur.up, degToRadians(angle), right)
	glm_vec3_rotate(&t.cur.forward, degToRadians(angle), right)
}

func (t *Turtle3d) Roll(angle float64) {
	glm_vec3_rotate(&t.cur.up, degToRadians(angle), t.cur.forward)
}

// GetForward gets the current forward vector
func (t *Turtle3d) GetForward() vectors.Vec3 {
	return t.cur.forward
}

// SetColor sets the current color.
func (t *Turtle3d) SetColor(c color.Color) {
	t.cur.color = c
}

// GetColor gets the current color.
func (t *Turtle3d) GetColor() color.Color {
	return t.cur.color
}

// SetWidth sets the current width.
func (t *Turtle3d) SetWidth(w float64) {
	t.cur.width = w
}

// GetWidth gets the current width.
func (t *Turtle3d) GetWidth() float64 {
	return t.cur.width
}

// UNLEASH THE TURTLE!
func (t *Turtle3d) Go(fname string, path []string) error {
	// draw lines
	for _, c := range path {
		if f, ok := t.rules[c]; ok {
			f(t)
		}
	}
	border := 5.0
	offx, offy, offz := t.boundary.minX-border, t.boundary.minY-border, t.boundary.minZ-border

	// Generate index for vertex indices.
	var vxs [][3]float64
	vtxIdx := make(map[[3]float64]int)
	for _, line := range t.lines {
		start := [3]float64{line.x1 - offx, line.y1 - offy, line.z1 - offz}
		if _, ok := vtxIdx[start]; !ok {
			vtxIdx[start] = len(vxs)
			vxs = append(vxs, start)
		}
		stop := [3]float64{line.x2 - offx, line.y2 - offy, line.z2 - offz}
		if _, ok := vtxIdx[stop]; !ok {
			vtxIdx[stop] = len(vxs)
			vxs = append(vxs, stop)
		}
	}
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	wr := bufio.NewWriter(f)
	for _, p := range vxs {
		wr.WriteString(fmt.Sprintf("v %f %f %f \n", p[0], p[1], p[2]))
	}
	for _, line := range t.lines {
		start := [3]float64{line.x1 - offx, line.y1 - offy, line.z1 - offz}
		stop := [3]float64{line.x2 - offx, line.y2 - offy, line.z2 - offz}
		wr.WriteString(fmt.Sprintf("l %d %d \n", vtxIdx[start]+1, vtxIdx[stop]+1))
	}
	wr.Flush()

	// cleanup
	t.Cleanup()
	return nil
}

func (t *Turtle3d) Cleanup() {
	t.cur = &stack3d{
		color:   color.RGBA{0x00, 0x00, 0x00, 0xFF},
		width:   1,
		forward: vectors.NewVec3(1, 0, 0),
		up:      vectors.NewVec3(0, 1, 0),
	}
	t.lines = nil
	t.boundary = new(Bounds3d)
}

func glm_vec3_rotate(v *vectors.Vec3, angle float64, axis vectors.Vec3) {
	c := math.Cos(angle)
	s := math.Sin(angle)
	k := axis.Normalize()

	/* Right Hand, Rodrigues' rotation formula:
	   v = v*cos(t) + (kxv)sin(t) + k*(k.v)(1 - cos(t))
	*/
	v1 := v.Mul(c)
	v2 := vectors.Cross3(k, *v)
	v2 = v2.Mul(s)
	v1 = vectors.Add3(v1, v2)
	v2 = k.Mul(vectors.Dot3(k, *v) * (1.0 - c))
	*v = vectors.Add3(v1, v2)
}
