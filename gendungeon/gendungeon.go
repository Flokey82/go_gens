// Package gendungeon implements a simple dungeon generator forked from
// https://github.com/brad811/go-dungeon
// which is in turn based on
// http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/
package gendungeon

import (
	"fmt"
	"math/rand"

	"github.com/Flokey82/go_gens/utils"
)

// Suggested default values for the dungeon generation.
var (
	RoomAttempts = 200
	MinRoomSize  = 5
	MaxRoomSize  = 15
)

// Config is a configuration for dungeon generation.
type Config struct {
	Width        int
	Height       int
	RoomAttempts int
	MinRoomSize  int
	MaxRoomSize  int
}

// Material represents the material of a tile.
type Material int

// The various valid materials.
const (
	MatWall       Material = iota // stone wall
	MatFloor                      // room floor
	MatDoor                       // door
	MatTunnel                     // tunnel / maze
	MatStairsUp                   // stairs up
	MatStairsDown                 // stairs down
)

// Point is a point at a specific x,y coordinate.
type Point struct {
	X int
	Y int
}

// Rect is a rectangle with a specific width and height.
type Rect struct {
	X, Y, Width, Height int
}

// Center returns the center of the rectangle.
func (r *Rect) Center() Point {
	return Point{r.X + r.Width/2, r.Y + r.Height/2}
}

// Tile represents a tile on the dungeon grid.
type Tile struct {
	Region   int      // Associates a region with a tile.
	Material Material // The material of the tile.
}

// Room represents a room in the dungeon.
type Room struct {
	Width    int     // width of the room
	Height   int     // height of the room
	Location Point   // top left corner of the room
	Edges    []Point // the edges of the room
}

// Overlap finds the rectangle representing the overlap between two rooms.
func (r *Room) Overlap(r2 Room) (Rect, bool) {
	// Find the edges of the overlap.
	left := utils.Max(r.Location.X, r2.Location.X)
	right := utils.Min(r.Location.X+r.Width, r2.Location.X+r2.Width)
	top := utils.Max(r.Location.Y, r2.Location.Y)
	bottom := utils.Min(r.Location.Y+r.Height, r2.Location.Y+r2.Height)

	// If the edges define a rectangle, then the rooms overlap.
	if left < right && top < bottom {
		return Rect{left, top, right - left, bottom - top}, true
	}

	return Rect{}, false
}

// Dungeon represents a generated dungeon.
type Dungeon struct {
	Tiles      [][]Tile   // dungeon grid
	Rooms      []Room     // rooms in the dungeon
	Width      int        // width of the dungeon
	Height     int        // height of the dungeon
	numRegions int        // number of regions in the dungeon
	rand       *rand.Rand // rand initialized with the seed
}

// Generate generates a new dungeon with the given width and height.
// Further parameters include the number of attempts to use to place rooms,
// the minimum and maximum room size, and the seed to use.
func Generate(width, height, roomAttempts, minRoomSize, maxRoomSize int, seed int64) *Dungeon {
	cfg := Config{
		Width:        width,
		Height:       height,
		RoomAttempts: roomAttempts,
		MinRoomSize:  minRoomSize,
		MaxRoomSize:  maxRoomSize,
	}
	return GenerateFromConfig(cfg, seed)
}

// GenerateFromConfig generates a new dungeon with the given configuration.
func GenerateFromConfig(cfg Config, seed int64) *Dungeon {
	dng := createEmptyDungeon(cfg.Width, cfg.Height, seed)
	dng.createRooms(cfg.MinRoomSize, cfg.MaxRoomSize, cfg.RoomAttempts)
	dng.createMaze()
	dng.identifyEdges()
	dng.connectRegions()
	dng.trimTunnels()
	return dng
}

func createEmptyDungeon(width, height int, seed int64) *Dungeon {
	fmt.Println("Creating empty dungeon...")
	dng := &Dungeon{
		Width:  width,
		Height: height,
		rand:   rand.New(rand.NewSource(seed)),
	}
	dng.Tiles = make([][]Tile, height)
	for i := range dng.Tiles {
		dng.Tiles[i] = make([]Tile, width)
	}
	return dng
}

