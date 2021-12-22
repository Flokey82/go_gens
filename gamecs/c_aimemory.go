package gamecs

import "github.com/Flokey82/go_gens/vectors"

type CAiMemory struct {
	ai        *CAi
	Locations map[string]vectors.Vec2
}

func newCAiMemory() *CAiMemory {
	return &CAiMemory{
		Locations: make(map[string]vectors.Vec2),
	}
}

func (c *CAiMemory) init(ai *CAi) {
	c.ai = ai
}

func (c *CAiMemory) GetLocation(tag string) vectors.Vec2 {
	return c.Locations[tag]
}

func (c *CAiMemory) SetLocation(tag string, pos vectors.Vec2) {
	c.Locations[tag] = pos
}

// Update ticks the AI memory by delta.
func (c *CAiMemory) Update(delta float64) {
	// TODO: Update short term memory from perception.
}
