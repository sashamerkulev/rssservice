package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

type UserRegisterRepository interface {
	FindUserIdByDeviceId(deviceId string, logger logger.Logger) (userId int64, err error)
	AddTokenForUserIdAndDeviceId(userId int64, deviceId string, token string, logger logger.Logger) error
	RegisterUser(deviceId string, firebaseId string, token string, logger logger.Logger) (userId int64, err error)
}

type UserRegister struct {
	DeviceId   string
	FirebaseId string
	Logger     logger.Logger
	Repository UserRegisterRepository
}

var signingKey = []byte("secret")

func (ru UserRegister) RegisterUser() (user model.User, err error) {
	userId, err := ru.Repository.FindUserIdByDeviceId(ru.DeviceId, ru.Logger)
	if err == nil {
		token, _ := getJwtToken(ru.DeviceId)
		err = ru.Repository.AddTokenForUserIdAndDeviceId(userId, ru.DeviceId, token, ru.Logger)
		if err != nil {
			return model.User{}, errors.UserRegistrationError()
		}
		return model.User{UserId: userId, Name: "", Phone: ""}, err
	} else {
		token, _ := getJwtToken(ru.DeviceId)
		userId, err := ru.Repository.RegisterUser(ru.DeviceId, ru.FirebaseId, token, ru.Logger)
		return model.User{UserId: userId, Name: "", Phone: ""}, err
	}
}

func getJwtToken(deviceId string) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["deviceId"] = deviceId
	claims["expired"] = time.Now().Add(time.Hour * 24).Unix()
	return jwtToken.SignedString(signingKey)
}
