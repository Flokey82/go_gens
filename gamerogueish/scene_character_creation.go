package gamerogueish

import (
	"log"

	"github.com/BigJk/ramen/components"
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type SceneCharacterCreation struct {
	*Game
	*console.ComponentBase
	tb *components.TextBox
}

func NewSceneCharacterCreation(rootView *console.Console, world *Game) *SceneCharacterCreation {
	return &SceneCharacterCreation{
		Game:          world,
		ComponentBase: console.NewComponentBase(10, 10, 10, 10),
	}
}

func (s *SceneCharacterCreation) NewTextbox(con *console.Console) {
	textBox := components.NewTextbox(7, 3, 20, 1)
	textBox.SetBackground(concolor.RGB(50, 50, 50).P(), concolor.RGB(70, 70, 70).P(), concolor.RGB(30, 30, 30).P())
	textBox.SetForeground(concolor.White.P(), concolor.RGB(90, 90, 90).P())
	textBox.SetText(s.player.Name)
	textBox.SetEnterCallback(func(text string) {
		s.player.Name = text
		log.Println("Name set to", text)
		textBox.Close()
		s.tb = nil
		con.RemoveComponent(textBox)
	})
	s.tb = textBox
	log.Println("New textbox created")
	con.AddComponent(textBox)
}

func (s *SceneCharacterCreation) Update(con *console.Console, timeElapsed float64) bool {
	// Logic for updating the scene.
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		s.NewTextbox(con)
		s.tb.SetFocus(true)
	}

	if s.tb != nil && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.setViewMode(ViewModeMap)
	}

	return true
}

func (s *SceneCharacterCreation) Draw(con *console.Console, timeElapsed float64) {
	if s.tb != nil {
		return
	}
	con.ClearAll()
	con.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))

	// Logic for drawing the scene.
	con.Print(1, 1, "Character Creation", t.Foreground(concolor.White))
	con.Print(1, 3, "Name: "+s.player.Name, t.Foreground(concolor.White))
	con.Print(1, 4, "Press TAB to change name", t.Foreground(concolor.White))
	con.Print(1, 5, "Press ENTER to start game", t.Foreground(concolor.White))
}

func (g *SceneCharacterCreation) Close() error {
	if g.tb != nil {
		g.tb.Close()
	}
	return nil
}

func (s *SceneCharacterCreation) FocusOnClick() bool { return false }
