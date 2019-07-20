package main

import (
	"github.com/sashamerkulev/rssservice/controllers"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql"
	"github.com/sashamerkulev/rssservice/reader"
	"time"
)

var repository domain.MainRepository

func read(logger logger.Logger) {
	ticker := time.NewTicker(time.Minute * 15)
	for _ = range ticker.C {
		reader.Do(mysql.AddArticles, logger)
	}
}

func wipeArticles(logger logger.Logger) {
	ticker := time.NewTicker(time.Hour * 24)
	for _ = range ticker.C {
		wipeTime := time.Now()
		wipeTime = wipeTime.Add(-24 * 7 * time.Hour)
		mysql.WipeOldArticles(wipeTime, logger)
	}
}

func wipeActivities(logger logger.Logger) {
	ticker := time.NewTicker(time.Hour * 24)
	for _ = range ticker.C {
		wipeTime := time.Now()
		wipeTime = wipeTime.Add(-24 * 30 * time.Hour)
		mysql.WipeOldActivities(wipeTime, logger)
	}
}

func main() {
	repository = mysql.MainRepositoryImpl{}
	err := repository.Open()
	if err != nil {
		var log = logger.ConsoleLogger{}
		log.Log("ERROR", "MAIN", "Can't open DB. The service will be closed.")
		return
	}
	defer repository.Close()
	var log = repository.GetLogger(-1, "")
	go read(log)
	go wipeArticles(log)
	go wipeActivities(log)
	controllers.Init(repository)
}
