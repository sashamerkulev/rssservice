package data

import (
	"fmt"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

func AddComment(userId int64, articleId int64, comments string, logger logger.Logger) (commentsId int64, err error) {
	tx, err := DB.Begin()
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		return
	}
	defer tx.Commit()
	res, err := DB.Exec("INSERT INTO userArticleComments(userId, articleId, timestamp, comment, status) VALUES(?,?,?,?,?)", userId, articleId, time.Now(), comments, 0)
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", "userId="+fmt.Sprint(userId))
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		tx.Rollback()
		return 0, err
	}
	commentId, err := res.LastInsertId()
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		tx.Rollback()
		return 0, err
	}
	return commentId, nil
}

func GetComments(userId int64, articleId int64, lastArticleReadDate time.Time, logger logger.Logger) (comments []model.UserArticleComment, err error) {
	return make([]model.UserArticleComment, 0), nil
}

func DeleteComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		logger.Log("ERROR", "DELETECOMMENT", err.Error())
		return
	}
	defer tx.Commit()
	_, err = DB.Exec("DELETE FROM userArticleComments where commentId =?", commentId)
	if err != nil {
		logger.Log("ERROR", "DELETECOMMENT", err.Error())
		tx.Rollback()
		return err
	}
	return nil
}

func LikeComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return model.UserArticleComment{}, nil
}

func DislikeComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	return model.UserArticleComment{}, nil
}

func GetComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	rows, err := DB.Query(`
SELECT uac.CommentId, uac.ArticleId, uac.UserId, ui.UserName, uac.Comment, uac.Timestamp, uac.Status,
(SELECT COUNT(*) from userCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND NOT ucl.Dislike) as Likes,
(SELECT COUNT(*) from userCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND ucl.Dislike) as Dislikes,
(SELECT COUNT(*) from userCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND NOT ucl.Dislike AND ucl.UserId = ?) as userlike,
(SELECT COUNT(*) from userCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND ucl.Dislike AND ucl.UserId = ?) as userdislike,
uac.UserId = ? as Owner
FROM dbnews.userarticlecomments uac
JOIN userInfo ui on ui.userId = uac.userId
JOIN article a on a.articleId = uac.articleId
WHERE uac.CommentId = ?
;
`, userId, userId, userId, commentId)
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		return model.UserArticleComment{}, err
	}
	if rows.Next() {
		result := model.UserArticleComment{}
		err := rows.Scan(&result.CommentId, &result.ArticleId, &result.UserId, &result.Name, &result.Comment, &result.PubDate,
			&result.Status, &result.Likes, &result.Dislikes, &result.Like, &result.Dislike, &result.Owner)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return result, nil
	}
	return model.UserArticleComment{}, errors.CommentNotFoundError
}
