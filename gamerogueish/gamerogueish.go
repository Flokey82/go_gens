package gamerogueish

import (
	"math/rand"

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

type ViewMode int

const (
	ViewModeMap ViewMode = iota
	ViewModeCharacterCreation
	ViewModeDeath
	ViewModeSuccess
	ViewModeMax
)

type Game struct {
	Seed         int64 // seed for the world
	Width        int
	Height       int
	*World                        // currently generated world
	*FOV                          // currently generated FOV
	generator    GenWorld         // world generator function
	player       *Entity          // player entity
	rootView     *console.Console // view for all sub views
	worldView    *console.Console // contains map
	sideView     *console.Console
	sideViews    []UIif           // contains all side views
	selectedUI   int              // currently selected UI
	messageView  *console.Console // contains messages
	Messages     []string         // messages to display
	view         ViewMode
	currentScene Scene
	// *textBox
}

func (g *Game) reset() {
	// Generate a new world.
	seed := g.Seed
	if seed == -1 {
		seed = rand.Int63()
	}

	g.World = g.generator(g.Width, g.Height, seed)

	// Create player.
	g.player = NewEntity(g.World.Width/2, g.World.Height/2, EntityPlayer) // Place the player in the middle.
	g.player.Name = "Glorbnorb"

	// Set up the FOV.
	g.FOV = NewFOV(g.World, 10)
	g.FOV.Update(g.player.X, g.player.Y) // Update FOV

	// Reset messages.
	g.Messages = nil

	// Reset textbox.
	// if g.textBox != nil {
	//	g.textBox.removeText()
	// }
}

// Initializes a new game.
// NOTE: A seed of -1 will generate a random seed each time.
func NewGame(gw GenWorld, width, height int, seed int64) (*Game, error) {
	g := &Game{
		Seed:      seed,
		Width:     width,
		Height:    height,
		generator: gw,
	}

	// Reset the game.
	g.reset()

	// Init views / UI.
	rootView, err := console.New(60, 35, font.DefaultFont, labelWindow)
	if err != nil {
		return nil, err
	}
	g.rootView = rootView

	// Initialize textbox.
	// g.textBox = NewTextbox(g.rootView, 32, 20)

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

	// Draw items.
	uiItems, err := g.newPlayerItems()
	if err != nil {
		return nil, err
	}
	g.sideViews = append(g.sideViews, uiItems)

	// Select inventory as default.
	g.selectedUI = 1

	// Draw messages.
	messageView, err := rootView.CreateSubConsole(0, rootView.Height-3, rootView.Width, 3)
	if err != nil {
		return nil, err
	}
	g.messageView = messageView

	// set default view mode.
	g.setViewMode(ViewModeMap)

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

func (g *Game) setViewMode(vm ViewMode) {
	g.view = vm
	// Unload current scene.
	if g.currentScene != nil {
		g.worldView.RemoveComponent(g.currentScene)
		g.currentScene.Close()
		g.currentScene = nil
	}
	var newScene Scene
	switch vm {
	case ViewModeMap:
		newScene = NewSceneMap(g.worldView, g)
	case ViewModeCharacterCreation:
		newScene = NewSceneCharacterCreation(g.worldView, g)
	case ViewModeDeath:
		newScene = NewSceneDeath(g.worldView, g)
	case ViewModeSuccess:
		newScene = NewSceneSuccess(g.worldView, g)
	}
	g.worldView.AddComponent(newScene)
	g.currentScene = newScene
}

func (g *Game) HandleInput(timeElapsed float64) error {
	// If we have an open textbox, don't do anything else.
	/*
		if g.textBox.isOpen() {
			log.Println("textbox is open")
			// TODO: Move this to the textbox.
			if inpututil.IsKeyJustPressed(ebiten.KeyY) {
				log.Println("Starting new game")
				g.textBox.removeText()
				g.reset()
				g.setViewMode(ViewModeMap)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyN) {
				log.Println("Quitting game")
				g.textBox.removeText()
			}
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				g.removeText()
			}
			// TODO: This doesn't prevent interactions with the inventory.
			g.textBox.handleInput()
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyL) {
			log.Println("showing intro text")
			g.textBox.displayText("You are a brave adventurer. You have been tasked with finding the legendary treasure of the ancient temple of doom. You have been warned that the temple is filled with traps and monsters. Do you want to continue?", "Press Y to continue, N to quit")
		}
	*/

	// TODO: Figure out a better way to handle side view focus.
	// MOVE THIS OUT OF HERE!
	// Inventory stuff.
	// TODO: On TAB, cycle through the UI elements.
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.selectedUI++
		g.selectedUI %= len(g.sideViews)
	}

	// TODO: Move this to a UI component.
	if ui := g.getCurrentActiveUI(); ui != nil {
		ui.HandleInput()
	}

	// For dev purposes we add a key to add potions.
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.player.Inventory.Items = append(g.player.Inventory.Items, ItemTypePotion.New())
	}
	// For dev purposes we add a key to add traps.
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		t := ItemTypeTrap.New()
		t.Hidden = false // We know where the trap is.
		g.player.Inventory.Items = append(g.player.Inventory.Items, t)
	}
	// END: MOVE THIS OUT OF HERE!

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.reset()
		g.setViewMode(ViewModeMap)
	}

	//if inpututil.IsKeyJustPressed(ebiten.KeyI) {
	//	// Cycle through view modes.
	//	log.Println("Cycle view mode", int(g.view), int(ViewModeMax))
	//	g.setViewMode(ViewMode((int(g.view) + 1) % int(ViewModeMax)))
	//}
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

	// Draw header.
	g.rootView.TransformArea(0, 0, g.rootView.Width, 1, t.Background(concolor.RGB(80, 80, 80)))
	g.rootView.Print(2, 0, labelWorldView, t.Foreground(concolor.White))
	g.rootView.Print(g.worldView.Width+2, 0, labelPlayerInfo, t.Foreground(concolor.White))

	// Draw world.
	//if g.currentScene == nil {
	//	g.drawMap()
	//}

	// Draw side menu.
	g.drawSideMenu()

	// List messages.
	g.drawMessages()

	return nil
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
