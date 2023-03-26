package genarchitecture

type RoofStyle struct {
	Shape string
	BaseStyle
}

func generateRoofStyle(availableMaterials []string) RoofStyle {
	return RoofStyle{
		Shape:     randomString(roofShapes),
		BaseStyle: generateBaseStyle(roofMaterials),
	}
}

func (s RoofStyle) Description() string {
	return "A " + s.Shape + " roof covered with " + s.BaseStyle.Description()
}

const (
	RoofShapeNone      = "none"
	RoofShapeGable     = "gable"
	RoofShapeHip       = "hip"
	RoofShapeDutch     = "dutch"
	RoofShapeJerkin    = "jerkin"
	RoofShapePyramid   = "pyramid"
	RoofShapeMansard   = "mansard"
	RoofShapeBonnet    = "bonnet"
	RoofShapeGambrel   = "gambrel"
	RoofShapeSilikon   = "silikon"
	RoofShapeCurved    = "curved"
	RoofShapeFlat      = "flat"
	RoofShapeSaltbox   = "saltbox"
	RoofShapeButterfly = "butterfly"
	RoofShapeSawtooth  = "sawtooth"
	RoofShapeDormer    = "dormer"
)

var roofShapes = []string{
	RoofShapeNone,
	RoofShapeGable,
	RoofShapeHip,
	RoofShapeDutch,
	RoofShapeJerkin,
	RoofShapePyramid,
	RoofShapeMansard,
	RoofShapeBonnet,
	RoofShapeGambrel,
	RoofShapeSilikon,
	RoofShapeCurved,
	RoofShapeFlat,
	RoofShapeSaltbox,
	RoofShapeButterfly,
	RoofShapeSawtooth,
	RoofShapeDormer,
}

const (
	RoofMaterialNone         = "none"
	RoofMaterialWoodPlanks   = "wood_planks"
	RoofMaterialWoodShingles = "wood_shingles"
	RoofMaterialHide         = "hide"
	RoofMaterialThatch       = "thatch"
	RoofMaterialStraw        = "straw"
	RoofMaterialSlate        = "slate"
	RoofMaterialTile         = "tile"
	RoofMaterialMetal        = "metal"
	RoofMaterialGlass        = "glass"
)

var roofMaterials = []string{
	RoofMaterialNone,
	RoofMaterialWoodPlanks,
	RoofMaterialWoodShingles,
	RoofMaterialHide,
	RoofMaterialThatch,
	RoofMaterialStraw,
	RoofMaterialSlate,
	RoofMaterialTile,
	RoofMaterialMetal,
	RoofMaterialGlass,
}

// TODO: Add weights to materials and other things
