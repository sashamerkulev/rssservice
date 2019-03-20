package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/mysql/data"
	"time"
)

type UserArticlesRepositoryImpl struct {
}

func (UserArticlesRepositoryImpl) GetUserArticles(userId int64, lastTime time.Time, logger logger.Logger) ([]model.ArticleUser, error) {
	return data.GetUserArticles(userId, lastTime, logger)
}

func (UserArticlesRepositoryImpl) GetUserFavoriteArticles(userId int64, logger logger.Logger) ([]model.ArticleUser, error) {
	return make([]model.ArticleUser, 0), nil
}
