package genmap2derosion

import (
	"bufio"
	"fmt"
	"github.com/Flokey82/go_gens/genheightmap"
	"github.com/Flokey82/go_gens/vectors"
	"math"
	"os"
)

func (w *World) ExportOBJ(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	wr := bufio.NewWriter(f)
	defer wr.Flush()

	// Write the vertex index.
	for i, h := range w.heightmap {
		x := float64(i/int(w.dim.Y)) / float64(w.dim.Y)
		y := float64(i%int(w.dim.Y)) / float64(w.dim.Y)
		wr.WriteString(fmt.Sprintf("v %f %f %f \n", x, h, y))
	}

	// Write the triangles.
	for x := 0; x < int(w.dim.X)-1; x++ {
		for y := 0; y < int(w.dim.Y-1); y++ {
			i1 := x*int(w.dim.Y) + y
			i2 := i1 + 1
			i3 := i2 + int(w.dim.Y)
			i4 := i1 + int(w.dim.Y)
			wr.WriteString(fmt.Sprintf("f %d %d %d \n", i1+1, i2+1, i3+1))
			wr.WriteString(fmt.Sprintf("f %d %d %d \n", i4+1, i1+1, i3+1))
		}
	}
	return nil
}

// Generate initial heightmap.
func (w *World) generate() {
	w.addCone(1.0)
	w.addNoise(0.5)
	w.addMountains(20, 40.0)
	w.addSlope(vectors.RandomVec2(1.0))
	// w.peakyHeight()
	// w.relaxHeight()
}

func (w *World) addNoise(amount float64) {
	w.ApplyGen(genheightmap.GenNoise(w.seed, amount))
}

func (w *World) addMountains(n int, r float64) {
	w.ApplyGen(genheightmap.GenMountains(float64(w.dim.X), float64(w.dim.Y), n, r))
}

func (w *World) addCone(slope float64) {
	w.ApplyGen(genheightmap.GenCone(slope))
}

func (w *World) addSlope(direction vectors.Vec2) {
	w.ApplyGen(genheightmap.GenSlope([2]float64{direction.X, direction.Y}))
}

func (w *World) ApplyGen(f genheightmap.GenFunc) {
	hm := make([]float64, w.dim.X*w.dim.Y)
	for i := int64(0); i < w.dim.X*w.dim.Y; i++ {
		x := float64(i / w.dim.Y)
		y := float64(i % w.dim.Y)
		hm[i] = f(x, y)
	}
	normalizeHeight(hm)
	for i, h := range hm {
		w.heightmap[i] += h
	}

	// Normalize
	w.normalizeHeight()
}

func getMinMax(hm []float64) (float64, float64) {
	var min, max float64
	for _, h := range hm {
		if h > max {
			max = h
		}

		if h < min {
			min = h
		}
	}
	return min, max
}

func (w *World) getMinMax() (float64, float64) {
	return getMinMax(w.heightmap)
}

func normalizeHeight(hm []float64) {
	min, max := getMinMax(hm)
	for i := 0; i < len(hm); i++ {
		hm[i] = (hm[i] - min) / (max - min)
	}
}

func (w *World) normalizeHeight() {
	normalizeHeight(w.heightmap)
}

func (w *World) peakyHeight() {
	peakyHeight(w.heightmap)

	// Normalize
	w.normalizeHeight()
}
func peakyHeight(hm []float64) {
	for i := 0; i < len(hm); i++ {
		hm[i] = math.Sqrt(hm[i])
	}
}
func (w *World) relaxHeight() {
	hm := relaxHeight(w.heightmap, int(w.dim.Y))
	normalizeHeight(hm)
	for i, h := range hm {
		w.heightmap[i] = h
	}

	// Normalize
	w.normalizeHeight()
}
func relaxHeight(hm []float64, dimY int) []float64 {
	nh := make([]float64, len(hm))
	getMeanNeighbor := func(i int) float64 {
		var vals []float64
		for _, nb := range getNeighbors(i, hm, dimY) {
			vals = append(vals, hm[nb])
		}
		return genheightmap.CalcMean(vals)
	}
	for i := 0; i < len(hm); i++ {
		nh[i] = getMeanNeighbor(i)
	}
	return nh
}

func getNeighbors(i int, hm []float64, dimY int) []int {
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
	if i < len(hm)-1 {
		nbs = append(nbs, i+1)
		if i < len(hm)-1-dimY {
			nbs = append(nbs, i+dimY-1, i+dimY)
		}
		if i+dimY+1 < len(hm) {
			nbs = append(nbs, i+dimY+1)
		}
	}
	return nbs
}

func (w *World) surfaceNormal(index int64) vectors.Vec3 {
	var n vectors.Vec3
	dimY := w.dim.Y
	dimX := w.dim.X
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
