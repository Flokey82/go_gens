package simnpcs2

import (
	"image/color"
)

var (
	colorWaypoint    = color.RGBA{0xFF, 0xFF, 0x00, 255}
	colorBeing       = color.RGBA{0xFF, 0x00, 0x00, 255}
	colorDestination = color.RGBA{0x00, 0xFF, 0x00, 255}
	colorObstacle    = color.RGBA{0x00, 0x00, 0xFF, 255}
	colorItem        = color.RGBA{0xFF, 0x00, 0xFF, 255}
)

type imageIf interface {
	Set(x, y int, c color.Color)
}

func (w *World) renderFrame(img imageIf) {
	// Draw all entities and their paths.
	for _, c := range w.Beings {
		ai := c.(*AI)

		// Draw the path.
		if ai.Pathfinding.Waypoints != nil {
			for _, wpIdx := range ai.Pathfinding.Waypoints[ai.Pathfinding.WaypointIdx:] {
				wp := w.CellIdxToPos(wpIdx)
				img.Set(int(wp.X), int(wp.Y), colorWaypoint)
			}
		}

		// Draw the being.
		pos := c.Pos()
		img.Set(int(pos.X), int(pos.Y), colorBeing)

		// Draw the destination.
		if dst := ai.Destination; dst != nil {
			img.Set(int(dst.X), int(dst.Y), colorDestination)
		}
	}

	// Draw all obstacles.
	for i, c := range w.Cells {
		if c {
			pos := w.CellIdxToPos(i)
			img.Set(int(pos.X), int(pos.Y), colorObstacle)
		}
	}

	// Draw all items that are visible.
	for _, c := range w.Items {
		pos := c.Pos()
		img.Set(int(pos.X), int(pos.Y), colorItem)
	}
}
