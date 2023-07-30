package gamerogueish

import (
	"fmt"
	"math/rand"
)

const (
	ItemWeapon = iota
	ItemPotion
	ItemArmor
	ItemHelmet
	ItemTrigger
	ItemTypeContainer
	ItemTypeDocument
	ItemTypeDecorative
	ItemTypeMax
)

// Item represents an item in the game.
type Item struct {
	*ItemType
	Hidden   bool // indicates if the item is hidden
	Equipped bool // indicates if the item is equipped
	X        int  // x position in the world (if dropped)
	Y        int  // y position in the world (if dropped)
	Contains []*Item
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
	Tile         byte
	Name         string
	Description  string
	Type         int
	Modifier     int
	Hidden       bool // indicates if the item is hidden by default
	Movable      bool // indicates if the item can be moved / picked up
	Rarity       *Rarity
	Capacity     int      // indicates the maximum number of items this item can contain
	PossibleLoot *ItemSet // TODO: Allow multiple sets.
	// TODO: Add variants with different rarity.
	// Variants     []*ItemType
	// TODO: Make this a map of event types to functions.
	OnTouch func(*Game, *Entity, *Item) // Trigger function called when item is used.
	OnUse   func(*Game, *Entity, *Item) // Trigger function called when item is used.
	// OnEquip   func(*Game, *Entity, *Item) // Trigger function called when item is equipped.
	// OnUnequip func(*Game, *Entity, *Item) // Trigger function called when item is unequipped.
	// OnConsume   func(*Game, *Entity, *Item) // Trigger function called when item is consumed.
	// OnDrop      func(*Game, *Entity, *Item) // Trigger function called when item is dropped.
	// OnPickup    func(*Game, *Entity, *Item) // Trigger function called when item is picked up.
	// OnOpen	  func(*Game, *Entity, *Item) // Trigger function called when item is opened.
}

// New returns a new item of the given type.
func (i ItemType) New() *Item {
	var contains []*Item
	if i.PossibleLoot != nil {
		numItems := rand.Intn(1 + i.Capacity) // TODO: improve picking of number of items based on capacity.
		for _, idx := range rand.Perm(len(i.PossibleLoot.Items)) {
			if len(contains) >= numItems {
				break
			}
			if it := i.PossibleLoot.Items[idx]; it.Rarity != nil && it.Rarity.Roll() {
				contains = append(contains, it.New())
			}
		}
	}
	return &Item{
		ItemType: &i,
		Hidden:   i.Hidden,
		Contains: contains,
	}
}

// Generate returns a new item of the given type.
func (i ItemType) Generate() *Item {
	return i.New()
}

