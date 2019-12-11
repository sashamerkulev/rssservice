package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/models"
	"time"
)

type MainRepository interface {
	Open(connectionString string) error
	Close() error
	GetLogger(userId int64, userIP string) logger.Logger
	GetUserIdByToken(token string) (int64, error)
	GetSources() ([]models.Link, error)
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
