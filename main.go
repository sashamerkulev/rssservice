package main

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/controllers"
	"github.com/sashamerkulev/rssservice/db"
	"github.com/sashamerkulev/rssservice/reader"
	"time"
)

func read(dbLogger logger.DbLogger) {
	ticker := time.NewTicker(time.Minute * 15)
	for _ = range ticker.C {
		reader.Do(db.AddArticles, dbLogger)
	}
}

func wipe(dbLogger logger.DbLogger) {
	ticker := time.NewTicker(time.Hour * 12)
	for _ = range ticker.C {
		wipeTime := time.Now()
		wipeTime = wipeTime.Add(-12 * time.Hour)
		db.WipeOldArticles(wipeTime, dbLogger)
	}
}

func main() {
	err := db.DBOpen()
	if err != nil {
		var consoleLogger = logger.ConsoleLogger{}
		consoleLogger.Log("ERROR", "MAIN", "Can't open DB. The service will be closed.")
		return
	}
	dbLogger := logger.DbLogger{DB: db.DB}
	go read(dbLogger)
	go wipe(dbLogger)
	controllers.Init(dbLogger)
	defer db.DBClose()
}
