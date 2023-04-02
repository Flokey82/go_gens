package gencitymap

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/Flokey82/go_gens/vectors"
	svgo "github.com/ajstarks/svgo"
	"github.com/llgcode/draw2d/draw2dimg"
)

type StreamlineIntegration struct {
	Seed          vectors.Vec2
	OriginalDir   vectors.Vec2
	PreviousDir   vectors.Vec2
	PreviousPoint vectors.Vec2
	Streamline    []vectors.Vec2
	Valid         bool
}

type StreamlineParams struct {
	Dsep              float64 // Streamline seed separating distance
	Dtest             float64 // Streamline integration separating distance
	Dstep             float64 // Step size
	Dcirclejoin       float64 // How far to look to join circles - (e.g. 2 x dstep)
	Dlookahead        float64 // How far to look ahead to join up dangling
	Joinangle         float64 // Angle to join roads in radians
	PathIterations    int     // Path integration iteration limit
	SeedTries         int     // Max failed seeds
	SimplifyTolerance float64
	CollideEarly      float64 // Chance of early collision 0-1
}

type StreamlineGenerator struct {
	SEED_AT_ENDPOINTS    bool
	NEAR_EDGE            int // Sample near edge
	majorGrid            *GridStorage
	minorGrid            *GridStorage
	nStreamlineStep      int
	nStreamlineLookBack  int
	dcollideselfSq       float64
	candidateSeedsMajor  []vectors.Vec2
	candidateSeedsMinor  []vectors.Vec2
	streamlinesDone      bool
	lastStreamlineMajor  bool
	resolve              func()
	allStreamlines       [][]vectors.Vec2
	streamlinesMajor     [][]vectors.Vec2
	streamlinesMinor     [][]vectors.Vec2
	allStreamlinesSimple [][]vectors.Vec2 // Reduced vertex count
	params               *StreamlineParams
	paramsSq             *StreamlineParams
	worldDimensions      vectors.Vec2
	origin               vectors.Vec2
	integrator           FieldIntegratorIf
	rng                  *rand.Rand
}

func NewStreamlineGenerator(seed int64, integrator FieldIntegratorIf, origin, worldDimensions vectors.Vec2, params *StreamlineParams) (*StreamlineGenerator, error) {
	if params.Dstep > params.Dsep {
		return nil, errors.New("STREAMLINE SAMPLE DISTANCE BIGGER THAN DSEP")
	}

	// Enforce test < sep
	if params.Dtest > params.Dsep {
		params.Dtest = params.Dsep
	}

	gen := &StreamlineGenerator{
		SEED_AT_ENDPOINTS:   true,
		NEAR_EDGE:           3, // Sample near edge
		majorGrid:           NewGridStorage(worldDimensions, origin, params.Dsep),
		minorGrid:           NewGridStorage(worldDimensions, origin, params.Dsep),
		streamlinesDone:     true,
		lastStreamlineMajor: true,
		params:              params,
		paramsSq:            params,
		worldDimensions:     worldDimensions,
		origin:              origin,
		integrator:          integrator,
		rng:                 rand.New(rand.NewSource(seed)),
	}

	// Needs to be less than circlejoin
	gen.dcollideselfSq = (params.Dcirclejoin / 2) * (params.Dcirclejoin / 2)
	gen.nStreamlineStep = int(params.Dcirclejoin / params.Dstep)
	gen.nStreamlineLookBack = 2 * gen.nStreamlineStep
	gen.setParamsSq()
	return gen, nil
}

