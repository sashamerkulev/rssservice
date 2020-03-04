package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/models"
)

type MainRepositoryImpl struct {
	Sources []models.Link
}

type UserDbLogger struct {
	UserId int64
	UserIP string
	DB     *sql.DB
}

var DB *sql.DB

func (MainRepositoryImpl) Open(connectionString string) error {
	mysql, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	DB = mysql
	return nil
}

func (MainRepositoryImpl) Close() error {
	err := DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (MainRepositoryImpl) GetLogger(userId int64, userIP string) logger.Logger {
	return UserDbLogger{DB: DB, UserId: userId, UserIP: userIP}
}

func (MainRepositoryImpl) GetUserIdByToken(token string) (int64, error) {
	return GetUserIdByToken(token)
}

func (r MainRepositoryImpl) GetSources() ([]models.Link, error) {
	return r.Sources, nil
}
