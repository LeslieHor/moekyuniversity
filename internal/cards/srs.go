package cards

import (
	"log"
	"time"
	"html/template"
)

type SrsData struct {
	DueCount int
	LearningCount int
	Card *Card
	MeaningMnemonicHtml template.HTML
	ReadingMnemonicHtml template.HTML
}

func (cd *CardData) GetNextSrsCard() SrsData {
	// Get all cards that are due
	c := cd.ToList()
	dueCards := filterCardsByDueBefore(c, time.Now())
	dueCards = filterOutCardsByTag(dueCards, "suspended")

	// Prioritise cards that are new
	// Prioritise cards that are in the learning stage
	// Don't sort by due date to add randomness to review order
	learningCards := filterCardsByLearningStage(dueCards, Learning)
	learnedCards := filterCardsByLearningStage(dueCards, Learned)

	upNextCards := cd.GetUpNextCards()

	// Up next cards are placed after you've reviewed everything
	srsDueCards := append(learningCards, learnedCards...)

	// Display the number of reviews.
	// If there are no due cards, then display the number of up next cards.
	// As this means the user has reviewed everything and is now looking at new cards.
	var l int
	if len(srsDueCards) == 0 {
		l = len(upNextCards)
	} else {
		l = len(srsDueCards)
	}

	srsDueCards = append(srsDueCards, upNextCards...)

	// Get the first card
	var card *Card
	if len(srsDueCards) == 0 {
		srsData := SrsData{
			DueCount: 0,
			LearningCount: 0,
			Card: nil,
			MeaningMnemonicHtml: template.HTML(""),
			ReadingMnemonicHtml: template.HTML(""),
		}
		return srsData
	}

	card = srsDueCards[0]
	// Create SRS data
	srsData := SrsData{
		DueCount: l,
		LearningCount: len(learningCards),
		Card: card,
		MeaningMnemonicHtml: template.HTML(
			customHtmlTagsToSpan(
				card.MeaningMnemonic)),
		ReadingMnemonicHtml: template.HTML(
			customHtmlTagsToSpan(
				card.ReadingMnemonic)),
	}

	return srsData
}

func (c *Card) CorrectAnswer() {
	// Check the next review date is in the past, otherwise this is a mistaken endpoint hit.
	t, err := time.Parse(time.RFC3339, c.NextReviewDate)
	if err != nil {
		log.Printf("Card %d has an invalid NextReviewDate: %s", c.ID, c.NextReviewDate)
		panic(err)
	}
	if time.Now().Before(t) {
		log.Printf("Card %d was reviewed too early. Next review date is %s", c.ID, c.NextReviewDate)
		return
	}
	
	c.ProcessCorrectAnswer()
}

func (c *Card) ProcessCorrectAnswer() {
	if c.LearningStage == Learning { // Learning stage
		c.LearningInterval *= 2

		// If the LearningInterval is more than 24 hours, then the card has graduated to the learned
		if c.LearningInterval >= 24 {
			c.LearningStage = Learned
			// Set the Interval to the max of LearningInterval and Interval
			// This accounts for cards that have been learned before, but are now being learned again due to forgetting the answer.
			if c.LearningInterval > c.Interval {
				c.Interval = c.LearningInterval
			} else {
				c.LearningInterval = c.Interval
			}
			// Then set the LearningInterval to 0, as it is no longer needed.
			c.LearningInterval = 0
		}
	} else if c.LearningStage == Learned { // Learned stage
		c.Interval *= 2

		// If the Interval is more than 365 days (8760 hours), then the card has graduated to the burned stage and will no longer be reviewed.
		if c.Interval >= 8760 {
			c.LearningStage = 4
		}
	} else if c.LearningStage == UpNext { // Up next stage
		// If the card is in the up next stage, then it is being reviewed for the first time.
		// Set the LearningStage to 2, and set the LearningInterval to 3 hours.
		c.LearningStage = Learning
		c.LearningInterval = 3
	}

	c.IncrementReviewCount()
	c.IncrementCorrectAnswerCount()
	c.SetNextReviewDate()
}

func (c *Card) IncorrectAnswer() {
	// Check the next review date is in the past, otherwise this is a mistaken endpoint hit.
	t, err := time.Parse(time.RFC3339, c.NextReviewDate)
	if err != nil {
		log.Printf("Card %d has an invalid NextReviewDate: %s", c.ID, c.NextReviewDate)
		panic(err)
	}
	if time.Now().Before(t) {
		log.Printf("Card %d was reviewed too early. Next review date is %s", c.ID, c.NextReviewDate)
		return
	}

	c.ProcessIncorrectAnswer()
}

func (c *Card) ProcessIncorrectAnswer() {
	if c.LearningStage == Learning { // Learning stage
		// Only affect the LearningInterval.
		// The Interval is not affected, to preserve progress.
		c.LearningInterval /= 2

		// LearniingInterval cannot be less than 3 hours.
		if c.LearningInterval < 3 {
			c.LearningInterval = 3
		}
		c.IncrementReviewCount()
		c.SetNextFailedReviewDate()
	} else if c.LearningStage == Learned { // Learned stage
		c.Interval /= 2

		// Card gets downgraded to the learning stage
		c.LearningStage = Learning
		c.LearningInterval = 3 // Initial LearningInterval is 3 hours

		c.IncrementReviewCount()
		c.SetNextFailedReviewDate()
	} else if c.LearningStage == UpNext { // Up next stage
		// If the card is in the up next stage, then it is being reviewed for the first time.
		// If the answer is incorrect, then the card is rescheduled for the NextReviewDate + 10 minutes.
		// Since the card's NextReviewDate is initialised to 1970-01-01, this will result in the card being reviewed immediately after the current queue of up next cards are done.
		nrd, err := time.Parse(time.RFC3339, c.NextReviewDate)
		if err != nil {
			log.Printf("Card %d has an invalid NextReviewDate: %s", c.ID, c.NextReviewDate)
			panic(err)
		}
		c.NextReviewDate = nrd.Add(10 * time.Minute).Format(time.RFC3339)
	}
}

func (c *Card) SetNextReviewDate() {
	// Set the NextReviewDate to the current date + the Interval rounded to the hour.
	// If the card is in the learning stage, use the LearningInterval instead.
	// Burned cards will not be reviewed.
	if c.LearningStage == 2 {
		c.NextReviewDate = time.Now().Add(time.Duration(c.LearningInterval) * time.Hour).Round(time.Hour).Format(time.RFC3339)
	} else if c.LearningStage == 4 {
		c.NextReviewDate = ""
	} else {
		c.NextReviewDate = time.Now().Add(time.Duration(c.Interval) * time.Hour).Round(time.Hour).Format(time.RFC3339)
	}
}

func (c *Card) SetNextFailedReviewDate() {
	// Set the NextReviewDate to the current date + 10 minutes.
	// User is forced to keep reviewing the card until they get it right.
	c.NextReviewDate = time.Now().Add(time.Duration(10) * time.Minute).Round(time.Minute).Format(time.RFC3339)
}