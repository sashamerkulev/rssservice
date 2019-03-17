package db

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/reader"
	"time"
)

func GetArticleUser(UserId int64, lastTime time.Time, logger logger.Logger) (results []model.ArticleUser, err error) {
	results = make([]model.ArticleUser, 0)
	// TODO improve SQL statements and remove this 'for'
	for i := 0; i < len(reader.Urls); i++ {
		rows, err := DB.Query("select a.*, "+
			" (select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike) as dislikes, "+
			" 	(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike) as likes, "+
			" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and aa.dislike and aa.userid = ?) as userdislike, "+
			" 		(select count(*) from userarticlelikes aa where aa.articleId = a.articleId and not aa.dislike and aa.userid = ?) as userlike "+
			" 		from article a "+
			" 		where a.sourcename = ? and a.PubDate >= ?"+
			" order by a.PubDate desc "+
			" limit 20",
			UserId, UserId, reader.Urls[i].Name, lastTime)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
			continue
		}
		for rows.Next() {
			article := model.ArticleUser{}
			err := rows.Scan(&article.ArticleId, &article.SourceName, &article.Title, &article.Link, &article.Description,
				&article.PubDate, &article.Category, &article.PictureUrl, &article.Dislikes, &article.Likes, &article.Dislike, &article.Like)
			if err != nil {
				logger.Log("ERROR", "GETARTICLEUSER", err.Error())
			}
			results = append(results, article)
		}
	}
	return results, nil
}

func GetUserIdByToken(userToken string, logger logger.Logger) (userId int64, err error) {
	rows, err := DB.Query("select userId from userDeviceToken WHERE token = ?", userToken)
	if err != nil {
		logger.Log("ERROR", "GETUSERIDBYTOKEN", err.Error())
		return -1, errors.UserNotFoundError()
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return -1, errors.UserNotFoundError()
		}
		return userId, nil
	}
	return -1, errors.UserNotFoundError()
}

func FindUserIdByDeviceId(deviceId string, logger logger.Logger) (userId int64, err error) {
	rows, err := DB.Query("select userId from userDeviceToken WHERE deviceId = ?", deviceId)
	if err != nil {
		logger.Log("ERROR", "FINDUSERIDBYDEVICEID", err.Error())
		return -1, errors.UserNotFoundError()
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return -1, errors.UserNotFoundError()
		}
		return userId, nil
	}
	return -1, errors.UserNotFoundError()
}

func AddTokenForUserIdAndDeviceId(userId int64, deviceId string, userToken string, logger logger.Logger) error {
	_, err := DB.Exec("insert into userDeviceToken(userId, deviceId, timestamp, token) values(?,?,?,?)", userId, deviceId, time.Now(), userToken)
	if err != nil {
		logger.Log("ERROR", "ADDTOKENFORUSERIDANDDEVICEID", err.Error())
		return errors.UserRegistrationError()
	}
	return nil
}

func RegisterUser(deviceId string, firebaseId string, userToken string, logger logger.Logger) (userId int64, err error) {
	tx, err := DB.Begin()
	if err != nil {
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.TransactionOpenError()
	}
	defer tx.Commit()
	result, err := DB.Exec("insert into userInfo(UserName, UserPhone) values(?,?)", "", "")
	if err != nil {
		tx.Rollback()
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.UserRegistrationError()
	}
	newUserId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.UserRegistrationError()
	}
	result, err = DB.Exec("insert into userDevices(userId, deviceId, firebaseId) values(?,?,?)", newUserId, deviceId, firebaseId)
	if err != nil {
		tx.Rollback()
		logger.Log("ERROR", "REGISTERUSER", err.Error())
		return -1, errors.UserRegistrationError()
	}
	err = AddTokenForUserIdAndDeviceId(newUserId, deviceId, userToken, logger)
	return newUserId, err
}

func UpdateUser(userId int64, name string, phone string, logger logger.Logger) (err error) {
	_, err = DB.Exec("update userInfo set UserName=?, UserPhone=? where userId = ?", name, phone, userId)
	if err != nil {
		logger.Log("ERROR", "UPDATEUSER", err.Error())
		return errors.UserUpdateError()
	}
	return nil
}

func LikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	found, err := findUserArticle(userId, articleId, logger)
	if err != nil {
		return err
	}
	if found == -1 {
		return addUserArticleLike(userId, articleId, false)
	} else {
		if found == 1 {
			return deleteUserArticleLike(userId, articleId) // remove like
		} else {
			return updateUserArticleLike(userId, articleId, false) // dislike is change to like
		}
	}
}

func DislikeArticle(userId int64, articleId int64, logger logger.Logger) error {
	found, err := findUserArticle(userId, articleId, logger)
	if err != nil {
		return err
	}
	if found == -1 {
		return addUserArticleLike(userId, articleId, true)
	} else {
		if found == 1 {
			return updateUserArticleLike(userId, articleId, true) // like is change to dislike
		} else {
			return deleteUserArticleLike(userId, articleId) // remove dislike
		}
	}
}

func addUserArticleLike(userId int64, articleId int64, dislike bool) error {
	_, err := DB.Exec("insert into userArticleLikes(userId, articleId, dislike) values(?,?,?)", userId, articleId, dislike)
	return err
}

func deleteUserArticleLike(userId int64, articleId int64) error {
	_, err := DB.Exec("delete from userArticleLikes where userId=? and articleId = ?", userId, articleId)
	return err
}

func updateUserArticleLike(userId int64, articleId int64, dislike bool) error {
	_, err := DB.Exec("update userArticleLikes set dislike = ? where userId=? and articleId = ?", dislike, userId, articleId)
	return err
}

func findUserArticle(userId int64, articleId int64, logger logger.Logger) (int, error) {
	rows, err := DB.Query("select dislike from userArticleLikes WHERE userId = ? and articleId = ?", userId, articleId)
	if err != nil {
		logger.Log("ERROR", "FINDUSERARTICLE", err.Error())
		return -1, errors.ArticleNotFoundError()
	}
	defer rows.Close()
	if rows.Next() {
		var dislike bool
		err = rows.Scan(&dislike)
		if err != nil {
			return -1, errors.ArticleNotFoundError()
		}
		if dislike { // already dislike (true)
			return 0, nil
		} else { // already like (false)
			return 1, nil
		}
	}
	return -1, nil
}
