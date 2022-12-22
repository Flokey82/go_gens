package genworldvoronoi

import (
	"log"
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/fogleman/delaunay"
)

// generateFibonacciSphere generates a number of points along a spiral on a sphere
// as a flat array of lat, lon coordinates.
func generateFibonacciSphere(seed int64, numPoints int, jitter float64) []float64 {
	rnd := rand.New(rand.NewSource(seed))
	var latLon []float64
	randomLat := make(map[int]float64)
	randomLon := make(map[int]float64)

	// Second algorithm from http://web.archive.org/web/20120421191837/http://www.cgafaq.info/wiki/Evenly_distributed_points_on_sphere
	s := 3.6 / math.Sqrt(float64(numPoints))
	dlong := math.Pi * (3 - math.Sqrt(5)) // ~2.39996323
	dz := 2.0 / float64(numPoints)

	for k, long, z := 0, 0.0, 1-(dz/2); k != numPoints; k++ {
		r := math.Sqrt(1 - z*z)
		latDeg := math.Asin(z) * 180 / math.Pi
		lonDeg := long * 180 / math.Pi
		if _, ok := randomLat[k]; !ok {
			randomLat[k] = rnd.Float64() - rnd.Float64()
		}
		if _, ok := randomLon[k]; !ok {
			randomLon[k] = rnd.Float64() - rnd.Float64()
		}

		latDeg += jitter * randomLat[k] * (latDeg - math.Asin(math.Max(-1, z-dz*2*math.Pi*r/s))*180/math.Pi)
		lonDeg += jitter * randomLon[k] * (s / r * 180 / math.Pi)
		latLon = append(latLon, latDeg, math.Mod(lonDeg, 360.0))
		long += dlong
		z -= dz
	}
	return latLon
}

// pushCartesianFromSpherical calculates x,y,z from spherical coordinates lat,lon and then push
// them onto out array; for one-offs pass nil as the first argument
func pushCartesianFromSpherical(out []float64, latDeg, lonDeg float64) []float64 {
	return append(out, latLonToCartesian(latDeg, lonDeg)...)
}

// latLonToCartesian converts latitude and longitude to x, y, z coordinates.
// See: https://rbrundritt.wordpress.com/2008/10/14/conversion-between-spherical-and-cartesian-coordinates-systems/
func latLonToCartesian(latDeg, lonDeg float64) []float64 {
	latRad := (latDeg / 180.0) * math.Pi
	lonRad := (lonDeg / 180.0) * math.Pi
	return []float64{
		math.Cos(latRad) * math.Cos(lonRad),
		math.Cos(latRad) * math.Sin(lonRad),
		math.Sin(latRad),
	}
}

// latLonFromVec3 converts a vectors.Vec3 to latitude and longitude.
// See: https://rbrundritt.wordpress.com/2008/10/14/conversion-between-spherical-and-cartesian-coordinates-systems/
func latLonFromVec3(position vectors.Vec3, sphereRadius float64) (float64, float64) {
	// See https://stackoverflow.com/questions/46247499/vector3-to-latitude-longitude
	lat := math.Asin(position.Z / sphereRadius) // theta
	lon := math.Atan2(position.Y, position.X)   // phi
	return radToDeg(lat), radToDeg(lon)
}

/** Add south pole back into the mesh.
 *
 * We run the Delaunay Triangulation on all points *except* the south
 * pole, which gets mapped to infinity with the stereographic
 * projection. This function adds the south pole into the
 * triangulation. The Delaunator guide explains how the halfedges have
 * to be connected to make the mesh work.
 * <https://mapbox.github.io/delaunator/>
 *
 * Returns the new {triangles, halfedges} for the triangulation with
 * one additional point added around the convex hull.
 */
