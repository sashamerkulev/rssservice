package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/reader"
	"time"
)

type MainRepository interface {
	Open() error
	Open(connectionString string) error
	Close() error
	GetLogger(userId int64, userIP string) logger.Logger
	GetUserIdByToken(token string) (int64, error)
	GetSources() ([]reader.Link, error)
}

func StringToDate(date string) time.Time {
	if date == "" {
		date = "2006-01-02T15:04:05"
	}
	datetime, err := time.Parse("2006-01-02T15:04:05", date)
	if err == nil {
		datetime = datetime.Add(time.Duration(-3) * time.Hour) // TODO
	}
	return datetime
}
