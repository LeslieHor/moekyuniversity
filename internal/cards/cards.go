package cards

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"foosoft.net/projects/jmdict"
)

type CardData struct {
	CardsFile          string
	DataDir            string
	BackupDir          string
	StaticDir          string
	UpNext             []*Card
	FuncMap            map[string]interface{}
	Cards              map[int]*Card
	Dictionary         jmdict.Jmdict
	DictionaryEntities map[string]string
	// Maps for fast dictionary searching
	DictionaryMap                map[int]*jmdict.JmdictEntry      // Sequence ID -> JmdictEntry
	DictionaryKanjiMap           map[string][]*jmdict.JmdictEntry // Kanji word -> JmdictEntry
	DictionaryReadingMap         map[string][]*jmdict.JmdictEntry // Reading (in hiragana) -> JmdictEntry
	DictionaryNonKanjiReadingMap map[string][]*jmdict.JmdictEntry // Reading -> JmdictEntry
	DictionaryMeaningMap         map[string][]*jmdict.JmdictEntry // Meaning -> JmdictEntry
}

func (cd *CardData) LoadCardJson() {
	log.Println("Loading card data...")

	cardsData := make(map[int]*Card)
	cardsJson, err := ioutil.ReadFile(cd.CardsFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(cardsJson, &cardsData)
	if err != nil {
		log.Fatal(err)
	}

	cd.Cards = cardsData

	log.Printf("Loaded %d cards", len(cd.Cards))

	// Backup cards on start up, then update them
	cd.BackupCardMap()
	cd.UpdateCardData()
	cd.SaveCardMap()
}

func (cd *CardData) UpdateCardData() {
	log.Println("Updating cards...")

	// Keep looping over the cards until no more changes are made
	changesMade := true
	for changesMade {
		log.Println("Iterating over cards")
		changesMade = false
		for _, c := range cd.Cards {
			oldLearningStage := c.LearningStage
			c.UpdateLearningStage(cd)
			if oldLearningStage != c.LearningStage {
				changesMade = true
			}
		}
	}

	// Check the UpNext list and remove any cards that are no longer in the UpNext stage
	for i := 0; i < len(cd.UpNext); i++ {
		if cd.UpNext[i].LearningStage != UpNext {
			cd.UpNext = append(cd.UpNext[:i], cd.UpNext[i+1:]...)
			i--
		}
	}

	log.Println("Updated cards")
}

func (cd *CardData) BackupCardMap() {
	// If the backup directory doesn't exist, create it
	if _, err := os.Stat(cd.BackupDir); os.IsNotExist(err) {
		log.Printf("Creating backup directory %s", cd.BackupDir)
		err := os.Mkdir(cd.BackupDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	t := time.Now()
	backupFilename := filepath.Join(cd.BackupDir, "cards-"+t.Format(time.RFC3339)+".json")
	log.Printf("Backing up cards to %s", backupFilename)
	cd.SaveCardMapToFilename(backupFilename)
}

func (cd *CardData) SaveCardMap() {
	log.Println("Saving cards")
	cd.SaveCardMapToFilename(cd.CardsFile)
}

func (cd *CardData) SaveCardMapToFilename(path string) {
	cardJson, err := json.Marshal(cd.Cards)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path, cardJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

type HistoricalData struct {
	HistoricalDataEntries []*HistoricalDataEntry
	RadicalCount          int
	KanjiCount            int
	VocabularyCount       int
	GrammarCount          int
	RadicalMax            int
	KanjiMax              int
	VocabularyMax         int
	GrammarMax            int
}

type HistoricalDataEntry struct {
	DateTime        string // RFC3339
	RadicalsKnown   int
	KanjiKnown      int
	VocabularyKnown int
	GrammarKnown    int
}

func (cd *CardData) GetHistoricalData() HistoricalData {
	// Load historical data csv
	historicalData := HistoricalData{}
	historicalDataFile, err := os.Open(filepath.Join(cd.DataDir, "historical-data.csv"))
	if err != nil {
		log.Fatal(err)
	}
	defer historicalDataFile.Close()

	historicalDataCsv := csv.NewReader(historicalDataFile)
	historicalDataCsv.Comma = ','

	radicalMax := 0
	kanjiMax := 0
	vocabularyMax := 0
	grammarMax := 0

	for {
		record, err := historicalDataCsv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		dateTime := record[0]
		radicalsKnown, _ := strconv.Atoi(record[1])
		kanjiKnown, _ := strconv.Atoi(record[2])
		vocabularyKnown, _ := strconv.Atoi(record[3])
		grammarKnown, _ := strconv.Atoi(record[4])

		if radicalsKnown > radicalMax {
			radicalMax = radicalsKnown
		}
		if kanjiKnown > kanjiMax {
			kanjiMax = kanjiKnown
		}
		if vocabularyKnown > vocabularyMax {
			vocabularyMax = vocabularyKnown
		}
		if grammarKnown > grammarMax {
			grammarMax = grammarKnown
		}

		historicalDataEntry := HistoricalDataEntry{
			DateTime:        dateTime,
			RadicalsKnown:   radicalsKnown,
			KanjiKnown:      kanjiKnown,
			VocabularyKnown: vocabularyKnown,
			GrammarKnown:    grammarKnown,
		}

		historicalData.HistoricalDataEntries = append(historicalData.HistoricalDataEntries, &historicalDataEntry)
	}

	historicalData.RadicalMax = radicalMax
	historicalData.KanjiMax = kanjiMax
	historicalData.VocabularyMax = vocabularyMax
	historicalData.GrammarMax = grammarMax

	// Calculate total counts
	radicalCount := 0
	kanjiCount := 0
	vocabularyCount := 0
	grammarCount := 0
	for _, c := range cd.Cards {
		switch c.Object {
		case "radical":
			radicalCount++
		case "kanji":
			kanjiCount++
		case "vocabulary":
			vocabularyCount++
		case "grammar":
			grammarCount++
		}
	}

	historicalData.RadicalCount = radicalCount
	historicalData.KanjiCount = kanjiCount
	historicalData.VocabularyCount = vocabularyCount
	historicalData.GrammarCount = grammarCount

	return historicalData
}

func (cd *CardData) SaveHistoricalData() {
	// Gather historical data
	dateTime := time.Now().Format("2006-01-02")
	radicalsKnown := 0
	kanjiKnown := 0
	vocabularyKnown := 0
	grammarKnown := 0

	for _, c := range cd.Cards {
		if c.LearningStage != Learned && c.LearningStage != Burned {
			continue
		}
		switch c.Object {
		case "radical":
			radicalsKnown++
		case "kanji":
			kanjiKnown++
		case "vocabulary":
			vocabularyKnown++
		case "grammar":
			grammarKnown++
		}
	}

	// Save historical data
	csvLine := fmt.Sprintf("%s,%d,%d,%d,%d", dateTime, radicalsKnown, kanjiKnown, vocabularyKnown, grammarKnown)
	historicalDataFile, err := os.OpenFile(filepath.Join(cd.DataDir, "historical-data.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer historicalDataFile.Close()

	log.Printf("Saving historical data: %s", csvLine)
	_, err = historicalDataFile.WriteString(csvLine + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func DoHistoricalData(cd *CardData) {
	for {
		// Wait until the next day at 00:00 and then save historical data
		nextMidnight := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		log.Printf("Waiting until %s to save historical data", nextMidnight.Format("2006-01-02 15:04:05"))
		time.Sleep(time.Until(nextMidnight))
		cd.SaveHistoricalData()
	}
}

func (cd *CardData) AddUpNextCards(n int) {
	// Add n cards to the up next list

	cs := cd.ToList()
	cs = filterCardsByLearningStage(cs, UpNext)
	cs = sortCardsByDue(cs)
	// Invert the card list to prioritise cards with alter due dates.
	// This is because when cards are reviewed, the due date is pushed
	// We want to review the cards we've seen more times, first.
	cs = reverseCards(cs)

	// Take the first n cards from the list and add them to the up next list
	for i := 0; i < n; i++ {
		if len(cs) > 0 {
			cd.UpNext = append(cd.UpNext, cs[0])
			cs = cs[1:]
		}
	}
}

func (cd *CardData) RemoveUpNextCard(id int) {
	for i, c := range cd.UpNext {
		if c.ID == id {
			cd.UpNext = append(cd.UpNext[:i], cd.UpNext[i+1:]...)
			break
		}
	}
}

func (cd *CardData) GetUpNextCards() []*Card {
	// Each time we get a card, rotate the list so that the most recently seen card is last
	if len(cd.UpNext) > 0 {
		cd.UpNext = append(cd.UpNext[1:], cd.UpNext[0])
	}
	return cd.UpNext
}

func (cd *CardData) AddCard(card *Card) {
	cd.Cards[card.ID] = card
}

func (cd *CardData) GetCard(id int) *Card {
	c := cd.Cards[id]
	return c
}

func (cd *CardData) FindVocabulary(vocabulary string) *Card {
	for _, c := range cd.Cards {
		if c.Object == "vocabulary" && (c.Characters == vocabulary || containsString(c.CharactersAlternateWritings, vocabulary)) {
			return c
		}
	}
	return nil
}

func (cd *CardData) FindKanji(kanji string) *Card {
	for _, c := range cd.Cards {
		if c.Object == "kanji" && c.Characters == kanji {
			return c
		}
	}
	return nil
}

func (cd *CardData) FindGrammar(grammar string) *Card {
	for _, c := range cd.Cards {
		if c.Object == "grammar" && c.Characters == grammar {
			return c
		}
	}
	return nil
}

func (cd *CardData) DeleteCard(id int) {
	cd.BackupCardMap()

	delete(cd.Cards, id)

	// Iterate through all cards and remove the deleted card from their
	// ComponentSubjectIDs and AmalgamationSubjectIDs
	for _, c := range cd.Cards {
		c.ComponentSubjectIDs = removeInt(c.ComponentSubjectIDs, id)
		c.AmalgamationSubjectIDs = removeInt(c.AmalgamationSubjectIDs, id)
	}

	log.Printf("Deleted card %d", id)
}

func filterCardsByLearned(cardData []*Card) []*Card {
	return append(filterCardsByLearningStage(cardData, Learned), filterCardsByLearningStage(cardData, Burned)...)
}

func filterCardsByLearningStage(cardData []*Card, learningStage LearningStage) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if card.LearningStage == learningStage {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByLevel(cardData []*Card, level int) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if card.Level == level {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByDueBefore(cardData []*Card, endTime time.Time) []*Card {
	return filterCardsByDueBetween(cardData, time.Time{}, endTime)
}

func filterCardsByDueBetween(cardData []*Card, startTime time.Time, endTime time.Time) []*Card {
	var cards []*Card
	for _, card := range cardData {
		// Convert ISO 8601 string to time.Time

		ct := card.NextReviewDate
		// Skip if there is no next review date
		if ct == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339, ct)
		if err != nil {
			panic(err)
		}

		// If the card is due between the start and end times, add it to the list
		// Inclusive of the start time, exclusive of the end time
		if (t.Equal(startTime) || t.After(startTime)) && t.Before(endTime) {
			cards = append(cards, card)
		}
	}
	return cards
}

// Currently unused
// func filterCardsByHasNextReviewDate(cardData []*Card) []*Card {
// 	var cards []*Card
// 	for _, card := range cardData {
// 		if card.NextReviewDate != "" {
// 			cards = append(cards, card)
// 		}
// 	}
// 	return cards
// }

func filterCardsByType(cardData []*Card, cardType string) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if card.Object == cardType {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByMissingCharacters(cardData []*Card) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if card.Characters == "" {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByMissingCharacterImage(cardData []*Card) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if card.CharacterImage == "" {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByPartsOfSpeech(cardData []*Card, partOfSpeech string) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if containsString(card.GetPartsOfSpeech(), partOfSpeech) {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByReviewPerformance(cardData []*Card, lowerBound float64, upperBound float64) []*Card {
	var cards []*Card
	for _, card := range cardData {
		// Calculate review performance
		performance := card.GetReviewPerformance()
		if performance >= lowerBound && performance <= upperBound {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByTag(cardData []*Card, tag string) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if containsString(card.Tags, tag) {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterCardsByCharacters(cardData []*Card, characters string) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if card.Characters == characters {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterOutCardsByLearningStage(cardData []*Card, learningStage LearningStage) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if card.LearningStage != learningStage {
			cards = append(cards, card)
		}
	}
	return cards
}

func filterOutCardsByTag(cardData []*Card, tag string) []*Card {
	var cards []*Card
	for _, card := range cardData {
		if !containsString(card.Tags, tag) {
			cards = append(cards, card)
		}
	}
	return cards
}

func sortCardsById(cards []*Card) []*Card {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].ID < cards[j].ID
	})
	return cards
}

func sortCardsByDue(cards []*Card) []*Card {
	sort.Slice(cards, func(i, j int) bool {
		// Convert ISO 8601 string to time.Time
		ct := cards[i].NextReviewDate
		// Skip if there is no next review date
		if ct == "" {
			return false
		}
		t, err := time.Parse(time.RFC3339, ct)
		if err != nil {
			panic(err)
		}

		// Do the same for the second card
		ct2 := cards[j].NextReviewDate
		// Skip if there is no next review date
		if ct2 == "" {
			// If the first card has a next review date, but the second doesn't, then the first card should come first
			return true
		}
		t2, err := time.Parse(time.RFC3339, ct2)
		if err != nil {
			panic(err)
		}

		return t.Before(t2)
	})
	return cards
}

func sortCardsByReviewPerformance(cards []*Card) []*Card {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].GetReviewPerformance() > cards[j].GetReviewPerformance()
	})
	return cards
}

func reverseCards(cards []*Card) []*Card {
	for i := len(cards)/2 - 1; i >= 0; i-- {
		opp := len(cards) - 1 - i
		cards[i], cards[opp] = cards[opp], cards[i]
	}
	return cards
}

func (cd *CardData) ToList() []*Card {
	var cardList []*Card
	for _, card := range cd.Cards {
		cardList = append(cardList, card)
	}
	return cardList
}

func (cd *CardData) Search(search string) []*Card {
	var cards []*Card

	// Search characters and meanings
	for _, card := range cd.Cards {
		if strings.Contains(card.Characters, search) {
			if !containsCard(cards, card) {
				cards = append(cards, card)
			}
		}

		for _, meaning := range card.Meanings {
			if strings.Contains(strings.ToLower(meaning.Meaning),
				strings.ToLower(search)) {
				if !containsCard(cards, card) {
					cards = append(cards, card)
				}
			}
		}
	}

	return cards
}

type PartOfSpeech struct {
	Text     string
	Selected bool
}

func (cd *CardData) GetPartsOfSpeech() []string {
	// Loop through all cards and get all parts of speech
	var partsOfSpeech []string
	for _, card := range cd.Cards {
		for _, partOfSpeech := range card.PartsOfSpeech {
			if !containsString(partsOfSpeech, partOfSpeech) {
				partsOfSpeech = append(partsOfSpeech, partOfSpeech)
			}
		}
	}

	// Sort parts of speech
	sort.Strings(partsOfSpeech)
	return partsOfSpeech
}

func (cd *CardData) GetAllPartsOfSpeech(id int) []PartOfSpeech {
	// Loop through all cards and get all parts of speech
	partsOfSpeech := cd.GetPartsOfSpeech()

	// Mark parts of speech in card as selected
	c := cd.Cards[id]
	returnData := []PartOfSpeech{}
	for _, partOfSpeech := range partsOfSpeech {
		ps := PartOfSpeech{
			Text:     partOfSpeech,
			Selected: false,
		}
		if containsString(c.PartsOfSpeech, partOfSpeech) {
			ps.Selected = true
		}
		returnData = append(returnData, ps)
	}

	return returnData
}

func (cd *CardData) GetTags() []string {
	// Loop through all cards and get all tags
	var tags []string
	for _, card := range cd.Cards {
		for _, tag := range card.Tags {
			if !containsString(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}

	// Sort tags
	sort.Strings(tags)
	return tags
}

func (cd *CardData) GetNewCardId() int {
	// Get the highest ID
	var highestId int
	for _, card := range cd.Cards {
		if card.ID > highestId {
			highestId = card.ID
		}
	}

	// Set a minimum ID of 100000 to avoid conflicts with existing cards
	if highestId < 100000 {
		highestId = 100000
	}

	return highestId + 1
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsCard(s []*Card, e *Card) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeInt(s []int, e int) []int {
	for i, a := range s {
		if a == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func removeCard(s []*Card, e *Card) []*Card {
	for i, a := range s {
		if a == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (cd *CardData) GetScheduleData() []ScheduleEntry {
	// Find the counts of reviews for each hour for the next 48 hours
	var scheduleData []ScheduleEntry
	var t1, t2 time.Time
	cards := cd.ToList()

	// Initialise t1 to the next XX:00
	// And set t2 to the next hour
	t1 = time.Now().Truncate(time.Hour)
	t2 = t1.Add(time.Hour)

	for i := 0; i < 300; i++ {
		// Get the number of cards that will be reviewed in this hour period
		cs := filterCardsByDueBetween(cards, t1, t2)
		count := len(cs)
		scheduleData = append(scheduleData, ScheduleEntry{
			Time:  t2.Format("2006-01-02 15:04"),
			Count: count,
		})

		// Increment t1 and t2 by one hour
		t1 = t1.Add(time.Hour)
		t2 = t2.Add(time.Hour)
	}

	return scheduleData
}

type KanjiFrequency struct {
	Name  string        `json:"name"`
	Total int           `json:"total"`
	Data  []interface{} `json:"data"` // <KANJI>,<COUNT>,<PERCENTAGE>
}

func (cd *CardData) GetKanjiFrequencyData() []KanjiFrequencyData {
	// For each file in the kanji frequency data directory,
	// Calculate the known percentage of the kanji.

	files, err := ioutil.ReadDir("data/kanji_frequencies")
	if err != nil {
		log.Fatal(err)
	}

	// Create a list of known kanji
	knownKanji := []string{}
	for _, c := range cd.ToList() {
		if c.Object == "kanji" && (c.LearningStage == Learned || c.LearningStage == Burned) {
			knownKanji = append(knownKanji, c.Characters)
		}
	}

	var kanjiFrequencyData []KanjiFrequencyData

	for _, f := range files {
		filepath := filepath.Join("data/kanji_frequencies", f.Name())
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Read and unmarshal the file
		var data KanjiFrequency
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}

		// Calculate the percentage of known kanji
		knownPercentage := 0.0
		for _, d := range data.Data {
			kanji := d.([]interface{})[0].(string)
			if containsString(knownKanji, kanji) {
				knownPercentage += d.([]interface{})[2].(float64)
			}
		}

		// Add the known percentage to the data
		kanjiFrequencyData = append(kanjiFrequencyData, KanjiFrequencyData{
			Name:         data.Name,
			TotalPercent: int(knownPercentage * 100),
		})
	}

	return kanjiFrequencyData
}
