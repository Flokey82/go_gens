package simvillage

type MgenModel interface {
	makeShortSentence(l int) string
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

func (m *MarkovGen) getChef() string {
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

func (m *MarkovGen) getDeath() string {
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
