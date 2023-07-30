package gamerogueish

import (
	"math/rand"
)

// World represents a game world.
type World struct {
	Cells     [][]rune    // 2D array of world cells
	Elevation [][]int     // 2D array of elevation values
	Objects   [][]*Object // 2D array of objects like furniture, chests, doors (TODO: Move this to cells)
	Width     int         // width of the world in cells
	Height    int         // height of the world in cells
	Entities  []*Entity   // entities in the world (creatures)
	Items     []*Item     // items in the world
}

// NewWorld returns a new world with the given width and height.
func NewWorld(width, height int) *World {
	w := &World{
		Width:  width,
		Height: height,
	}
	w.Cells = make([][]rune, height)
	w.Elevation = make([][]int, height)
	w.Objects = make([][]*Object, height)
	for y := range w.Cells {
		w.Cells[y] = make([]rune, width)
		w.Elevation[y] = make([]int, width)
		w.Objects[y] = make([]*Object, width)
	}
	return w
}

// NextToWall returns true if the given position is next to a wall.
func (w *World) NextToWall(x, y int) bool {
	return w.IsSolid(x-1, y) || w.IsSolid(x+1, y) || w.IsSolid(x, y-1) || w.IsSolid(x, y+1)
}

// IsEmtpy checks if a tile is empty (tile content is a space ' ' character and there is no object there).
func (w *World) IsEmpty(x int, y int) bool {
	return w.Cells[y][x] == CharFloor && w.Objects[y][x] == nil
}

// IsSolid checks if a tile is solid (tile content is not a space ' ' character).
func (w *World) IsSolid(x int, y int) bool {
	return w.Cells[y][x] != CharFloor
}

// CanMoveTo checks if a tile is solid and if it is not occupied by an entity.
func (w *World) CanMoveTo(x, y int) bool {
	// TODO: Bounds check.
	return w.Cells[y][x] == CharFloor && w.Objects[y][x] == nil
}

// Fill all cells with the given tile.
func (w *World) Fill(c rune) {
	for y := range w.Cells {
		for x := range w.Cells[y] {
			w.Cells[y][x] = c
		}
	}
}

// InBounds returns true if the given position is within the world bounds.
func (w *World) InBounds(x, y int) bool {
	return x >= 0 && x < w.Width && y >= 0 && y < w.Height
}

// Cardinal directions.
const (
	DirNorth = 0
	DirEast  = 1
	DirSouth = 2
	DirWest  = 3
)

