package cards

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type TextAnalysis struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Text   string  `json:"text"`
	Tokens []Token `json:"-"`
}

type Token struct {
	Surface             string   // As displayed to user
	BaseForm            string   // Actual word
	PartsOfSpeech       []string // Parts of speech
	Pronunciation       string   // Pronunciation in hiragana
	CardId              int      // ID of the matching card (if any)
	LearningStage       LearningStage
	LearningStageString string
	DictionaryEntries   []DictionaryEntry // Dictionary entries that match this token
	Token               tokenizer.Token   // The original token data
	Card                *Card             // The card that matches this token
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

func ConvertToken(token tokenizer.Token) Token {
	bf, _ := token.BaseForm()
	p, _ := token.Pronunciation()
	pos := token.POS()

	to := Token{
		Surface:       token.Surface,
		BaseForm:      bf,
		PartsOfSpeech: pos,
		Pronunciation: p,
		Token:         token,
	}

	return to
}

func (ta *TextAnalysis) Analyse(cd *CardData) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}

	log.Printf("Analyzing text: %s", ta.Name)
	tokens := t.Analyze(ta.Text, tokenizer.Normal)
	log.Printf("Found %d tokens", len(tokens))

	var ts []Token
	startTime := time.Now()
	dictCache := make(map[string][]DictionaryEntry)
	for _, token := range tokens {
		to := ConvertToken(token)

		c := cd.FindVocabulary(to.BaseForm)
		to.Card = c
		if c == nil {
			// Cache the dictionary entries that match this token
			// Improves performance by approximately 3x (depending on the text)
			key := to.BaseForm + to.Pronunciation
			if _, ok := dictCache[key]; ok {
				to.DictionaryEntries = dictCache[key]
			} else {
				to.AddDictionaryEntry(cd)
				dictCache[key] = to.DictionaryEntries
			}
		}

		ts = append(ts, to)
	}
	endTime := time.Now()
	log.Printf("Processed %d tokens in %s", len(tokens), endTime.Sub(startTime))
	// debugJsonPrint(ts)
	ta.Tokens = ts
}

// Currently unused
// func debugJsonPrint(v interface{}) {
// 	b, _ := json.Marshal(v)
// 	log.Printf("JSON: %s", b)
// }
