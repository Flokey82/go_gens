package genarchitecture

type WindowStyle struct {
	Shape string
	Size  int
	BaseStyle
	Glass    GlassStyle
	Curtains CurtainsStyle
}

func (s WindowStyle) LightValue() int {
	// No windows, no light
	if s.Size == SizeNone {
		return 0
	}

	baseValue := s.Size * 2
	switch s.Glass.Type {
	case GlassTypeClear:
		baseValue += 1
	case GlassTypeObscured:
		baseValue /= 2
	case GlassTypeStained, GlassTypeBottle, GlassTypeFrosted, GlassTypeTinted:
		baseValue -= 1
	}
	return s.Glass.Thickness
}

var shapeWindows = shapeDoors

type GlassStyle struct {
	Thickness int
	Color     string
	Type      string
	Shape     string
	BaseStyle
}

const (
	GlassTypeNone      = "none"
	GlassTypeClear     = "clear"
	GlassTypeFrosted   = "frosted"
	GlassTypeEtched    = "etched"
	GlassTypeTinted    = "tinted"
	GlassTypePatterned = "patterned"
	GlassTypeObscured  = "obscured"
	GlassTypeStained   = "stained"
	GlassTypeBottle    = "bottle"
)

type CurtainsStyle struct {
	Thickness int
	Size      int
	Shape     string
	BaseStyle
}

const (
	FabricTypeNone     = "none"
	FabricTypeSilk     = "silk"
	FabricTypeCotton   = "cotton"
	FabricTypeLinen    = "linen"
	FabricTypeWool     = "wool"
	FabricTypeVelvet   = "velvet"
	FabricTypeLeather  = "leather"
	FabricTypeFur      = "fur"
	FabricTypeSuede    = "suede"
	FabricTypeDenim    = "denim"
	FabricTypeCanvas   = "canvas"
	FabricTypeBurlap   = "burlap"
	FabricTypeCorduroy = "corduroy"
	FabricTypeChiffon  = "chiffon"
	FabricTypeSatin    = "satin"
	FabricTypeTaffeta  = "taffeta"
	FabricTypeTweed    = "tweed"
	FabricTypeLace     = "lace"
	FabricTypeOrganza  = "organza"
	FabricTypeJersey   = "jersey"
	FabricTypePique    = "pique"
	FabricTypePoplin   = "poplin"
)

func generateWindowStyle(availableMaterials []string) WindowStyle {
	return WindowStyle{
		Shape:     randomString(shapeWindows),
		Size:      randomInt(1, 3),
		BaseStyle: generateBaseStyle(availableMaterials),
		Glass: GlassStyle{
			Thickness: randomInt(1, 3),
			Color:     randomString(colors),
			Type:      randomString([]string{GlassTypeClear, GlassTypeObscured, GlassTypeStained, GlassTypeBottle, GlassTypeFrosted, GlassTypeTinted}),
			Shape:     randomString(shapeWindows),
			BaseStyle: generateBaseStyle(availableMaterials),
		},
		Curtains: CurtainsStyle{
			Thickness: randomInt(1, 3),
			Size:      randomInt(1, 3),
			Shape:     randomString(shapeWindows),
			BaseStyle: generateBaseStyle(availableMaterials),
		},
	}
}

var colors = []string{
	"red",
	"orange",
	"yellow",
	"green",
	"blue",
	"indigo",
	"violet",
	"black",
	"white",
	"grey",
	"brown",
}
