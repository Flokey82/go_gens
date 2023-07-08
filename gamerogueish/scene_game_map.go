package gamerogueish

import (
	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type SceneMap struct {
	*Game
	*console.ComponentBase
	turnTaken bool

	con *console.Console
	tb  *console.Console
}

func NewSceneMap(rootView *console.Console, world *Game) *SceneMap {
	g := &SceneMap{
		Game:          world,
		ComponentBase: console.NewComponentBase(0, 0, rootView.Height, rootView.Width),
		con:           rootView,
	}
	g.displayText(
		"This is the story of "+g.player.Name+", a warrior with a small brain but a big heart. "+
			"Years of intense aerobics had left him with cardiomegaly and a passion for sweatbands. "+
			"Now he is looking for love in the most unlikely of places... underground!", g.con)
	return g
}

func (g *SceneMap) Update(con *console.Console, timeElapsed float64) bool {
	// If we have an open textbox, don't do anything else.
	if g.tb != nil {
		// TODO: This doesn't prevent interactions with the inventory.
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.removeText(con)
		}
		return true
	}

	// TODO: Implement a proper turn based gameloop
	// See: http://journal.stuffwithstuff.com/2014/07/15/a-turn-based-game-loop/
	// Check if a turn has been taken.
	if g.turnTaken {
		g.turnTaken = false
		// Check for all items in range that might trigger something.
		for _, item := range g.World.Items {
			if item.X == g.player.X && item.Y == g.player.Y {
				if item.OnTouch != nil {
					item.OnTouch(g.Game, g.player, item)
				}
			}
		}

		// TODO: Find a better place to do this.
		// ... or rename / split this function.
		if g.player.IsDead() {
			g.setViewMode(ViewModeDeath)
		}
	}

	//if g.IsFocused() {
	//	log.Println("SceneMap is focused")
	//}

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
				g.player.Attack(g.Game, en)
				turnTaken = true // A successful attack counts as a turn.
			}
		}
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

		// Make sure that next turn we know that a turn has been taken.
		g.turnTaken = true
	}
	return true
}

func (g *SceneMap) removeText(con *console.Console) {
	if g.tb != nil {
		con.RemoveSubConsole(g.tb)
		g.tb = nil
	}
}

func (g *SceneMap) displayText(txt string, con *console.Console) {
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
	textBox.PrintBounded(1, boxHeight-2, boxWidth-2, boxHeight-2, "Press SPACE to close", t.Foreground(concolor.White))
	g.tb = textBox
}

// insertLineBreaks takes a string and a max width and inserts line breaks so that words don't get cut off.
func insertLineBreaks(txt string, maxWidth int) string {
	var result string
	var line string
	for _, word := range splitWords(txt) {
		if len(line)+len(word) >= maxWidth {
			result += line + "\n"
			line = ""
		}
		if line != "" {
			line += " "
		}
		line += word
	}
	result += line
	return result
}

func splitWords(txt string) []string {
	var result []string
	var word string
	for _, char := range txt {
		if char == ' ' {
			result = append(result, word)
			word = ""
		} else {
			word += string(char)
		}
	}
	result = append(result, word)
	return result
}

func (g *SceneMap) Draw(con *console.Console, timeElapsed float64) {
	// Clear world view.
	g.worldView.ClearAll()
	g.worldView.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))

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

	// Draw items.
	for _, it := range g.Items {
		// Draw only if we can see the items.
		if !g.IsInRadius(pX, pY, it.X, it.Y) {
			continue
		}
		g.worldView.Transform(midX-pX+it.X, midY-pY+it.Y, t.CharByte(it.Tile), t.Foreground(concolor.Green))
	}

	// Draw player in the middle.
	g.worldView.Transform(midX, midY, t.CharByte(g.player.Tile), t.Foreground(concolor.Green))
}

func (g *SceneMap) FocusOnClick() bool { return true }
