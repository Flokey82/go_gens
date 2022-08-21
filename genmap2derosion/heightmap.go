package genmap2derosion

import (
	"github.com/Flokey82/go_gens/genheightmap"
	"github.com/Flokey82/go_gens/vectors"
)

// Generate initial heightmap.
func (w *World) genTerrain() {
	w.addSlope(vectors.RandomVec2(4))
	w.addVolCone(-1.0)
	// w.addNoise(0.5)
	w.addMountains(50, 0.05)
	for i := 0; i < 10; i++ {
		w.heightRelax()
	}
	w.heightPeaky()
	w.heightNormalize()
}

func (w *World) addSlope(direction vectors.Vec2) {
	w.ApplyGen(genheightmap.GenSlope(direction))
}

func (w *World) addCone(slope float64) {
	w.ApplyGen(genheightmap.GenCone(slope))
}

func (w *World) addVolCone(slope float64) {
	w.ApplyGen(genheightmap.GenVolCone(slope))
}

func (w *World) addMountains(n int, r float64) {
	w.ApplyGen(genheightmap.GenMountains(1, 1, n, r)) // float64(w.dim.X), float64(w.dim.Y)
}

func (w *World) addNoise(amount float64) {
	w.ApplyGen(genheightmap.GenNoise(w.params.Seed, amount))
}

func (w *World) ApplyGen(f genheightmap.GenFunc) {
	x := w.params.Size.X
	y := w.params.Size.Y
	hm := make([]float64, x*y)
	for i := int64(0); i < int64(len(hm)); i++ {
		x := (float64(i/y) / float64(x)) - 0.5
		y := (float64(i%y) / float64(y)) - 0.5
		hm[i] = f(x, y)
	}
	heightNormalize(hm)
	for i, h := range hm {
		w.heightmap[i] += h
	}
	w.heightNormalize() // Normalize
}

func MapF(f genheightmap.Modify, hm []float64) {
	for i := 0; i < len(hm); i++ {
		hm[i] = f(hm[i])
	}
}

func (w *World) getMinMax() (float64, float64) {
	return genheightmap.MinMax(w.heightmap)
}

func (w *World) heightRelax() {
	hm := make([]float64, len(w.heightmap))
	f := genheightmap.ModRelax(w.getNeighbors, func(idx int) float64 {
		return w.heightmap[idx]
	})
	for i, h := range w.heightmap {
		hm[i] = f(i, h)
	}
	heightNormalize(hm)
	for i, h := range hm {
		w.heightmap[i] = h
	}
	w.heightNormalize() // Normalize
}

func (w *World) heightNormalize() {
	heightNormalize(w.heightmap)
}

func heightNormalize(hm []float64) {
	min, max := genheightmap.MinMax(hm)
	MapF(genheightmap.ModNormalize(min, max), hm)
}

func (w *World) heightPeaky() {
	MapF(genheightmap.ModPeaky(), w.heightmap)
	w.heightNormalize() // Normalize
}

func (w *World) getNeighbors(i int) []int {
	return getNeighbors(i, len(w.heightmap), int(w.params.Size.Y))
}

func getNeighbors(i, maxIdx int, dimY int) []int {
	var nbs []int
	if i > 0 {
		nbs = append(nbs, i-1)
		if i >= dimY {
			nbs = append(nbs, i-dimY+1, i-dimY)
		}
		if i-dimY-1 >= 0 {
			nbs = append(nbs, i-dimY-1)
		}
	}
	if i < maxIdx-1 {
		nbs = append(nbs, i+1)
		if i < maxIdx-1-dimY {
			nbs = append(nbs, i+dimY-1, i+dimY)
		}
		if i+dimY+1 < maxIdx {
			nbs = append(nbs, i+dimY+1)
		}
	}
	return nbs
}

func (w *World) surfaceNormal(index int64) vectors.Vec3 {
	var n vectors.Vec3
	dimY := w.params.Size.Y
	dimX := w.params.Size.X
	scale := w.scale
	// NOTE: Would be better to get x,y right in the arguments instead of index
	x := index / dimY
	y := index % dimY
	hIdx := w.heightmap[index]
	enableCross := false

	// Two large triangles adjacent to the plane (+Y -> +X) (-Y -> -X)
	var a, b bool
	if y < dimY-1 {
		if x < dimX-1 {
			//  _____
			// |_|_|_|
			// |_|_|A|
			// |_|B|_|
			n.AddToThis(vectors.Cross3XYZ(0.0, scale*(w.heightmap[index+1]-hIdx), 1.0, 1.0, scale*(w.heightmap[index+dimY]-hIdx), 0.0))
			a = true
		}
		if x > 0 {
			//  _____
			// |_|A|_|
			// |_|_|B|
			// |_|_|_|
			n.AddToThis(vectors.Cross3XYZ(-1.0, scale*(w.heightmap[index-dimY]-hIdx), 0.0, 0.0, scale*(w.heightmap[index+1]-hIdx), 1.0))
			b = true
		}
		if enableCross && a && b {
			//  _____
			// |_|_|A|
			// |_|_|_|
			// |_|_|B|
			n.AddToThis(vectors.Cross3XYZ(-1.0, scale*(w.heightmap[index-dimY+1]-hIdx), 1.0, 1.0, scale*(w.heightmap[index+dimY+1]-hIdx), 1.0))
		}
		a = false
		b = false
	}
	if y > 0 {
		if x > 0 {
			//  _____
			// |_|B|_|
			// |A|_|_|
			// |_|_|_|
			n.AddToThis(vectors.Cross3XYZ(0.0, scale*(w.heightmap[index-1]-hIdx), -1.0, -1.0, scale*(w.heightmap[index-dimY]-hIdx), 0.0))
			a = true
		}
		if x < dimX-1 {
			//  _____
			// |_|_|_|
			// |B|_|_|
			// |_|A|_|
			n.AddToThis(vectors.Cross3XYZ(1.0, scale*(w.heightmap[index+dimY]-hIdx), 0.0, 0.0, scale*(w.heightmap[index-1]-hIdx), -1.0))
			b = true
		}
		if enableCross && a && b {
			//  _____
			// |B|_|_|
			// |_|_|_|
			// |A|_|_|
			n.AddToThis(vectors.Cross3XYZ(1.0, scale*(w.heightmap[index+dimY-1]-hIdx), -1.0, -1.0, scale*(w.heightmap[index-dimY-1]-hIdx), -1.0))
		}
	}
	return n.Normalize()
}
