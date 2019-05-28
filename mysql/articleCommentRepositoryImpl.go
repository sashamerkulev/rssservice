package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
)

type ArticleCommentRepositoryImpl struct {
}

func (ArticleCommentRepositoryImpl) GetComments(userId int64, articleId int64, logger logger.Logger) (comments []model.UserArticleComment, err error) {
	return make([]model.UserArticleComment, 0), nil
}

func (ArticleCommentRepositoryImpl) AddComment(userId int64, articleId int64, comments string, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return model.UserArticleComment{}, nil
}

func (ArticleCommentRepositoryImpl) DeleteComment(userId int64, articleId int64, commentId int64, logger logger.Logger) (err error) {
	return nil
}

func (ArticleCommentRepositoryImpl) LikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return model.UserArticleComment{}, nil
}

func (ArticleCommentRepositoryImpl) DislikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return model.UserArticleComment{}, nil
}
