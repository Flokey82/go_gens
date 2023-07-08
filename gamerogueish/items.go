package gamerogueish

import "fmt"

const (
	ItemWeapon = iota
	ItemPotion
	ItemArmor
	ItemHelmet
	ItemTrigger
	ItemTypeMax
)

// Item represents an item in the game.
type Item struct {
	*ItemType
	Equipped bool // indicates if the item is equipped
	X        int  // x position in the world (if dropped)
	Y        int  // y position in the world (if dropped)
}

// String returns the name of the item.
func (i Item) String() string {
	return i.Name
}

// FullName returns the full name of the item including the modifier.
func (i Item) FullName() string {
	if i.Modifier == 0 {
		return i.Name
	}
	if i.Modifier > 0 {
		return fmt.Sprintf("%s (%+d)", i.Name, i.Modifier)
	}
	return fmt.Sprintf("%s (%d)", i.Name, i.Modifier)
}

// Equippable returns true if the item can be equipped.
func (i Item) Equippable() bool {
	return i.Type == ItemWeapon || i.Type == ItemArmor || i.Type == ItemHelmet
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

// Remove the given item from the inventory and return it.
func (i *Inventory) Remove(item *Item) *Item {
	for j, v := range i.Items {
		if v == item {
			return i.RemoveItem(j)
		}
	}
	return nil
}

// RemoveItem removes the item with the given index from the
// inventory and returns it.
func (i *Inventory) RemoveItem(index int) *Item {
	if index < 0 || index >= len(i.Items) {
		return nil
	}
	item := i.Items[index]
	i.Items = append(i.Items[:index], i.Items[index+1:]...)
	return item
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
	Tile        byte
	Name        string
	Description string
	Type        int
	Modifier    int
	OnTouch     func(*Game, *Entity, *Item) // Trigger function called when item is used.
	// OnEquip    func(*Game, *Entity, *Item) // Trigger function called when item is equipped.
	// OnUnequip  func(*Game, *Entity, *Item) // Trigger function called when item is unequipped.
}

// New returns a new item of the given type.
func (i ItemType) New() *Item {
	return &Item{
		ItemType: &i,
	}
}

var (
	ItemTypeWeaponFishingRod = &ItemType{
		Tile:        'f',
		Name:        "Fishing Rod",
		Description: "Baelin's fishing rod.",
		Type:        ItemWeapon,
		Modifier:    20,
	}
	ItemTypeWeaponSword = &ItemType{
		Tile:        '/',
		Name:        "Sword",
		Description: "A sharp sword.",
		Type:        ItemWeapon,
	}
	ItemTypeWeaponAxe = &ItemType{
		Tile:        'P',
		Name:        "Axe",
		Description: "A sharp axe.",
		Type:        ItemWeapon,
		Modifier:    1,
	}
	ItemTypePotion = &ItemType{
		Tile:        'Ö',
		Name:        "Potion",
		Description: "A healing potion.",
		Type:        ItemPotion,
	}
	ItemTypeTrollPoop = &ItemType{
		Tile:        '8',
		Name:        "Troll Poop",
		Description: "A troll poop.",
		Type:        ItemPotion,
		Modifier:    10,
	}
	ItemTypeArmorLeather = &ItemType{
		Tile:        'L',
		Name:        "Leather Armor",
		Description: "A leather armor.",
		Type:        ItemArmor,
	}
	ItemTypeArmorChain = &ItemType{
		Tile:        'C',
		Name:        "Chain Armor",
		Description: "A chain armor.",
		Type:        ItemArmor,
		Modifier:    2,
	}
	ItemTypeArmorPlate = &ItemType{
		Tile:        'P',
		Name:        "Plate Armor",
		Description: "A plate armor.",
		Type:        ItemArmor,
		Modifier:    4,
	}
	ItemTypeHelmetSweatband = &ItemType{
		Tile:        'S',
		Name:        "Sweatband",
		Description: "A stylish sweatband.",
		Type:        ItemHelmet,
		Modifier:    1,
	}
	ItemTypeExit = &ItemType{
		Tile:        '>',
		Name:        "Exit",
		Description: "An exit.",
		Type:        ItemTrigger,
		OnTouch: func(g *Game, e *Entity, i *Item) {
			g.setViewMode(ViewModeSuccess)
		},
	}
	ItemTypeTrap = &ItemType{
		Tile:        ' ',
		Name:        "Trap",
		Description: "A trap.",
		Type:        ItemTrigger,
		OnTouch: func(g *Game, e *Entity, i *Item) {
			// TODO: Mark this trap as sprung.
			g.AddMessage("You stepped on a trap!")
			e.Health -= 5
			if e.Health <= 0 {
				g.setViewMode(ViewModeDeath)
			}
		},
	}
)
