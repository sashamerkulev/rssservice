package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

type UserArticleRepositoryImpl struct {
	DB *sql.DB
}

func (db UserArticleRepositoryImpl) LikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	_, err := db.DB.Exec("insert into articleLikes (userId, articleId, dislike, timestamp) values(?,?,?,?)", userId, articleId, false, time.Now())
	return err
}

func (db UserArticleRepositoryImpl) DislikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	_, err := db.DB.Exec("insert into articleLikes (userId, articleId, dislike, timestamp) values(?,?,?,?)", userId, articleId, true, time.Now())
	return err
}

func (db UserArticleRepositoryImpl) FindUserArticleDislike(userId int64, articleId int64, logger logger.Logger) (bool, error) {
	rows, err := db.DB.Query("select dislike from articleLikes WHERE userId = ? and articleId = ?", userId, articleId)
	if err != nil {
		logger.Log("ERROR", "FINDUSERARTICLE", err.Error())
		return false, errors.ArticleNotFoundError
	}
	defer rows.Close()
	if rows.Next() {
		var dislike bool
		err = rows.Scan(&dislike)
		if err != nil {
			return false, errors.ArticleNotFoundError
		}
		return dislike, nil
	}
	return false, errors.ArticleNotFoundError
}

func (db UserArticleRepositoryImpl) SetUserArticleDislikeTo(userId int64, articleId int64, dislike bool, logger logger.Logger) error {
	_, err := db.DB.Exec("update articleLikes set dislike = ?, timestamp = ? where userId=? and articleId = ?", dislike, time.Now(), userId, articleId)
	return err
}

func (db UserArticleRepositoryImpl) RemoveUserArticleDislike(userId int64, articleId int64, logger logger.Logger) error {
	_, err := db.DB.Exec("delete from articleLikes where userId=? and articleId = ?", userId, articleId)
	return err
}

func (db UserArticleRepositoryImpl) GetUserArticle(userId int64, articleId int64, logger logger.Logger) (model.ArticleUser, error) {
	rows, err := db.DB.Query(`select a.*, 
		 (select max(ual.timestamp) from articleLikes ual where ual.articleId = a.articleId ) as lastUserLikeActivity, 
		 (select max(uac.timestamp) from articleComments uac where uac.articleId = a.articleId ) as lastUserCommentActivity, 
		 	(select max(ucl.timestamp) from articleComments uac join articleCommentLikes ucl on ucl.commentId = uac.commentId 
		 	where uac.articleId = a.articleId ) as lastUserLikeCommentActivity, 
		 (select count(*) from articleLikes aa where aa.articleId = a.articleId and aa.dislike) as dislikes, 
		 	(select count(*) from articleLikes aa where aa.articleId = a.articleId and not aa.dislike) as likes, 
		 	(select count(*) from articleComments aa where aa.articleId = a.articleId) as comments, 
		 		(select count(*) from articleLikes aa where aa.articleId = a.articleId and aa.dislike and aa.userid = ?) as userdislike, 
		 		(select count(*) from articleLikes aa where aa.articleId = a.articleId and not aa.dislike and aa.userid = ?) as userlike, 
		 		(select count(*) from articleComments aa where aa.articleId = a.articleId and aa.userid = ?) as usercomment 
		 		from articles a 
		 		where a.articleId = ?`,
		userId, userId, userId, articleId)
	if err != nil {
		logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		return model.ArticleUser{}, nil
	}
	if rows.Next() {
		var lastUserLikeActivity mysql.NullTime
		var lastUserCommentActivity mysql.NullTime
		var lastUserLikeCommentActivity mysql.NullTime
		article := model.ArticleUser{}
		err := rows.Scan(&article.ArticleId, &article.SourceName, &article.Title, &article.Link, &article.Description,
			&article.PubDate, &article.Category, &article.PictureUrl,
			&lastUserLikeActivity, &lastUserCommentActivity, &lastUserLikeCommentActivity,
			&article.Dislikes, &article.Likes, &article.Comments,
			&article.Dislike, &article.Like, &article.Comment)
		article.LastActivityDate = GetMaxTime(lastUserLikeActivity, lastUserCommentActivity, lastUserLikeCommentActivity, article.PubDate)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return article, nil
	}
	return model.ArticleUser{}, errors.ArticleNotFoundError
}
