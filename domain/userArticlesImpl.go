package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

type UserArticlesRepository interface {
	GetUserArticles(userId int64, lastTime time.Time, logger logger.Logger) ([]model.ArticleUser, error)
}

type UserArticles struct {
	UserId     int64
	LastTime   time.Time
	Logger     logger.Logger
	Repository UserArticlesRepository
}

func (ua UserArticles) GetUserArticles() ([]model.ArticleUser, error) {
	return ua.Repository.GetUserArticles(ua.UserId, ua.LastTime, ua.Logger)
}
