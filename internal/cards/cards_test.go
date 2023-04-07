package cards

import (
	"testing"
)

func checkSearchResults(t *testing.T, results []*Card, expectedIds []int) {
	if len(results) != len(expectedIds) {
		t.Errorf("Expected %d results, got %d", len(expectedIds), len(results))
	}
	for _, r := range results {
		if !containsInt(expectedIds, r.ID) {
			t.Errorf("Unexpected card with ID %d", r.ID)
		}
	}
}

func createCardWithMeaning(id int, characters string, meaning string) *Card {
	return &Card{ID: id,
				 Characters: characters,
				 Meanings: []struct {
					Meaning string `json:"meaning"`
					Primary bool `json:"primary"`
					AcceptedAnswer bool `json:"accepted_answer"`
				}{
					{meaning, true, true},
				}}
}

func TestSearch(t *testing.T) {
	c1 := createCardWithMeaning(1, "一", "One")
	c2 := createCardWithMeaning(2, "一つ", "One thing")
	c3 := createCardWithMeaning(3, "二", "Two")
	c4 := createCardWithMeaning(4, "二つ", "Two things")
	c5 := createCardWithMeaning(5, "三", "Three")
	c6 := createCardWithMeaning(6, "三つ", "Three things")
	cd := CreateCardDataFromSlice([]*Card{c1, c2, c3, c4, c5, c6})

	// Search for "一"
	results := cd.Search("一")
	checkSearchResults(t, results, []int{1, 2})

	// Search for "一つ"
	results = cd.Search("一つ")
	checkSearchResults(t, results, []int{2})

	// Search for "二"
	results = cd.Search("二")
	checkSearchResults(t, results, []int{3, 4})

	// Search for "Thing"
	results = cd.Search("Thing")
	checkSearchResults(t, results, []int{2, 4, 6})
}