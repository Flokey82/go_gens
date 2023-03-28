package simvillage

import (
	"fmt"
	"math/rand"
	"strings"
)

func pickRandString(strs []string) string {
	return strs[rand.Intn(len(strs))]
}

type Personality struct {
	personTraits []string
	bigFive      []int
	story        string
}

// Each villager gets a personality archetype which
// effects how social, work, and life events effect them.
func NewPersonality(name string) *Personality {
	p := &Personality{
		personTraits: nil,
		bigFive:      nil,
	}
	for i := 0; i <= 5; i++ {
		p.bigFive = append(p.bigFive, rand.Intn(100))
	}

	// Openness
	high_o := []string{"curious", "creative", "artsy"}
	low_o := []string{"cautious", "dogmatic"}

	// Conscientiousness
	high_c := []string{"organized", "efficient"}
	low_c := []string{"easy-going", "careless", "cheerful"}

	// Extraversion
	high_e := []string{"extroverted", "outgoing", "talkative"}
	low_e := []string{"introverted", "quiet", "shy"}

	// Agreeableness
	high_a := []string{"friendly", "compassionate"}
	low_a := []string{"difficult", "detatched", "challenging"}

	// Neuroticism
	high_n := []string{"sensitive", "neurotic"}
	low_n := []string{"secure", "confident"}

	// Openness
	if p.bigFive[0] > 50 {
		p.personTraits = append(p.personTraits, pickRandString(high_o))
	} else {
		p.personTraits = append(p.personTraits, pickRandString(low_o))
	}

	// Conscientiousness
	if p.bigFive[1] > 50 {
		p.personTraits = append(p.personTraits, pickRandString(high_c))
	} else {
		p.personTraits = append(p.personTraits, pickRandString(low_c))
	}

	// Extraversion
	if p.bigFive[2] > 50 {
		p.personTraits = append(p.personTraits, pickRandString(high_e))
	} else {
		p.personTraits = append(p.personTraits, pickRandString(low_e))
	}

	// Agreeableness
	if p.bigFive[3] > 50 {
		p.personTraits = append(p.personTraits, pickRandString(high_a))
	} else {
		p.personTraits = append(p.personTraits, pickRandString(low_a))
	}

	// Neuroticism
	if p.bigFive[4] > 50 {
		p.personTraits = append(p.personTraits, pickRandString(high_n))
	} else {
		p.personTraits = append(p.personTraits, pickRandString(low_n))
	}

	s_lines := []string{
		"%s is also %s.",
		"%s has a %s side.",
	}

	t_lines := []string{
		"%s is %s and %s.",
		"Friends know %s as %s and %s.",
		"Dont overlook %ss %s and %s side.",
	}

	p.story = strings.Join([]string{
		fmt.Sprintf(pickRandString(t_lines), name, p.personTraits[0], p.personTraits[1]),
		fmt.Sprintf(pickRandString(t_lines), name, p.personTraits[2], p.personTraits[3]),
		fmt.Sprintf(pickRandString(s_lines), name, p.personTraits[4]),
	}, " ")
	return p
}
func (p *Personality) getBackstory() string {
	return p.story
}