// ExportToPNG exports the streamlines to a PNG file.
func (sg *StreamlineGenerator) ExportToPNG(filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, int(sg.worldDimensions.X), int(sg.worldDimensions.Y)))

	// Fill the background with black.
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	// New draw2d graphics context on an RGBA image.
	gc := draw2dimg.NewGraphicContext(img)

	// Draw minor streamlines.
	gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})
	gc.SetLineWidth(4.0)
	for _, v := range sg.streamlinesMajor {
		// Draw a path.
		gc.MoveTo(v[0].X-sg.origin.X, v[0].Y-sg.origin.Y)
		for _, p := range v[1:] {
			gc.LineTo(p.X-sg.origin.X, p.Y-sg.origin.Y)
		}
		gc.Stroke()
	}

	gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})
	gc.SetLineWidth(2.0)
	for _, v := range sg.streamlinesMinor {
		// Draw a path.
		gc.BeginPath()
		gc.MoveTo(v[0].X-sg.origin.X, v[0].Y-sg.origin.Y)
		for _, p := range v[1:] {
			gc.LineTo(p.X-sg.origin.X, p.Y-sg.origin.Y)
		}
		gc.Stroke()
	}

	// Save to file.
	return draw2dimg.SaveToPngFile(filename, img)
}

// ExportToSVG exports the streamlines to an SVG file.
func (sg *StreamlineGenerator) ExportToSVG(filename string) error {
	svgFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	svg := svgo.New(svgFile)
	svg.Start(int(sg.worldDimensions.X), int(sg.worldDimensions.Y))

	svg.Gstyle("stroke:rgb(0,255,0);stroke-width:2")
	for _, v := range sg.streamlinesMajor {
		x, y := convToPairs(v, sg.origin)
		svg.Polyline(x, y, "fill:none")
	}
	for _, v := range sg.streamlinesMinor {
		x, y := convToPairs(v, sg.origin)
		svg.Polyline(x, y, "stroke:rgb(255,0,0);fill:none")
	}

	svg.Gend()
	svg.End()
	return svgFile.Close()
}

func (sg *StreamlineGenerator) ClearStreamlines() {
	sg.allStreamlinesSimple = nil
	sg.allStreamlines = nil
	sg.streamlinesMajor = nil
	sg.streamlinesMinor = nil
}

// joinDanglingStreamlines joins streamlines that are not closed.
func (sg *StreamlineGenerator) joinDanglingStreamlines() {
	// TODO do in update method
	for _, major := range []bool{true, false} {
		for _, streamline := range sg.streamlines(major) {
			// Ignore circles.
			if streamline[0].Equalish(streamline[len(streamline)-1]) {
				continue
			}

			newStart := sg.getBestNextPoint(streamline[0], streamline[4], streamline)
			if newStart != emptyVec2 {
				for _, p := range sg.pointsBetween(streamline[0], newStart, sg.params.Dstep) {
					streamline = append([]vectors.Vec2{p}, streamline...)
					sg.grid(major).AddSample(p, nil)
				}
			}

			newEnd := sg.getBestNextPoint(streamline[len(streamline)-1], streamline[len(streamline)-4], streamline)
			if newEnd != emptyVec2 {
				for _, p := range sg.pointsBetween(streamline[len(streamline)-1], newEnd, sg.params.Dstep) {
					streamline = append(streamline, p)
					sg.grid(major).AddSample(p, nil)
				}
			}
		}
	}

	// Reset simplified streamlines
	sg.allStreamlinesSimple = nil
	for _, s := range sg.allStreamlines {
		sg.allStreamlinesSimple = append(sg.allStreamlinesSimple, sg.simplifyStreamline(s))
	}
}

// pointsBetween returns array of points from v1 to v2 such that they are separated by at most dsep, not including v1.
func (sg *StreamlineGenerator) pointsBetween(v1, v2 vectors.Vec2, dstep float64) []vectors.Vec2 {
	d := v1.DistanceTo(v2)
	nPoints := int(d / dstep)
	if nPoints == 0 {
		return nil
	}

	stepVector := v2.Sub(v1)

	out := make([]vectors.Vec2, 0, nPoints)
	var i int
	var next vectors.Vec2
	for i = 1; i <= nPoints; i++ {
		if sg.integrator.Integrate(next, true).LengthSquared() > 0.001 { // Test for degenerate point
			out = append(out, next)
		} else {
			return out
		}

		next = v1.Add(stepVector.Mul(float64(i) / float64(nPoints)))
	}

	return out
}

