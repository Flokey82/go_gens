package gamerogueish

const (
	ItemWeapon = iota
	ItemPotion
	ItemArmor
	ItemTypeMax
)

// Item represents an item in the game.
type Item struct {
	*ItemType
	Equipped bool // indicates if the item is equipped
}

// String returns the name of the item.
func (i Item) String() string {
	return i.Name
}

// Equippable returns true if the item can be equipped.
func (i Item) Equippable() bool {
	return i.Type == ItemWeapon || i.Type == ItemArmor
}

// Consumable returns true if the item can be consumed.
func (i Item) Consumable() bool {
	return i.Type == ItemPotion
}

// Inventory represents a collection of items.
type Inventory struct {
	Items []*Item
}

// Add the given item to the inventory.
func (i *Inventory) Add(item *Item) {
	i.Items = append(i.Items, item)
}

// Remove the given item from the inventory.
func (i *Inventory) Remove(item *Item) {
	for j, v := range i.Items {
		if v == item {
			i.RemoveItem(j)
			break
		}
	}
}

// RemoveItem removes the item with the given index from the
// inventory.
func (i *Inventory) RemoveItem(index int) {
	if index < 0 || index >= len(i.Items) {
		return
	}
	i.Items = append(i.Items[:index], i.Items[index+1:]...)
}

// GetItem returns the item with the given index.
func (i *Inventory) GetItem(index int) *Item {
	if index < 0 || index >= len(i.Items) {
		return nil
	}
	return i.Items[index]
}

// Count returns the number of items in the inventory.
func (i *Inventory) Count() int {
	return len(i.Items)
}

// ItemType represents a type of item.
type ItemType struct {
	Name        string
	Description string
	Type        int
}

// New returns a new item of the given type.
func (i ItemType) New() *Item {
	return &Item{
		ItemType: &i,
	}
}

var (
	ItemTypeWeaponSword = &ItemType{
		Name:        "Sword",
		Description: "A sharp sword.",
		Type:        ItemWeapon,
	}
	ItemTypeWeaponAxe = &ItemType{
		Name:        "Axe",
		Description: "A sharp axe.",
		Type:        ItemWeapon,
	}
	ItemTypePotion = &ItemType{
		Name:        "Potion",
		Description: "A healing potion.",
		Type:        ItemPotion,
	}
	ItemTypeArmorLeather = &ItemType{
		Name:        "Leather Armor",
		Description: "A leather armor.",
		Type:        ItemArmor,
	}
	ItemTypeArmorChain = &ItemType{
		Name:        "Chain Armor",
		Description: "A chain armor.",
		Type:        ItemArmor,
	}
	ItemTypeArmorPlate = &ItemType{
		Name:        "Plate Armor",
		Description: "A plate armor.",
		Type:        ItemArmor,
	}
	ItemTypeTrollPoop = &ItemType{
		Name:        "Troll Poop",
		Description: "A troll poop.",
		Type:        ItemPotion,
	}
)
