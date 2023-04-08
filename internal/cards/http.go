package cards

import (
	"net/http"
	"html/template"
	"strconv"
	"encoding/json"
	"time"
	"log"
	"path/filepath"
	"io/ioutil"
	"os"
	"fmt"

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

	r.HandleFunc("/card/{id}", cd.CardHandler)
	r.HandleFunc("/card/{id}/raw", cd.CardRawHandler)
	r.HandleFunc("/card/{id}/json", cd.CardJsonHandler)
	r.HandleFunc("/card/{id}/edit", cd.CardJsonEditHandler)
	r.HandleFunc("/card/{id}/edit/save", cd.CardJsonEditSaveHandler)
	r.HandleFunc("/card/{id}/edit/characterimageupload", cd.CardCharacterImageUploadHandler)

	r.HandleFunc("/cardoverview", cd.OverviewByLearningStageHandler)
	r.HandleFunc("/cardoverview/bylearningstage", cd.OverviewByLearningStageHandler)
	r.HandleFunc("/cardoverview/bylevel", cd.OverviewByLevelHandler)
	r.HandleFunc("/cardoverview/bydue", cd.OverviewByDueHandler)
	r.HandleFunc("/cardoverview/bytype", cd.OverviewByTypeHandler)
	r.HandleFunc("/cardoverview/bypartsofspeech", cd.OverviewByPartsOfSpeechHandler)
	r.HandleFunc("/cardoverview/byreviewperformance", cd.OverviewByReviewPerformanceHandler)
	r.HandleFunc("/cardoverview/debug", cd.OverviewDebugHandler)

	r.HandleFunc("/srs", cd.SrsHandler)
	r.HandleFunc("/srs/correct/{id}", cd.SrsCorrectHandler)
	r.HandleFunc("/srs/incorrect/{id}", cd.SrsIncorrectHandler)

	r.HandleFunc("/search", cd.SearchHandler)

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
	CardDataTree CardDataTree
	LearningStages []struct{
		Value LearningStage
		Text string
	}
	PartsOfSpeech []struct{
		Text string
		Selected bool
	}
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
	editData := CardEditData{
		CardDataTree: dt,
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

type CardOverviewData struct {
	Title string
	Cards []*Card
	LearnedCount int
	ShowLearnedCount bool
}

func NewCardOverviewData(title string, cards []*Card, learnedCount int, showLearnedCount bool) CardOverviewData {
	return CardOverviewData{
		Title: title,
		Cards: cards,
		LearnedCount: learnedCount,
		ShowLearnedCount: showLearnedCount,
	}
}

func filterCards(cards []Card, filters ...func([]Card) []Card) []Card {
	for _, f := range filters {
		cards = f(cards)
	}

	return cards
}

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
		o := NewCardOverviewData("Level " + strconv.Itoa(i), cs, lc, true)
		codl = append(codl, o)
	}

	cd.doTemplate(w, r, "cardoverview.html", codl)
}

func (cd *CardData) OverviewByDueHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	// Due now
	cs := filterCardsByDueBefore(cl, time.Now())
	cs = sortCardsByDue(cs)
	o := NewCardOverviewData("Due now", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 24 hours
	cs = filterCardsByDueBetween(cl, time.Now(), time.Now().Add(24 * time.Hour))
	cs = sortCardsByDue(cs)
	o = NewCardOverviewData("Due in the next 24 hours", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 7 days
	cs = filterCardsByDueBetween(cl, time.Now().Add(24 * time.Hour), time.Now().Add(7 * 24 * time.Hour))
	cs = sortCardsByDue(cs)
	o = NewCardOverviewData("Due in the next week", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 30 days
	cs = filterCardsByDueBetween(cl, time.Now().Add(7 * 24 * time.Hour), time.Now().Add(30 * 24 * time.Hour))
	cs = sortCardsByDue(cs)
	o = NewCardOverviewData("Due in the next month", cs, 0, false)
	codl = append(codl, o)

	// Due in the next 365 days
	cs = filterCardsByDueBetween(cl, time.Now().Add(30 * 24 * time.Hour), time.Now().Add(365 * 24 * time.Hour))
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

func (cd *CardData) SearchHandler(w http.ResponseWriter, r *http.Request) {
	var pageData struct {
		SearchTerm string
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
	c.CorrectAnswer()

	cd.UpdateCardData()
	cd.SaveCardMap()
	
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