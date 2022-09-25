package gamecs

import "github.com/Flokey82/go_gens/vectors"

type CAiMemory struct {
	ai        *CompAi
	Locations map[string]*Location
	Positions map[string]vectors.Vec2
}

func newCAiMemory() *CAiMemory {
	return &CAiMemory{
		Locations: make(map[string]*Location),
		Positions: make(map[string]vectors.Vec2),
	}
}

func (c *CAiMemory) init(ai *CompAi) {
	c.ai = ai
}

func (c *CAiMemory) GetLocation(tag string) *Location {
	return c.Locations[tag]
}

func (c *CAiMemory) SetLocation(tag string, pos *Location) {
	c.Locations[tag] = pos
	c.SetPosition(tag, pos.Pos)
}

func (c *CAiMemory) SetPosition(tag string, pos vectors.Vec2) {
	c.Positions[tag] = pos
}

func (c *CAiMemory) GetPosition(tag string) vectors.Vec2 {
	return c.Positions[tag]
}

// Update ticks the AI memory by delta.
func (c *CAiMemory) Update(delta float64) {
	// TODO: Update short term memory from perception.
}
