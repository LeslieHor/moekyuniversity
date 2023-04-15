package cards

import (
	"time"
)

type LearningStage int

const (
	Unavailable LearningStage = iota
	Available
	Learning
	Learned
	Burned
	UpNext
	QueuedToLearn
)

var LearningStages = []LearningStage{
	Unavailable,
	Available,
	QueuedToLearn,
	UpNext,
	Learning,
	Learned,
	Burned,
}

type Meaning struct {
	Meaning        string `json:"meaning"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer bool   `json:"accepted_answer"`
}

type Reading struct {
	Reading        string `json:"reading"`
	Type           string `json:"type"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer bool   `json:"accepted_answer"`
}

type Card struct {
	ID                     int       `json:"id"`
	Object                 string    `json:"object"`
	Level                  int       `json:"level"`
	DocumentURL            string    `json:"document_url"`
	Characters             string    `json:"characters"`
	CharacterImage         string    `json:"character_image"`
	CharacterAlt           string    `json:"character_alt"`
	Meanings               []Meaning `json:"meanings"`
	Readings               []Reading `json:"readings"`
	PartsOfSpeech          []string  `json:"parts_of_speech"`
	ComponentSubjectIDs    []int     `json:"component_subject_ids"`
	AmalgamationSubjectIDs []int     `json:"amalgamation_subject_ids"`
	ReadingMnemonic        string    `json:"reading_mnemonic"`
	MeaningMnemonic        string    `json:"meaning_mnemonic"`

	Interval           int    `json:"interval"`          // Hours until next review
	LearningInterval   int    `json:"learning_interval"` // Hours until next review when in learning stage
	NextReviewDate     string `json:"next_review_date"`  // RFC3339 date string
	TotalTimesReviewed int    `json:"total_times_reviewed"`
	TotalTimesCorrect  int    `json:"total_times_correct"`
	QueuedToLearn      bool   `json:"queued_to_learn"`

	LearningStage LearningStage `json:"learning_stage"` // 0 = Unavailable, 1 = Available, 2 = Learning, 3 = Learned, 4 = Burned

	Tags []string `json:"tags"`

	// Below is for html output
	LearningStageString string `json:"learning_stage_string"` // Unavailable, Available, Learning, Learned, Burned
}

func (c *Card) UpdateLearningStage(cd *CardData) {
	if c.Interval > 8760 {
		c.LearningStage = Burned
		c.LearningStageString = LearningStageToString(c.LearningStage)
		return
	}

	if c.Interval >= 24 &&
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
	acl := dt.IsAllChildrenLearned()

	if c.QueuedToLearn && acl {
		// Card is queued to learn and all children are learned
		// So it should be on the UpNext queue
		c.LearningStage = UpNext
		c.NextReviewDate = time.Unix(0, 0).Format(time.RFC3339)
		c.LearningStageString = LearningStageToString(c.LearningStage)
		c.QueuedToLearn = false
		return

	} else if c.QueuedToLearn && !acl {
		// Card is queued to learn but not all children are learned
		// So it should stay as queued to learn
		c.LearningStage = QueuedToLearn
		c.LearningStageString = LearningStageToString(c.LearningStage)
		return
	}

	if acl {
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
	case QueuedToLearn:
		return "Queued To Learn"
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

func (c *Card) GetLearningStage() LearningStage {
	return c.LearningStage
}

func (c *Card) GetLearningStageString() string {
	return LearningStageToString(c.LearningStage)
}

func (c *Card) GetPartsOfSpeech() []string {
	return c.PartsOfSpeech
}

func (c *Card) GetReviewPerformance() float64 {
	if c.TotalTimesReviewed == 0 {
		return 0
	}

	return float64(c.TotalTimesCorrect) / float64(c.TotalTimesReviewed)
}

func (c *Card) IncrementReviewCount() {
	c.TotalTimesReviewed++
}

func (c *Card) IncrementCorrectAnswerCount() {
	c.TotalTimesCorrect++
}

func (c *Card) SetToUpNext() {
	// Check that the card is the in the Available stage
	if c.LearningStage != Available {
		return
	}

	// Set NextReviewDate to 1970-01-01T00:00:00Z
	c.NextReviewDate = time.Unix(0, 0).Format(time.RFC3339)
}

func (c *Card) TagSuspended() {
	c.Tags = []string{"suspended"}
}

func (c *Card) SetQueuedToLearn(cd *CardData) {
	if !c.IsQueueable() {
		return
	}

	c.QueuedToLearn = true
	for _, child := range c.ComponentSubjectIDs {
		cd.Cards[child].SetQueuedToLearn(cd)
	}
}

func (c *Card) IsQueueable() bool {
	if c.LearningStage == Learned || c.LearningStage == Burned || c.LearningStage == UpNext {
		return false
	}

	return true
}
