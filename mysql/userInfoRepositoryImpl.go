package mysql

import (
	"database/sql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
)

type UserInfoRepositoryImpl struct {
	DB *sql.DB
}

func (db UserInfoRepositoryImpl) GetUserInfo(userId int64, logger logger.Logger) (user model.User, err error) {
	rows, err := db.DB.Query("SELECT UserId, CASE WHEN LENGTH(UserName) = 0 OR UserName IS NULL THEN CONCAT('гость_', CONVERT(UserId, char)) ELSE UserName END AS UserName, UserPhone FROM userinfo where userId = ?", userId)
	if err != nil {
		logger.Log("ERROR", "GETUSERINFO", err.Error())
		return model.User{Name: ""}, nil
	}
	if rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.UserId, &user.Name, &user.Phone)
		if err != nil {
			logger.Log("ERROR", "GETARTICLEUSER", err.Error())
		}
		return user, nil
	}
	return model.User{Name: ""}, errors.UserNotFoundError
}
