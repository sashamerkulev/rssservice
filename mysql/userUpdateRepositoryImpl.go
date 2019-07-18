package mysql

import (
	"database/sql"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
)

type UserUpdateRepositoryImpl struct {
	DB *sql.DB
}

func (db UserUpdateRepositoryImpl) UpdateUser(userId int64, name string, phone string, logger logger.Logger) error {
	_, err := db.DB.Exec("update userInfo set UserName=?, UserPhone=? where userId = ?", name, phone, userId)
	if err != nil {
		logger.Log("ERROR", "UPDATEUSER", err.Error())
		return errors.UserUpdateError
	}
	return nil

}
