package cards

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func SetupRoutes(cd *CardData) {
	cd.SetupFuncMap()

	log.Printf("Data dir: %s", cd.DataDir)

	r := mux.NewRouter()

	r.HandleFunc("/", cd.IndexHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(cd.StaticDir))))
	r.HandleFunc("/stylesheet.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/css/stylesheet.css")
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/icon.png")
	})
	r.PathPrefix("/img/").Handler(http.FileServer(http.Dir(cd.StaticDir)))

	// Strip the /data prefix from the path and serve the file from the data directory
	r.PathPrefix("/data/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(cd.DataDir, r.URL.Path[6:]))
	})

	r.HandleFunc("/card/new", cd.CardNewHandler)
	r.HandleFunc("/card/{id}", cd.CardHandler)
	r.HandleFunc("/card/{id}/raw", cd.CardRawHandler)
	r.HandleFunc("/card/{id}/json", cd.CardJsonHandler)
	r.HandleFunc("/card/{id}/edit", cd.CardJsonEditHandler)
	r.HandleFunc("/card/{id}/edit/save", cd.CardJsonEditSaveHandler)
	r.HandleFunc("/card/{id}/edit/characterimageupload", cd.CardCharacterImageUploadHandler)
	r.HandleFunc("/card/{id}/delete", cd.CardDeleteHandler)
	r.HandleFunc("/card/{id}/tagsuspended", cd.CardTagSuspendedHandler)
	r.HandleFunc("/card/{id}/addtoqueue", cd.CardAddToQueueHandler)

	r.HandleFunc("/cardoverview", cd.OverviewByLearningStageHandler)
	r.HandleFunc("/cardoverview/bylearningstage", cd.OverviewByLearningStageHandler)
	r.HandleFunc("/cardoverview/bylevel", cd.OverviewByLevelHandler)
	r.HandleFunc("/cardoverview/bydue", cd.OverviewByDueHandler)
	r.HandleFunc("/cardoverview/bytype", cd.OverviewByTypeHandler)
	r.HandleFunc("/cardoverview/bypartsofspeech", cd.OverviewByPartsOfSpeechHandler)
	r.HandleFunc("/cardoverview/byreviewperformance", cd.OverviewByReviewPerformanceHandler)
	r.HandleFunc("/cardoverview/bytag", cd.OverviewByTagHandler)
	r.HandleFunc("/cardoverview/simulate/{correctRate}/{newCardsPerDay}", cd.OverviewSimulateHandler)
	r.HandleFunc("/cardoverview/debug", cd.OverviewDebugHandler)

	r.HandleFunc("/textanalysis", cd.TextAnalysisHandler)
	r.HandleFunc("/textanalysis/new", cd.TextAnalysisNewHandler)
	r.HandleFunc("/textanalysis/new/submit", cd.TextAnalysisNewSubmitHandler)
	r.HandleFunc("/textanalysis/{id}", cd.TextAnalysisIdHandler)
	r.HandleFunc("/textanalysis/{id}/delete", cd.TextAnalysisIdDeleteHandler)

	r.HandleFunc("/srs", cd.SrsHandler)
	r.HandleFunc("/srs/correct/{id}", cd.SrsCorrectHandler)
	r.HandleFunc("/srs/incorrect/{id}", cd.SrsIncorrectHandler)
	r.HandleFunc("/srs/addupnextcards/{n}", cd.SrsAddUpNextCardsHandler)

	r.HandleFunc("/search", cd.SearchHandler)
	r.HandleFunc("/dictionarysearch", cd.DictionarySearchHandler)

	http.ListenAndServe(":8080", r)
}

func (cd *CardData) SetupFuncMap() {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"div": func(a, b int) int {
			return a / b
		},
		"mod": func(a, b int) int {
			return a % b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"percent": func(a, b int) int {
			return a * 100 / b
		},
		"queueable": func(c *Card) bool {
			return c.IsQueueable()
		},
		"stripspaces": func(s string) string {
			return strings.ReplaceAll(s, " ", "")
		},
		"newlinetohtml": func(s string) template.HTML {
			return template.HTML(strings.ReplaceAll(s, "\r\n", "<br>"))
		},
	}
	cd.FuncMap = funcMap
}

