package genworldvoronoi

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
)

// GetTile returns the image of the tile at the given coordinates and zoom level.
// NOTE: This is definitely upside down. :(
func (m *Map) GetTile(x, y, zoom int) image.Image {
	// Wrap the tile coordinates.
	x = x % (1 << uint(zoom))
	if x < 0 {
		x += 1 << uint(zoom)
	}
	y = y % (1 << uint(zoom))
	if y < 0 {
		y += 1 << uint(zoom)
	}
	// log.Println("tile", x, y, "zoom", zoom)

	tbb := tileBoundingBox(x, y, zoom)
	la1, lo1, la2, lo2 := tbb.ToLatLon()
	latLonMargin := 20.0 / float64(zoom)

	// Since our mercator conversion gives us absolute pixel coordinates, we need to
	// remove the offset of the tile we are rendering from the path coordinates.
	dx, _ := latLonToPixels(la1, lo1, zoom)
	_, dy2 := latLonToPixels(la2, lo2, zoom)

	// Create a new image to draw the tile on.
	dest := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	gc := draw2dimg.NewGraphicContext(dest)

	// Calculate the min and max elevation and moisture values.
	min, max := minMax(m.Elevation)
	_, maxMois := minMax(m.Moisture)
	for i := 0; i < m.mesh.numRegions; i++ {
		rLat := m.LatLon[i][0]
		rLon := m.LatLon[i][1]

		// Check if we are within the tile with a small margin,
		// taking into account that we might have wrapped around the world.
		if rLat < la1-latLonMargin || rLat > la2+latLonMargin || rLon < lo1-latLonMargin || rLon > lo2+latLonMargin {
			// Check if the tile and the region we are looking at is adjecent to +/- 180 degrees and
			// NOTE: This could be improved by checking if one of the corners of the region is within the tile.
			if lo1 > -175 && lo2 < 175 || rLon < 175 && rLon > -175 {
				continue
			}
		}

		// Draw the path that outlines the region.
		var path [][2]float64
		for _, j := range m.mesh.r_circulate_t(nil, i) {
			tLat := m.t_latLon[j][0]
			tLon := m.t_latLon[j][1]

			// Check if we have wrapped around the world.
			if tLon-rLon > 120 {
				tLon -= 360
			} else if tLon-rLon < -120 {
				tLon += 360
			}

			// Calculate the coordinates of the path point.
			x, y := latLonToPixels(tLat, tLon, zoom)
			path = append(path, [2]float64{(x - dx), tileSize - (y - dy2)})
		}

		// Now check if the region we are looking at has wrapped around the world /
		// +- 180 degrees. If so, we need to adjust the points in the path.
		if lo1 < -175 && rLon > 175 {
			for i := range path {
				path[i][0] -= float64(sizeFromZoom(zoom))
			}
		} else if lo2 > 175 && rLon < -175 {
			for i := range path {
				path[i][0] += float64(sizeFromZoom(zoom))
			}
		}

		// Calculate the color of the region.
		elev := m.Elevation[i]
		val := (elev - min) / (max - min)
		var col color.NRGBA
		if elev <= 0 {
			col = genBlue(val)
		} else {
			valElev := elev / max
			valMois := m.Moisture[i] / maxMois
			col = getWhittakerModBiomeColor(rLat, valElev, valMois, val)
		}

		// If the path is empty, we can skip it.
		if len(path) == 0 {
			continue
		}

		// Draw the path.
		gc.SetFillColor(col)
		gc.SetStrokeColor(col)
		gc.SetLineWidth(0)
		if rLat > 60 {
			gc.SetStrokeColor(color.NRGBA{0, 0, 0, 255})
			gc.SetLineWidth(1)
		}
		gc.BeginPath()
		gc.MoveTo(path[0][0], path[0][1])
		for _, p := range path[1:] {
			gc.LineTo(p[0], p[1])
		}
		gc.Close()
		gc.FillStroke()
	}

	return dest
}
