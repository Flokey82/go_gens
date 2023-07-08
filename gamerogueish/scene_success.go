package gamerogueish

import (
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type SceneSuccess struct {
	*console.ComponentBase
	*Game
}

func NewSceneSuccess(rootView *console.Console, world *Game) *SceneSuccess {
	return &SceneSuccess{
		ComponentBase: console.NewComponentBase(10, 10, 10, 10),
		Game:          world,
	}
}

func (s *SceneSuccess) Update(con *console.Console, timeElapsed float64) bool {
	// Logic for updating the scene.
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.reset()
		s.setViewMode(ViewModeMap)
	}

	return true
}

func (s *SceneSuccess) Draw(con *console.Console, timeElapsed float64) {
	con.ClearAll()
	con.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))

	// Logic for drawing the scene.
	con.Print(1, 1, "You won! Press ESC to restart game", t.Foreground(concolor.White))
}

func (s *SceneSuccess) FocusOnClick() bool { return false }
