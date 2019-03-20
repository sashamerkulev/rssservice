package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
)

type UserInfoRepositoryImpl struct {
}

func (UserInfoRepositoryImpl) GetUserInfo(userId int64, logger logger.Logger) error {
	return nil
}
