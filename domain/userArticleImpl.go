package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
)

type UserArticleRepository interface {
	LikeArticle(userId int64, articleId int64, logger logger.Logger) error
	DislikeArticle(userId int64, articleId int64, logger logger.Logger) error
	FindUserArticleDislike(userId int64, articleId int64, logger logger.Logger) (bool, error)
	SetUserArticleDislikeTo(userId int64, articleId int64, dislike bool, logger logger.Logger) error
	RemoveUserArticleDislike(userId int64, articleId int64, logger logger.Logger) error
	GetUserArticle(userId int64, articleId int64, logger logger.Logger) (model.ArticleUser, error)
}

type UserArticle struct {
	UserId     int64
	ArticleId  int64
	Logger     logger.Logger
	Repository UserArticleRepository
}

func (ua UserArticle) Like() (model.ArticleUser, error) {
	// if no likes or dislike - add like
	// if dislike - change to like
	// if like - remove like
	dislike, err := ua.Repository.FindUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
	if err == nil {
		if dislike {
			err = ua.Repository.SetUserArticleDislikeTo(ua.UserId, ua.ArticleId, false, ua.Logger)
		} else {
			err = ua.Repository.RemoveUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
		}
	} else {
		err = ua.Repository.LikeArticle(ua.UserId, ua.ArticleId, ua.Logger)
	}
	if err != nil {
		return model.ArticleUser{}, err
	}
	return ua.Repository.GetUserArticle(ua.UserId, ua.ArticleId, ua.Logger)
}

func (ua UserArticle) Dislike() (model.ArticleUser, error) {
	// if no likes or dislike - add dislike
	// if dislike - remove dislike
	// if like - change to dislike
	dislike, err := ua.Repository.FindUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
	if err == nil {
		if dislike {
			err = ua.Repository.RemoveUserArticleDislike(ua.UserId, ua.ArticleId, ua.Logger)
		} else {
			err = ua.Repository.SetUserArticleDislikeTo(ua.UserId, ua.ArticleId, true, ua.Logger)
		}
	} else {
		err = ua.Repository.DislikeArticle(ua.UserId, ua.ArticleId, ua.Logger)
	}
	if err != nil {
		return model.ArticleUser{}, err
	}
	return ua.Repository.GetUserArticle(ua.UserId, ua.ArticleId, ua.Logger)
}

func (ua UserArticle) GetUserArticle() (model.ArticleUser, error) {
	return ua.Repository.GetUserArticle(ua.UserId, ua.ArticleId, ua.Logger)
}
