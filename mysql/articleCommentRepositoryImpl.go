package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
)

type ArticleCommentRepositoryImpl struct {
}

func (ArticleCommentRepositoryImpl) AddComment(userId int64, articleId int64, comments string, logger logger.Logger) error {
	return nil
}

func (ArticleCommentRepositoryImpl) LikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger) error {
	return nil
}

func (ArticleCommentRepositoryImpl) DislikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger) error {
	return nil
}
