package genfurnishing

import "math/rand"

// ItemBase represents a base item prototype.
type ItemBase struct {
	Name     string
	Rarity   *Rarity
	Capacity int
	Size     Size
	Type     string
	Variants []string
	Set      *ItemSet // TODO: Allow multiple sets.
}

// NewItemBase returns a new item base.
func NewItemBase(name, iType string, setFunc func(*ItemBase)) *ItemBase {
	item := &ItemBase{
		Name:   name,
		Size:   SizeMedium,
		Rarity: RarityAverage,
		Type:   iType,
	}
	if setFunc != nil {
		setFunc(item)
	}
	return item
}

// Generate returns a new item based on this item base.
func (it *ItemBase) Generate() *Item {
	var variant string
	if len(it.Variants) > 0 {
		variant = it.Variants[rand.Intn(len(it.Variants))]
	}
	// Generate a number of items from the set up to the capacity.
	var contains []*Item
	if it.Capacity > 0 && it.Set != nil {
		contains = append(contains, it.Set.Generate())
	}
	return &Item{
		ItemBase: it,
		Variant:  variant,
		Contains: contains,
	}
}

// Item represents a generated item.
type Item struct {
	*ItemBase
	Variant  string
	Contains []*Item
}

func (it *Item) Name() string {
	if it.Variant != "" {
		return it.ItemBase.Name + " (" + it.Variant + ")"
	}
	return it.ItemBase.Name
}

// Size is the size of an item.
type Size string

const (
	SizeTiny    Size = "tiny"
	SizeSmall   Size = "small"
	SizeMedium  Size = "medium"
	SizeLarge   Size = "large"
	SizeMassive Size = "massive"
)

// Storage capacities.
const (
	CapacitySmall  = 1
	CapacityMedium = 2
	CapacityLarge  = 3
)

// Rarity represents the rarity of an item.
type Rarity struct {
	Name           string // Name of this rarity
	Probability    int    // Probability of this rarity (the higher the more rare)
	IndicateRarity bool   // Indicate rarity in item name
}

// Roll returns true if the item should be generated.
func (r *Rarity) Roll() bool {
	return rand.Intn(101) >= r.Probability
}

var (
	RarityAbundant = &Rarity{
		Name:           "abundant",
		Probability:    25,
		IndicateRarity: false,
	}
	RarityCommon = &Rarity{
		Name:           "common",
		Probability:    45,
		IndicateRarity: false,
	}
	RarityAverage = &Rarity{
		Name:           "average",
		Probability:    65,
		IndicateRarity: false,
	}
	RarityUncommon = &Rarity{
		Name:           "uncommon",
		Probability:    80,
		IndicateRarity: true,
	}
	RarityRare = &Rarity{
		Name:           "rare",
		Probability:    93,
		IndicateRarity: true,
	}
	RarityExotic = &Rarity{
		Name:           "exotic",
		Probability:    99,
		IndicateRarity: true,
	}
	RarityLegendary = &Rarity{
		Name:           "legendary",
		Probability:    100,
		IndicateRarity: true,
	}
)

const (
	ItemTypeAmmo          = "ammo"
	ItemTypeArmor         = "armor"
	ItemTypeChancery      = "chancery"
	ItemTypeClothing      = "clothing"
	ItemTypeCoin          = "coin"
	ItemTypeContainer     = "container"
	ItemTypeFood          = "food"
	ItemTypeFurnishing    = "furnishing"
	ItemTypeKitchen       = "kitchen"
	ItemTypeLiquid        = "liquid"
	ItemTypeMiscellaneous = "miscellaneous"
	ItemTypeMysterious    = "mysterious"
	ItemTypeMystic        = "mystic"
	ItemTypePotion        = "potion"
	ItemTypeSurvival      = "survival"
	ItemTypeStationary    = "stationary"
	ItemTypeTack          = "tack"
	ItemTypeTool          = "tool"
	ItemTypeTreasure      = "treasure"
	ItemTypeTrinket       = "trinket"
	ItemTypeWeapon        = "weapon"
)

type ItemSet struct {
	Name  string
	Items []*ItemBase
}

func (is *ItemSet) Generate() *Item {
	if len(is.Items) == 0 {
		return nil
	}
	for _, i := range rand.Perm(len(is.Items)) {
		item := is.Items[i]
		if item.Rarity.Roll() {
			return item.Generate()
		}
	}
	return is.Items[rand.Intn(len(is.Items))].Generate()
}

// Clothing
const (
	ClothingNameCloak  = "cloak"
	ClothingNameCoat   = "coat"
	ClothingNameGloves = "gloves"
	ClothingNameHat    = "hat"
	ClothingNamePants  = "pants"
	ClothingNameShirt  = "shirt"
)

var ClothingSet = &ItemSet{
	Name: ItemTypeClothing,
	Items: []*ItemBase{
		NewItemBase(ClothingNameCloak, ItemTypeClothing, func(i *ItemBase) {
			i.Rarity = RarityUncommon
			i.Size = SizeLarge
		}),
		NewItemBase(ClothingNameCoat, ItemTypeClothing, func(i *ItemBase) {
			i.Rarity = RarityUncommon
			i.Size = SizeLarge
		}),
		NewItemBase(ClothingNameGloves, ItemTypeClothing, func(i *ItemBase) {
			i.Rarity = RarityCommon
			i.Size = SizeSmall
		}),
		NewItemBase(ClothingNameHat, ItemTypeClothing, func(i *ItemBase) {
			i.Rarity = RarityCommon
			i.Size = SizeSmall
		}),
		NewItemBase(ClothingNamePants, ItemTypeClothing, func(i *ItemBase) {
			i.Rarity = RarityCommon
			i.Size = SizeMedium
		}),
		NewItemBase(ClothingNameShirt, ItemTypeClothing, func(i *ItemBase) {
			i.Rarity = RarityCommon
			i.Size = SizeMedium
		}),
	},
}

// Stationary
const (
	StationaryNameLedger = "ledger"
	StationaryNameQuill  = "quill"
	StationaryNameScroll = "scroll"
)

var StationarySet = &ItemSet{
	Name: "stationary",
	Items: []*ItemBase{
		NewItemBase(StationaryNameLedger, ItemTypeStationary, func(i *ItemBase) {
			i.Rarity = RarityCommon
			i.Size = SizeMedium
		}),
		NewItemBase(StationaryNameQuill, ItemTypeStationary, func(i *ItemBase) {
			i.Rarity = RarityCommon
			i.Size = SizeTiny
		}),
		NewItemBase(StationaryNameScroll, ItemTypeStationary, func(i *ItemBase) {
			i.Rarity = RarityCommon
			i.Size = SizeTiny
		}),
	},
}
