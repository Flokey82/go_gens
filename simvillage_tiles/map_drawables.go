package simvillage_tiles

// drawable is a collection of tiles that can be drawn on a chunk.
// TODO: Use named layers instead of a fixed number of expected layers.
//
// NOTE:
// - A tile value of -1 sets the tile type of the dst layer to 0
// - A tile value of 0 does not affect the tile type of the dst layer
type drawable struct {
	Dimensions
	ground        []int
	groundOverlay []int
	structures    []int
	triggers      []Trigger // This should be done a bit better with multiple triggers and maybe data for what to do when triggered.
}

// Ground returns the Ground tiles as *Layer.
func (d *drawable) Ground() *Layer {
	return &Layer{d.Dimensions, d.ground}
}

// GroundOverlay returns the GroundOverlay tiles as *Layer.
func (d *drawable) GroundOverlay() *Layer {
	return &Layer{d.Dimensions, d.groundOverlay}
}

// Structures returns the Structures tiles as *Layer.
func (d *drawable) Structures() *Layer {
	return &Layer{d.Dimensions, d.structures}
}

// Trigger represents an interaction trigger.
type Trigger struct {
	Position    [2]int // Position of the trigger in the chunk.
	Destination        // Teleport the user to this destination.
}

// Destination is a destination for a trigger.
// TODO: Add also destination map.
type Destination struct {
	X, Y  int
	Chunk [2]int
}

var house1 = drawable{
	Dimensions: NewDimensions(4, 4),
	ground: []int{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	},
	groundOverlay: []int{
		-1, -1, -1, -1,
		-1, -1, -1, -1,
		-1, -1, -1, -1,
		-1, -1, -1, -1,
	},
	structures: []int{
		58, 59, 60, 61,
		83, 84, 85, 86,
		108, 109, 110, 111,
		133, 134, 135, 136,
	},
	triggers: []Trigger{
		{Position: [2]int{1, 4}, Destination: Destination{X: 2, Y: 4, Chunk: [2]int{10, 0}}},
		{Position: [2]int{2, 4}, Destination: Destination{X: 2, Y: 4, Chunk: [2]int{10, 0}}},
	},
}

var house2 = drawable{
	Dimensions: NewDimensions(6, 4),
	ground: []int{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	},
	groundOverlay: []int{
		-1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
	},
	structures: []int{
		63, 64, 65, 66, 67, 68,
		88, 89, 90, 91, 92, 93,
		113, 114, 115, 116, 117, 118,
		138, 139, 140, 141, 142, 143,
	},
	triggers: []Trigger{
		{Position: [2]int{2, 4}, Destination: Destination{X: 3, Y: 5, Chunk: [2]int{20, 0}}},
		{Position: [2]int{3, 4}, Destination: Destination{X: 3, Y: 5, Chunk: [2]int{20, 0}}},
	},
}

var house3 = drawable{
	Dimensions: NewDimensions(6, 5),
	ground: []int{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	},
	groundOverlay: []int{
		-1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
	},
	structures: []int{
		26, 27, 28, 29, 30, 31,
		51, 52, 53, 54, 55, 56,
		76, 77, 78, 79, 80, 81,
		101, 102, 103, 104, 105, 106,
		126, 127, 128, 129, 130, 131,
	},
	triggers: []Trigger{
		{Position: [2]int{2, 5}, Destination: Destination{X: 3, Y: 5, Chunk: [2]int{5, 0}}},
		{Position: [2]int{3, 5}, Destination: Destination{X: 3, Y: 5, Chunk: [2]int{5, 0}}},
	},
}

var hedge = drawable{
	Dimensions: NewDimensions(9, 5),
	ground: []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 202, 203, 204, 205, 206, 207, 208, 0,
		0, 227, 228, 229, 230, 231, 232, 233, 0,
		0, 0, 0, 0, 255, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	},
	groundOverlay: []int{
		-1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, 280, -1, -1, -1, -1,
	},
	structures: []int{
		176, 177, 178, 179, 180, 181, 182, 183, 184,
		201, 0, 0, 0, 0, 0, 0, 0, 209,
		226, 0, 0, 0, 0, 0, 0, 0, 234,
		251, 252, 253, 254, 0, 256, 257, 258, 259,
		276, 277, 278, 279, 0, 281, 282, 283, 284,
	},
}
