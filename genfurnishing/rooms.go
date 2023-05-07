package genfurnishing

import (
	"log"
)

// RoomConfig represents a room prototype and all the possible
// variations of items and furniture that can be found in it.
type RoomConfig struct {
	Name     string
	Required []string
	Possible []string
}

// NewRoomConfig returns a new room configuration.
func NewRoomConfig(name string, required, possible []string) *RoomConfig {
	return &RoomConfig{
		Name:     name,
		Required: required,
		Possible: possible,
	}
}

// Generate a room from the config.
func (rc *RoomConfig) Generate() *Room {
	return &Room{
		Name:  rc.Name,
		Items: generateItems(rc.Required, rc.Possible),
	}
}

func generateItems(required, possible []string) []*Item {
	items := []*Item{}
	for _, item := range required {
		items = append(items, furnishingIndex[item].Generate())
	}
	for _, item := range possible {
		itemProto := furnishingIndex[item]
		if itemProto.Rarity.Roll() {
			items = append(items, itemProto.Generate())
		}
	}
	return items
}

// Room is a generated room.
type Room struct {
	Name  string
	Items []*Item
}

func (r *Room) Log() {
	log.Printf("Room: %s", r.Name)
	for _, item := range r.Items {
		log.Printf("  %s", item.Name())
		for _, subItem := range item.Contains {
			log.Printf("    %s", subItem.Name())
		}
	}
}

const (
	RoomTypeArmory     = "armory"
	RoomTypeAtrium     = "atrium"
	RoomTypeBallroom   = "ballroom"
	RoomTypeBathhouse  = "bathhouse"
	RoomTypeBedroom    = "bedroom"
	RoomTypeChamber    = "chamber"
	RoomTypeDining     = "dining"
	RoomTypeDormitory  = "dormitory"
	RoomTypeGreatHall  = "greatHall"
	RoomTypeHallway    = "hallway"
	RoomTypeKitchen    = "kitchen"
	RoomTypeLaboratory = "laboratory"
	RoomTypeLibrary    = "library"
	RoomTypePantry     = "pantry"
	RoomTypeParlour    = "parlour"
	RoomTypePrison     = "prison"
	RoomTypeRoom       = "room"
	RoomTypeShrine     = "shrine"
	RoomTypeSmithy     = "smithy"
	RoomTypeStorage    = "storage"
	RoomTypeStudy      = "study"
	RoomTypeThrone     = "throne"
	RoomTypeTorture    = "torture"
	RoomTypeTreasury   = "treasury"
)

