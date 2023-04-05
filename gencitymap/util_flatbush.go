package gencitymap

import (
	"log"
	"math"

	"github.com/Flokey82/go_gens/utils"
)

// This code is based on https://github.com/mourner/flatbush
type Flatbush struct {
	NumItems    int
	NodeSize    int
	LevelBounds []int
	Boxes       []float64
	Indices     []int
	Pos         int
	MinX        float64
	MinY        float64
	MaxX        float64
	MaxY        float64
}

func NewFlatbush(numItems int, nodeSize int) *Flatbush {
	if numItems <= 0 {
		panic("numItems must be greater than zero")
	}

	if nodeSize < 2 {
		nodeSize = 2
	} else if nodeSize > 65535 {
		nodeSize = 65535
	}

	// calculate the total number of nodes in the R-tree to allocate space for
	// and the index of each tree level (used in search later)
	n := numItems
	numNodes := n
	levelBounds := make([]int, 0)
	levelBounds = append(levelBounds, n*4)
	for {
		n = int(math.Ceil(float64(n) / float64(nodeSize)))
		numNodes += n
		levelBounds = append(levelBounds, numNodes*4)
		if n == 1 {
			break
		}
	}

	boxes := make([]float64, numNodes*4)
	indices := make([]int, numNodes)
	pos := 0
	minX := math.Inf(1)
	minY := math.Inf(1)
	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	return &Flatbush{
		NumItems:    numItems,
		NodeSize:    nodeSize,
		LevelBounds: levelBounds,
		Boxes:       boxes,
		Indices:     indices,
		Pos:         pos,
		MinX:        minX,
		MinY:        minY,
		MaxX:        maxX,
		MaxY:        maxY,
	}
}

func (f *Flatbush) Add(minX float64, minY float64, maxX float64, maxY float64) {
	index := f.Pos >> 2
	f.Indices[index] = index
	f.Boxes[f.Pos] = minX
	f.Pos++
	f.Boxes[f.Pos] = minY
	f.Pos++
	f.Boxes[f.Pos] = maxX
	f.Pos++
	f.Boxes[f.Pos] = maxY
	f.Pos++

	if minX < f.MinX {
		f.MinX = minX
	}
	if minY < f.MinY {
		f.MinY = minY
	}
	if maxX > f.MaxX {
		f.MaxX = maxX
	}
	if maxY > f.MaxY {
		f.MaxY = maxY
	}
}

// / <summary>
// / Method to perform the indexing, to be called after adding all the boxes via <see cref="Add"/>.
// / </summary>
func (f *Flatbush) Finish() {
	if f.Pos != f.NumItems*4 {
		log.Println("f.Pos: ", f.Pos)
		log.Println("f.NumItems: ", f.NumItems)
		log.Println("f.NumItems: ", f.NumItems*4)
		panic("Added incorrect number of items")
	}

	// if number of items is less than node size then skip sorting since each node of boxes must be
	// fully scanned regardless and there is only one node
	if f.NumItems <= f.NodeSize {
		// fill root box with total extents
		f.Boxes[f.Pos] = f.MinX
		f.Pos++
		f.Boxes[f.Pos] = f.MinY
		f.Pos++
		f.Boxes[f.Pos] = f.MaxX
		f.Pos++
		f.Boxes[f.Pos] = f.MaxY
		f.Pos++
		return
	}

	width := f.MaxX - f.MinX
	height := f.MaxY - f.MinY
	hilbertValues := make([]uint, f.NumItems)
	pos := 0

	// map item centers into Hilbert coordinate space and calculate Hilbert values
	for i := 0; i < f.NumItems; i++ {
		pos = 4 * i
		minX := f.Boxes[pos]
		pos++
		minY := f.Boxes[pos]
		pos++
		maxX := f.Boxes[pos]
		pos++
		maxY := f.Boxes[pos]
		pos++

		const n = 1 << 16
		// hilbert max input value for x and y
		const hilbertMax = n - 1
		// mapping the x and y coordinates of the center of the box to values in the range [0 -> n - 1] such that
		// the min of the entire set of bounding boxes maps to 0 and the max of the entire set of bounding boxes maps to n - 1
		// our 2d space is x: [0 -> n-1] and y: [0 -> n-1], our 1d hilbert curve value space is d: [0 -> n^2 - 1]
		x := uint(math.Floor(hilbertMax * ((minX+maxX)/2 - f.MinX) / width))
		y := uint(math.Floor(hilbertMax * ((minY+maxY)/2 - f.MinY) / height))
		hilbertValues[i] = Hilbert(x, y)
	}

	// sort the hilbert values and the indices in parallel
	f.Sort(hilbertValues, f.Boxes, f.Indices, 0, f.NumItems-1, f.NodeSize)

	// generate nodes at each tree level, bottom-up
	pos = 0
	for i := 0; i < len(f.LevelBounds)-1; i++ {
		end := f.LevelBounds[i]

		// generate a parent node for each block of consecutive <nodeSize> nodes
		for pos < end {
			nodeIndex := pos

			// calculate bbox for the new node
			nodeMinX := f.Boxes[pos]
			pos++
			nodeMinY := f.Boxes[pos]
			pos++
			nodeMaxX := f.Boxes[pos]
			pos++
			nodeMaxY := f.Boxes[pos]
			pos++
			for j := 1; j < f.NodeSize && pos < end; j++ {
				nodeMinX = math.Min(nodeMinX, f.Boxes[pos])
				pos++
				nodeMinY = math.Min(nodeMinY, f.Boxes[pos])
				pos++
				nodeMaxX = math.Max(nodeMaxX, f.Boxes[pos])
				pos++
				nodeMaxY = math.Max(nodeMaxY, f.Boxes[pos])
				pos++
			}

			// add the new node to the tree data
			f.Indices[f.Pos>>2] = nodeIndex
			f.Boxes[f.Pos] = nodeMinX
			f.Pos++
			f.Boxes[f.Pos] = nodeMinY
			f.Pos++
			f.Boxes[f.Pos] = nodeMaxX
			f.Pos++
			f.Boxes[f.Pos] = nodeMaxY
			f.Pos++
		}
	}
}

