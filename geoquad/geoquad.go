// Package geoquad provides a simple way to generate a quadtree for a given set of lat/lon coordinates.
// The quadtree is laid out like described in:
// https://www.maptiler.com/google-maps-coordinates-tile-bounds-projection/
// to represent a tileset that is used by Google Maps, OpenStreetMap, etc.
// This will allow us to quickly find all coordinates within a given tile, or all coordinates
// within a given bounding box, as well as the closest coordinate to a given coordinate.
package geoquad

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// min and max latitude and longitude.
const (
	MinLongitude float64 = -180.0
	MaxLongitude float64 = 180.0
	MinLatitude  float64 = -90.0
	MaxLatitude  float64 = 90.0
)

// Point represents a point (latitude, longitude) in 2D space.
type Point struct {
	Lat  float64
	Lon  float64
	Data interface{}
}

// distance returns the great arc distance between the point and the other point.
func (p Point) distance(other Point) float64 {
	return distHaversine(p.Lat, p.Lon, other.Lat, other.Lon)
}

// Rect represents a rectangle in 2D space.
type Rect struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

// intersects returns true if the rectangle intersects the other rectangle.
func (r Rect) intersects(other Rect) bool {
	return r.MinLon <= other.MaxLon && r.MaxLon >= other.MinLon &&
		r.MinLat <= other.MaxLat && r.MaxLat >= other.MinLat
}

// contains returns true if the rectangle contains the other rectangle.
func (r Rect) contains(other Rect) bool {
	return r.MinLon <= other.MinLon && r.MaxLon >= other.MaxLon &&
		r.MinLat <= other.MinLat && r.MaxLat >= other.MaxLat
}

// containsPoint returns true if the rectangle contains the point.
func (r Rect) containsPoint(p Point) bool {
	return r.MinLon <= p.Lon && r.MaxLon >= p.Lon &&
		r.MinLat <= p.Lat && r.MaxLat >= p.Lat
}

func distHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLatSin := math.Sin(degToRad(lat2-lat1) / 2)
	dLonSin := math.Sin(degToRad(lon2-lon1) / 2)
	a := dLatSin*dLatSin + dLonSin*dLonSin*math.Cos(degToRad(lat1))*math.Cos(degToRad(lat2))
	return 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

type TileCoordinates struct {
	X, Y, Zoom int
}

func (t TileCoordinates) String() string {
	return fmt.Sprintf("x: %d, y: %d, zoom: %d", t.X, t.Y, t.Zoom)
}

// QuadTree represents a quadtree.
type QuadTree struct {
	TileCoordinates // The coordinates and zoom of the tile this quadtree represents.
	Points          []Point
	Bounds          Rect // The bounds of the tile this quadtree represents.
	// The four child quadtrees.
	// The first two are the top left and top right quadtrees.
	// The second two are the bottom left and bottom right quadtrees.
	Children [4]*QuadTree
	Parent   *QuadTree
}

// NewQuadTree creates a new quadtree for the given points.
func NewQuadTree(points []Point) *QuadTree {
	// Find the bounds of the points.
	bounds := Rect{
		MinLon: MinLongitude,
		MinLat: MinLatitude,
		MaxLon: MaxLongitude,
		MaxLat: MaxLatitude,
	}
	/*
		for _, p := range points {
			if p.Lon < bounds.MinLon {
				bounds.MinLon = p.Lon
			}
			if p.Lat < bounds.MinLat {
				bounds.MinLat = p.Lat
			}
			if p.Lon > bounds.MaxLon {
				bounds.MaxLon = p.Lon
			}
			if p.Lat > bounds.MaxLat {
				bounds.MaxLat = p.Lat
			}
		}*/

	// Create the quadtree.
	return newQuadTree(points, bounds, 0, 0, 0, nil)
}

