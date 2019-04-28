package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"time"
)

type SetupRegisterRepository interface {
	FindUserIdByDeviceId(deviceId string, logger logger.Logger) (userId int64, err error)
	AddTokenForUserIdAndDeviceId(userId int64, deviceId string, token string, logger logger.Logger) error
	RegisterUser(deviceId string, firebaseId string, token string, logger logger.Logger) (userId int64, err error)
}

type SetupRegister struct {
	SetupId   string
	FirebaseId string
	Logger     logger.Logger
	Repository SetupRegisterRepository
}

var signingKey = []byte("secret")

func (ru SetupRegister) RegisterUser() (user model.User, err error) {
	userId, err := ru.Repository.FindUserIdByDeviceId(ru.SetupId, ru.Logger)
	if err == nil {
		token, _ := getJwtToken(ru.SetupId)
		err = ru.Repository.AddTokenForUserIdAndDeviceId(userId, ru.SetupId, token, ru.Logger)
		if err != nil {
			return model.User{}, errors.UserRegistrationError()
		}
		return model.User{UserId: userId, Name: "", Phone: "", Token: token}, err
	} else {
		token, _ := getJwtToken(ru.SetupId)
		userId, err := ru.Repository.RegisterUser(ru.SetupId, ru.FirebaseId, token, ru.Logger)
		return model.User{UserId: userId, Name: "", Phone: "", Token: token}, err
	}
}

func getJwtToken(setupId string) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["setupId"] = setupId
	claims["expired"] = time.Now().Add(time.Hour * 24).Unix()
	return jwtToken.SignedString(signingKey)
}