func (dng *Dungeon) createRooms(minSize, maxSize, attempts int) {
	fmt.Println("Creating rooms...")
	var rooms []Room
	for i := 0; i < attempts; i++ {
		width := dng.rand.Intn(maxSize-minSize) + minSize
		height := dng.rand.Intn(maxSize-minSize) + minSize

		maxX := dng.Width - width - 2
		maxY := dng.Height - height - 2

		x := dng.rand.Intn(maxX-3) + 3
		y := dng.rand.Intn(maxY-3) + 3

		shouldAppend := true
		for _, r := range rooms {
			if x+width < r.Location.X || // to the left
				x > r.Location.X+r.Width || // to the right
				y+height < r.Location.Y || // fully above
				y > r.Location.Y+r.Height { // fully below
				continue // do nothing
			}
			shouldAppend = false
			break
		}

		if shouldAppend {
			rooms = append(rooms, Room{
				Width:    width,
				Height:   height,
				Location: Point{X: x, Y: y},
			})
		}
	}

	for _, r := range rooms {
		dng.numRegions++
		for i := r.Location.X; i < r.Location.X+r.Width; i++ {
			for j := r.Location.Y; j < r.Location.Y+r.Height; j++ {
				dng.Tiles[j][i].Material = MatFloor
				dng.Tiles[j][i].Region = dng.numRegions
			}
		}
	}

	dng.Rooms = rooms
}

func (dng *Dungeon) createMaze() {
	fmt.Println("Creating tunnels...")

	// Iterate through all tiles and start growing the maze from each elegible tile.
	for x := 1; x < dng.Width-1; x++ {
		for y := 1; y < dng.Height-1; y++ {
			// A maze is only generated at the given (stone) tile if
			// it is sourrounded by walls (stone).
			//    y
			//    |
			//   000
			// x-000
			//   000
			if dng.Tiles[y-1][x-1].Material == MatWall &&
				dng.Tiles[y][x-1].Material == MatWall &&
				dng.Tiles[y+1][x-1].Material == MatWall &&
				dng.Tiles[y-1][x].Material == MatWall &&
				dng.Tiles[y][x].Material == MatWall &&
				dng.Tiles[y+1][x].Material == MatWall &&
				dng.Tiles[y-1][x+1].Material == MatWall &&
				dng.Tiles[y][x+1].Material == MatWall &&
				dng.Tiles[y+1][x+1].Material == MatWall {

				dng.numRegions++
				dng.continueMaze(x, y)
			}
		}
	}
}

func (dng *Dungeon) continueMaze(x, y int) {
	var validTiles []Point

	// Check if we can dig left.
	//     x
	//     |
	//   00?
	// y-00?
	//   00?
	if x-2 >= 0 &&
		dng.Tiles[y][x-1].Material == MatWall &&
		dng.Tiles[y][x-2].Material == MatWall &&
		dng.Tiles[y+1][x-2].Material == MatWall &&
		dng.Tiles[y-1][x-2].Material == MatWall &&
		dng.Tiles[y+1][x-1].Material == MatWall &&
		dng.Tiles[y-1][x-1].Material == MatWall {
		validTiles = append(validTiles, Point{X: x - 1, Y: y})
	}

	// Check if we can dig right.
	//   x
	//   |
	//   ?00
	// y-?00
	//   ?00
	if x+2 < dng.Width &&
		dng.Tiles[y][x+1].Material == MatWall &&
		dng.Tiles[y][x+2].Material == MatWall &&
		dng.Tiles[y-1][x+2].Material == MatWall &&
		dng.Tiles[y+1][x+2].Material == MatWall &&
		dng.Tiles[y+1][x+1].Material == MatWall &&
		dng.Tiles[y-1][x+1].Material == MatWall {
		validTiles = append(validTiles, Point{X: x + 1, Y: y})
	}

	// Check if we can dig up.
	//    x
	//    |
	//   000
	//   000
	// y-???
	if y-2 >= 0 &&
		dng.Tiles[y-1][x].Material == MatWall &&
		dng.Tiles[y-2][x].Material == MatWall &&
		dng.Tiles[y-2][x-1].Material == MatWall &&
		dng.Tiles[y-2][x+1].Material == MatWall &&
		dng.Tiles[y-1][x-1].Material == MatWall &&
		dng.Tiles[y-1][x+1].Material == MatWall {
		validTiles = append(validTiles, Point{X: x, Y: y - 1})
	}

	// Check if we can dig down.
	//    x
	//    |
	// y-???
	//   000
	//   000
	if y+2 < dng.Height &&
		dng.Tiles[y+1][x].Material == MatWall &&
		dng.Tiles[y+2][x].Material == MatWall &&
		dng.Tiles[y+2][x-1].Material == MatWall &&
		dng.Tiles[y+2][x+1].Material == MatWall &&
		dng.Tiles[y+1][x-1].Material == MatWall &&
		dng.Tiles[y+1][x+1].Material == MatWall {
		validTiles = append(validTiles, Point{X: x, Y: y + 1})
	}

	// Check if we have any valid tiles for expanding the maze.
	if len(validTiles) == 0 {
		return // Nothing to do
	}

	// From all valid tiles, pick one (at random if we have more than one)
	// to dig a tunnel to.
	var idx int
	if len(validTiles) > 1 {
		idx = dng.rand.Intn(len(validTiles))
	}

	// Set the new tile to be a tunnel.
	point := validTiles[idx]
	dng.Tiles[point.Y][point.X].Material = MatTunnel
	dng.Tiles[point.Y][point.X].Region = dng.numRegions

	// Continue digging the maze...
	dng.continueMaze(point.X, point.Y) // ... from the new tile.
	dng.continueMaze(x, y)             // ... from the current tile.
}

