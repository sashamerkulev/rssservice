package data

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"sort"
	"strings"
	"time"
)

func AddArticles(articles []model.Article, logger logger.Logger) {
	tx, err := DB.Begin()
	if err != nil {
		logger.Log("ERROR", "ADDALL", err.Error())
		return
	}
	insertStmt, err := DB.Prepare("INSERT INTO article(SourceName, Title, Link, Description, PubDate, Category, PictureUrl) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		logger.Log("ERROR", "ADDALL", err.Error())
		tx.Rollback()
		return
	}
	defer tx.Commit()
	defer insertStmt.Close()
	for i := 0; i < len(articles); i++ {
		_, err = insertStmt.Exec(articles[i].SourceName, articles[i].Title, articles[i].Link, articles[i].Description, articles[i].PubDate, articles[i].Category, articles[i].PictureUrl)
		if err != nil {
			isDuplicate := strings.Contains(err.Error(), "Error 1062")
			if !isDuplicate {
				logger.Log("ERROR", "ADDARTICLES", err.Error())
			}
			continue
		}
	}
}

func WipeOldArticles(wipeTime time.Time, logger logger.Logger) {
	result, err := DB.Exec("DELETE FROM Article WHERE "+
		"ArticleId not in (SELECT * FROM (SELECT a1.ArticleId FROM Article a1 JOIN UserArticleLikes ual on ual.ArticleId = a1.ArticleId "+
		" UNION "+
		" SELECT a1.ArticleId FROM Article a1 JOIN UserArticleComments uac on uac.ArticleId = a1.ArticleId) as art) "+
		"AND PubDate <= ?", wipeTime)
	if err != nil {
		logger.Log("ERROR", "WIPE", err.Error())
		return
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		logger.Log("ERROR", "WIPE", err.Error())
		return
	}
	logger.Log("DEBUG", "WIPE", "Rows ("+fmt.Sprint(deleted)+") was deleted at "+wipeTime.Format(time.RFC3339))
}

func GetUserArticles(userId int64, lastTime time.Time, logger logger.Logger) (results []model.ArticleUser, err error) {
	results = make([]model.ArticleUser, 0)
	currentTime := time.Now()
	rows, err := DB.Query(`select * from (select a.*, 
			 (select max(ual.timestamp) from userarticlelikes ual where ual.articleId = a.articleId ) as lastUserLikeActivity, 
			 (select max(uac.timestamp) from userarticlecomments uac where uac.articleId = a.articleId ) as lastUserCommentActivity, 
			 	(select max(ucl.timestamp) from userarticlecomments uac join usercommentlikes ucl on ucl.commentId = uac.commentId 
			 	where uac.articleId = a.articleId ) as lastUserLikeCommentActivity, 
			 (select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike) as dislikes, 
			 	(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike) as likes, 
			 	(select count(*) from userarticlecomments aa where aa.articleId = a.articleId) as comments, 
			 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike and aa.userid = ?) as userdislike, 
			 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike and aa.userid = ?) as userlike, 
			 		(select count(*) from userarticlecomments aa where aa.articleId = a.articleId and aa.userid = ?) as usercomment 
			 		from article a) b
			 		where (b.PubDate >= ? and b.PubDate < ? or (b.lastUserLikeActivity >= ? or b.lastUserCommentActivity >= ? or b.lastUserLikeCommentActivity >= ?)) `,
		userId, userId, userId, lastTime, currentTime, lastTime, lastTime, lastTime)
	if err != nil {
		logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		return results, err
	}
	for rows.Next() {
		var lastUserLikeActivity mysql.NullTime
		var lastUserCommentActivity mysql.NullTime
		var lastUserLikeCommentActivity mysql.NullTime

		article := model.ArticleUser{}
		err := rows.Scan(&article.ArticleId, &article.SourceName, &article.Title, &article.Link, &article.Description,
			&article.PubDate, &article.Category, &article.PictureUrl,
			&lastUserLikeActivity, &lastUserCommentActivity, &lastUserLikeCommentActivity,
			&article.Dislikes, &article.Likes, &article.Comments,
			&article.Dislike, &article.Like, &article.Comment)
		article.LastActivityDate = getMaxTime(lastUserLikeActivity, lastUserCommentActivity, lastUserLikeCommentActivity, article.PubDate)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		results = append(results, article)
	}
	//}

	sort.Slice(results, func(i, j int) bool {
		return results[i].PubDate.After(results[j].PubDate)
	})

	return results, nil
}

