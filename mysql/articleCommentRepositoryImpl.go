package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/mysql/data"
	"time"
)

type ArticleCommentRepositoryImpl struct {
}

func (ArticleCommentRepositoryImpl) GetComments(userId int64, articleId int64, lastArticleReadDate time.Time, logger logger.Logger) (comments []model.UserArticleComment, err error) {
	list, err := data.GetComments(userId, articleId, lastArticleReadDate, logger)
	for i:= 0; i < len(list); i++ {
		list[i] = model.MakeSureCommenterExists(list[i])
	}
	return list, err
}

func (ArticleCommentRepositoryImpl) AddComment(userId int64, articleId int64, comments string, logger logger.Logger) (commentId int64, err error) {
	return data.AddComment(userId, articleId, comments, logger)
}

func (ArticleCommentRepositoryImpl) DeleteComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	return data.DeleteComment(userId, commentId, logger)
}

func (ArticleCommentRepositoryImpl) GetComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	uac, err := data.GetComment(userId, commentId, logger)
	return model.MakeSureCommenterExists(uac), err
}

func (ArticleCommentRepositoryImpl) FindCommentDislike(userId int64, commentId int64, logger logger.Logger) (bool, error) {
	return data.FindCommentDislike(userId, commentId, logger)
}

func (ArticleCommentRepositoryImpl) SetUserCommentDislikeTo(userId int64, commentId int64, dislike bool, logger logger.Logger) error {
	return data.SetUserCommentDislikeTo(userId, commentId, dislike, logger)
}

func (ArticleCommentRepositoryImpl) RemoveCommentDislike(userId int64, commentId int64, logger logger.Logger) error {
	return data.RemoveCommentDislike(userId, commentId, logger)
}

func (ArticleCommentRepositoryImpl) LikeComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	return data.LikeComment(userId, commentId, logger)
}

func (ArticleCommentRepositoryImpl) DislikeComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	return data.DislikeComment(userId, commentId, logger)
}
