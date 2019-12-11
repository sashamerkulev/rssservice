package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/models"
)

type UserInfoRepository interface {
	GetUserInfo(userId int64, logger logger.Logger) (user models.User, err error)
}

type UserInfo struct {
	UserId     int64
	Logger     logger.Logger
	Repository UserInfoRepository
}

func (ui UserInfo) GetUserInfo() (user models.User, err error) {
	return ui.Repository.GetUserInfo(ui.UserId, ui.Logger)
}
