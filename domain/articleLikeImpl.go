package domain

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
)

type ArticleUserLike struct {
	UserId    int64
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
	return action(aul.UserId, aul.ArticleId, aul.Logger)
}

func (aul ArticleUserLike) Comment(comments string) error {
	return nil
}
