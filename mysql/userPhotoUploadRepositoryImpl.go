package mysql

import (
	"fmt"
	"github.com/sashamerkulev/rssservice/logger"
	"io/ioutil"
)

type UserPhotoUploadRepositoryImpl struct {
}

func (UserPhotoUploadRepositoryImpl) UploadUserPhoto(userId int64, photo []byte, logger logger.Logger) error {
	return ioutil.WriteFile("pictures/"+fmt.Sprint(userId)+".png", photo, 0644)
	//return data.UploadUserPhoto(userId, photo, logger)
}

func (UserPhotoUploadRepositoryImpl) GetUserPhoto(userId int64, logger logger.Logger) (photo []byte, err error) {
	return ioutil.ReadFile("pictures/" + fmt.Sprint(userId) + ".png")
}
