package domain

import (
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
	"github.com/sashamerkulev/rssservice/model"
)

type UpdateUser struct {
	UserId int64
	Name   string
	Phone  string
	Logger logger.Logger
}

func (uu UpdateUser) UpdateUser() (user model.User, err error) {
	err = db.UpdateUser(uu.UserId, uu.Name, uu.Phone, uu.Logger)
	return model.User{UserId: uu.UserId, Name: uu.Name, Phone: uu.Phone}, err
}
