package genthesaurus

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"
)

// ConvertEAJSON converts the original ea-thesaurus.json file to a more
// convenient format.
//
// The original uses a map[string][]map[string]string structure, where
// every map in the slice only has one key-value pair.
//
// See: https://github.com/dariusk/ea-thesaurus
func ConvertEAJSON(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// The original json is structured as follows:
	//
	//	{
	//	 "mapkey": [{"word": "1"}],
	//	 "mapkey": [{"word": "1"}, {"word": "2"}],
	//	}
	t := make(map[string][]map[string]string)
	err = json.NewDecoder(f).Decode(&t)
	if err != nil {
		return err
	}

	// Fill the exorter and write it as a slightly
	// better structured json.
	export := make(map[string]map[string]int)
	for k, v := range t {
		export[k] = make(map[string]int)
		for _, e := range v {
			for k2, v2 := range e {
				count, err := strconv.Atoi(v2)
				if err != nil {
					return err
				}
				export[k][k2] = count
			}
		}
	}

	// Write the export as json.
	f, err = os.Create("data/ea-thesaurus-export.json")
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	return enc.Encode(export)
}

// Thesaurus is a thesaurus.
type Thesaurus struct {
	Entries []*Entry          // All entries.
	ByWord  map[string]*Entry // Entries by word.
}

// New returns a new thesaurus.
func New() *Thesaurus {
	return &Thesaurus{
		ByWord: make(map[string]*Entry),
	}
}

// NewFromJSON returns a new thesaurus from a json file.
func NewFromJSON(path string) (*Thesaurus, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	js := make(map[string]map[string]int)
	if err = json.NewDecoder(f).Decode(&js); err != nil {
		return nil, err
	}

	// New thesaurus.
	t := New()
	for word, associations := range js {
		for association, count := range associations {
			t.AddAssociation(word, association, count)
		}
	}
	t.Sort()
	return t, nil
}

// Add adds a word to the thesaurus.
// If the word exists, the supplied tags are appended.
func (t *Thesaurus) Add(word string, tags ...string) *Entry {
	if e, ok := t.ByWord[word]; ok {
		e.Tags.Append(tags...)
		return e
	}

	e := &Entry{
		Word: word,
		Tags: tags,
	}
	t.Entries = append(t.Entries, e)
	t.ByWord[word] = e
	return e
}

// AddAssociation adds an association to the thesaurus.
//
// 'word' is the word to add the association to.
// 'association' is the associated word to add.
// 'count' is the number of times the association was found / the strength of the association.
// 'relTags' are tags that describe the relationship between the word and the association.
func (t *Thesaurus) AddAssociation(word, association string, count int, relTags ...string) {
	e := t.Add(word)                       // Add the word if it doesn't exist.
	a := t.Add(association)                // Add the association if it doesn't exist.
	e.AddAssociation(a, count, relTags...) // Add the association to the word.
	a.AddAssociation(e, 1, relTags...)     // Add a backreference to the association.
}

// Sort sorts the thesaurus by word and the associations by count.
func (t *Thesaurus) Sort() {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Word < t.Entries[j].Word
	})

	for _, e := range t.Entries {
		e.SortAssociations()
	}
}

// Log logs the thesaurus.
func (t *Thesaurus) Log() {
	for _, e := range t.Entries {
		log.Println(e.Word)
		for _, a := range e.Associations {
			log.Printf("\t%s (%d)", a.Word, a.Count)
		}
	}
}

// Entry is a single entry in the thesaurus.
type Entry struct {
	Word         string         // The word.
	Associations []*Association // Associated words.
	Tags         Tags           // Tags (e.g. "noun", "verb", "color", "animal")
}

// AddAssociation adds an association to the entry (if it doesn't exist).
// If it does exist, the count is increased.
//
// 'word' is the associated word to add.
// 'count' is the number of times the association was found / the strength of the association.
// 'relTags' are tags that describe the relationship between the word and the association.
func (e *Entry) AddAssociation(word *Entry, count int, relTags ...string) {
	for _, a := range e.Associations {
		if a.Entry == word {
			a.Count += count
			a.Relation.Append(relTags...)
			return
		}
	}

	e.Associations = append(e.Associations, &Association{
		Entry:    word,
		Count:    count,
		Relation: relTags,
	})
}

// SortAssociations sorts the associations by count.
func (e *Entry) SortAssociations() {
	sort.Slice(e.Associations, func(i, j int) bool {
		return e.Associations[i].Count > e.Associations[j].Count
	})
}

// Association is a single word association in the thesaurus.
type Association struct {
	*Entry        // Reference to the associated entry.
	Count    int  // Number of times the association was found.
	Relation Tags // Relation type of the association (e.g. "antonym", "synonym")
}

// Tags is a list of tags.
type Tags []string

// Add adds a tag to the list (if it doesn't exist).
func (t *Tags) Add(tag string) {
	for _, t2 := range *t {
		if t2 == tag {
			return
		}
	}

	*t = append(*t, tag)
}

// Append adds a list of tags to the list (if they don't exist).
func (t *Tags) Append(tags ...string) {
	for _, tag := range tags {
		t.Add(tag)
	}
}

// Sort sorts the list of tags.
func (t *Tags) Sort() {
	sort.Slice(*t, func(i, j int) bool {
		return (*t)[i] < (*t)[j]
	})
}

// Has returns true if the list of tags contains the tag.
func (t *Tags) Has(tag string) bool {
	for _, t2 := range *t {
		if t2 == tag {
			return true
		}
	}
	return false
}
