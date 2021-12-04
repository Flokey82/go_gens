package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
)

type Item struct {
	Pos vectors.Vec2
	*ItemType
	// Owned bool
}

type ItemType struct {
	Name       string
	Tags       []string       // Food, Weapon
	Properties map[string]int // Price, weight, damage, ...
}
