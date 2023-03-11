package genreligion

var DeityMeaningApproaches []string

func init() {
	DeityMeaningApproaches = weightedToArray(GenMeaningApproaches)
}

// Deity represents a deity name.
type Deity struct {
	Name     string
	Meaning  string
	Approach string
}

// FullName returns the full name of the deity (including the meaning, if any).
func (d *Deity) FullName() string {
	if d == nil {
		return ""
	}
	if d.Meaning == "" {
		return d.Name
	}
	return d.Name + ", The " + d.Meaning
}

// GetDeity returns a deity name for the given culture.
// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/religions-generator.js
func (g *Generator) GetDeity() (*Deity, error) {
	return g.GetDeityWithApproach(g.RandDeityGenMethod())
}

// GetDeityWithApproach returns a deity name for the given culture.
func (g *Generator) GetDeityWithApproach(approach string) (*Deity, error) {
	meaning, err := g.GenerateDeityMeaning(approach)
	if err != nil {
		return nil, err
	}
	return &Deity{
		Name:     g.lang.MakeName(),
		Meaning:  meaning,
		Approach: approach,
	}, nil
}
