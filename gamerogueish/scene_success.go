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

	con *console.Console
	tb  *console.Console
}

func NewSceneSuccess(rootView *console.Console, world *Game) *SceneSuccess {
	g := &SceneSuccess{
		ComponentBase: console.NewComponentBase(10, 10, 10, 10),
		Game:          world,
		con:           rootView,
	}
	g.displayText(
		"You finally realize that the person you love most in the world is yourself! Congrats, you win!", g.con)
	return g
}

func (g *SceneSuccess) Update(con *console.Console, timeElapsed float64) bool {
	// Logic for updating the scene.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// If we have an open textbox, close it.
		g.removeText(con)
		g.reset()
		g.setViewMode(ViewModeMap)
	}

	return true
}

func (g *SceneSuccess) Draw(con *console.Console, timeElapsed float64) {
	con.ClearAll()
	con.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))
}

func (g *SceneSuccess) removeText(con *console.Console) {
	if g.tb != nil {
		con.RemoveSubConsole(g.tb)
		g.tb = nil
	}
}

func (g *SceneSuccess) displayText(txt string, con *console.Console) {
	// If we have an open textbox, close it.
	g.removeText(con)

	// Create a new textbox.
	// Center the console within parent.
	conWidth := con.Width
	conHeight := con.Height
	boxWidth := 30
	boxHeight := 20
	boxX := conWidth/2 - boxWidth/2
	boxY := conHeight/2 - boxHeight/2
	textBox, err := con.CreateSubConsole(boxX, boxY, boxWidth, boxHeight)
	if err != nil {
		panic(err)
	}
	textBox.TransformAll(t.Background(concolor.RGB(50, 50, 50)), t.Char(0))
	textBox.PrintBounded(1, 1, boxWidth-2, boxHeight-2, insertLineBreaks(txt, boxWidth-2), t.Foreground(concolor.White))

	// Print the close message.
	textBox.PrintBounded(1, boxHeight-2, boxWidth-2, boxHeight-2, "Press SPACE to restart", t.Foreground(concolor.White))
	g.tb = textBox
}

func (s *SceneSuccess) FocusOnClick() bool { return false }
