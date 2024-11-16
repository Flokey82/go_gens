package genempire

import (
	"testing"

	"github.com/Flokey82/go_gens/genlanguage"
	"github.com/Flokey82/go_gens/genstory"
)

func TestNewGenerator(t *testing.T) {
	seed := int64(1234)
	lang := genlanguage.GenLanguage(seed)
	g := NewGenerator(seed, lang)
	if g == nil {
		t.Errorf("NewGenerator(%d, %v) = nil, want not nil", seed, lang)
	}

	for i := 0; i < 20; i++ {
		// Let's set up a number of possible tokens and randomly pick one or more.
		tokenPlace := genstory.TokenReplacement{
			Token:       TokenPlace,
			Replacement: lang.MakeCityName(),
		}

		tokenFoundingFigure := genstory.TokenReplacement{
			Token:       TokenFoundingFigure,
			Replacement: lang.GetLastName(),
		}

		tokenRandom := genstory.TokenReplacement{
			Token:       TokenRandom,
			Replacement: lang.MakeName(),
		}

		availableTokens := []genstory.TokenReplacement{
			tokenPlace,
			tokenFoundingFigure,
			tokenRandom,
		}
		// Generate 3 variations of the name.
		for j := 0; j < 3; j++ {
			// Generate the name and Method.
			gen, err := g.GenEmpireName(availableTokens)
			if err != nil {
				t.Errorf("GenEmpireName(%v) = %v, %v, want not error", availableTokens, gen, err)
			}
			t.Logf("%d[%d]: %s", i, j, gen.Text)
		}
	}
}
