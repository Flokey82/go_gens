// Package genmapvoronoi is based on https://mewo2.com/notes/terrain/.
package genmapvoronoi

import (
	"github.com/Flokey82/go_gens/vectors"
	"github.com/Flokey82/go_gens/vmesh"
	"github.com/pzsz/voronoi"
	//"log"
	//"fmt"
)

type Extent struct {
	width  float64
	height float64
}

var defaultExtent = &Extent{
	width:  1,
	height: 1,
}

type Generator func(*Params) *vmesh.Heightmap

type Fontsizes struct {
	region int
	city   int
	town   int
}

type Params struct {
	extent    *Extent
	generator Generator
	npts      int
	ncities   int
	nterrs    int
	fontsizes Fontsizes
}

type Terrain struct {
	params *Params
	mesh   *vmesh.Mesh
	h      *vmesh.Heightmap

	cities          []int
	cityTerritories []int
	rivers          []int              // vertex to river id mapping
	riverPaths      [][]voronoi.Vertex // river paths
	coasts          [][]voronoi.Vertex
	terr            []int              // vertex to territory id mapping
	borders         [][]voronoi.Vertex // territory border paths
	cityBorders     [][]voronoi.Vertex
}

func DoMap() {
	r := &Terrain{
		params: &Params{
			extent: defaultExtent,
			//generator: generateCoast,
			npts:    16384,
			ncities: 15,
			nterrs:  5,
			fontsizes: Fontsizes{
				region: 40,
				city:   25,
				town:   20,
			},
		},
	}
	r.genTerrain()
	r.regenMapFeatures()
	r.ExportSVG("test.svg")

	if err := r.h.ExportOBJ("tmp.obj"); err != nil {
		panic(err)
	}
}

func (r *Terrain) regenMapFeatures() {
	riverThreshold := 0.005
	r.rivers = getRivers(r.h, riverThreshold)
	r.riverPaths = getRiverPaths(r.h, riverThreshold)
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

func (r *Terrain) genTerrain() {
	r.mesh = vmesh.GenerateGoodMesh(r.params.npts, &vmesh.Extent{
		Width:  r.params.extent.width,
		Height: r.params.extent.height,
	})

	r.h = vmesh.NewHeightmap(r.mesh)
	r.h.Add(
		MeshSlope(r.mesh, vectors.RandomVec2(4)),
		MeshVolCone(r.mesh, -1),
		//MeshCone(r.mesh, mesh, runif(-1, -1)),
		MeshMountains(r.mesh, 50, 0.05),
	)
	for i := 0; i < 10; i++ {
		r.h = HeightRelax(r.h)
	}
	r.h = HeightPeaky(r.h)
	r.h = HeightNormalize(r.h)
	r.h = doErosion(r.h, runif(0, 0.1), 5)
	r.h = HeightSetSeaLevel(r.h, runif(0.2, 0.6))
	r.h = HeightFillSinks(r.h)
	r.h = HeightCleanCoast(r.h, 5)
}
