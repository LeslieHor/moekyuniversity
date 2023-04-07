package cards

import (
	"testing"
	"time"
)

func TestCorrectLearningToLearning(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
		LearningInterval: 3, // 3 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.CorrectAnswer()
	if c.Interval != 0 {
		t.Errorf("Incorrect interval. Expected 0, got %d", c.Interval)
	}
	if c.LearningStage != 2 {
		t.Errorf("Incorrect learning stage. Expected 2, got %d", c.LearningStage)
	}
	if c.LearningInterval != 6 {
		t.Errorf("Incorrect learning interval. Expected 6, got %d", c.LearningInterval)
	}
	// Check the next review date is 6 hours in the future, rounded to the nearest hour.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(time.Hour * 6).Round(time.Hour).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestCorrectLearningToLearned(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
	    LearningInterval: 12, // 12 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.CorrectAnswer()
	if c.Interval != 24 {
		t.Errorf("Incorrect interval. Expected 24, got %d", c.Interval)
	}
	if c.LearningStage != 3 {
		t.Errorf("Incorrect learning stage. Expected 3, got %d", c.LearningStage)
	}
	if c.LearningInterval != 0 {
		t.Errorf("Incorrect learning interval. Expected 0, got %d", c.LearningInterval)
	}
	// Check the next review date is 24 hours in the future, rounded to the nearest hour.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(time.Hour * 24).Round(time.Hour).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestCorrectLearnedToLearned(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 48, // 2 days
		LearningInterval: 0,
		LearningStage: 3,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.CorrectAnswer()
	if c.Interval != 96 {
		t.Errorf("Incorrect interval. Expected 96, got %d", c.Interval)
	}
	if c.LearningStage != 3 {
		t.Errorf("Incorrect learning stage. Expected 3, got %d", c.LearningStage)
	}
	if c.LearningInterval != 0 {
		t.Errorf("Incorrect learning interval. Expected 0, got %d", c.LearningInterval)
	}
	// Check the next review date is 4 days in the future, rounded to the nearest hour.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(4 * 24 * time.Hour).Round(time.Hour).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestCorrectLearnedToBurned(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 24 * 200, // 200 days
		LearningInterval: 0,
		LearningStage: 3,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.CorrectAnswer()
	if c.Interval != 24 * 400 {
		t.Errorf("Incorrect interval. Expected 24 * 400, got %d", c.Interval)
	}
	if c.LearningStage != 4 {
		t.Errorf("Incorrect learning stage. Expected 4, got %d", c.LearningStage)
	}
	if c.LearningInterval != 0 {
		t.Errorf("Incorrect learning interval. Expected 0, got %d", c.LearningInterval)
	}
	// Burned cards are never reviewed again. Expect the next review date to be empty.
	if c.NextReviewDate != "" {
		t.Errorf("Incorrect next review date. Expected empty string, got %s", c.NextReviewDate)
	}
}

func TestIncorrectLearningToLearningMinimum(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
		LearningInterval: 3, // 3 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.IncorrectAnswer()
	if c.Interval != 0 {
		t.Errorf("Incorrect interval. Expected 0, got %d", c.Interval)
	}
	if c.LearningStage != 2 {
		t.Errorf("Incorrect learning stage. Expected 2, got %d", c.LearningStage)
	}
	// Learning interval should be reset to 3 hours minimum.
	if c.LearningInterval != 3 {
		t.Errorf("Incorrect learning interval. Expected 3, got %d", c.LearningInterval)
	}
	// Next review should be 10 minutes in the future, rounded to the nearest minute.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(10 * time.Minute).Round(time.Minute).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestIncorrectLearningToLearning(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
		LearningInterval: 12, // 12 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.IncorrectAnswer()
	if c.Interval != 0 {
		t.Errorf("Incorrect interval. Expected 0, got %d", c.Interval)
	}
	if c.LearningStage != 2 {
		t.Errorf("Incorrect learning stage. Expected 2, got %d", c.LearningStage)
	}
	// Learning interval should be halved.
	if c.LearningInterval != 6 {
		t.Errorf("Incorrect learning interval. Expected 6, got %d", c.LearningInterval)
	}

	// Next review should be 10 minutes in the future, rounded to the nearest minute.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(10 * time.Minute).Round(time.Minute).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestIncorrectLearnedToLearning(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 240, // 10 days
		LearningInterval: 0,
		LearningStage: 3,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.IncorrectAnswer()
	if c.Interval != 120 {
		t.Errorf("Incorrect interval. Expected 120, got %d", c.Interval)
	}
	if c.LearningStage != 2 {
		t.Errorf("Incorrect learning stage. Expected 2, got %d", c.LearningStage)
	}
	// Learning interval should be reset to 3 hours.
	if c.LearningInterval != 3 {
		t.Errorf("Incorrect learning interval. Expected 3, got %d", c.LearningInterval)
	}

	// Next review should be 10 minutes in the future, rounded to the nearest minute.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(10 * time.Minute).Round(time.Minute).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}

	// Getting the review wrong again should not halve the interval again
	// Set the next review date to the past
	c.NextReviewDate = "2020-01-01T00:00:00Z"
	c.IncorrectAnswer()
	if c.Interval != 120 {
		t.Errorf("Incorrect interval. Expected 120, got %d", c.Interval)
	}
	if c.LearningStage != 2 {
		t.Errorf("Incorrect learning stage. Expected 2, got %d", c.LearningStage)
	}
	// Learning interval should be reset to 3 hours.
	if c.LearningInterval != 3 {
		t.Errorf("Incorrect learning interval. Expected 3, got %d", c.LearningInterval)
	}

	// Next review should be 10 minutes in the future, rounded to the nearest minute.
	// Formatted as RFC3339.
	ExpectedNextReviewDate = time.Now().Add(10 * time.Minute).Round(time.Minute).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestCorrectRelearningToRelearning(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 120, // 5 days
		LearningInterval: 3, // 3 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.CorrectAnswer()

	// Interval should not change.
	if c.Interval != 120 {
		t.Errorf("Incorrect interval. Expected 120, got %d", c.Interval)
	}

	// Learning stage should not change.
	if c.LearningStage != 2 {
		t.Errorf("Incorrect learning stage. Expected 2, got %d", c.LearningStage)
	}

	// Learning interval should be doubled.
	if c.LearningInterval != 6 {
		t.Errorf("Incorrect learning interval. Expected 6, got %d", c.LearningInterval)
	}

	// Check the next review date is 6 hours in the future, rounded to the nearest hour.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(6 * time.Hour).Round(time.Hour).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestCorrectRelearningToLearned(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 120, // 5 days
		LearningInterval: 16, // 16 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
	}

	c.CorrectAnswer()

	// Interval should not be doubled.
	if c.Interval != 120 {
		t.Errorf("Incorrect interval. Expected 120, got %d", c.Interval)
	}

	// Learning stage should be set to 3.
	if c.LearningStage != 3 {
		t.Errorf("Incorrect learning stage. Expected 3, got %d", c.LearningStage)
	}

	// Learning interval should be reset to 0.
	if c.LearningInterval != 0 {
		t.Errorf("Incorrect learning interval. Expected 0, got %d", c.LearningInterval)
	}

	// Check the next review date is 5 days in the future, rounded to the nearest hour.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(5 * 24 * time.Hour).Round(time.Hour).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}
}

