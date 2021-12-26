package gamecs

import "log"

type CInventory struct {
	BodySlots map[int]*Item // Head, Feet, Upper Body, ...
	Slots     []*Item
	Size      int

	w  *World
	id int
}

func newCInventory(w *World, id, size int) *CInventory {
	return &CInventory{
		w:    w,
		id:   id,
		Size: size,
	}
}

func (in *CInventory) IsFull() bool {
	return len(in.Slots) >= in.Size
}

// Add item to inventory... return success.
func (in *CInventory) Add(it *Item) bool {
	if in.IsFull() {
		return false
	}
	log.Println("Add item to inventory")
	it.Location = LocInventory
	it.LocationID = in.id
	in.Slots = append(in.Slots, it)
	return true
}

// Has itemType in inventory... return success.
func (in *CInventory) Has(itt *ItemType) bool {
	for _, it := range in.Slots {
		if it.ItemType == itt {
			return true
		}
	}
	return false
}

// Equip item from inventory... return success.
func (in *CInventory) Equip(it *Item) bool {
	//for _, slt := range it.EquipSlots {
	//	if _, ok := in.BodySlot; !ok {
	//		return false
	//	}
	//}
	// TODO: Move all equipped items into inventory and assign item to slots.
	return false
}

// Find item in inventory... return item.
func (in *CInventory) Find(tag string) *Item {
	for _, it := range in.Slots {
		for _, t := range it.Tags {
			if t == tag {
				return it
			}
		}
	}
	return nil
}

// Drop item from inventory... return success.
func (in *CInventory) Drop(it *Item) bool {
	// TODO: Many nil pointer checks needed here...
	if in.RemoveID(it.id) {
		it.Location = LocWorld
		it.LocationID = -1
		it.Pos = in.w.mgr.GetEntityFromID(in.id).Pos
		return true
	}
	return false
}

// RemoveID removes the item with the given id from inventory... return success.
func (in *CInventory) RemoveID(id int) bool {
	for i, it := range in.Slots {
		if it.id == id {
			in.Slots = append(in.Slots[:i], in.Slots[i+1:]...)
			return true
		}
	}
	return false
}

func (in *CInventory) TransferAll(to *CInventory) bool {
	for len(in.Slots) > 0 {
		it := in.Slots[0]
		if !to.Add(it) { // Out of storage?
			return false
		}
		in.RemoveID(it.id)
	}
	return true
}

type BodySlot int

// TODO: What if this thing has six tentacles and no head?
const (
	BodySlotHead BodySlot = iota
	BodySlotNeck
	BodySlotChest
	BodySlotLowerBody
	BodySlotFeet
	BodySlotHandLeft
	BodySlotHandRight
	BodySlotMouth
)
