package data

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/reader"
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
		"ArticleId not in (SELECT * FROM (SELECT a1.ArticleId FROM Article a1 LEFT JOIN UserArticleLikes ual on ual.ArticleId = a1.ArticleId LEFT JOIN UserArticleComments uac on uac.ArticleId = a1.ArticleId) as art) "+
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
	// TODO improve SQL statements and remove this 'for'
	for i := 0; i < len(reader.Urls); i++ {
		rows, err := DB.Query("select a.*, "+
			" (select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike) as dislikes, "+
			" 	(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike) as likes, "+
			" 	(select count(*) from userarticlecomments aa where aa.articleId = a.articleId) as comments, "+
			" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike and aa.userid = ?) as userdislike, "+
			" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike and aa.userid = ?) as userlike, "+
			" 		(select count(*) from userarticlecomments aa where aa.articleId = a.articleId and aa.userid = ?) as usercomment "+
			" 		from article a "+
			" 		where a.sourcename = ? and a.PubDate >= ?"+
			" order by a.PubDate desc "+
			" limit 20",
			userId, userId, userId, reader.Urls[i].Name, lastTime)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
			continue
		}
		for rows.Next() {
			article := model.ArticleUser{}
			err := rows.Scan(&article.ArticleId, &article.SourceName, &article.Title, &article.Link, &article.Description,
				&article.PubDate, &article.Category, &article.PictureUrl,
				&article.Dislikes, &article.Likes, &article.Comments,
				&article.Dislike, &article.Like, &article.Comment)
			if err != nil {
				logger.Log("ERROR", "GETARTICLEUSER", err.Error())
			}
			results = append(results, article)
		}
	}

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
	_, err := DB.Exec("update userArticleLikes set dislike = ? where userId=? and articleId = ?", dislike, userId, articleId)
	return err
}

func LikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	_, err := DB.Exec("insert into userArticleLikes (userId, articleId, dislike) values(?,?,?)", userId, articleId, false)
	return err
}

func DislikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	_, err := DB.Exec("insert into userArticleLikes (userId, articleId, dislike) values(?,?,?)", userId, articleId, true)
	return err
}

func GetUserArticle(userId int64, articleId int64, logger logger.Logger) (model.ArticleUser, error) {
	rows, err := DB.Query("select a.*, "+
		" (select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike) as dislikes, "+
		" 	(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike) as likes, "+
		" 	(select count(*) from userarticlecomments aa where aa.articleId = a.articleId) as comments, "+
		" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike and aa.userid = ?) as userdislike, "+
		" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike and aa.userid = ?) as userlike, "+
		" 		(select count(*) from userarticlecomments aa where aa.articleId = a.articleId and aa.userid = ?) as usercomment "+
		" 		from article a "+
		" 		where a.articleId = ?"+
		" order by a.PubDate desc ",
		userId, userId, userId, articleId)
	if err != nil {
		logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		return model.ArticleUser{}, nil
	}
	if rows.Next() {
		article := model.ArticleUser{}
		err := rows.Scan(&article.ArticleId, &article.SourceName, &article.Title, &article.Link, &article.Description,
			&article.PubDate, &article.Category, &article.PictureUrl,
			&article.Dislikes, &article.Likes, &article.Comments,
			&article.Dislike, &article.Like, &article.Comment)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return article, nil
	}
	return model.ArticleUser{}, nil
}
