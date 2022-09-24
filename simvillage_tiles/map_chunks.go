package simvillage_tiles

// MapChunk represents a chunk of the map.
type MapChunk struct {
	Dimensions
	Ground        *Layer // Ground or terrain.
	GroundOverlay *Layer // Overlays like carpets, scratchmarks, etc.
	Objects       *Layer // Objects like stones, flowers, etc.
	Structures    *Layer // Structures like walls, gates, fences, etc.
	Roof          *Layer // Roof or ceiling.
}

// newMapChunk returns a new map chunk with the given width and height.
func newMapChunk(width, height int) *MapChunk {
	return &MapChunk{
		Dimensions:    NewDimensions(width, height),
		Ground:        newLayer(width, height),
		GroundOverlay: newLayer(width, height),
		Objects:       newLayer(width, height),
		Structures:    newLayer(width, height),
		Roof:          newLayer(width, height),
	}
}

// drawObject draws the given object at (dx, dy).
func (m *MapChunk) drawObject(h drawable, dx, dy int) {
	// Draw all object layers onto the the appropriate layers of the map chunk.
	// TODO: Maybe make this a property of *Layer.
	drawOnLayer := func(dst *Layer, src *Layer, dx, dy int) {
		for x := 0; x < src.Width; x++ {
			for y := 0; y < src.Height; y++ {
				if t := src.getTile(x, y); t > 0 {
					dst.setTile(x+dx, y+dy, t) // set the tile
				} else if t == -1 {
					dst.setTile(x+dx, y+dy, 0) // clear the tile
				}
			}
		}

	}
	drawOnLayer(m.Ground, h.Ground(), dx, dy)
	drawOnLayer(m.GroundOverlay, h.GroundOverlay(), dx, dy)
	drawOnLayer(m.Structures, h.Structures(), dx, dy)
}
