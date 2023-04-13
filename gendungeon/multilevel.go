package gendungeon

import "fmt"

// DungeonMultiLevel represents a multi level dungeon.
type DungeonMultiLevel struct {
	Levels []*Dungeon
}

// GenerateMultiLevelFromConfig generates a mult level dungeon.
func GenerateMultiLevelFromConfig(cfg Config, n int, seed int64) *DungeonMultiLevel {
	d := NewDungeonMultiLevel()
	d.CreateNLevels(cfg, n, seed)
	d.CreateStairs()
	return d
}

// NewDungeonMultiLevel creates a new multi level dungeon.
func NewDungeonMultiLevel() *DungeonMultiLevel {
	return &DungeonMultiLevel{}
}

// CreateNLevels creates n levels for the dungeon using the supplied config.
func (d *DungeonMultiLevel) CreateNLevels(cfg Config, n int, seed int64) {
	d.Levels = make([]*Dungeon, n)
	for i := range d.Levels {
		level := createEmptyDungeon(cfg.Width, cfg.Height, seed+int64(i))
		// TODO: Ensure that we have rooms overlapping between levels.
		level.createRooms(cfg.MinRoomSize, cfg.MaxRoomSize, cfg.RoomAttempts, cfg.AllowNonRect)
		level.createMaze()
		level.identifyEdges()
		level.connectRegions()
		level.trimTunnels()
		d.Levels[i] = level
	}
}

// AddLevel adds a level to the dungeon.
func (d *DungeonMultiLevel) AddLevel(level *Dungeon) {
	d.Levels = append(d.Levels, level)
}

// CreateStairs creates stairs between levels.
// Call this after all levels have been created.
func (d *DungeonMultiLevel) CreateStairs() {
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
