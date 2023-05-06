package genmap2derosion

import (
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
)

func (w *World) erodeWind(cycles int) {
	// Track the Movement of all Particles
	track := make([]bool, w.params.Size.X*w.params.Size.Y)
	sizeVec2 := vectors.Vec2{X: float64(w.params.Size.X), Y: float64(w.params.Size.Y)}

	// Do a series of iterations!
	for i := 0; i < cycles; i++ {

		// Spawn New Particle on Boundary
		var newpos vectors.Vec2
		shift := rand.Int() % (int(w.params.Size.X) + int(w.params.Size.Y))
		if shift < int(w.params.Size.X) {
			newpos = vectors.Vec2{X: float64(shift), Y: 1}
		} else {
			newpos = vectors.Vec2{X: 1, Y: float64(shift - int(w.params.Size.X))}
		}

		wind := NewWind(newpos)
		wind.Fly(w.heightmap, w.windpath, w.sediment, track, sizeVec2, w.scale)

		// If we should store GIF frames, do so.
		if w.params.StoreGIFFrames && i%100 == 0 {
			w.storeGifFrame(w.sediment, w.sediment, w.sediment)
		}
	}

	// Update Path
	lrate := 0.01
	for i := 0; i < int(w.params.Size.X*w.params.Size.Y); i++ {
		if track[i] {
			w.windpath[i] = (1.0-lrate)*w.windpath[i] + lrate*1.0
		} else {
			w.windpath[i] = (1.0-lrate)*w.windpath[i] + lrate*0.0
		}
	}
}

// Parameters
const (
	dtwind     = 0.25  // Delta time for wind erosion
	suspension = 0.002 // Affects transport rate
	abrasion   = 0.01  // Abrasion rate
	roughness  = 0.005
	settling   = 0.05
	gravity    = 0.09
)

type Wind struct {
	Index    int
	Pos      vectors.Vec2
	Height   float64
	PSpeed   vectors.Vec3
	Speed    vectors.Vec3
	Sediment float64 //Sediment Mass
}

func NewWind(pos vectors.Vec2) *Wind {
	pspeed := vectors.Vec3{X: 1.0, Y: 0.0, Z: 1.0}.Normalize().Mul(2.0)
	return &Wind{
		Pos:      pos,
		PSpeed:   pspeed,
		Speed:    pspeed,
		Sediment: 0.0,
	}
}

func (w *Wind) Fly(h, windpath, sediment []float64, track []bool, dim vectors.Vec2, scale float64) {
	ipos := w.Pos
	for {
		// Initial Position
		ipos = w.Pos
		ind := int(ipos.X)*int(dim.Y) + int(ipos.Y)

		// Set Height Correctly
		if w.Height < h[ind]+sediment[ind] {
			w.Height = h[ind] + sediment[ind]
		}

		// Surface Normal (Using Heightmap + Sediment Map)
		n := surfaceNormal3(w.Pos, h, sediment, dim, scale)

		// Movement Mechanics
		if w.Height > h[ind]+sediment[ind] { // Flying
			w.Speed.Y -= dtwind * gravity // Gravity
		} else { // Contact Movement
			track[ind] = true
			w.Speed = w.Speed.Add(n.Cross(n.Cross(w.Speed))).Mul(dtwind * 0.1).Add(w.Speed)
		}

		w.Speed = w.Speed.Add(w.PSpeed.Sub(w.Speed).Mul(dtwind * 0.1))
		w.Pos = w.Pos.Add(vectors.Vec2{X: w.Speed.X, Y: w.Speed.Z}.Mul(dtwind))
		w.Height += dtwind * w.Speed.Y

		// New Position
		nind := int(w.Pos.X)*int(dim.Y) + int(w.Pos.Y)

		// Out-Of-Bounds
		if w.Pos.X < 0 || w.Pos.Y < 0 || w.Pos.X >= dim.X || w.Pos.Y >= dim.Y {
			break
		}

		// Mass Transport

		// Surface Contact
		if w.Height <= h[nind]+sediment[nind] {
			force := w.Speed.Len() * (sediment[nind] + h[nind] - w.Height)
			if sediment[ind] <= 0 { // Abrasion
				sediment[ind] = 0
				h[ind] -= dtwind * abrasion * force * w.Sediment
				sediment[ind] += dtwind * abrasion * force * w.Sediment
			} else if sediment[ind] > dtwind*suspension*force { // Suspension
				sediment[ind] -= dtwind * suspension * force
				w.Sediment += dtwind * suspension * force
				w.Cascade(ind, h, sediment, dim)
			} else {
				sediment[ind] = 0 // Set to zero
			}
		} else { // Flying Particle
			w.Sediment -= dtwind * suspension * w.Sediment

			sediment[nind] += 0.5 * dtwind * suspension * w.Sediment
			sediment[ind] += 0.5 * dtwind * suspension * w.Sediment

			w.Cascade(nind, h, sediment, dim)
			w.Cascade(ind, h, sediment, dim)
		}

		// Particle has no speed (equilibrium movement)
		if w.Speed.Len() < 0.01 {
			break
		}
	}
}