// getBestNextPoint returns the next point to join a streamline and returns a zero vector if there are no good candidates.
func (sg *StreamlineGenerator) getBestNextPoint(point, previousPoint vectors.Vec2, streamline []vectors.Vec2) vectors.Vec2 {
	nearbyPoints := sg.majorGrid.GetNearbyPoints(point, sg.params.Dlookahead)
	nearbyPoints = append(nearbyPoints, sg.minorGrid.GetNearbyPoints(point, sg.params.Dlookahead)...)
	direction := point.Sub(previousPoint)

	var closestSample vectors.Vec2
	closestDistance := math.Inf(1)

	for _, sample := range nearbyPoints {
		if !sample.Equalish(point) && !sample.Equalish(previousPoint) && !streamlineIncludes(streamline, sample) {
			differenceVector := sample.Sub(point)
			if differenceVector.Dot(direction) < 0 {
				// Backwards
				continue
			}

			// Acute angle between vectors (agnostic of CW, ACW)
			distanceToSample := point.DistanceToSquared(sample)
			if distanceToSample < 2*sg.paramsSq.Dstep {
				closestSample = sample
				break
			}

			// Filter by angle
			angleBetween := math.Abs(vectors.AngleBetweenTwoVectors(direction, differenceVector))
			if angleBetween < sg.params.Joinangle && distanceToSample < closestDistance {
				closestDistance = distanceToSample
				closestSample = sample
			}
		}
	}

	// TODO is reimplement simplify-js to preserve intersection points
	//  - this is the primary reason polygons aren't found
	// If trying to find intersections in the simplified graph
	// prevent ends getting pulled away from simplified lines
	if closestSample != emptyVec2 {
		closestSample = closestSample.Add(direction.SetLength(sg.params.SimplifyTolerance * 4))
	}

	return closestSample
}

var emptyVec2 = vectors.Vec2{}

func streamlineIncludes(streamline []vectors.Vec2, point vectors.Vec2) bool {
	for _, p := range streamline {
		if p.Equalish(point) {
			return true
		}
	}
	return false
}

// addExistingStreamlines adds existing streamlines from another generator.
// NOTE: This assumes 's' has already generated streamlines
func (sg *StreamlineGenerator) addExistingStreamlines(s *StreamlineGenerator) {
	sg.majorGrid.AddAll(s.majorGrid)
	sg.minorGrid.AddAll(s.minorGrid)
}

func (sg *StreamlineGenerator) setGrid(s *StreamlineGenerator) {
	sg.majorGrid = s.majorGrid
	sg.minorGrid = s.minorGrid
}

// update updates the streamline generator if necessary and returns true if state updates.
func (sg *StreamlineGenerator) update() bool {
	if !sg.streamlinesDone {
		sg.lastStreamlineMajor = !sg.lastStreamlineMajor
		if !sg.createStreamline(sg.lastStreamlineMajor) {
			sg.streamlinesDone = true
			sg.resolve()
		}
		return true
	}
	return false
}

// createAllStreamlines creates all streamlines all at once - will freeze if dsep small.
func (sg *StreamlineGenerator) createAllStreamlines(animate bool) {
	sg.streamlinesDone = false
	if !animate {
		major := true
		for sg.createStreamline(major) {
			major = !major
		}
	}
	sg.joinDanglingStreamlines()
}

func (sg *StreamlineGenerator) simplifyStreamline(streamline []vectors.Vec2) []vectors.Vec2 {
	/*
		simplified := make([]vectors.Vec2, 0, len(streamline))
		for _, point := range simplify(streamline, sg.params.SimplifyTolerance) {
			simplified = append(simplified, vectors.Vec2{point[0], point[1]})
		}
		return simplified
	*/
	return streamline
}