func (cd *CardData) doTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	htmlDir := filepath.Join(cd.StaticDir, "html")
	templatemainFile := filepath.Join(htmlDir, "templatemain.html")
	templateFile := filepath.Join(htmlDir, templateName)
	// Load template with custom function map
	t, err := template.New("templatemain.html").Funcs(cd.FuncMap).ParseFiles(templatemainFile, templateFile)
	if err != nil {
		panic(err)
	}

	t.Execute(w, data)
}

func (cd *CardData) IndexHandler(w http.ResponseWriter, r *http.Request) {
	cd.doTemplate(w, r, "index.html", nil)
}

func (cd *CardData) ServeFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving file: %s", r.URL.Path[1:])
	http.ServeFile(w, r, filepath.Join(cd.StaticDir, r.URL.Path[1:]))
}

func (cd *CardData) CardHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	c := cd.GetCard(id)
	dt := c.GetDataTree(cd)

	cd.doTemplate(w, r, "card.html", dt)
}

func (cd *CardData) CardRawHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	c := cd.GetCard(id)
	dt := c.GetDataTree(cd)

	// Convert the card data tree to json and write it to the response
	json, err := json.Marshal(dt)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (cd *CardData) CardJsonHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	c := cd.GetCard(id)

	// Convert the card data tree to json and write it to the response
	json, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

type CardEditData struct {
	CardDataTree   CardDataTree
	LearningStages []struct {
		Value LearningStage
		Text  string
	}
	PartsOfSpeech []struct {
		Text     string
		Selected bool
	}
	SuggestedComponents []*Card
}

func (cd *CardData) CardJsonEditHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	c := cd.GetCard(id)
	dt := c.GetDataTree(cd)
	suggestedComponents := filterCardsByCharacters(cd.ToList(), c.Characters)
	suggestedComponents = removeCard(suggestedComponents, c)

	editData := CardEditData{
		CardDataTree:        dt,
		SuggestedComponents: suggestedComponents,
	}

	cd.doTemplate(w, r, "cardjsonedit.html", editData)
}

