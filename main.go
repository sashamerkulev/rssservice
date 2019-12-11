package main

import (
	"github.com/sashamerkulev/rssservice/config"
	"github.com/sashamerkulev/rssservice/controllers"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql"
)

var repository domain.MainRepository

func main() {
	repository = mysql.MainRepositoryImpl{}
	cfg := config.GetConfig()
	err := repository.Open(cfg.Connection.Mysql)
	if err != nil {
		var log = logger.ConsoleLogger{}
		log.Log("ERROR", "MAIN", "Can't open DB. The service will be closed.")
		return
	}
	defer repository.Close()
	controllers.Init(repository)
}
