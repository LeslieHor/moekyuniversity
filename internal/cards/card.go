package cards

import (

)

type LearningStage int

const (
	Unavailable LearningStage = iota
	Available
	Learning
	Learned
	Burned
	UpNext
)

var LearningStages = []LearningStage{
	Unavailable,
	Available,
	UpNext,
	Learning,
	Learned,
	Burned,
}

type Card struct {
	ID    int    `json:"id"`
	Object string `json:"object"`
	Level int    `json:"level"`
	DocumentURL string `json:"document_url"`
	Characters string `json:"characters"`
	Meanings []struct {
		Meaning string `json:"meaning"`
		Primary bool `json:"primary"`
		AcceptedAnswer bool `json:"accepted_answer"`
	} `json:"meanings"`
	Readings []struct {
		Reading string `json:"reading"`
		Type string `json:"type"`
		Primary bool `json:"primary"`
		AcceptedAnswer bool `json:"accepted_answer"`
	} `json:"readings"`
	PartsOfSpeech []string `json:"parts_of_speech"`
	ComponentSubjectIDs []int `json:"component_subject_ids"`
	AmalgamationSubjectIDs []int `json:"amalgamation_subject_ids"`
	ReadingMnemonic string `json:"reading_mnemonic"`
	MeaningMnemonic string `json:"meaning_mnemonic"`

	Interval int `json:"interval"` // Hours until next review
	LearningInterval int `json:"learning_interval"` // Hours until next review when in learning stage
	NextReviewDate string `json:"next_review_date"` // RFC3339 date string

	LearningStage LearningStage `json:"learning_stage"` // 0 = Unavailable, 1 = Available, 2 = Learning, 3 = Learned, 4 = Burned

	// Below is for html output
	LearningStageString string `json:"learning_stage_string"` // Unavailable, Available, Learning, Learned, Burned
}

func (c *Card) UpdateLearningStage(cd *CardData) {
	if c.Interval > 8760 {
		c.LearningStage = Burned
		c.LearningStageString = LearningStageToString(c.LearningStage)
		return
	}

	if c.Interval > 24 &&
		c.LearningInterval == 0 &&
		c.NextReviewDate != "" {
		c.LearningStage = Learned
		c.LearningStageString = LearningStageToString(c.LearningStage)
		return
	}

	if c.LearningInterval > 0 &&
		c.NextReviewDate != "" {
		c.LearningStage = Learning
		c.LearningStageString = LearningStageToString(c.LearningStage)
		return
	}

	if c.Interval == 0 &&
		c.LearningInterval == 0 &&
		c.NextReviewDate != "" {
		c.LearningStage = UpNext
		c.LearningStageString = LearningStageToString(c.LearningStage)
		return
	}

	dt := c.GetDataTree(cd)
	if dt.IsAllChildrenLearned() {
		c.LearningStage = Available
		c.LearningStageString = LearningStageToString(c.LearningStage)
		return
	}

	c.LearningStage = Unavailable
	c.LearningStageString = LearningStageToString(c.LearningStage)
}

func LearningStageToString(ls LearningStage) string {
	switch ls {
	case Unavailable:
		return "Unavailable"
	case Available:
		return "Available"
	case UpNext:
		return "Up Next"
	case Learning:
		return "Learning"
	case Learned:
		return "Learned"
	case Burned:
		return "Burned"
	default:
		return "Unknown"
	}
}

func (c *Card) GetLearningStageString() string {
	return LearningStageToString(c.LearningStage)
}

func (c *Card) GetPartsOfSpeech() []string {
	return c.PartsOfSpeech
}