// / <summary>
// / Returns a list of indices to boxes that intersect or overlap the bounding box given, <see cref="Finish"/> must be called before querying.
// / </summary>
// / <param name="minX">Min x value of the bounding box.</param>
// / <param name="minY">Min y value of the bounding box.</param>
// / <param name="maxX">Max x value of the bounding box.</param>
// / <param name="maxY">Max y value of the bounding box.</param>
// / <param name="filter">Optional filter function, if not null then only indices for which the filter function returns true will be included.</param>
// / <returns>List of indices that intersect or overlap with the bounding box given.</returns>
func (f *Flatbush) Query(minX float64, minY float64, maxX float64, maxY float64, filter func(index int) bool) []int {
	if f.Pos != len(f.Boxes) {
		log.Println("f.Pos: ", f.Pos)
		log.Println("f.NumItems: ", f.NumItems*4)

		panic("Data not yet indexed - call Finish()")
	}

	nodeIndex := len(f.Boxes) - 4
	level := len(f.LevelBounds) - 1

	// stack of nodes to search
	stack := make([]int, 0, 64)

	result := make([]int, 0, 64)

	done := false

	for !done {
		// find the end index of the node
		end := utils.Min(nodeIndex+f.NodeSize*4, f.LevelBounds[level])

		// search through the child nodes
		for pos := nodeIndex; pos < end; pos += 4 {
			index := f.Indices[pos>>2]

			// check if the child node intersects with query box
			if maxX < f.Boxes[pos] {
				continue
			}
			if maxY < f.Boxes[pos+1] {
				continue
			}
			if minX > f.Boxes[pos+2] {
				continue
			}
			if minY > f.Boxes[pos+3] {
				continue
			}

			if nodeIndex < f.NumItems*4 {
				if filter == nil || filter(index) {
					result = append(result, index) // leaf item
				}
			} else {
				// push node index and level for further traversal
				stack = append(stack, index)
				stack = append(stack, level-1)
			}
		}

		if len(stack) > 1 {
			level = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nodeIndex = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		} else {
			done = true
		}
	}

	return result
}

// / <summary>
// / Invokes a function on each of the indices of boxes that intersect or overlap with the bounding box given, <see cref="Finish"/> must be called before querying.
// / </summary>
// / <param name="minX">Min x value of the bounding box.</param>
// / <param name="minY">Min y value of the bounding box.</param>
// / <param name="maxX">Max x value of the bounding box.</param>
// / <param name="maxY">Max y value of the bounding box.</param>
// / <param name="visitor">The function to visit each of the result indices, if false is returned no more results will be visited.</param>
func (f *Flatbush) VisitQuery(minX float64, minY float64, maxX float64, maxY float64, visitor func(index int) bool) {
	if f.Pos != f.NumItems*4 {
		panic("Data not yet indexed - call Finish()")
	}

	if visitor == nil {
		panic("Visitor function cannot be nil")
	}

	nodeIndex := len(f.Boxes) - 4
	level := len(f.LevelBounds) - 1

	// stack of nodes to search
	stack := make([]int, 0, 64)

	done := false

	for !done {
		// find the end index of the node
		end := utils.Min(nodeIndex+f.NodeSize*4, f.LevelBounds[level])

		// search through the child nodes
		for pos := nodeIndex; pos < end; pos += 4 {
			index := f.Indices[pos>>2]

			// check if the child node intersects with query box
			if maxX < f.Boxes[pos] {
				continue
			}
			if maxY < f.Boxes[pos+1] {
				continue
			}
			if minX > f.Boxes[pos+2] {
				continue
			}
			if minY > f.Boxes[pos+3] {
				continue
			}

			if nodeIndex < f.NumItems*4 {
				if !visitor(index) {
					return // leaf item
				}
			} else {
				// push node index and level for further traversal
				stack = append(stack, index)
				stack = append(stack, level-1)
			}
		}

		if len(stack) > 1 {
			level = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nodeIndex = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		} else {
			done = true
		}
	}
}

