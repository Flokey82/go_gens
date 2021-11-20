package main

import (
	"github.com/Flokey82/go_gens/aifiver"
	"log"
)

func main() {
	t := aifiver.NewTraiter()
	setupTraits(t)
	p := t.NewPersonalityFromPreset(aifiver.PresetUnstable)
	for _, t := range p.Expressed {
		log.Println(t.Name)
	}
	p2 := t.NewPersonalityFromPreset(aifiver.PresetFool)
	for _, t := range p2.Expressed {
		log.Println(t.Name)
	}
}

func setupTraits(tt *aifiver.Traiter) {
	traitParanoid := aifiver.NewTrait("Paranoid", aifiver.TTypePersonality, func(p *aifiver.Personality) bool {
		return p.GetFacet(aifiver.FacetAgreTrust) < -6
	})
	traitParanoid.Stats.Skill[aifiver.TSkillDiplomacy] = -5
	traitParanoid.Stats.Skill[aifiver.TSkillIntrigue] = 5
	tt.AddTrait(traitParanoid)

	traitGullable := aifiver.NewTrait("Gullable", aifiver.TTypePersonality, func(p *aifiver.Personality) bool {
		return p.GetFacet(aifiver.FacetAgreTrust)+p.GetFacet(aifiver.FacetAgreCompliance) > 15
	})
	traitGullable.Stats.Skill[aifiver.TSkillDiplomacy] = -5
	traitGullable.Stats.Skill[aifiver.TSkillIntrigue] = -5
	tt.AddTrait(traitGullable)

	aifiver.MarkOppositeTraits(traitParanoid, traitGullable)

	traitChaste := aifiver.NewTrait("Chaste", aifiver.TTypePersonality, func(p *aifiver.Personality) bool {
		return p.GetFacet(aifiver.FacetNeurImpulsiveness) < -5
	})
	traitChaste.Stats.Skill[aifiver.TSkillLearning] = 2
	tt.AddTrait(traitChaste)

	traitLustful := aifiver.NewTrait("Lustful", aifiver.TTypePersonality, func(p *aifiver.Personality) bool {
		return p.GetFacet(aifiver.FacetNeurImpulsiveness) > 5
	})
	traitLustful.Stats.Skill[aifiver.TSkillIntrigue] = 2
	tt.AddTrait(traitLustful)

	aifiver.MarkOppositeTraits(traitChaste, traitLustful)

	traitHonest := aifiver.NewTrait("Honest", aifiver.TTypePersonality, func(p *aifiver.Personality) bool {
		return p.GetFacet(aifiver.FacetNeurImpulsiveness) < -6 &&
			p.GetFacet(aifiver.FacetConsDutifulness) > 6 &&
			p.GetFacet(aifiver.FacetAgreAltruism) > 2 &&
			p.GetFacet(aifiver.FacetAgreStraightforwardness) > 2
	})
	traitHonest.Stats.Skill[aifiver.TSkillDiplomacy] = 2
	traitHonest.Stats.Skill[aifiver.TSkillIntrigue] = -4
	tt.AddTrait(traitHonest)

	traitDeceitful := aifiver.NewTrait("Deceitful", aifiver.TTypePersonality, func(p *aifiver.Personality) bool {
		return p.GetFacet(aifiver.FacetAgreAltruism) < 0 &&
			p.GetFacet(aifiver.FacetAgreStraightforwardness) < 2 &&
			p.GetFacet(aifiver.FacetAgreModesty) < -5
	})
	traitDeceitful.Stats.Skill[aifiver.TSkillDiplomacy] = -2
	traitDeceitful.Stats.Skill[aifiver.TSkillIntrigue] = 4
	tt.AddTrait(traitDeceitful)

	aifiver.MarkOppositeTraits(traitHonest, traitDeceitful)
}
