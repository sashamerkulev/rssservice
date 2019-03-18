package domain

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

type ArticleUser struct {
	UserId   int64
	LastTime time.Time
	Logger   logger.Logger
}

func (au ArticleUser) GetArticleUser() ([]model.ArticleUser, error) {
	return db.GetArticleUser(au.UserId, au.LastTime, au.Logger)
}
