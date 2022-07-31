package main

import (
	"github.com/Flokey82/aifiver"
	"github.com/Flokey82/go_gens/simnpcs"
)

func main() {
	idx := simnpcs.New()

	// Setup items.
	itemResourceIronIngot := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "iron ingot",
		Buy:          6,
		Sell:         2,
		Durability:   100,
		RequiresTime: 2,
	}

	itemResourceLeather := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "leather",
		Buy:          6,
		Sell:         2,
		Durability:   100,
		RequiresTime: 10,
	}

	itemResourceWood := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "wood",
		Buy:          2,
		Sell:         1,
		Durability:   100,
		RequiresTime: 10,
	}

	itemLeatherArmor := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "leather armor",
		Buy:          20,
		Sell:         10,
		Durability:   100,
		RequiresTime: 20,
		Consumes:     []*simnpcs.Item{itemResourceLeather},
	}

	itemIronSword := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "basic sword",
		Buy:          10,
		Sell:         5,
		Durability:   50,
		RequiresTime: 10,
		Consumes:     []*simnpcs.Item{itemResourceIronIngot},
	}

	itemIronPickaxe := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "basic pickaxe",
		Buy:          6,
		Sell:         3,
		Durability:   50,
		RequiresTime: 4,
		Consumes:     []*simnpcs.Item{itemResourceIronIngot, itemResourceWood},
	}

	itemFoodWheat := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "Wheat",
		Buy:          5,
		Sell:         1,
		Durability:   100,
		RequiresTime: 20,
	}

	itemFoodBread := &simnpcs.Item{
		ID:           idx.GetID(),
		Name:         "Bread",
		Buy:          2,
		Sell:         1,
		Durability:   100,
		RequiresTime: 5,
		Consumes:     []*simnpcs.Item{itemFoodWheat},
	}

	// Setup blacksmith profession.
	profBlacksmith := idx.NewProfession("Blacksmith", simnpcs.LocTypeSmith)
	profBlacksmith.AddSkill(idx.GetID(), "basic smithing",
		[]*simnpcs.Item{itemLeatherArmor, itemIronSword, itemIronPickaxe}, 0)

	// Setup shopkeeper profession.
	profShopkeeper := idx.NewProfession("Shopkeeper", simnpcs.LocTypeShop)
	profShopkeeper.AddSkill(idx.GetID(), "basic baking", []*simnpcs.Item{itemFoodBread}, 0)

	// Setup farmer profession.
	profFarmer := idx.NewProfession("Farmer", simnpcs.LocTypeFarm)
	profFarmer.AddSkill(idx.GetID(), "basic farming", []*simnpcs.Item{itemFoodWheat, itemResourceLeather, itemResourceWood}, 0)

	// Setup miner profession.
	profMiner := idx.NewProfession("Miner", simnpcs.LocTypeMine)
	profMiner.AddSkill(idx.GetID(), "basic mining", []*simnpcs.Item{itemResourceIronIngot}, 0)

	// Setup town.
	locTown := idx.NewLocation("Hometownville", nil, simnpcs.LocTypeTown, simnpcs.LocScaleTown)

	// Setup general store.
	locShop := idx.NewLocation("General Store", locTown, simnpcs.LocTypeShop, simnpcs.LocScaleMerchantShopTavern)

	// Setup inn.
	locInn := idx.NewLocation("Local Inn", locTown, simnpcs.LocTypeInn, simnpcs.LocScaleMerchantShopTavern)

	// Setup farm.
	locFarm := idx.NewLocation("Local Farm", locTown, simnpcs.LocTypeFarm, simnpcs.LocScaleMerchantShopTavern)

	// Setup blacksmith.
	locBlacksmith := idx.NewLocation("Blacksmith", locTown, simnpcs.LocTypeSmith, simnpcs.LocScaleMerchantShopTavern)

	// Setup fancy townhouse.
	locTownhouse := idx.NewLocation("Fancy Townhouse", locTown, simnpcs.LocTypeHome, simnpcs.LocScaleHouse)

	// Setup local mine.
	locMine := idx.NewLocation("Scary Mine", locTown, simnpcs.LocTypeMine, simnpcs.LocScaleMerchantShopTavern)

	// Setup shopkeeper.
	charShopKeep := simnpcs.NewCharacter(idx.GetID(), "ShopKeep", "LastName",
		aifiver.SmallModel{6, 4, 3, 4, 0})
	idx.StartCareer(charShopKeep, profShopkeeper, locShop)
	idx.Entries = append(idx.Entries, charShopKeep)
	charShopKeep.Home = locShop

	// Setup hard working farmer.

	charFarmer := simnpcs.NewCharacter(idx.GetID(), "Farmer", "Fritz",
		aifiver.SmallModel{2, 5, 1, 2, 0})
	idx.StartCareer(charFarmer, profFarmer, locFarm)
	idx.Entries = append(idx.Entries, charFarmer)
	charFarmer.Home = locFarm

	// Setup friendly blacksmith.
	charBlacksmith := simnpcs.NewCharacter(idx.GetID(), "Frieda", "Friendly",
		aifiver.SmallModel{3, 1, 3, 4, 0})
	idx.StartCareer(charBlacksmith, profBlacksmith, locBlacksmith)
	idx.Entries = append(idx.Entries, charBlacksmith)
	charBlacksmith.AddRoutine(&simnpcs.Routine{
		DayOfWeek: 2,
		Hour:      12,
		C:         charBlacksmith,
		Location:  locShop,
	})
	charBlacksmith.AddRoutine(&simnpcs.Routine{
		DayOfWeek: 3,
		Hour:      13,
		C:         charBlacksmith,
		Location:  locFarm,
	})
	charBlacksmith.Home = locInn

	// Setup unfriendly person.
	charBaron := simnpcs.NewCharacter(idx.GetID(), "Gunther", "Grouchy",
		aifiver.SmallModel{-2, 0, -1, -5, -1})
	charBaron.Title = "Baron"
	idx.Entries = append(idx.Entries, charBaron)
	charBaron.AddRoutine(&simnpcs.Routine{
		DayOfWeek: 2,
		Hour:      12,
		C:         charBaron,
		Location:  locShop,
	})
	charBaron.Home = locTownhouse

	// Setup miner.
	charMiner := simnpcs.NewCharacter(idx.GetID(), "Miner", "Mike",
		aifiver.SmallModel{2, 0, 1, 5, 1})
	idx.StartCareer(charMiner, profMiner, locMine)
	idx.Entries = append(idx.Entries, charMiner)
	charMiner.AddRoutine(&simnpcs.Routine{
		DayOfWeek: 2,
		Hour:      12,
		C:         charMiner,
		Location:  locShop,
	})
	charMiner.Home = locInn
	for i := 0; i < 2000; i++ {
		idx.Tick()
	}
}