var (
	RoomArmory = NewRoomConfig(RoomTypeArmory,
		[]string{
			RoomFurnishingRack},
		[]string{
			RoomFurnishingAnvil,
			RoomFurnishingBench,
			RoomFurnishingCabinet,
			RoomFurnishingForge,
			RoomFurnishingLamp,
			RoomFurnishingRack,
			RoomFurnishingTableLg,
			RoomFurnishingShelf,
			RoomFurnishingTorch,
			RoomFurnishingWorkbench},
	)
	RoomAtrium = NewRoomConfig(RoomTypeAtrium,
		nil,
		[]string{
			RoomFurnishingBench,
			RoomFurnishingCarpet,
			RoomFurnishingPillar},
	)
	RoomBallroom = NewRoomConfig(RoomTypeBallroom,
		nil,
		[]string{
			RoomFurnishingBench,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingFireplace,
			RoomFurnishingLamp,
			RoomFurnishingTableLg,
			RoomFurnishingTableSm},
	)
	RoomBathhouse = NewRoomConfig(RoomTypeBathhouse,
		nil,
		[]string{
			RoomFurnishingBench,
			RoomFurnishingRack,
			RoomFurnishingShelf},
	)
	RoomBedroom = NewRoomConfig(RoomTypeBedroom,
		[]string{
			RoomFurnishingBed},
		[]string{
			RoomFurnishingBed,
			RoomFurnishingBench,
			RoomFurnishingBookcase,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingDesk,
			RoomFurnishingDresser,
			RoomFurnishingFireplace,
			RoomFurnishingLamp,
			RoomFurnishingMirror,
			RoomFurnishingTableSm,
			RoomFurnishingShelf,
			RoomFurnishingShrine,
			RoomFurnishingWardrobe},
	)
	RoomChamber = NewRoomConfig(RoomTypeChamber,
		nil,
		[]string{
			RoomFurnishingBookcase,
			RoomFurnishingCabinet,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingDesk,
			RoomFurnishingFireplace,
			RoomFurnishingLamp,
			RoomFurnishingTableSm,
			RoomFurnishingShelf,
			RoomFurnishingTorch},
	)
	RoomDining = NewRoomConfig(RoomTypeDining,
		[]string{
			RoomFurnishingTableLg},
		[]string{
			RoomFurnishingBench,
			RoomFurnishingCabinet,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingCupboard,
			RoomFurnishingFireplace,
			RoomFurnishingLamp,
			RoomFurnishingTableLg,
			RoomFurnishingTableSm,
			RoomFurnishingSpit,
			RoomFurnishingTorch},
	)
	RoomDormitory = NewRoomConfig(RoomTypeDormitory,
		[]string{
			RoomFurnishingBed},
		[]string{
			RoomFurnishingBed,
			RoomFurnishingCarpet,
			RoomFurnishingBench,
			RoomFurnishingBookcase,
			RoomFurnishingChair,
			RoomFurnishingCupboard,
			RoomFurnishingDesk,
			RoomFurnishingDresser,
			RoomFurnishingFireplace,
			RoomFurnishingPillar,
			RoomFurnishingRack,
			RoomFurnishingTableSm,
			RoomFurnishingShelf,
			RoomFurnishingTorch},
	)
	RoomGreatHall = NewRoomConfig(RoomTypeGreatHall,
		nil,
		[]string{
			RoomFurnishingBench,
			RoomFurnishingCarpet,
			RoomFurnishingBookcase,
			RoomFurnishingFireplace,
			RoomFurnishingForge,
			RoomFurnishingLamp,
			RoomFurnishingPillar,
			RoomFurnishingRack,
			RoomFurnishingTableLg,
			RoomFurnishingThrone,
			RoomFurnishingShrine,
			RoomFurnishingTorch},
	)
	RoomHallway = NewRoomConfig(RoomTypeHallway,
		nil,
		[]string{
			RoomFurnishingCarpet,
			RoomFurnishingShelf,
			RoomFurnishingTorch},
	)
	RoomKitchen = NewRoomConfig(RoomTypeKitchen,
		[]string{
			RoomFurnishingTableSm,
			RoomFurnishingSpit},
		[]string{
			RoomFurnishingFirePit,
			RoomFurnishingFireplace,
			RoomFurnishingLamp,
			RoomFurnishingRack,
			RoomFurnishingTableLg,
			RoomFurnishingTableSm,
			RoomFurnishingShelf,
			RoomFurnishingSpit,
			RoomFurnishingWorkbench},
	)
	RoomLaboratory = NewRoomConfig(RoomTypeLaboratory,
		[]string{
			RoomFurnishingAlchemy,
			RoomFurnishingWorkbench},
		[]string{
			RoomFurnishingAlchemy,
			RoomFurnishingBench,
			RoomFurnishingBookcase,
			RoomFurnishingCabinet,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingDesk,
			RoomFurnishingFireplace,
			RoomFurnishingLamp,
			RoomFurnishingMirror,
			RoomFurnishingRack,
			RoomFurnishingTableSm,
			RoomFurnishingTableLg,
			RoomFurnishingShelf,
			RoomFurnishingTorch,
			RoomFurnishingWorkbench},
	)
	RoomLibrary = NewRoomConfig(RoomTypeLibrary,
		[]string{
			RoomFurnishingBookcase},
		[]string{
			RoomFurnishingBench,
			RoomFurnishingBookcase,
			RoomFurnishingCabinet,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingDesk,
			RoomFurnishingFireplace,
			RoomFurnishingLamp,
			RoomFurnishingTableLg,
			RoomFurnishingTableSm,
			RoomFurnishingShelf},
	)
	RoomPantry = NewRoomConfig(RoomTypePantry,
		[]string{
			RoomFurnishingShelf},
		[]string{
			RoomFurnishingCabinet,
			RoomFurnishingCupboard,
			RoomFurnishingRack,
			RoomFurnishingShelf},
	)
	RoomParlour = NewRoomConfig(RoomTypeParlour,
		nil,
		[]string{
			RoomFurnishingBench,
			RoomFurnishingBookcase,
			RoomFurnishingCabinet,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingDesk,
			RoomFurnishingTableSm},
	)
	RoomPrison = NewRoomConfig(RoomTypePrison,
		nil, nil)
	RoomRoom = NewRoomConfig(RoomTypeRoom,
		nil,
		[]string{
			RoomFurnishingCarpet,
			RoomFurnishingFirePit,
			RoomFurnishingTableSm,
			RoomFurnishingTorch},
	)
	RoomShrine = NewRoomConfig(RoomTypeShrine,
		[]string{
			RoomFurnishingShrine},
		[]string{
			RoomFurnishingCarpet,
			RoomFurnishingLamp,
			RoomFurnishingShrine,
			RoomFurnishingTorch},
	)
	RoomSmithy = NewRoomConfig(RoomTypeSmithy,
		[]string{
			RoomFurnishingAnvil,
			RoomFurnishingForge,
			RoomFurnishingWorkbench},
		nil,
	)
	RoomStorage = NewRoomConfig(RoomTypeStorage,
		[]string{
			RoomFurnishingRack},
		[]string{
			RoomFurnishingCabinet,
			RoomFurnishingCupboard,
			RoomFurnishingRack,
			RoomFurnishingTableSm,
			RoomFurnishingShelf},
	)
	RoomStudy = NewRoomConfig(RoomTypeStudy,
		[]string{
			RoomFurnishingChair,
			RoomFurnishingDesk},
		[]string{
			RoomFurnishingBookcase,
			RoomFurnishingCabinet,
			RoomFurnishingCarpet,
			RoomFurnishingChair,
			RoomFurnishingDesk,
			RoomFurnishingLamp,
			RoomFurnishingTableSm,
			RoomFurnishingShelf},
	)
	RoomThrone = NewRoomConfig(RoomTypeThrone,
		[]string{
			RoomFurnishingThrone},
		[]string{
			RoomFurnishingBench,
			RoomFurnishingCarpet,
			RoomFurnishingLamp,
			RoomFurnishingPillar,
			RoomFurnishingTableLg,
			RoomFurnishingThrone,
			RoomFurnishingTorch},
	)
	RoomTorture = NewRoomConfig(RoomTypeTorture,
		[]string{
			RoomFurnishingChair},
		[]string{
			RoomFurnishingFireplace,
			RoomFurnishingTorch,
			RoomFurnishingWorkbench},
	)
	RoomTreasury = NewRoomConfig(RoomTypeTreasury,
		nil,
		[]string{
			RoomFurnishingCarpet,
			RoomFurnishingDesk,
			RoomFurnishingLamp,
			RoomFurnishingMirror,
			RoomFurnishingRack,
			RoomFurnishingTableLg,
			RoomFurnishingTableSm},
	)
)

