package cards

import (
	"log"
	"os"

	"foosoft.net/projects/jmdict"
	"github.com/mochi-co/kana-tools"
)

type DictionaryEntry struct {
	Expressions []string // Word in kanji
	Readings    []string // Pronunciations in hiragana
	Definitions []DictionaryDefinition
	JmdictEntry jmdict.JmdictEntry // The original JMdict entry
}

type DictionaryDefinition struct {
	PartsOfSpeech []string
	Definitions   []string // List of translations
}

func (cd *CardData) LoadDictionary() {
	log.Printf("Loading dictionary...")

	f, err := os.Open(cd.DataDir + "/JMdict_e")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dict, entities, err := jmdict.LoadJmdict(f)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded dictionary with %d entries", len(dict.Entries))

	log.Printf("Building dictionary index...")
	// create map by kanji and by readings
	var kanjiMap = make(map[string][]*jmdict.JmdictEntry)
	var readingMap = make(map[string][]*jmdict.JmdictEntry)
	var nonKanjiReadingMap = make(map[string][]*jmdict.JmdictEntry)
	for i := range dict.Entries {
		entry := &dict.Entries[i]
		for _, kanji := range entry.Kanji {
			kanjiMap[kanji.Expression] = append(kanjiMap[kanji.Expression], entry)
		}

		usuallyKana := false
		for _, sense := range entry.Sense {
			for _, m := range sense.Misc {
				if m == "word usually written using kana alone" {
					usuallyKana = true
					break
				}
			}
		}

		IsNonKanji := usuallyKana || (len(entry.Kanji) == 0)
		for _, reading := range entry.Readings {
			if IsNonKanji {
				nonKanjiReadingMap[reading.Reading] = append(readingMap[reading.Reading], entry)
			}
			readingMap[kana.ToHiragana(reading.Reading)] = append(readingMap[reading.Reading], entry)
		}
	}

	log.Printf("Index built")
	cd.Dictionary = dict
	cd.DictionaryKanjiMap = kanjiMap
	cd.DictionaryReadingMap = readingMap
	cd.DictionaryNonKanjiReadingMap = nonKanjiReadingMap
	cd.DictionaryEntities = entities
}

func (t *Token) AddDictionaryEntry(cd *CardData) {
	var matches []*jmdict.JmdictEntry

	// Search on the kanji word
	matches = cd.DictionaryKanjiMap[t.BaseForm]
	if len(matches) > 0 {
		var convertedEntries []DictionaryEntry
		for _, entry := range matches {
			convertedEntries = append(convertedEntries, convertJmdictEntryToDictionaryEntry(*entry))
		}
		t.DictionaryEntries = append(t.DictionaryEntries, convertedEntries...)
	}

	// Search on the pure hiragana word (in case the word is written in hiragana)
	matches = cd.DictionaryNonKanjiReadingMap[t.BaseForm]
	if len(matches) > 0 {
		var convertedEntries []DictionaryEntry
		for _, entry := range matches {
			convertedEntries = append(convertedEntries, convertJmdictEntryToDictionaryEntry(*entry))
		}
		t.DictionaryEntries = append(t.DictionaryEntries, convertedEntries...)
	}

	// If there have been matches for the word itself, we don't need to look for
	// matches for the pronunciation
	if len(t.DictionaryEntries) > 0 {
		return
	}

	// Search on the pronunciation
	// This is a fallback, in case the word is not written in the typical kanji
	matches = cd.DictionaryReadingMap[kana.ToHiragana(t.Pronunciation)]
	if len(matches) > 0 {
		var convertedEntries []DictionaryEntry
		for _, entry := range matches {
			convertedEntries = append(convertedEntries, convertJmdictEntryToDictionaryEntry(*entry))
		}
		t.DictionaryEntries = append(t.DictionaryEntries, convertedEntries...)
		return
	}
}

// Convert a JMdict entry to a dictionary entry for ease of use
func convertJmdictEntryToDictionaryEntry(entry jmdict.JmdictEntry) DictionaryEntry {
	var de DictionaryEntry

	de.JmdictEntry = entry

	for _, kanji := range entry.Kanji {
		de.Expressions = append(de.Expressions, kanji.Expression)
	}

	for _, reading := range entry.Readings {
		de.Readings = append(de.Readings, reading.Reading)
	}

	for _, sense := range entry.Sense {
		var dd DictionaryDefinition

		dd.PartsOfSpeech = sense.PartsOfSpeech
		for _, gloss := range sense.Glossary {
			dd.Definitions = append(dd.Definitions, gloss.Content)
		}

		de.Definitions = append(de.Definitions, dd)
	}

	return de
}
