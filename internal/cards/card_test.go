package cards

import (
	"testing"
)

func CreateCardDataFromSlice(cards []*Card) *CardData {
	cardMap := make(map[int]*Card)
	for _, c := range cards {
		cardMap[c.ID] = c
	}
	return &CardData{
		Cards: cardMap,
	}
}

func CreateCard(id int, interval int, learningInterval int, NextReviewDate string) *Card {
	return &Card{
		ID: id,
		Interval: interval,
		LearningInterval: learningInterval,
		NextReviewDate: NextReviewDate,
	}
}

func TestCardBurned(t *testing.T) {
	c := CreateCard(1, 9600, 0, "")
	cd := CreateCardDataFromSlice([]*Card{c})
	cd.UpdateCardData()
	if c.LearningStage != Burned {
		t.Errorf("Card stage. Expected: %d, Got: %d", Burned, c.LearningStage)
	}
	if c.GetLearningStageString() != "Burned" {
		t.Errorf("Card stage string. Expected: %s, Got: %s", "Burned", c.LearningStageString)
	}
}

func TestCardLearned(t *testing.T) {
	c := CreateCard(1, 48, 0, "2020-01-01T00:00:00Z")
	cd := CreateCardDataFromSlice([]*Card{c})
	cd.UpdateCardData()
	if c.LearningStage != Learned {
		t.Errorf("Card stage. Expected: %d, Got: %d", Learned, c.LearningStage)
	}
	if c.GetLearningStageString() != "Learned" {
		t.Errorf("Card stage string. Expected: %s, Got: %s", "Learned", c.LearningStageString)
	}
}

func TestCardLearning(t *testing.T) {
	c := CreateCard(1, 0, 3, "2020-01-01T00:00:00Z")
	cd := CreateCardDataFromSlice([]*Card{c})
	cd.UpdateCardData()
	if c.LearningStage != Learning {
		t.Errorf("Card stage. Expected: %d, Got: %d", Learning, c.LearningStage)
	}
	if c.GetLearningStageString() != "Learning" {
		t.Errorf("Card stage string. Expected: %s, Got: %s", "Learning", c.LearningStageString)
	}
}

func TestCardUpNext(t *testing.T) {
	c := CreateCard(1, 0, 0, "2020-01-01T00:00:00Z")
	cd := CreateCardDataFromSlice([]*Card{c})
	cd.UpdateCardData()
	if c.LearningStage != UpNext {
		t.Errorf("Card stage. Expected: %d, Got: %d", UpNext, c.LearningStage)
	}
	if c.GetLearningStageString() != "Up Next" {
		t.Errorf("Card stage string. Expected: %s, Got: %s", "Up Next", c.LearningStageString)
	}
}

func TestCardAvailable(t *testing.T) {
	c1 := CreateCard(1, 9600, 0, "2020-01-01T00:00:00Z") // Burned card
	c2 := CreateCard(2, 48, 0, "2020-01-01T00:00:00Z") // Learned card
	c := CreateCard(3, 0, 0, "")
	c.ComponentSubjectIDs = []int{1, 2}
	cd := CreateCardDataFromSlice([]*Card{c1, c2, c})
	cd.UpdateCardData()
	if c.LearningStage != Available {
		t.Errorf("Card stage. Expected: %d, Got: %d", Available, c.LearningStage)
	}
	if c.GetLearningStageString() != "Available" {
		t.Errorf("Card stage string. Expected: %s, Got: %s", "Available", c.LearningStageString)
	}
}

func TestCardUnavailable(t *testing.T) {
	c1 := CreateCard(1, 48, 0, "2020-01-01T00:00:00Z") // Learned card
	c2 := CreateCard(2, 3, 0, "2020-01-01T00:00:00Z") // Learning card
	c := CreateCard(3, 0, 0, "")
	c.ComponentSubjectIDs = []int{1, 2}
	cd := CreateCardDataFromSlice([]*Card{c1, c2, c})
	cd.UpdateCardData()
	if c.LearningStage != Unavailable {
		t.Errorf("Card stage. Expected: %d, Got: %d", Unavailable, c.LearningStage)
	}
	if c.GetLearningStageString() != "Unavailable" {
		t.Errorf("Card stage string. Expected: %s, Got: %s", "Unavailable", c.LearningStageString)
	}
}

func TestCardSetToUpNext(t *testing.T) {
	c := CreateCard(1, 0, 0, "")
	cd := CreateCardDataFromSlice([]*Card{c})
	cd.UpdateCardData()
	if c.LearningStage != Available {
		t.Errorf("Card stage. Expected: %d, Got: %d", Available, c.LearningStage)
	}

	c.SetToUpNext()
	cd.UpdateCardData()

	if c.LearningStage != UpNext {
		t.Errorf("Card stage. Expected: %d, Got: %d", UpNext, c.LearningStage)
	}
}

func TestCardSetToUpNextFailure(t *testing.T) {
	c1 := CreateCard(1, 0, 0, "")
	c2 := CreateCard(2, 0, 0, "")
	c2.ComponentSubjectIDs = []int{1}
	cd := CreateCardDataFromSlice([]*Card{c1, c2})
	cd.UpdateCardData()

	if c2.LearningStage != Unavailable {
		t.Errorf("Card stage. Expected: %d, Got: %d", Unavailable, c2.LearningStage)
	}

	c2.SetToUpNext()
	cd.UpdateCardData()

	if c2.LearningStage != Unavailable {
		t.Errorf("Card stage. Expected: %d, Got: %d", Unavailable, c2.LearningStage)
	}
}