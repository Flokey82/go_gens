package simnpcs

// TopicPool represents a pool of topics npc's can talk about.
type TopicPool struct {
	ByProfession map[*Profession][]Topic
	ByLocation   map[*Location][]Topic
	ByEducation  map[*Education][]Topic
}

// NewTopicPool returns a new, empty topic pool.
func NewTopicPool() *TopicPool {
	return &TopicPool{
		ByProfession: make(map[*Profession][]Topic),
		ByLocation:   make(map[*Location][]Topic),
		ByEducation:  make(map[*Education][]Topic),
	}
}

// Topic represents a topic of interest.
type Topic struct {
	ID        uint64        // Unique id
	Name      string        // Name of the topic
	Facts     []*Fact       // Facts related to this topic
	RelatedTo TopicRelation // Topics related to this topic
	// Connected to:
	// - Place
	// - Time
	// - People
	// - Profession
}

// NewTopic returns a new topic.
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
	// - Common knowledge (90%)
	// - Specialist knowledge (50%)
	// - Expert knowledge (20%)
	// - First hand knowledge (0%)
	// Locked by:
	// - Profession | Experience
	// - Passion | Interest
}

// TopicRelation represents locations, profession, etc. related to a topic.
type TopicRelation struct {
	Professions []*Profession
	Characters  []*Character
	Locations   []*Location
	// Passions    []*Passion
}
