package gamehex

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type pathFinder struct {
	pathMode  bool     // Path mode.
	pathStart [2]int   // Start of path.
	pathEnd   [2]int   // End of path.
	path      [][2]int // Current path.
}

func newPathFinder() pathFinder {
	return pathFinder{
		pathMode:  false,
		pathStart: [2]int{-1, -1},
		pathEnd:   [2]int{-1, -1},
		path:      nil,
	}
}

func (g *Game) handlePathInput() {
	// If p is pressed, we get ready to calculate a path.
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.pathMode = !g.pathMode
		if !g.pathMode {
			log.Println("switched to normal mode")
			// Clear path.
			g.path = nil
		} else {
			log.Println("switched to path mode")
		}

		// Unset path start and end.
		g.pathStart = [2]int{-1, -1}
		g.pathEnd = [2]int{-1, -1}
	}

	// Handle pathfinding input.
	if !g.pathMode {
		return
	}

	// Start of path if mouse is clicked.
	selectedTile := [2]int{-1, -1}
	selected := false
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		selectedTile = g.getTileAtCursor()
		selected = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		selectedTile = g.clickedTile
		selected = true
	}
	if selected {
		// Clear the path and start over if we have a start and end.
		if g.pathEnd != [2]int{-1, -1} {
			// Unset path start and end.
			g.pathStart = [2]int{-1, -1}
			g.pathEnd = [2]int{-1, -1}

			// Clear path.
			g.path = nil
		}
		if g.pathStart[0] == -1 && g.pathStart[1] == -1 {
			log.Println("setting path start to", selectedTile)
			g.pathStart = selectedTile
		} else if g.pathEnd[0] == -1 && g.pathEnd[1] == -1 {
			log.Println("setting path end to", selectedTile)
			g.pathEnd = selectedTile

			// Calculate path.
			log.Println("calculating path from", g.pathStart, "to", g.pathEnd, "...")
			g.path = g.currentLevel.findPath(g.pathStart, g.pathEnd)
			log.Println("path:", g.path)
		}
	}
}

func (g *Game) drawPathInfo(screen *ebiten.Image) {
	// Print info if we are in path mode.
	if !g.pathMode {
		return
	}

	// Print path info.
	x := screen.Bounds().Dx() - 200
	ebitenutil.DebugPrintAt(screen, "Path mode: ON\n"+
		fmt.Sprintf("Path start: %v\n", g.pathStart)+
		fmt.Sprintf("Path end: %v\n", g.pathEnd)+
		fmt.Sprintf("Path: %v\n", g.path), x, 0)
}

type node struct {
	pos    [2]int
	parent *node // The parent node for retracing the path.
}

func (l *Level) findPath(start, end [2]int) [][2]int {
	// NOTE: This is pretty bad code, so you might want to replace it with something better.
	// The open list contains the nodes that we have to check.
	openList := []*node{{pos: start}}

	// The closed list contains the nodes that we have already checked.
	closedList := make(map[[2]int]bool)

	// While we have nodes to check.
	for len(openList) > 0 {
		// Get the first node from the open list.
		current := openList[0]
		openList = openList[1:]

		// If we have reached the end, we are done.
		if current.pos == end {
			// Retrace the path.
			var path [][2]int
			for current != nil {
				path = append(path, current.pos)
				current = current.parent
			}
			return path
		}

		// Add the current node to the closed list.
		closedList[current.pos] = true

		// Get the neighbors of the current node.
		for _, n := range l.getNeighbors(current.pos[0], current.pos[1]) {
			// If the neighbor is already in the closed list, we skip it.
			if closedList[n] {
				continue
			}

			// If the tile is water, we skip it since we cannot cross it.
			if l.Tiles[n[1]*l.Width+n[0]].Type() == TileTypeWater {
				continue
			}

			// If the neighbor is already in the open list, we skip it.
			if contains(openList, n) {
				continue
			}

			// Add the neighbor to the open list.
			openList = append(openList, &node{pos: n, parent: current})
		}
	}

	// We have not found a path.
	return nil
}

func contains(nodes []*node, n [2]int) bool {
	for _, node := range nodes {
		if node.pos == n {
			return true
		}
	}
	return false
}
