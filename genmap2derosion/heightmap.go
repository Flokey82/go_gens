package genmap2derosion

import (
	"bufio"
	"fmt"
	"github.com/Flokey82/go_gens/vectors"
	"log"
	"math"
	"os"

	opensimplex "github.com/ojrac/opensimplex-go"
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
	// Initialize Heightmap
	perlin := opensimplex.New(w.seed)

	// Mountainy:
	// perlin.SetOctaveCount(8)
	// perlin.SetFrequency(1.0)
	// perlin.SetPersistence(0.5)
	// fmb := &fractalBrownianMotionOptions{8, 1.0, 0.5}//{octaves, lacunarity, gain}

	var e float64
	mult := 15.0
	pow := 1.0

	hm := make([]float64, w.dim.X*w.dim.Y)
	for i := int64(0); i < w.dim.X*w.dim.Y; i++ {
		x := (float64(i) / float64(w.dim.Y)) * (float64(1.0) / float64(w.dim.X)) * mult
		y := float64(int64(i)%int64(w.dim.Y)) * (float64(1.0) / float64(w.dim.Y)) * mult
		e = 1 * perlin.Eval2(x, y)
		e += 0.5 * perlin.Eval2(x*2, y*2)
		e += 0.25 * perlin.Eval2(x*4, y*4)
		e /= (1 + 0.5 + 0.25)
		hm[i] += math.Pow(e, pow)
	}
	normalizeHeight(hm)
	for i, h := range hm {
		w.heightmap[i] += h * amount
	}

	// Normalize
	w.normalizeHeight()
}

func (w *World) addMountains(n int, r float64) {
	var mounts [][2]float64
	for i := 0; i < n; i++ {
		mounts = append(mounts, [2]float64{float64(w.dim.X) * w.r.Float64(), float64(w.dim.Y) * w.r.Float64()})
	}
	hm := make([]float64, w.dim.X*w.dim.Y)

	for i := int64(0); i < w.dim.X*w.dim.Y; i++ {
		x := float64(i / w.dim.Y)
		y := float64(i % w.dim.Y)
		for j := 0; j < n; j++ {
			m := mounts[j]
			hm[i] += math.Pow(math.Exp(-((x-m[0])*(x-m[0])+(y-m[1])*(y-m[1]))/(2*r*r)), 2)
		}
	}
	normalizeHeight(hm)
	for i, h := range hm {
		w.heightmap[i] += h
	}

	// Normalize
	w.normalizeHeight()
}

func (w *World) addCone(slope float64) {
	hm := make([]float64, w.dim.X*w.dim.Y)
	for i := int64(0); i < w.dim.X*w.dim.Y; i++ {
		x := float64(i/w.dim.Y) - float64(w.dim.X/2)
		y := float64(i%w.dim.Y) - float64(w.dim.Y/2)
		hm[i] += math.Pow(x*x+y*y, 0.5) * slope
	}
	normalizeHeight(hm)
	for i, h := range hm {
		w.heightmap[i] += h
	}

	// Normalize
	w.normalizeHeight()
}

func (w *World) addSlope(direction vectors.Vec2) {
	hm := make([]float64, w.dim.X*w.dim.Y)
	for i := int64(0); i < w.dim.X*w.dim.Y; i++ {
		x := float64(i / w.dim.Y)
		y := float64(i % w.dim.Y)
		hm[i] += x*direction.X + y*direction.Y
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
	return getMinMax(w.heightmap[:])
}

func normalizeHeight(hm []float64) {
	min, max := getMinMax(hm)
	for i := 0; i < len(hm); i++ {
		hm[i] = (hm[i] - min) / (max - min)
	}
}

func (w *World) normalizeHeight() {
	normalizeHeight(w.heightmap[:])
}

func (w *World) peakyHeight() {
	peakyHeight(w.heightmap[:])

	// Normalize
	w.normalizeHeight()
}
func peakyHeight(hm []float64) {
	for i := 0; i < len(hm); i++ {
		hm[i] = math.Sqrt(hm[i])
	}
}
func (w *World) relaxHeight() {
	hm := relaxHeight(w.heightmap[:], int(w.dim.Y))
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
		return CalcMean(vals)
	}
	for i := 0; i < len(hm); i++ {
		nh[i] = getMeanNeighbor(i)
	}
	return nh
}

func CalcMean(nums []float64) float64 {
	total := 0.0
	for _, v := range nums {
		total += v
	}
	return total / float64(len(nums))
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