// GenWorldSimpleDungeon generates a simple random-walk-ish dungeon.
// - A starting room is placed in the center of the world.
// - Rooms are then placed in random directions neighboring a randomly selectd room.
// - Rooms are not placed if they would overlap with an existing room.
func GenWorldSimpleDungeon(width, height int, seed int64) *World {
	const (
		attempts    = 200
		maxRooms    = 200
		minRoomSize = 2
		maxRoomSize = 20
	)
	w := NewWorld(width, height)
	w.Fill(CharWall)

	ssrc := rand.NewSource(seed)
	rng := rand.New(ssrc)

	// Start with a single room in the middle of the map.
	rl := randInt(rng, minRoomSize, maxRoomSize)
	rw := randInt(rng, minRoomSize, maxRoomSize)
	rooms := []*Room{{
		X: (width / 2) - rw/2,
		Y: (height / 2) - rl/2,
		W: rw,
		H: rl,
	}}

	// Carve out the starting room.
	w.CarveRoom(rooms[0])

	doesOverlap := func(newRoom *Room) bool {
		return newRoom.X < 1 || newRoom.Y < 1 || newRoom.X+newRoom.W > width-1 || newRoom.Y+newRoom.H > height-1 || newRoom.Overlaps(rooms)
	}

	placeRoom := func(room, newRoom *Room, dir int) {
		// Append the room to the dungeon.
		rooms = append(rooms, newRoom)

		// Draw room.
		w.CarveRoom(newRoom)

		// Draw a tunnel between the rooms and connect them.
		// TODO: Allow alignment of the corridor to be specified (e.g. center, left, right).
		var corridorX, corridorY int
		switch dir {
		case DirNorth:
			corridorX = newRoom.X + newRoom.W/2
			corridorY = newRoom.Y + newRoom.H
		case DirEast:
			corridorX = newRoom.X - 1
			corridorY = newRoom.Y + newRoom.H/2
		case DirSouth:
			corridorX = newRoom.X + newRoom.W/2
			corridorY = newRoom.Y - 1
		case DirWest:
			corridorX = newRoom.X + newRoom.W
			corridorY = newRoom.Y + newRoom.H/2
		}

		newRoom.Connect(room, corridorX, corridorY)
		room.Connect(newRoom, corridorX, corridorY)

		// Set the cell to a floor tile and set the elevation to the average of the two rooms.
		w.Cells[corridorY][corridorX] = CharFloor
		w.Elevation[corridorY][corridorX] = (room.E + newRoom.E) / 2
	}

	// Place rooms until we run out of attempts or reach the max room count.
	for i := 0; i < attempts; i++ {
		// Pick a room and place a neighboring room.
		room := rooms[rng.Intn(len(rooms))]

		// Pick a random direction in which to place the room (north, east, south, west).
		dir := rng.Intn(4)

		// Pick a random length and width.
		rl := randInt(rng, minRoomSize, maxRoomSize)
		rw := randInt(rng, minRoomSize, maxRoomSize)

		// Pick a randome elevation based on the room elevation (up to 2 higher or lower).
		re := room.E + rng.Intn(5) - 3

		// Calculate position based on direction.
		// NOTE: Right now we center the neighboring room.
		//
		// This could be changed to use a random offset to
		// make the rooms more varied.
		var x, y int
		var misaligned bool
		if randInt(rng, 0, 100) < 50 {
			// Randomly we do not center the room, but place it randomly along the edge.
			//
			// ###############
			// #####         #
			// #####         #
			// #####         #
			// ####### #######
			// #             #
			// #             #
			// #             #
			// ###############
			//
			// TODO: Prevent edge misalignment on split rooms.
			//
			// ##############
			// #####  #  #  #
			// #####  #  #  #
			// #####  #  #  #
			// ###### ## ## # <- misaligned
			// #        #####
			// #        #####
			// #        #####
			// ##############
			switch dir {
			case DirNorth:
				x = room.X + randInt(rng, 0, room.W) - rw/2
				y = room.Y - rl - 1
			case DirEast:
				x = room.X + room.W + 1
				y = room.Y + randInt(rng, 0, room.H) - rl/2
			case DirSouth:
				x = room.X + randInt(rng, 0, room.W) - rw/2
				y = room.Y + room.H + 1
			case DirWest:
				x = room.X - rw - 1
				y = room.Y + randInt(rng, 0, room.H) - rl/2
			}

			// Remark that the room is misaligned.
			misaligned = true
		} else {
			// This will center the new room next to the room.
			//
			// ###############
			// ###         ###
			// ###         ###
			// ###         ###
			// ####### #######
			// #             #
			// #             #
			// #             #
			// ###############
			switch dir {
			case DirNorth:
				x = room.X + room.W/2 - rw/2
				y = room.Y - rl - 1
			case DirEast:
				x = room.X + room.W + 1
				y = room.Y + room.H/2 - rl/2
			case DirSouth:
				x = room.X + room.W/2 - rw/2
				y = room.Y + room.H + 1
			case DirWest:
				x = room.X - rw - 1
				y = room.Y + room.H/2 - rl/2
			}
		}

		// Create a new room with the calculated position and size.
		newRoom := &Room{
			X: x,
			Y: y,
			W: rw,
			H: rl,
			E: re,
		}

		// Check if the new room overlaps with any existing rooms.
		if doesOverlap(newRoom) {
			continue
		}

		// TODO: Take note of child-parent relationships which will allow us to
		// determine which rooms are connected to each other and what type of
		// room they are (e.g. kitchen, bedroom, corridor, etc.).

		// Place the room or randomly create smaller rooms within the bounds of the new room.
		if rng.Intn(100) < 90 || misaligned {
			placeRoom(room, newRoom, dir)
		} else {
			// MULTI ROOM SPLIT
			//
			// Create smaller rooms within the bounds of the new room.
			// Depending on the orientation, we either create a horizontal or vertical split.

			// Pick the shortest side shared between old and new room that we can split,
			// so we are guaranteed to connect the rooms.
			var splitLength int
			if dir == DirNorth || dir == DirSouth {
				splitLength = newRoom.W
				if splitLength > room.W {
					splitLength = room.W
				}
			} else {
				splitLength = newRoom.H
				if splitLength > room.H {
					splitLength = room.H
				}
			}

			// Calculate the maximum number of rooms we can place.
			maxNumRooms := (splitLength + 1) / (minRoomSize + 1)
			roomSize := minRoomSize

			// If we can place more than 2 rooms, we can decrease the number of rooms
			// and increase the room size.
			if maxNumRooms > 3 {
				// Calculate the room size based on the number of rooms we want to place.
				maxNumRooms = randInt(rng, 3, maxNumRooms)
				roomSize = ((splitLength + 1) / maxNumRooms) - 1
			}

			// Calculate the starting position of the rooms.
			startX := newRoom.X
			startY := newRoom.Y
			if dir == DirNorth || dir == DirSouth {
				startX += (newRoom.W - splitLength - 1) / 2 // Floor division.
			} else {
				startY += (newRoom.H - splitLength - 1) / 2 // Floor division.
			}

			// Place the rooms.
			for j := 0; j < maxNumRooms; j++ {
				// Use minRoomSize as the room size depending on the orientation.
				var roomW, roomH int
				var roomX, roomY int
				if dir == DirNorth || dir == DirSouth {
					roomH = newRoom.H
					roomW = roomSize
					roomX = startX + j*(roomSize+1)
					roomY = startY
				} else {
					roomW = newRoom.W
					roomH = roomSize
					roomX = startX
					roomY = startY + j*(roomSize+1)
				}

				// Create a new room.
				newRoom := &Room{
					X: roomX,
					Y: roomY,
					W: roomW,
					H: roomH,
					E: re,
				}

				// Check if the new room overlaps with any existing rooms.
				if doesOverlap(newRoom) {
					continue
				}

				// Place the room.
				placeRoom(room, newRoom, dir)
			}
		}

		// Stop if we have enough rooms.
		if len(rooms) > maxRooms {
			break
		}
	}

	// Pick the last room as the exit.
	// TODO: Improve this and make sure that this is reachable.
	exit := ItemTypeExit.New()
	exit.X = rooms[len(rooms)-1].X + rooms[len(rooms)-1].W/2
	exit.Y = rooms[len(rooms)-1].Y + rooms[len(rooms)-1].H/2
	w.Items = append(w.Items, exit)

	// Place furnishings, columns and puddles in rooms.
	// TODO: Avoid placing items in the same position.
	for i, room := range rooms {
		// Small rooms have furnishings (which is currently a bedroom)
		if room.Size() <= 20 {
			w.AddRoomFurnishings(room)
		}
		// Larger rooms may have puddles and columns.
		if room.W >= 4 && room.H >= 4 {
			if rng.Intn(100) < 20 {
				// There is a chance that a puddle of water is placed randomly
				// in the room.
				w.AddRoomPuddle(room)
			} else if rng.Intn(100) < 5 {
				// There is a chance that columns are placed randomly
				// in the room.
				w.AddRoomColumns(room)
			} else if rng.Intn(100) < 10 && i != 0 && i != len(rooms)-1 {
				// There is a chance that a fountain is placed randomly
				// in the room.
				w.AddRoomFountain(room)
			}
		}

		// Place a chest randomly in the room.
		if rng.Intn(100) < 20 {
			chest := ItemTypeChest.New()
			chest.X, chest.Y = room.RandAlongWall()
			w.Items = append(w.Items, chest)
		}

		// Place a bookshelf randomly in the room.
		if rng.Intn(100) < 20 {
			book := ItemTypeBookshelf.New()
			book.X, book.Y = room.RandAlongWall()
			w.Items = append(w.Items, book)
		}

		// Place a trap randomly in the room.
		if rng.Intn(100) < 20 {
			trap := ItemTypeTrap.New()
			trap.X = randInt(rng, room.X, room.X+room.W)
			trap.Y = randInt(rng, room.Y, room.Y+room.H)
			w.Items = append(w.Items, trap)
		}

		// There is a chance that a creature entity is placed randomly
		// in the room.
		if rng.Intn(100) < 20 {
			// Pick a random location in the room.
			cx := randInt(rng, room.X, room.X+room.W)
			cy := randInt(rng, room.Y, room.Y+room.H)
			w.Entities = append(w.Entities, NewEntity(cx, cy, MonsterEntities[rng.Intn(len(MonsterEntities))]))
		}
	}

	// Place some random stairs in rooms.
	//
	// TODO: Instead, assign different room heights and add stairs to the
	// connecting rooms accounting for 50% of the height difference in
	// each room.
	//
	// Example:
	// 		  Room 1
	//        ___
	//     __/
	// ___/
	// Room 2
	//
	return w
}

// randInt returns a random integer between min and max using the given rng.
func randInt(rng *rand.Rand, min, max int) int {
	return min + rng.Intn(max-min)
}

// GenWorldBigBox generates a big box world.
func GenWorldBigBox(width, height int, seed int64) *World {
	w := NewWorld(width, height)
	w.Fill(CharWall)
	w.CarveRoom(&Room{
		X: 1,
		Y: 1,
		W: width - 2,
		H: height - 2,
	})
	return w
}

const (
	CharWall   = '#'
	CharWater  = '~'
	CharFloor  = ' '
	CharTree   = 'T'
	CharShelf  = 'S'
	CharColumn = 'o'
)