func FindUserArticleDislike(userId int64, articleId int64, logger logger.Logger) (bool, error) {
	rows, err := DB.Query("select dislike from userArticleLikes WHERE userId = ? and articleId = ?", userId, articleId)
	if err != nil {
		logger.Log("ERROR", "FINDUSERARTICLE", err.Error())
		return false, errors.ArticleNotFoundError()
	}
	defer rows.Close()
	if rows.Next() {
		var dislike bool
		err = rows.Scan(&dislike)
		if err != nil {
			return false, errors.ArticleNotFoundError()
		}
		return dislike, nil
	}
	return false, errors.ArticleNotFoundError()
}

func RemoveUserArticleDislike(userId int64, articleId int64, logger logger.Logger) error {
	_, err := DB.Exec("delete from userArticleLikes where userId=? and articleId = ?", userId, articleId)
	return err
}

func SetUserArticleDislikeTo(userId int64, articleId int64, dislike bool, logger logger.Logger) error {
	_, err := DB.Exec("update userArticleLikes set dislike = ?, timestamp = ? where userId=? and articleId = ?", dislike, time.Now(), userId, articleId)
	return err
}

func LikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	_, err := DB.Exec("insert into userArticleLikes (userId, articleId, dislike, timestamp) values(?,?,?,?)", userId, articleId, false, time.Now())
	return err
}

func DislikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	_, err := DB.Exec("insert into userArticleLikes (userId, articleId, dislike, timestamp) values(?,?,?,?)", userId, articleId, true, time.Now())
	return err
}

func GetUserArticle(userId int64, articleId int64, logger logger.Logger) (model.ArticleUser, error) {
	rows, err := DB.Query("select a.*, "+
		" (select max(ual.timestamp) from userarticlelikes ual where ual.articleId = a.articleId ) as lastUserLikeActivity, "+
		" (select max(uac.timestamp) from userarticlecomments uac where uac.articleId = a.articleId ) as lastUserCommentActivity, "+
		" 	(select max(ucl.timestamp) from userarticlecomments uac join usercommentlikes ucl on ucl.commentId = uac.commentId "+
		" 	where uac.articleId = a.articleId ) as lastUserLikeCommentActivity, "+
		" (select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike) as dislikes, "+
		" 	(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike) as likes, "+
		" 	(select count(*) from userarticlecomments aa where aa.articleId = a.articleId) as comments, "+
		" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike and aa.userid = ?) as userdislike, "+
		" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike and aa.userid = ?) as userlike, "+
		" 		(select count(*) from userarticlecomments aa where aa.articleId = a.articleId and aa.userid = ?) as usercomment "+
		" 		from article a "+
		" 		where a.articleId = ?",
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
		article.LastActivityDate = getMaxTime(lastUserLikeActivity, lastUserCommentActivity, lastUserLikeCommentActivity, article.PubDate)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return article, nil
	}
	return model.ArticleUser{}, errors.ArticleNotFoundError()
}

func getMaxTime(t1 mysql.NullTime, t2 mysql.NullTime, t3 mysql.NullTime, defaultTime time.Time) time.Time {
	max12 := getMaxTime2(t1, t2)
	max3 := getMaxTime2(max12, t3)
	if max3.Valid {
		return max3.Time
	}
	return defaultTime
}

func getMaxTime2(t1 mysql.NullTime, t2 mysql.NullTime) mysql.NullTime {
	if t1.Valid && t2.Valid {
		if t1.Time.After(t2.Time) {
			return t1
		}
		return t2
	}
	if t1.Valid {
		return t1
	}
	if t2.Valid {
		return t2
	}
	return mysql.NullTime{Valid: false}
}
