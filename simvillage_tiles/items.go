package simvillage_tiles

// Tile IDs of items.
const (
	TileItemRock = 141
)

// ItemType represents an type of an item.
type ItemType struct {
	Name string
	Tile int
}

// NewInstance returns a new item with the given type.
func (i *ItemType) NewInstance() *Item {
	return &Item{
		ItemType: i,
	}
}

// Item represents an item in the game.
//
// This item has a base type but has additional
// properties like a position, etc.
type Item struct {
	*ItemType
	Position [2]int
}

var (
	ItemTypeRock = &ItemType{
		Name: "Rock",
		Tile: TileItemRock,
	}
)
