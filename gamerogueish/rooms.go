package gamerogueish

import (
	"math/rand"
)

// CarveRoom sets all tiles occupied by the room to ' '.
func (w *World) CarveRoom(room *Room) {
	enableRoundedCorners := false
	// randomize rounded corners
	// TODO: This should be part of the room config.
	if enableRoundedCorners {
		for i := range room.RoundedCorner {
			room.RoundedCorner[i] = rand.Intn(2) == 0
		}
	}
	for y := room.Y; y < room.Y+room.H; y++ {
		for x := room.X; x < room.X+room.W; x++ {
			// Make sure that if we have rounded corners, we don't overwrite them.
			if room.RoundedCorner[DirNorth] && y == room.Y && (x == room.X || x == room.X+room.W-1) ||
				room.RoundedCorner[DirEast] && x == room.X+room.W-1 && (y == room.Y || y == room.Y+room.H-1) ||
				room.RoundedCorner[DirSouth] && y == room.Y+room.H-1 && (x == room.X || x == room.X+room.W-1) ||
				room.RoundedCorner[DirWest] && x == room.X && (y == room.Y || y == room.Y+room.H-1) {
				continue
			}
			w.Cells[y][x] = CharFloor
			w.Elevation[y][x] = room.E
		}
	}
}

// AddRoomFurnishings adds furnishings to the given room.
func (w *World) AddRoomFurnishings(room *Room) {
	// Check if we have a room config for the given room.
	// In this case we just have a bedroom.
	if RoomConfigBedroom.Condition(room, w) {
		RoomConfigBedroom.Apply(room, w)
	}
	if RoomConfigAltar.Condition(room, w) {
		RoomConfigAltar.Apply(room, w)
	}
	if RoomConfigLibrary.Condition(room, w) {
		RoomConfigLibrary.Apply(room, w)
	}
}

// AddRoomPuddle adds a puddle of water to the given room.
func (w *World) AddRoomPuddle(room *Room) {
	// Place a small puddle of water in the room.
	// We pick a random location which is not blocking any entrance to the room.
	// Then we random walk a random number of steps for each poddle cell.
	var puddleCells [][2]int

	// Random position in the room that is not blocking the entrance.
	var px, py int
	for i := 0; i < 10; i++ {
		px = rand.Intn(room.W-2) + room.X + 1
		py = rand.Intn(room.H-2) + room.Y + 1
		if !room.NextToDoor(px, py) {
			break
		}
	}

	puddleCells = append(puddleCells, [2]int{px, py})
	// Random number of steps (max 10% of the room size)
	numSteps := rand.Intn(room.W*room.H/10) + 1
	for i := 0; i < numSteps; i++ {
		// Pick a random cell.
		idx := rand.Intn(len(puddleCells))
		px, py = puddleCells[idx][0], puddleCells[idx][1]
		// Pick a random direction.
		dir := rand.Intn(4)
		switch dir {
		case DirNorth:
			py--
		case DirEast:
			px++
		case DirSouth:
			py++
		case DirWest:
			px--
		}
		// Check if the position is valid.
		if w.InBounds(px, py) && w.IsEmpty(px, py) && !room.NextToDoor(px, py) {
			puddleCells = append(puddleCells, [2]int{px, py})
		}
	}

	// Place the puddle cells.
	for _, cell := range puddleCells {
		w.Cells[cell[1]][cell[0]] = CharWater
	}
}

// AddRoomFountain adds a fountain to the given room.
func (w *World) AddRoomFountain(room *Room) {
	// TODO: Randomize fountain shape.
	fountainSideLenX := room.W - 2
	fountainSideLenY := room.H - 2
	if fountainSideLenX < 3 || fountainSideLenY < 3 {
		return
	}
	minX := room.X + (room.W-fountainSideLenX)/2
	minY := room.Y + (room.H-fountainSideLenY)/2

	// TODO: Make this configurable.
	hasWalls := true
	hasColumns := false
	hasArcades := true
	hasBrickWalls := false

	// isCorner returns true if the given position is a corner.
	isCorner := func(dx, dy int) bool {
		return dx == 0 && dy == 0 || dx == 0 && dy == fountainSideLenY-1 || dx == fountainSideLenX-1 && dy == 0 || dx == fountainSideLenX-1 && dy == fountainSideLenY-1
	}

	// calculate the interval for evenly spaced columns.
	ivalX := calcInterval(fountainSideLenX)
	ivalY := calcInterval(fountainSideLenY)

	// isArcade returns true if the given position is an arcade column.
	isArcade := func(dx, dy int) bool {
		return dx%ivalX == 0 && dy%ivalY == 0
	}

	// isOuterWall returns true if the given position is an outer wall.
	isOuterWall := func(dx, dy int) bool {
		return dx == 0 || dy == 0 || dx == fountainSideLenX-1 || dy == fountainSideLenY-1
	}

	for dx := 0; dx < fountainSideLenX; dx++ {
		for dy := 0; dy < fountainSideLenY; dy++ {
			x := minX + dx
			y := minY + dy
			if hasWalls && isOuterWall(dx, dy) {
				if hasColumns && isCorner(dx, dy) || hasArcades && isArcade(dx, dy) {
					w.Cells[y][x] = CharColumn
				} else if hasBrickWalls { // Brick wall (can't be crossed)
					w.Cells[y][x] = CharWall
				} else { // Raised floor (can be crossed)
					w.Elevation[y][x] += 1
				}
			} else {
				w.Cells[y][x] = CharWater
			}
		}
	}
}

