package simnpcs

type TopicPool struct {
	ByProfession map[*Profession][]Topic
	ByLocation   map[*Location][]Topic
	ByEducation  map[*Education][]Topic
}

func NewTopicPool() *TopicPool {
	return &TopicPool{
		ByProfession: make(map[*Profession][]Topic),
		ByLocation:   make(map[*Location][]Topic),
		ByEducation:  make(map[*Education][]Topic),
	}
}

type Topic struct {
	ID   uint64
	Name string

	RelatedTo TopicRelation
	// Connected to:
	//// Place
	//// Time
	//// People
	//// Profession

	// Facts related to this topic
	Facts []*Fact
}

func NewTopic(id uint64, name string) *Topic {
	return &Topic{
		ID:   id,
		Name: name,
	}
}

type Fact struct {
	ID    uint64
	Topic *Topic
	// Chance of knowledge:
	//// Common knowledge (90%)
	//// Specialist knowledge (50%)
	//// Expert knowledge (20%)
	//// First hand knowledge (0%)

	// Locked by:
	//// Profession | Experience
	//// Passion | Interest
}

type TopicRelation struct {
	Professions []*Profession
	Characters  []*Character
	//Passions    []*Passion
	Locations []*Location
}
