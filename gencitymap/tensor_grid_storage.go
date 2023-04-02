package gencitymap

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

// GridStorage is a grid-based storage for samples.
// NOTE: Cartesian grid accelerated data structure, Grid of cells, each containing a list of vectors.
type GridStorage struct {
	worldDimensions vectors.Vec2
	gridDimensions  vectors.Vec2
	grid            [][][]vectors.Vec2
	dsepSq          float64
	dsep            float64
	origin          vectors.Vec2
}

// NewGridStorage creates a new grid storage.
// NOTE: worldDimensions assumes origin of 0,0, dsep is the separation distance between samples.
func NewGridStorage(worldDimensions vectors.Vec2, origin vectors.Vec2, dsep float64) *GridStorage {
	gs := &GridStorage{
		dsepSq:          dsep * dsep,
		worldDimensions: worldDimensions,
		gridDimensions:  worldDimensions.Mul(1 / dsep),
		grid:            make([][][]vectors.Vec2, int(worldDimensions.X/dsep)),
		dsep:            dsep,
		origin:          origin,
	}
	for x := range gs.grid {
		gs.grid[x] = make([][]vectors.Vec2, int(worldDimensions.Y/dsep))
		for y := range gs.grid[x] {
			gs.grid[x][y] = make([]vectors.Vec2, 0)
		}
	}
	return gs
}

// AddAll adds all samples from another grid to this one.
func (gs *GridStorage) AddAll(gridStorage *GridStorage) {
	for x := range gridStorage.grid {
		for y := range gridStorage.grid[x] {
			for _, sample := range gridStorage.grid[x][y] {
				gs.AddSample(sample, nil)
			}
		}
	}
}

func (gs *GridStorage) AddPolyline(line []vectors.Vec2) {
	for _, v := range line {
		gs.AddSample(v, nil)
	}
}

// AddSample adds a sample to the grid.
// NOTE: Does not enforce separation, does not clone.
func (gs *GridStorage) AddSample(v vectors.Vec2, coords *vectors.Vec2) {
	if coords == nil {
		coord := gs.GetSampleCoords(v)
		coords = &coord
	}
	gs.grid[int(coords.X)][int(coords.Y)] = append(gs.grid[int(coords.X)][int(coords.Y)], v)
}

// IsValidSample returns true if the sample is valid and within the world dimensions.
// Tests whether v is at least d away from samples.
// NOTE: Performance is very important - this is called at every integration step!
// dSq=this.dsepSq squared test distance
// Could be dtest if we are integrating a streamline
func (gs *GridStorage) IsValidSample(v vectors.Vec2, dSq float64) bool {
	// Code duplication with this.getNearbyPoints but much slower when calling
	// this.getNearbyPoints due to array creation in that method

	coords := gs.GetSampleCoords(v)

	// Check samples in 9 cells in 3x3 grid
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			cell := coords.Add(vectors.Vec2{X: float64(x), Y: float64(y)})
			if !gs.VectorOutOfBounds(cell, gs.gridDimensions) {
				if !gs.vectorFarFromVectors(v, gs.grid[int(cell.X)][int(cell.Y)], dSq) {
					return false
				}
			}
		}
	}

	return true
}

// vectorFarFromVectors returns true if the vector is far from the vectors.
// Test whether v is at least d away from vectors.
// NOTE: Performance very important - this is called at every integration step
// dSq squared test distance
func (gs *GridStorage) vectorFarFromVectors(v vectors.Vec2, vectors []vectors.Vec2, dSq float64) bool {
	for _, sample := range vectors {
		if sample != v {
			distanceSq := sample.DistanceToSquared(v)
			if distanceSq < dSq {
				return false
			}
		}
	}

	return true
}

// GetNearbyPoints returns points in cells surrounding v.
// Results include v, if it exists in the grid.
// NOTE: Returns samples (kind of) closer than distance - returns all samples in
// cells so approximation (square to approximate circle).
func (gs *GridStorage) GetNearbyPoints(v vectors.Vec2, distance float64) []vectors.Vec2 {
	radius := int(math.Ceil((distance / gs.dsep) - 0.5))
	coords := gs.GetSampleCoords(v)
	out := make([]vectors.Vec2, 0)
	for x := -1 * radius; x <= 1*radius; x++ {
		for y := -1 * radius; y <= 1*radius; y++ {
			cell := coords.Add(vectors.Vec2{X: float64(x), Y: float64(y)})
			if !gs.VectorOutOfBounds(cell, gs.gridDimensions) {
				for _, v2 := range gs.grid[int(cell.X)][int(cell.Y)] {
					out = append(out, v2)
				}
			}
		}
	}

	return out
}

func (gs *GridStorage) worldToGrid(v vectors.Vec2) vectors.Vec2 {
	return v.Sub(gs.origin)
}

func (gs *GridStorage) gridToWorld(v vectors.Vec2) vectors.Vec2 {
	return v.Add(gs.origin)
}

func (gs *GridStorage) VectorOutOfBounds(v vectors.Vec2, bounds vectors.Vec2) bool {
	return (v.X < 0 || v.Y < 0 ||
		v.X >= bounds.X || v.Y >= bounds.Y)
}

// GetSampleCoords returns the cell coords corresponding to the vector.
// NOTE: Performance is important here - this is called at every integration step.
func (gs *GridStorage) GetSampleCoords(worldV vectors.Vec2) vectors.Vec2 {
	v := gs.worldToGrid(worldV)
	if gs.VectorOutOfBounds(v, gs.worldDimensions) {
		// log.error("Tried to access out-of-bounds sample in grid");
		return vectors.Vec2{}
	}

	return vectors.Vec2{
		X: math.Floor(v.X / gs.dsep),
		Y: math.Floor(v.Y / gs.dsep),
	}
}
