package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql/data"
)

type UserUpdateRepositoryImpl struct {
}

func (UserUpdateRepositoryImpl) UpdateUser(userId int64, name string, phone string, logger logger.Logger) error {
	return data.UpdateUser(userId, name, phone, logger)
}
