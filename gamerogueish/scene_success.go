package gamerogueish

import (
	"strings"

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
		textBox:       NewTextbox(rootView, 32, 20),
	}
	g.displayText(
		"You finally realize that the person you love most in the world is yourself! Congrats, you win! "+
			"Still a better love story than Twilight.", "Press SPACE to restart")
	return g
}

func (g *SceneSuccess) Update(con *console.Console, timeElapsed float64) bool {
	// Logic for updating the scene.
	// TODO: Move this to the textbox.
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

func (g *SceneSuccess) Close() error {
	g.textBox.removeText()
	return nil
}

func (s *SceneSuccess) FocusOnClick() bool { return false }

type textBox struct {
	con        *console.Console
	tb         *console.Console
	width      int
	height     int
	background concolor.Color
	text       string // text to display
	textLine   int    // line offset
	textBottom string // text to display at the bottom
	margin     int    // margin between text and border
}

func NewTextbox(con *console.Console, width, height int) *textBox {
	return &textBox{
		con:        con,
		width:      width,
		height:     height,
		background: concolor.RGB(50, 50, 50),
		margin:     1,
	}
}

func (g *textBox) isOpen() bool {
	return g.tb != nil
}

func (g *textBox) removeText() {
	if g.isOpen() {
		g.con.RemoveSubConsole(g.tb)
		g.tb = nil
		g.text = ""
		g.textBottom = ""
		g.textLine = 0
	}
}

func (g *textBox) handleInput() bool {
	// TODO: Handle closing the textbox.
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.nextPage()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.prevPage()
	}
	return true
}

func (g *textBox) nextPage() {
	// Get number of lines per page.
	numLines := g.numLinesInBox()

	// Advance the text by one page.
	// Prevent jumping to the last line if we are already on the last page.
	// NOTE: This is a bit hacky.
	if g.textLine+numLines < strings.Count(g.text, "\n") {
		g.jumpToLine(g.textLine + numLines)
	}
}

func (g *textBox) prevPage() {
	// Get number of lines per page.
	numLines := g.numLinesInBox()

	// Rewind the text by one page.
	g.jumpToLine(g.textLine - numLines)
}

func (g *textBox) jumpToLine(lineNr int) {
	g.textLine = lineNr
	if numLines := strings.Count(g.text, "\n"); g.textLine > numLines {
		g.textLine = numLines
	} else if g.textLine < 0 {
		g.textLine = 0
	}
	g.drawText()
}

func (g *textBox) numLinesInBox() int {
	if !g.isOpen() {
		return 0
	}
	return g.tb.Height - 2*g.margin
}

func (g *textBox) drawText() {
	numLines := g.numLinesInBox()
	var curLines int
	var moreText bool
	var newText string
	for i, line := range strings.Split(g.text, "\n") {
		// If we have reached the end of the box,
		// stop and indicate that there is more text.
		if curLines >= numLines-2 { // We skip the last line to make room for the "more text" indicator.
			moreText = true
			break
		}
		if i < g.textLine {
			continue
		}
		newText += line + "\n"
		curLines++
	}

	boxWidth := g.width
	boxHeight := g.height
	textWidth := boxWidth - 2*g.margin
	textHeight := boxHeight - 2*g.margin
	g.tb.TransformAll(t.Background(g.background), t.Char(0))
	g.tb.PrintBounded(g.margin, g.margin, textWidth, textHeight, newText, t.Foreground(concolor.White))

	// Print the close message.
	g.tb.PrintBounded(g.margin, textHeight, textWidth, textHeight, g.textBottom, t.Foreground(concolor.White))

	// Print the "more text" indicator.
	if moreText {
		g.tb.PrintBounded(textWidth, textHeight, textWidth, textHeight, ">", t.Foreground(concolor.Green))
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
	// TODO: Add pagination (maybe per line?)
	g.text = insertLineBreaks(txt, boxWidth-2)
	g.textBottom = bottomStr
	g.tb = textBox
	g.jumpToLine(0)
}

// insertLineBreaks takes a string and a max width and inserts line breaks so that words don't get cut off.
func insertLineBreaks(txt string, maxWidth int) string {
	var result string
	var line string
	// TODO: First, split by line breaks and then by words.
	// This will preserve paragraphs.
	for _, paragraph := range strings.Split(txt, "\n") {
		for _, word := range splitWords(paragraph) {
			if len(line)+len(word) >= maxWidth {
				result += line + "\n"
				line = ""
			}
			if line != "" {
				line += " "
			}
			line += word
		}
		result += line + "\n"
		line = ""
	}
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