var roomTypeToConfig = map[string]*RoomConfig{
	RoomTypeArmory:     RoomArmory,
	RoomTypeAtrium:     RoomAtrium,
	RoomTypeBallroom:   RoomBallroom,
	RoomTypeBathhouse:  RoomBathhouse,
	RoomTypeBedroom:    RoomBedroom,
	RoomTypeChamber:    RoomChamber,
	RoomTypeDining:     RoomDining,
	RoomTypeDormitory:  RoomDormitory,
	RoomTypeGreatHall:  RoomGreatHall,
	RoomTypeHallway:    RoomHallway,
	RoomTypeKitchen:    RoomKitchen,
	RoomTypeLaboratory: RoomLaboratory,
	RoomTypeLibrary:    RoomLibrary,
	RoomTypePantry:     RoomPantry,
	RoomTypeParlour:    RoomParlour,
	RoomTypePrison:     RoomPrison,
	RoomTypeRoom:       RoomRoom,
	RoomTypeShrine:     RoomShrine,
	RoomTypeSmithy:     RoomSmithy,
	RoomTypeStorage:    RoomStorage,
	RoomTypeStudy:      RoomStudy,
	RoomTypeThrone:     RoomThrone,
	RoomTypeTorture:    RoomTorture,
	RoomTypeTreasury:   RoomTreasury,
}