func (cd *CardData) CardJsonEditSaveHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	log.Printf("Saving card %d", id)

	// Get the json data from the request
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Log the json data
	log.Printf("JSON data: %s", jsonData)

	// Attempt to unmarshal the json data into a card
	var c Card
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		log.Printf("Error unmarshalling json data: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate the NextReviewDate is a valid date
	if c.NextReviewDate != "" {
		_, err = time.Parse(time.RFC3339, c.NextReviewDate)
		if err != nil {
			log.Printf("Error parsing NextReviewDate: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Save the card
	cd.Cards[id] = &c
	cd.UpdateCardData()
	cd.SaveCardMap()

	// Send an ok response
	w.WriteHeader(http.StatusOK)
}

func (cd *CardData) CardCharacterImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	log.Printf("Uploading image for card %d", id)

	// Get the image data from the POST request
	imageData, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error getting image data: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the image data
	imageDataBytes, err := ioutil.ReadAll(imageData)
	if err != nil {
		log.Printf("Error reading image data: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := cd.GetCard(id)
	filename := fmt.Sprintf("%d_characterimage.png", c.ID)
	file, err := os.Create(filepath.Join(cd.DataDir, "img", filename))
	if err != nil {
		log.Printf("Error creating image file: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	log.Printf("Writing image data to file %s", file.Name())
	file.Write(imageDataBytes)

	// Save the image data
	c.CharacterImage = filename
	cd.UpdateCardData()
	cd.SaveCardMap()

	// Send an ok response
	w.WriteHeader(http.StatusOK)
}

func (cd *CardData) CardNewHandler(w http.ResponseWriter, r *http.Request) {
	m := Meaning{
		Meaning:        "Placeholder meaning",
		Primary:        true,
		AcceptedAnswer: true,
	}
	r1 := Reading{
		Reading:        "Placeholder reading",
		Type:           "onyomi",
		Primary:        true,
		AcceptedAnswer: true,
	}
	r2 := Reading{
		Reading:        "Placeholder reading",
		Type:           "kunyomi",
		Primary:        false,
		AcceptedAnswer: false,
	}

	c := Card{
		ID:       cd.GetNewCardId(),
		Meanings: []Meaning{m},
		Readings: []Reading{r1, r2},
	}

	cd.Cards[c.ID] = &c
	cd.UpdateCardData()
	cd.SaveCardMap()

	http.Redirect(w, r, fmt.Sprintf("/card/%d", c.ID), http.StatusFound)
}

func (cd *CardData) CardDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	log.Printf("Deleting card %d", id)

	// Delete the card
	cd.DeleteCard(id)
	cd.UpdateCardData()
	cd.SaveCardMap()

	http.Redirect(w, r, "/", http.StatusFound)
}

func (cd *CardData) CardTagSuspendedHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	log.Printf("Tagging card %d as suspended", id)

	c := cd.GetCard(id)
	c.TagSuspended()
	cd.UpdateCardData()
	cd.SaveCardMap()

	// Redirect to the card page
	http.Redirect(w, r, fmt.Sprintf("/card/%d", id), http.StatusFound)
}

func (cd *CardData) CardAddToQueueHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	log.Printf("Adding card %d to queue", id)

	// Check that the card is available or unavailable
	c := cd.GetCard(id)
	if !c.IsQueueable() {
		log.Printf("Card %d is not queuable", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add the card to the queue
	c.SetQueuedToLearn(cd)
	cd.UpdateCardData()
	cd.SaveCardMap()

	// Redirect to the card page
	http.Redirect(w, r, fmt.Sprintf("/card/%d", id), http.StatusFound)
}

type CardOverviewData struct {
	Title            string
	Cards            []*Card
	LearnedCount     int
	ShowLearnedCount bool
}

func NewCardOverviewData(title string, cards []*Card, learnedCount int, showLearnedCount bool) CardOverviewData {
	return CardOverviewData{
		Title:            title,
		Cards:            cards,
		LearnedCount:     learnedCount,
		ShowLearnedCount: showLearnedCount,
	}
}

// Currently unused
// Filter a list of cards by the given functions
// func filterCards(cards []Card, filters ...func([]Card) []Card) []Card {
// 	for _, f := range filters {
// 		cards = f(cards)
// 	}

// 	return cards
// }

func getOverviewLearningStage(cards []*Card, ls LearningStage) CardOverviewData {
	cs := filterCardsByLearningStage(cards, ls)
	cs = sortCardsById(cs)
	o := NewCardOverviewData(LearningStageToString(ls), cs, 0, false)
	return o
}

func (cd *CardData) OverviewByLearningStageHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	codl = append(codl, getOverviewLearningStage(cl, Unavailable))
	codl = append(codl, getOverviewLearningStage(cl, Available))
	codl = append(codl, getOverviewLearningStage(cl, QueuedToLearn))
	codl = append(codl, getOverviewLearningStage(cl, UpNext))
	codl = append(codl, getOverviewLearningStage(cl, Learning))
	codl = append(codl, getOverviewLearningStage(cl, Learned))
	codl = append(codl, getOverviewLearningStage(cl, Burned))

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewByLevelHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	for i := 1; i <= 60; i++ {
		cs := filterCardsByLevel(cl, i)
		cs = sortCardsById(cs)
		lc := len(filterCardsByLearned(cs))
		o := NewCardOverviewData("Level "+strconv.Itoa(i), cs, lc, true)
		codl = append(codl, o)
	}

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewByDueHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	// Due now
	cs := filterCardsByDueBefore(cl, time.Now())
	cs = filterOutCardsByLearningStage(cs, UpNext)
	cs = filterOutCardsByTag(cs, "suspended")
	cs = sortCardsByDue(cs)
	o := NewCardOverviewData("Due now", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 24 hours
	cs = filterCardsByDueBetween(cl, time.Now(), time.Now().Add(24*time.Hour))
	cs = filterOutCardsByTag(cs, "suspended")
	cs = sortCardsByDue(cs)
	o = NewCardOverviewData("Due in the next 24 hours", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 7 days
	cs = filterCardsByDueBetween(cl, time.Now().Add(24*time.Hour), time.Now().Add(7*24*time.Hour))
	cs = filterOutCardsByTag(cs, "suspended")
	cs = sortCardsByDue(cs)
	o = NewCardOverviewData("Due in the next week", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 30 days
	cs = filterCardsByDueBetween(cl, time.Now().Add(7*24*time.Hour), time.Now().Add(30*24*time.Hour))
	cs = filterOutCardsByTag(cs, "suspended")
	cs = sortCardsByDue(cs)
	o = NewCardOverviewData("Due in the next month", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 365 days
	cs = filterCardsByDueBetween(cl, time.Now().Add(30*24*time.Hour), time.Now().Add(365*24*time.Hour))
	cs = filterOutCardsByTag(cs, "suspended")
	cs = sortCardsByDue(cs)
	o = NewCardOverviewData("Due in the next year", cs, 0, false)
	codl = append(codl, o)

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewByTypeHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	cs := filterCardsByType(cl, "radical")
	cs = sortCardsById(cs)
	lc := len(filterCardsByLearned(cs))
	o := NewCardOverviewData("Radicals", cs, lc, true)
	codl = append(codl, o)

	cs = filterCardsByType(cl, "kanji")
	cs = sortCardsById(cs)
	lc = len(filterCardsByLearned(cs))
	o = NewCardOverviewData("Kanji", cs, lc, true)
	codl = append(codl, o)

	cs = filterCardsByType(cl, "vocabulary")
	cs = sortCardsById(cs)
	lc = len(filterCardsByLearned(cs))
	o = NewCardOverviewData("Vocabulary", cs, lc, true)
	codl = append(codl, o)

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewByPartsOfSpeechHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	partsOfSpeech := cd.GetPartsOfSpeech()

	for _, pos := range partsOfSpeech {
		cs := filterCardsByPartsOfSpeech(cl, pos)
		cs = sortCardsById(cs)
		lc := len(filterCardsByLearned(cs))
		o := NewCardOverviewData(pos, cs, lc, true)
		codl = append(codl, o)
	}

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewByReviewPerformanceHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cardList := cd.ToList()
	cl := append(filterCardsByLearningStage(cardList, Learning), filterCardsByLearningStage(cardList, Learned)...)
	cl = append(cl, filterCardsByLearningStage(cardList, Burned)...)

	cs := filterCardsByReviewPerformance(cl, 0, 0.5)
	cs = sortCardsByReviewPerformance(cs)
	o := NewCardOverviewData("0% - 50%", cs, 0, false)
	codl = append(codl, o)

	cs = filterCardsByReviewPerformance(cl, 0.5, 0.75)
	cs = sortCardsByReviewPerformance(cs)
	o = NewCardOverviewData("50% - 75%", cs, 0, false)
	codl = append(codl, o)

	cs = filterCardsByReviewPerformance(cl, 0.75, 0.95)
	cs = sortCardsByReviewPerformance(cs)
	o = NewCardOverviewData("75% - 95%", cs, 0, false)
	codl = append(codl, o)

	cs = filterCardsByReviewPerformance(cl, 0.95, 1)
	cs = sortCardsByReviewPerformance(cs)
	o = NewCardOverviewData("95% - 100%", cs, 0, false)
	codl = append(codl, o)

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewByTagHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	tags := cd.GetTags()

	for _, tag := range tags {
		cs := filterCardsByTag(cl, tag)
		cs = sortCardsById(cs)
		lc := len(filterCardsByLearned(cs))
		o := NewCardOverviewData(tag, cs, lc, true)
		codl = append(codl, o)
	}

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewSimulateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	correctRate := vars["correctRate"]
	correctRateFloat, err := strconv.ParseFloat(correctRate, 64)
	if err != nil {
		log.Println(err)
		return
	}
	newCardsPerDay := vars["newCardsPerDay"]
	newCardsPerDayInt, err := strconv.Atoi(newCardsPerDay)
	if err != nil {
		log.Println(err)
		return
	}

	codl := []CardOverviewData{}
	cl := cd.ToList()
	cl = filterCardsByLearned(cl)
	cl = filterOutCardsByLearningStage(cl, Burned)

	// Create a copy of the card list
	for i, c := range cl {
		if c == nil {
			continue
		}
		v := *c
		cl[i] = &v
	}

	// Begin with all cards due now
	t := time.Now()

	// Simulate the next 60 days
	newCardCounter := 0
	for i := 0; i < 60; i++ {
		// Add new cards to the list
		var nc []*Card
		for j := 0; j < newCardsPerDayInt; j++ {
			newCardCounter++
			nc = append(nc, &Card{
				ID:         newCardCounter,
				Characters: fmt.Sprintf("Card %d", newCardCounter),
				Meanings: []Meaning{
					{
						Meaning:        fmt.Sprintf("Card %d", newCardCounter),
						Primary:        true,
						AcceptedAnswer: true,
					},
				},
				Interval:         0,
				LearningInterval: 3,
				NextReviewDate:   "1970-01-01T00:00:00Z",
			})
		}
		cl = append(cl, nc...)

		// Update all learning stages
		for _, c := range cl {
			// TODO: This is a hack to get the learning stage to update. I don't really want to use the real card data here.
			c.UpdateLearningStage(cd)
		}

		// Remove burned cards
		cl = filterOutCardsByLearningStage(cl, Burned)
		// Get the cards due today
		cs := filterCardsByDueBefore(cl, t)

		// Fake review the cards
		for _, c := range cs {
			if rand.Float64() < correctRateFloat {
				c.ProcessCorrectAnswer()
			} else {
				c.ProcessIncorrectAnswer()
			}
		}

		// Add those cards we just reviewed to the overview
		codl = append(codl, CardOverviewData{
			Title:            fmt.Sprintf("Day %d", i),
			Cards:            cs,
			ShowLearnedCount: false,
		})

		// The calculated next review date is based on time.now(), so we need to fake it into the future.
		// Get the time difference between now and the next review date and add on the t time., but if the time is in the past, don't change it, because it is an upnext learning card.
		for _, c := range cs {
			if c.NextReviewDate == "" {
				// Burned cards don't have a next review date
				continue
			}
			t1, err := time.Parse(time.RFC3339, c.NextReviewDate)
			if err != nil {
				panic(err)
			}
			t2 := time.Now()
			if t1.Before(t2) {
				continue
			}
			t3 := t1.Sub(t2)
			c.NextReviewDate = t.Add(t3).Format(time.RFC3339)
		}

		// Move to the next day
		t = t.Add(24 * time.Hour)
	}

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewDebugHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	cs := filterCardsByMissingCharacters(cl)
	cs = filterCardsByMissingCharacterImage(cs)
	cs = sortCardsById(cs)
	o := NewCardOverviewData("Missing characters", cs, 0, false)
	codl = append(codl, o)

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) TextAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	// Read all the files in the text analysis directory
	files, err := ioutil.ReadDir("data/text_analysis")
	if err != nil {
		log.Fatal(err)
	}

	var taList []TextAnalysis
	for _, f := range files {
		// Get the json data for the id
		filepath := filepath.Join(cd.DataDir, "text_analysis", f.Name())
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Decode the json data
		var ta TextAnalysis
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&ta)
		if err != nil {
			log.Fatal(err)
		}

		taList = append(taList, ta)
	}

	// Sort the list by Name
	sort.Slice(taList, func(i, j int) bool {
		return taList[i].Name < taList[j].Name
	})

	pageData := struct {
		TextAnalysisList []TextAnalysis
	}{
		TextAnalysisList: taList,
	}

	cd.doTemplate(w, r, "textanalysisoverview.html", pageData)
}

func (cd *CardData) TextAnalysisIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the json data for the id
	filepath := filepath.Join(cd.DataDir, "text_analysis", id+".json")
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Decode the json data
	var ta TextAnalysis
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&ta)
	if err != nil {
		log.Fatal(err)
	}
	ta.Analyse(cd)

	// Replace newlines with <br>
	htmlText := strings.Replace(ta.Text, "\r\n", "<br>", -1)
	htmlText = strings.Replace(htmlText, "\r", "<br>", -1)
	htmlText = strings.Replace(htmlText, "\n", "<br>", -1)

	pageData := struct {
		TextAnalysis TextAnalysis
		HTMLSafeText template.HTML
	}{
		TextAnalysis: ta,
		HTMLSafeText: template.HTML(htmlText),
	}

	cd.doTemplate(w, r, "textanalysis.html", pageData)
}

func (cd *CardData) TextAnalysisNewHandler(w http.ResponseWriter, r *http.Request) {
	cd.doTemplate(w, r, "textanalysisnew.html", nil)
}

func (cd *CardData) TextAnalysisNewSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Get the text from the form
	name := r.FormValue("name")
	text := r.FormValue("text")

	log.Printf("New text analysis: %s", name)
	log.Printf("Text: %s", text)

	// Create a new text analysis
	ta := TextAnalysis{
		ID:   uuid.New().String(),
		Name: name,
		Text: text,
	}

	// Save the text analysis
	filepath := filepath.Join(cd.DataDir, "text_analysis", ta.ID+".json")
	ta.Save(filepath)

	// Redirect to the text analysis page
	http.Redirect(w, r, "/textanalysis/"+ta.ID, http.StatusFound)
}

func (cd *CardData) TextAnalysisIdDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete the file
	filepath := filepath.Join(cd.DataDir, "text_analysis", id+".json")
	err := os.Remove(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// Redirect to the text analysis overview page
	http.Redirect(w, r, "/textanalysis", http.StatusFound)
}

func (cd *CardData) SearchHandler(w http.ResponseWriter, r *http.Request) {
	var pageData struct {
		SearchTerm    string
		SearchResults []*Card
	}
	// Get search query "q"
	values := r.URL.Query()
	q := values.Get("q")

	searchResults := cd.Search(q)

	pageData.SearchTerm = q
	pageData.SearchResults = searchResults

	log.Printf("Search for %s returned %d results", q, len(searchResults))

	cd.doTemplate(w, r, "search.html", pageData)
}

type DictionarySearchData struct {
	Tokens            []string
	DictSearchTerm    string
	DictSearchResults []DictionaryEntry
}

func (cd *CardData) DictionarySearchHandler(w http.ResponseWriter, r *http.Request) {
	// Get search query "q"
	values := r.URL.Query()
	q := values.Get("q")

	searchResults := SearchDictionary(cd, q)

	log.Printf("Dictionary search for %s returned %d results", q, len(searchResults.DictSearchResults))

	cd.doTemplate(w, r, "dictionarysearch.html", searchResults)
}

func (cd *CardData) SrsHandler(w http.ResponseWriter, r *http.Request) {
	srsData := cd.GetNextSrsCard()
	cd.doTemplate(w, r, "srs.html", srsData)
}

func (cd *CardData) SrsCorrectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting id to int: %s", err)
		return
	}

	log.Printf("Correct answer for card %d", cardId)
	c := cd.GetCard(cardId)
	prevState := c.GetLearningStageString()

	c.CorrectAnswer()

	cd.UpdateCardData()
	cd.SaveCardMap()

	currentState := c.GetLearningStageString()
	if currentState != prevState {
		log.Printf("Card %d changed from %s to %s", cardId, prevState, currentState)
		s := cd.GetNextSrsCard()
		pageData := struct {
			Card          *Card
			DueCount      int
			LearningCount int
		}{
			Card:          c,
			DueCount:      s.DueCount,
			LearningCount: s.LearningCount,
		}

		cd.doTemplate(w, r, "congratulationssrs.html", pageData)
		return
	}

	http.Redirect(w, r, "/srs", http.StatusFound)
}

func (cd *CardData) SrsIncorrectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting id to int: %s", err)
		return
	}

	log.Printf("Incorrect answer for card %d", cardId)
	c := cd.GetCard(cardId)
	c.IncorrectAnswer()

	cd.UpdateCardData()
	cd.SaveCardMap()

	http.Redirect(w, r, "/srs", http.StatusFound)
}

func (cd *CardData) SrsAddUpNextCardsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		log.Printf("Error converting number to int: %s", err)
		return
	}

	log.Printf("Adding %d cards to up next", n)
	cd.AddUpNextCards(n)

	http.Redirect(w, r, "/srs", http.StatusFound)
}