// custom quicksort that partially sorts bbox data alongside the hilbert values
func (f *Flatbush) Sort(values []uint, boxes []float64, indices []int, left int, right int, nodeSize int) {
	// check against nodeSize (only need to sort down to nodeSize buckets)
	if left/nodeSize >= right/nodeSize {
		return
	}

	pivot := values[(left+right)>>1]
	i := left - 1
	j := right + 1

	for {
		i++
		for values[i] < pivot {
			i++
		}

		j--
		for values[j] > pivot {
			j--
		}

		if i >= j {
			break
		}

		f.Swap(values, boxes, indices, i, j)
	}

	f.Sort(values, boxes, indices, left, j, nodeSize)
	f.Sort(values, boxes, indices, j+1, right, nodeSize)
}

// swap two values and two corresponding boxes
func (f *Flatbush) Swap(values []uint, boxes []float64, indices []int, i int, j int) {
	temp := values[i]
	values[i] = values[j]
	values[j] = temp

	k := 4 * i
	m := 4 * j

	a := boxes[k]
	b := boxes[k+1]
	c := boxes[k+2]
	d := boxes[k+3]
	boxes[k] = boxes[m]
	boxes[k+1] = boxes[m+1]
	boxes[k+2] = boxes[m+2]
	boxes[k+3] = boxes[m+3]
	boxes[m] = a
	boxes[m+1] = b
	boxes[m+2] = c
	boxes[m+3] = d

	e := indices[i]
	indices[i] = indices[j]
	indices[j] = e
}

// Fast Hilbert curve algorithm by http://threadlocalmutex.com/
// Ported from C++ https://github.com/rawrunprotected/hilbert_curves (public domain)
func Hilbert(x uint, y uint) uint {
	var a, b, c, d uint

	a = x ^ y
	b = 0xFFFF ^ a
	c = 0xFFFF ^ (x | y)
	d = x & (y ^ 0xFFFF)

	A := a | (b >> 1)
	B := (a >> 1) ^ a
	C := ((c >> 1) ^ (b & (d >> 1))) ^ c
	D := ((a & (c >> 1)) ^ (d >> 1)) ^ d

	a = A
	b = B
	c = C
	d = D
	A = ((a & (a >> 2)) ^ (b & (b >> 2)))
	B = ((a & (b >> 2)) ^ (b & ((a ^ b) >> 2)))
	C ^= ((a & (c >> 2)) ^ (b & (d >> 2)))
	D ^= ((b & (c >> 2)) ^ ((a ^ b) & (d >> 2)))

	a = A
	b = B
	c = C
	d = D
	A = ((a & (a >> 4)) ^ (b & (b >> 4)))
	B = ((a & (b >> 4)) ^ (b & ((a ^ b) >> 4)))
	C ^= ((a & (c >> 4)) ^ (b & (d >> 4)))
	D ^= ((b & (c >> 4)) ^ ((a ^ b) & (d >> 4)))

	a = A
	b = B
	c = C
	d = D
	C ^= ((a & (c >> 8)) ^ (b & (d >> 8)))
	D ^= ((b & (c >> 8)) ^ ((a ^ b) & (d >> 8)))

	a = C ^ (C >> 1)
	b = D ^ (D >> 1)

	i0 := x ^ y
	i1 := b | (0xFFFF ^ (i0 | a))

	i0 = (i0 | (i0 << 8)) & 0x00FF00FF
	i0 = (i0 | (i0 << 4)) & 0x0F0F0F0F
	i0 = (i0 | (i0 << 2)) & 0x33333333
	i0 = (i0 | (i0 << 1)) & 0x55555555

	i1 = (i1 | (i1 << 8)) & 0x00FF00FF
	i1 = (i1 | (i1 << 4)) & 0x0F0F0F0F
	i1 = (i1 | (i1 << 2)) & 0x33333333
	i1 = (i1 | (i1 << 1)) & 0x55555555

	return (i1 << 1) | i0
}