// createStreamline finds a seed, creates a streamline from that point, and pushes new candidate seeds to the queue.
// Returns false if seed isn't found within params.seedTries.
func (sg *StreamlineGenerator) createStreamline(major bool) bool {
	seed := sg.getSeed(major)
	if seed == emptyVec2 {
		return false
	}
	streamline := sg.integrateStreamline(seed, major)
	if sg.validStreamline(streamline) {
		sg.grid(major).AddPolyline(streamline)
		sg.appendStreamlines(major, streamline)
		sg.allStreamlines = append(sg.allStreamlines, streamline)
		sg.allStreamlinesSimple = append(sg.allStreamlinesSimple, sg.simplifyStreamline(streamline))

		// Add candidate seeds
		if !streamline[0].Equalish(streamline[len(streamline)-1]) {
			sg.appendCandidateSeeds(!major, streamline[0])
			sg.appendCandidateSeeds(!major, streamline[len(streamline)-1])
		}
	}

	return true
}

func (sg *StreamlineGenerator) validStreamline(s []vectors.Vec2) bool {
	return len(s) > 5
}

func (sg *StreamlineGenerator) setParamsSq() {
	sg.paramsSq = &StreamlineParams{}
	*sg.paramsSq = *sg.params
	sg.paramsSq.Dsep *= sg.paramsSq.Dsep
	sg.paramsSq.Dtest *= sg.paramsSq.Dtest
	sg.paramsSq.Dstep *= sg.paramsSq.Dstep
	sg.paramsSq.Dcirclejoin *= sg.paramsSq.Dcirclejoin
	sg.paramsSq.Dlookahead *= sg.paramsSq.Dlookahead
	sg.paramsSq.Joinangle *= sg.paramsSq.Joinangle
	sg.paramsSq.SimplifyTolerance *= sg.paramsSq.SimplifyTolerance
	sg.paramsSq.CollideEarly *= sg.paramsSq.CollideEarly
}

func (sg *StreamlineGenerator) samplePoint() vectors.Vec2 {
	// TODO better seeding scheme
	return vectors.Vec2{
		X: sg.rng.Float64() * sg.worldDimensions.X,
		Y: sg.rng.Float64() * sg.worldDimensions.Y,
	}.Add(sg.origin)
}

// getSeed returns a seed point for a streamline.
// Tries this.candidateSeeds first, then samples using this.samplePoint.
func (sg *StreamlineGenerator) getSeed(major bool) vectors.Vec2 {
	// Candidate seeds first
	if sg.SEED_AT_ENDPOINTS && len(sg.candidateSeeds(major)) > 0 {
		for len(sg.candidateSeeds(major)) > 0 {
			seed := sg.candidateSeeds(major)[len(sg.candidateSeeds(major))-1]
			sg.setCandidateSeeds(major, sg.candidateSeeds(major)[:len(sg.candidateSeeds(major))-1])
			if sg.isValidSample(major, seed, sg.paramsSq.Dsep, false) {
				return seed
			}
		}
	}

	seed := sg.samplePoint()
	log.Println("seed", seed)
	i := 0
	for !sg.isValidSample(major, seed, sg.paramsSq.Dsep, false) {
		log.Println("seed", seed, "invalid", i)
		if i >= sg.params.SeedTries {
			return emptyVec2
		}
		seed = sg.samplePoint()
		i++
	}
	return seed
}

func (sg *StreamlineGenerator) isValidSample(major bool, point vectors.Vec2, dSq float64, bothGrids bool) bool {
	// dSq = dSq * point.distanceToSquared(Vector.zeroVector());
	gridValid := sg.grid(major).IsValidSample(point, dSq)
	if bothGrids {
		gridValid = gridValid && sg.grid(!major).IsValidSample(point, dSq)
	}
	return sg.integrator.OnLand(point) && gridValid
}

