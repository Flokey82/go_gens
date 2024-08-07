package simnpcs

// Item represents an item that can be bought, sold, crafted, etc.
// TODO: Allow passive production (growing crops).
type Item struct {
	ID           uint64        // Unique id
	Name         string        // Name of the item
	Buy          int           // Purchase price
	Sell         int           // Sale price
	Durability   int           // Base durability in cycles
	Consumes     []*Item       // Items consumed to produce
	RequiresTime int           // Number of cycles required to produce
	Requires     *Skill        // Skill required to produce
	ProducedBy   []*Profession // Professions that can produce this item
}

// newInstance creates a new instance of the item.
func (item *Item) newInstance(id uint64) *ItemInstance {
	return &ItemInstance{
		ID:         id,
		Item:       item,
		Durability: item.Durability,
	}
}

// ItemInstance represents an instance of an item.
type ItemInstance struct {
	ID         uint64 // Unique id
	Item       *Item  // Base item
	Durability int    // Remaining durability
}
