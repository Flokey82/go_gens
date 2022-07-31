package simnpcs

type Inventory struct {
	Items []*ItemInstance
}

func newInventory() *Inventory {
	return &Inventory{}
}

func (inv *Inventory) List() []*ItemInstance {
	return inv.Items
}

func (inv *Inventory) Add(item *ItemInstance) {
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) Move(item *ItemInstance, to *Inventory) {
	to.Add(item)
	inv.Remove(inv.FindIndex(item))
}

func (inv *Inventory) FindIndex(item *ItemInstance) int {
	for i, it := range inv.Items {
		if it == item {
			return i
		}
	}
	return -1
}

func (inv *Inventory) Find(item *Item) *ItemInstance {
	for _, it := range inv.Items {
		if it.Item == item {
			return it
		}
	}
	return nil
}

func (inv *Inventory) Remove(s int) {
	inv.Items = append(inv.Items[:s], inv.Items[s+1:]...)
}
