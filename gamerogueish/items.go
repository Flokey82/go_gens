package gamerogueish

const (
	ItemWeapon = iota
	ItemPotion
	ItemArmor
	ItemTypeMax
)

type Item struct {
	Name        string
	Description string
	Type        int
	Equipped    bool
}

func (i Item) String() string {
	return i.Name
}

func (i Item) Equippable() bool {
	return i.Type == ItemWeapon || i.Type == ItemArmor
}

func (i Item) Consumable() bool {
	return i.Type == ItemPotion
}

type Inventory struct {
	Items []Item

	// This should be part of a UI component.
	selectedItem int
}

func (i *Inventory) SelectItem(index int) {
	if index < 0 {
		index = 0
	} else if index >= len(i.Items) {
		index = len(i.Items) - 1
	}
	i.selectedItem = index
}

func (i *Inventory) Add(item Item) {
	i.Items = append(i.Items, item)
}

func (i *Inventory) Remove(item Item) {
	for j, v := range i.Items {
		if v == item {
			i.Items = append(i.Items[:j], i.Items[j+1:]...)
			break
		}
	}
}

func (i *Inventory) Get(index int) Item {
	return i.Items[index]
}

func (i *Inventory) Count() int {
	return len(i.Items)
}
