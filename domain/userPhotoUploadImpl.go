package domain

import (
	"github.com/sashamerkulev/rssservice/logger"
)

type UserPhotoUploadRepository interface {
	UploadUserPhoto(userId int64, photo []byte, logger logger.Logger) error
}

type UserPhotoUpload struct {
	UserId     int64
	Photo      []byte
	Logger     logger.Logger
	Repository UserPhotoUploadRepository
}

func (upu UserPhotoUpload) UploadUserPhoto() error {
	return upu.Repository.UploadUserPhoto(upu.UserId, upu.Photo, upu.Logger)
}
