package simvillage

type MgenModel interface {
	make_short_sentence(l int) string
}

type MgenIf interface {
	Text(test string) MgenModel
}

type MarkovGen struct {
	generator MgenIf
	prefix    string
}

func NewMarkovGen() *MarkovGen {
	return &MarkovGen{
		prefix: "wordlists",
	}
}

func (m *MarkovGen) get_chef() string {
	/*
		fdat, err := os.ReadFile(m.prefix + "/cooking.txt")
		if err != nil {
			panic(err)
		}
		text_model := markovify.Text(string(fdat))
		return text_model.make_short_sentence(100)
	*/
	return "cooked stuff"
}

func (m *MarkovGen) get_death() string {
	/*
		fdat, err := os.ReadFile(m.prefix + "/death.txt")
		if err != nil {
			panic(err)
		}
		text_model := markovify.Text(string(fdat))
		return text_model.make_short_sentence(100)
	*/
	return "die"
}
