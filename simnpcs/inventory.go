package simnpcs

// Inventory represents a storage of items.
type Inventory struct {
	Items []*ItemInstance
}

// newInventory creates a new inventory.
func newInventory() *Inventory {
	return &Inventory{}
}

// List returns a list of all items in the inventory.
func (inv *Inventory) List() []*ItemInstance {
	return inv.Items
}

// Add adds an item to the inventory.
func (inv *Inventory) Add(item *ItemInstance) {
	inv.Items = append(inv.Items, item)
}

// Move moves an item from one inventory to another.
func (inv *Inventory) Move(item *ItemInstance, to *Inventory) {
	to.Add(item)
	inv.Remove(inv.FindIndex(item))
}

// FindIndex returns the index of the given item instance in the inventory.
func (inv *Inventory) FindIndex(item *ItemInstance) int {
	for idx, it := range inv.Items {
		if it == item {
			return idx
		}
	}
	return -1
}

// Find returns the first item instance of a given item in the inventory.
func (inv *Inventory) Find(item *Item) *ItemInstance {
	for _, it := range inv.Items {
		if it.Item == item {
			return it
		}
	}
	return nil
}

// Remove removes an item instance with the given index from the inventory.
func (inv *Inventory) Remove(idx int) {
	inv.Items = append(inv.Items[:idx], inv.Items[idx+1:]...)
}
