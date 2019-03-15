package domain

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
	"github.com/sashamerkulev/rssservice/model"
)

type UpdateUser struct {
	UserToken string
	Name      string
	Phone     string
	Logger    logger.Logger
}

func (updateUser UpdateUser) UpdateUser() (user model.User, err error) {
	userId, err := db.GetUserIdByToken(updateUser.UserToken, updateUser.Logger)
	if err != nil {
		return model.User{}, err
	}
	err = db.UpdateUser(userId, updateUser.Name, updateUser.Phone, updateUser.Logger)
	return model.User{UserToken: updateUser.UserToken, UserId: userId, Name: updateUser.Name, Phone: updateUser.Phone}, err
}
