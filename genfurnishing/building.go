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

// BuildingConfig represents a building prototype and all the possible
// variations of rooms that can be found in it.
type BuildingConfig struct {
	Name     string   // Name of the building type.
	Required []string // Required rooms
	Possible []string // Possible rooms
	Size     [2]int   // Minimum/Maximum size of the building.
}

// NewBuildingConfig returns a new building configuration for generating a specific type of building.
func NewBuildingConfig(name string, required, possible []string, minSize, maxSize int) *BuildingConfig {
	return &BuildingConfig{
		Name:     name,
		Required: required,
		Possible: possible,
		Size:     [2]int{minSize, maxSize},
	}
}

// Generate a building from the config.
func (bc *BuildingConfig) Generate() *Building {
	building := &Building{
		Name: bc.Name,
		Size: bc.Size[0] + rand.Intn(bc.Size[1]-bc.Size[0]),
	}

	// Generate all required rooms.
	for _, room := range bc.Required {
		building.Rooms = append(building.Rooms, roomTypeToConfig[room].Generate())
	}

	// Generate random rooms.
	for _, i := range rand.Perm(len(bc.Possible)) {
		room := roomTypeToConfig[bc.Possible[i]].Generate()
		building.Rooms = append(building.Rooms, room)
		if len(building.Rooms) >= building.Size {
			break
		}
	}
	return building
}

// Log the building.
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
	BuildingKeep = NewBuildingConfig(TypeKeep, []string{
		RoomTypeGreatHall,
		RoomTypeKitchen,
		RoomTypePantry,
		RoomTypeBedroom,
		RoomTypeArmory,
		RoomTypePrison,
	}, []string{
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
	}, 10, 15)
	BuildingGreatHouse = NewBuildingConfig(TypeGreatHouse, []string{
		RoomTypeGreatHall,
		RoomTypeKitchen,
		RoomTypePantry,
		RoomTypeBedroom,
	}, []string{
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
	}, 6, 12)
)
