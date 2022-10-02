package gamerogueish

type Entity struct {
	X         int  // x position in the world
	Y         int  // y position in the world
	Tile      byte // tile character
	Inventory      // inventory component

	// This should be part of the entity.
	Slots [ItemTypeMax]*Item
}

// NewEntity returns a new entity with the given position and tile.
func NewEntity(x, y int, tile byte) *Entity {
	return &Entity{
		X:    x,
		Y:    y,
		Tile: tile,
	}
}

// Equip equips the item at the given inventory index.
func (e *Entity) Equip(index int) {
	if index < 0 || index >= len(e.Items) || !e.Items[index].Equippable() {
		return
	}

	// Toggle equipped state.
	it := e.Items[index]
	it.Equipped = !it.Equipped

	// If we unequip the item, unset the equipped item.
	if !it.Equipped {
		e.Slots[it.Type] = nil
		return
	}

	// If there is already an item in the slot, unequip it.
	if e.Slots[it.Type] != nil {
		e.Slots[it.Type].Equipped = false
	}
	e.Slots[it.Type] = it
}

// Consume consumes the item at the given inventory index.
func (e *Entity) Consume(index int) {
	if index < 0 || index >= len(e.Items) || !e.Items[index].Consumable() {
		return
	}
	e.RemoveItem(index)
}
