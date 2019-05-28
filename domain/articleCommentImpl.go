package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
)

type ArticleCommentRepository interface {
	GetComments(userId int64, articleId int64, logger logger.Logger) (comments []model.UserArticleComment, err error)
	AddComment(userId int64, articleId int64, comments string, logger logger.Logger)  (comment model.UserArticleComment, err error)
	DeleteComment(userId int64, articleId int64, commentId int64, logger logger.Logger)  (err error)
	LikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger)  (comment model.UserArticleComment, err error)
	DislikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger)  (comment model.UserArticleComment, err error)
}

type ArticleComment struct {
	UserId     int64
	ArticleId  int64
	CommentId  int64
	Logger     logger.Logger
	Repository ArticleCommentRepository
}

func (ac ArticleComment) GetComments() (comments []model.UserArticleComment, err error) {
	return ac.Repository.GetComments(ac.UserId, ac.ArticleId, ac.Logger)
}

func (ac ArticleComment) AddComment(comments string) (comment model.UserArticleComment, err error) {
	return ac.Repository.AddComment(ac.UserId, ac.ArticleId, comments, ac.Logger)
}

func (ac ArticleComment) DeleteComment() (err error) {
	return ac.Repository.DeleteComment(ac.UserId, ac.ArticleId, ac.CommentId, ac.Logger)
}

func (ac ArticleComment) LikeComment() (comment model.UserArticleComment, err error) {
	return ac.Repository.LikeComment(ac.UserId, ac.ArticleId, ac.CommentId, ac.Logger)
}

func (ac ArticleComment) DislikeComment() (comment model.UserArticleComment, err error) {
	return ac.Repository.DislikeComment(ac.UserId, ac.ArticleId, ac.CommentId, ac.Logger)
}
