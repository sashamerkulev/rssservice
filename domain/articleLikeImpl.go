package domain

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
)

type ArticleUserLike struct {
	UserToken string
	ArticleId int64
	Logger    logger.Logger
}

func (aul ArticleUserLike) Like() error {
	return aul.likeOrDislike(db.LikeArticle)
}

func (aul ArticleUserLike) Dislike() error {
	return aul.likeOrDislike(db.DislikeArticle)
}

func (aul ArticleUserLike) likeOrDislike(action func(userId int64, articleId int64, logger logger.Logger) error) error {
	userId, err := db.GetUserIdByToken(aul.UserToken, aul.Logger)
	if err != nil {
		return err
	}
	return action(userId, aul.ArticleId, aul.Logger)
}
