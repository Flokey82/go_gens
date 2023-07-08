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

	*textBox
}

func NewSceneSuccess(rootView *console.Console, world *Game) *SceneSuccess {
	g := &SceneSuccess{
		ComponentBase: console.NewComponentBase(10, 10, 10, 10),
		Game:          world,
		textBox:       NewTextbox(rootView, 30, 20),
	}
	g.displayText(
		"You finally realize that the person you love most in the world is yourself! Congrats, you win! "+
			"Still a better love story than Twilight.", "Press SPACE to restart")
	return g
}

func (g *SceneSuccess) Update(con *console.Console, timeElapsed float64) bool {
	// Logic for updating the scene.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// If we have an open textbox, close it.
		g.removeText()
		g.reset()
		g.setViewMode(ViewModeMap)
	}

	return true
}

func (g *SceneSuccess) Draw(con *console.Console, timeElapsed float64) {
	con.ClearAll()
	con.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))
}

func (s *SceneSuccess) FocusOnClick() bool { return false }

type textBox struct {
	con        *console.Console
	tb         *console.Console
	width      int
	height     int
	background concolor.Color
}

func NewTextbox(con *console.Console, width, height int) *textBox {
	return &textBox{
		con:        con,
		width:      width,
		height:     height,
		background: concolor.RGB(50, 50, 50),
	}
}

func (g *textBox) removeText() {
	if g.tb != nil {
		g.con.RemoveSubConsole(g.tb)
		g.tb = nil
	}
}

func (g *textBox) displayText(txt, bottomStr string) {
	// If we have an open textbox, close it.
	g.removeText()

	// Create a new textbox.
	// Center the console within parent.
	conWidth := g.con.Width
	conHeight := g.con.Height
	boxWidth := g.width
	boxHeight := g.height
	boxX := conWidth/2 - boxWidth/2
	boxY := conHeight/2 - boxHeight/2
	textBox, err := g.con.CreateSubConsole(boxX, boxY, boxWidth, boxHeight)
	if err != nil {
		panic(err)
	}
	textBox.TransformAll(t.Background(g.background), t.Char(0))
	textBox.PrintBounded(1, 1, boxWidth-2, boxHeight-2, insertLineBreaks(txt, boxWidth-2), t.Foreground(concolor.White))

	// Print the close message.
	textBox.PrintBounded(1, boxHeight-2, boxWidth-2, boxHeight-2, bottomStr, t.Foreground(concolor.White))
	g.tb = textBox
}
