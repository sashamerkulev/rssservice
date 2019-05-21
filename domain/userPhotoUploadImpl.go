package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
)

type UserPhotoUploadRepository interface {
	UploadUserPhoto(userId int64, photo []byte, logger logger.Logger) error
	GetUserPhoto(userId int64, logger logger.Logger) (photo []byte, err error)
}

type UserPhoto struct {
	UserId     int64
	Logger     logger.Logger
	Repository UserPhotoUploadRepository
}

func (upu UserPhoto) UploadUserPhoto(photo []byte) error {
	return upu.Repository.UploadUserPhoto(upu.UserId, photo, upu.Logger)
}

func (upu UserPhoto) GetUserPhoto() (photo []byte, err error) {
	return upu.Repository.GetUserPhoto(upu.UserId, upu.Logger)
}
