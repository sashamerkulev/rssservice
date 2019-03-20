package mysql

import (
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql/data"
)

type UserPhotoUploadRepositoryImpl struct {
}

func (UserPhotoUploadRepositoryImpl) UploadUserPhoto(userId int64, photo []byte, logger logger.Logger) error {
	return data.UploadUserPhoto(userId, photo, logger)
}