// newQuadTree creates a new quadtree for the given points.
// The points are sorted by X coordinate.
func newQuadTree(points []Point, bounds Rect, x, y, zoom int, parent *QuadTree) *QuadTree {
	// Create the quadtree.
	qt := &QuadTree{
		TileCoordinates: TileCoordinates{
			X:    x,
			Y:    y,
			Zoom: zoom,
		},
		Points: points,
		Bounds: bounds,
		Parent: parent,
	}

	// If there are no points, return.
	if len(points) == 0 {
		return qt
	}

	// If there is only one point, return.
	if len(points) == 1 {
		return qt
	}

	// Find the median point.
	medianLon := (bounds.MinLon + bounds.MaxLon) / 2
	medianLat := (bounds.MinLat + bounds.MaxLat) / 2

	// Split the points into four groups.
	var topLeft, topRight, bottomLeft, bottomRight []Point
	for i := range points {
		p := points[i]
		if p.Lat >= medianLat {
			if p.Lon < medianLon {
				topLeft = append(topLeft, p)
			} else {
				topRight = append(topRight, p)
			}
		} else {
			if p.Lon < medianLon {
				bottomLeft = append(bottomLeft, p)
			} else {
				bottomRight = append(bottomRight, p)
			}
		}
	}

	// Create the child quadtrees.
	if len(topLeft) > 0 {
		qt.Children[0] = newQuadTree(topLeft, Rect{
			MinLon: bounds.MinLon,
			MinLat: medianLat,
			MaxLon: medianLon,
			MaxLat: bounds.MaxLat,
		}, x*2, y*2, zoom+1, qt)
		qt.Children[0].Verify()
	}
	if len(topRight) > 0 {
		qt.Children[1] = newQuadTree(topRight, Rect{
			MinLon: medianLon,
			MinLat: medianLat,
			MaxLon: bounds.MaxLon,
			MaxLat: bounds.MaxLat,
		}, x*2+1, y*2, zoom+1, qt)
		qt.Children[1].Verify()
	}
	if len(bottomLeft) > 0 {
		qt.Children[2] = newQuadTree(bottomLeft, Rect{
			MinLon: bounds.MinLon,
			MinLat: bounds.MinLat,
			MaxLon: medianLon,
			MaxLat: medianLat,
		}, x*2, y*2+1, zoom+1, qt)
		qt.Children[2].Verify()
	}
	if len(bottomRight) > 0 {
		qt.Children[3] = newQuadTree(bottomRight, Rect{
			MinLon: medianLon,
			MinLat: bounds.MinLat,
			MaxLon: bounds.MaxLon,
			MaxLat: medianLat,
		}, x*2+1, y*2+1, zoom+1, qt)
		qt.Children[3].Verify()
	}

	return qt
}

func (qt *QuadTree) Verify() {
	// Verify the points are all in the bounds.
	for _, p := range qt.Points {
		if p.Lon < qt.Bounds.MinLon || p.Lon > qt.Bounds.MaxLon ||
			p.Lat < qt.Bounds.MinLat || p.Lat > qt.Bounds.MaxLat {
			log.Println("Apoint is not in bounds:", p)
			panic("point is not in bounds")
		}
		if !qt.Bounds.containsPoint(p) {
			log.Println("Bpoint is not in bounds:", p)
			panic("point is not in bounds")
		}
	}
}

// FindPointsInTile returns all points in the given tile.
func (qt *QuadTree) FindPointsInTile(x, y, zoom int) []Point {
	// If the tile is not in the quadtree, return.
	if zoom < qt.Zoom {
		return nil
	}

	// If the tile is exactly the one requested, return.
	if zoom == qt.Zoom && x == qt.X && y == qt.Y {
		return qt.Points
	}

	// Check if the tile is a parent of the tile we are looking for.
	xAtCurrentTile := x >> uint(zoom-qt.Zoom)
	yAtCurrentTile := y >> uint(zoom-qt.Zoom)
	if xAtCurrentTile != qt.X || yAtCurrentTile != qt.Y {
		return nil
	}

	// Find the points in the child quadtrees.
	// Get the index of the child that will contain the tile that we are looking for.
	childIndex := 0
	xChild := x >> uint(zoom-qt.Zoom-1)
	yChild := y >> uint(zoom-qt.Zoom-1)

	if xChild%2 == 1 {
		childIndex += 1
	}
	if yChild%2 == 1 {
		childIndex += 2
	}

	// If the child is not nil, find the points in the child.
	if qt.Children[childIndex] != nil {
		return qt.Children[childIndex].FindPointsInTile(x, y, zoom)
	}

	// If the child is nil, return the points in the current tile.
	allowInaccurateLeafNodes := true

	// If we don't have a child node with the next level
	if allowInaccurateLeafNodes {
		// If the child is nil, return the points in the quadtree.
		return qt.Points
	}

	return nil
}

// FindPointsInRect returns all points in the given rectangle.
func (qt *QuadTree) FindPointsInRect(rect Rect) []Point {
	// If the rectangle is not in the quadtree, return.
	if !rect.intersects(qt.Bounds) {
		return nil
	}

	// If the rectangle is in the quadtree, return the points.
	if rect.contains(qt.Bounds) {
		return qt.Points
	}

	if len(qt.Points) == 1 { // TODO VERIFY OR OPTIMIZE
		if rect.containsPoint(qt.Points[0]) {
			return qt.Points
		}
		return nil
	}

	// Find the points in the child quadtrees.
	var points []Point
	for _, child := range qt.Children {
		if child == nil {
			continue
		}
		points = append(points, child.FindPointsInRect(rect)...)
	}

	return points
}

func getNeighboringTileCoords(x, y, zoom int) []TileCoordinates {
	// If zoom is 0, there are no neighboring tiles.
	if zoom == 0 {
		return nil
	}

	// If zoom is 1, there are 4 neighboring tiles.
	if zoom == 1 {
		return []TileCoordinates{
			{X: 0, Y: 0, Zoom: 1},
			{X: 1, Y: 0, Zoom: 1},
			{X: 0, Y: 1, Zoom: 1},
			{X: 1, Y: 1, Zoom: 1},
		}
	}

	maxTileAtZoom := 1 << uint(zoom)

	// Get the neighboring tile coordinates including diagonals.
	var coords []TileCoordinates
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Wrap the coordinates if they are out of bounds.
			nx := (x + i + maxTileAtZoom) % maxTileAtZoom
			ny := (y + j + maxTileAtZoom) % maxTileAtZoom
			coords = append(coords, TileCoordinates{X: nx, Y: ny, Zoom: zoom})
		}
	}

	return coords
}

