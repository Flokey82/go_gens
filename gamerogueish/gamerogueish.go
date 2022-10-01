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

const (
	labelWindow     = "ramen - roguelike example"
	labelWorldView  = "World View"
	labelPlayerInfo = "Player Info"
)

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) *Position {
	return &Position{X: x, Y: y}
}

type Game struct {
	player         *Position        // player entity
	cells          [][]byte         // world cells
	rootView       *console.Console // view for all sub views
	worldView      *console.Console // contains map
	playerInfoView *console.Console // contains player info
}

func NewGame() (*Game, error) {
	g := &Game{
		player: NewPosition(3, 3),
	}
	if err := g.init(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) init() error {
	// Init world cells.
	wWidth := 50
	wHeight := 50
	w := genFancyWorld(wWidth, wHeight)
	g.cells = [][]byte(*w)

	// Place the player in the middle.
	g.player.X = wWidth / 2
	g.player.Y = wHeight / 2

	// Init views / UI.
	rootView, err := console.New(60, 35, font.DefaultFont, labelWindow)
	if err != nil {
		return err
	}
	g.rootView = rootView

	worldView, err := rootView.CreateSubConsole(0, 1, rootView.Width-20, rootView.Height-1)
	if err != nil {
		return err
	}
	g.worldView = worldView

	playerInfoView, err := rootView.CreateSubConsole(worldView.Width, 1, 20, rootView.Height-1)
	if err != nil {
		return err
	}
	g.playerInfoView = playerInfoView

	return nil
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
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && g.canMoveTo(g.player.X, g.player.Y-1) {
		g.player.Y -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) && g.canMoveTo(g.player.X, g.player.Y+1) {
		g.player.Y += 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) && g.canMoveTo(g.player.X-1, g.player.Y) {
		g.player.X -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) && g.canMoveTo(g.player.X+1, g.player.Y) {
		g.player.X += 1
	}

	return nil
}

// checks if a tile is solid (tile content is not a space ' ' character)
func (g *Game) isSolid(x int, y int) bool {
	return g.cells[y][x] != ' '
}

func (g *Game) canMoveTo(x, y int) bool {
	return g.cells[y][x] == ' '
}

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
	for y := range g.cells {
		for x := range g.cells[y] {
			if g.cells[y][x] == ' ' {
				continue
			}
			g.worldView.Transform(midX-g.player.X+x, midY-g.player.Y+y, t.CharByte(g.cells[y][x]))
		}
	}

	// draw player in the middle
	g.worldView.Transform(midX, midY, t.CharByte('@'), t.Foreground(concolor.Green))

	// draw player info
	g.playerInfoView.PrintBounded(1, 1, g.playerInfoView.Width-2, 2, fmt.Sprintf("X=%d Y=%d", g.player.X, g.player.Y))

	return nil
}
