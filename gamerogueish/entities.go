package gamerogueish

type Entity struct {
	X    int
	Y    int
	Tile byte
	Inventory

	// This should be part of the entity.
	Slots [ItemTypeMax]*Item
}

func NewEntity(x, y int, tile byte) *Entity {
	return &Entity{
		X:    x,
		Y:    y,
		Tile: tile,
	}
}

func (e *Entity) Equip(index int) {
	if index < 0 || index >= len(e.Items) || !e.Items[index].Equippable() {
		return
	}

	// Toggle equipped state.
	it := &e.Items[index]
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
