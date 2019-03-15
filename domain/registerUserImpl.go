package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
	"github.com/sashamerkulev/rssservice/model"
)

type RegisterUser struct {
	DeviceId   string
	FirebaseId string
	Logger     logger.Logger
}

func (registerUser RegisterUser) RegisterUser() (user model.User, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tok, _ := token.SigningString()
	userId, err := db.RegisterUser(registerUser.DeviceId, registerUser.FirebaseId, tok, registerUser.Logger)
	return model.User{UserToken: tok, UserId: userId, Name: "", Phone: ""}, err
}