func TestNextSrsCardUpNext(t *testing.T) {
	c1 := &Card{ID: 1, LearningStage: UpNext, NextReviewDate: "2020-01-01T00:00:00Z"}
	cd := CreateCardDataFromSlice([]*Card{c1})

	// The card should be returned as the next SRS card.
	srsCard := cd.GetNextSrsCard()
	if srsCard.Card.ID != c1.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c1.ID, srsCard.Card.ID)
	}

	// Set the next review date to the future.
	c1.NextReviewDate = time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	// The card should still be returned as the next SRS card.
	srsCard = cd.GetNextSrsCard()
	if srsCard.Card.ID != c1.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c1.ID, srsCard.Card.ID)
	}
}

func TestNextSrsCardLearning(t *testing.T) {
	c1 := &Card{ID: 1, LearningStage: Learning, NextReviewDate: "2020-01-01T00:00:00Z"}
	cd := CreateCardDataFromSlice([]*Card{c1})

	// The card should be returned as the next SRS card.
	srsCard := cd.GetNextSrsCard()
	if srsCard.Card.ID != c1.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c1.ID, srsCard.Card.ID)
	}

	// Set the next review date to the future.
	c1.NextReviewDate = time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	// The card should not be returned as the next SRS card.
	srsCard = cd.GetNextSrsCard()
	if srsCard.Card != nil {
		t.Errorf("Incorrect card returned. Expected nil, got a card")
	}
}

func TestNextSrsCardLearned(t *testing.T) {
	c1 := &Card{ID: 1, LearningStage: Learned, NextReviewDate: "2020-01-01T00:00:00Z"}
	cd := CreateCardDataFromSlice([]*Card{c1})

	// The card should be returned as the next SRS card.
	srsCard := cd.GetNextSrsCard()
	if srsCard.Card.ID != c1.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c1.ID, srsCard.Card.ID)
	}

	// Set the next review date to the future.
	c1.NextReviewDate = time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	// The card should not be returned as the next SRS card.
	srsCard = cd.GetNextSrsCard()
	if srsCard.Card != nil {
		t.Errorf("Incorrect card returned. Expected nil, got a card")
	}
}

func TestNextSrsCardUpNextPriority(t *testing.T) {
	c1 := &Card{ID: 1, LearningStage: UpNext, NextReviewDate: "3020-01-01T00:00:00Z"}
	c2 := &Card{ID: 2, LearningStage: Learning, NextReviewDate: "2020-01-01T00:00:00Z"}
	cd := CreateCardDataFromSlice([]*Card{c1, c2})

	// The card with the earlier ID should be returned as the next SRS card.
	srsCard := cd.GetNextSrsCard()
	if srsCard.Card.ID != c1.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c1.ID, srsCard.Card.ID)
	}

	// Set the next review date to the future.
	c1.NextReviewDate = time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	// The UpNext card should still be returned as the next SRS card.
	srsCard = cd.GetNextSrsCard()
	if srsCard.Card.ID != c1.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c1.ID, srsCard.Card.ID)
	}

	// Set the status to Learning.
	c1.LearningStage = Learning
	
	// The c2 Learning card should now be returned as the next SRS card, as it has the earliest next review date.
	srsCard = cd.GetNextSrsCard()
	if srsCard.Card.ID != c2.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c2.ID, srsCard.Card.ID)
	}
}