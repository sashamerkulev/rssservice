package domain

import "github.com/sashamerkulev/rssservice/logger"

type MainRepository interface {
	Open() error
	Close() error
	GetLogger(userId int64, userIP string) logger.Logger
	GetUserIdByToken(token string) (int64, error)
}
