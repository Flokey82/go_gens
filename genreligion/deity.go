package genreligion

import (
	"github.com/Flokey82/go_gens/genstory"
)

var DeityMeaningApproaches []string

func init() {
	DeityMeaningApproaches = weightedToArray(GenMeaningApproaches)
}

// Deity represents a deity name.
type Deity struct {
	Name    string
	Meaning *genstory.Generated
}

// FullName returns the full name of the deity (including the meaning, if any).
func (d *Deity) FullName() string {
	if d == nil {
		return ""
	}
	if d.Meaning == nil || d.Meaning.Text == "" {
		return d.Name
	}
	return d.Name + ", The " + d.Meaning.Text
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
		Name:    g.lang.MakeName(),
		Meaning: meaning,
	}, nil
}

// GetDeityWithAntonyms returns a deity name for the given culture.
/*
func (g *Generator) GetDeityWithAntonyms(d *Deity) (*Deity, error) {
	var provided []genstory.TokenReplacement
	for _, token := range d.Meaning.Tokens {
		// Try to find an antonym for this token.
		candidates := genlanguage.AllAntonyms[token.Replacement]
		if len(candidates) == 0 {
			continue
		}

		provided = append(provided, genstory.TokenReplacement{
			Token:       token.Token,
			Replacement: g.Ra(candidates),
		})
	}
	meaning, err := g.GenerateDeityMeaningV2(provided, d.Meaning.Template)
	if err != nil {
		return nil, err
	}
	return &Deity{
		Name:    g.lang.MakeName(),
		Meaning: meaning,
	}, nil
}
*/
