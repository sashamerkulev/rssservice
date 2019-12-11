package models

import (
	"github.com/sashamerkulev/rssservice/logger"
	"time"
)

type Article struct {
	ArticleId        int64
	SourceName       string
	Title            string
	Link             string
	Description      string
	PubDate          time.Time
	LastActivityDate time.Time
	Category         string
	PictureUrl       string
}

type ConsumerArticleFunc func(articles []Article, logger logger.Logger)
