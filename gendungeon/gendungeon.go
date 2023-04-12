// Package gendungeon implements a simple dungeon generator forked from
// https://github.com/brad811/go-dungeon
// which is in turn based on
// http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/
package gendungeon

import (
	"fmt"
	"math"
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
	AllowOval    bool
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

// RoomStyle represents the style of a room.
type RoomStyle int

// The various valid room styles.
const (
	RoomStyleRect RoomStyle = iota
	RoomStyleOval
)

// Room represents a room in the dungeon.
type Room struct {
	Width    int       // width of the room
	Height   int       // height of the room
	Location Point     // top left corner of the room
	Edges    []Point   // the edges of the room
	Style    RoomStyle // style / shape of the room
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

// Center returns the center of the room.
func (r *Room) Center() Point {
	return Point{r.Location.X + r.Width/2, r.Location.Y + r.Height/2}
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
	dng.createRooms(cfg.MinRoomSize, cfg.MaxRoomSize, cfg.RoomAttempts, cfg.AllowOval)
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

func (dng *Dungeon) createRooms(minSize, maxSize, attempts int, allowOval bool) {
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
			r := Room{
				Width:    width,
				Height:   height,
				Location: Point{X: x, Y: y},
				Style:    RoomStyleRect,
			}

			// 20% chance of making the room oval (if allowed).
			if allowOval && dng.rand.Intn(100) < 20 {
				r.Style = RoomStyleOval
			}
			rooms = append(rooms, r)
		}
	}

	// Draw the rooms.
	ovalMargin := 0.5 // margin for oval rooms to make the ends less pointy.
	for _, r := range rooms {
		dng.numRegions++

		switch r.Style {
		case RoomStyleRect:
			// Draw the room as a rectangle.
			for i := r.Location.X; i < r.Location.X+r.Width; i++ {
				for j := r.Location.Y; j < r.Location.Y+r.Height; j++ {
					dng.Tiles[j][i].Material = MatFloor
					dng.Tiles[j][i].Region = dng.numRegions
				}
			}
		case RoomStyleOval:
			// Get center of room.
			center := r.Center()

			// Draw the room as an oval / ellipse.
			for i := r.Location.X; i < r.Location.X+r.Width; i++ {
				for j := r.Location.Y; j < r.Location.Y+r.Height; j++ {
					// Check if the tile is within the oval / ellipse.
					if (math.Pow(float64(i-center.X), 2)/math.Pow(float64(r.Width/2)+ovalMargin, 2) +
						math.Pow(float64(j-center.Y), 2)/math.Pow(float64(r.Height/2)+ovalMargin, 2)) <= 1 {
						dng.Tiles[j][i].Material = MatFloor
						dng.Tiles[j][i].Region = dng.numRegions
					}
				}
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
		switch dng.Rooms[i].Style {
		case RoomStyleRect:
			x := dng.Rooms[i].Location.X
			y := dng.Rooms[i].Location.Y

			// Iterate along the width (top and bottom) of the room to identify edges.
			//
			//    ? ? ?
			//  0 ^ ^ ^ 0
			//  0 = = = 0
			//  0 = = = 0 ^, v = Possible edges, if ? is a tunnel or floor
			//  0 = = = 0
			//  0 v v v 0
			//    ? ? ?
			//
			for j := x; j < x+dng.Rooms[i].Width; j++ {
				// Check if we border on a tunnel or a floor.
				if dng.Tiles[y-2][j].Material == MatTunnel ||
					dng.Tiles[y-2][j].Material == MatFloor {
					//   ?
					// 0 ^ 0 ^ = Edge (y-2), j = Tile (j, y), ? = y-2
					// 0 j 0
					dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: j, Y: y - 1})
				}
				if dng.Tiles[y+dng.Rooms[i].Height+1][j].Material == MatTunnel ||
					dng.Tiles[y+dng.Rooms[i].Height+1][j].Material == MatFloor {
					// 0 j 0
					// 0 v 0 v = Edge (y+height), j = Tile (j, y), ? = y+height+1
					//   ?
					dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: j, Y: y + dng.Rooms[i].Height})
				}
			}

			// Iterate along the height (left and right) of the room to identify edges.
			//
			//   0 0 0 0 0
			// ? < = = = > ?
			// ? < = = = > ? <, > = Possible edges, if ? is a tunnel or floor
			// ? < = = = > ?
			//   0 0 0 0 0
			//
			for k := y; k < y+dng.Rooms[i].Height; k++ {
				// Check if we border on a tunnel or a floor.
				if dng.Tiles[k][x-2].Material == MatTunnel ||
					dng.Tiles[k][x-2].Material == MatFloor {
					//   0 0
					// ? < k   < = Edge (x-1), k = Tile (x, k), ? = x-2
					//   0 0
					dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: x - 1, Y: k})
				}
				if dng.Tiles[k][x+dng.Rooms[i].Width+1].Material == MatTunnel ||
					dng.Tiles[k][x+dng.Rooms[i].Width+1].Material == MatFloor {
					//   0 0
					//   k > ? > = Edge (x+width), k = Tile (x, k), ? = x+width+1
					//   0 0
					dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, Point{X: x + dng.Rooms[i].Width, Y: k})
				}
			}
		case RoomStyleOval:
			x := dng.Rooms[i].Location.X
			y := dng.Rooms[i].Location.Y

			// Iterate along the width (top and bottom) of the room.
			for j := x; j < x+dng.Rooms[i].Width; j++ {
				// Iterate along the height (left and right) of the room.
				for k := y; k < y+dng.Rooms[i].Height; k++ {
					// Skip tiles that are not floors.
					if dng.Tiles[k][j].Material != MatFloor {
						continue
					}

					var wallCount int
					// Count how many walls we border on.
					for nbX := j - 1; nbX <= j+1; nbX++ {
						for nbY := k - 1; nbY <= k+1; nbY++ {
							// Skip the current tile.
							if nbX == j && nbY == k {
								continue
							}
							if dng.Tiles[nbY][nbX].Material == MatWall {
								wallCount++
							}
						}
					}

					// If we have at least 3 walls as direct neighbors, we are an edge.
					if wallCount <= 2 {
						continue
					}
					if ok, edge := dng.hasTunnelOrFloorNearby(j, k); ok {
						dng.Rooms[i].Edges = append(dng.Rooms[i].Edges, edge)
					}
				}
			}
		}
	}
}

