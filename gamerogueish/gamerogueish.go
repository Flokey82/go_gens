package gamerogueish

import (
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type GenWorld func(width, height int, seed int64) *World

const (
	labelWindow     = "rogue-ish"
	labelWorldView  = "World View"
	labelPlayerInfo = "Player Info"
)

type Game struct {
	*World                       // currently generated world
	*FOV                         // currently generated FOV
	generator   GenWorld         // world generator function
	player      *Entity          // player entity
	rootView    *console.Console // view for all sub views
	worldView   *console.Console // contains map
	sideView    *console.Console
	sideViews   []UIif           // contains all side views
	selectedUI  int              // currently selected UI
	messageView *console.Console // contains messages
	Messages    []string         // messages to display
}

func NewGame(gw GenWorld, width, height int, seed int64) (*Game, error) {
	g := &Game{
		generator: gw,
		World:     gw(width, height, seed),
		player:    NewEntity(width/2, height/2, EntityPlayer), // Place the player in the middle.
	}

	// Seed the player inventory with some items.
	// NOTE: This is just for testing purposes.
	g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypeWeaponSword.New())
	g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypeWeaponAxe.New())
	g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypePotion.New())
	g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypePotion.New())
	g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypePotion.New())
	g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypeArmorLeather.New())
	g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypeArmorPlate.New())

	// Set up the FOV.
	g.FOV = NewFOV(g.World, 10)
	g.FOV.Update(g.player.X, g.player.Y) // Update FOV

	// Init views / UI.
	rootView, err := console.New(60, 35, font.DefaultFont, labelWindow)
	if err != nil {
		return nil, err
	}
	g.rootView = rootView

	worldView, err := rootView.CreateSubConsole(0, 1, rootView.Width-20, rootView.Height-4)
	if err != nil {
		return nil, err
	}
	g.worldView = worldView

	sideView, err := rootView.CreateSubConsole(worldView.Width, 1, 20, rootView.Height-4)
	if err != nil {
		return nil, err
	}
	g.sideView = sideView

	// Draw player info.
	uiInfo, err := g.newPlayerInfo()
	if err != nil {
		return nil, err
	}
	g.sideViews = append(g.sideViews, uiInfo)

	// Draw inventory.
	uiInventory, err := g.newPlayerInventory()
	if err != nil {
		return nil, err
	}
	g.sideViews = append(g.sideViews, uiInventory)

	// Draw enemies.
	uiEnemy, err := g.newPlayerEnemies()
	if err != nil {
		return nil, err
	}
	g.sideViews = append(g.sideViews, uiEnemy)

	messageView, err := rootView.CreateSubConsole(0, rootView.Height-3, rootView.Width, 3)
	if err != nil {
		return nil, err
	}
	g.messageView = messageView

	return g, nil
}

func (g *Game) Start() {
	// Setup input.
	g.rootView.SetTickHook(g.HandleInput)

	// Set up renderer.
	g.rootView.SetPreRenderHook(g.Update)

	// Start!
	g.rootView.Start(2)
}

func (g *Game) HandleInput(timeElapsed float64) error {
	// TODO: Implement a proper turn based gameloop
	// See: http://journal.stuffwithstuff.com/2014/07/15/a-turn-based-game-loop/

	// Player movement.
	var turnTaken bool
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && g.CanMoveTo(g.player.X, g.player.Y-1) {
		g.player.Y -= 1
		turnTaken = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) && g.CanMoveTo(g.player.X, g.player.Y+1) {
		g.player.Y += 1
		turnTaken = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) && g.CanMoveTo(g.player.X-1, g.player.Y) {
		g.player.X -= 1
		turnTaken = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) && g.CanMoveTo(g.player.X+1, g.player.Y) {
		g.player.X += 1
		turnTaken = true
	}

	// Attack entities if space is pressed.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		pX := g.player.X
		pY := g.player.Y

		// TODO: Factor this out into a function.
		var entities []*Entity
		for _, e := range g.Entities {
			if e.X == pX && e.Y == pY {
				entities = append(entities, e)
			}
		}

		// Select Random enemy.
		for _, en := range entities {
			if !en.IsDead() {
				g.player.Attack(g, en)
			}
		}
		turnTaken = true
	}

	// Inventory stuff.
	// TODO: On TAB, cycle through the UI elements.
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.selectedUI++
		g.selectedUI %= len(g.sideViews)
	}

	// TODO: Move this to a UI component.
	if ui := g.getCurrentActiveUI(); ui != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			ui.Prev()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			ui.Next()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			ui.Select()
		}
	}

	// For dev purposes we add a key to add potions.
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypePotion.New())
	}

	// If we have taken an actual turn, refresh the FOV
	// and allow the creatures to take their turn.
	if turnTaken {
		// If we move, update seen tiles.
		g.Compute(g.player.X, g.player.Y)

		// Handle entity AI.
		// TODO: Make this energy based or something.
		for _, e := range g.Entities {
			g.decideAction(e)
		}
	}

	return nil
}