func (sg *StreamlineGenerator) candidateSeeds(major bool) []vectors.Vec2 {
	if major {
		return sg.candidateSeedsMajor
	}
	return sg.candidateSeedsMinor
}

func (sg *StreamlineGenerator) appendCandidateSeeds(major bool, seeds ...vectors.Vec2) {
	if major {
		sg.candidateSeedsMajor = append(sg.candidateSeedsMajor, seeds...)
	} else {
		sg.candidateSeedsMinor = append(sg.candidateSeedsMinor, seeds...)
	}
}

func (sg *StreamlineGenerator) setCandidateSeeds(major bool, seeds []vectors.Vec2) {
	if major {
		sg.candidateSeedsMajor = seeds
	} else {
		sg.candidateSeedsMinor = seeds
	}
}

func (sg *StreamlineGenerator) streamlines(major bool) [][]vectors.Vec2 {
	if major {
		return sg.streamlinesMajor
	}
	return sg.streamlinesMinor
}

func (sg *StreamlineGenerator) appendStreamlines(major bool, streamlines ...[]vectors.Vec2) {
	if major {
		sg.streamlinesMajor = append(sg.streamlinesMajor, streamlines...)
	} else {
		sg.streamlinesMinor = append(sg.streamlinesMinor, streamlines...)
	}
}

func (sg *StreamlineGenerator) grid(major bool) *GridStorage {
	if major {
		return sg.majorGrid
	}
	return sg.minorGrid
}

func (sg *StreamlineGenerator) pointInBounds(v vectors.Vec2) bool {
	return v.X >= sg.origin.X &&
		v.Y >= sg.origin.Y &&
		v.X < sg.worldDimensions.X+sg.origin.X &&
		v.Y < sg.worldDimensions.Y+sg.origin.Y
}

// doesStreamlineCollideSelf stops spirals from forming, uses 0.5 dcirclejoin so that circles are still joined up.
// TestSample is candidate to pushed on end of streamlineForwards, returns true if streamling collides with itself
// NOTE: Currently unused - bit expensive, used streamlineTurned instead.
func (sg *StreamlineGenerator) doesStreamlineCollideSelf(testSample vectors.Vec2, streamlineForwards, streamlineBackwards []vectors.Vec2) bool {
	// Streamline long enough
	if len(streamlineForwards) > sg.nStreamlineLookBack {
		// Forwards check
		for i := 0; i < len(streamlineForwards)-sg.nStreamlineLookBack; i += sg.nStreamlineStep {
			if testSample.DistanceToSquared(streamlineForwards[i]) < sg.dcollideselfSq {
				return true
			}
		}
		// Backwards check
		for i := 0; i < len(streamlineBackwards); i += sg.nStreamlineStep {
			if testSample.DistanceToSquared(streamlineBackwards[i]) < sg.dcollideselfSq {
				return true
			}
		}
	}
	return false
}

// streamlineTurned tests whether streamline has turned through greater than 180 degrees and
// stops spirals from forming.
func (sg *StreamlineGenerator) streamlineTurned(seed, originalDir, point, direction vectors.Vec2) bool {
	if originalDir.Dot(direction) < 0 {
		// TODO: Optimize!
		perpendicularVector := vectors.Vec2{X: originalDir.Y, Y: -originalDir.X}
		isLeft := point.Sub(seed).Dot(perpendicularVector) < 0
		directionUp := direction.Dot(perpendicularVector) > 0
		return isLeft == directionUp
	}
	return false
}

