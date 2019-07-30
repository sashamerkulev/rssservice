package mysql

import (
	"fmt"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"strings"
	"time"
)

func AddArticles(articles []model.Article, logger logger.Logger) {
	tx, err := DB.Begin()
	if err != nil {
		logger.Log("ERROR", "ADDARTICLES", err.Error())
		return
	}
	insertStmt, err := DB.Prepare("INSERT INTO articles(SourceName, Title, Link, Description, PubDate, Category, PictureUrl) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		logger.Log("ERROR", "ADDARTICLES", err.Error())
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
	result, err := DB.Exec("DELETE FROM articles WHERE "+
		"ArticleId not in (SELECT * FROM (SELECT a1.ArticleId FROM articles a1 JOIN articleLikes ual on ual.ArticleId = a1.ArticleId "+
		" UNION "+
		" SELECT a1.ArticleId FROM articles a1 JOIN articleComments uac on uac.ArticleId = a1.ArticleId) as art) "+
		"AND PubDate <= ?", wipeTime)
	if err != nil {
		logger.Log("ERROR", "WIPEOLDARTICLES", err.Error())
		return
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		logger.Log("ERROR", "WIPEOLDARTICLES", err.Error())
		return
	}
	logger.Log("DEBUG", "WIPEOLDARTICLES", "Rows ("+fmt.Sprint(deleted)+") was deleted at "+wipeTime.Format(time.RFC3339))
}

func WipeOldActivities(wipeTime time.Time, logger logger.Logger) {
	result, err := DB.Exec(`DELETE FROM articles WHERE ArticleId in ( 
select b.articleId
from (
select a.articleId, max(coalesce(uac.timestamp, a.pubdate)) as timestamp
from articles a
left join articleComments uac on uac.articleId = a.articleId
group by a.ArticleId
union
select a.articleId, max(coalesce(uac.timestamp, a.pubdate)) as timestamp
from articles a
left join articleLikes uac on uac.articleId = a.articleId
group by a.ArticleId
) b WHERE b.timestamp < ?)
`, wipeTime)
	if err != nil {
		logger.Log("ERROR", "WIPEOLDACTIVITIES", err.Error())
		return
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		logger.Log("ERROR", "WIPEOLDACTIVITIES", err.Error())
		return
	}
	logger.Log("DEBUG", "WIPEOLDACTIVITIES", "Rows ("+fmt.Sprint(deleted)+") was deleted at "+wipeTime.Format(time.RFC3339))
}
