package cards

import (
	"testing"
	"time"
)

func LearningStagesCardSlice() []*Card {
	c1 := CreateCard(1, 9600, 0, "") // Burned
	c2 := CreateCard(2, 48, 0, "2020-01-01T00:00:00Z") // Learned
	c3 := CreateCard(3, 0, 3, "2020-01-01T00:00:00Z") // Learning
	c4 := CreateCard(4, 0, 0, "2020-01-01T00:00:00Z") // Up Next
	c5 := CreateCard(5, 0, 0, "") // Available
	c5.ComponentSubjectIDs = []int{1, 2}
	c6 := CreateCard(6, 0, 0, "") // Unavailable
	c6.ComponentSubjectIDs = []int{5}

	s := []*Card{c1, c2, c3, c4, c5, c6}
	cd := CreateCardDataFromSlice(s)
	cd.UpdateCardData()

	return s
}

func TestFilterUnavailable(t *testing.T) {
	cs := LearningStagesCardSlice()
	codl := getOverviewLearningStage(cs, Unavailable)
	if len(codl.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl.Cards))
	}
	if codl.Cards[0].ID != 6 {
		t.Errorf("Expected card with ID 6, got %d", codl.Cards[0].ID)
	}
	if codl.Cards[0].LearningStage != Unavailable {
		t.Errorf("Expected learning stage to be Unavailable, got %d", codl.Cards[0].LearningStage)
	}
}

func TestFilterAvailable(t *testing.T) {
	cs := LearningStagesCardSlice()
	codl := getOverviewLearningStage(cs, Available)
	if len(codl.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl.Cards))
	}
	if codl.Cards[0].ID != 5 {
		t.Errorf("Expected card with ID 5, got %d", codl.Cards[0].ID)
	}
	if codl.Cards[0].LearningStage != Available {
		t.Errorf("Expected learning stage to be Available, got %d", codl.Cards[0].LearningStage)
	}
}

func TestFilterUpNext(t *testing.T) {
	cs := LearningStagesCardSlice()
	codl := getOverviewLearningStage(cs, UpNext)
	if len(codl.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl.Cards))
	}
	if codl.Cards[0].ID != 4 {
		t.Errorf("Expected card with ID 4, got %d", codl.Cards[0].ID)
	}
	if codl.Cards[0].LearningStage != UpNext {
		t.Errorf("Expected learning stage to be UpNext, got %d", codl.Cards[0].LearningStage)
	}
}

func TestFilterLearning(t *testing.T) {
	cs := LearningStagesCardSlice()
	codl := getOverviewLearningStage(cs, Learning)
	if len(codl.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl.Cards))
	}
	if codl.Cards[0].ID != 3 {
		t.Errorf("Expected card with ID 3, got %d", codl.Cards[0].ID)
	}
	if codl.Cards[0].LearningStage != Learning {
		t.Errorf("Expected learning stage to be Learning, got %d", codl.Cards[0].LearningStage)
	}
}

func TestFilterLearned(t *testing.T) {
	cs := LearningStagesCardSlice()
	codl := getOverviewLearningStage(cs, Learned)
	if len(codl.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl.Cards))
	}
	if codl.Cards[0].ID != 2 {
		t.Errorf("Expected card with ID 2, got %d", codl.Cards[0].ID)
	}
	if codl.Cards[0].LearningStage != Learned {
		t.Errorf("Expected learning stage to be Learned, got %d", codl.Cards[0].LearningStage)
	}
}

func TestFilterBurned(t *testing.T) {
	cs := LearningStagesCardSlice()
	codl := getOverviewLearningStage(cs, Burned)
	if len(codl.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl.Cards))
	}
	if codl.Cards[0].ID != 1 {
		t.Errorf("Expected card with ID 1, got %d", codl.Cards[0].ID)
	}
	if codl.Cards[0].LearningStage != Burned {
		t.Errorf("Expected learning stage to be Burned, got %d", codl.Cards[0].LearningStage)
	}
}

func LevelCardsSlice() []*Card {
	c1 := CreateCard(1, 9600, 0, "") // Burned
	c2 := CreateCard(2, 48, 0, "2020-01-01T00:00:00Z") // Learned
	c3 := CreateCard(3, 0, 3, "2020-01-01T00:00:00Z") // Learning
	c4 := CreateCard(4, 0, 0, "2020-01-01T00:00:00Z") // Up Next
	c5 := CreateCard(5, 0, 0, "") // Available
	c5.ComponentSubjectIDs = []int{1, 2}
	c6 := CreateCard(6, 0, 0, "") // Unavailable
	c6.ComponentSubjectIDs = []int{5}

	c1.Level = 1
	c2.Level = 1
	c3.Level = 2
	c4.Level = 3
	c5.Level = 4
	c6.Level = 3
	s := []*Card{c1, c2, c3, c4, c5, c6}
	cd := CreateCardDataFromSlice(s)
	cd.UpdateCardData()

	return s
}