// streamlineIntegrationStep performs one step of the streamline integration process.
// TODO: this doesn't work well - consider something disallowing one direction (F/B)
// to turn more than 180 degrees.
func (sg *StreamlineGenerator) streamlineIntegrationStep(params *StreamlineIntegration, major bool, collideBoth bool) {
	if params.Valid {
		params.Streamline = append(params.Streamline, params.PreviousPoint)
		nextDirection := sg.integrator.Integrate(params.PreviousPoint, major)

		// Stop at degenerate point
		if nextDirection.LengthSquared() < 0.01 {
			params.Valid = false
			return
		}

		// Make sure we travel in the same direction
		if nextDirection.Dot(params.PreviousDir) < 0 {
			nextDirection = nextDirection.Negate()
		}

		nextPoint := params.PreviousPoint.Add(nextDirection)

		// Visualise stopping points
		// if (this.streamlineTurned(params.seed, params.originalDir, nextPoint, nextDirection)) {
		//     params.valid = false;
		//     params.streamline.push(Vector.zeroVector());
		// }

		if sg.pointInBounds(nextPoint) &&
			sg.isValidSample(major, nextPoint, sg.paramsSq.Dtest, collideBoth) &&
			!sg.streamlineTurned(params.Seed, params.OriginalDir, nextPoint, nextDirection) {
			params.PreviousPoint = nextPoint
			params.PreviousDir = nextDirection
		} else {
			// One more step
			params.Streamline = append(params.Streamline, nextPoint)
			params.Valid = false
		}
	}
}

// integrateStreamline integrates a streamline from a seed point in both directions.
// By simultaneously integrating in both directions we reduce the impact of circles
// not joining up as the error matches at the join.
func (sg *StreamlineGenerator) integrateStreamline(seed vectors.Vec2, major bool) []vectors.Vec2 {
	count := 0
	pointsEscaped := false // True once two integration fronts have moved dlookahead away

	// Whether or not to test validity using both grid storages
	// (Collide with both major and minor)
	collideBoth := sg.rng.Float64() < sg.params.CollideEarly

	d := sg.integrator.Integrate(seed, major)

	forwardParams := StreamlineIntegration{
		Seed:          seed,
		OriginalDir:   d,
		Streamline:    []vectors.Vec2{seed},
		PreviousDir:   d,
		PreviousPoint: seed.Add(d),
		Valid:         true,
	}
	forwardParams.Valid = sg.pointInBounds(forwardParams.PreviousPoint)

	negD := d.Negate()
	backwardParams := StreamlineIntegration{
		Seed:          seed,
		OriginalDir:   negD,
		Streamline:    []vectors.Vec2{},
		PreviousDir:   negD,
		PreviousPoint: seed.Add(negD),
		Valid:         true,
	}
	backwardParams.Valid = sg.pointInBounds(backwardParams.PreviousPoint)

	for count < sg.params.PathIterations && (forwardParams.Valid || backwardParams.Valid) {
		sg.streamlineIntegrationStep(&forwardParams, major, collideBoth)
		sg.streamlineIntegrationStep(&backwardParams, major, collideBoth)

		// Join up circles
		sqDistanceBetweenPoints := forwardParams.PreviousPoint.DistanceToSquared(backwardParams.PreviousPoint)
		if !pointsEscaped && sqDistanceBetweenPoints > sg.paramsSq.Dcirclejoin {
			pointsEscaped = true
		}
		if pointsEscaped && sqDistanceBetweenPoints <= sg.paramsSq.Dcirclejoin {
			forwardParams.Streamline = append(forwardParams.Streamline, forwardParams.PreviousPoint)
			forwardParams.Streamline = append(forwardParams.Streamline, backwardParams.PreviousPoint)
			backwardParams.Streamline = append(backwardParams.Streamline, backwardParams.PreviousPoint)
			break
		}
		count++
	}

	// Reverse backwards streamline
	for i := len(backwardParams.Streamline)/2 - 1; i >= 0; i-- {
		opp := len(backwardParams.Streamline) - 1 - i
		backwardParams.Streamline[i], backwardParams.Streamline[opp] = backwardParams.Streamline[opp], backwardParams.Streamline[i]
	}

	// Append forwards to backwards
	backwardParams.Streamline = append(backwardParams.Streamline, forwardParams.Streamline...)
	return backwardParams.Streamline
}
