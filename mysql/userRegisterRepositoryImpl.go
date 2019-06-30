package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql/data"
)

type UserRegisterRepositoryImpl struct {
}

func (UserRegisterRepositoryImpl) FindUserIdByDeviceId(deviceId string, logger logger.Logger) (userId int64, err error) {
	return data.FindUserIdByDeviceId(deviceId, logger)
}

func (UserRegisterRepositoryImpl) AddTokenForUserIdAndDeviceId(userId int64, deviceId string, token string, logger logger.Logger) error {
	return data.AddTokenForUserIdAndDeviceId(userId, deviceId, token, logger)
}

func (UserRegisterRepositoryImpl) RegisterUser(deviceId string, firebaseId string, token string, logger logger.Logger) (userId int64, err error) {
	return data.RegisterUser(deviceId, firebaseId, token, logger)
}

func (UserRegisterRepositoryImpl) UpdateFirebaseId(userId int64, firebaseId string, logger logger.Logger) error {
	return data.UpdateFirebaseId(userId, firebaseId, logger)
}