func TestFilterLevel(t *testing.T) {
	cs := LevelCardsSlice()
	codl := filterCardsByLevel(cs, 1)
	if len(codl) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(codl))
	}
	codl = filterCardsByLevel(cs, 2)
	if len(codl) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl))
	}
	codl = filterCardsByLevel(cs, 3)
	if len(codl) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(codl))
	}
	codl = filterCardsByLevel(cs, 4)
	if len(codl) != 1 {
		t.Errorf("Expected 1 card, got %d", len(codl))
	}
}

func DueCardsSlice() []*Card {
	c1 := CreateCard(1, 9600, 0, "") // Burned
	c2 := CreateCard(2, 48, 0, "2020-01-01T00:00:00Z") // In the past
	c3 := CreateCard(3, 0, 3, time.Now().Add(-10*time.Minute).Format(time.RFC3339)) // 10 minutes ago
	c4 := CreateCard(4, 0, 0, time.Now().Add(2*time.Hour).Format(time.RFC3339)) // 2 hours from now 
	c5 := CreateCard(5, 0, 0, time.Now().Add(24*3*time.Hour).Format(time.RFC3339)) // 3 days from now
	c5.ComponentSubjectIDs = []int{1, 2}
	c6 := CreateCard(6, 0, 0, time.Now().Add(24*15*time.Hour).Format(time.RFC3339)) // 15 days from now
	c6.ComponentSubjectIDs = []int{5}

	s := []*Card{c1, c2, c3, c4, c5, c6}
	cd := CreateCardDataFromSlice(s)
	cd.UpdateCardData()

	return s
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func checkCardListLenAndIds(t *testing.T, length int, codl []*Card, expectedIds []int) {
	if len(codl) != length {
		t.Errorf("Expected %d cards, got %d", length, len(codl))
	}
	for _, c := range codl {
		if !containsInt(expectedIds, c.ID) {
			t.Errorf("Unexpected card with ID %d", c.ID)
		}
	}
}

func TestFilterDue(t *testing.T) {
	cs := DueCardsSlice()

	codl := filterCardsByDueBefore(cs, time.Now())
	checkCardListLenAndIds(t, 2, codl, []int{2, 3})

	codl = filterCardsByDueBefore(cs, time.Now().Add(24*time.Hour))
	checkCardListLenAndIds(t, 3, codl, []int{2, 3, 4})

	codl = filterCardsByDueBefore(cs, time.Now().Add(24*7*time.Hour))
	checkCardListLenAndIds(t, 4, codl, []int{2, 3, 4, 5})

	codl = filterCardsByDueBefore(cs, time.Now().Add(24*30*time.Hour))
	checkCardListLenAndIds(t, 5, codl, []int{2, 3, 4, 5, 6})

	codl = filterCardsByDueBetween(cs, time.Now(), time.Now().Add(24*time.Hour))
	checkCardListLenAndIds(t, 1, codl, []int{4})

	codl = filterCardsByDueBetween(cs, time.Now().Add(24*time.Hour), time.Now().Add(24*7*time.Hour))
	checkCardListLenAndIds(t, 1, codl, []int{5})

	codl = filterCardsByDueBetween(cs, time.Now().Add(24*7*time.Hour), time.Now().Add(24*30*time.Hour))
	checkCardListLenAndIds(t, 1, codl, []int{6})
}

func TypesCardSlice() []*Card {
	c1 := Card{ID: 1, Object: "radical"}
	c2 := Card{ID: 2, Object: "kanji"}
	c3 := Card{ID: 3, Object: "vocabulary"}
	c4 := Card{ID: 4, Object: "radical"}
	c5 := Card{ID: 5, Object: "kanji"}
	c6 := Card{ID: 6, Object: "vocabulary"}

	s := []*Card{&c1, &c2, &c3, &c4, &c5, &c6}
	return s
}

func TestFilterType(t *testing.T) {
	cd := TypesCardSlice()
	codl := filterCardsByType(cd, "radical")
	checkCardListLenAndIds(t, 2, codl, []int{1, 4})

	codl = filterCardsByType(cd, "kanji")
	checkCardListLenAndIds(t, 2, codl, []int{2, 5})

	codl = filterCardsByType(cd, "vocabulary")
	checkCardListLenAndIds(t, 2, codl, []int{3, 6})
}