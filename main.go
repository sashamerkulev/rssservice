package main

import (
	"fmt"
	"github.com/sashamerkulev/rssservice/config"
	"github.com/sashamerkulev/rssservice/controllers"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/models"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
	"time"
)

var repository domain.MainRepository

var Urls = []models.Link{
	{Link: "http://lenta.ru/rss", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "lenta", Title: "Лента"},
	//{Link: "http://static.feed.rbc.ru/rbc/internal/rss.rbc.ru/rbc.ru/mainnews.rss", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "rbc", Title: "РБК"},
	{Link: "http://worldoftanks.ru/ru/rss/news/", Layout: "Mon, _2 Jan 2006 15:04:05 MST", Name: "wot", Title: "World of Tanks"},
	{Link: "https://topwar.ru/rss.xml", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "topwar", Title: "Topwar"},
	{Link: "http://www.interfax.ru/rss.asp", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "interfax", Title: "Интерфакс"}, // Fri, 8 Mar 2019 13:58:00 +0300
	{Link: "http://www.vesti.ru/vesti.rss", Layout: time.RFC1123Z, Name: "vesti", Title: "Вести"},
	{Link: "http://russian.rt.com/rss/", Layout: time.RFC1123Z, Name: "rt", Title: "RT"},
	{Link: "http://www.planetanovosti.com/news/rss/", Layout: time.RFC1123, Name: "planetanovosti", Title: "Планета Новости"},
	{Link: "https://news.rambler.ru/rss/world/", Layout: time.RFC1123Z, Name: "rambler", Title: "Rambler.news"},
	{Link: "http://rss.newsru.com/world/", Layout: time.RFC1123Z, Name: "newsru", Title: "News.ru"},
	{Link: "http://mixednews.ru/feed/", Layout: time.RFC1123Z, Name: "mixednews", Title: "MixedNews"},
	{Link: "http://rg.ru/xml/index.xml", Layout: "_2 Jan 2006 15:04:05 MST", Name: "rg", Title: "Российская газета"}, //8 Mar 2019 12:40:17 GMT
	{Link: "http://www.ng.ru/rss/", Layout: time.RFC1123Z, Name: "ng", Title: "Независимая газета"},
	{Link: "http://www.kp.ru/rss/allsections.xml", Layout: time.RFC1123Z, Name: "kp", Title: "Комсомольская правда"},
	{Link: "http://www.km.ru/rss/main", Layout: time.RFC1123Z, Name: "km", Title: "Кирил и Мефодий"},
	{Link: "http://feeds.feedburner.com/aftershock/news", Layout: time.RFC1123Z, Name: "aftershock", Title: "Aftershock"},
	//{Link: "http://otredakcii.odnako.org/rss/", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "odnako", Title: "Однако"},
	{Link: "http://www.aif.ru/rss/all.php", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "aif", Title: "Аргументы и факты"},
	{Link: "http://feeds.bbci.co.uk/russian/rss.xml", Layout: "Mon, _2 Jan 2006 15:04:05 MST", Name: "bbcrussian", Title: "BBC Russian"},
	{Link: "http://tass.ru/rss/v2.xml", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "tass", Title: "ТАСС"},
	{Link: "http://www.nkj.ru/rss/", Layout: time.RFC1123Z, Name: "nkj", Title: "Наука и Жизнь"},
	//{Link: "http://www.mk.ru/rss/index.xml", Layout: time.RFC1123Z, Name: "mk", Title: "Московский комсомолец"},
	{Link: "http://www.cnews.ru/inc/rss/news.xml", Layout: time.RFC1123Z, Name: "cnews", Title: "CNews"},
	{Link: "https://news.mail.ru/rss/91/", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "mailnews", Title: "Mail.news"}, //Fri, 8 Mar 2019 15:49:34 +0300
	{Link: "https://www.sport-express.ru/services/materials/news/se/", Layout: time.RFC1123Z, Name: "sportexpress", Title: "Спорт-Экспресс"},
	{Link: "http://dp.ru/rss", Layout: time.RFC1123Z, Name: "dp", Title: "Деловой Петербург"},
	{Link: "https://sdelanounas.ru/index/rss/", Layout: time.RFC1123Z, Name: "sdelanounas", Title: "Сделано у нас"},
	//{Link:"http://www.rosbalt.ru/feed/", Layout: time.RFC1123Z, Name: "rosbalt", Title: "Росбалт"},
	//{Link:"http://feeds.feedburner.com/iarex/?utm_source=social-icons&utm_medium=cpc&utm_campaign=web", Layout: time.RFC1123Z, Name: "iarex", Title: "ИА REX"},
}

func main() {
	repository = mysql.MainRepositoryImpl{Sources: Urls}
	cfg := config.GetConfig()
	fmt.Print(cfg.Connection.Mysql)
	fmt.Print(cfg.Server.Port)
	err := repository.Open(cfg.Connection.Mysql)
	if err != nil {
		var log = logger.ConsoleLogger{}
		log.Log("ERROR", "MAIN", "Can't open DB. The service will be closed.")
		return
	}
	defer repository.Close()
	var r = controllers.Init(repository)
	panic(http.ListenAndServe(cfg.Server.Port, r))
}
