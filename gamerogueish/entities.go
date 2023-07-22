package gamerogueish

import (
	"fmt"
	"math/rand"

	"github.com/Flokey82/gamedice"
)

type EntityType struct {
	Tile              byte
	Name              string
	Description       string // TODO: Generate description to create variety.
	BaseHealth        int
	BaseAttack        int
	BaseDefense       int
	Equipment         []*ItemType
	OptionalEquipment []*ItemType // TODO: Change to item collections.
}

var (
	EntityPlayer = &EntityType{
		Tile:        '@',
		Name:        "player",
		BaseHealth:  10,
		BaseAttack:  2,
		BaseDefense: 10,
		Equipment: []*ItemType{
			ItemTypeWeaponFishingRod,
			ItemTypeArmorPlate,
			ItemTypeHelmetSweatband,
			ItemTypePotion,
			ItemTypeNote,
		},
	}
	EntityGoblin = &EntityType{
		Tile:        'g',
		Name:        "goblin",
		Description: "A small goblin.",
		BaseHealth:  5,
		BaseAttack:  1,
		BaseDefense: 5,
		Equipment:   []*ItemType{ItemTypeWeaponAxe, ItemTypeArmorChain},
		OptionalEquipment: []*ItemType{
			ItemTypePotion,
			ItemTypeGoblinToe,
		},
	}
	EntityOrc = &EntityType{
		Tile:        'o',
		Name:        "orc",
		Description: "A big orc.",
		BaseHealth:  10,
		BaseAttack:  5,
		BaseDefense: 14,
		Equipment:   []*ItemType{ItemTypeWeaponSword, ItemTypeArmorLeather},
		OptionalEquipment: []*ItemType{
			ItemTypePotion,
		},
	}
	EntityTroll = &EntityType{
		Tile:        't',
		Name:        "troll",
		Description: "A huge troll.",
		BaseHealth:  15,
		BaseAttack:  7,
		BaseDefense: 15,
		Equipment:   []*ItemType{ItemTypeTrollPoop},
	}
)

var MonsterEntities = []*EntityType{
	EntityGoblin,
	EntityOrc,
	EntityTroll,
}

type Entity struct {
	Name        string             // name of entity (might be different from EntityType.Name)
	*EntityType                    // type of entity
	Inventory                      // inventory component
	X           int                // x position in the world
	Y           int                // y position in the world
	Health      int                // health points
	Slots       [ItemTypeMax]*Item // Equipped items.
}

// NewEntity returns a new entity with the given position and tile.
func NewEntity(x, y int, e *EntityType) *Entity {
	entity := &Entity{
		Name:       e.Name,
		X:          x,
		Y:          y,
		EntityType: e,
		Health:     e.BaseHealth,
	}
	// Add equipment.
	for _, it := range e.Equipment {
		entity.Inventory.Add(it.New())
	}

	// Equip all items that can be equipped.
	for i, it := range entity.Items {
		if it.Equippable() {
			entity.Equip(i)
		}
	}

	// Add optional equipment.
	for _, it := range e.OptionalEquipment {
		if rand.Intn(2) == 0 {
			entity.Inventory.Add(it.New())
		}
	}
	return entity
}

// Equip equips the item at the given inventory index.
func (e *Entity) Equip(index int) {
	if index < 0 || index >= len(e.Items) || !e.Items[index].Equippable() {
		return
	}

	// Toggle equipped state.
	it := e.Items[index]
	it.Equipped = !it.Equipped

	// If we unequip the item, unset the equipped item.
	if !it.Equipped {
		e.Slots[it.Type] = nil
		return
	}

	// If there is already an item in the slot, unequip it.
	if e.Slots[it.Type] != nil {
		e.Slots[it.Type].Equipped = false
	}
	e.Slots[it.Type] = it
}

func (e *Entity) Attack(g *Game, target *Entity) {
	// Check if attack roll is successful.
	if roll := gamedice.D20.Roll(); roll >= target.DefenseValue() {
		g.AddMessage(fmt.Sprintf("%s hit %s (%d/%d)", e.Name, target.Name, roll, target.DefenseValue()))
		target.TakeDamage(g, e.AttackDamage())
	} else {
		g.AddMessage(fmt.Sprintf("%s missed %s (%d/%d)", e.Name, target.Name, roll, target.DefenseValue()))
	}
}

func (e *Entity) AttackDamage() int {
	damage := e.BaseAttack // Unarmed attack.
	// Check if we have a weapon equipped.
	// TODO: Allow weapon specific damage.
	if e.Slots[ItemWeapon] != nil {
		damage = 5 + e.Slots[ItemWeapon].Modifier
	}
	return damage
}

func (e *Entity) DefenseValue() int {
	defense := e.BaseDefense // Unarmored defense.
	// Check if we have armor equipped.
	// TODO: Allow armor specific defense.
	if e.Slots[ItemArmor] != nil {
		defense += 2 + e.Slots[ItemArmor].Modifier
	}
	return defense
}

func (e *Entity) TakeDamage(g *Game, damage int) {
	g.AddMessage(fmt.Sprintf("%s took %d damage", e.Name, damage))
	e.Health -= damage
	if e.Health <= 0 {
		g.AddMessage(fmt.Sprintf("%s died", e.Name))
		e.Health = 0
	}
}

func (e *Entity) IsDead() bool {
	return e.Health <= 0
}

// Consume consumes the item at the given inventory index.
func (e *Entity) Consume(index int) {
	if index < 0 || index >= len(e.Items) || !e.Items[index].Consumable() {
		return
	}
	// For now, we assume this is a health potion.
	// If we are full health, we do nothing.
	if e.Health == e.BaseHealth {
		return
	}
	// TODO: Add more potion types.
	e.Health += 5 + e.Items[index].Modifier
	if e.Health > e.BaseHealth {
		e.Health = e.BaseHealth
	}
	e.RemoveItem(index)
}
