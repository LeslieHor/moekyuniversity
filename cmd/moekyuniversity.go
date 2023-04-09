package main

import (
	"flag"
	"log"

	"moekyuniversity/internal/cards"
)

var (
	cardsFile = flag.String("cards-file", "data/cards.json", "Cards file")
	dataDir   = flag.String("data-dir", "data", "Data directory")
	backupDir = flag.String("backup-dir", "data/backup", "Backup directory")
	staticDir = flag.String("static-dir", "static", "Static directory")
)

func main() {
	// Read flags and log
	flag.Parse()
	log.Printf("Cards file: %s", *cardsFile)
	log.Printf("Data directory: %s", *dataDir)
	log.Printf("Backup directory: %s", *backupDir)
	log.Printf("Static directory: %s", *staticDir)

	cardData := cards.CardData{
		CardsFile: *cardsFile,
		DataDir: *dataDir,
		BackupDir: *backupDir,
		StaticDir: *staticDir,
	}
	cardData.LoadCardJson()

	cards.SetupRoutes(&cardData)
}