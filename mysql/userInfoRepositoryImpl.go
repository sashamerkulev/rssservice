package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/mysql/data"
)

type UserInfoRepositoryImpl struct {
}

func (UserInfoRepositoryImpl) GetUserInfo(userId int64, logger logger.Logger) (user model.User, err error) {
	return data.GetUserInfo(userId, logger)
}