var (
	ItemTypeWeaponFishingRod = &ItemType{
		Tile:        'f',
		Name:        "Fishing Rod",
		Description: "Baelin's fishing rod.",
		Type:        ItemWeapon,
		Modifier:    20,
		Rarity:      RarityLegendary,
	}
	ItemTypeWeaponSword = &ItemType{
		Tile:        '/',
		Name:        "Sword",
		Description: "A sharp sword.",
		Type:        ItemWeapon,
		Rarity:      RarityUncommon,
	}
	ItemTypeWeaponAxe = &ItemType{
		Tile:        'P',
		Name:        "Axe",
		Description: "A sharp axe.",
		Type:        ItemWeapon,
		Modifier:    1,
		Rarity:      RarityUncommon,
	}
	ItemTypePotion = &ItemType{
		Tile:        'Ö',
		Name:        "Potion",
		Description: "A healing potion.",
		Type:        ItemPotion,
		Rarity:      RarityUncommon,
	}
	ItemTypeTrollPoop = &ItemType{
		Tile:        '8',
		Name:        "Troll Poop",
		Description: "... with sprinkles!",
		Type:        ItemPotion,
		Modifier:    10,
		Rarity:      RarityLegendary,
	}
	ItemTypeGoblinToe = &ItemType{
		Tile:        't',
		Name:        "Goblin Toe",
		Description: "It's a bit smelly. A popular snack on TikTok.",
		Type:        ItemPotion,
		Modifier:    -1,
		Rarity:      RarityRare,
	}
	ItemTypeArmorLeather = &ItemType{
		Tile:        'L',
		Name:        "Leather Armor",
		Description: "A leather armor.",
		Type:        ItemArmor,
		Rarity:      RarityUncommon,
	}
	ItemTypeArmorChain = &ItemType{
		Tile:        'C',
		Name:        "Chain Armor",
		Description: "A chain armor.",
		Type:        ItemArmor,
		Modifier:    2,
		Rarity:      RarityRare,
	}
	ItemTypeArmorPlate = &ItemType{
		Tile:        'P',
		Name:        "Plate Armor",
		Description: "A plate armor.",
		Type:        ItemArmor,
		Modifier:    4,
		Rarity:      RarityExotic,
	}
	ItemTypeHelmetSweatband = &ItemType{
		Tile:        'S',
		Name:        "Sweatband",
		Description: "A stylish sweatband.",
		Type:        ItemHelmet,
		Modifier:    1,
		Rarity:      RarityUncommon,
	}
	ItemTypeExit = &ItemType{
		Tile:        '>',
		Name:        "Exit",
		Description: "An exit.",
		Type:        ItemTrigger,
		OnTouch: func(g *Game, e *Entity, i *Item) {
			// TODO: Improve the way we check who triggered OnTouch.
			if e.EntityType == EntityPlayer {
				g.setViewMode(ViewModeSuccess)
			}
		},
	}
	ItemTypeTrap = &ItemType{
		Tile:        '^',
		Name:        "Trap",
		Description: "A trap.",
		Hidden:      true,
		Type:        ItemTrigger,
		OnTouch: func(g *Game, e *Entity, i *Item) {
			// TODO: Improve the way we check who triggered OnTouch.
			g.AddMessage(e.Name + " stepped on a trap!")
			e.Health -= 5
			// Mark this trap as revealed.
			i.Hidden = false
			// TODO: Mark this trap as sprung?
		},
	}
	ItemTypeNote = &ItemType{
		Tile:        'N',
		Name:        "Note",
		Description: "A note.",
		Type:        ItemTypeDocument,
		Rarity:      RarityRare,
		OnUse: func(g *Game, e *Entity, i *Item) {
			// TODO: Use text box and load text from item data.
			g.AddMessage("You read the note: \"You are a wizard, Harry!\"")
		},
	}
	ItemTypeLetter = &ItemType{
		Tile:        'L',
		Name:        "Letter",
		Description: "A letter.",
		Type:        ItemTypeDocument,
		Rarity:      RarityLegendary,
		OnUse: func(g *Game, e *Entity, i *Item) {
			// TODO: Use text box and load text from item data.
			g.AddMessage("You read the letter: \"Dear Mr. Bigglesworth!\"")
		},
	}
	ItemTypeBook = &ItemType{
		Tile:         'b',
		Name:         "Book",
		Description:  "A book.",
		Type:         ItemTypeDocument,
		Capacity:     1,
		Rarity:       RarityCommon,
		PossibleLoot: ItemSetBookLoot,
	}
	ItemTypeChest = &ItemType{
		Tile:         'c',
		Name:         "Chest",
		Description:  "A chest.",
		Type:         ItemTypeContainer,
		Capacity:     4,
		Movable:      false,
		PossibleLoot: ItemSetChestLoot,
	}
	ItemTypeBed = &ItemType{
		Tile:         'b',
		Name:         "Bed",
		Description:  "A bed.",
		Type:         ItemTypeContainer,
		Capacity:     2,
		Movable:      false,
		PossibleLoot: ItemSetBedLoot,
	}
	ItemTypeSideTable = &ItemType{
		Tile:        't',
		Name:        "Side Table",
		Description: "A side table.",
		Type:        ItemTypeContainer,
		Capacity:    1,
		Movable:     false,
	}
	ItemTypeAltar = &ItemType{
		Tile:        'a',
		Name:        "Altar",
		Description: "An altar.",
		Type:        ItemTypeContainer,
		Capacity:    1,
		Movable:     false,
	}
	ItemTypeDesk = &ItemType{
		Tile:         'd',
		Name:         "Desk",
		Description:  "A desk.",
		Type:         ItemTypeContainer,
		Capacity:     2,
		Movable:      false,
		PossibleLoot: ItemSetStationary,
	}
	ItemTypeBookshelf = &ItemType{
		Tile:         'b',
		Name:         "Book Shelf",
		Description:  "A book shelf.",
		Type:         ItemTypeContainer,
		Capacity:     4,
		Movable:      false,
		Rarity:       RarityCommon,
		PossibleLoot: ItemSetBookshelfLoot,
	}
	ItemTypeCandle = &ItemType{
		Tile:        'c',
		Name:        "Candle",
		Description: "A candle.",
		Type:        ItemTypeDecorative,
		Movable:     true,
	}
	ItemTypeQuill = &ItemType{
		Tile:        'q',
		Name:        "Quill",
		Description: "A quill.",
		Type:        ItemTypeDecorative,
		Rarity:      RarityUncommon,
		Movable:     true,
	}
	ItemTypeInk = &ItemType{
		Tile:        'i',
		Name:        "Ink",
		Description: "An ink.",
		Type:        ItemTypeDecorative,
		Rarity:      RarityUncommon,
		Movable:     true,
	}
	ItemTypeParchment = &ItemType{
		Tile:        'p',
		Name:        "Parchment",
		Description: "A parchment.",
		Type:        ItemTypeDecorative,
		Rarity:      RarityUncommon,
		Movable:     true,
	}
	ItemTypeSealingWax = &ItemType{
		Tile:        'w',
		Name:        "Sealing Wax",
		Description: "A sealing wax.",
		Type:        ItemTypeDecorative,
		Rarity:      RarityUncommon,
		Movable:     true,
	}
	ItemTypeEnvelope = &ItemType{
		Tile:         'e',
		Name:         "Envelope",
		Description:  "An envelope.",
		Type:         ItemTypeDecorative,
		Rarity:       RarityUncommon,
		Capacity:     1,
		Movable:      true,
		PossibleLoot: ItemSetBookLoot,
	}
	ItemTypeBedSheet = &ItemType{
		Tile:        's',
		Name:        "Bed Sheet",
		Description: "A bed sheet.",
		Type:        ItemTypeDecorative,
		Rarity:      RarityCommon,
		Movable:     true,
	}
	ItemTypePillow = &ItemType{
		Tile:        'p',
		Name:        "Pillow",
		Description: "A pillow.",
		Type:        ItemTypeDecorative,
		Rarity:      RarityUncommon,
		Movable:     true,
	}
	ItemTypeBlanket = &ItemType{
		Tile:        'b',
		Name:        "Blanket",
		Description: "A blanket.",
		Type:        ItemTypeDecorative,
		Rarity:      RarityCommon,
		Movable:     true,
	}
)
