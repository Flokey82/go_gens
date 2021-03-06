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
	person_traits []string
	big_five      []int
	story         string
}

// Each villager gets a personality archetype which
// effects how social, work, and life events effect them.
func NewPersonality(name string) *Personality {
	p := &Personality{
		person_traits: nil,
		big_five:      nil,
	}
	for i := 0; i <= 5; i++ {
		p.big_five = append(p.big_five, rand.Intn(100))
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
	if p.big_five[0] > 50 {
		p.person_traits = append(p.person_traits, pickRandString(high_o))
	} else {
		p.person_traits = append(p.person_traits, pickRandString(low_o))
	}

	// Conscientiousness
	if p.big_five[1] > 50 {
		p.person_traits = append(p.person_traits, pickRandString(high_c))
	} else {
		p.person_traits = append(p.person_traits, pickRandString(low_c))
	}

	// Extraversion
	if p.big_five[2] > 50 {
		p.person_traits = append(p.person_traits, pickRandString(high_e))
	} else {
		p.person_traits = append(p.person_traits, pickRandString(low_e))
	}

	// Agreeableness
	if p.big_five[3] > 50 {
		p.person_traits = append(p.person_traits, pickRandString(high_a))
	} else {
		p.person_traits = append(p.person_traits, pickRandString(low_a))
	}

	// Neuroticism
	if p.big_five[4] > 50 {
		p.person_traits = append(p.person_traits, pickRandString(high_n))
	} else {
		p.person_traits = append(p.person_traits, pickRandString(low_n))
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
		fmt.Sprintf(pickRandString(t_lines), name, p.person_traits[0], p.person_traits[1]),
		fmt.Sprintf(pickRandString(t_lines), name, p.person_traits[2], p.person_traits[3]),
		fmt.Sprintf(pickRandString(s_lines), name, p.person_traits[4]),
	}, " ")
	return p
}
func (p *Personality) get_backstory() string {
	return p.story
}
