package genfurnishing

const (
	RoomFurnishingAlchemy   = "alchemy"
	RoomFurnishingAnvil     = "anvil"
	RoomFurnishingBed       = "bed"
	RoomFurnishingBench     = "bench"
	RoomFurnishingBookcase  = "bookcase"
	RoomFurnishingCabinet   = "cabinet"
	RoomFurnishingCarpet    = "carpet"
	RoomFurnishingChair     = "chair"
	RoomFurnishingCupboard  = "cupboard"
	RoomFurnishingDesk      = "desk"
	RoomFurnishingDresser   = "dresser"
	RoomFurnishingFirePit   = "firePit"
	RoomFurnishingFireplace = "fireplace"
	RoomFurnishingForge     = "forge"
	RoomFurnishingLamp      = "lamp"
	RoomFurnishingMirror    = "mirror"
	RoomFurnishingPainting  = "painting"
	RoomFurnishingPillar    = "pillar"
	RoomFurnishingRack      = "rack"
	RoomFurnishingShelf     = "shelf"
	RoomFurnishingShrine    = "shrine"
	RoomFurnishingSpit      = "spit"
	RoomFurnishingTableLg   = "tableLg"
	RoomFurnishingTableSm   = "tableSm"
	RoomFurnishingTapestry  = "tapestry"
	RoomFurnishingThrone    = "throne"
	RoomFurnishingTorch     = "torch"
	RoomFurnishingWardrobe  = "wardrobe"
	RoomFurnishingWorkbench = "workbench"
)

var furnishingIndex = map[string]*ItemBase{
	RoomFurnishingAlchemy: NewItemBase("Alchemy equipment", ItemTypeFurnishing, nil),
	RoomFurnishingAnvil:   NewItemBase("Anvil", ItemTypeFurnishing, nil),
	RoomFurnishingBed: NewItemBase("Bed", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Variants = []string{"single", "double", "queen", "king", "bedroll", "cot", "rag pile"}
	}),
	RoomFurnishingBench: NewItemBase("Bench", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Variants = []string{"wood", "cushioned", "stone"}
	}),
	RoomFurnishingBookcase: NewItemBase("Bookcase", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
		ib.Variants = []string{"wood", "metal"}
	}),
	RoomFurnishingCabinet: NewItemBase("Cabinet", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
	}),
	RoomFurnishingCarpet: NewItemBase("Carpet", ItemTypeFurnishing, nil),
	RoomFurnishingChair: NewItemBase("Chair", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Variants = []string{"armchair", "wood", "cushioned", "stone", "stool"}
	}),
	RoomFurnishingCupboard: NewItemBase("Cupboard", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
	}),
	RoomFurnishingDesk: NewItemBase("Desk", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
		ib.Variants = []string{"wood", "stone", "metal"}
		ib.Set = StationarySet
	}),
	RoomFurnishingDresser: NewItemBase("Dresser", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
		ib.Set = ClothingSet
	}),
	RoomFurnishingFirePit:   NewItemBase("Fire pit", ItemTypeFurnishing, nil),
	RoomFurnishingFireplace: NewItemBase("Fireplace", ItemTypeFurnishing, nil),
	RoomFurnishingForge:     NewItemBase("Forge", ItemTypeFurnishing, nil),
	RoomFurnishingLamp:      NewItemBase("Oil lamp", ItemTypeFurnishing, nil),
	RoomFurnishingMirror:    NewItemBase("Mirror, large", ItemTypeFurnishing, nil),
	RoomFurnishingPainting:  NewItemBase("Painting", ItemTypeFurnishing, nil),
	RoomFurnishingPillar:    NewItemBase("Pillar", ItemTypeFurnishing, nil),
	RoomFurnishingRack: NewItemBase("Rack", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
		ib.Variants = []string{"wood", "metal"}
	}),
	RoomFurnishingShelf: NewItemBase("Table, small", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacitySmall
	}),
	RoomFurnishingShrine: NewItemBase("Shrine", ItemTypeFurnishing, nil),
	RoomFurnishingSpit:   NewItemBase("Cooking spit", ItemTypeFurnishing, nil),
	RoomFurnishingTableLg: NewItemBase("Table, large", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
		ib.Variants = []string{"wood", "stone", "metal"}
	}),
	RoomFurnishingTableSm: NewItemBase("Table, small", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacitySmall
		ib.Variants = []string{"wood", "stone", "metal"}
	}),
	RoomFurnishingTapestry: NewItemBase("Tapestry", ItemTypeFurnishing, nil),
	RoomFurnishingThrone:   NewItemBase("Throne", ItemTypeFurnishing, nil),
	RoomFurnishingTorch:    NewItemBase("Torch", ItemTypeFurnishing, nil),
	RoomFurnishingWardrobe: NewItemBase("Wardrobe", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
	}),
	RoomFurnishingWorkbench: NewItemBase("Workbench", ItemTypeFurnishing, func(ib *ItemBase) {
		ib.Capacity = CapacityMedium
		ib.Variants = []string{"wood", "stone", "metal"}
	}),
}
