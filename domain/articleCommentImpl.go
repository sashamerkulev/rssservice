package domain

import "github.com/sashamerkulev/rssservice/logger"

type ArticleCommentRepository interface {
	AddComment(userId int64, articleId int64, comments string, logger logger.Logger) error
	LikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger) error
	DislikeComment(userId int64, articleId int64, commentId int64, logger logger.Logger) error
}

type ArticleComment struct {
	UserId     int64
	ArticleId  int64
	CommentId  int64
	Logger     logger.Logger
	Repository ArticleCommentRepository
}

func (ac ArticleComment) AddComment(comments string) error {
	return ac.Repository.AddComment(ac.UserId, ac.ArticleId, comments, ac.Logger)
}

func (ac ArticleComment) LikeComment() error {
	return ac.Repository.LikeComment(ac.UserId, ac.ArticleId, ac.CommentId, ac.Logger)
}

func (ac ArticleComment) DislikeComment() error {
	return ac.Repository.DislikeComment(ac.UserId, ac.ArticleId, ac.CommentId, ac.Logger)
}