func (w *Wind) Cascade(i int, h, sediment []float64, dim vectors.Vec2) {
	// Neighbor Positions (8-Way)
	nx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	ny := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	n := []int{i - int(dim.Y) - 1, i - int(dim.Y), i - int(dim.Y) + 1, i - 1, i + 1,
		i + int(dim.Y) - 1, i + int(dim.Y), i + int(dim.Y) + 1}

	// Iterate over all Neighbors
	for m := 0; m < 8; m++ {
		if n[m] < 0 || n[m] >= len(h) {
			continue
		}
		if int(w.Pos.X)+nx[m] >= int(dim.X) || int(w.Pos.Y)+ny[m] >= int(dim.Y) {
			continue
		}
		if int(w.Pos.X)+nx[m] < 0 || int(w.Pos.Y)+ny[m] < 0 {
			continue
		}

		// Pile Size Difference
		diff := (h[i] + sediment[i]) - (h[n[m]] + sediment[n[m]])
		excess := math.Abs(diff) - roughness
		if excess <= 0 {
			continue
		}

		// Transfer Mass
		var transfer float64
		if diff > 0 { // Pile is Larger
			transfer = math.Min(sediment[i], excess/2.0)
		} else { // Neighbor is Larger
			transfer = -math.Min(sediment[n[m]], excess/2.0)
		}

		sediment[i] -= dtwind * settling * transfer
		sediment[n[m]] += dtwind * settling * transfer
	}
}

func surfaceNormal(index int, h, sediment []float64, dim vectors.Vec2, scale float64) vectors.Vec3 {
	n := vectors.Vec3{X: 0.0, Y: 0.0, Z: 0.0}
	dimY := int(dim.Y)
	dimX := int(dim.X)
	x := index / int(dimY)
	y := index % int(dimY)

	if x < 1 || x >= dimX-1 || y < 1 || y >= dimY-1 {
		return n
	}

	// Two large triangles adjacent to the plane (+Y -> +X) (-Y -> -X)
	for i := 1; i <= 1; i++ {
		n = n.Add(vectors.Vec3{X: 0.0, Y: scale * (h[index+i] - h[index] + sediment[index+i] - sediment[index]), Z: float64(i)}.Cross(vectors.Vec3{X: float64(i), Y: scale * (h[index+i*int(dim.Y)] - h[index] + sediment[index+i*int(dim.Y)] - sediment[index]), Z: 0.0})).Mul(1.0 / float64(i*i))
		n = n.Add(vectors.Vec3{X: 0.0, Y: scale * (h[index-i] - h[index] + sediment[index-i] - sediment[index]), Z: float64(-i)}.Cross(vectors.Vec3{X: float64(-i), Y: scale * (h[index-i*int(dim.Y)] - h[index] + sediment[index-i*int(dim.Y)] - sediment[index]), Z: 0.0})).Mul(1.0 / float64(i*i))
		n = n.Add(vectors.Vec3{X: float64(i), Y: scale * (h[index+i*int(dim.Y)] - h[index] + sediment[index+i*int(dim.Y)] - sediment[index]), Z: 0.0}.Cross(vectors.Vec3{X: 0.0, Y: scale * (h[index-i] - h[index] + sediment[index-i] - sediment[index]), Z: float64(-i)})).Mul(1.0 / float64(i*i))
		n = n.Add(vectors.Vec3{X: float64(-i), Y: scale * (h[index-i*int(dim.Y)] - h[index] + sediment[index-i*int(dim.Y)] - sediment[index]), Z: 0.0}.Cross(vectors.Vec3{X: 0.0, Y: scale * (h[index+i] - h[index] + sediment[index+i] - sediment[index]), Z: float64(i)})).Mul(1.0 / float64(i*i))
	}

	return n.Normalize()
}

func surfaceNormal3(pos vectors.Vec2, h, sediment []float64, dim vectors.Vec2, scale float64) vectors.Vec3 {
	P00 := floorVec2(pos) // Floored Position

	P10 := P00.Add(vectors.Vec2{X: 1.0, Y: 0.0})
	P01 := P00.Add(vectors.Vec2{X: 0.0, Y: 1.0})
	P11 := P00.Add(vectors.Vec2{X: 1.0, Y: 1.0})

	N00 := surfaceNormal(int(P00.X)*int(dim.Y)+int(P00.Y), h, sediment, dim, scale)
	N10 := surfaceNormal(int(P10.X)*int(dim.Y)+int(P10.Y), h, sediment, dim, scale)
	N01 := surfaceNormal(int(P01.X)*int(dim.Y)+int(P01.Y), h, sediment, dim, scale)
	N11 := surfaceNormal(int(P11.X)*int(dim.Y)+int(P11.Y), h, sediment, dim, scale)

	// Weights (modulo position)
	w := vectors.Vec2{X: 1.0, Y: 1.0}.Sub(modVec2(pos, vectors.Vec2{X: 1.0, Y: 1.0}))
	return N00.Mul(w.X * w.Y).Add(N10.Mul(1.0 - w.X).Mul(w.Y)).Add(N01.Mul(w.X).Mul(1.0 - w.Y)).Add(N11.Mul(1.0 - w.X).Mul(1.0 - w.Y))
}

func floorVec3(v vectors.Vec3) vectors.Vec3 {
	return vectors.Vec3{X: math.Floor(v.X), Y: math.Floor(v.Y), Z: math.Floor(v.Z)}
}

func floorVec2(v vectors.Vec2) vectors.Vec2 {
	return vectors.Vec2{X: math.Floor(v.X), Y: math.Floor(v.Y)}
}

func modVec2(v, vmod vectors.Vec2) vectors.Vec2 {
	return vectors.Vec2{X: math.Mod(v.X, vmod.X), Y: math.Mod(v.Y, vmod.Y)}
}
