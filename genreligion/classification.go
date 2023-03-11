package genreligion

// Classification represents a religion type.
type Classification struct {
	Group string
	Form  string
	Type  string
}

// NewClassification returns a religion classification.
// NOTE: If form or type are empty, a random one is chosen.
func (g *Generator) NewClassification(group string) *Classification {
	return g.NewClassificationWithForm(group, g.RandFormFromGroup(group))
}

// NewClassificationWithForm returns a religion classification.
func (g *Generator) NewClassificationWithForm(group, form string) *Classification {
	return &Classification{
		Group: group,
		Form:  form,
		Type:  g.RandTypeFromForm(form),
	}
}

// HasDeity returns true if the religion has one or multiple deities.
func (c *Classification) HasDeity() bool {
	return c.Form != FormNontheism && c.Form != FormAnimism
}