func addSouthPoleToMesh(southPoleId int, d *delaunay.Triangulation) *delaunay.Triangulation {
	// This logic is from <https://github.com/redblobgames/dual-mesh>,
	// where I use it to insert a "ghost" region on the "back" side of
	// the planar map. The same logic works here. In that code I use
	// "s" for edges ("sides"), "r" for regions ("points"), t for triangles
	triangles := d.Triangles
	numSides := len(triangles)
	s_next_s := func(s int) int {
		if s%3 == 2 {
			return s - 2
		}
		return s + 1
	}

	halfedges := d.Halfedges
	numUnpairedSides := 0
	firstUnpairedSide := -1
	pointIdToSideId := make(map[int]int) // seed to side
	for s := 0; s < numSides; s++ {
		if halfedges[s] == -1 {
			numUnpairedSides++
			pointIdToSideId[triangles[s]] = s
			firstUnpairedSide = s
		}
	}

	newTriangles := make([]int, numSides+3*numUnpairedSides)
	copy(newTriangles, triangles)

	newHalfedges := make([]int, numSides+3*numUnpairedSides)
	copy(newHalfedges, halfedges)

	for i, s := 0, firstUnpairedSide; i < numUnpairedSides; i++ {
		log.Println(i)
		// Construct a pair for the unpaired side s
		newSide := numSides + 3*i
		newHalfedges[s] = newSide
		newHalfedges[newSide] = s
		newTriangles[newSide] = newTriangles[s_next_s(s)]

		// Construct a triangle connecting the new side to the south pole
		newTriangles[newSide+1] = newTriangles[s]
		newTriangles[newSide+2] = southPoleId
		k := numSides + (3*i+4)%(3*numUnpairedSides)
		newHalfedges[newSide+2] = k
		newHalfedges[k] = newSide + 2
		s = pointIdToSideId[newTriangles[s_next_s(s)]]
	}

	return &delaunay.Triangulation{
		Triangles: newTriangles,
		Halfedges: newHalfedges,
	}
}

// stereographicProjection converts 3d coordinates into two dimensions.
func stereographicProjection(xyz []float64) []float64 {
	// See <https://en.wikipedia.org/wiki/Stereographic_projection>
	numPoints := len(xyz) / 3
	xy := make([]float64, 0, 2*numPoints)
	for r := 0; r < numPoints; r++ {
		x := xyz[3*r]
		y := xyz[3*r+1]
		z := xyz[3*r+2]
		xy = append(xy, x/(1-z), y/(1-z)) // Append projected 2d coordinates.
	}
	return xy
}

type SphereMesh struct {
	mesh   *TriangleMesh
	xyz    []float64    // Region coordinates
	latLon [][2]float64 // Region latitude and longitude
}

func MakeSphere(seed int64, numPoints int, jitter float64) (*SphereMesh, error) {
	latlong := generateFibonacciSphere(seed, numPoints, jitter)
	var xyz []float64
	var latLon [][2]float64
	for r := 0; r < len(latlong); r += 2 {
		// HACKY! Fix this properly!
		nla, nlo := latLonFromVec3(convToVec3(latLonToCartesian(latlong[r], latlong[r+1])).Normalize(), 1.0)
		latLon = append(latLon, [2]float64{nla, nlo})
		xyz = pushCartesianFromSpherical(xyz, latlong[r], latlong[r+1])
	}

	xy := stereographicProjection(xyz)
	pts := make([]delaunay.Point, 0, len(xy)/2)
	for i := 0; i < len(xy); i += 2 {
		pts = append(pts, delaunay.Point{X: xy[i], Y: xy[i+1]})
	}

	tri, err := delaunay.Triangulate(pts)
	if err != nil {
		return nil, err
	}

	// TODO: rotate an existing point into this spot instead of creating one.
	xyz = append(xyz, 0, 0, 1)
	latLon = append(latLon, [2]float64{-90.0, 45.0})

	tri = addSouthPoleToMesh((len(xyz)/3)-1, tri)

	mesh := NewTriangleMesh(0, len(tri.Triangles), make([]Vertex, numPoints+1), tri.Triangles, tri.Halfedges)
	return &SphereMesh{
		mesh:   mesh,
		xyz:    xyz,
		latLon: latLon,
	}, nil
}
