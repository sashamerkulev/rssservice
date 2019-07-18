package mysql

import (
	"database/sql"
	"fmt"
	"github.com/sashamerkulev/rssservice/logger"
	"io/ioutil"
)

type UserPhotoUploadRepositoryImpl struct {
	DB *sql.DB
}

func (db UserPhotoUploadRepositoryImpl) UploadUserPhoto(userId int64, photo []byte, logger logger.Logger) error {
	return ioutil.WriteFile("pictures/"+fmt.Sprint(userId)+".png", photo, 0644)
}

func (db UserPhotoUploadRepositoryImpl) GetUserPhoto(userId int64, logger logger.Logger) (photo []byte, err error) {
	return ioutil.ReadFile("pictures/" + fmt.Sprint(userId) + ".png")
}
