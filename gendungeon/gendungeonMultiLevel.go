package gendungeon

import "fmt"

// ConfigMultiLevel represents the configuration for a multi level dungeon.
type ConfigMultiLevel struct {
	Depth int // Depth is the number of levels in the dungeon.
	Config
}

// DungeonMultiLevel represents a multi level dungeon.
type DungeonMultiLevel struct {
	Levels []*Dungeon
	*ConfigMultiLevel
	Seed int64
}

// GenerateMultiLevelFromConfig generates a 3D dungeon.
func GenerateMultiLevelFromConfig(cfg *ConfigMultiLevel, seed int64) *DungeonMultiLevel {
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
func NewDungeon3D(cfg *ConfigMultiLevel, seed int64) *DungeonMultiLevel {
	d := &DungeonMultiLevel{
		ConfigMultiLevel: cfg,
		Seed:             seed,
	}

	return d
}

func (d *DungeonMultiLevel) createLevels() {
	d.Levels = make([]*Dungeon, d.Depth)
	for i := 0; i < d.Depth; i++ {
		d.Levels[i] = createEmptyDungeon(d.Width, d.Height, d.Seed+int64(i))
	}
}

func (d *DungeonMultiLevel) createRooms(minSize, maxSize, attempts int) {
	for _, level := range d.Levels {
		// TODO: Ensure that we have rooms overlapping between levels.
		level.createRooms(minSize, maxSize, attempts)
	}
}

func (d *DungeonMultiLevel) createMaze() {
	for _, level := range d.Levels {
		level.createMaze()
	}
}

func (d *DungeonMultiLevel) identifyEdges() {
	for _, level := range d.Levels {
		level.identifyEdges()
	}
}

func (d *DungeonMultiLevel) connectRegions() {
	for _, level := range d.Levels {
		level.connectRegions()
	}
}

func (d *DungeonMultiLevel) trimTunnels() {
	for _, level := range d.Levels {
		level.trimTunnels()
	}
}

func (d *DungeonMultiLevel) createStairs() {
	for i := 0; i < len(d.Levels)-1; i++ {
		d.Levels[i].createStairs(d.Levels[i+1])
	}
}

// RenderToConsole renders the dungeon to the console.
func (d *DungeonMultiLevel) RenderToConsole() {
	for i, level := range d.Levels {
		fmt.Println("Level", i)
		level.RenderToConsole()
	}
}
