package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/reader"
)

type MainRepositoryImpl struct {
}

type UserDbLogger struct {
	UserId int64
	UserIP string
	DB     *sql.DB
}

var DB *sql.DB

func (MainRepositoryImpl) Open() error {
	mysql, err := sql.Open("mysql", "news:News,News@/dbnews?parseTime=true")
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

func (MainRepositoryImpl) GetSources() ([]reader.Link, error) {
	return reader.Urls, nil
}
