package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
)

type UserArticleRepository interface {
	LikeArticle(userId int64, articleId int64, logger logger.Logger) error
	DislikeArticle(userId int64, articleId int64, logger logger.Logger) error
	FindUserArticleDislike(userId int64, articleId int64, logger logger.Logger) (bool, error)
	SetUserArticleDislikeTo(userId int64, articleId int64, dislike bool, logger logger.Logger) error
	RemoveUserArticleDislike(userId int64, articleId int64, logger logger.Logger) error
}

type UserArticle struct {
	UserId     int64
	ArticleId  int64
	Logger     logger.Logger
	Repository UserArticleRepository
}

func (ua UserArticle) Like() error {
	// if no likes or dislike - add like
	// if dislike - change to like
	// if like - remove like
	dislike, err := ua.Repository.FindUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
	if err == nil {
		if dislike {
			return ua.Repository.SetUserArticleDislikeTo(ua.UserId, ua.ArticleId, false, ua.Logger)
		} else {
			return ua.Repository.RemoveUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
		}
	} else {
		return ua.Repository.LikeArticle(ua.UserId, ua.ArticleId, ua.Logger)
	}
}

func (ua UserArticle) Dislike() error {
	// if no likes or dislike - add dislike
	// if dislike - remove dislike
	// if like - change to dislike
	dislike, err := ua.Repository.FindUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
	if err == nil {
		if dislike {
			return ua.Repository.RemoveUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
		} else {
			return ua.Repository.SetUserArticleDislikeTo(ua.UserId, ua.ArticleId, true, ua.Logger)
		}
	} else {
		return ua.Repository.DislikeArticle(ua.UserId, ua.ArticleId, ua.Logger)
	}
}
