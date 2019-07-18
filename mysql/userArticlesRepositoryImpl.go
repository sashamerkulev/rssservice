package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"sort"
	"time"
)

type UserArticlesRepositoryImpl struct {
	DB *sql.DB
}

func (db UserArticlesRepositoryImpl) GetUserArticles(userId int64, lastTime time.Time, logger logger.Logger) ([]model.ArticleUser, error) {
	results := make([]model.ArticleUser, 0)
	currentTime := time.Now()
	rows, err := db.DB.Query(`select * from (select a.*, 
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
		article.LastActivityDate = GetMaxTime(lastUserLikeActivity, lastUserCommentActivity, lastUserLikeCommentActivity, article.PubDate)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		results = append(results, article)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].PubDate.After(results[j].PubDate)
	})
	return results, nil
}
