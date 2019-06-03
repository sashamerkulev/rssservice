package reader

import (
	"encoding/xml"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"io/ioutil"
	"net/http"
	"time"
)

type Link struct {
	Link   string
	Layout string
	Name   string
	Title  string
}

var Urls = []Link{
	{Link: "http://lenta.ru/rss", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "lenta", Title: "Лента"},
	{Link: "http://static.feed.rbc.ru/rbc/internal/rss.rbc.ru/rbc.ru/mainnews.rss", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "rbc", Title: "РБК"},
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
	{Link: "http://otredakcii.odnako.org/rss/", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "odnako", Title: "Однако"},
	{Link: "http://www.aif.ru/rss/all.php", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "aif", Title: "Аргументы и факты"},
	{Link: "http://feeds.bbci.co.uk/russian/rss.xml", Layout: "Mon, _2 Jan 2006 15:04:05 MST", Name: "bbcrussian", Title: "BBC Russian"},
	{Link: "http://tass.ru/rss/v2.xml", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "tass", Title: "ТАСС"},
	{Link: "http://www.nkj.ru/rss/", Layout: time.RFC1123Z, Name: "nkj", Title: "Наука и Жизнь"},
	//{Link: "http://www.mk.ru/rss/index.xml", Layout: time.RFC1123Z, Name: "mk", Title: "Московский комсомолец"},
	{Link: "http://www.cnews.ru/inc/rss/news.xml", Layout: time.RFC1123Z, Name: "cnews", Title: "CNews"},
	{Link: "https://news.mail.ru/rss/91/", Layout: "Mon, _2 Jan 2006 15:04:05 -0700", Name: "mailnews", Title: "Mail.news"}, //Fri, 8 Mar 2019 15:49:34 +0300
	{Link: "http://www.sport-express.ru/controllers/materials/news/se/", Layout: time.RFC1123Z, Name: "sportexpress", Title: "Спорт-Экспресс"},
	{Link: "http://dp.ru/rss", Layout: time.RFC1123Z, Name: "dp", Title: "Деловой Петербург"},
	{Link: "https://sdelanounas.ru/index/rss/", Layout: time.RFC1123Z, Name: "sdelanounas", Title: "Сделано у нас"},
}

type rss2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	// Optional
	PubDate  string `xml:"channel>pubDate"`
	ItemList []item `xml:"channel>item"`
}

type enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type item struct {
	// Required
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	// Optional
	Content   string    `xml:"encoded"`
	PubDate   string    `xml:"pubDate"`
	Comments  string    `xml:"comments"`
	Enclosure enclosure `xml:"enclosure"`
	Category  string    `xml:"category"`
	Guid      string    `xml:"guid"`
}

func read(url string, bytes chan []byte, logger logger.Logger) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Log("ERROR", "READ ("+url+")", err.Error())
		bytes <- make([]byte, 0)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log("ERROR", "READ ("+url+")", err.Error())
		bytes <- make([]byte, 0)
		return
	}
	bytes <- body
}

func parse(name string, bytes []byte, items chan []item, logger logger.Logger) {
	if len(bytes) <= 0 {
		items <- make([]item, 0)
		return
	}
	r := rss2{}
	err := xml.Unmarshal(bytes, &r)
	if err != nil {
		logger.Log("ERROR", "PARSE ("+name+")", err.Error())
		items <- make([]item, 0)
		return
	}
	items <- r.ItemList
}

func save(saver model.ConsumerArticleFunc, id int, items []item, logger logger.Logger) {
	if len(items) <= 0 {
		return
	}
	var articles = make([]model.Article, 0)
	for i := 0; i < len(items); i++ {
		date, err := time.Parse(Urls[id].Layout, items[i].PubDate)
		if err == nil {
			article := model.Article{
				SourceName:  Urls[id].Name,
				Category:    items[i].Category,
				Description: string(items[i].Description),
				Link:        items[i].Link,
				PubDate:     date,
				Title:       items[i].Title,
				PictureUrl:  items[i].Enclosure.Url,
			}
			articles = append(articles, article)
		} else {
			logger.Log("ERROR", "SAVE ("+Urls[i].Name+")", err.Error())
		}
	}
	saver(articles, logger)
}

func Do(consumer model.ConsumerArticleFunc, logger logger.Logger) {
	for i := 0; i < len(Urls); i++ {
		bytes := make(chan []byte)
		items := make(chan []item)
		go read(Urls[i].Link, bytes, logger)
		go parse(Urls[i].Name, <-bytes, items, logger)
		go save(consumer, i, <-items, logger)
	}
}
