package gendungeon

import "fmt"

// Config3D represents the configuration for a 3D dungeon.
type Config3D struct {
	Depth int // Depth is the number of levels in the dungeon.
	Config
}

// Dungeon3D represents a 3D dungeon.
type Dungeon3D struct {
	Levels []*Dungeon
	*Config3D
	Seed int64
}

// Generate3DFromConfig generates a 3D dungeon.
func Generate3DFromConfig(cfg *Config3D, seed int64) *Dungeon3D {
	d := NewDungeon3D(cfg, seed)
	d.createLevels()
	d.createRooms(cfg.MinRoomSize, cfg.MaxRoomSize, cfg.RoomAttempts)
	d.createMaze()
	d.identifyEdges()
	d.connectRegions()
	d.trimTunnels()
	d.createStairs()
	return d
}

// NewDungeon3D creates a new 3D dungeon.
func NewDungeon3D(cfg *Config3D, seed int64) *Dungeon3D {
	d := &Dungeon3D{
		Config3D: cfg,
		Seed:     seed,
	}

	return d
}

func (d *Dungeon3D) createLevels() {
	d.Levels = make([]*Dungeon, d.Depth)
	for i := 0; i < d.Depth; i++ {
		d.Levels[i] = createEmptyDungeon(d.Width, d.Height, d.Seed+int64(i))
	}
}

func (d *Dungeon3D) createRooms(minSize, maxSize, attempts int) {
	for _, level := range d.Levels {
		// TODO: Ensure that we have rooms overlapping between levels.
		level.createRooms(minSize, maxSize, attempts)
	}
}

func (d *Dungeon3D) createMaze() {
	for _, level := range d.Levels {
		level.createMaze()
	}
}

func (d *Dungeon3D) identifyEdges() {
	for _, level := range d.Levels {
		level.identifyEdges()
	}
}

func (d *Dungeon3D) connectRegions() {
	for _, level := range d.Levels {
		level.connectRegions()
	}
}

func (d *Dungeon3D) trimTunnels() {
	for _, level := range d.Levels {
		level.trimTunnels()
	}
}

func (d *Dungeon3D) createStairs() {
	for i := 0; i < len(d.Levels)-1; i++ {
		d.Levels[i].createStairs(d.Levels[i+1])
	}
}

// RenderToConsole renders the dungeon to the console.
func (d *Dungeon3D) RenderToConsole() {
	for i, level := range d.Levels {
		fmt.Println("Level", i)
		level.RenderToConsole()
	}
}