// hasTunnelOrFloorNearby returns true if the tile at x, y has a tunnel or floor
// tile within 2 tiles of it, and returns the tile that separates the two regions.
func (dng *Dungeon) hasTunnelOrFloorNearby(x, y int) (bool, Point) {
	region := dng.Tiles[y][x].Region
	coords := [4]Point{
		{X: x, Y: y - 2},
		{X: x - 2, Y: y},
		{X: x + 2, Y: y},
		{X: x, Y: y + 2},
	}

	for _, coord := range coords {
		nb := dng.Tiles[coord.Y][coord.X]
		if (nb.Material == MatTunnel || nb.Material == MatFloor) && nb.Region != region {
			return true, Point{X: (coord.X + x) / 2, Y: (coord.Y + y) / 2}
		}
	}

	return false, Point{}
}

func (dng *Dungeon) connectRegions() {
	fmt.Println("Connecting regions...")

	// NOTE: This code is a bit... weird. I'll have to think this through a bit more since
	// it seems kinda unnecessary to have a two-pass approach.

	// Iterate through all rooms and connect them to the corridors or other rooms.
	connectedTo := make(map[int]int)
	for _, room := range dng.Rooms {
		roomCenter := room.Center()
	Loop:
		for _, i := range dng.rand.Perm(len(room.Edges)) {
			// Pick a random edge to connect to.
			edge := room.Edges[i]
			roomRegion := dng.Tiles[roomCenter.Y][roomCenter.X].Region

			// The neighboring tiles.
			nbs := [8]Tile{
				dng.Tiles[edge.Y-1][edge.X],
				dng.Tiles[edge.Y][edge.X-1],
				dng.Tiles[edge.Y][edge.X+1],
				dng.Tiles[edge.Y+1][edge.X],
				dng.Tiles[edge.Y-1][edge.X-1],
				dng.Tiles[edge.Y-1][edge.X+1],
				dng.Tiles[edge.Y+1][edge.X-1],
				dng.Tiles[edge.Y+1][edge.X+1],
			}

			// Check if edge is unconnected.
			for _, j := range dng.rand.Perm(len(nbs)) {
				// A suitable tile is one that is a tunnel or floor and is not part of the room.
				if (nbs[j].Material == MatFloor || nbs[j].Material == MatTunnel) &&
					nbs[j].Region != roomRegion {
					// We found a suitable unexpected edge, so set it to a door.
					dng.Tiles[edge.Y][edge.X].Material = MatDoor

					// Make sure we remember which region this room is connected to.
					connectedTo[roomRegion] = nbs[j].Region
					break Loop // We found a suitable edge, so stop looking.
				}
			}
		}
	}

	// Iterate through the rooms and their edges in random order
	// to see if any of them are still a separate region.
	connectedRegions := make(map[int]bool)
RoomsLoop:
	for _, i := range dng.rand.Perm(len(dng.Rooms)) {
		room := dng.Rooms[i]
		roomCenter := room.Center()
		roomRegion := dng.Tiles[roomCenter.Y][roomCenter.X].Region

		// Find a suitable edge to connect to and make sure we only connect to a new region.
		for _, j := range dng.rand.Perm(len(dng.Rooms[i].Edges)) {
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

				// Make sure we don't connect to the same region or a region that is not set.
				if connectedTo[roomRegion] == tile.Region || tile.Region == 0 {
					continue
				}

				// If we haven't found a region yet, set it.
				if curRegion == -1 {
					curRegion = tile.Region
					continue
				}

				// Now check if we have found a new, different region.
				if !connectedRegions[tile.Region] {
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
	var foundPoint bool
Loop:
	for _, room := range dng.Rooms {
		for _, roomUp := range dngUp.Rooms {
			overlap, ok := room.Overlap(roomUp)
			if ok {
				// We found a matching pair, so add stairs.
				stairsUpDown = overlap.Center()
				// Make sure that the center is a floor tile for both levels.
				if dng.Tiles[stairsUpDown.Y][stairsUpDown.X].Material == MatFloor &&
					dngUp.Tiles[stairsUpDown.Y][stairsUpDown.X].Material == MatFloor {
					foundPoint = true
					break Loop // We found a suitable pair, so stop looking.
				}
			}
		}
	}

	// If we didn't find any matching rooms, then complain.
	if !foundPoint {
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
