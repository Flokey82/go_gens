package main

import (
	"github.com/Flokey82/go_gens/genfurnishing"
)

func main() {
	// Generate a new room.
	room := genfurnishing.RoomBedroom.Generate()
	room.Log()

	// Generate a new building.
	building := genfurnishing.BuildingKeep.Generate()
	building.Log()
}
