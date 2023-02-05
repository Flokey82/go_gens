package genreligion

// Classification represents a religion type.
type Classification struct {
	Group string
	Form  string
	Type  string
}

// NewClassification returns a religion classification.
// NOTE: If form or type are empty, a random one is chosen.
func (g *Generator) NewClassification(group, form, rType string) *Classification {
	if form == "" {
		form = g.RandFormFromGroup(group)
	}
	if rType == "" {
		rType = g.RandTypeFromForm(form)
	}
	return &Classification{
		Group: group,
		Form:  form,
		Type:  rType,
	}
}

// HasDeity returns true if the religion has one or multiple deities.
func (c *Classification) HasDeity() bool {
	return c.Form != FormNontheism && c.Form != FormAnimism
}
