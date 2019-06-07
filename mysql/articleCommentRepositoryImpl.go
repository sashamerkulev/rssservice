package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/mysql/data"
)

type ArticleCommentRepositoryImpl struct {
}

func (ArticleCommentRepositoryImpl) GetComments(userId int64, articleId int64, logger logger.Logger) (comments []model.UserArticleComment, err error) {
	return data.GetComments(userId, articleId, logger)
}

func (ArticleCommentRepositoryImpl) AddComment(userId int64, articleId int64, comments string, logger logger.Logger) (commentId int64, err error) {
	return data.AddComment(userId, articleId, comments, logger)
}

func (ArticleCommentRepositoryImpl) DeleteComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	return data.DeleteComment(userId, commentId, logger)
}

func (ArticleCommentRepositoryImpl) LikeComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return data.LikeComment(userId, commentId, logger)
}

func (ArticleCommentRepositoryImpl) DislikeComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return data.DislikeComment(userId, commentId, logger)
}

func (ArticleCommentRepositoryImpl) GetComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return data.GetComment(userId, commentId, logger)
}
