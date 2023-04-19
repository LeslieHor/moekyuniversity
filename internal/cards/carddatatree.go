package cards

import (
	"html/template"
	"strings"
)

type CardDataTree struct {
	Card                    *Card
	ComponentSubjects       []CardDataTree
	MeaningMnemonicHtml     template.HTML
	ReadingMnemonicHtml     template.HTML
	AmalgamationSubjectData []AmalgamationSubjectData
	SentencesHtml           []SentenceHtml
}

type AmalgamationSubjectData struct {
	ID                  int
	Object              string
	Characters          string
	Meaning             string
	LearningStage       LearningStage
	LearningStageString string
}

type SentenceHtml struct {
	Japanese template.HTML
	English  template.HTML
}

func (c *Card) GetDataTree(cd *CardData) CardDataTree {
	dt := recursivelyGenerateCardDataTree(cd, c)

	// Convert custom html tags to span tags
	dt.MeaningMnemonicHtml = template.HTML(customHtmlTagsToSpan(c.MeaningMnemonic))
	dt.ReadingMnemonicHtml = template.HTML(customHtmlTagsToSpan(c.ReadingMnemonic))

	var sentencesHtml []SentenceHtml
	for _, s := range c.Sentences {
		sentencesHtml = append(sentencesHtml, SentenceHtml{
			Japanese: template.HTML(customHtmlTagsToSpan(s.Japanese)),
			English:  template.HTML(customHtmlTagsToSpan(s.English)),
		})
	}
	dt.SentencesHtml = sentencesHtml

	// Generate amalgamation subject data
	for _, id := range c.AmalgamationSubjectIDs {
		dt.AmalgamationSubjectData = append(dt.AmalgamationSubjectData, AmalgamationSubjectData{
			ID:                  id,
			Object:              cd.GetCard(id).Object,
			Characters:          cd.GetCard(id).Characters,
			Meaning:             cd.GetCard(id).Meanings[0].Meaning,
			LearningStage:       cd.GetCard(id).LearningStage,
			LearningStageString: cd.GetCard(id).LearningStageString,
		})
	}

	return dt
}

// Recursively generate card data tree, where ComponentSubjects is a list of dependencies
func recursivelyGenerateCardDataTree(cd *CardData, c *Card) CardDataTree {
	dt := CardDataTree{
		Card: c,
	}

	for _, id := range c.ComponentSubjectIDs {
		dt.ComponentSubjects = append(dt.ComponentSubjects, recursivelyGenerateCardDataTree(cd, cd.GetCard(id)))
	}

	return dt
}

// Convert custom html tags to span tags
// e.g. <radical>text</radical> -> <span class="inline-highlight radical-highlight">text</span>
func customHtmlTagsToSpan(str string) string {
	str = strings.Replace(str, "<radical>", "<span class=\"inline-highlight radical-highlight\">", -1)
	str = strings.Replace(str, "</radical>", "</span>", -1)
	str = strings.Replace(str, "<kanji>", "<span class=\"inline-highlight kanji-highlight\">", -1)
	str = strings.Replace(str, "</kanji>", "</span>", -1)
	str = strings.Replace(str, "<vocabulary>", "<span class=\"inline-highlight vocabulary-highlight\">", -1)
	str = strings.Replace(str, "</vocabulary>", "</span>", -1)
	str = strings.Replace(str, "<grammar>", "<span class=\"inline-highlight grammar-highlight\">", -1)
	str = strings.Replace(str, "</grammar>", "</span>", -1)

	return str
}

// CardDataTree methods

func (dt *CardDataTree) IsAllChildrenLearned() bool {
	for _, child := range dt.ComponentSubjects {
		if child.Card.LearningStage != Learned &&
			child.Card.LearningStage != Burned {
			return false
		}
	}

	return true
}
