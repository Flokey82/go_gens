// Package aifiver implements a simplified representation of the five factor model.
package aifiver

// Factor represents one of the five factors in the five factor model.
type Factor int

// The five factors.
const (
	FactorOpenness Factor = iota
	FactorConscientiousness
	FactorExtraversion
	FactorAgreeableness
	FactorNeuroticism
)

// Facet represents specific facets within the five factors.
type Facet int

// The facets.
const (
	FacetOpenFantasy             Facet = iota // Proneness to imagination, day-dreaming, and creating
	FacetOpenAesthetics                       // Appreciation for beauty in e.g. art, music, poetry, or nature
	FacetOpenFeelings                         // Receptivity to and intensity of experienced emotions
	FacetOpenAdventurousness                  // The tendency to choose novelty over the familiar
	FacetOpenIdeas                            // The degree of interest and curiosity in entertaining new thoughts and ideas
	FacetOpenValues                           // Willingness to re-evaluate norms and values
	FacetConsCompetence                       // Belief in one's own capacity to handle life's many challenges
	FacetConsOrder                            // Degree of neatness and orderliness
	FacetConsDutifulness                      // How strongly ethical principles guide action
	FacetConsAchievementStriving              // Aspiration-level, the willingness to work towards goals
	FacetConsSelfDicipline                    // The ability to follow through on tasks despite boredom
	FacetConsCautiousness                     // The care and thought put into actions
	FacetExtrWarmth                           // The degree of displayed affection and closeness in relationships
	FacetExtrGregariousness                   // The tendency of seeking the company of others
	FacetExtrAssertiveness                    // The degree of dominance in social interaction
	FacetExtrActivity                         // The level of energy and activity in daily life
	FacetExtrExcitementSeeking                // The need for thrills and intense stimulation
	FacetExtrPositiveEmotions                 // The tendency to be happy, excited, and cheerful
	FacetAgreTrust                            // The general level of wariness or suspicion in contact with other people
	FacetAgreStraightforwardness              // Degree of sincerity vs shrewdness
	FacetAgreAltruism                         // Active concern for the well-being of others
	FacetAgreCompliance                       // Inhibiting vs expressing agression towards others in conflict
	FacetAgreModesty                          // Degree of humility vs arrogance
	FacetAgreTenderMindedness                 // Propensity to empathize with others
	FacetNeurAnxiety                          // Proneness to worry and rumination
	FacetNeurAngryHostility                   // The readiness to experience frustration, anger, and bitterness
	FacetNeurDepression                       // The tendency for guilt, sadness, lonliness, and hopelessness
	FacetNeurSelfConciousness                 // Sensitivity in social situations, such as ridicule, rejection, or awkwardness
	FacetNeurImpulsiveness                    // The ability to tolerate frustration and to control urges, cravings, and desires
	FacetNeurVulnerability                    // The ability to cope with stress
)

// Factor returns the factor associated with the given facet.
func (f *Facet) Factor() Factor {
	return facetToFactor[*f]
}

// facetToFactor maps a facet to its respective factor.
var facetToFactor = map[Facet]Factor{
	FacetOpenFantasy:             FactorOpenness,
	FacetOpenAesthetics:          FactorOpenness,
	FacetOpenFeelings:            FactorOpenness,
	FacetOpenAdventurousness:     FactorOpenness,
	FacetOpenIdeas:               FactorOpenness,
	FacetOpenValues:              FactorOpenness,
	FacetConsCompetence:          FactorConscientiousness,
	FacetConsOrder:               FactorConscientiousness,
	FacetConsDutifulness:         FactorConscientiousness,
	FacetConsAchievementStriving: FactorConscientiousness,
	FacetConsSelfDicipline:       FactorConscientiousness,
	FacetConsCautiousness:        FactorConscientiousness,
	FacetExtrWarmth:              FactorExtraversion,
	FacetExtrGregariousness:      FactorExtraversion,
	FacetExtrAssertiveness:       FactorExtraversion,
	FacetExtrActivity:            FactorExtraversion,
	FacetExtrExcitementSeeking:   FactorExtraversion,
	FacetExtrPositiveEmotions:    FactorExtraversion,
	FacetAgreTrust:               FactorAgreeableness,
	FacetAgreStraightforwardness: FactorAgreeableness,
	FacetAgreAltruism:            FactorAgreeableness,
	FacetAgreCompliance:          FactorAgreeableness,
	FacetAgreModesty:             FactorAgreeableness,
	FacetAgreTenderMindedness:    FactorAgreeableness,
	FacetNeurAnxiety:             FactorNeuroticism,
	FacetNeurAngryHostility:      FactorNeuroticism,
	FacetNeurDepression:          FactorNeuroticism,
	FacetNeurSelfConciousness:    FactorNeuroticism,
	FacetNeurImpulsiveness:       FactorNeuroticism,
	FacetNeurVulnerability:       FactorNeuroticism,
}

// Note to self:
// Here are some interesting approaches..
// - https://github.com/JacobStoneman/NPCProcGen/blob/master/Library/Collab/Download/Assets/Scripts/NPCController.cs

type Model interface {
	Get(factor Factor) int
	GetFacet(facet Facet) int
}

// SmallModel represents a simplified five factor model.
// Please note that this is not granular enough to simulate real personalities.
type SmallModel [5]int

func (p *SmallModel) Get(factor Factor) int {
	return p[factor]
}

func (p *SmallModel) GetFacet(facet Facet) int {
	return p.Get(facet.Factor())
}

type BigModel struct {
	Factors [5]int
	Facets  [30]int
}

func (p *BigModel) Get(f Factor) int {
	return p.Factors[f]
}

func (p *BigModel) GetFacet(f Facet) int {
	return p.Facets[f]
}

// Compatibility returns a compatibility value for the given personalities.
// TODO: Provide function to calculate conflict potential of larger groups.
// https://github.com/FZarattini/PirateShip/blob/master/PirateShip/Assets/Scripts/AI/Empathy.cs
func Compatibility(p1, p2 Model) int {
	// Likelyhood of initiating communication is low if two introverts meet.
	ex := (p1.Get(FactorExtraversion) + p2.Get(FactorExtraversion)) / 2
	if ex <= 0 {
		// Well, that's two people that don't want to interact.
		return 0
	}

	// Potential for conflict if the total value for agreeableness is below zero.
	ag := (p1.Get(FactorAgreeableness) + p2.Get(FactorAgreeableness)) / 2
	if ag < 0 {
		// We are in a very disagreeable situation: Conflict!
		return -1
	}
	return (ex + ag) / 2
}
