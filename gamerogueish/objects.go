package gamerogueish

import "fmt"

// ObjectType represents the type of an object.
// TODO: Unify with items?
// A sword might have a hidden compartment...
// A book might contain a letter
// A shelf might contain various things...
// A bed might contain a hidden dagger.
type ObjectType struct {
	Tile        byte
	Name        string
	Description string
	Size        int
	Capacity    int
	Actions     []*Action
}

// Object represents an object in the game.
type Object struct {
	*ObjectType
	Contains []*Item
}

// NewObject returns a new object of the given type.
func (t *ObjectType) New() *Object {
	return &Object{
		ObjectType: t,
	}
}

var (
	ObjectTypeDoor = &ObjectType{
		Tile:        '+',
		Name:        "door",
		Description: "A door.",
		Actions: []*Action{
			&Action{
				Name: "open",
				Func: func(o *Object) {
					fmt.Println("You open the door.")
					// TODO: Make object passable.
				},
			},
			&Action{
				Name: "close",
				Func: func(o *Object) {
					fmt.Println("You close the door.")
					// TODO: Make object impassable.
				},
			},
		},
	}
	ObjectTypeChest = &ObjectType{
		Tile:        'c',
		Name:        "chest",
		Description: "A chest.",
		Actions: []*Action{
			&Action{
				Name: "open",
				Func: func(o *Object) {
					fmt.Println("You open the chest.")
					// TODO: Show inventory.
				},
			},
			&Action{
				Name: "close",
				Func: func(o *Object) {
					fmt.Println("You close the chest.")
					// TODO: Hide inventory.
				},
			},
		},
	}
	ObjectTypeBookshelf = &ObjectType{
		Tile:        'b',
		Name:        "bookshelf",
		Description: "A bookshelf.",
		Actions: []*Action{
			&Action{
				Name: "browse",
				Func: func(o *Object) {
					fmt.Println("You browse the bookshelf.")
					// TODO: Show inventory.
				},
			},
		},
	}
)

type Action struct {
	Name string
	Func func(*Object)
}
