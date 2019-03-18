package domain

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
)

type UploadPhotoUser struct {
	UserId int64
	Photo  []byte
	Logger logger.Logger
}

func (upu UploadPhotoUser) UploadPhoto() error {
	return db.UploadUserPhoto(upu.UserId, upu.Photo, upu.Logger)
}