// FindNearestNeighbor returns the closest point to the given point within the quadtree.
func (qt *QuadTree) FindNearestNeighbor(point Point) (Point, bool) {
	// If the point is not in the quadtree, return.
	if !qt.Bounds.containsPoint(point) {
		return Point{}, false
	}

	// If there are no points, return.
	if len(qt.Points) == 0 {
		return Point{}, false
	}

	// If there is only one point, return the point.
	if len(qt.Points) == 1 {
		return qt.Points[0], true
	}

	// Find the leaf tile that contains the point.
	leafTile := qt.FindLeafTile(point)

	// Find the point in the leaf tile that is closest to the given point.
	var closestPoint Point
	minDistance := math.Inf(1)
	for i := range leafTile.Points {
		p := leafTile.Points[i]
		distance := point.distance(p)
		if distance < minDistance {
			closestPoint = p
			minDistance = distance
		}
	}

	// Get all neighboring tiles.
	coords := getNeighboringTileCoords(leafTile.X, leafTile.Y, leafTile.Zoom)

	// Find the closest point in the neighboring tiles.
	for _, coord := range coords {
		// Get the tile.
		points := qt.FindPointsInTile(coord.X, coord.Y, coord.Zoom)
		// Find the closest point in the tile.
		for i := range points {
			p := points[i]
			//log.Println("p", p)
			distance := point.distance(p)
			if distance < minDistance {
				closestPoint = p
				minDistance = distance
			}
		}
	}
	return closestPoint, true
}

// FindLeafTile returns the leaf tile that contains the given point.
func (qt *QuadTree) FindLeafTile(point Point) *QuadTree {
	// If the point is not in the quadtree, return.
	if !qt.Bounds.containsPoint(point) {
		return nil
	}

	// If there are no points, return.
	if len(qt.Points) == 0 {
		return nil
	}

	// If there is only one point, return the quadtree.
	if len(qt.Points) == 1 {
		return qt
	}

	// Find the leaf tile in the child quadtrees.
	for _, child := range qt.Children {
		if child == nil {
			continue
		}
		leafTile := child.FindLeafTile(point)
		if leafTile != nil {
			return leafTile
		}
	}

	// If no leaf tile was found, return the quadtree.
	return qt
}

// FindPoint returns the QuadTree that contains the given point.
func (qt *QuadTree) FindPoint(point Point) *QuadTree {
	// If the point is not in the quadtree, return.
	if !qt.Bounds.containsPoint(point) {
		return nil
	}

	// If there are no points, return.
	if len(qt.Points) == 0 {
		return nil
	}

	// If there is only one point, return it.
	if len(qt.Points) == 1 {
		return qt
	}

	// Find the point in the child quadtrees.
	for _, child := range qt.Children {
		if child == nil {
			continue
		}
		if cc := child.FindPoint(point); cc != nil {
			return cc
		}
	}

	return nil
}

const exportPNGScale = 2

// ExportToPNG exports the quadtree to a PNG image.
func (qt *QuadTree) ExportToPNG(filename string) error {
	// Create the image.
	height := int(qt.Bounds.MaxLat-qt.Bounds.MinLat) * exportPNGScale
	width := int(qt.Bounds.MaxLon-qt.Bounds.MinLon) * exportPNGScale
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// Fill the image with white.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// Draw the quadtree.
	qt.draw(img)

	// Save the image.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}
	return nil
}

// draw draws the quadtree to the image.
func (qt *QuadTree) draw(img *image.RGBA) {
	// Draw the quadtree bounds.
	drawRect(img, qt.Bounds, color.RGBA{0, 0, 0, 255})

	// Draw the points.
	for _, p := range qt.Points {
		drawPoint(img, p, color.RGBA{255, 0, 0, 255})
	}

	// Draw the child quadtrees.
	for _, child := range qt.Children {
		if child == nil {
			continue
		}
		child.draw(img)
	}
}

// drawPoint draws a point to the image.
func drawPoint(img *image.RGBA, p Point, c color.RGBA) {
	dx := int(MinLongitude)
	dy := int(MinLatitude)
	x := (int(p.Lon) - dx) * exportPNGScale
	y := (int(p.Lat) - dy) * exportPNGScale
	img.Set(x, y, c)
}

// drawRect draws a rectangle to the image.
func drawRect(img *image.RGBA, r Rect, c color.RGBA) {
	dx := int(MinLongitude)
	dy := int(MinLatitude)

	minX := (int(r.MinLon) - dx) * exportPNGScale
	minY := (int(r.MinLat) - dy) * exportPNGScale
	maxX := (int(r.MaxLon) - dx) * exportPNGScale
	maxY := (int(r.MaxLat) - dy) * exportPNGScale

	for x := minX; x < maxX; x++ {
		img.Set(x, minY, c)
		img.Set(x, maxY, c)
	}
	for y := minY; y < maxY; y++ {
		img.Set(minX, y, c)
		img.Set(maxX, y, c)
	}
}
