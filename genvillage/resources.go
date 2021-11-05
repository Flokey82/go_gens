package genvillage

// Production provides information on transformation per cycle.
type Production struct {
	Requires map[string]int
	Provides map[string]int
}

// NewProduction returns a new production index.
// This represents somewhat a business or process that requires a
// number of resources and in turn provides other resources.
//
// Example:
// A bakery requires water, flour, fire-wood, and one worker to produce bread.
func NewProduction() *Production {
	return &Production{
		Requires: make(map[string]int),
		Provides: make(map[string]int),
	}
}

// GetMissing returns all resources that are needed and not provided.
func (p *Production) GetMissing() map[string]int {
	return compMaps(p.Requires, p.Provides)
}

// GetExcess returns all resources that are not needed but provided.
func (p *Production) GetExcess() map[string]int {
	return compMaps(p.Provides, p.Requires)
}

func (p *Production) reset() {
	p.Requires = make(map[string]int)
	p.Provides = make(map[string]int)
}
