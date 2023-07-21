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
	*textBox
	turnTaken bool
}

func NewSceneMap(rootView *console.Console, world *Game) *SceneMap {
	g := &SceneMap{
		Game:          world,
		ComponentBase: console.NewComponentBase(0, 0, rootView.Height, rootView.Width),
		textBox:       NewTextbox(rootView, 32, 20),
	}
	g.displayText(
		"This is the story of "+g.player.Name+", a warrior with a small brain but a big heart. "+
			"Years of intense aerobics had left him with cardiomegaly and a passion for sweatbands. "+
			"Now he is looking for love in the most unlikely of places... \n \n"+
			"... underground!", "Press SPACE to close")
	return g
}

func (g *SceneMap) Update(con *console.Console, timeElapsed float64) bool {
	// If we have an open textbox, don't do anything else.
	if g.textBox.isOpen() {
		// TODO: Move this to the textbox.
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.removeText()
		}
		// TODO: This doesn't prevent interactions with the inventory.
		g.textBox.handleInput()
		return true
	}

	// TODO: Implement a proper turn based gameloop
	// See: http://journal.stuffwithstuff.com/2014/07/15/a-turn-based-game-loop/
	// Check if a turn has been taken.
	if g.turnTaken {
		g.turnTaken = false
		// Check for all items in range that might trigger something.
		for _, item := range g.World.Items {
			// TODO: Deduplicate this code and make it generic.
			if item.X == g.player.X && item.Y == g.player.Y {
				if item.OnTouch != nil {
					item.OnTouch(g.Game, g.player, item)
				}
			}
			// Do the same for all entities that are in range.
			for _, e := range g.Entities {
				// This might a bit cheaty, but we don't want to have all entities
				// kill themselves while wandering around.
				// TODO: Make this more robust and less hacky.
				if e.IsDead() || !g.IsInRadius(g.player.X, g.player.Y, e.X, e.Y) {
					continue
				}
				if e.X == item.X && e.Y == item.Y {
					if item.OnTouch != nil {
						item.OnTouch(g.Game, e, item)
					}
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

	// Calculate the bounds of the view.
	minX := pX - midX
	if minX < 0 {
		minX = 0
	}
	maxX := pX + midX
	if maxX >= g.World.Width {
		maxX = g.World.Width
	}
	minY := pY - midY
	if minY < 0 {
		minY = 0
	}
	maxY := pY + midY
	if maxY >= g.World.Height {
		maxY = g.World.Height
	}

	// Draw everything in the view.
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			cv := g.World.Cells[y][x]
			// Skip empty cells and cells we haven't seen.
			if cv == CharFloor || !g.Seen[y][x] {
				continue
			}

			// Previously seen tiles that we can't see right now are greyed out.
			// TODO: Move this somewhere else.
			var col concolor.Color
			inRadius := g.IsInRadius(pX, pY, x, y)
			if !inRadius {
				col = colGrey
			} else {
				switch cv {
				case CharWall, CharColumn: // Wall or column
					col = concolor.White
				case CharWater: // Water
					col = concolor.Blue
				case CharTree: // Tree
					col = concolor.Green
				}
			}

			g.worldView.Transform(midX-pX+x, midY-pY+y, t.CharRune(cv), t.Foreground(col))
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
		// TODO: Reveal hidden items if we have the ability to see them.
		if !g.IsInRadius(pX, pY, it.X, it.Y) || it.Hidden {
			continue
		}
		g.worldView.Transform(midX-pX+it.X, midY-pY+it.Y, t.CharByte(it.Tile), t.Foreground(concolor.Green))
	}

	// Draw player in the middle.
	g.worldView.Transform(midX, midY, t.CharByte(g.player.Tile), t.Foreground(concolor.Green))
}

func (g *SceneMap) FocusOnClick() bool { return true }
