package gamerogueish

import (
	"log"
	"math/rand"
)

// RoomConfig represents a (furnishing) generator configuration for a room type.
// TODO: Allow for random attributes like, columns, puddles, shape, etc.
type RoomConfig struct {
	Name       string                   // Name of the room type.
	HasColumns bool                     // Indicates if the room has columns.
	HasPuddle  bool                     // Indicates if the room has a puddle.
	Condition  func(*Room, *World) bool // Condition to check if the room can be placed.
	Placers    []*Placer
}

// Apply applies the room config to the given room.
func (rc *RoomConfig) Apply(room *Room, w *World) {
	// TODO: Also do puddles, columns, etc.
	if rc.HasColumns {
		w.AddRoomColumns(room)
	}

	if rc.HasPuddle {
		w.AddRoomPuddle(room)
	}

	// Add furnishings.
	for _, placer := range rc.Placers {
		if rand.Intn(100) < placer.Chance {
			if placer.Func(room, w) {
				for _, child := range placer.Children {
					child.Func(room, w)
				}
			}
		}
	}
}

var (
	// RoomConfigBedroom is a bedroom.
	RoomConfigBedroom = &RoomConfig{
		Name: "bedroom",
		Condition: func(room *Room, w *World) bool {
			// Make sure we are a room with a single entrance and we are not too big or too small.
			if len(room.Connections) != 1 || room.W < 2 || room.H < 2 || room.W > 10 || room.H > 10 {
				return false
			}
			return true
		},
		Placers: []*Placer{
			placerBed,
			placerDesk,
		},
	}
	// RoomConfigAltar is an altar room.
	RoomConfigAltar = &RoomConfig{
		Name:       "altar",
		HasColumns: true,
		Condition: func(room *Room, w *World) bool {
			// Make sure we are not too big or too small.
			if room.W < 10 && room.H < 10 {
				return false
			}
			log.Printf("Found altar room at %d,%d", room.X, room.Y)
			return true
		},
		Placers: []*Placer{
			placerAltar,
		},
	}
	// RoomConfigLibrary is a library.
	RoomConfigLibrary = &RoomConfig{
		Name:       "library",
		HasColumns: true,
		Condition: func(room *Room, w *World) bool {
			// Make sure we are not too big or too small.
			if room.W < 4 || room.H < 4 {
				return false
			}
			log.Printf("Found library at %d,%d", room.X, room.Y)
			return true
		},
		Placers: []*Placer{
			placerLibrary,
			placerDesk,
		},
	}
)

type Placer struct {
	Chance   int
	Func     func(*Room, *World) bool
	ItemType *ItemType
	Children []*Placer // Possible accompanying placers.
}

var placerBed = &Placer{
	Chance: 100,
	Func: func(room *Room, w *World) bool {
		// If the room is too small, return.
		if room.W < 2 || room.H < 2 {
			return false
		}

		// TODO: Place bed opposite of entrance.
		if len(room.Connections) == 1 {
			// Pick the opposite wall.
			connX := room.Connections[0].X
			connY := room.Connections[0].Y
			var x, y int
			switch {
			case connX < room.X:
				// West
				x = room.X + room.W - 1
				y = rand.Intn(room.H) + room.Y
			case connX > room.X+room.W:
				// East
				x = room.X
				y = rand.Intn(room.H) + room.Y
			case connY < room.Y:
				// North
				x = rand.Intn(room.W) + room.X
				y = room.Y + room.H - 1
			case connY > room.Y+room.H:
				// South
				x = rand.Intn(room.W) + room.X
				y = room.Y
			}

			// Check if the cell is empty and make sure that at least one neighbor is a wall.
			if !w.IsEmpty(x, y) || !w.NextToWall(x, y) {
				return false
			}

			bed := ItemTypeBed.Generate()
			bed.X = x
			bed.Y = y
			log.Printf("Placing bed at %d,%d", x, y)
			w.Items = append(w.Items, bed)
			return true
		}

		// Place a bed in the room near a wall.
		for i := 0; i < 20; i++ {
			// Check if the cell is empty and make sure that at least one neighbor is a wall.
			x, y := room.RandAlongWall()
			if !w.IsEmpty(x, y) || !w.NextToWall(x, y) {
				continue
			}

			bed := ItemTypeBed.Generate()
			bed.X = x
			bed.Y = y
			log.Printf("Placing bed at %d,%d", x, y)
			w.Items = append(w.Items, bed)
			return true
		}
		return false
	},
	Children: []*Placer{
		placerBedSideTable,
	},
}

