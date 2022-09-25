package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
)

// ItemLocation indicates the type of location where an item
// is located.
type ItemLocation int

const (
	LocWorld ItemLocation = iota
	LocContainer
	LocInventory
)

// Item represents a movable item in the world.
type Item struct {
	id         int          // ID of the item
	Location   ItemLocation // Type of location (world, inventory, ...)
	LocationID int          // ID of the entity if in inventory
	Pos        vectors.Vec2 // World position if in world
	*ItemType               // Base information
}

// ID returns the unique identifier for this item.
func (it *Item) ID() int {
	return it.id
}

// ItemType represents the base type of an item.
type ItemType struct {
	Name       string
	Tags       []string       // Food, Weapon
	Properties map[string]int // Price, weight, damage, ...
	Requires   []*ItemType    // Requires items to craft
}

// NewItemType returns a new item type with the given name and tags.
func NewItemType(name string, tags ...string) *ItemType {
	return &ItemType{
		Name:       name,
		Properties: make(map[string]int),
		Tags:       tags,
	}
}

// New returns a new item of the current type.
func (i *ItemType) New(w *World, pos vectors.Vec2) *Item {
	return &Item{
		ItemType: i,
		Pos:      pos,
		id:       w.mgr.NextID(),
	}
}
