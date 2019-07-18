package mysql

import (
	"database/sql"
	"fmt"
	"github.com/sashamerkulev/rssservice/logger"
	"time"
)

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