package gamerogueish

import (
	"math/rand"
	"strings"
)

var levelLayout = `
#####################
#         #    #    #
#    #    #         #
#    ######    #    #
#              #    #
##  #############  ##
#    #    #    #    #
#    #         #    #
#    ######         #
#              #    #
#####################
`

// genWorld returns the default world as world cells.
func genWorld(width, height int) [][]byte {
	var cells [][]byte
	lines := strings.Split(levelLayout, "\n")
	for i := range lines {
		if len(lines[i]) == 0 {
			continue
		}
		cells = append(cells, []byte(lines[i]))
	}
	return cells
}

type World struct {
	Cells  [][]byte // 2D array of world cells
	Width  int      // width of the world in cells
	Height int      // height of the world in cells
}

func newWorld(width, height int) *World {
	w := &World{
		Width:  width,
		Height: height,
	}
	w.Cells = make([][]byte, height)
	for y := range w.Cells {
		w.Cells[y] = make([]byte, width)
	}
	return w
}

// IsSolid checks if a tile is solid (tile content is not a space ' ' character).
func (w *World) IsSolid(x int, y int) bool {
	return w.Cells[y][x] != ' '
}

// CanMoveTo checks if a tile is solid and if it is not occupied by an entity.
func (w *World) CanMoveTo(x, y int) bool {
	// TODO: Bounds check.
	return w.Cells[y][x] == ' '
}

func (w *World) Fill(c byte) {
	for y := range w.Cells {
		for x := range w.Cells[y] {
			w.Cells[y][x] = c
		}
	}
}

func (w *World) CarveRoom(room *Room) {
	for y := room.Y; y < room.Y+room.H; y++ {
		for x := room.X; x < room.X+room.W; x++ {
			w.Cells[y][x] = ' '
		}
	}
}

// Cardinal directions.
const (
	DirNorth = 0
	DirEast  = 1
	DirSouth = 2
	DirWest  = 3
)

// Room generation algorithm:
func GenSimpleDungeon(width, height int) *World {
	const (
		attempts    = 200
		maxRooms    = 100
		minRoomSize = 4
		maxRoomSize = 20
	)
	w := newWorld(width, height)
	w.Fill('#')

	var rooms []*Room

	// Start with a single room.
	rooms = append(rooms, &Room{
		X: (width / 2) - 3,
		Y: (height / 2) - 2,
		W: 6,
		H: 4,
	})

	w.CarveRoom(rooms[0])
	for i := 0; i < attempts; i++ {
		// Pick a room and place a neighboring room.
		room := rooms[rand.Intn(len(rooms))]
		// Pick a random direction.
		dir := rand.Intn(4)
		// Pick a random length and width.
		rl := randInt(minRoomSize, maxRoomSize)
		rw := randInt(minRoomSize, maxRoomSize)

		// Calculate position based on direction.
		// NOTE: Right now we center the neighboring room.
		var x, y int
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

		// Create a new room.
		newRoom := &Room{
			X: x,
			Y: y,
			W: rw,
			H: rl,
		}

		// Check if the new room overlaps with any existing rooms.
		if newRoom.X < 1 || newRoom.Y < 1 || newRoom.X+newRoom.W > width-1 || newRoom.Y+newRoom.H > height-1 || newRoom.Overlaps(rooms) {
			continue
		}

		// Append the room to the dungeon.
		rooms = append(rooms, newRoom)

		// Draw room.
		w.CarveRoom(newRoom)

		// Draw a tunnel between the rooms.
		// NOTE: Right now, we just place the door in the middle.
		switch dir {
		case DirNorth:
			w.Cells[room.Y-1][room.X+room.W/2] = ' '
		case DirEast:
			w.Cells[room.Y+room.H/2][room.X+room.W] = ' '
		case DirSouth:
			w.Cells[room.Y+room.H][room.X+room.W/2] = ' '
		case DirWest:
			w.Cells[room.Y+room.H/2][room.X-1] = ' '
		}
		// Stop if we have enough rooms.
		if len(rooms) > maxRooms {
			break
		}
	}
	return w
}

// Room represents a room in the world.
// TODO: Store connecting rooms
type Room struct {
	X, Y, W, H int
}

func (r *Room) Overlaps(rooms []*Room) bool {
	for _, room := range rooms {
		if r.X+r.W < room.X || r.X > room.X+room.W {
			continue
		}
		if r.Y+r.H < room.Y || r.Y > room.Y+room.H {
			continue
		}
		return true
	}
	return false
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func GenBigBox(width, height int) *World {
	w := newWorld(width, height)
	w.Fill('#')
	w.CarveRoom(&Room{
		X: 1,
		Y: 1,
		W: width - 2,
		H: height - 2,
	})
	return w
}
