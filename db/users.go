package db

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

func GetArticleUser(UserId int64, lastTime time.Time, logger logger.Logger) []model.ArticleUser {
	results := make([]model.ArticleUser, 0)
	rows, err := DB.Query("select SourceName, Title, Link, Description, PubDate, Category, PictureUrl from article WHERE PubDate >= ?", lastTime)
	if err != nil {
		logger.Log("ERROR", "GET", err.Error())
		return results
	}
	defer rows.Close()
	for rows.Next() {
		article := model.Article{}
		err := rows.Scan(&article.SourceName, &article.Title, &article.Link, &article.Description, &article.PubDate, &article.Category, &article.PictureUrl)
		if err != nil {
			logger.Log("ERROR", "GET", err.Error())
		}
		//results = append(results, article)
	}
	return results
}

func GetUserIdByToken(userToken string, logger logger.Logger) (int64, error) {
	return 1, nil
}

func RegisterUser(deviceId string, firebaseId string, userToken string, logger logger.Logger) (userId int64, err error) {
	return 1, nil
}

func UpdateUser(userId int64, name string, phone string, logger logger.Logger) (err error) {
	return nil
}