func (dng *Dungeon) identifyEdges() {
	fmt.Println("Identifying edges...")

	// Iterate through all rooms and identify edges.
	for i := range dng.Rooms {
		x := dng.Rooms[i].Location.X
		y := dng.Rooms[i].Location.Y

		// Iterate along the width (top and bottom) of the room.
		for j := x; j < x+dng.Rooms[i].Width; j++ {
			// Check if we border on a tunnel or a floor.
			if dng.Tiles[y-2][j].Material == MatTunnel ||
				dng.Tiles[y-2][j].Material == MatFloor {
				dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: j, Y: y - 1})
			}
			if dng.Tiles[y+dng.Rooms[i].Height+1][j].Material == MatTunnel ||
				dng.Tiles[y+dng.Rooms[i].Height+1][j].Material == MatFloor {
				dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: j, Y: y + dng.Rooms[i].Height})
			}
		}

		// Iterate along the height (left and right) of the room.
		for k := y; k < y+dng.Rooms[i].Height; k++ {
			// Check if we border on a tunnel or a floor.
			if dng.Tiles[k][x-2].Material == MatTunnel ||
				dng.Tiles[k][x-2].Material == MatFloor {
				dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: x - 1, Y: k})
			}
			if dng.Tiles[k][x+dng.Rooms[i].Width+1].Material == MatTunnel ||
				dng.Tiles[k][x+dng.Rooms[i].Width+1].Material == MatFloor {
				dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: x + dng.Rooms[i].Width, Y: k})
			}
		}
	}
}

func (dng *Dungeon) connectRegions() {
	fmt.Println("Conneting regions...")

	// Iterate through all rooms and connect them to the corridors or other rooms.
	for i, room := range dng.Rooms {
		// Pick a random edge to connect to.
		edge := room.Edges[dng.rand.Intn(len(dng.Rooms[i].Edges))]
		roomRegion := dng.Tiles[dng.Rooms[i].Location.Y][dng.Rooms[i].Location.X].Region

		// The neighboring tiles.
		nbs := [8]Tile{
			dng.Tiles[edge.Y-1][edge.X-1],
			dng.Tiles[edge.Y-1][edge.X],
			dng.Tiles[edge.Y-1][edge.X+1],
			dng.Tiles[edge.Y][edge.X-1],
			dng.Tiles[edge.Y][edge.X+1],
			dng.Tiles[edge.Y+1][edge.X-1],
			dng.Tiles[edge.Y+1][edge.X],
			dng.Tiles[edge.Y+1][edge.X+1],
		}

		// Check if edge is unconnected.
		for _, j := range dng.rand.Perm(len(nbs)) {
			// A suitable tile is one that is a tunnel or floor and is not part of the room.
			if (nbs[j].Material == MatFloor || nbs[j].Material == MatTunnel) &&
				nbs[j].Region != roomRegion {
				// We found a suitable unexpected edge, so set it to a door.
				dng.Tiles[edge.Y][edge.X].Material = MatDoor
				for x := room.Location.X; x < room.Location.X+room.Width; x++ {
					for y := room.Location.Y; y < room.Location.Y+room.Height; y++ {
						dng.Tiles[y][x].Region = nbs[j].Region
					}
				}
				break // We found a suitable edge, so stop looking.
			}
		}
	}

	// Iterate through the rooms and their edges in random order
	// to see if any of them are still a separate region.
	connectedRegions := make(map[int]bool)
RoomsLoop:
	for _, i := range dng.rand.Perm(len(dng.Rooms)) {
		for _, j := range dng.rand.Perm(len(dng.Rooms[i].Edges)) {
			room := dng.Rooms[i]
			edge := room.Edges[j]
			x := edge.X
			y := edge.Y

			// The neighboring tiles (only NSEW).
			//  N
			// WxE
			//  S
			nbs := [4]Point{
				{X: x - 1, Y: y},
				{X: x + 1, Y: y},
				{X: x, Y: y - 1},
				{X: x, Y: y + 1},
			}

			curRegion := -1
			for _, sp := range nbs {
				tile := dng.Tiles[sp.Y][sp.X]
				if curRegion == -1 && tile.Region != 0 {
					curRegion = tile.Region
				} else if tile.Region != curRegion &&
					tile.Region != 0 &&
					!connectedRegions[tile.Region] {
					dng.Tiles[y][x].Material = MatDoor
					connectedRegions[tile.Region] = true
					connectedRegions[curRegion] = true
					continue RoomsLoop
				}
			}
		}
	}
}

