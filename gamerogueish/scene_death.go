package gamerogueish

import (
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type SceneDeath struct {
	*Game
	*console.ComponentBase
	*textBox
}

func NewSceneDeath(rootView *console.Console, world *Game) *SceneDeath {
	g := &SceneDeath{
		Game:          world,
		ComponentBase: console.NewComponentBase(10, 10, 10, 10),
		textBox:       NewTextbox(rootView, 32, 20),
	}
	g.displayText("You died.", "Press ESC to restart game")
	return g
}

func (s *SceneDeath) Update(con *console.Console, timeElapsed float64) bool {
	// Logic for updating the scene.
	// TODO: Move this to the textbox.
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.removeText()
		s.reset()
		s.setViewMode(ViewModeMap)
	}

	return true
}

func (s *SceneDeath) Draw(con *console.Console, timeElapsed float64) {
	con.ClearAll()
	con.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))
}

func (g *SceneDeath) Close() error {
	g.textBox.removeText()
	return nil
}

func (s *SceneDeath) FocusOnClick() bool { return false }
