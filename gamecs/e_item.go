package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
)

type ItemLocation int

const (
	LocWorld ItemLocation = iota
	LocContainer
	LocInventory
)

type Item struct {
	id         int          // ID of the item
	Location   ItemLocation // Type of location (world, inventory, ...)
	LocationID int          // ID of the entity if in inventory
	Pos        vectors.Vec2 // World position if in world
	*ItemType               // Base information
}

func (it *Item) ID() int {
	return it.id
}

type ItemType struct {
	Name       string
	Tags       []string       // Food, Weapon
	Properties map[string]int // Price, weight, damage, ...
	Requires   []*ItemType    // Requires items to craft
}

func NewItemType(name string, tags ...string) *ItemType {
	return &ItemType{
		Name:       name,
		Properties: make(map[string]int),
		Tags:       tags,
	}
}

func (i *ItemType) New(w *World, pos vectors.Vec2) *Item {
	return &Item{
		ItemType: i,
		Pos:      pos,
		id:       w.mgr.NextID(),
	}
}
