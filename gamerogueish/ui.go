package gamerogueish

import (
	"fmt"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
)

type UIif interface {
	Draw()
	Next()   // Next selects the next UI element.
	Prev()   // Prev selects the previous UI element.
	Select() // Select selects the current UI element.
}

type uiInventory struct {
	*Game
	view         *console.Console
	selectedItem int
}

func (g *Game) newPlayerInventory() (*uiInventory, error) {
	playerInventoryView, err := g.sideView.CreateSubConsole(0, 5, 20, 10)
	if err != nil {
		return nil, err
	}

	uiInventory := &uiInventory{
		Game: g,
		view: playerInventoryView,
	}
	return uiInventory, nil
}

func (ui *uiInventory) Draw() {
	ui.view.ClearAll()
	ui.view.PrintBounded(1, 0, ui.view.Width-2, 2, fmt.Sprintf("Inventory (%d)", ui.player.Inventory.Count()))
	var idx int
	for i, item := range ui.player.Items {
		var entry string
		if item.Equipped {
			entry = fmt.Sprintf("%d:*%s", i, item.Name)
		} else {
			entry = fmt.Sprintf("%d: %s", i, item.Name)
		}
		var transformers []t.Transformer
		if i == ui.selectedItem && ui.isUIActive(ui) {
			transformers = append(transformers, t.Foreground(concolor.Green))
		}
		ui.view.PrintBounded(2, 1+idx, ui.view.Width-2, 2, entry, transformers...)
		idx++
	}
}

func (ui *uiInventory) Next() {
	ui.SelectItem(ui.selectedItem + 1)
}

func (ui *uiInventory) Prev() {
	ui.SelectItem(ui.selectedItem - 1)
}

// Selected returns the currently selected item.
func (ui *uiInventory) Selected() *Item {
	if len(ui.player.Inventory.Items) == 0 {
		return nil
	}
	return ui.player.Inventory.Items[ui.selectedItem]
}

// SelectItem selects the given item index while clamping
// it to the inventory.
func (ui *uiInventory) SelectItem(index int) {
	if index < 0 {
		index = 0
	} else if index >= len(ui.player.Inventory.Items) {
		index = len(ui.player.Inventory.Items) - 1
	}
	ui.selectedItem = index
}

func (ui *uiInventory) Select() {
	if sel := ui.Selected(); sel != nil {
		if sel.Equippable() {
			ui.player.Equip(ui.selectedItem)
		} else if sel.Consumable() {
			ui.player.Consume(ui.selectedItem)
		}
		ui.SelectItem(ui.selectedItem) // Update the selected item.
	}
}

type uiEnemies struct {
	*Game
	view          *console.Console
	selectedEnemy int
}

func (g *Game) newPlayerEnemies() (*uiEnemies, error) {
	playerEnemiesView, err := g.sideView.CreateSubConsole(0, 15, 20, 10)
	if err != nil {
		return nil, err
	}

	uiEnemies := &uiEnemies{
		Game: g,
		view: playerEnemiesView,
	}
	return uiEnemies, nil
}

func (ui *uiEnemies) Draw() {
	ui.view.ClearAll()
	// Player position.
	pX := ui.player.X
	pY := ui.player.Y

	// Draw what can be found at the current position.
	// List entities first.
	// TODO: Factor this out into a function.
	var entities []*Entity
	for _, e := range ui.Entities {
		if e.X == pX && e.Y == pY {
			entities = append(entities, e)
		}
	}

	if ui.selectedEnemy < 0 || ui.selectedEnemy >= len(entities) {
		ui.selectedEnemy = 0
	}

	idx := 0
	ui.view.ClearAll()
	ui.view.PrintBounded(1, 0, ui.view.Width-2, 2, fmt.Sprintf("In Range (%d)", len(entities)))
	for i, e := range entities {
		entry := e.Name
		var transformers []t.Transformer
		if e.IsDead() {
			entry += " (dead)"
			if ui.isUIActive(ui) && i == ui.selectedEnemy {
				transformers = append(transformers, t.Foreground(concolor.Red))
			} else {
				transformers = append(transformers, t.Foreground(colDarkRed))
			}
		}
		ui.view.PrintBounded(2, 1+idx, ui.view.Width-2, 2, fmt.Sprintf("%d: %s", i, entry), transformers...)
		idx++

		// List inventory items if the entity is dead. (Loot)
		if e.IsDead() {
			for _, it := range e.Items {
				ui.view.PrintBounded(2+3, 1+idx, ui.view.Width-2, 2, it.Name, transformers...)
				idx++
			}
		}
	}
}

func (ui *uiEnemies) Next() {
	ui.selectedEnemy++
	if ui.selectedEnemy >= len(ui.Entities) {
		ui.selectedEnemy = 0
	}
}

func (ui *uiEnemies) Prev() {
	ui.selectedEnemy--
	if ui.selectedEnemy < 0 {
		ui.selectedEnemy = len(ui.Entities) - 1
	}
}

func (ui *uiEnemies) Select() {
	// Player position.
	pX := ui.player.X
	pY := ui.player.Y
	// TODO: Factor this out into a function.
	var entities []*Entity
	for _, e := range ui.Entities {
		if e.X == pX && e.Y == pY {
			entities = append(entities, e)
		}
	}
	if ui.selectedEnemy < 0 || ui.selectedEnemy >= len(entities) {
		ui.selectedEnemy = 0
		return
	}
	e := entities[ui.selectedEnemy]
	if e.IsDead() {
		// Loot
		for _, it := range e.Items {
			ui.player.Inventory.Add(it)
		}
		// Find the actual index and remove the entity.
		// TODO: Instead we should just let the entity be dead or
		// let it decay after some time.
		var idx int
		for i, e2 := range ui.Entities {
			if e == e2 {
				idx = i
				break
			}
		}

		ui.Entities = append(ui.Entities[:idx], ui.Entities[idx+1:]...)
	} else {
		// Attack
		ui.player.Attack(ui.Game, e)
	}
}

type uiPlayerInfo struct {
	*Game
	view *console.Console
}

func (g *Game) newPlayerInfo() (*uiPlayerInfo, error) {
	playerInfoView, err := g.sideView.CreateSubConsole(0, 0, 20, 4)
	if err != nil {
		return nil, err
	}

	uiInfo := &uiPlayerInfo{
		Game: g,
		view: playerInfoView,
	}
	return uiInfo, nil
}

func (ui *uiPlayerInfo) Draw() {
	ui.view.ClearAll()
	// Player position.
	pX := ui.player.X
	pY := ui.player.Y

	// Draw player info.
	ui.view.PrintBounded(1, 1, ui.view.Width-2, 2, fmt.Sprintf("Health: %d/%d", ui.player.Health, ui.player.BaseHealth))
	ui.view.PrintBounded(1, 2, ui.view.Width-2, 2, fmt.Sprintf("Def: %d Att: %d", ui.player.DefenseValue(), ui.player.AttackDamage()))
	ui.view.PrintBounded(1, 3, ui.view.Width-2, 2, fmt.Sprintf("X=%d Y=%d", pX, pY), t.Foreground(colGrey))
}

func (ui *uiPlayerInfo) Next() {
}

func (ui *uiPlayerInfo) Prev() {
}

func (ui *uiPlayerInfo) Select() {
}
