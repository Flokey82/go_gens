package gamerogueish

import "math/rand"

type ItemSet struct {
	Name  string
	Items []*ItemType
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
		if item.Rarity == nil || item.Rarity.Roll() {
			return item.Generate()
		}
	}
	return is.Items[rand.Intn(len(is.Items))].Generate()
}

var (
	ItemSetWeapons = &ItemSet{
		Name: "weapons",
		Items: []*ItemType{
			ItemTypeWeaponFishingRod,
			ItemTypeWeaponSword,
			ItemTypeWeaponAxe,
		},
	}
	ItemSetArmor = &ItemSet{
		Name: "armor",
		Items: []*ItemType{
			ItemTypeArmorLeather,
			ItemTypeArmorChain,
			ItemTypeArmorPlate,
			ItemTypeHelmetSweatband,
		},
	}
	ItemSetPotions = &ItemSet{
		Name: "potions",
		Items: []*ItemType{
			ItemTypePotion,
			ItemTypeTrollPoop,
			ItemTypeGoblinToe,
		},
	}
	ItemSetChestLoot = &ItemSet{
		Name: "chest loot",
		Items: []*ItemType{
			ItemTypeWeaponFishingRod,
			ItemTypeWeaponSword,
			ItemTypeWeaponAxe,
			ItemTypePotion,
			ItemTypeArmorLeather,
			ItemTypeArmorChain,
			ItemTypeArmorPlate,
			ItemTypeHelmetSweatband,
			ItemTypeExit,
			ItemTypeTrap,
			ItemTypeNote,
		},
	}
	ItemSetBookLoot = &ItemSet{
		Name: "book loot",
		Items: []*ItemType{
			ItemTypeNote,
			ItemTypeLetter,
		},
	}
	ItemSetBookshelfLoot = &ItemSet{
		Name: "bookshelf loot",
		Items: []*ItemType{
			ItemTypeBook,
			ItemTypeNote,
			ItemTypeLetter,
			ItemTypeBook,
			ItemTypeBook,
			ItemTypeBook,
		},
	}
	ItemSetStationary = &ItemSet{
		Name: "stationary",
		Items: []*ItemType{
			ItemTypeBook,
			ItemTypeNote,
			ItemTypeLetter,
			ItemTypeQuill,
			ItemTypeInk,
			ItemTypeParchment,
			ItemTypeSealingWax,
		},
	}
	ItemSetBedLoot = &ItemSet{
		Name: "bed loot",
		Items: []*ItemType{
			ItemTypeBedSheet,
			ItemTypePillow,
			ItemTypeBlanket,
		},
	}
)
