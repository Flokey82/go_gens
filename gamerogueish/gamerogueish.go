package gamerogueish

import (
	"fmt"

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
	*World                          // currently generated world
	*FOV                            // currently generated FOV
	generator      GenWorld         // world generator function
	player         *Entity          // player entity
	rootView       *console.Console // view for all sub views
	worldView      *console.Console // contains map
	playerInfoView *console.Console // contains player info
}

func NewGame(gw GenWorld, width, height int, seed int64) (*Game, error) {
	g := &Game{
		generator: gw,
		World:     gw(width, height, seed),
		player:    NewEntity(width/2, height/2, '@'), // Place the player in the middle.
	}

	g.player.Inventory.Items = append(g.player.Inventory.Items, Item{
		Name: "Sword",
		Type: ItemWeapon,
	})
	g.player.Inventory.Items = append(g.player.Inventory.Items, Item{
		Name: "Axe",
		Type: ItemWeapon,
	})
	g.player.Inventory.Items = append(g.player.Inventory.Items, Item{
		Name: "Potion",
		Type: ItemPotion,
	})
	g.player.Inventory.Items = append(g.player.Inventory.Items, Item{
		Name: "Leather Armor",
		Type: ItemArmor,
	})
	g.player.Inventory.Items = append(g.player.Inventory.Items, Item{
		Name: "Plate Armor",
		Type: ItemArmor,
	})

	g.Entities = append(g.Entities, NewEntity(10, 10, 'R'))

	g.FOV = NewFOV(g.World, 10)
	g.FOV.Update(g.player.X, g.player.Y) // Update FOV

	// Init views / UI.
	rootView, err := console.New(60, 35, font.DefaultFont, labelWindow)
	if err != nil {
		return nil, err
	}
	g.rootView = rootView

	worldView, err := rootView.CreateSubConsole(0, 1, rootView.Width-20, rootView.Height-1)
	if err != nil {
		return nil, err
	}
	g.worldView = worldView

	playerInfoView, err := rootView.CreateSubConsole(worldView.Width, 1, 20, rootView.Height-1)
	if err != nil {
		return nil, err
	}
	g.playerInfoView = playerInfoView

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

	// TODO: Move this to a UI component.
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.player.Inventory.SelectItem(g.player.Inventory.selectedItem - 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.player.Inventory.SelectItem(g.player.Inventory.selectedItem + 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.player.Equip(g.player.Inventory.selectedItem)
	}

	if turnTaken {
		// If we move, update seen tiles.
		g.Compute(g.player.X, g.player.Y)

		// Handle entity AI.
	}

	return nil
}

var colGrey = concolor.RGB(128, 128, 128)

func (g *Game) Update(screen *ebiten.Image, timeDelta float64) error {
	// clear console
	g.rootView.ClearAll()
	g.rootView.TransformAll(t.Background(concolor.RGB(50, 50, 50)))

	g.worldView.ClearAll()
	g.worldView.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))

	g.playerInfoView.ClearAll()

	// draw header
	g.rootView.TransformArea(0, 0, g.rootView.Width, 1, t.Background(concolor.RGB(80, 80, 80)))
	g.rootView.Print(2, 0, labelWorldView, t.Foreground(concolor.White))
	g.rootView.Print(g.worldView.Width+2, 0, labelPlayerInfo, t.Foreground(concolor.White))

	// draw world
	midX := g.worldView.Width / 2
	midY := g.worldView.Height / 2

	// TODO: Skip drawing everything outside of the view.
	for y := range g.Cells {
		for x := range g.Cells[y] {
			if g.Cells[y][x] == ' ' || !g.Seen[y][x] {
				continue
			}
			if !g.IsInRadius(g.player.X, g.player.Y, x, y) {
				g.worldView.Transform(midX-g.player.X+x, midY-g.player.Y+y, t.CharByte(g.Cells[y][x]), t.Foreground(colGrey))
			} else {
				g.worldView.Transform(midX-g.player.X+x, midY-g.player.Y+y, t.CharByte(g.Cells[y][x]))
			}
		}
	}

	// draw player in the middle
	g.worldView.Transform(midX, midY, t.CharByte(g.player.Tile), t.Foreground(concolor.Green))

	// draw entities
	for _, e := range g.Entities {
		if !g.IsInRadius(g.player.X, g.player.Y, e.X, e.Y) {
			continue
		}
		g.worldView.Transform(midX-g.player.X+e.X, midY-g.player.Y+e.Y, t.CharByte(e.Tile), t.Foreground(concolor.Blue))
	}

	// draw player info
	g.playerInfoView.PrintBounded(1, 1, g.playerInfoView.Width-2, 2, fmt.Sprintf("X=%d Y=%d", g.player.X, g.player.Y))

	// draw inventory
	// TODO:
	// - Move this to a UI component.
	// - Render equipped armor and weapon.
	g.playerInfoView.PrintBounded(1, 3, g.playerInfoView.Width-2, 2, fmt.Sprintf("Inventory (%d)", g.player.Inventory.Count()))
	for i, item := range g.player.Items {
		var entry string
		if item.Equipped {
			entry = fmt.Sprintf("%d:*%s", i, item.Name)
		} else {
			entry = fmt.Sprintf("%d: %s", i, item.Name)
		}
		var transformers []t.Transformer
		if i == g.player.selectedItem {
			transformers = append(transformers, t.Foreground(concolor.Green))
		}
		g.playerInfoView.PrintBounded(2, 5+i, g.playerInfoView.Width-2, 2, entry, transformers...)
	}
	return nil
}
