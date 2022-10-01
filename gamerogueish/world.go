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
func genWorld() [][]byte {
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

type World [][]byte

func newWorld(width, height int) *World {
	w := World(genWorldBoxy(width, height))
	return &w
}

func (w *World) fill(c byte) {
	for y := range *w {
		for x := range (*w)[y] {
			(*w)[y][x] = c
		}
	}
}

func (w *World) carveRoom(room *Room) {
	for y := room.y; y < room.y+room.h; y++ {
		for x := room.x; x < room.x+room.w; x++ {
			(*w)[y][x] = ' '
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
func genFancyWorld(width, height int) *World {
	const (
		attempts    = 200
		maxRooms    = 100
		minRoomSize = 4
		maxRoomSize = 20
	)
	w := newWorld(width, height)
	w.fill('#')

	var rooms []*Room

	// Start with a single room.
	rooms = append(rooms, &Room{
		x: (width / 2) - 3,
		y: (height / 2) - 2,
		w: 6,
		h: 4,
	})

	w.carveRoom(rooms[0])
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
			x = room.x + room.w/2 - rw/2
			y = room.y - rl - 1
		case DirEast:
			x = room.x + room.w + 1
			y = room.y + room.h/2 - rl/2
		case DirSouth:
			x = room.x + room.w/2 - rw/2
			y = room.y + room.h + 1
		case DirWest:
			x = room.x - rw - 1
			y = room.y + room.h/2 - rl/2
		}

		// Create a new room.
		newRoom := &Room{
			x: x,
			y: y,
			w: rw,
			h: rl,
		}

		// Check if the new room overlaps with any existing rooms.
		if newRoom.x < 1 || newRoom.y < 1 || newRoom.x+newRoom.w > width-1 || newRoom.y+newRoom.h > height-1 || newRoom.overlaps(rooms) {
			continue
		}

		// Append the room to the dungeon.
		rooms = append(rooms, newRoom)

		// Draw room.
		w.carveRoom(newRoom)

		// Draw a tunnel between the rooms.
		// NOTE: Right now, we just place the door in the middle.
		switch dir {
		case DirNorth:
			(*w)[room.y-1][room.x+room.w/2] = ' '
		case DirEast:
			(*w)[room.y+room.h/2][room.x+room.w] = ' '
		case DirSouth:
			(*w)[room.y+room.h][room.x+room.w/2] = ' '
		case DirWest:
			(*w)[room.y+room.h/2][room.x-1] = ' '
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
	x, y, w, h int
}

func (r *Room) overlaps(rooms []*Room) bool {
	for _, room := range rooms {
		if r.x+r.w < room.x || r.x > room.x+room.w {
			continue
		}
		if r.y+r.h < room.y || r.y > room.y+room.h {
			continue
		}
		return true
	}
	return false
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func genWorldBoxy(width, height int) [][]byte {
	var cells [][]byte
	for y := 0; y < height; y++ {
		var row []byte
		for x := 0; x < width; x++ {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				row = append(row, '#')
			} else {
				row = append(row, '.')
			}
		}
		cells = append(cells, row)
	}
	return cells
}
