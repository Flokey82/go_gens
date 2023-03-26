package genarchitecture

import "github.com/Flokey82/go_gens/genlanguage"

const (
	ComplexityNone      = 0
	ComplexitySimple    = 1
	ComplexityIntricate = 2
	ComplexityDetailed  = 3
	ComplexityMasterful = 4
)

func complexityToString(c int) string {
	switch c {
	case ComplexityNone:
		return "plain"
	case ComplexitySimple:
		return "simple"
	case ComplexityIntricate:
		return "intricate"
	case ComplexityDetailed:
		return "detailed"
	case ComplexityMasterful:
		return "masterful"
	}
	return "unknown"
}

type Decoration struct {
	Type       string
	Complexity int
	Motif      string
}

func (d Decoration) Description() string {
	if d.Type == DecorationTypeNone {
		return "undecorated"
	}
	return genlanguage.GetDegreeAdverbFromAdjective(complexityToString(d.Complexity)) + " " + d.Type +
		" depicting " + d.Motif
}

func genDecoration() *Decoration {
	dec := &Decoration{
		Type: randomString(decorationTypes),
	}
	if dec.Type == DecorationTypeNone {
		return dec
	}

	dec.Complexity = randomInt(ComplexityNone, ComplexityMasterful+1)
	dec.Motif = randomString(motifs) + " " + randomString(motifsActions)
	return dec
}

var motifs = []string{
	"birds",
	"butterflies",
	"cats",
	"dogs",
	"dragons",
	"flowers",
	"geometric patterns",
	"hearts",
	"leaves",
	"lions",
	"monkeys",
	"mountains",
	"people",
	"pigs",
	"rivers",
}

var motifsActions = []string{
	"fighting",
	"playing",
	"sleeping",
	"swimming",
	"walking",
	"working",
}

const (
	DecorationTypeNone    = "undecorated"
	DecorationTypePainted = "painted"
	DecorationTypeCarved  = "carved"
	DecorationTypeInlaid  = "inlaid"
	DecorationTypeGilded  = "gilded"
)

var decorationTypes = []string{
	DecorationTypeNone,
	DecorationTypePainted,
	DecorationTypeCarved,
	DecorationTypeInlaid,
	DecorationTypeGilded,
}
