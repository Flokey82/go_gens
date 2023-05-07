package genfurnishing

import (
	"log"
	"math/rand"
)

// A building is a collection of rooms serving a purpose with a total size and a number of inhabitants with
// a certain social status, profession, and hobbies.
type Building struct {
	Name  string  // Name of the building.
	Rooms []*Room // Rooms in the building.
	Size  int     // Size of the building.
}

type BuildingConfig struct {
	Name             string   // Name of the building type.
	MinSize          int      // Minimum size of the building.
	MaxSize          int      // Maximum size of the building.
	Required         []string // Required rooms
	Possible         []string // Possible rooms
	PossibleMultiple []string // Possible rooms that can be added multiple times.
}

func (bc *BuildingConfig) Generate() *Building {
	building := &Building{
		Name: bc.Name,
		Size: bc.MinSize + rand.Intn(bc.MaxSize-bc.MinSize),
	}
	for _, room := range bc.Required {
		building.Rooms = append(building.Rooms, roomTypeToConfig[room].Generate())
	}
	for _, i := range rand.Perm(len(bc.Possible)) {
		room := roomTypeToConfig[bc.Possible[i]].Generate()
		if len(building.Rooms) >= building.Size {
			break
		}
		building.Rooms = append(building.Rooms, room)
	}
	return building
}

func (b *Building) Log() {
	log.Printf("Building: %s", b.Name)
	for _, room := range b.Rooms {
		room.Log()
	}
}

const (
	TypeKeep       = "Keep"
	TypeGreatHouse = "Great House"
	// TypeManor      = "Manor"
	// TypeTownhouse  = "Townhouse"
	// TypeCottage    = "Cottage"
	// TypeChurch     = "Church"
	// TypeTemple     = "Temple"
	// TypeGuardhouse = "Guardhouse"
	// TypeInn        = "Inn"
)

var (
	BuildingKeep = &BuildingConfig{
		Name:    TypeKeep,
		MinSize: 10,
		MaxSize: 15,
		Required: []string{
			RoomTypeGreatHall,
			RoomTypeKitchen,
			RoomTypePantry,
			RoomTypeBedroom,
			RoomTypeArmory,
			RoomTypePrison,
		},
		Possible: []string{
			RoomTypeAtrium,
			RoomTypeBallroom,
			RoomTypeBathhouse,
			RoomTypeChamber,
			RoomTypeDining,
			RoomTypeDormitory,
			RoomTypeHallway,
			RoomTypeLaboratory,
			RoomTypeLibrary,
			RoomTypeParlour,
			RoomTypeThrone,
			RoomTypeTorture,
			RoomTypeTreasury,
		},
	}
	BuildingGreatHouse = &BuildingConfig{
		Name:    TypeGreatHouse,
		MinSize: 6,
		MaxSize: 12,
		Required: []string{
			RoomTypeGreatHall,
			RoomTypeKitchen,
			RoomTypePantry,
			RoomTypeBedroom,
		},
		Possible: []string{
			RoomTypeArmory,
			RoomTypeAtrium,
			RoomTypeBallroom,
			RoomTypeBathhouse,
			RoomTypeChamber,
			RoomTypeDining,
			RoomTypeDormitory,
			RoomTypeHallway,
			RoomTypeLaboratory,
			RoomTypeLibrary,
		},
	}
)