var (
	colGrey    = concolor.RGB(128, 128, 128)
	colDarkRed = concolor.RGB(128, 0, 0)
)

func (g *Game) Update(screen *ebiten.Image, timeDelta float64) error {
	// Clear console.
	g.rootView.ClearAll()
	g.rootView.TransformAll(t.Background(concolor.RGB(50, 50, 50)))

	// Draw world.
	g.drawMap()

	// Draw side menu.
	g.drawSideMenu()

	// List messages.
	g.drawMessages()

	return nil
}

func (g *Game) drawMap() {
	g.worldView.ClearAll()
	g.worldView.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))

	// Draw header.
	g.rootView.TransformArea(0, 0, g.rootView.Width, 1, t.Background(concolor.RGB(80, 80, 80)))
	g.rootView.Print(2, 0, labelWorldView, t.Foreground(concolor.White))
	g.rootView.Print(g.worldView.Width+2, 0, labelPlayerInfo, t.Foreground(concolor.White))

	// Draw world centered around the player.
	midX := g.worldView.Width / 2
	midY := g.worldView.Height / 2

	// Player position.
	pX := g.player.X
	pY := g.player.Y

	// TODO: Skip drawing everything outside of the view.
	for y := range g.Cells {
		for x, cv := range g.Cells[y] {
			// Skip empty cells and cells we haven't seen.
			if cv == ' ' || !g.Seen[y][x] {
				continue
			}

			// Previously seen tiles that we can't see right now are greyed out.
			if !g.IsInRadius(pX, pY, x, y) {
				g.worldView.Transform(midX-pX+x, midY-pY+y, t.CharByte(cv), t.Foreground(colGrey))
			} else {
				g.worldView.Transform(midX-pX+x, midY-pY+y, t.CharByte(cv))
			}
		}
	}

	// draw player in the middle
	g.worldView.Transform(midX, midY, t.CharByte(g.player.Tile), t.Foreground(concolor.Green))

	// Draw entities.
	for _, e := range g.Entities {
		// Draw only if we can see the creatures.
		if !g.IsInRadius(pX, pY, e.X, e.Y) {
			continue
		}
		transformer := t.Foreground(concolor.Red)
		if e.IsDead() {
			transformer = t.Foreground(colDarkRed)
		}
		g.worldView.Transform(midX-pX+e.X, midY-pY+e.Y, t.CharByte(e.Tile), transformer)
	}
}

func (g *Game) isUIActive(ui UIif) bool {
	return g.getCurrentActiveUI() == ui
}

func (g *Game) getCurrentActiveUI() UIif {
	if g.selectedUI < 0 || g.selectedUI >= len(g.sideViews) {
		return nil
	}
	return g.sideViews[g.selectedUI]
}

func (g *Game) drawSideMenu() {
	for _, ui := range g.sideViews {
		ui.Draw()
	}
}

func (g *Game) drawMessages() {
	g.messageView.ClearAll()
	for i, m := range g.Messages {
		g.messageView.PrintBounded(1, i, g.messageView.Width-2, 2, m)
	}
}

func (g *Game) AddMessage(msg string) {
	// TODO: Move this to a messaging component.
	const maxMessages = 3
	g.Messages = append(g.Messages, msg)
	if len(g.Messages) > maxMessages {
		g.Messages = g.Messages[1:]
	}
}
