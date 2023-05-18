package genfurnishing

import "math/rand"

type ItemSet struct {
	Name  string
	Items []*ItemBase
}

func (is *ItemSet) GenerateN(n int) []*Item {
	if len(is.Items) == 0 {
		return nil
	}
	var res []*Item
	for len(res) < n {
		res = append(res, is.Generate())
	}
	return res
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
