// Package genmapvoronoi is based on https://mewo2.com/notes/terrain/.
package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/vmesh"
	"github.com/pzsz/voronoi"
	//"log"
	//"fmt"
)

type Extent struct {
	Width  float64
	Height float64
}

var DefaultExtent = &Extent{
	Width:  1,
	Height: 1,
}

type Params struct {
	Extent         *Extent
	NumPoints      int
	NumCities      int
	NumTerritories int
	RiverThreshold float64
}

type Terrain struct {
	params   *Params
	mesh     *vmesh.Mesh
	h        *vmesh.Heightmap
	sediment *vmesh.Heightmap

	cities          []int
	cityTerritories []int
	rivers          []int              // vertex to river id mapping
	riverPaths      [][]voronoi.Vertex // river paths
	coasts          [][]voronoi.Vertex
	terr            []int              // vertex to territory id mapping
	borders         [][]voronoi.Vertex // territory border paths
	cityBorders     [][]voronoi.Vertex
}

var DefaultParams = &Params{
	Extent:         DefaultExtent,
	NumPoints:      16384,
	NumCities:      15,
	NumTerritories: 5,
	RiverThreshold: 0.005,
}

func NewTerrain(params *Params) *Terrain {
	r := &Terrain{
		params: params,
	}

	r.genTerrain()
	r.regenMapFeatures()

	return r
}

func (r *Terrain) regenMapFeatures() {
	r.rivers = getRivers(r.h, r.params.RiverThreshold)
	r.riverPaths = getRiverPaths(r.h, r.params.RiverThreshold)
	r.coasts = contour(r.h, 0)

	// Place cities.
	placeCities(r)

	// Establish territories.
	r.terr = getTerritories(r)
	r.borders = getBorders(r)

	// Establish city territories.
	r.cityTerritories = getCityTerritories(r, r.terr)
	r.cityBorders = getCityBorders(r)
}
