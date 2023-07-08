package gamerogueish

import (
	"fmt"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type UIif interface {
	Draw()
	Select()      // Select selects the current UI element.
	HandleInput() // HandleInput handles input for the current UI element.
}

type uiInventory struct {
	*Game
	view         *console.Console
	selectedItem int
}

func (g *Game) newPlayerInventory() (*uiInventory, error) {
	playerInventoryView, err := g.sideView.CreateSubConsole(0, 6, 20, 10)
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
	ui.view.PrintBounded(1, 0, ui.view.Width-1, 2, fmt.Sprintf("Inventory (%d)", ui.player.Inventory.Count()), t.Background(colGrey))
	var idx int
	for i, item := range ui.player.Items {
		var entry string
		if item.Equipped {
			entry = fmt.Sprintf("%d:*%s", i, item.FullName())
		} else {
			entry = fmt.Sprintf("%d: %s", i, item.FullName())
		}
		var transformers []t.Transformer
		if i == ui.selectedItem && ui.isUIActive(ui) {
			transformers = append(transformers, t.Foreground(concolor.Green))
		}
		ui.view.PrintBounded(1, 1+idx, ui.view.Width-1, 2, entry, transformers...)
		idx++
	}
}

// Selected returns the currently selected item.
func (ui *uiInventory) Selected() *Item {
	if len(ui.player.Inventory.Items) == 0 {
		return nil
	}
	return ui.player.Inventory.Items[ui.selectedItem]
}

func (ui *uiInventory) HandleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ui.SelectItem(ui.selectedItem - 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ui.SelectItem(ui.selectedItem + 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.Select()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		// TODO: Ask for confirmation and drop the item into the world.
		item := ui.player.Inventory.Remove(ui.Selected())
		if item != nil {
			ui.AddMessage(fmt.Sprintf("Dropped %s", item.Name))
			item.X = ui.player.X
			item.Y = ui.player.Y
			ui.World.Items = append(ui.World.Items, item)
		}
		ui.SelectItem(ui.selectedItem)
	}
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
	view         *console.Console
	selectedItem int
}

func (g *Game) newPlayerEnemies() (*uiEnemies, error) {
	playerEnemiesView, err := g.sideView.CreateSubConsole(0, 16, 20, 6)
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

	// Draw what can be found at the current position.
	// List entities first.
	entities := ui.inRange()

	if ui.selectedItem < 0 || ui.selectedItem >= len(entities) {
		ui.selectedItem = 0
	}

	idx := 0
	ui.view.ClearAll()
	ui.view.PrintBounded(1, 0, ui.view.Width-2, 2, fmt.Sprintf("Enemies (%d)", len(entities)), t.Background(colGrey))
	for i, e := range entities {
		entry := e.Name
		var transformers []t.Transformer
		if e.IsDead() {
			entry += " (dead)"
			if ui.isUIActive(ui) && i == ui.selectedItem {
				transformers = append(transformers, t.Foreground(concolor.Green))
			}
		}
		ui.view.PrintBounded(1, 1+idx, ui.view.Width-2, 2, fmt.Sprintf("%d: %s", i, entry), transformers...)
		idx++

		// List inventory items if the entity is dead. (Loot)
		if e.IsDead() {
			for _, it := range e.Items {
				ui.view.PrintBounded(1+3, 1+idx, ui.view.Width-2, 2, it.Name, transformers...)
				idx++
			}
		}
	}
}

func (ui *uiEnemies) HandleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ui.SelectItem(ui.selectedItem - 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ui.SelectItem(ui.selectedItem + 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.Select()
	}
}

// SelectItem selects the given item index while clamping
// it to the entities in range.
func (ui *uiEnemies) SelectItem(index int) {
	entities := ui.inRange()
	if index < 0 {
		index = 0
	} else if index >= len(entities) {
		index = len(entities) - 1
	}
	ui.selectedItem = index
}

func (ui *uiEnemies) inRange() []*Entity {
	// Player position.
	pX := ui.player.X
	pY := ui.player.Y
	var entities []*Entity
	for _, e := range ui.Entities {
		if e.X == pX && e.Y == pY {
			entities = append(entities, e)
		}
	}
	return entities
}

func (ui *uiEnemies) Select() {
	entities := ui.inRange()
	if ui.selectedItem < 0 || ui.selectedItem >= len(entities) {
		ui.selectedItem = 0
		return
	}
	e := entities[ui.selectedItem]
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

type uiItems struct {
	*Game
	view         *console.Console
	selectedItem int
}

func (g *Game) newPlayerItems() (*uiItems, error) {
	playerItemsView, err := g.sideView.CreateSubConsole(0, 22, 20, 6)
	if err != nil {
		return nil, err
	}

	uiItems := &uiItems{
		Game: g,
		view: playerItemsView,
	}
	return uiItems, nil
}

func (ui *uiItems) Draw() {
	ui.view.ClearAll()

	// Draw what can be found at the current position.
	// List items first.
	items := ui.inRange()

	if ui.selectedItem < 0 || ui.selectedItem >= len(items) {
		ui.selectedItem = 0
	}

	idx := 0
	ui.view.ClearAll()

	ui.view.PrintBounded(1, 0, ui.view.Width-2, 2, fmt.Sprintf("Items (%d)", len(items)), t.Background(colGrey))
	for i, it := range items {
		var transformers []t.Transformer
		if ui.isUIActive(ui) && i == ui.selectedItem {
			transformers = append(transformers, t.Foreground(concolor.Green))
		}
		ui.view.PrintBounded(1, 1+idx, ui.view.Width-2, 2, fmt.Sprintf("%d: %s", i, it.Name), transformers...)
		idx++
	}
}

func (ui *uiItems) HandleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ui.SelectItem(ui.selectedItem - 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ui.SelectItem(ui.selectedItem + 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.Select()
	}
}

// SelectItem selects the given item index while clamping
// it to the items in range.
func (ui *uiItems) SelectItem(index int) {
	items := ui.inRange()
	if index < 0 {
		index = 0
	} else if index >= len(items) {
		index = len(items) - 1
	}
	ui.selectedItem = index
}

func (ui *uiItems) inRange() []*Item {
	// Player position.
	pX := ui.player.X
	pY := ui.player.Y
	var items []*Item
	for _, it := range ui.Items {
		if it.X == pX && it.Y == pY {
			items = append(items, it)
		}
	}
	return items
}

func (ui *uiItems) Select() {
	items := ui.inRange()
	if ui.selectedItem < 0 || ui.selectedItem >= len(items) {
		ui.selectedItem = 0
		return
	}
	it := items[ui.selectedItem]
	ui.player.Inventory.Add(it)
	// Find the actual index and remove the item.
	var idx int
	for i, it2 := range ui.Items {
		if it == it2 {
			idx = i
			break
		}
	}
	ui.Items = append(ui.Items[:idx], ui.Items[idx+1:]...)
}

type uiPlayerInfo struct {
	*Game
	view *console.Console
}

func (g *Game) newPlayerInfo() (*uiPlayerInfo, error) {
	playerInfoView, err := g.sideView.CreateSubConsole(0, 0, 20, 5)
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
	ui.view.PrintBounded(1, 1, ui.view.Width, 1, "Player: "+ui.player.Name, t.Background(colGrey))
	ui.view.PrintBounded(1, 2, ui.view.Width-2, 2, fmt.Sprintf("Health: %d/%d", ui.player.Health, ui.player.BaseHealth))
	ui.view.PrintBounded(1, 3, ui.view.Width-2, 2, fmt.Sprintf("Def: %d Att: %d", ui.player.DefenseValue(), ui.player.AttackDamage()))
	ui.view.PrintBounded(1, 4, ui.view.Width-2, 2, fmt.Sprintf("X=%d Y=%d", pX, pY), t.Foreground(colGrey))
}

func (ui *uiPlayerInfo) HandleInput() {
}

func (ui *uiPlayerInfo) Select() {
}
