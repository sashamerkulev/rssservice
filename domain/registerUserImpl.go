package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/model"
)

type RegisterUser struct {
	DeviceId   string
	FirebaseId string
	Logger     logger.Logger
}

func (registerUser RegisterUser) RegisterUser() (user model.User, err error) {
	userId, err := db.FindUserIdByDeviceId(registerUser.DeviceId, registerUser.Logger)
	if err == nil {
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		token, _ := jwtToken.SigningString()
		err = db.AddTokenForUserIdAndDeviceId(userId, registerUser.DeviceId, token, registerUser.Logger)
		if err != nil {
			return model.User{}, errors.UserRegistrationError()
		}
		return model.User{UserToken: token, UserId: userId, Name: "", Phone: ""}, err
	} else {
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		token, _ := jwtToken.SigningString()
		userId, err := db.RegisterUser(registerUser.DeviceId, registerUser.FirebaseId, token, registerUser.Logger)
		return model.User{UserToken: token, UserId: userId, Name: "", Phone: ""}, err
	}
}