var placerBedSideTable = &Placer{
	Chance: 100,
	Func: func(room *Room, w *World) bool {
		// TODO: Avoid blocking exits.
		// Find a bed in the room.
		var bed *Item
		for _, item := range w.Items {
			// TODO: Find better way to identify bed.
			if item.ItemType.Name == ItemTypeBed.Name && room.Contains(item.X, item.Y) {
				log.Printf("Found bed at %d,%d", item.X, item.Y)
				bed = item
				break
			}
		}
		if bed == nil {
			return false
		}
		log.Printf("Found bed at %d,%d", bed.X, bed.Y)

		// Place a side table next to the bed.
		for _, i := range rand.Perm(len(tileNeighboursCardinal)) {
			// Get the neighbor cell.
			x := bed.X + tileNeighboursCardinal[i][0]
			y := bed.Y + tileNeighboursCardinal[i][1]

			// Check if the cell is inside the room, is empty and make sure that at least one neighbor is a wall.
			if !room.Contains(x, y) || !w.IsEmpty(x, y) || !w.NextToWall(x, y) || room.NextToDoor(x, y) {
				continue
			}

			// Place the side table.
			table := ItemTypeSideTable.Generate()
			table.X = x
			table.Y = y
			log.Printf("Placing side table at %d,%d", x, y)
			w.Items = append(w.Items, table)
			return true
		}
		return false
	},
}

var placerAltar = &Placer{
	Chance: 100,
	Func: func(room *Room, w *World) bool {
		// Either put the altar in the center of the room or next to a wall.
		var x, y int
		if rand.Intn(2) == 0 {
			// Center
			x = room.X + room.W/2
			y = room.Y + room.H/2
		} else {
			// Wall
			x, y = room.RandAlongWall()
		}

		// Check if the cell is empty.
		if !w.IsEmpty(x, y) {
			return false
		}

		// TODO: Make sure we don't block any exits.
		// Place the altar.
		altar := ItemTypeAltar.Generate()
		altar.X = x
		altar.Y = y
		log.Printf("Placing altar at %d,%d", x, y)
		w.Items = append(w.Items, altar)
		// Randomly place a candle on the altar.
		if rand.Intn(2) == 0 {
			candle := ItemTypeCandle.Generate()
			candle.X = x
			candle.Y = y
			log.Printf("Placing candle at %d,%d", x, y)
			w.Items = append(w.Items, candle)
		}
		return true
	},
}

var placerDesk = &Placer{
	Chance: 100,
	Func: func(room *Room, w *World) bool {
		// If the room is too small, return.
		if room.W < 2 || room.H < 2 {
			return false
		}

		x, y := room.RandAlongWall()

		// Check if the cell is empty and make sure that at least one neighbor is a wall.
		if !w.IsEmpty(x, y) || !w.NextToWall(x, y) {
			return false
		}

		desk := ItemTypeDesk.Generate()
		desk.X = x
		desk.Y = y
		log.Printf("Placing desk at %d,%d", x, y)
		w.Items = append(w.Items, desk)
		return true
	},
}

var placerLibrary = &Placer{
	Chance: 100,
	Func: func(room *Room, w *World) bool {
		log.Printf("Placing library in room at %d,%d", room.X, room.Y)
		// Line the walls with bookshelves.
		for dx := 0; dx < room.W; dx++ {
			for dy := 0; dy < room.H; dy++ {
				if dx == 0 || dx == room.W-1 || dy == 0 || dy == room.H-1 {
					x, y := room.X+dx, room.Y+dy
					// Check if the cell is empty and that we don't block any exits.
					if !w.IsEmpty(x, y) || room.NextToDoor(x, y) {
						continue
					}

					// Place the bookshelf.
					bookshelf := ItemTypeBookshelf.Generate()
					bookshelf.X = x
					bookshelf.Y = y
					log.Printf("Placing bookshelf at %d,%d", x, y)
					w.Items = append(w.Items, bookshelf)
				}
			}
		}
		return true
	},
}

var tileNeighboursCardinal = [][2]int{
	{0, -1}, // north
	{1, 0},  // east
	{0, 1},  // south
	{-1, 0}, // west
}
