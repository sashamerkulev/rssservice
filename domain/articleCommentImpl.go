package domain

import (
	"github.com/sashamerkulev/rssservice/fcm"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

type ArticleCommentRepository interface {
	GetComments(userId int64, articleId int64, lastArticleReadDate time.Time, logger logger.Logger) (comments []model.UserArticleComment, err error)
	AddComment(userId int64, articleId int64, comments string, logger logger.Logger) (commentId int64, err error)
	DeleteComment(userId int64, commentId int64, logger logger.Logger) (err error)
	GetComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error)
	FindCommentDislike(userId int64, commentId int64, logger logger.Logger) (bool, error)
	SetUserCommentDislikeTo(userId int64, commentId int64, dislike bool, logger logger.Logger) error
	RemoveCommentDislike(userId int64, commentId int64, logger logger.Logger) error
	LikeComment(userId int64, commentId int64, logger logger.Logger) (err error)
	DislikeComment(userId int64, commentId int64, logger logger.Logger) (err error)
}

type ArticleComment struct {
	UserId              int64
	ArticleId           int64
	CommentId           int64
	LastArticleReadDate time.Time
	Logger              logger.Logger
	Repository          ArticleCommentRepository
}

func (ac ArticleComment) GetComments() (comments []model.UserArticleComment, err error) {
	return ac.Repository.GetComments(ac.UserId, ac.ArticleId, ac.LastArticleReadDate, ac.Logger)
}

func (ac ArticleComment) AddComment(comments string) (comment model.UserArticleComment, err error) {
	commentId, err := ac.Repository.AddComment(ac.UserId, ac.ArticleId, comments, ac.Logger)
	return ac.Repository.GetComment(ac.UserId, commentId, ac.Logger)
}

func (ac ArticleComment) DeleteComment() (err error) {
	return ac.Repository.DeleteComment(ac.UserId, ac.CommentId, ac.Logger)
}

func (ac ArticleComment) LikeComment() (comment model.UserArticleComment, err error) {
	// if no likes or dislike - add like
	// if dislike - change to like
	// if like - remove like
	dislike, err := ac.Repository.FindCommentDislike(ac.UserId, ac.CommentId, ac.Logger)
	if err == nil {
		if dislike {
			err = ac.Repository.SetUserCommentDislikeTo(ac.UserId, ac.CommentId, false, ac.Logger)
		} else {
			err = ac.Repository.RemoveCommentDislike(ac.UserId, ac.CommentId, ac.Logger)
		}
	} else {
		err = ac.Repository.LikeComment(ac.UserId, ac.CommentId, ac.Logger)
	}
	if err != nil {
		return model.UserArticleComment{}, err
	}
	c, e := ac.Repository.GetComment(ac.UserId, ac.CommentId, ac.Logger)
	if e == nil {
		fcm.SendNotificationLikeArticleComment(c.Name, ac.CommentId, false, ac.Logger)
	}
	return c, e
}

func (ac ArticleComment) DislikeComment() (comment model.UserArticleComment, err error) {
	// if no likes or dislike - add dislike
	// if dislike - remove dislike
	// if like - change to dislike
	dislike, err := ac.Repository.FindCommentDislike(ac.UserId, ac.CommentId, ac.Logger)
	if err == nil {
		if dislike {
			err = ac.Repository.RemoveCommentDislike(ac.UserId, ac.CommentId, ac.Logger)
		} else {
			err = ac.Repository.SetUserCommentDislikeTo(ac.UserId, ac.CommentId, true, ac.Logger)
		}
	} else {
		err = ac.Repository.DislikeComment(ac.UserId, ac.CommentId, ac.Logger)
	}
	if err != nil {
		return model.UserArticleComment{}, err
	}
	c, e := ac.Repository.GetComment(ac.UserId, ac.CommentId, ac.Logger)
	if e == nil {
		fcm.SendNotificationLikeArticleComment(c.Name, ac.CommentId, true, ac.Logger)
	}
	return c, e
}
