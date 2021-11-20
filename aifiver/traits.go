package aifiver

// Trait represents a character trait that influences behavior and decisions.
type Trait struct {
	Type          TraitType
	Opposite      *Trait
	Name          string
	Description   string
	IsExpressedBy func(*Personality) bool
	*Stats        // Contains all modifiers.
}

// NewTrait returns a new trait.
func NewTrait(name string, t TraitType, cond func(p *Personality) bool) *Trait {
	return &Trait{
		Name:          name,
		Type:          t,
		Stats:         NewStats(),
		IsExpressedBy: cond,
	}
}

// MarkOppositeTraits defines two traits as their polar opposite.
func MarkOppositeTraits(a, b *Trait) {
	a.Opposite = b
	b.Opposite = a
}

// mergeTraitSets returns a combined set of traits.
// TODO: Filter out opposites that can not co-exist.
func mergeTraitSets(a, b []*Trait) []*Trait {
	seenTrait := make(map[*Trait]bool)

	var result []*Trait
	for _, t := range a {
		if seenTrait[t] {
			continue
		}
		seenTrait[t] = true
		result = append(result, t)
	}

	for _, t := range b {
		if seenTrait[t] {
			continue
		}
		seenTrait[t] = true
		result = append(result, t)
	}

	return result
}

// TraitType defining how a trait can be aquired
type TraitType int

// Various trait types.
const (
	TTypePersonality TraitType = iota
	TTypeEducation
	TTypeCongenital
	TTypeCongenitalLeveled
	TTypePhysical
	TTypeLifestyle
	TTypeLifestyleLeveled
	TTypeCommander
	TTypeInfamous
	TTypeCopingMechanism
	TTypeChildhood
	TTypeHealth
	TTypeDiseases
	TTypeDynasty
	TTypeDescendant
	TTypeDecision
	TTypeOther
)

// Traiter is an index used for assigning traits to personalities.
type Traiter struct {
	Traits []*Trait
}

// NewTraiter returns a new Traiter index.
func NewTraiter() *Traiter {
	return &Traiter{}
}

// NewPersonality returns a new, pretty boring personality.
func (tt *Traiter) NewPersonality() *Personality {
	return NewPersonality(tt)
}

// NewPersonalityFromPreset returns a new personality based on the given preset.
func (tt *Traiter) NewPersonalityFromPreset(ps map[Facet]int) *Personality {
	return NewPersonalityFromPreset(tt, ps)
}

// AddTrait registers the given trait.
func (tt *Traiter) AddTrait(t *Trait) {
	tt.Traits = append(tt.Traits, t)
}
