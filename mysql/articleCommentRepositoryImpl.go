package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"sort"
	"time"
)

type ArticleCommentRepositoryImpl struct {
	DB *sql.DB
}

func (db ArticleCommentRepositoryImpl) GetComments(userId int64, articleId int64, lastArticleReadDate time.Time, logger logger.Logger) (comments []model.UserArticleComment, err error) {
	comments = make([]model.UserArticleComment, 0)
	rows, err := db.DB.Query(`SELECT * FROM (
SELECT uac.CommentId, uac.ArticleId, uac.UserId, 
CASE WHEN LENGTH(ui.UserName) = 0 OR ui.UserName IS NULL THEN CONCAT('гость_', CONVERT(ui.UserId, char)) ELSE ui.UserName END AS UserName, 
uac.Comment, uac.Timestamp, uac.Status,
(select max(ucl.timestamp) from articleCommentLikes ucl where ucl.CommentId = uac.CommentId ) as lastActivityDate, 
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND NOT ucl.Dislike) as Likes,
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND ucl.Dislike) as Dislikes,
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND NOT ucl.Dislike AND ucl.UserId = ?) as userlike,
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND ucl.Dislike AND ucl.UserId = ?) as userdislike,
uac.UserId = ? as Owner
FROM articleComments uac
JOIN users ui on ui.userId = uac.userId
JOIN articles a on a.articleId = uac.articleId) b
WHERE b.articleId = ? AND (b.Timestamp >= ? or (b.lastActivityDate >= ?))
;
`, userId, userId, userId, articleId, lastArticleReadDate, lastArticleReadDate)
	if err != nil {
		logger.Log("ERROR", "GETCOMMENTS", err.Error())
		return comments, err
	}
	for rows.Next() {
		var lastActivityDate mysql.NullTime
		result := model.UserArticleComment{}
		err := rows.Scan(&result.CommentId, &result.ArticleId, &result.UserId, &result.Name, &result.Comment, &result.PubDate,
			&result.Status, &lastActivityDate, &result.Likes, &result.Dislikes, &result.Like, &result.Dislike, &result.Owner)
		if lastActivityDate.Valid {
			result.LastActivityDate = lastActivityDate.Time
		} else {
			result.LastActivityDate = result.PubDate
		}
		if err != nil {
			logger.Log("ERROR", "GETCOMMENTS", err.Error())
		}
		comments = append(comments, result)
	}
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].PubDate.After(comments[j].PubDate)
	})
	return comments, nil
}

func (db ArticleCommentRepositoryImpl) AddComment(userId int64, articleId int64, comments string, logger logger.Logger) (commentId int64, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		return
	}
	defer tx.Commit()
	res, err := db.DB.Exec("INSERT INTO articleComments(userId, articleId, timestamp, comment, status) VALUES(?,?,?,?,?)", userId, articleId, time.Now(), comments, 0)
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", "userId="+fmt.Sprint(userId))
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		tx.Rollback()
		return 0, err
	}
	commentId, err = res.LastInsertId()
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		tx.Rollback()
		return 0, err
	}
	return commentId, nil

}

func (db ArticleCommentRepositoryImpl) DeleteComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		logger.Log("ERROR", "DELETECOMMENT", err.Error())
		return
	}
	defer tx.Commit()
	_, err = db.DB.Exec("DELETE FROM articleComments where commentId =?", commentId)
	if err != nil {
		logger.Log("ERROR", "DELETECOMMENT", err.Error())
		tx.Rollback()
		return err
	}
	return nil

}

func (db ArticleCommentRepositoryImpl) GetComment(userId int64, commentId int64, logger logger.Logger) (comment model.UserArticleComment, err error) {
	rows, err := db.DB.Query(`
SELECT uac.CommentId, uac.ArticleId, uac.UserId, 
CASE WHEN LENGTH(ui.UserName) = 0 OR ui.UserName IS NULL THEN CONCAT('гость_', CONVERT(ui.UserId, char)) ELSE ui.UserName END AS UserName, 
uac.Comment, uac.Timestamp, uac.Status,
(select max(ucl.timestamp) from articleCommentLikes ucl where ucl.CommentId = uac.CommentId ) as lastActivityDate, 
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND NOT ucl.Dislike) as Likes,
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND ucl.Dislike) as Dislikes,
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND NOT ucl.Dislike AND ucl.UserId = ?) as userlike,
(SELECT COUNT(*) from articleCommentLikes ucl WHERE ucl.CommentId = uac.CommentId AND ucl.Dislike AND ucl.UserId = ?) as userdislike,
uac.UserId = ? as Owner
FROM articleComments uac
JOIN users ui on ui.userId = uac.userId
JOIN articles a on a.articleId = uac.articleId
WHERE uac.CommentId = ?
;
`, userId, userId, userId, commentId)
	if err != nil {
		logger.Log("ERROR", "ADDCOMMENT", err.Error())
		return model.UserArticleComment{}, err
	}
	if rows.Next() {
		var lastActivityDate mysql.NullTime
		result := model.UserArticleComment{}
		err := rows.Scan(&result.CommentId, &result.ArticleId, &result.UserId, &result.Name, &result.Comment, &result.PubDate,
			&result.Status, &lastActivityDate, &result.Likes, &result.Dislikes, &result.Like, &result.Dislike, &result.Owner)
		if lastActivityDate.Valid {
			result.LastActivityDate = lastActivityDate.Time
		} else {
			result.LastActivityDate = result.PubDate
		}
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return result, nil
	}
	return model.UserArticleComment{}, errors.CommentNotFoundError
}

func (db ArticleCommentRepositoryImpl) FindCommentDislike(userId int64, commentId int64, logger logger.Logger) (bool, error) {
	rows, err := db.DB.Query("select dislike from articleCommentLikes WHERE userId = ? and commentId = ?", userId, commentId)
	if err != nil {
		logger.Log("ERROR", "FINDCOMMENTDISLIKE", err.Error())
		return false, errors.CommentNotFoundError
	}
	defer rows.Close()
	if rows.Next() {
		var dislike bool
		err = rows.Scan(&dislike)
		if err != nil {
			return false, errors.CommentNotFoundError
		}
		return dislike, nil
	}
	return false, errors.CommentNotFoundError
}

func (db ArticleCommentRepositoryImpl) SetUserCommentDislikeTo(userId int64, commentId int64, dislike bool, logger logger.Logger) error {
	_, err := db.DB.Exec("update articleCommentLikes set dislike = ?, timestamp = ? where userId=? and commentId = ?", dislike, time.Now(), userId, commentId)
	return err
}

func (db ArticleCommentRepositoryImpl) RemoveCommentDislike(userId int64, commentId int64, logger logger.Logger) error {
	_, err := db.DB.Exec("delete from articleCommentLikes where userId=? and commentId = ?", userId, commentId)
	return err
}

func (db ArticleCommentRepositoryImpl) LikeComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	_, err = db.DB.Exec("insert into articleCommentLikes (userId, commentId, dislike, timestamp) values(?,?,?,?)", userId, commentId, false, time.Now())
	return err
}

func (db ArticleCommentRepositoryImpl) DislikeComment(userId int64, commentId int64, logger logger.Logger) (err error) {
	_, err = db.DB.Exec("insert into articleCommentLikes (userId, commentId, dislike, timestamp) values(?,?,?,?)", userId, commentId, true, time.Now())
	return err
}
