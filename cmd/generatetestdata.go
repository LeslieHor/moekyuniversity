package main

import (
	"flag"
	"time"
	"moekyuniversity/internal/cards"
)

var (
	cardsFile = flag.String("cards-file", "data/cards.json", "Cards file")
)

var Kanji []string = []string{
	"一", "二", "三",
}

func setBurned(c *cards.Card) {
	c.Interval = 8761
	c.LearningInterval = 0
	c.NextReviewDate = ""
}

func setLearned(c *cards.Card) {
	c.Interval = 500
	c.LearningInterval = 0
	c.NextReviewDate = time.Now().Format(time.RFC3339)
}

func setLearning(c *cards.Card) {
	c.Interval = 0
	c.LearningInterval = 3
	c.NextReviewDate = time.Now().Format(time.RFC3339)
}

func setUpNext(c *cards.Card) {
	c.Interval = 0
	c.LearningInterval = 0
	c.NextReviewDate = time.Now().Format(time.RFC3339)
}

func main() {
	flag.Parse()
	cardData := cards.CardData{CardsFile: *cardsFile}

	cs := make(map[int]*cards.Card)
	var i int
	types := []string{"radical", "kanji", "vocabulary"}
	
	// Generate cards

	radicals := []*cards.Card{}
	kanji := []*cards.Card{}
	vocabulary := []*cards.Card{}

	for _, t := range types {
		for _, k := range Kanji {
			i++
			cs[i] = &cards.Card{
				ID: i,
				Object: t,
				Characters: k,
				Meanings: []cards.Meaning {
				   {"one", true, true},
				   {"two", false, false},
				   {"three", false, false},
			   },
			}

			if t == "radical" {
				radicals = append(radicals, cs[i])
			} else if t == "kanji" {
				kanji = append(kanji, cs[i])
			} else if t == "vocabulary" {
				vocabulary = append(vocabulary, cs[i])
			}
		}
	}

	// Create dependencies
	for i, r := range radicals {
		r.AmalgamationSubjectIDs = []int{kanji[i].ID}
	}
	for i, k := range kanji {
		k.ComponentSubjectIDs = []int{radicals[i].ID}
		k.AmalgamationSubjectIDs = []int{vocabulary[i].ID}
	}
	for i, v := range vocabulary {
		v.ComponentSubjectIDs = []int{kanji[i].ID}
	}

	// Custom dependencies
	vocabulary[2].ComponentSubjectIDs = append(vocabulary[2].ComponentSubjectIDs, vocabulary[0].ID)
	vocabulary[2].ComponentSubjectIDs = append(vocabulary[2].ComponentSubjectIDs, vocabulary[1].ID)

	// Set some cards to burned
	setBurned(radicals[0])
	setBurned(radicals[1])
	setBurned(radicals[2])

	// Set some cards to learned
	setLearned(kanji[0])
	setLearned(kanji[1])

	// Set some cards to learning
	setLearning(vocabulary[0])

	// Set some cards to up next
	setUpNext(vocabulary[1])

	// Set review performance
	radicals[0].TotalTimesReviewed = 101
	radicals[0].TotalTimesCorrect = 95
	
	// A card that has no characters
	i++
	cs[i] = &cards.Card{
		ID: i,
		Object: "radical",
		Characters: "",
		Meanings: []cards.Meaning {
		   {"Stick", true, true},
		   {"two", false, false},
		   {"three", false, false},
	   },
	}

	// A card with a character that is different for Japanese vs Chinese
	i++
	cs[i] = &cards.Card{
		ID: i,
		Object: "radical",
		Characters: "令",
		Meanings: []cards.Meaning {
			{"Orders", true, true},
			{"two", false, false},
		},
	}

	cardData.Cards = cs
	cardData.SaveCardMap()
}