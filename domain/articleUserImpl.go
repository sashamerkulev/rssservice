package domain

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

type ArticleUser struct {
	UserToken string
	LastTime  time.Time
	Logger    logger.Logger
}

func (au ArticleUser) GetArticleUser() ([]model.ArticleUser, error) {
	userId, err := db.GetUserIdByToken(au.UserToken, au.Logger)
	if err != nil {
		return make([]model.ArticleUser, 0), err
	}
	return db.GetArticleUser(userId, au.LastTime, au.Logger), nil
}
