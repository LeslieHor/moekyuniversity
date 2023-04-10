package cards

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type TextAnalysis struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
	Tokens []Token
}

type Token struct {
	Surface string
	BaseForm string
	POS []string
	Pronunciation string
	LearningStage LearningStage
	LearningStageString string
}

func (ta *TextAnalysis) Save(filepath string) {
	taJson, err := json.Marshal(ta)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath, taJson, 0644)
	if err != nil {
		panic(err)
	}
}

func (ta *TextAnalysis) Analyse() {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	tokens := t.Tokenize(ta.Text)

	var ts []Token
	for _, token := range tokens {
		bf, _ := token.BaseForm()
		p, _ := token.Pronunciation()
		ts = append(ts, Token{
			Surface: token.Surface,
			BaseForm: bf,
			POS: token.POS(),
			Pronunciation: p,
		})
	}
	ta.Tokens = ts
}

func (cd *CardData) AddLearningStages(ta *TextAnalysis) {
	for i, token := range ta.Tokens {
		c := cd.FindVocabulary(token.BaseForm)
		if c != nil {
			ta.Tokens[i].LearningStage = c.LearningStage
			ta.Tokens[i].LearningStageString = LearningStageToString(c.LearningStage)
		}
	}
}