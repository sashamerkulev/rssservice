package mysql

import (
	"database/sql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/models"
)

type UserInfoRepositoryImpl struct {
	DB *sql.DB
}

func (db UserInfoRepositoryImpl) GetUserInfo(userId int64, logger logger.Logger) (user models.User, err error) {
	rows, err := db.DB.Query("SELECT UserId, CASE WHEN LENGTH(UserName) = 0 OR UserName IS NULL THEN CONCAT('гость_', CONVERT(UserId, char)) ELSE UserName END AS UserName, "+
		"COALESCE(UserPhone,'') AS UserPhone FROM users where userId = ?", userId)
	if err != nil {
		logger.Log("ERROR", "GETUSERINFO", err.Error())
		return models.User{Name: ""}, nil
	}
	if rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.UserId, &user.Name, &user.Phone)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return user, nil
	}
	return models.User{Name: ""}, errors.UserNotFoundError
}
