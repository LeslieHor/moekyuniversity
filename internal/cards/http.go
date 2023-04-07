package cards

import (
	"net/http"
	"html/template"
	"strconv"
	"encoding/json"
	"time"
	"log"
	"regexp"
	"strings"
	"path/filepath"

	"github.com/gorilla/mux"
)

func SetupRoutes(cd *CardData) {
	r := mux.NewRouter()

	r.HandleFunc("/", cd.IndexHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(cd.StaticDir))))
	r.HandleFunc("/stylesheet.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/css/stylesheet.css")
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/icon.png")
	})
	r.HandleFunc("/img/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.ServeFile(w, r, "static/img/" + vars["name"])
	})

	r.HandleFunc("/card/{id}", cd.CardHandler)
	r.HandleFunc("/card/{id}/json", cd.CardRawHandler)
	r.HandleFunc("/card/{id}/edit", cd.CardEditHandler)
	r.HandleFunc("/card/{id}/edit/save", cd.CardEditPostHandler)

	r.HandleFunc("/cardoverview", cd.OverviewByLearningStageHandler)
	r.HandleFunc("/cardoverview/bylearningstage", cd.OverviewByLearningStageHandler)
	r.HandleFunc("/cardoverview/bylevel", cd.OverviewByLevelHandler)
	r.HandleFunc("/cardoverview/bydue", cd.OverviewByDueHandler)
	r.HandleFunc("/cardoverview/bytype", cd.OverviewByTypeHandler)
	r.HandleFunc("/cardoverview/bypartsofspeech", cd.OverviewByPartsOfSpeechHandler)
	r.HandleFunc("/cardoverview/debug", cd.OverviewDebugHandler)

	r.HandleFunc("/srs", cd.SrsHandler)
	r.HandleFunc("/srs/correct/{id}", cd.SrsCorrectHandler)
	r.HandleFunc("/srs/incorrect/{id}", cd.SrsIncorrectHandler)

	r.HandleFunc("/search", cd.SearchHandler)

	http.ListenAndServe(":8080", r)
}

func (cd *CardData) doTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	htmlDir := filepath.Join(cd.StaticDir, "html")
	templatemainFile := filepath.Join(htmlDir, "templatemain.html")
	templateFile := filepath.Join(htmlDir, templateName)
	t, err := template.ParseFiles(templatemainFile, templateFile)
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

func (cd *CardData) CardEditHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	c := cd.GetCard(id)
	dt := c.GetDataTree(cd)

	// Get the learning stages and their string representation
	var learningStages []struct{
		Value LearningStage
		Text string
	}
	for _, ls := range LearningStages {
		learningStages = append(learningStages, struct{
			Value LearningStage
			Text string
		}{
			Value: ls,
			Text: LearningStageToString(ls),
		})
	}

	// Get all the parts of speech whether the current card has them or not
	var partsOfSpeech []struct{
		Text string
		Selected bool
	}
	for _, pos := range cd.GetPartsOfSpeech() {
		// Check if the current card has the part of speech
		selected := containsString(c.PartsOfSpeech, pos)

		partsOfSpeech = append(partsOfSpeech, struct{
			Text string
			Selected bool
		}{
			Text: pos,
			Selected: selected,
		})
	}

	editData := CardEditData{
		CardDataTree: dt,
		LearningStages: learningStages,
		PartsOfSpeech: partsOfSpeech,
	}

	cd.doTemplate(w, r, "cardedit.html", editData)
}

func (cd *CardData) CardEditPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get card ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	c := cd.GetCard(id)

	// Get the card data from the form
	r.ParseForm()
	
	// Get component subject ids and parse into a slice of ints
	// Example input from form: '["1", "2", "3"]'
	componentSubjectIds := r.FormValue("componentSubjectIds")
	log.Printf("componentSubjectIds: %s", componentSubjectIds)
	
	// Strip any non numeric or non-comma characters from the string
	re := regexp.MustCompile("[^0-9,]")
	componentSubjectIds = re.ReplaceAllString(componentSubjectIds, "")

	// Split the string by commas
	s := strings.Split(componentSubjectIds, ",")

	var componentSubjectIdsInt []int
	// Parse the string slice into an int slice
	for _, id := range s {
		log.Printf("id: %s", id)
		idInt, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		componentSubjectIdsInt = append(componentSubjectIdsInt, idInt)
	}
	log.Printf("componentSubjectIdsInt: %v", componentSubjectIdsInt)

	// Assign the card data
	c.ComponentSubjectIDs = componentSubjectIdsInt

	cd.UpdateCardData()
	cd.SaveCardMap()
	// Redirect to the card page
	http.Redirect(w, r, "/card/" + strconv.Itoa(id), http.StatusSeeOther)
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

func (cd *CardData) OverviewDebugHandler(w http.ResponseWriter, r *http.Request) {
	codl := []CardOverviewData{}
	cl := cd.ToList()

	cs := filterCardsByMissingCharacters(cl)
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