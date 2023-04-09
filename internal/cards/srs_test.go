package cards

import (
	"testing"
	"time"
)

func TestCorrectUpNextToLearning(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
		LearningInterval: 0,
		LearningStage: UpNext,
		NextReviewDate: "1970-01-01T00:00:00Z", // 0 unix time
	}

	c.CorrectAnswer()
	if c.Interval != 0 {
		t.Errorf("Incorrect interval. Expected 0, got %d", c.Interval)
	}
	if c.LearningStage != Learning {
		t.Errorf("Incorrect learning stage. Expected 1, got %d", c.LearningStage)
	}
	if c.LearningInterval != 3 {
		t.Errorf("Incorrect learning interval. Expected 3, got %d", c.LearningInterval)
	}
	// Check the next review date is 3 hours in the future, rounded to the nearest hour.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Now().Add(time.Hour * 3).Round(time.Hour).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}

	if c.TotalTimesReviewed != 1 {
		t.Errorf("Incorrect total times reviewed. Expected 1, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 1 {
		t.Errorf("Incorrect total times correct. Expected 1, got %d", c.TotalTimesCorrect)
	}
}

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

	if c.TotalTimesReviewed != 1 {
		t.Errorf("Incorrect total times reviewed. Expected 1, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 1 {
		t.Errorf("Incorrect total times correct. Expected 1, got %d", c.TotalTimesCorrect)
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

	if c.TotalTimesReviewed != 1 {
		t.Errorf("Incorrect total times reviewed. Expected 1, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 1 {
		t.Errorf("Incorrect total times correct. Expected 1, got %d", c.TotalTimesCorrect)
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

	if c.TotalTimesReviewed != 1 {
		t.Errorf("Incorrect total times reviewed. Expected 1, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 1 {
		t.Errorf("Incorrect total times correct. Expected 1, got %d", c.TotalTimesCorrect)
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

	if c.TotalTimesReviewed != 1 {
		t.Errorf("Incorrect total times reviewed. Expected 1, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 1 {
		t.Errorf("Incorrect total times correct. Expected 1, got %d", c.TotalTimesCorrect)
	}
}

func TestIncorrectUpNextToUpNext(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
		LearningInterval: 0,
		LearningStage: UpNext,
		NextReviewDate: "1970-01-01T00:00:00Z", // Any date in the past
	}

	c.IncorrectAnswer()
	if c.Interval != 0 {
		t.Errorf("Incorrect interval. Expected 0, got %d", c.Interval)
	}
	if c.LearningStage != UpNext {
		t.Errorf("Incorrect learning stage. Expected 0, got %d", c.LearningStage)
	}
	if c.LearningInterval != 0 {
		t.Errorf("Incorrect learning interval. Expected 0, got %d", c.LearningInterval)
	}
	// Check the next review date is 10 minutes after 1970-01-01T00:00:00Z.
	// Formatted as RFC3339.
	ExpectedNextReviewDate := time.Date(1970, 1, 1, 0, 10, 0, 0, time.UTC).Format(time.RFC3339)
	if c.NextReviewDate != ExpectedNextReviewDate {
		t.Errorf("Incorrect next review date. Expected %s, got %s", ExpectedNextReviewDate, c.NextReviewDate)
	}

	// We don't count incorrect answers for UpNext cards.
	if c.TotalTimesReviewed != 0 {
		t.Errorf("Incorrect total times reviewed. Expected 0, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 0 {
		t.Errorf("Incorrect total times correct. Expected 0, got %d", c.TotalTimesCorrect)
	}
}

func TestIncorrectLearningToLearningMinimum(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
		LearningInterval: 3, // 3 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
		TotalTimesReviewed: 5,
		TotalTimesCorrect: 3,
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

	if c.TotalTimesReviewed != 6 {
		t.Errorf("Incorrect total times reviewed. Expected 6, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 3 {
		t.Errorf("Incorrect total times correct. Expected 3, got %d", c.TotalTimesCorrect)
	}
}

func TestIncorrectLearningToLearning(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 0,
		LearningInterval: 12, // 12 hours
		LearningStage: 2,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
		TotalTimesReviewed: 5,
		TotalTimesCorrect: 3,
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

	if c.TotalTimesReviewed != 6 {
		t.Errorf("Incorrect total times reviewed. Expected 6, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 3 {
		t.Errorf("Incorrect total times correct. Expected 3, got %d", c.TotalTimesCorrect)
	}
}

func TestIncorrectLearnedToLearning(t *testing.T) {
	c := Card{
		ID: 1,
		Interval: 240, // 10 days
		LearningInterval: 0,
		LearningStage: 3,
		NextReviewDate: "2020-01-01T00:00:00Z", // Any date in the past
		TotalTimesReviewed: 5,
		TotalTimesCorrect: 3,
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

	if c.TotalTimesReviewed != 6 {
		t.Errorf("Incorrect total times reviewed. Expected 6, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 3 {
		t.Errorf("Incorrect total times correct. Expected 3, got %d", c.TotalTimesCorrect)
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

	if c.TotalTimesReviewed != 7 {
		t.Errorf("Incorrect total times reviewed. Expected 7, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 3 {
		t.Errorf("Incorrect total times correct. Expected 3, got %d", c.TotalTimesCorrect)
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

	if c.TotalTimesReviewed != 1 {
		t.Errorf("Incorrect total times reviewed. Expected 1, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 1 {
		t.Errorf("Incorrect total times correct. Expected 1, got %d", c.TotalTimesCorrect)
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

	if c.TotalTimesReviewed != 1 {
		t.Errorf("Incorrect total times reviewed. Expected 1, got %d", c.TotalTimesReviewed)
	}
	if c.TotalTimesCorrect != 1 {
		t.Errorf("Incorrect total times correct. Expected 1, got %d", c.TotalTimesCorrect)
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

	// The card that is Learning should be returned as the next SRS card.
	srsCard := cd.GetNextSrsCard()
	if srsCard.Card.ID != c2.ID {
		t.Errorf("Incorrect card returned. Expected %d, got %d", c2.ID, srsCard.Card.ID)
	}
}

func TestNextSrsUpNext(t *testing.T) {
	c1 := &Card{ID: 1, LearningStage: UpNext, NextReviewDate: "1970-01-01T00:00:00Z"}
	c2 := &Card{ID: 2, LearningStage: UpNext, NextReviewDate: "1970-01-01T00:00:00Z"}
	c3 := &Card{ID: 3, LearningStage: UpNext, NextReviewDate: "1970-01-01T00:00:00Z"}
	c4 := &Card{ID: 4, LearningStage: UpNext, NextReviewDate: "1970-01-01T00:00:00Z"}
	c5 := &Card{ID: 5, LearningStage: UpNext, NextReviewDate: "1970-01-01T00:00:00Z"}
	c6 := &Card{ID: 6, LearningStage: UpNext, NextReviewDate: "1970-01-01T00:00:00Z"}
	
	cd := CreateCardDataFromSlice([]*Card{c1, c2, c3, c4, c5, c6})

	// UpNext cards should be limited to 5.
	srsCard := cd.GetNextSrsCard()
	if srsCard.DueCount != 5 {
		t.Errorf("Incorrect number of cards returned. Expected %d, got %d", 5, srsCard.DueCount)
	}
}

func TestSuspended(t *testing.T) {
	c1 := &Card{ID: 1, LearningStage: UpNext, NextReviewDate: "1970-01-01T00:00:00Z", Tags: []string{"suspended"}}
	cd := CreateCardDataFromSlice([]*Card{c1})

	// The card should not be returned as the next SRS card.
	srsCard := cd.GetNextSrsCard()
	if srsCard.Card != nil {
		t.Errorf("Incorrect card returned. Expected nil, got a card")
	}
}

func TestAutoUpNext(t *testing.T) {
	c1 := &Card{ID: 1, Interval: 500, LearningInterval: 0, NextReviewDate: time.Now().Add(24 * 30 * time.Hour).Format(time.RFC3339)}
	c2 := &Card{ID: 2, Interval: 0, LearningInterval: 16, NextReviewDate: time.Now().Add(-2 * time.Hour).Format(time.RFC3339)}
	c3 := &Card{ID: 3, Interval: 0, LearningInterval: 0, QueuedToLearn: true, NextReviewDate: "", ComponentSubjectIDs: []int{1, 2}}
	cd := CreateCardDataFromSlice([]*Card{c1, c2, c3})
	cd.UpdateCardData()

	// c1 should be Learned
	if c1.LearningStage != Learned {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Learned, c1.LearningStage)
	}
	// c2 should be Learning
	if c2.LearningStage != Learning {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Learning, c2.LearningStage)
	}
	// c3 should be QueuedToLearn
	if c3.LearningStage != QueuedToLearn {
		t.Errorf("Incorrect card stage. Expected %d, got %d", QueuedToLearn, c3.LearningStage)
	}

	// Set c2 to learned
	c2.CorrectAnswer()
	cd.UpdateCardData()

	// c2 should be Learned
	if c2.LearningStage != Learned {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Learned, c2.LearningStage)
	}

	// c3 should now be UpNext
	if c3.LearningStage != UpNext {
		t.Errorf("Incorrect card stage. Expected %d, got %d", UpNext, c3.LearningStage)
	}
}

func TestAutoUpNextMultipleDependencies(t *testing.T) {
	c1 := &Card{ID: 1, Interval: 0, LearningInterval: 16, NextReviewDate: time.Now().Add(-2 * time.Hour).Format(time.RFC3339)}
	c2 := &Card{ID: 2, Interval: 0, LearningInterval: 0, ComponentSubjectIDs: []int{1}}
	c3 := &Card{ID: 3, Interval: 0, LearningInterval: 0, ComponentSubjectIDs: []int{2}}
	cd := CreateCardDataFromSlice([]*Card{c1, c2, c3})
	cd.UpdateCardData()

	// c1 should be learning
	if c1.LearningStage != Learning {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Learning, c1.LearningStage)
	}
	// c2 should be unavailable
	if c2.LearningStage != Unavailable {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Unavailable, c2.LearningStage)
	}
	// c3 should be unavailable
	if c3.LearningStage != Unavailable {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Unavailable, c3.LearningStage)
	}
	
	// Set c3 to queued to learn
	c3.SetQueuedToLearn(cd)
	cd.UpdateCardData()

	// c2 and c3 should be queued to learn
	if !c2.QueuedToLearn {
		t.Errorf("Card should be queued to learn")
	}
	if c2.LearningStage != QueuedToLearn {
		t.Errorf("Incorrect card stage. Expected %d, got %d", QueuedToLearn, c2.LearningStage)
	}
	if !c3.QueuedToLearn {
		t.Errorf("Card should be queued to learn")
	}
	if c3.LearningStage != QueuedToLearn {
		t.Errorf("Incorrect card stage. Expected %d, got %d", QueuedToLearn, c3.LearningStage)
	}

	// Mark c1 as answered. It should now be learned.
	c1.CorrectAnswer()
	cd.UpdateCardData()

	// c1 should be learned
	if c1.LearningStage != Learned {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Learned, c1.LearningStage)
	}

	// c2 should now be upnext
	if c2.LearningStage != UpNext {
		t.Errorf("Incorrect card stage. Expected %d, got %d", UpNext, c2.LearningStage)
	}
	if c2.QueuedToLearn {
		t.Errorf("Card should not be queued to learn")
	}
	// c3 should remain queued to learn
	if c3.LearningStage != QueuedToLearn {
		t.Errorf("Incorrect card stage. Expected %d, got %d", QueuedToLearn, c3.LearningStage)
	}
	if !c3.QueuedToLearn {
		t.Errorf("Card should be queued to learn")
	}

	// Mark c2 as answered. It should now be learning
	c2.CorrectAnswer()
	cd.UpdateCardData()

	// c2 should be learning
	if c2.LearningStage != Learning {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Learning, c2.LearningStage)
	}

	// Fake c2 reviews
	c2.NextReviewDate = time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	c2.LearningInterval = 16
	c2.CorrectAnswer()
	cd.UpdateCardData()

	// c2 should be learned
	if c2.LearningStage != Learned {
		t.Errorf("Incorrect card stage. Expected %d, got %d", Learned, c2.LearningStage)
	}

	// c3 should now be upnext
	if c3.LearningStage != UpNext {
		t.Errorf("Incorrect card stage. Expected %d, got %d", UpNext, c3.LearningStage)
	}
	if c3.QueuedToLearn {
		t.Errorf("Card should not be queued to learn")
	}
}