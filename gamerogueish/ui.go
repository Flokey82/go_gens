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
	view *console.Console
	selectableList[Item]
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
	// Update the selectable list.
	ui.selectableList.setItems(ui.player.Inventory.Items)

	ui.view.ClearAll()
	// Make sure the active UI is highlighted.
	ts := []t.Transformer{t.Background(colGrey)}
	if ui.isUIActive(ui) {
		ts = append(ts, t.Foreground(concolor.Green))
	}
	ui.view.PrintBounded(1, 0, ui.view.Width-1, 2, fmt.Sprintf("Inventory (%d)", ui.player.Inventory.Count()), ts...)

	var idx int
	start, end := calcVisibleRange(ui.view.Height-2, len(ui.player.Inventory.Items), ui.selectedItem)
	for i := start; i < end; i++ {
		item := ui.player.Inventory.Items[i]
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

func (ui *uiInventory) HandleInput() {
	// Update the selectable list.
	ui.selectableList.handleInput()

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

	// Show information about the selected item (if any).
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		// Display information about the item.
		// TODO: Display in a text box.
		if sel := ui.Selected(); sel != nil {
			ui.AddMessage(sel.Name + ": " + sel.Description)
		}
	}
}

func (ui *uiInventory) Select() {
	if sel := ui.Selected(); sel != nil {
		if sel.Equippable() {
			ui.player.Equip(ui.selectedItem)
		} else if sel.Consumable() {
			ui.player.Consume(ui.selectedItem)
		} else if sel.OnUse != nil {
			sel.OnUse(ui.Game, ui.player, sel)
		}
		ui.SelectItem(ui.selectedItem) // Update the selected item.
	}
}

type uiEnemies struct {
	*Game
	view *console.Console
	selectableList[Entity]
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
	// Update the selectable list.
	ui.selectableList.setItems(ui.inRange())

	// Draw what can be found at the current position.
	// List entities first.
	entities := ui.items

	ui.view.ClearAll()

	// Make sure the active UI is highlighted.
	ts := []t.Transformer{t.Background(colGrey)}
	if ui.isUIActive(ui) {
		ts = append(ts, t.Foreground(concolor.Green))
	}
	ui.view.PrintBounded(1, 0, ui.view.Width-2, 2, fmt.Sprintf("Enemies (%d)", len(entities)), ts...)

	// TODO: Fix the scrolling and unify with other UIs.
	var idx int
	for i := ui.selectedItem; i < len(entities); i++ {
		e := entities[i]
		entry := e.Name
		var transformers []t.Transformer
		if ui.isUIActive(ui) && i == ui.selectedItem {
			transformers = append(transformers, t.Foreground(concolor.Green))
		}
		if e.IsDead() {
			entry += " (dead)"
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
	// Update the selectable list.
	ui.selectableList.handleInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.Select()
	}
	// Show information about the selected enemy (if any).
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		// Display information about the enemy.
		// TODO: Display in a text box.
		if sel := ui.Selected(); sel != nil {
			ui.AddMessage(sel.Name + ": " + sel.Description)
		}
	}
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
	e := ui.Selected()
	if e == nil {
		return
	}
	if e.IsDead() {
		// Loot
		for _, it := range e.Items {
			it.Equipped = false // Mark all items as unequipped.
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
	view *console.Console
	selectableList[Item]
	ItemOpen *Item // The item that is currently open.
}

func (g *Game) newPlayerItems() (*uiItems, error) {
	playerItemsView, err := g.sideView.CreateSubConsole(0, 23, 20, 8)
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
	// Update the selectable list.
	ui.selectableList.setItems(ui.inRange())

	// Draw what can be found at the current position.
	// List items first.
	items := ui.items

	ui.view.ClearAll()
	headline := fmt.Sprintf("Items (%d)", len(items))
	if ui.ItemOpen != nil {
		headline = fmt.Sprintf("Items %s (%d)", ui.ItemOpen.Name, len(items))
	}
	// Make sure the active UI is highlighted.
	ts := []t.Transformer{t.Background(colGrey)}
	if ui.isUIActive(ui) {
		ts = append(ts, t.Foreground(concolor.Green))
	}
	ui.view.PrintBounded(1, 0, ui.view.Width-2, 2, headline, ts...)

	var idx int
	start, end := calcVisibleRange(ui.view.Height-2, len(ui.inRange()), ui.selectedItem)
	for i := start; i < end; i++ {
		it := items[i]
		var transformers []t.Transformer
		if ui.isUIActive(ui) && i == ui.selectedItem {
			transformers = append(transformers, t.Foreground(concolor.Green))
		}
		ui.view.PrintBounded(1, 1+idx, ui.view.Width-2, 2, fmt.Sprintf("%d: %s", i, it.Name), transformers...)
		idx++
	}
}

func (ui *uiItems) HandleInput() {
	// Update the selectable list.
	ui.selectableList.handleInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.Select()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		// Close the item if it is open.
		if ui.ItemOpen != nil {
			ui.ItemOpen = nil
		}
	}

	// Show information about the selected item (if any).
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		// Display information about the item.
		// TODO: Display in a text box.
		if sel := ui.Selected(); sel != nil {
			ui.AddMessage(sel.Name + ": " + sel.Description)
		}
	}
}

