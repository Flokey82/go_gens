package gencitymap

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/ojrac/opensimplex-go"
)

// https://github.com/ProbableTrain/MapGenerator/blob/f487e4cee321d105d5b0e14258363fa8b4d004ef/src/ts/impl/basis_field.ts

type NoiseParams struct {
	Seed             int64
	globalNoise      bool
	noiseSizePark    float64
	noiseAnglePark   float64
	noiseSizeGlobal  float64
	noiseAngleGlobal float64
}

type TensorField struct {
	basisFields []BasisFieldInterface
	noise       opensimplex.Noise
	parks       [][]vectors.Vec2
	sea         []vectors.Vec2
	river       []vectors.Vec2
	ignoreRiver bool
	smooth      bool
	noiseParams *NoiseParams
}

// NewTensorField creates a new tensor field.
func NewTensorField(n *NoiseParams) *TensorField {
	return &TensorField{
		basisFields: []BasisFieldInterface{},
		noise:       opensimplex.New(n.Seed),
		parks:       [][]vectors.Vec2{},
		sea:         []vectors.Vec2{},
		river:       []vectors.Vec2{},
		ignoreRiver: false,
		smooth:      false,
		noiseParams: n,
	}
}

func (t *TensorField) AddGrid(centre vectors.Vec2, size float64, decay float64, theta float64) {
	t.basisFields = append(t.basisFields, NewGridField(centre, size, decay, theta))
}

func (t *TensorField) AddRadial(centre vectors.Vec2, size float64, decay float64, theta float64) {
	t.basisFields = append(t.basisFields, NewRadialField(centre, size, decay, theta))
}

func (t *TensorField) addField(field BasisFieldInterface) {
	t.basisFields = append(t.basisFields, field)
}

func (t *TensorField) removeField(field BasisFieldInterface) {
	for i, v := range t.basisFields {
		if v == field {
			t.basisFields = append(t.basisFields[:i], t.basisFields[i+1:]...)
			return
		}
	}
}

func (t *TensorField) getCentrePoints() []vectors.Vec2 {
	centrePoints := []vectors.Vec2{}
	for _, field := range t.basisFields {
		centrePoints = append(centrePoints, field.GetCentre())
	}
	return centrePoints
}

func (t *TensorField) samplePoint(point vectors.Vec2) *Tensor {
	if !t.onLand(point) {
		// Degenerate point
		return newZeroTensor()
	}

	// Default field is a grid
	if len(t.basisFields) == 0 {
		return newTensor(1, [2]float64{0, 0})
	}

	tensorAcc := newZeroTensor()
	for _, field := range t.basisFields {
		tensorAcc.add(field.GetWeightedTensor(point, t.smooth), t.smooth)
	}

	// Add rotational noise for parks - range -pi/2 to pi/2
	for _, park := range t.parks {
		if InsidePolygon(point, park) {
			// TODO optimise insidePolygon e.g. distance
			tensorAcc.rotate(t.getRotationalNoise(point, t.noiseParams.noiseSizePark, t.noiseParams.noiseAnglePark))
			break
		}
	}

	if t.noiseParams.globalNoise {
		tensorAcc.rotate(t.getRotationalNoise(point, t.noiseParams.noiseSizeGlobal, t.noiseParams.noiseAngleGlobal))
	}

	return tensorAcc
}

/**
 * Noise Angle is in degrees
 */

func (t *TensorField) getRotationalNoise(point vectors.Vec2, noiseSize float64, noiseAngle float64) float64 {
	return t.noise.Eval2(point.X/noiseSize, point.Y/noiseSize) * noiseAngle * math.Pi / 180
}

func (t *TensorField) onLand(point vectors.Vec2) bool {
	//log.Println("Hack! onLand always returns true")
	return true
	/*
		if t.ignoreRiver {
			return InsidePolygon(point, t.sea)
		}
		return InsidePolygon(point, t.sea) && !InsidePolygon(point, t.river)*/
}

func (t *TensorField) inParks(point vectors.Vec2) bool {
	for _, p := range t.parks {
		if InsidePolygon(point, p) {
			return true
		}
	}
	return false
}

func InsidePolygon(point vectors.Vec2, polygon []vectors.Vec2) bool {
	// TODO optimise insidePolygon e.g. distance
	for i := 0; i < len(polygon); i++ {
		j := (i + 1) % len(polygon)
		if (polygon[i].Y > point.Y) != (polygon[j].Y > point.Y) &&
			point.X < (polygon[j].X-polygon[i].X)*(point.Y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X {
			return true
		}
	}
	return false
}
