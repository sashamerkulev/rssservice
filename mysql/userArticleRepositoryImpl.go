package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql/data"
)

type UserArticleRepositoryImpl struct {
}

func (UserArticleRepositoryImpl) LikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	return data.LikeArticle(userId, articleId, logger)
}

func (UserArticleRepositoryImpl) DislikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	return data.DislikeArticle(userId, articleId, logger)
}

func (UserArticleRepositoryImpl) FindUserArticleDislike(userId int64, articleId int64, logger logger.Logger) (bool, error) {
	return data.FindUserArticleDislike(userId, articleId, logger)
}

func (UserArticleRepositoryImpl) SetUserArticleDislikeTo(userId int64, articleId int64, dislike bool, logger logger.Logger) error {
	return data.SetUserArticleDislikeTo(userId, articleId, dislike, logger)
}

func (UserArticleRepositoryImpl) RemoveUserArticleDislike(userId int64, articleId int64, logger logger.Logger) error {
	return data.RemoveUserArticleDislike(userId, articleId, logger)
}
