package simnpcs

// TODO: Allow passive production (growing crops).
type Item struct {
	ID           uint64
	Name         string
	Buy          int     // Purchase price
	Sell         int     // Sale price
	Durability   int     // Base durability in cycles
	Consumes     []*Item // Items consumed to produce
	RequiresTime int     // Number of cycles required to produce
	Requires     *Skill  // Skill required to produce
	ProducedBy   []*Profession
}

func (item *Item) newInstance(id uint64) *ItemInstance {
	return &ItemInstance{
		ID:         id,
		Item:       item,
		Durability: item.Durability,
	}
}

type ItemInstance struct {
	ID         uint64 // Unique id
	Item       *Item  // Base item
	Durability int    // Remaining durability
}
