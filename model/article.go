package model

import (
	"github.com/sashamerkulev/logger"
	"time"
)

type Article struct {
	ArticleId   int64
	SourceName  string
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	Category    string
	PictureUrl  string
}

type ConsumerArticleFunc func(articles []Article, logger logger.Logger)
