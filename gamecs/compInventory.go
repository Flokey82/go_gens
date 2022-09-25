package gamecs

// CompInventory represents a storage for items with a specified
// capacity.
type CompInventory struct {
	id    int
	w     *World
	Slots []*Item
	Size  int
}

// newCompInventory returns a new CompInventory instance.
func newCompInventory(w *World, id, size int) *CompInventory {
	return &CompInventory{
		w:    w,
		id:   id,
		Size: size,
	}
}

// IsFull returns true if we have exceeded or have reached
// the inventory capacity
func (in *CompInventory) IsFull() bool {
	return len(in.Slots) >= in.Size
}

// Add adds the given item to the inventory and return success.
// NOTE: This will fail if the inventory is full.
func (in *CompInventory) Add(it *Item) bool {
	if in.IsFull() {
		return false
	}
	it.Location = LocInventory
	it.LocationID = in.id
	in.Slots = append(in.Slots, it)
	return true
}

// Has returns true if an item with the given ItemType is present in
// the inventory.
// TODO: Maybe return count instead of bool or introduce a second
// method that returns the count.
func (in *CompInventory) Has(itt *ItemType) bool {
	for _, it := range in.Slots {
		if it.ItemType == itt {
			return true
		}
	}
	return false
}

// Find finds an item with a given tag in the inventory and
// returns the found item.
func (in *CompInventory) Find(tag string) *Item {
	for _, it := range in.Slots {
		for _, t := range it.Tags {
			if t == tag {
				return it
			}
		}
	}
	return nil
}

// Drop drops the specified item from inventory and returns
// true on success.
func (in *CompInventory) Drop(it *Item) bool {
	// TODO: Many nil pointer checks needed here...
	if in.RemoveID(it.id) {
		it.Location = LocWorld
		it.LocationID = -1
		it.Pos = in.w.mgr.GetEntityFromID(in.id).Pos
		return true
	}
	return false
}

// RemoveID removes the item with the given id from inventory
// and returns true on success.
func (in *CompInventory) RemoveID(id int) bool {
	for i, it := range in.Slots {
		if it.id == id {
			in.Slots = append(in.Slots[:i], in.Slots[i+1:]...)
			return true
		}
	}
	return false
}

// TransferAll transfers all items from the current inventory
// to the target inventory.
func (in *CompInventory) TransferAll(to *CompInventory) bool {
	for len(in.Slots) > 0 {
		it := in.Slots[0]
		if !to.Add(it) { // Out of storage?
			return false
		}
		in.RemoveID(it.id)
	}
	return true
}
