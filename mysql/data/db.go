package data

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/rssservice/logger"
	"time"
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

func (ul UserDbLogger) Log(severity string, tag string, message string) {
	go dbLog(ul.DB, severity, fmt.Sprint(ul.UserId), ul.UserIP, tag, message)
}

func dbLog(DB *sql.DB, severity string, userId string, userIp string, tag string, message string) {
	_, err := DB.Exec("INSERT INTO log(Severity, UserId, UserIP, Timestamp, Tag, Message) VALUES(?,?,?,?,?,?)", severity, userId, userIp, time.Now(), tag, message)
	if err != nil {
		logger.ConsoleLog(severity, userId, userIp, tag, message)
		logger.ConsoleLog(severity, userId, userIp, "DBLOGGER", err.Error())
	}
}