func (ui *uiItems) inRange() []*Item {
	if ui.ItemOpen != nil {
		// If the open item is in range, return its contents.
		if ui.ItemOpen.X == ui.player.X && ui.ItemOpen.Y == ui.player.Y {
			return ui.ItemOpen.Contains
		}

		// If the item is not in range, close it.
		ui.ItemOpen = nil
	}

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
	it := ui.Selected()
	if it == nil {
		return
	}

	// If the item is a container, open it.
	// TODO: Store parent for nested containers.
	if it.Type == ItemTypeContainer || len(it.Contains) > 0 {
		// Open container.
		ui.ItemOpen = it
		return
	}

	// Add the item to the inventory.
	ui.player.Inventory.Add(it)

	// TODO: Remove the item from the container if one is open.
	if ui.ItemOpen != nil {
		// Remove the item from the container.
		ui.ItemOpen.Contains = removeItemFromSlice(ui.ItemOpen.Contains, it)
	} else {
		// Remove the item from the world.
		ui.Items = removeItemFromSlice(ui.Items, it)
	}
}

func removeItemFromSlice(items []*Item, item *Item) []*Item {
	// Find the actual index and remove the item.
	for i, it := range items {
		if it == item {
			return append(items[:i], items[i+1:]...)
		}
	}
	return items
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
	pE := ui.World.Elevation[pY][pX]

	// Draw player info.
	// Make sure the active UI is highlighted.
	ts := []t.Transformer{t.Background(colGrey)}
	if ui.isUIActive(ui) {
		ts = append(ts, t.Foreground(concolor.Green))
	}
	ui.view.PrintBounded(1, 1, ui.view.Width, 1, "Player: "+ui.player.Name, ts...)
	ui.view.PrintBounded(1, 2, ui.view.Width-2, 2, fmt.Sprintf("Health: %d/%d", ui.player.Health, ui.player.BaseHealth))
	ui.view.PrintBounded(1, 3, ui.view.Width-2, 2, fmt.Sprintf("Def: %d Att: %d", ui.player.DefenseValue(), ui.player.AttackDamage()))
	ui.view.PrintBounded(1, 4, ui.view.Width-2, 2, fmt.Sprintf("X=%d Y=%d E=%d", pX, pY, pE), t.Foreground(colGrey))
}

func (ui *uiPlayerInfo) HandleInput() {
}

func (ui *uiPlayerInfo) Select() {
}

type selectableList[T any] struct {
	selectedItem int
	items        []*T
}

func (ui *selectableList[T]) setItems(items []*T) {
	ui.items = items
	ui.SelectItem(ui.selectedItem)
}

// SelectItem selects the given item index while clamping
// it to the items in range.
func (ui *selectableList[T]) SelectItem(index int) {
	items := len(ui.items)
	if items == 0 {
		index = 0
	} else if index < 0 {
		index = 0
	} else if index >= items {
		index = items - 1
	}
	ui.selectedItem = index
}

func (ui *selectableList[T]) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ui.SelectItem(ui.selectedItem - 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ui.SelectItem(ui.selectedItem + 1)
	}
}

// Selected returns the currently selected item.
func (ui *selectableList[T]) Selected() *T {
	items := ui.items
	// Verify selectedItem is in range.
	if len(items) == 0 || ui.selectedItem < 0 || ui.selectedItem >= len(items) {
		return nil
	}
	return items[ui.selectedItem]
}

// NOTE; This doesn't work well if an item takes more than one line.
func calcVisibleRange(numVisible, numItems, selectedIdx int) (int, int) {
	if numItems < numVisible {
		return 0, numItems
	}
	// Get the selected item index and make sure it is visible.
	sel := selectedIdx
	if sel < 0 {
		sel = 0
	}
	if sel >= numItems {
		sel = numItems - 1
	}

	// Calculate the start and end index.
	start := sel - numVisible/2
	end := start + numVisible
	if start < 0 {
		start = 0
		end = numVisible
	}
	if end > numItems {
		end = numItems
		start = end - numVisible
	}
	return start, end
}