// calcInterval calculates the interval for evenly spaced columns.
func calcInterval(l int) int {
	if l < 3 {
		return 1
	}
	// We're looking for the smallest interval that is greater than 1 and divides the length evenly.
	// Essentially we're looking for prime numbers (numbers that are only divisible by 1 and themselves).
	ival := 2
	for l%ival != 1 {
		ival++
	}
	return ival
}

// AddRoomColumns adds columns to the given room.
func (w *World) AddRoomColumns(room *Room) {
	// Outer arcades
	//
	// ###########
	// #o o o o o#
	// #         #
	// #o       o#
	// #         #
	// #o       o#
	// #         #
	// #o o o o o#
	// ###########
	if rand.Intn(2) == 0 {
		// Calculate ideal interval for evenly spaced columns given the room dimensions.
		ivalX := calcInterval(room.W)
		ivalY := calcInterval(room.H)
		for dx := 0; dx < room.W; dx++ {
			for dy := 0; dy < room.H; dy++ {
				if dx == 0 || dy == 0 || dx == room.W-1 || dy == room.H-1 {
					if dx%ivalX == 0 && dy%ivalY == 0 && !room.NextToDoor(room.X+dx, room.Y+dy) {
						w.Cells[room.Y+dy][room.X+dx] = CharColumn
					}
				}
			}
		}
		return
	}

	// Inner arcades
	//
	// ###########
	// #         #
	// # o o o o #
	// #         #
	// # o o o o #
	// #         #
	// # o o o o #
	// #         #
	// ###########

	// We evenly space the columns in the room, leaving at least one cell between each column
	// and the room walls.
	minX := room.X + 1
	maxX := room.X + room.W - 1
	minY := room.Y + 1
	maxY := room.Y + room.H - 1

	// Determine the x and y interval between each column.
	xInterval := calcInterval(maxX - minX)
	yInterval := calcInterval(maxY - minY)
	if xInterval == 1 || yInterval == 1 {
		// No space for columns.
		return
	}

	// Place the columns.
	for y := minY; y < maxY; y += yInterval {
		for x := minX; x < maxX; x += xInterval {
			w.Cells[y][x] = CharColumn
		}
	}
}

// Room represents a room in the world.
// TODO: Store connecting rooms
type Room struct {
	X, Y          int     // top left corner
	W, H          int     // width and height
	E             int     // elevation
	RoundedCorner [4]bool // rounded corners
	Connections   []*Connection
}

// Connect connects the given room to the given room.
func (r *Room) Connect(to *Room, x, y int) {
	r.Connections = append(r.Connections, &Connection{
		X:  x,
		Y:  y,
		To: to,
	})
}

// RandAlongWall returns a random position along the wall of the room.
func (r *Room) RandAlongWall() (x, y int) {
	// Pick a random wall.
	for i := 0; i < 100; i++ {
		/*
			if rand.Intn(2) == 0 {
				// Vertical
				x = rand.Intn(2)*(room.W-1) + room.X
				y = rand.Intn(room.H) + room.Y
			} else {
				// Horizontal
				x = rand.Intn(room.W) + room.X
				y = rand.Intn(2)*(room.H-1) + room.Y
			}
		*/
		switch rand.Intn(4) {
		case DirNorth:
			x = rand.Intn(r.W) + r.X
			y = r.Y
		case DirEast:
			x = r.X + r.W - 1
			y = rand.Intn(r.H) + r.Y
		case DirSouth:
			x = rand.Intn(r.W) + r.X
			y = r.Y + r.H - 1
		case DirWest:
			x = r.X
			y = rand.Intn(r.H) + r.Y
		}
		// Check if we are not in front of a door.
		if r.Contains(x, y) && !r.NextToDoor(x, y) {
			break
		}
	}
	return x, y
}

// Size returns the size of the room in square tiles.
func (r *Room) Size() int {
	return r.W * r.H
}

// NextToDoor returns true if the given position is next to a door.
func (r *Room) NextToDoor(x, y int) bool {
	for _, conn := range r.Connections {
		for _, n := range tileNeighboursCardinal {
			if conn.X+n[0] == x && conn.Y+n[1] == y {
				return true
			}
		}
	}
	return false
}

// Contains returns true if the given position is inside the room.
func (r *Room) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

// Overlaps returns true if the given room overlaps with any of the rooms in the list.
func (r *Room) Overlaps(rooms []*Room) bool {
	for _, room := range rooms {
		if r.X+r.W < room.X || r.X > room.X+room.W || r.Y+r.H < room.Y || r.Y > room.Y+room.H {
			continue
		}
		return true
	}
	return false
}

type Connection struct {
	X, Y int // position of the connection
	To   *Room
}
