package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/reader"
)

type MainRepository interface {
	Open() error
	Close() error
	GetLogger(userId int64, userIP string) logger.Logger
	GetUserIdByToken(token string) (int64, error)
	GetSources() ([]reader.Link, error)
}