func (dng *Dungeon) trimTunnels() {
	fmt.Println("Trimming tunnels...")
	for x := 1; x < dng.Width-1; x++ {
		for y := 1; y < dng.Height-1; y++ {
			dng.continueTrimTunnels(x, y)
		}
	}
}

func (dng *Dungeon) continueTrimTunnels(x int, y int) {
	if dng.Tiles[y][x].Material == MatTunnel || dng.Tiles[y][x].Material == MatDoor {
		var wallCount int
		var nextPoint Point

		// The neighboring tiles (only NSEW).
		//  N
		// WxE
		//  S
		nbs := [4]Point{
			{X: x - 1, Y: y},
			{X: x + 1, Y: y},
			{X: x, Y: y - 1},
			{X: x, Y: y + 1},
		}

		for _, nb := range nbs {
			tile := dng.Tiles[nb.Y][nb.X]
			switch tile.Material {
			case MatWall:
				wallCount++
			case MatTunnel, MatDoor:
				nextPoint = Point{
					X: nb.X,
					Y: nb.Y,
				}
			}
		}

		if wallCount >= 3 {
			dng.Tiles[y][x].Material = MatWall
			dng.Tiles[y][x].Region = 0
			if nextPoint.X != 0 || nextPoint.Y != 0 {
				dng.continueTrimTunnels(nextPoint.X, nextPoint.Y)
			}
		}
	}
}

func (dng *Dungeon) createStairs(dngUp *Dungeon) {
	fmt.Println("Creating stairs...")

	// Find rooms that overlap with the previous dungeon,
	// then add stairs in each matching pair.
	var stairsUpDown Point
	for _, room := range dng.Rooms {
		for _, roomUp := range dngUp.Rooms {
			overlap, ok := room.Overlap(roomUp)
			if ok {
				// We found a matching pair, so add stairs.
				stairsUpDown = overlap.Center()
				break
			}
		}
	}

	// If we didn't find any matching rooms, then complain.
	if stairsUpDown.X == 0 && stairsUpDown.Y == 0 {
		fmt.Println("ERROR: No matching rooms found for stairs!")
		return

	}

	// Add the stairs to the dungeon.
	dng.Tiles[stairsUpDown.Y][stairsUpDown.X].Material = MatStairsUp
	dngUp.Tiles[stairsUpDown.Y][stairsUpDown.X].Material = MatStairsDown
}

// RenderToConsole prints the dungeon layout to the console.
func (dng *Dungeon) RenderToConsole() {
	fmt.Println("Dungeon: (", dng.Width, ",", dng.Height, ") Regions: ", dng.numRegions)

	for y := 0; y < dng.Height; y++ {
		for x := 0; x < dng.Width; x++ {
			switch dng.Tiles[y][x].Material {
			case MatWall:
				fmt.Print("0 ")
			case MatFloor:
				fmt.Print("= ")
			case MatDoor:
				fmt.Print("| ")
			case MatTunnel:
				fmt.Print("- ")
			case MatStairsUp:
				fmt.Print("U ")
			case MatStairsDown:
				fmt.Print("D ")
			default:
				fmt.Print("ER")
			}
		}
		fmt.Println()
	}
}
